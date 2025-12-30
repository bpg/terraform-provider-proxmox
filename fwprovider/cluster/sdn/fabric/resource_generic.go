/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package fabric

import (
	"context"
	"errors"
	"fmt"
	"maps"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/config"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/validators"
	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/sdn"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/sdn/fabrics"
	"github.com/bpg/terraform-provider-proxmox/proxmox/helpers/ptr"
	proxmoxtypes "github.com/bpg/terraform-provider-proxmox/proxmox/types"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type genericModel struct {
	ID types.String `tfsdk:"id"`
}

func (m *genericModel) fromAPI(name string, data *fabrics.FabricData, diags *diag.Diagnostics) {
	m.ID = types.StringValue(name)
}

func (m *genericModel) toAPI(ctx context.Context, diags *diag.Diagnostics) *fabrics.Fabric {
	data := &fabrics.Fabric{}

	data.ID = m.ID.ValueString()

	return data
}

func (m *genericModel) getID() string {
	return m.ID.ValueString()
}

func (m *genericModel) handleDeletedStringValue(value *string) types.String {
	if value == nil {
		return types.StringNull()
	}

	if *value == "deleted" {
		return types.StringNull()
	}

	return types.StringValue(*value)
}

func (m *genericModel) handleDeletedInt64Value(value *int64) types.Int64 {
	if value == nil {
		return types.Int64Null()
	}

	return types.Int64Value(*value)
}

func checkDeletedFields(state, plan *genericModel) []string {
	var toDelete []string

	return toDelete
}

func genericAttributesWith(extraAttributes map[string]schema.Attribute) map[string]schema.Attribute {
	// Start with generic attributes as the base
	result := map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The unique identifier of the SDN fabric.",
			Required:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
			Validators: validators.SDNID(),
		},
	}

	// Add extra attributes, allowing them to override generic ones if needed
	if extraAttributes != nil {
		maps.Copy(result, extraAttributes)
	}

	return result
}

type fabricModel interface {
	fromAPI(name string, data *fabrics.FabricData, diags *diag.Diagnostics)
	toAPI(ctx context.Context, diags *diag.Diagnostics) *fabrics.Fabric
	getID() string
	getGenericModel() *genericModel
}

type fabricResourceConfig struct {
	typeNameSuffix string
	fabricProtocol string
	modelFunc      func() fabricModel
}

type genericFabricResource struct {
	client *fabrics.Client
	config fabricResourceConfig
}

func newGenericFabricResource(cfg fabricResourceConfig) resource.Resource {
	return &genericFabricResource{config: cfg}
}

func (r *genericFabricResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + r.config.typeNameSuffix
}

func (r *genericFabricResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = cfg.Client.Cluster().SDNFabrics(r.config.fabricProtocol)
}

func (r *genericFabricResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := r.config.modelFunc()
	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	newFabric := plan.toAPI(ctx, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	newFabric.Protocol = ptr.Ptr(r.config.fabricProtocol)

	if err := r.client.CreateFabric(ctx, newFabric); err != nil {
		resp.Diagnostics.AddError("Unable to Create SDN Fabric", err.Error())
		return
	}

	r.readAndSetState(ctx, plan.getID(), &resp.State, &resp.Diagnostics)
}

func (r *genericFabricResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := r.config.modelFunc()
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	fabric, err := r.client.GetFabricWithParams(ctx, state.getID(), &sdn.QueryParams{Pending: proxmoxtypes.CustomBool(true).Pointer()})
	if err != nil {
		if errors.Is(err, api.ErrResourceDoesNotExist) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError("Unable to Read SDN Fabric", err.Error())

		return
	}

	r.setModelFromFabric(ctx, fabric, &resp.State, &resp.Diagnostics)
}

func (r *genericFabricResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	plan := r.config.modelFunc()
	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)

	state := r.config.modelFunc()
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	updateFabric := plan.toAPI(ctx, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	toDelete := checkDeletedFields(state.getGenericModel(), plan.getGenericModel())
	update := &fabrics.FabricUpdate{
		Fabric: *updateFabric,
		Delete: toDelete,
	}

	if err := r.client.UpdateFabric(ctx, update); err != nil {
		resp.Diagnostics.AddError("Unable to Update SDN Fabric", err.Error())

		return
	}

	r.readAndSetState(ctx, plan.getID(), &resp.State, &resp.Diagnostics)
}

func (r *genericFabricResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	state := r.config.modelFunc()
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteFabric(ctx, state.getID()); err != nil &&
		!errors.Is(err, api.ErrResourceDoesNotExist) {
		resp.Diagnostics.AddError("Unable to Delete SDN Fabric", err.Error())
	}
}

func (r *genericFabricResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	fabric, err := r.client.GetFabric(ctx, req.ID)
	if err != nil {
		if errors.Is(err, api.ErrResourceDoesNotExist) {
			resp.Diagnostics.AddError(fmt.Sprintf("Fabric %s does not exist", req.ID), err.Error())
			return
		}

		resp.Diagnostics.AddError(fmt.Sprintf("Unable to Import SDN Fabric %s", req.ID), err.Error())

		return
	}

	r.setModelFromFabric(ctx, fabric, &resp.State, &resp.Diagnostics)
}

// Schema is required to satisfy the resource.Resource interface. It should be implemented by the specific resource.
func (r *genericFabricResource) Schema(_ context.Context, _ resource.SchemaRequest, _ *resource.SchemaResponse) {
	// Intentionally left blank. Should be set by the specific resource.
}

func (r *genericFabricResource) readAndSetState(ctx context.Context, fabricID string, state *tfsdk.State, diags *diag.Diagnostics) {
	fabric, err := r.client.GetFabricWithParams(ctx, fabricID, &sdn.QueryParams{Pending: proxmoxtypes.CustomBool(true).Pointer()})
	if err != nil {
		diags.AddError("Unable to Read SDN Fabric", err.Error())
		return
	}

	r.setModelFromFabric(ctx, fabric, state, diags)
}

func (r *genericFabricResource) setModelFromFabric(ctx context.Context, fabric *fabrics.FabricData, state *tfsdk.State, diags *diag.Diagnostics) {
	readModel := r.config.modelFunc()
	readModel.fromAPI(fabric.ID, fabric, diags)

	if diags.HasError() {
		return
	}

	diags.Append(state.Set(ctx, readModel)...)
}
