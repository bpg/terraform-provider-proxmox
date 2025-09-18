/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package zone

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/sdn/zones"
)

var (
	_ resource.ResourceWithConfigure   = &QinQResource{}
	_ resource.ResourceWithImportState = &QinQResource{}
)

type qinqModel struct {
	vlanModel

	ServiceVLAN         types.Int64  `tfsdk:"service_vlan"`
	ServiceVLANProtocol types.String `tfsdk:"service_vlan_protocol"`
}

func (m *qinqModel) fromAPI(name string, data *zones.ZoneData, diags *diag.Diagnostics) {
	m.vlanModel.fromAPI(name, data, diags)

	m.ServiceVLAN = types.Int64PointerValue(data.ServiceVLAN)
	m.ServiceVLANProtocol = types.StringPointerValue(data.ServiceVLANProtocol)

	if data.Pending != nil {
		m.ServiceVLAN = types.Int64PointerValue(data.Pending.ServiceVLAN)
		m.ServiceVLANProtocol = types.StringPointerValue(data.Pending.ServiceVLANProtocol)
	}
}

func (m *qinqModel) toAPI(ctx context.Context, diags *diag.Diagnostics) *zones.Zone {
	data := m.vlanModel.toAPI(ctx, diags)

	data.ServiceVLAN = m.ServiceVLAN.ValueInt64Pointer()
	data.ServiceVLANProtocol = m.ServiceVLANProtocol.ValueStringPointer()

	return data
}

type QinQResource struct {
	generic *genericZoneResource
}

func NewQinQResource() resource.Resource {
	return &QinQResource{
		generic: newGenericZoneResource(zoneResourceConfig{
			typeNameSuffix: "_sdn_zone_qinq",
			zoneType:       zones.TypeQinQ,
			modelFunc:      func() zoneModel { return &qinqModel{} },
		}).(*genericZoneResource),
	}
}

func (r *QinQResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "QinQ Zone in Proxmox SDN.",
		MarkdownDescription: "QinQ Zone in Proxmox SDN. QinQ also known as VLAN stacking, that uses multiple layers of " +
			"VLAN tags for isolation. The QinQ zone defines the outer VLAN tag (the Service VLAN) whereas the inner " +
			"VLAN tag is defined by the VNet. Your physical network switches must support stacked VLANs for this " +
			"configuration. Due to the double stacking of tags, you need 4 more bytes for QinQ VLANs. " +
			"For example, you must reduce the MTU to 1496 if you physical interface MTU is 1500.",
		Attributes: genericAttributesWith(map[string]schema.Attribute{
			"bridge": schema.StringAttribute{
				Description: "A local, VLAN-aware bridge that is already configured on each local node",
				Required:    true,
			},
			"service_vlan": schema.Int64Attribute{
				Description:         "Service VLAN tag for QinQ.",
				MarkdownDescription: "Service VLAN tag for QinQ. The tag must be between `1` and `4094`.",
				Validators: []validator.Int64{
					int64validator.Between(int64(1), int64(4094)),
				},
				Required: true,
			},
			"service_vlan_protocol": schema.StringAttribute{
				Description:         "Service VLAN protocol for QinQ.",
				MarkdownDescription: "Service VLAN protocol for QinQ. The protocol must be `802.1ad` or `802.1q`.",
				Validators: []validator.String{
					stringvalidator.OneOf("802.1ad", "802.1q"),
				},
				Optional: true,
			},
		}),
	}
}

func (r *QinQResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.generic.Metadata(ctx, req, resp)
}

func (r *QinQResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.generic.Configure(ctx, req, resp)
}

func (r *QinQResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	r.generic.Create(ctx, req, resp)
}

func (r *QinQResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	r.generic.Read(ctx, req, resp)
}

func (r *QinQResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	r.generic.Update(ctx, req, resp)
}

func (r *QinQResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	r.generic.Delete(ctx, req, resp)
}

func (r *QinQResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	r.generic.ImportState(ctx, req, resp)
}
