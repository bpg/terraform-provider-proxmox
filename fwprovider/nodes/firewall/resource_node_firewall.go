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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
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

	body.Enable = proxmoxtypes.CustomBoolPtr(m.Enable.ValueBoolPointer())
	body.LogLevelIn = m.LogLevelIn.ValueStringPointer()
	body.LogLevelOut = m.LogLevelOut.ValueStringPointer()
	body.LogLevelForward = m.LogLevelForward.ValueStringPointer()
	body.NDP = proxmoxtypes.CustomBoolPtr(m.NDP.ValueBoolPointer())
	body.NFConntrackMax = m.NFConntrackMax.ValueInt64Pointer()
	body.NFConntrackTCPTimeoutEstablished = m.NFConntrackTCPTimeoutEstablished.ValueInt64Pointer()
	body.NFTables = proxmoxtypes.CustomBoolPtr(m.NFTables.ValueBoolPointer())
	body.NoSMURFs = proxmoxtypes.CustomBoolPtr(m.NoSMURFs.ValueBoolPointer())
	body.SMURFLogLevel = m.SMURFLogLevel.ValueStringPointer()
	body.TCPFlagsLogLevel = m.TCPFlagsLogLevel.ValueStringPointer()

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
		MarkdownDescription: "Manages Proxmox VE Node Firewall options.\n\n" +
			"~> This resource in fact updates existing node firewall configuration created by PVE on bootstrap. " +
			"All optional attributes have explicit defaults for deterministic behavior (PVE may change defaults in the future). " +
			"See [API documentation](https://pve.proxmox.com/pve-docs/api-viewer/index.html#/nodes/{node}/firewall/options).",
		Attributes: map[string]schema.Attribute{
			"id": attribute.ResourceID(),
			"node_name": schema.StringAttribute{
				Description: "The cluster node name.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"enabled": schema.BoolAttribute{
				Description: "Enable host firewall rules (defaults to `true`).",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"log_level_in": schema.StringAttribute{
				Description: "Log level for incoming traffic.",
				MarkdownDescription: "Log level for incoming traffic. Must be one of: " +
					"`emerg`, `alert`, `crit`, `err`, `warning`, `notice`, `info`, `debug`, `nolog` " +
					"(defaults to `nolog`).",
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf("emerg", "alert", "crit", "err", "warning", "notice", "info", "debug", "nolog"),
				},
				Default: stringdefault.StaticString("nolog"),
			},
			"log_level_out": schema.StringAttribute{
				Description: "Log level for outgoing traffic.",
				MarkdownDescription: "Log level for outgoing traffic. Must be one of: " +
					"`emerg`, `alert`, `crit`, `err`, `warning`, `notice`, `info`, `debug`, `nolog` " +
					"(defaults to `nolog`).",
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf("emerg", "alert", "crit", "err", "warning", "notice", "info", "debug", "nolog"),
				},
				Default: stringdefault.StaticString("nolog"),
			},
			"log_level_forward": schema.StringAttribute{
				Description: "Log level for forwarded traffic.",
				MarkdownDescription: "Log level for forwarded traffic. Must be one of: " +
					"`emerg`, `alert`, `crit`, `err`, `warning`, `notice`, `info`, `debug`, `nolog` " +
					"(defaults to `nolog`).",
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf("emerg", "alert", "crit", "err", "warning", "notice", "info", "debug", "nolog"),
				},
				Default: stringdefault.StaticString("nolog"),
			},
			"ndp": schema.BoolAttribute{
				Description: "Enable NDP - Neighbor Discovery Protocol (defaults to `true`).",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"nf_conntrack_max": schema.Int64Attribute{
				Description: "Maximum number of tracked connections (defaults to `262144`). Minimum value " +
					"is `32768`.",
				Optional: true,
				Computed: true,
				Validators: []validator.Int64{
					int64validator.AtLeast(32768),
				},
				Default: int64default.StaticInt64(262144),
			},
			"nf_conntrack_tcp_timeout_established": schema.Int64Attribute{
				Description: "Conntrack established timeout in seconds (defaults to `432000` - 5 days). " +
					"Minimum value is `7875`.",
				Optional: true,
				Computed: true,
				Validators: []validator.Int64{
					int64validator.AtLeast(7875),
				},
				Default: int64default.StaticInt64(432000),
			},
			"nftables": schema.BoolAttribute{
				Description: "Enable nftables based firewall (tech preview, defaults to `false`).",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"nosmurfs": schema.BoolAttribute{
				Description: "Enable SMURFS filter (defaults to `true`).",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"smurf_log_level": schema.StringAttribute{
				Description: "Log level for SMURFS filter.",
				MarkdownDescription: "Log level for SMURFS filter. Must be one of: " +
					"`emerg`, `alert`, `crit`, `err`, `warning`, `notice`, `info`, `debug`, `nolog` " +
					"(defaults to `nolog`).",
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf("emerg", "alert", "crit", "err", "warning", "notice", "info", "debug", "nolog"),
				},
				Default: stringdefault.StaticString("nolog"),
			},
			"tcp_flags_log_level": schema.StringAttribute{
				Description: "Log level for illegal tcp flags filter.",
				MarkdownDescription: "Log level for illegal tcp flags filter. Must be one of: " +
					"`emerg`, `alert`, `crit`, `err`, `warning`, `notice`, `info`, `debug`, `nolog` " +
					"(defaults to `nolog`).",
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf("emerg", "alert", "crit", "err", "warning", "notice", "info", "debug", "nolog"),
				},
				Default: stringdefault.StaticString("nolog"),
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
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var state nodeFirewallOptionsModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	nodeName := state.NodeName.ValueString()

	toDelete := []string{
		"enable", "log_level_in", "log_level_out", "log_level_forward",
		"ndp", "nf_conntrack_max", "nf_conntrack_tcp_timeout_established", "nftables",
		"nosmurfs", "smurf_log_level", "tcp_flags_log_level",
	}

	body := &nodefirewall.OptionsPutRequestBody{
		Delete: &toDelete,
	}

	err := r.client.Node(nodeName).Firewall().SetNodeOptions(ctx, body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting node firewall options",
			"Could not delete node firewall options, unexpected error: "+err.Error(),
		)
	}
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
