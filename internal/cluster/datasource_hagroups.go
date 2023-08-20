/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package cluster

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/internal/structure"
	"github.com/bpg/terraform-provider-proxmox/proxmox"
	hagroups "github.com/bpg/terraform-provider-proxmox/proxmox/cluster/ha/groups"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &hagroupsDatasource{}
	_ datasource.DataSourceWithConfigure = &hagroupsDatasource{}
)

// NewHAGroupsDataSource is a helper function to simplify the provider implementation.
func NewHAGroupsDataSource() datasource.DataSource {
	return &hagroupsDatasource{}
}

// hagroupsDatasource is the data source implementation for High Availability groups.
type hagroupsDatasource struct {
	client *hagroups.Client
}

// hagroupsModel maps the schema data for the High Availability groups data source.
type hagroupsModel struct {
	Groups types.Set    `tfsdk:"group_ids"`
	ID     types.String `tfsdk:"id"`
}

// Metadata returns the data source type name.
func (d *hagroupsDatasource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_hagroups"
}

// Schema returns the schema for the data source.
func (d *hagroupsDatasource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves the list of High Availability groups.",
		Attributes: map[string]schema.Attribute{
			"id": structure.IDAttribute(),
			"group_ids": schema.SetAttribute{
				Description: "The identifiers of the High Availability groups.",
				ElementType: types.StringType,
				Computed:    true,
			},
		},
	}
}

// Configure adds the provider-configured client to the data source.
func (d *hagroupsDatasource) Configure(
	_ context.Context,
	req datasource.ConfigureRequest,
	resp *datasource.ConfigureResponse,
) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(proxmox.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *proxmox.Client, got: %T. Please report this issue to the provider developers.",
				req.ProviderData),
		)

		return
	}

	d.client = client.Cluster().HA().Groups()
}

// Read fetches the list of HA groups from the Proxmox cluster then converts it to a list of strings.
func (d *hagroupsDatasource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state hagroupsModel

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
