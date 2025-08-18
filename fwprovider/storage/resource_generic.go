package storage

import (
	"context"
	"fmt"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/config"
	"github.com/bpg/terraform-provider-proxmox/proxmox"
	"github.com/bpg/terraform-provider-proxmox/proxmox/storage"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// storageModel is an interface that all storage resource models must implement.
// This allows a generic resource implementation to handle the CRUD operations.
type storageModel interface {
	// GetID returns the storage identifier from the model.
	GetID() types.String

	// toCreateAPIRequest converts the Terraform model to the specific API request body for creation.
	toCreateAPIRequest(ctx context.Context) (interface{}, error)

	// toUpdateAPIRequest converts the Terraform model to the specific API request body for updates.
	toUpdateAPIRequest(ctx context.Context) (interface{}, error)

	// fromAPI populates the model from the Proxmox API response.
	fromAPI(ctx context.Context, datastore *storage.DatastoreGetResponseData) error
}

// storageResource is a generic implementation for all storage resources.
// It uses a generic type parameter 'T' which must be a pointer to a struct
// that implements the storageModel interface.
type storageResource[T interface {
	*M
	storageModel
}, M any] struct {
	client       proxmox.Client
	storageType  string
	resourceName string
}

// Configure is the generic configuration function.
func (r *storageResource[T, M]) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	cfg, ok := req.ProviderData.(config.Resource)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Resource Configure Type", fmt.Sprintf("Expected config.Resource, got: %T", req.ProviderData))
		return
	}
	r.client = cfg.Client
}

// Create is the generic create function.
func (r *storageResource[T, M]) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan T = new(M)
	diags := req.Plan.Get(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	requestBody, err := plan.toCreateAPIRequest(ctx)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error creating API request for %s storage", r.storageType), err.Error())
		return
	}

	err = r.client.Storage().CreateDatastore(ctx, requestBody)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error creating %s storage", r.storageType), err.Error())
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Read is the generic read function.
func (r *storageResource[T, M]) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state T = new(M)
	diags := req.State.Get(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	datastoreID := state.GetID().ValueString()
	datastore, err := r.client.Storage().GetDatastore(ctx, &storage.DatastoreGetRequest{ID: &datastoreID})
	if err != nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.fromAPI(ctx, datastore)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// Update is the generic update function.
func (r *storageResource[T, M]) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan T = new(M)
	diags := req.Plan.Get(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	requestBody, err := plan.toUpdateAPIRequest(ctx)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error creating API request for %s storage", r.storageType), err.Error())
		return
	}

	err = r.client.Storage().UpdateDatastore(ctx, plan.GetID().ValueString(), requestBody)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error updating %s storage", r.storageType), err.Error())
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Delete is the generic delete function.
func (r *storageResource[T, M]) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state T = new(M)
	diags := req.State.Get(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Storage().DeleteDatastore(ctx, state.GetID().ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error deleting %s storage", r.storageType), err.Error())
		return
	}
}
