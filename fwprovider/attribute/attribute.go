/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package attribute

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

// ResourceID generates an attribute definition suitable for the always-present resource `id` attribute.
func ResourceID(desc ...string) schema.StringAttribute {
	a := schema.StringAttribute{
		Computed:    true,
		Description: "The unique identifier of this resource.",
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	}

	if len(desc) > 0 {
		a.Description = desc[0]
	}

	return a
}

// ShouldBeRemoved evaluates if an attribute should be removed from the plan during update.
func ShouldBeRemoved(plan attr.Value, state attr.Value, isClone bool) bool {
	return !IsDefined(plan) && IsDefined(state) && !isClone
}

// IsDefined returns true if attribute is known and not null.
func IsDefined(v attr.Value) bool {
	return !v.IsNull() && !v.IsUnknown()
}
