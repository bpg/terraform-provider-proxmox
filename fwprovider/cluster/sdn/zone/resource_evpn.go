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
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

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

func (m *evpnModel) importFromAPI(name string, data *zones.ZoneData, diags *diag.Diagnostics) {
	m.genericModel.importFromAPI(name, data, diags)

	m.AdvertiseSubnets = types.BoolPointerValue(data.AdvertiseSubnets.PointerBool())
	m.Controller = types.StringPointerValue(data.Controller)
	m.DisableARPNDSuppression = types.BoolPointerValue(data.DisableARPNDSuppression.PointerBool())
	m.ExitNodes = stringset.NewValueString(data.ExitNodes, diags, stringset.WithSeparator(","))
	m.ExitNodesLocalRouting = types.BoolPointerValue(data.ExitNodesLocalRouting.PointerBool())
	m.PrimaryExitNode = types.StringPointerValue(data.ExitNodesPrimary)
	m.RouteTargetImport = types.StringPointerValue(data.RouteTargetImport)
	m.VRFVXLANID = types.Int64PointerValue(data.VRFVXLANID)
}

func (m *evpnModel) toAPIRequestBody(ctx context.Context, diags *diag.Diagnostics) *zones.ZoneRequestData {
	data := m.genericModel.toAPIRequestBody(ctx, diags)

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

type EVPNResource struct {
	generic *genericZoneResource
}

func NewEVPNResource() resource.Resource {
	return &EVPNResource{
		generic: newGenericZoneResource(zoneResourceConfig{
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
			},
			"controller": schema.StringAttribute{
				Description: "EVPN controller address.",
				Required:    true,
			},
			"disable_arp_nd_suppression": schema.BoolAttribute{
				Description: "Disable ARP/ND suppression for EVPN.",
				Optional:    true,
			},
			"exit_nodes": stringset.ResourceAttribute("List of exit nodes for EVPN.", ""),
			"exit_nodes_local_routing": schema.BoolAttribute{
				Description: "Enable local routing for EVPN exit nodes.",
				Optional:    true,
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

func (r *EVPNResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.generic.Metadata(ctx, req, resp)
}

func (r *EVPNResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.generic.Configure(ctx, req, resp)
}

func (r *EVPNResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	r.generic.Create(ctx, req, resp)
}

func (r *EVPNResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	r.generic.Read(ctx, req, resp)
}

func (r *EVPNResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	r.generic.Update(ctx, req, resp)
}

func (r *EVPNResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	r.generic.Delete(ctx, req, resp)
}

func (r *EVPNResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	r.generic.ImportState(ctx, req, resp)
}
