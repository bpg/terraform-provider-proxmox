/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package zone

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/attribute"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/types/stringset"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/sdn/zones"
	proxmoxtypes "github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

var (
	_ resource.ResourceWithConfigure   = &EVPNResource{}
	_ resource.ResourceWithImportState = &EVPNResource{}
)

type evpnModel struct {
	genericModel

	AdvertiseSubnets        types.Bool      `tfsdk:"advertise_subnets"`
	Controller              types.String    `tfsdk:"controller"`
	DisableARPNDSuppression types.Bool      `tfsdk:"disable_arp_nd_suppression"`
	ExitNodes               stringset.Value `tfsdk:"exit_nodes"`
	ExitNodesLocalRouting   types.Bool      `tfsdk:"exit_nodes_local_routing"`
	PrimaryExitNode         types.String    `tfsdk:"primary_exit_node"`
	RouteTargetImport       types.String    `tfsdk:"rt_import"`
	VRFVXLANID              types.Int64     `tfsdk:"vrf_vxlan"`
}

func (m *evpnModel) fromAPI(name string, data *zones.ZoneData, diags *diag.Diagnostics) {
	m.genericModel.fromAPI(name, data, diags)

	m.AdvertiseSubnets = types.BoolPointerValue(data.AdvertiseSubnets.PointerBool())
	m.Controller = types.StringPointerValue(data.Controller)
	m.DisableARPNDSuppression = types.BoolPointerValue(data.DisableARPNDSuppression.PointerBool())
	m.ExitNodes = stringset.NewValueString(data.ExitNodes, diags, stringset.WithSeparator(","))
	m.ExitNodesLocalRouting = types.BoolPointerValue(data.ExitNodesLocalRouting.PointerBool())
	m.PrimaryExitNode = types.StringPointerValue(data.ExitNodesPrimary)
	m.RouteTargetImport = types.StringPointerValue(data.RouteTargetImport)
	m.VRFVXLANID = types.Int64PointerValue(data.VRFVXLANID)

	if data.Pending != nil {
		hasPendingChanges := false

		if data.Pending.AdvertiseSubnets != nil {
			m.applyPendingBool(data.Pending.AdvertiseSubnets, &m.AdvertiseSubnets)

			hasPendingChanges = true
		}

		if data.Pending.Controller != nil && *data.Pending.Controller != "" {
			m.applyPendingString(data.Pending.Controller, &m.Controller)

			hasPendingChanges = true
		}

		if data.Pending.DisableARPNDSuppression != nil {
			m.applyPendingBool(data.Pending.DisableARPNDSuppression, &m.DisableARPNDSuppression)

			hasPendingChanges = true
		}

		if data.Pending.ExitNodes != nil && *data.Pending.ExitNodes != "" {
			m.ExitNodes = stringset.NewValueString(data.Pending.ExitNodes, diags, stringset.WithSeparator(","))
			hasPendingChanges = true
		}

		if data.Pending.ExitNodesLocalRouting != nil {
			m.applyPendingBool(data.Pending.ExitNodesLocalRouting, &m.ExitNodesLocalRouting)

			hasPendingChanges = true
		}

		if data.Pending.ExitNodesPrimary != nil && *data.Pending.ExitNodesPrimary != "" {
			m.applyPendingString(data.Pending.ExitNodesPrimary, &m.PrimaryExitNode)

			hasPendingChanges = true
		}

		if data.Pending.RouteTargetImport != nil && *data.Pending.RouteTargetImport != "" {
			m.applyPendingString(data.Pending.RouteTargetImport, &m.RouteTargetImport)

			hasPendingChanges = true
		}

		if data.Pending.VRFVXLANID != nil && *data.Pending.VRFVXLANID != 0 {
			m.VRFVXLANID = types.Int64Value(*data.Pending.VRFVXLANID)
			hasPendingChanges = true
		}

		if hasPendingChanges {
			m.Pending = types.BoolValue(true)
		}
	}
}

