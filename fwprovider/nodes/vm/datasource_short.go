/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package vm

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

const shortDataSourceTypeName = "proxmox_vm"

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &datasourceShort{}
	_ datasource.DataSourceWithConfigure = &datasourceShort{}
)

// datasourceShort is the short-name alias for the VM2 datasource (ADR-007).
type datasourceShort struct {
	Datasource
}

// NewShortDataSource creates a new short-named VM2 datasource.
func NewShortDataSource() datasource.DataSource {
	return &datasourceShort{}
}

// Metadata defines the short datasource type name.
func (d *datasourceShort) Metadata(
	_ context.Context,
	_ datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = shortDataSourceTypeName
}

// Schema returns the schema with no deprecation message (this is the canonical name).
func (d *datasourceShort) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	d.Datasource.Schema(ctx, req, resp)
	resp.Schema.DeprecationMessage = ""
}
