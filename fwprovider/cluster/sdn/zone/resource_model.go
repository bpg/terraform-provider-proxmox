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
	// m.Bridge = types.StringPointerValue(data.Bridge)
	// m.ServiceVLAN = types.Int64PointerValue(data.ServiceVLAN)
	// m.ServiceVLANProtocol = types.StringPointerValue(data.ServiceVLANProtocol)
	// m.Peers = stringset.NewValueString(data.Peers, diags, comaSeparated)
	// m.Controller = types.StringPointerValue(data.Controller)
	// m.ExitNodes = stringset.NewValueString(data.ExitNodes, diags, comaSeparated)
	// m.PrimaryExitNode = types.StringPointerValue(data.ExitNodesPrimary)
	// m.RouteTargetImport = types.StringPointerValue(data.RouteTargetImport)
	// m.VRFVXLANID = types.Int64PointerValue(data.VRFVXLANID)
	// m.ExitNodesLocalRouting = types.BoolPointerValue(ptrConversion.Int64ToBoolPtr(data.ExitNodesLocalRouting))
	// m.AdvertiseSubnets = types.BoolPointerValue(ptrConversion.Int64ToBoolPtr(data.AdvertiseSubnets))
	// m.DisableARPNDSuppression = types.BoolPointerValue(ptrConversion.Int64ToBoolPtr(data.DisableARPNDSuppression))
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
	// data.Bridge = m.Bridge.ValueStringPointer()
	// data.ServiceVLAN = m.ServiceVLAN.ValueInt64Pointer()
	// data.ServiceVLANProtocol = m.ServiceVLANProtocol.ValueStringPointer()
	// data.Peers = m.Peers.ValueStringPointer(ctx, diags, comaSeparated)
	// data.Controller = m.Controller.ValueStringPointer()
	// data.ExitNodes = m.ExitNodes.ValueStringPointer(ctx, diags, comaSeparated)
	// data.ExitNodesPrimary = m.PrimaryExitNode.ValueStringPointer()
	// data.RouteTargetImport = m.RouteTargetImport.ValueStringPointer()
	// data.VRFVXLANID = m.VRFVXLANID.ValueInt64Pointer()
	// data.ExitNodesLocalRouting = ptrConversion.BoolToInt64Ptr(m.ExitNodesLocalRouting.ValueBoolPointer())
	// data.AdvertiseSubnets = ptrConversion.BoolToInt64Ptr(m.AdvertiseSubnets.ValueBoolPointer())
	// data.DisableARPNDSuppression = ptrConversion.BoolToInt64Ptr(m.DisableARPNDSuppression.ValueBoolPointer())

	return data
}
