package storage

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

// Ensure the implementation satisfies the expected interfaces.
var _ resource.Resource = &nfsStorageResource{}

// NewNFSStorageResource is a helper function to simplify the provider implementation.
func NewNFSStorageResource() resource.Resource {
	return &nfsStorageResource{
		storageResource: &storageResource[
			*NFSStorageModel, // The pointer to our model
			NFSStorageModel,  // The struct type of our model
		]{
			storageType:  "nfs",
			resourceName: "proxmox_virtual_environment_storage_nfs",
		},
	}
}

// nfsStorageResource is the resource implementation.
type nfsStorageResource struct {
	*storageResource[*NFSStorageModel, NFSStorageModel]
}

// Metadata returns the resource type name.
func (r *nfsStorageResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = r.resourceName
}

// Schema defines the schema for the NFS storage resource.
func (r *nfsStorageResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	s := storageSchemaFactory(map[string]schema.Attribute{
		"server": schema.StringAttribute{
			Description: "The IP address or DNS name of the NFS server.",
			Required:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"export": schema.StringAttribute{
			Description: "The path of the NFS export.",
			Required:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"preallocation": schema.StringAttribute{
			Description: "The preallocation mode for raw and qcow2 images.",
			Optional:    true,
		},
		"options": schema.StringAttribute{
			Description: "The options to pass to the NFS service.",
			Optional:    true,
		},
		"snapshot_as_volume_chain": schema.BoolAttribute{
			Description: "Enable support for creating snapshots through volume backing-chains.",
			Optional:    true,
		},
	})
	s.Description = "Manages an NFS-based storage in Proxmox VE."
	resp.Schema = s
}
