/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package zone

import (
	"context"
	"errors"
	"fmt"
	"maps"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/config"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/types/stringset"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/validators"
	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/sdn"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/sdn/zones"
	"github.com/bpg/terraform-provider-proxmox/proxmox/helpers/ptr"
	proxmoxtypes "github.com/bpg/terraform-provider-proxmox/proxmox/types"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type genericModel struct {
	ID         types.String    `tfsdk:"id"`
	IPAM       types.String    `tfsdk:"ipam"`
	DNS        types.String    `tfsdk:"dns"`
	ReverseDNS types.String    `tfsdk:"reverse_dns"`
	DNSZone    types.String    `tfsdk:"dns_zone"`
	Nodes      stringset.Value `tfsdk:"nodes"`
	MTU        types.Int64     `tfsdk:"mtu"`
	Pending    types.Bool      `tfsdk:"pending"`
	State      types.String    `tfsdk:"state"`
}

func (m *genericModel) fromAPI(name string, data *zones.ZoneData, diags *diag.Diagnostics) {
	m.ID = types.StringValue(name)

	m.DNS = types.StringPointerValue(data.DNS)
	m.DNSZone = types.StringPointerValue(data.DNSZone)
	m.IPAM = types.StringPointerValue(data.IPAM)
	m.MTU = types.Int64PointerValue(data.MTU)
	m.Nodes = stringset.NewValueString(data.Nodes, diags, stringset.WithSeparator(","))
	m.ReverseDNS = types.StringPointerValue(data.ReverseDNS)
	m.State = types.StringPointerValue(data.State)

	if data.Pending != nil {
		m.Pending = types.BoolValue(true)

		if data.Pending.DNS != nil && *data.Pending.DNS != "" {
			m.DNS = types.StringValue(*data.Pending.DNS)
		}

		if data.Pending.DNSZone != nil && *data.Pending.DNSZone != "" {
			m.DNSZone = types.StringValue(*data.Pending.DNSZone)
		}

		if data.Pending.IPAM != nil && *data.Pending.IPAM != "" {
			m.IPAM = types.StringValue(*data.Pending.IPAM)
		}

		if data.Pending.MTU != nil && *data.Pending.MTU != 0 {
			m.MTU = types.Int64Value(*data.Pending.MTU)
		}

		if data.Pending.Nodes != nil && len(*data.Pending.Nodes) > 0 {
			m.Nodes = stringset.NewValueString(data.Pending.Nodes, diags, stringset.WithSeparator(","))
		}

		if data.Pending.ReverseDNS != nil && *data.Pending.ReverseDNS != "" {
			m.ReverseDNS = types.StringValue(*data.Pending.ReverseDNS)
		}

		if data.Pending.State != nil && *data.Pending.State != "" {
			m.State = types.StringValue(*data.Pending.State)
		}
	}
}

func (m *genericModel) toAPI(ctx context.Context, diags *diag.Diagnostics) *zones.Zone {
	data := &zones.Zone{}

	data.ID = m.ID.ValueString()

	data.IPAM = m.IPAM.ValueStringPointer()
	data.DNS = m.DNS.ValueStringPointer()
	data.ReverseDNS = m.ReverseDNS.ValueStringPointer()
	data.DNSZone = m.DNSZone.ValueStringPointer()
	data.Nodes = m.Nodes.ValueStringPointer(ctx, diags, stringset.WithSeparator(","))
	data.MTU = m.MTU.ValueInt64Pointer()

	return data
}

func (m *genericModel) getID() string {
	return m.ID.ValueString()
}

func genericAttributesWith(extraAttributes map[string]schema.Attribute) map[string]schema.Attribute {
	// Start with generic attributes as the base
	result := map[string]schema.Attribute{
		"dns": schema.StringAttribute{
			Optional:    true,
			Description: "DNS API server address.",
		},
		"dns_zone": schema.StringAttribute{
			Optional:    true,
			Description: "DNS domain name. The DNS zone must already exist on the DNS server.",
			MarkdownDescription: "DNS domain name. Used to register hostnames, such as `<hostname>.<domain>`. " +
				"The DNS zone must already exist on the DNS server.",
		},
		"id": schema.StringAttribute{
			Description: "The unique identifier of the SDN zone.",
			Required:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
			Validators: validators.SDNID(),
		},
		"ipam": schema.StringAttribute{
			Optional:    true,
			Description: "IP Address Management system.",
		},
		"mtu": schema.Int64Attribute{
			Optional:    true,
			Description: "MTU value for the zone.",
		},
		"nodes": stringset.ResourceAttribute("The Proxmox nodes which the zone and associated VNets should be deployed on", "", stringset.WithRequired()),
		"pending": schema.BoolAttribute{
			Computed:    true,
			Description: "Indicates if the zone has pending configuration changes that need to be applied.",
		},
		"state": schema.StringAttribute{
			Computed:    true,
			Description: "Indicates the current state of the zone.",
		},
		"reverse_dns": schema.StringAttribute{
			Optional:    true,
			Description: "Reverse DNS API server address.",
		},
	}

	// Add extra attributes, allowing them to override generic ones if needed
	if extraAttributes != nil {
		maps.Copy(result, extraAttributes)
	}

	return result
}

type zoneModel interface {
	fromAPI(name string, data *zones.ZoneData, diags *diag.Diagnostics)
	toAPI(ctx context.Context, diags *diag.Diagnostics) *zones.Zone
	getID() string
}

type zoneResourceConfig struct {
	typeNameSuffix string
	zoneType       string
	modelFunc      func() zoneModel
}

type genericZoneResource struct {
	client *zones.Client
	config zoneResourceConfig
}

func newGenericZoneResource(cfg zoneResourceConfig) resource.Resource {
	return &genericZoneResource{config: cfg}
}

func (r *genericZoneResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + r.config.typeNameSuffix
}

func (r *genericZoneResource) Configure(
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
			fmt.Sprintf(
				"Expected config.Resource, got: %T",
				req.ProviderData,
			),
		)

		return
	}

	r.client = cfg.Client.Cluster().SDNZones()
}

