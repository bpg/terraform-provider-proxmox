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
	"strconv"
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

	"github.com/bpg/terraform-provider-proxmox/fwprovider/attribute"
	customtypes "github.com/bpg/terraform-provider-proxmox/fwprovider/types"

	"github.com/bpg/terraform-provider-proxmox/proxmox"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes"
	proxmoxtypes "github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

var (
	_ resource.Resource                = &linuxBridgeResource{}
	_ resource.ResourceWithConfigure   = &linuxBridgeResource{}
	_ resource.ResourceWithImportState = &linuxBridgeResource{}
)

type linuxBridgeResourceModel struct {
	// Base attributes
	ID        types.String            `tfsdk:"id"`
	NodeName  types.String            `tfsdk:"node_name"`
	Name      types.String            `tfsdk:"name"`
	Address   customtypes.IPCIDRValue `tfsdk:"address"`
	Gateway   customtypes.IPAddrValue `tfsdk:"gateway"`
	Address6  customtypes.IPCIDRValue `tfsdk:"address6"`
	Gateway6  customtypes.IPAddrValue `tfsdk:"gateway6"`
	Autostart types.Bool              `tfsdk:"autostart"`
	MTU       types.Int64             `tfsdk:"mtu"`
	Comment   types.String            `tfsdk:"comment"`
	// Linux bridge attributes
	Ports     []types.String `tfsdk:"ports"`
	VLANAware types.Bool     `tfsdk:"vlan_aware"`
}

//nolint:lll
func (m *linuxBridgeResourceModel) exportToNetworkInterfaceCreateUpdateBody() *nodes.NetworkInterfaceCreateUpdateRequestBody {
	body := &nodes.NetworkInterfaceCreateUpdateRequestBody{
		Iface:     m.Name.ValueString(),
		Type:      "bridge",
		Autostart: proxmoxtypes.CustomBool(m.Autostart.ValueBool()).Pointer(),
	}

	body.CIDR = m.Address.ValueStringPointer()
	body.Gateway = m.Gateway.ValueStringPointer()
	body.CIDR6 = m.Address6.ValueStringPointer()
	body.Gateway6 = m.Gateway6.ValueStringPointer()

	if !m.MTU.IsUnknown() {
		body.MTU = m.MTU.ValueInt64Pointer()
	}

	body.Comments = m.Comment.ValueStringPointer()

	var sanitizedPorts []string

	for _, port := range m.Ports {
		port := strings.TrimSpace(port.ValueString())
		if len(port) > 0 {
			sanitizedPorts = append(sanitizedPorts, port)
		}
	}

	sort.Strings(sanitizedPorts)
	bridgePorts := strings.Join(sanitizedPorts, " ")

	if len(bridgePorts) > 0 {
		body.BridgePorts = &bridgePorts
	}

	if m.VLANAware.ValueBool() {
		body.BridgeVLANAware = proxmoxtypes.CustomBool(true).Pointer()
	}

	return body
}

func (m *linuxBridgeResourceModel) importFromNetworkInterfaceList(
	ctx context.Context,
	iface *nodes.NetworkInterfaceListResponseData,
) error {
	m.Address = customtypes.NewIPCIDRPointerValue(iface.CIDR)
	m.Gateway = customtypes.NewIPAddrPointerValue(iface.Gateway)
	m.Address6 = customtypes.NewIPCIDRPointerValue(iface.CIDR6)
	m.Gateway6 = customtypes.NewIPAddrPointerValue(iface.Gateway6)

	m.Autostart = types.BoolPointerValue(iface.Autostart.PointerBool())
	if m.Autostart.IsNull() {
		m.Autostart = types.BoolValue(false)
	}

	if iface.MTU != nil {
		if v, err := strconv.Atoi(*iface.MTU); err == nil {
			m.MTU = types.Int64Value(int64(v))
		}
	} else {
		m.MTU = types.Int64Null()
	}

	// Comments can be set to an empty string in plant, which will translate to a "no value" in PVE
	// So we don't want to set it to null if it's empty, as this will be indicated as a plan drift
	if iface.Comments != nil {
		m.Comment = types.StringValue(strings.TrimSpace(*iface.Comments))
	}

	if iface.BridgeVLANAware != nil {
		m.VLANAware = types.BoolPointerValue(iface.BridgeVLANAware.PointerBool())
	} else {
		m.VLANAware = types.BoolValue(false)
	}

	if iface.BridgePorts != nil && len(*iface.BridgePorts) > 0 {
		ports, diags := types.ListValueFrom(ctx, types.StringType, strings.Split(*iface.BridgePorts, " "))
		if diags.HasError() {
			return fmt.Errorf("failed to parse bridge ports: %s", *iface.BridgePorts)
		}

		diags = ports.ElementsAs(ctx, &m.Ports, false)
		if diags.HasError() {
			return fmt.Errorf("failed to build bridge ports list: %s", *iface.BridgePorts)
		}
	}

	return nil
}

