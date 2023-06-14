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
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	pvetypes "github.com/bpg/terraform-provider-proxmox/internal/types"
	"github.com/bpg/terraform-provider-proxmox/proxmox"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes"
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

type linuxBridgeResourceModel struct {
	// Base attributes
	ID        types.String           `tfsdk:"id"`
	NodeName  types.String           `tfsdk:"node_name"`
	Iface     types.String           `tfsdk:"iface"`
	Address   pvetypes.IPv4CIDRValue `tfsdk:"address"`
	Gateway   pvetypes.IPv4CIDRValue `tfsdk:"gateway"`
	Autostart types.Bool             `tfsdk:"autostart"`
	Comment   types.String           `tfsdk:"comment"`
	// Linux bridge attributes
	BridgePorts     []types.String `tfsdk:"bridge_ports"`
	BridgeVLANAware types.Bool     `tfsdk:"bridge_vlan_aware"`
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
			},
			"address": schema.StringAttribute{
				Description: "The interface IPv4/CIDR address.",
				CustomType:  pvetypes.IPv4CIDRType{},
				Optional:    true,
			},
			"gateway": schema.StringAttribute{
				Description: "Default gateway address.",
				CustomType:  pvetypes.IPv4CIDRType{},
				Optional:    true,
			},
			"autostart": schema.BoolAttribute{
				Description: "Automatically start interface on boot.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
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

	body := &nodes.NetworkInterfaceCreateUpdateRequestBody{
		Iface:     plan.Iface.ValueString(),
		Type:      "bridge",
		Autostart: pvetypes.CustomBool(plan.Autostart.ValueBool()).Pointer(),
	}

	if !plan.Address.IsUnknown() {
		body.CIDR = plan.Address.ValueStringPointer()
	}

	if !plan.Gateway.IsUnknown() {
		body.Gateway = plan.Gateway.ValueStringPointer()
	}

	if !plan.Comment.IsUnknown() {
		body.Comments = plan.Comment.ValueStringPointer()
	}

	var sanitizedPorts []string

	for i := 0; i < len(plan.BridgePorts); i++ {
		port := strings.TrimSpace(plan.BridgePorts[i].ValueString())
		if len(port) > 0 {
			sanitizedPorts = append(sanitizedPorts, port)
		}
	}
	sort.Strings(sanitizedPorts)
	bridgePorts := strings.Join(sanitizedPorts, " ")

	if len(bridgePorts) > 0 {
		body.BridgePorts = &bridgePorts
	}

	body.BridgeVLANAware = pvetypes.CustomBool(plan.BridgeVLANAware.ValueBool()).Pointer()

	err := r.client.Node(plan.NodeName.ValueString()).CreateNetworkInterface(ctx, body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Linux Bridge interface",
			"Could not create Linux Bridge, unexpected error: "+err.Error(),
		)

		return
	}

	plan.ID = types.StringValue(plan.NodeName.ValueString() + ":" + plan.Iface.ValueString())

	err = r.read(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Linux Bridge interface",
			"Could not read Linux Bridge, unexpected error: "+err.Error(),
		)

		return
	}

	resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *linuxBridgeResource) read(ctx context.Context, model *linuxBridgeResourceModel) error {
	ifaces, err := r.client.Node(model.NodeName.ValueString()).ListNetworkInterfaces(ctx)
	if err != nil {
		return fmt.Errorf("error listing network interfaces: %w", err)
	}

	for _, iface := range ifaces {
		if iface.Iface != model.Iface.ValueString() {
			continue
		}

		if iface.CIDR != nil {
			model.Address = pvetypes.NewIPv4CIDRPointerValue(iface.CIDR)
		}

		if iface.Gateway != nil {
			model.Gateway = pvetypes.NewIPv4CIDRPointerValue(iface.Gateway)
		}

		if iface.Autostart != nil {
			model.Autostart = types.BoolPointerValue(iface.Autostart.PointerBool())
		}

		if iface.Comments != nil {
			model.Comment = types.StringValue(strings.TrimSpace(*iface.Comments))
		}

		if iface.BridgeVLANAware != nil {
			model.BridgeVLANAware = types.BoolPointerValue(iface.BridgeVLANAware.PointerBool())
		}

		// model.BridgePorts = types.NewStringListValue(strings.Split(iface.BridgePorts, " "))
		break
	}

	return nil
}

// Read reads a Linux Bridge interface.
func (r *linuxBridgeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}

// Update updates a Linux Bridge interface.
func (r *linuxBridgeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
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
		resp.Diagnostics.AddError(
			"Error deleting Linux Bridge interface",
			"Could not delete Linux Bridge, unexpected error: "+err.Error(),
		)

		return
	}
}

func (r *linuxBridgeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
