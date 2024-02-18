/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package vm

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/vms"
	"github.com/bpg/terraform-provider-proxmox/proxmox/pools"
	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
	"github.com/bpg/terraform-provider-proxmox/proxmoxtf"
	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/structure"
	"github.com/bpg/terraform-provider-proxmox/utils"
)

func vmCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clone := d.Get(mkClone).([]interface{})

	if len(clone) > 0 {
		return vmCreateClone(ctx, d, m)
	}

	return vmCreateCustom(ctx, d, m)
}

// Check for an existing CloudInit IDE drive. If no such drive is found, return the specified `defaultValue`.
func findExistingCloudInitDrive(vmConfig *vms.GetResponseData, vmID int, defaultValue string) string {
	devices := []*vms.CustomStorageDevice{
		vmConfig.IDEDevice0, vmConfig.IDEDevice1, vmConfig.IDEDevice2, vmConfig.IDEDevice3,
	}
	for i, device := range devices {
		if device != nil && device.Enabled && device.IsCloudInitDrive(vmID) {
			return fmt.Sprintf("ide%d", i)
		}
	}

	return defaultValue
}

// Return a pointer to the IDE device configuration based on its name. The device name is assumed to be a
// valid IDE interface name.
func getIdeDevice(vmConfig *vms.GetResponseData, deviceName string) *vms.CustomStorageDevice {
	ideDevice := vmConfig.IDEDevice3

	switch deviceName {
	case "ide0":
		ideDevice = vmConfig.IDEDevice0
	case "ide1":
		ideDevice = vmConfig.IDEDevice1
	case "ide2":
		ideDevice = vmConfig.IDEDevice2
	}

	return ideDevice
}

// Delete IDE interfaces that can then be used for CloudInit. The first interface will always
// be deleted. The second will be deleted only if it isn't empty and isn't the same as the
// first.
func deleteIdeDrives(ctx context.Context, vmAPI *vms.Client, itf1 string, itf2 string) diag.Diagnostics {
	ddUpdateBody := &vms.UpdateRequestBody{}
	ddUpdateBody.Delete = append(ddUpdateBody.Delete, itf1)
	tflog.Debug(ctx, fmt.Sprintf("Deleting IDE interface '%s'", itf1))

	if itf2 != "" && itf2 != itf1 {
		ddUpdateBody.Delete = append(ddUpdateBody.Delete, itf2)
		tflog.Debug(ctx, fmt.Sprintf("Deleting IDE interface '%s'", itf2))
	}

	e := vmAPI.UpdateVM(ctx, ddUpdateBody)
	if e != nil {
		return diag.FromErr(e)
	}

	return nil
}

// Start the VM, then wait for it to actually start; it may not be started immediately if running in HA mode.
func vmStart(ctx context.Context, vmAPI *vms.Client, d *schema.ResourceData) diag.Diagnostics {
	var diags diag.Diagnostics

	tflog.Debug(ctx, "Starting VM")

	startVMTimeout := d.Get(mkTimeoutStartVM).(int)

	log, e := vmAPI.StartVM(ctx, startVMTimeout)
	if e != nil {
		return append(diags, diag.FromErr(e)...)
	}

	if len(log) > 0 {
		lines := "\n\t| " + strings.Join(log, "\n\t| ")
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  fmt.Sprintf("the VM startup task finished with a warning, task log:\n%s", lines),
		})
	}

	return append(diags, diag.FromErr(vmAPI.WaitForVMStatus(ctx, "running", startVMTimeout, 1))...)
}

// Shutdown the VM, then wait for it to actually shut down (it may not be shut down immediately if
// running in HA mode).
func vmShutdown(ctx context.Context, vmAPI *vms.Client, d *schema.ResourceData) diag.Diagnostics {
	tflog.Debug(ctx, "Shutting down VM")

	forceStop := types.CustomBool(true)
	shutdownTimeout := d.Get(mkTimeoutShutdownVM).(int)

	e := vmAPI.ShutdownVM(ctx, &vms.ShutdownRequestBody{
		ForceStop: &forceStop,
		Timeout:   &shutdownTimeout,
	}, shutdownTimeout+30)
	if e != nil {
		return diag.FromErr(e)
	}

	return diag.FromErr(vmAPI.WaitForVMStatus(ctx, "stopped", shutdownTimeout, 1))
}

// Forcefully stop the VM, then wait for it to actually stop.
func vmStop(ctx context.Context, vmAPI *vms.Client, d *schema.ResourceData) diag.Diagnostics {
	tflog.Debug(ctx, "Stopping VM")

	stopTimeout := d.Get(mkTimeoutStopVM).(int)

	e := vmAPI.StopVM(ctx, stopTimeout+30)
	if e != nil {
		return diag.FromErr(e)
	}

	return diag.FromErr(vmAPI.WaitForVMStatus(ctx, "stopped", stopTimeout, 1))
}

