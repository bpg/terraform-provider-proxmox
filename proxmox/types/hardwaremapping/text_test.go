/*
	This Source Code Form is subject to the terms of the Mozilla Public
	License, v. 2.0. If a copy of the MPL was not distributed with this
	file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/

package hardwaremapping

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncodeText(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		encoded string
	}{
		{"plain ASCII unchanged", "NVIDIA RTX 2000, GPU: primary", "NVIDIA RTX 2000, GPU: primary"},
		{"em dash", "a — b", "a %E2%80%94 b"},
		{"percent sign", "50% shared", "50%25 shared"},
		{"multi-byte unicode", "café ✓", "caf%C3%A9 %E2%9C%93"},
		{"control character", "a\nb", "a%0Ab"},
		{"empty", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, tt.encoded, EncodeText(tt.input))
			require.Equal(t, tt.input, DecodeText(tt.encoded), "decode must reverse encode")
		})
	}
}

func TestDecodeTextInvalidInput(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
	}{
		{"stray percent", "100% GPU"},
		{"invalid escape", "a%ZZb"},
		{"invalid UTF-8 after decode", "a%FFb"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, tt.input, DecodeText(tt.input), "undecodable values must be returned unchanged")
		})
	}
}

func TestMapStringDescriptionRoundTrip(t *testing.T) {
	t.Parallel()

	description := "café — 50% shared"
	hm := Map{
		ID:          "8086:5916",
		Node:        "pve",
		Description: &description,
	}

	encoded := hm.String()
	require.Contains(t, encoded, "description=caf%C3%A9 %E2%80%94 50%25 shared")

	parsed, err := ParseMap(encoded)
	require.NoError(t, err)
	require.NotNil(t, parsed.Description)
	require.Equal(t, description, *parsed.Description)
}
