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
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes"
)

var (
	_ resource.Resource                = &linuxVLANResource{}
	_ resource.ResourceWithConfigure   = &linuxVLANResource{}
	_ resource.ResourceWithImportState = &linuxVLANResource{}
)

type linuxVLANResourceModel struct {
	// Base attributes
	ID        types.String         `tfsdk:"id"`
	NodeName  types.String         `tfsdk:"node_name"`
	Name      types.String         `tfsdk:"name"`
	Address   pvetypes.IPCIDRValue `tfsdk:"address"`
	Gateway   pvetypes.IPAddrValue `tfsdk:"gateway"`
	Address6  pvetypes.IPCIDRValue `tfsdk:"address6"`
	Gateway6  pvetypes.IPAddrValue `tfsdk:"gateway6"`
	Autostart types.Bool           `tfsdk:"autostart"`
	MTU       types.Int64          `tfsdk:"mtu"`
	Comment   types.String         `tfsdk:"comment"`
	// Linux VLAN attributes
	Interface types.String `tfsdk:"interface"`
	VLAN      types.Int64  `tfsdk:"vlan"`
}

//nolint:lll
func (m *linuxVLANResourceModel) exportToNetworkInterfaceCreateUpdateBody() *nodes.NetworkInterfaceCreateUpdateRequestBody {
	body := &nodes.NetworkInterfaceCreateUpdateRequestBody{
		Iface:     m.Name.ValueString(),
		Type:      "vlan",
		Autostart: pvetypes.CustomBool(m.Autostart.ValueBool()).Pointer(),
	}

	body.CIDR = m.Address.ValueStringPointer()
	body.Gateway = m.Gateway.ValueStringPointer()
	body.CIDR6 = m.Address6.ValueStringPointer()
	body.Gateway6 = m.Gateway6.ValueStringPointer()
	body.Comments = m.Comment.ValueStringPointer()

	if !m.MTU.IsUnknown() {
		body.MTU = m.MTU.ValueInt64Pointer()
	}

	if !m.Interface.IsUnknown() {
		body.VLANRawDevice = m.Interface.ValueStringPointer()
	}

	if !m.VLAN.IsUnknown() {
		body.VLANID = m.VLAN.ValueInt64Pointer()
	}

	return body
}

func (m *linuxVLANResourceModel) importFromNetworkInterfaceList(iface *nodes.NetworkInterfaceListResponseData) {
	m.Address = pvetypes.NewIPCIDRPointerValue(iface.CIDR)
	m.Gateway = pvetypes.NewIPAddrPointerValue(iface.Gateway)
	m.Address6 = pvetypes.NewIPCIDRPointerValue(iface.CIDR6)
	m.Gateway6 = pvetypes.NewIPAddrPointerValue(iface.Gateway6)
	m.Autostart = types.BoolPointerValue(iface.Autostart.PointerBool())

	if iface.MTU != nil {
		if v, err := strconv.Atoi(*iface.MTU); err == nil {
			m.MTU = types.Int64Value(int64(v))
		}
	} else {
		m.MTU = types.Int64Null()
	}

	if iface.Comments != nil {
		m.Comment = types.StringValue(strings.TrimSpace(*iface.Comments))
	} else {
		m.Comment = types.StringNull()
	}

	if iface.VLANID != nil {
		if v, err := strconv.Atoi(*iface.VLANID); err == nil {
			m.VLAN = types.Int64Value(int64(v))
		}
	} else {
		// in reality, this should never happen
		m.VLAN = types.Int64Unknown()
	}

	if iface.VLANRawDevice != nil {
		m.Interface = types.StringValue(strings.TrimSpace(*iface.VLANRawDevice))
	} else {
		m.Interface = types.StringNull()
	}
}

// NewLinuxVLANResource creates a new resource for managing Linux VLAN network interfaces.
func NewLinuxVLANResource() resource.Resource {
	return &linuxVLANResource{}
}

type linuxVLANResource struct {
	client proxmox.Client
}

func (r *linuxVLANResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_network_linux_vlan"
}

// Schema defines the schema for the resource.
func (r *linuxVLANResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Manages a Linux VLAN network interface in a Proxmox VE node.",
		Attributes: map[string]schema.Attribute{
			// Base attributes
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "A unique identifier with format '<node name>:<iface>'.",
			},
			"node_name": schema.StringAttribute{
				Description: "The name of the node.",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "The interface name.",
				MarkdownDescription: "The interface name. Either add the VLAN tag number to an existing interface name, " +
					"e.g. `ens18.21` (and do not set `interface` and `vlan`), or use custom name, e.g. `vlan_lab` " +
					"(`interface` and `vlan` are then required).",
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(3),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"address": schema.StringAttribute{
				Description: "The interface IPv4/CIDR address.",
				CustomType:  pvetypes.IPCIDRType{},
				Optional:    true,
			},
			"gateway": schema.StringAttribute{
				Description: "Default gateway address.",
				CustomType:  pvetypes.IPAddrType{},
				Optional:    true,
			},
			"address6": schema.StringAttribute{
				Description: "The interface IPv6/CIDR address.",
				CustomType:  pvetypes.IPCIDRType{},
				Optional:    true,
			},
			"gateway6": schema.StringAttribute{
				Description: "Default IPv6 gateway address.",
				CustomType:  pvetypes.IPAddrType{},
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
				Computed:    true,
			},
			"comment": schema.StringAttribute{
				Description: "Comment for the interface.",
				Optional:    true,
			},
			// Linux VLAN attributes
			"interface": schema.StringAttribute{
				Description: "The VLAN raw device. See also `name`.",
				Optional:    true,
				Computed:    true,
			},
			"vlan": schema.Int64Attribute{
				Description: "The VLAN tag. See also `name`.",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func (r *linuxVLANResource) Configure(
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

//nolint:dupl
func (r *linuxVLANResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan linuxVLANResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	body := plan.exportToNetworkInterfaceCreateUpdateBody()

	err := r.client.Node(plan.NodeName.ValueString()).CreateNetworkInterface(ctx, body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Linux VLAN interface",
			"Could not create Linux VLAN, unexpected error: "+err.Error(),
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

func (r *linuxVLANResource) read(ctx context.Context, model *linuxVLANResourceModel, diags *diag.Diagnostics) {
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

		model.importFromNetworkInterfaceList(iface)

		break
	}
}

// Read reads a Linux VLAN interface.
func (r *linuxVLANResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state linuxVLANResourceModel
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

// Update updates a Linux VLAN interface.
func (r *linuxVLANResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state linuxVLANResourceModel

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

	if len(toDelete) > 0 {
		body.Delete = &toDelete
	}

	err := r.client.Node(plan.NodeName.ValueString()).UpdateNetworkInterface(ctx, plan.Name.ValueString(), body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Linux VLAN interface",
			"Could not update Linux VLAN, unexpected error: "+err.Error(),
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

// Delete deletes a Linux VLAN interface.
//
//nolint:dupl
func (r *linuxVLANResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state linuxVLANResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Node(state.NodeName.ValueString()).DeleteNetworkInterface(ctx, state.Name.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "interface does not exist") {
			resp.Diagnostics.AddWarning(
				"Linux VLAN interface does not exist",
				fmt.Sprintf("Could not delete Linux VLAN '%s', interface does not exist, "+
					"or has already been deleted outside of Terraform.", state.Name.ValueString()),
			)
		} else {
			resp.Diagnostics.AddError(
				"Error deleting Linux VLAN interface",
				fmt.Sprintf("Could not delete Linux VLAN '%s', unexpected error: %s",
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

func (r *linuxVLANResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
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

	state := linuxVLANResourceModel{
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
