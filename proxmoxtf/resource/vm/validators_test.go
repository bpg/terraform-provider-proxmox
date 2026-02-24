/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package resource

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCPUType(t *testing.T) {
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
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			f := CPUTypeValidator()
			res := f(tt.value, nil)

			if tt.valid {
				require.Empty(t, res, "validate: '%s'", tt.value)
			} else {
				require.NotEmpty(t, res, "validate: '%s'", tt.value)
			}
		})
	}
}

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
		{"valid q35 with viommu", "q35,viommu=virtio", true},
		{"invalid q35 with viommu", "q35,viommu=invalid", false},
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

func TestVmHostname(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value string
		valid bool
	}{
		{"empty", "", false},
		{"underscores", "my_name", false},
		{"trailing dot", "my-name.com.", false},
		{"starts with alphanumeric", "-my-name.com", false},
		{"ends with alphanumeric", "my-name.com!", false},
		{"single letter", "a", true},
		{"domain name", "my-name.com", true},
		{"multi domain", "my-name.com.edu.net.xyz.dev", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			f := HostnameValidator()
			res := f(tt.value, nil)

			valid := !res.HasError()
			assert.Equal(t, tt.valid, valid)
		})
	}
}
