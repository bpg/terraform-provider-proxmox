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

func Test_getCPUTypeValidator(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value string
		valid bool
	}{
		{"empty", "", false},
		{"invalid", "invalid", false},
		{"valid", "host", true},
		{"valid", "qemu64", true},
		{"valid", "custom-abc", true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require := require.New(t)

			f := CPUType()
			res := f(tt.value, nil)

			if tt.valid {
				require.Empty(res, "validate: '%s'", tt.value)
			} else {
				require.NotEmpty(res, "validate: '%s'", tt.value)
			}
		})
	}
}
