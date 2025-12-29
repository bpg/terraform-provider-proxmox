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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

// Ensure the implementation satisfies the expected interfaces.
var _ resource.Resource = &cifsStorageResource{}

// NewCIFSStorageResource is a helper function to simplify the provider implementation.
func NewCIFSStorageResource() resource.Resource {
	return &cifsStorageResource{
		storageResource: &storageResource[
			*CIFSStorageModel,
			CIFSStorageModel,
		]{
			storageType:  "cifs",
			resourceName: "proxmox_virtual_environment_storage_cifs",
		},
	}
}

// cifsStorageResource is the resource implementation.
type cifsStorageResource struct {
	*storageResource[*CIFSStorageModel, CIFSStorageModel]
}

// Metadata returns the resource type name.
func (r *cifsStorageResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = r.resourceName
}

// Schema defines the schema for the CIFS storage resource.
func (r *cifsStorageResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	attributes := map[string]schema.Attribute{
		"server": schema.StringAttribute{
			Description: "The IP address or DNS name of the SMB/CIFS server.",
			Required:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"username": schema.StringAttribute{
			Description: "The username for authenticating with the SMB/CIFS server.",
			Required:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"password": schema.StringAttribute{
			Description: "The password for authenticating with the SMB/CIFS server.",
			Required:    true,
			Sensitive:   true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"share": schema.StringAttribute{
			Description: "The name of the SMB/CIFS share.",
			Required:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"domain": schema.StringAttribute{
			Description: "The SMB/CIFS domain.",
			Optional:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"subdirectory": schema.StringAttribute{
			Description: "A subdirectory to mount within the share.",
			Optional:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"preallocation": schema.StringAttribute{
			Description: "The preallocation mode for raw and qcow2 images.",
			Optional:    true,
		},
		"snapshot_as_volume_chain": schema.BoolAttribute{
			Description: "Enable support for creating snapshots through volume backing-chains.",
			Optional:    true,
			PlanModifiers: []planmodifier.Bool{
				boolplanmodifier.RequiresReplace(),
			},
		},
		"shared": schema.BoolAttribute{
			Description: "Whether the storage is shared across all nodes.",
			Computed:    true,
			Default:     booldefault.StaticBool(true),
		},
	}

	factory := NewStorageSchemaFactory()
	factory.WithAttributes(attributes)
	factory.WithDescription("Manages an SMB/CIFS based storage server in Proxmox VE.")
	resp.Schema = *factory.Schema
}
