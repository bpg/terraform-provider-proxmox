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

type schemaFactoryOptions struct {
	IsSharedByDefault bool
}

// storageSchemaFactory generates the schema for a storage resource.
func storageSchemaFactory(specificAttributes map[string]schema.Attribute, opt ...*schemaFactoryOptions) schema.Schema {
	options := &schemaFactoryOptions{}
	if opt != nil && len(opt) > 0 {
		options = opt[0]
	}
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

	if options.IsSharedByDefault {
		// For types like NFS, 'shared' is a computed, read-only attribute. The user cannot set it.
		attributes["shared"] = schema.BoolAttribute{
			Description: "Whether the storage is shared across all nodes. This is inherent to the storage type.",
			Computed:    true,
			Default:     booldefault.StaticBool(true),
		}
	} else {
		attributes["shared"] = schema.BoolAttribute{
			Description: "Whether the storage is shared across all nodes.",
			Optional:    true,
			Computed:    true,
		}
	}

	// Merge provided attributes for the given storage type
	for k, v := range specificAttributes {
		attributes[k] = v
	}

	return schema.Schema{
		Attributes: attributes,
	}
}
