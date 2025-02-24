/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package types

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func Test_IPCIDRTypeValueFromTerraform(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		val         tftypes.Value
		expected    func(val IPCIDRValue) bool
		expectError bool
	}{
		"null value": {
			val: tftypes.NewValue(tftypes.String, nil),
			expected: func(val IPCIDRValue) bool {
				return val.IsNull()
			},
		},
		"unknown value": {
			val: tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
			expected: func(val IPCIDRValue) bool {
				return val.IsUnknown()
			},
		},
		"valid IPv4/CIDR": {
			val: tftypes.NewValue(tftypes.String, "1.2.3.4/32"),
			expected: func(val IPCIDRValue) bool {
				return val.ValueString() == "1.2.3.4/32"
			},
		},
		"valid IPv6/CIDR": {
			val: tftypes.NewValue(tftypes.String, "2001:db8::/32"),
			expected: func(val IPCIDRValue) bool {
				return val.ValueString() == "2001:db8::/32"
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			val, err := IPCIDRType{}.ValueFromTerraform(ctx, test.val)

			if err == nil && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if err != nil && !test.expectError {
				t.Fatalf("got unexpected error: %s", err)
			}

			if !test.expected(val.(IPCIDRValue)) {
				t.Errorf("unexpected result")
			}
		})
	}
}

func Test_IPCIDRTypeValidate(t *testing.T) {
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
		"valid IPv4 string": {
			val: tftypes.NewValue(tftypes.String, "1.2.3.4/32"),
		},
		"valid IPv6 string": {
			val: tftypes.NewValue(tftypes.String, "2001:db8::/32"),
		},
		"invalid string": {
			val:         tftypes.NewValue(tftypes.String, "not ok"),
			expectError: true,
		},
		"invalid IPv4 string no CIDR": {
			val:         tftypes.NewValue(tftypes.String, "1.2.3.4"),
			expectError: true,
		},
		"invalid IPv6 string no CIDR": {
			val:         tftypes.NewValue(tftypes.String, "2001:db8::"),
			expectError: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()

			diags := IPCIDRType{}.Validate(ctx, test.val, path.Root("test"))

			if !diags.HasError() && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if diags.HasError() && !test.expectError {
				t.Fatalf("got unexpected error: %s", diags)
			}
		})
	}
}
