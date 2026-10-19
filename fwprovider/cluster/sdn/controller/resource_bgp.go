/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package controller

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/attribute"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/types/stringset"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/sdn/controllers"
	proxmoxtypes "github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

var (
	_ resource.ResourceWithConfigure   = &BGPResource{}
	_ resource.ResourceWithImportState = &BGPResource{}
)

type bgpModel struct {
	genericModel

	ASNumber       types.Int64     `tfsdk:"asn"`
	NodeName       types.String    `tfsdk:"node_name"`
	Peers          stringset.Value `tfsdk:"peers"`
	EBGP           types.Bool      `tfsdk:"ebgp"`
	EBGPMultiHop   types.Int64     `tfsdk:"ebgp_multihop"`
	Loopback       types.String    `tfsdk:"loopback"`
	MultiPathRelax types.Bool      `tfsdk:"bgp_multipath_as_path_relax"`
}

func (m *bgpModel) fromAPI(name string, data *controllers.ControllerData, diags *diag.Diagnostics) {
	m.genericModel.fromAPI(name, data, diags)

	m.ASNumber = types.Int64PointerValue(data.ASNumber)
	m.NodeName = m.nodeFromAPI(data.Nodes)
	m.EBGPMultiHop = multiHopFromAPI(data.EBPGMultiHop)
	m.EBGP = types.BoolPointerValue(data.EBGP.PointerBool())
	m.MultiPathRelax = types.BoolPointerValue(data.BgpMultiPathAsRelax.PointerBool())
	m.Loopback = m.handleDeletedStringValue(data.Loopback)

	if data.Peers != nil {
		peers := make([]string, len(*data.Peers))
		copy(peers, *data.Peers)
		m.Peers = m.handleDeletedStringSetValue(peers, diags)
	} else {
		m.Peers = stringset.NullValue()
	}

	if data.Pending != nil {
		if data.Pending.ASNumber != nil {
			m.ASNumber = types.Int64PointerValue(data.Pending.ASNumber)
		}

		if data.Pending.Nodes != nil {
			m.NodeName = m.nodeFromAPI(data.Pending.Nodes)
		}

		if data.Pending.EBPGMultiHop != nil {
			m.EBGPMultiHop = multiHopFromAPI(data.Pending.EBPGMultiHop)
		}

		if data.Pending.EBGP != nil {
			m.EBGP = types.BoolPointerValue(data.Pending.EBGP.PointerBool())
		}

		if data.Pending.BgpMultiPathAsRelax != nil {
			m.MultiPathRelax = types.BoolPointerValue(data.Pending.BgpMultiPathAsRelax.PointerBool())
		}

		if data.Pending.Loopback != nil {
			m.Loopback = m.handleDeletedStringValue(data.Pending.Loopback)
		}

		if data.Pending.Peers != nil {
			peers := make([]string, len(*data.Pending.Peers))
			copy(peers, *data.Pending.Peers)
			m.Peers = m.handleDeletedStringSetValue(peers, diags)
		}
	}
}

// nodeFromAPI reads the PVE `node` property, which the API layer models as a
// comma-separated list even though the BGP plugin binds a controller to a
// single node. The literal "deleted" pending marker maps to null.
func (m *bgpModel) nodeFromAPI(nodes *proxmoxtypes.CustomCommaSeparatedList) types.String {
	if nodes == nil {
		return types.StringNull()
	}

	if len(*nodes) == 0 {
		return types.StringNull()
	}

	if len(*nodes) == 1 && (*nodes)[0] == "deleted" {
		return types.StringNull()
	}

	if slices.Contains(*nodes, "") {
		return types.StringNull()
	}

	return types.StringValue(strings.Join(*nodes, ","))
}

func (m *bgpModel) fromAPIForDatasource(name string, data *controllers.ControllerData, diags *diag.Diagnostics) {
	m.genericModel.fromAPIForDatasource(name, data, diags)

	m.ASNumber = attribute.Int64ValueFromPtr(data.ASNumber)

	if data.Nodes != nil && len(*data.Nodes) > 0 {
		m.NodeName = types.StringValue(strings.Join(*data.Nodes, ","))
	} else {
		m.NodeName = types.StringValue("")
	}

	m.EBGP = attribute.BoolValueFromCustomBoolPtr(data.EBGP)
	m.MultiPathRelax = attribute.BoolValueFromCustomBoolPtr(data.BgpMultiPathAsRelax)
	m.Loopback = attribute.StringValueFromPtr(data.Loopback)

	if data.EBPGMultiHop != nil {
		m.EBGPMultiHop = types.Int64Value(int64(*data.EBPGMultiHop))
	} else {
		m.EBGPMultiHop = types.Int64Value(0)
	}

	if data.Peers != nil && len(*data.Peers) > 0 {
		peers := make([]string, len(*data.Peers))
		copy(peers, *data.Peers)
		m.Peers = stringset.NewValueList(peers, diags)
	} else {
		m.Peers = stringset.NewValueList([]string{}, diags)
	}
}

