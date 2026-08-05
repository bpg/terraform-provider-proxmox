/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetAnyStringMapEnv(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		want    map[string]string
		wantErr string
	}{
		{
			name:  "unset",
			value: "",
			want:  nil,
		},
		{
			name:  "single pair",
			value: "CF-Access-Client-Id=abc.access",
			want:  map[string]string{"CF-Access-Client-Id": "abc.access"},
		},
		{
			name:  "multiple pairs with whitespace",
			value: " CF-Access-Client-Id = abc.access , CF-Access-Client-Secret = s3cret ",
			want: map[string]string{
				"CF-Access-Client-Id":     "abc.access",
				"CF-Access-Client-Secret": "s3cret",
			},
		},
		{
			name:  "value containing equals",
			value: "X-Token=a=b=c",
			want:  map[string]string{"X-Token": "a=b=c"},
		},
		{
			name:  "empty value",
			value: "X-Foo=",
			want:  map[string]string{"X-Foo": ""},
		},
		{
			name:  "trailing comma",
			value: "X-Foo=bar,",
			want:  map[string]string{"X-Foo": "bar"},
		},
		{
			name:    "missing equals",
			value:   "X-Foo",
			wantErr: "pair 1 is not in the `Name=Value` format",
		},
		{
			name:    "missing equals in second pair",
			value:   "X-Foo=bar,X-Baz",
			wantErr: "pair 2 is not in the `Name=Value` format",
		},
		{
			name:    "empty name",
			value:   "=bar",
			wantErr: "pair 1 has an empty name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value != "" {
				t.Setenv("PROXMOX_VE_TEST_HEADERS", tt.value)
			}

			got, err := GetAnyStringMapEnv("PROXMOX_VE_TEST_HEADERS")

			if tt.wantErr != "" {
				require.EqualError(t, err, tt.wantErr)
				// the variable holds credentials, so its content must not leak into the error
				require.NotContains(t, err.Error(), "bar")

				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
