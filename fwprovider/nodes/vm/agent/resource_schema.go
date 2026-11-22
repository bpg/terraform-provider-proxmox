/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL can be obtained
 * at http://mozilla.org/MPL/2.0/.
 */

package agent

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func ResourceSchema() schema.Attribute {
	return schema.SingleNestedAttribute{
		CustomType: basetypes.ObjectType{
			AttrTypes: attributeTypes(),
		},
		Description: "The QEMU guest agent configuration.",
		Optional:    true,
		Attributes: map[string]schema.Attribute{
			"enabled": schema.BoolAttribute{
				Description: "Whether the QEMU guest agent is enabled.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"trim": schema.BoolAttribute{
				Description: "Whether to enable fstrim after cloning.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"type": schema.StringAttribute{
				Description: "The QEMU guest agent interface type.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("virtio"),
				Validators: []validator.String{
					stringvalidator.OneOf("virtio", "isa"),
				},
			},
			"freeze_fs": schema.BoolAttribute{
				Description: "Whether filesystem freeze is enabled for backup operations.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
		},
	}
}
