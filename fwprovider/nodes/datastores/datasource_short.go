/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package datastores

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

type datasourceShort struct{ Datasource }

var (
	_ datasource.DataSource              = &datasourceShort{}
	_ datasource.DataSourceWithConfigure = &datasourceShort{}
)

// NewShortDataSource creates the short-name alias proxmox_datastores (data source).
func NewShortDataSource() datasource.DataSource {
	return &datasourceShort{}
}

func (d *datasourceShort) Metadata(
	_ context.Context,
	_ datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = "proxmox_datastores"
}

func (d *datasourceShort) Schema(
	ctx context.Context,
	req datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	d.Datasource.Schema(ctx, req, resp)
	resp.Schema.DeprecationMessage = ""
}
