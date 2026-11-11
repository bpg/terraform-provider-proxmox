/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package controller

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/attribute"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/types/stringset"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/sdn/controllers"
)

var (
	_ resource.ResourceWithConfigure   = &EVPNResource{}
	_ resource.ResourceWithImportState = &EVPNResource{}
)

type evpnModel struct {
	genericModel

	ASNumber types.Int64     `tfsdk:"asn"`
	FabricID types.String    `tfsdk:"fabric"`
	Peers    stringset.Value `tfsdk:"peers"`
}

func (m *evpnModel) fromAPI(name string, data *controllers.ControllerData, diags *diag.Diagnostics) {
	m.genericModel.fromAPI(name, data, diags)

	m.ASNumber = types.Int64PointerValue(data.ASNumber)
	m.FabricID = m.handleDeletedStringValue(data.Fabric)

	m.Peers = stringset.NewValueCommaSeparatedList(data.Peers, diags)

	if data.Pending != nil {
		if data.Pending.ASNumber != nil {
			m.ASNumber = types.Int64PointerValue(data.Pending.ASNumber)
		}

		if data.Pending.Fabric != nil {
			m.FabricID = m.handleDeletedStringValue(data.Pending.Fabric)
		}

		if data.Pending.Peers != nil {
			m.Peers = stringset.NewValueCommaSeparatedList(data.Pending.Peers, diags)
		}
	}
}

func (m *evpnModel) fromAPIForDatasource(name string, data *controllers.ControllerData, diags *diag.Diagnostics) {
	m.genericModel.fromAPIForDatasource(name, data, diags)

	m.ASNumber = attribute.Int64ValueFromPtr(data.ASNumber)
	m.FabricID = attribute.StringValueFromPtr(data.Fabric)

	m.Peers = stringset.NewValueCommaSeparatedListOrEmpty(data.Peers, diags)
}

func (m *evpnModel) toAPI(ctx context.Context, diags *diag.Diagnostics) *controllers.Controller {
	data := m.genericModel.toAPI(ctx, diags)

	data.ASNumber = attribute.Int64PtrFromValue(m.ASNumber)
	data.Fabric = attribute.StringPtrFromValue(m.FabricID)

	data.Peers = m.Peers.ValueCommaSeparatedList(ctx, diags)

	return data
}

func checkDeletedEVPNFields(plan, state *evpnModel) []string {
	var toDelete []string

	if plan.ASNumber.IsNull() && !state.ASNumber.IsNull() {
		toDelete = append(toDelete, "asn")
	}

	if plan.FabricID.IsNull() && !state.FabricID.IsNull() {
		toDelete = append(toDelete, "fabric")
	}

	if plan.Peers.IsNull() && !state.Peers.IsNull() {
		toDelete = append(toDelete, "peers")
	}

	return toDelete
}

type EVPNResource struct {
	*genericControllerResource
}

func NewEVPNResource() resource.Resource {
	return &EVPNResource{
		genericControllerResource: newGenericControllerResource(controllerResourceConfig{
			typeNameSuffix: "_sdn_controller_evpn",
			controllerType: controllers.TypeEVPN,
			modelFunc:      func() controllerModel { return &evpnModel{} },
		}).(*genericControllerResource),
	}
}

func (r *EVPNResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan evpnModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var state evpnModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	toDelete := checkDeletedEVPNFields(&plan, &state)

	updateEVPNController := plan.toAPI(ctx, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	update := &controllers.ControllerUpdate{
		Controller: *updateEVPNController,
		Delete:     toDelete,
	}

	err := r.client.UpdateController(ctx, update)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to Update SDN Controller %q", plan.ID.ValueString()),
			err.Error(),
		)

		return
	}

	r.readAndSetState(ctx, plan.getID(), &resp.State, &resp.Diagnostics)
}

func (r *EVPNResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The EVPN controller plugin configures the Free Range Routing (frr) router.",
		MarkdownDescription: "The EVPN zone requires an external controller to manage the control plane." +
			" The EVPN controller plugin configures the Free Range Routing (frr) router.",
		Attributes: genericAttributesWith(map[string]schema.Attribute{
			"asn": schema.Int64Attribute{
				Description: "Autonomous System Number for the EVPN controller.",
				Required:    true,
			},
			"fabric": schema.StringAttribute{
				Description: "ID of the fabric this EVPN controller belongs to.",
				Optional:    true,
			},
			"peers": stringset.ResourceAttribute(
				"Set of BGP peer IP addresses for the EVPN controller.",
				"",
				stringset.WithOptionalNonComputed(),
				stringset.WithIPAddresses(),
			),
		}),
	}
}

func (r *EVPNResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.AtLeastOneOf(
			path.MatchRoot("fabric"),
			path.MatchRoot("peers"),
		),
	}
}

func (m *evpnModel) getGenericModel() *genericModel {
	return &m.genericModel
}
