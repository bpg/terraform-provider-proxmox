/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package status

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/attribute"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/config"
	"github.com/bpg/terraform-provider-proxmox/proxmox"
	clusterceph "github.com/bpg/terraform-provider-proxmox/proxmox/cluster/ceph"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &DataSource{}
	_ datasource.DataSourceWithConfigure = &DataSource{}
)

// NewDataSource creates the proxmox_ceph_status data source.
func NewDataSource() datasource.DataSource {
	return &DataSource{}
}

// DataSource is the proxmox_ceph_status data source.
type DataSource struct {
	client proxmox.Client
}

// Metadata returns the data source type name.
func (d *DataSource) Metadata(
	_ context.Context,
	_ datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = "proxmox_ceph_status"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	_ context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Retrieves Ceph status from Proxmox. Queries the per-node endpoint when " +
			"`node_name` is set, otherwise the cluster-wide endpoint.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The Ceph cluster `fsid`.",
				Computed:    true,
			},
			"node_name": schema.StringAttribute{
				Description: "Optional node name. When set, the data source queries the per-node endpoint; " +
					"when omitted, it queries the cluster-wide endpoint.",
				Optional: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"fsid": schema.StringAttribute{
				Description: "Ceph cluster UUID.",
				Computed:    true,
			},
			"health_status": schema.StringAttribute{
				Description: "Overall cluster health: `HEALTH_OK`, `HEALTH_WARN`, or `HEALTH_ERR`.",
				Computed:    true,
			},
			"quorum_names": schema.ListAttribute{
				Description: "Monitor names currently participating in the quorum.",
				ElementType: types.StringType,
				Computed:    true,
			},
		},
	}
}

// Configure adds the provider-configured client to the data source.
func (d *DataSource) Configure(
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

	d.client = cfg.Client
}

// Read fetches Ceph status from either the cluster or node endpoint.
func (d *DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state model

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var (
		data *clusterceph.StatusResponseData
		err  error
	)

	if attribute.IsDefined(state.NodeName) {
		data, err = d.client.Node(state.NodeName.ValueString()).Ceph().GetStatus(ctx)
	} else {
		data, err = d.client.Cluster().Ceph().GetStatus(ctx)
	}

	if err != nil {
		resp.Diagnostics.AddError("Unable to Read Ceph Status", err.Error())
		return
	}

	resp.Diagnostics.Append(state.fromAPI(ctx, data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
