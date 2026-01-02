/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package storage

import (
	"net/url"
	"testing"

	"github.com/google/go-querystring/query"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

func intPtr(i int) *int {
	return &i
}

func customInt64Ptr(i int64) *types.CustomInt64 {
	c := types.CustomInt64(i)
	return &c
}

// TestDataStoreWithBackups_String tests backup settings are encoded correctly into a string.
func TestDataStoreWithBackups_String(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		input    DataStoreWithBackups
		expected string
	}{
		{
			name:     "Empty struct",
			input:    DataStoreWithBackups{},
			expected: "",
		},
		{
			name:     "KeepLast only",
			input:    DataStoreWithBackups{KeepLast: intPtr(5)},
			expected: "keep-last=5",
		},
		{
			name:     "KeepHourly only",
			input:    DataStoreWithBackups{KeepHourly: intPtr(24)},
			expected: "keep-hourly=24",
		},
		{
			name:     "KeepDaily only",
			input:    DataStoreWithBackups{KeepDaily: intPtr(7)},
			expected: "keep-daily=7",
		},
		{
			name:     "KeepWeekly only",
			input:    DataStoreWithBackups{KeepWeekly: intPtr(4)},
			expected: "keep-weekly=4",
		},
		{
			name:     "KeepMonthly only",
			input:    DataStoreWithBackups{KeepMonthly: intPtr(12)},
			expected: "keep-monthly=12",
		},
		{
			name:     "KeepYearly only",
			input:    DataStoreWithBackups{KeepYearly: intPtr(3)},
			expected: "keep-yearly=3",
		},
		{
			name: "Multiple values",
			input: DataStoreWithBackups{
				KeepDaily:  intPtr(30),
				KeepWeekly: intPtr(8),
				KeepYearly: intPtr(10),
			},
			expected: "keep-daily=30,keep-weekly=8,keep-yearly=10",
		},
		{
			name: "All values set",
			input: DataStoreWithBackups{
				KeepLast:    intPtr(1),
				KeepHourly:  intPtr(2),
				KeepDaily:   intPtr(3),
				KeepWeekly:  intPtr(4),
				KeepMonthly: intPtr(5),
				KeepYearly:  intPtr(6),
			},
			expected: "keep-last=1,keep-hourly=2,keep-daily=3,keep-weekly=4,keep-monthly=5,keep-yearly=6",
		},
		{
			name:     "MaxProtectedBackups should be ignored",
			input:    DataStoreWithBackups{MaxProtectedBackups: customInt64Ptr(10)},
			expected: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			result := tc.input.String()
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestDataStoreWithBackups_EncodeValues(t *testing.T) {
	t.Parallel()

	input := DataStoreWithBackups{
		MaxProtectedBackups: customInt64Ptr(10),
		KeepDaily:           intPtr(7),
	}

	values := url.Values{}
	err := input.EncodeValues("", &values)
	require.NoError(t, err)

	assert.Equal(t, "10", values.Get("max-protected-backups"))
	assert.Equal(t, "keep-daily=7", values.Get("prune-backups"))
}

func TestStorageCreateRequest_QueryEncoding_Backups(t *testing.T) {
	t.Parallel()

	req := NFSStorageCreateRequest{
		NFSStorageMutableFields: NFSStorageMutableFields{
			Backups: DataStoreWithBackups{
				MaxProtectedBackups: customInt64Ptr(9),
				KeepHourly:          intPtr(24),
			},
		},
	}

	values, err := query.Values(req)
	require.NoError(t, err)

	assert.Equal(t, "9", values.Get("max-protected-backups"))
	assert.Equal(t, "keep-hourly=24", values.Get("prune-backups"))
	assert.Empty(t, values.Get("max"))
}

func TestStorageCreateRequest_QueryEncoding_Shared(t *testing.T) {
	t.Parallel()

	shared := true

	req := LVMStorageCreateRequest{
		LVMStorageMutableFields: LVMStorageMutableFields{
			DataStoreCommonMutableFields: DataStoreCommonMutableFields{
				Shared: types.CustomBoolPtr(&shared),
			},
		},
	}

	values, err := query.Values(req)
	require.NoError(t, err)

	assert.Equal(t, "1", values.Get("shared"))
}
