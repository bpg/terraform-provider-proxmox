/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package types

import (
	"encoding/json"
	"slices"
	"testing"
)

func TestCustomCommaSeparatedList_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		json    string
		want    []string
		wantErr bool
	}{
		{"comma-separated string", `"/tmp,/var"`, []string{"/tmp", "/var"}, false},
		{"single string", `"/tmp"`, []string{"/tmp"}, false},
		{"string array", `["/tmp","/var"]`, []string{"/tmp", "/var"}, false},
		{"single-element array", `["/tmp"]`, []string{"/tmp"}, false},
		{"number", `1`, nil, true},
		{"array of numbers", `[1,2]`, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var v CustomCommaSeparatedList

			err := v.UnmarshalJSON([]byte(tt.json))
			if (err != nil) != tt.wantErr {
				t.Fatalf("error=%v, wantErr=%v", err, tt.wantErr)
			}

			if !tt.wantErr && !slices.Equal(v, tt.want) {
				t.Fatalf("got=%q, want=%q", v, tt.want)
			}
		})
	}
}

func TestCustomCommaSeparatedList_UnmarshalJSON_InStruct(t *testing.T) {
	t.Parallel()

	// PVE returns a backup job's exclude-path as an array (GET /cluster/backup/{id}).
	var data struct {
		ExcludePath *CustomCommaSeparatedList `json:"exclude-path,omitempty"`
	}

	err := json.Unmarshal([]byte(`{"exclude-path":["/var/lib/systemd/credential.secret"]}`), &data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if data.ExcludePath == nil || !slices.Equal(*data.ExcludePath, []string{"/var/lib/systemd/credential.secret"}) {
		t.Fatalf("got=%v", data.ExcludePath)
	}
}
