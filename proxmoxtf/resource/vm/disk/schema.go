package disk

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/resource/validators"
)

const (
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

	// MkDisk is the name of the disk resource.
	MkDisk                    = "disk"
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

	// MkTimeoutMoveDisk is the name of the timeout_move_disk attribute.
	MkTimeoutMoveDisk = "timeout_move_disk"
)

// Schema returns the schema for the disk resource.
func Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		MkDisk: {
			Type:        schema.TypeList,
			Description: "The disk devices",
			Optional:    true,
			ForceNew:    true,
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
						ValidateDiagFunc: validators.FileFormat(),
					},
					mkDiskFileID: {
						Type:             schema.TypeString,
						Description:      "The file id for a disk image",
						Optional:         true,
						ForceNew:         true,
						Default:          dvDiskFileID,
						ValidateDiagFunc: validators.FileID(),
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
			MaxItems: 14,
			MinItems: 0,
		},
	}
}
