/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package vnet

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

var (
	_ datasource.DataSource              = &sdnVnetsDSShort{}
	_ datasource.DataSourceWithConfigure = &sdnVnetsDSShort{}
)

type sdnVnetsDSShort struct{ vnetsDataSource }

// NewShortVNetsDataSource creates the short-name alias proxmox_sdn_vnets (data source).
func NewShortVNetsDataSource() datasource.DataSource {
	return &sdnVnetsDSShort{}
}

func (d *sdnVnetsDSShort) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "proxmox_sdn_vnets"
}

func (d *sdnVnetsDSShort) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	d.vnetsDataSource.Schema(ctx, req, resp)
	resp.Schema.DeprecationMessage = ""
}
