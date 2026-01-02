/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package fabric_node

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	customtypes "github.com/bpg/terraform-provider-proxmox/fwprovider/types"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/sdn/fabric_nodes"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/sdn/fabrics"
)

var (
	_ resource.ResourceWithConfigure   = &OSPFResource{}
	_ resource.ResourceWithImportState = &OSPFResource{}
)

type ospfModel struct {
	genericModel

	IPv4Address customtypes.IPAddrValue `tfsdk:"ip"`
}

func (m *ospfModel) fromAPI(name string, data *fabric_nodes.FabricNodeData, diags *diag.Diagnostics) {
	m.genericModel.fromAPI(name, data, diags)

	m.IPv4Address = m.handleDeletedIPAddrValue(data.IPv4Address)
}

func (m *ospfModel) toAPI(ctx context.Context, diags *diag.Diagnostics) *fabric_nodes.FabricNode {
	data := m.genericModel.toAPI(ctx, diags)

	data.IPv4Address = m.IPv4Address.ValueStringPointer()
	return data
}

func checkDeletedOspfFields(state, plan *ospfModel) []string {
	var toDelete []string

	if plan.IPv4Address.IsNull() && !state.IPv4Address.IsNull() {
		toDelete = append(toDelete, "ip")
	}

	return toDelete
}

type OSPFResource struct {
	*genericFabricNodeResource
}

func NewOSPFResource() resource.Resource {
	return &OSPFResource{
		genericFabricNodeResource: newGenericFabricNodeResource(fabricNodeResourceConfig{
			typeNameSuffix: "_sdn_fabric_node_ospf",
			fabricProtocol: fabrics.ProtocolOSPF,
			modelFunc:      func() fabricNodeModel { return &ospfModel{} },
		}).(*genericFabricNodeResource),
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

	update := &fabric_nodes.FabricNodeUpdate{
		FabricNode: *updateFabric,
		Delete:     toDelete,
	}

	client := r.client.SDNFabricNodes(plan.getGenericModel().FabricID.ValueString(), r.config.fabricProtocol)
	err := client.UpdateFabricNode(ctx, update)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating OSPF SDN Fabric Node",
			fmt.Sprintf("Could not update OSPF SDN Fabric Node %q: %v", plan.getID(), err),
		)

		return
	}

	r.readAndSetState(ctx, client, plan.getID(), &resp.State, &resp.Diagnostics)
}

func (r *OSPFResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "OSPF Fabric Node in Proxmox SDN. Fabrics in Proxmox VE SDN provide automated routing between nodes in a cluster.",
		MarkdownDescription: "OSPF Fabric Node in Proxmox SDN. Fabrics in Proxmox VE SDN provide automated routing between nodes in a cluster.",
		Attributes: genericAttributesWith(map[string]schema.Attribute{
			"ip": schema.StringAttribute{
				Description: "IPv4 address for the fabric node.",
				Required:    true,
				CustomType:  customtypes.IPAddrType{},
			},
		}),
	}
}

func (m *ospfModel) getGenericModel() *genericModel {
	return &m.genericModel
}
