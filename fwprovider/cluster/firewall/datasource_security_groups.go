/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package firewall

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/attribute"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/config"
	clusterfirewall "github.com/bpg/terraform-provider-proxmox/proxmox/cluster/firewall"
)

var (
	_ datasource.DataSource              = &securityGroupsDataSource{}
	_ datasource.DataSourceWithConfigure = &securityGroupsDataSource{}
)

type securityGroupsDataSource struct {
	client clusterfirewall.API
}

type securityGroupsDataSourceModel struct {
	Groups []securityGroupsDataSourceGroup `tfsdk:"groups"`
}

type securityGroupsDataSourceGroup struct {
	Name    types.String `tfsdk:"name"`
	Comment types.String `tfsdk:"comment"`
}

// NewSecurityGroupsDataSource creates the proxmox_cluster_firewall_security_groups data source.
func NewSecurityGroupsDataSource() datasource.DataSource {
	return &securityGroupsDataSource{}
}

func (d *securityGroupsDataSource) Metadata(
	_ context.Context,
	_ datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = "proxmox_cluster_firewall_security_groups"
}

func (d *securityGroupsDataSource) Configure(
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

	d.client = cfg.Client.Cluster().Firewall()
}

func (d *securityGroupsDataSource) Schema(
	_ context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Retrieves all cluster-level firewall security groups.",
		Attributes: map[string]schema.Attribute{
			"groups": schema.ListNestedAttribute{
				Description: "Security groups, sorted by name.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "Security group name.",
							Computed:    true,
						},
						"comment": schema.StringAttribute{
							Description: "Security group comment. Empty when not set.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *securityGroupsDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	groups, err := d.client.ListGroups(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Unable to List Security Groups", err.Error())
		return
	}

	state := securityGroupsDataSourceModel{
		Groups: make([]securityGroupsDataSourceGroup, len(groups)),
	}

	for i, g := range groups {
		state.Groups[i] = securityGroupsDataSourceGroup{
			Name:    types.StringValue(g.Group),
			Comment: attribute.StringValueFromPtr(g.Comment),
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
