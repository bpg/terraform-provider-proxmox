/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package subnet

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

var (
	_ datasource.DataSource              = &sdnSubnetDSShort{}
	_ datasource.DataSourceWithConfigure = &sdnSubnetDSShort{}
)

type sdnSubnetDSShort struct{ DataSource }

// NewShortDataSource creates the short-name alias proxmox_sdn_subnet (data source).
func NewShortDataSource() datasource.DataSource {
	return &sdnSubnetDSShort{}
}

func (d *sdnSubnetDSShort) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "proxmox_sdn_subnet"
}

func (d *sdnSubnetDSShort) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	d.DataSource.Schema(ctx, req, resp)
	resp.Schema.DeprecationMessage = ""
}
