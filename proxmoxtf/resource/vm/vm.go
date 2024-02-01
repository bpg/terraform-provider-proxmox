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
	"unicode"

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
	clone := d.Get(mkResourceVirtualEnvironmentVMClone).([]interface{})

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

	startVMTimeout := d.Get(mkResourceVirtualEnvironmentVMTimeoutStartVM).(int)

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

	return append(diags, diag.FromErr(vmAPI.WaitForVMState(ctx, "running", startVMTimeout, 1))...)
}

// Shutdown the VM, then wait for it to actually shut down (it may not be shut down immediately if
// running in HA mode).
func vmShutdown(ctx context.Context, vmAPI *vms.Client, d *schema.ResourceData) diag.Diagnostics {
	tflog.Debug(ctx, "Shutting down VM")

	forceStop := types.CustomBool(true)
	shutdownTimeout := d.Get(mkResourceVirtualEnvironmentVMTimeoutShutdownVM).(int)

	e := vmAPI.ShutdownVM(ctx, &vms.ShutdownRequestBody{
		ForceStop: &forceStop,
		Timeout:   &shutdownTimeout,
	}, shutdownTimeout+30)
	if e != nil {
		return diag.FromErr(e)
	}

	return diag.FromErr(vmAPI.WaitForVMState(ctx, "stopped", shutdownTimeout, 1))
}

// Forcefully stop the VM, then wait for it to actually stop.
func vmStop(ctx context.Context, vmAPI *vms.Client, d *schema.ResourceData) diag.Diagnostics {
	tflog.Debug(ctx, "Stopping VM")

	stopTimeout := d.Get(mkResourceVirtualEnvironmentVMTimeoutStopVM).(int)

	e := vmAPI.StopVM(ctx, stopTimeout+30)
	if e != nil {
		return diag.FromErr(e)
	}

	return diag.FromErr(vmAPI.WaitForVMState(ctx, "stopped", stopTimeout, 1))
}

