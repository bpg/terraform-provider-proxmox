/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package metrics

import (
	"context"
	"fmt"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/attribute"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/config"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/metrics"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &metricsServerDatasource{}
	_ datasource.DataSourceWithConfigure = &metricsServerDatasource{}
)

type metricsServerDatasource struct {
	client *metrics.Client
}

// NewMetricsServerDatasource creates new metrics server data source.
func NewMetricsServerDatasource() datasource.DataSource {
	return &metricsServerDatasource{}
}

func (r *metricsServerDatasource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_metrics_server"
}

func (r *metricsServerDatasource) Configure(
	_ context.Context,
	req datasource.ConfigureRequest,
	resp *datasource.ConfigureResponse,
) {
	if req.ProviderData == nil {
		return
	}

	cfg, ok := req.ProviderData.(config.Resource)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *proxmox.Client, got: %T", req.ProviderData),
		)

		return
	}

	r.client = cfg.Client.Cluster().Metrics()
}

func (r *metricsServerDatasource) Schema(
	_ context.Context,
	req datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Retrieves information about a specific PVE metric server.",
		Attributes: map[string]schema.Attribute{
			"id": attribute.ID(),
			"disable": schema.BoolAttribute{
				Description: "Indicates if the metric server is disabled.",
				Computed:    true,
			},
			"port": schema.Int64Attribute{
				Description: "Server network port.",
				Computed:    true,
			},
			"server": schema.StringAttribute{
				Description: "Server dns name or IP address.",
				Computed:    true,
			},
			"type": schema.StringAttribute{
				Description: "Plugin type. Either `graphite` or `influxdb`.",
				Computed:    true,
			},
		},
	}
}

func (r *metricsServerDatasource) Read(
	_ context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
}
