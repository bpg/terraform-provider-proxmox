/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package proxmoxtf

import (
	"testing"
)

func Test_getCPUTypeValidator(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		value string
		valid bool
	}{
		{"empty", "", false},
		{"invalid", "invalid", false},
		{"valid", "host", true},
		{"valid", "qemu64", true},
		{"valid", "custom-abc", true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			f := getCPUTypeValidator()
			res := f(tt.value, nil)
			if !res.HasError() != tt.valid {
				t.Errorf("validate: '%s', want %v got %v", tt.value, tt.valid, res)
			}
		})
	}
}
