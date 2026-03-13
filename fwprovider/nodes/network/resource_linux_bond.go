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
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/attribute"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/config"
	customtypes "github.com/bpg/terraform-provider-proxmox/fwprovider/types"
	"github.com/bpg/terraform-provider-proxmox/proxmox"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes"
	proxmoxtypes "github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

var (
	_ resource.Resource                = &linuxBondResource{}
	_ resource.ResourceWithConfigure   = &linuxBondResource{}
	_ resource.ResourceWithImportState = &linuxBondResource{}
)

type linuxBondResourceModel struct {
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
	Timeout   types.Int64             `tfsdk:"timeout_reload"`
	// Linux bond attributes
	Slaves             []types.String `tfsdk:"slaves"`
	BondMode           types.String   `tfsdk:"bond_mode"`
	BondPrimary        types.String   `tfsdk:"bond_primary"`
	BondXmitHashPolicy types.String   `tfsdk:"bond_xmit_hash_policy"`
}

func (m *linuxBondResourceModel) exportToNetworkInterfaceCreateUpdateBody() *nodes.NetworkInterfaceCreateUpdateRequestBody {
	body := &nodes.NetworkInterfaceCreateUpdateRequestBody{
		Iface:     m.Name.ValueString(),
		Type:      "bond",
		Autostart: proxmoxtypes.CustomBool(m.Autostart.ValueBool()).Pointer(),
	}

	body.CIDR = m.Address.ValueStringPointer()
	body.Gateway = m.Gateway.ValueStringPointer()
	body.CIDR6 = m.Address6.ValueStringPointer()
	body.Gateway6 = m.Gateway6.ValueStringPointer()
	body.Comments = m.Comment.ValueStringPointer()

	if !m.MTU.IsUnknown() {
		body.MTU = m.MTU.ValueInt64Pointer()
	}

	var sanitizedSlaves []string

	for _, slave := range m.Slaves {
		s := strings.TrimSpace(slave.ValueString())
		if len(s) > 0 {
			sanitizedSlaves = append(sanitizedSlaves, s)
		}
	}

	sort.Strings(sanitizedSlaves)
	slaves := strings.Join(sanitizedSlaves, " ")

	if len(slaves) > 0 {
		body.Slaves = &slaves
	}

	body.BondMode = m.BondMode.ValueStringPointer()
	body.BondPrimary = m.BondPrimary.ValueStringPointer()
	body.BondXmitHashPolicy = m.BondXmitHashPolicy.ValueStringPointer()

	return body
}

func (m *linuxBondResourceModel) importFromNetworkInterfaceList(
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
		} else {
			m.MTU = types.Int64Null()
		}
	} else {
		m.MTU = types.Int64Null()
	}

	if iface.Comments != nil {
		m.Comment = types.StringValue(strings.TrimSpace(*iface.Comments))
	}

	if iface.BondMode != nil {
		m.BondMode = types.StringValue(*iface.BondMode)
	} else {
		m.BondMode = types.StringNull()
	}

	if iface.BondPrimary != nil && *iface.BondPrimary != "" {
		m.BondPrimary = types.StringValue(*iface.BondPrimary)
	} else {
		m.BondPrimary = types.StringNull()
	}

	if iface.BondXmitHashPolicy != nil && *iface.BondXmitHashPolicy != "" {
		m.BondXmitHashPolicy = types.StringValue(*iface.BondXmitHashPolicy)
	} else {
		m.BondXmitHashPolicy = types.StringNull()
	}

	if iface.Slaves != nil && len(*iface.Slaves) > 0 {
		slaves, diags := types.ListValueFrom(ctx, types.StringType, strings.Split(*iface.Slaves, " "))
		if diags.HasError() {
			return fmt.Errorf("failed to parse bond slaves: %s", *iface.Slaves)
		}

		diags = slaves.ElementsAs(ctx, &m.Slaves, false)
		if diags.HasError() {
			return fmt.Errorf("failed to build bond slaves list: %s", *iface.Slaves)
		}
	}

	return nil
}

// NewLinuxBondResource creates a new resource for managing Linux Bond network interfaces.
func NewLinuxBondResource() resource.Resource {
	return &linuxBondResource{}
}

type linuxBondResource struct {
	client proxmox.Client
}

func (r *linuxBondResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_network_linux_bond"
}