func (r *genericZoneResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := r.config.modelFunc()
	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	reqData := plan.toAPI(ctx, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	reqData.Type = ptr.Ptr(r.config.zoneType)

	if err := r.client.CreateZone(ctx, reqData); err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create SDN Zone",
			err.Error(),
		)

		return
	}

	zone, err := r.client.GetZoneWithParams(ctx, plan.getID(), &sdn.QueryParams{Pending: proxmoxtypes.CustomBool(true).Pointer()})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read SDN Zone",
			err.Error(),
		)

		return
	}

	readModel := r.config.modelFunc()
	readModel.fromAPI(zone.ID, zone, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, readModel)...)
}

func (r *genericZoneResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := r.config.modelFunc()
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	zone, err := r.client.GetZoneWithParams(ctx, state.getID(), &sdn.QueryParams{Pending: proxmoxtypes.CustomBool(true).Pointer()})
	if err != nil {
		if errors.Is(err, api.ErrResourceDoesNotExist) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Unable to Read SDN Zone",
			err.Error(),
		)

		return
	}

	readModel := r.config.modelFunc()
	readModel.fromAPI(zone.ID, zone, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, readModel)...)
}

func (r *genericZoneResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	plan := r.config.modelFunc()
	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	diags := &diag.Diagnostics{}
	reqData := plan.toAPI(ctx, diags)
	resp.Diagnostics.Append(*diags...)

	update := &zones.ZoneUpdate{
		Zone: *reqData,
	}

	if err := r.client.UpdateZone(ctx, update); err != nil {
		resp.Diagnostics.AddError(
			"Unable to Update SDN Zone",
			err.Error(),
		)

		return
	}

	zone, err := r.client.GetZoneWithParams(ctx, plan.getID(), &sdn.QueryParams{Pending: proxmoxtypes.CustomBool(true).Pointer()})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read SDN Zone After Update",
			err.Error(),
		)

		return
	}

	state := r.config.modelFunc()
	state.fromAPI(zone.ID, zone, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *genericZoneResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	state := r.config.modelFunc()
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteZone(ctx, state.getID()); err != nil &&
		!errors.Is(err, api.ErrResourceDoesNotExist) {
		resp.Diagnostics.AddError(
			"Unable to Delete SDN Zone",
			err.Error(),
		)
	}
}

func (r *genericZoneResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	zone, err := r.client.GetZone(ctx, req.ID)
	if err != nil {
		if errors.Is(err, api.ErrResourceDoesNotExist) {
			resp.Diagnostics.AddError(fmt.Sprintf("Zone %s does not exist", req.ID), err.Error())
			return
		}

		resp.Diagnostics.AddError(fmt.Sprintf("Unable to Import SDN Zone %s", req.ID), err.Error())

		return
	}

	readModel := r.config.modelFunc()
	readModel.fromAPI(zone.ID, zone, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, readModel)...)
}

// Schema is required to satisfy the resource.Resource interface. It should be implemented by the specific resource.
func (r *genericZoneResource) Schema(_ context.Context, _ resource.SchemaRequest, _ *resource.SchemaResponse) {
	// Intentionally left blank. Should be set by the specific resource.
}
