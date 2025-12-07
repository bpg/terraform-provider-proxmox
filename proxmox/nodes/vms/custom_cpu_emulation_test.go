/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package vms

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/proxmox/helpers/ptr"
	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

func TestCustomCPUEmulation_EncodeValues(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		emulation *CustomCPUEmulation
		key       string
		expected  string
	}{
		{
			name: "type only - should output just type value",
			emulation: &CustomCPUEmulation{
				Type: "x86-64-v4",
			},
			key:      "cpu",
			expected: "x86-64-v4",
		},
		{
			name: "type with flags - should output cputype= format",
			emulation: &CustomCPUEmulation{
				Type:  "x86-64-v4",
				Flags: &[]string{"+avx", "+sse"},
			},
			key:      "cpu",
			expected: "cputype=x86-64-v4,flags=+avx;+sse",
		},
		{
			name: "type with hidden - should output cputype= format",
			emulation: &CustomCPUEmulation{
				Type:   "x86-64-v4",
				Hidden: types.CustomBool(true).Pointer(),
			},
			key:      "cpu",
			expected: "cputype=x86-64-v4,hidden=1",
		},
		{
			name: "type with hv-vendor-id - should output cputype= format",
			emulation: &CustomCPUEmulation{
				Type:       "x86-64-v4",
				HVVendorID: ptr.Ptr("vendor123"),
			},
			key:      "cpu",
			expected: "cputype=x86-64-v4,hv-vendor-id=vendor123",
		},
		{
			name: "type with all options - should output cputype= format",
			emulation: &CustomCPUEmulation{
				Type:       "x86-64-v4",
				Flags:      &[]string{"+avx"},
				Hidden:     types.CustomBool(false).Pointer(),
				HVVendorID: ptr.Ptr("vendor123"),
			},
			key:      "cpu",
			expected: "cputype=x86-64-v4,flags=+avx,hidden=0,hv-vendor-id=vendor123",
		},
		{
			name: "type only - kvm64",
			emulation: &CustomCPUEmulation{
				Type: "kvm64",
			},
			key:      "cpu",
			expected: "kvm64",
		},
		{
			name: "type only - qemu64",
			emulation: &CustomCPUEmulation{
				Type: "qemu64",
			},
			key:      "cpu",
			expected: "qemu64",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			values := &url.Values{}
			err := tt.emulation.EncodeValues(tt.key, values)
			require.NoError(t, err)
			require.Equal(t, tt.expected, values.Get(tt.key))
		})
	}
}

func TestCustomCPUEmulation_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		line    string
		want    *CustomCPUEmulation
		wantErr bool
	}{
		{
			name: "type only format (GUI format)",
			line: `"x86-64-v4"`,
			want: &CustomCPUEmulation{
				Type: "x86-64-v4",
			},
		},
		{
			name: "cputype= format (provider format)",
			line: `"cputype=x86-64-v4"`,
			want: &CustomCPUEmulation{
				Type: "x86-64-v4",
			},
		},
		{
			name: "cputype= with flags",
			line: `"cputype=x86-64-v4,flags=+avx;+sse"`,
			want: &CustomCPUEmulation{
				Type:  "x86-64-v4",
				Flags: &[]string{"+avx", "+sse"},
			},
		},
		{
			name: "cputype= with hidden",
			line: `"cputype=x86-64-v4,hidden=1"`,
			want: &CustomCPUEmulation{
				Type:   "x86-64-v4",
				Hidden: types.CustomBool(true).Pointer(),
			},
		},
		{
			name: "cputype= with hv-vendor-id",
			line: `"cputype=x86-64-v4,hv-vendor-id=vendor123"`,
			want: &CustomCPUEmulation{
				Type:       "x86-64-v4",
				HVVendorID: ptr.Ptr("vendor123"),
			},
		},
		{
			name: "cputype= with all options",
			line: `"cputype=x86-64-v4,flags=+avx,hidden=0,hv-vendor-id=vendor123"`,
			want: &CustomCPUEmulation{
				Type:       "x86-64-v4",
				Flags:      &[]string{"+avx"},
				Hidden:     types.CustomBool(false).Pointer(),
				HVVendorID: ptr.Ptr("vendor123"),
			},
		},
		{
			name: "type only - kvm64",
			line: `"kvm64"`,
			want: &CustomCPUEmulation{
				Type: "kvm64",
			},
		},
		{
			name: "type only - qemu64",
			line: `"qemu64"`,
			want: &CustomCPUEmulation{
				Type: "qemu64",
			},
		},
		{
			name:    "invalid JSON",
			line:    `{"cputype": "x86-64-v4"}`,
			wantErr: true,
		},
		{
			name:    "empty string",
			line:    `""`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := &CustomCPUEmulation{}
			err := r.UnmarshalJSON([]byte(tt.line))

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want.Type, r.Type)

			if tt.want.Flags != nil {
				require.NotNil(t, r.Flags)
				require.Equal(t, *tt.want.Flags, *r.Flags)
			} else {
				require.Nil(t, r.Flags)
			}

			if tt.want.Hidden != nil {
				require.NotNil(t, r.Hidden)
				require.Equal(t, *tt.want.Hidden, *r.Hidden)
			} else {
				require.Nil(t, r.Hidden)
			}

			if tt.want.HVVendorID != nil {
				require.NotNil(t, r.HVVendorID)
				require.Equal(t, *tt.want.HVVendorID, *r.HVVendorID)
			} else {
				require.Nil(t, r.HVVendorID)
			}
		})
	}
}

func TestCustomCPUEmulation_RoundTrip(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "type only format round trip",
			input:    `"x86-64-v4"`,
			expected: "x86-64-v4",
		},
		{
			name:     "cputype= format round trip",
			input:    `"cputype=x86-64-v4"`,
			expected: "x86-64-v4",
		},
		{
			name:     "cputype= with flags round trip",
			input:    `"cputype=x86-64-v4,flags=+avx"`,
			expected: "cputype=x86-64-v4,flags=+avx",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := &CustomCPUEmulation{}
			err := r.UnmarshalJSON([]byte(tt.input))
			require.NoError(t, err)

			values := &url.Values{}
			err = r.EncodeValues("cpu", values)
			require.NoError(t, err)

			require.Equal(t, tt.expected, values.Get("cpu"))
		})
	}
}
