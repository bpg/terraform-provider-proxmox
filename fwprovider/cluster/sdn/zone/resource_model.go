/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package zone

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/types/stringset"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/sdn/zones"

	proxmoxtypes "github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

type baseModel struct {
	ID         types.String    `tfsdk:"id"`
	IPAM       types.String    `tfsdk:"ipam"`
	DNS        types.String    `tfsdk:"dns"`
	ReverseDNS types.String    `tfsdk:"reverse_dns"`
	DNSZone    types.String    `tfsdk:"dns_zone"`
	Nodes      stringset.Value `tfsdk:"nodes"`
	MTU        types.Int64     `tfsdk:"mtu"`
	// // VLAN.
	// Bridge types.String `tfsdk:"bridge"`
	// // QinQ.
	// ServiceVLAN         types.Int64  `tfsdk:"service_vlan"`
	// ServiceVLANProtocol types.String `tfsdk:"service_vlan_protocol"`
	// // VXLAN.
	// Peers stringset.Value `tfsdk:"peers"`
	// // EVPN.
	// Controller              types.String    `tfsdk:"controller"`
	// ExitNodes               stringset.Value `tfsdk:"exit_nodes"`
	// PrimaryExitNode         types.String    `tfsdk:"primary_exit_node"`
	// RouteTargetImport       types.String    `tfsdk:"rt_import"`
	// VRFVXLANID              types.Int64     `tfsdk:"vrf_vxlan"`
	// ExitNodesLocalRouting   types.Bool      `tfsdk:"exit_nodes_local_routing"`
	// AdvertiseSubnets        types.Bool      `tfsdk:"advertise_subnets"`
	// DisableARPNDSuppression types.Bool      `tfsdk:"disable_arp_nd_suppression"`
}

func (m *baseModel) importFromAPI(name string, data *zones.ZoneData, diags *diag.Diagnostics) {
	m.ID = types.StringValue(name)

	m.DNS = types.StringPointerValue(data.DNS)
	m.DNSZone = types.StringPointerValue(data.DNSZone)
	m.IPAM = types.StringPointerValue(data.IPAM)
	m.MTU = types.Int64PointerValue(data.MTU)
	m.Nodes = stringset.NewValueString(data.Nodes, diags, stringset.WithSeparator(","))
	m.ReverseDNS = types.StringPointerValue(data.ReverseDNS)
}

func (m *baseModel) toAPIRequestBody(ctx context.Context, diags *diag.Diagnostics) *zones.ZoneRequestData {
	data := &zones.ZoneRequestData{}

	data.ID = m.ID.ValueString()

	data.IPAM = m.IPAM.ValueStringPointer()
	data.DNS = m.DNS.ValueStringPointer()
	data.ReverseDNS = m.ReverseDNS.ValueStringPointer()
	data.DNSZone = m.DNSZone.ValueStringPointer()
	data.Nodes = m.Nodes.ValueStringPointer(ctx, diags, stringset.WithSeparator(","))
	data.MTU = m.MTU.ValueInt64Pointer()

	return data
}

type simpleModel struct {
	baseModel
}

type vlanModel struct {
	baseModel

	Bridge types.String `tfsdk:"bridge"`
}

func (m *vlanModel) importFromAPI(name string, data *zones.ZoneData, diags *diag.Diagnostics) {
	m.baseModel.importFromAPI(name, data, diags)

	m.Bridge = types.StringPointerValue(data.Bridge)
}

func (m *vlanModel) toAPIRequestBody(ctx context.Context, diags *diag.Diagnostics) *zones.ZoneRequestData {
	data := m.baseModel.toAPIRequestBody(ctx, diags)

	data.Bridge = m.Bridge.ValueStringPointer()

	return data
}

type qinqModel struct {
	vlanModel

	ServiceVLAN         types.Int64  `tfsdk:"service_vlan"`
	ServiceVLANProtocol types.String `tfsdk:"service_vlan_protocol"`
}

func (m *qinqModel) importFromAPI(name string, data *zones.ZoneData, diags *diag.Diagnostics) {
	m.vlanModel.importFromAPI(name, data, diags)

	m.ServiceVLAN = types.Int64PointerValue(data.ServiceVLAN)
	m.ServiceVLANProtocol = types.StringPointerValue(data.ServiceVLANProtocol)
}

