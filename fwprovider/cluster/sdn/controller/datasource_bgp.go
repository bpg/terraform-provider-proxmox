/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package controller

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/types/stringset"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/sdn/controllers"
)

var (
	_ datasource.DataSource              = &BGPControllerDataSource{}
	_ datasource.DataSourceWithConfigure = &BGPControllerDataSource{}
)

type BGPControllerDataSource struct {
	generic *genericControllerDataSource
}

func NewBGPControllerDataSource() datasource.DataSource {
	return &BGPControllerDataSource{
		generic: newGenericControllerDataSource(controllerDataSourceConfig{
			typeNameSuffix: "_sdn_controller_bgp",
			controllerType: controllers.TypeBGP,
			modelFunc:      func() controllerModel { return &bgpModel{} },
		}),
	}
}

func (d *BGPControllerDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Reads the BGP controller configuration of a node in Proxmox VE.",
		MarkdownDescription: "Reads the BGP controller configuration, which configures the Free Range Routing (frr)" +
			" BGP router of a node. The data source reports the configuration desired by the SDN API," +
			" not the operational state of the BGP sessions.",
		Attributes: genericDataSourceAttributesWith(map[string]schema.Attribute{
			"asn": schema.Int64Attribute{
				Description: "Autonomous System Number of the local BGP router.",
				Computed:    true,
			},
			"node_name": schema.StringAttribute{
				Description: "Name of the node this BGP controller is bound to.",
				Computed:    true,
			},
			"peers": schema.SetAttribute{
				CustomType: stringset.Type{
					SetType: types.SetType{
						ElemType: types.StringType,
					},
				},
				Description: "Set of BGP peer IP addresses for the BGP controller.",
				Computed:    true,
			},
			"ebgp": schema.BoolAttribute{
				Description: "Whether eBGP is used (remote-as external) instead of iBGP.",
				Computed:    true,
			},
			"ebgp_multihop": schema.Int64Attribute{
				Description: "Maximum amount of hops for eBGP peers. Zero when not configured.",
				Computed:    true,
			},
			"loopback": schema.StringAttribute{
				Description: "Name of the loopback/dummy interface that provides the Router-IP.",
				Computed:    true,
			},
			"bgp_multipath_as_path_relax": schema.BoolAttribute{
				Description: "Whether different AS paths of equal length are considered for multipath computation.",
				Computed:    true,
			},
		}),
	}
}

func (d *BGPControllerDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	d.generic.Metadata(ctx, req, resp)
}

func (d *BGPControllerDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.generic.Configure(ctx, req, resp)
}

func (d *BGPControllerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	d.generic.Read(ctx, req, resp)
}
