/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package attribute

import (
	"github.com/bpg/terraform-provider-proxmox/fwprovider/types/stringset"
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

// ShouldBeRemoved evaluates if an attribute should be removed during update.
func ShouldBeRemoved(plan attr.Value, state attr.Value) bool {
	return !IsDefined(plan) && IsDefined(state)
}

// IsDefined returns true if attribute is known and not null.
func IsDefined(v attr.Value) bool {
	return !v.IsNull() && !v.IsUnknown()
}

// CheckDelete adds an API field name to the delete list if the plan field is null but the state field is not null.
// This is used to handle attribute deletion in API calls.
func CheckDelete(planField, stateField attr.Value, toDelete *[]string, apiName string) {
	planIsEmpty := planField.IsNull()
	stateIsEmpty := stateField.IsNull()

	// Special handling for stringset.Value: treat empty set as null
	if planSet, ok := planField.(stringset.Value); ok {
		planIsEmpty = planIsEmpty || len(planSet.Elements()) == 0
	}

	if planIsEmpty && !stateIsEmpty {
		*toDelete = append(*toDelete, apiName)
	}
}
