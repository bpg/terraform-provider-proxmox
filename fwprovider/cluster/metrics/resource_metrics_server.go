/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package metrics

import (
	"context"
	"errors"
	"fmt"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/attribute"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/config"
	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/metrics"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
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

func (r *metricsServerResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_metrics_server"
}

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
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Manages PVE metrics server.",
		Attributes: map[string]schema.Attribute{
			"id": attribute.ID(),
			"name": schema.StringAttribute{
				Description: "Unique name that will be ID of this metric server in PVE.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"disable": schema.BoolAttribute{
				Description: "Set this to `true` to disable this metric server.",
				Optional:    true,
				Default:     nil,
			},
			"mtu": schema.Int64Attribute{
				Description: "MTU (maximum transmission unit) for metrics transmission over UDP. " +
					"If not set, PVE default is `1500` (allowed `512` - `65536`).",
				Validators: []validator.Int64{int64validator.Between(512, 65536)},
				Optional:   true,
				Default:    nil,
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
				Description: "TCP socket timeout in seconds. If not set, PVE default is `1`.",
				Optional:    true,
				Default:     nil,
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
				Description: "An API path prefix inserted between `<host>:<port>/` and `/api2/`." +
					" Can be useful if the InfluxDB service runs behind a reverse proxy.",
				Optional: true,
				Default:  nil,
			},
			"influx_bucket": schema.StringAttribute{
				Description: "The InfluxDB bucket/db. Only necessary when using the http v2 api.",
				Optional:    true,
				Default:     nil,
			},
			"influx_db_proto": schema.StringAttribute{
				Description: "Protocol for InfluxDB. Choice is between `udp` | `http` | `https`. " +
					"If not set, PVE default is `udp`.",
				Validators: []validator.String{stringvalidator.OneOf("udp", "http", "https")},
				Optional:   true,
				Default:    nil,
			},
			"influx_max_body_size": schema.Int64Attribute{
				Description: "InfluxDB max-body-size in bytes. Requests are batched up to this " +
					"size. If not set, PVE default is `25000000`.",
				Optional: true,
				Default:  nil,
			},
			"influx_organization": schema.StringAttribute{
				Description: "The InfluxDB organization. Only necessary when using the http v2 " +
					"api. Has no meaning when using v2 compatibility api.",
				Optional: true,
				Default:  nil,
			},
			"influx_token": schema.StringAttribute{
				Description: "The InfluxDB access token. Only necessary when using the http v2 " +
					"api. If the v2 compatibility api is used, use `user:password` instead.",
				Optional:  true,
				Default:   nil,
				Sensitive: true,
			},
			"influx_verify": schema.BoolAttribute{
				Description: "Set to `false` to disable certificate verification for https " +
					"endpoints.",
				Optional: true,
				Default:  nil,
			},
			"graphite_path": schema.StringAttribute{
				Description: "Root graphite path (ex: `proxmox.mycluster.mykey`).",
				Optional:    true,
				Default:     nil,
			},
			"graphite_proto": schema.StringAttribute{
				Description: "Protocol to send graphite data. Choice is between `udp` | `tcp`. " +
					"If not set, PVE default is `udp`.",
				Validators: []validator.String{stringvalidator.OneOf("udp", "tcp")},
				Optional:   true,
				Default:    nil,
			},
		},
	}
}

func (r *metricsServerResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var state metricsServerModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.GetServer(ctx, state.ID.ValueString())
	if err != nil {
		if errors.Is(err, api.ErrResourceDoesNotExist) {
			resp.State.RemoveResource(ctx)

			return
		}

		resp.Diagnostics.AddError(
			"Unable to Refresh Resource",
			"An unexpected error occurred while attempting to refresh resource state. "+
				"Please retry the operation or report this issue to the provider developers.\n\n"+
				"Error: "+err.Error(),
		)

		return
	}

	readModel := &metricsServerModel{}
	readModel.importFromAPI(state.ID.ValueString(), data)

	resp.Diagnostics.Append(resp.State.Set(ctx, readModel)...)
}

func (r *metricsServerResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan metricsServerModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	reqData := plan.toAPIRequestBody()

	err := r.client.CreateServer(ctx, reqData)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Resource",
			"An unexpected error occurred while creating the resource create request.\n\n"+
				"Error: "+err.Error(),
		)

		return
	}

	plan.ID = plan.Name

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func checkDelete(planField, stateField attr.Value, toDelete *[]string, apiName string) {
	// we need to remove field via api field if there is value in state
	// but someone decided to use PVE default and removed value from resource
	if planField.IsNull() && !stateField.IsNull() {
		*toDelete = append(*toDelete, apiName)
	}
}

func (r *metricsServerResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan metricsServerModel

	var state metricsServerModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var toDelete []string

	checkDelete(plan.Disable, state.Disable, &toDelete, "disable")
	checkDelete(plan.MTU, state.MTU, &toDelete, "mtu")
	checkDelete(plan.Timeout, state.Timeout, &toDelete, "timeout")
	checkDelete(plan.InfluxAPIPathPrefix, state.InfluxAPIPathPrefix, &toDelete, "api-path-prefix")
	checkDelete(plan.InfluxBucket, state.InfluxBucket, &toDelete, "bucket")
	checkDelete(plan.InfluxDBProto, state.InfluxDBProto, &toDelete, "influxdbproto")
	checkDelete(plan.InfluxMaxBodySize, state.InfluxMaxBodySize, &toDelete, "max-body-size")
	checkDelete(plan.InfluxOrganization, state.InfluxOrganization, &toDelete, "organization")
	checkDelete(plan.InfluxToken, state.InfluxToken, &toDelete, "token")
	checkDelete(plan.InfluxVerify, state.InfluxVerify, &toDelete, "verify-certificate")
	checkDelete(plan.GraphitePath, state.GraphitePath, &toDelete, "path")
	checkDelete(plan.GraphiteProto, state.GraphiteProto, &toDelete, "proto")

	reqData := plan.toAPIRequestBody()
	reqData.Delete = &toDelete

	err := r.client.UpdateServer(ctx, reqData)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Update Resource",
			"An unexpected error occurred while creating the resource update request.\n\n"+
				"Error: "+err.Error(),
		)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *metricsServerResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var state metricsServerModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteServer(ctx, state.ID.ValueString())
	if err != nil {
		if errors.Is(err, api.ErrResourceDoesNotExist) {
			return
		}

		resp.Diagnostics.AddError(
			"Unable to Delete Resource",
			"An unexpected error occurred while creating the resource delete request.\n\n"+
				"Error: "+err.Error(),
		)

		return
	}
}

func (r *metricsServerResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	data, err := r.client.GetServer(ctx, req.ID)
	if err != nil {
		if errors.Is(err, api.ErrResourceDoesNotExist) {
			resp.Diagnostics.AddError(
				"Resource does not exist",
				"Resource you try to import does not exist.\n\n"+
					"Error: "+err.Error(),
			)

			return
		}

		resp.Diagnostics.AddError(
			"Unable to Import Resource",
			"An unexpected error occurred while attempting to import resource state.\n\n"+
				"Error: "+err.Error(),
		)

		return
	}

	readModel := &metricsServerModel{}
	readModel.importFromAPI(req.ID, data)

	resp.Diagnostics.Append(resp.State.Set(ctx, readModel)...)
}
