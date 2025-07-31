/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package zone

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/types/stringset"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/sdn/zones"
)

var (
	_ datasource.DataSource              = &EVPNDataSource{}
	_ datasource.DataSourceWithConfigure = &EVPNDataSource{}
)

type EVPNDataSource struct {
	generic *genericZoneDataSource
}

func NewEVPNDataSource() datasource.DataSource {
	return &EVPNDataSource{
		generic: newGenericZoneDataSource(zoneDataSourceConfig{
			typeNameSuffix: "_sdn_zone_evpn",
			zoneType:       zones.TypeEVPN,
			modelFunc:      func() zoneModel { return &evpnModel{} },
		}),
	}
}

func (d *EVPNDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves information about an EVPN Zone in Proxmox SDN.",
		MarkdownDescription: "Retrieves information about an EVPN Zone in Proxmox SDN. The EVPN zone creates a routable Layer 3 network, capable of " +
			"spanning across multiple clusters.",
		Attributes: genericDataSourceAttributesWith(map[string]schema.Attribute{
			"advertise_subnets": schema.BoolAttribute{
				Description: "Enable subnet advertisement for EVPN.",
				Computed:    true,
			},
			"controller": schema.StringAttribute{
				Description: "EVPN controller address.",
				Computed:    true,
			},
			"disable_arp_nd_suppression": schema.BoolAttribute{
				Description: "Disable ARP/ND suppression for EVPN.",
				Computed:    true,
			},
			"exit_nodes": schema.SetAttribute{
				CustomType: stringset.Type{
					SetType: types.SetType{
						ElemType: types.StringType,
					},
				},
				Description: "List of exit nodes for EVPN.",
				ElementType: types.StringType,
				Computed:    true,
			},
			"exit_nodes_local_routing": schema.BoolAttribute{
				Description: "Enable local routing for EVPN exit nodes.",
				Computed:    true,
			},
			"primary_exit_node": schema.StringAttribute{
				Description: "Primary exit node for EVPN.",
				Computed:    true,
			},
			"rt_import": schema.StringAttribute{
				Description:         "Route target import for EVPN.",
				MarkdownDescription: "Route target import for EVPN. Must be in the format '<ASN>:<number>' (e.g., '65000:65000').",
				Computed:            true,
			},
			"vrf_vxlan": schema.Int64Attribute{
				Description: "VRF VXLAN-ID used for dedicated routing interconnect between VNets. It must be different " +
					"than the VXLAN-ID of the VNets.",
				Computed: true,
			},
		}),
	}
}

func (d *EVPNDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	d.generic.Metadata(ctx, req, resp)
}

func (d *EVPNDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.generic.Configure(ctx, req, resp)
}

func (d *EVPNDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	d.generic.Read(ctx, req, resp)
}
