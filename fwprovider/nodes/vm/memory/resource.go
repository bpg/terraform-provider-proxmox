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

	// Map Proxmox API fields to Terraform schema:
	// API 'memory' (DedicatedMemory) → our 'size'
	// API 'balloon' (FloatingMemory) → our 'balloon'
	// API 'shares' (FloatingMemoryShares) → our 'shares'
	// API 'hugepages' → our 'hugepages'
	// API 'keephugepages' → our 'keep_hugepages'

	// Size (Proxmox API: 'memory')
	if config.DedicatedMemory != nil {
		mem.Size = types.Int64Value(int64(*config.DedicatedMemory))
	} else {
		// Default to 512 MiB if not specified
		mem.Size = types.Int64Value(512)
	}

	// Balloon (Proxmox API: 'balloon')
	if config.FloatingMemory != nil {
		mem.Balloon = types.Int64Value(int64(*config.FloatingMemory))
	} else {
		// Default to 0 (balloon disabled) if not specified
		mem.Balloon = types.Int64Value(0)
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
// This function maps Terraform schema fields to Proxmox API fields for API calls.
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

	// Map Terraform schema fields to Proxmox API fields:
	// our 'size' → API 'memory' (DedicatedMemory)
	// our 'balloon' → API 'balloon' (FloatingMemory)
	// our 'shares' → API 'shares' (FloatingMemoryShares)

	// Size → Proxmox 'memory' parameter
	if !plan.Size.IsUnknown() && !plan.Size.IsNull() {
		size := int(plan.Size.ValueInt64())
		body.DedicatedMemory = &size
	}

	// Balloon → Proxmox 'balloon' parameter
	if !plan.Balloon.IsUnknown() && !plan.Balloon.IsNull() {
		balloon := int(plan.Balloon.ValueInt64())
		body.FloatingMemory = &balloon
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