func vmCreateClone(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(proxmoxtf.ProviderConfiguration)

	api, e := config.GetClient()
	if e != nil {
		return diag.FromErr(e)
	}

	clone := d.Get(mkClone).([]interface{})
	cloneBlock := clone[0].(map[string]interface{})
	cloneRetries := cloneBlock[mkCloneRetries].(int)
	cloneDatastoreID := cloneBlock[mkCloneDatastoreID].(string)
	cloneNodeName := cloneBlock[mkCloneNodeName].(string)
	cloneVMID := cloneBlock[mkCloneVMID].(int)
	cloneFull := cloneBlock[mkCloneFull].(bool)

	description := d.Get(mkDescription).(string)
	name := d.Get(mkName).(string)
	tags := d.Get(mkTags).([]interface{})
	nodeName := d.Get(mkNodeName).(string)
	poolID := d.Get(mkPoolID).(string)
	vmIDUntyped, hasVMID := d.GetOk(mkVMID)
	vmID := vmIDUntyped.(int)

	if !hasVMID {
		vmIDNew, err := api.Cluster().GetVMID(ctx)
		if err != nil {
			return diag.FromErr(err)
		}

		vmID = *vmIDNew

		err = d.Set(mkVMID, vmID)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	fullCopy := types.CustomBool(cloneFull)

	cloneBody := &vms.CloneRequestBody{
		FullCopy: &fullCopy,
		VMIDNew:  vmID,
	}

	if cloneDatastoreID != "" {
		cloneBody.TargetStorage = &cloneDatastoreID
	}

	if description != "" {
		cloneBody.Description = &description
	}

	if name != "" {
		cloneBody.Name = &name
	}

	if poolID != "" {
		cloneBody.PoolID = &poolID
	}

	cloneTimeout := d.Get(mkTimeoutClone).(int)

	if cloneNodeName != "" && cloneNodeName != nodeName {
		// Check if any used datastores of the source VM are not shared
		vmConfig, err := api.Node(cloneNodeName).VM(cloneVMID).GetVM(ctx)
		if err != nil {
			return diag.FromErr(err)
		}

		datastores := getDiskDatastores(vmConfig, d)

		onlySharedDatastores := true

		for _, datastore := range datastores {
			datastoreStatus, err2 := api.Node(cloneNodeName).Storage(datastore).GetDatastoreStatus(ctx)
			if err2 != nil {
				return diag.FromErr(err2)
			}

			if datastoreStatus.Shared != nil && !*datastoreStatus.Shared {
				onlySharedDatastores = false
				break
			}
		}

		if onlySharedDatastores {
			// If the source and the target node are not the same, only clone directly to the target node if
			//  all used datastores in the source VM are shared. Directly cloning to non-shared storage
			//  on a different node is currently not supported by proxmox.
			cloneBody.TargetNodeName = &nodeName

			err = api.Node(cloneNodeName).VM(cloneVMID).CloneVM(
				ctx,
				cloneRetries,
				cloneBody,
				cloneTimeout,
			)
			if err != nil {
				return diag.FromErr(err)
			}
		} else {
			// If the source and the target node are not the same and any used datastore in the source VM is
			//  not shared, clone to the source node and then migrate to the target node. This is a workaround
			//  for missing functionality in the proxmox api as recommended per
			//  https://forum.proxmox.com/threads/500-cant-clone-to-non-shared-storage-local.49078/#post-229727

			// Temporarily clone to local node
			err = api.Node(cloneNodeName).VM(cloneVMID).CloneVM(ctx, cloneRetries, cloneBody, cloneTimeout)
			if err != nil {
				return diag.FromErr(err)
			}

			// Wait for the virtual machine to be created and its configuration lock to be released before migrating.

			err = api.Node(cloneNodeName).VM(vmID).WaitForVMConfigUnlock(ctx, 600, 5, true)
			if err != nil {
				return diag.FromErr(err)
			}

			// Migrate to target node
			withLocalDisks := types.CustomBool(true)
			migrateBody := &vms.MigrateRequestBody{
				TargetNode:     nodeName,
				WithLocalDisks: &withLocalDisks,
			}

			if cloneDatastoreID != "" {
				migrateBody.TargetStorage = &cloneDatastoreID
			}

			err = api.Node(cloneNodeName).VM(vmID).MigrateVM(ctx, migrateBody, cloneTimeout)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	} else {
		e = api.Node(nodeName).VM(cloneVMID).CloneVM(ctx, cloneRetries, cloneBody, cloneTimeout)
	}

	if e != nil {
		return diag.FromErr(e)
	}

	d.SetId(strconv.Itoa(vmID))

	vmAPI := api.Node(nodeName).VM(vmID)

	// Wait for the virtual machine to be created and its configuration lock to be released.
	e = vmAPI.WaitForVMConfigUnlock(ctx, 600, 5, true)
	if e != nil {
		return diag.FromErr(e)
	}

	//// UPDATE AFTER CLONE, can we just call update?

	// Now that the virtual machine has been cloned, we need to perform some modifications.
	acpi := types.CustomBool(d.Get(mkACPI).(bool))
	audioDevices := vmGetAudioDeviceList(d)

	bios := d.Get(mkBIOS).(string)
	kvmArguments := d.Get(mkKVMArguments).(string)
	scsiHardware := d.Get(mkSCSIHardware).(string)
	cdrom := d.Get(mkCDROM).([]interface{})
	cpu := d.Get(mkCPU).([]interface{})
	initialization := d.Get(mkInitialization).([]interface{})
	hostPCI := d.Get(mkHostPCI).([]interface{})
	hostUSB := d.Get(mkHostUSB).([]interface{})
	keyboardLayout := d.Get(mkKeyboardLayout).(string)
	memory := d.Get(mkMemory).([]interface{})
	networkDevice := d.Get(mkNetworkDevice).([]interface{})
	operatingSystem := d.Get(mkOperatingSystem).([]interface{})
	serialDevice := d.Get(mkSerialDevice).([]interface{})
	onBoot := types.CustomBool(d.Get(mkOnBoot).(bool))
	tabletDevice := types.CustomBool(d.Get(mkTabletDevice).(bool))
	template := types.CustomBool(d.Get(mkTemplate).(bool))
	vga := d.Get(mkVGA).([]interface{})

	updateBody := &vms.UpdateRequestBody{
		AudioDevices: audioDevices,
	}

	ideDevices := vms.CustomStorageDevices{}

	var del []string

	//nolint:gosimple
	if acpi != dvACPI {
		updateBody.ACPI = &acpi
	}

	createAgent(d, updateBody)

	if kvmArguments != dvKVMArguments {
		updateBody.KVMArguments = &kvmArguments
	}

	if bios != dvBIOS {
		updateBody.BIOS = &bios
	}

	if scsiHardware != dvSCSIHardware {
		updateBody.SCSIHardware = &scsiHardware
	}

	if len(cdrom) > 0 || len(initialization) > 0 {
		ideDevices = vms.CustomStorageDevices{
			"ide0": &vms.CustomStorageDevice{
				Enabled: false,
			},
			"ide1": &vms.CustomStorageDevice{
				Enabled: false,
			},
			"ide2": &vms.CustomStorageDevice{
				Enabled: false,
			},
			"ide3": &vms.CustomStorageDevice{
				Enabled: false,
			},
		}
	}

	if len(cdrom) > 0 {
		cdromBlock := cdrom[0].(map[string]interface{})

		cdromEnabled := cdromBlock[mkCDROMEnabled].(bool)
		cdromFileID := cdromBlock[mkCDROMFileID].(string)
		cdromInterface := cdromBlock[mkCDROMInterface].(string)

		if cdromFileID == "" {
			cdromFileID = "cdrom"
		}

		cdromMedia := "cdrom"

		ideDevices[cdromInterface] = &vms.CustomStorageDevice{
			Enabled:    cdromEnabled,
			FileVolume: cdromFileID,
			Media:      &cdromMedia,
		}
	}

	if len(cpu) > 0 {
		cpuBlock := cpu[0].(map[string]interface{})

		cpuArchitecture := cpuBlock[mkCPUArchitecture].(string)
		cpuCores := cpuBlock[mkCPUCores].(int)
		cpuFlags := cpuBlock[mkCPUFlags].([]interface{})
		cpuHotplugged := cpuBlock[mkCPUHotplugged].(int)
		cpuLimit := cpuBlock[mkCPULimit].(int)
		cpuNUMA := types.CustomBool(cpuBlock[mkCPUNUMA].(bool))
		cpuSockets := cpuBlock[mkCPUSockets].(int)
		cpuType := cpuBlock[mkCPUType].(string)
		cpuUnits := cpuBlock[mkCPUUnits].(int)

		cpuFlagsConverted := make([]string, len(cpuFlags))

		for fi, flag := range cpuFlags {
			cpuFlagsConverted[fi] = flag.(string)
		}

		// Only the root account is allowed to change the CPU architecture, which makes this check necessary.
		if api.API().IsRootTicket() ||
			cpuArchitecture != dvCPUArchitecture {
			updateBody.CPUArchitecture = &cpuArchitecture
		}

		updateBody.CPUCores = &cpuCores
		updateBody.CPUEmulation = &vms.CustomCPUEmulation{
			Flags: &cpuFlagsConverted,
			Type:  cpuType,
		}
		updateBody.NUMAEnabled = &cpuNUMA
		updateBody.CPUSockets = &cpuSockets
		updateBody.CPUUnits = &cpuUnits

		if cpuHotplugged > 0 {
			updateBody.VirtualCPUCount = &cpuHotplugged
		}

		if cpuLimit > 0 {
			updateBody.CPULimit = &cpuLimit
		}
	}

	vmConfig, err := vmAPI.GetVM(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(initialization) > 0 {
		tflog.Trace(ctx, "Preparing the CloudInit configuration")

		initializationBlock := initialization[0].(map[string]interface{})
		initializationDatastoreID := initializationBlock[mkInitializationDatastoreID].(string)
		initializationInterface := initializationBlock[mkInitializationInterface].(string)

		existingInterface := findExistingCloudInitDrive(vmConfig, vmID, "ide2")
		if initializationInterface == "" {
			initializationInterface = existingInterface
		}

		tflog.Trace(ctx, fmt.Sprintf("CloudInit IDE interface is '%s'", initializationInterface))

		const cdromCloudInitEnabled = true

		cdromCloudInitFileID := fmt.Sprintf("%s:cloudinit", initializationDatastoreID)
		cdromCloudInitMedia := "cdrom"
		ideDevices[initializationInterface] = &vms.CustomStorageDevice{
			Enabled:    cdromCloudInitEnabled,
			FileVolume: cdromCloudInitFileID,
			Media:      &cdromCloudInitMedia,
		}

		if err := deleteIdeDrives(ctx, vmAPI, initializationInterface, existingInterface); err != nil {
			return err
		}

		updateBody.CloudInitConfig = vmGetCloudInitConfig(d)
	}

	if len(hostPCI) > 0 {
		updateBody.PCIDevices = vmGetHostPCIDeviceObjects(d)
	}

	if len(hostUSB) > 0 {
		updateBody.USBDevices = vmGetHostUSBDeviceObjects(d)
	}

	if len(cdrom) > 0 || len(initialization) > 0 {
		updateBody.IDEDevices = ideDevices
	}

	if keyboardLayout != dvKeyboardLayout {
		updateBody.KeyboardLayout = &keyboardLayout
	}

	if len(memory) > 0 {
		memoryBlock := memory[0].(map[string]interface{})

		memoryDedicated := memoryBlock[mkMemoryDedicated].(int)
		memoryFloating := memoryBlock[mkMemoryFloating].(int)
		memoryShared := memoryBlock[mkMemoryShared].(int)

		updateBody.DedicatedMemory = &memoryDedicated
		updateBody.FloatingMemory = &memoryFloating

		if memoryShared > 0 {
			memorySharedName := fmt.Sprintf("vm-%d-ivshmem", vmID)

			updateBody.SharedMemory = &vms.CustomSharedMemory{
				Name: &memorySharedName,
				Size: memoryShared,
			}
		}
	}

	if len(networkDevice) > 0 {
		updateBody.NetworkDevices = vmGetNetworkDeviceObjects(d)

		for i := 0; i < len(updateBody.NetworkDevices); i++ {
			if !updateBody.NetworkDevices[i].Enabled {
				del = append(del, fmt.Sprintf("net%d", i))
			}
		}

		for i := len(updateBody.NetworkDevices); i < maxResourceVirtualEnvironmentVMNetworkDevices; i++ {
			del = append(del, fmt.Sprintf("net%d", i))
		}
	}

	if len(operatingSystem) > 0 {
		operatingSystemBlock := operatingSystem[0].(map[string]interface{})
		operatingSystemType := operatingSystemBlock[mkOperatingSystemType].(string)

		updateBody.OSType = &operatingSystemType
	}

	if len(serialDevice) > 0 {
		updateBody.SerialDevices = vmGetSerialDeviceList(d)

		for i := len(updateBody.SerialDevices); i < maxResourceVirtualEnvironmentVMSerialDevices; i++ {
			del = append(del, fmt.Sprintf("serial%d", i))
		}
	}

	updateBody.StartOnBoot = &onBoot

	updateBody.SMBIOS = vmGetSMBIOS(d)

	updateBody.StartupOrder = vmGetStartupOrder(d)

	//nolint:gosimple
	if tabletDevice != dvTabletDevice {
		updateBody.TabletDeviceEnabled = &tabletDevice
	}

	if len(tags) > 0 {
		tagString := vmGetTagsString(d)
		updateBody.Tags = &tagString
	}

	//nolint:gosimple
	if template != dvTemplate {
		updateBody.Template = &template
	}

	if len(vga) > 0 {
		vgaDevice, err := vmGetVGADeviceObject(d)
		if err != nil {
			return diag.FromErr(err)
		}

		updateBody.VGADevice = vgaDevice
	}

	hookScript := d.Get(mkHookScriptFileID).(string)
	currentHookScript := vmConfig.HookScript

	if len(hookScript) > 0 {
		updateBody.HookScript = &hookScript
	} else if currentHookScript != nil {
		del = append(del, "hookscript")
	}

	updateBody.Delete = del

	e = vmAPI.UpdateVM(ctx, updateBody)
	if e != nil {
		return diag.FromErr(e)
	}

	vmConfig, e = vmAPI.GetVM(ctx)
	if e != nil {
		if strings.Contains(e.Error(), "HTTP 404") ||
			(strings.Contains(e.Error(), "HTTP 500") && strings.Contains(e.Error(), "does not exist")) {
			d.SetId("")

			return nil
		}

		return diag.FromErr(e)
	}

	allDiskInfo, err := createDisks(ctx, vmConfig, d, vmAPI)
	if err != nil {
		return diag.FromErr(err)
	}

	efiDisk := d.Get(mkEFIDisk).([]interface{})
	efiDiskInfo := vmGetEfiDisk(d, nil) // from the resource config

	for i := range efiDisk {
		diskBlock := efiDisk[i].(map[string]interface{})
		diskInterface := "efidisk0"
		dataStoreID := diskBlock[mkEFIDiskDatastoreID].(string)
		efiType := diskBlock[mkEFIDiskType].(string)

		currentDiskInfo := vmConfig.EFIDisk
		configuredDiskInfo := efiDiskInfo

		if currentDiskInfo == nil {
			diskUpdateBody := &vms.UpdateRequestBody{}

			diskUpdateBody.EFIDisk = configuredDiskInfo

			e = vmAPI.UpdateVM(ctx, diskUpdateBody)
			if e != nil {
				return diag.FromErr(e)
			}

			continue
		}

		if &efiType != currentDiskInfo.Type {
			return diag.Errorf(
				"resizing of efidisks is not supported.",
			)
		}

		deleteOriginalDisk := types.CustomBool(true)

		diskMoveBody := &vms.MoveDiskRequestBody{
			DeleteOriginalDisk: &deleteOriginalDisk,
			Disk:               diskInterface,
			TargetStorage:      dataStoreID,
		}

		moveDisk := false

		if dataStoreID != "" {
			moveDisk = true

			if allDiskInfo[diskInterface] != nil {
				fileIDParts := strings.Split(allDiskInfo[diskInterface].FileVolume, ":")
				moveDisk = dataStoreID != fileIDParts[0]
			}
		}

		if moveDisk {
			moveDiskTimeout := d.Get(mkTimeoutMoveDisk).(int)

			e = vmAPI.MoveVMDisk(ctx, diskMoveBody, moveDiskTimeout)
			if e != nil {
				return diag.FromErr(e)
			}
		}
	}

	tpmState := d.Get(mkTPMState).([]interface{})
	tpmStateInfo := vmGetTPMState(d, nil) // from the resource config

	for i := range tpmState {
		diskBlock := tpmState[i].(map[string]interface{})
		diskInterface := "tpmstate0"
		dataStoreID := diskBlock[mkTPMStateDatastoreID].(string)

		currentTPMState := vmConfig.TPMState
		configuredTPMStateInfo := tpmStateInfo

		if currentTPMState == nil {
			diskUpdateBody := &vms.UpdateRequestBody{}

			diskUpdateBody.TPMState = configuredTPMStateInfo

			e = vmAPI.UpdateVM(ctx, diskUpdateBody)
			if e != nil {
				return diag.FromErr(e)
			}

			continue
		}

		deleteOriginalDisk := types.CustomBool(true)

		diskMoveBody := &vms.MoveDiskRequestBody{
			DeleteOriginalDisk: &deleteOriginalDisk,
			Disk:               diskInterface,
			TargetStorage:      dataStoreID,
		}

		moveDisk := false

		if dataStoreID != "" {
			moveDisk = true

			if allDiskInfo[diskInterface] != nil {
				fileIDParts := strings.Split(allDiskInfo[diskInterface].FileVolume, ":")
				moveDisk = dataStoreID != fileIDParts[0]
			}
		}

		if moveDisk {
			moveDiskTimeout := d.Get(mkTimeoutMoveDisk).(int)

			e = vmAPI.MoveVMDisk(ctx, diskMoveBody, moveDiskTimeout)
			if e != nil {
				return diag.FromErr(e)
			}
		}
	}

	return vmCreateStart(ctx, d, m)
}

func vmCreateCustom(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(proxmoxtf.ProviderConfiguration)

	api, err := config.GetClient()
	if err != nil {
		return diag.FromErr(err)
	}

	resource := VM()

	acpi := types.CustomBool(d.Get(mkACPI).(bool))

	customAgent, err := customAgent(d, resource)
	if err != nil {
		return diag.FromErr(err)
	}

	kvmArguments := d.Get(mkKVMArguments).(string)

	audioDevices := vmGetAudioDeviceList(d)

	bios := d.Get(mkBIOS).(string)

	cdromBlock, err := structure.GetSchemaBlock(
		resource,
		d,
		[]string{mkCDROM},
		0,
		true,
	)
	if err != nil {
		return diag.FromErr(err)
	}

	cdromEnabled := cdromBlock[mkCDROMEnabled].(bool)
	cdromFileID := cdromBlock[mkCDROMFileID].(string)
	cdromInterface := cdromBlock[mkCDROMInterface].(string)

	cdromCloudInitEnabled := false
	cdromCloudInitFileID := ""
	cdromCloudInitInterface := ""

	if cdromFileID == "" {
		cdromFileID = "cdrom"
	}

	cpuBlock, err := structure.GetSchemaBlock(
		resource,
		d,
		[]string{mkCPU},
		0,
		true,
	)
	if err != nil {
		return diag.FromErr(err)
	}

	cpuArchitecture := cpuBlock[mkCPUArchitecture].(string)
	cpuCores := cpuBlock[mkCPUCores].(int)
	cpuFlags := cpuBlock[mkCPUFlags].([]interface{})
	cpuHotplugged := cpuBlock[mkCPUHotplugged].(int)
	cpuLimit := cpuBlock[mkCPULimit].(int)
	cpuSockets := cpuBlock[mkCPUSockets].(int)
	cpuNUMA := types.CustomBool(cpuBlock[mkCPUNUMA].(bool))
	cpuType := cpuBlock[mkCPUType].(string)
	cpuUnits := cpuBlock[mkCPUUnits].(int)

	description := d.Get(mkDescription).(string)

	var efiDisk *vms.CustomEFIDisk

	efiDiskBlock := d.Get(mkEFIDisk).([]interface{})
	if len(efiDiskBlock) > 0 {
		block := efiDiskBlock[0].(map[string]interface{})

		datastoreID, _ := block[mkEFIDiskDatastoreID].(string)
		fileFormat, _ := block[mkEFIDiskFileFormat].(string)
		efiType, _ := block[mkEFIDiskType].(string)
		preEnrolledKeys := types.CustomBool(block[mkEFIDiskPreEnrolledKeys].(bool))

		if fileFormat == "" {
			fileFormat = dvEFIDiskFileFormat
		}

		efiDisk = &vms.CustomEFIDisk{
			Type:            &efiType,
			FileVolume:      fmt.Sprintf("%s:1", datastoreID),
			Format:          &fileFormat,
			PreEnrolledKeys: &preEnrolledKeys,
		}
	}

	var tpmState *vms.CustomTPMState

	tpmStateBlock := d.Get(mkTPMState).([]interface{})
	if len(tpmStateBlock) > 0 {
		block := tpmStateBlock[0].(map[string]interface{})

		datastoreID, _ := block[mkTPMStateDatastoreID].(string)
		version, _ := block[mkTPMStateVersion].(string)

		if version == "" {
			version = dvTPMStateVersion
		}

		tpmState = &vms.CustomTPMState{
			FileVolume: fmt.Sprintf("%s:1", datastoreID),
			Version:    &version,
		}
	}

	initializationConfig := vmGetCloudInitConfig(d)

	if initializationConfig != nil {
		initialization := d.Get(mkInitialization).([]interface{})
		initializationBlock := initialization[0].(map[string]interface{})
		initializationDatastoreID := initializationBlock[mkInitializationDatastoreID].(string)

		cdromCloudInitEnabled = true
		cdromCloudInitFileID = fmt.Sprintf("%s:cloudinit", initializationDatastoreID)

		cdromCloudInitInterface = initializationBlock[mkInitializationInterface].(string)
		if cdromCloudInitInterface == "" {
			cdromCloudInitInterface = "ide2"
		}
	}

	pciDeviceObjects := vmGetHostPCIDeviceObjects(d)

	usbDeviceObjects := vmGetHostUSBDeviceObjects(d)

	keyboardLayout := d.Get(mkKeyboardLayout).(string)

	memoryBlock, err := structure.GetSchemaBlock(
		resource,
		d,
		[]string{mkMemory},
		0,
		true,
	)
	if err != nil {
		return diag.FromErr(err)
	}

	memoryDedicated := memoryBlock[mkMemoryDedicated].(int)
	memoryFloating := memoryBlock[mkMemoryFloating].(int)
	memoryShared := memoryBlock[mkMemoryShared].(int)

	machine := d.Get(mkMachine).(string)
	name := d.Get(mkName).(string)
	tags := d.Get(mkTags).([]interface{})

	networkDeviceObjects := vmGetNetworkDeviceObjects(d)

	nodeName := d.Get(mkNodeName).(string)

	operatingSystem, err := structure.GetSchemaBlock(
		resource,
		d,
		[]string{mkOperatingSystem},
		0,
		true,
	)
	if err != nil {
		return diag.FromErr(err)
	}

	operatingSystemType := operatingSystem[mkOperatingSystemType].(string)

	poolID := d.Get(mkPoolID).(string)

	serialDevices := vmGetSerialDeviceList(d)

	smbios := vmGetSMBIOS(d)

	startupOrder := vmGetStartupOrder(d)

	onBoot := types.CustomBool(d.Get(mkOnBoot).(bool))
	tabletDevice := types.CustomBool(d.Get(mkTabletDevice).(bool))
	template := types.CustomBool(d.Get(mkTemplate).(bool))

	vgaDevice, err := vmGetVGADeviceObject(d)
	if err != nil {
		return diag.FromErr(err)
	}

	vmIDUntyped, hasVMID := d.GetOk(mkVMID)
	vmID := vmIDUntyped.(int)

	if !hasVMID {
		vmIDNew, e := api.Cluster().GetVMID(ctx)
		if e != nil {
			return diag.FromErr(e)
		}

		vmID = *vmIDNew
		e = d.Set(mkVMID, vmID)

		if e != nil {
			return diag.FromErr(e)
		}
	}

	var memorySharedObject *vms.CustomSharedMemory

	var bootOrderConverted []string
	if cdromEnabled {
		bootOrderConverted = []string{cdromInterface}
	}

	planDisks, err := getStorageDevicesFromResource(d)
	if err != nil {
		return diag.FromErr(err)
	}

	bootOrder := d.Get(mkBootOrder).([]interface{})
	if len(bootOrder) == 0 {
		for _, disk := range planDisks {
			if *disk.Interface == "sata0" {
				bootOrderConverted = append(bootOrderConverted, "sata0")
			}

			if *disk.Interface == "scsi0" {
				bootOrderConverted = append(bootOrderConverted, "scsi0")
			}

			if *disk.Interface == "virtio0" {
				bootOrderConverted = append(bootOrderConverted, "virtio0")
			}
		}

		if networkDeviceObjects != nil {
			bootOrderConverted = append(bootOrderConverted, "net0")
		}
	} else {
		bootOrderConverted = make([]string, len(bootOrder))
		for i, device := range bootOrder {
			bootOrderConverted[i] = device.(string)
		}
	}

	cpuFlagsConverted := make([]string, len(cpuFlags))
	for fi, flag := range cpuFlags {
		cpuFlagsConverted[fi] = flag.(string)
	}

	ideDevice2Media := "cdrom"
	ideDevices := vms.CustomStorageDevices{
		cdromCloudInitInterface: &vms.CustomStorageDevice{
			Enabled:    cdromCloudInitEnabled,
			FileVolume: cdromCloudInitFileID,
			Media:      &ideDevice2Media,
		},
		cdromInterface: &vms.CustomStorageDevice{
			Enabled:    cdromEnabled,
			FileVolume: cdromFileID,
			Media:      &ideDevice2Media,
		},
	}

	if memoryShared > 0 {
		memorySharedName := fmt.Sprintf("vm-%d-ivshmem", vmID)
		memorySharedObject = &vms.CustomSharedMemory{
			Name: &memorySharedName,
			Size: memoryShared,
		}
	}

	scsiHardware := d.Get(mkSCSIHardware).(string)

	createBody := &vms.CreateRequestBody{
		ACPI:         &acpi,
		Agent:        customAgent,
		AudioDevices: audioDevices,
		BIOS:         &bios,
		Boot: &vms.CustomBoot{
			Order: &bootOrderConverted,
		},
		CloudInitConfig: initializationConfig,
		CPUCores:        &cpuCores,
		CPUEmulation: &vms.CustomCPUEmulation{
			Flags: &cpuFlagsConverted,
			Type:  cpuType,
		},
		CPUSockets:          &cpuSockets,
		CPUUnits:            &cpuUnits,
		DedicatedMemory:     &memoryDedicated,
		EFIDisk:             efiDisk,
		TPMState:            tpmState,
		FloatingMemory:      &memoryFloating,
		IDEDevices:          ideDevices,
		KeyboardLayout:      &keyboardLayout,
		NetworkDevices:      networkDeviceObjects,
		NUMAEnabled:         &cpuNUMA,
		OSType:              &operatingSystemType,
		PCIDevices:          pciDeviceObjects,
		SCSIHardware:        &scsiHardware,
		SerialDevices:       serialDevices,
		SharedMemory:        memorySharedObject,
		StartOnBoot:         &onBoot,
		SMBIOS:              smbios,
		StartupOrder:        startupOrder,
		TabletDeviceEnabled: &tabletDevice,
		Template:            &template,
		USBDevices:          usbDeviceObjects,
		VGADevice:           vgaDevice,
		VMID:                &vmID,
	}

	sataDeviceObjects := planDisks.ByStorageInterface("sata")
	if len(sataDeviceObjects) > 0 {
		createBody.SATADevices = sataDeviceObjects
	}

	scsiDeviceObjects := planDisks.ByStorageInterface("scsi")
	if len(scsiDeviceObjects) > 0 {
		createBody.SCSIDevices = scsiDeviceObjects
	}

	virtioDeviceObjects := planDisks.ByStorageInterface("virtio")
	if len(virtioDeviceObjects) > 0 {
		createBody.VirtualIODevices = virtioDeviceObjects
	}

	// Only the root account is allowed to change the CPU architecture, which makes this check necessary.
	if api.API().IsRootTicket() ||
		cpuArchitecture != dvCPUArchitecture {
		createBody.CPUArchitecture = &cpuArchitecture
	}

	if cpuHotplugged > 0 {
		createBody.VirtualCPUCount = &cpuHotplugged
	}

	if cpuLimit > 0 {
		createBody.CPULimit = &cpuLimit
	}

	if description != "" {
		createBody.Description = &description
	}

	if len(tags) > 0 {
		tagsString := vmGetTagsString(d)
		createBody.Tags = &tagsString
	}

	if kvmArguments != "" {
		createBody.KVMArguments = &kvmArguments
	}

	if machine != "" {
		createBody.Machine = &machine
	}

	if name != "" {
		createBody.Name = &name
	}

	if poolID != "" {
		createBody.PoolID = &poolID
	}

	hookScript := d.Get(mkHookScriptFileID).(string)
	if len(hookScript) > 0 {
		createBody.HookScript = &hookScript
	}

	createTimeout := d.Get(mkTimeoutClone).(int)

	err = api.Node(nodeName).VM(0).CreateVM(ctx, createBody, createTimeout)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(vmID))

	err = vmImportCustomDisks(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	return vmCreateStart(ctx, d, m)
}

func vmCreateStart(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	started := d.Get(mkStarted).(bool)
	template := d.Get(mkTemplate).(bool)
	reboot := d.Get(mkRebootAfterCreation).(bool)

	if !started || template {
		return vmRead(ctx, d, m)
	}

	config := m.(proxmoxtf.ProviderConfiguration)

	api, err := config.GetClient()
	if err != nil {
		return diag.FromErr(err)
	}

	nodeName := d.Get(mkNodeName).(string)

	vmID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	vmAPI := api.Node(nodeName).VM(vmID)

	// Start the virtual machine and wait for it to reach a running state before continuing.
	if diags := vmStart(ctx, vmAPI, d); diags != nil {
		return diags
	}

	if reboot {
		rebootTimeout := d.Get(mkTimeoutReboot).(int)

		err := vmAPI.RebootVM(
			ctx,
			&vms.RebootRequestBody{
				Timeout: &rebootTimeout,
			},
			rebootTimeout+30,
		)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return vmRead(ctx, d, m)
}

func vmGetAudioDeviceList(d *schema.ResourceData) vms.CustomAudioDevices {
	devices := d.Get(mkAudioDevice).([]interface{})
	list := make(vms.CustomAudioDevices, len(devices))

	for i, v := range devices {
		block := v.(map[string]interface{})

		device, _ := block[mkAudioDeviceDevice].(string)
		driver, _ := block[mkAudioDeviceDriver].(string)
		enabled, _ := block[mkAudioDeviceEnabled].(bool)

		list[i].Device = device
		list[i].Driver = &driver
		list[i].Enabled = enabled
	}

	return list
}

func vmGetAudioDeviceValidator() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.StringInSlice([]string{
		"AC97",
		"ich9-intel-hda",
		"intel-hda",
	}, false))
}

func vmGetAudioDriverValidator() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.StringInSlice([]string{
		"spice",
	}, false))
}

