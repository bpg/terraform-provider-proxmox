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
var _ resource.Resource = &lvmPoolStorageResource{}

// NewLVMPoolStorageResource is a helper function to simplify the provider implementation.
func NewLVMPoolStorageResource() resource.Resource {
	return &lvmPoolStorageResource{
		storageResource: &storageResource[
			*LVMStorageModel, // The pointer to our model
			LVMStorageModel,  // The struct type of our model
		]{
			storageType:  "lvm",
			resourceName: "proxmox_virtual_environment_storage_lvm",
		},
	}
}

// lvmPoolStorageResource is the resource implementation.
type lvmPoolStorageResource struct {
	*storageResource[*LVMStorageModel, LVMStorageModel]
}

// Metadata returns the resource type name.
func (r *lvmPoolStorageResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = r.resourceName
}

// Schema defines the schema for the NFS storage resource.
func (r *lvmPoolStorageResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	attributes := map[string]schema.Attribute{
		"volume_group": schema.StringAttribute{
			Description: "The name of the volume group to use.",
			Required:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"wipe_removed_volumes": schema.BoolAttribute{
			Description: "Whether to zero-out data when removing LVMs.",
			Optional:    true,
			Computed:    true,
			Default:     booldefault.StaticBool(false),
		},
	}
	s := storageSchemaFactory(attributes)
	s.Description = "Manages LVM-based storage in Proxmox VE."
	resp.Schema = s
}
