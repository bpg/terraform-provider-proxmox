/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package config

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/config"
	"github.com/bpg/terraform-provider-proxmox/proxmox"
)

var (
	_ datasource.DataSource              = &nodeConfigDataSource{}
	_ datasource.DataSourceWithConfigure = &nodeConfigDataSource{}
)

func NewNodeConfigDataSource() datasource.DataSource {
	return &nodeConfigDataSource{}
}

type nodeConfigDataSource struct {
	client proxmox.Client
}

type nodeConfigDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	NodeName    types.String `tfsdk:"node_name"`
	Description types.String `tfsdk:"description"`
}

func (d *nodeConfigDataSource) Metadata(
	_ context.Context,
	_ datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = "proxmox_node_config"
}

func (d *nodeConfigDataSource) Schema(
	_ context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Retrieves configuration of a Proxmox VE node.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of this resource.",
				Computed:    true,
			},
			"node_name": schema.StringAttribute{
				Description: "The name of the node.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "Description of the node. Shown in the web-interface node notes panel.",
				Computed:    true,
			},
		},
	}
}

func (d *nodeConfigDataSource) Configure(
	_ context.Context,
	req datasource.ConfigureRequest,
	resp *datasource.ConfigureResponse,
) {
	if req.ProviderData == nil {
		return
	}

	cfg, ok := req.ProviderData.(config.DataSource)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected DataSource Configure Type",
			fmt.Sprintf("Expected config.DataSource, got: %T", req.ProviderData),
		)

		return
	}

	d.client = cfg.Client
}

func (d *nodeConfigDataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var state nodeConfigDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	nodeName := state.NodeName.ValueString()

	data, err := d.client.Node(nodeName).GetConfig(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to Read Node Config %q", nodeName),
			err.Error(),
		)

		return
	}

	state.ID = types.StringValue(nodeName)

	if data.Description != nil && *data.Description != "" {
		trimmed := strings.TrimRight(*data.Description, "\n")
		state.Description = types.StringValue(trimmed)
	} else {
		state.Description = types.StringValue("")
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