func vmGetCloudInitConfig(d *schema.ResourceData) *vms.CustomCloudInitConfig {
	initialization := d.Get(mkInitialization).([]interface{})

	if len(initialization) == 0 || initialization[0] == nil {
		return nil
	}

	var initializationConfig *vms.CustomCloudInitConfig

	initializationBlock := initialization[0].(map[string]interface{})
	initializationConfig = &vms.CustomCloudInitConfig{}
	initializationDNS := initializationBlock[mkInitializationDNS].([]interface{})

	if len(initializationDNS) > 0 && initializationDNS[0] != nil {
		initializationDNSBlock := initializationDNS[0].(map[string]interface{})
		domain := initializationDNSBlock[mkInitializationDNSDomain].(string)

		if domain != "" {
			initializationConfig.SearchDomain = &domain
		}

		servers := initializationDNSBlock[mkInitializationDNSServers].([]interface{})
		deprecatedServer := initializationDNSBlock[mkInitializationDNSServer].(string)

		if len(servers) > 0 {
			nameserver := strings.Join(utils.ConvertToStringSlice(servers), " ")

			initializationConfig.Nameserver = &nameserver
		} else if deprecatedServer != "" {
			initializationConfig.Nameserver = &deprecatedServer
		}
	}

	initializationIPConfig := initializationBlock[mkInitializationIPConfig].([]interface{})
	initializationConfig.IPConfig = make(
		[]vms.CustomCloudInitIPConfig,
		len(initializationIPConfig),
	)

	for i, c := range initializationIPConfig {
		configBlock := c.(map[string]interface{})
		ipv4 := configBlock[mkInitializationIPConfigIPv4].([]interface{})

		if len(ipv4) > 0 {
			ipv4Block := ipv4[0].(map[string]interface{})
			ipv4Address := ipv4Block[mkInitializationIPConfigIPv4Address].(string)

			if ipv4Address != "" {
				initializationConfig.IPConfig[i].IPv4 = &ipv4Address
			}

			ipv4Gateway := ipv4Block[mkInitializationIPConfigIPv4Gateway].(string)

			if ipv4Gateway != "" {
				initializationConfig.IPConfig[i].GatewayIPv4 = &ipv4Gateway
			}
		}

		ipv6 := configBlock[mkInitializationIPConfigIPv6].([]interface{})

		if len(ipv6) > 0 {
			ipv6Block := ipv6[0].(map[string]interface{})
			ipv6Address := ipv6Block[mkInitializationIPConfigIPv6Address].(string)

			if ipv6Address != "" {
				initializationConfig.IPConfig[i].IPv6 = &ipv6Address
			}

			ipv6Gateway := ipv6Block[mkInitializationIPConfigIPv6Gateway].(string)

			if ipv6Gateway != "" {
				initializationConfig.IPConfig[i].GatewayIPv6 = &ipv6Gateway
			}
		}
	}

	initializationUserAccount := initializationBlock[mkInitializationUserAccount].([]interface{})

	if len(initializationUserAccount) > 0 {
		initializationUserAccountBlock := initializationUserAccount[0].(map[string]interface{})
		keys := initializationUserAccountBlock[mkInitializationUserAccountKeys].([]interface{})

		if len(keys) > 0 {
			sshKeys := make(vms.CustomCloudInitSSHKeys, len(keys))

			for i, k := range keys {
				sshKeys[i] = k.(string)
			}

			initializationConfig.SSHKeys = &sshKeys
		}

		password := initializationUserAccountBlock[mkInitializationUserAccountPassword].(string)

		if password != "" {
			initializationConfig.Password = &password
		}

		username := initializationUserAccountBlock[mkInitializationUserAccountUsername].(string)

		initializationConfig.Username = &username
	}

	initializationUserDataFileID := initializationBlock[mkInitializationUserDataFileID].(string)

	if initializationUserDataFileID != "" {
		initializationConfig.Files = &vms.CustomCloudInitFiles{
			UserVolume: &initializationUserDataFileID,
		}
	}

	initializationVendorDataFileID := initializationBlock[mkInitializationVendorDataFileID].(string)

	if initializationVendorDataFileID != "" {
		if initializationConfig.Files == nil {
			initializationConfig.Files = &vms.CustomCloudInitFiles{}
		}

		initializationConfig.Files.VendorVolume = &initializationVendorDataFileID
	}

	initializationNetworkDataFileID := initializationBlock[mkInitializationNetworkDataFileID].(string)

	if initializationNetworkDataFileID != "" {
		if initializationConfig.Files == nil {
			initializationConfig.Files = &vms.CustomCloudInitFiles{}
		}

		initializationConfig.Files.NetworkVolume = &initializationNetworkDataFileID
	}

	initializationMetaDataFileID := initializationBlock[mkInitializationMetaDataFileID].(string)

	if initializationMetaDataFileID != "" {
		if initializationConfig.Files == nil {
			initializationConfig.Files = &vms.CustomCloudInitFiles{}
		}

		initializationConfig.Files.MetaVolume = &initializationMetaDataFileID
	}

	initializationType := initializationBlock[mkInitializationType].(string)

	if initializationType != "" {
		initializationConfig.Type = &initializationType
	}

	return initializationConfig
}

func vmGetCPUArchitectureValidator() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.StringInSlice([]string{
		"aarch64",
		"x86_64",
	}, false))
}

func vmGetEfiDisk(d *schema.ResourceData, disk []interface{}) *vms.CustomEFIDisk {
	var efiDisk []interface{}

	if disk != nil {
		efiDisk = disk
	} else {
		efiDisk = d.Get(mkEFIDisk).([]interface{})
	}

	var efiDiskConfig *vms.CustomEFIDisk

	if len(efiDisk) > 0 {
		efiDiskConfig = &vms.CustomEFIDisk{}

		block := efiDisk[0].(map[string]interface{})
		datastoreID, _ := block[mkEFIDiskDatastoreID].(string)
		fileFormat, _ := block[mkEFIDiskFileFormat].(string)
		efiType, _ := block[mkEFIDiskType].(string)
		preEnrolledKeys := types.CustomBool(block[mkEFIDiskPreEnrolledKeys].(bool))

		// use the special syntax STORAGE_ID:SIZE_IN_GiB to allocate a new volume.
		// NB SIZE_IN_GiB is ignored, see docs for more info.
		efiDiskConfig.FileVolume = fmt.Sprintf("%s:1", datastoreID)
		efiDiskConfig.Format = &fileFormat
		efiDiskConfig.Type = &efiType
		efiDiskConfig.PreEnrolledKeys = &preEnrolledKeys
	}

	return efiDiskConfig
}

func vmGetEfiDiskAsStorageDevice(d *schema.ResourceData, disk []interface{}) (*vms.CustomStorageDevice, error) {
	efiDisk := vmGetEfiDisk(d, disk)

	var storageDevice *vms.CustomStorageDevice

	if efiDisk != nil {
		id := "0"
		baseDiskInterface := "efidisk"
		diskInterface := fmt.Sprint(baseDiskInterface, id)

		storageDevice = &vms.CustomStorageDevice{
			Enabled:     true,
			FileVolume:  efiDisk.FileVolume,
			Format:      efiDisk.Format,
			Interface:   &diskInterface,
			DatastoreID: &id,
		}

		if efiDisk.Type != nil {
			ds, err := types.ParseDiskSize(*efiDisk.Type)
			if err != nil {
				return nil, fmt.Errorf("invalid efi disk type: %s", err.Error())
			}

			storageDevice.Size = &ds
		}
	}

	return storageDevice, nil
}

func vmGetTPMState(d *schema.ResourceData, disk []interface{}) *vms.CustomTPMState {
	var tpmState []interface{}

	if disk != nil {
		tpmState = disk
	} else {
		tpmState = d.Get(mkTPMState).([]interface{})
	}

	var tpmStateConfig *vms.CustomTPMState

	if len(tpmState) > 0 {
		tpmStateConfig = &vms.CustomTPMState{}

		block := tpmState[0].(map[string]interface{})
		datastoreID, _ := block[mkTPMStateDatastoreID].(string)
		version, _ := block[mkTPMStateVersion].(string)

		// use the special syntax STORAGE_ID:SIZE_IN_GiB to allocate a new volume.
		// NB SIZE_IN_GiB is ignored, see docs for more info.
		tpmStateConfig.FileVolume = fmt.Sprintf("%s:1", datastoreID)
		tpmStateConfig.Version = &version
	}

	return tpmStateConfig
}

func vmGetTPMStateAsStorageDevice(d *schema.ResourceData, disk []interface{}) *vms.CustomStorageDevice {
	tpmState := vmGetTPMState(d, disk)

	var storageDevice *vms.CustomStorageDevice

	if tpmState != nil {
		id := "0"
		baseDiskInterface := "tpmstate"
		diskInterface := fmt.Sprint(baseDiskInterface, id)

		storageDevice = &vms.CustomStorageDevice{
			Enabled:     true,
			FileVolume:  tpmState.FileVolume,
			Interface:   &diskInterface,
			DatastoreID: &id,
		}
	}

	return storageDevice
}

func vmGetHostPCIDeviceObjects(d *schema.ResourceData) vms.CustomPCIDevices {
	pciDevice := d.Get(mkHostPCI).([]interface{})
	pciDeviceObjects := make(vms.CustomPCIDevices, len(pciDevice))

	for i, pciDeviceEntry := range pciDevice {
		block := pciDeviceEntry.(map[string]interface{})

		ids, _ := block[mkHostPCIDeviceID].(string)
		mdev, _ := block[mkHostPCIDeviceMDev].(string)
		pcie := types.CustomBool(block[mkHostPCIDevicePCIE].(bool))
		rombar := types.CustomBool(
			block[mkHostPCIDeviceROMBAR].(bool),
		)
		romfile, _ := block[mkHostPCIDeviceROMFile].(string)
		xvga := types.CustomBool(block[mkHostPCIDeviceXVGA].(bool))
		mapping, _ := block[mkHostPCIDeviceMapping].(string)

		device := vms.CustomPCIDevice{
			PCIExpress: &pcie,
			ROMBAR:     &rombar,
			XVGA:       &xvga,
		}

		if ids != "" {
			dIDs := strings.Split(ids, ";")
			device.DeviceIDs = &dIDs
		}

		if mdev != "" {
			device.MDev = &mdev
		}

		if romfile != "" {
			device.ROMFile = &romfile
		}

		if mapping != "" {
			device.Mapping = &mapping
		}

		pciDeviceObjects[i] = device
	}

	return pciDeviceObjects
}

func vmGetHostUSBDeviceObjects(d *schema.ResourceData) vms.CustomUSBDevices {
	usbDevice := d.Get(mkHostUSB).([]interface{})
	usbDeviceObjects := make(vms.CustomUSBDevices, len(usbDevice))

	for i, usbDeviceEntry := range usbDevice {
		block := usbDeviceEntry.(map[string]interface{})

		host, _ := block[mkHostUSBDevice].(string)
		usb3 := types.CustomBool(block[mkHostUSBDeviceUSB3].(bool))
		mapping, _ := block[mkHostUSBDeviceMapping].(string)

		device := vms.CustomUSBDevice{
			HostDevice: &host,
			USB3:       &usb3,
		}
		if mapping != "" {
			device.Mapping = &mapping
		}

		usbDeviceObjects[i] = device
	}

	return usbDeviceObjects
}

func vmGetNetworkDeviceObjects(d *schema.ResourceData) vms.CustomNetworkDevices {
	networkDevice := d.Get(mkNetworkDevice).([]interface{})
	networkDeviceObjects := make(vms.CustomNetworkDevices, len(networkDevice))

	for i, networkDeviceEntry := range networkDevice {
		block := networkDeviceEntry.(map[string]interface{})

		bridge := block[mkNetworkDeviceBridge].(string)
		enabled := block[mkNetworkDeviceEnabled].(bool)
		firewall := types.CustomBool(block[mkNetworkDeviceFirewall].(bool))
		macAddress := block[mkNetworkDeviceMACAddress].(string)
		model := block[mkNetworkDeviceModel].(string)
		queues := block[mkNetworkDeviceQueues].(int)
		rateLimit := block[mkNetworkDeviceRateLimit].(float64)
		vlanID := block[mkNetworkDeviceVLANID].(int)
		mtu := block[mkNetworkDeviceMTU].(int)

		device := vms.CustomNetworkDevice{
			Enabled:  enabled,
			Firewall: &firewall,
			Model:    model,
		}

		if bridge != "" {
			device.Bridge = &bridge
		}

		if macAddress != "" {
			device.MACAddress = &macAddress
		}

		if queues != 0 {
			device.Queues = &queues
		}

		if rateLimit != 0 {
			device.RateLimit = &rateLimit
		}

		if vlanID != 0 {
			device.Tag = &vlanID
		}

		if mtu != 0 {
			device.MTU = &mtu
		}

		networkDeviceObjects[i] = device
	}

	return networkDeviceObjects
}

