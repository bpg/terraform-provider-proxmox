/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package vga

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
