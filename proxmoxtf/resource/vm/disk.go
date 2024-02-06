package vm

import (
	"context"
	"fmt"
	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/resource/ssh"
	"regexp"
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

const (
	mkDisk                    = "disk"
	mkDiskInterface           = "interface"
	mkDiskDatastoreID         = "datastore_id"
	mkDiskPathInDatastore     = "path_in_datastore"
	mkDiskFileFormat          = "file_format"
	mkDiskFileID              = "file_id"
	mkDiskSize                = "size"
	mkDiskIOThread            = "iothread"
	mkDiskSSD                 = "ssd"
	mkDiskDiscard             = "discard"
	mkDiskCache               = "cache"
	mkDiskSpeed               = "speed"
	mkDiskSpeedRead           = "read"
	mkDiskSpeedReadBurstable  = "read_burstable"
	mkDiskSpeedWrite          = "write"
	mkDiskSpeedWriteBurstable = "write_burstable"

	dvDiskInterface           = "scsi0"
	dvDiskDatastoreID         = "local-lvm"
	dvDiskFileFormat          = "qcow2"
	dvDiskFileID              = ""
	dvDiskSize                = 8
	dvDiskIOThread            = false
	dvDiskSSD                 = false
	dvDiskDiscard             = "ignore"
	dvDiskCache               = "none"
	dvDiskSpeedRead           = 0
	dvDiskSpeedReadBurstable  = 0
	dvDiskSpeedWrite          = 0
	dvDiskSpeedWriteBurstable = 0
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

// called from vmCreateClone
func createDisks(
	ctx context.Context, vmConfig *vms.GetResponseData, d *schema.ResourceData, vmAPI *vms.Client,
) (vms.CustomStorageDevices, error) {
	// this is what VM has at the moment: map of interface name (virtio1) -> disk object
	currentDisks := populateFileIDs(mapStorageDevices(vmConfig), d)

	// map of interface name (virtio1) -> disk object
	planDisks, e := getStorageDevicesFromResource(d)
	if e != nil {
		return nil, e
	}

	for iface, planDisk := range planDisks {
		currentDisk := currentDisks[iface]

		// create disks that are not present in the current configuration
		if currentDisk == nil {
			err := createDisk(ctx, planDisk, vmAPI)
			if err != nil {
				return nil, err
			}

			continue
		}

		// disk is present, i.e. when cloning a template, but we need to check if it needs to be moved or resized

		timeoutSec := d.Get(mkTimeoutMoveDisk).(int)

		err := resizeDiskIfRequired(ctx, currentDisk, planDisk, vmAPI, timeoutSec)
		if err != nil {
			return nil, err
		}

		err = moveDiskIfRequired(ctx, currentDisk, planDisk, vmAPI, timeoutSec)
		if err != nil {
			return nil, err
		}
	}

	return currentDisks, nil
}

func resizeDiskIfRequired(
	ctx context.Context,
	currentDisk *vms.CustomStorageDevice, planDisk *vms.CustomStorageDevice,
	vmAPI *vms.Client, timeoutSec int,
) error {
	if planDisk.Size.InGigabytes() < currentDisk.Size.InGigabytes() {
		return fmt.Errorf("the planned disk size (%dG) is lower than the current size (%s)",
			planDisk.Size.InGigabytes(),
			*currentDisk.Size,
		)
	}

	if planDisk.Size.InGigabytes() > currentDisk.Size.InGigabytes() {
		diskResizeBody := &vms.ResizeDiskRequestBody{
			Disk: *planDisk.Interface,
			Size: *planDisk.Size,
		}

		err := vmAPI.ResizeVMDisk(ctx, diskResizeBody, timeoutSec)
		if err != nil {
			return err
		}
	}

	return nil
}

func moveDiskIfRequired(
	ctx context.Context,
	currentDisk *vms.CustomStorageDevice, planDisk *vms.CustomStorageDevice,
	vmAPI *vms.Client, timeoutSec int,
) error {
	needToMove := false

	if *planDisk.ID != "" {
		fileIDParts := strings.Split(currentDisk.FileVolume, ":")
		needToMove = *planDisk.ID != fileIDParts[0]
	}

	if needToMove {
		diskMoveBody := &vms.MoveDiskRequestBody{
			DeleteOriginalDisk: types.CustomBool(true).Pointer(),
			Disk:               *planDisk.Interface,
			TargetStorage:      *planDisk.ID,
		}

		err := vmAPI.MoveVMDisk(ctx, diskMoveBody, timeoutSec)
		if err != nil {
			return err
		}
	}

	return nil
}

func createDisk(ctx context.Context, disk *vms.CustomStorageDevice, vmAPI *vms.Client) error {
	addToDevices := func(ds vms.CustomStorageDevices, disk *vms.CustomStorageDevice) vms.CustomStorageDevices {
		if ds == nil {
			ds = vms.CustomStorageDevices{}
		}

		ds[*disk.Interface] = disk

		return ds
	}

	diskUpdateBody := &vms.UpdateRequestBody{}

	switch disk.StorageInterface() {
	case "virtio":
		diskUpdateBody.VirtualIODevices = addToDevices(diskUpdateBody.VirtualIODevices, disk)
	case "sata":
		diskUpdateBody.SATADevices = addToDevices(diskUpdateBody.SATADevices, disk)
	case "scsi":
		diskUpdateBody.SCSIDevices = addToDevices(diskUpdateBody.SCSIDevices, disk)
	}

	err := vmAPI.UpdateVM(ctx, diskUpdateBody)
	if err != nil {
		return err
	}

	return nil
}

func vmImportCustomDisks(ctx context.Context, d *schema.ResourceData, m interface{}) error {
	vmID, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	planDisks, err := getStorageDevicesFromResource(d)
	if err != nil {
		return err
	}

	diskCount := 0

	for _, d := range planDisks {
		if *d.FileID == "" {
			diskCount++
		}
	}

	// Generate the commands required to import the specified disks.
	commands := []string{}
	importedDiskCount := 0

	for _, d := range planDisks {
		if *d.FileID == "" {
			continue
		}

		diskOptions := d.EncodeOptions()
		if diskOptions != "" {
			diskOptions = "," + diskOptions
		}

		filePathTmp := fmt.Sprintf(
			"/tmp/vm-%d-disk-%d.%s",
			vmID,
			diskCount+importedDiskCount,
			*d.Format,
		)

		//nolint:lll
		commands = append(
			commands,
			`set -e`,
			ssh.TrySudo,
			fmt.Sprintf(`file_id="%s"`, *d.FileID),
			fmt.Sprintf(`file_format="%s"`, *d.Format),
			fmt.Sprintf(`datastore_id_target="%s"`, *d.ID),
			fmt.Sprintf(`disk_options="%s"`, diskOptions),
			fmt.Sprintf(`disk_size="%d"`, d.Size.InGigabytes()),
			fmt.Sprintf(`disk_interface="%s"`, *d.Interface),
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
			if matches, e := regexp.Match(`pvesm: .* not found`, out); e == nil && matches {
				return ssh.NewErrSSHUserNoPermission(api.SSH().Username())
			}

			return err
		}

		tflog.Debug(ctx, "vmCreateCustomDisks", map[string]interface{}{
			"output": string(out),
		})
	}

	return nil
}

func getStorageDevicesFromResource(d *schema.ResourceData) (vms.CustomStorageDevices, error) {
	return getDiskDeviceObjects1(d, d.Get(mkDisk).([]interface{}))
}

func getDiskDeviceObjects1(d *schema.ResourceData, disks []interface{}) (vms.CustomStorageDevices, error) {
	diskDeviceObjects := vms.CustomStorageDevices{}
	resource := VM()

	for _, diskEntry := range disks {
		diskDevice := vms.CustomStorageDevice{
			Enabled: true,
		}

		block := diskEntry.(map[string]interface{})
		diskInterface, _ := block[mkDiskInterface].(string)
		datastoreID, _ := block[mkDiskDatastoreID].(string)
		size, _ := block[mkDiskSize].(int)
		fileFormat, _ := block[mkDiskFileFormat].(string)
		fileID, _ := block[mkDiskFileID].(string)
		ioThread := types.CustomBool(block[mkDiskIOThread].(bool))
		ssd := types.CustomBool(block[mkDiskSSD].(bool))
		discard := block[mkDiskDiscard].(string)
		cache := block[mkDiskCache].(string)

		pathInDatastore := ""
		if untyped, hasPathInDatastore := block[mkDiskPathInDatastore]; hasPathInDatastore {
			pathInDatastore = untyped.(string)
		}

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
		diskDevice.Size = types.DiskSizeFromGigabytes(int64(size))
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

		storageInterface := diskDevice.StorageInterface()

		if storageInterface != "virtio" && storageInterface != "scsi" && storageInterface != "sata" {
			return diskDeviceObjects, fmt.Errorf(
				"The disk interface '%s' is not supported, should be one of 'virtioN', 'sataN', or 'scsiN'",
				diskInterface,
			)
		}

		diskDeviceObjects[diskInterface] = &diskDevice
	}

	return diskDeviceObjects, nil
}

func readDisk1(ctx context.Context, d *schema.ResourceData,
	vmConfig *vms.GetResponseData, vmID int, api proxmox.Client, nodeName string, clone []interface{},
) diag.Diagnostics {
	currentDiskList := d.Get(mkDisk).([]interface{})
	diskMap := map[string]interface{}{}
	diskObjects := populateFileIDs(mapStorageDevices(vmConfig), d)

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

	// currentDisks := populateFileIDs(mapStorageDevices(vmConfig), d)

	planDisks, err := getStorageDevicesFromResource(d)
	if err != nil {
		return err
	}

	addToDevices := func(ds vms.CustomStorageDevices, disk *vms.CustomStorageDevice) vms.CustomStorageDevices {
		if ds == nil {
			ds = vms.CustomStorageDevices{}
		}

		ds[*disk.Interface] = disk

		return ds
	}

	for _, disk := range planDisks {
		// for diskInterface, disk := range planDisks {
		// if currentDisks[diskInterface] == nil {
		// 	// TODO: create a new disk here
		// 	return fmt.Errorf("missing device %s", diskInterface)
		// }

		tmp := *disk

		// copy the current disk and update the fields
		// tmp := *currentDisks[diskInterface]
		// tmp.BurstableReadSpeedMbps = disk.BurstableReadSpeedMbps
		// tmp.BurstableWriteSpeedMbps = disk.BurstableWriteSpeedMbps
		// tmp.MaxReadSpeedMbps = disk.MaxReadSpeedMbps
		// tmp.MaxWriteSpeedMbps = disk.MaxWriteSpeedMbps
		// tmp.Cache = disk.Cache
		// tmp.Discard = disk.Discard
		// tmp.IOThread = disk.IOThread
		// tmp.SSD = disk.SSD

		switch disk.StorageInterface() {
		case "virtio":
			updateBody.VirtualIODevices = addToDevices(updateBody.VirtualIODevices, &tmp)
		case "sata":
			updateBody.SATADevices = addToDevices(updateBody.SATADevices, &tmp)
		case "scsi":
			updateBody.SCSIDevices = addToDevices(updateBody.SCSIDevices, &tmp)
		case "ide":
			{
				// Investigate whether to support IDE mapping.
			}
		default:
			return fmt.Errorf("device storage interface %s not supported", disk.StorageInterface())
		}
	}

	return nil
}

// mapStorageDevices maps the current VM storage devices by their interface names.
func mapStorageDevices(resp *vms.GetResponseData) map[string]*vms.CustomStorageDevice {
	storageDevices := map[string]*vms.CustomStorageDevice{}

	fillMap := func(iface string, dev *vms.CustomStorageDevice) {
		if dev != nil {
			d := *dev

			if d.Size == nil {
				d.Size = new(types.DiskSize)
			}

			d.Interface = &iface

			storageDevices[iface] = &d
		}
	}

	fillMap("ide0", resp.IDEDevice0)
	fillMap("ide1", resp.IDEDevice1)
	fillMap("ide2", resp.IDEDevice2)
	fillMap("ide3", resp.IDEDevice3)

	fillMap("sata0", resp.SATADevice0)
	fillMap("sata1", resp.SATADevice1)
	fillMap("sata2", resp.SATADevice2)
	fillMap("sata3", resp.SATADevice3)
	fillMap("sata4", resp.SATADevice4)
	fillMap("sata5", resp.SATADevice5)

	fillMap("scsi0", resp.SCSIDevice0)
	fillMap("scsi1", resp.SCSIDevice1)
	fillMap("scsi2", resp.SCSIDevice2)
	fillMap("scsi3", resp.SCSIDevice3)
	fillMap("scsi4", resp.SCSIDevice4)
	fillMap("scsi5", resp.SCSIDevice5)
	fillMap("scsi6", resp.SCSIDevice6)
	fillMap("scsi7", resp.SCSIDevice7)
	fillMap("scsi8", resp.SCSIDevice8)
	fillMap("scsi9", resp.SCSIDevice9)
	fillMap("scsi10", resp.SCSIDevice10)
	fillMap("scsi11", resp.SCSIDevice11)
	fillMap("scsi12", resp.SCSIDevice12)
	fillMap("scsi13", resp.SCSIDevice13)

	fillMap("virtio0", resp.VirtualIODevice0)
	fillMap("virtio1", resp.VirtualIODevice1)
	fillMap("virtio2", resp.VirtualIODevice2)
	fillMap("virtio3", resp.VirtualIODevice3)
	fillMap("virtio4", resp.VirtualIODevice4)
	fillMap("virtio5", resp.VirtualIODevice5)
	fillMap("virtio6", resp.VirtualIODevice6)
	fillMap("virtio7", resp.VirtualIODevice7)
	fillMap("virtio8", resp.VirtualIODevice8)
	fillMap("virtio9", resp.VirtualIODevice9)
	fillMap("virtio10", resp.VirtualIODevice10)
	fillMap("virtio11", resp.VirtualIODevice11)
	fillMap("virtio12", resp.VirtualIODevice12)
	fillMap("virtio13", resp.VirtualIODevice13)
	fillMap("virtio14", resp.VirtualIODevice14)
	fillMap("virtio15", resp.VirtualIODevice15)

	return storageDevices
}

// mapStorageDevices maps the current VM storage devices by their interface names.
func populateFileIDs(devices vms.CustomStorageDevices, d *schema.ResourceData) vms.CustomStorageDevices {
	planDisk := d.Get(mkDisk)

	planDiskList := planDisk.([]interface{})
	planDiskMap := map[string]map[string]interface{}{}

	for _, v := range planDiskList {
		dm := v.(map[string]interface{})
		iface := dm[mkDiskInterface].(string)

		planDiskMap[iface] = dm
	}

	for k, v := range devices {
		if v != nil && planDiskMap[k] != nil {
			if planDiskMap[k][mkDiskFileID] != nil {
				fileID := planDiskMap[k][mkDiskFileID].(string)
				v.FileID = &fileID
			}
		}
	}

	return devices
}

// getDiskDatastores returns a list of the used datastores in a VM.
func getDiskDatastores(vm *vms.GetResponseData, d *schema.ResourceData) []string {
	storageDevices := populateFileIDs(mapStorageDevices(vm), d)
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
