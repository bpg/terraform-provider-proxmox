/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package fabric

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	customtypes "github.com/bpg/terraform-provider-proxmox/fwprovider/types"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/sdn/fabrics"
)

var (
	_ resource.ResourceWithConfigure   = &OSPFResource{}
	_ resource.ResourceWithImportState = &OSPFResource{}
)

type ospfModel struct {
	genericModel

	Area       types.String            `tfsdk:"area"`
	IPv4Prefix customtypes.IPCIDRValue `tfsdk:"ip_prefix"`
}

func (m *ospfModel) fromAPI(name string, data *fabrics.FabricData, diags *diag.Diagnostics) {
	m.genericModel.fromAPI(name, data, diags)

	m.Area = m.handleDeletedStringValue(data.Area)
	m.IPv4Prefix = m.handleDeletedIPCIDRValue(data.IPv4Prefix)
}

func (m *ospfModel) toAPI(ctx context.Context, diags *diag.Diagnostics) *fabrics.Fabric {
	data := m.genericModel.toAPI(ctx, diags)

	data.Area = m.Area.ValueStringPointer()
	data.IPv4Prefix = m.IPv4Prefix.ValueStringPointer()

	return data
}

func checkDeletedOspfFields(state, plan *ospfModel) []string {
	var toDelete []string

	if plan.IPv4Prefix.IsNull() && !state.IPv4Prefix.IsNull() {
		toDelete = append(toDelete, "ip_prefix")
	}

	if plan.Area.IsNull() && !state.Area.IsNull() {
		toDelete = append(toDelete, "area")
	}

	return toDelete
}

type OSPFResource struct {
	*genericFabricResource
}

func NewOSPFResource() resource.Resource {
	return &OSPFResource{
		genericFabricResource: newGenericFabricResource(fabricResourceConfig{
			typeNameSuffix: "_sdn_fabric_ospf",
			fabricProtocol: fabrics.ProtocolOSPF,
			modelFunc:      func() fabricModel { return &ospfModel{} },
		}).(*genericFabricResource),
	}
}

func (r *OSPFResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ospfModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var state ospfModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	toDelete := checkDeletedOspfFields(&state, &plan)

	updateFabric := plan.toAPI(ctx, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	update := &fabrics.FabricUpdate{
		Fabric: *updateFabric,
		Delete: toDelete,
	}

	err := r.client.UpdateFabric(ctx, update)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating OSPF SDN Fabric",
			fmt.Sprintf("Could not update OSPF SDN Fabric %q: %v", plan.ID.ValueString(), err),
		)

		return
	}

	r.readAndSetState(ctx, plan.getID(), &resp.State, &resp.Diagnostics)
}

func (r *OSPFResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "OSPF Fabric in Proxmox SDN. Fabrics in Proxmox VE SDN provide automated routing between nodes in a cluster.",
		MarkdownDescription: "OSPF Fabric in Proxmox SDN. Fabrics in Proxmox VE SDN provide automated routing between nodes in a cluster.",
		Attributes: genericAttributesWith(map[string]schema.Attribute{
			"area": schema.StringAttribute{
				Description: "OSPF area. Either a IPv4 address or a 32-bit number. Gets validated in rust.",
				Required:    true,
			},
			"ip_prefix": schema.StringAttribute{
				Description: "IPv4 prefix cidr for the fabric.",
				Required:    true,
				CustomType:  customtypes.IPCIDRType{},
			},
		}),
	}
}

func (m *ospfModel) getGenericModel() *genericModel {
	return &m.genericModel
}
