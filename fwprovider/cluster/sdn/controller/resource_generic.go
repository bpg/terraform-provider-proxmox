/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package controller

import (
	"context"
	"errors"
	"fmt"
	"maps"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/attribute"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/config"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/types/stringset"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/validators"
	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/sdn"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/sdn/controllers"
	proxmoxtypes "github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

type genericModel struct {
	ID     types.String `tfsdk:"id"`
	Digest types.String `tfsdk:"digest"`
}

func (m *genericModel) fromAPI(name string, data *controllers.ControllerData, _ *diag.Diagnostics) {
	m.ID = types.StringValue(name)
	m.Digest = m.handleDeletedStringValue(data.Digest)
}

func (m *genericModel) fromAPIForDatasource(name string, data *controllers.ControllerData, _ *diag.Diagnostics) {
	m.ID = types.StringValue(name)
	m.Digest = attribute.StringValueFromPtr(data.Digest)
}

func (m *genericModel) toAPI(_ context.Context, _ *diag.Diagnostics) *controllers.Controller {
	data := &controllers.Controller{}

	data.ID = m.ID.ValueString()

	return data
}

func (m *genericModel) handleDeletedStringValue(value *string) types.String {
	if value == nil {
		return types.StringNull()
	}

	if *value == "deleted" {
		return types.StringNull()
	}

	return attribute.StringValueFromPtr(value)
}

func (m *genericModel) handleDeletedStringSetValue(value []string, diags *diag.Diagnostics) stringset.Value {
	if value == nil {
		return stringset.NullValue()
	}

	if len(value) == 1 && value[0] == "deleted" {
		return stringset.NullValue()
	}

	return stringset.NewValueList(value, diags)
}

func checkDeletedFields(_, _ *genericModel) []string {
	var toDelete []string

	return toDelete
}

func (m *genericModel) getID() string {
	return m.ID.ValueString()
}

func genericAttributesWith(extraAttributes map[string]schema.Attribute) map[string]schema.Attribute {
	// Start with generic attributes as the base
	result := map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The SDN controller object identifier.",
			Required:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
			Validators: validators.SDNID(),
		},
		"digest": schema.StringAttribute{
			Description: "Digest of the controller section.",
			Computed:    true,
		},
	}

	// Add extra attributes, allowing them to override generic ones if needed
	if extraAttributes != nil {
		maps.Copy(result, extraAttributes)
	}

	return result
}

type controllerModel interface {
	fromAPI(name string, data *controllers.ControllerData, diags *diag.Diagnostics)
	fromAPIForDatasource(name string, data *controllers.ControllerData, diags *diag.Diagnostics)
	toAPI(ctx context.Context, diags *diag.Diagnostics) *controllers.Controller
	getID() string
	getGenericModel() *genericModel
}

type controllerResourceConfig struct {
	typeNameSuffix string
	controllerType string
	modelFunc      func() controllerModel
}

type genericControllerResource struct {
	client *controllers.Client
	config controllerResourceConfig
}

func newGenericControllerResource(cfg controllerResourceConfig) resource.Resource {
	return &genericControllerResource{config: cfg}
}

func (r *genericControllerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "proxmox" + r.config.typeNameSuffix
}

func (r *genericControllerResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = cfg.Client.Cluster().SDNControllers()
}

func (r *genericControllerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := r.config.modelFunc()
	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	newController := plan.toAPI(ctx, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	newController.Type = new(r.config.controllerType)

	if err := r.client.CreateController(ctx, newController); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to Create SDN Controller %q", plan.getID()),
			err.Error(),
		)

		return
	}

	r.readAndSetState(ctx, plan.getID(), &resp.State, &resp.Diagnostics)
}

func (r *genericControllerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := r.config.modelFunc()
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	controller, err := r.client.GetControllerWithParams(ctx, state.getID(), &sdn.QueryParams{Pending: proxmoxtypes.CustomBool(true).Pointer()})
	if err != nil {
		if errors.Is(err, api.ErrResourceDoesNotExist) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to Read SDN Controller %q", state.getID()),
			err.Error(),
		)

		return
	}

	r.setModelFromController(ctx, controller, &resp.State, &resp.Diagnostics)
}

func (r *genericControllerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	plan := r.config.modelFunc()
	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)

	state := r.config.modelFunc()
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	updateController := plan.toAPI(ctx, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	toDelete := checkDeletedFields(plan.getGenericModel(), state.getGenericModel())
	update := &controllers.ControllerUpdate{
		Controller: *updateController,
		Delete:     toDelete,
	}

	if err := r.client.UpdateController(ctx, update); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to Update SDN Controller %q", state.getID()),
			err.Error(),
		)

		return
	}

	r.readAndSetState(ctx, plan.getID(), &resp.State, &resp.Diagnostics)
}

func (r *genericControllerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	state := r.config.modelFunc()
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteController(ctx, state.getID()); err != nil &&
		!errors.Is(err, api.ErrResourceDoesNotExist) {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to Delete SDN Controller %q", state.getID()),
			err.Error(),
		)

		return
	}
}

func (r *genericControllerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	controller, err := r.client.GetControllerWithParams(ctx, req.ID, &sdn.QueryParams{Pending: proxmoxtypes.CustomBool(true).Pointer()})
	if err != nil {
		if errors.Is(err, api.ErrResourceDoesNotExist) {
			resp.Diagnostics.AddError(fmt.Sprintf("Controller %s does not exist", req.ID), err.Error())
			return
		}

		resp.Diagnostics.AddError(fmt.Sprintf("Unable to Import SDN Controller %s", req.ID), err.Error())

		return
	}

	if controller.Type != nil && *controller.Type != r.config.controllerType {
		resp.Diagnostics.AddError(
			"SDN Controller Type Mismatch",
			fmt.Sprintf("Expected controller type %q but found %q for id %q",
				r.config.controllerType, *controller.Type, req.ID),
		)

		return
	}

	r.setModelFromController(ctx, controller, &resp.State, &resp.Diagnostics)
}

// Schema is required to satisfy the resource.Resource interface. It should be implemented by the specific resource.
func (r *genericControllerResource) Schema(_ context.Context, _ resource.SchemaRequest, _ *resource.SchemaResponse) {
	// Intentionally left blank. Should be set by the specific resource.
}

func (r *genericControllerResource) readAndSetState(ctx context.Context, controllerID string, state *tfsdk.State, diags *diag.Diagnostics) {
	controller, err := r.client.GetControllerWithParams(ctx, controllerID, &sdn.QueryParams{Pending: proxmoxtypes.CustomBool(true).Pointer()})
	if err != nil {
		diags.AddError(fmt.Sprintf("Unable to Read SDN Controller %q", controllerID), err.Error())
		return
	}

	r.setModelFromController(ctx, controller, state, diags)
}

func (r *genericControllerResource) setModelFromController(
	ctx context.Context,
	controller *controllers.ControllerData,
	state *tfsdk.State,
	diags *diag.Diagnostics,
) {
	readModel := r.config.modelFunc()
	readModel.fromAPI(controller.ID, controller, diags)

	if diags.HasError() {
		return
	}

	diags.Append(state.Set(ctx, readModel)...)
}
