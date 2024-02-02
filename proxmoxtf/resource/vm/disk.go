package vm

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/bpg/terraform-provider-proxmox/proxmox"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/vms"
	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
	"github.com/bpg/terraform-provider-proxmox/proxmoxtf"
	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/resource/validator"
	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/structure"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func diskSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Description: "The disk devices",
		Optional:    true,
		DefaultFunc: func() (interface{}, error) {
			return []interface{}{
				map[string]interface{}{
					mkDiskDatastoreID:     dvDiskDatastoreID,
					mkDiskPathInDatastore: nil,
					mkDiskFileID:          dvDiskFileID,
					mkDiskInterface:       dvDiskInterface,
					mkDiskSize:            dvDiskSize,
					mkDiskIOThread:        dvDiskIOThread,
					mkDiskSSD:             dvDiskSSD,
					mkDiskDiscard:         dvDiskDiscard,
					mkDiskCache:           dvDiskCache,
				},
			}, nil
		},
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				mkDiskInterface: {
					Type:        schema.TypeString,
					Description: "The datastore name",
					Required:    true,
				},
				mkDiskDatastoreID: {
					Type:        schema.TypeString,
					Description: "The datastore id",
					Optional:    true,
					Default:     dvDiskDatastoreID,
				},
				mkDiskPathInDatastore: {
					Type:        schema.TypeString,
					Description: "The in-datastore path to disk image",
					Computed:    true,
					Optional:    true,
					Default:     nil,
				},
				mkDiskFileFormat: {
					Type:             schema.TypeString,
					Description:      "The file format",
					Optional:         true,
					ForceNew:         true,
					Computed:         true,
					ValidateDiagFunc: validator.FileFormat(),
				},
				mkDiskFileID: {
					Type:             schema.TypeString,
					Description:      "The file id for a disk image",
					Optional:         true,
					ForceNew:         true,
					Default:          dvDiskFileID,
					ValidateDiagFunc: validator.FileID(),
				},
				mkDiskSize: {
					Type:             schema.TypeInt,
					Description:      "The disk size in gigabytes",
					Optional:         true,
					Default:          dvDiskSize,
					ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(1)),
				},
				mkDiskIOThread: {
					Type:        schema.TypeBool,
					Description: "Whether to use iothreads for this disk drive",
					Optional:    true,
					Default:     dvDiskIOThread,
				},
				mkDiskSSD: {
					Type:        schema.TypeBool,
					Description: "Whether to use ssd for this disk drive",
					Optional:    true,
					Default:     dvDiskSSD,
				},
				mkDiskDiscard: {
					Type:        schema.TypeString,
					Description: "Whether to pass discard/trim requests to the underlying storage.",
					Optional:    true,
					Default:     dvDiskDiscard,
				},
				mkDiskCache: {
					Type:        schema.TypeString,
					Description: "The drive’s cache mode",
					Optional:    true,
					Default:     dvDiskCache,
					ValidateDiagFunc: validation.ToDiagFunc(
						validation.StringInSlice([]string{
							"none",
							"writethrough",
							"writeback",
							"unsafe",
							"directsync",
						}, false),
					),
				},
				mkDiskSpeed: {
					Type:        schema.TypeList,
					Description: "The speed limits",
					Optional:    true,
					DefaultFunc: func() (interface{}, error) {
						return []interface{}{
							map[string]interface{}{
								mkDiskSpeedRead:           dvDiskSpeedRead,
								mkDiskSpeedReadBurstable:  dvDiskSpeedReadBurstable,
								mkDiskSpeedWrite:          dvDiskSpeedWrite,
								mkDiskSpeedWriteBurstable: dvDiskSpeedWriteBurstable,
							},
						}, nil
					},
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							mkDiskSpeedRead: {
								Type:        schema.TypeInt,
								Description: "The maximum read speed in megabytes per second",
								Optional:    true,
								Default:     dvDiskSpeedRead,
							},
							mkDiskSpeedReadBurstable: {
								Type:        schema.TypeInt,
								Description: "The maximum burstable read speed in megabytes per second",
								Optional:    true,
								Default:     dvDiskSpeedReadBurstable,
							},
							mkDiskSpeedWrite: {
								Type:        schema.TypeInt,
								Description: "The maximum write speed in megabytes per second",
								Optional:    true,
								Default:     dvDiskSpeedWrite,
							},
							mkDiskSpeedWriteBurstable: {
								Type:        schema.TypeInt,
								Description: "The maximum burstable write speed in megabytes per second",
								Optional:    true,
								Default:     dvDiskSpeedWriteBurstable,
							},
						},
					},
					MaxItems: 1,
					MinItems: 0,
				},
			},
		},
	}
}

