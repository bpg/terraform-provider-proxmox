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
	"regexp"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/config"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/types/stringset"
	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/sdn/zones"
	"github.com/bpg/terraform-provider-proxmox/proxmox/helpers/ptr"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
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
}

func (m *genericModel) importFromAPI(name string, data *zones.ZoneData, diags *diag.Diagnostics) {
	m.ID = types.StringValue(name)

	m.DNS = types.StringPointerValue(data.DNS)
	m.DNSZone = types.StringPointerValue(data.DNSZone)
	m.IPAM = types.StringPointerValue(data.IPAM)
	m.MTU = types.Int64PointerValue(data.MTU)
	m.Nodes = stringset.NewValueString(data.Nodes, diags, stringset.WithSeparator(","))
	m.ReverseDNS = types.StringPointerValue(data.ReverseDNS)
}

func (m *genericModel) toAPIRequestBody(ctx context.Context, diags *diag.Diagnostics) *zones.ZoneRequestData {
	data := &zones.ZoneRequestData{}

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
			Validators: []validator.String{
				// https://github.com/proxmox/pve-network/blob/faaf96a8378a3e41065018562c09c3de0aa434f5/src/PVE/Network/SDN/Zones/Plugin.pm#L34
				stringvalidator.RegexMatches(
					regexp.MustCompile(`^[A-Za-z][A-Za-z0-9]*[A-Za-z0-9]$`),
					"must be a valid zone identifier",
				),
				stringvalidator.LengthAtMost(8),
			},
		},
		"ipam": schema.StringAttribute{
			Optional:    true,
			Description: "IP Address Management system.",
		},
		"mtu": schema.Int64Attribute{
			Optional:    true,
			Description: "MTU value for the zone.",
		},
		"nodes": stringset.ResourceAttribute("Proxmox node names.", ""),
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
	importFromAPI(name string, data *zones.ZoneData, diags *diag.Diagnostics)
	toAPIRequestBody(ctx context.Context, diags *diag.Diagnostics) *zones.ZoneRequestData
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

	diags := &diag.Diagnostics{}
	reqData := plan.toAPIRequestBody(ctx, diags)
	resp.Diagnostics.Append(*diags...)

	reqData.Type = ptr.Ptr(r.config.zoneType)

	if err := r.client.CreateZone(ctx, reqData); err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create SDN Zone",
			err.Error(),
		)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *genericZoneResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := r.config.modelFunc()
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	zone, err := r.client.GetZone(ctx, state.getID())
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
	diags := &diag.Diagnostics{}
	readModel.importFromAPI(zone.ID, zone, diags)
	resp.Diagnostics.Append(*diags...)
	resp.Diagnostics.Append(resp.State.Set(ctx, readModel)...)
}

func (r *genericZoneResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	plan := r.config.modelFunc()
	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	diags := &diag.Diagnostics{}
	reqData := plan.toAPIRequestBody(ctx, diags)
	resp.Diagnostics.Append(*diags...)

	if err := r.client.UpdateZone(ctx, reqData); err != nil {
		resp.Diagnostics.AddError(
			"Unable to Update SDN Zone",
			err.Error(),
		)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
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
	diags := &diag.Diagnostics{}
	readModel.importFromAPI(zone.ID, zone, diags)
	resp.Diagnostics.Append(*diags...)
	resp.Diagnostics.Append(resp.State.Set(ctx, readModel)...)
}

// Schema is required to satisfy the resource.Resource interface. It should be implemented by the specific resource.
func (r *genericZoneResource) Schema(_ context.Context, _ resource.SchemaRequest, _ *resource.SchemaResponse) {
	// Intentionally left blank. Should be set by the specific resource.
}
