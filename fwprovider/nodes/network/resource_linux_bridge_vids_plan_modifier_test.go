/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package network

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/require"
)

// TestVIDsPlanModifier exercises every decision branch of vidsPlanModifier
// without requiring a Proxmox endpoint: the modifier is pure logic over its
// inputs and is the place a regression is most likely to land silently
// (since acceptance tests don't run in fast feedback loops).
func TestVIDsPlanModifier(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	modifier := vidsPlanModifier{}

	// Minimal schema — vidsPlanModifier only reads vlan_aware via Plan.GetAttribute.
	planSchema := schema.Schema{
		Attributes: map[string]schema.Attribute{
			"vlan_aware": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
		},
	}
	planObjectType := tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"vlan_aware": tftypes.Bool,
		},
	}

	// makePlan builds a tfsdk.Plan whose only attribute is vlan_aware. Pass
	// a *bool value (nil for null), or an explicit tftypes.UnknownValue token.
	makePlan := func(vlanAware tftypes.Value) tfsdk.Plan {
		return tfsdk.Plan{
			Schema: planSchema,
			Raw: tftypes.NewValue(planObjectType, map[string]tftypes.Value{
				"vlan_aware": vlanAware,
			}),
		}
	}

	t.Run("config has explicit value: no-op", func(t *testing.T) {
		t.Parallel()

		req := planmodifier.StringRequest{
			Path:        path.Root("vids"),
			Plan:        makePlan(tftypes.NewValue(tftypes.Bool, true)),
			ConfigValue: types.StringValue("1 20 130"),
			StateValue:  types.StringNull(),
			PlanValue:   types.StringValue("1 20 130"),
		}
		resp := &planmodifier.StringResponse{PlanValue: req.PlanValue}

		modifier.PlanModifyString(ctx, req, resp)

		require.False(t, resp.Diagnostics.HasError())
		// PlanValue should be unchanged from the explicit config.
		require.True(t, resp.PlanValue.Equal(types.StringValue("1 20 130")))
	})

	t.Run("vlan_aware unknown: no-op (leave plan as-is)", func(t *testing.T) {
		t.Parallel()

		req := planmodifier.StringRequest{
			Path:        path.Root("vids"),
			Plan:        makePlan(tftypes.NewValue(tftypes.Bool, tftypes.UnknownValue)),
			ConfigValue: types.StringNull(),
			StateValue:  types.StringNull(),
			PlanValue:   types.StringUnknown(),
		}
		resp := &planmodifier.StringResponse{PlanValue: req.PlanValue}

		modifier.PlanModifyString(ctx, req, resp)

		require.False(t, resp.Diagnostics.HasError())
		require.True(t, resp.PlanValue.IsUnknown())
	})

	t.Run("vlan_aware = false: vids forced to null", func(t *testing.T) {
		t.Parallel()

		req := planmodifier.StringRequest{
			Path:        path.Root("vids"),
			Plan:        makePlan(tftypes.NewValue(tftypes.Bool, false)),
			ConfigValue: types.StringNull(),
			StateValue:  types.StringValue("10 20 30"), // would persist if not for the false branch
			PlanValue:   types.StringUnknown(),
		}
		resp := &planmodifier.StringResponse{PlanValue: req.PlanValue}

		modifier.PlanModifyString(ctx, req, resp)

		require.False(t, resp.Diagnostics.HasError())
		require.True(t, resp.PlanValue.IsNull())
	})

	t.Run("vlan_aware = true, state null: leave unknown for create", func(t *testing.T) {
		t.Parallel()

		req := planmodifier.StringRequest{
			Path:        path.Root("vids"),
			Plan:        makePlan(tftypes.NewValue(tftypes.Bool, true)),
			ConfigValue: types.StringNull(),
			StateValue:  types.StringNull(),
			PlanValue:   types.StringUnknown(),
		}
		resp := &planmodifier.StringResponse{PlanValue: req.PlanValue}

		modifier.PlanModifyString(ctx, req, resp)

		require.False(t, resp.Diagnostics.HasError())
		require.True(t, resp.PlanValue.IsUnknown())
	})

	t.Run("vlan_aware = true, state set: preserve state on refresh", func(t *testing.T) {
		t.Parallel()

		req := planmodifier.StringRequest{
			Path:        path.Root("vids"),
			Plan:        makePlan(tftypes.NewValue(tftypes.Bool, true)),
			ConfigValue: types.StringNull(),
			StateValue:  types.StringValue("10 20 30"),
			PlanValue:   types.StringUnknown(),
		}
		resp := &planmodifier.StringResponse{PlanValue: req.PlanValue}

		modifier.PlanModifyString(ctx, req, resp)

		require.False(t, resp.Diagnostics.HasError())
		require.True(t, resp.PlanValue.Equal(types.StringValue("10 20 30")))
	})
}
