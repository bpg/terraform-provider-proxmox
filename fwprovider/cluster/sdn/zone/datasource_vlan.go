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

	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/sdn/zones"
)

var (
	_ datasource.DataSource              = &VLANDataSource{}
	_ datasource.DataSourceWithConfigure = &VLANDataSource{}
)

type VLANDataSource struct {
	generic *genericZoneDataSource
}

func NewVLANDataSource() datasource.DataSource {
	return &VLANDataSource{
		generic: newGenericZoneDataSource(zoneDataSourceConfig{
			typeNameSuffix: "_sdn_zone_vlan",
			zoneType:       zones.TypeVLAN,
			modelFunc:      func() zoneModel { return &vlanModel{} },
		}),
	}
}

func (d *VLANDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves information about a VLAN Zone in Proxmox SDN.",
		MarkdownDescription: "Retrieves information about a VLAN Zone in Proxmox SDN. It uses an existing local Linux or OVS bridge to connect to the " +
			"node's physical interface. It uses VLAN tagging defined in the VNet to isolate the network segments. " +
			"This allows connectivity of VMs between different nodes.",
		Attributes: genericDataSourceAttributesWith(map[string]schema.Attribute{
			"bridge": schema.StringAttribute{
				Description: "Bridge interface for VLAN.",
				MarkdownDescription: "The local bridge or OVS switch, already configured on _each_ node that allows " +
					"node-to-node connection.",
				Computed: true,
			},
		}),
	}
}

func (d *VLANDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	d.generic.Metadata(ctx, req, resp)
}

func (d *VLANDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.generic.Configure(ctx, req, resp)
}

func (d *VLANDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	d.generic.Read(ctx, req, resp)
}
