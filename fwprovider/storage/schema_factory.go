package storage

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// storageSchemaFactory generates the schema for a storage resource.
func storageSchemaFactory(specificAttributes map[string]schema.Attribute) schema.Schema {
	attributes := map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The unique identifier of the storage.",
			Required:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"nodes": schema.SetAttribute{
			Description: "A list of nodes where this storage is available.",
			ElementType: types.StringType,
			Optional:    true,
			Computed:    true,
			Default: setdefault.StaticValue(
				types.SetValueMust(types.StringType, []attr.Value{}),
			),
		},
		"content": schema.SetAttribute{
			Description: "The content types that can be stored on this storage.",
			ElementType: types.StringType,
			Optional:    true,
			Computed:    true,
			Default: setdefault.StaticValue(
				types.SetValueMust(types.StringType, []attr.Value{}),
			),
		},
		"disable": schema.BoolAttribute{
			Description: "Whether the storage is disabled.",
			Optional:    true,
			Default:     booldefault.StaticBool(false),
			Computed:    true,
		},
	}

	// Merge provided attributes for the given storage type
	for k, v := range specificAttributes {
		attributes[k] = v
	}

	return schema.Schema{
		Attributes: attributes,
	}
}