func (m *evpnModel) toAPI(ctx context.Context, diags *diag.Diagnostics) *zones.Zone {
	data := m.genericModel.toAPI(ctx, diags)

	data.AdvertiseSubnets = proxmoxtypes.CustomBoolPtr(m.AdvertiseSubnets.ValueBoolPointer())
	data.Controller = m.Controller.ValueStringPointer()
	data.DisableARPNDSuppression = proxmoxtypes.CustomBoolPtr(m.DisableARPNDSuppression.ValueBoolPointer())
	data.ExitNodes = m.ExitNodes.ValueStringPointer(ctx, diags, stringset.WithSeparator(","))
	data.ExitNodesLocalRouting = proxmoxtypes.CustomBoolPtr(m.ExitNodesLocalRouting.ValueBoolPointer())
	data.ExitNodesPrimary = m.PrimaryExitNode.ValueStringPointer()
	data.RouteTargetImport = m.RouteTargetImport.ValueStringPointer()
	data.VRFVXLANID = m.VRFVXLANID.ValueInt64Pointer()

	return data
}

func (m *evpnModel) checkDeletedFields(state zoneModel) []string {
	evpnState := state.(*evpnModel)
	toDelete := m.genericModel.checkDeletedFields(evpnState.getGenericModel())

	// Add EVPN-specific deleted fields
	attribute.CheckDelete(m.AdvertiseSubnets, evpnState.AdvertiseSubnets, &toDelete, "advertise-subnets")
	attribute.CheckDelete(m.DisableARPNDSuppression, evpnState.DisableARPNDSuppression, &toDelete, "disable-arp-nd-suppression")
	attribute.CheckDelete(m.ExitNodes, evpnState.ExitNodes, &toDelete, "exitnodes")
	attribute.CheckDelete(m.ExitNodesLocalRouting, evpnState.ExitNodesLocalRouting, &toDelete, "exitnodes-local-routing")
	attribute.CheckDelete(m.PrimaryExitNode, evpnState.PrimaryExitNode, &toDelete, "exitnodes-primary")
	attribute.CheckDelete(m.RouteTargetImport, evpnState.RouteTargetImport, &toDelete, "rt-import")

	return toDelete
}

type EVPNResource struct {
	*genericZoneResource
}

func NewEVPNResource() resource.Resource {
	return &EVPNResource{
		genericZoneResource: newGenericZoneResource(zoneResourceConfig{
			typeNameSuffix: "_sdn_zone_evpn",
			zoneType:       zones.TypeEVPN,
			modelFunc:      func() zoneModel { return &evpnModel{} },
		}).(*genericZoneResource),
	}
}

func (r *EVPNResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "EVPN Zone in Proxmox SDN.",
		MarkdownDescription: "EVPN Zone in Proxmox SDN. The EVPN zone creates a routable Layer 3 network, capable of " +
			"spanning across multiple clusters.",
		Attributes: genericAttributesWith(map[string]schema.Attribute{
			"advertise_subnets": schema.BoolAttribute{
				Description: "Enable subnet advertisement for EVPN.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"controller": schema.StringAttribute{
				Description: "EVPN controller address.",
				Required:    true,
			},
			"disable_arp_nd_suppression": schema.BoolAttribute{
				Description: "Disable ARP/ND suppression for EVPN.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"exit_nodes": stringset.ResourceAttribute("List of exit nodes for EVPN.", ""),
			"exit_nodes_local_routing": schema.BoolAttribute{
				Description: "Enable local routing for EVPN exit nodes.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"primary_exit_node": schema.StringAttribute{
				Description: "Primary exit node for EVPN.",
				Optional:    true,
			},
			"rt_import": schema.StringAttribute{
				Description: "Route target import for EVPN.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^(\d+):(\d+)$`),
						"must be in the format '<ASN>:<number>' (e.g., '65000:65000')",
					),
				},
				Optional: true,
			},
			"vrf_vxlan": schema.Int64Attribute{
				Description: "VRF VXLAN-ID used for dedicated routing interconnect between VNets. It must be different " +
					"than the VXLAN-ID of the VNets.",
				Required: true,
			},
		}),
	}
}

func (m *evpnModel) getGenericModel() *genericModel {
	return &m.genericModel
}
