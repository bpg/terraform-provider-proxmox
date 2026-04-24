/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package cdrom

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/validators"
)

// ResourceSchema defines the schema for the CD-ROM resource.
func ResourceSchema() schema.Attribute {
	return schema.MapNestedAttribute{
		Description: "The CD-ROM configuration",
		MarkdownDescription: "The CD-ROM configuration. The key is the interface of the CD-ROM, " +
			"could be one of `ideN`, `sataN`, `scsiN`, where N is the index of the interface. " +
			"Note that `q35` machine type only supports `ide0` and `ide2` of IDE interfaces.",
		// Optional only (not Computed) per ADR-004 §Provider Defaults vs PVE Defaults: PVE does
		// not auto-attach CD-ROM devices to a VM, so the map-level Read value should be null when
		// the user has no `cdrom` block in HCL.
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
				"file_id": schema.StringAttribute{
					Description: "The file ID of the CD-ROM",
					// `Default("cdrom")` is kept as provider UX so `cdrom = { ide2 = {} }` declares an empty
					// slot without repeating the `file_id = "cdrom"` ritual ("cdrom" is PVE's literal "no
					// media inserted" storage path). This is an explicit carve-out from ADR-004 §Provider
					// Defaults vs PVE Defaults: PVE always surfaces `file_id` when the slot exists, so the
					// default is not masking a PVE auto-populate — it is a pure HCL shorthand.
					MarkdownDescription: "The file ID of the CD-ROM, or `cdrom|none`." +
						" Defaults to `cdrom` (i.e. empty CD-ROM drive — `cdrom` is PVE's literal \"no media inserted\"" +
						" storage path). Use `none` to leave the CD-ROM unplugged, or a storage path like" +
						" `local:iso/debian.iso` to insert an image.",
					Optional: true,
					Computed: true,
					Default:  stringdefault.StaticString("cdrom"),
					Validators: []validator.String{
						stringvalidator.Any(
							stringvalidator.OneOf("cdrom", "none"),
							validators.FileID(),
						),
					},
				},
			},
		},
	}
}
