/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package ha

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/attribute"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/config"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/migration"
	hagroups "github.com/bpg/terraform-provider-proxmox/proxmox/cluster/ha/groups"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &haGroupsDatasource{}
	_ datasource.DataSourceWithConfigure = &haGroupsDatasource{}
)

// NewHAGroupsDataSource is a helper function to simplify the provider implementation.
func NewHAGroupsDataSource() datasource.DataSource {
	return &haGroupsDatasource{}
}

// haGroupsDatasource is the data source implementation for High Availability groups.
type haGroupsDatasource struct {
	client *hagroups.Client
}

// haGroupsModel maps the schema data for the High Availability groups data source.
type haGroupsModel struct {
	Groups types.Set    `tfsdk:"group_ids"`
	ID     types.String `tfsdk:"id"`
}

// Metadata returns the data source type name.
func (d *haGroupsDatasource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_hagroups"
}

// Schema returns the schema for the data source.
func (d *haGroupsDatasource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:        "Retrieves the list of High Availability groups.",
		DeprecationMessage: migration.DeprecationMessage("proxmox_hagroups"),
		Attributes: map[string]schema.Attribute{
			"id": attribute.ResourceID(),
			"group_ids": schema.SetAttribute{
				Description: "The identifiers of the High Availability groups.",
				ElementType: types.StringType,
				Computed:    true,
			},
		},
	}
}

// Configure adds the provider-configured client to the data source.
func (d *haGroupsDatasource) Configure(
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
			"Unexpected DataSource Configure Type",
			fmt.Sprintf("Expected config.DataSource, got: %T", req.ProviderData),
		)

		return
	}

	d.client = cfg.Client.Cluster().HA().Groups()
}

// Read fetches the list of HA groups from the Proxmox cluster then converts it to a list of strings.
func (d *haGroupsDatasource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state haGroupsModel

	list, err := d.client.List(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read High Availability groups",
			err.Error(),
		)

		return
	}

	groups := make([]attr.Value, len(list))
	for i, v := range list {
		groups[i] = types.StringValue(v.ID)
	}

	groupsValue, diags := types.SetValue(types.StringType, groups)
	resp.Diagnostics.Append(diags...)

	state.ID = types.StringValue("hagroups")
	state.Groups = groupsValue

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Short-name alias for proxmox_hagroups data source (ADR-007).

var (
	_ datasource.DataSource              = &haGroupsDSShort{}
	_ datasource.DataSourceWithConfigure = &haGroupsDSShort{}
)

type haGroupsDSShort struct{ haGroupsDatasource }

// NewHAGroupsShortDataSource creates the short-name version of the HA groups data source.
func NewHAGroupsShortDataSource() datasource.DataSource {
	return &haGroupsDSShort{}
}

func (d *haGroupsDSShort) Metadata(
	_ context.Context,
	_ datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = "proxmox_hagroups"
}

func (d *haGroupsDSShort) Schema(
	ctx context.Context,
	req datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	d.haGroupsDatasource.Schema(ctx, req, resp)
	resp.Schema.DeprecationMessage = ""
}
