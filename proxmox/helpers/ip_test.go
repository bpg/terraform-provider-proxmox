/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidGlobalUnicast(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		ipStr       string
		want        bool
		description string
	}{
		{
			name:        "accepts valid IPv4 addresses",
			ipStr:       "192.168.1.9",
			want:        true,
			description: "valid IPv4 address should be accepted",
		},
		{
			name:        "accepts public IPv4 addresses",
			ipStr:       "8.8.8.8",
			want:        true,
			description: "public IPv4 address should be accepted",
		},
		{
			name:        "accepts private IPv4 addresses",
			ipStr:       "10.0.0.1",
			want:        true,
			description: "private IPv4 address (10.0.0.0/8) should be accepted",
		},
		{
			name:        "accepts valid IPv6 global unicast addresses",
			ipStr:       "2001:db8::1",
			want:        true,
			description: "valid IPv6 global unicast address should be accepted",
		},
		{
			name:        "accepts IPv6 unique local addresses",
			ipStr:       "fc00::1",
			want:        true,
			description: "IPv6 unique local address (fc00::/7) should be accepted",
		},
		{
			name:        "rejects loopback addresses",
			ipStr:       "127.0.0.1",
			want:        false,
			description: "loopback address should be rejected",
		},
		{
			name:        "rejects IPv6 loopback address",
			ipStr:       "::1",
			want:        false,
			description: "IPv6 loopback address should be rejected",
		},
		{
			name:        "rejects IPv4 link-local address",
			ipStr:       "169.254.1.1",
			want:        false,
			description: "IPv4 link-local address (169.254.0.0/16) should be rejected",
		},
		{
			name:        "rejects link-local IPv6 addresses",
			ipStr:       "fe80::be24:11ff:fe87:d827",
			want:        false,
			description: "link-local IPv6 address should be rejected",
		},
		{
			name:        "rejects IPv6 link-local multicast",
			ipStr:       "ff02::1",
			want:        false,
			description: "IPv6 link-local multicast address should be rejected",
		},
		{
			name:        "rejects IPv4 multicast addresses",
			ipStr:       "224.0.0.1",
			want:        false,
			description: "IPv4 multicast address should be rejected",
		},
		{
			name:        "rejects IPv6 site-local multicast addresses",
			ipStr:       "ff05::1",
			want:        false,
			description: "IPv6 site-local multicast address should be rejected",
		},
		{
			name:        "rejects invalid IP addresses",
			ipStr:       "invalid-ip",
			want:        false,
			description: "invalid IP address should be rejected",
		},
		{
			name:        "rejects empty string",
			ipStr:       "",
			want:        false,
			description: "empty string should be rejected",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := IsValidGlobalUnicast(tt.ipStr)
			assert.Equal(t, tt.want, result, tt.description)
		})
	}
}
