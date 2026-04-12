/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package vms

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

func TestCustomCloudInitConfig_EncodeValues_Upgrade(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		config        CustomCloudInitConfig
		wantCiupgrade string // expected value, or "" if key should be absent
		wantPresent   bool   // whether ciupgrade key should be present at all
	}{
		{
			name:        "upgrade nil does not send ciupgrade",
			config:      CustomCloudInitConfig{},
			wantPresent: false,
		},
		{
			name:          "upgrade false sends ciupgrade=0",
			config:        CustomCloudInitConfig{Upgrade: types.CustomBool(false).Pointer()},
			wantCiupgrade: "0",
			wantPresent:   true,
		},
		{
			name:          "upgrade true sends ciupgrade=1",
			config:        CustomCloudInitConfig{Upgrade: types.CustomBool(true).Pointer()},
			wantCiupgrade: "1",
			wantPresent:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			values := &url.Values{}
			err := tt.config.EncodeValues("", values)
			require.NoError(t, err)

			if tt.wantPresent {
				require.Equal(t, tt.wantCiupgrade, values.Get("ciupgrade"),
					"ciupgrade should be %q", tt.wantCiupgrade)
			} else {
				require.False(t, values.Has("ciupgrade"),
					"ciupgrade should not be sent when Upgrade is nil (non-root users cannot set it)")
			}
		})
	}
}
