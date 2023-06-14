/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package types

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func Test_IPv4TypeValueFromTerraform(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		val         tftypes.Value
		expected    func(val IPv4Value) bool
		expectError bool
	}{
		"null value": {
			val: tftypes.NewValue(tftypes.String, nil),
			expected: func(val IPv4Value) bool {
				return val.IsNull()
			},
		},
		"unknown value": {
			val: tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
			expected: func(val IPv4Value) bool {
				return val.IsUnknown()
			},
		},
		"valid IPV4": {
			val: tftypes.NewValue(tftypes.String, "1.2.3.4"),
			expected: func(val IPv4Value) bool {
				return val.ValueString() == "1.2.3.4"
			},
		},
	}

	for name, test := range tests {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx := context.TODO()
			val, err := IPv4Type{}.ValueFromTerraform(ctx, test.val)

			if err == nil && test.expectError {
				t.Fatal("expected error, got no error")
			}
			if err != nil && !test.expectError {
				t.Fatalf("got unexpected error: %s", err)
			}

			if !test.expected(val.(IPv4Value)) {
				t.Errorf("unexpected result")
			}
		})
	}
}

func Test_IPv4TypeValidate(t *testing.T) {
	t.Parallel()

	type testCase struct {
		val         tftypes.Value
		expectError bool
	}

	tests := map[string]testCase{
		"not a string": {
			val:         tftypes.NewValue(tftypes.Bool, true),
			expectError: true,
		},
		"unknown string": {
			val: tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
		},
		"null string": {
			val: tftypes.NewValue(tftypes.String, nil),
		},
		"valid string": {
			val: tftypes.NewValue(tftypes.String, "1.2.3.4"),
		},
		"invalid string": {
			val:         tftypes.NewValue(tftypes.String, "not ok"),
			expectError: true,
		},
	}

	for name, test := range tests {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx := context.TODO()

			diags := IPv4Type{}.Validate(ctx, test.val, path.Root("test"))

			if !diags.HasError() && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if diags.HasError() && !test.expectError {
				t.Fatalf("got unexpected error: %s", diags)
			}
		})
	}
}
