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

	"github.com/bpg/terraform-provider-proxmox/internal/structure"
	"github.com/bpg/terraform-provider-proxmox/proxmox"
	cluster "github.com/bpg/terraform-provider-proxmox/proxmox/cluster"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &clusterOptionsResource{}
	_ resource.ResourceWithConfigure   = &clusterOptionsResource{}
	_ resource.ResourceWithImportState = &clusterOptionsResource{}
)

type clusterOptionsModel struct {
	ID                      types.String `tfsdk:"id"`
	BandwidthLimitClone     types.Int64  `tfsdk:"bandwidth_limit_clone"`
	BandwidthLimitDefault   types.Int64  `tfsdk:"bandwidth_limit_default"`
	BandwidthLimitMigration types.Int64  `tfsdk:"bandwidth_limit_migration"`
	BandwidthLimitMove      types.Int64  `tfsdk:"bandwidth_limit_move"`
	BandwidthLimitRestore   types.Int64  `tfsdk:"bandwidth_limit_restore"`
	EmailFrom               types.String `tfsdk:"email_from"`
	Keyboard                types.String `tfsdk:"keyboard"`
	Language                types.String `tfsdk:"language"`
	MaxWorkers              types.Int64  `tfsdk:"max_workers"`
}

func (m *clusterOptionsModel) bandwithData() string {
	var bandwidthParams []string

	if !m.BandwidthLimitClone.IsNull() && m.BandwidthLimitClone.ValueInt64() != 0 {
		bandwidthParams = append(bandwidthParams, fmt.Sprintf("clone=%d", m.BandwidthLimitClone.ValueInt64()))
	}

	if !m.BandwidthLimitDefault.IsNull() && m.BandwidthLimitDefault.ValueInt64() != 0 {
		bandwidthParams = append(bandwidthParams, fmt.Sprintf("default=%d", m.BandwidthLimitDefault.ValueInt64()))
	}

	if !m.BandwidthLimitMigration.IsNull() && m.BandwidthLimitMigration.ValueInt64() != 0 {
		bandwidthParams = append(bandwidthParams, fmt.Sprintf("migration=%d", m.BandwidthLimitMigration.ValueInt64()))
	}

	if !m.BandwidthLimitMove.IsNull() && m.BandwidthLimitMove.ValueInt64() != 0 {
		bandwidthParams = append(bandwidthParams, fmt.Sprintf("move=%d", m.BandwidthLimitMove.ValueInt64()))
	}

	if !m.BandwidthLimitRestore.IsNull() && m.BandwidthLimitRestore.ValueInt64() != 0 {
		bandwidthParams = append(bandwidthParams, fmt.Sprintf("restore=%d", m.BandwidthLimitRestore.ValueInt64()))
	}

	return strings.Join(bandwidthParams, ",")
}

func (m *clusterOptionsModel) toOptionsRequestBody() *cluster.OptionsRequestData {
	body := &cluster.OptionsRequestData{}

	body.EmailFrom = m.EmailFrom.ValueStringPointer()
	body.Keyboard = m.Keyboard.ValueStringPointer()
	body.Language = m.Language.ValueStringPointer()
	body.MaxWorkers = m.MaxWorkers.ValueInt64Pointer()

	bandwidth := m.bandwithData()
	body.BandwidthLimit = &bandwidth

	return body
}

func (m *clusterOptionsModel) importFromOptionsAPI(
	_ context.Context,
	iface *cluster.OptionsResponseData,
) error {
	m.BandwidthLimitClone = types.Int64Null()
	m.BandwidthLimitDefault = types.Int64Null()
	m.BandwidthLimitMigration = types.Int64Null()
	m.BandwidthLimitMove = types.Int64Null()
	m.BandwidthLimitRestore = types.Int64Null()

	//nolint:nestif
	if *iface.BandwidthLimit != "" {
		for _, bandwidth := range strings.Split(*iface.BandwidthLimit, ",") {
			bandwidthData := strings.SplitN(bandwidth, "=", 2)
			bandwidthName := bandwidthData[0]

			bandwidthLimit, err := strconv.ParseInt(bandwidthData[1], 10, 64)
			if err != nil {
				return fmt.Errorf("failed to parse bandwidth limit: %s", *iface.BandwidthLimit)
			}

			if bandwidthName == "clone" {
				m.BandwidthLimitClone = types.Int64Value(bandwidthLimit)
			}

			if bandwidthName == "default" {
				m.BandwidthLimitDefault = types.Int64Value(bandwidthLimit)
			}

			if bandwidthName == "migration" {
				m.BandwidthLimitMigration = types.Int64Value(bandwidthLimit)
			}

			if bandwidthName == "move" {
				m.BandwidthLimitMove = types.Int64Value(bandwidthLimit)
			}

			if bandwidthName == "restore" {
				m.BandwidthLimitRestore = types.Int64Value(bandwidthLimit)
			}
		}
	}

	m.EmailFrom = types.StringPointerValue(iface.EmailFrom)
	m.Keyboard = types.StringPointerValue(iface.Keyboard)
	m.Language = types.StringPointerValue(iface.Language)

	if iface.MaxWorkers != nil {
		m.MaxWorkers = types.Int64Value(int64(*iface.MaxWorkers))
	}

	return nil
}

