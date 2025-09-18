/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package vnet

import (
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/sdn/vnets"
	"github.com/hashicorp/terraform-plugin-framework/types"

	proxmoxtypes "github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

type model struct {
	ID           types.String `tfsdk:"id"`
	Zone         types.String `tfsdk:"zone"`
	Alias        types.String `tfsdk:"alias"`
	IsolatePorts types.Bool   `tfsdk:"isolate_ports"`
	Tag          types.Int64  `tfsdk:"tag"`
	VlanAware    types.Bool   `tfsdk:"vlan_aware"`
}

func (m *model) fromAPI(id string, data *vnets.Vnet) {
	m.ID = types.StringValue(id)

	m.Zone = types.StringPointerValue(data.Zone)
	m.Alias = types.StringPointerValue(data.Alias)

	m.IsolatePorts = types.BoolPointerValue(data.IsolatePorts.PointerBool())
	m.Tag = types.Int64PointerValue(data.Tag)
	m.VlanAware = types.BoolPointerValue(data.VlanAware.PointerBool())
}

func (m *model) toAPI() *vnets.Vnet {
	data := &vnets.Vnet{}

	data.Zone = m.Zone.ValueStringPointer()
	data.Alias = m.Alias.ValueStringPointer()
	data.IsolatePorts = proxmoxtypes.CustomBoolPtr(m.IsolatePorts.ValueBoolPointer())
	data.Tag = m.Tag.ValueInt64Pointer()
	data.VlanAware = proxmoxtypes.CustomBoolPtr(m.VlanAware.ValueBoolPointer())

	return data
}