func vmGetOperatingSystemTypeValidator() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.StringInSlice([]string{
		"l24",
		"l26",
		"other",
		"solaris",
		"w2k",
		"w2k3",
		"w2k8",
		"win7",
		"win8",
		"win10",
		"win11",
		"wvista",
		"wxp",
	}, false))
}

func vmGetSerialDeviceList(d *schema.ResourceData) vms.CustomSerialDevices {
	device := d.Get(mkSerialDevice).([]interface{})
	list := make(vms.CustomSerialDevices, len(device))

	for i, v := range device {
		block := v.(map[string]interface{})

		device, _ := block[mkSerialDeviceDevice].(string)

		list[i] = device
	}

	return list
}

func vmGetSMBIOS(d *schema.ResourceData) *vms.CustomSMBIOS {
	smbiosSections := d.Get(mkSMBIOS).([]interface{})
	//nolint:nestif
	if len(smbiosSections) > 0 {
		smbiosBlock := smbiosSections[0].(map[string]interface{})
		b64 := types.CustomBool(true)
		family, _ := smbiosBlock[mkSMBIOSFamily].(string)
		manufacturer, _ := smbiosBlock[mkSMBIOSManufacturer].(string)
		product, _ := smbiosBlock[mkSMBIOSProduct].(string)
		serial, _ := smbiosBlock[mkSMBIOSSerial].(string)
		sku, _ := smbiosBlock[mkSMBIOSSKU].(string)
		version, _ := smbiosBlock[mkSMBIOSVersion].(string)
		uid, _ := smbiosBlock[mkSMBIOSUUID].(string)

		smbios := vms.CustomSMBIOS{
			Base64: &b64,
		}

		if family != "" {
			v := base64.StdEncoding.EncodeToString([]byte(family))
			smbios.Family = &v
		}

		if manufacturer != "" {
			v := base64.StdEncoding.EncodeToString([]byte(manufacturer))
			smbios.Manufacturer = &v
		}

		if product != "" {
			v := base64.StdEncoding.EncodeToString([]byte(product))
			smbios.Product = &v
		}

		if serial != "" {
			v := base64.StdEncoding.EncodeToString([]byte(serial))
			smbios.Serial = &v
		}

		if sku != "" {
			v := base64.StdEncoding.EncodeToString([]byte(sku))
			smbios.SKU = &v
		}

		if version != "" {
			v := base64.StdEncoding.EncodeToString([]byte(version))
			smbios.Version = &v
		}

		if uid != "" {
			smbios.UUID = &uid
		}

		if smbios.UUID == nil || *smbios.UUID == "" {
			smbios.UUID = types.StrPtr(uuid.New().String())
		}

		return &smbios
	}

	return nil
}

func vmGetStartupOrder(d *schema.ResourceData) *vms.CustomStartupOrder {
	startup := d.Get(mkStartup).([]interface{})
	if len(startup) > 0 {
		startupBlock := startup[0].(map[string]interface{})
		startupOrder := startupBlock[mkStartupOrder].(int)
		startupUpDelay := startupBlock[mkStartupUpDelay].(int)
		startupDownDelay := startupBlock[mkStartupDownDelay].(int)

		order := vms.CustomStartupOrder{}

		if startupUpDelay >= 0 {
			order.Up = &startupUpDelay
		}

		if startupDownDelay >= 0 {
			order.Down = &startupDownDelay
		}

		if startupOrder >= 0 {
			order.Order = &startupOrder
		}

		return &order
	}

	return nil
}

func vmGetTagsString(d *schema.ResourceData) string {
	var sanitizedTags []string

	tags := d.Get(mkTags).([]interface{})
	for i := 0; i < len(tags); i++ {
		tag := strings.TrimSpace(tags[i].(string))
		if len(tag) > 0 {
			sanitizedTags = append(sanitizedTags, tag)
		}
	}

	sort.Strings(sanitizedTags)

	return strings.Join(sanitizedTags, ";")
}

func vmGetSerialDeviceValidator() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(func(i interface{}, k string) ([]string, []error) {
		v, ok := i.(string)

		var es []error

		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be string", k))
			return nil, es
		}

		if !strings.HasPrefix(v, "/dev/") && v != "socket" {
			es = append(es, fmt.Errorf("expected %s to be '/dev/*' or 'socket'", k))
			return nil, es
		}

		return nil, es
	})
}

func vmGetVGADeviceObject(d *schema.ResourceData) (*vms.CustomVGADevice, error) {
	resource := VM()

	vgaBlock, err := structure.GetSchemaBlock(
		resource,
		d,
		[]string{mkVGA},
		0,
		true,
	)
	if err != nil {
		return nil, fmt.Errorf("error reading VGA block: %w", err)
	}

	vgaEnabled := types.CustomBool(vgaBlock[mkVGAEnabled].(bool))
	vgaMemory := vgaBlock[mkVGAMemory].(int)
	vgaType := vgaBlock[mkVGAType].(string)

	vgaDevice := &vms.CustomVGADevice{}

	if vgaEnabled {
		if vgaMemory > 0 {
			vgaDevice.Memory = &vgaMemory
		}

		vgaDevice.Type = &vgaType
	} else {
		vgaType = "none"

		vgaDevice = &vms.CustomVGADevice{
			Type: &vgaType,
		}
	}

	return vgaDevice, nil
}

func vmRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(proxmoxtf.ProviderConfiguration)

	api, err := config.GetClient()
	if err != nil {
		return diag.FromErr(err)
	}

	vmID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	vmNodeName, err := api.Cluster().GetVMNodeName(ctx, vmID)
	if err != nil {
		if errors.Is(err, cluster.ErrVMDoesNotExist) {
			d.SetId("")

			return nil
		}

		return diag.FromErr(err)
	}

	if vmNodeName != d.Get(mkNodeName) {
		err = d.Set(mkNodeName, vmNodeName)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	nodeName := d.Get(mkNodeName).(string)

	vmAPI := api.Node(nodeName).VM(vmID)

	// Retrieve the entire configuration in order to compare it to the state.
	vmConfig, err := vmAPI.GetVM(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "HTTP 404") ||
			(strings.Contains(err.Error(), "HTTP 500") && strings.Contains(err.Error(), "does not exist")) {
			d.SetId("")

			return nil
		}

		return diag.FromErr(err)
	}

	vmStatus, err := vmAPI.GetVMStatus(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	return vmReadCustom(ctx, d, m, vmID, vmConfig, vmStatus)
}

// orderedListFromMap generates a list from a map's values. The values are sorted based on the map's keys.
func orderedListFromMap(inputMap map[string]interface{}) []interface{} {
	itemCount := len(inputMap)
	keyList := make([]string, itemCount)
	i := 0

	for key := range inputMap {
		keyList[i] = key
		i++
	}

	sort.Strings(keyList)

	orderedList := make([]interface{}, itemCount)
	for i, k := range keyList {
		orderedList[i] = inputMap[k]
	}

	return orderedList
}

