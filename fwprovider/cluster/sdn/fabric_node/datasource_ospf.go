/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package fabric_node

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	customtypes "github.com/bpg/terraform-provider-proxmox/fwprovider/types"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/sdn/fabrics"
)

var (
	_ datasource.DataSource              = &OSPFDataSource{}
	_ datasource.DataSourceWithConfigure = &OSPFDataSource{}
)

type OSPFDataSource struct {
	generic *genericFabricNodeDataSource
}

func NewOSPFDataSource() datasource.DataSource {
	return &OSPFDataSource{
		generic: newGenericFabricNodeDataSource(fabricNodeDataSourceConfig{
			typeNameSuffix: "_sdn_fabric_node_ospf",
			fabricProtocol: fabrics.ProtocolOSPF,
			modelFunc:      func() fabricNodeModel { return &ospfModel{} },
		}),
	}
}

func (d *OSPFDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "OSPF Fabric Node in Proxmox SDN. Fabrics in Proxmox VE SDN provide automated routing between nodes in a cluster.",
		MarkdownDescription: "OSPF Fabric Node in Proxmox SDN. Fabrics in Proxmox VE SDN provide automated routing between nodes in a cluster.",
		Attributes: genericDataSourceAttributesWith(map[string]schema.Attribute{
			"ip": schema.StringAttribute{
				Description: "IPv4 address for the fabric node.",
				Computed:    true,
				CustomType:  customtypes.IPAddrType{},
			},
		}),
	}
}

func (d *OSPFDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	d.generic.Metadata(ctx, req, resp)
}

func (d *OSPFDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.generic.Configure(ctx, req, resp)
}

func (d *OSPFDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	d.generic.Read(ctx, req, resp)
}
