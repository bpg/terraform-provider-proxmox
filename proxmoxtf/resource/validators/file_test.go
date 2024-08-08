/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package validators

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFileID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value string
		valid bool
	}{
		{"empty", "", true},
		{"invalid", "invalid", false},
		{"valid", "local:vztmpl/zen-dns-0.1.tar.zst", true},
		{"valid when datastore name has dots", "terraform.proxmox.storage.compute.zen:vztmpl/zen-dns-0.1.tar.zst", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			f := FileID()
			res := f(tt.value, nil)

			if tt.valid {
				require.Empty(t, res, "validate: '%s'", tt.value)
			} else {
				require.NotEmpty(t, res, "validate: '%s'", tt.value)
			}
		})
	}
}

func TestFileMode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value string
		valid bool
	}{
		{"valid", "0700", true},
		{"invalid", "invalid", false},
		// Even though Go supports octal prefixes, we should not allow them in the string value to reduce the complexity.
		{"invalid", "0o700", false},
		{"invalid", "0x700", false},
		// Maximum value for uint32, incremented by one.
		{"too large", "4294967296", false},
		{"too small", "0", false},
		{"negative", "-1", false},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			f := FileMode()
			res := f(tt.value, nil)

			if tt.valid {
				require.Empty(t, res, "validate: '%s'", tt.value)
			} else {
				require.NotEmpty(t, res, "validate: '%s'", tt.value)
			}
		})
	}
}
