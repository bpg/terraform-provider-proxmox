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

func (m *model) fromAPI(id string, data *vnets.VNetData) {
	m.ID = types.StringValue(id)

	m.Zone = m.handleDeletedValue(data.Zone)
	m.Alias = m.handleDeletedValue(data.Alias)
	m.IsolatePorts = types.BoolPointerValue(data.IsolatePorts.PointerBool())
	m.Tag = types.Int64PointerValue(data.Tag)
	m.VlanAware = types.BoolPointerValue(data.VlanAware.PointerBool())

	if data.Pending != nil {
		if data.Pending.Zone != nil {
			m.Zone = m.handleDeletedValue(data.Pending.Zone)
		}

		if data.Pending.Alias != nil {
			m.Alias = m.handleDeletedValue(data.Pending.Alias)
		}

		if data.Pending.IsolatePorts != nil {
			m.IsolatePorts = types.BoolPointerValue(data.Pending.IsolatePorts.PointerBool())
		}

		if data.Pending.Tag != nil {
			m.Tag = types.Int64Value(*data.Pending.Tag)
		}

		if data.Pending.VlanAware != nil {
			m.VlanAware = types.BoolPointerValue(data.Pending.VlanAware.PointerBool())
		}
	}
}

func (m *model) handleDeletedValue(value *string) types.String {
	if value == nil {
		return types.StringNull()
	}

	if *value == "deleted" {
		return types.StringNull()
	}

	return types.StringValue(*value)
}

func (m *model) toAPI() *vnets.VNet {
	data := &vnets.VNet{}

	data.Zone = m.Zone.ValueStringPointer()
	data.Alias = m.Alias.ValueStringPointer()
	data.IsolatePorts = proxmoxtypes.CustomBoolPtr(m.IsolatePorts.ValueBoolPointer())
	data.Tag = m.Tag.ValueInt64Pointer()
	data.VlanAware = proxmoxtypes.CustomBoolPtr(m.VlanAware.ValueBoolPointer())

	return data
}
