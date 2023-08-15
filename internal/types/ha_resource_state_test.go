/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package types

import (
	"testing"
)

func TestParseHAResourceState(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    HAResourceState
		wantErr bool
	}{
		{"valid value started", "started", HAResourceStateStarted, false},
		{"valid value enabled", "enabled", HAResourceStateStarted, false},
		{"valid value stopped", "stopped", HAResourceStateStopped, false},
		{"valid value disabled", "disabled", HAResourceStateDisabled, false},
		{"valid value ignored", "ignored", HAResourceStateIgnored, false},
		{"empty value", "", HAResourceStateIgnored, true},
		{"invalid value", "blah", HAResourceStateIgnored, true},
	}
	for _, test := range tests {
		tt := test
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := ParseHAResourceState(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseHAResourceState() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("ParseHAResourceState() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHAResourceStateToString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		state HAResourceState
		want  string
	}{
		{"stringify started", HAResourceStateStarted, "started"},
		{"stringify stopped", HAResourceStateStopped, "stopped"},
		{"stringify disabled", HAResourceStateDisabled, "disabled"},
		{"stringify ignored", HAResourceStateIgnored, "ignored"},
	}
	for _, test := range tests {
		tt := test
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.state.String(); got != tt.want {
				t.Errorf("HAResourceState.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
