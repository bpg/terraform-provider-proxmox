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

func Test_IPAddrTypeValueFromTerraform(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		val         tftypes.Value
		expected    func(val IPAddrValue) bool
		expectError bool
	}{
		"null value": {
			val: tftypes.NewValue(tftypes.String, nil),
			expected: func(val IPAddrValue) bool {
				return val.IsNull()
			},
		},
		"unknown value": {
			val: tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
			expected: func(val IPAddrValue) bool {
				return val.IsUnknown()
			},
		},
		"valid IPv4": {
			val: tftypes.NewValue(tftypes.String, "1.2.3.4"),
			expected: func(val IPAddrValue) bool {
				return val.ValueString() == "1.2.3.4"
			},
		},
		"valid IPv6": {
			val: tftypes.NewValue(tftypes.String, "2001:0db8:85a3:0000:0000:8a2e:0370:7334"),
			expected: func(val IPAddrValue) bool {
				return val.ValueString() == "2001:0db8:85a3:0000:0000:8a2e:0370:7334"
			},
		},
	}

	for name, test := range tests {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx := context.TODO()
			val, err := IPAddrType{}.ValueFromTerraform(ctx, test.val)

			if err == nil && test.expectError {
				t.Fatal("expected error, got no error")
			}
			if err != nil && !test.expectError {
				t.Fatalf("got unexpected error: %s", err)
			}

			if !test.expected(val.(IPAddrValue)) {
				t.Errorf("unexpected result")
			}
		})
	}
}

func Test_IPAddrTypeValidate(t *testing.T) {
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
			val: tftypes.NewValue(tftypes.String, "1.2.3.4"),
		},
		"valid IPv6 string": {
			val: tftypes.NewValue(tftypes.String, "2001:0db8:85a3:0000:0000:8a2e:0370:7334"),
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

			diags := IPAddrType{}.Validate(ctx, test.val, path.Root("test"))

			if !diags.HasError() && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if diags.HasError() && !test.expectError {
				t.Fatalf("got unexpected error: %s", diags)
			}
		})
	}
}
