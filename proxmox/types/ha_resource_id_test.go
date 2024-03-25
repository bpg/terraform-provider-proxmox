/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package types

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestParseHAResourceID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    HAResourceID
		wantErr bool
	}{
		{"VM value", "vm:123", HAResourceID{HAResourceTypeVM, "123"}, false},
		{"container value", "ct:123", HAResourceID{HAResourceTypeContainer, "123"}, false},
		{"no semicolon", "ct", HAResourceID{}, true},
		{"invalid type", "blah:123", HAResourceID{}, true},
		{"invalid VM name", "vm:moo", HAResourceID{}, true},
		{"invalid container name", "ct:moo", HAResourceID{}, true},
		{"VM name too low", "vm:99", HAResourceID{}, true},
		{"VM name too high", "vm:1000000000", HAResourceID{}, true},
		{"container name too low", "ct:99", HAResourceID{}, true},
		{"container name too high", "ct:1000000000", HAResourceID{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := ParseHAResourceID(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseHAResourceID() error = %v, wantErr %v", err, tt.wantErr)
			}

			if got != tt.want {
				t.Errorf("ParseHAResourceID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHAResourceIDToString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		state HAResourceID
		want  string
	}{
		{"stringify VM", HAResourceID{HAResourceTypeVM, "123"}, "vm:123"},
		{"stringify CT", HAResourceID{HAResourceTypeContainer, "123"}, "ct:123"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.state.String(); got != tt.want {
				t.Errorf("HAResourceID.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHAResourceIDToJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		state HAResourceID
		want  string
	}{
		{"jsonify VM", HAResourceID{HAResourceTypeVM, "123"}, `"vm:123"`},
		{"jsonify CT", HAResourceID{HAResourceTypeContainer, "123"}, `"ct:123"`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := json.Marshal(tt.state)
			if err != nil {
				t.Errorf("json.Marshal(HAResourceID): err = %v", err)
			} else if !bytes.Equal(got, []byte(tt.want)) {
				t.Errorf("json.Marshal(HAResourceID) = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHAResourceIDFromJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		json    string
		want    HAResourceID
		wantErr bool
	}{
		{"VM", `"vm:123"`, HAResourceID{HAResourceTypeVM, "123"}, false},
		{"container", `"ct:123"`, HAResourceID{HAResourceTypeContainer, "123"}, false},
		{"invalid JSON", `\\/yo`, HAResourceID{}, true},
		{"incompatible type", `["yo"]`, HAResourceID{}, true},
		{"invalid content", `"nope:notatall"`, HAResourceID{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var got HAResourceID

			err := json.Unmarshal([]byte(tt.json), &got)
			if (err != nil) != tt.wantErr {
				t.Errorf("json.Unmarshal(HAResourceID) error = %v, wantErr %v", err, tt.wantErr)
			}

			if got != tt.want {
				t.Errorf("json.Unmarshal(HAResourceID) got = %v, want %v", got, tt.want)
			}
		})
	}
}