func (m *bgpModel) toAPI(ctx context.Context, diags *diag.Diagnostics) *controllers.Controller {
	data := m.genericModel.toAPI(ctx, diags)

	data.ASNumber = attribute.Int64PtrFromValue(m.ASNumber)

	if attribute.IsDefined(m.NodeName) {
		node := proxmoxtypes.CustomCommaSeparatedList([]string{m.NodeName.ValueString()})
		data.Nodes = &node
	}

	data.EBGP = attribute.CustomBoolPtrFromValue(m.EBGP)
	data.EBPGMultiHop = proxmoxtypes.CustomInt64OptionPtr(attribute.Int64PtrFromValue(m.EBGPMultiHop))
	data.Loopback = attribute.StringPtrFromValue(m.Loopback)
	data.BgpMultiPathAsRelax = attribute.CustomBoolPtrFromValue(m.MultiPathRelax)

	var peers []string
	diags.Append(m.Peers.ElementsAs(ctx, &peers, false)...)

	if len(peers) > 0 {
		p := proxmoxtypes.CustomCommaSeparatedList(peers)
		data.Peers = &p
	} else {
		data.Peers = nil
	}

	return data
}

func checkDeletedBGPFields(plan, state *bgpModel) []string {
	var toDelete []string

	attribute.CheckDelete(plan.EBGP, state.EBGP, &toDelete, "ebgp")
	attribute.CheckDelete(plan.EBGPMultiHop, state.EBGPMultiHop, &toDelete, "ebgp-multihop")
	attribute.CheckDelete(plan.Loopback, state.Loopback, &toDelete, "loopback")
	attribute.CheckDelete(plan.MultiPathRelax, state.MultiPathRelax, &toDelete, "bgp-multipath-as-path-relax")

	return toDelete
}

func multiHopFromAPI(v *proxmoxtypes.CustomInt64Option) types.Int64 {
	if v == nil || *v == 0 {
		return types.Int64Null()
	}

	return types.Int64Value(int64(*v))
}

// ebgpMultiHopRequiresEBGP returns true when ebgp_multihop is set but ebgp is
// not enabled: PVE ignores ebgp-multihop unless ebgp is true. Unknown values are
// not rejected, they are checked again once known.
func ebgpMultiHopRequiresEBGP(ebgp types.Bool, ebgpMultiHop types.Int64) bool {
	return attribute.IsDefined(ebgpMultiHop) && !ebgp.IsUnknown() && !ebgp.ValueBool()
}

type bgpEBGPMultiHopValidator struct{}

func (v bgpEBGPMultiHopValidator) Description(_ context.Context) string {
	return "Verifies that ebgp_multihop is only set when ebgp is set to true."
}

func (v bgpEBGPMultiHopValidator) MarkdownDescription(_ context.Context) string {
	return "Verifies that `ebgp_multihop` is only set when `ebgp` is set to `true`."
}

func (v bgpEBGPMultiHopValidator) ValidateResource(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var config bgpModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if ebgpMultiHopRequiresEBGP(config.EBGP, config.EBGPMultiHop) {
		resp.Diagnostics.AddAttributeError(
			path.Root("ebgp_multihop"),
			"Invalid Attribute Combination",
			"Attribute \"ebgp_multihop\" requires \"ebgp\" set to true",
		)
	}
}

type BGPResource struct {
	*genericControllerResource
}

func NewBGPResource() resource.Resource {
	return &BGPResource{
		genericControllerResource: newGenericControllerResource(controllerResourceConfig{
			typeNameSuffix: "_sdn_controller_bgp",
			controllerType: controllers.TypeBGP,
			modelFunc:      func() controllerModel { return &bgpModel{} },
		}).(*genericControllerResource),
	}
}

func (r *BGPResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan bgpModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var state bgpModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	updateController := plan.toAPI(ctx, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	update := &controllers.ControllerUpdate{
		Controller: *updateController,
		Delete:     checkDeletedBGPFields(&plan, &state),
	}

	if err := r.client.UpdateController(ctx, update); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to Update SDN Controller %q", plan.getID()),
			err.Error(),
		)

		return
	}

	r.readAndSetState(ctx, plan.getID(), &resp.State, &resp.Diagnostics)
}

func (r *BGPResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The BGP controller plugin configures the Free Range Routing (frr) BGP router of a node.",
		MarkdownDescription: "The BGP controller configures the Free Range Routing (frr) BGP router of a single node." +
			" Proxmox allows at most one BGP controller per node, and none while a BGP fabric is configured." +
			" SDN changes must be applied with the proxmox_sdn_applier resource before they affect the network.",
		Attributes: genericAttributesWith(map[string]schema.Attribute{
			"asn": schema.Int64Attribute{
				Description: "Autonomous System Number of the local BGP router.",
				Required:    true,
				Validators: []validator.Int64{
					int64validator.Between(0, 4294967295),
				},
			},
			"node_name": schema.StringAttribute{
				Description: "Name of the node this BGP controller is bound to. Required by the BGP controller plugin.",
				Required:    true,
			},
			"peers": stringset.ResourceAttribute(
				"Set of BGP peer IP addresses for the BGP controller.",
				"",
				stringset.WithRequired(),
			),
			"ebgp": schema.BoolAttribute{
				Description: "Whether eBGP is used (remote-as external) instead of iBGP.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"ebgp_multihop": schema.Int64Attribute{
				Description: "Maximum amount of hops for eBGP peers. Requires ebgp set to true.",
				Optional:    true,
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},
			"loopback": schema.StringAttribute{
				Description: "Name of the loopback/dummy interface that provides the Router-IP.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"bgp_multipath_as_path_relax": schema.BoolAttribute{
				Description: "Whether different AS paths of equal length are considered for multipath computation.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
		}),
	}
}

func (r *BGPResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		bgpEBGPMultiHopValidator{},
	}
}

func (m *bgpModel) getGenericModel() *genericModel {
	return &m.genericModel
}
