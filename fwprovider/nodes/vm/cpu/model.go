/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package cpu

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Model represents the CPU model.
type Model struct {
	Affinity     types.String `tfsdk:"affinity"`
	Architecture types.String `tfsdk:"architecture"`
	Cores        types.Int64  `tfsdk:"cores"`
	Flags        types.Set    `tfsdk:"flags"`
	Hotplugged   types.Int64  `tfsdk:"hotplugged"`
	Limit        types.Int64  `tfsdk:"limit"`
	Numa         types.Bool   `tfsdk:"numa"`
	Sockets      types.Int64  `tfsdk:"sockets"`
	Type         types.String `tfsdk:"type"`
	Units        types.Int64  `tfsdk:"units"`
}

func attributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"affinity":     types.StringType,
		"architecture": types.StringType,
		"cores":        types.Int64Type,
		"flags":        types.SetType{ElemType: types.StringType},
		"hotplugged":   types.Int64Type,
		"limit":        types.Int64Type,
		"numa":         types.BoolType,
		"sockets":      types.Int64Type,
		"type":         types.StringType,
		"units":        types.Int64Type,
	}
}
