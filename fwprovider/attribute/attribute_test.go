/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package attribute_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/attribute"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/types/stringset"
	proxmoxtypes "github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

// fakeBody is a minimal DeleteAppender used by CheckDeleteBody tests.
type fakeBody struct {
	deletes []string
}

func (b *fakeBody) AppendDelete(apiName string) {
	b.deletes = append(b.deletes, apiName)
}

func newStringSet(t *testing.T, elements []string) stringset.Value {
	t.Helper()

	var diags diag.Diagnostics

	v := stringset.NewValueList(elements, &diags)
	require.False(t, diags.HasError(), "stringset.NewValueList diags: %v", diags)

	return v
}

func TestCheckDelete(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name       string
		plan       attr.Value
		state      attr.Value
		apiName    string
		wantDelete []string
	}{
		{
			name:       "null plan, known state adds delete",
			plan:       types.StringNull(),
			state:      types.StringValue("foo"),
			apiName:    "alias",
			wantDelete: []string{"alias"},
		},
		{
			name:       "null plan, null state does nothing",
			plan:       types.StringNull(),
			state:      types.StringNull(),
			apiName:    "alias",
			wantDelete: nil,
		},
		{
			name:       "known plan, known state does nothing (population handled elsewhere)",
			plan:       types.StringValue("bar"),
			state:      types.StringValue("foo"),
			apiName:    "alias",
			wantDelete: nil,
		},
		{
			name:       "unknown plan, known state does nothing (unknown != deletion)",
			plan:       types.StringUnknown(),
			state:      types.StringValue("foo"),
			apiName:    "alias",
			wantDelete: nil,
		},
		{
			name:       "empty stringset plan, non-empty state adds delete",
			plan:       newStringSet(t, nil),
			state:      newStringSet(t, []string{"a", "b"}),
			apiName:    "tags",
			wantDelete: []string{"tags"},
		},
		{
			name:       "non-empty stringset plan, non-empty state does nothing",
			plan:       newStringSet(t, []string{"c"}),
			state:      newStringSet(t, []string{"a", "b"}),
			apiName:    "tags",
			wantDelete: nil,
		},
		{
			name:       "empty list plan, non-empty state adds delete",
			plan:       types.ListValueMust(types.StringType, nil),
			state:      types.ListValueMust(types.StringType, []attr.Value{types.StringValue("x")}),
			apiName:    "vmid",
			wantDelete: []string{"vmid"},
		},
		{
			name:       "empty map plan, non-empty state adds delete",
			plan:       types.MapValueMust(types.StringType, nil),
			state:      types.MapValueMust(types.StringType, map[string]attr.Value{"k": types.StringValue("v")}),
			apiName:    "ide0",
			wantDelete: []string{"ide0"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var toDelete []string

			attribute.CheckDelete(tc.plan, tc.state, &toDelete, tc.apiName)
			assert.Equal(t, tc.wantDelete, toDelete, "slice form")

			body := &fakeBody{}
			attribute.CheckDeleteBody(tc.plan, tc.state, body, tc.apiName)
			assert.Equal(t, tc.wantDelete, body.deletes, "body form")
		})
	}
}

func TestCheckDeleteBody_AppendsInOrder(t *testing.T) {
	t.Parallel()

	body := &fakeBody{}

	attribute.CheckDeleteBody(types.StringNull(), types.StringValue("a"), body, "first")
	attribute.CheckDeleteBody(types.StringValue("kept"), types.StringValue("kept"), body, "second") // no-op
	attribute.CheckDeleteBody(types.StringNull(), types.StringValue("b"), body, "third")

	assert.Equal(t, []string{"first", "third"}, body.deletes)
}

func TestStringValueFromPtr(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		val  *string
		want string
	}{
		{"nil", nil, ""},
		{"zero", new(string), ""},
		{"non-empty", new("hello"), "hello"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := attribute.StringValueFromPtr(tt.val)
			if got.ValueString() != tt.want {
				t.Errorf("StringValueFromPtr() = %q, want %q", got.ValueString(), tt.want)
			}
		})
	}
}

func TestInt64ValueFromPtr(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		val  *int64
		want int64
	}{
		{"nil", nil, 0},
		{"zero", new(int64), 0},
		{"positive", new(int64(42)), 42},
		{"negative", new(int64(-1)), -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := attribute.Int64ValueFromPtr(tt.val)
			if got.ValueInt64() != tt.want {
				t.Errorf("Int64ValueFromPtr() = %d, want %d", got.ValueInt64(), tt.want)
			}
		})
	}
}

func TestFloat64ValueFromPtr(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		val  *float64
		want float64
	}{
		{"nil", nil, 0},
		{"zero", new(float64), 0},
		{"positive", new(3.14), 3.14},
		{"negative", new(-2.5), -2.5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := attribute.Float64ValueFromPtr(tt.val)
			if got.ValueFloat64() != tt.want {
				t.Errorf("Float64ValueFromPtr() = %f, want %f", got.ValueFloat64(), tt.want)
			}
		})
	}
}

func TestBoolValueFromPtr(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		val  *bool
		want bool
	}{
		{"nil", nil, false},
		{"false", new(bool), false},
		{"true", new(true), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := attribute.BoolValueFromPtr(tt.val)
			if got.ValueBool() != tt.want {
				t.Errorf("BoolValueFromPtr() = %t, want %t", got.ValueBool(), tt.want)
			}
		})
	}
}

func TestBoolValueFromCustomBoolPtr(t *testing.T) {
	t.Parallel()

	trueVal := proxmoxtypes.CustomBool(true)
	falseVal := proxmoxtypes.CustomBool(false)

	tests := []struct {
		name string
		val  *proxmoxtypes.CustomBool
		want bool
	}{
		{"nil", nil, false},
		{"false", &falseVal, false},
		{"true", &trueVal, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := attribute.BoolValueFromCustomBoolPtr(tt.val)
			if got.ValueBool() != tt.want {
				t.Errorf("BoolValueFromCustomBoolPtr() = %t, want %t", got.ValueBool(), tt.want)
			}
		})
	}
}
