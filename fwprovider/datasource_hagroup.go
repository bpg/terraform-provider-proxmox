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

	"github.com/bpg/terraform-provider-proxmox/fwprovider/structure"

	"github.com/bpg/terraform-provider-proxmox/proxmox"
	hagroups "github.com/bpg/terraform-provider-proxmox/proxmox/cluster/ha/groups"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &haGroupDatasource{}
	_ datasource.DataSourceWithConfigure = &haGroupDatasource{}
)

// NewHAGroupDataSource is a helper function to simplify the provider implementation.
func NewHAGroupDataSource() datasource.DataSource {
	return &haGroupDatasource{}
}

// haGroupDatasource is the data source implementation for full information about
// specific High Availability groups.
type haGroupDatasource struct {
	client *hagroups.Client
}

// Metadata returns the data source type name.
func (d *haGroupDatasource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_hagroup"
}

// Schema returns the schema for the data source.
func (d *haGroupDatasource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves information about a specific High Availability group.",
		Attributes: map[string]schema.Attribute{
			"id": structure.IDAttribute(),
			"group": schema.StringAttribute{
				Description: "The identifier of the High Availability group to read.",
				Required:    true,
			},
			"comment": schema.StringAttribute{
				Description: "The comment associated with this group",
				Computed:    true,
			},
			"nodes": schema.MapAttribute{
				Description: "The member nodes for this group. They are provided as a map, where the keys are the node " +
					"names and the values represent their priority: integers for known priorities or `null` for unset " +
					"priorities.",
				Computed:    true,
				ElementType: types.Int64Type,
			},
			"no_failback": schema.BoolAttribute{
				Description: "A flag that indicates that failing back to a higher priority node is disabled for this HA group.",
				Computed:    true,
			},
			"restricted": schema.BoolAttribute{
				Description: "A flag that indicates that other nodes may not be used to run resources associated to this HA group.",
				Computed:    true,
			},
		},
	}
}

// Configure adds the provider-configured client to the data source.
func (d *haGroupDatasource) Configure(
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
func (d *haGroupDatasource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state haGroupModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	groupID := state.Group.ValueString()

	group, err := d.client.Get(ctx, groupID)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to read High Availability group '%s'", groupID),
			err.Error(),
		)

		return
	}

	state.ID = types.StringValue(groupID)

	resp.Diagnostics.Append(state.importFromAPI(*group)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
