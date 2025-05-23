/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package vga

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// DataSourceSchema defines the schema for the VGA datasource.
func DataSourceSchema() schema.Attribute {
	return schema.SingleNestedAttribute{
		CustomType: basetypes.ObjectType{
			AttrTypes: attributeTypes(),
		},
		Description: "The VGA configuration.",
		Optional:    true,
		Computed:    true,
		Attributes: map[string]schema.Attribute{
			"clipboard": schema.StringAttribute{
				Description: "Enable a specific clipboard.",
				Optional:    true,
				Computed:    true,
			},
			"type": schema.StringAttribute{
				Description: "The VGA type.",
				Optional:    true,
				Computed:    true,
			},
			"memory": schema.Int64Attribute{
				Description:         "The VGA memory in megabytes (4-512 MB)",
				MarkdownDescription: "The VGA memory in megabytes (4-512 MB). Has no effect with serial display. ",
				Optional:            true,
				Computed:            true,
			},
		},
	}
}
