/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package tffwk

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// IDAttribute generates an attribute definition suitable for the always-present `id` attribute.
func IDAttribute() schema.StringAttribute {
	return schema.StringAttribute{
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	}
}

// StringOrNull returns either a string pointer or a `nil`, depending on whether the input attribute value
// was null/an empty string or some actual string.
func StringOrNull(inval types.String) *string {
	value := inval.ValueString()
	if value == "" {
		return nil
	}

	return &value
}
