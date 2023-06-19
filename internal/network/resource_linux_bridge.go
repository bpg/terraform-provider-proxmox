/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package network

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	pvetypes "github.com/bpg/terraform-provider-proxmox/internal/types"
	"github.com/bpg/terraform-provider-proxmox/proxmox"
)

var (
	_ resource.Resource                = &linuxBridgeResource{}
	_ resource.ResourceWithConfigure   = &linuxBridgeResource{}
	_ resource.ResourceWithImportState = &linuxBridgeResource{}
)

// NewInterfaceLinuxBridgeResource creates a new resource for managing Linux Bridge network interfaces.
func NewInterfaceLinuxBridgeResource() resource.Resource {
	return &linuxBridgeResource{}
}

type linuxBridgeResource struct {
	client proxmox.Client
}

func (r *linuxBridgeResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_network_linux_bridge"
}

// Schema defines the schema for the resource.
func (r *linuxBridgeResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Manages an Linux Bridge network interface in a Proxmox VE node.",
		Attributes: map[string]schema.Attribute{
			// Base attributes
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "A unique identifier with format '<node name>:<iface>'",
			},
			"node_name": schema.StringAttribute{
				Description: "The name of the node.",
				Required:    true,
			},
			"iface": schema.StringAttribute{
				Description: "The interface name.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^vmbr(\d{1,4})$`),
						`must be "vmbrN", where N is a number between 0 and 9999`,
					),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"address": schema.StringAttribute{
				Description: "The interface IPv4/CIDR address.",
				CustomType:  pvetypes.IPv4CIDRType{},
				Optional:    true,
			},
			"gateway": schema.StringAttribute{
				Description: "Default gateway address.",
				CustomType:  pvetypes.IPv4Type{},
				Optional:    true,
			},
			"address6": schema.StringAttribute{
				Description: "The interface IPv6/CIDR address.",
				CustomType:  pvetypes.IPv6CIDRType{},
				Optional:    true,
			},
			"gateway6": schema.StringAttribute{
				Description: "Default IPv6 gateway address.",
				CustomType:  pvetypes.IPv6Type{},
				Optional:    true,
			},
			"autostart": schema.BoolAttribute{
				Description: "Automatically start interface on boot.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"mtu": schema.Int64Attribute{
				Description: "The interface MTU.",
				Optional:    true,
				Computed:    true,
			},
			"comment": schema.StringAttribute{
				Description: "Comment for the interface.",
				Optional:    true,
			},
			// Linux Bridge attributes
			"bridge_ports": schema.ListAttribute{
				Description: "The interface bridge ports.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"bridge_vlan_aware": schema.BoolAttribute{
				Description: "Whether the interface bridge is VLAN aware.",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func (r *linuxBridgeResource) Configure(
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

func (r *linuxBridgeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan linuxBridgeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	body := plan.exportToNetworkInterfaceCreateUpdateBody()

	err := r.client.Node(plan.NodeName.ValueString()).CreateNetworkInterface(ctx, body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Linux Bridge interface",
			"Could not create Linux Bridge, unexpected error: "+err.Error(),
		)

		return
	}

	plan.ID = types.StringValue(plan.NodeName.ValueString() + ":" + plan.Iface.ValueString())

	r.read(ctx, &plan, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *linuxBridgeResource) read(ctx context.Context, model *linuxBridgeResourceModel, diags *diag.Diagnostics) {
	ifaces, err := r.client.Node(model.NodeName.ValueString()).ListNetworkInterfaces(ctx)
	if err != nil {
		diags.AddError(
			"Error listing network interfaces",
			"Could not list network interfaces, unexpected error: "+err.Error(),
		)

		return
	}

	for _, iface := range ifaces {
		if iface.Iface != model.Iface.ValueString() {
			continue
		}

		err = model.importFromNetworkInterfaceList(ctx, iface)
		if err != nil {
			diags.AddError(
				"Error converting network interface to a model",
				"Could not import network interface from API response, unexpected error: "+err.Error(),
			)

			return
		}

		break
	}
}

// Read reads a Linux Bridge interface.
func (r *linuxBridgeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state linuxBridgeResourceModel
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

// Update updates a Linux Bridge interface.
func (r *linuxBridgeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state linuxBridgeResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	body := plan.exportToNetworkInterfaceCreateUpdateBody()

	var toDelete []string

	if !plan.MTU.Equal(state.MTU) && (plan.MTU.IsUnknown() || plan.MTU.ValueInt64() == 0) {
		toDelete = append(toDelete, "mtu")
		body.MTU = nil
	}

	// BridgeVLANAware is computed, will never be null
	if !plan.BridgeVLANAware.Equal(state.BridgeVLANAware) && !plan.BridgeVLANAware.ValueBool() {
		toDelete = append(toDelete, "bridge_vlan_aware")
		body.BridgeVLANAware = nil
	}

	if len(toDelete) > 0 {
		body.Delete = &toDelete
	}

	err := r.client.Node(plan.NodeName.ValueString()).UpdateNetworkInterface(ctx, plan.Iface.ValueString(), body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Linux Bridge interface",
			"Could not update Linux Bridge, unexpected error: "+err.Error(),
		)

		return
	}

	r.read(ctx, &plan, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// Delete deletes a Linux Bridge interface.
func (r *linuxBridgeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state linuxBridgeResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Node(state.NodeName.ValueString()).DeleteNetworkInterface(ctx, state.Iface.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "interface does not exist") {
			resp.Diagnostics.AddWarning(
				"Linux Bridge interface does not exist",
				fmt.Sprintf("Could not delete Linux Bridge '%s', interface does not exist, "+
					"or has already been deleted outside of Terraform.", state.Iface.ValueString()),
			)
		} else {
			resp.Diagnostics.AddError(
				"Error deleting Linux Bridge interface",
				fmt.Sprintf("Could not delete Linux Bridge '%s', unexpected error: %s",
					state.Iface.ValueString(), err.Error()),
			)
		}

		return
	}
}

func (r *linuxBridgeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ":")
	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: node_name:iface. Got: %q", req.ID),
		)

		return
	}

	nodeName := idParts[0]
	iface := idParts[1]

	state := linuxBridgeResourceModel{
		ID:       types.StringValue(req.ID),
		NodeName: types.StringValue(nodeName),
		Iface:    types.StringValue(iface),
	}
	r.read(ctx, &state, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	diags := resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}
