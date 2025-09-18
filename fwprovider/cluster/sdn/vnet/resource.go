/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package vnet

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/config"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/validators"
	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/sdn/vnets"
)

var (
	_ resource.Resource                = &Resource{}
	_ resource.ResourceWithConfigure   = &Resource{}
	_ resource.ResourceWithImportState = &Resource{}
)

type Resource struct {
	client *cluster.Client
}

func NewResource() resource.Resource {
	return &Resource{}
}

func (r *Resource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_sdn_vnet"
}

func (r *Resource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = cfg.Client.Cluster()
}

func (r *Resource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages Proxmox VE SDN VNet.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the SDN VNet.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: validators.SDNID(),
			},
			"zone": schema.StringAttribute{
				Description: "The zone to which this VNet belongs.",
				Required:    true,
				Validators:  validators.SDNID(),
			},
			"alias": schema.StringAttribute{
				Optional:    true,
				Description: "An optional alias for this VNet.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[()._a-zA-Z0-9\s-]+$`),
						"alias must contain only alphanumeric characters, spaces, hyphens, underscores, dots, and parentheses",
					),
					stringvalidator.LengthAtMost(256),
				},
			},
			"isolate_ports": schema.BoolAttribute{
				Optional:    true,
				Description: "Isolate ports within this VNet.",
			},
			"tag": schema.Int64Attribute{
				Optional:    true,
				Description: "Tag value for VLAN/VXLAN (can't be used with other zone types).",
			},
			"vlan_aware": schema.BoolAttribute{
				Optional:    true,
				Description: "Allow VM VLANs to pass through this VNet.",
			},
		},
	}
}

func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan model

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	vnet := plan.toAPI()

	err := r.client.SDNVnets(plan.ID.ValueString()).CreateVnet(ctx, vnet)
	if err != nil {
		resp.Diagnostics.AddError("Unable to Create SDN VNet", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state model

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.SDNVnets(state.ID.ValueString()).GetVnet(ctx)
	if err != nil {
		if errors.Is(err, api.ErrResourceDoesNotExist) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError("Unable to Read SDN VNet", err.Error())

		return
	}

	readModel := &model{}
	readModel.fromAPI(state.ID.ValueString(), data)

	resp.Diagnostics.Append(resp.State.Set(ctx, readModel)...)
}

func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan model

	var state model

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var toDelete []string

	checkDelete(plan.Alias, state.Alias, &toDelete, "alias")
	checkDelete(plan.IsolatePorts, state.IsolatePorts, &toDelete, "isolate-ports")
	checkDelete(plan.Tag, state.Tag, &toDelete, "tag")
	checkDelete(plan.VlanAware, state.VlanAware, &toDelete, "vlanaware")

	vnet := plan.toAPI()
	reqData := &vnets.VNetUpdate{
		VNet:   *vnet,
		Delete: toDelete,
	}

	err := r.client.SDNVnets(plan.ID.ValueString()).UpdateVnet(ctx, reqData)
	if err != nil {
		resp.Diagnostics.AddError("Unable to Update SDN VNet", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state model

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.SDNVnets(state.ID.ValueString()).DeleteVnet(ctx)
	if err != nil && !errors.Is(err, api.ErrResourceDoesNotExist) {
		resp.Diagnostics.AddError("Unable to Delete SDN VNet", err.Error())
	}
}

func (r *Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	data, err := r.client.SDNVnets(req.ID).GetVnet(ctx)
	if err != nil {
		if errors.Is(err, api.ErrResourceDoesNotExist) {
			resp.Diagnostics.AddError("SDN VNet Not Found", fmt.Sprintf("SDN VNet with ID '%s' was not found", req.ID))
			return
		}

		resp.Diagnostics.AddError("Unable to Import SDN VNet", err.Error())

		return
	}

	readModel := &model{}
	readModel.fromAPI(req.ID, data)

	resp.Diagnostics.Append(resp.State.Set(ctx, readModel)...)
}

func checkDelete(planField, stateField attr.Value, toDelete *[]string, apiName string) {
	if planField.IsNull() && !stateField.IsNull() {
		*toDelete = append(*toDelete, apiName)
	}
}
