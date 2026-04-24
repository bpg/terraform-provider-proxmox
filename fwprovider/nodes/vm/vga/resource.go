/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package vga

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/attribute"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/vms"
)

// Value represents the type for VGA settings.
type Value = types.Object

// NewValue returns a new Value with the given VGA settings from the PVE API.
//
// Returns NullValue() when PVE has no vga device on the VM — audit Section 4 confirmed PVE
// returns the vga key absent unless the user set it. Returning a non-null Object with null
// subfields would produce a permanent plan-vs-state diff now that the schema is Optional only.
func NewValue(ctx context.Context, config *vms.GetResponseData, diags *diag.Diagnostics) Value {
	if config.VGADevice == nil {
		return NullValue()
	}

	vga := Model{
		Clipboard: types.StringPointerValue(config.VGADevice.Clipboard),
		Type:      types.StringPointerValue(config.VGADevice.Type),
		Memory:    types.Int64PointerValue(config.VGADevice.Memory),
	}

	obj, d := types.ObjectValueFrom(ctx, attributeTypes(), vga)
	diags.Append(d...)

	return obj
}

// FillCreateBody fills the CreateRequestBody with the VGA settings from the plan Value.
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

	body.VGADevice = plan.toAPI()
}

// FillUpdateBody fills the UpdateRequestBody with the VGA settings from the plan Value.
//
// The `vga` PVE API parameter is a compound property string — all subfields round-trip through a
// single key. Deletion semantics therefore follow the block, not the subfield: when the user
// removes the vga block from HCL, the provider emits `delete=vga`; when the block is present,
// the provider sends the whole CustomVGADevice and PVE replaces the key atomically. Subfields
// omitted from the plan are implicitly cleared on the wire (omitempty in EncodeValues).
func FillUpdateBody(
	ctx context.Context,
	planValue, stateValue Value,
	updateBody *vms.UpdateRequestBody,
	diags *diag.Diagnostics,
) {
	attribute.CheckDeleteBody(planValue, stateValue, updateBody, "vga")

	if planValue.IsNull() || planValue.IsUnknown() || planValue.Equal(stateValue) {
		return
	}

	var plan Model

	d := planValue.As(ctx, &plan, basetypes.ObjectAsOptions{})
	diags.Append(d...)

	if diags.HasError() {
		return
	}

	updateBody.VGADevice = plan.toAPI()
}
