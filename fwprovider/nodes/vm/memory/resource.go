/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package memory

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/attribute"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/vms"
)

// Value represents the type for memory settings.
type Value = types.Object

// NewValue returns a new Value with the given memory settings from the PVE API.
//
// Returns NullValue() when the API has no memory-related keys set. Unlike rng/vga, memory maps
// to five independent top-level PVE keys rather than a compound property, so we inspect each
// pointer to decide whether the block is wholly absent.
func NewValue(ctx context.Context, config *vms.GetResponseData, diags *diag.Diagnostics) Value {
	if config.DedicatedMemory == nil &&
		config.FloatingMemory == nil &&
		config.FloatingMemoryShares == nil &&
		config.Hugepages == nil &&
		config.KeepHugepages == nil {
		return NullValue()
	}

	mem := Model{
		Size:          types.Int64PointerValue(config.DedicatedMemory.PointerInt64()),
		Balloon:       types.Int64PointerValue(config.FloatingMemory.PointerInt64()),
		Shares:        types.Int64PointerValue(asInt64Ptr(config.FloatingMemoryShares)),
		Hugepages:     types.StringPointerValue(config.Hugepages),
		KeepHugepages: types.BoolPointerValue(config.KeepHugepages.PointerBool()),
	}

	obj, d := types.ObjectValueFrom(ctx, attributeTypes(), mem)
	diags.Append(d...)

	return obj
}

// FillCreateBody fills the CreateRequestBody with the memory settings from the plan Value.
func FillCreateBody(ctx context.Context, planValue Value, body *vms.CreateRequestBody, diags *diag.Diagnostics) {
	if planValue.IsNull() || planValue.IsUnknown() {
		return
	}

	var plan Model

	d := planValue.As(ctx, &plan, basetypes.ObjectAsOptions{})
	diags.Append(d...)

	if d.HasError() {
		return
	}

	plan.toAPI(body)
}

// FillUpdateBody fills the UpdateRequestBody with the memory settings diff from state → plan.
//
// Memory exposes five independent PVE API keys (`memory`, `balloon`, `shares`, `hugepages`,
// `keephugepages`), so deletion and population are per-field: when a previously-set field is
// removed from the plan we emit `delete=<key>`, and when it is still set we (re)send its value.
// Fields absent from both state and plan stay off the wire entirely.
func FillUpdateBody(
	ctx context.Context,
	planValue, stateValue Value,
	updateBody *vms.UpdateRequestBody,
	diags *diag.Diagnostics,
) {
	plan := unpackOrEmpty(ctx, planValue, diags)
	state := unpackOrEmpty(ctx, stateValue, diags)

	if diags.HasError() {
		return
	}

	attribute.CheckDeleteBody(plan.Size, state.Size, updateBody, "memory")
	attribute.CheckDeleteBody(plan.Balloon, state.Balloon, updateBody, "balloon")
	attribute.CheckDeleteBody(plan.Shares, state.Shares, updateBody, "shares")
	attribute.CheckDeleteBody(plan.Hugepages, state.Hugepages, updateBody, "hugepages")
	attribute.CheckDeleteBody(plan.KeepHugepages, state.KeepHugepages, updateBody, "keephugepages")

	if planValue.IsNull() || planValue.IsUnknown() || planValue.Equal(stateValue) {
		return
	}

	plan.toAPI(updateBody)
}

// unpackOrEmpty returns a Model decoded from the Object, or a Model with all-null fields when
// the Object itself is null or unknown. Works around the lack of a natural "empty" Model for
// the CheckDeleteBody per-field diff when the whole block is absent on one side.
func unpackOrEmpty(ctx context.Context, value Value, diags *diag.Diagnostics) Model {
	if value.IsNull() || value.IsUnknown() {
		return Model{
			Size:          types.Int64Null(),
			Balloon:       types.Int64Null(),
			Shares:        types.Int64Null(),
			Hugepages:     types.StringNull(),
			KeepHugepages: types.BoolNull(),
		}
	}

	var m Model

	diags.Append(value.As(ctx, &m, basetypes.ObjectAsOptions{})...)

	return m
}

// asInt64Ptr converts a *int from the PVE API response into a *int64 for framework conversion.
func asInt64Ptr(v *int) *int64 {
	if v == nil {
		return nil
	}

	n := int64(*v)

	return &n
}
