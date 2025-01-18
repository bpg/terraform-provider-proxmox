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
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var (
	_ resource.Resource                = &metricsServerResource{}
	_ resource.ResourceWithConfigure   = &metricsServerResource{}
	_ resource.ResourceWithImportState = &metricsServerResource{}
)

type metricsServerResource struct {
	client *metrics.Client
}

// NewMetricsServerResource creates new metrics server resource.
func NewMetricsServerResource() resource.Resource {
	return &metricsServerResource{}
}

// Metadata returns the resource type name.
func (r *metricsServerResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_metrics_server"
}

// Configure adds the provider-configured client to the resource.
func (r *metricsServerResource) Configure(
	_ context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
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

func (r *metricsServerResource) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Manages PVE metrics server.",
		Attributes: map[string]schema.Attribute{
			"id": attribute.ID(),
			"disable": schema.BoolAttribute{
				Description: "Set this to `true` to disable this metric server.",
				Optional:    true,
				Computed:    true,
			},
			"mtu": schema.Int64Attribute{
				Description: "MTU (maximum transmission unit) for metrics transmission over UDP. " +
					"If not set, PVE default is `1500` (allowed `512` - `65536`).",
				Optional: true,
				Computed: true,
			},
			"port": schema.Int64Attribute{
				Description: "Server network port.",
				Required:    true,
				Validators:  []validator.Int64{int64validator.Between(1, 65536)},
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"server": schema.StringAttribute{
				Description: "Server dns name or IP address.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"timeout": schema.Int64Attribute{
				Description: "TCP socket timeout in seconds. If not set, PVE default is `1`.)",
				Optional:    true,
				Computed:    true,
			},
			"type": schema.StringAttribute{
				Description: "Plugin type. Choice is between `graphite` | `influxdb`.",
				Required:    true,
				Validators:  []validator.String{stringvalidator.OneOf("graphite", "influxdb")},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"influx_api_path_prefix": schema.StringAttribute{
				Description: "An API path prefix inserted between '<host>:<port>/' and '/api2/'." +
					" Can be useful if the InfluxDB service runs behind a reverse proxy.",
				Optional: true,
				Computed: true,
			},
			"influx_bucket": schema.StringAttribute{
				Description: "The InfluxDB bucket/db. Only necessary when using the http v2 api.",
				Optional:    true,
				Computed:    true,
			},
			"influx_db_proto": schema.StringAttribute{
				Description: "Protocol for InfluxDB. Choice is between `udp` | `http` | `https`." +
					"If not set, PVE default is `udp`.",
				Optional: true,
				Computed: true,
			},
			"influx_max_body_size": schema.StringAttribute{
				Description: "InfluxDB max-body-size in bytes. Requests are batched up to this " +
					"size. If not set, PVE default is `25000000`.",
				Optional: true,
				Computed: true,
			},
			"influx_organization": schema.StringAttribute{
				Description: "The InfluxDB organization. Only necessary when using the http v2 " +
					"api. Has no meaning when using v2 compatibility api.",
				Optional: true,
				Computed: true,
			},
			"influx_token": schema.StringAttribute{
				Description: "The InfluxDB access token. Only necessary when using the http v2 " +
					"api. If the v2 compatibility api is used, use 'user:password' instead.",
				Optional:  true,
				Computed:  true,
				Sensitive: true,
			},
			"influx_verify": schema.BoolAttribute{
				Description: "Set to `false` to disable certificate verification for https " +
					"endpoints.",
				Optional: true,
				Computed: true,
			},
			"graphite_path": schema.StringAttribute{
				Description: "Root graphite path (ex: `proxmox.mycluster.mykey`).",
				Optional:    true,
				Computed:    true,
			},
			"graphite_proto": schema.StringAttribute{
				Description: "Protocol to send graphite data. Choice is between `udp` | `tcp`. " +
					"If not set, PVE default is `udp`.",
				Optional: true,
				Computed: true,
			},
		},
	}
}

func (r *metricsServerResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
}

func (r *metricsServerResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
}

func (r *metricsServerResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
}

func (r *metricsServerResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
}

func (r *metricsServerResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
}
