/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package planmodifiers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/attribute"
)

// UseUnknownForNullConfigList returns a plan modifier that sets the value of an attribute
// to Unknown if the attribute is missing from the plan and the config is null AND the resource is not a clone.
//
// Use this for optional computed attributes that can be reset / removed by the user. If the resource is a clone,
// the value will be copied from the prior state (e.g. the clone source).
//
// The behavior for Terraform for Optional + Computed attributes is to copy the prior state
// if there is no configuration for it. This plan modifier will instead set the value to Unknown,
// so the provider can handle the attribute as needed.
func UseUnknownForNullConfigList(elementType attr.Type) planmodifier.List {
	return useUnknownForNullConfigList{elementType}
}

// useUnknownForNullConfigList implements the plan modifier.
type useUnknownForNullConfigList struct {
	elementType attr.Type
}

// Description returns a human-readable description of the plan modifier.
func (m useUnknownForNullConfigList) Description(_ context.Context) string {
	return "Value of this attribute will be set to Unknown if missing from the plan."
}

// MarkdownDescription returns a markdown description of the plan modifier.
func (m useUnknownForNullConfigList) MarkdownDescription(_ context.Context) string {
	return "Value of this attribute will be set to Unknown if missing from the plan. " +
		"Use for optional computed attributes that can be reset / removed by user."
}

// PlanModifyList implements the plan modification logic.
func (m useUnknownForNullConfigList) PlanModifyList(
	ctx context.Context,
	req planmodifier.ListRequest,
	resp *planmodifier.ListResponse,
) {
	if !m.isClone(ctx, req) {
		if req.PlanValue.IsNull() {
			return
		}

		if !req.ConfigValue.IsNull() {
			return
		}

		resp.PlanValue = types.ListUnknown(m.elementType)
	}
}

func (m useUnknownForNullConfigList) isClone(ctx context.Context, req planmodifier.ListRequest) bool {
	var cloneID types.Int64
	_ = req.Plan.GetAttribute(ctx, path.Root("clone").AtName("id"), &cloneID)

	return attribute.IsDefined(cloneID)
}
