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

// Float64PtrFromValue returns a *float64 from a types.Float64, returning nil for null or unknown values.
func Float64PtrFromValue(v types.Float64) *float64 {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}

	return v.ValueFloat64Pointer()
}

// DeleteAppender is implemented by API request body types that accumulate PVE API parameter
// names scheduled for deletion. CheckDeleteBody uses it to record removals without leaking the
// underlying `Delete []string` slice into caller code.
//
// Implementations (e.g. *proxmox/nodes/vms.UpdateRequestBody) append the apiName to their
// wire-level `delete=...` parameter. The name is the PVE API parameter, not the Go field name.
type DeleteAppender interface {
	AppendDelete(apiName string)
}

// planRemovesField reports whether the transition from state to plan represents a user deletion:
// state had a value, plan is empty (null, or an empty collection for stringset/List/Map).
func planRemovesField(planField, stateField attr.Value) bool {
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

	return planIsEmpty && !stateIsEmpty
}

// CheckDelete adds an API field name to the delete list if the plan field is null but the state field is not null.
// This is used to handle attribute deletion in API calls.
func CheckDelete(planField, stateField attr.Value, toDelete *[]string, apiName string) {
	if planRemovesField(planField, stateField) {
		*toDelete = append(*toDelete, apiName)
	}
}

// CheckDeleteBody is the body-taking counterpart to CheckDelete documented in ADR-008
// §FillCreateBody and FillUpdateBody. It records the PVE API name on the body's own delete list
// via AppendDelete, keeping sub-package call-sites free of local `[]string` plumbing.
//
// Use from VM sub-packages (cpu, vga, rng, memory, ...) whose body type exposes AppendDelete.
// Non-VM Framework resources with plain `[]string` delete plumbing keep using CheckDelete.
func CheckDeleteBody[B DeleteAppender](planField, stateField attr.Value, body B, apiName string) {
	if planRemovesField(planField, stateField) {
		body.AppendDelete(apiName)
	}
}

// StringValueFromPtr returns a types.String from a *string, returning an empty string for nil.
func StringValueFromPtr(p *string) types.String {
	if p == nil {
		return types.StringValue("")
	}

	return types.StringValue(*p)
}

// Int64ValueFromPtr returns a types.Int64 from a *int64, returning 0 for nil.
func Int64ValueFromPtr(p *int64) types.Int64 {
	if p == nil {
		return types.Int64Value(0)
	}

	return types.Int64Value(*p)
}

// Float64ValueFromPtr returns a types.Float64 from a *float64, returning 0 for nil.
func Float64ValueFromPtr(p *float64) types.Float64 {
	if p == nil {
		return types.Float64Value(0)
	}

	return types.Float64Value(*p)
}

// BoolValueFromPtr returns a types.Bool from a *bool, returning false for nil.
func BoolValueFromPtr(p *bool) types.Bool {
	if p == nil {
		return types.BoolValue(false)
	}

	return types.BoolValue(*p)
}

// BoolValueFromCustomBoolPtr returns a types.Bool from a *proxmoxtypes.CustomBool, returning false for nil.
func BoolValueFromCustomBoolPtr(p *proxmoxtypes.CustomBool) types.Bool {
	if p == nil {
		return types.BoolValue(false)
	}

	return p.ToValue()
}
