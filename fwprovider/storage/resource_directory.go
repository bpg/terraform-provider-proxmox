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
var _ resource.Resource = &directoryStorageResource{}

// NewDirectoryStorageResource is a helper function to simplify the provider implementation.
func NewDirectoryStorageResource() resource.Resource {
	return &directoryStorageResource{
		storageResource: &storageResource[
			*DirectoryStorageModel,
			DirectoryStorageModel,
		]{
			storageType:  "dir",
			resourceName: "proxmox_virtual_environment_storage_directory",
		},
	}
}

// directoryStorageResource is the resource implementation.
type directoryStorageResource struct {
	*storageResource[*DirectoryStorageModel, DirectoryStorageModel]
}

// Metadata returns the resource type name.
func (r *directoryStorageResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = r.resourceName
}

func (r *directoryStorageResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	attributes := map[string]schema.Attribute{
		"path": schema.StringAttribute{
			Description: "The path to the directory on the Proxmox node.",
			Required:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"preallocation": schema.StringAttribute{
			Description: "The preallocation mode for raw and qcow2 images.",
			Optional:    true,
		},
		"shared": schema.BoolAttribute{
			Description: "Whether the storage is shared across all nodes.",
			Optional:    true,
			Default:     booldefault.StaticBool(true),
			Computed:    true,
		},
	}

	factory := newStorageSchemaFactory()
	factory.WithAttributes(attributes)
	factory.WithDescription("Manages directory-based storage in Proxmox VE.")
	factory.WithBackupBlock()
	resp.Schema = *factory.Schema
}
