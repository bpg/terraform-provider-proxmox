package sdn

import (
	"github.com/bpg/terraform-provider-proxmox/fwprovider/helpers/ptrConversion"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/sdn/zones"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type sdnZoneModel struct {
	ID         types.String `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	Type       types.String `tfsdk:"type"`
	IPAM       types.String `tfsdk:"ipam"`
	DNS        types.String `tfsdk:"dns"`
	ReverseDNS types.String `tfsdk:"reversedns"`
	DNSZone    types.String `tfsdk:"dns_zone"`
	Nodes      types.String `tfsdk:"nodes"`
	MTU        types.Int64  `tfsdk:"mtu"`
	// VLAN.
	Bridge types.String `tfsdk:"bridge"`
	// QinQ.
	ServiceVLAN         types.Int64  `tfsdk:"tag"`
	ServiceVLANProtocol types.String `tfsdk:"vlan_protocol"`
	// VXLAN.
	Peers types.String `tfsdk:"peers"`
	// EVPN.
	Controller              types.String `tfsdk:"controller"`
	ExitNodes               types.String `tfsdk:"exit_nodes"`
	PrimaryExitNode         types.String `tfsdk:"primary_exit_node"`
	RouteTargetImport       types.String `tfsdk:"rt_import"`
	VRFVXLANID              types.Int64  `tfsdk:"vrf_vxlan"`
	ExitNodesLocalRouting   types.Bool   `tfsdk:"exit_nodes_local_routing"`
	AdvertiseSubnets        types.Bool   `tfsdk:"advertise_subnets"`
	DisableARPNDSuppression types.Bool   `tfsdk:"disable_arp_nd_suppression"`
}

func (m *sdnZoneModel) importFromAPI(name string, data *zones.ZoneData) {
	m.ID = types.StringValue(name)
	m.Name = types.StringValue(name)

	m.Type = types.StringPointerValue(data.Type)
	m.IPAM = types.StringPointerValue(data.IPAM)
	m.DNS = types.StringPointerValue(data.DNS)
	m.ReverseDNS = types.StringPointerValue(data.ReverseDNS)
	m.DNSZone = types.StringPointerValue(data.DNSZone)
	m.Nodes = types.StringPointerValue(data.Nodes)
	m.MTU = types.Int64PointerValue(data.MTU)
	m.Bridge = types.StringPointerValue(data.Bridge)
	m.ServiceVLAN = types.Int64PointerValue(data.ServiceVLAN)
	m.ServiceVLANProtocol = types.StringPointerValue(data.ServiceVLANProtocol)
	m.Peers = types.StringPointerValue(data.Peers)
	m.Controller = types.StringPointerValue(data.Controller)
	m.ExitNodes = types.StringPointerValue(data.ExitNodes)
	m.PrimaryExitNode = types.StringPointerValue(data.PrimaryExitNode)
	m.RouteTargetImport = types.StringPointerValue(data.RouteTargetImport)
	m.VRFVXLANID = types.Int64PointerValue(data.VRFVXLANID)
	m.ExitNodesLocalRouting = types.BoolPointerValue(ptrConversion.Int64ToBoolPtr(data.ExitNodesLocalRouting))
	m.AdvertiseSubnets = types.BoolPointerValue(ptrConversion.Int64ToBoolPtr(data.AdvertiseSubnets))
	m.DisableARPNDSuppression = types.BoolPointerValue(ptrConversion.Int64ToBoolPtr(data.DisableARPNDSuppression))
}

func (m *sdnZoneModel) toAPIRequestBody() *zones.ZoneRequestData {
	data := &zones.ZoneRequestData{}

	data.ID = m.Name.ValueString()

	data.Type = m.Type.ValueStringPointer()
	data.IPAM = m.IPAM.ValueStringPointer()
	data.DNS = m.DNS.ValueStringPointer()
	data.ReverseDNS = m.ReverseDNS.ValueStringPointer()
	data.DNSZone = m.DNSZone.ValueStringPointer()
	data.Nodes = m.Nodes.ValueStringPointer()
	data.MTU = m.MTU.ValueInt64Pointer()
	data.Bridge = m.Bridge.ValueStringPointer()
	data.ServiceVLAN = m.ServiceVLAN.ValueInt64Pointer()
	data.ServiceVLANProtocol = m.ServiceVLANProtocol.ValueStringPointer()
	data.Peers = m.Peers.ValueStringPointer()
	data.Controller = m.Controller.ValueStringPointer()
	data.ExitNodes = m.ExitNodes.ValueStringPointer()
	data.PrimaryExitNode = m.PrimaryExitNode.ValueStringPointer()
	data.RouteTargetImport = m.RouteTargetImport.ValueStringPointer()
	data.VRFVXLANID = m.VRFVXLANID.ValueInt64Pointer()
	data.ExitNodesLocalRouting = ptrConversion.BoolToInt64Ptr(m.ExitNodesLocalRouting.ValueBoolPointer())
	data.AdvertiseSubnets = ptrConversion.BoolToInt64Ptr(m.AdvertiseSubnets.ValueBoolPointer())
	data.DisableARPNDSuppression = ptrConversion.BoolToInt64Ptr(m.DisableARPNDSuppression.ValueBoolPointer())

	return data
}
