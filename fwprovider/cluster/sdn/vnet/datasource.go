/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package vnet

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/config"
	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster"
)

var _ datasource.DataSource = &DataSource{}

var _ datasource.DataSourceWithConfigure = &DataSource{}

type DataSource struct {
	client *cluster.Client
}

func NewDataSource() datasource.DataSource {
	return &DataSource{}
}

func (d *DataSource) Metadata(
	ctx context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_sdn_vnet"
}

func (d *DataSource) Configure(
	ctx context.Context,
	req datasource.ConfigureRequest,
	resp *datasource.ConfigureResponse,
) {
	if req.ProviderData == nil {
		return
	}

	cfg, ok := req.ProviderData.(config.DataSource)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Provider Data",
			fmt.Sprintf("Expected config.DataSource, got: %T", req.ProviderData),
		)

		return
	}

	d.client = cfg.Client.Cluster()
}

func (d *DataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves information about an existing SDN VNet.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the SDN VNet.",
				Computed:    true,
			},
			"zone": schema.StringAttribute{
				Computed:    true,
				Description: "The zone to which this VNet belongs.",
			},
			"alias": schema.StringAttribute{
				Computed:    true,
				Description: "An optional alias for this VNet.",
			},
			"isolate_ports": schema.BoolAttribute{
				Computed:    true,
				Description: "Isolate ports within this VNet.",
			},
			"tag": schema.Int64Attribute{
				Computed:    true,
				Description: "VLAN/VXLAN tag.",
			},
			"vlan_aware": schema.BoolAttribute{
				Computed:    true,
				Description: "Allow VM VLANs to pass through this VNet.",
			},
		},
	}
}

func (d *DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config model

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	vnet, err := d.client.SDNVnets(config.ID.ValueString()).GetVnet(ctx)
	if err != nil {
		if errors.Is(err, api.ErrResourceDoesNotExist) {
			resp.Diagnostics.AddError("SDN VNet Not Found", fmt.Sprintf("SDN VNet with ID '%s' was not found", config.ID.ValueString()))
			return
		}

		resp.Diagnostics.AddError("Unable to Read SDN VNet", err.Error())

		return
	}

	state := model{}
	state.fromAPI(config.ID.ValueString(), &vnet.Vnet)
	state.ID = types.StringValue(config.ID.ValueString())

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
