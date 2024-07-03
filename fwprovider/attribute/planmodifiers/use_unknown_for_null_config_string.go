/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package planmodifiers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/attribute"
)

// UseUnknownForNullConfigString returns a plan modifier that sets the value of an attribute
// to Unknown if the attribute is missing from the plan and the config is null AND the resource is not a clone.
//
// Use this for optional computed attributes that can be reset / removed by the user. If the resource is a clone,
// the value will be copied from the prior state (e.g. the clone source).
//
// The behavior for Terraform for Optional + Computed attributes is to copy the prior state
// if there is no configuration for it. This plan modifier will instead set the value to Unknown,
// so the provider can handle the attribute as needed.
func UseUnknownForNullConfigString() planmodifier.String {
	return useUnknownForNullConfigString{}
}

// useUnknownForNullConfigString implements the plan modifier.
type useUnknownForNullConfigString struct{}

// Description returns a human-readable description of the plan modifier.
func (m useUnknownForNullConfigString) Description(_ context.Context) string {
	return "Value of this attribute will be set to Unknown if missing from the plan."
}

// MarkdownDescription returns a markdown description of the plan modifier.
func (m useUnknownForNullConfigString) MarkdownDescription(_ context.Context) string {
	return "Value of this attribute will be set to Unknown if missing from the plan. " +
		"Use for optional computed attributes that can be reset / removed by user."
}

// PlanModifyString implements the plan modification logic.
func (m useUnknownForNullConfigString) PlanModifyString(
	ctx context.Context,
	req planmodifier.StringRequest,
	resp *planmodifier.StringResponse,
) {
	if !m.isClone(ctx, req) {
		if req.PlanValue.IsNull() {
			return
		}

		if !req.ConfigValue.IsNull() {
			return
		}

		resp.PlanValue = types.StringUnknown()
	}
}

func (m useUnknownForNullConfigString) isClone(ctx context.Context, req planmodifier.StringRequest) bool {
	var cloneID types.Int64
	_ = req.Plan.GetAttribute(ctx, path.Root("clone").AtName("id"), &cloneID)

	return attribute.IsDefined(cloneID)
}
