/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package vms

import (
	"testing"

	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

func TestCustomUSBDevice_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		line    string
		want    *CustomUSBDevice
		wantErr bool
	}{
		{
			name: "id only usb device",
			line: `"host=0000:81"`,
			want: &CustomUSBDevice{
				HostDevice: new("0000:81"),
			},
		},
		{
			name: "usb device with more details",
			line: `"host=81:00,usb3=0"`,
			want: &CustomUSBDevice{
				HostDevice: new("81:00"),
				USB3:       types.CustomBool(false).Pointer(),
			},
		},
		{
			name: "usb device with mapping",
			line: `"mapping=mappeddevice,usb=0"`,
			want: &CustomUSBDevice{
				HostDevice: nil,
				Mapping:    new("mappeddevice"),
				USB3:       types.CustomBool(false).Pointer(),
			},
		},
		{
			name: "spice usb device",
			line: `"spice"`,
			want: &CustomUSBDevice{
				HostDevice: new("spice"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := &CustomUSBDevice{}
			if err := r.UnmarshalJSON([]byte(tt.line)); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
