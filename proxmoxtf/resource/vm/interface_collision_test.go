/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package resource

import (
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/stretchr/testify/require"
)

func diskBlock(iface cty.Value) cty.Value {
	return cty.ObjectVal(map[string]cty.Value{"interface": iface})
}

func vmConfig(disks, efi, init cty.Value) cty.Value {
	return cty.ObjectVal(map[string]cty.Value{
		"disk":           disks,
		"efi_disk":       efi,
		"initialization": init,
	})
}

func initBlocks(iface cty.Value) cty.Value {
	return cty.ListVal([]cty.Value{cty.ObjectVal(map[string]cty.Value{"interface": iface})})
}

func TestValidateInitializationInterface(t *testing.T) {
	t.Parallel()

	str := cty.StringVal
	diskType := cty.Object(map[string]cty.Type{"interface": cty.String})
	disks := func(ifaces ...string) cty.Value {
		vals := make([]cty.Value, len(ifaces))
		for i, v := range ifaces {
			vals[i] = diskBlock(str(v))
		}

		return cty.ListVal(vals)
	}
	noEFI := cty.ListValEmpty(cty.Object(map[string]cty.Type{"type": cty.String}))
	efi := cty.ListVal([]cty.Value{cty.ObjectVal(map[string]cty.Value{"type": str("4m")})})

	tests := []struct {
		name    string
		config  cty.Value
		wantErr string
	}{
		{"no collision", vmConfig(disks("scsi0", "scsi1"), noEFI, initBlocks(str("ide2"))), ""},
		{"disk collision", vmConfig(disks("scsi0", "ide2"), noEFI, initBlocks(str("ide2"))),
			`initialization.0.interface "ide2" collides with disk.1.interface`},
		{"efi collision", vmConfig(disks("scsi0"), efi, initBlocks(str("efidisk0"))),
			`collides with the efi_disk device`},
		{"efidisk0 without efi block", vmConfig(disks("scsi0"), noEFI, initBlocks(str("efidisk0"))), ""},
		{"empty initialization interface", vmConfig(disks("ide2"), noEFI, initBlocks(str(""))), ""},
		{"unknown initialization interface",
			vmConfig(disks("ide2"), noEFI, initBlocks(cty.UnknownVal(cty.String))), ""},
		{"unknown disk interface",
			vmConfig(cty.ListVal([]cty.Value{diskBlock(cty.UnknownVal(cty.String))}), noEFI, initBlocks(str("ide2"))), ""},
		{"unknown disk list", vmConfig(cty.UnknownVal(cty.List(diskType)), noEFI, initBlocks(str("ide2"))), ""},
		{"null disk interface", vmConfig(cty.ListVal([]cty.Value{diskBlock(cty.NullVal(cty.String))}), noEFI,
			initBlocks(str("ide2"))), ""},
		{"no initialization block", vmConfig(disks("ide2"), noEFI,
			cty.ListValEmpty(cty.Object(map[string]cty.Type{"interface": cty.String}))), ""},
		{"marked initialization interface", vmConfig(disks("ide2"), noEFI, initBlocks(str("ide2").Mark("sensitive"))),
			`collides with disk.0.interface`},
		{"marked disk interface", vmConfig(cty.ListVal([]cty.Value{diskBlock(str("ide2").Mark("sensitive"))}), noEFI,
			initBlocks(str("ide2"))), `collides with disk.0.interface`},
		{"null config", cty.NullVal(cty.EmptyObject), ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := validateInitializationInterface(tt.config)
			if tt.wantErr == "" {
				require.NoError(t, err)
				return
			}

			require.ErrorContains(t, err, tt.wantErr)
		})
	}
}
