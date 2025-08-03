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
	_ datasource.DataSource              = &SimpleDataSource{}
	_ datasource.DataSourceWithConfigure = &SimpleDataSource{}
)

type SimpleDataSource struct {
	generic *genericZoneDataSource
}

func NewSimpleDataSource() datasource.DataSource {
	return &SimpleDataSource{
		generic: newGenericZoneDataSource(zoneDataSourceConfig{
			typeNameSuffix: "_sdn_zone_simple",
			zoneType:       zones.TypeSimple,
			modelFunc:      func() zoneModel { return &simpleModel{} },
		}),
	}
}

func (d *SimpleDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves information about a Simple Zone in Proxmox SDN.",
		MarkdownDescription: "Retrieves information about a Simple Zone in Proxmox SDN. It will create an isolated VNet bridge. " +
			"This bridge is not linked to a physical interface, and VM traffic is only local on each the node. " +
			"It can be used in NAT or routed setups.",
		Attributes: genericDataSourceAttributesWith(nil),
	}
}

func (d *SimpleDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	d.generic.Metadata(ctx, req, resp)
}

func (d *SimpleDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.generic.Configure(ctx, req, resp)
}

func (d *SimpleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	d.generic.Read(ctx, req, resp)
}
