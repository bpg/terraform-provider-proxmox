/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package cpu

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/attribute"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/vms"
)

// Value represents the type for CPU settings.
type Value = types.Object

// NewValue returns a new Value with the given CPU settings from the PVE API.
//
// Returns NullValue() when none of the cpu-related API keys are present (block wholly absent).
// Otherwise builds the Object via Model.fromAPI. Empirical testing showed PVE does not
// auto-populate cores/sockets on Read when other cpu.* fields are set (contrary to the
// originally-proposed Optional+Computed carve-out), so all 10 inner attributes stay Optional-only
// and a null inner value simply reflects "user didn't set it".
func NewValue(ctx context.Context, config *vms.GetResponseData, diags *diag.Diagnostics) Value {
	if config.CPUAffinity == nil &&
		config.CPUArchitecture == nil &&
		config.CPUCores == nil &&
		config.CPUSockets == nil &&
		config.CPULimit == nil &&
		config.CPUUnits == nil &&
		config.CPUEmulation == nil &&
		config.NUMAEnabled == nil &&
		config.VirtualCPUCount == nil {
		return NullValue()
	}

	var m Model

	m.fromAPI(ctx, config, diags)

	obj, d := types.ObjectValueFrom(ctx, attributeTypes(), m)
	diags.Append(d...)

	return obj
}

// FillCreateBody fills the CreateRequestBody with the CPU settings from the plan Value.
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

	plan.toAPI(ctx, body, diags)
}

// FillUpdateBody fills the UpdateRequestBody with the CPU settings diff from state → plan.
//
// Each scalar CPU field is an independent top-level PVE API key (`affinity`, `arch`, `cores`,
// `cpulimit`, `sockets`, `cpuunits`), so deletion and population are per-field: when a
// previously-set field is removed from the plan we emit `delete=<apiName>`, and when it is still
// set we (re)send its value. The compound `cpu` key (CPUEmulation) is atomic — partial deletion
// of just `flags` without `type` is rejected on the PVE side.
func FillUpdateBody(
	ctx context.Context,
	planValue, stateValue Value,
	updateBody *vms.UpdateRequestBody,
	diags *diag.Diagnostics,
) {
	// Skip when plan is unknown (inner-field unknowns would be mis-read as deletions by the
	// null-based diff below) or identical to state.
	if planValue.IsUnknown() || planValue.Equal(stateValue) {
		return
	}

	plan := unpackOrEmpty(ctx, planValue, diags)
	state := unpackOrEmpty(ctx, stateValue, diags)

	if diags.HasError() {
		return
	}

	// Per-field scalar deletions.
	attribute.CheckDeleteBody(plan.Affinity, state.Affinity, updateBody, "affinity")
	attribute.CheckDeleteBody(plan.Architecture, state.Architecture, updateBody, "arch")
	attribute.CheckDeleteBody(plan.Cores, state.Cores, updateBody, "cores")
	attribute.CheckDeleteBody(plan.Limit, state.Limit, updateBody, "cpulimit")
	attribute.CheckDeleteBody(plan.Numa, state.Numa, updateBody, "numa")
	attribute.CheckDeleteBody(plan.Sockets, state.Sockets, updateBody, "sockets")
	attribute.CheckDeleteBody(plan.Units, state.Units, updateBody, "cpuunits")
	attribute.CheckDeleteBody(plan.Vcpus, state.Vcpus, updateBody, "vcpus")

	// Compound CPUEmulation: delete when the user explicitly removes `type` (flags alone isn't
	// valid). PVE rejects `cpu=...` without a cputype, so the block is all-or-nothing on the
	// wire. Gate on IsNull rather than !IsDefined so a plan.Type that's unknown from a dynamic
	// reference doesn't get mis-read as a deletion.
	if !state.Type.IsNull() && plan.Type.IsNull() {
		if attribute.IsDefined(plan.Flags) {
			diags.AddError("Cannot have CPU flags without explicit definition of CPU type", "")

			return
		}

		updateBody.AppendDelete("cpu")
	}

	if planValue.IsNull() {
		return
	}

	plan.toAPI(ctx, updateBody, diags)
}

// unpackOrEmpty returns a Model decoded from the Object, or a Model with all-null fields when the
// Object itself is null or unknown. Gives CheckDeleteBody a stable per-field view regardless of
// whether the whole block is absent on one side.
func unpackOrEmpty(ctx context.Context, value Value, diags *diag.Diagnostics) Model {
	if value.IsNull() || value.IsUnknown() {
		return Model{
			Affinity:     types.StringNull(),
			Architecture: types.StringNull(),
			Cores:        types.Int64Null(),
			Flags:        types.SetNull(basetypes.StringType{}),
			Limit:        types.Float64Null(),
			Numa:         types.BoolNull(),
			Sockets:      types.Int64Null(),
			Type:         types.StringNull(),
			Units:        types.Int64Null(),
			Vcpus:        types.Int64Null(),
		}
	}

	var m Model

	diags.Append(value.As(ctx, &m, basetypes.ObjectAsOptions{})...)

	return m
}
