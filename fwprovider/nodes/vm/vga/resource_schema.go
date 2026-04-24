/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package vga

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// ResourceSchema defines the schema for the VGA resource.
func ResourceSchema() schema.Attribute {
	return schema.SingleNestedAttribute{
		CustomType: basetypes.ObjectType{
			AttrTypes: attributeTypes(),
		},
		Description: "The VGA configuration.",
		MarkdownDescription: "Configure the VGA Hardware. If you want to use high resolution modes (>= 1280x1024x16) " +
			"you may need to increase the vga memory option. Since QEMU 2.9 the default VGA display type is `std` " +
			"for all OS types besides some Windows versions (XP and older) which use `cirrus`. The `qxl` option " +
			"enables the SPICE display server. For win* OS you can select how many independent displays you want, " +
			"Linux guests can add displays themself. You can also run without any graphic card, using a serial device " +
			"as terminal. See the [Proxmox documentation](https://pve.proxmox.com/pve-docs/pve-admin-guide.html#" +
			"qm_virtual_machines_settings) section 10.2.8 for more information and available configuration parameters.",
		// Optional only (not Computed) per ADR-004 §Provider Defaults vs PVE Defaults: audit Section 4
		// confirms PVE does not auto-populate vga subfields on Read, so block-level Read is null when
		// the user has no `vga` block in HCL.
		Optional: true,
		Attributes: map[string]schema.Attribute{
			"clipboard": schema.StringAttribute{
				Description: "Enable a specific clipboard.",
				MarkdownDescription: "Enable a specific clipboard. If not set, depending on the display type the SPICE " +
					"one will be added. Currently only `vnc` is available. Migration with VNC clipboard is not " +
					"supported by Proxmox.",
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf("vnc"),
				},
			},
			"type": schema.StringAttribute{
				Description:         "The VGA type.",
				MarkdownDescription: "The VGA type (defaults to `std`).",
				Optional:            true,
				// Long, version-evolving PVE enum — per ADR-004 §Enum Validators the provider defers to
				// PVE apply-time validation rather than shipping a release each time PVE adds a type.
			},
			"memory": schema.Int64Attribute{
				Description:         "The VGA memory in megabytes (4-512 MB)",
				MarkdownDescription: "The VGA memory in megabytes (4-512 MB). Has no effect with serial display. ",
				Optional:            true,
				Validators: []validator.Int64{
					int64validator.Between(4, 512),
				},
			},
		},
	}
}
