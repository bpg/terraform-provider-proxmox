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
			"must be one of `ideN`, `sataN`, `scsiN`, where N is the index of the interface. " +
			"Note that `q35` machine type only supports `ide0` and `ide2` of IDE interfaces.",
		Optional: true,
		Computed: true,
		Validators: []validator.Map{
			mapvalidator.KeysAre(
				stringvalidator.RegexMatches(
					regexp.MustCompile(`^(ide[0-3]|sata[0-5]|scsi([0-9]|1[0-3]))$`),
					"one of `ide[0-3]`, `sata[0-5]`, `scsi[0-13]`",
				),
			),
		},
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"file_id": schema.StringAttribute{
					Description: "The file ID of the CD-ROM",
					MarkdownDescription: "The file ID of the CD-ROM, or `cdrom|none`." +
						" Defaults to `none` to leave the CD-ROM empty. Use `cdrom` to connect to the physical drive.",
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
