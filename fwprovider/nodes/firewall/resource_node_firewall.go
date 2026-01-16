/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package firewall

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/attribute"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/config"
	"github.com/bpg/terraform-provider-proxmox/proxmox"
	nodefirewall "github.com/bpg/terraform-provider-proxmox/proxmox/nodes/firewall"
	proxmoxtypes "github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

var (
	_ resource.Resource                = &nodeFirewallOptionsResource{}
	_ resource.ResourceWithConfigure   = &nodeFirewallOptionsResource{}
	_ resource.ResourceWithImportState = &nodeFirewallOptionsResource{}
)

type nodeFirewallOptionsModel struct {
	ID                               types.String `tfsdk:"id"`
	NodeName                         types.String `tfsdk:"node_name"`
	Enable                           types.Bool   `tfsdk:"enabled"`
	LogLevelIn                       types.String `tfsdk:"log_level_in"`
	LogLevelOut                      types.String `tfsdk:"log_level_out"`
	LogLevelForward                  types.String `tfsdk:"log_level_forward"`
	NDP                              types.Bool   `tfsdk:"ndp"`
	NFConntrackMax                   types.Int64  `tfsdk:"nf_conntrack_max"`
	NFConntrackTCPTimeoutEstablished types.Int64  `tfsdk:"nf_conntrack_tcp_timeout_established"`
	NFTables                         types.Bool   `tfsdk:"nftables"`
	NoSMURFs                         types.Bool   `tfsdk:"nosmurfs"`
	SMURFLogLevel                    types.String `tfsdk:"smurf_log_level"`
	TCPFlagsLogLevel                 types.String `tfsdk:"tcp_flags_log_level"`
}

func (m *nodeFirewallOptionsModel) toOptionsRequestBody() *nodefirewall.OptionsPutRequestBody {
	body := &nodefirewall.OptionsPutRequestBody{}

	if !m.Enable.IsUnknown() {
		body.Enable = proxmoxtypes.CustomBoolPtr(m.Enable.ValueBoolPointer())
	}

	if !m.LogLevelIn.IsUnknown() {
		body.LogLevelIn = m.LogLevelIn.ValueStringPointer()
	}

	if !m.LogLevelOut.IsUnknown() {
		body.LogLevelOut = m.LogLevelOut.ValueStringPointer()
	}

	if !m.LogLevelForward.IsUnknown() {
		body.LogLevelForward = m.LogLevelForward.ValueStringPointer()
	}

	if !m.NDP.IsUnknown() {
		body.NDP = proxmoxtypes.CustomBoolPtr(m.NDP.ValueBoolPointer())
	}

	if !m.NFConntrackMax.IsUnknown() && !m.NFConntrackMax.IsNull() {
		body.NFConntrackMax = m.NFConntrackMax.ValueInt64Pointer()
	}

	if !m.NFConntrackTCPTimeoutEstablished.IsUnknown() && !m.NFConntrackTCPTimeoutEstablished.IsNull() {
		body.NFConntrackTCPTimeoutEstablished = m.NFConntrackTCPTimeoutEstablished.ValueInt64Pointer()
	}

	if !m.NFTables.IsUnknown() {
		body.NFTables = proxmoxtypes.CustomBoolPtr(m.NFTables.ValueBoolPointer())
	}

	if !m.NoSMURFs.IsUnknown() {
		body.NoSMURFs = proxmoxtypes.CustomBoolPtr(m.NoSMURFs.ValueBoolPointer())
	}

	if !m.SMURFLogLevel.IsUnknown() {
		body.SMURFLogLevel = m.SMURFLogLevel.ValueStringPointer()
	}

	if !m.TCPFlagsLogLevel.IsUnknown() {
		body.TCPFlagsLogLevel = m.TCPFlagsLogLevel.ValueStringPointer()
	}

	return body
}

func (m *nodeFirewallOptionsModel) importFromOptionsAPI(opts *nodefirewall.OptionsGetResponseData) {
	m.Enable = types.BoolPointerValue(opts.Enable.PointerBool())
	m.LogLevelIn = types.StringPointerValue(opts.LogLevelIn)
	m.LogLevelOut = types.StringPointerValue(opts.LogLevelOut)
	m.LogLevelForward = types.StringPointerValue(opts.LogLevelForward)
	m.NDP = types.BoolPointerValue(opts.NDP.PointerBool())
	m.NFConntrackMax = types.Int64PointerValue(opts.NFConntrackMax)
	m.NFConntrackTCPTimeoutEstablished = types.Int64PointerValue(opts.NFConntrackTCPTimeoutEstablished)
	m.NFTables = types.BoolPointerValue(opts.NFTables.PointerBool())
	m.NoSMURFs = types.BoolPointerValue(opts.NoSMURFs.PointerBool())
	m.SMURFLogLevel = types.StringPointerValue(opts.SMURFLogLevel)
	m.TCPFlagsLogLevel = types.StringPointerValue(opts.TCPFlagsLogLevel)
}

// NewNodeFirewallOptionsResource manages node firewall options resource.
func NewNodeFirewallOptionsResource() resource.Resource {
	return &nodeFirewallOptionsResource{}
}

type nodeFirewallOptionsResource struct {
	client proxmox.Client
}

func (r *nodeFirewallOptionsResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_node_firewall"
}

