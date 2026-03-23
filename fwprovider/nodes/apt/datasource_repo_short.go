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

type repositoryDataSourceShort struct{ repositoryDataSource }

var (
	_ datasource.DataSource              = &repositoryDataSourceShort{}
	_ datasource.DataSourceWithConfigure = &repositoryDataSourceShort{}
)

// NewShortRepositoryDataSource creates the short-name alias proxmox_apt_repository (data source).
func NewShortRepositoryDataSource() datasource.DataSource {
	return &repositoryDataSourceShort{}
}

func (d *repositoryDataSourceShort) Metadata(
	_ context.Context,
	_ datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = "proxmox_apt_repository"
}

func (d *repositoryDataSourceShort) Schema(
	ctx context.Context,
	req datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	d.repositoryDataSource.Schema(ctx, req, resp)
	resp.Schema.DeprecationMessage = ""
}
