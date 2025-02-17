/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package rng

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// ResourceSchema defines the schema for the RNG resource.
func ResourceSchema() schema.Attribute {
	return schema.SingleNestedAttribute{
		CustomType: basetypes.ObjectType{
			AttrTypes: attributeTypes(),
		},
		Description: "The RNG (Random Number Generator) configuration. Can only be set by `root@pam.`",
		MarkdownDescription: "Configure the RNG (Random Number Generator) device. The RNG device provides entropy " +
			"to guests to ensure good quality random numbers for guest applications that require them. " +
			"Can only be set by `root@pam.`" +
			"See the [Proxmox documentation](https://pve.proxmox.com/pve-docs/pve-admin-guide.html#qm_virtual_machines_settings) " +
			"for more information.",
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.Object{
			objectplanmodifier.UseStateForUnknown(),
		},
		Attributes: map[string]schema.Attribute{
			"source": schema.StringAttribute{
				Description: "The entropy source for the RNG device.",
				MarkdownDescription: "The file on the host to gather entropy from. " +
					"In most cases `/dev/urandom` should be preferred over `/dev/random` " +
					"to avoid entropy-starvation issues on the host.",
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"max_bytes": schema.Int64Attribute{
				Description: "Maximum bytes of entropy allowed to get injected into the guest every period.",
				MarkdownDescription: "Maximum bytes of entropy allowed to get injected into the guest every period. " +
					"Use 0 to disable limiting (potentially dangerous).",
				Optional: true,
				Computed: true,
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
				},
			},
			"period": schema.Int64Attribute{
				Description: "Period in milliseconds to limit entropy injection to the guest.",
				MarkdownDescription: "Period in milliseconds to limit entropy injection to the guest. " +
					"Use 0 to disable limiting (potentially dangerous).",
				Optional: true,
				Computed: true,
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
				},
			},
		},
	}
}
