/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package agent

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/attribute"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/vms"
)

// Value represents the type for agent settings.
type Value = types.Object

// NewValue returns a new Value with the given agent settings from the PVE API.
//
// Returns NullValue() when no agent is configured — the `agent` key is absent from the GET
// response unless the user set it.
func NewValue(ctx context.Context, config *vms.GetResponseData, diags *diag.Diagnostics) Value {
	if config.Agent == nil {
		return NullValue()
	}

	m := Model{
		Enabled: types.BoolPointerValue(config.Agent.Enabled.PointerBool()),
		Trim:    types.BoolPointerValue(config.Agent.TrimClonedDisks.PointerBool()),
		Type:    types.StringPointerValue(config.Agent.Type),
	}

	obj, d := types.ObjectValueFrom(ctx, attributeTypes(), m)
	diags.Append(d...)

	return obj
}

// FillCreateBody fills the CreateRequestBody with the agent settings from the plan Value.
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

	body.Agent = plan.toAPI()
}

// FillUpdateBody fills the UpdateRequestBody with the agent settings diff from state → plan.
//
// The `agent` PVE API parameter is a compound property string — all subfields round-trip through
// a single key. Deletion is block-level (emit `delete=agent`); updates send the whole CustomAgent.
func FillUpdateBody(
	ctx context.Context,
	planValue, stateValue Value,
	updateBody *vms.UpdateRequestBody,
	diags *diag.Diagnostics,
) {
	attribute.CheckDeleteBody(planValue, stateValue, updateBody, "agent")

	if planValue.IsNull() || planValue.IsUnknown() || planValue.Equal(stateValue) {
		return
	}

	var plan Model

	d := planValue.As(ctx, &plan, basetypes.ObjectAsOptions{})
	diags.Append(d...)

	if diags.HasError() {
		return
	}

	updateBody.Agent = plan.toAPI()
}
