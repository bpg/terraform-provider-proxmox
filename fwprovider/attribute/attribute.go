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
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/types/stringset"
	proxmoxtypes "github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

// ResourceID generates a Computed string attribute for server-assigned resource IDs.
// It includes UseStateForUnknown() so the ID is preserved across plan/apply cycles.
// Use this for resources where the server generates the ID (e.g., backup jobs, metrics servers).
// For user-provided IDs, define the attribute manually with Required + RequiresReplace + validators.
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

// StringPtrFromValue returns a *string from a types.String, returning nil for null or unknown values.
// Use this instead of ValueStringPointer() when the field is Optional+Computed without a Default,
// because ValueStringPointer() returns &"" for unknown values which sends empty strings to the API.
func StringPtrFromValue(v types.String) *string {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}

	return v.ValueStringPointer()
}

// CustomBoolPtrFromValue returns a *CustomBool from a types.Bool, returning nil for null or unknown values.
func CustomBoolPtrFromValue(v types.Bool) *proxmoxtypes.CustomBool {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}

	return proxmoxtypes.CustomBoolPtr(v.ValueBoolPointer())
}

// Int64PtrFromValue returns a *int64 from a types.Int64, returning nil for null or unknown values.
func Int64PtrFromValue(v types.Int64) *int64 {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}

	return v.ValueInt64Pointer()
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

	// Special handling for types.List: treat empty list as null
	if planList, ok := planField.(types.List); ok {
		planIsEmpty = planIsEmpty || len(planList.Elements()) == 0
	}

	// Special handling for types.Map: treat empty map as null
	if planMap, ok := planField.(types.Map); ok {
		planIsEmpty = planIsEmpty || len(planMap.Elements()) == 0
	}

	if planIsEmpty && !stateIsEmpty {
		*toDelete = append(*toDelete, apiName)
	}
}