// NewLinuxBridgeResource creates a new resource for managing Linux Bridge network interfaces.
func NewLinuxBridgeResource() resource.Resource {
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
		Description: "Manages a Linux Bridge network interface in a Proxmox VE node.",
		Attributes: map[string]schema.Attribute{
			// Base attributes
			"id": attribute.ID("A unique identifier with format `<node name>:<iface>`"),
			"node_name": schema.StringAttribute{
				Description: "The name of the node.",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description:         "The interface name.",
				MarkdownDescription: "The interface name. Must be `vmbrN`, where N is a number between 0 and 9999.",
				Required:            true,
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
				CustomType:  customtypes.IPCIDRType{},
				Optional:    true,
			},
			"gateway": schema.StringAttribute{
				Description: "Default gateway address.",
				CustomType:  customtypes.IPAddrType{},
				Optional:    true,
			},
			"address6": schema.StringAttribute{
				Description: "The interface IPv6/CIDR address.",
				CustomType:  customtypes.IPCIDRType{},
				Optional:    true,
			},
			"gateway6": schema.StringAttribute{
				Description: "Default IPv6 gateway address.",
				CustomType:  customtypes.IPAddrType{},
				Optional:    true,
			},
			"autostart": schema.BoolAttribute{
				Description: "Automatically start interface on boot (defaults to `true`).",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"mtu": schema.Int64Attribute{
				Description: "The interface MTU.",
				Optional:    true,
			},
			"comment": schema.StringAttribute{
				Description: "Comment for the interface.",
				Optional:    true,
			},
			// Linux Bridge attributes
			"ports": schema.ListAttribute{
				Description: "The interface bridge ports.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"vlan_aware": schema.BoolAttribute{
				Description: "Whether the interface bridge is VLAN aware (defaults to `false`).",
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
			fmt.Sprintf("Expected *proxmox.Client, got: %T", req.ProviderData),
		)

		return
	}

	r.client = client
}

//nolint:dupl
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

	plan.ID = types.StringValue(plan.NodeName.ValueString() + ":" + plan.Name.ValueString())

	r.read(ctx, &plan, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)

	err = r.client.Node(plan.NodeName.ValueString()).ReloadNetworkConfiguration(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reloading network configuration",
			fmt.Sprintf("Could not reload network configuration on node '%s', unexpected error: %s",
				plan.NodeName.ValueString(), err.Error()),
		)
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
		if iface.Iface != model.Name.ValueString() {
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

	if !plan.MTU.Equal(state.MTU) && plan.MTU.ValueInt64() == 0 {
		toDelete = append(toDelete, "mtu")
		body.MTU = nil
	}

	if !plan.Gateway.Equal(state.Gateway) && plan.Gateway.ValueString() == "" {
		toDelete = append(toDelete, "gateway")
		body.Gateway = nil
	}

	if !plan.Gateway6.Equal(state.Gateway6) && plan.Gateway6.ValueString() == "" {
		toDelete = append(toDelete, "gateway6")
		body.Gateway6 = nil
	}

	// VLANAware is computed, will never be null
	if !plan.VLANAware.Equal(state.VLANAware) && !plan.VLANAware.ValueBool() {
		toDelete = append(toDelete, "bridge_vlan_aware")
		body.BridgeVLANAware = nil
	}

	if len(toDelete) > 0 {
		body.Delete = &toDelete
	}

	err := r.client.Node(plan.NodeName.ValueString()).UpdateNetworkInterface(ctx, plan.Name.ValueString(), body)
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

	err = r.client.Node(state.NodeName.ValueString()).ReloadNetworkConfiguration(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reloading network configuration",
			fmt.Sprintf("Could not reload network configuration on node '%s', unexpected error: %s",
				state.NodeName.ValueString(), err.Error()),
		)
	}
}

// Delete deletes a Linux Bridge interface.
//
//nolint:dupl
func (r *linuxBridgeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state linuxBridgeResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Node(state.NodeName.ValueString()).DeleteNetworkInterface(ctx, state.Name.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "interface does not exist") {
			resp.Diagnostics.AddWarning(
				"Linux Bridge interface does not exist",
				fmt.Sprintf("Could not delete Linux Bridge '%s', interface does not exist, "+
					"or has already been deleted outside of Terraform.", state.Name.ValueString()),
			)
		} else {
			resp.Diagnostics.AddError(
				"Error deleting Linux Bridge interface",
				fmt.Sprintf("Could not delete Linux Bridge '%s', unexpected error: %s",
					state.Name.ValueString(), err.Error()),
			)
		}

		return
	}

	err = r.client.Node(state.NodeName.ValueString()).ReloadNetworkConfiguration(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reloading network configuration",
			fmt.Sprintf("Could not reload network configuration on node '%s', unexpected error: %s",
				state.NodeName.ValueString(), err.Error()),
		)
	}
}

func (r *linuxBridgeResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	idParts := strings.Split(req.ID, ":")
	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: `node_name:iface`. Got: %q", req.ID),
		)

		return
	}

	nodeName := idParts[0]
	iface := idParts[1]

	state := linuxBridgeResourceModel{
		ID:       types.StringValue(req.ID),
		NodeName: types.StringValue(nodeName),
		Name:     types.StringValue(iface),
	}
	r.read(ctx, &state, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	diags := resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}
