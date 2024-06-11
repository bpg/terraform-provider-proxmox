/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package disk

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/resource/validators"
	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/structure"
)

const (
	dvDiskInterface   = "scsi0"
	dvDiskDatastoreID = "local-lvm"
	dvDiskFileFormat  = "qcow2"
	dvDiskSize        = 8
	dvDiskAIO         = "io_uring"
	dvDiskDiscard     = "ignore"
	dvDiskCache       = "none"

	// MkDisk is the name of the disk resource.
	MkDisk                    = "disk"
	mkDiskAIO                 = "aio"
	mkDiskBackup              = "backup"
	mkDiskCache               = "cache"
	mkDiskDatastoreID         = "datastore_id"
	mkDiskDiscard             = "discard"
	mkDiskFileFormat          = "file_format"
	mkDiskFileID              = "file_id"
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
			DefaultFunc: func() (interface{}, error) {
				return []interface{}{
					map[string]interface{}{
						mkDiskAIO:             dvDiskAIO,
						mkDiskBackup:          true,
						mkDiskCache:           dvDiskCache,
						mkDiskDatastoreID:     dvDiskDatastoreID,
						mkDiskDiscard:         dvDiskDiscard,
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
						Type:             schema.TypeString,
						Description:      "The file id for a disk image",
						Optional:         true,
						ForceNew:         true,
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
						DefaultFunc: func() (interface{}, error) {
							return []interface{}{
								map[string]interface{}{
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
			MaxItems: 14,
			MinItems: 0,
		},
	}
}
