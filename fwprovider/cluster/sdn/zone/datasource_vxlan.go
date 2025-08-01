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
	_ datasource.DataSource              = &VXLANDataSource{}
	_ datasource.DataSourceWithConfigure = &VXLANDataSource{}
)

type VXLANDataSource struct {
	generic *genericZoneDataSource
}

func NewVXLANDataSource() datasource.DataSource {
	return &VXLANDataSource{
		generic: newGenericZoneDataSource(zoneDataSourceConfig{
			typeNameSuffix: "_sdn_zone_vxlan",
			zoneType:       zones.TypeVXLAN,
			modelFunc:      func() zoneModel { return &vxlanModel{} },
		}),
	}
}

func (d *VXLANDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves information about a VXLAN Zone in Proxmox SDN.",
		MarkdownDescription: "Retrieves information about a VXLAN Zone in Proxmox SDN. It establishes a tunnel (overlay) on top of an existing network " +
			"(underlay). This encapsulates layer 2 Ethernet frames within layer 4 UDP datagrams using the default " +
			"destination port 4789. You have to configure the underlay network yourself to enable UDP connectivity " +
			"between all peers. Because VXLAN encapsulation uses 50 bytes, the MTU needs to be 50 bytes lower than the " +
			"outgoing physical interface.",
		Attributes: genericDataSourceAttributesWith(map[string]schema.Attribute{
			"peers": schema.SetAttribute{
				CustomType: stringset.Type{
					SetType: types.SetType{
						ElemType: types.StringType,
					},
				},
				Description: "A list of IP addresses of each node in the VXLAN zone.",
				MarkdownDescription: "A list of IP addresses of each node in the VXLAN zone. " +
					"This can be external nodes reachable at this IP address. All nodes in the cluster need to be " +
					"mentioned here",
				ElementType: types.StringType,
				Computed:    true,
			},
		}),
	}
}

func (d *VXLANDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	d.generic.Metadata(ctx, req, resp)
}

func (d *VXLANDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.generic.Configure(ctx, req, resp)
}

func (d *VXLANDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	d.generic.Read(ctx, req, resp)
}
