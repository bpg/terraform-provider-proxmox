/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package vga

import (
	"context"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/attribute"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/vms"
)

// Value represents the type for VGA settings.
type Value = types.Object

// NewValue returns a new Value with the given VGA settings from the PVE API.
func NewValue(ctx context.Context, config *vms.GetResponseData, diags *diag.Diagnostics) Value {
	vga := Model{}

	if config.VGADevice != nil {
		vga.Clipboard = types.StringPointerValue(config.VGADevice.Clipboard)
		vga.Type = types.StringPointerValue(config.VGADevice.Type)
		vga.Memory = types.Int64PointerValue(config.VGADevice.Memory)
	}

	obj, d := types.ObjectValueFrom(ctx, attributeTypes(), vga)
	diags.Append(d...)

	return obj
}

// FillCreateBody fills the CreateRequestBody with the VGA settings from the Value.
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

	vgaDevice := &vms.CustomVGADevice{}

	// for computed fields, we need to check if they are unknown
	if !plan.Clipboard.IsUnknown() {
		vgaDevice.Clipboard = plan.Clipboard.ValueStringPointer()
	}

	if !plan.Type.IsUnknown() {
		vgaDevice.Type = plan.Type.ValueStringPointer()
	}

	if !plan.Memory.IsUnknown() {
		vgaDevice.Memory = plan.Memory.ValueInt64Pointer()
	}

	if !reflect.DeepEqual(vgaDevice, &vms.CustomVGADevice{}) {
		body.VGADevice = vgaDevice
	}
}

// FillUpdateBody fills the UpdateRequestBody with the VGA settings from the Value.
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

	vgaDevice := &vms.CustomVGADevice{
		Clipboard: state.Clipboard.ValueStringPointer(),
		Type:      state.Type.ValueStringPointer(),
		Memory:    state.Memory.ValueInt64Pointer(),
	}

	if !plan.Clipboard.Equal(state.Clipboard) {
		if attribute.ShouldBeRemoved(plan.Clipboard, state.Clipboard, isClone) {
			vgaDevice.Clipboard = nil
		} else if attribute.IsDefined(plan.Clipboard) {
			vgaDevice.Clipboard = plan.Clipboard.ValueStringPointer()
		}
	}

	if !plan.Type.Equal(state.Type) {
		if attribute.ShouldBeRemoved(plan.Type, state.Type, isClone) {
			vgaDevice.Type = nil
		} else if attribute.IsDefined(plan.Type) {
			vgaDevice.Type = plan.Type.ValueStringPointer()
		}
	}

	if !plan.Memory.Equal(state.Memory) {
		if attribute.ShouldBeRemoved(plan.Memory, state.Memory, isClone) {
			vgaDevice.Memory = nil
		} else if attribute.IsDefined(plan.Memory) {
			vgaDevice.Memory = plan.Memory.ValueInt64Pointer()
		}
	}

	if !reflect.DeepEqual(vgaDevice, &vms.CustomVGADevice{}) {
		updateBody.VGADevice = vgaDevice
	}
}