func vmReadCustom(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
	vmID int,
	vmConfig *vms.GetResponseData,
	vmStatus *vms.GetStatusResponseData,
) diag.Diagnostics {
	config := m.(proxmoxtf.ProviderConfiguration)

	api, err := config.GetClient()
	if err != nil {
		return diag.FromErr(err)
	}

	diags := vmReadPrimitiveValues(d, vmConfig, vmStatus)
	if diags.HasError() {
		return diags
	}

	// Fix terraform.tfstate, by replacing '-1' (the old default value) with actual vm_id value
	if storedVMID := d.Get(mkVMID).(int); storedVMID == -1 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary: fmt.Sprintf("VM %s has stored legacy vm_id %d, setting vm_id to its correct value %d.",
				d.Id(), storedVMID, vmID),
		})

		err = d.Set(mkVMID, vmID)
		diags = append(diags, diag.FromErr(err)...)
	}

	nodeName := d.Get(mkNodeName).(string)
	clone := d.Get(mkClone).([]interface{})

	err = setAgent(d, len(clone) > 0, vmConfig)
	diags = append(diags, diag.FromErr(err)...)

	// Compare the audio devices to those stored in the state.
	currentAudioDevice := d.Get(mkAudioDevice).([]interface{})

	audioDevices := make([]interface{}, 1)
	audioDevicesArray := []*vms.CustomAudioDevice{
		vmConfig.AudioDevice,
	}
	audioDevicesCount := 0

	for adi, ad := range audioDevicesArray {
		m := map[string]interface{}{}

		if ad != nil {
			m[mkAudioDeviceDevice] = ad.Device

			if ad.Driver != nil {
				m[mkAudioDeviceDriver] = *ad.Driver
			} else {
				m[mkAudioDeviceDriver] = ""
			}

			m[mkAudioDeviceEnabled] = true

			audioDevicesCount = adi + 1
		} else {
			m[mkAudioDeviceDevice] = ""
			m[mkAudioDeviceDriver] = ""
			m[mkAudioDeviceEnabled] = false
		}

		audioDevices[adi] = m
	}

	if len(clone) == 0 || len(currentAudioDevice) > 0 {
		err := d.Set(mkAudioDevice, audioDevices[:audioDevicesCount])
		diags = append(diags, diag.FromErr(err)...)
	}

	// Compare the IDE devices to the CD-ROM configurations stored in the state.
	currentInterface := dvCDROMInterface

	currentCDROM := d.Get(mkCDROM).([]interface{})
	if len(currentCDROM) > 0 {
		currentBlock := currentCDROM[0].(map[string]interface{})
		currentInterface = currentBlock[mkCDROMInterface].(string)
	}

	cdromIDEDevice := getIdeDevice(vmConfig, currentInterface)

	//nolint:nestif
	if cdromIDEDevice != nil {
		cdrom := make([]interface{}, 1)
		cdromBlock := map[string]interface{}{}

		if len(clone) == 0 || len(currentCDROM) > 0 {
			cdromBlock[mkCDROMEnabled] = cdromIDEDevice.Enabled
			cdromBlock[mkCDROMFileID] = cdromIDEDevice.FileVolume
			cdromBlock[mkCDROMInterface] = currentInterface

			if len(currentCDROM) > 0 {
				currentBlock := currentCDROM[0].(map[string]interface{})

				if currentBlock[mkCDROMFileID] == "" {
					cdromBlock[mkCDROMFileID] = ""
				}

				if currentBlock[mkCDROMEnabled] == false {
					cdromBlock[mkCDROMEnabled] = false
				}
			}

			cdrom[0] = cdromBlock

			err := d.Set(mkCDROM, cdrom)
			diags = append(diags, diag.FromErr(err)...)
		}
	} else {
		err := d.Set(mkCDROM, []interface{}{})
		diags = append(diags, diag.FromErr(err)...)
	}

	// Compare the CPU configuration to the one stored in the state.
	cpu := map[string]interface{}{}

	if vmConfig.CPUArchitecture != nil {
		cpu[mkCPUArchitecture] = *vmConfig.CPUArchitecture
	} else {
		// Default value of "arch" is "" according to the API documentation.
		// However, assume the provider's default value as a workaround when the root account is not being used.
		if !api.API().IsRootTicket() {
			cpu[mkCPUArchitecture] = dvCPUArchitecture
		} else {
			cpu[mkCPUArchitecture] = ""
		}
	}

	if vmConfig.CPUCores != nil {
		cpu[mkCPUCores] = *vmConfig.CPUCores
	} else {
		// Default value of "cores" is "1" according to the API documentation.
		cpu[mkCPUCores] = 1
	}

	if vmConfig.VirtualCPUCount != nil {
		cpu[mkCPUHotplugged] = *vmConfig.VirtualCPUCount
	} else {
		// Default value of "vcpus" is "1" according to the API documentation.
		cpu[mkCPUHotplugged] = 0
	}

	if vmConfig.CPULimit != nil {
		cpu[mkCPULimit] = *vmConfig.CPULimit
	} else {
		// Default value of "cpulimit" is "0" according to the API documentation.
		cpu[mkCPULimit] = 0
	}

	if vmConfig.NUMAEnabled != nil {
		cpu[mkCPUNUMA] = *vmConfig.NUMAEnabled
	} else {
		// Default value of "numa" is "false" according to the API documentation.
		cpu[mkCPUNUMA] = false
	}

	if vmConfig.CPUSockets != nil {
		cpu[mkCPUSockets] = *vmConfig.CPUSockets
	} else {
		// Default value of "sockets" is "1" according to the API documentation.
		cpu[mkCPUSockets] = 1
	}

	if vmConfig.CPUEmulation != nil {
		if vmConfig.CPUEmulation.Flags != nil {
			convertedFlags := make([]interface{}, len(*vmConfig.CPUEmulation.Flags))

			for fi, fv := range *vmConfig.CPUEmulation.Flags {
				convertedFlags[fi] = fv
			}

			cpu[mkCPUFlags] = convertedFlags
		} else {
			cpu[mkCPUFlags] = []interface{}{}
		}

		cpu[mkCPUType] = vmConfig.CPUEmulation.Type
	} else {
		cpu[mkCPUFlags] = []interface{}{}
		// Default value of "cputype" is "qemu64" according to the QEMU documentation.
		cpu[mkCPUType] = "qemu64"
	}

	if vmConfig.CPUUnits != nil {
		cpu[mkCPUUnits] = *vmConfig.CPUUnits
	} else {
		// Default value of "cpuunits" is "1024" according to the API documentation.
		cpu[mkCPUUnits] = 1024
	}

	currentCPU := d.Get(mkCPU).([]interface{})

	if len(clone) > 0 {
		if len(currentCPU) > 0 {
			err := d.Set(mkCPU, []interface{}{cpu})
			diags = append(diags, diag.FromErr(err)...)
		}
	} else if len(currentCPU) > 0 ||
		cpu[mkCPUArchitecture] != dvCPUArchitecture ||
		cpu[mkCPUCores] != dvCPUCores ||
		len(cpu[mkCPUFlags].([]interface{})) > 0 ||
		cpu[mkCPUHotplugged] != dvCPUHotplugged ||
		cpu[mkCPULimit] != dvCPULimit ||
		cpu[mkCPUSockets] != dvCPUSockets ||
		cpu[mkCPUType] != dvCPUType ||
		cpu[mkCPUUnits] != dvCPUUnits {
		err := d.Set(mkCPU, []interface{}{cpu})
		diags = append(diags, diag.FromErr(err)...)
	}

	diags = append(diags, readDisk1(ctx, d, vmConfig, vmID, api, nodeName, clone)...)

	//nolint:nestif
	if vmConfig.EFIDisk != nil {
		efiDisk := map[string]interface{}{}

		fileIDParts := strings.Split(vmConfig.EFIDisk.FileVolume, ":")

		efiDisk[mkEFIDiskDatastoreID] = fileIDParts[0]

		if vmConfig.EFIDisk.Format != nil {
			efiDisk[mkEFIDiskFileFormat] = *vmConfig.EFIDisk.Format
		} else {
			// disk format may not be returned by config API if it is default for the storage, and that may be different
			// from the default qcow2, so we need to read it from the storage API to make sure we have the correct value
			volume, err := api.Node(nodeName).Storage(fileIDParts[0]).GetDatastoreFile(ctx, vmConfig.EFIDisk.FileVolume)
			if err != nil {
				diags = append(diags, diag.FromErr(err)...)
			} else {
				efiDisk[mkEFIDiskFileFormat] = volume.FileFormat
			}
		}

		if vmConfig.EFIDisk.Type != nil {
			efiDisk[mkEFIDiskType] = *vmConfig.EFIDisk.Type
		} else {
			efiDisk[mkEFIDiskType] = dvEFIDiskType
		}

		if vmConfig.EFIDisk.PreEnrolledKeys != nil {
			efiDisk[mkEFIDiskPreEnrolledKeys] = *vmConfig.EFIDisk.PreEnrolledKeys
		} else {
			efiDisk[mkEFIDiskPreEnrolledKeys] = false
		}

		currentEfiDisk := d.Get(mkEFIDisk).([]interface{})

		if len(clone) > 0 {
			if len(currentEfiDisk) > 0 {
				err := d.Set(mkEFIDisk, []interface{}{efiDisk})
				diags = append(diags, diag.FromErr(err)...)
			}
		} else if len(currentEfiDisk) > 0 ||
			efiDisk[mkEFIDiskDatastoreID] != dvEFIDiskDatastoreID ||
			efiDisk[mkEFIDiskType] != dvEFIDiskType ||
			efiDisk[mkEFIDiskPreEnrolledKeys] != dvEFIDiskPreEnrolledKeys ||
			efiDisk[mkEFIDiskFileFormat] != dvEFIDiskFileFormat {
			err := d.Set(mkEFIDisk, []interface{}{efiDisk})
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	if vmConfig.TPMState != nil {
		tpmState := map[string]interface{}{}

		fileIDParts := strings.Split(vmConfig.TPMState.FileVolume, ":")

		tpmState[mkTPMStateDatastoreID] = fileIDParts[0]
		tpmState[mkTPMStateVersion] = dvTPMStateVersion

		currentTPMState := d.Get(mkTPMState).([]interface{})

		if len(clone) > 0 {
			if len(currentTPMState) > 0 {
				err := d.Set(mkTPMState, []interface{}{tpmState})
				diags = append(diags, diag.FromErr(err)...)
			}
		} else if len(currentTPMState) > 0 ||
			tpmState[mkTPMStateDatastoreID] != dvTPMStateDatastoreID ||
			tpmState[mkTPMStateVersion] != dvTPMStateVersion {
			err := d.Set(mkTPMState, []interface{}{tpmState})
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	currentPCIList := d.Get(mkHostPCI).([]interface{})
	pciMap := map[string]interface{}{}

	pciDevices := getPCIInfo(vmConfig, d)
	for pi, pp := range pciDevices {
		if (pp == nil) || (pp.DeviceIDs == nil && pp.Mapping == nil) {
			continue
		}

		pci := map[string]interface{}{}

		pci[mkHostPCIDevice] = pi
		if pp.DeviceIDs != nil {
			pci[mkHostPCIDeviceID] = strings.Join(*pp.DeviceIDs, ";")
		} else {
			pci[mkHostPCIDeviceID] = ""
		}

		if pp.MDev != nil {
			pci[mkHostPCIDeviceMDev] = *pp.MDev
		} else {
			pci[mkHostPCIDeviceMDev] = ""
		}

		if pp.PCIExpress != nil {
			pci[mkHostPCIDevicePCIE] = *pp.PCIExpress
		} else {
			pci[mkHostPCIDevicePCIE] = false
		}

		if pp.ROMBAR != nil {
			pci[mkHostPCIDeviceROMBAR] = *pp.ROMBAR
		} else {
			pci[mkHostPCIDeviceROMBAR] = true
		}

		if pp.ROMFile != nil {
			pci[mkHostPCIDeviceROMFile] = *pp.ROMFile
		} else {
			pci[mkHostPCIDeviceROMFile] = ""
		}

		if pp.XVGA != nil {
			pci[mkHostPCIDeviceXVGA] = *pp.XVGA
		} else {
			pci[mkHostPCIDeviceXVGA] = false
		}

		if pp.Mapping != nil {
			pci[mkHostPCIDeviceMapping] = *pp.Mapping
		} else {
			pci[mkHostPCIDeviceMapping] = ""
		}

		pciMap[pi] = pci
	}

	if len(clone) == 0 || len(currentPCIList) > 0 {
		orderedPCIList := orderedListFromMap(pciMap)
		err := d.Set(mkHostPCI, orderedPCIList)
		diags = append(diags, diag.FromErr(err)...)
	}

	currentUSBList := d.Get(mkHostUSB).([]interface{})
	usbMap := map[string]interface{}{}

	usbDevices := getUSBInfo(vmConfig, d)
	for pi, pp := range usbDevices {
		if (pp == nil) || (pp.HostDevice == nil && pp.Mapping == nil) {
			continue
		}

		usb := map[string]interface{}{}

		if pp.HostDevice != nil {
			usb[mkHostUSBDevice] = *pp.HostDevice
		} else {
			usb[mkHostUSBDevice] = ""
		}

		if pp.USB3 != nil {
			usb[mkHostUSBDeviceUSB3] = *pp.USB3
		} else {
			usb[mkHostUSBDeviceUSB3] = false
		}

		if pp.Mapping != nil {
			usb[mkHostUSBDeviceMapping] = *pp.Mapping
		} else {
			usb[mkHostUSBDeviceMapping] = ""
		}

		usbMap[pi] = usb
	}

	if len(clone) == 0 || len(currentUSBList) > 0 {
		// todo: reordering of devices by PVE may cause an issue here
		orderedUSBList := orderedListFromMap(usbMap)
		err := d.Set(mkHostUSB, orderedUSBList)
		diags = append(diags, diag.FromErr(err)...)
	}

	// Compare the initialization configuration to the one stored in the state.
	initialization := map[string]interface{}{}

	initializationInterface := findExistingCloudInitDrive(vmConfig, vmID, "")
	if initializationInterface != "" {
		initializationDevice := getIdeDevice(vmConfig, initializationInterface)
		fileVolumeParts := strings.Split(initializationDevice.FileVolume, ":")

		initialization[mkInitializationInterface] = initializationInterface
		initialization[mkInitializationDatastoreID] = fileVolumeParts[0]
	}

	if vmConfig.CloudInitDNSDomain != nil || vmConfig.CloudInitDNSServer != nil {
		initializationDNS := map[string]interface{}{}

		if vmConfig.CloudInitDNSDomain != nil {
			initializationDNS[mkInitializationDNSDomain] = *vmConfig.CloudInitDNSDomain
		} else {
			initializationDNS[mkInitializationDNSDomain] = ""
		}

		// check what we have in the plan
		currentInitializationDNSBlock := map[string]interface{}{}
		currentInitialization := d.Get(mkInitialization).([]interface{})

		if len(currentInitialization) > 0 {
			currentInitializationBlock := currentInitialization[0].(map[string]interface{})
			currentInitializationDNS := currentInitializationBlock[mkInitializationDNS].([]interface{})

			if len(currentInitializationDNS) > 0 {
				currentInitializationDNSBlock = currentInitializationDNS[0].(map[string]interface{})
			}
		}

		currentInitializationDNSServer, ok := currentInitializationDNSBlock[mkInitializationDNSServer]
		if vmConfig.CloudInitDNSServer != nil {
			if ok && currentInitializationDNSServer != "" {
				// the template is using deprecated attribute mkInitializationDNSServer
				initializationDNS[mkInitializationDNSServer] = *vmConfig.CloudInitDNSServer
			} else {
				dnsServer := strings.Split(*vmConfig.CloudInitDNSServer, " ")
				initializationDNS[mkInitializationDNSServers] = dnsServer
			}
		} else {
			initializationDNS[mkInitializationDNSServer] = ""
			initializationDNS[mkInitializationDNSServers] = []string{}
		}

		initialization[mkInitializationDNS] = []interface{}{
			initializationDNS,
		}
	}

	ipConfigLast := -1
	ipConfigObjects := []*vms.CustomCloudInitIPConfig{
		vmConfig.IPConfig0,
		vmConfig.IPConfig1,
		vmConfig.IPConfig2,
		vmConfig.IPConfig3,
		vmConfig.IPConfig4,
		vmConfig.IPConfig5,
		vmConfig.IPConfig6,
		vmConfig.IPConfig7,
		vmConfig.IPConfig7,
		vmConfig.IPConfig8,
		vmConfig.IPConfig9,
		vmConfig.IPConfig10,
		vmConfig.IPConfig11,
		vmConfig.IPConfig12,
		vmConfig.IPConfig13,
		vmConfig.IPConfig14,
		vmConfig.IPConfig15,
		vmConfig.IPConfig16,
		vmConfig.IPConfig17,
		vmConfig.IPConfig18,
		vmConfig.IPConfig19,
		vmConfig.IPConfig20,
		vmConfig.IPConfig21,
		vmConfig.IPConfig22,
		vmConfig.IPConfig23,
		vmConfig.IPConfig24,
		vmConfig.IPConfig25,
		vmConfig.IPConfig26,
		vmConfig.IPConfig27,
		vmConfig.IPConfig28,
		vmConfig.IPConfig29,
		vmConfig.IPConfig30,
		vmConfig.IPConfig31,
	}
	ipConfigList := make([]interface{}, len(ipConfigObjects))

	for ipConfigIndex, ipConfig := range ipConfigObjects {
		ipConfigItem := map[string]interface{}{}

		if ipConfig != nil {
			ipConfigLast = ipConfigIndex

			if ipConfig.GatewayIPv4 != nil || ipConfig.IPv4 != nil {
				ipv4 := map[string]interface{}{}

				if ipConfig.IPv4 != nil {
					ipv4[mkInitializationIPConfigIPv4Address] = *ipConfig.IPv4
				} else {
					ipv4[mkInitializationIPConfigIPv4Address] = ""
				}

				if ipConfig.GatewayIPv4 != nil {
					ipv4[mkInitializationIPConfigIPv4Gateway] = *ipConfig.GatewayIPv4
				} else {
					ipv4[mkInitializationIPConfigIPv4Gateway] = ""
				}

				ipConfigItem[mkInitializationIPConfigIPv4] = []interface{}{
					ipv4,
				}
			} else {
				ipConfigItem[mkInitializationIPConfigIPv4] = []interface{}{}
			}

			if ipConfig.GatewayIPv6 != nil || ipConfig.IPv6 != nil {
				ipv6 := map[string]interface{}{}

				if ipConfig.IPv6 != nil {
					ipv6[mkInitializationIPConfigIPv6Address] = *ipConfig.IPv6
				} else {
					ipv6[mkInitializationIPConfigIPv6Address] = ""
				}

				if ipConfig.GatewayIPv6 != nil {
					ipv6[mkInitializationIPConfigIPv6Gateway] = *ipConfig.GatewayIPv6
				} else {
					ipv6[mkInitializationIPConfigIPv6Gateway] = ""
				}

				ipConfigItem[mkInitializationIPConfigIPv6] = []interface{}{
					ipv6,
				}
			} else {
				ipConfigItem[mkInitializationIPConfigIPv6] = []interface{}{}
			}
		} else {
			ipConfigItem[mkInitializationIPConfigIPv4] = []interface{}{}
			ipConfigItem[mkInitializationIPConfigIPv6] = []interface{}{}
		}

		ipConfigList[ipConfigIndex] = ipConfigItem
	}

	if ipConfigLast >= 0 {
		initialization[mkInitializationIPConfig] = ipConfigList[:ipConfigLast+1]
	}

	//nolint:nestif
	if vmConfig.CloudInitPassword != nil || vmConfig.CloudInitSSHKeys != nil ||
		vmConfig.CloudInitUsername != nil {
		initializationUserAccount := map[string]interface{}{}

		if vmConfig.CloudInitSSHKeys != nil {
			initializationUserAccount[mkInitializationUserAccountKeys] = []string(
				*vmConfig.CloudInitSSHKeys,
			)
		} else {
			initializationUserAccount[mkInitializationUserAccountKeys] = []string{}
		}

		if vmConfig.CloudInitPassword != nil {
			initializationUserAccount[mkInitializationUserAccountPassword] = *vmConfig.CloudInitPassword
		} else {
			initializationUserAccount[mkInitializationUserAccountPassword] = ""
		}

		if vmConfig.CloudInitUsername != nil {
			initializationUserAccount[mkInitializationUserAccountUsername] = *vmConfig.CloudInitUsername
		} else {
			initializationUserAccount[mkInitializationUserAccountUsername] = ""
		}

		initialization[mkInitializationUserAccount] = []interface{}{
			initializationUserAccount,
		}
	}

	if vmConfig.CloudInitFiles != nil {
		if vmConfig.CloudInitFiles.UserVolume != nil {
			initialization[mkInitializationUserDataFileID] = *vmConfig.CloudInitFiles.UserVolume
		} else {
			initialization[mkInitializationUserDataFileID] = ""
		}

		if vmConfig.CloudInitFiles.VendorVolume != nil {
			initialization[mkInitializationVendorDataFileID] = *vmConfig.CloudInitFiles.VendorVolume
		} else {
			initialization[mkInitializationVendorDataFileID] = ""
		}

		if vmConfig.CloudInitFiles.NetworkVolume != nil {
			initialization[mkInitializationNetworkDataFileID] = *vmConfig.CloudInitFiles.NetworkVolume
		} else {
			initialization[mkInitializationNetworkDataFileID] = ""
		}

		if vmConfig.CloudInitFiles.MetaVolume != nil {
			initialization[mkInitializationMetaDataFileID] = *vmConfig.CloudInitFiles.MetaVolume
		} else {
			initialization[mkInitializationMetaDataFileID] = ""
		}
	} else if len(initialization) > 0 {
		initialization[mkInitializationUserDataFileID] = ""
		initialization[mkInitializationVendorDataFileID] = ""
		initialization[mkInitializationNetworkDataFileID] = ""
		initialization[mkInitializationMetaDataFileID] = ""
	}

	if vmConfig.CloudInitType != nil {
		initialization[mkInitializationType] = *vmConfig.CloudInitType
	} else if len(initialization) > 0 {
		initialization[mkInitializationType] = ""
	}

	currentInitialization := d.Get(mkInitialization).([]interface{})

	switch {
	case len(clone) > 0:
		if len(currentInitialization) > 0 {
			if len(initialization) > 0 {
				err := d.Set(
					mkInitialization,
					[]interface{}{initialization},
				)
				diags = append(diags, diag.FromErr(err)...)
			} else {
				err := d.Set(mkInitialization, []interface{}{})
				diags = append(diags, diag.FromErr(err)...)
			}
		}
	case len(initialization) > 0:
		err := d.Set(mkInitialization, []interface{}{initialization})
		diags = append(diags, diag.FromErr(err)...)
	default:
		err := d.Set(mkInitialization, []interface{}{})
		diags = append(diags, diag.FromErr(err)...)
	}

	// Compare the operating system configuration to the one stored in the state.
	kvmArguments := map[string]interface{}{}

	if vmConfig.KVMArguments != nil {
		kvmArguments[mkKVMArguments] = *vmConfig.KVMArguments
	} else {
		kvmArguments[mkKVMArguments] = ""
	}

	// Compare the memory configuration to the one stored in the state.
	memory := map[string]interface{}{}

	if vmConfig.DedicatedMemory != nil {
		memory[mkMemoryDedicated] = int(*vmConfig.DedicatedMemory)
	} else {
		memory[mkMemoryDedicated] = 0
	}

	if vmConfig.FloatingMemory != nil {
		memory[mkMemoryFloating] = int(*vmConfig.FloatingMemory)
	} else {
		memory[mkMemoryFloating] = 0
	}

	if vmConfig.SharedMemory != nil {
		memory[mkMemoryShared] = vmConfig.SharedMemory.Size
	} else {
		memory[mkMemoryShared] = 0
	}

	currentMemory := d.Get(mkMemory).([]interface{})

	if len(clone) > 0 {
		if len(currentMemory) > 0 {
			err := d.Set(mkMemory, []interface{}{memory})
			diags = append(diags, diag.FromErr(err)...)
		}
	} else if len(currentMemory) > 0 ||
		memory[mkMemoryDedicated] != dvMemoryDedicated ||
		memory[mkMemoryFloating] != dvMemoryFloating ||
		memory[mkMemoryShared] != dvMemoryShared {
		err := d.Set(mkMemory, []interface{}{memory})
		diags = append(diags, diag.FromErr(err)...)
	}

	// Compare the network devices to those stored in the state.
	currentNetworkDeviceList := d.Get(mkNetworkDevice).([]interface{})

	macAddresses := make([]interface{}, maxResourceVirtualEnvironmentVMNetworkDevices)
	networkDeviceLast := -1
	networkDeviceList := make([]interface{}, maxResourceVirtualEnvironmentVMNetworkDevices)
	networkDeviceObjects := []*vms.CustomNetworkDevice{
		vmConfig.NetworkDevice0,
		vmConfig.NetworkDevice1,
		vmConfig.NetworkDevice2,
		vmConfig.NetworkDevice3,
		vmConfig.NetworkDevice4,
		vmConfig.NetworkDevice5,
		vmConfig.NetworkDevice6,
		vmConfig.NetworkDevice7,
		vmConfig.NetworkDevice8,
		vmConfig.NetworkDevice9,
		vmConfig.NetworkDevice10,
		vmConfig.NetworkDevice11,
		vmConfig.NetworkDevice12,
		vmConfig.NetworkDevice13,
		vmConfig.NetworkDevice14,
		vmConfig.NetworkDevice15,
		vmConfig.NetworkDevice16,
		vmConfig.NetworkDevice17,
		vmConfig.NetworkDevice18,
		vmConfig.NetworkDevice19,
		vmConfig.NetworkDevice20,
		vmConfig.NetworkDevice21,
		vmConfig.NetworkDevice22,
		vmConfig.NetworkDevice23,
		vmConfig.NetworkDevice24,
		vmConfig.NetworkDevice25,
		vmConfig.NetworkDevice26,
		vmConfig.NetworkDevice27,
		vmConfig.NetworkDevice28,
		vmConfig.NetworkDevice29,
		vmConfig.NetworkDevice30,
		vmConfig.NetworkDevice31,
	}

	for ni, nd := range networkDeviceObjects {
		networkDevice := map[string]interface{}{}

		if nd != nil {
			networkDeviceLast = ni

			if nd.Bridge != nil {
				networkDevice[mkNetworkDeviceBridge] = *nd.Bridge
			} else {
				networkDevice[mkNetworkDeviceBridge] = ""
			}

			networkDevice[mkNetworkDeviceEnabled] = nd.Enabled

			if nd.Firewall != nil {
				networkDevice[mkNetworkDeviceFirewall] = *nd.Firewall
			} else {
				networkDevice[mkNetworkDeviceFirewall] = false
			}

			if nd.MACAddress != nil {
				macAddresses[ni] = *nd.MACAddress
			} else {
				macAddresses[ni] = ""
			}

			networkDevice[mkNetworkDeviceMACAddress] = macAddresses[ni]
			networkDevice[mkNetworkDeviceModel] = nd.Model

			if nd.Queues != nil {
				networkDevice[mkNetworkDeviceQueues] = *nd.Queues
			} else {
				networkDevice[mkNetworkDeviceQueues] = 0
			}

			if nd.RateLimit != nil {
				networkDevice[mkNetworkDeviceRateLimit] = *nd.RateLimit
			} else {
				networkDevice[mkNetworkDeviceRateLimit] = 0
			}

			if nd.Tag != nil {
				networkDevice[mkNetworkDeviceVLANID] = nd.Tag
			} else {
				networkDevice[mkNetworkDeviceVLANID] = 0
			}

			if nd.MTU != nil {
				networkDevice[mkNetworkDeviceMTU] = nd.MTU
			} else {
				networkDevice[mkNetworkDeviceMTU] = 0
			}
		} else {
			macAddresses[ni] = ""
			networkDevice[mkNetworkDeviceEnabled] = false
		}

		networkDeviceList[ni] = networkDevice
	}

	if len(currentNetworkDeviceList) == 0 {
		err := d.Set(mkMACAddresses, []interface{}{})
		diags = append(diags, diag.FromErr(err)...)
		err = d.Set(mkNetworkDevice, []interface{}{})
		diags = append(diags, diag.FromErr(err)...)
	} else {
		err := d.Set(mkMACAddresses, macAddresses[0:len(currentNetworkDeviceList)])
		diags = append(diags, diag.FromErr(err)...)

		if len(clone) > 0 {
			err = d.Set(mkNetworkDevice, networkDeviceList[:networkDeviceLast+1])
			diags = append(diags, diag.FromErr(err)...)
		} else if len(currentNetworkDeviceList) > 0 || networkDeviceLast > -1 {
			err := d.Set(mkNetworkDevice, networkDeviceList[:networkDeviceLast+1])
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	// Compare the operating system configuration to the one stored in the state.
	operatingSystem := map[string]interface{}{}

	if vmConfig.OSType != nil {
		operatingSystem[mkOperatingSystemType] = *vmConfig.OSType
	} else {
		operatingSystem[mkOperatingSystemType] = ""
	}

	currentOperatingSystem := d.Get(mkOperatingSystem).([]interface{})

	switch {
	case len(clone) > 0:
		if len(currentOperatingSystem) > 0 {
			err := d.Set(
				mkOperatingSystem,
				[]interface{}{operatingSystem},
			)
			diags = append(diags, diag.FromErr(err)...)
		}
	case len(currentOperatingSystem) > 0 ||
		operatingSystem[mkOperatingSystemType] != dvOperatingSystemType:
		err := d.Set(mkOperatingSystem, []interface{}{operatingSystem})
		diags = append(diags, diag.FromErr(err)...)
	default:
		err := d.Set(mkOperatingSystem, []interface{}{})
		diags = append(diags, diag.FromErr(err)...)
	}

	// Compare the pool ID to the value stored in the state.
	currentPoolID := d.Get(mkPoolID).(string)

	if len(clone) == 0 || currentPoolID != dvPoolID {
		if vmConfig.PoolID != nil {
			err := d.Set(mkPoolID, *vmConfig.PoolID)
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	// Compare the serial devices to those stored in the state.
	serialDevices := make([]interface{}, 4)
	serialDevicesArray := []*string{
		vmConfig.SerialDevice0,
		vmConfig.SerialDevice1,
		vmConfig.SerialDevice2,
		vmConfig.SerialDevice3,
	}
	serialDevicesCount := 0

	for sdi, sd := range serialDevicesArray {
		m := map[string]interface{}{}

		if sd != nil {
			m[mkSerialDeviceDevice] = *sd
			serialDevicesCount = sdi + 1
		} else {
			m[mkSerialDeviceDevice] = ""
		}

		serialDevices[sdi] = m
	}

	currentSerialDevice := d.Get(mkSerialDevice).([]interface{})

	if len(clone) == 0 || len(currentSerialDevice) > 0 {
		err := d.Set(mkSerialDevice, serialDevices[:serialDevicesCount])
		diags = append(diags, diag.FromErr(err)...)
	}

	// Compare the SMBIOS to the one stored in the state.
	var smbios map[string]interface{}

	//nolint:nestif
	if vmConfig.SMBIOS != nil {
		smbios = map[string]interface{}{}

		if vmConfig.SMBIOS.Family != nil {
			b, err := base64.StdEncoding.DecodeString(*vmConfig.SMBIOS.Family)
			diags = append(diags, diag.FromErr(err)...)
			smbios[mkSMBIOSFamily] = string(b)
		} else {
			smbios[mkSMBIOSFamily] = dvSMBIOSFamily
		}

		if vmConfig.SMBIOS.Manufacturer != nil {
			b, err := base64.StdEncoding.DecodeString(*vmConfig.SMBIOS.Manufacturer)
			diags = append(diags, diag.FromErr(err)...)
			smbios[mkSMBIOSManufacturer] = string(b)
		} else {
			smbios[mkSMBIOSManufacturer] = dvSMBIOSManufacturer
		}

		if vmConfig.SMBIOS.Product != nil {
			b, err := base64.StdEncoding.DecodeString(*vmConfig.SMBIOS.Product)
			diags = append(diags, diag.FromErr(err)...)
			smbios[mkSMBIOSProduct] = string(b)
		} else {
			smbios[mkSMBIOSProduct] = dvSMBIOSProduct
		}

		if vmConfig.SMBIOS.Serial != nil {
			b, err := base64.StdEncoding.DecodeString(*vmConfig.SMBIOS.Serial)
			diags = append(diags, diag.FromErr(err)...)
			smbios[mkSMBIOSSerial] = string(b)
		} else {
			smbios[mkSMBIOSSerial] = dvSMBIOSSerial
		}

		if vmConfig.SMBIOS.SKU != nil {
			b, err := base64.StdEncoding.DecodeString(*vmConfig.SMBIOS.SKU)
			diags = append(diags, diag.FromErr(err)...)
			smbios[mkSMBIOSSKU] = string(b)
		} else {
			smbios[mkSMBIOSSKU] = dvSMBIOSSKU
		}

		if vmConfig.SMBIOS.Version != nil {
			b, err := base64.StdEncoding.DecodeString(*vmConfig.SMBIOS.Version)
			diags = append(diags, diag.FromErr(err)...)
			smbios[mkSMBIOSVersion] = string(b)
		} else {
			smbios[mkSMBIOSVersion] = dvSMBIOSVersion
		}

		if vmConfig.SMBIOS.UUID != nil {
			smbios[mkSMBIOSUUID] = *vmConfig.SMBIOS.UUID
		} else {
			smbios[mkSMBIOSUUID] = nil
		}
	}

	currentSMBIOS := d.Get(mkSMBIOS).([]interface{})

	//nolint:gocritic
	if len(clone) > 0 {
		if len(currentSMBIOS) > 0 {
			err := d.Set(mkSMBIOS, currentSMBIOS)
			diags = append(diags, diag.FromErr(err)...)
		}
	} else if len(smbios) == 0 {
		err := d.Set(mkSMBIOS, []interface{}{})
		diags = append(diags, diag.FromErr(err)...)
	} else if len(currentSMBIOS) > 0 ||
		smbios[mkSMBIOSFamily] != dvSMBIOSFamily ||
		smbios[mkSMBIOSManufacturer] != dvSMBIOSManufacturer ||
		smbios[mkSMBIOSProduct] != dvSMBIOSProduct ||
		smbios[mkSMBIOSSerial] != dvSMBIOSSerial ||
		smbios[mkSMBIOSSKU] != dvSMBIOSSKU ||
		smbios[mkSMBIOSVersion] != dvSMBIOSVersion {
		err := d.Set(mkSMBIOS, []interface{}{smbios})
		diags = append(diags, diag.FromErr(err)...)
	}

	// Compare the startup order to the one stored in the state.
	var startup map[string]interface{}

	//nolint:nestif
	if vmConfig.StartupOrder != nil {
		startup = map[string]interface{}{}

		if vmConfig.StartupOrder.Order != nil {
			startup[mkStartupOrder] = *vmConfig.StartupOrder.Order
		} else {
			startup[mkStartupOrder] = dvStartupOrder
		}

		if vmConfig.StartupOrder.Up != nil {
			startup[mkStartupUpDelay] = *vmConfig.StartupOrder.Up
		} else {
			startup[mkStartupUpDelay] = dvStartupUpDelay
		}

		if vmConfig.StartupOrder.Down != nil {
			startup[mkStartupDownDelay] = *vmConfig.StartupOrder.Down
		} else {
			startup[mkStartupDownDelay] = dvStartupDownDelay
		}
	}

	currentStartup := d.Get(mkStartup).([]interface{})

	//nolint:gocritic
	if len(clone) > 0 {
		if len(currentStartup) > 0 {
			err := d.Set(mkStartup, []interface{}{startup})
			diags = append(diags, diag.FromErr(err)...)
		}
	} else if len(startup) == 0 {
		err := d.Set(mkStartup, []interface{}{})
		diags = append(diags, diag.FromErr(err)...)
	} else if len(currentStartup) > 0 ||
		startup[mkStartupOrder] != mkStartupOrder ||
		startup[mkStartupUpDelay] != dvStartupUpDelay ||
		startup[mkStartupDownDelay] != dvStartupDownDelay {
		err := d.Set(mkStartup, []interface{}{startup})
		diags = append(diags, diag.FromErr(err)...)
	}

	// Compare the VGA configuration to the one stored in the state.
	vga := map[string]interface{}{}

	if vmConfig.VGADevice != nil {
		vgaEnabled := true

		if vmConfig.VGADevice.Type != nil {
			vgaEnabled = *vmConfig.VGADevice.Type != "none"
		}

		vga[mkVGAEnabled] = vgaEnabled

		if vmConfig.VGADevice.Memory != nil {
			vga[mkVGAMemory] = *vmConfig.VGADevice.Memory
		} else {
			vga[mkVGAMemory] = 0
		}

		if vgaEnabled {
			if vmConfig.VGADevice.Type != nil {
				vga[mkVGAType] = *vmConfig.VGADevice.Type
			} else {
				vga[mkVGAType] = ""
			}
		}
	} else {
		vga[mkVGAEnabled] = true
		vga[mkVGAMemory] = 0
		vga[mkVGAType] = ""
	}

	currentVGA := d.Get(mkVGA).([]interface{})

	switch {
	case len(clone) > 0:
		if len(currentVGA) > 0 {
			err := d.Set(mkVGA, []interface{}{vga})
			diags = append(diags, diag.FromErr(err)...)
		}
	case len(currentVGA) > 0 ||
		vga[mkVGAEnabled] != dvVGAEnabled ||
		vga[mkVGAMemory] != dvVGAMemory ||
		vga[mkVGAType] != dvVGAType:
		err := d.Set(mkVGA, []interface{}{vga})
		diags = append(diags, diag.FromErr(err)...)
	default:
		err := d.Set(mkVGA, []interface{}{})
		diags = append(diags, diag.FromErr(err)...)
	}

	// Compare SCSI hardware type
	scsiHardware := d.Get(mkSCSIHardware).(string)

	if len(clone) == 0 || scsiHardware != dvSCSIHardware {
		if vmConfig.SCSIHardware != nil {
			err := d.Set(mkSCSIHardware, *vmConfig.SCSIHardware)
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	diags = append(
		diags,
		vmReadNetworkValues(ctx, d, m, vmID, vmConfig)...)

	// during import these core attributes might not be set, so set them explicitly here
	d.SetId(strconv.Itoa(vmID))
	e := d.Set(mkVMID, vmID)
	diags = append(diags, diag.FromErr(e)...)
	e = d.Set(mkNodeName, nodeName)
	diags = append(diags, diag.FromErr(e)...)

	return diags
}

func vmReadNetworkValues(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
	vmID int,
	vmConfig *vms.GetResponseData,
) diag.Diagnostics {
	var diags diag.Diagnostics

	config := m.(proxmoxtf.ProviderConfiguration)

	api, e := config.GetClient()
	if e != nil {
		return diag.FromErr(e)
	}

	nodeName := d.Get(mkNodeName).(string)

	vmAPI := api.Node(nodeName).VM(vmID)

	started := d.Get(mkStarted).(bool)

	var ipv4Addresses []interface{}

	var ipv6Addresses []interface{}

	var networkInterfaceNames []interface{}

	if started {
		if vmConfig.Agent != nil && vmConfig.Agent.Enabled != nil && *vmConfig.Agent.Enabled {
			resource := VM()

			agentBlock, err := structure.GetSchemaBlock(
				resource,
				d,
				[]string{mkAgent},
				0,
				true,
			)
			if err != nil {
				return diag.FromErr(err)
			}

			agentTimeout, err := time.ParseDuration(
				agentBlock[mkAgentTimeout].(string),
			)
			if err != nil {
				return diag.FromErr(err)
			}

			var macAddresses []interface{}

			networkInterfaces, err := vmAPI.WaitForNetworkInterfacesFromVMAgent(ctx, int(agentTimeout.Seconds()), 5, true)
			if err == nil && networkInterfaces.Result != nil {
				ipv4Addresses = make([]interface{}, len(*networkInterfaces.Result))
				ipv6Addresses = make([]interface{}, len(*networkInterfaces.Result))
				macAddresses = make([]interface{}, len(*networkInterfaces.Result))
				networkInterfaceNames = make([]interface{}, len(*networkInterfaces.Result))

				for ri, rv := range *networkInterfaces.Result {
					var rvIPv4Addresses []interface{}

					var rvIPv6Addresses []interface{}

					if rv.IPAddresses != nil {
						for _, ip := range *rv.IPAddresses {
							switch ip.Type {
							case "ipv4":
								rvIPv4Addresses = append(rvIPv4Addresses, ip.Address)
							case "ipv6":
								rvIPv6Addresses = append(rvIPv6Addresses, ip.Address)
							}
						}
					}

					ipv4Addresses[ri] = rvIPv4Addresses
					ipv6Addresses[ri] = rvIPv6Addresses
					macAddresses[ri] = strings.ToUpper(rv.MACAddress)
					networkInterfaceNames[ri] = rv.Name
				}
			}

			err = d.Set(mkMACAddresses, macAddresses)
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	e = d.Set(mkIPv4Addresses, ipv4Addresses)
	diags = append(diags, diag.FromErr(e)...)
	e = d.Set(mkIPv6Addresses, ipv6Addresses)
	diags = append(diags, diag.FromErr(e)...)
	e = d.Set(mkNetworkInterfaceNames, networkInterfaceNames)
	diags = append(diags, diag.FromErr(e)...)

	return diags
}

func vmReadPrimitiveValues(
	d *schema.ResourceData,
	vmConfig *vms.GetResponseData,
	vmStatus *vms.GetStatusResponseData,
) diag.Diagnostics {
	var diags diag.Diagnostics

	var err error

	clone := d.Get(mkClone).([]interface{})
	currentACPI := d.Get(mkACPI).(bool)

	//nolint:gosimple
	if len(clone) == 0 || currentACPI != dvACPI {
		if vmConfig.ACPI != nil {
			err = d.Set(mkACPI, bool(*vmConfig.ACPI))
		} else {
			// Default value of "acpi" is "1" according to the API documentation.
			err = d.Set(mkACPI, true)
		}

		diags = append(diags, diag.FromErr(err)...)
	}

	currentkvmArguments := d.Get(mkKVMArguments).(string)

	if len(clone) == 0 || currentkvmArguments != dvKVMArguments {
		// PVE API returns "args" as " " if it is set to empty.
		if vmConfig.KVMArguments != nil && len(strings.TrimSpace(*vmConfig.KVMArguments)) > 0 {
			err = d.Set(mkKVMArguments, *vmConfig.KVMArguments)
		} else {
			// Default value of "args" is "" according to the API documentation.
			err = d.Set(mkKVMArguments, "")
		}

		diags = append(diags, diag.FromErr(err)...)
	}

	currentBIOS := d.Get(mkBIOS).(string)

	if len(clone) == 0 || currentBIOS != dvBIOS {
		if vmConfig.BIOS != nil {
			err = d.Set(mkBIOS, *vmConfig.BIOS)
		} else {
			// Default value of "bios" is "seabios" according to the API documentation.
			err = d.Set(mkBIOS, "seabios")
		}

		diags = append(diags, diag.FromErr(err)...)
	}

	currentDescription := d.Get(mkDescription).(string)

	if len(clone) == 0 || currentDescription != dvDescription {
		if vmConfig.Description != nil {
			err = d.Set(mkDescription, *vmConfig.Description)
		} else {
			// Default value of "description" is "" according to the API documentation.
			err = d.Set(mkDescription, "")
		}

		diags = append(diags, diag.FromErr(err)...)
	}

	currentTags := d.Get(mkTags).([]interface{})

	if len(clone) == 0 || len(currentTags) > 0 {
		var tags []string

		if vmConfig.Tags != nil {
			for _, tag := range strings.Split(*vmConfig.Tags, ";") {
				t := strings.TrimSpace(tag)
				if len(t) > 0 {
					tags = append(tags, t)
				}
			}

			sort.Strings(tags)
		}

		err = d.Set(mkTags, tags)
		diags = append(diags, diag.FromErr(err)...)
	}

	currentKeyboardLayout := d.Get(mkKeyboardLayout).(string)

	if len(clone) == 0 || currentKeyboardLayout != dvKeyboardLayout {
		if vmConfig.KeyboardLayout != nil {
			err = d.Set(mkKeyboardLayout, *vmConfig.KeyboardLayout)
		} else {
			// Default value of "keyboard" is "" according to the API documentation.
			err = d.Set(mkKeyboardLayout, "")
		}

		diags = append(diags, diag.FromErr(err)...)
	}

	currentMachine := d.Get(mkMachine).(string)

	if len(clone) == 0 || currentMachine != dvMachineType {
		if vmConfig.Machine != nil {
			err = d.Set(mkMachine, *vmConfig.Machine)
		} else {
			err = d.Set(mkMachine, "")
		}

		diags = append(diags, diag.FromErr(err)...)
	}

	currentName := d.Get(mkName).(string)

	if len(clone) == 0 || currentName != dvName {
		if vmConfig.Name != nil {
			err = d.Set(mkName, *vmConfig.Name)
		} else {
			// Default value of "name" is "" according to the API documentation.
			err = d.Set(mkName, "")
		}

		diags = append(diags, diag.FromErr(err)...)
	}

	if !d.Get(mkTemplate).(bool) {
		err = d.Set(mkStarted, vmStatus.Status == "running")
		diags = append(diags, diag.FromErr(err)...)
	}

	currentTabletDevice := d.Get(mkTabletDevice).(bool)

	//nolint:gosimple
	if len(clone) == 0 || currentTabletDevice != dvTabletDevice {
		if vmConfig.TabletDeviceEnabled != nil {
			err = d.Set(
				mkTabletDevice,
				bool(*vmConfig.TabletDeviceEnabled),
			)
		} else {
			// Default value of "tablet" is "1" according to the API documentation.
			err = d.Set(mkTabletDevice, true)
		}

		diags = append(diags, diag.FromErr(err)...)
	}

	currentTemplate := d.Get(mkTemplate).(bool)

	//nolint:gosimple
	if len(clone) == 0 || currentTemplate != dvTemplate {
		if vmConfig.Template != nil {
			err = d.Set(mkTemplate, bool(*vmConfig.Template))
		} else {
			// Default value of "template" is "0" according to the API documentation.
			err = d.Set(mkTemplate, false)
		}

		diags = append(diags, diag.FromErr(err)...)
	}

	return diags
}

// vmUpdatePool moves the VM to the pool it is supposed to be in if the pool ID changed.
func vmUpdatePool(
	ctx context.Context,
	d *schema.ResourceData,
	api *pools.Client,
	vmID int,
) error {
	oldPoolValue, newPoolValue := d.GetChange(mkPoolID)
	if cmp.Equal(newPoolValue, oldPoolValue) {
		return nil
	}

	oldPool := oldPoolValue.(string)
	newPool := newPoolValue.(string)
	vmList := (types.CustomCommaSeparatedList)([]string{strconv.Itoa(vmID)})

	tflog.Debug(ctx, fmt.Sprintf("Moving VM %d from pool '%s' to pool '%s'", vmID, oldPool, newPool))

	if oldPool != "" {
		trueValue := types.CustomBool(true)
		poolUpdate := &pools.PoolUpdateRequestBody{
			VMs:    &vmList,
			Delete: &trueValue,
		}

		err := api.UpdatePool(ctx, oldPool, poolUpdate)
		if err != nil {
			return fmt.Errorf("while removing VM %d from pool %s: %w", vmID, oldPool, err)
		}
	}

	if newPool != "" {
		poolUpdate := &pools.PoolUpdateRequestBody{VMs: &vmList}

		err := api.UpdatePool(ctx, newPool, poolUpdate)
		if err != nil {
			return fmt.Errorf("while adding VM %d to pool %s: %w", vmID, newPool, err)
		}
	}

	return nil
}

func vmUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(proxmoxtf.ProviderConfiguration)

	api, e := config.GetClient()
	if e != nil {
		return diag.FromErr(e)
	}

	nodeName := d.Get(mkNodeName).(string)
	rebootRequired := false

	vmID, e := strconv.Atoi(d.Id())
	if e != nil {
		return diag.FromErr(e)
	}

	e = vmUpdatePool(ctx, d, api.Pool(), vmID)
	if e != nil {
		return diag.FromErr(e)
	}

	// If the node name has changed we need to migrate the VM to the new node before we do anything else.
	if d.HasChange(mkNodeName) {
		oldNodeNameValue, _ := d.GetChange(mkNodeName)
		oldNodeName := oldNodeNameValue.(string)
		vmAPI := api.Node(oldNodeName).VM(vmID)

		migrateTimeout := d.Get(mkTimeoutMigrate).(int)
		trueValue := types.CustomBool(true)
		migrateBody := &vms.MigrateRequestBody{
			TargetNode:      nodeName,
			WithLocalDisks:  &trueValue,
			OnlineMigration: &trueValue,
		}

		err := vmAPI.MigrateVM(ctx, migrateBody, migrateTimeout)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	vmAPI := api.Node(nodeName).VM(vmID)

	updateBody := &vms.UpdateRequestBody{
		IDEDevices: vms.CustomStorageDevices{
			"ide0": &vms.CustomStorageDevice{
				Enabled: false,
			},
			"ide1": &vms.CustomStorageDevice{
				Enabled: false,
			},
			"ide2": &vms.CustomStorageDevice{
				Enabled: false,
			},
			"ide3": &vms.CustomStorageDevice{
				Enabled: false,
			},
		},
	}

	var del []string

	resource := VM()

	// Retrieve the entire configuration as we need to process certain values.
	vmConfig, e := vmAPI.GetVM(ctx)
	if e != nil {
		return diag.FromErr(e)
	}

	// Prepare the new primitive configuration values.
	if d.HasChange(mkACPI) {
		acpi := types.CustomBool(d.Get(mkACPI).(bool))
		updateBody.ACPI = &acpi
		rebootRequired = true
	}

	if d.HasChange(mkKVMArguments) {
		kvmArguments := d.Get(mkKVMArguments).(string)
		updateBody.KVMArguments = &kvmArguments
		rebootRequired = true
	}

	if d.HasChange(mkBIOS) {
		bios := d.Get(mkBIOS).(string)
		updateBody.BIOS = &bios
		rebootRequired = true
	}

	if d.HasChange(mkDescription) {
		description := d.Get(mkDescription).(string)
		updateBody.Description = &description
	}

	if d.HasChange(mkOnBoot) {
		startOnBoot := types.CustomBool(d.Get(mkOnBoot).(bool))
		updateBody.StartOnBoot = &startOnBoot
	}

	if d.HasChange(mkTags) {
		tagString := vmGetTagsString(d)
		updateBody.Tags = &tagString
	}

	if d.HasChange(mkKeyboardLayout) {
		keyboardLayout := d.Get(mkKeyboardLayout).(string)
		updateBody.KeyboardLayout = &keyboardLayout
		rebootRequired = true
	}

	if d.HasChange(mkMachine) {
		machine := d.Get(mkMachine).(string)
		updateBody.Machine = &machine
		rebootRequired = true
	}

	name := d.Get(mkName).(string)

	if name == "" {
		del = append(del, "name")
	} else {
		updateBody.Name = &name
	}

	if d.HasChange(mkTabletDevice) {
		tabletDevice := types.CustomBool(d.Get(mkTabletDevice).(bool))
		updateBody.TabletDeviceEnabled = &tabletDevice
		rebootRequired = true
	}

	template := types.CustomBool(d.Get(mkTemplate).(bool))

	if d.HasChange(mkTemplate) {
		updateBody.Template = &template
		rebootRequired = true
	}

	// Prepare the new agent configuration.
	if d.HasChange(mkAgent) {
		agentBlock, err := structure.GetSchemaBlock(
			resource,
			d,
			[]string{mkAgent},
			0,
			true,
		)
		if err != nil {
			return diag.FromErr(err)
		}

		agentEnabled := types.CustomBool(
			agentBlock[mkAgentEnabled].(bool),
		)
		agentTrim := types.CustomBool(agentBlock[mkAgentTrim].(bool))
		agentType := agentBlock[mkAgentType].(string)

		updateBody.Agent = &vms.CustomAgent{
			Enabled:         &agentEnabled,
			TrimClonedDisks: &agentTrim,
			Type:            &agentType,
		}

		rebootRequired = true
	}

	// Prepare the new audio devices.
	if d.HasChange(mkAudioDevice) {
		updateBody.AudioDevices = vmGetAudioDeviceList(d)

		for i := 0; i < len(updateBody.AudioDevices); i++ {
			if !updateBody.AudioDevices[i].Enabled {
				del = append(del, fmt.Sprintf("audio%d", i))
			}
		}

		for i := len(updateBody.AudioDevices); i < maxResourceVirtualEnvironmentVMAudioDevices; i++ {
			del = append(del, fmt.Sprintf("audio%d", i))
		}

		rebootRequired = true
	}

	// Prepare the new boot configuration.
	if d.HasChange(mkBootOrder) {
		bootOrder := d.Get(mkBootOrder).([]interface{})
		bootOrderConverted := make([]string, len(bootOrder))

		for i, device := range bootOrder {
			bootOrderConverted[i] = device.(string)
		}

		updateBody.Boot = &vms.CustomBoot{
			Order: &bootOrderConverted,
		}
		rebootRequired = true
	}

	// Prepare the new CD-ROM configuration.
	if d.HasChange(mkCDROM) {
		cdromBlock, err := structure.GetSchemaBlock(
			resource,
			d,
			[]string{mkCDROM},
			0,
			true,
		)
		if err != nil {
			return diag.FromErr(err)
		}

		cdromEnabled := cdromBlock[mkCDROMEnabled].(bool)
		cdromFileID := cdromBlock[mkCDROMFileID].(string)
		cdromInterface := cdromBlock[mkCDROMInterface].(string)

		old, _ := d.GetChange(mkCDROM)

		if len(old.([]interface{})) > 0 {
			oldList := old.([]interface{})[0]
			oldBlock := oldList.(map[string]interface{})

			// If the interface is not set, use the default, for backward compatibility.
			oldInterface, ok := oldBlock[mkCDROMInterface].(string)
			if !ok || oldInterface == "" {
				oldInterface = dvCDROMInterface
			}

			if oldInterface != cdromInterface {
				del = append(del, oldInterface)
			}
		}

		if !cdromEnabled && cdromFileID == "" {
			del = append(del, cdromInterface)
		}

		if cdromFileID == "" {
			cdromFileID = "cdrom"
		}

		cdromMedia := "cdrom"

		updateBody.IDEDevices[cdromInterface] = &vms.CustomStorageDevice{
			Enabled:    cdromEnabled,
			FileVolume: cdromFileID,
			Media:      &cdromMedia,
		}
	}

	// Prepare the new CPU configuration.
	if d.HasChange(mkCPU) {
		cpuBlock, err := structure.GetSchemaBlock(
			resource,
			d,
			[]string{mkCPU},
			0,
			true,
		)
		if err != nil {
			return diag.FromErr(err)
		}

		cpuArchitecture := cpuBlock[mkCPUArchitecture].(string)
		cpuCores := cpuBlock[mkCPUCores].(int)
		cpuFlags := cpuBlock[mkCPUFlags].([]interface{})
		cpuHotplugged := cpuBlock[mkCPUHotplugged].(int)
		cpuLimit := cpuBlock[mkCPULimit].(int)
		cpuNUMA := types.CustomBool(cpuBlock[mkCPUNUMA].(bool))
		cpuSockets := cpuBlock[mkCPUSockets].(int)
		cpuType := cpuBlock[mkCPUType].(string)
		cpuUnits := cpuBlock[mkCPUUnits].(int)

		// Only the root account is allowed to change the CPU architecture, which makes this check necessary.
		if api.API().IsRootTicket() ||
			cpuArchitecture != dvCPUArchitecture {
			updateBody.CPUArchitecture = &cpuArchitecture
		}

		updateBody.CPUCores = &cpuCores
		updateBody.CPUSockets = &cpuSockets
		updateBody.CPUUnits = &cpuUnits
		updateBody.NUMAEnabled = &cpuNUMA

		if cpuHotplugged > 0 {
			updateBody.VirtualCPUCount = &cpuHotplugged
		} else {
			del = append(del, "vcpus")
		}

		if cpuLimit > 0 {
			updateBody.CPULimit = &cpuLimit
		} else {
			del = append(del, "cpulimit")
		}

		cpuFlagsConverted := make([]string, len(cpuFlags))

		for fi, flag := range cpuFlags {
			cpuFlagsConverted[fi] = flag.(string)
		}

		updateBody.CPUEmulation = &vms.CustomCPUEmulation{
			Flags: &cpuFlagsConverted,
			Type:  cpuType,
		}

		rebootRequired = true
	}

	err := updateDisk(d, updateBody)
	if err != nil {
		return diag.FromErr(err)
	}

	// Prepare the new efi disk configuration.
	if d.HasChange(mkEFIDisk) {
		efiDisk := vmGetEfiDisk(d, nil)

		updateBody.EFIDisk = efiDisk

		rebootRequired = true
	}

	// Prepare the new tpm state configuration.
	if d.HasChange(mkTPMState) {
		tpmState := vmGetTPMState(d, nil)

		updateBody.TPMState = tpmState

		rebootRequired = true
	}

	// Prepare the new cloud-init configuration.
	stoppedBeforeUpdate := false

	if d.HasChange(mkInitialization) {
		initializationConfig := vmGetCloudInitConfig(d)

		updateBody.CloudInitConfig = initializationConfig

		if updateBody.CloudInitConfig != nil {
			var fileVolume string

			initialization := d.Get(mkInitialization).([]interface{})
			initializationBlock := initialization[0].(map[string]interface{})
			initializationDatastoreID := initializationBlock[mkInitializationDatastoreID].(string)
			initializationInterface := initializationBlock[mkInitializationInterface].(string)
			cdromMedia := "cdrom"

			existingInterface := findExistingCloudInitDrive(vmConfig, vmID, "")
			if initializationInterface == "" && existingInterface == "" {
				initializationInterface = "ide2"
			} else if initializationInterface == "" {
				initializationInterface = existingInterface
			}

			mustMove := existingInterface != "" && initializationInterface != existingInterface
			if mustMove {
				tflog.Debug(ctx, fmt.Sprintf("CloudInit must be moved from %s to %s", existingInterface, initializationInterface))
			}

			oldInit, _ := d.GetChange(mkInitialization)
			oldInitBlock := oldInit.([]interface{})[0].(map[string]interface{})
			prevDatastoreID := oldInitBlock[mkInitializationDatastoreID].(string)

			mustChangeDatastore := prevDatastoreID != initializationDatastoreID
			if mustChangeDatastore {
				tflog.Debug(ctx, fmt.Sprintf("CloudInit must be moved from datastore %s to datastore %s",
					prevDatastoreID, initializationDatastoreID))
			}

			if mustMove || mustChangeDatastore || existingInterface == "" {
				// CloudInit must be moved, either from a device to another or from a datastore
				// to another (or both). This requires the VM to be stopped.
				if err := vmShutdown(ctx, vmAPI, d); err != nil {
					return err
				}

				if err := deleteIdeDrives(ctx, vmAPI, initializationInterface, existingInterface); err != nil {
					return err
				}

				stoppedBeforeUpdate = true
				fileVolume = fmt.Sprintf("%s:cloudinit", initializationDatastoreID)
			} else {
				ideDevice := getIdeDevice(vmConfig, existingInterface)
				fileVolume = ideDevice.FileVolume
			}

			updateBody.IDEDevices[initializationInterface] = &vms.CustomStorageDevice{
				Enabled:    true,
				FileVolume: fileVolume,
				Media:      &cdromMedia,
			}
		}

		rebootRequired = true
	}

	// Prepare the new hostpci devices configuration.
	if d.HasChange(mkHostPCI) {
		updateBody.PCIDevices = vmGetHostPCIDeviceObjects(d)

		for i := len(updateBody.PCIDevices); i < maxResourceVirtualEnvironmentVMHostPCIDevices; i++ {
			del = append(del, fmt.Sprintf("hostpci%d", i))
		}

		rebootRequired = true
	}

	// Prepare the new usb devices configuration.
	if d.HasChange(mkHostUSB) {
		updateBody.USBDevices = vmGetHostUSBDeviceObjects(d)

		for i := len(updateBody.USBDevices); i < maxResourceVirtualEnvironmentVMHostUSBDevices; i++ {
			del = append(del, fmt.Sprintf("usb%d", i))
		}

		rebootRequired = true
	}

	// Prepare the new memory configuration.
	if d.HasChange(mkMemory) {
		memoryBlock, err := structure.GetSchemaBlock(
			resource,
			d,
			[]string{mkMemory},
			0,
			true,
		)
		if err != nil {
			return diag.FromErr(err)
		}

		memoryDedicated := memoryBlock[mkMemoryDedicated].(int)
		memoryFloating := memoryBlock[mkMemoryFloating].(int)
		memoryShared := memoryBlock[mkMemoryShared].(int)

		updateBody.DedicatedMemory = &memoryDedicated
		updateBody.FloatingMemory = &memoryFloating

		if memoryShared > 0 {
			memorySharedName := fmt.Sprintf("vm-%d-ivshmem", vmID)

			updateBody.SharedMemory = &vms.CustomSharedMemory{
				Name: &memorySharedName,
				Size: memoryShared,
			}
		}

		rebootRequired = true
	}

	// Prepare the new network device configuration.
	if d.HasChange(mkNetworkDevice) {
		updateBody.NetworkDevices = vmGetNetworkDeviceObjects(d)

		for i := 0; i < len(updateBody.NetworkDevices); i++ {
			if !updateBody.NetworkDevices[i].Enabled {
				del = append(del, fmt.Sprintf("net%d", i))
			}
		}

		for i := len(updateBody.NetworkDevices); i < maxResourceVirtualEnvironmentVMNetworkDevices; i++ {
			del = append(del, fmt.Sprintf("net%d", i))
		}

		rebootRequired = true
	}

	// Prepare the new operating system configuration.
	if d.HasChange(mkOperatingSystem) {
		operatingSystem, err := structure.GetSchemaBlock(
			resource,
			d,
			[]string{mkOperatingSystem},
			0,
			true,
		)
		if err != nil {
			return diag.FromErr(err)
		}

		operatingSystemType := operatingSystem[mkOperatingSystemType].(string)

		updateBody.OSType = &operatingSystemType

		rebootRequired = true
	}

	// Prepare the new serial devices.
	if d.HasChange(mkSerialDevice) {
		updateBody.SerialDevices = vmGetSerialDeviceList(d)

		for i := len(updateBody.SerialDevices); i < maxResourceVirtualEnvironmentVMSerialDevices; i++ {
			del = append(del, fmt.Sprintf("serial%d", i))
		}

		rebootRequired = true
	}

	if d.HasChange(mkSMBIOS) {
		updateBody.SMBIOS = vmGetSMBIOS(d)
		if updateBody.SMBIOS == nil {
			del = append(del, "smbios1")
		}
	}

	if d.HasChange(mkStartup) {
		updateBody.StartupOrder = vmGetStartupOrder(d)
		if updateBody.StartupOrder == nil {
			del = append(del, "startup")
		}
	}

	// Prepare the new VGA configuration.
	if d.HasChange(mkVGA) {
		updateBody.VGADevice, e = vmGetVGADeviceObject(d)
		if e != nil {
			return diag.FromErr(e)
		}

		rebootRequired = true
	}

	// Prepare the new SCSI hardware type
	if d.HasChange(mkSCSIHardware) {
		scsiHardware := d.Get(mkSCSIHardware).(string)
		updateBody.SCSIHardware = &scsiHardware

		rebootRequired = true
	}

	if d.HasChanges(mkHookScriptFileID) {
		hookScript := d.Get(mkHookScriptFileID).(string)
		if len(hookScript) > 0 {
			updateBody.HookScript = &hookScript
		} else {
			del = append(del, "hookscript")
		}
	}

	// Update the configuration now that everything has been prepared.
	updateBody.Delete = del

	e = vmAPI.UpdateVM(ctx, updateBody)
	if e != nil {
		return diag.FromErr(e)
	}

	// Determine if the state of the virtual machine state needs to be changed.
	//nolint: nestif
	if (d.HasChange(mkStarted) || stoppedBeforeUpdate) && !bool(template) {
		started := d.Get(mkStarted).(bool)
		if started {
			if diags := vmStart(ctx, vmAPI, d); diags != nil {
				return diags
			}
		} else {
			if e := vmShutdown(ctx, vmAPI, d); e != nil {
				return e
			}

			rebootRequired = false
		}
	}

	// Change the disk locations and/or sizes, if necessary.
	return vmUpdateDiskLocationAndSize(
		ctx,
		d,
		m,
		!bool(template) && rebootRequired,
	)
}

func vmUpdateDiskLocationAndSize(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
	reboot bool,
) diag.Diagnostics {
	config := m.(proxmoxtf.ProviderConfiguration)

	api, err := config.GetClient()
	if err != nil {
		return diag.FromErr(err)
	}

	nodeName := d.Get(mkNodeName).(string)
	started := d.Get(mkStarted).(bool)
	template := d.Get(mkTemplate).(bool)

	vmID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	vmAPI := api.Node(nodeName).VM(vmID)

	// Determine if any of the disks are changing location and/or size, and initiate the necessary actions.
	if d.HasChange(mkDisk) {
		diskOld, diskNew := d.GetChange(mkDisk)

		diskOldEntries, err := getDiskDeviceObjects1(d, diskOld.([]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}

		diskNewEntries, err := getDiskDeviceObjects1(d, diskNew.([]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}

		// Add efidisk if it has changes
		if d.HasChange(mkEFIDisk) {
			diskOld, diskNew := d.GetChange(mkEFIDisk)

			oldEfiDisk, e := vmGetEfiDiskAsStorageDevice(d, diskOld.([]interface{}))
			if e != nil {
				return diag.FromErr(e)
			}

			newEfiDisk, e := vmGetEfiDiskAsStorageDevice(d, diskNew.([]interface{}))
			if e != nil {
				return diag.FromErr(e)
			}

			if oldEfiDisk != nil {
				diskOldEntries[*oldEfiDisk.Interface] = oldEfiDisk
			}

			if newEfiDisk != nil {
				diskNewEntries[*newEfiDisk.Interface] = newEfiDisk
			}

			if oldEfiDisk != nil && newEfiDisk != nil && oldEfiDisk.Size != newEfiDisk.Size {
				return diag.Errorf(
					"resizing of efidisks is not supported.",
				)
			}
		}

		// Add tpm state if it has changes
		if d.HasChange(mkTPMState) {
			diskOld, diskNew := d.GetChange(mkTPMState)

			oldTPMState := vmGetTPMStateAsStorageDevice(d, diskOld.([]interface{}))
			newTPMState := vmGetTPMStateAsStorageDevice(d, diskNew.([]interface{}))

			if oldTPMState != nil {
				diskOldEntries[*oldTPMState.Interface] = oldTPMState
			}

			if newTPMState != nil {
				diskNewEntries[*newTPMState.Interface] = newTPMState
			}

			if oldTPMState != nil && newTPMState != nil && oldTPMState.Size != newTPMState.Size {
				return diag.Errorf(
					"resizing of tpm state is not supported.",
				)
			}
		}

		// TODO: move to disks.go

		var diskMoveBodies []*vms.MoveDiskRequestBody

		var diskResizeBodies []*vms.ResizeDiskRequestBody

		shutdownForDisksRequired := false

		for oldKey, oldDisk := range diskOldEntries {
			if _, present := diskNewEntries[oldKey]; !present {
				return diag.Errorf(
					"deletion of disks not supported. Please delete disk by hand. Old Interface was %s",
					*oldDisk.Interface,
				)
			}

			if *oldDisk.DatastoreID != *diskNewEntries[oldKey].DatastoreID {
				if oldDisk.IsOwnedBy(vmID) {
					deleteOriginalDisk := types.CustomBool(true)

					diskMoveBodies = append(
						diskMoveBodies,
						&vms.MoveDiskRequestBody{
							DeleteOriginalDisk: &deleteOriginalDisk,
							Disk:               *oldDisk.Interface,
							TargetStorage:      *diskNewEntries[oldKey].DatastoreID,
						},
					)

					// Cannot be done while VM is running.
					shutdownForDisksRequired = true
				} else {
					return diag.Errorf(
						"Cannot move %s:%s to datastore %s in VM %d configuration, it is not owned by this VM!",
						*oldDisk.DatastoreID,
						*oldDisk.PathInDatastore(),
						*diskNewEntries[oldKey].DatastoreID,
						vmID,
					)
				}
			}

			if *oldDisk.Size < *diskNewEntries[oldKey].Size {
				if oldDisk.IsOwnedBy(vmID) {
					diskResizeBodies = append(
						diskResizeBodies,
						&vms.ResizeDiskRequestBody{
							Disk: *oldDisk.Interface,
							Size: *diskNewEntries[oldKey].Size,
						},
					)
				} else {
					return diag.Errorf(
						"Cannot resize %s:%s in VM %d configuration, it is not owned by this VM!",
						*oldDisk.DatastoreID,
						*oldDisk.PathInDatastore(),
						vmID,
					)
				}
			}
		}

		if shutdownForDisksRequired && !template {
			if e := vmShutdown(ctx, vmAPI, d); e != nil {
				return e
			}
		}

		for _, reqBody := range diskMoveBodies {
			moveDiskTimeout := d.Get(mkTimeoutMoveDisk).(int)

			err = vmAPI.MoveVMDisk(ctx, reqBody, moveDiskTimeout)
			if err != nil {
				return diag.FromErr(err)
			}
		}

		for _, reqBody := range diskResizeBodies {
			moveDiskTimeout := d.Get(mkTimeoutMoveDisk).(int)

			err = vmAPI.ResizeVMDisk(ctx, reqBody, moveDiskTimeout)
			if err != nil {
				return diag.FromErr(err)
			}
		}

		if shutdownForDisksRequired && started && !template {
			if diags := vmStart(ctx, vmAPI, d); diags != nil {
				return diags
			}

			// This concludes an equivalent of a reboot, avoid doing another.
			reboot = false
		}
	}

	// Perform a regular reboot in case it's necessary and haven't already been done.
	if reboot {
		vmStatus, err := vmAPI.GetVMStatus(ctx)
		if err != nil {
			return diag.FromErr(err)
		}

		if vmStatus.Status != "stopped" {
			rebootTimeout := d.Get(mkTimeoutReboot).(int)

			err := vmAPI.RebootVM(
				ctx,
				&vms.RebootRequestBody{
					Timeout: &rebootTimeout,
				},
				rebootTimeout+30,
			)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return vmRead(ctx, d, m)
}

func vmDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(proxmoxtf.ProviderConfiguration)

	api, err := config.GetClient()
	if err != nil {
		return diag.FromErr(err)
	}

	nodeName := d.Get(mkNodeName).(string)

	vmID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	vmAPI := api.Node(nodeName).VM(vmID)

	// Stop or shut down the virtual machine before deleting it.
	status, err := vmAPI.GetVMStatus(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	stop := d.Get(mkStopOnDestroy).(bool)

	if status.Status != "stopped" {
		if stop {
			if e := vmStop(ctx, vmAPI, d); e != nil {
				return e
			}
		} else {
			if e := vmShutdown(ctx, vmAPI, d); e != nil {
				return e
			}
		}
	}

	err = vmAPI.DeleteVM(ctx)

	if err != nil {
		if strings.Contains(err.Error(), "HTTP 404") ||
			(strings.Contains(err.Error(), "HTTP 500") && strings.Contains(err.Error(), "does not exist")) {
			d.SetId("")

			return nil
		}

		return diag.FromErr(err)
	}

	// Wait for the state to become unavailable as that clearly indicates the destruction of the VM.
	err = vmAPI.WaitForVMStatus(ctx, "", 60, 2)
	if err == nil {
		return diag.Errorf("failed to delete VM \"%d\"", vmID)
	}

	d.SetId("")

	return nil
}

func getPCIInfo(resp *vms.GetResponseData, _ *schema.ResourceData) map[string]*vms.CustomPCIDevice {
	pciDevices := map[string]*vms.CustomPCIDevice{}

	pciDevices["hostpci0"] = resp.PCIDevice0
	pciDevices["hostpci1"] = resp.PCIDevice1
	pciDevices["hostpci2"] = resp.PCIDevice2
	pciDevices["hostpci3"] = resp.PCIDevice3

	return pciDevices
}

func getUSBInfo(resp *vms.GetResponseData, _ *schema.ResourceData) map[string]*vms.CustomUSBDevice {
	usbDevices := map[string]*vms.CustomUSBDevice{}

	usbDevices["usb0"] = resp.USBDevice0
	usbDevices["usb1"] = resp.USBDevice1
	usbDevices["usb2"] = resp.USBDevice2
	usbDevices["usb3"] = resp.USBDevice3

	return usbDevices
}

func parseImportIDWithNodeName(id string) (string, string, error) {
	nodeName, id, found := strings.Cut(id, "/")

	if !found {
		return "", "", fmt.Errorf("unexpected format of ID (%s), expected node/id", id)
	}

	return nodeName, id, nil
}
