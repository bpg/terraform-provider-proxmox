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

	"github.com/bpg/terraform-provider-proxmox/fwprovider/types/stringset"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/sdn/zones"
)

var (
	_ resource.ResourceWithConfigure   = &VXLANResource{}
	_ resource.ResourceWithImportState = &VXLANResource{}
)

type vxlanModel struct {
	genericModel

	Peers stringset.Value `tfsdk:"peers"`
}

func (m *vxlanModel) fromAPI(name string, data *zones.ZoneData, diags *diag.Diagnostics) {
	m.genericModel.fromAPI(name, data, diags)

	m.Peers = stringset.NewValueString(data.Peers, diags, stringset.WithSeparator(","))

	if data.Pending != nil {
		m.Peers = stringset.NewValueString(data.Pending.Peers, diags, stringset.WithSeparator(","))
	}
}

func (m *vxlanModel) toAPI(ctx context.Context, diags *diag.Diagnostics) *zones.Zone {
	data := m.genericModel.toAPI(ctx, diags)

	data.Peers = m.Peers.ValueStringPointer(ctx, diags, stringset.WithSeparator(","))

	return data
}

type VXLANResource struct {
	*genericZoneResource
}

func NewVXLANResource() resource.Resource {
	return &VXLANResource{
		genericZoneResource: newGenericZoneResource(zoneResourceConfig{
			typeNameSuffix: "_sdn_zone_vxlan",
			zoneType:       zones.TypeVXLAN,
			modelFunc:      func() zoneModel { return &vxlanModel{} },
		}).(*genericZoneResource),
	}
}

func (r *VXLANResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "VXLAN Zone in Proxmox SDN.",
		MarkdownDescription: "VXLAN Zone in Proxmox SDN. It establishes a tunnel (overlay) on top of an existing network " +
			"(underlay). This encapsulates layer 2 Ethernet frames within layer 4 UDP datagrams using the default " +
			"destination port 4789. You have to configure the underlay network yourself to enable UDP connectivity " +
			"between all peers. Because VXLAN encapsulation uses 50 bytes, the MTU needs to be 50 bytes lower than the " +
			"outgoing physical interface.",
		Attributes: genericAttributesWith(map[string]schema.Attribute{
			"peers": stringset.ResourceAttribute(
				"A list of IP addresses of each node in the VXLAN zone.",
				"A list of IP addresses of each node in the VXLAN zone. "+
					"This can be external nodes reachable at this IP address. All nodes in the cluster need to be "+
					"mentioned here",
				stringset.WithRequired(),
			),
		}),
	}
}

func (m *vxlanModel) getGenericModel() *genericModel {
	return &m.genericModel
}
