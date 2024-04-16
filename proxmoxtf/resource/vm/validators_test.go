/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package resource

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMachineType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value string
		valid bool
	}{
		{"empty is valid", "", true},
		{"invalid", "invalid", false},
		{"valid q35", "q35", true},
		{"valid pc-q35", "pc-q35-2.3", true},
		{"valid i440fx", "pc-i440fx-3.1+pve0", true},
		{"valid virt", "virt", true},
		{"invalid i440fx", "i440fx", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			f := MachineTypeValidator()
			res := f(tt.value, nil)

			if tt.valid {
				require.Empty(t, res, "validate: '%s'", tt.value)
			} else {
				require.NotEmpty(t, res, "validate: '%s'", tt.value)
			}
		})
	}
}
