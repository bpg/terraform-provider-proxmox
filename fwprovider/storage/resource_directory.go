package storage

import (
	"context"
	"fmt"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/config"
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
			*DirectoryStorageModel, // The pointer to our model
			DirectoryStorageModel,  // The struct type of our model
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

	factory := NewStorageSchemaFactory()
	factory.WithAttributes(attributes)
	factory.WithDescription("Manages directory-based storage in Proxmox VE.")
	factory.WithBackupBlock()
	resp.Schema = *factory.Schema
}

// Configure adds the provider configured client to the resource.
func (r *directoryStorageResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	cfg, ok := req.ProviderData.(config.Resource)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected config.Resource, got: %T", req.ProviderData),
		)

		return
	}

	r.client = cfg.Client
}
