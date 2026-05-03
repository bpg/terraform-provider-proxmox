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
	"github.com/hashicorp/terraform-plugin-framework/path"
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
	"github.com/bpg/terraform-provider-proxmox/fwprovider/migration"
	customtypes "github.com/bpg/terraform-provider-proxmox/fwprovider/types"
	"github.com/bpg/terraform-provider-proxmox/proxmox"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes"
	proxmoxtypes "github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

var (
	_ resource.Resource                   = &linuxBridgeResource{}
	_ resource.ResourceWithConfigure      = &linuxBridgeResource{}
	_ resource.ResourceWithImportState    = &linuxBridgeResource{}
	_ resource.ResourceWithValidateConfig = &linuxBridgeResource{}
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
	Timeout   types.Int64             `tfsdk:"timeout_reload"`
	// Linux bridge attributes
	Ports     types.List   `tfsdk:"ports"`
	VLANAware types.Bool   `tfsdk:"vlan_aware"`
	VIDs      types.String `tfsdk:"vids"`
}

func (m *linuxBridgeResourceModel) exportToNetworkInterfaceCreateUpdateBody(
	ctx context.Context,
	diags *diag.Diagnostics,
) *nodes.NetworkInterfaceCreateUpdateRequestBody {
	body := &nodes.NetworkInterfaceCreateUpdateRequestBody{
		Iface:     m.Name.ValueString(),
		Type:      "bridge",
		Autostart: proxmoxtypes.CustomBool(m.Autostart.ValueBool()).Pointer(),
	}

	body.CIDR = m.Address.ValueStringPointer()
	body.Gateway = m.Gateway.ValueStringPointer()
	body.CIDR6 = m.Address6.ValueStringPointer()
	body.Gateway6 = m.Gateway6.ValueStringPointer()

	body.MTU = attribute.Int64PtrFromValue(m.MTU)

	body.Comments = attribute.StringPtrFromValue(m.Comment)

	var sanitizedPorts []string

	if attribute.IsDefined(m.Ports) {
		var portsList []string

		d := m.Ports.ElementsAs(ctx, &portsList, false)
		diags.Append(d...)

		if d.HasError() {
			return body
		}

		for _, port := range portsList {
			port = strings.TrimSpace(port)
			if len(port) > 0 {
				sanitizedPorts = append(sanitizedPorts, port)
			}
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

	body.BridgeVIDs = attribute.StringPtrFromValue(m.VIDs)

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

	// When comments is nil in PVE (e.g. after setting comment=""), we preserve the model's
	// current value to avoid "inconsistent result after apply" errors. The CheckDelete in
	// Update handles actual comment removal.
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

		m.Ports = ports
	} else if m.Ports.ElementType(ctx) == nil {
		// state.Set rejects a zero-value List (no element type) on fresh ImportState.
		m.Ports = types.ListNull(types.StringType)
	}

	if iface.BridgeVIDs != nil {
		m.VIDs = types.StringValue(*iface.BridgeVIDs)
	} else {
		m.VIDs = types.StringNull()
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
		DeprecationMessage: migration.DeprecationMessage("proxmox_network_linux_bridge"),
		Description:        "Manages a Linux Bridge network interface in a Proxmox VE node.",
		Attributes: map[string]schema.Attribute{
			// Base attributes
			"id": attribute.ResourceID("A unique identifier with format `<node name>:<iface>`"),
			"node_name": schema.StringAttribute{
				Description: "The name of the node.",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "The interface name.",
				MarkdownDescription: "The interface name. Commonly vmbr[N], where 0 ≤ N ≤ 4094 (vmbr0 - vmbr4094), but " +
					"can be any string containing only letters, numbers, and underscores (_), starting with a letter and at most 10 characters long.",
				Required: true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[A-Za-z][A-Za-z0-9_]{0,9}$`),
						`must contain only letters, numbers, and underscores (_), start with a letter, and be no longer than 10 characters`,
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
			"vids": schema.StringAttribute{
				Description: "VLAN IDs allowed on the bridge (Linux Bridge `bridge-vids`).",
				MarkdownDescription: "VLAN IDs allowed on the bridge (Linux Bridge `bridge-vids`). " +
					"Space-separated list of VLAN IDs and/or hyphenated ranges " +
					"(e.g. `\"2-4094\"`, `\"1 20 130\"`, or `\"1 10-20 30\"`). " +
					"Requires `vlan_aware = true`. PVE/ifupdown2 fills in `2-4094` as the " +
					"implicit default for VLAN-aware bridges when this attribute is omitted; " +
					"the provider surfaces that default in state.",
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^\d+(-\d+)?( \d+(-\d+)?)*$`),
						`must be a space-separated list of VLAN IDs or ranges (e.g. "1 20 130" or "2-4094")`,
					),
				},
				PlanModifiers: []planmodifier.String{
					vidsPlanModifier{},
				},
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

// ValidateConfig validates the resource configuration.
func (r *linuxBridgeResource) ValidateConfig(
	ctx context.Context,
	req resource.ValidateConfigRequest,
	resp *resource.ValidateConfigResponse,
) {
	var data linuxBridgeResourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// `vids` (bridge-vids) only applies to VLAN-aware bridges. Setting it
	// without `vlan_aware = true` has no effect on PVE — surface the
	// dependency at plan time so the misconfiguration is loud, not silent.
	// Skip only when `vlan_aware` is unknown (e.g. a cross-resource reference);
	// null means the user omitted it from config and PVE defaults to false, which
	// is still a misconfiguration when `vids` is set.
	if attribute.IsDefined(data.VIDs) && !data.VLANAware.IsUnknown() && !data.VLANAware.ValueBool() {
		resp.Diagnostics.AddAttributeError(
			path.Root("vids"),
			"Invalid attribute combination",
			"The `vids` attribute requires `vlan_aware = true`. A bridge that is "+
				"not VLAN-aware does not perform VLAN filtering, so `vids` has no "+
				"effect. Either set `vlan_aware = true`, or remove `vids`.",
		)
	}
}

func (r *linuxBridgeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan linuxBridgeResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	body := plan.exportToNetworkInterfaceCreateUpdateBody(ctx, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Node(plan.NodeName.ValueString()).CreateNetworkInterface(ctx, body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Linux Bridge interface",
			"Could not create Linux Bridge, unexpected error: "+err.Error(),
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
			"Linux Bridge interface not found after creation",
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

func (r *linuxBridgeResource) read(ctx context.Context, model *linuxBridgeResourceModel, diags *diag.Diagnostics) bool {
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

// Read reads a Linux Bridge interface.
func (r *linuxBridgeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state linuxBridgeResourceModel

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

// Update updates a Linux Bridge interface.
func (r *linuxBridgeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state linuxBridgeResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	body := plan.exportToNetworkInterfaceCreateUpdateBody(ctx, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	var toDelete []string

	attribute.CheckDelete(plan.Address, state.Address, &toDelete, "cidr")
	attribute.CheckDelete(plan.Address6, state.Address6, &toDelete, "cidr6")
	attribute.CheckDelete(plan.MTU, state.MTU, &toDelete, "mtu")
	attribute.CheckDelete(plan.Gateway, state.Gateway, &toDelete, "gateway")
	attribute.CheckDelete(plan.Gateway6, state.Gateway6, &toDelete, "gateway6")
	attribute.CheckDelete(plan.Comment, state.Comment, &toDelete, "comments")
	attribute.CheckDelete(plan.VIDs, state.VIDs, &toDelete, "bridge_vids")

	// VLANAware is computed with a default, will never be null
	if !plan.VLANAware.Equal(state.VLANAware) && !plan.VLANAware.ValueBool() {
		toDelete = append(toDelete, "bridge_vlan_aware")
		body.BridgeVLANAware = nil
	}

	if len(toDelete) > 0 {
		body.Delete = toDelete
	}

	err := r.client.Node(plan.NodeName.ValueString()).UpdateNetworkInterface(ctx, plan.Name.ValueString(), body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Linux Bridge interface",
			"Could not update Linux Bridge, unexpected error: "+err.Error(),
		)

		return
	}

	found := r.read(ctx, &plan, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	if !found {
		resp.Diagnostics.AddError(
			"Linux Bridge interface not found after update",
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

// Delete deletes a Linux Bridge interface.
//

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

//nolint:dupl // ImportState mirrors linux_vlan and linux_bond but is bound to a distinct resource type
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
		Timeout:  types.Int64Value(int64(nodes.NetworkReloadTimeout.Seconds())),
	}
	found := r.read(ctx, &state, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	if !found {
		resp.Diagnostics.AddError(
			"Linux Bridge interface not found",
			fmt.Sprintf("Interface %q on node %q could not be imported", iface, nodeName),
		)

		return
	}

	diags := resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}
