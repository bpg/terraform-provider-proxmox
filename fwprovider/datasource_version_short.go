/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package fwprovider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

const shortVersionDataSourceTypeName = "proxmox_version"

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &versionDatasourceShort{}
	_ datasource.DataSourceWithConfigure = &versionDatasourceShort{}
)

// versionDatasourceShort is the short-name alias for the version datasource (ADR-007).
type versionDatasourceShort struct {
	versionDatasource
}

// NewShortVersionDataSource creates a new short-named version datasource.
func NewShortVersionDataSource() datasource.DataSource {
	return &versionDatasourceShort{}
}

// Metadata defines the short datasource type name.
func (d *versionDatasourceShort) Metadata(
	_ context.Context,
	_ datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = shortVersionDataSourceTypeName
}

// Schema returns the schema with no deprecation message (this is the canonical name).
func (d *versionDatasourceShort) Schema(
	ctx context.Context,
	req datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	d.versionDatasource.Schema(ctx, req, resp)
	resp.Schema.DeprecationMessage = ""
}
