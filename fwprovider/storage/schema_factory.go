package storage

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type StorageSchemaFactory struct {
	Schema *schema.Schema

	description string
}

func NewStorageSchemaFactory() *StorageSchemaFactory {
	s := &schema.Schema{
		Attributes: map[string]schema.Attribute{
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
		},
		Blocks: map[string]schema.Block{},
	}
	return &StorageSchemaFactory{
		Schema: s,
	}
}

func (s *StorageSchemaFactory) WithDescription(description string) *StorageSchemaFactory {
	s.Schema.Description = description
	return s
}

func (s *StorageSchemaFactory) WithAttributes(attributes map[string]schema.Attribute) *StorageSchemaFactory {
	for k, v := range attributes {
		s.Schema.Attributes[k] = v
	}
	return s
}

func (s *StorageSchemaFactory) WithBlocks(blocks map[string]schema.Block) *StorageSchemaFactory {
	for k, v := range blocks {
		s.Schema.Blocks[k] = v
	}
	return s
}

func (s *StorageSchemaFactory) WithBackupBlock() *StorageSchemaFactory {
	return s.WithBlocks(map[string]schema.Block{
		"backups": schema.SingleNestedBlock{
			Attributes: map[string]schema.Attribute{
				"max_protected_backups": schema.Int64Attribute{
					Description: "The maximum number of protected backups per guest. Use '-1' for unlimited.",
					Optional:    true,
				},
				"keep_last": schema.Int64Attribute{
					Description: "Specifies the number of the most recent backups to keep, regardless of their age.",
					Optional:    true,
					Validators: []validator.Int64{
						int64validator.AtLeast(0),
					},
				},
				"keep_hourly": schema.Int64Attribute{
					Description: "The number of hourly backups to keep. Older backups will be removed.",
					Optional:    true,
					Validators: []validator.Int64{
						int64validator.AtLeast(0),
					},
				},
				"keep_daily": schema.Int64Attribute{
					Description: "The number of daily backups to keep. Older backups will be removed.",
					Optional:    true,
					Validators: []validator.Int64{
						int64validator.AtLeast(0),
					},
				},
				"keep_weekly": schema.Int64Attribute{
					Description: "The number of weekly backups to keep. Older backups will be removed.",
					Optional:    true,
					Validators: []validator.Int64{
						int64validator.AtLeast(0),
					},
				},
				"keep_monthly": schema.Int64Attribute{
					Description: "The number of monthly backups to keep. Older backups will be removed.",
					Optional:    true,
					Validators: []validator.Int64{
						int64validator.AtLeast(0),
					},
				},
				"keep_yearly": schema.Int64Attribute{
					Description: "The number of yearly backups to keep. Older backups will be removed.",
					Optional:    true,
					Validators: []validator.Int64{
						int64validator.AtLeast(0),
					},
				},
				"keep_all": schema.BoolAttribute{
					Description: "Specifies if all backups should be kept, regardless of their age.",
					Optional:    true,
					Computed:    true,
					Default:     booldefault.StaticBool(false),
				},
			},
			Description: "Configure backup retention settings for the storage type.",
		},
	})
}
