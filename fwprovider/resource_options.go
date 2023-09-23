/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package fwprovider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/structure"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/validators"

	"github.com/bpg/terraform-provider-proxmox/proxmox"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster"
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
	Console                 types.String `tfsdk:"console"`
	CrsHA                   types.String `tfsdk:"crs_ha"`
	CrsHARebalanceOnStart   types.Bool   `tfsdk:"crs_ha_rebalance_on_start"`
	Description             types.String `tfsdk:"description"`
	EmailFrom               types.String `tfsdk:"email_from"`
	HAShutdownPolicy        types.String `tfsdk:"ha_shutdown_policy"`
	HTTPProxy               types.String `tfsdk:"http_proxy"`
	Keyboard                types.String `tfsdk:"keyboard"`
	Language                types.String `tfsdk:"language"`
	MacPrefix               types.String `tfsdk:"mac_prefix"`
	MaxWorkers              types.Int64  `tfsdk:"max_workers"`
	MigrationNetwork        types.String `tfsdk:"migration_cidr"`
	MigrationType           types.String `tfsdk:"migration_type"`
}

// haData returns HA settings parameter string for API, HA settings are
// defined, otherwise empty string is returned.
func (m *clusterOptionsModel) haData() string {
	var haDataParams []string

	if !m.HAShutdownPolicy.IsNull() && m.HAShutdownPolicy.ValueString() != "" {
		haDataParams = append(haDataParams, fmt.Sprintf("shutdown_policy=%s", m.HAShutdownPolicy.ValueString()))
	}

	if len(haDataParams) > 0 {
		return strings.Join(haDataParams, ",")
	}

	return ""
}

// migrationData returns migration settings parameter string for API, if any of migration
// settings are defined, otherwise empty string is returned.
func (m *clusterOptionsModel) migrationData() string {
	var migrationDataParams []string

	if !m.MigrationType.IsNull() && m.MigrationType.ValueString() != "" {
		migrationDataParams = append(migrationDataParams, fmt.Sprintf("type=%s", m.MigrationType.ValueString()))
	}

	if !m.MigrationNetwork.IsNull() && m.MigrationNetwork.ValueString() != "" {
		migrationDataParams = append(migrationDataParams, fmt.Sprintf("network=%s", m.MigrationNetwork.ValueString()))
	}

	if len(migrationDataParams) > 0 {
		return strings.Join(migrationDataParams, ",")
	}

	return ""
}

// crsData returns cluster resource scheduling settings parameter string for API, if any of cluster resource scheduling
// settings are defined, otherwise empty string is returned.
func (m *clusterOptionsModel) crsData() string {
	var crsDataParams []string

	if !m.CrsHA.IsNull() && m.CrsHA.ValueString() != "" {
		crsDataParams = append(crsDataParams, fmt.Sprintf("ha=%s", m.CrsHA.ValueString()))
	}

	if !m.CrsHARebalanceOnStart.IsNull() {
		var haRebalanceOnStart string
		if m.CrsHARebalanceOnStart.ValueBool() {
			haRebalanceOnStart = "1"
		} else {
			haRebalanceOnStart = "0"
		}

		crsDataParams = append(crsDataParams, fmt.Sprintf("ha-rebalance-on-start=%s", haRebalanceOnStart))
	}

	if len(crsDataParams) > 0 {
		return strings.Join(crsDataParams, ",")
	}

	return ""
}

// bandwidthData returns bandwidth limit settings parameter string for API, if any of bandwidth
// limit settings are defined, otherwise empty string is returned.
func (m *clusterOptionsModel) bandwidthData() string {
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

	if len(bandwidthParams) > 0 {
		return strings.Join(bandwidthParams, ",")
	}

	return ""
}

