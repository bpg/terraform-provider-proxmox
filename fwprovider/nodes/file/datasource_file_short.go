/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package file

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

type fileDataSourceShort struct{ fileDataSource }

var (
	_ datasource.DataSource              = &fileDataSourceShort{}
	_ datasource.DataSourceWithConfigure = &fileDataSourceShort{}
)

// NewShortFileDataSource creates the short-name alias proxmox_file (data source).
func NewShortFileDataSource() datasource.DataSource {
	return &fileDataSourceShort{}
}

func (d *fileDataSourceShort) Metadata(
	_ context.Context,
	_ datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = "proxmox_file"
}

func (d *fileDataSourceShort) Schema(
	ctx context.Context,
	req datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	d.fileDataSource.Schema(ctx, req, resp)
	resp.Schema.DeprecationMessage = ""
}
