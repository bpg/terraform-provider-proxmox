package sdn

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/attribute"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/config"
	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/sdn/zones"
	"github.com/hashicorp/terraform-plugin-framework/attr"
)

var (
	_ resource.Resource                = &sdnZoneResource{}
	_ resource.ResourceWithConfigure   = &sdnZoneResource{}
	_ resource.ResourceWithImportState = &sdnZoneResource{}
)

type sdnZoneResource struct {
	client *zones.Client
}

func NewSDNZoneResource() resource.Resource {
	return &sdnZoneResource{}
}

func (r *sdnZoneResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_sdn_zone"
}

func (r *sdnZoneResource) Configure(
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

	r.client = cfg.Client.Cluster().SDNZones()
}

func (r *sdnZoneResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Manages SDN Zones in Proxmox VE.",
		Attributes: map[string]schema.Attribute{
			"id": attribute.ResourceID(),
			"name": schema.StringAttribute{
				Description: "The unique ID of the SDN zone.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				Description: "Zone type (e.g. simple, vlan, qinq, vxlan, evpn).",
				Required:    true,
			},
			"ipam": schema.StringAttribute{
				Optional:    true,
				Description: "IP Address Management system.",
			},
			"dns": schema.StringAttribute{
				Optional:    true,
				Description: "DNS server address.",
			},
			"reversedns": schema.StringAttribute{
				Optional:    true,
				Description: "Reverse DNS settings.",
			},
			"dns_zone": schema.StringAttribute{
				Optional:    true,
				Description: "DNS zone name.",
			},
			"nodes": schema.StringAttribute{
				Optional:    true,
				Description: "Comma-separated list of Proxmox node names.",
			},
			"mtu": schema.Int64Attribute{
				Optional:    true,
				Description: "MTU value for the zone.",
			},
			"bridge": schema.StringAttribute{
				Optional:    true,
				Description: "Bridge interface for VLAN/QinQ.",
			},
			"tag": schema.Int64Attribute{
				Optional:    true,
				Description: "Service VLAN tag for QinQ.",
			},
			"vlan_protocol": schema.StringAttribute{
				Optional:    true,
				Description: "Service VLAN protocol for QinQ.",
			},
			"peers": schema.StringAttribute{
				Optional:    true,
				Description: "Peers list for VXLAN.",
			},
			"controller": schema.StringAttribute{
				Optional:    true,
				Description: "EVPN controller address.",
			},
			"vrf_vxlan": schema.Int64Attribute{
				Optional:    true,
				Description: "EVPN VRF VXLAN ID.",
			},
			"exit_nodes": schema.StringAttribute{
				Optional:    true,
				Description: "Comma-separated list of exit nodes for EVPN.",
			},
			"primary_exit_node": schema.StringAttribute{
				Optional:    true,
				Description: "Primary exit node for EVPN.",
			},
			"exit_nodes_local_routing": schema.BoolAttribute{
				Optional:    true,
				Description: "Enable local routing for EVPN exit nodes.",
			},
			"advertise_subnets": schema.BoolAttribute{
				Optional:    true,
				Description: "Enable subnet advertisement for EVPN.",
			},
			"disable_arp_nd_suppression": schema.BoolAttribute{
				Optional:    true,
				Description: "Disable ARP/ND suppression for EVPN.",
			},
			"rt_import": schema.StringAttribute{
				Optional:    true,
				Description: "Route target import for EVPN.",
			},
		},
	}
}

func (r *sdnZoneResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan sdnZoneModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqData := plan.toAPIRequestBody()
	err := r.client.CreateZone(ctx, reqData)
	if err != nil {
		resp.Diagnostics.AddError("Unable to Create SDN Zone", err.Error())
		return
	}

	plan.ID = plan.Name
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *sdnZoneResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var state sdnZoneModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	zone, err := r.client.GetZone(ctx, state.ID.ValueString())
	if err != nil {
		if errors.Is(err, api.ErrResourceDoesNotExist) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError("Unable to Read SDN Zone", err.Error())
		return
	}

	readModel := &sdnZoneModel{}
	readModel.importFromAPI(zone.ID, zone)
	resp.Diagnostics.Append(resp.State.Set(ctx, readModel)...)
}

func (r *sdnZoneResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan sdnZoneModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqData := plan.toAPIRequestBody()
	err := r.client.UpdateZone(ctx, reqData)
	if err != nil {
		resp.Diagnostics.AddError("Unable to Update SDN Zone", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *sdnZoneResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var state sdnZoneModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteZone(ctx, state.ID.ValueString())
	if err != nil && !errors.Is(err, api.ErrResourceDoesNotExist) {
		resp.Diagnostics.AddError("Unable to Delete SDN Zone", err.Error())
	}
}

func (r *sdnZoneResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	zone, err := r.client.GetZone(ctx, req.ID)
	if err != nil {
		if errors.Is(err, api.ErrResourceDoesNotExist) {
			resp.Diagnostics.AddError("Zone does not exist", err.Error())
			return
		}

		resp.Diagnostics.AddError("Unable to Import SDN Zone", err.Error())
		return
	}

	readModel := &sdnZoneModel{}
	readModel.importFromAPI(zone.ID, zone)
	resp.Diagnostics.Append(resp.State.Set(ctx, readModel)...)
}

func (r *sdnZoneResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data sdnZoneModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Check the type field
	if data.Type.IsNull() || data.Type.IsUnknown() {
		return
	}

	required := map[string][]string{
		"vlan":  {"bridge"},
		"qinq":  {"bridge", "service_vlan"},
		"vxlan": {"peers"},
		"evpn":  {"controller", "vrf_vxlan"},
	}

	zoneType := data.Type.ValueString()

	// Extracts required fields and at the same time checks zone type validity
	fields, ok := required[zoneType]
	if !ok {
		return
	}

	// Map of field names to their values from data
	fieldMap := map[string]attr.Value{
		"bridge":       data.Bridge,
		"service_vlan": data.ServiceVLAN,
		"peers":        data.Peers,
		"controller":   data.Controller,
		"vrf_vxlan":    data.VRFVXLANID,
	}

	for _, field := range fields {
		val, exists := fieldMap[field]
		if !exists || val.IsNull() || val.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root(field),
				"Missing Required Field",
				fmt.Sprintf("Attribute %q is required when type is %q.", field, zoneType),
			)
		}
	}
}
