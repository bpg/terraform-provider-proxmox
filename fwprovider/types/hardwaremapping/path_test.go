/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package hardwaremapping

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestPathValueFromTerraform(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		val         tftypes.Value
		expected    func(val PathValue) bool
		expectError bool
	}{
		"null value": {
			val: tftypes.NewValue(tftypes.String, nil),
			expected: func(val PathValue) bool {
				return val.IsNull()
			},
		},
		"unknown value": {
			val: tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
			expected: func(val PathValue) bool {
				return val.IsUnknown()
			},
		},
		"valid for PCI type": {
			val: tftypes.NewValue(tftypes.String, "8086:5916"),
			expected: func(val PathValue) bool {
				return val.ValueString() == "8086:5916"
			},
		},
		"valid for USB type": {
			val: tftypes.NewValue(tftypes.String, "1-5.2"),
			expected: func(val PathValue) bool {
				return val.ValueString() == "1-5.2"
			},
		},
	}

	for name, test := range tests {
		t.Run(
			name, func(t *testing.T) {
				t.Parallel()

				ctx := context.TODO()
				val, err := PathType{}.ValueFromTerraform(ctx, test.val)

				if err == nil && test.expectError {
					t.Fatal("expected error, got no error")
				}

				if err != nil && !test.expectError {
					t.Fatalf("got unexpected error: %s", err)
				}

				if !test.expected(val.(PathValue)) {
					t.Errorf("unexpected result")
				}
			},
		)
	}
}
