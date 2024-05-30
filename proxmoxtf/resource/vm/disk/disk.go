package disk

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/bpg/terraform-provider-proxmox/proxmox"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/vms"
	"github.com/bpg/terraform-provider-proxmox/proxmox/ssh"
	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/structure"
	"github.com/bpg/terraform-provider-proxmox/utils"
)

const supportedDiskInterfaces = "virtio, sata, scsi, ide"

// GetInfo returns the disk information for a VM.
func GetInfo(resp *vms.GetResponseData, d *schema.ResourceData) vms.CustomStorageDevices {
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

	currentDisk := d.Get(MkDisk)

	diskMap := utils.MapResourcesByAttribute(currentDisk.([]interface{}), mkDiskInterface)

	for k, v := range storageDevices {
		if v != nil && diskMap[k] != nil {
			if d, ok := diskMap[k].(map[string]interface{}); ok {
				if fileID, ok := d[mkDiskFileID].(string); ok && fileID != "" {
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

			case "ide":
				if diskUpdateBody.IDEDevices == nil {
					diskUpdateBody.IDEDevices = vms.CustomStorageDevices{}
				}

				diskUpdateBody.IDEDevices[diskInterface] = configuredDiskInfo
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

		if moveDisk {
			err := vmAPI.MoveVMDisk(ctx, diskMoveBody)
			if err != nil {
				return fmt.Errorf("disk move fails: %w", err)
			}
		}

		if diskSize > currentDiskInfo.Size.InGigabytes() {
			err := vmAPI.ResizeVMDisk(ctx, diskResizeBody)
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

		aio := block[mkDiskAIO].(string)
		backup := types.CustomBool(block[mkDiskBackup].(bool))
		cache := block[mkDiskCache].(string)
		discard := block[mkDiskDiscard].(string)
		diskInterface, _ := block[mkDiskInterface].(string)
		fileFormat, _ := block[mkDiskFileFormat].(string)
		fileID, _ := block[mkDiskFileID].(string)
		ioThread := types.CustomBool(block[mkDiskIOThread].(bool))
		replicate := types.CustomBool(block[mkDiskReplicate].(bool))
		size, _ := block[mkDiskSize].(int)
		ssd := types.CustomBool(block[mkDiskSSD].(bool))

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

		// Explicitly disable the disk, so it won't be encoded in "Create" or "Update" operations.
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

		diskDevice.AIO = &aio
		diskDevice.Backup = &backup
		diskDevice.Cache = &cache
		diskDevice.DatastoreID = &datastoreID
		diskDevice.Discard = &discard
		diskDevice.FileID = &fileID
		diskDevice.Format = &fileFormat
		diskDevice.Interface = &diskInterface
		diskDevice.Replicate = &replicate
		diskDevice.Size = types.DiskSizeFromGigabytes(int64(size))

		if !strings.HasPrefix(diskInterface, "virtio") {
			diskDevice.SSD = &ssd
		}

		if !strings.HasPrefix(diskInterface, "sata") && !strings.HasPrefix(diskInterface, "ide") {
			diskDevice.IOThread = &ioThread
		}

		if len(speedBlock) > 0 {
			iopsRead := speedBlock[mkDiskIopsRead].(int)
			iopsReadBurstable := speedBlock[mkDiskIopsReadBurstable].(int)
			iopsWrite := speedBlock[mkDiskIopsWrite].(int)
			iopsWriteBurstable := speedBlock[mkDiskIopsWriteBurstable].(int)
			speedLimitRead := speedBlock[mkDiskSpeedRead].(int)
			speedLimitReadBurstable := speedBlock[mkDiskSpeedReadBurstable].(int)
			speedLimitWrite := speedBlock[mkDiskSpeedWrite].(int)
			speedLimitWriteBurstable := speedBlock[mkDiskSpeedWriteBurstable].(int)

			if iopsRead > 0 {
				diskDevice.IopsRead = &iopsRead
			}

			if iopsReadBurstable > 0 {
				diskDevice.MaxIopsRead = &iopsReadBurstable
			}

			if iopsWrite > 0 {
				diskDevice.IopsWrite = &iopsWrite
			}

			if iopsWriteBurstable > 0 {
				diskDevice.MaxIopsWrite = &iopsWriteBurstable
			}

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
		if !strings.Contains(supportedDiskInterfaces, baseDiskInterface) {
			errorMsg := fmt.Sprintf(
				"Defined disk interface not supported. Interface was %s, but only %s are supported",
				diskInterface, supportedDiskInterfaces,
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
	client proxmox.Client,
	nodeName string,
	vmID int,
	storageDevices map[string]vms.CustomStorageDevices,
) diag.Diagnostics {
	for _, diskDevice := range storageDevices {
		for _, disk := range diskDevice {
			if disk != nil && disk.FileID != nil && *disk.FileID != "" {
				// only custom disks with defined file ID
				err := createCustomDisk(ctx, client, nodeName, vmID, *disk)
				if err != nil {
					return diag.FromErr(err)
				}
			}
		}
	}

	return nil
}

func createCustomDisk(
	ctx context.Context,
	client proxmox.Client,
	nodeName string,
	vmID int,
	d vms.CustomStorageDevice,
) error {
	fileFormat := dvDiskFileFormat
	if d.Format != nil && *d.Format != "" {
		fileFormat = *d.Format
	}

	//nolint:lll
	commands := []string{
		`set -e`,
		ssh.TrySudo,
		fmt.Sprintf(`file_id="%s"`, *d.FileID),
		fmt.Sprintf(`file_format="%s"`, fileFormat),
		fmt.Sprintf(`datastore_id_target="%s"`, *d.DatastoreID),
		fmt.Sprintf(`vm_id="%d"`, vmID),
		fmt.Sprintf(`disk_options="%s"`, d.EncodeOptions()),
		fmt.Sprintf(`disk_interface="%s"`, *d.Interface),
		`source_image=$(try_sudo "pvesm path $file_id")`,
		`imported_disk="$(try_sudo "qm importdisk $vm_id $source_image $datastore_id_target -format $file_format" | grep "unused0" | cut -d ":" -f 3 | cut -d "'" -f 1)"`,
		`disk_id="${datastore_id_target}:$imported_disk,${disk_options}"`,
		`try_sudo "qm set $vm_id -${disk_interface} $disk_id"`,
	}

	out, err := client.SSH().ExecuteNodeCommands(ctx, nodeName, commands)
	if err != nil {
		if matches, e := regexp.Match(`pvesm: .* not found`, out); e == nil && matches {
			err = ssh.NewErrUserHasNoPermission(client.SSH().Username())
		}

		return fmt.Errorf("creating custom disk: %w", err)
	}

	tflog.Debug(ctx, "vmCreateCustomDisks: commands", map[string]interface{}{
		"output": string(out),
	})

	err = client.Node(nodeName).VM(vmID).ResizeVMDisk(ctx, &vms.ResizeDiskRequestBody{
		Disk: *d.Interface,
		Size: *d.Size,
	})
	if err != nil {
		return fmt.Errorf("resizing disk: %w", err)
	}

	return nil
}

// Read reads the disk configuration of a VM.
func Read(
	ctx context.Context,
	d *schema.ResourceData,
	diskObjects vms.CustomStorageDevices,
	vmID int,
	client proxmox.Client,
	nodeName string,
	isClone bool,
) diag.Diagnostics {
	currentDiskList := d.Get(MkDisk).([]interface{})
	diskMap := map[string]interface{}{}

	var diags diag.Diagnostics

	for di, dd := range diskObjects {
		if dd == nil || dd.FileVolume == "none" {
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
				volume, e := client.Node(nodeName).Storage(datastoreID).GetDatastoreFile(ctx, dd.FileVolume)
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

		if dd.AIO != nil {
			disk[mkDiskAIO] = *dd.AIO
		} else {
			disk[mkDiskAIO] = dvDiskAIO
		}

		if dd.Backup != nil {
			disk[mkDiskBackup] = *dd.Backup
		} else {
			disk[mkDiskBackup] = true
		}

		if dd.IOThread != nil {
			disk[mkDiskIOThread] = *dd.IOThread
		} else {
			disk[mkDiskIOThread] = false
		}

		if dd.Replicate != nil {
			disk[mkDiskReplicate] = *dd.Replicate
		} else {
			disk[mkDiskReplicate] = true
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

		if dd.IopsRead != nil ||
			dd.MaxIopsRead != nil ||
			dd.IopsWrite != nil ||
			dd.MaxIopsWrite != nil ||
			dd.BurstableReadSpeedMbps != nil ||
			dd.BurstableWriteSpeedMbps != nil ||
			dd.MaxReadSpeedMbps != nil ||
			dd.MaxWriteSpeedMbps != nil {
			speed := map[string]interface{}{}

			if dd.IopsRead != nil {
				speed[mkDiskIopsRead] = *dd.IopsRead
			} else {
				speed[mkDiskIopsRead] = 0
			}

			if dd.MaxIopsRead != nil {
				speed[mkDiskIopsReadBurstable] = *dd.MaxIopsRead
			} else {
				speed[mkDiskIopsReadBurstable] = 0
			}

			if dd.IopsWrite != nil {
				speed[mkDiskIopsWrite] = *dd.IopsWrite
			} else {
				speed[mkDiskIopsWrite] = 0
			}

			if dd.MaxIopsWrite != nil {
				speed[mkDiskIopsWriteBurstable] = *dd.MaxIopsWrite
			} else {
				speed[mkDiskIopsWriteBurstable] = 0
			}

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

		diskMap[di] = disk
	}

	if !isClone || len(currentDiskList) > 0 {
		var diskList []interface{}

		if len(currentDiskList) > 0 {
			interfaces := utils.ListResourcesAttributeValue(currentDiskList, mkDiskInterface)
			diskList = utils.OrderedListFromMapByKeyValues(diskMap, interfaces)
		} else {
			diskList = utils.OrderedListFromMap(diskMap)
		}

		err := d.Set(MkDisk, diskList)
		diags = append(diags, diag.FromErr(err)...)
	}

	return diags
}

// Update updates the disk configuration of a VM.
func Update(
	ctx context.Context,
	client proxmox.Client,
	nodeName string,
	vmID int,
	d *schema.ResourceData,
	planDisks map[string]vms.CustomStorageDevices,
	currentDisks vms.CustomStorageDevices,
	updateBody *vms.UpdateRequestBody,
) (bool, error) {
	rebootRequired := false

	if d.HasChange(MkDisk) {
		for prefix, diskMap := range planDisks {
			if diskMap == nil {
				continue
			}

			for key, disk := range diskMap {
				var tmp *vms.CustomStorageDevice

				switch {
				case currentDisks[key] == nil && disk != nil:
					if disk.FileID != nil && *disk.FileID != "" {
						// only disks with defined file ID are custom image disks that need to be created via import.
						err := createCustomDisk(ctx, client, nodeName, vmID, *disk)
						if err != nil {
							return false, fmt.Errorf("creating custom disk: %w", err)
						}
					} else {
						// otherwise this is a blank disk that can be added directly via update API
						tmp = disk
					}
				case currentDisks[key] != nil:
					// update existing disk
					tmp = currentDisks[key]
				default:
					// something went wrong
					return false, fmt.Errorf("missing %s device %s", prefix, key)
				}

				if tmp == nil || disk == nil {
					continue
				}

				if tmp.AIO != disk.AIO {
					rebootRequired = true
					tmp.AIO = disk.AIO
				}

				tmp.Backup = disk.Backup
				tmp.BurstableReadSpeedMbps = disk.BurstableReadSpeedMbps
				tmp.BurstableWriteSpeedMbps = disk.BurstableWriteSpeedMbps
				tmp.Cache = disk.Cache
				tmp.Discard = disk.Discard
				tmp.IOThread = disk.IOThread
				tmp.IopsRead = disk.IopsRead
				tmp.IopsWrite = disk.IopsWrite
				tmp.MaxIopsRead = disk.MaxIopsRead
				tmp.MaxIopsWrite = disk.MaxIopsWrite
				tmp.MaxReadSpeedMbps = disk.MaxReadSpeedMbps
				tmp.MaxWriteSpeedMbps = disk.MaxWriteSpeedMbps
				tmp.Replicate = disk.Replicate
				tmp.SSD = disk.SSD

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
						if updateBody.IDEDevices == nil {
							updateBody.IDEDevices = vms.CustomStorageDevices{}
						}

						updateBody.IDEDevices[key] = tmp
					}
				default:
					return false, fmt.Errorf("device prefix %s not supported", prefix)
				}
			}
		}
	}

	return rebootRequired, nil
}
