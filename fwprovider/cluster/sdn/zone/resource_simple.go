/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package zone

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/sdn/zones"
)

var (
	_ resource.ResourceWithConfigure   = &SimpleResource{}
	_ resource.ResourceWithImportState = &SimpleResource{}
)

type simpleModel struct {
	genericModel

	DHCP types.String `tfsdk:"dhcp"`
}

func (m *simpleModel) fromAPI(name string, data *zones.ZoneData, diags *diag.Diagnostics) {
	m.genericModel.fromAPI(name, data, diags)

	m.DHCP = types.StringPointerValue(data.DHCP)

	if data.Pending != nil {
		if data.Pending.DHCP != nil && *data.Pending.DHCP != "" {
			m.DHCP = types.StringValue(*data.Pending.DHCP)
		}
	}
}

func (m *simpleModel) toAPI(ctx context.Context, diags *diag.Diagnostics) *zones.Zone {
	data := m.genericModel.toAPI(ctx, diags)

	data.DHCP = m.DHCP.ValueStringPointer()

	return data
}

type SimpleResource struct {
	generic *genericZoneResource
}

func NewSimpleResource() resource.Resource {
	return &SimpleResource{
		generic: newGenericZoneResource(zoneResourceConfig{
			typeNameSuffix: "_sdn_zone_simple",
			zoneType:       zones.TypeSimple,
			modelFunc:      func() zoneModel { return &simpleModel{} },
		}).(*genericZoneResource),
	}
}

func (r *SimpleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Simple Zone in Proxmox SDN.",
		MarkdownDescription: "Simple Zone in Proxmox SDN. It will create an isolated VNet bridge. " +
			"This bridge is not linked to a physical interface, and VM traffic is only local on each the node. " +
			"It can be used in NAT or routed setups.",
		Attributes: genericAttributesWith(map[string]schema.Attribute{
			"dhcp": schema.StringAttribute{
				Optional: true,
				Description: "The type of the DHCP backend for this zone. " +
					"Currently supported values are `none` (default) and `dnsmasq`.",
			},
		}),
	}
}

func (r *SimpleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.generic.Metadata(ctx, req, resp)
}

func (r *SimpleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.generic.Configure(ctx, req, resp)
}

func (r *SimpleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	r.generic.Create(ctx, req, resp)
}

func (r *SimpleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	r.generic.Read(ctx, req, resp)
}

func (r *SimpleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	r.generic.Update(ctx, req, resp)
}

func (r *SimpleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	r.generic.Delete(ctx, req, resp)
}

func (r *SimpleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	r.generic.ImportState(ctx, req, resp)
}
