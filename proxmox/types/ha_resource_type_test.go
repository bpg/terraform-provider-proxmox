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

func TestParseHAResourceType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    HAResourceType
		wantErr bool
	}{
		{"valid value vm", "vm", HAResourceTypeVM, false},
		{"valid value ct", "ct", HAResourceTypeContainer, false},
		{"empty value", "", _haResourceTypeValue, true},
		{"invalid value", "blah", _haResourceTypeValue, true},
	}
	for _, test := range tests {
		tt := test
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := ParseHAResourceType(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseHAResourceType() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("ParseHAResourceType() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHAResourceTypeToString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		resType HAResourceType
		want    string
	}{
		{"stringify vm", HAResourceTypeVM, "vm"},
		{"stringify ct", HAResourceTypeContainer, "ct"},
	}
	for _, test := range tests {
		tt := test
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.resType.String(); got != tt.want {
				t.Errorf("HAResourceType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHAResourceTypeToJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		state HAResourceType
		want  string
	}{
		{"jsonify vm", HAResourceTypeVM, `"vm"`},
		{"jsonify container", HAResourceTypeContainer, `"ct"`},
	}
	for _, test := range tests {
		tt := test
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := json.Marshal(tt.state)
			if err != nil {
				t.Errorf("json.Marshal(HAResourceType): err = %v", err)
			} else if !bytes.Equal(got, []byte(tt.want)) {
				t.Errorf("json.Marshal(HAResourceType) = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHAResourceTypeFromJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		json    string
		want    HAResourceType
		wantErr bool
	}{
		{"started", `"vm"`, HAResourceTypeVM, false},
		{"container", `"ct"`, HAResourceTypeContainer, false},
		{"invalid JSON", `\\/yo`, HAResourceTypeVM, true},
		{"incompatible type", `["yo"]`, HAResourceTypeVM, true},
		{"invalid content", `"nope"`, HAResourceTypeVM, true},
	}
	for _, test := range tests {
		tt := test
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var got HAResourceType

			err := json.Unmarshal([]byte(tt.json), &got)
			if (err != nil) != tt.wantErr {
				t.Errorf("json.Unmarshal(HAResourceType) error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && got != tt.want {
				t.Errorf("json.Unmarshal(HAResourceType) got = %v, want %v", got, tt.want)
			}
		})
	}
}
