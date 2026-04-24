/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package vga

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/attribute"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/vms"
)

// Model represents the VGA model.
type Model struct {
	Clipboard types.String `tfsdk:"clipboard"`
	Type      types.String `tfsdk:"type"`
	Memory    types.Int64  `tfsdk:"memory"`
}

func attributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"clipboard": types.StringType,
		"type":      types.StringType,
		"memory":    types.Int64Type,
	}
}

// NullValue returns a properly typed null Value.
func NullValue() Value {
	return types.ObjectNull(attributeTypes())
}

// toAPI builds the PVE wire struct from the plan-side Model. The `vga` PVE parameter is a
// compound property, so subfields round-trip through a single key and omitempty in EncodeValues
// clears anything the plan leaves null.
func (m *Model) toAPI() *vms.CustomVGADevice {
	return &vms.CustomVGADevice{
		Clipboard: attribute.StringPtrFromValue(m.Clipboard),
		Type:      attribute.StringPtrFromValue(m.Type),
		Memory:    attribute.Int64PtrFromValue(m.Memory),
	}
}
