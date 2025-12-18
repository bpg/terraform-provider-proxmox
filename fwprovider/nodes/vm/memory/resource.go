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

	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/vms"
	proxmoxtypes "github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

// Value represents the type for memory settings.
type Value = types.Object

// NewValue returns a new Value with the given memory settings from the PVE API.
func NewValue(ctx context.Context, config *vms.GetResponseData, diags *diag.Diagnostics) Value {
	mem := Model{}

	// Map Proxmox API fields to our clearer naming:
	// API 'memory' (DedicatedMemory) → our 'maximum'
	// API 'balloon' (FloatingMemory) → our 'minimum'
	// API 'shares' (FloatingMemoryShares) → our 'shares'
	// API 'hugepages' → our 'hugepages'
	// API 'keephugepages' → our 'keep_hugepages'

	// Maximum memory (Proxmox API: 'memory')
	if config.DedicatedMemory != nil {
		mem.Maximum = types.Int64Value(int64(*config.DedicatedMemory))
	} else {
		// Default to 512 MiB if not specified
		mem.Maximum = types.Int64Value(512)
	}

	// Minimum memory (Proxmox API: 'balloon')
	if config.FloatingMemory != nil {
		mem.Minimum = types.Int64Value(int64(*config.FloatingMemory))
	} else {
		// Default to 0 (balloon disabled) if not specified
		mem.Minimum = types.Int64Value(0)
	}

	// Shares (CPU scheduler priority)
	if config.FloatingMemoryShares != nil {
		mem.Shares = types.Int64Value(int64(*config.FloatingMemoryShares))
	} else {
		// Default to 1000 if not specified
		mem.Shares = types.Int64Value(1000)
	}

	// Hugepages
	mem.Hugepages = types.StringPointerValue(config.Hugepages)

	// Keep hugepages
	mem.KeepHugepages = types.BoolPointerValue(config.KeepHugepages.PointerBool())

	obj, d := types.ObjectValueFrom(ctx, attributeTypes(), mem)
	diags.Append(d...)

	return obj
}

// FillUpdateBody fills the UpdateRequestBody with the memory settings from the Value.
//
// This function converts our clearer naming convention (maximum/minimum) back to
// the Proxmox API naming (memory/balloon) for API calls.
//
// In the 'update' context, planValue is the plan (desired state).
func FillUpdateBody(ctx context.Context, planValue Value, body *vms.UpdateRequestBody, diags *diag.Diagnostics) {
	var plan Model

	if planValue.IsNull() || planValue.IsUnknown() {
		return
	}

	d := planValue.As(ctx, &plan, basetypes.ObjectAsOptions{})
	diags.Append(d...)

	if d.HasError() {
		return
	}

	// Map our clearer naming back to Proxmox API fields:
	// our 'maximum' → API 'memory' (DedicatedMemory)
	// our 'minimum' → API 'balloon' (FloatingMemory)
	// our 'shares' → API 'shares' (FloatingMemoryShares)

	// Maximum memory → Proxmox 'memory' parameter
	if !plan.Maximum.IsUnknown() && !plan.Maximum.IsNull() {
		maximum := int(plan.Maximum.ValueInt64())
		body.DedicatedMemory = &maximum
	}

	// Minimum memory → Proxmox 'balloon' parameter
	if !plan.Minimum.IsUnknown() && !plan.Minimum.IsNull() {
		minimum := int(plan.Minimum.ValueInt64())
		body.FloatingMemory = &minimum
	}

	// Shares → Proxmox 'shares' parameter
	if !plan.Shares.IsUnknown() && !plan.Shares.IsNull() {
		shares := int(plan.Shares.ValueInt64())
		body.FloatingMemoryShares = &shares
	}

	// Hugepages
	if !plan.Hugepages.IsUnknown() && !plan.Hugepages.IsNull() {
		body.Hugepages = plan.Hugepages.ValueStringPointer()
	}

	// Keep hugepages
	if !plan.KeepHugepages.IsUnknown() && !plan.KeepHugepages.IsNull() {
		keepHugepages := proxmoxtypes.CustomBool(plan.KeepHugepages.ValueBool())
		body.KeepHugepages = &keepHugepages
	}
}
