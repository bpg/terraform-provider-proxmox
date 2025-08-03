/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package zone

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/config"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/types/stringset"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/sdn/zones"
)

// Ensure the implementation satisfies the required interfaces.
var (
	_ datasource.DataSource              = &zonesDataSource{}
	_ datasource.DataSourceWithConfigure = &zonesDataSource{}
)

// zonesDataSource is the data source implementation for SDN zones.
type zonesDataSource struct {
	client *zones.Client
}

// zonesDataSourceModel represents the data source model for listing zones.
type zonesDataSourceModel struct {
	Type  types.String `tfsdk:"type"`
	Zones types.List   `tfsdk:"zones"`
}

// zoneDataModel represents individual zone data in the list.
type zoneDataModel struct {
	ID                      types.String    `tfsdk:"id"`
	Type                    types.String    `tfsdk:"type"`
	IPAM                    types.String    `tfsdk:"ipam"`
	DNS                     types.String    `tfsdk:"dns"`
	ReverseDNS              types.String    `tfsdk:"reverse_dns"`
	DNSZone                 types.String    `tfsdk:"dns_zone"`
	Nodes                   stringset.Value `tfsdk:"nodes"`
	MTU                     types.Int64     `tfsdk:"mtu"`
	Bridge                  types.String    `tfsdk:"bridge"`
	ServiceVLAN             types.Int64     `tfsdk:"service_vlan"`
	ServiceVLANProtocol     types.String    `tfsdk:"service_vlan_protocol"`
	Peers                   stringset.Value `tfsdk:"peers"`
	AdvertiseSubnets        types.Bool      `tfsdk:"advertise_subnets"`
	Controller              types.String    `tfsdk:"controller"`
	DisableARPNDSuppression types.Bool      `tfsdk:"disable_arp_nd_suppression"`
	ExitNodes               stringset.Value `tfsdk:"exit_nodes"`
	ExitNodesLocalRouting   types.Bool      `tfsdk:"exit_nodes_local_routing"`
	PrimaryExitNode         types.String    `tfsdk:"primary_exit_node"`
	RouteTargetImport       types.String    `tfsdk:"rt_import"`
	VRFVXLANID              types.Int64     `tfsdk:"vrf_vxlan"`
}

// Configure adds the provider-configured client to the data source.
func (d *zonesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	cfg, ok := req.ProviderData.(config.DataSource)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected DataSource Configure Type",
			fmt.Sprintf("Expected config.DataSource, got: %T", req.ProviderData),
		)

		return
	}

	d.client = cfg.Client.Cluster().SDNZones()
}

// Metadata returns the data source type name.
func (d *zonesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sdn_zones"
}

// Schema defines the schema for the data source.
func (d *zonesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves information about all SDN Zones in Proxmox.",
		MarkdownDescription: "Retrieves information about all SDN Zones in Proxmox. " +
			"This data source can optionally filter zones by type.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "Filter zones by type (simple, vlan, qinq, vxlan, evpn).",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("simple", "vlan", "qinq", "vxlan", "evpn"),
				},
			},
			"zones": schema.ListAttribute{
				Description: "List of SDN zones.",
				Computed:    true,
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"id":          types.StringType,
						"type":        types.StringType,
						"ipam":        types.StringType,
						"dns":         types.StringType,
						"reverse_dns": types.StringType,
						"dns_zone":    types.StringType,
						"nodes": types.SetType{
							ElemType: types.StringType,
						},
						"mtu":                   types.Int64Type,
						"bridge":                types.StringType,
						"service_vlan":          types.Int64Type,
						"service_vlan_protocol": types.StringType,
						"peers": types.SetType{
							ElemType: types.StringType,
						},
						"advertise_subnets":          types.BoolType,
						"controller":                 types.StringType,
						"disable_arp_nd_suppression": types.BoolType,
						"exit_nodes": types.SetType{
							ElemType: types.StringType,
						},
						"exit_nodes_local_routing": types.BoolType,
						"primary_exit_node":        types.StringType,
						"rt_import":                types.StringType,
						"vrf_vxlan":                types.Int64Type,
					},
				},
			},
		},
	}
}

