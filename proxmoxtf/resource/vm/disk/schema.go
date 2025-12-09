/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package disk

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/resource/validators"
	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/structure"
	"github.com/bpg/terraform-provider-proxmox/utils"
)

const (
	dvDiskInterface   = "scsi0"
	dvDiskDatastoreID = "local-lvm"
	dvDiskSize        = 8
	dvDiskAIO         = "io_uring"
	dvDiskDiscard     = "ignore"
	dvDiskCache       = "none"

	// see /usr/share/perl5/PVE/QemuServer/Drive.pm
	// Using SCSI limit (31) as the highest value (IDE: 4, SCSI: 31, VirtIO: 16, SATA: 6).
	maxResourceVirtualEnvironmentVMDiskDevices = 31

	// MkDisk is the name of the disk resource.
	MkDisk                    = "disk"
	mkDiskAIO                 = "aio"
	mkDiskBackup              = "backup"
	mkDiskCache               = "cache"
	mkDiskDatastoreID         = "datastore_id"
	mkDiskDiscard             = "discard"
	mkDiskFileFormat          = "file_format"
	mkDiskFileID              = "file_id"
	mkDiskImportFrom          = "import_from"
	mkDiskInterface           = "interface"
	mkDiskIopsRead            = "iops_read"
	mkDiskIopsReadBurstable   = "iops_read_burstable"
	mkDiskIopsWrite           = "iops_write"
	mkDiskIopsWriteBurstable  = "iops_write_burstable"
	mkDiskIOThread            = "iothread"
	mkDiskPathInDatastore     = "path_in_datastore"
	mkDiskReplicate           = "replicate"
	mkDiskSerial              = "serial"
	mkDiskSize                = "size"
	mkDiskSpeed               = "speed"
	mkDiskSpeedRead           = "read"
	mkDiskSpeedReadBurstable  = "read_burstable"
	mkDiskSpeedWrite          = "write"
	mkDiskSpeedWriteBurstable = "write_burstable"
	mkDiskSSD                 = "ssd"
)

// Schema returns the schema for the disk resource.
func Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		MkDisk: {
			Type:        schema.TypeList,
			Description: "The disk devices",
			Optional:    true,
			DefaultFunc: func() (any, error) {
				return []any{
					map[string]any{
						mkDiskAIO:             dvDiskAIO,
						mkDiskBackup:          true,
						mkDiskCache:           dvDiskCache,
						mkDiskDatastoreID:     dvDiskDatastoreID,
						mkDiskDiscard:         dvDiskDiscard,
						mkDiskImportFrom:      "",
						mkDiskFileID:          "",
						mkDiskInterface:       dvDiskInterface,
						mkDiskIOThread:        false,
						mkDiskPathInDatastore: nil,
						mkDiskReplicate:       true,
						mkDiskSerial:          "",
						mkDiskSize:            dvDiskSize,
						mkDiskSSD:             false,
					},
				}, nil
			},
			DiffSuppressFunc: structure.SuppressIfListsOfMapsAreEqualIgnoringOrderByKey(
				mkDiskInterface, mkDiskPathInDatastore,
			),
			DiffSuppressOnRefresh: true,
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
						Computed:         true,
						ValidateDiagFunc: validators.FileFormat(),
					},
					mkDiskAIO: {
						Type:        schema.TypeString,
						Description: "The disk AIO mode",
						Optional:    true,
						Default:     dvDiskAIO,
						ValidateDiagFunc: validation.ToDiagFunc(
							validation.StringInSlice([]string{
								"io_uring",
								"native",
								"threads",
							}, false),
						),
					},
					mkDiskBackup: {
						Type:        schema.TypeBool,
						Description: "Whether the drive should be included when making backups",
						Optional:    true,
						Default:     true,
					},
					mkDiskFileID: {
						Type:        schema.TypeString,
						Description: "The file id for a disk image",
						Optional:    true,
						ForceNew:    true,
						Default:     "",
						DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
							if old == "" && new == "" {
								return true
							}

							if old == new {
								return true
							}

							diskAttrPath := k[:strings.LastIndex(k, ".")]
							diskIndexStr := diskAttrPath[strings.LastIndex(diskAttrPath, ".")+1:]
							diskIndex, err := strconv.Atoi(diskIndexStr)
							if err != nil {
								return false
							}

							oldData, newData := d.GetChange(MkDisk)
							oldArray, ok := oldData.([]any)
							if !ok {
								return false
							}
							newArray, ok := newData.([]any)
							if !ok {
								return false
							}

							if diskIndex >= len(newArray) {
								return false
							}

							newDisk, ok := newArray[diskIndex].(map[string]any)
							if !ok {
								return false
							}

							newInterface, _ := newDisk[mkDiskInterface].(string)
							if newInterface == "" {
								return false
							}

							oldMap := utils.MapResourcesByAttribute(oldArray, mkDiskInterface)

							oldDiskByInterface, oldExists := oldMap[newInterface].(map[string]any)
							if !oldExists {
								return false
							}

							oldFileIDByInterface, _ := oldDiskByInterface[mkDiskFileID].(string)

							return oldFileIDByInterface == new
						},
						ValidateDiagFunc: validators.FileID(),
					},
					mkDiskImportFrom: {
						Type:             schema.TypeString,
						Description:      "The file id of a disk image to import from storage.",
						Optional:         true,
						ForceNew:         false,
						Default:          "",
						ValidateDiagFunc: validators.FileID(),
					},
					mkDiskSerial: {
						Type:             schema.TypeString,
						Description:      "The drive’s reported serial number",
						Optional:         true,
						Default:          "",
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(0, 20)),
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
						Default:     false,
					},
					mkDiskReplicate: {
						Type:        schema.TypeBool,
						Description: "Whether the drive should be considered for replication jobs",
						Optional:    true,
						Default:     true,
					},
					mkDiskSSD: {
						Type:        schema.TypeBool,
						Description: "Whether to use ssd for this disk drive",
						Optional:    true,
						Default:     false,
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
						DefaultFunc: func() (any, error) {
							return []any{
								map[string]any{
									mkDiskIopsRead:            0,
									mkDiskIopsWrite:           0,
									mkDiskIopsReadBurstable:   0,
									mkDiskIopsWriteBurstable:  0,
									mkDiskSpeedRead:           0,
									mkDiskSpeedReadBurstable:  0,
									mkDiskSpeedWrite:          0,
									mkDiskSpeedWriteBurstable: 0,
								},
							}, nil
						},
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								mkDiskIopsRead: {
									Type:        schema.TypeInt,
									Description: "The maximum read I/O in operations per second",
									Optional:    true,
									Default:     0,
								},
								mkDiskIopsWrite: {
									Type:        schema.TypeInt,
									Description: "The maximum write I/O in operations per second",
									Optional:    true,
									Default:     0,
								},
								mkDiskIopsReadBurstable: {
									Type:        schema.TypeInt,
									Description: "The maximum unthrottled read I/O pool in operations per second",
									Optional:    true,
									Default:     0,
								},
								mkDiskIopsWriteBurstable: {
									Type:        schema.TypeInt,
									Description: "The maximum unthrottled write I/O pool in operations per second",
									Optional:    true,
									Default:     0,
								},
								mkDiskSpeedRead: {
									Type:        schema.TypeInt,
									Description: "The maximum read speed in megabytes per second",
									Optional:    true,
									Default:     0,
								},
								mkDiskSpeedReadBurstable: {
									Type:        schema.TypeInt,
									Description: "The maximum burstable read speed in megabytes per second",
									Optional:    true,
									Default:     0,
								},
								mkDiskSpeedWrite: {
									Type:        schema.TypeInt,
									Description: "The maximum write speed in megabytes per second",
									Optional:    true,
									Default:     0,
								},
								mkDiskSpeedWriteBurstable: {
									Type:        schema.TypeInt,
									Description: "The maximum burstable write speed in megabytes per second",
									Optional:    true,
									Default:     0,
								},
							},
						},
						MaxItems: 1,
						MinItems: 0,
					},
				},
			},
			MaxItems: maxResourceVirtualEnvironmentVMDiskDevices,
			MinItems: 0,
		},
	}
}

