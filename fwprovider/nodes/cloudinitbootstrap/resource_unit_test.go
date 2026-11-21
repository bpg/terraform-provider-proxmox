/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package cloudinitbootstrap

import (
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/vms"
)

func TestChangedSeedAttrs(t *testing.T) {
	t.Parallel()

	state := bootstrapModel{
		UserData:      types.StringValue("#cloud-config\n"),
		MetaData:      types.StringNull(),
		NetworkConfig: types.StringValue("version: 2\n"),
		VendorData:    types.StringNull(),
		Timeout:       types.Int64Value(1800),
	}

	same := state
	same.MetaData = types.StringValue("") // empty and unset build the same seed
	same.Timeout = types.Int64Value(60)   // not seed content
	assert.Empty(t, changedSeedAttrs(same, state))

	changed := state
	changed.UserData = types.StringValue("#cloud-config\npackages: [htop]\n")
	changed.VendorData = types.StringValue("#cloud-config\n")
	changed.NetworkConfig = types.StringNull()
	assert.Equal(t, []string{"user_data", "network_config", "vendor_data"}, changedSeedAttrs(changed, state))
}

func TestCheckVMConfig(t *testing.T) {
	t.Parallel()

	// VM configs as PVE 9 returns them.
	tests := []struct {
		name, config, iface, wantErr string
	}{
		{"empty CD-ROM", `{"agent":"1","ide0":"none,media=cdrom"}`, "ide0", ""},
		{"agent type only", `{"agent":"type=virtio","ide0":"none,media=cdrom"}`, "ide0", "guest agent is disabled"},
		{"agent disabled", `{"agent":"enabled=0","ide0":"none,media=cdrom"}`, "ide0", "guest agent is disabled"},
		{"no agent", `{"ide0":"none,media=cdrom"}`, "ide0", "guest agent is disabled"},
		{"no slot", `{"agent":"1"}`, "ide0", ""},
		{"disk on slot", `{"agent":"1","ide0":"local:100/vm-100-disk-0.qcow2,size=8G"}`, "ide0", "a disk (local:100/vm-100-disk-0.qcow2)"},
		{"ISO in slot", `{"agent":"1","ide0":"local:iso/other.iso,media=cdrom"}`, "ide0", "a CD-ROM holding local:iso/other.iso"},
		{"host drive", `{"agent":"1","ide0":"cdrom,media=cdrom"}`, "ide0", "the host's CD drive"},
		{
			"cloud-init drive elsewhere",
			`{"agent":"1","ide0":"none,media=cdrom","ide2":"local:100/vm-100-cloudinit.qcow2,media=cdrom"}`,
			"ide0", "cloud-init drive on ide2",
		},
		{"seed on sata", `{"agent":"1","sata1":"none,media=cdrom"}`, "sata1", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var cfg vms.GetResponseData
			require.NoError(t, json.Unmarshal([]byte(tt.config), &cfg))

			err := checkVMConfig(&cfg, 100, tt.iface)
			if tt.wantErr == "" {
				assert.NoError(t, err)

				return
			}

			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}
}
