/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package fwprovider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/config"
	"github.com/bpg/terraform-provider-proxmox/proxmox"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &versionDatasource{}
	_ datasource.DataSourceWithConfigure = &versionDatasource{}
)

// NewVersionDataSource is a helper function to simplify the provider implementation.
func NewVersionDataSource() datasource.DataSource {
	return &versionDatasource{}
}

// versionDatasource is the data source implementation.
type versionDatasource struct {
	client proxmox.Client
}

// versionDataSourceModel maps the data source schema data.
type versionDataSourceModel struct {
	Release      types.String `tfsdk:"release"`
	RepositoryID types.String `tfsdk:"repository_id"`
	Version      types.String `tfsdk:"version"`
	ID           types.String `tfsdk:"id"`
}

// Metadata returns the data source type name.
func (d *versionDatasource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_version"
}

// Schema defines the schema for the data source.
func (d *versionDatasource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves API version details.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Placeholder identifier attribute.",
				Computed:    true,
			},
			"release": schema.StringAttribute{
				Description: "The current Proxmox VE point release in `x.y` format.",
				Computed:    true,
			},
			"repository_id": schema.StringAttribute{
				Description: "The short git revision from which this version was build.",
				Computed:    true,
			},
			"version": schema.StringAttribute{
				Description: "The full pve-manager package version of this node.",
				Computed:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *versionDatasource) Configure(
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
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected config.DataSource, got: %T", req.ProviderData),
		)

		return
	}

	d.client = cfg.Client
}

// Read refreshes the Terraform state with the latest data.
func (d *versionDatasource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state versionDataSourceModel

	version, err := d.client.Version().Version(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Version",
			err.Error(),
		)

		return
	}

	state.Release = types.StringValue(version.Release)
	state.RepositoryID = types.StringValue(version.RepositoryID)
	state.Version = types.StringValue(version.Version.String())

	state.ID = types.StringValue("version")

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
