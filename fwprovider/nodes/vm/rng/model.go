/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package rng

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/attribute"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/vms"
)

// Model represents the RNG model.
type Model struct {
	Source   types.String `tfsdk:"source"`
	MaxBytes types.Int64  `tfsdk:"max_bytes"`
	Period   types.Int64  `tfsdk:"period"`
}

func attributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"source":    types.StringType,
		"max_bytes": types.Int64Type,
		"period":    types.Int64Type,
	}
}

// NullValue returns a properly typed null Value.
func NullValue() Value {
	return types.ObjectNull(attributeTypes())
}

// toAPI builds the PVE wire struct from the plan-side Model. Null/unknown framework fields map
// to zero/nil on the API struct; EncodeValues skips zero-valued fields when serializing the
// compound `rng0=...` property string.
//
// MaxBytes/Period carry zero through deliberately (F37 fix): the old code dropped user-set 0 via
// a `ValueInt64() != 0` guard, making the schema's documented "Use 0 to disable limiting" a lie.
// Now if the user writes 0, it reaches PVE.
func (m *Model) toAPI() *vms.CustomRNGDevice {
	dev := &vms.CustomRNGDevice{}

	if attribute.IsDefined(m.Source) {
		dev.Source = m.Source.ValueString()
	}

	if v := attribute.Int64PtrFromValue(m.MaxBytes); v != nil {
		n := int(*v)
		dev.MaxBytes = &n
	}

	if v := attribute.Int64PtrFromValue(m.Period); v != nil {
		n := int(*v)
		dev.Period = &n
	}

	return dev
}
