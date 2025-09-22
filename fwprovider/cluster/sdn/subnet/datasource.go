/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package subnet

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/config"
	customtypes "github.com/bpg/terraform-provider-proxmox/fwprovider/types"
	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/sdn/subnets"
)

var _ datasource.DataSource = &DataSource{}

var _ datasource.DataSourceWithConfigure = &DataSource{}

type DataSource struct {
	client *cluster.Client
}

func NewDataSource() datasource.DataSource {
	return &DataSource{}
}

func (d *DataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sdn_subnet"
}

func (d *DataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	cfg, ok := req.ProviderData.(config.DataSource)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Provider Configuration",
			fmt.Sprintf("Expected config.DataSource, got: %T", req.ProviderData),
		)

		return
	}

	d.client = cfg.Client.Cluster()
}

func (d *DataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieve details about a specific SDN Subnet in Proxmox VE.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The full ID in the format 'vnet-id/subnet-id'.",
			},
			"subnet": schema.StringAttribute{
				Required:   true,
				CustomType: customtypes.IPCIDRType{},
			},
			"canonical_name": schema.StringAttribute{
				Computed: true,
			},
			"type": schema.StringAttribute{
				Computed: true,
			},
			"vnet": schema.StringAttribute{
				Required:    true,
				Description: "The VNet this subnet belongs to.",
			},
			"dhcp_dns_server": schema.StringAttribute{
				Computed:    true,
				Description: "The DNS server used for DHCP.",
			},
			"dhcp_range": schema.SingleNestedAttribute{
				Optional:    true,
				Computed:    true,
				Description: "DHCP range (start and end IPs).",
				Attributes: map[string]schema.Attribute{
					"start_address": schema.StringAttribute{
						Computed:    true,
						Description: "Start of the DHCP range.",
					},
					"end_address": schema.StringAttribute{
						Computed:    true,
						Description: "End of the DHCP range.",
					},
				},
			},
			"dns_zone_prefix": schema.StringAttribute{
				Computed:    true,
				Description: "Prefix used for DNS zone delegation.",
			},
			"gateway": schema.StringAttribute{
				Computed:    true,
				Description: "The gateway address for the subnet.",
			},
			"snat": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether SNAT is enabled for the subnet.",
			},
		},
	}
}

type datasourceModel struct {
	ID     types.String            `tfsdk:"id"`
	VNet   types.String            `tfsdk:"vnet"`
	Subnet customtypes.IPCIDRValue `tfsdk:"subnet"`

	CanonicalName types.String    `tfsdk:"canonical_name"`
	Type          types.String    `tfsdk:"type"`
	DhcpDnsServer types.String    `tfsdk:"dhcp_dns_server"`
	DhcpRange     *dhcpRangeModel `tfsdk:"dhcp_range"`
	DnsZonePrefix types.String    `tfsdk:"dns_zone_prefix"`
	Gateway       types.String    `tfsdk:"gateway"`
	SNAT          types.Bool      `tfsdk:"snat"`
}

func (d *DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data datasourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	client := d.client.SDNVnets(data.VNet.ValueString()).Subnets()

	canonicalID, err := resolveCanonicalSubnetID(ctx, client, data.Subnet.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to Resolve SDN Subnet ID", err.Error())
		return
	}

	subnet, err := client.GetSubnet(ctx, canonicalID)
	if err != nil {
		if errors.Is(err, api.ErrResourceDoesNotExist) {
			resp.Diagnostics.AddError(
				"SDN Subnet Not Found",
				fmt.Sprintf("SDN Subnet with ID '%s' in VNet '%s' was not found",
					data.Subnet.ValueString(), data.VNet.ValueString()),
			)

			return
		}

		resp.Diagnostics.AddError("Unable to Read SDN Subnet", err.Error())

		return
	}

	state := &datasourceModel{}
	state.Subnet = data.Subnet
	state.VNet = data.VNet
	state.fromAPI(&subnet.Subnet)

	state.ID = types.StringValue(fmt.Sprintf("%s/%s", data.VNet.ValueString(), data.Subnet.ValueString()))
	state.CanonicalName = types.StringValue(subnet.ID)
	state.Type = types.StringValue("subnet")

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (m *datasourceModel) fromAPI(subnet *subnets.Subnet) {
	m.VNet = types.StringPointerValue(subnet.VNet)
	cidr := strings.SplitN(subnet.ID, "-", 2)[1]
	m.Subnet = customtypes.NewIPCIDRValue(strings.ReplaceAll(cidr, "-", "/"))

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
