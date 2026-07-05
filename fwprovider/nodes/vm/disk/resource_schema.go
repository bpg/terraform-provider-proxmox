/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package disk

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func ResourceSchema() schema.Attribute {
	return schema.SingleNestedAttribute{
		Description: "(Optional) A Disk (multiple blocks supported)",
		Optional:    true,
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
				Default:             stringdefault.StaticString("none"),
				Validators: []validator.String{
					stringvalidator.OneOf("none", "directsync", "writethrough", "writeback", "unsafe"),
				},
			},
			"datastore_id": schema.StringAttribute{
				Description:         "The identifier for the datastore",
				MarkdownDescription: "The identifier for the datastore to create the disk in (defaults to `local-lvm`).",
				Optional:            true,
				Default:             stringdefault.StaticString("local-lvm"),
			},
			"discard": schema.StringAttribute{
				Description: "Enable/disable discard",
				MarkdownDescription: "Whether to pass discard/trim requests to the underlying storage." +
					" Supported values are `on`/`ignore` (defaults to `ignore`).",
				Optional: true,
				Default:  stringdefault.StaticString("ignore"),
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
			"interface": schema.StringAttribute{
				Description: "The disk interface",
				MarkdownDescription: "The disk interface for Proxmox," +
					" currently, `scsi`, `sata` and `virtio` are supported." +
					" Append the disk index at the end, for example, `virtio0` for the" +
					" first virtio disk, `virtio1` for the second, etc.",
				Required: true,
				//TODO: implement interface validation
			},
			"size": schema.Int32Attribute{
				Description:         "The disk size in gigabytes",
				MarkdownDescription: "The disk size in gigabytes",
				Optional:            true,
				Default:             int32default.StaticInt32(8),
			},
		},
	}
}
