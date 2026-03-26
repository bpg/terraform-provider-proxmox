/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package apt

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

type standardRepositoryDataSourceShort struct{ standardRepositoryDataSource }

var (
	_ datasource.DataSource              = &standardRepositoryDataSourceShort{}
	_ datasource.DataSourceWithConfigure = &standardRepositoryDataSourceShort{}
)

// NewShortStandardRepositoryDataSource creates the short-name alias proxmox_apt_standard_repository (data source).
func NewShortStandardRepositoryDataSource() datasource.DataSource {
	return &standardRepositoryDataSourceShort{}
}

func (d *standardRepositoryDataSourceShort) Metadata(
	_ context.Context,
	_ datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = "proxmox_apt_standard_repository"
}

func (d *standardRepositoryDataSourceShort) Schema(
	ctx context.Context,
	req datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	d.standardRepositoryDataSource.Schema(ctx, req, resp)
	resp.Schema.DeprecationMessage = ""
}
