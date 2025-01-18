/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package metrics

import (
	"context"
	"fmt"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/config"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/metrics"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
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

// Metadata returns the data source type name.
func (r *metricsServerDatasource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_metrics_server"
}

// Configure adds the provider-configured client to the data source.
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

// Schema returns the schema for the data source.
func (r *metricsServerDatasource) Schema(
	_ context.Context,
	req datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
}

// Read fetches the metrics server data from Proxmox VE.
func (r *metricsServerDatasource) Read(
	_ context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
}
