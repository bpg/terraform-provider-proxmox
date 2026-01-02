/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package fabric_node

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	customtypes "github.com/bpg/terraform-provider-proxmox/fwprovider/types"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/sdn/fabric_nodes"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/sdn/fabrics"
)

var (
	_ resource.ResourceWithConfigure   = &OpenFabricResource{}
	_ resource.ResourceWithImportState = &OpenFabricResource{}
)

type openFabricModel struct {
	genericModel

	IPv4Address customtypes.IPAddrValue `tfsdk:"ip"`
	IPv6Address customtypes.IPAddrValue `tfsdk:"ip6"`
}

func (m *openFabricModel) fromAPI(name string, data *fabric_nodes.FabricNodeData, diags *diag.Diagnostics) {
	m.genericModel.fromAPI(name, data, diags)

	m.IPv4Address = m.handleDeletedIPAddrValue(data.IPv4Address)
	m.IPv6Address = m.handleDeletedIPAddrValue(data.IPv6Address)
}

func (m *openFabricModel) toAPI(ctx context.Context, diags *diag.Diagnostics) *fabric_nodes.FabricNode {
	data := m.genericModel.toAPI(ctx, diags)

	data.IPv4Address = m.IPv4Address.ValueStringPointer()
	data.IPv6Address = m.IPv6Address.ValueStringPointer()

	return data
}

func checkDeletedOpenFabricFields(state, plan *openFabricModel) []string {
	var toDelete []string

	if plan.IPv4Address.IsNull() && !state.IPv4Address.IsNull() {
		toDelete = append(toDelete, "ip")
	}

	if plan.IPv6Address.IsNull() && !state.IPv6Address.IsNull() {
		toDelete = append(toDelete, "ip6")
	}

	toDelete = append(toDelete, checkDeletedFields(state.getGenericModel(), plan.getGenericModel())...)

	return toDelete
}

type OpenFabricResource struct {
	*genericFabricNodeResource
}

func NewOpenFabricResource() resource.Resource {
	return &OpenFabricResource{
		genericFabricNodeResource: newGenericFabricNodeResource(fabricNodeResourceConfig{
			typeNameSuffix: "_sdn_fabric_node_openfabric",
			fabricProtocol: fabrics.ProtocolOpenFabric,
			modelFunc:      func() fabricNodeModel { return &openFabricModel{} },
		}).(*genericFabricNodeResource),
	}
}

func (r *OpenFabricResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan openFabricModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var state openFabricModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	toDelete := checkDeletedOpenFabricFields(&state, &plan)

	updateFabricNode := plan.toAPI(ctx, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	update := &fabric_nodes.FabricNodeUpdate{
		FabricNode: *updateFabricNode,
		Delete:     toDelete,
	}

	client := r.client.SDNFabricNodes(plan.getGenericModel().FabricID.ValueString(), r.config.fabricProtocol)

	err := client.UpdateFabricNode(ctx, update)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating OpenFabric SDN Fabric Node",
			fmt.Sprintf("Could not update OpenFabric SDN Fabric Node %q: %v", plan.getID(), err),
		)

		return
	}

	// Read updated state
	r.readAndSetState(ctx, client, plan.getID(), &resp.State, &resp.Diagnostics)
}

func (r *OpenFabricResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "OpenFabric Fabric Node in Proxmox SDN. Fabrics in Proxmox VE SDN provide automated routing between nodes in a cluster.",
		MarkdownDescription: "OpenFabric Fabric Node in Proxmox SDN. Fabrics in Proxmox VE SDN provide automated routing between nodes in a cluster.",
		Attributes: genericAttributesWith(map[string]schema.Attribute{
			"ip": schema.StringAttribute{
				Description: "IPv4 address for the fabric node.",
				Optional:    true,
				CustomType:  customtypes.IPAddrType{},
			},
			"ip6": schema.StringAttribute{
				Description: "IPv6 address for the fabric node.",
				Optional:    true,
				CustomType:  customtypes.IPAddrType{},
			},
		}),
	}
}

func (r *OpenFabricResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.AtLeastOneOf(
			path.MatchRoot("ip"),
			path.MatchRoot("ip6"),
		),
	}
}

func (m *openFabricModel) getGenericModel() *genericModel {
	return &m.genericModel
}
