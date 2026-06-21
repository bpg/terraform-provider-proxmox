/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package agent

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/attribute"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/vms"
	proxmoxtypes "github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

// Model represents the QEMU guest agent configuration.
//
// Mapping to Proxmox API (CustomAgent):
//   - Enabled → enabled
//   - Trim → fstrim_cloned_disks
//   - Type → type
type Model struct {
	Enabled types.Bool   `tfsdk:"enabled"`
	Trim    types.Bool   `tfsdk:"trim"`
	Type    types.String `tfsdk:"type"`
}

func attributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"enabled": types.BoolType,
		"trim":    types.BoolType,
		"type":    types.StringType,
	}
}

// NullValue returns a properly typed null Value.
func NullValue() Value {
	return types.ObjectNull(attributeTypes())
}

// toAPI builds the PVE wire struct from the plan-side Model.
func (m *Model) toAPI() *vms.CustomAgent {
	a := &vms.CustomAgent{}

	if attribute.IsDefined(m.Enabled) {
		v := proxmoxtypes.CustomBool(m.Enabled.ValueBool())
		a.Enabled = &v
	}

	if attribute.IsDefined(m.Trim) {
		v := proxmoxtypes.CustomBool(m.Trim.ValueBool())
		a.TrimClonedDisks = &v
	}

	a.Type = attribute.StringPtrFromValue(m.Type)

	return a
}
