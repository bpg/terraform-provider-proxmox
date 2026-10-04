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

func TestValidateInitializationInterface(t *testing.T) {
	t.Parallel()

	str := cty.StringVal
	diskType := cty.Object(map[string]cty.Type{"interface": cty.String})
	initType := cty.Object(map[string]cty.Type{"interface": cty.String})

	disks := func(ifaces ...cty.Value) cty.Value {
		vals := make([]cty.Value, len(ifaces))
		for i, v := range ifaces {
			vals[i] = cty.ObjectVal(map[string]cty.Value{"interface": v})
		}

		return cty.ListVal(vals)
	}

	initBlocks := func(iface cty.Value) cty.Value {
		return cty.ListVal([]cty.Value{cty.ObjectVal(map[string]cty.Value{"interface": iface})})
	}

	config := func(disksVal, initVal cty.Value) cty.Value {
		return cty.ObjectVal(map[string]cty.Value{"disk": disksVal, "initialization": initVal})
	}

	noInit := cty.ListValEmpty(initType)
	collides := `collides with disk.`

	tests := []struct {
		name     string
		config   cty.Value
		isCreate bool
		wantErr  string
	}{
		{
			"no collision",
			config(disks(str("scsi0"), str("scsi1")), initBlocks(str("ide2"))),
			true, "",
		},
		{
			"disk collision",
			config(disks(str("scsi0"), str("ide2")), initBlocks(str("ide2"))),
			true, `initialization.0.interface "ide2" collides with disk.1.interface`,
		},
		{
			"disk collision on update",
			config(disks(str("sata1")), initBlocks(str("sata1"))),
			false, collides,
		},
		{
			"unset interface defaults to ide2 on create",
			config(disks(str("ide2")), initBlocks(cty.NullVal(cty.String))),
			true, `initialization.0.interface "ide2" collides with disk.0.interface`,
		},
		{
			"empty interface defaults to ide2 on create",
			config(disks(str("ide2")), initBlocks(str(""))),
			true, collides,
		},
		{
			"unset interface, no ide2 disk, on create",
			config(disks(str("scsi0")), initBlocks(cty.NullVal(cty.String))),
			true, "",
		},
		{
			"unset interface is skipped on update",
			config(disks(str("ide2")), initBlocks(cty.NullVal(cty.String))),
			false, "",
		},
		{
			"unknown initialization interface",
			config(disks(str("ide2")), initBlocks(cty.UnknownVal(cty.String))),
			true, "",
		},
		{
			"unknown disk interface",
			config(disks(cty.UnknownVal(cty.String)), initBlocks(str("ide2"))),
			true, "",
		},
		{
			"unknown disk list",
			config(cty.UnknownVal(cty.List(diskType)), initBlocks(str("ide2"))),
			true, "",
		},
		{
			"null disk interface",
			config(disks(cty.NullVal(cty.String)), initBlocks(str("ide2"))),
			true, "",
		},
		{
			"no initialization block",
			config(disks(str("ide2")), noInit),
			true, "",
		},
		{
			"marked initialization interface",
			config(disks(str("ide2")), initBlocks(str("ide2").Mark("sensitive"))),
			true, collides,
		},
		{
			"marked disk interface",
			config(disks(str("ide2").Mark("sensitive")), initBlocks(str("ide2"))),
			true, collides,
		},
		{"null config", cty.NullVal(cty.EmptyObject), true, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := validateInitializationInterface(tt.config, tt.isCreate)
			if tt.wantErr == "" {
				require.NoError(t, err)
				return
			}

			require.ErrorContains(t, err, tt.wantErr)
		})
	}
}