func (m *clusterOptionsModel) toOptionsRequestBody() *cluster.OptionsRequestData {
	body := &cluster.OptionsRequestData{}

	if !m.EmailFrom.IsUnknown() {
		body.EmailFrom = m.EmailFrom.ValueStringPointer()
	}

	if !m.Keyboard.IsUnknown() {
		body.Keyboard = m.Keyboard.ValueStringPointer()
	}

	if !m.Language.IsUnknown() {
		body.Language = m.Language.ValueStringPointer()
	}

	if !m.MaxWorkers.IsUnknown() {
		body.MaxWorkers = m.MaxWorkers.ValueInt64Pointer()
	}

	if !m.Console.IsUnknown() {
		body.Console = m.Console.ValueStringPointer()
	}

	if !m.HTTPProxy.IsUnknown() {
		body.HTTPProxy = m.HTTPProxy.ValueStringPointer()
	}

	if !m.MacPrefix.IsUnknown() {
		body.MacPrefix = m.MacPrefix.ValueStringPointer()
	}

	if !m.MacPrefix.IsUnknown() {
		body.Description = m.Description.ValueStringPointer()
	}

	haData := m.haData()
	if haData != "" {
		body.HASettings = &haData
	}

	bandwidthData := m.bandwidthData()
	if bandwidthData != "" {
		body.BandwidthLimit = &bandwidthData
	}

	crsData := m.crsData()
	if crsData != "" {
		body.ClusterResourceScheduling = &crsData
	}

	migrationData := m.migrationData()
	if migrationData != "" {
		body.Migration = &migrationData
	}

	return body
}

