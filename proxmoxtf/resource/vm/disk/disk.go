/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package disk

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strings"
	"unicode"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/bpg/terraform-provider-proxmox/proxmox"
	"github.com/bpg/terraform-provider-proxmox/proxmox/helpers/ptr"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/vms"
	"github.com/bpg/terraform-provider-proxmox/proxmox/ssh"
	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
	"github.com/bpg/terraform-provider-proxmox/utils"
)

// GetDiskInfoWithFileID returns the disk information for a VM.
//
// Use only when you need to get the disk information with the file ID. Otherwise, use vmConfig.StorageDevices instead.
func GetDiskInfoWithFileID(resp *vms.GetResponseData, d *schema.ResourceData) vms.CustomStorageDevices {
	storageDevices := resp.StorageDevices

	currentDisk := d.Get(MkDisk)

	diskMap := utils.MapResourcesByAttribute(currentDisk.([]any), mkDiskInterface)

	for k, v := range storageDevices {
		if v != nil && diskMap[k] != nil {
			if disk, ok := diskMap[k].(map[string]any); ok {
				if fileID, ok := disk[mkDiskFileID].(string); ok && fileID != "" {
					v.FileID = &fileID
				}
			}

			if v.Size == nil {
				v.Size = new(types.DiskSize)
			}
		}
	}

	return storageDevices
}

