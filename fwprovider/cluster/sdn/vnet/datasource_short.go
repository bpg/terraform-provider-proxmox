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
	_ datasource.DataSource              = &sdnVnetDSShort{}
	_ datasource.DataSourceWithConfigure = &sdnVnetDSShort{}
)

type sdnVnetDSShort struct{ DataSource }

// NewShortDataSource creates the short-name alias proxmox_sdn_vnet (data source).
func NewShortDataSource() datasource.DataSource {
	return &sdnVnetDSShort{}
}

func (d *sdnVnetDSShort) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "proxmox_sdn_vnet"
}

func (d *sdnVnetDSShort) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	d.DataSource.Schema(ctx, req, resp)
	resp.Schema.DeprecationMessage = ""
}
