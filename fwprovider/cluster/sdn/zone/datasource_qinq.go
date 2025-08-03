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
	_ datasource.DataSource              = &QinQDataSource{}
	_ datasource.DataSourceWithConfigure = &QinQDataSource{}
)

type QinQDataSource struct {
	generic *genericZoneDataSource
}

func NewQinQDataSource() datasource.DataSource {
	return &QinQDataSource{
		generic: newGenericZoneDataSource(zoneDataSourceConfig{
			typeNameSuffix: "_sdn_zone_qinq",
			zoneType:       zones.TypeQinQ,
			modelFunc:      func() zoneModel { return &qinqModel{} },
		}),
	}
}

func (d *QinQDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves information about a QinQ Zone in Proxmox SDN.",
		MarkdownDescription: "Retrieves information about a QinQ Zone in Proxmox SDN. QinQ also known as VLAN stacking, that uses multiple layers of " +
			"VLAN tags for isolation. The QinQ zone defines the outer VLAN tag (the Service VLAN) whereas the inner " +
			"VLAN tag is defined by the VNet. Your physical network switches must support stacked VLANs for this " +
			"configuration. Due to the double stacking of tags, you need 4 more bytes for QinQ VLANs. " +
			"For example, you must reduce the MTU to 1496 if you physical interface MTU is 1500.",
		Attributes: genericDataSourceAttributesWith(map[string]schema.Attribute{
			"bridge": schema.StringAttribute{
				Description: "A local, VLAN-aware bridge that is already configured on each local node",
				Computed:    true,
			},
			"service_vlan": schema.Int64Attribute{
				Description:         "Service VLAN tag for QinQ.",
				MarkdownDescription: "Service VLAN tag for QinQ. The tag must be between `1` and `4094`.",
				Computed:            true,
			},
			"service_vlan_protocol": schema.StringAttribute{
				Description:         "Service VLAN protocol for QinQ.",
				MarkdownDescription: "Service VLAN protocol for QinQ. The protocol must be `802.1ad` or `802.1q`.",
				Computed:            true,
			},
		}),
	}
}

func (d *QinQDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	d.generic.Metadata(ctx, req, resp)
}

func (d *QinQDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.generic.Configure(ctx, req, resp)
}

func (d *QinQDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	d.generic.Read(ctx, req, resp)
}
