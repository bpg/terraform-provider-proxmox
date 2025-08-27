/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package storage

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

// Ensure the implementation satisfies the expected interfaces.
var _ resource.Resource = &lvmThinPoolStorageResource{}

// NewLVMThinPoolStorageResource is a helper function to simplify the provider implementation.
func NewLVMThinPoolStorageResource() resource.Resource {
	return &lvmThinPoolStorageResource{
		storageResource: &storageResource[
			*LVMThinStorageModel,
			LVMThinStorageModel,
		]{
			storageType:  "lvmthin",
			resourceName: "proxmox_virtual_environment_storage_lvmthin",
		},
	}
}

// lvmThinPoolStorageResource is the resource implementation.
type lvmThinPoolStorageResource struct {
	*storageResource[*LVMThinStorageModel, LVMThinStorageModel]
}

// Metadata returns the resource type name.
func (r *lvmThinPoolStorageResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = r.resourceName
}

// Schema defines the schema for the NFS storage resource.
func (r *lvmThinPoolStorageResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	attributes := map[string]schema.Attribute{
		"volume_group": schema.StringAttribute{
			Description: "The name of the volume group to use.",
			Required:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"thin_pool": schema.StringAttribute{
			Description: "The name of the LVM thin pool to use.",
			Required:    true,
		},
		"shared": schema.BoolAttribute{
			Description: "Whether the storage is shared across all nodes.",
			Optional:    true,
			Default:     booldefault.StaticBool(false),
			Computed:    true,
		},
	}
	factory := NewStorageSchemaFactory()
	factory.WithAttributes(attributes)
	factory.WithDescription("Manages thin LVM-based storage in Proxmox VE.")
	resp.Schema = *factory.Schema
}