func (m *clusterOptionsModel) importFromOptionsAPI(
	_ context.Context,
	opts *cluster.OptionsResponseData,
) error {
	m.BandwidthLimitClone = types.Int64Null()
	m.BandwidthLimitDefault = types.Int64Null()
	m.BandwidthLimitMigration = types.Int64Null()
	m.BandwidthLimitMove = types.Int64Null()
	m.BandwidthLimitRestore = types.Int64Null()

	//nolint:nestif
	if opts.BandwidthLimit != nil {
		for _, bandwidth := range strings.Split(*opts.BandwidthLimit, ",") {
			bandwidthData := strings.SplitN(bandwidth, "=", 2)
			bandwidthName := bandwidthData[0]

			bandwidthLimit, err := strconv.ParseInt(bandwidthData[1], 10, 64)
			if err != nil {
				return fmt.Errorf("failed to parse bandwidth limit: %s", *opts.BandwidthLimit)
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

	m.EmailFrom = types.StringPointerValue(opts.EmailFrom)
	m.Keyboard = types.StringPointerValue(opts.Keyboard)
	m.Language = types.StringPointerValue(opts.Language)

	if opts.MaxWorkers != nil {
		m.MaxWorkers = types.Int64PointerValue(opts.MaxWorkers.PointerInt64())
	} else {
		m.MaxWorkers = types.Int64Null()
	}

	m.Console = types.StringPointerValue(opts.Console)
	m.HTTPProxy = types.StringPointerValue(opts.HTTPProxy)
	m.MacPrefix = types.StringPointerValue(opts.MacPrefix)

	if opts.Description != nil && *opts.Description != "" {
		m.Description = types.StringPointerValue(opts.Description)
	} else {
		m.Description = types.StringNull()
	}

	if opts.HASettings != nil {
		m.HAShutdownPolicy = types.StringPointerValue(opts.HASettings.ShutdownPolicy)
	} else {
		m.HAShutdownPolicy = types.StringNull()
	}

	if opts.Migration != nil {
		m.MigrationType = types.StringPointerValue(opts.Migration.Type)
		m.MigrationNetwork = types.StringPointerValue(opts.Migration.Network)
	} else {
		m.MigrationType = types.StringNull()
		m.MigrationNetwork = types.StringNull()
	}

	if opts.ClusterResourceScheduling != nil {
		m.CrsHARebalanceOnStart = types.BoolPointerValue(opts.ClusterResourceScheduling.HaRebalanceOnStart.PointerBool())
		m.CrsHA = types.StringPointerValue(opts.ClusterResourceScheduling.HA)
	} else {
		m.CrsHARebalanceOnStart = types.BoolNull()
		m.CrsHA = types.StringNull()
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

// Schema defines the schema for the resource.
func (r *clusterOptionsResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Manages Proxmox VE Cluster Datacenter options.",
		Attributes: map[string]schema.Attribute{
			"id": structure.IDAttribute(),
			"email_from": schema.StringAttribute{
				Description: "email address to send notification from (default is root@$hostname).",
				Optional:    true,
			},
			"keyboard": schema.StringAttribute{
				Description: "Default keyboard layout for vnc server.",
				MarkdownDescription: "Default keyboard layout for vnc server. Must be `de` | " +
					"`de-ch` | `da` | `en-gb` | `en-us` | `es` | `fi` | `fr` | `fr-be` | `fr-ca` " +
					"| `fr-ch` | `hu` | `is` | `it` | `ja` | `lt` | `mk` | `nl` | `no` | `pl` | " +
					"`pt` | `pt-br` | `sv` | `sl` | `tr`.",
				Optional: true,
				Validators: []validator.String{
					validators.KeyboardLayoutValidator(),
				},
			},
			"max_workers": schema.Int64Attribute{
				Description: "Defines how many workers (per node) are maximal started on" +
					" actions like 'stopall VMs' or task from the ha-manager.",
				Optional: true,
			},
			"language": schema.StringAttribute{
				Description: "Default GUI language.",
				MarkdownDescription: "Default GUI language. Must be `ca` | `da` | `de` " +
					"| `en` | `es` | `eu` | `fa` | `fr` | `he` | `it` | `ja` | `nb` | " +
					"`nn` | `pl` | `pt_BR` | `ru` | `sl` | `sv` | `tr` | `zh_CN` | `zh_TW`.",
				Optional: true,
				Validators: []validator.String{
					validators.LanguageValidator(),
				},
			},
			"console": schema.StringAttribute{
				Description: "Select the default Console viewer.",
				MarkdownDescription: "Select the default Console viewer. " +
					"Must be `applet` | `vv`| `html5` | `xtermjs`. " +
					"You can either use the builtin java applet (VNC; deprecated and maps to html5), " +
					"an external virt-viewer compatible application (SPICE), " +
					"an HTML5 based vnc viewer (noVNC), " +
					"or an HTML5 based console client (xtermjs). " +
					"If the selected viewer is not available " +
					"(e.g. SPICE not activated for the VM), " +
					"the fallback is noVNC.",
				Optional: true,
				Validators: []validator.String{stringvalidator.OneOf([]string{
					"applet",
					"vv",
					"html5",
					"xtermjs",
				}...)},
			},
			"http_proxy": schema.StringAttribute{
				Description: "Specify external http proxy which is used for downloads.",
				MarkdownDescription: "Specify external http proxy which is used for downloads " +
					"(example: `http://username:password@host:port/`).",
				Optional: true,
			},
			"mac_prefix": schema.StringAttribute{
				Description: "Prefix for autogenerated MAC addresses.",
				Optional:    true,
			},
			"description": schema.StringAttribute{
				Description: "Datacenter description. Shown in the web-interface datacenter notes panel. " +
					"This is saved as comment inside the configuration file.",
				Optional: true,
			},
			"ha_shutdown_policy": schema.StringAttribute{
				Description: "Cluster wide HA shutdown policy.",
				MarkdownDescription: "Cluster wide HA shutdown policy (). " +
					"Must be `freeze` | `failover` | `migrate` | `conditional` (default is `conditional`).",
				Optional: true,
				Validators: []validator.String{stringvalidator.OneOf([]string{
					"freeze",
					"failover",
					"migrate",
					"conditional",
				}...)},
			},
			"migration_type": schema.StringAttribute{
				Description: "Cluster wide migration type.",
				MarkdownDescription: "Cluster wide migration type. Must be `secure` | `unsecure` " +
					"(default is `secure`).",
				Optional: true,
				Validators: []validator.String{stringvalidator.OneOf([]string{
					"secure",
					"unsecure",
				}...)},
			},
			"migration_cidr": schema.StringAttribute{
				Description: "Cluster wide migration network CIDR.",
				Optional:    true,
			},
			"crs_ha": schema.StringAttribute{
				Description:         "Cluster resource scheduling setting for HA.",
				MarkdownDescription: "Cluster resource scheduling setting for HA. Must be `static` | `basic` (default is `basic`).",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{stringvalidator.OneOf([]string{
					"static",
					"basic",
				}...)},
			},
			"crs_ha_rebalance_on_start": schema.BoolAttribute{
				Description: "Cluster resource scheduling setting for HA rebalance on start.",
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

// Create update must-existing cluster options.
func (r *clusterOptionsResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
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
			"Error creating cluster options",
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
			"Error converting cluster options to a model",
			"Could not import cluster options from API response, unexpected error: "+err.Error(),
		)

		return
	}
}

// Read reads cluster options.
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

// Update updates cluster options.
func (r *clusterOptionsResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
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

	if plan.bandwidthData() != state.bandwidthData() && plan.bandwidthData() == "" {
		toDelete = append(toDelete, "bwlimit")
	}

	if plan.crsData() != state.crsData() && plan.crsData() == "" {
		toDelete = append(toDelete, "crs")
	}

	if plan.haData() != state.haData() && plan.haData() == "" {
		toDelete = append(toDelete, "ha")
	}

	if plan.migrationData() != state.migrationData() && plan.migrationData() == "" {
		toDelete = append(toDelete, "migration")
	}

	if !plan.EmailFrom.Equal(state.EmailFrom) && plan.EmailFrom.ValueString() == "" {
		toDelete = append(toDelete, "email_from")
	}

	if !plan.Language.Equal(state.Language) && plan.Language.ValueString() == "" {
		toDelete = append(toDelete, "language")
	}

	if !plan.Console.Equal(state.Console) && plan.Console.ValueString() == "" {
		toDelete = append(toDelete, "console")
	}

	if !plan.HTTPProxy.Equal(state.HTTPProxy) && plan.HTTPProxy.ValueString() == "" {
		toDelete = append(toDelete, "http_proxy")
	}

	if !plan.MacPrefix.Equal(state.MacPrefix) && plan.MacPrefix.ValueString() == "" {
		toDelete = append(toDelete, "mac_prefix")
	}

	if !plan.Description.Equal(state.Description) && plan.Description.ValueString() == "" {
		toDelete = append(toDelete, "description")
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
			"Error updating cluster options",
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

// Delete deletes cluster options.
func (r *clusterOptionsResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var state clusterOptionsModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	var toDelete []string

	if !state.Keyboard.IsNull() && state.Keyboard.ValueString() != "" {
		toDelete = append(toDelete, "keyboard")
	}

	if state.bandwidthData() != "" {
		toDelete = append(toDelete, "bwlimit")
	}

	if state.crsData() != "" {
		toDelete = append(toDelete, "crs")
	}

	if state.haData() != "" {
		toDelete = append(toDelete, "ha")
	}

	if state.migrationData() != "" {
		toDelete = append(toDelete, "migration")
	}

	if !state.EmailFrom.IsNull() && state.EmailFrom.ValueString() != "" {
		toDelete = append(toDelete, "email_from")
	}

	if !state.Language.IsNull() && state.Language.ValueString() != "" {
		toDelete = append(toDelete, "language")
	}

	if !state.Console.IsNull() && state.Console.ValueString() != "" {
		toDelete = append(toDelete, "console")
	}

	if !state.HTTPProxy.IsNull() && state.HTTPProxy.ValueString() != "" {
		toDelete = append(toDelete, "http_proxy")
	}

	if !state.MacPrefix.IsNull() && state.MacPrefix.ValueString() != "" {
		toDelete = append(toDelete, "mac_prefix")
	}

	if !state.Description.IsNull() && state.Description.ValueString() != "" {
		toDelete = append(toDelete, "description")
	}

	if !state.MaxWorkers.IsNull() && state.MaxWorkers.ValueInt64() != 0 {
		toDelete = append(toDelete, "max_workers")
	}

	if len(toDelete) > 0 {
		d := strings.Join(toDelete, ",")
		body := &cluster.OptionsRequestData{
			Delete: &d,
		}

		err := r.client.Cluster().CreateUpdateOptions(ctx, body)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error updating cluster options",
				"Could not update cluster options, unexpected error: "+err.Error(),
			)
		}
	}
}

// ImportState imports cluster options.
func (r *clusterOptionsResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	state := clusterOptionsModel{ID: types.StringValue(req.ID)}
	r.read(ctx, &state, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	diags := resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}
