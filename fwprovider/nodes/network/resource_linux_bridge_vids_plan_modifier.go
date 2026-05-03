/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package network

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// vidsPlanModifier resolves the `vids` plan value from the bridge's vlan_aware
// state, avoiding noisy "(known after apply)" output on every plan.
//
// Why a custom modifier instead of `stringplanmodifier.UseStateForUnknown()`:
// PVE/ifupdown2 silently substitutes "2-4094" as the implicit `bridge-vids`
// default whenever a bridge is VLAN-aware and the field is unset. That forces
// the provider to surface the value via Computed (otherwise Read on a freshly
// created bridge produces "inconsistent result after apply: was null, but now
// cty.StringVal(\"2-4094\")"). UseStateForUnknown alone is then unsafe in the
// other direction: when a user toggles vlan_aware from true to false the bridge
// stops storing bridge_vids and PVE returns nil, but the stock modifier would
// have pinned the prior "2-4094" value into the plan, producing the mirror
// inconsistency on apply.
//
// The vlan_aware-aware branches below cover both transitions:
//
//   - vlan_aware = false → vids is null (PVE doesn't store bridge_vids on
//     non-VLAN-aware bridges).
//   - vlan_aware = true with config null and a non-null state → preserve state
//     (no spurious diff on refresh).
//   - vlan_aware = true on create with no config or state → leave unknown so
//     apply can pick up PVE's implicit default ("2-4094").
type vidsPlanModifier struct{}

func (vidsPlanModifier) Description(_ context.Context) string {
	return "Resolves vids from vlan_aware: null when vlan_aware is false; " +
		"otherwise uses state when config is unset to keep refreshes quiet."
}

func (m vidsPlanModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (vidsPlanModifier) PlanModifyString(
	ctx context.Context,
	req planmodifier.StringRequest,
	resp *planmodifier.StringResponse,
) {
	// Defer to an explicit config value.
	if !req.ConfigValue.IsNull() {
		return
	}

	var vlanAware types.Bool

	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("vlan_aware"), &vlanAware)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// vlan_aware itself is unresolved; nothing useful to compute yet.
	if vlanAware.IsUnknown() {
		return
	}

	if !vlanAware.ValueBool() {
		resp.PlanValue = types.StringNull()

		return
	}

	// VLAN-aware bridge with no config value: prefer state to avoid a
	// "(known after apply)" diff on every refresh.
	if !req.StateValue.IsNull() && !req.StateValue.IsUnknown() {
		resp.PlanValue = req.StateValue
	}
}
