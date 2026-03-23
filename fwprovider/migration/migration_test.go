/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package migration

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestPrefixMoveState_MatchingSourceTypeName(t *testing.T) {
	t.Parallel()

	testSchema := schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Required: true},
		},
	}

	raw := tftypes.NewValue(tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"id": tftypes.String,
		},
	}, map[string]tftypes.Value{
		"id": tftypes.NewValue(tftypes.String, "test-id"),
	})

	sourceState := tfsdk.State{Schema: testSchema, Raw: raw}

	schemaCopy := testSchema

	targetState := tfsdk.State{Schema: schemaCopy, Raw: tftypes.NewValue(tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"id": tftypes.String,
		},
	}, nil)}

	req := resource.MoveStateRequest{
		SourceTypeName: "proxmox_virtual_environment_example",
		SourceState:    &sourceState,
	}

	resp := &resource.MoveStateResponse{
		TargetState: targetState,
	}

	mover := PrefixMoveState("proxmox_virtual_environment_example", testSchema)
	mover.StateMover(context.Background(), req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected errors: %s", resp.Diagnostics.Errors())
	}

	// Verify the target state received the source raw value by comparing Raw directly.
	if !resp.TargetState.Raw.Equal(raw) {
		t.Errorf("expected target Raw to equal source Raw")
	}
}

func TestPrefixMoveState_NonMatchingSourceTypeName(t *testing.T) {
	t.Parallel()

	testSchema := schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Required: true},
		},
	}

	raw := tftypes.NewValue(tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"id": tftypes.String,
		},
	}, map[string]tftypes.Value{
		"id": tftypes.NewValue(tftypes.String, "should-not-copy"),
	})

	sourceState := tfsdk.State{Schema: testSchema, Raw: raw}

	// Target state starts with a null raw value.
	targetRaw := tftypes.NewValue(tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"id": tftypes.String,
		},
	}, nil)

	targetState := tfsdk.State{Schema: testSchema, Raw: targetRaw}

	req := resource.MoveStateRequest{
		SourceTypeName: "proxmox_virtual_environment_other",
		SourceState:    &sourceState,
	}

	resp := &resource.MoveStateResponse{
		TargetState: targetState,
	}

	mover := PrefixMoveState("proxmox_virtual_environment_example", testSchema)
	mover.StateMover(context.Background(), req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected errors: %s", resp.Diagnostics.Errors())
	}

	// Target state should remain unchanged (null).
	if !resp.TargetState.Raw.IsNull() {
		t.Error("expected target state to remain null for non-matching source type name")
	}
}

func TestDeprecationMessage(t *testing.T) {
	t.Parallel()

	got := DeprecationMessage("proxmox_example")
	expected := "Use proxmox_example instead. This resource will be removed in v1.0."

	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}
