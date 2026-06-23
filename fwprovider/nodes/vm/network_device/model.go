/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package network_device

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/attribute"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/vms"
	proxmoxtypes "github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

// Value is the type alias for a list of network device objects.
type Value = types.List

// Model represents a single network device's Terraform model.
type Model struct {
	Bridge       types.String  `tfsdk:"bridge"`
	Disconnected types.Bool    `tfsdk:"disconnected"`
	Firewall     types.Bool    `tfsdk:"firewall"`
	MACAddress   types.String  `tfsdk:"mac_address"`
	ModelAttr    types.String  `tfsdk:"model"`
	MTU          types.Int64   `tfsdk:"mtu"`
	Queues       types.Int64   `tfsdk:"queues"`
	RateLimit    types.Float64 `tfsdk:"rate_limit"`
	Trunks       types.List    `tfsdk:"trunks"`
	VlanID       types.Int64   `tfsdk:"vlan_id"`
}

// attributeTypes returns the attribute types for a single network device object.
func attributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"bridge":       types.StringType,
		"disconnected": types.BoolType,
		"firewall":     types.BoolType,
		"mac_address":  types.StringType,
		"model":        types.StringType,
		"mtu":          types.Int64Type,
		"queues":       types.Int64Type,
		"rate_limit":   types.Float64Type,
		"trunks":       types.ListType{ElemType: types.Int64Type},
		"vlan_id":      types.Int64Type,
	}
}

// elementType returns the object type for a single network device element.
func elementType() attr.Type {
	return types.ObjectType{AttrTypes: attributeTypes()}
}

// NullValue returns a null network device list.
func NullValue() Value {
	return types.ListNull(elementType())
}

// toAPI converts the model to the PVE CustomNetworkDevice struct.
func (m *Model) toAPI() vms.CustomNetworkDevice {
	dev := vms.CustomNetworkDevice{
		Enabled: true,
		Model:   m.ModelAttr.ValueString(),
	}

	dev.Bridge = attribute.StringPtrFromValue(m.Bridge)
	dev.MACAddress = attribute.StringPtrFromValue(m.MACAddress)
	dev.RateLimit = attribute.Float64PtrFromValue(m.RateLimit)

	if attribute.IsDefined(m.Firewall) {
		v := proxmoxtypes.CustomBool(m.Firewall.ValueBool())
		dev.Firewall = &v
	}

	if attribute.IsDefined(m.Disconnected) {
		v := proxmoxtypes.CustomBool(m.Disconnected.ValueBool())
		dev.LinkDown = &v
	}

	if attribute.IsDefined(m.MTU) {
		v := int(m.MTU.ValueInt64())
		dev.MTU = &v
	}

	if attribute.IsDefined(m.Queues) {
		v := int(m.Queues.ValueInt64())
		dev.Queues = &v
	}

	if attribute.IsDefined(m.VlanID) {
		v := int(m.VlanID.ValueInt64())
		dev.Tag = &v
	}

	if !m.Trunks.IsNull() && !m.Trunks.IsUnknown() && len(m.Trunks.Elements()) > 0 {
		elems := m.Trunks.Elements()
		trunks := make([]int, len(elems))

		for i, e := range elems {
			if v, ok := e.(types.Int64); ok {
				trunks[i] = int(v.ValueInt64())
			}
		}

		dev.Trunks = trunks
	}

	return dev
}

// fromAPI converts a PVE CustomNetworkDevice to a Model.
func fromAPI(dev *vms.CustomNetworkDevice) Model {
	m := Model{
		Bridge:     types.StringPointerValue(dev.Bridge),
		MACAddress: types.StringPointerValue(dev.MACAddress),
		ModelAttr:  types.StringValue(dev.Model),
		RateLimit:  types.Float64PointerValue(dev.RateLimit),
	}

	if dev.Firewall != nil {
		m.Firewall = types.BoolValue(bool(*dev.Firewall))
	}

	if dev.LinkDown != nil {
		m.Disconnected = types.BoolValue(bool(*dev.LinkDown))
	}

	if dev.MTU != nil {
		m.MTU = types.Int64Value(int64(*dev.MTU))
	}

	if dev.Queues != nil {
		m.Queues = types.Int64Value(int64(*dev.Queues))
	}

	if dev.Tag != nil {
		m.VlanID = types.Int64Value(int64(*dev.Tag))
	}

	if len(dev.Trunks) > 0 {
		elems := make([]attr.Value, len(dev.Trunks))
		for i, t := range dev.Trunks {
			elems[i] = types.Int64Value(int64(t))
		}

		m.Trunks = types.ListValueMust(types.Int64Type, elems)
	} else {
		m.Trunks = types.ListNull(types.Int64Type)
	}

	return m
}
