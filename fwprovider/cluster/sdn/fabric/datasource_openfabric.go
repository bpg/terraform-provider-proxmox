/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package fabric

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	customtypes "github.com/bpg/terraform-provider-proxmox/fwprovider/types"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/sdn/fabrics"
)

var (
	_ datasource.DataSource              = &OpenFabricDataSource{}
	_ datasource.DataSourceWithConfigure = &OpenFabricDataSource{}
)

type OpenFabricDataSource struct {
	generic *genericFabricDataSource
}

func NewOpenFabricDataSource() datasource.DataSource {
	return &OpenFabricDataSource{
		generic: newGenericFabricDataSource(fabricDataSourceConfig{
			typeNameSuffix: "_sdn_fabric_openfabric",
			fabricProtocol: fabrics.ProtocolOpenFabric,
			modelFunc:      func() fabricModel { return &openFabricModel{} },
		}),
	}
}

func (d *OpenFabricDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "OpenFabric Fabric in Proxmox SDN. Fabrics in Proxmox VE SDN provide automated routing between nodes in a cluster.",
		MarkdownDescription: "OpenFabric Fabric in Proxmox SDN. Fabrics in Proxmox VE SDN provide automated routing between nodes in a cluster.",
		Attributes: genericDataSourceAttributesWith(map[string]schema.Attribute{
			"ip_prefix": schema.StringAttribute{
				Description: "IPv4 prefix cidr for the fabric.",
				Computed:    true,
				CustomType:  customtypes.IPCIDRType{},
			},
			"ip6_prefix": schema.StringAttribute{
				Description: "IPv6 prefix cidr for the fabric.",
				Computed:    true,
				CustomType:  customtypes.IPCIDRType{},
			},
			"csnp_interval": schema.Int64Attribute{
				Description: "The csnp_interval property for OpenFabric.",
				Computed:    true,
			},
			"hello_interval": schema.Int64Attribute{
				Description: "The hello_interval property for OpenFabric.",
				Computed:    true,
			},
		}),
	}
}

func (d *OpenFabricDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	d.generic.Metadata(ctx, req, resp)
}

func (d *OpenFabricDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.generic.Configure(ctx, req, resp)
}

func (d *OpenFabricDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	d.generic.Read(ctx, req, resp)
}
