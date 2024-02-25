package disk

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/bpg/terraform-provider-proxmox/proxmox"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/vms"
	"github.com/bpg/terraform-provider-proxmox/proxmox/ssh"
	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
	"github.com/bpg/terraform-provider-proxmox/proxmoxtf"
	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/structure"
	"github.com/bpg/terraform-provider-proxmox/utils"
)

// GetInfo returns the disk information for a VM.
func GetInfo(resp *vms.GetResponseData, d *schema.ResourceData) vms.CustomStorageDevices {
	currentDisk := d.Get(MkDisk)

	currentDiskList := currentDisk.([]interface{})
	currentDiskMap := map[string]map[string]interface{}{}

	for _, v := range currentDiskList {
		diskMap := v.(map[string]interface{})
		diskInterface := diskMap[mkDiskInterface].(string)

		currentDiskMap[diskInterface] = diskMap
	}

	storageDevices := vms.CustomStorageDevices{}

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
				if currentDiskMap[k][mkDiskFileID] != nil {
					fileID := currentDiskMap[k][mkDiskFileID].(string)
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

// CreateClone creates disks for a cloned VM.
func CreateClone(
	ctx context.Context,
	d *schema.ResourceData,
	planDisks map[string]vms.CustomStorageDevices,
	allDiskInfo vms.CustomStorageDevices,
	vmAPI *vms.Client,
) error {
	disk := d.Get(MkDisk).([]interface{})
	for i := range disk {
		diskBlock := disk[i].(map[string]interface{})
		diskInterface := diskBlock[mkDiskInterface].(string)
		dataStoreID := diskBlock[mkDiskDatastoreID].(string)
		diskSize := int64(diskBlock[mkDiskSize].(int))
		prefix := DigitPrefix(diskInterface)

		currentDiskInfo := allDiskInfo[diskInterface]
		configuredDiskInfo := planDisks[prefix][diskInterface]

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

			err := vmAPI.UpdateVM(ctx, diskUpdateBody)
			if err != nil {
				return fmt.Errorf("disk create fails: %w", err)
			}

			continue
		}

		if diskSize < currentDiskInfo.Size.InGigabytes() {
			return fmt.Errorf("disk resize fails requests size (%dG) is lower than current size (%d)",
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
			Size: *types.DiskSizeFromGigabytes(diskSize),
		}

		moveDisk := false

		if dataStoreID != "" {
			moveDisk = true

			if allDiskInfo[diskInterface] != nil {
				fileIDParts := strings.Split(allDiskInfo[diskInterface].FileVolume, ":")
				moveDisk = dataStoreID != fileIDParts[0]
			}
		}

		timeout := d.Get(MkTimeoutMoveDisk).(int)

		if moveDisk {
			err := vmAPI.MoveVMDisk(ctx, diskMoveBody, timeout)
			if err != nil {
				return fmt.Errorf("disk move fails: %w", err)
			}
		}

		if diskSize > currentDiskInfo.Size.InGigabytes() {
			err := vmAPI.ResizeVMDisk(ctx, diskResizeBody, timeout)
			if err != nil {
				return fmt.Errorf("disk resize fails: %w", err)
			}
		}
	}

	return nil
}

// DigitPrefix returns the prefix of a string that is not a digit.
func DigitPrefix(s string) string {
	for i, r := range s {
		if unicode.IsDigit(r) {
			return s[:i]
		}
	}

	return s
}

// GetDiskDeviceObjects returns a map of disk devices for a VM.
func GetDiskDeviceObjects(
	d *schema.ResourceData,
	resource *schema.Resource,
	disks []interface{},
) (map[string]vms.CustomStorageDevices, error) {
	var diskDevices []interface{}

	if disks != nil {
		diskDevices = disks
	} else {
		diskDevices = d.Get(MkDisk).([]interface{})
	}

	diskDeviceObjects := map[string]vms.CustomStorageDevices{}

	for _, diskEntry := range diskDevices {
		diskDevice := &vms.CustomStorageDevice{
			Enabled: true,
		}

		block := diskEntry.(map[string]interface{})
		datastoreID, _ := block[mkDiskDatastoreID].(string)
		pathInDatastore := ""

		if untyped, hasPathInDatastore := block[mkDiskPathInDatastore]; hasPathInDatastore {
			pathInDatastore = untyped.(string)
		}

		fileFormat, _ := block[mkDiskFileFormat].(string)
		fileID, _ := block[mkDiskFileID].(string)
		size, _ := block[mkDiskSize].(int)
		diskInterface, _ := block[mkDiskInterface].(string)
		ioThread := types.CustomBool(block[mkDiskIOThread].(bool))
		ssd := types.CustomBool(block[mkDiskSSD].(bool))
		discard := block[mkDiskDiscard].(string)
		cache := block[mkDiskCache].(string)

		speedBlock, err := structure.GetSchemaBlock(
			resource,
			d,
			[]string{MkDisk, mkDiskSpeed},
			0,
			false,
		)
		if err != nil {
			return diskDeviceObjects, fmt.Errorf("error getting disk speed block: %w", err)
		}

		if fileFormat == "" {
			fileFormat = dvDiskFileFormat
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

		diskDevice.DatastoreID = &datastoreID
		diskDevice.Interface = &diskInterface
		diskDevice.Format = &fileFormat
		diskDevice.FileID = &fileID
		diskSize := types.DiskSizeFromGigabytes(int64(size))
		diskDevice.Size = diskSize
		diskDevice.IOThread = &ioThread
		diskDevice.Discard = &discard
		diskDevice.Cache = &cache

		if !strings.HasPrefix(diskInterface, "virtio") {
			diskDevice.SSD = &ssd
		}

		if len(speedBlock) > 0 {
			speedLimitRead := speedBlock[mkDiskSpeedRead].(int)
			speedLimitReadBurstable := speedBlock[mkDiskSpeedReadBurstable].(int)
			speedLimitWrite := speedBlock[mkDiskSpeedWrite].(int)
			speedLimitWriteBurstable := speedBlock[mkDiskSpeedWriteBurstable].(int)

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

		baseDiskInterface := DigitPrefix(diskInterface)

		if baseDiskInterface != "virtio" && baseDiskInterface != "scsi" &&
			baseDiskInterface != "sata" {
			errorMsg := fmt.Sprintf(
				"Defined disk interface not supported. Interface was %s, but only virtio, sata and scsi are supported",
				diskInterface,
			)

			return diskDeviceObjects, errors.New(errorMsg)
		}

		if _, present := diskDeviceObjects[baseDiskInterface]; !present {
			diskDeviceObjects[baseDiskInterface] = vms.CustomStorageDevices{}
		}

		diskDeviceObjects[baseDiskInterface][diskInterface] = diskDevice
	}

	return diskDeviceObjects, nil
}

// CreateCustomDisks creates custom disks for a VM.
func CreateCustomDisks(
	ctx context.Context,
	nodeName string,
	d *schema.ResourceData,
	resource *schema.Resource,
	m interface{},
) diag.Diagnostics {
	vmID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	//nolint:prealloc
	var commands []string

	// Determine the ID of the next disk.
	disk := d.Get(MkDisk).([]interface{})
	diskCount := 0

	for _, d := range disk {
		block := d.(map[string]interface{})
		fileID, _ := block[mkDiskFileID].(string)

		if fileID == "" {
			diskCount++
		}
	}

	// Retrieve some information about the disk schema.
	resourceSchema := resource.Schema
	diskSchemaElem := resourceSchema[MkDisk].Elem
	diskSchemaResource := diskSchemaElem.(*schema.Resource)
	diskSpeedResource := diskSchemaResource.Schema[mkDiskSpeed]

	// Generate the commands required to import the specified disks.
	importedDiskCount := 0

	for _, d := range disk {
		block := d.(map[string]interface{})

		fileID, _ := block[mkDiskFileID].(string)

		if fileID == "" {
			continue
		}

		datastoreID, _ := block[mkDiskDatastoreID].(string)
		fileFormat, _ := block[mkDiskFileFormat].(string)
		size, _ := block[mkDiskSize].(int)
		speed := block[mkDiskSpeed].([]interface{})
		diskInterface, _ := block[mkDiskInterface].(string)
		ioThread := types.CustomBool(block[mkDiskIOThread].(bool))
		ssd := types.CustomBool(block[mkDiskSSD].(bool))
		discard, _ := block[mkDiskDiscard].(string)
		cache, _ := block[mkDiskCache].(string)

		if fileFormat == "" {
			fileFormat = dvDiskFileFormat
		}

		if len(speed) == 0 {
			diskSpeedDefault, err := diskSpeedResource.DefaultValue()
			if err != nil {
				return diag.FromErr(err)
			}

			speed = diskSpeedDefault.([]interface{})
		}

		speedBlock := speed[0].(map[string]interface{})
		speedLimitRead := speedBlock[mkDiskSpeedRead].(int)
		speedLimitReadBurstable := speedBlock[mkDiskSpeedReadBurstable].(int)
		speedLimitWrite := speedBlock[mkDiskSpeedWrite].(int)
		speedLimitWriteBurstable := speedBlock[mkDiskSpeedWriteBurstable].(int)

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
			ssh.TrySudo,
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

		out, err := api.SSH().ExecuteNodeCommands(ctx, nodeName, commands)
		if err != nil {
			if matches, e := regexp.Match(`pvesm: .* not found`, out); e == nil && matches {
				return diag.FromErr(ssh.NewErrUserHasNoPermission(api.SSH().Username()))
			}

			return diag.FromErr(err)
		}

		tflog.Debug(ctx, "vmCreateCustomDisks", map[string]interface{}{
			"output": string(out),
		})
	}

	return nil
}

// Read reads the disk configuration of a VM.
func Read(
	ctx context.Context,
	d *schema.ResourceData,
	diskObjects vms.CustomStorageDevices,
	vmID int,
	api proxmox.Client,
	nodeName string,
	isClone bool,
) diag.Diagnostics {
	currentDiskList := d.Get(MkDisk).([]interface{})
	diskMap := map[string]interface{}{}

	var diags diag.Diagnostics

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

		disk[mkDiskDatastoreID] = datastoreID
		disk[mkDiskPathInDatastore] = pathInDatastore

		if dd.Format == nil {
			disk[mkDiskFileFormat] = dvDiskFileFormat

			if datastoreID != "" {
				// disk format may not be returned by config API if it is default for the storage, and that may be different
				// from the default qcow2, so we need to read it from the storage API to make sure we have the correct value
				volume, e := api.Node(nodeName).Storage(datastoreID).GetDatastoreFile(ctx, dd.FileVolume)
				if e != nil {
					diags = append(diags, diag.FromErr(e)...)
					continue
				}

				disk[mkDiskFileFormat] = volume.FileFormat
			}
		} else {
			disk[mkDiskFileFormat] = dd.Format
		}

		if dd.FileID != nil {
			disk[mkDiskFileID] = dd.FileID
		}

		disk[mkDiskInterface] = di
		disk[mkDiskSize] = dd.Size.InGigabytes()

		if dd.BurstableReadSpeedMbps != nil ||
			dd.BurstableWriteSpeedMbps != nil ||
			dd.MaxReadSpeedMbps != nil ||
			dd.MaxWriteSpeedMbps != nil {
			speed := map[string]interface{}{}

			if dd.MaxReadSpeedMbps != nil {
				speed[mkDiskSpeedRead] = *dd.MaxReadSpeedMbps
			} else {
				speed[mkDiskSpeedRead] = 0
			}

			if dd.BurstableReadSpeedMbps != nil {
				speed[mkDiskSpeedReadBurstable] = *dd.BurstableReadSpeedMbps
			} else {
				speed[mkDiskSpeedReadBurstable] = 0
			}

			if dd.MaxWriteSpeedMbps != nil {
				speed[mkDiskSpeedWrite] = *dd.MaxWriteSpeedMbps
			} else {
				speed[mkDiskSpeedWrite] = 0
			}

			if dd.BurstableWriteSpeedMbps != nil {
				speed[mkDiskSpeedWriteBurstable] = *dd.BurstableWriteSpeedMbps
			} else {
				speed[mkDiskSpeedWriteBurstable] = 0
			}

			disk[mkDiskSpeed] = []interface{}{speed}
		} else {
			disk[mkDiskSpeed] = []interface{}{}
		}

		if dd.IOThread != nil {
			disk[mkDiskIOThread] = *dd.IOThread
		} else {
			disk[mkDiskIOThread] = false
		}

		if dd.SSD != nil {
			disk[mkDiskSSD] = *dd.SSD
		} else {
			disk[mkDiskSSD] = false
		}

		if dd.Discard != nil {
			disk[mkDiskDiscard] = *dd.Discard
		} else {
			disk[mkDiskDiscard] = dvDiskDiscard
		}

		if dd.Cache != nil {
			disk[mkDiskCache] = *dd.Cache
		} else {
			disk[mkDiskCache] = dvDiskCache
		}

		diskMap[di] = disk
	}

	if !isClone || len(currentDiskList) > 0 {
		orderedDiskList := utils.OrderedListFromMap(diskMap)
		err := d.Set(MkDisk, orderedDiskList)
		diags = append(diags, diag.FromErr(err)...)
	}

	return diags
}

// Update updates the disk configuration of a VM.
func Update(
	d *schema.ResourceData,
	planDisks map[string]vms.CustomStorageDevices,
	allDiskInfo vms.CustomStorageDevices,
	updateBody *vms.UpdateRequestBody,
) error {
	if d.HasChange(MkDisk) {
		for prefix, diskMap := range planDisks {
			if diskMap == nil {
				continue
			}

			for key, value := range diskMap {
				if allDiskInfo[key] == nil {
					return fmt.Errorf("missing %s device %s", prefix, key)
				}

				tmp := allDiskInfo[key]
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
				//nolint:revive
				case "ide":
					{
						// Investigate whether to support IDE mapping.
					}
				default:
					return fmt.Errorf("device prefix %s not supported", prefix)
				}
			}
		}
	}

	return nil
}