func vmCreateClone(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(proxmoxtf.ProviderConfiguration)

	api, e := config.GetClient()
	if e != nil {
		return diag.FromErr(e)
	}

	clone := d.Get(mkResourceVirtualEnvironmentVMClone).([]interface{})
	cloneBlock := clone[0].(map[string]interface{})
	cloneRetries := cloneBlock[mkResourceVirtualEnvironmentVMCloneRetries].(int)
	cloneDatastoreID := cloneBlock[mkResourceVirtualEnvironmentVMCloneDatastoreID].(string)
	cloneNodeName := cloneBlock[mkResourceVirtualEnvironmentVMCloneNodeName].(string)
	cloneVMID := cloneBlock[mkResourceVirtualEnvironmentVMCloneVMID].(int)
	cloneFull := cloneBlock[mkResourceVirtualEnvironmentVMCloneFull].(bool)

	description := d.Get(mkResourceVirtualEnvironmentVMDescription).(string)
	name := d.Get(mkResourceVirtualEnvironmentVMName).(string)
	tags := d.Get(mkResourceVirtualEnvironmentVMTags).([]interface{})
	nodeName := d.Get(mkResourceVirtualEnvironmentVMNodeName).(string)
	poolID := d.Get(mkResourceVirtualEnvironmentVMPoolID).(string)
	vmIDUntyped, hasVMID := d.GetOk(mkResourceVirtualEnvironmentVMVMID)
	vmID := vmIDUntyped.(int)

	if !hasVMID {
		vmIDNew, err := api.Cluster().GetVMID(ctx)
		if err != nil {
			return diag.FromErr(err)
		}

		vmID = *vmIDNew
		err = d.Set(mkResourceVirtualEnvironmentVMVMID, vmID)

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

	cloneTimeout := d.Get(mkResourceVirtualEnvironmentVMTimeoutClone).(int)

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

	// Now that the virtual machine has been cloned, we need to perform some modifications.
	acpi := types.CustomBool(d.Get(mkResourceVirtualEnvironmentVMACPI).(bool))
	audioDevices := vmGetAudioDeviceList(d)

	bios := d.Get(mkResourceVirtualEnvironmentVMBIOS).(string)
	kvmArguments := d.Get(mkResourceVirtualEnvironmentVMKVMArguments).(string)
	scsiHardware := d.Get(mkResourceVirtualEnvironmentVMSCSIHardware).(string)
	cdrom := d.Get(mkResourceVirtualEnvironmentVMCDROM).([]interface{})
	cpu := d.Get(mkResourceVirtualEnvironmentVMCPU).([]interface{})
	initialization := d.Get(mkResourceVirtualEnvironmentVMInitialization).([]interface{})
	hostPCI := d.Get(mkResourceVirtualEnvironmentVMHostPCI).([]interface{})
	hostUSB := d.Get(mkResourceVirtualEnvironmentVMHostUSB).([]interface{})
	keyboardLayout := d.Get(mkResourceVirtualEnvironmentVMKeyboardLayout).(string)
	memory := d.Get(mkResourceVirtualEnvironmentVMMemory).([]interface{})
	networkDevice := d.Get(mkResourceVirtualEnvironmentVMNetworkDevice).([]interface{})
	operatingSystem := d.Get(mkResourceVirtualEnvironmentVMOperatingSystem).([]interface{})
	serialDevice := d.Get(mkResourceVirtualEnvironmentVMSerialDevice).([]interface{})
	onBoot := types.CustomBool(d.Get(mkResourceVirtualEnvironmentVMOnBoot).(bool))
	tabletDevice := types.CustomBool(d.Get(mkResourceVirtualEnvironmentVMTabletDevice).(bool))
	template := types.CustomBool(d.Get(mkResourceVirtualEnvironmentVMTemplate).(bool))
	vga := d.Get(mkResourceVirtualEnvironmentVMVGA).([]interface{})

	updateBody := &vms.UpdateRequestBody{
		AudioDevices: audioDevices,
	}

	ideDevices := vms.CustomStorageDevices{}

	var del []string

	//nolint:gosimple
	if acpi != dvResourceVirtualEnvironmentVMACPI {
		updateBody.ACPI = &acpi
	}

	createAgent(d, updateBody)

	if kvmArguments != dvResourceVirtualEnvironmentVMKVMArguments {
		updateBody.KVMArguments = &kvmArguments
	}

	if bios != dvResourceVirtualEnvironmentVMBIOS {
		updateBody.BIOS = &bios
	}

	if scsiHardware != dvResourceVirtualEnvironmentVMSCSIHardware {
		updateBody.SCSIHardware = &scsiHardware
	}

	if len(cdrom) > 0 || len(initialization) > 0 {
		ideDevices = vms.CustomStorageDevices{
			"ide0": vms.CustomStorageDevice{
				Enabled: false,
			},
			"ide1": vms.CustomStorageDevice{
				Enabled: false,
			},
			"ide2": vms.CustomStorageDevice{
				Enabled: false,
			},
			"ide3": vms.CustomStorageDevice{
				Enabled: false,
			},
		}
	}

	if len(cdrom) > 0 {
		cdromBlock := cdrom[0].(map[string]interface{})

		cdromEnabled := cdromBlock[mkResourceVirtualEnvironmentVMCDROMEnabled].(bool)
		cdromFileID := cdromBlock[mkResourceVirtualEnvironmentVMCDROMFileID].(string)
		cdromInterface := cdromBlock[mkResourceVirtualEnvironmentVMCDROMInterface].(string)

		if cdromFileID == "" {
			cdromFileID = "cdrom"
		}

		cdromMedia := "cdrom"

		ideDevices[cdromInterface] = vms.CustomStorageDevice{
			Enabled:    cdromEnabled,
			FileVolume: cdromFileID,
			Media:      &cdromMedia,
		}
	}

	if len(cpu) > 0 {
		cpuBlock := cpu[0].(map[string]interface{})

		cpuArchitecture := cpuBlock[mkResourceVirtualEnvironmentVMCPUArchitecture].(string)
		cpuCores := cpuBlock[mkResourceVirtualEnvironmentVMCPUCores].(int)
		cpuFlags := cpuBlock[mkResourceVirtualEnvironmentVMCPUFlags].([]interface{})
		cpuHotplugged := cpuBlock[mkResourceVirtualEnvironmentVMCPUHotplugged].(int)
		cpuLimit := cpuBlock[mkResourceVirtualEnvironmentVMCPULimit].(int)
		cpuNUMA := types.CustomBool(cpuBlock[mkResourceVirtualEnvironmentVMCPUNUMA].(bool))
		cpuSockets := cpuBlock[mkResourceVirtualEnvironmentVMCPUSockets].(int)
		cpuType := cpuBlock[mkResourceVirtualEnvironmentVMCPUType].(string)
		cpuUnits := cpuBlock[mkResourceVirtualEnvironmentVMCPUUnits].(int)

		cpuFlagsConverted := make([]string, len(cpuFlags))

		for fi, flag := range cpuFlags {
			cpuFlagsConverted[fi] = flag.(string)
		}

		// Only the root account is allowed to change the CPU architecture, which makes this check necessary.
		if api.API().IsRootTicket() ||
			cpuArchitecture != dvResourceVirtualEnvironmentVMCPUArchitecture {
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
		initializationDatastoreID := initializationBlock[mkResourceVirtualEnvironmentVMInitializationDatastoreID].(string)
		initializationInterface := initializationBlock[mkResourceVirtualEnvironmentVMInitializationInterface].(string)

		existingInterface := findExistingCloudInitDrive(vmConfig, vmID, "ide2")
		if initializationInterface == "" {
			initializationInterface = existingInterface
		}

		tflog.Trace(ctx, fmt.Sprintf("CloudInit IDE interface is '%s'", initializationInterface))

		const cdromCloudInitEnabled = true

		cdromCloudInitFileID := fmt.Sprintf("%s:cloudinit", initializationDatastoreID)
		cdromCloudInitMedia := "cdrom"
		ideDevices[initializationInterface] = vms.CustomStorageDevice{
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

	if keyboardLayout != dvResourceVirtualEnvironmentVMKeyboardLayout {
		updateBody.KeyboardLayout = &keyboardLayout
	}

	if len(memory) > 0 {
		memoryBlock := memory[0].(map[string]interface{})

		memoryDedicated := memoryBlock[mkResourceVirtualEnvironmentVMMemoryDedicated].(int)
		memoryFloating := memoryBlock[mkResourceVirtualEnvironmentVMMemoryFloating].(int)
		memoryShared := memoryBlock[mkResourceVirtualEnvironmentVMMemoryShared].(int)

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
		operatingSystemType := operatingSystemBlock[mkResourceVirtualEnvironmentVMOperatingSystemType].(string)

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
	if tabletDevice != dvResourceVirtualEnvironmentVMTabletDevice {
		updateBody.TabletDeviceEnabled = &tabletDevice
	}

	if len(tags) > 0 {
		tagString := vmGetTagsString(d)
		updateBody.Tags = &tagString
	}

	//nolint:gosimple
	if template != dvResourceVirtualEnvironmentVMTemplate {
		updateBody.Template = &template
	}

	if len(vga) > 0 {
		vgaDevice, err := vmGetVGADeviceObject(d)
		if err != nil {
			return diag.FromErr(err)
		}

		updateBody.VGADevice = vgaDevice
	}

	hookScript := d.Get(mkResourceVirtualEnvironmentVMHookScriptFileID).(string)
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

	disk := d.Get(mkResourceVirtualEnvironmentVMDisk).([]interface{})
	efiDisk := d.Get(mkResourceVirtualEnvironmentVMEFIDisk).([]interface{})

	vmConfig, e = vmAPI.GetVM(ctx)
	if e != nil {
		if strings.Contains(e.Error(), "HTTP 404") ||
			(strings.Contains(e.Error(), "HTTP 500") && strings.Contains(e.Error(), "does not exist")) {
			d.SetId("")

			return nil
		}

		return diag.FromErr(e)
	}

	allDiskInfo := getDiskInfo(vmConfig, d) // from the cloned VM

	diskDeviceObjects, e := vmGetDiskDeviceObjects(d, nil) // from the resource config
	if e != nil {
		return diag.FromErr(e)
	}

	for i := range disk {
		diskBlock := disk[i].(map[string]interface{})
		diskInterface := diskBlock[mkResourceVirtualEnvironmentVMDiskInterface].(string)
		dataStoreID := diskBlock[mkResourceVirtualEnvironmentVMDiskDatastoreID].(string)
		diskSize := int64(diskBlock[mkResourceVirtualEnvironmentVMDiskSize].(int))
		prefix := diskDigitPrefix(diskInterface)

		currentDiskInfo := allDiskInfo[diskInterface]
		configuredDiskInfo := diskDeviceObjects[prefix][diskInterface]

		if currentDiskInfo == nil {
			diskUpdateBody := &vms.UpdateRequestBody{}

			switch prefix {
			case "virtio":
				if diskUpdateBody.VirtualIODevices == nil {
					diskUpdateBody.VirtualIODevices = vms.CustomStorageDevices{}
				}

				diskUpdateBody.VirtualIODevices[diskInterface] = configuredDiskInfo
			case "sata":
				if diskUpdateBody.SATADevices == nil {
					diskUpdateBody.SATADevices = vms.CustomStorageDevices{}
				}

				diskUpdateBody.SATADevices[diskInterface] = configuredDiskInfo
			case "scsi":
				if diskUpdateBody.SCSIDevices == nil {
					diskUpdateBody.SCSIDevices = vms.CustomStorageDevices{}
				}

				diskUpdateBody.SCSIDevices[diskInterface] = configuredDiskInfo
			}

			e = vmAPI.UpdateVM(ctx, diskUpdateBody)
			if e != nil {
				return diag.FromErr(e)
			}

			continue
		}

		if diskSize < currentDiskInfo.Size.InGigabytes() {
			return diag.Errorf(
				"disk resize fails requests size (%dG) is lower than current size (%s)",
				diskSize,
				*currentDiskInfo.Size,
			)
		}

		deleteOriginalDisk := types.CustomBool(true)

		diskMoveBody := &vms.MoveDiskRequestBody{
			DeleteOriginalDisk: &deleteOriginalDisk,
			Disk:               diskInterface,
			TargetStorage:      dataStoreID,
		}

		diskResizeBody := &vms.ResizeDiskRequestBody{
			Disk: diskInterface,
			Size: types.DiskSizeFromGigabytes(diskSize),
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
			moveDiskTimeout := d.Get(mkResourceVirtualEnvironmentVMTimeoutMoveDisk).(int)

			e = vmAPI.MoveVMDisk(ctx, diskMoveBody, moveDiskTimeout)
			if e != nil {
				return diag.FromErr(e)
			}
		}

		if diskSize > currentDiskInfo.Size.InGigabytes() {
			e = vmAPI.ResizeVMDisk(ctx, diskResizeBody)
			if e != nil {
				return diag.FromErr(e)
			}
		}
	}

	efiDiskInfo := vmGetEfiDisk(d, nil) // from the resource config

	for i := range efiDisk {
		diskBlock := efiDisk[i].(map[string]interface{})
		diskInterface := "efidisk0"
		dataStoreID := diskBlock[mkResourceVirtualEnvironmentVMEFIDiskDatastoreID].(string)
		efiType := diskBlock[mkResourceVirtualEnvironmentVMEFIDiskType].(string)

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
			moveDiskTimeout := d.Get(mkResourceVirtualEnvironmentVMTimeoutMoveDisk).(int)

			e = vmAPI.MoveVMDisk(ctx, diskMoveBody, moveDiskTimeout)
			if e != nil {
				return diag.FromErr(e)
			}
		}
	}

	tpmState := d.Get(mkResourceVirtualEnvironmentVMTPMState).([]interface{})
	tpmStateInfo := vmGetTPMState(d, nil) // from the resource config

	for i := range tpmState {
		diskBlock := tpmState[i].(map[string]interface{})
		diskInterface := "tpmstate0"
		dataStoreID := diskBlock[mkResourceVirtualEnvironmentVMTPMStateDatastoreID].(string)

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
			moveDiskTimeout := d.Get(mkResourceVirtualEnvironmentVMTimeoutMoveDisk).(int)

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

	acpi := types.CustomBool(d.Get(mkResourceVirtualEnvironmentVMACPI).(bool))

	customAgent, err := customAgent(d, resource)
	if err != nil {
		return diag.FromErr(err)
	}

	kvmArguments := d.Get(mkResourceVirtualEnvironmentVMKVMArguments).(string)

	audioDevices := vmGetAudioDeviceList(d)

	bios := d.Get(mkResourceVirtualEnvironmentVMBIOS).(string)

	cdromBlock, err := structure.GetSchemaBlock(
		resource,
		d,
		[]string{mkResourceVirtualEnvironmentVMCDROM},
		0,
		true,
	)
	if err != nil {
		return diag.FromErr(err)
	}

	cdromEnabled := cdromBlock[mkResourceVirtualEnvironmentVMCDROMEnabled].(bool)
	cdromFileID := cdromBlock[mkResourceVirtualEnvironmentVMCDROMFileID].(string)
	cdromInterface := cdromBlock[mkResourceVirtualEnvironmentVMCDROMInterface].(string)

	cdromCloudInitEnabled := false
	cdromCloudInitFileID := ""
	cdromCloudInitInterface := ""

	if cdromFileID == "" {
		cdromFileID = "cdrom"
	}

	cpuBlock, err := structure.GetSchemaBlock(
		resource,
		d,
		[]string{mkResourceVirtualEnvironmentVMCPU},
		0,
		true,
	)
	if err != nil {
		return diag.FromErr(err)
	}

	cpuArchitecture := cpuBlock[mkResourceVirtualEnvironmentVMCPUArchitecture].(string)
	cpuCores := cpuBlock[mkResourceVirtualEnvironmentVMCPUCores].(int)
	cpuFlags := cpuBlock[mkResourceVirtualEnvironmentVMCPUFlags].([]interface{})
	cpuHotplugged := cpuBlock[mkResourceVirtualEnvironmentVMCPUHotplugged].(int)
	cpuLimit := cpuBlock[mkResourceVirtualEnvironmentVMCPULimit].(int)
	cpuSockets := cpuBlock[mkResourceVirtualEnvironmentVMCPUSockets].(int)
	cpuNUMA := types.CustomBool(cpuBlock[mkResourceVirtualEnvironmentVMCPUNUMA].(bool))
	cpuType := cpuBlock[mkResourceVirtualEnvironmentVMCPUType].(string)
	cpuUnits := cpuBlock[mkResourceVirtualEnvironmentVMCPUUnits].(int)

	description := d.Get(mkResourceVirtualEnvironmentVMDescription).(string)
	diskDeviceObjects, err := vmGetDiskDeviceObjects(d, nil)

	if err != nil {
		return diag.FromErr(err)
	}

	var efiDisk *vms.CustomEFIDisk

	efiDiskBlock := d.Get(mkResourceVirtualEnvironmentVMEFIDisk).([]interface{})
	if len(efiDiskBlock) > 0 {
		block := efiDiskBlock[0].(map[string]interface{})

		datastoreID, _ := block[mkResourceVirtualEnvironmentVMEFIDiskDatastoreID].(string)
		fileFormat, _ := block[mkResourceVirtualEnvironmentVMEFIDiskFileFormat].(string)
		efiType, _ := block[mkResourceVirtualEnvironmentVMEFIDiskType].(string)
		preEnrolledKeys := types.CustomBool(block[mkResourceVirtualEnvironmentVMEFIDiskPreEnrolledKeys].(bool))

		if fileFormat == "" {
			fileFormat = dvResourceVirtualEnvironmentVMEFIDiskFileFormat
		}

		efiDisk = &vms.CustomEFIDisk{
			Type:            &efiType,
			FileVolume:      fmt.Sprintf("%s:1", datastoreID),
			Format:          &fileFormat,
			PreEnrolledKeys: &preEnrolledKeys,
		}
	}

	var tpmState *vms.CustomTPMState

	tpmStateBlock := d.Get(mkResourceVirtualEnvironmentVMTPMState).([]interface{})
	if len(tpmStateBlock) > 0 {
		block := tpmStateBlock[0].(map[string]interface{})

		datastoreID, _ := block[mkResourceVirtualEnvironmentVMTPMStateDatastoreID].(string)
		version, _ := block[mkResourceVirtualEnvironmentVMTPMStateVersion].(string)

		if version == "" {
			version = dvResourceVirtualEnvironmentVMTPMStateVersion
		}

		tpmState = &vms.CustomTPMState{
			FileVolume: fmt.Sprintf("%s:1", datastoreID),
			Version:    &version,
		}
	}

	virtioDeviceObjects := diskDeviceObjects["virtio"]
	scsiDeviceObjects := diskDeviceObjects["scsi"]
	// ideDeviceObjects := getOrderedDiskDeviceList(diskDeviceObjects, "ide")
	sataDeviceObjects := diskDeviceObjects["sata"]

	initializationConfig := vmGetCloudInitConfig(d)

	if initializationConfig != nil {
		initialization := d.Get(mkResourceVirtualEnvironmentVMInitialization).([]interface{})
		initializationBlock := initialization[0].(map[string]interface{})
		initializationDatastoreID := initializationBlock[mkResourceVirtualEnvironmentVMInitializationDatastoreID].(string)

		cdromCloudInitEnabled = true
		cdromCloudInitFileID = fmt.Sprintf("%s:cloudinit", initializationDatastoreID)

		cdromCloudInitInterface = initializationBlock[mkResourceVirtualEnvironmentVMInitializationInterface].(string)
		if cdromCloudInitInterface == "" {
			cdromCloudInitInterface = "ide2"
		}
	}

	pciDeviceObjects := vmGetHostPCIDeviceObjects(d)

	usbDeviceObjects := vmGetHostUSBDeviceObjects(d)

	keyboardLayout := d.Get(mkResourceVirtualEnvironmentVMKeyboardLayout).(string)
	memoryBlock, err := structure.GetSchemaBlock(
		resource,
		d,
		[]string{mkResourceVirtualEnvironmentVMMemory},
		0,
		true,
	)
	if err != nil {
		return diag.FromErr(err)
	}

	memoryDedicated := memoryBlock[mkResourceVirtualEnvironmentVMMemoryDedicated].(int)
	memoryFloating := memoryBlock[mkResourceVirtualEnvironmentVMMemoryFloating].(int)
	memoryShared := memoryBlock[mkResourceVirtualEnvironmentVMMemoryShared].(int)

	machine := d.Get(mkResourceVirtualEnvironmentVMMachine).(string)
	name := d.Get(mkResourceVirtualEnvironmentVMName).(string)
	tags := d.Get(mkResourceVirtualEnvironmentVMTags).([]interface{})

	networkDeviceObjects := vmGetNetworkDeviceObjects(d)

	nodeName := d.Get(mkResourceVirtualEnvironmentVMNodeName).(string)

	operatingSystem, err := structure.GetSchemaBlock(
		resource,
		d,
		[]string{mkResourceVirtualEnvironmentVMOperatingSystem},
		0,
		true,
	)
	if err != nil {
		return diag.FromErr(err)
	}

	operatingSystemType := operatingSystem[mkResourceVirtualEnvironmentVMOperatingSystemType].(string)

	poolID := d.Get(mkResourceVirtualEnvironmentVMPoolID).(string)

	serialDevices := vmGetSerialDeviceList(d)

	smbios := vmGetSMBIOS(d)

	startupOrder := vmGetStartupOrder(d)

	onBoot := types.CustomBool(d.Get(mkResourceVirtualEnvironmentVMOnBoot).(bool))
	tabletDevice := types.CustomBool(d.Get(mkResourceVirtualEnvironmentVMTabletDevice).(bool))
	template := types.CustomBool(d.Get(mkResourceVirtualEnvironmentVMTemplate).(bool))

	vgaDevice, err := vmGetVGADeviceObject(d)
	if err != nil {
		return diag.FromErr(err)
	}

	vmIDUntyped, hasVMID := d.GetOk(mkResourceVirtualEnvironmentVMVMID)
	vmID := vmIDUntyped.(int)

	if !hasVMID {
		vmIDNew, e := api.Cluster().GetVMID(ctx)
		if e != nil {
			return diag.FromErr(e)
		}

		vmID = *vmIDNew
		e = d.Set(mkResourceVirtualEnvironmentVMVMID, vmID)

		if e != nil {
			return diag.FromErr(e)
		}
	}

	var memorySharedObject *vms.CustomSharedMemory

	var bootOrderConverted []string
	if cdromEnabled {
		bootOrderConverted = []string{cdromInterface}
	}

	bootOrder := d.Get(mkResourceVirtualEnvironmentVMBootOrder).([]interface{})
	//nolint:nestif
	if len(bootOrder) == 0 {
		if sataDeviceObjects != nil {
			bootOrderConverted = append(bootOrderConverted, "sata0")
		}

		if scsiDeviceObjects != nil {
			bootOrderConverted = append(bootOrderConverted, "scsi0")
		}

		if virtioDeviceObjects != nil {
			bootOrderConverted = append(bootOrderConverted, "virtio0")
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
		cdromCloudInitInterface: vms.CustomStorageDevice{
			Enabled:    cdromCloudInitEnabled,
			FileVolume: cdromCloudInitFileID,
			Media:      &ideDevice2Media,
		},
		cdromInterface: vms.CustomStorageDevice{
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

	scsiHardware := d.Get(mkResourceVirtualEnvironmentVMSCSIHardware).(string)

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

	if sataDeviceObjects != nil {
		createBody.SATADevices = sataDeviceObjects
	}

	if scsiDeviceObjects != nil {
		createBody.SCSIDevices = scsiDeviceObjects
	}

	if virtioDeviceObjects != nil {
		createBody.VirtualIODevices = virtioDeviceObjects
	}

	// Only the root account is allowed to change the CPU architecture, which makes this check necessary.
	if api.API().IsRootTicket() ||
		cpuArchitecture != dvResourceVirtualEnvironmentVMCPUArchitecture {
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

	hookScript := d.Get(mkResourceVirtualEnvironmentVMHookScriptFileID).(string)
	if len(hookScript) > 0 {
		createBody.HookScript = &hookScript
	}

	createTimeout := d.Get(mkResourceVirtualEnvironmentVMTimeoutClone).(int)

	err = api.Node(nodeName).VM(0).CreateVM(ctx, createBody, createTimeout)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(vmID))

	return vmCreateCustomDisks(ctx, d, m)
}

func vmCreateCustomDisks(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	vmID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	var commands []string

	// Determine the ID of the next disk.
	disk := d.Get(mkResourceVirtualEnvironmentVMDisk).([]interface{})
	diskCount := 0

	for _, d := range disk {
		block := d.(map[string]interface{})
		fileID, _ := block[mkResourceVirtualEnvironmentVMDiskFileID].(string)

		if fileID == "" {
			diskCount++
		}
	}

	// Retrieve some information about the disk schema.
	resourceSchema := VM().Schema
	diskSchemaElem := resourceSchema[mkResourceVirtualEnvironmentVMDisk].Elem
	diskSchemaResource := diskSchemaElem.(*schema.Resource)
	diskSpeedResource := diskSchemaResource.Schema[mkResourceVirtualEnvironmentVMDiskSpeed]

	// Generate the commands required to import the specified disks.
	importedDiskCount := 0

	for _, d := range disk {
		block := d.(map[string]interface{})

		fileID, _ := block[mkResourceVirtualEnvironmentVMDiskFileID].(string)

		if fileID == "" {
			continue
		}

		datastoreID, _ := block[mkResourceVirtualEnvironmentVMDiskDatastoreID].(string)
		fileFormat, _ := block[mkResourceVirtualEnvironmentVMDiskFileFormat].(string)
		size, _ := block[mkResourceVirtualEnvironmentVMDiskSize].(int)
		speed := block[mkResourceVirtualEnvironmentVMDiskSpeed].([]interface{})
		diskInterface, _ := block[mkResourceVirtualEnvironmentVMDiskInterface].(string)
		ioThread := types.CustomBool(block[mkResourceVirtualEnvironmentVMDiskIOThread].(bool))
		ssd := types.CustomBool(block[mkResourceVirtualEnvironmentVMDiskSSD].(bool))
		discard, _ := block[mkResourceVirtualEnvironmentVMDiskDiscard].(string)
		cache, _ := block[mkResourceVirtualEnvironmentVMDiskCache].(string)

		if fileFormat == "" {
			fileFormat = dvResourceVirtualEnvironmentVMDiskFileFormat
		}

		if len(speed) == 0 {
			diskSpeedDefault, err := diskSpeedResource.DefaultValue()
			if err != nil {
				return diag.FromErr(err)
			}
			speed = diskSpeedDefault.([]interface{})
		}

		speedBlock := speed[0].(map[string]interface{})
		speedLimitRead := speedBlock[mkResourceVirtualEnvironmentVMDiskSpeedRead].(int)
		speedLimitReadBurstable := speedBlock[mkResourceVirtualEnvironmentVMDiskSpeedReadBurstable].(int)
		speedLimitWrite := speedBlock[mkResourceVirtualEnvironmentVMDiskSpeedWrite].(int)
		speedLimitWriteBurstable := speedBlock[mkResourceVirtualEnvironmentVMDiskSpeedWriteBurstable].(int)

		diskOptions := ""

		if ioThread {
			diskOptions += ",iothread=1"
		}

		if ssd {
			diskOptions += ",ssd=1"
		}

		if discard != "" {
			diskOptions += fmt.Sprintf(",discard=%s", discard)
		}

		if cache != "" {
			diskOptions += fmt.Sprintf(",cache=%s", cache)
		}

		if speedLimitRead > 0 {
			diskOptions += fmt.Sprintf(",mbps_rd=%d", speedLimitRead)
		}

		if speedLimitReadBurstable > 0 {
			diskOptions += fmt.Sprintf(",mbps_rd_max=%d", speedLimitReadBurstable)
		}

		if speedLimitWrite > 0 {
			diskOptions += fmt.Sprintf(",mbps_wr=%d", speedLimitWrite)
		}

		if speedLimitWriteBurstable > 0 {
			diskOptions += fmt.Sprintf(",mbps_wr_max=%d", speedLimitWriteBurstable)
		}

		filePathTmp := fmt.Sprintf(
			"/tmp/vm-%d-disk-%d.%s",
			vmID,
			diskCount+importedDiskCount,
			fileFormat,
		)

		//nolint:lll
		commands = append(
			commands,
			`set -e`,
			`try_sudo(){ if [ $(sudo -n echo tfpve 2>&1 | grep "tfpve" | wc -l) -gt 0 ]; then sudo $1; else $1; fi }`,
			fmt.Sprintf(`file_id="%s"`, fileID),
			fmt.Sprintf(`file_format="%s"`, fileFormat),
			fmt.Sprintf(`datastore_id_target="%s"`, datastoreID),
			fmt.Sprintf(`disk_options="%s"`, diskOptions),
			fmt.Sprintf(`disk_size="%d"`, size),
			fmt.Sprintf(`disk_interface="%s"`, diskInterface),
			fmt.Sprintf(`file_path_tmp="%s"`, filePathTmp),
			fmt.Sprintf(`vm_id="%d"`, vmID),
			`source_image=$(try_sudo "pvesm path $file_id")`,
			`imported_disk="$(try_sudo "qm importdisk $vm_id $source_image $datastore_id_target -format $file_format" | grep "unused0" | cut -d ":" -f 3 | cut -d "'" -f 1)"`,
			`disk_id="${datastore_id_target}:$imported_disk${disk_options}"`,
			`try_sudo "qm set $vm_id -${disk_interface} $disk_id"`,
			`try_sudo "qm resize $vm_id ${disk_interface} ${disk_size}G"`,
		)

		importedDiskCount++
	}

	// Execute the commands on the node and wait for the result.
	// This is a highly experimental approach to disk imports and is not recommended by Proxmox.
	if len(commands) > 0 {
		config := m.(proxmoxtf.ProviderConfiguration)

		api, err := config.GetClient()
		if err != nil {
			return diag.FromErr(err)
		}

		nodeName := d.Get(mkResourceVirtualEnvironmentVMNodeName).(string)

		out, err := api.SSH().ExecuteNodeCommands(ctx, nodeName, commands)
		if err != nil {
			if strings.Contains(err.Error(), "pvesm: not found") {
				return diag.Errorf("The configured SSH user '%s' does not have the required permissions to import disks. "+
					"Make sure `sudo` is installated and the user is a member of sudoers.", api.SSH().Username())
			}

			return diag.FromErr(err)
		}

		tflog.Debug(ctx, "vmCreateCustomDisks", map[string]interface{}{
			"output": string(out),
		})
	}

	return vmCreateStart(ctx, d, m)
}

func vmCreateStart(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	started := d.Get(mkResourceVirtualEnvironmentVMStarted).(bool)
	template := d.Get(mkResourceVirtualEnvironmentVMTemplate).(bool)
	reboot := d.Get(mkResourceVirtualEnvironmentVMRebootAfterCreation).(bool)

	if !started || template {
		return vmRead(ctx, d, m)
	}

	config := m.(proxmoxtf.ProviderConfiguration)
	api, err := config.GetClient()
	if err != nil {
		return diag.FromErr(err)
	}

	nodeName := d.Get(mkResourceVirtualEnvironmentVMNodeName).(string)
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
		rebootTimeout := d.Get(mkResourceVirtualEnvironmentVMTimeoutReboot).(int)

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
	devices := d.Get(mkResourceVirtualEnvironmentVMAudioDevice).([]interface{})
	list := make(vms.CustomAudioDevices, len(devices))

	for i, v := range devices {
		block := v.(map[string]interface{})

		device, _ := block[mkResourceVirtualEnvironmentVMAudioDeviceDevice].(string)
		driver, _ := block[mkResourceVirtualEnvironmentVMAudioDeviceDriver].(string)
		enabled, _ := block[mkResourceVirtualEnvironmentVMAudioDeviceEnabled].(bool)

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

//nolint:lll
func vmGetCloudInitConfig(d *schema.ResourceData) *vms.CustomCloudInitConfig {
	initialization := d.Get(mkResourceVirtualEnvironmentVMInitialization).([]interface{})

	if len(initialization) == 0 || initialization[0] == nil {
		return nil
	}

	var initializationConfig *vms.CustomCloudInitConfig

	initializationBlock := initialization[0].(map[string]interface{})
	initializationConfig = &vms.CustomCloudInitConfig{}
	initializationDNS := initializationBlock[mkResourceVirtualEnvironmentVMInitializationDNS].([]interface{})

	if len(initializationDNS) > 0 && initializationDNS[0] != nil {
		initializationDNSBlock := initializationDNS[0].(map[string]interface{})
		domain := initializationDNSBlock[mkResourceVirtualEnvironmentVMInitializationDNSDomain].(string)

		if domain != "" {
			initializationConfig.SearchDomain = &domain
		}

		servers := initializationDNSBlock[mkResourceVirtualEnvironmentVMInitializationDNSServers].([]interface{})
		deprecatedServer := initializationDNSBlock[mkResourceVirtualEnvironmentVMInitializationDNSServer].(string)

		if len(servers) > 0 {
			nameserver := strings.Join(utils.ConvertToStringSlice(servers), " ")

			initializationConfig.Nameserver = &nameserver
		} else if deprecatedServer != "" {
			initializationConfig.Nameserver = &deprecatedServer
		}
	}

	initializationIPConfig := initializationBlock[mkResourceVirtualEnvironmentVMInitializationIPConfig].([]interface{})
	initializationConfig.IPConfig = make(
		[]vms.CustomCloudInitIPConfig,
		len(initializationIPConfig),
	)

	for i, c := range initializationIPConfig {
		configBlock := c.(map[string]interface{})
		ipv4 := configBlock[mkResourceVirtualEnvironmentVMInitializationIPConfigIPv4].([]interface{})

		if len(ipv4) > 0 {
			ipv4Block := ipv4[0].(map[string]interface{})
			ipv4Address := ipv4Block[mkResourceVirtualEnvironmentVMInitializationIPConfigIPv4Address].(string)

			if ipv4Address != "" {
				initializationConfig.IPConfig[i].IPv4 = &ipv4Address
			}

			ipv4Gateway := ipv4Block[mkResourceVirtualEnvironmentVMInitializationIPConfigIPv4Gateway].(string)

			if ipv4Gateway != "" {
				initializationConfig.IPConfig[i].GatewayIPv4 = &ipv4Gateway
			}
		}

		ipv6 := configBlock[mkResourceVirtualEnvironmentVMInitializationIPConfigIPv6].([]interface{})

		if len(ipv6) > 0 {
			ipv6Block := ipv6[0].(map[string]interface{})
			ipv6Address := ipv6Block[mkResourceVirtualEnvironmentVMInitializationIPConfigIPv6Address].(string)

			if ipv6Address != "" {
				initializationConfig.IPConfig[i].IPv6 = &ipv6Address
			}

			ipv6Gateway := ipv6Block[mkResourceVirtualEnvironmentVMInitializationIPConfigIPv6Gateway].(string)

			if ipv6Gateway != "" {
				initializationConfig.IPConfig[i].GatewayIPv6 = &ipv6Gateway
			}
		}
	}

	initializationUserAccount := initializationBlock[mkResourceVirtualEnvironmentVMInitializationUserAccount].([]interface{})

	if len(initializationUserAccount) > 0 {
		initializationUserAccountBlock := initializationUserAccount[0].(map[string]interface{})
		keys := initializationUserAccountBlock[mkResourceVirtualEnvironmentVMInitializationUserAccountKeys].([]interface{})

		if len(keys) > 0 {
			sshKeys := make(vms.CustomCloudInitSSHKeys, len(keys))

			for i, k := range keys {
				sshKeys[i] = k.(string)
			}

			initializationConfig.SSHKeys = &sshKeys
		}

		password := initializationUserAccountBlock[mkResourceVirtualEnvironmentVMInitializationUserAccountPassword].(string)

		if password != "" {
			initializationConfig.Password = &password
		}

		username := initializationUserAccountBlock[mkResourceVirtualEnvironmentVMInitializationUserAccountUsername].(string)

		initializationConfig.Username = &username
	}

	initializationUserDataFileID := initializationBlock[mkResourceVirtualEnvironmentVMInitializationUserDataFileID].(string)

	if initializationUserDataFileID != "" {
		initializationConfig.Files = &vms.CustomCloudInitFiles{
			UserVolume: &initializationUserDataFileID,
		}
	}

	initializationVendorDataFileID := initializationBlock[mkResourceVirtualEnvironmentVMInitializationVendorDataFileID].(string)

	if initializationVendorDataFileID != "" {
		if initializationConfig.Files == nil {
			initializationConfig.Files = &vms.CustomCloudInitFiles{}
		}

		initializationConfig.Files.VendorVolume = &initializationVendorDataFileID
	}

	initializationNetworkDataFileID := initializationBlock[mkResourceVirtualEnvironmentVMInitializationNetworkDataFileID].(string)

	if initializationNetworkDataFileID != "" {
		if initializationConfig.Files == nil {
			initializationConfig.Files = &vms.CustomCloudInitFiles{}
		}

		initializationConfig.Files.NetworkVolume = &initializationNetworkDataFileID
	}

	//nolint:lll
	initializationMetaDataFileID := initializationBlock[mkResourceVirtualEnvironmentVMInitializationMetaDataFileID].(string)

	if initializationMetaDataFileID != "" {
		if initializationConfig.Files == nil {
			initializationConfig.Files = &vms.CustomCloudInitFiles{}
		}

		initializationConfig.Files.MetaVolume = &initializationMetaDataFileID
	}

	initializationType := initializationBlock[mkResourceVirtualEnvironmentVMInitializationType].(string)

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

func vmGetDiskDeviceObjects(
	d *schema.ResourceData,
	disks []interface{},
) (map[string]map[string]vms.CustomStorageDevice, error) {
	var diskDevice []interface{}

	if disks != nil {
		diskDevice = disks
	} else {
		diskDevice = d.Get(mkResourceVirtualEnvironmentVMDisk).([]interface{})
	}

	diskDeviceObjects := map[string]map[string]vms.CustomStorageDevice{}
	resource := VM()

	for _, diskEntry := range diskDevice {
		diskDevice := vms.CustomStorageDevice{
			Enabled: true,
		}

		block := diskEntry.(map[string]interface{})
		datastoreID, _ := block[mkResourceVirtualEnvironmentVMDiskDatastoreID].(string)
		pathInDatastore := ""

		if untyped, hasPathInDatastore := block[mkResourceVirtualEnvironmentVMDiskPathInDatastore]; hasPathInDatastore {
			pathInDatastore = untyped.(string)
		}

		fileFormat, _ := block[mkResourceVirtualEnvironmentVMDiskFileFormat].(string)
		fileID, _ := block[mkResourceVirtualEnvironmentVMDiskFileID].(string)
		size, _ := block[mkResourceVirtualEnvironmentVMDiskSize].(int)
		diskInterface, _ := block[mkResourceVirtualEnvironmentVMDiskInterface].(string)
		ioThread := types.CustomBool(block[mkResourceVirtualEnvironmentVMDiskIOThread].(bool))
		ssd := types.CustomBool(block[mkResourceVirtualEnvironmentVMDiskSSD].(bool))
		discard := block[mkResourceVirtualEnvironmentVMDiskDiscard].(string)
		cache := block[mkResourceVirtualEnvironmentVMDiskCache].(string)

		speedBlock, err := structure.GetSchemaBlock(
			resource,
			d,
			[]string{mkResourceVirtualEnvironmentVMDisk, mkResourceVirtualEnvironmentVMDiskSpeed},
			0,
			false,
		)
		if err != nil {
			return diskDeviceObjects, err
		}

		if fileFormat == "" {
			fileFormat = dvResourceVirtualEnvironmentVMDiskFileFormat
		}
		if fileID != "" {
			diskDevice.Enabled = false
		}

		if pathInDatastore != "" {
			if datastoreID != "" {
				diskDevice.FileVolume = fmt.Sprintf("%s:%s", datastoreID, pathInDatastore)
			} else {
				// FileVolume is absolute path in the host filesystem
				diskDevice.FileVolume = pathInDatastore
			}
		} else {
			diskDevice.FileVolume = fmt.Sprintf("%s:%d", datastoreID, size)
		}

		diskDevice.ID = &datastoreID
		diskDevice.Interface = &diskInterface
		diskDevice.Format = &fileFormat
		diskDevice.FileID = &fileID
		diskSize := types.DiskSizeFromGigabytes(int64(size))
		diskDevice.Size = &diskSize
		diskDevice.IOThread = &ioThread
		diskDevice.Discard = &discard
		diskDevice.Cache = &cache

		if !strings.HasPrefix(diskInterface, "virtio") {
			diskDevice.SSD = &ssd
		}

		if len(speedBlock) > 0 {
			speedLimitRead := speedBlock[mkResourceVirtualEnvironmentVMDiskSpeedRead].(int)
			speedLimitReadBurstable := speedBlock[mkResourceVirtualEnvironmentVMDiskSpeedReadBurstable].(int)
			speedLimitWrite := speedBlock[mkResourceVirtualEnvironmentVMDiskSpeedWrite].(int)
			speedLimitWriteBurstable := speedBlock[mkResourceVirtualEnvironmentVMDiskSpeedWriteBurstable].(int)

			if speedLimitRead > 0 {
				diskDevice.MaxReadSpeedMbps = &speedLimitRead
			}

			if speedLimitReadBurstable > 0 {
				diskDevice.BurstableReadSpeedMbps = &speedLimitReadBurstable
			}

			if speedLimitWrite > 0 {
				diskDevice.MaxWriteSpeedMbps = &speedLimitWrite
			}

			if speedLimitWriteBurstable > 0 {
				diskDevice.BurstableWriteSpeedMbps = &speedLimitWriteBurstable
			}
		}

		baseDiskInterface := diskDigitPrefix(diskInterface)

		if baseDiskInterface != "virtio" && baseDiskInterface != "scsi" &&
			baseDiskInterface != "sata" {
			errorMsg := fmt.Sprintf(
				"Defined disk interface not supported. Interface was %s, but only virtio, sata and scsi are supported",
				diskInterface,
			)
			return diskDeviceObjects, errors.New(errorMsg)
		}

		if _, present := diskDeviceObjects[baseDiskInterface]; !present {
			diskDeviceObjects[baseDiskInterface] = map[string]vms.CustomStorageDevice{}
		}

		diskDeviceObjects[baseDiskInterface][diskInterface] = diskDevice
	}

	return diskDeviceObjects, nil
}

func vmGetEfiDisk(d *schema.ResourceData, disk []interface{}) *vms.CustomEFIDisk {
	var efiDisk []interface{}

	if disk != nil {
		efiDisk = disk
	} else {
		efiDisk = d.Get(mkResourceVirtualEnvironmentVMEFIDisk).([]interface{})
	}

	var efiDiskConfig *vms.CustomEFIDisk

	if len(efiDisk) > 0 {
		efiDiskConfig = &vms.CustomEFIDisk{}

		block := efiDisk[0].(map[string]interface{})
		datastoreID, _ := block[mkResourceVirtualEnvironmentVMEFIDiskDatastoreID].(string)
		fileFormat, _ := block[mkResourceVirtualEnvironmentVMEFIDiskFileFormat].(string)
		efiType, _ := block[mkResourceVirtualEnvironmentVMEFIDiskType].(string)
		preEnrolledKeys := types.CustomBool(block[mkResourceVirtualEnvironmentVMEFIDiskPreEnrolledKeys].(bool))

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
			Enabled:    true,
			FileVolume: efiDisk.FileVolume,
			Format:     efiDisk.Format,
			Interface:  &diskInterface,
			ID:         &id,
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
		tpmState = d.Get(mkResourceVirtualEnvironmentVMTPMState).([]interface{})
	}

	var tpmStateConfig *vms.CustomTPMState

	if len(tpmState) > 0 {
		tpmStateConfig = &vms.CustomTPMState{}

		block := tpmState[0].(map[string]interface{})
		datastoreID, _ := block[mkResourceVirtualEnvironmentVMTPMStateDatastoreID].(string)
		version, _ := block[mkResourceVirtualEnvironmentVMTPMStateVersion].(string)

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
			Enabled:    true,
			FileVolume: tpmState.FileVolume,
			Interface:  &diskInterface,
			ID:         &id,
		}
	}

	return storageDevice
}

func vmGetHostPCIDeviceObjects(d *schema.ResourceData) vms.CustomPCIDevices {
	pciDevice := d.Get(mkResourceVirtualEnvironmentVMHostPCI).([]interface{})
	pciDeviceObjects := make(vms.CustomPCIDevices, len(pciDevice))

	for i, pciDeviceEntry := range pciDevice {
		block := pciDeviceEntry.(map[string]interface{})

		ids, _ := block[mkResourceVirtualEnvironmentVMHostPCIDeviceID].(string)
		mdev, _ := block[mkResourceVirtualEnvironmentVMHostPCIDeviceMDev].(string)
		pcie := types.CustomBool(block[mkResourceVirtualEnvironmentVMHostPCIDevicePCIE].(bool))
		rombar := types.CustomBool(
			block[mkResourceVirtualEnvironmentVMHostPCIDeviceROMBAR].(bool),
		)
		romfile, _ := block[mkResourceVirtualEnvironmentVMHostPCIDeviceROMFile].(string)
		xvga := types.CustomBool(block[mkResourceVirtualEnvironmentVMHostPCIDeviceXVGA].(bool))
		mapping, _ := block[mkResourceVirtualEnvironmentVMHostPCIDeviceMapping].(string)

		device := vms.CustomPCIDevice{
			PCIExpress: &pcie,
			ROMBAR:     &rombar,
			XVGA:       &xvga,
		}
		if ids != "" {
			dIds := strings.Split(ids, ";")
			device.DeviceIDs = &dIds
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
	usbDevice := d.Get(mkResourceVirtualEnvironmentVMHostUSB).([]interface{})
	usbDeviceObjects := make(vms.CustomUSBDevices, len(usbDevice))

	for i, usbDeviceEntry := range usbDevice {
		block := usbDeviceEntry.(map[string]interface{})

		host, _ := block[mkResourceVirtualEnvironmentVMHostUSBDevice].(string)
		usb3 := types.CustomBool(block[mkResourceVirtualEnvironmentVMHostUSBDeviceUSB3].(bool))
		mapping, _ := block[mkResourceVirtualEnvironmentVMHostUSBDeviceMapping].(string)

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
	networkDevice := d.Get(mkResourceVirtualEnvironmentVMNetworkDevice).([]interface{})
	networkDeviceObjects := make(vms.CustomNetworkDevices, len(networkDevice))

	for i, networkDeviceEntry := range networkDevice {
		block := networkDeviceEntry.(map[string]interface{})

		bridge := block[mkResourceVirtualEnvironmentVMNetworkDeviceBridge].(string)
		enabled := block[mkResourceVirtualEnvironmentVMNetworkDeviceEnabled].(bool)
		firewall := types.CustomBool(block[mkResourceVirtualEnvironmentVMNetworkDeviceFirewall].(bool))
		macAddress := block[mkResourceVirtualEnvironmentVMNetworkDeviceMACAddress].(string)
		model := block[mkResourceVirtualEnvironmentVMNetworkDeviceModel].(string)
		queues := block[mkResourceVirtualEnvironmentVMNetworkDeviceQueues].(int)
		rateLimit := block[mkResourceVirtualEnvironmentVMNetworkDeviceRateLimit].(float64)
		vlanID := block[mkResourceVirtualEnvironmentVMNetworkDeviceVLANID].(int)
		mtu := block[mkResourceVirtualEnvironmentVMNetworkDeviceMTU].(int)

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
	device := d.Get(mkResourceVirtualEnvironmentVMSerialDevice).([]interface{})
	list := make(vms.CustomSerialDevices, len(device))

	for i, v := range device {
		block := v.(map[string]interface{})

		device, _ := block[mkResourceVirtualEnvironmentVMSerialDeviceDevice].(string)

		list[i] = device
	}

	return list
}

func vmGetSMBIOS(d *schema.ResourceData) *vms.CustomSMBIOS {
	smbiosSections := d.Get(mkResourceVirtualEnvironmentVMSMBIOS).([]interface{})
	//nolint:nestif
	if len(smbiosSections) > 0 {
		smbiosBlock := smbiosSections[0].(map[string]interface{})
		b64 := types.CustomBool(true)
		family, _ := smbiosBlock[mkResourceVirtualEnvironmentVMSMBIOSFamily].(string)
		manufacturer, _ := smbiosBlock[mkResourceVirtualEnvironmentVMSMBIOSManufacturer].(string)
		product, _ := smbiosBlock[mkResourceVirtualEnvironmentVMSMBIOSProduct].(string)
		serial, _ := smbiosBlock[mkResourceVirtualEnvironmentVMSMBIOSSerial].(string)
		sku, _ := smbiosBlock[mkResourceVirtualEnvironmentVMSMBIOSSKU].(string)
		version, _ := smbiosBlock[mkResourceVirtualEnvironmentVMSMBIOSVersion].(string)
		uid, _ := smbiosBlock[mkResourceVirtualEnvironmentVMSMBIOSUUID].(string)

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
	startup := d.Get(mkResourceVirtualEnvironmentVMStartup).([]interface{})
	if len(startup) > 0 {
		startupBlock := startup[0].(map[string]interface{})
		startupOrder := startupBlock[mkResourceVirtualEnvironmentVMStartupOrder].(int)
		startupUpDelay := startupBlock[mkResourceVirtualEnvironmentVMStartupUpDelay].(int)
		startupDownDelay := startupBlock[mkResourceVirtualEnvironmentVMStartupDownDelay].(int)

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

	tags := d.Get(mkResourceVirtualEnvironmentVMTags).([]interface{})
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
		[]string{mkResourceVirtualEnvironmentVMVGA},
		0,
		true,
	)
	if err != nil {
		return nil, err
	}

	vgaEnabled := types.CustomBool(vgaBlock[mkResourceVirtualEnvironmentVMVGAEnabled].(bool))
	vgaMemory := vgaBlock[mkResourceVirtualEnvironmentVMVGAMemory].(int)
	vgaType := vgaBlock[mkResourceVirtualEnvironmentVMVGAType].(string)

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

	if vmNodeName != d.Get(mkResourceVirtualEnvironmentVMNodeName) {
		err = d.Set(mkResourceVirtualEnvironmentVMNodeName, vmNodeName)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	nodeName := d.Get(mkResourceVirtualEnvironmentVMNodeName).(string)

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
	if storedVMID := d.Get(mkResourceVirtualEnvironmentVMVMID).(int); storedVMID == -1 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary: fmt.Sprintf("VM %s has stored legacy vm_id %d, setting vm_id to its correct value %d.",
				d.Id(), storedVMID, vmID),
		})

		err = d.Set(mkResourceVirtualEnvironmentVMVMID, vmID)
		diags = append(diags, diag.FromErr(err)...)
	}

	nodeName := d.Get(mkResourceVirtualEnvironmentVMNodeName).(string)
	clone := d.Get(mkResourceVirtualEnvironmentVMClone).([]interface{})

	err = setAgent(d, len(clone) > 0, vmConfig)
	diags = append(diags, diag.FromErr(err)...)

	// Compare the audio devices to those stored in the state.
	currentAudioDevice := d.Get(mkResourceVirtualEnvironmentVMAudioDevice).([]interface{})

	audioDevices := make([]interface{}, 1)
	audioDevicesArray := []*vms.CustomAudioDevice{
		vmConfig.AudioDevice,
	}
	audioDevicesCount := 0

	for adi, ad := range audioDevicesArray {
		m := map[string]interface{}{}

		if ad != nil {
			m[mkResourceVirtualEnvironmentVMAudioDeviceDevice] = ad.Device

			if ad.Driver != nil {
				m[mkResourceVirtualEnvironmentVMAudioDeviceDriver] = *ad.Driver
			} else {
				m[mkResourceVirtualEnvironmentVMAudioDeviceDriver] = ""
			}

			m[mkResourceVirtualEnvironmentVMAudioDeviceEnabled] = true

			audioDevicesCount = adi + 1
		} else {
			m[mkResourceVirtualEnvironmentVMAudioDeviceDevice] = ""
			m[mkResourceVirtualEnvironmentVMAudioDeviceDriver] = ""
			m[mkResourceVirtualEnvironmentVMAudioDeviceEnabled] = false
		}

		audioDevices[adi] = m
	}

	if len(clone) == 0 || len(currentAudioDevice) > 0 {
		err := d.Set(mkResourceVirtualEnvironmentVMAudioDevice, audioDevices[:audioDevicesCount])
		diags = append(diags, diag.FromErr(err)...)
	}

	// Compare the IDE devices to the CD-ROM configurations stored in the state.
	currentInterface := dvResourceVirtualEnvironmentVMCDROMInterface

	currentCDROM := d.Get(mkResourceVirtualEnvironmentVMCDROM).([]interface{})
	if len(currentCDROM) > 0 {
		currentBlock := currentCDROM[0].(map[string]interface{})
		currentInterface = currentBlock[mkResourceVirtualEnvironmentVMCDROMInterface].(string)
	}

	cdromIDEDevice := getIdeDevice(vmConfig, currentInterface)

	//nolint:nestif
	if cdromIDEDevice != nil {
		cdrom := make([]interface{}, 1)
		cdromBlock := map[string]interface{}{}

		if len(clone) == 0 || len(currentCDROM) > 0 {
			cdromBlock[mkResourceVirtualEnvironmentVMCDROMEnabled] = cdromIDEDevice.Enabled
			cdromBlock[mkResourceVirtualEnvironmentVMCDROMFileID] = cdromIDEDevice.FileVolume
			cdromBlock[mkResourceVirtualEnvironmentVMCDROMInterface] = currentInterface

			if len(currentCDROM) > 0 {
				currentBlock := currentCDROM[0].(map[string]interface{})

				if currentBlock[mkResourceVirtualEnvironmentVMCDROMFileID] == "" {
					cdromBlock[mkResourceVirtualEnvironmentVMCDROMFileID] = ""
				}

				if currentBlock[mkResourceVirtualEnvironmentVMCDROMEnabled] == false {
					cdromBlock[mkResourceVirtualEnvironmentVMCDROMEnabled] = false
				}
			}

			cdrom[0] = cdromBlock

			err := d.Set(mkResourceVirtualEnvironmentVMCDROM, cdrom)
			diags = append(diags, diag.FromErr(err)...)
		}

	} else {
		err := d.Set(mkResourceVirtualEnvironmentVMCDROM, []interface{}{})
		diags = append(diags, diag.FromErr(err)...)
	}

	// Compare the CPU configuration to the one stored in the state.
	cpu := map[string]interface{}{}

	if vmConfig.CPUArchitecture != nil {
		cpu[mkResourceVirtualEnvironmentVMCPUArchitecture] = *vmConfig.CPUArchitecture
	} else {
		// Default value of "arch" is "" according to the API documentation.
		// However, assume the provider's default value as a workaround when the root account is not being used.
		if !api.API().IsRootTicket() {
			cpu[mkResourceVirtualEnvironmentVMCPUArchitecture] = dvResourceVirtualEnvironmentVMCPUArchitecture
		} else {
			cpu[mkResourceVirtualEnvironmentVMCPUArchitecture] = ""
		}
	}

	if vmConfig.CPUCores != nil {
		cpu[mkResourceVirtualEnvironmentVMCPUCores] = *vmConfig.CPUCores
	} else {
		// Default value of "cores" is "1" according to the API documentation.
		cpu[mkResourceVirtualEnvironmentVMCPUCores] = 1
	}

	if vmConfig.VirtualCPUCount != nil {
		cpu[mkResourceVirtualEnvironmentVMCPUHotplugged] = *vmConfig.VirtualCPUCount
	} else {
		// Default value of "vcpus" is "1" according to the API documentation.
		cpu[mkResourceVirtualEnvironmentVMCPUHotplugged] = 0
	}

	if vmConfig.CPULimit != nil {
		cpu[mkResourceVirtualEnvironmentVMCPULimit] = *vmConfig.CPULimit
	} else {
		// Default value of "cpulimit" is "0" according to the API documentation.
		cpu[mkResourceVirtualEnvironmentVMCPULimit] = 0
	}

	if vmConfig.NUMAEnabled != nil {
		cpu[mkResourceVirtualEnvironmentVMCPUNUMA] = *vmConfig.NUMAEnabled
	} else {
		// Default value of "numa" is "false" according to the API documentation.
		cpu[mkResourceVirtualEnvironmentVMCPUNUMA] = false
	}

	if vmConfig.CPUSockets != nil {
		cpu[mkResourceVirtualEnvironmentVMCPUSockets] = *vmConfig.CPUSockets
	} else {
		// Default value of "sockets" is "1" according to the API documentation.
		cpu[mkResourceVirtualEnvironmentVMCPUSockets] = 1
	}

	if vmConfig.CPUEmulation != nil {
		if vmConfig.CPUEmulation.Flags != nil {
			convertedFlags := make([]interface{}, len(*vmConfig.CPUEmulation.Flags))

			for fi, fv := range *vmConfig.CPUEmulation.Flags {
				convertedFlags[fi] = fv
			}

			cpu[mkResourceVirtualEnvironmentVMCPUFlags] = convertedFlags
		} else {
			cpu[mkResourceVirtualEnvironmentVMCPUFlags] = []interface{}{}
		}

		cpu[mkResourceVirtualEnvironmentVMCPUType] = vmConfig.CPUEmulation.Type
	} else {
		cpu[mkResourceVirtualEnvironmentVMCPUFlags] = []interface{}{}
		// Default value of "cputype" is "qemu64" according to the QEMU documentation.
		cpu[mkResourceVirtualEnvironmentVMCPUType] = "qemu64"
	}

	if vmConfig.CPUUnits != nil {
		cpu[mkResourceVirtualEnvironmentVMCPUUnits] = *vmConfig.CPUUnits
	} else {
		// Default value of "cpuunits" is "1024" according to the API documentation.
		cpu[mkResourceVirtualEnvironmentVMCPUUnits] = 1024
	}

	currentCPU := d.Get(mkResourceVirtualEnvironmentVMCPU).([]interface{})

	if len(clone) > 0 {
		if len(currentCPU) > 0 {
			err := d.Set(mkResourceVirtualEnvironmentVMCPU, []interface{}{cpu})
			diags = append(diags, diag.FromErr(err)...)
		}
	} else if len(currentCPU) > 0 ||
		cpu[mkResourceVirtualEnvironmentVMCPUArchitecture] != dvResourceVirtualEnvironmentVMCPUArchitecture ||
		cpu[mkResourceVirtualEnvironmentVMCPUCores] != dvResourceVirtualEnvironmentVMCPUCores ||
		len(cpu[mkResourceVirtualEnvironmentVMCPUFlags].([]interface{})) > 0 ||
		cpu[mkResourceVirtualEnvironmentVMCPUHotplugged] != dvResourceVirtualEnvironmentVMCPUHotplugged ||
		cpu[mkResourceVirtualEnvironmentVMCPULimit] != dvResourceVirtualEnvironmentVMCPULimit ||
		cpu[mkResourceVirtualEnvironmentVMCPUSockets] != dvResourceVirtualEnvironmentVMCPUSockets ||
		cpu[mkResourceVirtualEnvironmentVMCPUType] != dvResourceVirtualEnvironmentVMCPUType ||
		cpu[mkResourceVirtualEnvironmentVMCPUUnits] != dvResourceVirtualEnvironmentVMCPUUnits {
		err := d.Set(mkResourceVirtualEnvironmentVMCPU, []interface{}{cpu})
		diags = append(diags, diag.FromErr(err)...)
	}

	currentDiskList := d.Get(mkResourceVirtualEnvironmentVMDisk).([]interface{})
	diskMap := map[string]interface{}{}
	diskObjects := getDiskInfo(vmConfig, d)

	for di, dd := range diskObjects {
		if dd == nil || dd.FileVolume == "none" || strings.HasPrefix(di, "ide") {
			continue
		}

		if dd.IsCloudInitDrive(vmID) {
			continue
		}

		disk := map[string]interface{}{}

		datastoreID, pathInDatastore, hasDatastoreID := strings.Cut(dd.FileVolume, ":")
		if !hasDatastoreID {
			// when no ':' separator is found, 'Cut' places the whole string to 'datastoreID',
			// we want it in 'pathInDatastore' (it is absolute filesystem path)
			pathInDatastore = datastoreID
			datastoreID = ""
		}

		disk[mkResourceVirtualEnvironmentVMDiskDatastoreID] = datastoreID
		disk[mkResourceVirtualEnvironmentVMDiskPathInDatastore] = pathInDatastore

		if dd.Format == nil {
			disk[mkResourceVirtualEnvironmentVMDiskFileFormat] = dvResourceVirtualEnvironmentVMDiskFileFormat

			if datastoreID != "" {
				// disk format may not be returned by config API if it is default for the storage, and that may be different
				// from the default qcow2, so we need to read it from the storage API to make sure we have the correct value
				volume, err := api.Node(nodeName).Storage(datastoreID).GetDatastoreFile(ctx, dd.FileVolume)
				if err != nil {
					diags = append(diags, diag.FromErr(err)...)
					continue
				}

				disk[mkResourceVirtualEnvironmentVMDiskFileFormat] = volume.FileFormat
			}
		} else {
			disk[mkResourceVirtualEnvironmentVMDiskFileFormat] = dd.Format
		}

		if dd.FileID != nil {
			disk[mkResourceVirtualEnvironmentVMDiskFileID] = dd.FileID
		}

		disk[mkResourceVirtualEnvironmentVMDiskInterface] = di
		disk[mkResourceVirtualEnvironmentVMDiskSize] = dd.Size.InGigabytes()

		if dd.BurstableReadSpeedMbps != nil ||
			dd.BurstableWriteSpeedMbps != nil ||
			dd.MaxReadSpeedMbps != nil ||
			dd.MaxWriteSpeedMbps != nil {
			speed := map[string]interface{}{}

			if dd.MaxReadSpeedMbps != nil {
				speed[mkResourceVirtualEnvironmentVMDiskSpeedRead] = *dd.MaxReadSpeedMbps
			} else {
				speed[mkResourceVirtualEnvironmentVMDiskSpeedRead] = 0
			}

			if dd.BurstableReadSpeedMbps != nil {
				speed[mkResourceVirtualEnvironmentVMDiskSpeedReadBurstable] = *dd.BurstableReadSpeedMbps
			} else {
				speed[mkResourceVirtualEnvironmentVMDiskSpeedReadBurstable] = 0
			}

			if dd.MaxWriteSpeedMbps != nil {
				speed[mkResourceVirtualEnvironmentVMDiskSpeedWrite] = *dd.MaxWriteSpeedMbps
			} else {
				speed[mkResourceVirtualEnvironmentVMDiskSpeedWrite] = 0
			}

			if dd.BurstableWriteSpeedMbps != nil {
				speed[mkResourceVirtualEnvironmentVMDiskSpeedWriteBurstable] = *dd.BurstableWriteSpeedMbps
			} else {
				speed[mkResourceVirtualEnvironmentVMDiskSpeedWriteBurstable] = 0
			}

			disk[mkResourceVirtualEnvironmentVMDiskSpeed] = []interface{}{speed}
		} else {
			disk[mkResourceVirtualEnvironmentVMDiskSpeed] = []interface{}{}
		}

		if dd.IOThread != nil {
			disk[mkResourceVirtualEnvironmentVMDiskIOThread] = *dd.IOThread
		} else {
			disk[mkResourceVirtualEnvironmentVMDiskIOThread] = false
		}

		if dd.SSD != nil {
			disk[mkResourceVirtualEnvironmentVMDiskSSD] = *dd.SSD
		} else {
			disk[mkResourceVirtualEnvironmentVMDiskSSD] = false
		}

		if dd.Discard != nil {
			disk[mkResourceVirtualEnvironmentVMDiskDiscard] = *dd.Discard
		} else {
			disk[mkResourceVirtualEnvironmentVMDiskDiscard] = dvResourceVirtualEnvironmentVMDiskDiscard
		}

		if dd.Cache != nil {
			disk[mkResourceVirtualEnvironmentVMDiskCache] = *dd.Cache
		} else {
			disk[mkResourceVirtualEnvironmentVMDiskCache] = dvResourceVirtualEnvironmentVMDiskCache
		}

		diskMap[di] = disk
	}

	if len(clone) == 0 || len(currentDiskList) > 0 {
		orderedDiskList := orderedListFromMap(diskMap)
		err := d.Set(mkResourceVirtualEnvironmentVMDisk, orderedDiskList)
		diags = append(diags, diag.FromErr(err)...)
	}

	//nolint:nestif
	if vmConfig.EFIDisk != nil {
		efiDisk := map[string]interface{}{}

		fileIDParts := strings.Split(vmConfig.EFIDisk.FileVolume, ":")

		efiDisk[mkResourceVirtualEnvironmentVMEFIDiskDatastoreID] = fileIDParts[0]

		if vmConfig.EFIDisk.Format != nil {
			efiDisk[mkResourceVirtualEnvironmentVMEFIDiskFileFormat] = *vmConfig.EFIDisk.Format
		} else {
			// disk format may not be returned by config API if it is default for the storage, and that may be different
			// from the default qcow2, so we need to read it from the storage API to make sure we have the correct value
			volume, err := api.Node(nodeName).Storage(fileIDParts[0]).GetDatastoreFile(ctx, vmConfig.EFIDisk.FileVolume)
			if err != nil {
				diags = append(diags, diag.FromErr(err)...)
			} else {
				efiDisk[mkResourceVirtualEnvironmentVMEFIDiskFileFormat] = volume.FileFormat
			}
		}

		if vmConfig.EFIDisk.Type != nil {
			efiDisk[mkResourceVirtualEnvironmentVMEFIDiskType] = *vmConfig.EFIDisk.Type
		} else {
			efiDisk[mkResourceVirtualEnvironmentVMEFIDiskType] = dvResourceVirtualEnvironmentVMEFIDiskType
		}

		if vmConfig.EFIDisk.PreEnrolledKeys != nil {
			efiDisk[mkResourceVirtualEnvironmentVMEFIDiskPreEnrolledKeys] = *vmConfig.EFIDisk.PreEnrolledKeys
		} else {
			efiDisk[mkResourceVirtualEnvironmentVMEFIDiskPreEnrolledKeys] = false
		}

		currentEfiDisk := d.Get(mkResourceVirtualEnvironmentVMEFIDisk).([]interface{})

		if len(clone) > 0 {
			if len(currentEfiDisk) > 0 {
				err := d.Set(mkResourceVirtualEnvironmentVMEFIDisk, []interface{}{efiDisk})
				diags = append(diags, diag.FromErr(err)...)
			}
		} else if len(currentEfiDisk) > 0 ||
			efiDisk[mkResourceVirtualEnvironmentVMEFIDiskDatastoreID] != dvResourceVirtualEnvironmentVMEFIDiskDatastoreID ||
			efiDisk[mkResourceVirtualEnvironmentVMEFIDiskType] != dvResourceVirtualEnvironmentVMEFIDiskType ||
			efiDisk[mkResourceVirtualEnvironmentVMEFIDiskPreEnrolledKeys] != dvResourceVirtualEnvironmentVMEFIDiskPreEnrolledKeys || //nolint:lll
			efiDisk[mkResourceVirtualEnvironmentVMEFIDiskFileFormat] != dvResourceVirtualEnvironmentVMEFIDiskFileFormat {
			err := d.Set(mkResourceVirtualEnvironmentVMEFIDisk, []interface{}{efiDisk})
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	if vmConfig.TPMState != nil {
		tpmState := map[string]interface{}{}

		fileIDParts := strings.Split(vmConfig.TPMState.FileVolume, ":")

		tpmState[mkResourceVirtualEnvironmentVMTPMStateDatastoreID] = fileIDParts[0]
		tpmState[mkResourceVirtualEnvironmentVMTPMStateVersion] = dvResourceVirtualEnvironmentVMTPMStateVersion

		currentTPMState := d.Get(mkResourceVirtualEnvironmentVMTPMState).([]interface{})

		if len(clone) > 0 {
			if len(currentTPMState) > 0 {
				err := d.Set(mkResourceVirtualEnvironmentVMTPMState, []interface{}{tpmState})
				diags = append(diags, diag.FromErr(err)...)
			}
		} else if len(currentTPMState) > 0 ||
			tpmState[mkResourceVirtualEnvironmentVMTPMStateDatastoreID] != dvResourceVirtualEnvironmentVMTPMStateDatastoreID ||
			tpmState[mkResourceVirtualEnvironmentVMTPMStateVersion] != dvResourceVirtualEnvironmentVMTPMStateVersion {
			err := d.Set(mkResourceVirtualEnvironmentVMTPMState, []interface{}{tpmState})
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	currentPCIList := d.Get(mkResourceVirtualEnvironmentVMHostPCI).([]interface{})
	pciMap := map[string]interface{}{}

	pciDevices := getPCIInfo(vmConfig, d)
	for pi, pp := range pciDevices {
		if (pp == nil) || (pp.DeviceIDs == nil && pp.Mapping == nil) {
			continue
		}

		pci := map[string]interface{}{}

		pci[mkResourceVirtualEnvironmentVMHostPCIDevice] = pi
		if pp.DeviceIDs != nil {
			pci[mkResourceVirtualEnvironmentVMHostPCIDeviceID] = strings.Join(*pp.DeviceIDs, ";")
		} else {
			pci[mkResourceVirtualEnvironmentVMHostPCIDeviceID] = ""
		}

		if pp.MDev != nil {
			pci[mkResourceVirtualEnvironmentVMHostPCIDeviceMDev] = *pp.MDev
		} else {
			pci[mkResourceVirtualEnvironmentVMHostPCIDeviceMDev] = ""
		}

		if pp.PCIExpress != nil {
			pci[mkResourceVirtualEnvironmentVMHostPCIDevicePCIE] = *pp.PCIExpress
		} else {
			pci[mkResourceVirtualEnvironmentVMHostPCIDevicePCIE] = false
		}

		if pp.ROMBAR != nil {
			pci[mkResourceVirtualEnvironmentVMHostPCIDeviceROMBAR] = *pp.ROMBAR
		} else {
			pci[mkResourceVirtualEnvironmentVMHostPCIDeviceROMBAR] = true
		}

		if pp.ROMFile != nil {
			pci[mkResourceVirtualEnvironmentVMHostPCIDeviceROMFile] = *pp.ROMFile
		} else {
			pci[mkResourceVirtualEnvironmentVMHostPCIDeviceROMFile] = ""
		}

		if pp.XVGA != nil {
			pci[mkResourceVirtualEnvironmentVMHostPCIDeviceXVGA] = *pp.XVGA
		} else {
			pci[mkResourceVirtualEnvironmentVMHostPCIDeviceXVGA] = false
		}

		if pp.Mapping != nil {
			pci[mkResourceVirtualEnvironmentVMHostPCIDeviceMapping] = *pp.Mapping
		} else {
			pci[mkResourceVirtualEnvironmentVMHostPCIDeviceMapping] = ""
		}

		pciMap[pi] = pci
	}

	if len(clone) == 0 || len(currentPCIList) > 0 {
		orderedPCIList := orderedListFromMap(pciMap)
		err := d.Set(mkResourceVirtualEnvironmentVMHostPCI, orderedPCIList)
		diags = append(diags, diag.FromErr(err)...)
	}

	currentUSBList := d.Get(mkResourceVirtualEnvironmentVMHostUSB).([]interface{})
	usbMap := map[string]interface{}{}

	usbDevices := getUSBInfo(vmConfig, d)
	for pi, pp := range usbDevices {
		if (pp == nil) || (pp.HostDevice == nil && pp.Mapping == nil) {
			continue
		}

		usb := map[string]interface{}{}

		if pp.HostDevice != nil {
			usb[mkResourceVirtualEnvironmentVMHostUSBDevice] = *pp.HostDevice
		} else {
			usb[mkResourceVirtualEnvironmentVMHostUSBDevice] = ""
		}

		if pp.USB3 != nil {
			usb[mkResourceVirtualEnvironmentVMHostUSBDeviceUSB3] = *pp.USB3
		} else {
			usb[mkResourceVirtualEnvironmentVMHostUSBDeviceUSB3] = false
		}

		if pp.Mapping != nil {
			usb[mkResourceVirtualEnvironmentVMHostUSBDeviceMapping] = *pp.Mapping
		} else {
			usb[mkResourceVirtualEnvironmentVMHostUSBDeviceMapping] = ""
		}

		usbMap[pi] = usb
	}

	if len(clone) == 0 || len(currentUSBList) > 0 {
		// todo: reordering of devices by PVE may cause an issue here
		orderedUSBList := orderedListFromMap(usbMap)
		err := d.Set(mkResourceVirtualEnvironmentVMHostUSB, orderedUSBList)
		diags = append(diags, diag.FromErr(err)...)
	}

	// Compare the initialization configuration to the one stored in the state.
	initialization := map[string]interface{}{}

	initializationInterface := findExistingCloudInitDrive(vmConfig, vmID, "")
	if initializationInterface != "" {
		initializationDevice := getIdeDevice(vmConfig, initializationInterface)
		fileVolumeParts := strings.Split(initializationDevice.FileVolume, ":")

		initialization[mkResourceVirtualEnvironmentVMInitializationInterface] = initializationInterface
		initialization[mkResourceVirtualEnvironmentVMInitializationDatastoreID] = fileVolumeParts[0]
	}

	if vmConfig.CloudInitDNSDomain != nil || vmConfig.CloudInitDNSServer != nil {
		initializationDNS := map[string]interface{}{}

		if vmConfig.CloudInitDNSDomain != nil {
			initializationDNS[mkResourceVirtualEnvironmentVMInitializationDNSDomain] = *vmConfig.CloudInitDNSDomain
		} else {
			initializationDNS[mkResourceVirtualEnvironmentVMInitializationDNSDomain] = ""
		}

		// check what we have in the plan
		currentInitializationDNSBlock := map[string]interface{}{}
		currentInitialization := d.Get(mkResourceVirtualEnvironmentVMInitialization).([]interface{})

		if len(currentInitialization) > 0 {
			currentInitializationBlock := currentInitialization[0].(map[string]interface{})
			//nolint:lll
			currentInitializationDNS := currentInitializationBlock[mkResourceVirtualEnvironmentVMInitializationDNS].([]interface{})
			if len(currentInitializationDNS) > 0 {
				currentInitializationDNSBlock = currentInitializationDNS[0].(map[string]interface{})
			}
		}

		//nolint:lll
		currentInitializationDNSServer, ok := currentInitializationDNSBlock[mkResourceVirtualEnvironmentVMInitializationDNSServer]
		if vmConfig.CloudInitDNSServer != nil {
			if ok && currentInitializationDNSServer != "" {
				// the template is using deprecated attribute mkResourceVirtualEnvironmentVMInitializationDNSServer
				initializationDNS[mkResourceVirtualEnvironmentVMInitializationDNSServer] = *vmConfig.CloudInitDNSServer
			} else {
				dnsServer := strings.Split(*vmConfig.CloudInitDNSServer, " ")
				initializationDNS[mkResourceVirtualEnvironmentVMInitializationDNSServers] = dnsServer
			}
		} else {
			initializationDNS[mkResourceVirtualEnvironmentVMInitializationDNSServer] = ""
			initializationDNS[mkResourceVirtualEnvironmentVMInitializationDNSServers] = []string{}
		}

		initialization[mkResourceVirtualEnvironmentVMInitializationDNS] = []interface{}{
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
					ipv4[mkResourceVirtualEnvironmentVMInitializationIPConfigIPv4Address] = *ipConfig.IPv4
				} else {
					ipv4[mkResourceVirtualEnvironmentVMInitializationIPConfigIPv4Address] = ""
				}

				if ipConfig.GatewayIPv4 != nil {
					ipv4[mkResourceVirtualEnvironmentVMInitializationIPConfigIPv4Gateway] = *ipConfig.GatewayIPv4
				} else {
					ipv4[mkResourceVirtualEnvironmentVMInitializationIPConfigIPv4Gateway] = ""
				}

				ipConfigItem[mkResourceVirtualEnvironmentVMInitializationIPConfigIPv4] = []interface{}{
					ipv4,
				}
			} else {
				ipConfigItem[mkResourceVirtualEnvironmentVMInitializationIPConfigIPv4] = []interface{}{}
			}

			if ipConfig.GatewayIPv6 != nil || ipConfig.IPv6 != nil {
				ipv6 := map[string]interface{}{}

				if ipConfig.IPv6 != nil {
					ipv6[mkResourceVirtualEnvironmentVMInitializationIPConfigIPv6Address] = *ipConfig.IPv6
				} else {
					ipv6[mkResourceVirtualEnvironmentVMInitializationIPConfigIPv6Address] = ""
				}

				if ipConfig.GatewayIPv6 != nil {
					ipv6[mkResourceVirtualEnvironmentVMInitializationIPConfigIPv6Gateway] = *ipConfig.GatewayIPv6
				} else {
					ipv6[mkResourceVirtualEnvironmentVMInitializationIPConfigIPv6Gateway] = ""
				}

				ipConfigItem[mkResourceVirtualEnvironmentVMInitializationIPConfigIPv6] = []interface{}{
					ipv6,
				}
			} else {
				ipConfigItem[mkResourceVirtualEnvironmentVMInitializationIPConfigIPv6] = []interface{}{}
			}
		} else {
			ipConfigItem[mkResourceVirtualEnvironmentVMInitializationIPConfigIPv4] = []interface{}{}
			ipConfigItem[mkResourceVirtualEnvironmentVMInitializationIPConfigIPv6] = []interface{}{}
		}

		ipConfigList[ipConfigIndex] = ipConfigItem
	}

	if ipConfigLast >= 0 {
		initialization[mkResourceVirtualEnvironmentVMInitializationIPConfig] = ipConfigList[:ipConfigLast+1]
	}

	//nolint:nestif
	if vmConfig.CloudInitPassword != nil || vmConfig.CloudInitSSHKeys != nil ||
		vmConfig.CloudInitUsername != nil {
		initializationUserAccount := map[string]interface{}{}

		if vmConfig.CloudInitSSHKeys != nil {
			initializationUserAccount[mkResourceVirtualEnvironmentVMInitializationUserAccountKeys] = []string(
				*vmConfig.CloudInitSSHKeys,
			)
		} else {
			initializationUserAccount[mkResourceVirtualEnvironmentVMInitializationUserAccountKeys] = []string{}
		}

		if vmConfig.CloudInitPassword != nil {
			initializationUserAccount[mkResourceVirtualEnvironmentVMInitializationUserAccountPassword] = *vmConfig.CloudInitPassword
		} else {
			initializationUserAccount[mkResourceVirtualEnvironmentVMInitializationUserAccountPassword] = ""
		}

		if vmConfig.CloudInitUsername != nil {
			initializationUserAccount[mkResourceVirtualEnvironmentVMInitializationUserAccountUsername] = *vmConfig.CloudInitUsername
		} else {
			initializationUserAccount[mkResourceVirtualEnvironmentVMInitializationUserAccountUsername] = ""
		}

		initialization[mkResourceVirtualEnvironmentVMInitializationUserAccount] = []interface{}{
			initializationUserAccount,
		}
	}

	if vmConfig.CloudInitFiles != nil {
		if vmConfig.CloudInitFiles.UserVolume != nil {
			initialization[mkResourceVirtualEnvironmentVMInitializationUserDataFileID] = *vmConfig.CloudInitFiles.UserVolume
		} else {
			initialization[mkResourceVirtualEnvironmentVMInitializationUserDataFileID] = ""
		}

		if vmConfig.CloudInitFiles.VendorVolume != nil {
			initialization[mkResourceVirtualEnvironmentVMInitializationVendorDataFileID] = *vmConfig.CloudInitFiles.VendorVolume
		} else {
			initialization[mkResourceVirtualEnvironmentVMInitializationVendorDataFileID] = ""
		}

		if vmConfig.CloudInitFiles.NetworkVolume != nil {
			initialization[mkResourceVirtualEnvironmentVMInitializationNetworkDataFileID] = *vmConfig.CloudInitFiles.NetworkVolume
		} else {
			initialization[mkResourceVirtualEnvironmentVMInitializationNetworkDataFileID] = ""
		}

		if vmConfig.CloudInitFiles.MetaVolume != nil {
			initialization[mkResourceVirtualEnvironmentVMInitializationMetaDataFileID] = *vmConfig.CloudInitFiles.MetaVolume
		} else {
			initialization[mkResourceVirtualEnvironmentVMInitializationMetaDataFileID] = ""
		}
	} else if len(initialization) > 0 {
		initialization[mkResourceVirtualEnvironmentVMInitializationUserDataFileID] = ""
		initialization[mkResourceVirtualEnvironmentVMInitializationVendorDataFileID] = ""
		initialization[mkResourceVirtualEnvironmentVMInitializationNetworkDataFileID] = ""
		initialization[mkResourceVirtualEnvironmentVMInitializationMetaDataFileID] = ""
	}

	if vmConfig.CloudInitType != nil {
		initialization[mkResourceVirtualEnvironmentVMInitializationType] = *vmConfig.CloudInitType
	} else if len(initialization) > 0 {
		initialization[mkResourceVirtualEnvironmentVMInitializationType] = ""
	}

	currentInitialization := d.Get(mkResourceVirtualEnvironmentVMInitialization).([]interface{})

	if len(clone) > 0 {
		if len(currentInitialization) > 0 {
			if len(initialization) > 0 {
				err := d.Set(
					mkResourceVirtualEnvironmentVMInitialization,
					[]interface{}{initialization},
				)
				diags = append(diags, diag.FromErr(err)...)
			} else {
				err := d.Set(mkResourceVirtualEnvironmentVMInitialization, []interface{}{})
				diags = append(diags, diag.FromErr(err)...)
			}
		}
	} else if len(initialization) > 0 {
		err := d.Set(mkResourceVirtualEnvironmentVMInitialization, []interface{}{initialization})
		diags = append(diags, diag.FromErr(err)...)
	} else {
		err := d.Set(mkResourceVirtualEnvironmentVMInitialization, []interface{}{})
		diags = append(diags, diag.FromErr(err)...)
	}

	// Compare the operating system configuration to the one stored in the state.
	kvmArguments := map[string]interface{}{}

	if vmConfig.KVMArguments != nil {
		kvmArguments[mkResourceVirtualEnvironmentVMKVMArguments] = *vmConfig.KVMArguments
	} else {
		kvmArguments[mkResourceVirtualEnvironmentVMKVMArguments] = ""
	}

	// Compare the memory configuration to the one stored in the state.
	memory := map[string]interface{}{}

	if vmConfig.DedicatedMemory != nil {
		memory[mkResourceVirtualEnvironmentVMMemoryDedicated] = int(*vmConfig.DedicatedMemory)
	} else {
		memory[mkResourceVirtualEnvironmentVMMemoryDedicated] = 0
	}

	if vmConfig.FloatingMemory != nil {
		memory[mkResourceVirtualEnvironmentVMMemoryFloating] = int(*vmConfig.FloatingMemory)
	} else {
		memory[mkResourceVirtualEnvironmentVMMemoryFloating] = 0
	}

	if vmConfig.SharedMemory != nil {
		memory[mkResourceVirtualEnvironmentVMMemoryShared] = vmConfig.SharedMemory.Size
	} else {
		memory[mkResourceVirtualEnvironmentVMMemoryShared] = 0
	}

	currentMemory := d.Get(mkResourceVirtualEnvironmentVMMemory).([]interface{})

	if len(clone) > 0 {
		if len(currentMemory) > 0 {
			err := d.Set(mkResourceVirtualEnvironmentVMMemory, []interface{}{memory})
			diags = append(diags, diag.FromErr(err)...)
		}
	} else if len(currentMemory) > 0 ||
		memory[mkResourceVirtualEnvironmentVMMemoryDedicated] != dvResourceVirtualEnvironmentVMMemoryDedicated ||
		memory[mkResourceVirtualEnvironmentVMMemoryFloating] != dvResourceVirtualEnvironmentVMMemoryFloating ||
		memory[mkResourceVirtualEnvironmentVMMemoryShared] != dvResourceVirtualEnvironmentVMMemoryShared {
		err := d.Set(mkResourceVirtualEnvironmentVMMemory, []interface{}{memory})
		diags = append(diags, diag.FromErr(err)...)
	}

	// Compare the network devices to those stored in the state.
	currentNetworkDeviceList := d.Get(mkResourceVirtualEnvironmentVMNetworkDevice).([]interface{})

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
				networkDevice[mkResourceVirtualEnvironmentVMNetworkDeviceBridge] = *nd.Bridge
			} else {
				networkDevice[mkResourceVirtualEnvironmentVMNetworkDeviceBridge] = ""
			}

			networkDevice[mkResourceVirtualEnvironmentVMNetworkDeviceEnabled] = nd.Enabled

			if nd.Firewall != nil {
				networkDevice[mkResourceVirtualEnvironmentVMNetworkDeviceFirewall] = *nd.Firewall
			} else {
				networkDevice[mkResourceVirtualEnvironmentVMNetworkDeviceFirewall] = false
			}

			if nd.MACAddress != nil {
				macAddresses[ni] = *nd.MACAddress
			} else {
				macAddresses[ni] = ""
			}

			networkDevice[mkResourceVirtualEnvironmentVMNetworkDeviceMACAddress] = macAddresses[ni]
			networkDevice[mkResourceVirtualEnvironmentVMNetworkDeviceModel] = nd.Model

			if nd.Queues != nil {
				networkDevice[mkResourceVirtualEnvironmentVMNetworkDeviceQueues] = *nd.Queues
			} else {
				networkDevice[mkResourceVirtualEnvironmentVMNetworkDeviceQueues] = 0
			}

			if nd.RateLimit != nil {
				networkDevice[mkResourceVirtualEnvironmentVMNetworkDeviceRateLimit] = *nd.RateLimit
			} else {
				networkDevice[mkResourceVirtualEnvironmentVMNetworkDeviceRateLimit] = 0
			}

			if nd.Tag != nil {
				networkDevice[mkResourceVirtualEnvironmentVMNetworkDeviceVLANID] = nd.Tag
			} else {
				networkDevice[mkResourceVirtualEnvironmentVMNetworkDeviceVLANID] = 0
			}
			if nd.MTU != nil {
				networkDevice[mkResourceVirtualEnvironmentVMNetworkDeviceMTU] = nd.MTU
			} else {
				networkDevice[mkResourceVirtualEnvironmentVMNetworkDeviceMTU] = 0
			}
		} else {
			macAddresses[ni] = ""
			networkDevice[mkResourceVirtualEnvironmentVMNetworkDeviceEnabled] = false
		}

		networkDeviceList[ni] = networkDevice
	}

	if len(clone) > 0 {
		if len(currentNetworkDeviceList) > 0 {
			err := d.Set(
				mkResourceVirtualEnvironmentVMMACAddresses,
				macAddresses[0:len(currentNetworkDeviceList)],
			)
			diags = append(diags, diag.FromErr(err)...)
			err = d.Set(
				mkResourceVirtualEnvironmentVMNetworkDevice,
				networkDeviceList[:networkDeviceLast+1],
			)
			diags = append(diags, diag.FromErr(err)...)
		}
	} else {
		err := d.Set(mkResourceVirtualEnvironmentVMMACAddresses, macAddresses[0:len(currentNetworkDeviceList)])
		diags = append(diags, diag.FromErr(err)...)

		if len(currentNetworkDeviceList) > 0 || networkDeviceLast > -1 {
			err := d.Set(mkResourceVirtualEnvironmentVMNetworkDevice, networkDeviceList[:networkDeviceLast+1])
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	// Compare the operating system configuration to the one stored in the state.
	operatingSystem := map[string]interface{}{}

	if vmConfig.OSType != nil {
		operatingSystem[mkResourceVirtualEnvironmentVMOperatingSystemType] = *vmConfig.OSType
	} else {
		operatingSystem[mkResourceVirtualEnvironmentVMOperatingSystemType] = ""
	}

	currentOperatingSystem := d.Get(mkResourceVirtualEnvironmentVMOperatingSystem).([]interface{})

	if len(clone) > 0 {
		if len(currentOperatingSystem) > 0 {
			err := d.Set(
				mkResourceVirtualEnvironmentVMOperatingSystem,
				[]interface{}{operatingSystem},
			)
			diags = append(diags, diag.FromErr(err)...)
		}
	} else if len(currentOperatingSystem) > 0 ||
		operatingSystem[mkResourceVirtualEnvironmentVMOperatingSystemType] != dvResourceVirtualEnvironmentVMOperatingSystemType {
		err := d.Set(mkResourceVirtualEnvironmentVMOperatingSystem, []interface{}{operatingSystem})
		diags = append(diags, diag.FromErr(err)...)
	} else {
		err := d.Set(mkResourceVirtualEnvironmentVMOperatingSystem, []interface{}{})
		diags = append(diags, diag.FromErr(err)...)
	}

	// Compare the pool ID to the value stored in the state.
	currentPoolID := d.Get(mkResourceVirtualEnvironmentVMPoolID).(string)

	if len(clone) == 0 || currentPoolID != dvResourceVirtualEnvironmentVMPoolID {
		if vmConfig.PoolID != nil {
			err := d.Set(mkResourceVirtualEnvironmentVMPoolID, *vmConfig.PoolID)
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
			m[mkResourceVirtualEnvironmentVMSerialDeviceDevice] = *sd
			serialDevicesCount = sdi + 1
		} else {
			m[mkResourceVirtualEnvironmentVMSerialDeviceDevice] = ""
		}

		serialDevices[sdi] = m
	}

	currentSerialDevice := d.Get(mkResourceVirtualEnvironmentVMSerialDevice).([]interface{})

	if len(clone) == 0 || len(currentSerialDevice) > 0 {
		err := d.Set(mkResourceVirtualEnvironmentVMSerialDevice, serialDevices[:serialDevicesCount])
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
			smbios[mkResourceVirtualEnvironmentVMSMBIOSFamily] = string(b)
		} else {
			smbios[mkResourceVirtualEnvironmentVMSMBIOSFamily] = dvResourceVirtualEnvironmentVMSMBIOSFamily
		}

		if vmConfig.SMBIOS.Manufacturer != nil {
			b, err := base64.StdEncoding.DecodeString(*vmConfig.SMBIOS.Manufacturer)
			diags = append(diags, diag.FromErr(err)...)
			smbios[mkResourceVirtualEnvironmentVMSMBIOSManufacturer] = string(b)
		} else {
			smbios[mkResourceVirtualEnvironmentVMSMBIOSManufacturer] = dvResourceVirtualEnvironmentVMSMBIOSManufacturer
		}

		if vmConfig.SMBIOS.Product != nil {
			b, err := base64.StdEncoding.DecodeString(*vmConfig.SMBIOS.Product)
			diags = append(diags, diag.FromErr(err)...)
			smbios[mkResourceVirtualEnvironmentVMSMBIOSProduct] = string(b)
		} else {
			smbios[mkResourceVirtualEnvironmentVMSMBIOSProduct] = dvResourceVirtualEnvironmentVMSMBIOSProduct
		}

		if vmConfig.SMBIOS.Serial != nil {
			b, err := base64.StdEncoding.DecodeString(*vmConfig.SMBIOS.Serial)
			diags = append(diags, diag.FromErr(err)...)
			smbios[mkResourceVirtualEnvironmentVMSMBIOSSerial] = string(b)
		} else {
			smbios[mkResourceVirtualEnvironmentVMSMBIOSSerial] = dvResourceVirtualEnvironmentVMSMBIOSSerial
		}

		if vmConfig.SMBIOS.SKU != nil {
			b, err := base64.StdEncoding.DecodeString(*vmConfig.SMBIOS.SKU)
			diags = append(diags, diag.FromErr(err)...)
			smbios[mkResourceVirtualEnvironmentVMSMBIOSSKU] = string(b)
		} else {
			smbios[mkResourceVirtualEnvironmentVMSMBIOSSKU] = dvResourceVirtualEnvironmentVMSMBIOSSKU
		}

		if vmConfig.SMBIOS.Version != nil {
			b, err := base64.StdEncoding.DecodeString(*vmConfig.SMBIOS.Version)
			diags = append(diags, diag.FromErr(err)...)
			smbios[mkResourceVirtualEnvironmentVMSMBIOSVersion] = string(b)
		} else {
			smbios[mkResourceVirtualEnvironmentVMSMBIOSVersion] = dvResourceVirtualEnvironmentVMSMBIOSVersion
		}

		if vmConfig.SMBIOS.UUID != nil {
			smbios[mkResourceVirtualEnvironmentVMSMBIOSUUID] = *vmConfig.SMBIOS.UUID
		} else {
			smbios[mkResourceVirtualEnvironmentVMSMBIOSUUID] = nil
		}
	}

	currentSMBIOS := d.Get(mkResourceVirtualEnvironmentVMSMBIOS).([]interface{})

	//nolint:gocritic
	if len(clone) > 0 {
		if len(currentSMBIOS) > 0 {
			err := d.Set(mkResourceVirtualEnvironmentVMSMBIOS, currentSMBIOS)
			diags = append(diags, diag.FromErr(err)...)
		}
	} else if len(smbios) == 0 {
		err := d.Set(mkResourceVirtualEnvironmentVMSMBIOS, []interface{}{})
		diags = append(diags, diag.FromErr(err)...)
	} else if len(currentSMBIOS) > 0 ||
		smbios[mkResourceVirtualEnvironmentVMSMBIOSFamily] != dvResourceVirtualEnvironmentVMSMBIOSFamily ||
		smbios[mkResourceVirtualEnvironmentVMSMBIOSManufacturer] != dvResourceVirtualEnvironmentVMSMBIOSManufacturer ||
		smbios[mkResourceVirtualEnvironmentVMSMBIOSProduct] != dvResourceVirtualEnvironmentVMSMBIOSProduct ||
		smbios[mkResourceVirtualEnvironmentVMSMBIOSSerial] != dvResourceVirtualEnvironmentVMSMBIOSSerial ||
		smbios[mkResourceVirtualEnvironmentVMSMBIOSSKU] != dvResourceVirtualEnvironmentVMSMBIOSSKU ||
		smbios[mkResourceVirtualEnvironmentVMSMBIOSVersion] != dvResourceVirtualEnvironmentVMSMBIOSVersion {
		err := d.Set(mkResourceVirtualEnvironmentVMSMBIOS, []interface{}{smbios})
		diags = append(diags, diag.FromErr(err)...)
	}

	// Compare the startup order to the one stored in the state.
	var startup map[string]interface{}

	//nolint:nestif
	if vmConfig.StartupOrder != nil {
		startup = map[string]interface{}{}

		if vmConfig.StartupOrder.Order != nil {
			startup[mkResourceVirtualEnvironmentVMStartupOrder] = *vmConfig.StartupOrder.Order
		} else {
			startup[mkResourceVirtualEnvironmentVMStartupOrder] = dvResourceVirtualEnvironmentVMStartupOrder
		}

		if vmConfig.StartupOrder.Up != nil {
			startup[mkResourceVirtualEnvironmentVMStartupUpDelay] = *vmConfig.StartupOrder.Up
		} else {
			startup[mkResourceVirtualEnvironmentVMStartupUpDelay] = dvResourceVirtualEnvironmentVMStartupUpDelay
		}

		if vmConfig.StartupOrder.Down != nil {
			startup[mkResourceVirtualEnvironmentVMStartupDownDelay] = *vmConfig.StartupOrder.Down
		} else {
			startup[mkResourceVirtualEnvironmentVMStartupDownDelay] = dvResourceVirtualEnvironmentVMStartupDownDelay
		}
	}

	currentStartup := d.Get(mkResourceVirtualEnvironmentVMStartup).([]interface{})

	//nolint:gocritic
	if len(clone) > 0 {
		if len(currentStartup) > 0 {
			err := d.Set(mkResourceVirtualEnvironmentVMStartup, []interface{}{startup})
			diags = append(diags, diag.FromErr(err)...)
		}
	} else if len(startup) == 0 {
		err := d.Set(mkResourceVirtualEnvironmentVMStartup, []interface{}{})
		diags = append(diags, diag.FromErr(err)...)
	} else if len(currentStartup) > 0 ||
		startup[mkResourceVirtualEnvironmentVMStartupOrder] != mkResourceVirtualEnvironmentVMStartupOrder ||
		startup[mkResourceVirtualEnvironmentVMStartupUpDelay] != dvResourceVirtualEnvironmentVMStartupUpDelay ||
		startup[mkResourceVirtualEnvironmentVMStartupDownDelay] != dvResourceVirtualEnvironmentVMStartupDownDelay {
		err := d.Set(mkResourceVirtualEnvironmentVMStartup, []interface{}{startup})
		diags = append(diags, diag.FromErr(err)...)
	}

	// Compare the VGA configuration to the one stored in the state.
	vga := map[string]interface{}{}

	if vmConfig.VGADevice != nil {
		vgaEnabled := true

		if vmConfig.VGADevice.Type != nil {
			vgaEnabled = *vmConfig.VGADevice.Type != "none"
		}

		vga[mkResourceVirtualEnvironmentVMVGAEnabled] = vgaEnabled

		if vmConfig.VGADevice.Memory != nil {
			vga[mkResourceVirtualEnvironmentVMVGAMemory] = *vmConfig.VGADevice.Memory
		} else {
			vga[mkResourceVirtualEnvironmentVMVGAMemory] = 0
		}

		if vgaEnabled {
			if vmConfig.VGADevice.Type != nil {
				vga[mkResourceVirtualEnvironmentVMVGAType] = *vmConfig.VGADevice.Type
			} else {
				vga[mkResourceVirtualEnvironmentVMVGAType] = ""
			}
		}
	} else {
		vga[mkResourceVirtualEnvironmentVMVGAEnabled] = true
		vga[mkResourceVirtualEnvironmentVMVGAMemory] = 0
		vga[mkResourceVirtualEnvironmentVMVGAType] = ""
	}

	currentVGA := d.Get(mkResourceVirtualEnvironmentVMVGA).([]interface{})

	if len(clone) > 0 {
		if len(currentVGA) > 0 {
			err := d.Set(mkResourceVirtualEnvironmentVMVGA, []interface{}{vga})
			diags = append(diags, diag.FromErr(err)...)
		}
	} else if len(currentVGA) > 0 ||
		vga[mkResourceVirtualEnvironmentVMVGAEnabled] != dvResourceVirtualEnvironmentVMVGAEnabled ||
		vga[mkResourceVirtualEnvironmentVMVGAMemory] != dvResourceVirtualEnvironmentVMVGAMemory ||
		vga[mkResourceVirtualEnvironmentVMVGAType] != dvResourceVirtualEnvironmentVMVGAType {
		err := d.Set(mkResourceVirtualEnvironmentVMVGA, []interface{}{vga})
		diags = append(diags, diag.FromErr(err)...)
	} else {
		err := d.Set(mkResourceVirtualEnvironmentVMVGA, []interface{}{})
		diags = append(diags, diag.FromErr(err)...)
	}

	// Compare SCSI hardware type
	scsiHardware := d.Get(mkResourceVirtualEnvironmentVMSCSIHardware).(string)

	if len(clone) == 0 || scsiHardware != dvResourceVirtualEnvironmentVMSCSIHardware {
		if vmConfig.SCSIHardware != nil {
			err := d.Set(mkResourceVirtualEnvironmentVMSCSIHardware, *vmConfig.SCSIHardware)
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	diags = append(
		diags,
		vmReadNetworkValues(ctx, d, m, vmID, vmConfig)...)

	// during import these core attributes might not be set, so set them explicitly here
	d.SetId(strconv.Itoa(vmID))
	e := d.Set(mkResourceVirtualEnvironmentVMVMID, vmID)
	diags = append(diags, diag.FromErr(e)...)
	e = d.Set(mkResourceVirtualEnvironmentVMNodeName, nodeName)
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

	nodeName := d.Get(mkResourceVirtualEnvironmentVMNodeName).(string)

	vmAPI := api.Node(nodeName).VM(vmID)

	started := d.Get(mkResourceVirtualEnvironmentVMStarted).(bool)

	var ipv4Addresses []interface{}
	var ipv6Addresses []interface{}
	var networkInterfaceNames []interface{}

	if started {
		if vmConfig.Agent != nil && vmConfig.Agent.Enabled != nil && *vmConfig.Agent.Enabled {
			resource := VM()
			agentBlock, err := structure.GetSchemaBlock(
				resource,
				d,
				[]string{mkResourceVirtualEnvironmentVMAgent},
				0,
				true,
			)
			if err != nil {
				return diag.FromErr(err)
			}

			agentTimeout, err := time.ParseDuration(
				agentBlock[mkResourceVirtualEnvironmentVMAgentTimeout].(string),
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

			err = d.Set(mkResourceVirtualEnvironmentVMMACAddresses, macAddresses)
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	e = d.Set(mkResourceVirtualEnvironmentVMIPv4Addresses, ipv4Addresses)
	diags = append(diags, diag.FromErr(e)...)
	e = d.Set(mkResourceVirtualEnvironmentVMIPv6Addresses, ipv6Addresses)
	diags = append(diags, diag.FromErr(e)...)
	e = d.Set(mkResourceVirtualEnvironmentVMNetworkInterfaceNames, networkInterfaceNames)
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

	clone := d.Get(mkResourceVirtualEnvironmentVMClone).([]interface{})
	currentACPI := d.Get(mkResourceVirtualEnvironmentVMACPI).(bool)

	//nolint:gosimple
	if len(clone) == 0 || currentACPI != dvResourceVirtualEnvironmentVMACPI {
		if vmConfig.ACPI != nil {
			err = d.Set(mkResourceVirtualEnvironmentVMACPI, bool(*vmConfig.ACPI))
		} else {
			// Default value of "acpi" is "1" according to the API documentation.
			err = d.Set(mkResourceVirtualEnvironmentVMACPI, true)
		}
		diags = append(diags, diag.FromErr(err)...)
	}

	currentkvmArguments := d.Get(mkResourceVirtualEnvironmentVMKVMArguments).(string)

	if len(clone) == 0 || currentkvmArguments != dvResourceVirtualEnvironmentVMKVMArguments {
		// PVE API returns "args" as " " if it is set to empty.
		if vmConfig.KVMArguments != nil && len(strings.TrimSpace(*vmConfig.KVMArguments)) > 0 {
			err = d.Set(mkResourceVirtualEnvironmentVMKVMArguments, *vmConfig.KVMArguments)
		} else {
			// Default value of "args" is "" according to the API documentation.
			err = d.Set(mkResourceVirtualEnvironmentVMKVMArguments, "")
		}

		diags = append(diags, diag.FromErr(err)...)
	}

	currentBIOS := d.Get(mkResourceVirtualEnvironmentVMBIOS).(string)

	if len(clone) == 0 || currentBIOS != dvResourceVirtualEnvironmentVMBIOS {
		if vmConfig.BIOS != nil {
			err = d.Set(mkResourceVirtualEnvironmentVMBIOS, *vmConfig.BIOS)
		} else {
			// Default value of "bios" is "seabios" according to the API documentation.
			err = d.Set(mkResourceVirtualEnvironmentVMBIOS, "seabios")
		}
		diags = append(diags, diag.FromErr(err)...)
	}

	currentDescription := d.Get(mkResourceVirtualEnvironmentVMDescription).(string)

	if len(clone) == 0 || currentDescription != dvResourceVirtualEnvironmentVMDescription {
		if vmConfig.Description != nil {
			err = d.Set(mkResourceVirtualEnvironmentVMDescription, *vmConfig.Description)
		} else {
			// Default value of "description" is "" according to the API documentation.
			err = d.Set(mkResourceVirtualEnvironmentVMDescription, "")
		}
		diags = append(diags, diag.FromErr(err)...)
	}

	currentTags := d.Get(mkResourceVirtualEnvironmentVMTags).([]interface{})

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
		err = d.Set(mkResourceVirtualEnvironmentVMTags, tags)
		diags = append(diags, diag.FromErr(err)...)
	}

	currentKeyboardLayout := d.Get(mkResourceVirtualEnvironmentVMKeyboardLayout).(string)

	if len(clone) == 0 || currentKeyboardLayout != dvResourceVirtualEnvironmentVMKeyboardLayout {
		if vmConfig.KeyboardLayout != nil {
			err = d.Set(mkResourceVirtualEnvironmentVMKeyboardLayout, *vmConfig.KeyboardLayout)
		} else {
			// Default value of "keyboard" is "" according to the API documentation.
			err = d.Set(mkResourceVirtualEnvironmentVMKeyboardLayout, "")
		}
		diags = append(diags, diag.FromErr(err)...)
	}

	currentMachine := d.Get(mkResourceVirtualEnvironmentVMMachine).(string)

	if len(clone) == 0 || currentMachine != dvResourceVirtualEnvironmentVMMachineType {
		if vmConfig.Machine != nil {
			err = d.Set(mkResourceVirtualEnvironmentVMMachine, *vmConfig.Machine)
		} else {
			err = d.Set(mkResourceVirtualEnvironmentVMMachine, "")
		}
		diags = append(diags, diag.FromErr(err)...)
	}

	currentName := d.Get(mkResourceVirtualEnvironmentVMName).(string)

	if len(clone) == 0 || currentName != dvResourceVirtualEnvironmentVMName {
		if vmConfig.Name != nil {
			err = d.Set(mkResourceVirtualEnvironmentVMName, *vmConfig.Name)
		} else {
			// Default value of "name" is "" according to the API documentation.
			err = d.Set(mkResourceVirtualEnvironmentVMName, "")
		}
		diags = append(diags, diag.FromErr(err)...)
	}

	if !d.Get(mkResourceVirtualEnvironmentVMTemplate).(bool) {
		err = d.Set(mkResourceVirtualEnvironmentVMStarted, vmStatus.Status == "running")
		diags = append(diags, diag.FromErr(err)...)
	}

	currentTabletDevice := d.Get(mkResourceVirtualEnvironmentVMTabletDevice).(bool)

	//nolint:gosimple
	if len(clone) == 0 || currentTabletDevice != dvResourceVirtualEnvironmentVMTabletDevice {
		if vmConfig.TabletDeviceEnabled != nil {
			err = d.Set(
				mkResourceVirtualEnvironmentVMTabletDevice,
				bool(*vmConfig.TabletDeviceEnabled),
			)
		} else {
			// Default value of "tablet" is "1" according to the API documentation.
			err = d.Set(mkResourceVirtualEnvironmentVMTabletDevice, true)
		}
		diags = append(diags, diag.FromErr(err)...)
	}

	currentTemplate := d.Get(mkResourceVirtualEnvironmentVMTemplate).(bool)

	//nolint:gosimple
	if len(clone) == 0 || currentTemplate != dvResourceVirtualEnvironmentVMTemplate {
		if vmConfig.Template != nil {
			err = d.Set(mkResourceVirtualEnvironmentVMTemplate, bool(*vmConfig.Template))
		} else {
			// Default value of "template" is "0" according to the API documentation.
			err = d.Set(mkResourceVirtualEnvironmentVMTemplate, false)
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
	oldPoolValue, newPoolValue := d.GetChange(mkResourceVirtualEnvironmentVMPoolID)
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

	nodeName := d.Get(mkResourceVirtualEnvironmentVMNodeName).(string)
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
	if d.HasChange(mkResourceVirtualEnvironmentVMNodeName) {
		oldNodeNameValue, _ := d.GetChange(mkResourceVirtualEnvironmentVMNodeName)
		oldNodeName := oldNodeNameValue.(string)
		vmAPI := api.Node(oldNodeName).VM(vmID)

		migrateTimeout := d.Get(mkResourceVirtualEnvironmentVMTimeoutMigrate).(int)
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
			"ide0": vms.CustomStorageDevice{
				Enabled: false,
			},
			"ide1": vms.CustomStorageDevice{
				Enabled: false,
			},
			"ide2": vms.CustomStorageDevice{
				Enabled: false,
			},
			"ide3": vms.CustomStorageDevice{
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
	if d.HasChange(mkResourceVirtualEnvironmentVMACPI) {
		acpi := types.CustomBool(d.Get(mkResourceVirtualEnvironmentVMACPI).(bool))
		updateBody.ACPI = &acpi
		rebootRequired = true
	}

	if d.HasChange(mkResourceVirtualEnvironmentVMKVMArguments) {
		kvmArguments := d.Get(mkResourceVirtualEnvironmentVMKVMArguments).(string)
		updateBody.KVMArguments = &kvmArguments
		rebootRequired = true
	}

	if d.HasChange(mkResourceVirtualEnvironmentVMBIOS) {
		bios := d.Get(mkResourceVirtualEnvironmentVMBIOS).(string)
		updateBody.BIOS = &bios
		rebootRequired = true
	}

	if d.HasChange(mkResourceVirtualEnvironmentVMDescription) {
		description := d.Get(mkResourceVirtualEnvironmentVMDescription).(string)
		updateBody.Description = &description
	}

	if d.HasChange(mkResourceVirtualEnvironmentVMOnBoot) {
		startOnBoot := types.CustomBool(d.Get(mkResourceVirtualEnvironmentVMOnBoot).(bool))
		updateBody.StartOnBoot = &startOnBoot
	}

	if d.HasChange(mkResourceVirtualEnvironmentVMTags) {
		tagString := vmGetTagsString(d)
		updateBody.Tags = &tagString
	}

	if d.HasChange(mkResourceVirtualEnvironmentVMKeyboardLayout) {
		keyboardLayout := d.Get(mkResourceVirtualEnvironmentVMKeyboardLayout).(string)
		updateBody.KeyboardLayout = &keyboardLayout
		rebootRequired = true
	}

	if d.HasChange(mkResourceVirtualEnvironmentVMMachine) {
		machine := d.Get(mkResourceVirtualEnvironmentVMMachine).(string)
		updateBody.Machine = &machine
		rebootRequired = true
	}

	name := d.Get(mkResourceVirtualEnvironmentVMName).(string)

	if name == "" {
		del = append(del, "name")
	} else {
		updateBody.Name = &name
	}

	if d.HasChange(mkResourceVirtualEnvironmentVMTabletDevice) {
		tabletDevice := types.CustomBool(d.Get(mkResourceVirtualEnvironmentVMTabletDevice).(bool))
		updateBody.TabletDeviceEnabled = &tabletDevice
		rebootRequired = true
	}

	template := types.CustomBool(d.Get(mkResourceVirtualEnvironmentVMTemplate).(bool))

	if d.HasChange(mkResourceVirtualEnvironmentVMTemplate) {
		updateBody.Template = &template
		rebootRequired = true
	}

	// Prepare the new agent configuration.
	if d.HasChange(mkResourceVirtualEnvironmentVMAgent) {
		agentBlock, err := structure.GetSchemaBlock(
			resource,
			d,
			[]string{mkResourceVirtualEnvironmentVMAgent},
			0,
			true,
		)
		if err != nil {
			return diag.FromErr(err)
		}

		agentEnabled := types.CustomBool(
			agentBlock[mkResourceVirtualEnvironmentVMAgentEnabled].(bool),
		)
		agentTrim := types.CustomBool(agentBlock[mkResourceVirtualEnvironmentVMAgentTrim].(bool))
		agentType := agentBlock[mkResourceVirtualEnvironmentVMAgentType].(string)

		updateBody.Agent = &vms.CustomAgent{
			Enabled:         &agentEnabled,
			TrimClonedDisks: &agentTrim,
			Type:            &agentType,
		}

		rebootRequired = true
	}

	// Prepare the new audio devices.
	if d.HasChange(mkResourceVirtualEnvironmentVMAudioDevice) {
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
	if d.HasChange(mkResourceVirtualEnvironmentVMBootOrder) {
		bootOrder := d.Get(mkResourceVirtualEnvironmentVMBootOrder).([]interface{})
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
	if d.HasChange(mkResourceVirtualEnvironmentVMCDROM) {
		cdromBlock, err := structure.GetSchemaBlock(
			resource,
			d,
			[]string{mkResourceVirtualEnvironmentVMCDROM},
			0,
			true,
		)
		if err != nil {
			return diag.FromErr(err)
		}

		cdromEnabled := cdromBlock[mkResourceVirtualEnvironmentVMCDROMEnabled].(bool)
		cdromFileID := cdromBlock[mkResourceVirtualEnvironmentVMCDROMFileID].(string)
		cdromInterface := cdromBlock[mkResourceVirtualEnvironmentVMCDROMInterface].(string)

		old, _ := d.GetChange(mkResourceVirtualEnvironmentVMCDROM)

		if len(old.([]interface{})) > 0 {
			oldList := old.([]interface{})[0]
			oldBlock := oldList.(map[string]interface{})

			// If the interface is not set, use the default, for backward compatibility.
			oldInterface, ok := oldBlock[mkResourceVirtualEnvironmentVMCDROMInterface].(string)
			if !ok || oldInterface == "" {
				oldInterface = dvResourceVirtualEnvironmentVMCDROMInterface
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

		updateBody.IDEDevices[cdromInterface] = vms.CustomStorageDevice{
			Enabled:    cdromEnabled,
			FileVolume: cdromFileID,
			Media:      &cdromMedia,
		}
	}

	// Prepare the new CPU configuration.
	if d.HasChange(mkResourceVirtualEnvironmentVMCPU) {
		cpuBlock, err := structure.GetSchemaBlock(
			resource,
			d,
			[]string{mkResourceVirtualEnvironmentVMCPU},
			0,
			true,
		)
		if err != nil {
			return diag.FromErr(err)
		}

		cpuArchitecture := cpuBlock[mkResourceVirtualEnvironmentVMCPUArchitecture].(string)
		cpuCores := cpuBlock[mkResourceVirtualEnvironmentVMCPUCores].(int)
		cpuFlags := cpuBlock[mkResourceVirtualEnvironmentVMCPUFlags].([]interface{})
		cpuHotplugged := cpuBlock[mkResourceVirtualEnvironmentVMCPUHotplugged].(int)
		cpuLimit := cpuBlock[mkResourceVirtualEnvironmentVMCPULimit].(int)
		cpuNUMA := types.CustomBool(cpuBlock[mkResourceVirtualEnvironmentVMCPUNUMA].(bool))
		cpuSockets := cpuBlock[mkResourceVirtualEnvironmentVMCPUSockets].(int)
		cpuType := cpuBlock[mkResourceVirtualEnvironmentVMCPUType].(string)
		cpuUnits := cpuBlock[mkResourceVirtualEnvironmentVMCPUUnits].(int)

		// Only the root account is allowed to change the CPU architecture, which makes this check necessary.
		if api.API().IsRootTicket() ||
			cpuArchitecture != dvResourceVirtualEnvironmentVMCPUArchitecture {
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

	// Prepare the new disk device configuration.
	if d.HasChange(mkResourceVirtualEnvironmentVMDisk) {
		diskDeviceObjects, err := vmGetDiskDeviceObjects(d, nil)
		if err != nil {
			return diag.FromErr(err)
		}

		diskDeviceInfo := getDiskInfo(vmConfig, d)

		for prefix, diskMap := range diskDeviceObjects {
			if diskMap == nil {
				continue
			}

			for key, value := range diskMap {
				if diskDeviceInfo[key] == nil {
					return diag.Errorf("missing %s device %s", prefix, key)
				}

				tmp := *diskDeviceInfo[key]
				tmp.BurstableReadSpeedMbps = value.BurstableReadSpeedMbps
				tmp.BurstableWriteSpeedMbps = value.BurstableWriteSpeedMbps
				tmp.MaxReadSpeedMbps = value.MaxReadSpeedMbps
				tmp.MaxWriteSpeedMbps = value.MaxWriteSpeedMbps
				tmp.Cache = value.Cache

				switch prefix {
				case "virtio":
					{
						if updateBody.VirtualIODevices == nil {
							updateBody.VirtualIODevices = vms.CustomStorageDevices{}
						}

						updateBody.VirtualIODevices[key] = tmp
					}
				case "sata":
					{
						if updateBody.SATADevices == nil {
							updateBody.SATADevices = vms.CustomStorageDevices{}
						}

						updateBody.SATADevices[key] = tmp
					}
				case "scsi":
					{
						if updateBody.SCSIDevices == nil {
							updateBody.SCSIDevices = vms.CustomStorageDevices{}
						}

						updateBody.SCSIDevices[key] = tmp
					}
				case "ide":
					{
						// Investigate whether to support IDE mapping.
					}
				default:
					return diag.Errorf("device prefix %s not supported", prefix)
				}
			}
		}
	}

	// Prepare the new efi disk configuration.
	if d.HasChange(mkResourceVirtualEnvironmentVMEFIDisk) {
		efiDisk := vmGetEfiDisk(d, nil)

		updateBody.EFIDisk = efiDisk

		rebootRequired = true
	}

	// Prepare the new tpm state configuration.
	if d.HasChange(mkResourceVirtualEnvironmentVMTPMState) {
		tpmState := vmGetTPMState(d, nil)

		updateBody.TPMState = tpmState

		rebootRequired = true
	}

	// Prepare the new cloud-init configuration.
	stoppedBeforeUpdate := false

	if d.HasChange(mkResourceVirtualEnvironmentVMInitialization) {
		initializationConfig := vmGetCloudInitConfig(d)

		updateBody.CloudInitConfig = initializationConfig

		if updateBody.CloudInitConfig != nil {
			var fileVolume string

			initialization := d.Get(mkResourceVirtualEnvironmentVMInitialization).([]interface{})
			initializationBlock := initialization[0].(map[string]interface{})
			initializationDatastoreID := initializationBlock[mkResourceVirtualEnvironmentVMInitializationDatastoreID].(string)
			initializationInterface := initializationBlock[mkResourceVirtualEnvironmentVMInitializationInterface].(string)
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

			oldInit, _ := d.GetChange(mkResourceVirtualEnvironmentVMInitialization)
			oldInitBlock := oldInit.([]interface{})[0].(map[string]interface{})
			prevDatastoreID := oldInitBlock[mkResourceVirtualEnvironmentVMInitializationDatastoreID].(string)

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

			updateBody.IDEDevices[initializationInterface] = vms.CustomStorageDevice{
				Enabled:    true,
				FileVolume: fileVolume,
				Media:      &cdromMedia,
			}
		}

		rebootRequired = true
	}

	// Prepare the new hostpci devices configuration.
	if d.HasChange(mkResourceVirtualEnvironmentVMHostPCI) {
		updateBody.PCIDevices = vmGetHostPCIDeviceObjects(d)

		for i := len(updateBody.PCIDevices); i < maxResourceVirtualEnvironmentVMHostPCIDevices; i++ {
			del = append(del, fmt.Sprintf("hostpci%d", i))
		}

		rebootRequired = true
	}

	// Prepare the new usb devices configuration.
	if d.HasChange(mkResourceVirtualEnvironmentVMHostUSB) {
		updateBody.USBDevices = vmGetHostUSBDeviceObjects(d)

		for i := len(updateBody.USBDevices); i < maxResourceVirtualEnvironmentVMHostUSBDevices; i++ {
			del = append(del, fmt.Sprintf("usb%d", i))
		}

		rebootRequired = true
	}

	// Prepare the new memory configuration.
	if d.HasChange(mkResourceVirtualEnvironmentVMMemory) {
		memoryBlock, err := structure.GetSchemaBlock(
			resource,
			d,
			[]string{mkResourceVirtualEnvironmentVMMemory},
			0,
			true,
		)
		if err != nil {
			return diag.FromErr(err)
		}

		memoryDedicated := memoryBlock[mkResourceVirtualEnvironmentVMMemoryDedicated].(int)
		memoryFloating := memoryBlock[mkResourceVirtualEnvironmentVMMemoryFloating].(int)
		memoryShared := memoryBlock[mkResourceVirtualEnvironmentVMMemoryShared].(int)

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
	if d.HasChange(mkResourceVirtualEnvironmentVMNetworkDevice) {
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
	if d.HasChange(mkResourceVirtualEnvironmentVMOperatingSystem) {
		operatingSystem, err := structure.GetSchemaBlock(
			resource,
			d,
			[]string{mkResourceVirtualEnvironmentVMOperatingSystem},
			0,
			true,
		)
		if err != nil {
			return diag.FromErr(err)
		}

		operatingSystemType := operatingSystem[mkResourceVirtualEnvironmentVMOperatingSystemType].(string)

		updateBody.OSType = &operatingSystemType

		rebootRequired = true
	}

	// Prepare the new serial devices.
	if d.HasChange(mkResourceVirtualEnvironmentVMSerialDevice) {
		updateBody.SerialDevices = vmGetSerialDeviceList(d)

		for i := len(updateBody.SerialDevices); i < maxResourceVirtualEnvironmentVMSerialDevices; i++ {
			del = append(del, fmt.Sprintf("serial%d", i))
		}

		rebootRequired = true
	}

	if d.HasChange(mkResourceVirtualEnvironmentVMSMBIOS) {
		updateBody.SMBIOS = vmGetSMBIOS(d)
		if updateBody.SMBIOS == nil {
			del = append(del, "smbios1")
		}
	}

	if d.HasChange(mkResourceVirtualEnvironmentVMStartup) {
		updateBody.StartupOrder = vmGetStartupOrder(d)
		if updateBody.StartupOrder == nil {
			del = append(del, "startup")
		}
	}

	// Prepare the new VGA configuration.
	if d.HasChange(mkResourceVirtualEnvironmentVMVGA) {
		updateBody.VGADevice, e = vmGetVGADeviceObject(d)
		if e != nil {
			return diag.FromErr(e)
		}

		rebootRequired = true
	}

	// Prepare the new SCSI hardware type
	if d.HasChange(mkResourceVirtualEnvironmentVMSCSIHardware) {
		scsiHardware := d.Get(mkResourceVirtualEnvironmentVMSCSIHardware).(string)
		updateBody.SCSIHardware = &scsiHardware

		rebootRequired = true
	}

	if d.HasChanges(mkResourceVirtualEnvironmentVMHookScriptFileID) {
		hookScript := d.Get(mkResourceVirtualEnvironmentVMHookScriptFileID).(string)
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
	if (d.HasChange(mkResourceVirtualEnvironmentVMStarted) || stoppedBeforeUpdate) && !bool(template) {
		started := d.Get(mkResourceVirtualEnvironmentVMStarted).(bool)
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

	nodeName := d.Get(mkResourceVirtualEnvironmentVMNodeName).(string)
	started := d.Get(mkResourceVirtualEnvironmentVMStarted).(bool)
	template := d.Get(mkResourceVirtualEnvironmentVMTemplate).(bool)

	vmID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	vmAPI := api.Node(nodeName).VM(vmID)

	// Determine if any of the disks are changing location and/or size, and initiate the necessary actions.
	if d.HasChange(mkResourceVirtualEnvironmentVMDisk) {
		diskOld, diskNew := d.GetChange(mkResourceVirtualEnvironmentVMDisk)

		diskOldEntries, err := vmGetDiskDeviceObjects(
			d,
			diskOld.([]interface{}),
		)
		if err != nil {
			return diag.FromErr(err)
		}

		diskNewEntries, err := vmGetDiskDeviceObjects(
			d,
			diskNew.([]interface{}),
		)
		if err != nil {
			return diag.FromErr(err)
		}

		// Add efidisk if it has changes
		if d.HasChange(mkResourceVirtualEnvironmentVMEFIDisk) {
			diskOld, diskNew := d.GetChange(mkResourceVirtualEnvironmentVMEFIDisk)

			oldEfiDisk, e := vmGetEfiDiskAsStorageDevice(d, diskOld.([]interface{}))
			if e != nil {
				return diag.FromErr(e)
			}

			newEfiDisk, e := vmGetEfiDiskAsStorageDevice(d, diskNew.([]interface{}))
			if e != nil {
				return diag.FromErr(e)
			}

			if oldEfiDisk != nil {
				baseDiskInterface := diskDigitPrefix(*oldEfiDisk.Interface)
				diskOldEntries[baseDiskInterface][*oldEfiDisk.Interface] = *oldEfiDisk
			}

			if newEfiDisk != nil {
				baseDiskInterface := diskDigitPrefix(*newEfiDisk.Interface)
				diskNewEntries[baseDiskInterface][*newEfiDisk.Interface] = *newEfiDisk
			}

			if oldEfiDisk != nil && newEfiDisk != nil && oldEfiDisk.Size != newEfiDisk.Size {
				return diag.Errorf(
					"resizing of efidisks is not supported.",
				)
			}
		}

		// Add tpm state if it has changes
		if d.HasChange(mkResourceVirtualEnvironmentVMTPMState) {
			diskOld, diskNew := d.GetChange(mkResourceVirtualEnvironmentVMTPMState)

			oldTPMState := vmGetTPMStateAsStorageDevice(d, diskOld.([]interface{}))
			newTPMState := vmGetTPMStateAsStorageDevice(d, diskNew.([]interface{}))

			if oldTPMState != nil {
				baseDiskInterface := diskDigitPrefix(*oldTPMState.Interface)
				diskOldEntries[baseDiskInterface][*oldTPMState.Interface] = *oldTPMState
			}

			if newTPMState != nil {
				baseDiskInterface := diskDigitPrefix(*newTPMState.Interface)
				diskNewEntries[baseDiskInterface][*newTPMState.Interface] = *newTPMState
			}

			if oldTPMState != nil && newTPMState != nil && oldTPMState.Size != newTPMState.Size {
				return diag.Errorf(
					"resizing of tpm state is not supported.",
				)
			}
		}

		var diskMoveBodies []*vms.MoveDiskRequestBody

		var diskResizeBodies []*vms.ResizeDiskRequestBody

		shutdownForDisksRequired := false

		for prefix, diskMap := range diskOldEntries {
			for oldKey, oldDisk := range diskMap {
				if _, present := diskNewEntries[prefix][oldKey]; !present {
					return diag.Errorf(
						"deletion of disks not supported. Please delete disk by hand. Old Interface was %s",
						*oldDisk.Interface,
					)
				}

				if *oldDisk.ID != *diskNewEntries[prefix][oldKey].ID {
					if oldDisk.IsOwnedBy(vmID) {
						deleteOriginalDisk := types.CustomBool(true)

						diskMoveBodies = append(
							diskMoveBodies,
							&vms.MoveDiskRequestBody{
								DeleteOriginalDisk: &deleteOriginalDisk,
								Disk:               *oldDisk.Interface,
								TargetStorage:      *diskNewEntries[prefix][oldKey].ID,
							},
						)

						// Cannot be done while VM is running.
						shutdownForDisksRequired = true
					} else {
						return diag.Errorf(
							"Cannot move %s:%s to datastore %s in VM %d configuration, it is not owned by this VM!",
							*oldDisk.ID,
							*oldDisk.PathInDatastore(),
							*diskNewEntries[prefix][oldKey].ID,
							vmID,
						)
					}
				}

				if *oldDisk.Size < *diskNewEntries[prefix][oldKey].Size {
					if oldDisk.IsOwnedBy(vmID) {
						diskResizeBodies = append(
							diskResizeBodies,
							&vms.ResizeDiskRequestBody{
								Disk: *oldDisk.Interface,
								Size: *diskNewEntries[prefix][oldKey].Size,
							},
						)
					} else {
						return diag.Errorf(
							"Cannot resize %s:%s in VM %d configuration, it is not owned by this VM!",
							*oldDisk.ID,
							*oldDisk.PathInDatastore(),
							vmID,
						)
					}
				}
			}
		}

		if shutdownForDisksRequired && !template {
			if e := vmShutdown(ctx, vmAPI, d); e != nil {
				return e
			}
		}

		for _, reqBody := range diskMoveBodies {
			moveDiskTimeout := d.Get(mkResourceVirtualEnvironmentVMTimeoutMoveDisk).(int)

			err = vmAPI.MoveVMDisk(ctx, reqBody, moveDiskTimeout)
			if err != nil {
				return diag.FromErr(err)
			}
		}

		for _, reqBody := range diskResizeBodies {
			err = vmAPI.ResizeVMDisk(ctx, reqBody)
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
			rebootTimeout := d.Get(mkResourceVirtualEnvironmentVMTimeoutReboot).(int)

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

	nodeName := d.Get(mkResourceVirtualEnvironmentVMNodeName).(string)

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

	stop := d.Get(mkResourceVirtualEnvironmentVMStopOnDestroy).(bool)

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
	err = vmAPI.WaitForVMState(ctx, "", 60, 2)
	if err == nil {
		return diag.Errorf("failed to delete VM \"%d\"", vmID)
	}

	d.SetId("")

	return nil
}

func diskDigitPrefix(s string) string {
	for i, r := range s {
		if unicode.IsDigit(r) {
			return s[:i]
		}
	}

	return s
}

func getDiskInfo(resp *vms.GetResponseData, d *schema.ResourceData) map[string]*vms.CustomStorageDevice {
	currentDisk := d.Get(mkResourceVirtualEnvironmentVMDisk)

	currentDiskList := currentDisk.([]interface{})
	currentDiskMap := map[string]map[string]interface{}{}

	for _, v := range currentDiskList {
		diskMap := v.(map[string]interface{})
		diskInterface := diskMap[mkResourceVirtualEnvironmentVMDiskInterface].(string)

		currentDiskMap[diskInterface] = diskMap
	}

	storageDevices := map[string]*vms.CustomStorageDevice{}

	storageDevices["ide0"] = resp.IDEDevice0
	storageDevices["ide1"] = resp.IDEDevice1
	storageDevices["ide2"] = resp.IDEDevice2
	storageDevices["ide3"] = resp.IDEDevice3

	storageDevices["sata0"] = resp.SATADevice0
	storageDevices["sata1"] = resp.SATADevice1
	storageDevices["sata2"] = resp.SATADevice2
	storageDevices["sata3"] = resp.SATADevice3
	storageDevices["sata4"] = resp.SATADevice4
	storageDevices["sata5"] = resp.SATADevice5

	storageDevices["scsi0"] = resp.SCSIDevice0
	storageDevices["scsi1"] = resp.SCSIDevice1
	storageDevices["scsi2"] = resp.SCSIDevice2
	storageDevices["scsi3"] = resp.SCSIDevice3
	storageDevices["scsi4"] = resp.SCSIDevice4
	storageDevices["scsi5"] = resp.SCSIDevice5
	storageDevices["scsi6"] = resp.SCSIDevice6
	storageDevices["scsi7"] = resp.SCSIDevice7
	storageDevices["scsi8"] = resp.SCSIDevice8
	storageDevices["scsi9"] = resp.SCSIDevice9
	storageDevices["scsi10"] = resp.SCSIDevice10
	storageDevices["scsi11"] = resp.SCSIDevice11
	storageDevices["scsi12"] = resp.SCSIDevice12
	storageDevices["scsi13"] = resp.SCSIDevice13

	storageDevices["virtio0"] = resp.VirtualIODevice0
	storageDevices["virtio1"] = resp.VirtualIODevice1
	storageDevices["virtio2"] = resp.VirtualIODevice2
	storageDevices["virtio3"] = resp.VirtualIODevice3
	storageDevices["virtio4"] = resp.VirtualIODevice4
	storageDevices["virtio5"] = resp.VirtualIODevice5
	storageDevices["virtio6"] = resp.VirtualIODevice6
	storageDevices["virtio7"] = resp.VirtualIODevice7
	storageDevices["virtio8"] = resp.VirtualIODevice8
	storageDevices["virtio9"] = resp.VirtualIODevice9
	storageDevices["virtio10"] = resp.VirtualIODevice10
	storageDevices["virtio11"] = resp.VirtualIODevice11
	storageDevices["virtio12"] = resp.VirtualIODevice12
	storageDevices["virtio13"] = resp.VirtualIODevice13
	storageDevices["virtio14"] = resp.VirtualIODevice14
	storageDevices["virtio15"] = resp.VirtualIODevice15

	for k, v := range storageDevices {
		if v != nil {
			if currentDiskMap[k] != nil {
				if currentDiskMap[k][mkResourceVirtualEnvironmentVMDiskFileID] != nil {
					fileID := currentDiskMap[k][mkResourceVirtualEnvironmentVMDiskFileID].(string)
					v.FileID = &fileID
				}
			}

			if v.Size == nil {
				v.Size = new(types.DiskSize)
			}

			// defensive copy of the loop variable
			iface := k
			v.Interface = &iface
		}
	}

	return storageDevices
}

// getDiskDatastores returns a list of the used datastores in a VM.
func getDiskDatastores(vm *vms.GetResponseData, d *schema.ResourceData) []string {
	storageDevices := getDiskInfo(vm, d)
	datastoresSet := map[string]int{}

	for _, diskInfo := range storageDevices {
		// Ignore empty storage devices and storage devices (like ide) which may not have any media mounted
		if diskInfo == nil || diskInfo.FileVolume == "none" {
			continue
		}

		fileIDParts := strings.Split(diskInfo.FileVolume, ":")
		datastoresSet[fileIDParts[0]] = 1
	}

	if vm.EFIDisk != nil {
		fileIDParts := strings.Split(vm.EFIDisk.FileVolume, ":")
		datastoresSet[fileIDParts[0]] = 1
	}

	if vm.TPMState != nil {
		fileIDParts := strings.Split(vm.TPMState.FileVolume, ":")
		datastoresSet[fileIDParts[0]] = 1
	}

	datastores := []string{}
	for datastore := range datastoresSet {
		datastores = append(datastores, datastore)
	}

	return datastores
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