// Schema defines the schema for the resource.
func (r *linuxBondResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description:         "Manages a Linux Bond network interface in a Proxmox VE node.",
		MarkdownDescription: "Manages a Linux Bond network interface in a Proxmox VE node.",
		Attributes: map[string]schema.Attribute{
			// Base attributes
			"id": attribute.ResourceID("A unique identifier with format `<node name>:<iface>`"),
			"node_name": schema.StringAttribute{
				Description: "The name of the node.",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "The interface name.",
				MarkdownDescription: "The interface name. Must be `bond[N]`, where 0 ≤ N (e.g. bond0, bond1), " +
					"or any alphanumeric string that starts with a character and is at most 10 characters long.",
				Required: true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[A-Za-z][A-Za-z0-9]{0,9}$`),
						`must be an alphanumeric string that starts with a character and is at most 10 characters long`,
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
			"timeout_reload": schema.Int64Attribute{
				Description: "Timeout for network reload operations in seconds (defaults to `100`).",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(int64(nodes.NetworkReloadTimeout.Seconds())),
				Validators: []validator.Int64{
					int64validator.AtLeast(5),
				},
			},
			// Linux Bond attributes
			"slaves": schema.ListAttribute{
				Description:         "The interface bond slaves (member interfaces).",
				MarkdownDescription: "The interface bond slaves (member interfaces).",
				Required:            true,
				ElementType:         types.StringType,
			},
			"bond_mode": schema.StringAttribute{
				Description: "The bonding mode (defaults to `balance-rr`).",
				MarkdownDescription: "The bonding mode. Possible values are `balance-rr`, `active-backup`, `balance-xor`, " +
					"`broadcast`, `802.3ad`, `balance-tlb`, `balance-alb` (defaults to `balance-rr`).",
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"balance-rr",
						"active-backup",
						"balance-xor",
						"broadcast",
						"802.3ad",
						"balance-tlb",
						"balance-alb",
					),
				},
			},
			"bond_primary": schema.StringAttribute{
				Description: "The primary interface for active-backup bond mode.",
				MarkdownDescription: "The primary interface for `active-backup` bond mode. " +
					"Specifies which slave interface should be the active one.",
				Optional: true,
			},
			"bond_xmit_hash_policy": schema.StringAttribute{
				Description: "The transmit hash policy for balance-xor and 802.3ad bond modes.",
				MarkdownDescription: "The transmit hash policy for `balance-xor` and `802.3ad` bond modes. " +
					"Possible values are `layer2`, `layer2+3`, `layer3+4`.",
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"layer2",
						"layer2+3",
						"layer3+4",
					),
				},
			},
		},
	}
}

func (r *linuxBondResource) Configure(
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
			fmt.Sprintf("Expected config.Resource, got: %T", req.ProviderData),
		)

		return
	}

	r.client = cfg.Client
}

//nolint:dupl
func (r *linuxBondResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan linuxBondResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	body := plan.exportToNetworkInterfaceCreateUpdateBody()

	err := r.client.Node(plan.NodeName.ValueString()).CreateNetworkInterface(ctx, body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Linux Bond interface",
			"Could not create Linux Bond, unexpected error: "+err.Error(),
		)

		return
	}

	plan.ID = types.StringValue(plan.NodeName.ValueString() + ":" + plan.Name.ValueString())

	found := r.read(ctx, &plan, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	if !found {
		resp.Diagnostics.AddError(
			"Linux Bond interface not found after creation",
			fmt.Sprintf(
				"Interface %q on node %q could not be read after creation",
				plan.Name.ValueString(), plan.NodeName.ValueString()),
		)

		return
	}

	resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)

	reloadCtx, cancel := context.WithTimeout(ctx, time.Duration(plan.Timeout.ValueInt64())*time.Second)
	defer cancel()

	err = r.client.Node(plan.NodeName.ValueString()).ReloadNetworkConfiguration(reloadCtx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reloading network configuration",
			fmt.Sprintf("Could not reload network configuration on node '%s', unexpected error: %s",
				plan.NodeName.ValueString(), err.Error()),
		)
	}
}

func (r *linuxBondResource) read(ctx context.Context, model *linuxBondResourceModel, diags *diag.Diagnostics) bool {
	ifaces, err := r.client.Node(model.NodeName.ValueString()).ListNetworkInterfaces(ctx)
	if err != nil {
		diags.AddError(
			"Error listing network interfaces",
			"Could not list network interfaces, unexpected error: "+err.Error(),
		)

		return false
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

			return false
		}

		return true
	}

	return false
}

// Read reads a Linux Bond interface.
func (r *linuxBondResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state linuxBondResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	found := r.read(ctx, &state, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	if !found {
		resp.State.RemoveResource(ctx)
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// Update updates a Linux Bond interface.
func (r *linuxBondResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state linuxBondResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	body := plan.exportToNetworkInterfaceCreateUpdateBody()

	var toDelete []string

	attribute.CheckDelete(plan.Address, state.Address, &toDelete, "cidr")
	attribute.CheckDelete(plan.Address6, state.Address6, &toDelete, "cidr6")
	attribute.CheckDelete(plan.MTU, state.MTU, &toDelete, "mtu")
	attribute.CheckDelete(plan.Gateway, state.Gateway, &toDelete, "gateway")
	attribute.CheckDelete(plan.Gateway6, state.Gateway6, &toDelete, "gateway6")
	attribute.CheckDelete(plan.BondPrimary, state.BondPrimary, &toDelete, "bond-primary")
	attribute.CheckDelete(plan.BondXmitHashPolicy, state.BondXmitHashPolicy, &toDelete, "bond_xmit_hash_policy")

	if len(toDelete) > 0 {
		body.Delete = toDelete
	}

	err := r.client.Node(plan.NodeName.ValueString()).UpdateNetworkInterface(ctx, plan.Name.ValueString(), body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Linux Bond interface",
			"Could not update Linux Bond, unexpected error: "+err.Error(),
		)

		return
	}

	found := r.read(ctx, &plan, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	if !found {
		resp.Diagnostics.AddError(
			"Linux Bond interface not found after update",
			fmt.Sprintf(
				"Interface %q on node %q could not be read after update",
				plan.Name.ValueString(), plan.NodeName.ValueString()),
		)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)

	reloadCtx, cancel := context.WithTimeout(ctx, time.Duration(plan.Timeout.ValueInt64())*time.Second)
	defer cancel()

	err = r.client.Node(plan.NodeName.ValueString()).ReloadNetworkConfiguration(reloadCtx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reloading network configuration",
			fmt.Sprintf("Could not reload network configuration on node '%s', unexpected error: %s",
				plan.NodeName.ValueString(), err.Error()),
		)
	}
}

// Delete deletes a Linux Bond interface.
//
//nolint:dupl
func (r *linuxBondResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state linuxBondResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Node(state.NodeName.ValueString()).DeleteNetworkInterface(ctx, state.Name.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "interface does not exist") {
			resp.Diagnostics.AddWarning(
				"Linux Bond interface does not exist",
				fmt.Sprintf("Could not delete Linux Bond '%s', interface does not exist, "+
					"or has already been deleted outside of Terraform.", state.Name.ValueString()),
			)
		} else {
			resp.Diagnostics.AddError(
				"Error deleting Linux Bond interface",
				fmt.Sprintf("Could not delete Linux Bond '%s', unexpected error: %s",
					state.Name.ValueString(), err.Error()),
			)
		}

		return
	}

	reloadCtx, cancel := context.WithTimeout(ctx, time.Duration(state.Timeout.ValueInt64())*time.Second)
	defer cancel()

	err = r.client.Node(state.NodeName.ValueString()).ReloadNetworkConfiguration(reloadCtx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reloading network configuration",
			fmt.Sprintf("Could not reload network configuration on node '%s', unexpected error: %s",
				state.NodeName.ValueString(), err.Error()),
		)
	}
}

func (r *linuxBondResource) ImportState(
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

	state := linuxBondResourceModel{
		ID:       types.StringValue(req.ID),
		NodeName: types.StringValue(nodeName),
		Name:     types.StringValue(iface),
		Timeout:  types.Int64Value(int64(nodes.NetworkReloadTimeout.Seconds())),
	}
	found := r.read(ctx, &state, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	if !found {
		resp.Diagnostics.AddError(
			"Linux Bond interface not found",
			fmt.Sprintf("Interface %q on node %q could not be imported", iface, nodeName),
		)

		return
	}

	diags := resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}