// CustomizeDiff returns the custom diff functions for the disk resource.
func CustomizeDiff() []schema.CustomizeDiffFunc {
	return []schema.CustomizeDiffFunc{
		customdiff.If(
			func(_ context.Context, d *schema.ResourceDiff, _ any) bool {
				return d.HasChange(MkDisk)
			},
			func(ctx context.Context, d *schema.ResourceDiff, _ any) error {
				oldData, newData := d.GetChange(MkDisk)

				oldArray, ok := oldData.([]any)
				if !ok {
					return nil
				}

				newArray, ok := newData.([]any)
				if !ok {
					return nil
				}

				if len(oldArray) != len(newArray) {
					return nil
				}

				oldMap := utils.MapResourcesByAttribute(oldArray, mkDiskInterface)
				newMap := utils.MapResourcesByAttribute(newArray, mkDiskInterface)

				if len(oldMap) != len(newMap) {
					return nil
				}

				copyWithoutPath := func(disk map[string]any) map[string]any {
					diskCopy := make(map[string]any)

					for key, val := range disk {
						if key != mkDiskPathInDatastore {
							diskCopy[key] = val
						}
					}

					return diskCopy
				}

				allEqual := true
				for k, v := range oldMap {
					if _, ok := newMap[k]; !ok {
						allEqual = false
						break
					}

					oldDisk := v.(map[string]any)
					newDisk := newMap[k].(map[string]any)

					if !reflect.DeepEqual(copyWithoutPath(oldDisk), copyWithoutPath(newDisk)) {
						allEqual = false
					}
				}

				if allEqual {
					for i := range oldArray {
						diskMap := oldArray[i].(map[string]any)
						for key := range diskMap {
							attrPath := MkDisk + "." + strconv.Itoa(i) + "." + key
							if d.HasChange(attrPath) {
								if err := d.Clear(attrPath); err != nil {
									return fmt.Errorf("failed to clear diff for %s: %w", attrPath, err)
								}
							}
						}
					}

					if d.HasChange(MkDisk + ".#") {
						if err := d.Clear(MkDisk + ".#"); err != nil {
							return fmt.Errorf("failed to clear diff for %s: %w", MkDisk+".#", err)
						}
					}
				}

				return nil
			},
		),
	}
}
