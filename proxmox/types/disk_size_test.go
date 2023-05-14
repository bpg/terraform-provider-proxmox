/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseDiskSize(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		size    *string
		want    int64
		wantErr bool
	}{
		{"handle null size", nil, 0, false},
		{"parse TB", StrPtr("2TB"), 2199023255552, false},
		{"parse T", StrPtr("2T"), 2199023255552, false},
		{"parse fraction T", StrPtr("2.2T"), 2418925581108, false},
		{"parse GB", StrPtr("2GB"), 2147483648, false},
		{"parse G", StrPtr("2G"), 2147483648, false},
		{"parse M", StrPtr("2048M"), 2147483648, false},
		{"parse MB", StrPtr("2048MB"), 2147483648, false},
		{"parse MiB", StrPtr("2048MiB"), 2147483648, false},
		{"parse K", StrPtr("1K"), 1024, false},
		{"parse KB", StrPtr("2KB"), 2048, false},
		{"parse KiB", StrPtr("4KiB"), 4096, false},
		{"parse no units as bytes", StrPtr("12345"), 12345, false},
		{"error on bad format string", StrPtr("20l8G"), -1, true},
		{"error on unknown unit string", StrPtr("2048W"), -1, true},
		{"error on arbitrary string", StrPtr("something"), -1, true},
	}
	for _, test := range tests {
		tt := test
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := parseDiskSize(tt.size)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseDiskSize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseDiskSize() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatDiskSize(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		size int64
		want string
	}{
		{"handle 0 size", 0, "0"},
		{"handle bytes", 1001, "1001"},
		{"handle kilobytes", 1234, "1.21K"},
		{"handle megabytes", 2097152, "2M"},
		{"handle gigabytes", 2147483648, "2G"},
		{"handle terabytes", 2199023255552, "2T"},
	}
	for _, test := range tests {
		tt := test
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := formatDiskSize(tt.size); got != tt.want {
				t.Errorf("formatDiskSize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToFromGigabytes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		size int
		want string
	}{
		{"handle 0 size", 0, "0"},
		{"handle 99 GB", 99, "99G"},
		{"handle 100 GB", 100, "100G"},
		{"handle 101 GB", 101, "101G"},
		{"handle 1023 GB", 1023, "1023G"},
		{"handle 1024 GB", 1024, "1T"},
		{"handle 1025 GB", 1025, "1.01T"},
	}
	for _, test := range tests {
		tt := test
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ds := DiskSizeFromGigabytes(tt.size)
			gb := ds.InGigabytes()
			assert.Equal(t, tt.size, gb)
			if got := ds.String(); got != tt.want {
				t.Errorf("DiskSize.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
