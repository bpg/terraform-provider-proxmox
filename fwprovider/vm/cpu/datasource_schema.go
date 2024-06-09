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
		Optional:    true,
		Computed:    true,
		Attributes: map[string]schema.Attribute{
			"affinity": schema.StringAttribute{
				Description: "List of host cores used to execute guest processes, for example: '0,5,8-11'",
				Optional:    true,
				Computed:    true,
			},
			"architecture": schema.StringAttribute{
				Description: "The CPU architecture.",
				Optional:    true,
				Computed:    true,
			},
			"cores": schema.Int64Attribute{
				Description: "The number of CPU cores per socket.",
				Optional:    true,
				Computed:    true,
			},
			"flags": schema.SetAttribute{
				Description: "Set of additional CPU flags.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"hotplugged": schema.Int64Attribute{
				Description: "The number of hotplugged vCPUs.",
				Optional:    true,
				Computed:    true,
			},
			"limit": schema.Int64Attribute{
				Description: "Limit of CPU usage.",
				Optional:    true,
				Computed:    true,
			},
			"numa": schema.BoolAttribute{
				Description: "Enable NUMA.",
				Optional:    true,
				Computed:    true,
			},
			"sockets": schema.Int64Attribute{
				Description: "The number of CPU sockets.",
				Optional:    true,
				Computed:    true,
			},
			"type": schema.StringAttribute{
				Description: "Emulated CPU type.",
				Optional:    true,
				Computed:    true,
			},
			"units": schema.Int64Attribute{
				Description: "CPU weight for a VM",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}