// UpdateClone updates disks in a cloned VM.
func UpdateClone(
	ctx context.Context,
	planDisks vms.CustomStorageDevices,
	allDiskInfo vms.CustomStorageDevices,
	vmAPI *vms.Client,
) error {
	for diskInterface, planDisk := range planDisks {
		currentDisk := allDiskInfo[diskInterface]

		if currentDisk == nil {
			diskUpdateBody := &vms.UpdateRequestBody{}
			diskUpdateBody.AddCustomStorageDevice(diskInterface, *planDisk)

			if err := vmAPI.UpdateVM(ctx, diskUpdateBody); err != nil {
				return fmt.Errorf("disk update fails: %w", err)
			}

			continue
		}

		if planDisk.Size.InMegabytes() < currentDisk.Size.InMegabytes() {
			return fmt.Errorf("disk resize failure: requested size (%s) is lower than current size (%s)",
				planDisk.Size.String(),
				currentDisk.Size.String(),
			)
		}

		// update other disk parameters
		// we have to do it before moving the disk, because the disk volume and location may change
		if currentDisk.MergeWith(*planDisk) {
			diskUpdateBody := &vms.UpdateRequestBody{}
			diskUpdateBody.AddCustomStorageDevice(diskInterface, *currentDisk)

			if err := vmAPI.UpdateVM(ctx, diskUpdateBody); err != nil {
				return fmt.Errorf("disk update fails: %w", err)
			}
		}

		moveDisk := false

		if *planDisk.DatastoreID != "" {
			fileIDParts := strings.Split(currentDisk.FileVolume, ":")
			moveDisk = *planDisk.DatastoreID != fileIDParts[0]
		}

		if moveDisk {
			deleteOriginalDisk := types.CustomBool(true)

			diskMoveBody := &vms.MoveDiskRequestBody{
				DeleteOriginalDisk: &deleteOriginalDisk,
				Disk:               diskInterface,
				TargetStorage:      *planDisk.DatastoreID,
			}

			// Note: after disk move, the actual disk volume ID will be different: both datastore id *and*
			// path in datastore will change.
			err := vmAPI.MoveVMDisk(ctx, diskMoveBody)
			if err != nil {
				return fmt.Errorf("disk move fails: %w", err)
			}
		}

		if planDisk.Size.InMegabytes() > currentDisk.Size.InMegabytes() {
			diskResizeBody := &vms.ResizeDiskRequestBody{
				Disk: diskInterface,
				Size: *planDisk.Size,
			}

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
	_ *schema.Resource,
	disks []any,
) (vms.CustomStorageDevices, error) {
	var diskDevices []any

	if disks != nil {
		diskDevices = disks
	} else {
		diskDevices = d.Get(MkDisk).([]any)
	}

	diskDeviceObjects := vms.CustomStorageDevices{}

	for _, diskEntry := range diskDevices {
		diskDevice := &vms.CustomStorageDevice{}

		block := diskEntry.(map[string]any)
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
		importFrom, _ := block[mkDiskImportFrom].(string)
		ioThread := types.CustomBool(block[mkDiskIOThread].(bool))
		replicate := types.CustomBool(block[mkDiskReplicate].(bool))
		serial := block[mkDiskSerial].(string)
		ssd := types.CustomBool(block[mkDiskSSD].(bool))

		// Handle disk size: prefer disk_size (string with units) over deprecated size (int in GB)
		var diskSize *types.DiskSize

		if diskSizeStr, ok := block[MkDiskSizeStr].(string); ok && diskSizeStr != "" {
			ds, err := types.ParseDiskSize(diskSizeStr)
			if err != nil {
				return diskDeviceObjects, fmt.Errorf("invalid disk_size %q: %w", diskSizeStr, err)
			}

			diskSize = &ds
		} else {
			size, _ := block[mkDiskSize].(int)
			diskSize = types.DiskSizeFromGigabytes(int64(size))
		}

		// get speed block directly from the current disk entry
		var speedBlock map[string]any

		if speedList, ok := block[mkDiskSpeed].([]any); ok && len(speedList) > 0 {
			if sb, ok := speedList[0].(map[string]any); ok {
				speedBlock = sb
			}
		}

		if pathInDatastore != "" {
			if datastoreID != "" {
				diskDevice.FileVolume = fmt.Sprintf("%s:%s", datastoreID, pathInDatastore)
			} else {
				// FileVolume is the absolute path in the host filesystem
				diskDevice.FileVolume = pathInDatastore
			}
		} else {
			diskDevice.FileVolume = fmt.Sprintf("%s:%d", datastoreID, diskSize.InGigabytes())
		}

		diskDevice.AIO = &aio
		diskDevice.Backup = &backup
		diskDevice.Cache = &cache
		diskDevice.DatastoreID = &datastoreID
		diskDevice.Discard = &discard
		diskDevice.FileID = &fileID
		diskDevice.ImportFrom = &importFrom
		diskDevice.Replicate = &replicate
		diskDevice.Serial = &serial
		diskDevice.Size = diskSize

		if fileFormat != "" {
			diskDevice.Format = &fileFormat
		}

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

		if !slices.Contains(vms.StorageInterfaces, DigitPrefix(diskInterface)) {
			errorMsg := fmt.Sprintf(
				"Defined disk interface not supported. Interface was %s, but only %v are supported",
				diskInterface, vms.StorageInterfaces,
			)

			return diskDeviceObjects, errors.New(errorMsg)
		}

		diskDeviceObjects[diskInterface] = diskDevice
	}

	return diskDeviceObjects, nil
}

// CreateCustomDisks creates custom disks for a VM.
func CreateCustomDisks(
	ctx context.Context,
	client proxmox.Client,
	nodeName string,
	vmID int,
	storageDevices vms.CustomStorageDevices,
) diag.Diagnostics {
	for iface, disk := range storageDevices {
		if disk != nil && disk.FileID != nil && *disk.FileID != "" {
			// only custom disks with defined file ID
			err := createCustomDisk(ctx, client, nodeName, vmID, iface, *disk)
			if err != nil {
				return diag.FromErr(err)
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
	iface string,
	disk vms.CustomStorageDevice,
) error {
	// use "old default" specifically here.
	fileFormat := "qcow2"
	if disk.Format != nil && *disk.Format != "" {
		fileFormat = *disk.Format
	}

	//nolint:lll
	commands := []string{
		`set -e`,
		ssh.TrySudo,
		`file_id=` + *disk.FileID,
		`file_format=` + fileFormat,
		`datastore_id_target=` + *disk.DatastoreID,
		fmt.Sprintf(`vm_id=%d`, vmID),
		fmt.Sprintf(`disk_options=%s`, disk.EncodeOptions()),
		fmt.Sprintf(`disk_interface=%s`, iface),
		`source_image=$(try_sudo /usr/sbin/pvesm path $file_id)`,
		`imported_disk=$(try_sudo /usr/sbin/qm disk import $vm_id $source_image $datastore_id_target -format $file_format | grep unused0 | cut -d : -f 3 | cut -d \' -f 1)`,
		`disk_id=${datastore_id_target}:$imported_disk,$disk_options`,
		`try_sudo /usr/sbin/qm set $vm_id -${disk_interface} $disk_id`,
	}

	out, err := client.SSH().ExecuteNodeCommands(ctx, nodeName, commands)
	if err != nil {
		if matches, e := regexp.Match(`pvesm: .* not found`, out); e == nil && matches {
			err = ssh.NewErrUserHasNoPermission(client.SSH().Username())
		}

		return fmt.Errorf("creating custom disk: %w", err)
	}

	tflog.Debug(ctx, "vmCreateCustomDisks: commands", map[string]any{
		"output": string(out),
	})

	err = client.Node(nodeName).VM(vmID).ResizeVMDisk(ctx, &vms.ResizeDiskRequestBody{
		Disk: iface,
		Size: *disk.Size,
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
	currentDiskList := d.Get(MkDisk).([]any)
	diskMap := map[string]any{}

	var diags diag.Diagnostics

	for di, dd := range diskObjects {
		if dd == nil || dd.FileVolume == "none" {
			continue
		}

		if dd.IsCloudInitDrive(vmID) {
			continue
		}

		disk := map[string]any{}

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

		// note that PVE does not return back the 'import-from' attribute for the disks that are imported,
		// but we'll keep it here for consistency. the actual value is set later down
		if dd.ImportFrom != nil {
			disk[mkDiskImportFrom] = dd.ImportFrom
		}

		disk[mkDiskInterface] = di
		disk[mkDiskSize] = dd.Size.InGigabytes()
		disk[MkDiskSizeStr] = dd.Size.String()

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

		if dd.Serial != nil {
			disk[mkDiskSerial] = *dd.Serial
		} else {
			disk[mkDiskSerial] = ""
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
			speed := map[string]any{}

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

			disk[mkDiskSpeed] = []any{speed}
		} else {
			disk[mkDiskSpeed] = []any{}
		}

		diskMap[di] = disk
	}

	if !isClone || len(currentDiskList) > 0 {
		var diskList []any

		if len(currentDiskList) > 0 {
			currentDiskMap := utils.MapResourcesByAttribute(currentDiskList, mkDiskInterface)
			// copy import_from and size from the current disk if it exists
			for k, v := range currentDiskMap {
				if disk, ok := v.(map[string]any); ok {
					if _, exists := diskMap[k]; exists {
						if importFrom, ok := disk[mkDiskImportFrom].(string); ok && importFrom != "" {
							diskMap[k].(map[string]any)[mkDiskImportFrom] = importFrom
						}
						// preserve size from state when API returns zero size (for disks with import_from or file_id)
						if currentSize, ok := disk[mkDiskSize].(int); ok && currentSize > 0 {
							if apiSize, ok := diskMap[k].(map[string]any)[mkDiskSize].(int64); ok && apiSize == 0 {
								diskMap[k].(map[string]any)[mkDiskSize] = currentSize
							}
						}
						// Handle disk_size: only preserve if user configured it
						currentDiskSize, hasDiskSize := disk[MkDiskSizeStr].(string)
						if hasDiskSize && currentDiskSize != "" {
							// User configured disk_size, preserve from state if API returns empty
							if apiDiskSize, ok := diskMap[k].(map[string]any)[MkDiskSizeStr].(string); ok && apiDiskSize == "" {
								diskMap[k].(map[string]any)[MkDiskSizeStr] = currentDiskSize
							}
						} else {
							// User didn't configure disk_size (used deprecated size instead), clear it to avoid diff
							diskMap[k].(map[string]any)[MkDiskSizeStr] = ""
						}
					}
				}
			}

			disks := utils.ListResourcesAttributeValue(currentDiskList, mkDiskInterface)
			diskList = utils.OrderedListFromMapByKeyValues(diskMap, disks)
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
	planDisks vms.CustomStorageDevices,
	currentDisks vms.CustomStorageDevices,
	updateBody *vms.UpdateRequestBody,
) (bool, bool, error) {
	rebootRequired := false

	if d.HasChange(MkDisk) {
		for iface, disk := range planDisks {
			var tmp *vms.CustomStorageDevice

			switch {
			case currentDisks[iface] == nil && disk != nil:
				if disk.FileID != nil && *disk.FileID != "" {
					// only disks with defined file ID are custom image disks that need to be created via import.
					err := createCustomDisk(ctx, client, nodeName, vmID, iface, *disk)
					if err != nil {
						return false, false, fmt.Errorf("creating custom disk: %w", err)
					}
				} else {
					// otherwise this is a blank disk that can be added directly via update API
					tmp = disk
				}
			case currentDisks[iface] != nil:
				// Check if the disk has actually changed before updating
				if currentDisks[iface].Equals(disk) {
					// Disk hasn't changed, skip update
					continue
				}
				// update existing disk
				tmp = currentDisks[iface]
			default:
				// something went wrong
				return false, false, fmt.Errorf("missing device %s", iface)
			}

			if tmp == nil || disk == nil {
				continue
			}

			if !ptr.Eq(tmp.AIO, disk.AIO) {
				rebootRequired = true
				tmp.AIO = disk.AIO
			}

			// Never re-import existing disks - import_from is only for initial disk creation.
			// See https://github.com/bpg/terraform-provider-proxmox/issues/2385

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
			tmp.Serial = disk.Serial
			tmp.SSD = disk.SSD

			updateBody.AddCustomStorageDevice(iface, *tmp)
		}
	}

	return false, rebootRequired, nil
}
