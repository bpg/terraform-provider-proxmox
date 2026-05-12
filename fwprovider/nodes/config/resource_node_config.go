/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package config

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/attribute"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/config"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/migration"
	"github.com/bpg/terraform-provider-proxmox/proxmox"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes"
)

var (
	_ resource.Resource                = &nodeConfigResource{}
	_ resource.ResourceWithConfigure   = &nodeConfigResource{}
	_ resource.ResourceWithImportState = &nodeConfigResource{}
)

func NewNodeConfigResource() resource.Resource {
	return &nodeConfigResource{}
}

type nodeConfigResource struct {
	client proxmox.Client
}

func (r *nodeConfigResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_node_config"
}

func (r *nodeConfigResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		DeprecationMessage: migration.DeprecationMessage("proxmox_node_config"),
		Description:        "Manages Proxmox VE node configuration.",
		Attributes: map[string]schema.Attribute{
			"id": attribute.ResourceID(),
			"node_name": schema.StringAttribute{
				Description: "The name of the node.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				Description: "Description of the node. Shown in the web-interface node notes panel. " +
					"This is saved as a comment inside the configuration file.",
				Optional: true,
			},
		},
	}
}

func (r *nodeConfigResource) Configure(
	_ context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
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

func (r *nodeConfigResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan nodeConfigModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	nodeName := plan.NodeName.ValueString()

	err := r.client.Node(nodeName).UpdateConfig(ctx, plan.toAPI())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to Create Node Config %q", nodeName),
			err.Error(),
		)

		return
	}

	plan.ID = types.StringValue(nodeName)

	r.read(ctx, &plan, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *nodeConfigResource) read(ctx context.Context, model *nodeConfigModel, diags *diag.Diagnostics) {
	nodeName := model.NodeName.ValueString()

	data, err := r.client.Node(nodeName).GetConfig(ctx)
	if err != nil {
		diags.AddError(
			fmt.Sprintf("Unable to Read Node Config %q", nodeName),
			err.Error(),
		)

		return
	}

	model.fromAPI(data)
}

func (r *nodeConfigResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var state nodeConfigModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	r.read(ctx, &state, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *nodeConfigResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan, state nodeConfigModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	nodeName := plan.NodeName.ValueString()
	body := plan.toAPI()

	var toDelete []string

	attribute.CheckDelete(plan.Description, state.Description, &toDelete, "description")

	if len(toDelete) > 0 {
		body.Delete = toDelete
	}

	err := r.client.Node(nodeName).UpdateConfig(ctx, body)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to Update Node Config %q", nodeName),
			err.Error(),
		)

		return
	}

	r.read(ctx, &plan, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *nodeConfigResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var state nodeConfigModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if !attribute.IsDefined(state.Description) {
		return
	}

	nodeName := state.NodeName.ValueString()

	err := r.client.Node(nodeName).UpdateConfig(ctx, &nodes.ConfigUpdateRequestBody{
		Delete: []string{"description"},
	})
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to Delete Node Config %q", nodeName),
			err.Error(),
		)
	}
}

func (r *nodeConfigResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	nodeName := req.ID
	state := nodeConfigModel{
		ID:       types.StringValue(nodeName),
		NodeName: types.StringValue(nodeName),
	}

	r.read(ctx, &state, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
