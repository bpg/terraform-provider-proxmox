/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package storage

import (
	"maps"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/storage"
)

type schemaFactory struct {
	Schema *schema.Schema
}

func newStorageSchemaFactory() *schemaFactory {
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
			},
			"content": schema.SetAttribute{
				Description: "The content types that can be stored on this storage",
				MarkdownDescription: "The content types that can be stored on this storage. " +
					"Valid values: `backup` (VM backups), `images` (VM disk images), " +
					"`import` (VM disk images for import), `iso` (ISO images), " +
					"`rootdir` (container root directories), `snippets` (cloud-init, hook scripts, etc.), " +
					"`vztmpl` (container templates).",
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(
						stringvalidator.OneOf(storage.ValidContentTypes()...),
					),
				},
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

	return &schemaFactory{
		Schema: s,
	}
}

func (s *schemaFactory) WithDescription(description string) *schemaFactory {
	s.Schema.Description = description
	return s
}

func (s *schemaFactory) WithAttributes(attributes map[string]schema.Attribute) *schemaFactory {
	maps.Copy(s.Schema.Attributes, attributes)

	return s
}

func (s *schemaFactory) WithBlocks(blocks map[string]schema.Block) *schemaFactory {
	maps.Copy(s.Schema.Blocks, blocks)

	return s
}

func (s *schemaFactory) WithBackupBlock() *schemaFactory {
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
					Description: "Specifies if all backups should be kept, regardless of their age. " +
						"When set to true, other keep_* attributes must not be set.",
					Optional: true,
					Computed: true,
					Default:  booldefault.StaticBool(false),
				},
			},
			Validators: []validator.Object{
				backupsKeepAllExcludesOtherKeepSettingsValidator{},
			},
			Description: "Configure backup retention settings for the storage type.",
		},
	})
}
