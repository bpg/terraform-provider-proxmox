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
var _ resource.Resource = &zfsPoolStorageResource{}

// NewZFSPoolStorageResource is a helper function to simplify the provider implementation.
func NewZFSPoolStorageResource() resource.Resource {
	return &zfsPoolStorageResource{
		storageResource: &storageResource[
			*ZFSStorageModel, // The pointer to our model
			ZFSStorageModel,  // The struct type of our model
		]{
			storageType:  "zfspool",
			resourceName: "proxmox_virtual_environment_storage_zfspool",
		},
	}
}

// zfsPoolStorageResource is the resource implementation.
type zfsPoolStorageResource struct {
	*storageResource[*ZFSStorageModel, ZFSStorageModel]
}

// Metadata returns the resource type name.
func (r *zfsPoolStorageResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = r.resourceName
}

// Schema defines the schema for the NFS storage resource.
func (r *zfsPoolStorageResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	attributes := map[string]schema.Attribute{
		"zfs_pool": schema.StringAttribute{
			Description: "The name of the ZFS storage pool to use (e.g. `tank`, `rpool/data`).",
			Required:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"thin_provision": schema.BoolAttribute{
			Description: "Whether to enable thin provisioning (`on` or `off`). Thin provisioning allows flexible disk allocation without pre-allocating full space.",
			Optional:    true,
		},
		"blocksize": schema.StringAttribute{
			Description: "Block size for newly created volumes (e.g. `4k`, `8k`, `16k`). Larger values may improve throughput for large I/O, while smaller values optimize space efficiency.",
			Optional:    true,
		},
		"shared": schema.BoolAttribute{
			Description: "Whether the storage is shared across all nodes.",
			Computed:    true,
			Default:     booldefault.StaticBool(false),
		},
	}
	s := storageSchemaFactory(attributes)
	s.Description = "Manages ZFS-based storage in Proxmox VE."
	resp.Schema = s
}
