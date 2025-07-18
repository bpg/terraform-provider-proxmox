/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package zone

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/config"
	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/sdn/zones"
	"github.com/bpg/terraform-provider-proxmox/proxmox/helpers/ptr"
)

var (
	_ resource.ResourceWithConfigure   = &VXLANResource{}
	_ resource.ResourceWithImportState = &VXLANResource{}
)

type VXLANResource struct {
	client *zones.Client
}

func NewVXLANResource() resource.Resource {
	return &VXLANResource{}
}

func (r *VXLANResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_sdn_zone_vxlan"
}

func (r *VXLANResource) Configure(
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
			fmt.Sprintf(
				"Expected config.Resource, got: %T",
				req.ProviderData,
			),
		)
		return
	}

	r.client = cfg.Client.Cluster().SDNZones()
}

func (r *VXLANResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan vxlanModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	reqData := plan.toAPIRequestBody(ctx, &resp.Diagnostics)
	reqData.Type = ptr.Ptr(zones.TypeVXLAN)

	if err := r.client.CreateZone(ctx, reqData); err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create SDN VXLAN Zone",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *VXLANResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var state vxlanModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	zone, err := r.client.GetZone(ctx, state.ID.ValueString())
	if err != nil {
		if errors.Is(err, api.ErrResourceDoesNotExist) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Unable to Read SDN VXLAN Zone",
			err.Error(),
		)
		return
	}

	readModel := &vxlanModel{}
	readModel.importFromAPI(zone.ID, zone, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, readModel)...)
}

func (r *VXLANResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan vxlanModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	reqData := plan.toAPIRequestBody(ctx, &resp.Diagnostics)

	if err := r.client.UpdateZone(ctx, reqData); err != nil {
		resp.Diagnostics.AddError(
			"Unable to Update SDN VXLAN Zone",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *VXLANResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var state vxlanModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteZone(ctx, state.ID.ValueString()); err != nil &&
		!errors.Is(err, api.ErrResourceDoesNotExist) {
		resp.Diagnostics.AddError(
			"Unable to Delete SDN VXLAN Zone",
			err.Error(),
		)
	}
}

func (r *VXLANResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	zone, err := r.client.GetZone(ctx, req.ID)
	if err != nil {
		if errors.Is(err, api.ErrResourceDoesNotExist) {
			resp.Diagnostics.AddError(fmt.Sprintf("Zone %s does not exist", req.ID), err.Error())
			return
		}
		resp.Diagnostics.AddError(fmt.Sprintf("Unable to Import SDN VXLAN Zone %s", req.ID), err.Error())
		return
	}
	readModel := &vxlanModel{}
	readModel.importFromAPI(zone.ID, zone, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, readModel)...)
}
