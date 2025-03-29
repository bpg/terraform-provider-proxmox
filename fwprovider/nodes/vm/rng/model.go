/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package rng

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
