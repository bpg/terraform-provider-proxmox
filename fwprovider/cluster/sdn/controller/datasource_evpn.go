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
	_ datasource.DataSource              = &EVPNControllerDataSource{}
	_ datasource.DataSourceWithConfigure = &EVPNControllerDataSource{}
)

type EVPNControllerDataSource struct {
	generic *genericControllerDataSource
}

func NewEVPNControllerDataSource() datasource.DataSource {
	return &EVPNControllerDataSource{
		generic: newGenericControllerDataSource(controllerDataSourceConfig{
			typeNameSuffix: "_sdn_controller_evpn",
			controllerType: controllers.TypeEVPN,
			modelFunc:      func() controllerModel { return &evpnModel{} },
		}),
	}
}

func (d *EVPNControllerDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The EVPN controller plugin configures the Free Range Routing (frr) router.",
		MarkdownDescription: "The EVPN zone requires an external controller to manage the control plane." +
			" The EVPN controller plugin configures the Free Range Routing (frr) router.",
		Attributes: genericDataSourceAttributesWith(map[string]schema.Attribute{
			"asn": schema.Int64Attribute{
				Description: "Autonomous System Number for the EVPN controller.",
				Computed:    true,
			},
			"fabric": schema.StringAttribute{
				Description: "ID of the fabric this EVPN controller belongs to.",
				Computed:    true,
			},
			"peers": schema.SetAttribute{
				CustomType: stringset.Type{
					SetType: types.SetType{
						ElemType: types.StringType,
					},
				},
				Description: "Set of BGP peer IP addresses for the EVPN controller.",
				Computed:    true,
			},
		}),
	}
}

func (d *EVPNControllerDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	d.generic.Metadata(ctx, req, resp)
}

func (d *EVPNControllerDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.generic.Configure(ctx, req, resp)
}

func (d *EVPNControllerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	d.generic.Read(ctx, req, resp)
}
