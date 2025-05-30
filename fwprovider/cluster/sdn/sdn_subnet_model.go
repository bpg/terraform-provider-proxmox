package sdn

/*
--------------------------------- Subnet Model Terraform ---------------------------------

Note: Currently in the API there are Delete and Digest options which are not available
in the UI so the choice was made to remove them temporary, waiting for a fix.
Also, it is not really in the way of working with terraform to use such parameters.
----------------------------------------------------------------------------------------
*/
import (
	"context"
	"fmt"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/helpers/ptrConversion"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/sdn/subnets"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type sdnSubnetModel struct {
	ID            types.String     `tfsdk:"id"`
	Subnet        types.String     `tfsdk:"subnet"`
	CanonicalName types.String     `tfsdk:"canonical_name"`
	Type          types.String     `tfsdk:"type"`
	Vnet          types.String     `tfsdk:"vnet"`
	DhcpDnsServer types.String     `tfsdk:"dhcp_dns_server"`
	DhcpRange     []dhcpRangeModel `tfsdk:"dhcp_range"`
	DnsZonePrefix types.String     `tfsdk:"dnszoneprefix"`
	Gateway       types.String     `tfsdk:"gateway"`
	Snat          types.Bool       `tfsdk:"snat"`
}

type dhcpRangeModel struct {
	StartAddress types.String `tfsdk:"start_address"`
	EndAddress   types.String `tfsdk:"end_address"`
}

func (m *sdnSubnetModel) importFromAPI(name string, data *subnets.SubnetData) {
	m.ID = types.StringValue(name)
	m.CanonicalName = types.StringValue(name)

	m.Type = types.StringPointerValue(data.Type)
	m.Vnet = types.StringPointerValue(data.Vnet)
	m.DhcpDnsServer = types.StringPointerValue(data.DHCPDNSServer)
	if data.DHCPRange != nil {
		var ranges []dhcpRangeModel
		for _, r := range data.DHCPRange {
			ranges = append(ranges, dhcpRangeModel{
				StartAddress: types.StringValue(r.StartAddress),
				EndAddress:   types.StringValue(r.EndAddress),
			})
		}
		m.DhcpRange = ranges
	}

	m.DnsZonePrefix = types.StringPointerValue(data.DNSZonePrefix)
	m.Gateway = types.StringPointerValue(data.Gateway)
	m.Snat = types.BoolPointerValue(ptrConversion.Int64ToBoolPtr(data.SNAT))
}

func (m *sdnSubnetModel) toAPIRequestBody() *subnets.SubnetRequestData {
	data := &subnets.SubnetRequestData{}

	// When creating the subnet it is ok to pass subnet cidr, but when updating need to pass canonical name
	if m.CanonicalName.ValueString() == "" {
		data.ID = m.Subnet.ValueString()
	} else {
		data.ID = m.CanonicalName.ValueString()
	}
	tflog.Warn(context.Background(), "TO API", map[string]any{
		"canonical name": m.CanonicalName.ValueString(),
		"ID":             m.ID.ValueString(),
	})
	data.Type = m.Type.ValueStringPointer()
	data.Vnet = m.Vnet.ValueStringPointer()
	data.DHCPDNSServer = m.DhcpDnsServer.ValueStringPointer()
	if m.DhcpRange != nil {
		var dhcpRanges []string
		for _, r := range m.DhcpRange {
			dhcpRanges = append(dhcpRanges, fmt.Sprintf("start-address=%s,end-address=%s", r.StartAddress.ValueString(), r.EndAddress.ValueString()))
		}
		data.DHCPRange = dhcpRanges
	}
	data.DNSZonePrefix = m.DnsZonePrefix.ValueStringPointer()
	data.Gateway = m.Gateway.ValueStringPointer()
	data.SNAT = ptrConversion.BoolToInt64Ptr(m.Snat.ValueBoolPointer())
	return data
}