func updateDisk1(
	ctx context.Context, vmConfig *vms.GetResponseData, d *schema.ResourceData, vmAPI *vms.Client,
) (map[string]*vms.CustomStorageDevice, error) {
	allDiskInfo := getDiskInfo(vmConfig, d)

	diskDeviceObjects, e := vmGetDiskDeviceObjects(d, nil)
	if e != nil {
		return nil, e
	}

	disk := d.Get(mkDisk).([]interface{})
	for i := range disk {
		diskBlock := disk[i].(map[string]interface{})
		diskInterface := diskBlock[mkDiskInterface].(string)
		dataStoreID := diskBlock[mkDiskDatastoreID].(string)
		diskSize := int64(diskBlock[mkDiskSize].(int))
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
				return nil, e
			}

			continue
		}

		if diskSize < currentDiskInfo.Size.InGigabytes() {
			return nil, fmt.Errorf("disk resize fails requests size (%dG) is lower than current size (%s)",
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
			moveDiskTimeout := d.Get(mkTimeoutMoveDisk).(int)

			e = vmAPI.MoveVMDisk(ctx, diskMoveBody, moveDiskTimeout)
			if e != nil {
				return nil, e
			}
		}

		if diskSize > currentDiskInfo.Size.InGigabytes() {
			e = vmAPI.ResizeVMDisk(ctx, diskResizeBody)
			if e != nil {
				return nil, e
			}
		}
	}

	return allDiskInfo, nil
}

func vmCreateCustomDisks(ctx context.Context, d *schema.ResourceData, m interface{}) error {
	vmID, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	// Determine the ID of the next disk.
	disk := d.Get(mkDisk).([]interface{})
	diskCount := 0

	for _, d := range disk {
		block := d.(map[string]interface{})
		fileID, _ := block[mkDiskFileID].(string)

		if fileID == "" {
			diskCount++
		}
	}

	// Retrieve some information about the disk schema.
	resourceSchema := VM().Schema
	diskSchemaElem := resourceSchema[mkDisk].Elem
	diskSchemaResource := diskSchemaElem.(*schema.Resource)
	diskSpeedResource := diskSchemaResource.Schema[mkDiskSpeed]

	// Generate the commands required to import the specified disks.
	commands := []string{}
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
				return err
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
			return err
		}

		nodeName := d.Get(mkNodeName).(string)

		out, err := api.SSH().ExecuteNodeCommands(ctx, nodeName, commands)
		if err != nil {
			if strings.Contains(err.Error(), "pvesm: not found") {
				return fmt.Errorf("The configured SSH user '%s' does not have the required permissions to import disks. "+
					"Make sure `sudo` is installed and the user is a member of sudoers.", api.SSH().Username())
			}

			return err
		}

		tflog.Debug(ctx, "vmCreateCustomDisks", map[string]interface{}{
			"output": string(out),
		})
	}

	return nil
}

func vmGetDiskDeviceObjects(
	d *schema.ResourceData,
	disks []interface{},
) (map[string]map[string]vms.CustomStorageDevice, error) {
	var diskDevice []interface{}

	if disks != nil {
		diskDevice = disks
	} else {
		diskDevice = d.Get(mkDisk).([]interface{})
	}

	diskDeviceObjects := map[string]map[string]vms.CustomStorageDevice{}
	resource := VM()

	for _, diskEntry := range diskDevice {
		diskDevice := vms.CustomStorageDevice{
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
			[]string{mkDisk, mkDiskSpeed},
			0,
			false,
		)
		if err != nil {
			return diskDeviceObjects, err
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

func readDisk1(ctx context.Context, d *schema.ResourceData,
	vmConfig *vms.GetResponseData, vmID int, api proxmox.Client, nodeName string, clone []interface{},
) diag.Diagnostics {
	currentDiskList := d.Get(mkDisk).([]interface{})
	diskMap := map[string]interface{}{}
	diskObjects := getDiskInfo(vmConfig, d)

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
				volume, err := api.Node(nodeName).Storage(datastoreID).GetDatastoreFile(ctx, dd.FileVolume)
				if err != nil {
					diags = append(diags, diag.FromErr(err)...)
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

	if len(clone) == 0 || len(currentDiskList) > 0 {
		orderedDiskList := orderedListFromMap(diskMap)
		err := d.Set(mkDisk, orderedDiskList)
		diags = append(diags, diag.FromErr(err)...)
	}

	return diags
}

func updateDisk(d *schema.ResourceData, vmConfig *vms.GetResponseData, updateBody *vms.UpdateRequestBody) error {
	// Prepare the new disk device configuration.
	if !d.HasChange(mkDisk) {
		return nil
	}

	diskDeviceObjects, err := vmGetDiskDeviceObjects(d, nil)
	if err != nil {
		return err
	}

	diskDeviceInfo := getDiskInfo(vmConfig, d)

	for prefix, diskMap := range diskDeviceObjects {
		if diskMap == nil {
			continue
		}

		for key, value := range diskMap {
			if diskDeviceInfo[key] == nil {
				// TODO: create a new disk here
				return fmt.Errorf("missing %s device %s", prefix, key)
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
				return fmt.Errorf("device prefix %s not supported", prefix)
			}
		}
	}

	return nil
}
