/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package fabric

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/sdn/fabrics"
)

var (
	_ resource.ResourceWithConfigure   = &OSPFResource{}
	_ resource.ResourceWithImportState = &OSPFResource{}
)

type ospfModel struct {
	genericModel

	Area       types.String `tfsdk:"area"`
	IPv4Prefix types.String `tfsdk:"ip_prefix"`
}

func (m *ospfModel) fromAPI(name string, data *fabrics.FabricData, diags *diag.Diagnostics) {
	m.genericModel.fromAPI(name, data, diags)

	m.Area = types.StringPointerValue(data.Area)
	m.IPv4Prefix = types.StringPointerValue(data.IPv4Prefix)
}

func (m *ospfModel) toAPI(ctx context.Context, diags *diag.Diagnostics) *fabrics.Fabric {
	data := m.genericModel.toAPI(ctx, diags)

	data.Area = m.Area.ValueStringPointer()
	data.IPv4Prefix = m.IPv4Prefix.ValueStringPointer()

	return data
}

type OSPFResource struct {
	*genericFabricResource
}

func NewOSPFResource() resource.Resource {
	return &OSPFResource{
		genericFabricResource: newGenericFabricResource(fabricResourceConfig{
			typeNameSuffix: "_sdn_fabric_ospf",
			fabricProtocol: fabrics.ProtocolOSPF,
			modelFunc:      func() fabricModel { return &ospfModel{} },
		}).(*genericFabricResource),
	}
}

func (r *OSPFResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "OSPF Fabric in Proxmox SDN.",
		MarkdownDescription: "OSPF Fabric in Proxmox SDN.",
		Attributes: genericAttributesWith(map[string]schema.Attribute{
			"area": schema.StringAttribute{
				Description: "OSPF area. Either a IPv4 address or a 32-bit number. Gets validated in rust.",
				Required:    true,
			},
			"ip_prefix": schema.StringAttribute{
				Description: "IPv4 prefix cidr for the fabric.",
				Required:    true,
			},
		}),
	}
}

func (m *ospfModel) getGenericModel() *genericModel {
	return &m.genericModel
}
