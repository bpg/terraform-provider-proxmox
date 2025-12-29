/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package zone

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
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
	*genericZoneResource
}

func NewSimpleResource() resource.Resource {
	return &SimpleResource{
		genericZoneResource: newGenericZoneResource(zoneResourceConfig{
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
					"Currently the only supported value is `dnsmasq`.",
				Validators: []validator.String{
					stringvalidator.OneOf("dnsmasq"),
				},
			},
		}),
	}
}

func (m *simpleModel) getGenericModel() *genericModel {
	return &m.genericModel
}
