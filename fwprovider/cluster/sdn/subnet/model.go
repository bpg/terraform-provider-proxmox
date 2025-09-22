/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package subnet

import (
	"strings"

	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/sdn/subnets"
	"github.com/hashicorp/terraform-plugin-framework/types"

	customtypes "github.com/bpg/terraform-provider-proxmox/fwprovider/types"
	proxmoxtypes "github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

type model struct {
	ID   types.String            `tfsdk:"id"`
	VNet types.String            `tfsdk:"vnet"`
	CIDR customtypes.IPCIDRValue `tfsdk:"cidr"`

	DhcpDnsServer types.String    `tfsdk:"dhcp_dns_server"`
	DhcpRange     *dhcpRangeModel `tfsdk:"dhcp_range"`
	DnsZonePrefix types.String    `tfsdk:"dns_zone_prefix"`
	Gateway       types.String    `tfsdk:"gateway"`
	SNAT          types.Bool      `tfsdk:"snat"`
}

type dhcpRangeModel struct {
	StartAddress customtypes.IPAddrValue `tfsdk:"start_address"`
	EndAddress   customtypes.IPAddrValue `tfsdk:"end_address"`
}

func (m *model) fromAPI(subnet *subnets.Subnet) {
	m.ID = types.StringValue(subnet.ID)
	m.VNet = types.StringPointerValue(subnet.VNet)
	cidr := strings.SplitN(subnet.ID, "-", 2)[1]
	m.CIDR = customtypes.NewIPCIDRValue(strings.ReplaceAll(cidr, "-", "/"))

	m.DhcpDnsServer = types.StringPointerValue(subnet.DHCPDNSServer)

	if len(subnet.DHCPRange) == 0 {
		m.DhcpRange = nil
	} else {
		r := subnet.DHCPRange[0]
		m.DhcpRange = &dhcpRangeModel{
			StartAddress: customtypes.NewIPAddrPointerValue(&r.StartAddress),
			EndAddress:   customtypes.NewIPAddrPointerValue(&r.EndAddress),
		}
	}

	m.DnsZonePrefix = types.StringPointerValue(subnet.DNSZonePrefix)
	m.Gateway = types.StringPointerValue(subnet.Gateway)
	m.SNAT = types.BoolPointerValue(subnet.SNAT.PointerBool())
}

func (m *model) toAPI() *subnets.Subnet {
	subnet := &subnets.Subnet{}
	subnet.VNet = m.VNet.ValueStringPointer()
	subnet.ID = m.ID.ValueString()

	subnet.DHCPDNSServer = m.DhcpDnsServer.ValueStringPointer()

	if m.DhcpRange != nil {
		subnet.DHCPRange = subnets.DHCPRange{
			{
				StartAddress: m.DhcpRange.StartAddress.ValueString(),
				EndAddress:   m.DhcpRange.EndAddress.ValueString(),
			},
		}
	}

	subnet.DNSZonePrefix = m.DnsZonePrefix.ValueStringPointer()
	subnet.Gateway = m.Gateway.ValueStringPointer()
	subnet.SNAT = proxmoxtypes.CustomBoolPtr(m.SNAT.ValueBoolPointer())

	return subnet
}
