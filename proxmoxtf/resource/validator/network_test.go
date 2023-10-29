/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package validator

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMACAddress(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value string
		valid bool
	}{
		{"empty", "", true},
		{"invalid", "invalid", false},
		{"invalid: no dashes", "38-f9-d3-4b-f5-51", false},
		{"valid", "38:f9:d3:4b:f5:51", true},
		{"valid", "38:F9:D3:4B:F5:51", true},
		{"valid", "00:15:5d:00:09:03", true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			f := MACAddress()
			res := f(tt.value, nil)

			if tt.valid {
				require.Empty(t, res, "validate: '%s'", tt.value)
			} else {
				require.NotEmpty(t, res, "validate: '%s'", tt.value)
			}
		})
	}
}
