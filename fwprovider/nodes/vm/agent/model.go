/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL can be obtained
 * at http://mozilla.org/MPL/2.0/.
 */

package agent

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/attribute"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/vms"
	proxmoxtypes "github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

type Value = types.Object

type Model struct {
	Enabled  types.Bool   `tfsdk:"enabled"`
	Trim     types.Bool   `tfsdk:"trim"`
	Type     types.String `tfsdk:"type"`
	FreezeFS types.Bool   `tfsdk:"freeze_fs"`
}

func attributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"enabled":   types.BoolType,
		"trim":      types.BoolType,
		"type":      types.StringType,
		"freeze_fs": types.BoolType,
	}
}

func NullValue() Value {
	return types.ObjectNull(attributeTypes())
}

func NewValue(ctx context.Context, config *vms.GetResponseData, diags *diag.Diagnostics) Value {
	if config.Agent == nil {
		return NullValue()
	}

	model := Model{
		Enabled:  types.BoolPointerValue(config.Agent.Enabled.PointerBool()),
		Trim:     types.BoolPointerValue(config.Agent.TrimClonedDisks.PointerBool()),
		Type:     types.StringPointerValue(config.Agent.Type),
		FreezeFS: types.BoolValue(true),
	}

	if config.Agent.FreezeFS != nil {
		model.FreezeFS = types.BoolPointerValue(config.Agent.FreezeFS.PointerBool())
	}

	obj, d := types.ObjectValueFrom(ctx, attributeTypes(), model)
	diags.Append(d...)

	return obj
}

func FillCreateBody(ctx context.Context, planValue Value, body *vms.CreateRequestBody, diags *diag.Diagnostics) {
	if planValue.IsNull() || planValue.IsUnknown() {
		return
	}

	var plan Model
	d := planValue.As(ctx, &plan, basetypes.ObjectAsOptions{})
	diags.Append(d...)

	if diags.HasError() {
		return
	}

	body.Agent = plan.toAPI()
}

func FillUpdateBody(ctx context.Context, planValue, stateValue Value, body *vms.UpdateRequestBody, diags *diag.Diagnostics) {
	attribute.CheckDeleteBody(planValue, stateValue, body, "agent")

	if planValue.IsNull() || planValue.IsUnknown() || planValue.Equal(stateValue) {
		return
	}

	var plan Model
	d := planValue.As(ctx, &plan, basetypes.ObjectAsOptions{})
	diags.Append(d...)

	if diags.HasError() {
		return
	}

	body.Agent = plan.toAPI()
}

func (m *Model) toAPI() *vms.CustomAgent {
	enabled := proxmoxtypes.CustomBool(m.Enabled.ValueBool())
	trim := proxmoxtypes.CustomBool(m.Trim.ValueBool())
	freezeFS := proxmoxtypes.CustomBool(m.FreezeFS.ValueBool())
	agentType := m.Type.ValueString()

	return &vms.CustomAgent{
		Enabled:         &enabled,
		TrimClonedDisks: &trim,
		Type:            &agentType,
		FreezeFS:        &freezeFS,
	}
}
