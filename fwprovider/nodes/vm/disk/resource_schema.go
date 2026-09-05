/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package disk

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func ResourceSchema() schema.Attribute {
	return schema.MapNestedAttribute{
		Description: "The disk configuration",
		MarkdownDescription: "The disk configuration. The key is the interface of the Disk," +
			" could be one of `ide[0-3]`, `sata[0-5]`, `scsi[0-30]`, where the number is" +
			" the index of the interface.",
		Optional: true,
		Validators: []validator.Map{
			mapvalidator.KeysAre(
				stringvalidator.RegexMatches(
					// Slot bounds per qemu-server.git: MAX_IDE_DISKS=4, MAX_SATA_DISKS=6, MAX_SCSI_DISKS=31.
					regexp.MustCompile(`^(ide[0-3]|sata[0-5]|scsi([0-9]|[12][0-9]|30))$`),
					"one of `ide[0-3]`, `sata[0-5]`, `scsi[0-30]`",
				),
			),
		},
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"aio": schema.StringAttribute{
					Description:         "The disk AIO mode",
					MarkdownDescription: "The disk AIO mode `<io_uring | native | thread>` (defaults to `io_uring` when unset).",
					Optional:            true,
					Validators: []validator.String{
						stringvalidator.OneOf("io_uring", "native", "threads"),
					},
				},
				"backup": schema.BoolAttribute{
					Description:         "Enable disk Backup",
					MarkdownDescription: "Whether the drive should be included when making backups (defaults to `true`).",
					Optional:            true,
				},
				"cache": schema.StringAttribute{
					Description:         "The cache type",
					MarkdownDescription: "The cache type. `<none | directsync | writethrough | writeback | unsafe>` (defaults to `none`).",
					Optional:            true,
					Validators: []validator.String{
						stringvalidator.OneOf("none", "directsync", "writethrough", "writeback", "unsafe"),
					},
				},
				"datastore_id": schema.StringAttribute{
					Description:         "The identifier for the datastore",
					MarkdownDescription: "The identifier for the datastore to create the disk in (defaults to `local-lvm`).",
					Optional:            true,
				},
				"discard": schema.StringAttribute{
					Description: "Enable/disable discard",
					MarkdownDescription: "Whether to pass discard/trim requests to the underlying storage." +
						" Supported values are `on`/`ignore` (defaults to `ignore`).",
					Optional: true,
				},
				"file_format": schema.StringAttribute{
					Description:         "The file format.",
					MarkdownDescription: "The file format `<qcow2 | raw | vmdk>`",
					Optional:            true,
					Validators: []validator.String{
						stringvalidator.OneOf("qcow2", "raw", "vmdk"),
					},
				},
				"import_from": schema.StringAttribute{
					Description: "The file ID for a disk image to import into VM.",
					MarkdownDescription: "The file ID for a disk image to import into VM." +
						" The image must be of `import` content type (uncompressed images only)." +
						" The ID format is `<datastore_id>:import/<file_name>`," +
						" Can be also taken from `proxmox_download_file` resource." +
						" Note: compressed images downloaded with `decompression_algorithm` cannot" +
						" use `import_from`; use `file_id`instead.",
					Optional: true,
				},
				"iothread": schema.BoolAttribute{
					Description:         "Enable IOThread",
					MarkdownDescription: "Whether to use IOThreads for this disk. (defaults to `false`).",
					Optional:            true,
				},
				"size": schema.Int64Attribute{
					Description:         "The disk size in gigabytes",
					MarkdownDescription: "The disk size in gigabytes",
					Optional:            true,
				},
			},
		},
	}
}