// Schema defines the schema for the resource.
func (r *nodeFirewallOptionsResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Manages Proxmox VE Node Firewall options.",
		Attributes: map[string]schema.Attribute{
			"id": attribute.ResourceID(),
			"node_name": schema.StringAttribute{
				Description: "The cluster node name.",
				Required:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Enable host firewall rules.",
				Optional:    true,
				Computed:    true,
			},
			"log_level_in": schema.StringAttribute{
				Description: "Log level for incoming traffic.",
				MarkdownDescription: "Log level for incoming traffic. Must be one of: " +
					"`emerg`, `alert`, `crit`, `err`, `warning`, `notice`, `info`, `debug`, `nolog`.",
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf("emerg", "alert", "crit", "err", "warning", "notice", "info", "debug", "nolog"),
				},
			},
			"log_level_out": schema.StringAttribute{
				Description: "Log level for outgoing traffic.",
				MarkdownDescription: "Log level for outgoing traffic. Must be one of: " +
					"`emerg`, `alert`, `crit`, `err`, `warning`, `notice`, `info`, `debug`, `nolog`.",
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf("emerg", "alert", "crit", "err", "warning", "notice", "info", "debug", "nolog"),
				},
			},
			"log_level_forward": schema.StringAttribute{
				Description: "Log level for forwarded traffic.",
				MarkdownDescription: "Log level for forwarded traffic. Must be one of: " +
					"`emerg`, `alert`, `crit`, `err`, `warning`, `notice`, `info`, `debug`, `nolog`.",
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf("emerg", "alert", "crit", "err", "warning", "notice", "info", "debug", "nolog"),
				},
			},
			"ndp": schema.BoolAttribute{
				Description: "Enable NDP (Neighbor Discovery Protocol).",
				Optional:    true,
				Computed:    true,
			},
			"nf_conntrack_max": schema.Int64Attribute{
				Description: "Maximum number of tracked connections.",
				Optional:    true,
				Computed:    true,
				Validators: []validator.Int64{
					int64validator.Between(32768, 999999999),
				},
			},
			"nf_conntrack_tcp_timeout_established": schema.Int64Attribute{
				Description: "Conntrack established timeout.",
				Optional:    true,
				Computed:    true,
				Validators: []validator.Int64{
					int64validator.Between(7875, 999999999),
				},
			},
			"nftables": schema.BoolAttribute{
				Description: "Enable nftables based firewall (tech preview).",
				Optional:    true,
				Computed:    true,
			},
			"nosmurfs": schema.BoolAttribute{
				Description: "Enable SMURFS filter.",
				Optional:    true,
				Computed:    true,
			},
			"smurf_log_level": schema.StringAttribute{
				Description: "Log level for SMURFS filter.",
				MarkdownDescription: "Log level for SMURFS filter. Must be one of: " +
					"`emerg`, `alert`, `crit`, `err`, `warning`, `notice`, `info`, `debug`, `nolog`.",
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf("emerg", "alert", "crit", "err", "warning", "notice", "info", "debug", "nolog"),
				},
			},
			"tcp_flags_log_level": schema.StringAttribute{
				Description: "Log level for illegal tcp flags filter.",
				MarkdownDescription: "Log level for illegal tcp flags filter. Must be one of: " +
					"`emerg`, `alert`, `crit`, `err`, `warning`, `notice`, `info`, `debug`, `nolog`.",
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf("emerg", "alert", "crit", "err", "warning", "notice", "info", "debug", "nolog"),
				},
			},
		},
	}
}

// Configure configures the resource.
func (r *nodeFirewallOptionsResource) Configure(
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

	r.client = cfg.Client
}

// Create creates node firewall options.
func (r *nodeFirewallOptionsResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan nodeFirewallOptionsModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	body := plan.toOptionsRequestBody()

	nodeName := plan.NodeName.ValueString()

	err := r.client.Node(nodeName).Firewall().SetNodeOptions(ctx, body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating node firewall options",
			"Could not create node firewall options, unexpected error: "+err.Error(),
		)

		return
	}

	plan.ID = types.StringValue(nodeName)

	r.read(ctx, &plan, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *nodeFirewallOptionsResource) read(
	ctx context.Context,
	model *nodeFirewallOptionsModel,
	diags *diag.Diagnostics,
) {
	nodeName := model.NodeName.ValueString()

	options, err := r.client.Node(nodeName).Firewall().GetNodeOptions(ctx)
	if err != nil {
		diags.AddError(
			"Error getting node firewall options",
			"Could not get node firewall options, unexpected error: "+err.Error(),
		)

		return
	}

	model.importFromOptionsAPI(options)
}

// Read reads node firewall options.
func (r *nodeFirewallOptionsResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var state nodeFirewallOptionsModel

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

// Update updates node firewall options.
func (r *nodeFirewallOptionsResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan nodeFirewallOptionsModel
	var state nodeFirewallOptionsModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	body := plan.toOptionsRequestBody()

	nodeName := plan.NodeName.ValueString()

	var toDelete []string
	attribute.CheckDeleteComputed(plan.NFConntrackMax, state.NFConntrackMax, &toDelete, "nf_conntrack_max")
	attribute.CheckDeleteComputed(
		plan.NFConntrackTCPTimeoutEstablished, state.NFConntrackTCPTimeoutEstablished, &toDelete, "nf_conntrack_tcp_timeout_established",
	)

	body.Delete = &toDelete

	err := r.client.Node(nodeName).Firewall().SetNodeOptions(ctx, body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating node firewall options",
			"Could not update node firewall options, unexpected error: "+err.Error(),
		)

		return
	}

	r.read(ctx, &plan, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// Delete deletes node firewall options.
func (r *nodeFirewallOptionsResource) Delete(
	_ context.Context,
	_ resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	resp.Diagnostics.AddWarning(
		"Node firewall options cannot be deleted",
		"Node firewall options are a configuration of the node and cannot be deleted. "+
			"The resource will be removed from the Terraform state only.",
	)
}

// ImportState imports node firewall options.
func (r *nodeFirewallOptionsResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	state := nodeFirewallOptionsModel{
		ID:       types.StringValue(req.ID),
		NodeName: types.StringValue(req.ID),
	}
	r.read(ctx, &state, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	diags := resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}
