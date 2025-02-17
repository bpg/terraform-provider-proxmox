/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package rng

import (
	"context"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/attribute"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/vms"
)

// Value represents the type for RNG settings.
type Value = types.Object

// NewValue returns a new Value with the given RNG settings from the PVE API.
func NewValue(ctx context.Context, config *vms.GetResponseData, diags *diag.Diagnostics) Value {
	rng := Model{}

	if config.RNGDevice != nil {
		rng.Source = types.StringValue(config.RNGDevice.Source)

		if config.RNGDevice.MaxBytes != nil {
			rng.MaxBytes = types.Int64Value(int64(*config.RNGDevice.MaxBytes))
		}

		if config.RNGDevice.Period != nil {
			rng.Period = types.Int64Value(int64(*config.RNGDevice.Period))
		}
	}

	obj, d := types.ObjectValueFrom(ctx, attributeTypes(), rng)
	diags.Append(d...)

	return obj
}

// FillCreateBody fills the CreateRequestBody with the RNG settings from the Value.
//
// In the 'create' context, v is the plan.
func FillCreateBody(ctx context.Context, planValue Value, body *vms.CreateRequestBody, diags *diag.Diagnostics) {
	var plan Model

	if planValue.IsNull() || planValue.IsUnknown() {
		return
	}

	d := planValue.As(ctx, &plan, basetypes.ObjectAsOptions{})
	diags.Append(d...)

	if d.HasError() {
		return
	}

	rngDevice := &vms.CustomRNGDevice{}

	if !plan.Source.IsUnknown() {
		rngDevice.Source = plan.Source.ValueString()
	}

	if !plan.MaxBytes.IsUnknown() {
		maxBytes := int(plan.MaxBytes.ValueInt64())
		rngDevice.MaxBytes = &maxBytes
	}

	if !plan.Period.IsUnknown() {
		period := int(plan.Period.ValueInt64())
		rngDevice.Period = &period
	}

	if !reflect.DeepEqual(rngDevice, &vms.CustomRNGDevice{}) {
		body.RNGDevice = rngDevice
	}
}

// FillUpdateBody fills the UpdateRequestBody with the RNG settings from the Value.
//
// In the 'update' context, v is the plan and stateValue is the current state.
func FillUpdateBody(
	ctx context.Context,
	planValue, stateValue Value,
	updateBody *vms.UpdateRequestBody,
	isClone bool,
	diags *diag.Diagnostics,
) {
	var plan, state Model

	if planValue.IsNull() || planValue.IsUnknown() || planValue.Equal(stateValue) {
		return
	}

	d := planValue.As(ctx, &plan, basetypes.ObjectAsOptions{})
	diags.Append(d...)
	d = stateValue.As(ctx, &state, basetypes.ObjectAsOptions{})
	diags.Append(d...)

	if diags.HasError() {
		return
	}

	rngDevice := &vms.CustomRNGDevice{
		Source: state.Source.ValueString(),
	}

	if state.MaxBytes.ValueInt64() != 0 {
		maxBytes := int(state.MaxBytes.ValueInt64())
		rngDevice.MaxBytes = &maxBytes
	}

	if state.Period.ValueInt64() != 0 {
		period := int(state.Period.ValueInt64())
		rngDevice.Period = &period
	}

	if !plan.Source.Equal(state.Source) {
		if attribute.ShouldBeRemoved(plan.Source, state.Source, isClone) {
			rngDevice.Source = ""
		} else if attribute.IsDefined(plan.Source) {
			rngDevice.Source = plan.Source.ValueString()
		}
	}

	if !plan.MaxBytes.Equal(state.MaxBytes) {
		if attribute.ShouldBeRemoved(plan.MaxBytes, state.MaxBytes, isClone) {
			rngDevice.MaxBytes = nil
		} else if attribute.IsDefined(plan.MaxBytes) {
			maxBytes := int(plan.MaxBytes.ValueInt64())
			rngDevice.MaxBytes = &maxBytes
		}
	}

	if !plan.Period.Equal(state.Period) {
		if attribute.ShouldBeRemoved(plan.Period, state.Period, isClone) {
			rngDevice.Period = nil
		} else if attribute.IsDefined(plan.Period) {
			period := int(plan.Period.ValueInt64())
			rngDevice.Period = &period
		}
	}

	if !reflect.DeepEqual(rngDevice, &vms.CustomRNGDevice{}) {
		updateBody.RNGDevice = rngDevice
	}
}
