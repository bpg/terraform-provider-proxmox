/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package cluster

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/bpg/terraform-provider-proxmox/internal/tffwk"
	"github.com/bpg/terraform-provider-proxmox/proxmox"
	hagroups "github.com/bpg/terraform-provider-proxmox/proxmox/cluster/ha/groups"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &hagroupDatasource{}
	_ datasource.DataSourceWithConfigure = &hagroupDatasource{}
)

// NewHAGroupDataSource is a helper function to simplify the provider implementation.
func NewHAGroupDataSource() datasource.DataSource {
	return &hagroupDatasource{}
}

// hagroupDatasource is the data source implementation for full information about
// specific High Availability groups.
type hagroupDatasource struct {
	client *hagroups.Client
}

// hagroupModel maps the schema data for the High Availability group data source.
type hagroupModel struct {
	ID         types.String `tfsdk:"id"`
	Group      types.String `tfsdk:"group"`
	Comment    types.String `tfsdk:"comment"`
	Members    types.Map    `tfsdk:"members"`
	NoFailback types.Bool   `tfsdk:"no_failback"`
	Restricted types.Bool   `tfsdk:"restricted"`
}

// Metadata returns the data source type name.
func (d *hagroupDatasource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_hagroup"
}

// Schema returns the schema for the data source.
func (d *hagroupDatasource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves information about a specific High Availability group.",
		Attributes: map[string]schema.Attribute{
			"id": tffwk.IDAttribute(),
			"group": schema.StringAttribute{
				Description: "The identifier of the High Availability group to read.",
				Required:    true,
			},
			"comment": schema.StringAttribute{
				Description: "The comment associated with this group",
				Computed:    true,
			},
			"members": schema.MapAttribute{
				Description: "The member nodes for this group, associated with their priority or to null if no priority is set.",
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
func (d *hagroupDatasource) Configure(
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
func (d *hagroupDatasource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state hagroupModel

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

	if group.Comment != nil {
		state.Comment = types.StringValue(*group.Comment)
	} else {
		state.Comment = types.StringNull()
	}

	state.ID = types.StringValue(groupID)
	state.NoFailback = types.BoolValue(group.NoFailback != 0)
	state.Restricted = types.BoolValue(group.Restricted != 0)

	members, diags := parseHAGroupMembers(groupID, group.Nodes)
	resp.Diagnostics.Append(diags...)

	state.Members = members

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Parse the list of member nodes. The list is received from the Proxmox API as a string. It must
// be converted into a map value. Errors will be returned as Terraform diagnostics.
func parseHAGroupMembers(groupID string, nodes string) (basetypes.MapValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	membersIn := strings.Split(nodes, ",")
	membersOut := make(map[string]attr.Value)

	for _, nodeDescStr := range membersIn {
		nodeDesc := strings.Split(nodeDescStr, ":")
		if len(nodeDesc) > 2 {
			diags.AddWarning(
				"Could not parse HA group member",
				fmt.Sprintf("Received group member '%s' for HA group '%s'",
					nodeDescStr, groupID),
			)

			continue
		}

		priority := types.Int64Null()

		if len(nodeDesc) == 2 {
			prio, err := strconv.Atoi(nodeDesc[1])
			if err == nil {
				priority = types.Int64Value(int64(prio))
			} else {
				diags.AddWarning(
					"Could not parse HA group member priority",
					fmt.Sprintf("Node priority string '%s' for node %s of HA group '%s'",
						nodeDesc[1], nodeDesc[0], groupID),
				)
			}
		}

		membersOut[nodeDesc[0]] = priority
	}

	value, mbDiags := types.MapValue(types.Int64Type, membersOut)
	diags.Append(mbDiags...)

	return value, diags
}
