package storage

import (
	"context"
	"fmt"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/config"
	"github.com/bpg/terraform-provider-proxmox/proxmox"
	"github.com/bpg/terraform-provider-proxmox/proxmox/storage"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &directoryStorageResource{}
	_ resource.ResourceWithConfigure = &directoryStorageResource{}
)

var allowedStorageTypes = []string{
	"btrfs", "cephfs", "cifs", "dir", "esxi", "iscsi", "iscsidirect",
	"lvm", "lvmthin", "nfs", "pbs", "rbd", "zfs", "zfspool",
}

// NewDirectoryStorageResource is a helper function to simplify the provider implementation.
func NewDirectoryStorageResource() resource.Resource {
	return &directoryStorageResource{}
}

// directoryStorageResource is the resource implementation.
type directoryStorageResource struct {
	client proxmox.Client
}

func (r *directoryStorageResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a directory-based storage in Proxmox VE.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the storage.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				Description: "The type of storage to create.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(allowedStorageTypes...),
				},
			},
			"path": schema.StringAttribute{
				Description: "The path to the directory on the Proxmox node.",
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
			"shared": schema.BoolAttribute{
				Description: "Whether the storage is shared across all nodes.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"preallocation": schema.StringAttribute{
				Description: "The preallocation mode for raw and qcow2 images.",
				Optional:    true,
			},
		},
	}
}

// Create creates the resource and sets the initial state.
func (r *directoryStorageResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan DirectoryStorageModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	requestBody, err := plan.toCreateAPIRequest(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error creating create request for directory storage", err.Error())
		return
	}

	err = r.client.Storage().CreateDatastore(ctx, &requestBody)
	if err != nil {
		resp.Diagnostics.AddError("Error creating directory storage", err.Error())
		return
	}

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Read refreshes the resource state from the Proxmox API.
func (r *directoryStorageResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state DirectoryStorageModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	requestBody := &storage.DatastoreGetRequest{ID: state.ID.ValueStringPointer()}
	datastore, err := r.client.Storage().GetDatastore(ctx, requestBody)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Proxmox Storage",
			"Could not read storage ("+state.ID.ValueString()+"): "+err.Error(),
		)
		return
	}

	state.importFromAPI(ctx, *datastore)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the new state.
func (r *directoryStorageResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan DirectoryStorageModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	requestBody, err := plan.toUpdateAPIRequest(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error creating update request for directory storage", err.Error())
		return
	}

	err = r.client.Storage().UpdateDatastore(ctx, plan.ID.ValueString(), &requestBody)
	if err != nil {
		resp.Diagnostics.AddError("Error updating directory storage", err.Error())
		return
	}

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes the resource and removes it from the state.
func (r *directoryStorageResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state DirectoryStorageModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Storage().DeleteDatastore(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting directory storage",
			"Could not delete directory storage, unexpected error: "+err.Error(),
		)
		return
	}
}

// Metadata returns the resource type name.
func (r *directoryStorageResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_storage_directory"
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