func (m *qinqModel) toAPIRequestBody(ctx context.Context, diags *diag.Diagnostics) *zones.ZoneRequestData {
	data := m.vlanModel.toAPIRequestBody(ctx, diags)

	data.ServiceVLAN = m.ServiceVLAN.ValueInt64Pointer()
	data.ServiceVLANProtocol = m.ServiceVLANProtocol.ValueStringPointer()

	return data
}

type vxlanModel struct {
	baseModel

	Peers stringset.Value `tfsdk:"peers"`
}

func (m *vxlanModel) importFromAPI(name string, data *zones.ZoneData, diags *diag.Diagnostics) {
	m.baseModel.importFromAPI(name, data, diags)
	m.Peers = stringset.NewValueString(data.Peers, diags, stringset.WithSeparator(","))
}

func (m *vxlanModel) toAPIRequestBody(ctx context.Context, diags *diag.Diagnostics) *zones.ZoneRequestData {
	data := m.baseModel.toAPIRequestBody(ctx, diags)

	data.Peers = m.Peers.ValueStringPointer(ctx, diags, stringset.WithSeparator(","))

	return data
}

type evpnModel struct {
	baseModel

	AdvertiseSubnets        types.Bool      `tfsdk:"advertise_subnets"`
	Controller              types.String    `tfsdk:"controller"`
	DisableARPNDSuppression types.Bool      `tfsdk:"disable_arp_nd_suppression"`
	ExitNodes               stringset.Value `tfsdk:"exit_nodes"`
	ExitNodesLocalRouting   types.Bool      `tfsdk:"exit_nodes_local_routing"`
	PrimaryExitNode         types.String    `tfsdk:"primary_exit_node"`
	RouteTargetImport       types.String    `tfsdk:"rt_import"`
	VRFVXLANID              types.Int64     `tfsdk:"vrf_vxlan"`
}

func (m *evpnModel) importFromAPI(name string, data *zones.ZoneData, diags *diag.Diagnostics) {
	m.baseModel.importFromAPI(name, data, diags)

	m.AdvertiseSubnets = types.BoolPointerValue(data.AdvertiseSubnets.PointerBool())
	m.Controller = types.StringPointerValue(data.Controller)
	m.DisableARPNDSuppression = types.BoolPointerValue(data.DisableARPNDSuppression.PointerBool())
	m.ExitNodes = stringset.NewValueString(data.ExitNodes, diags, stringset.WithSeparator(","))
	m.ExitNodesLocalRouting = types.BoolPointerValue(data.ExitNodesLocalRouting.PointerBool())
	m.PrimaryExitNode = types.StringPointerValue(data.ExitNodesPrimary)
	m.RouteTargetImport = types.StringPointerValue(data.RouteTargetImport)
	m.VRFVXLANID = types.Int64PointerValue(data.VRFVXLANID)
}

func (m *evpnModel) toAPIRequestBody(ctx context.Context, diags *diag.Diagnostics) *zones.ZoneRequestData {
	data := m.baseModel.toAPIRequestBody(ctx, diags)

	data.AdvertiseSubnets = proxmoxtypes.CustomBoolPtr(m.AdvertiseSubnets.ValueBoolPointer())
	data.Controller = m.Controller.ValueStringPointer()
	data.DisableARPNDSuppression = proxmoxtypes.CustomBoolPtr(m.DisableARPNDSuppression.ValueBoolPointer())
	data.ExitNodes = m.ExitNodes.ValueStringPointer(ctx, diags, stringset.WithSeparator(","))
	data.ExitNodesLocalRouting = proxmoxtypes.CustomBoolPtr(m.ExitNodesLocalRouting.ValueBoolPointer())
	data.ExitNodesPrimary = m.PrimaryExitNode.ValueStringPointer()
	data.RouteTargetImport = m.RouteTargetImport.ValueStringPointer()
	data.VRFVXLANID = m.VRFVXLANID.ValueInt64Pointer()

	return data
}
