/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package cpu

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// DataSourceSchema defines the schema for the CPU datasource.
func DataSourceSchema() schema.Attribute {
	return schema.SingleNestedAttribute{
		CustomType: basetypes.ObjectType{
			AttrTypes: attributeTypes(),
		},
		Description: "The CPU configuration.",
		Computed:    true,
		Attributes: map[string]schema.Attribute{
			"affinity": schema.StringAttribute{
				Description: "List of host cores used to execute guest processes, for example: '0,5,8-11'",
				Computed:    true,
			},
			"architecture": schema.StringAttribute{
				Description: "The CPU architecture.",
				Computed:    true,
			},
			"cores": schema.Int64Attribute{
				Description: "The number of CPU cores per socket.",
				Computed:    true,
			},
			"flags": schema.SetAttribute{
				Description: "Set of additional CPU flags.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"limit": schema.Float64Attribute{
				Description: "Limit of CPU usage.",
				Computed:    true,
			},
			"numa": schema.BoolAttribute{
				Description: "Whether NUMA emulation is enabled.",
				Computed:    true,
			},
			"sockets": schema.Int64Attribute{
				Description: "The number of CPU sockets.",
				Computed:    true,
			},
			"type": schema.StringAttribute{
				Description: "Emulated CPU type.",
				Computed:    true,
			},
			"units": schema.Int64Attribute{
				Description: "CPU weight for a VM",
				Computed:    true,
			},
			"vcpus": schema.Int64Attribute{
				Description: "Number of active vCPUs.",
				Computed:    true,
			},
		},
	}
}
