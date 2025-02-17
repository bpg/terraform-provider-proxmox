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
)

func TestCustomRNGDevice_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		line    string
		want    *CustomRNGDevice
		wantErr bool
	}{
		{
			name: "source only",
			line: `"source=urandom"`,
			want: &CustomRNGDevice{
				Source: "urandom",
			},
		},
		{
			name: "all options",
			line: `"source=/dev/random,max_bytes=1024,period=1000"`,
			want: &CustomRNGDevice{
				Source:   "/dev/random",
				MaxBytes: ptr.Ptr(1024),
				Period:   ptr.Ptr(1000),
			},
		},
		{
			name: "source with max_bytes",
			line: `"source=urandom,max_bytes=2048"`,
			want: &CustomRNGDevice{
				Source:   "urandom",
				MaxBytes: ptr.Ptr(2048),
			},
		},
		{
			name: "source with period",
			line: `"source=urandom,period=2000"`,
			want: &CustomRNGDevice{
				Source: "urandom",
				Period: ptr.Ptr(2000),
			},
		},
		{
			name:    "invalid JSON",
			line:    `{"source": "urandom"}`,
			wantErr: true,
		},
		{
			name:    "invalid max_bytes",
			line:    `"source=urandom,max_bytes=invalid"`,
			wantErr: true,
		},
		{
			name:    "invalid period",
			line:    `"source=urandom,period=invalid"`,
			wantErr: true,
		},
		{
			name: "single value source",
			line: `"urandom"`,
			want: &CustomRNGDevice{
				Source: "urandom",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := &CustomRNGDevice{}
			err := r.UnmarshalJSON([]byte(tt.line))

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, r)
		})
	}
}

func TestCustomRNGDevice_EncodeValues(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		device   *CustomRNGDevice
		key      string
		expected string
	}{
		{
			name: "source only",
			device: &CustomRNGDevice{
				Source: "urandom",
			},
			key:      "rng0",
			expected: "source=urandom",
		},
		{
			name: "all options",
			device: &CustomRNGDevice{
				Source:   "/dev/random",
				MaxBytes: ptr.Ptr(1024),
				Period:   ptr.Ptr(1000),
			},
			key:      "rng0",
			expected: "source=/dev/random,max_bytes=1024,period=1000",
		},
		{
			name: "source with max_bytes",
			device: &CustomRNGDevice{
				Source:   "urandom",
				MaxBytes: ptr.Ptr(2048),
			},
			key:      "rng0",
			expected: "source=urandom,max_bytes=2048",
		},
		{
			name: "source with period",
			device: &CustomRNGDevice{
				Source: "urandom",
				Period: ptr.Ptr(2000),
			},
			key:      "rng0",
			expected: "source=urandom,period=2000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			values := &url.Values{}
			err := tt.device.EncodeValues(tt.key, values)
			require.NoError(t, err)
			require.Equal(t, tt.expected, values.Get(tt.key))
		})
	}
}
