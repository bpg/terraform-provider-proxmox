/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package vms

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/proxmox/helpers/ptr"
	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

func TestCustomPCIDevice_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		line    string
		want    *CustomPCIDevice
		wantErr bool
	}{
		{
			name: "id only pci device",
			line: `"0000:81:00.2"`,
			want: &CustomPCIDevice{
				DeviceIDs: &[]string{"0000:81:00.2"},
			},
		},
		{
			name: "pci device with more details",
			line: `"host=81:00.4,pcie=0,rombar=1,x-vga=0"`,
			want: &CustomPCIDevice{
				DeviceIDs:  &[]string{"81:00.4"},
				MDev:       nil,
				PCIExpress: types.CustomBool(false).Pointer(),
				ROMBAR:     types.CustomBool(true).Pointer(),
				ROMFile:    nil,
				XVGA:       types.CustomBool(false).Pointer(),
			},
		},
		{
			name: "pci device with mapping",
			line: `"mapping=mappeddevice,pcie=0,rombar=1,x-vga=0"`,
			want: &CustomPCIDevice{
				DeviceIDs:  nil,
				Mapping:    ptr.Ptr("mappeddevice"),
				MDev:       nil,
				PCIExpress: types.CustomBool(false).Pointer(),
				ROMBAR:     types.CustomBool(true).Pointer(),
				ROMFile:    nil,
				XVGA:       types.CustomBool(false).Pointer(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := &CustomPCIDevice{}
			if err := r.UnmarshalJSON([]byte(tt.line)); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}

			require.Equal(t, tt.want, r)
		})
	}
}
