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
	_ resource.ResourceWithConfigure   = &VLANResource{}
	_ resource.ResourceWithImportState = &VLANResource{}
)

type vlanModel struct {
	genericModel

	Bridge types.String `tfsdk:"bridge"`
}

func (m *vlanModel) fromAPI(name string, data *zones.ZoneData, diags *diag.Diagnostics) {
	m.genericModel.fromAPI(name, data, diags)

	m.Bridge = types.StringPointerValue(data.Bridge)

	if data.Pending != nil {
		if data.Pending.Bridge != nil && *data.Pending.Bridge != "" {
			m.Bridge = types.StringValue(*data.Pending.Bridge)
		}
	}
}

func (m *vlanModel) toAPI(ctx context.Context, diags *diag.Diagnostics) *zones.Zone {
	data := m.genericModel.toAPI(ctx, diags)

	data.Bridge = m.Bridge.ValueStringPointer()

	return data
}

type VLANResource struct {
	generic *genericZoneResource
}

func NewVLANResource() resource.Resource {
	return &VLANResource{
		generic: newGenericZoneResource(zoneResourceConfig{
			typeNameSuffix: "_sdn_zone_vlan",
			zoneType:       zones.TypeVLAN,
			modelFunc:      func() zoneModel { return &vlanModel{} },
		}).(*genericZoneResource),
	}
}

func (r *VLANResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "VLAN Zone in Proxmox SDN.",
		MarkdownDescription: "VLAN Zone in Proxmox SDN. It uses an existing local Linux or OVS bridge to connect to the " +
			"node's physical interface. It uses VLAN tagging defined in the VNet to isolate the network segments. " +
			"This allows connectivity of VMs between different nodes.",
		Attributes: genericAttributesWith(map[string]schema.Attribute{
			"bridge": schema.StringAttribute{
				Description: "Bridge interface for VLAN.",
				MarkdownDescription: "The local bridge or OVS switch, already configured on _each_ node that allows " +
					"node-to-node connection.",
				Required: true,
			},
		}),
	}
}

func (r *VLANResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.generic.Metadata(ctx, req, resp)
}

func (r *VLANResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.generic.Configure(ctx, req, resp)
}

func (r *VLANResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	r.generic.Create(ctx, req, resp)
}

func (r *VLANResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	r.generic.Read(ctx, req, resp)
}

func (r *VLANResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	r.generic.Update(ctx, req, resp)
}

func (r *VLANResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	r.generic.Delete(ctx, req, resp)
}

func (r *VLANResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	r.generic.ImportState(ctx, req, resp)
}
