/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package vms

import (
	"testing"
)

func TestCustomNUMADevice_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		line    string
		want    *CustomNUMADevice
		wantErr bool
	}{
		{
			name: "numa device all options",
			line: `"cpus=1-2;3-4,hostnodes=1-2,memory=1024,policy=preferred"`,
			want: &CustomNUMADevice{
				CPUIDs:        []string{"1-2", "3-4"},
				HostNodeNames: &[]string{"1-2"},
				Memory:        new(1024),
				Policy:        new("preferred"),
			},
		},
		{
			name: "numa device cpus/memory only",
			line: `"cpus=1-2,memory=1024"`,
			want: &CustomNUMADevice{
				CPUIDs: []string{"1-2"},
				Memory: new(1024),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := &CustomNUMADevice{}
			if err := r.UnmarshalJSON([]byte(tt.line)); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
