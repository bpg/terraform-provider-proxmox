/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package containers

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCustomIDMaps_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		json     string
		expected CustomIDMaps
	}{
		{
			name:     "null value",
			json:     `null`,
			expected: nil,
		},
		{
			name:     "empty array",
			json:     `[]`,
			expected: nil,
		},
		{
			name: "single uid mapping",
			json: `[["lxc.idmap", "u 0 100000 65536"]]`,
			expected: CustomIDMaps{
				{Type: "uid", ContainerID: 0, HostID: 100000, Size: 65536},
			},
		},
		{
			name: "uid and gid mappings",
			json: `[["lxc.idmap", "u 0 100000 65536"], ["lxc.idmap", "g 0 100000 44"], ["lxc.idmap", "g 44 44 1"]]`,
			expected: CustomIDMaps{
				{Type: "uid", ContainerID: 0, HostID: 100000, Size: 65536},
				{Type: "gid", ContainerID: 0, HostID: 100000, Size: 44},
				{Type: "gid", ContainerID: 44, HostID: 44, Size: 1},
			},
		},
		{
			name: "ignores non-idmap entries",
			json: `[["lxc.idmap", "u 0 100000 65536"], ["lxc.cap.drop", ""], ["lxc.idmap", "g 0 100000 65536"]]`,
			expected: CustomIDMaps{
				{Type: "uid", ContainerID: 0, HostID: 100000, Size: 65536},
				{Type: "gid", ContainerID: 0, HostID: 100000, Size: 65536},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var idmaps CustomIDMaps
			err := json.Unmarshal([]byte(tt.json), &idmaps)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, idmaps)
		})
	}

	t.Run("error on unknown type character", func(t *testing.T) {
		t.Parallel()

		var idmaps CustomIDMaps
		err := json.Unmarshal([]byte(`[["lxc.idmap", "x 0 100000 65536"]]`), &idmaps)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "unknown idmap type")
	})

	t.Run("error on trailing extra fields", func(t *testing.T) {
		t.Parallel()

		var idmaps CustomIDMaps
		err := json.Unmarshal([]byte(`[["lxc.idmap", "u 0 100000 65536 extra-field"]]`), &idmaps)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "expected 4 fields, got 5")
	})

	t.Run("error on too few fields", func(t *testing.T) {
		t.Parallel()

		var idmaps CustomIDMaps
		err := json.Unmarshal([]byte(`[["lxc.idmap", "u 0 100000"]]`), &idmaps)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "expected 4 fields, got 3")
	})

	t.Run("error on non-integer container_id", func(t *testing.T) {
		t.Parallel()

		var idmaps CustomIDMaps
		err := json.Unmarshal([]byte(`[["lxc.idmap", "u abc 100000 65536"]]`), &idmaps)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "unable to parse container_id")
	})

	t.Run("error on non-integer host_id", func(t *testing.T) {
		t.Parallel()

		var idmaps CustomIDMaps
		err := json.Unmarshal([]byte(`[["lxc.idmap", "u 0 xyz 65536"]]`), &idmaps)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "unable to parse host_id")
	})

	t.Run("error on non-integer size", func(t *testing.T) {
		t.Parallel()

		var idmaps CustomIDMaps
		err := json.Unmarshal([]byte(`[["lxc.idmap", "u 0 100000 abc"]]`), &idmaps)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "unable to parse size")
	})
}

func TestGetResponseData_UnmarshalJSON_WithIDMaps(t *testing.T) {
	t.Parallel()

	jsonData := `{
		"digest": "abc123",
		"lxc": [
			["lxc.idmap", "u 0 100000 65536"],
			["lxc.idmap", "g 0 100000 65536"]
		]
	}`

	var data GetResponseData
	err := json.Unmarshal([]byte(jsonData), &data)
	require.NoError(t, err)

	require.NotNil(t, data.IDMaps)
	require.Len(t, data.IDMaps, 2)
	assert.Equal(t, "uid", data.IDMaps[0].Type)
	assert.Equal(t, 0, data.IDMaps[0].ContainerID)
	assert.Equal(t, 100000, data.IDMaps[0].HostID)
	assert.Equal(t, 65536, data.IDMaps[0].Size)
	assert.Equal(t, "gid", data.IDMaps[1].Type)
	assert.Equal(t, 0, data.IDMaps[1].ContainerID)
	assert.Equal(t, 100000, data.IDMaps[1].HostID)
	assert.Equal(t, 65536, data.IDMaps[1].Size)
}