// Read fetches all SDN zones from the Proxmox VE API.
func (d *zonesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data zonesDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	zonesList, err := d.client.GetZones(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read SDN Zones",
			err.Error(),
		)

		return
	}

	filteredZones := zonesList

	if !data.Type.IsNull() && !data.Type.IsUnknown() {
		filterType := data.Type.ValueString()
		filteredZones = make([]zones.ZoneData, 0)

		for _, zone := range zonesList {
			if zone.Type != nil && *zone.Type == filterType {
				filteredZones = append(filteredZones, zone)
			}
		}
	}

	// Convert zones to list elements
	zoneElements := make([]attr.Value, len(filteredZones))
	for i, zone := range filteredZones {
		diags := &resp.Diagnostics

		zoneData := zoneDataModel{
			ID:                      types.StringValue(zone.ID),
			Type:                    types.StringPointerValue(zone.Type),
			IPAM:                    types.StringPointerValue(zone.IPAM),
			DNS:                     types.StringPointerValue(zone.DNS),
			ReverseDNS:              types.StringPointerValue(zone.ReverseDNS),
			DNSZone:                 types.StringPointerValue(zone.DNSZone),
			Nodes:                   stringset.NewValueString(zone.Nodes, diags, stringset.WithSeparator(",")),
			MTU:                     types.Int64PointerValue(zone.MTU),
			Bridge:                  types.StringPointerValue(zone.Bridge),
			ServiceVLAN:             types.Int64PointerValue(zone.ServiceVLAN),
			ServiceVLANProtocol:     types.StringPointerValue(zone.ServiceVLANProtocol),
			Peers:                   stringset.NewValueString(zone.Peers, diags, stringset.WithSeparator(",")),
			AdvertiseSubnets:        types.BoolPointerValue(zone.AdvertiseSubnets.PointerBool()),
			Controller:              types.StringPointerValue(zone.Controller),
			DisableARPNDSuppression: types.BoolPointerValue(zone.DisableARPNDSuppression.PointerBool()),
			ExitNodes:               stringset.NewValueString(zone.ExitNodes, diags, stringset.WithSeparator(",")),
			ExitNodesLocalRouting:   types.BoolPointerValue(zone.ExitNodesLocalRouting.PointerBool()),
			PrimaryExitNode:         types.StringPointerValue(zone.ExitNodesPrimary),
			RouteTargetImport:       types.StringPointerValue(zone.RouteTargetImport),
			VRFVXLANID:              types.Int64PointerValue(zone.VRFVXLANID),
		}

		objValue, objDiag := types.ObjectValueFrom(ctx, map[string]attr.Type{
			"id":          types.StringType,
			"type":        types.StringType,
			"ipam":        types.StringType,
			"dns":         types.StringType,
			"reverse_dns": types.StringType,
			"dns_zone":    types.StringType,
			"nodes": types.SetType{
				ElemType: types.StringType,
			},
			"mtu":                   types.Int64Type,
			"bridge":                types.StringType,
			"service_vlan":          types.Int64Type,
			"service_vlan_protocol": types.StringType,
			"peers": types.SetType{
				ElemType: types.StringType,
			},
			"advertise_subnets":          types.BoolType,
			"controller":                 types.StringType,
			"disable_arp_nd_suppression": types.BoolType,
			"exit_nodes": types.SetType{
				ElemType: types.StringType,
			},
			"exit_nodes_local_routing": types.BoolType,
			"primary_exit_node":        types.StringType,
			"rt_import":                types.StringType,
			"vrf_vxlan":                types.Int64Type,
		}, zoneData)
		resp.Diagnostics.Append(objDiag...)

		if resp.Diagnostics.HasError() {
			return
		}

		zoneElements[i] = objValue
	}

	listValue, listDiag := types.ListValue(types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":          types.StringType,
			"type":        types.StringType,
			"ipam":        types.StringType,
			"dns":         types.StringType,
			"reverse_dns": types.StringType,
			"dns_zone":    types.StringType,
			"nodes": types.SetType{
				ElemType: types.StringType,
			},
			"mtu":                   types.Int64Type,
			"bridge":                types.StringType,
			"service_vlan":          types.Int64Type,
			"service_vlan_protocol": types.StringType,
			"peers": types.SetType{
				ElemType: types.StringType,
			},
			"advertise_subnets":          types.BoolType,
			"controller":                 types.StringType,
			"disable_arp_nd_suppression": types.BoolType,
			"exit_nodes": types.SetType{
				ElemType: types.StringType,
			},
			"exit_nodes_local_routing": types.BoolType,
			"primary_exit_node":        types.StringType,
			"rt_import":                types.StringType,
			"vrf_vxlan":                types.Int64Type,
		},
	}, zoneElements)
	resp.Diagnostics.Append(listDiag...)

	if resp.Diagnostics.HasError() {
		return
	}

	data.Zones = listValue
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// NewZonesDataSource returns a new data source for SDN zones.
func NewZonesDataSource() datasource.DataSource {
	return &zonesDataSource{}
}
