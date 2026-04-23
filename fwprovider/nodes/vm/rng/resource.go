/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package rng

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/attribute"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/vms"
)

// Value represents the type for RNG settings.
type Value = types.Object

// NewValue returns a new Value with the given RNG settings from the PVE API.
//
// Returns NullValue() when PVE has no rng device configured — audit Section 4 confirmed the
// `rng0` key is absent from the GET response unless the user set it. Returning a non-null
// Object with null inner fields would produce a permanent plan-vs-state diff now that the
// schema is Optional only.
func NewValue(ctx context.Context, config *vms.GetResponseData, diags *diag.Diagnostics) Value {
	if config.RNGDevice == nil {
		return NullValue()
	}

	rng := Model{}

	// CustomRNGDevice.Source is string (not *string); treat "" as null so plans without
	// an explicit source don't drift after Read.
	if config.RNGDevice.Source != "" {
		rng.Source = types.StringValue(config.RNGDevice.Source)
	} else {
		rng.Source = types.StringNull()
	}

	if config.RNGDevice.MaxBytes != nil {
		rng.MaxBytes = types.Int64Value(int64(*config.RNGDevice.MaxBytes))
	}

	if config.RNGDevice.Period != nil {
		rng.Period = types.Int64Value(int64(*config.RNGDevice.Period))
	}

	obj, d := types.ObjectValueFrom(ctx, attributeTypes(), rng)
	diags.Append(d...)

	return obj
}

// FillCreateBody fills the CreateRequestBody with the RNG settings from the plan Value.
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

	body.RNGDevice = plan.toAPI()
}

// FillUpdateBody fills the UpdateRequestBody with the RNG settings from the plan Value.
//
// The `rng0` PVE API parameter is a compound property string — all subfields round-trip through
// a single key, same pattern as `vga`. Deletion is block-level (emit `delete=rng0`); updates
// send the whole CustomRNGDevice and PVE replaces the key atomically. Subfields omitted from
// the plan are implicitly cleared on the wire via EncodeValues' zero-check.
func FillUpdateBody(
	ctx context.Context,
	planValue, stateValue Value,
	updateBody *vms.UpdateRequestBody,
	diags *diag.Diagnostics,
) {
	attribute.CheckDeleteBody(planValue, stateValue, updateBody, "rng0")

	if planValue.IsNull() || planValue.IsUnknown() || planValue.Equal(stateValue) {
		return
	}

	var plan Model

	d := planValue.As(ctx, &plan, basetypes.ObjectAsOptions{})
	diags.Append(d...)

	if diags.HasError() {
		return
	}

	updateBody.RNGDevice = plan.toAPI()
}