// NewClusterOptionsResource manages cluster options resource.
func NewClusterOptionsResource() resource.Resource {
	return &clusterOptionsResource{}
}

type clusterOptionsResource struct {
	client proxmox.Client
}

func (r *clusterOptionsResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_cluster_options"
}

// // Schema defines the schema for the resource.
func (r *clusterOptionsResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Manages Proxmox VE cluster options.",
		Attributes: map[string]schema.Attribute{
			// Base attributes
			"id": structure.IDAttribute("A unique identifier"),
			"email_from": schema.StringAttribute{
				Description: "email address to send notification from (default is root@$hostname).",
				Optional:    true,
			},
			"keyboard": schema.StringAttribute{
				Description: "Default keybord layout for vnc server.",
				Optional:    true,
			},
			"max_workers": schema.Int64Attribute{
				Description: "Defines how many workers (per node) are maximal started on" +
					" actions like 'stopall VMs' or task from the ha-manager.",
				Optional: true,
			},
			"language": schema.StringAttribute{
				Description: "Default GUI language.",
				Optional:    true,
			},
			"bandwidth_limit_clone": schema.Int64Attribute{
				Description: "Clone I/O bandwidth limit in KiB/s.",
				Optional:    true,
			},
			"bandwidth_limit_default": schema.Int64Attribute{
				Description: "Default I/O bandwidth limit in KiB/s.",
				Optional:    true,
			},
			"bandwidth_limit_migration": schema.Int64Attribute{
				Description: "Migration I/O bandwidth limit in KiB/s.",
				Optional:    true,
			},
			"bandwidth_limit_move": schema.Int64Attribute{
				Description: "Move I/O bandwidth limit in KiB/s.",
				Optional:    true,
			},
			"bandwidth_limit_restore": schema.Int64Attribute{
				Description: "Restore I/O bandwidth limit in KiB/s.",
				Optional:    true,
			},
		},
	}
}

func (r *clusterOptionsResource) Configure(
	_ context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
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

	r.client = client
}

// Create update must-existing cluster options interface.
//
//nolint:lll
func (r *clusterOptionsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan clusterOptionsModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	body := plan.toOptionsRequestBody()

	err := r.client.Cluster().CreateUpdateOptions(ctx, body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating cluster options interface",
			"Could not create cluster options, unexpected error: "+err.Error(),
		)

		return
	}

	plan.ID = types.StringValue("cluster")

	r.read(ctx, &plan, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *clusterOptionsResource) read(ctx context.Context, model *clusterOptionsModel, diags *diag.Diagnostics) {
	options, err := r.client.Cluster().GetOptions(ctx)
	if err != nil {
		diags.AddError(
			"Error get cluster options",
			"Could not get cluster options, unexpected error: "+err.Error(),
		)

		return
	}

	err = model.importFromOptionsAPI(ctx, options)

	if err != nil {
		diags.AddError(
			"Error converting cluster options interface to a model",
			"Could not import cluster options from API response, unexpected error: "+err.Error(),
		)

		return
	}
}

// Read reads a cluster options interface.
func (r *clusterOptionsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state clusterOptionsModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	r.read(ctx, &state, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// Update updates a cluster options interface.
//
//nolint:lll
func (r *clusterOptionsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state clusterOptionsModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	body := plan.toOptionsRequestBody()

	var toDelete []string

	if !plan.Keyboard.Equal(state.Keyboard) && plan.Keyboard.ValueString() == "" {
		toDelete = append(toDelete, "keyboard")
	}

	if plan.bandwithData() != state.bandwithData() && plan.bandwithData() == "" {
		toDelete = append(toDelete, "bwlimit")
	}

	if !plan.EmailFrom.Equal(state.EmailFrom) && plan.EmailFrom.ValueString() == "" {
		toDelete = append(toDelete, "email_from")
	}

	if !plan.Language.Equal(state.Language) && plan.Language.ValueString() == "" {
		toDelete = append(toDelete, "language")
	}

	if !plan.MaxWorkers.Equal(state.MaxWorkers) && plan.MaxWorkers.ValueInt64() == 0 {
		toDelete = append(toDelete, "max_workers")
	}

	if len(toDelete) > 0 {
		d := strings.Join(toDelete, ",")
		body.Delete = &d
	}

	err := r.client.Cluster().CreateUpdateOptions(ctx, body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating cluster options interface",
			"Could not update cluster options, unexpected error: "+err.Error(),
		)

		return
	}

	r.read(ctx, &plan, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// Delete deletes a cluster options interface.
//
//nolint:lll
func (r *clusterOptionsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state clusterOptionsModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
}

// Imports a cluster options interface.
func (r *clusterOptionsResource) ImportState(
	ctx context.Context,
	_ resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	state := clusterOptionsModel{}
	r.read(ctx, &state, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	diags := resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}
