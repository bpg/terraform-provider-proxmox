/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package network

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/config"
	"github.com/bpg/terraform-provider-proxmox/proxmox"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes"
)

var (
	_ resource.Resource              = &applierResource{}
	_ resource.ResourceWithConfigure = &applierResource{}
)

type applierResourceModel struct {
	// ID is opaque, set to the Unix timestamp (milliseconds) of the last apply.
	ID        types.String `tfsdk:"id"`
	NodeName  types.String `tfsdk:"node_name"`
	OnCreate  types.Bool   `tfsdk:"on_create"`
	OnDestroy types.Bool   `tfsdk:"on_destroy"`
	Timeout   types.Int64  `tfsdk:"timeout_reload"`
	Triggers  types.Map    `tfsdk:"triggers"`
}

// NewApplierResource creates a new resource for applying staged node network configuration.
func NewApplierResource() resource.Resource {
	return &applierResource{}
}

type applierResource struct {
	client proxmox.Client
}

func (r *applierResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "proxmox_network_applier"
}

func (r *applierResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Applies staged (pending) node network configuration.",
		MarkdownDescription: "**EXPERIMENTAL** Reloads the network configuration of a single node, activating " +
			"changes staged by interface resources that set `reload = false` (equivalent to " +
			"`PUT /nodes/{node}/network`).\n\n" +
			"The reload happens on create and on destroy. Use `triggers` to force a new apply when " +
			"something it depends on changes.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Opaque identifier set to the Unix timestamp (milliseconds) when the apply was executed.",
			},
			"node_name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the node.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"on_create": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
				Description: "Whether to reload on resource creation (defaults to `true`).",
			},
			"on_destroy": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
				Description: "Whether to reload on resource destruction (defaults to `true`).",
			},
			"timeout_reload": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(int64(nodes.NetworkReloadTimeout.Seconds())),
				Description: "Timeout for network reload operations in seconds (defaults to `100`).",
				Validators: []validator.Int64{
					int64validator.AtLeast(5),
				},
			},
			"triggers": schema.MapAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "Arbitrary values that force a new apply when they change.",
				PlanModifiers: []planmodifier.Map{
					mapplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *applierResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	cfg, ok := req.ProviderData.(config.Resource)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected config.Resource, got: %T", req.ProviderData),
		)

		return
	}

	r.client = cfg.Client
}

// reload activates the staged network configuration on the node.
func (r *applierResource) reload(ctx context.Context, m *applierResourceModel) error {
	reloadCtx, cancel := context.WithTimeout(ctx, time.Duration(m.Timeout.ValueInt64())*time.Second)
	defer cancel()

	return r.client.Node(m.NodeName.ValueString()).ReloadNetworkConfiguration(reloadCtx)
}

func (r *applierResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan applierResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if plan.OnCreate.ValueBool() {
		if err := r.reload(ctx, &plan); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Unable to Apply Network Configuration on Node %q", plan.NodeName.ValueString()),
				err.Error(),
			)

			return
		}
	}

	plan.ID = types.StringValue(strconv.FormatInt(time.Now().UTC().UnixMilli(), 10))

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read is a no-op: the applier owns no server-side object, only the record of an apply.
func (r *applierResource) Read(_ context.Context, _ resource.ReadRequest, _ *resource.ReadResponse) {
}

func (r *applierResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan applierResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.reload(ctx, &plan); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to Apply Network Configuration on Node %q", plan.NodeName.ValueString()),
			err.Error(),
		)

		return
	}

	plan.ID = types.StringValue(strconv.FormatInt(time.Now().UTC().UnixMilli(), 10))

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *applierResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state applierResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if state.OnDestroy.ValueBool() {
		if err := r.reload(ctx, &state); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Unable to Apply Network Configuration on Node %q", state.NodeName.ValueString()),
				err.Error(),
			)

			return
		}
	}
}
