/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package vms

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

func TestCustomStorageDevice_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	ds8gig := types.DiskSizeFromGigabytes(8)
	tests := []struct {
		name    string
		line    string
		want    *CustomStorageDevice
		wantErr bool
	}{
		{
			name: "simple volume",
			line: `"local-lvm:vm-2041-disk-0,discard=on,ssd=1,iothread=1,size=8G,cache=writeback"`,
			want: &CustomStorageDevice{
				Cache:      types.StrPtr("writeback"),
				Discard:    types.StrPtr("on"),
				Enabled:    true,
				FileVolume: "local-lvm:vm-2041-disk-0",
				IOThread:   types.BoolPtr(true),
				Size:       &ds8gig,
				SSD:        types.BoolPtr(true),
			},
		},
		{
			name: "raw volume type",
			line: `"nfs:2041/vm-2041-disk-0.raw,discard=ignore,ssd=1,iothread=1,size=8G"`,
			want: &CustomStorageDevice{
				Discard:    types.StrPtr("ignore"),
				Enabled:    true,
				FileVolume: "nfs:2041/vm-2041-disk-0.raw",
				Format:     types.StrPtr("raw"),
				IOThread:   types.BoolPtr(true),
				Size:       &ds8gig,
				SSD:        types.BoolPtr(true),
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := &CustomStorageDevice{}
			if err := r.UnmarshalJSON([]byte(tt.line)); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
			require.Equal(t, tt.want, r)
		})
	}
}

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
				DeviceIDs:  &[]string{"0000:81:00.2"},
				MDev:       nil,
				PCIExpress: types.BoolPtr(false),
				ROMBAR:     types.BoolPtr(true),
				ROMFile:    nil,
				XVGA:       types.BoolPtr(false),
			},
		},
		{
			name: "pci device with more details",
			line: `"host=81:00.4,pcie=0,rombar=1,x-vga=0"`,
			want: &CustomPCIDevice{
				DeviceIDs:  &[]string{"81:00.4"},
				MDev:       nil,
				PCIExpress: types.BoolPtr(false),
				ROMBAR:     types.BoolPtr(true),
				ROMFile:    nil,
				XVGA:       types.BoolPtr(false),
			},
		},
		{
			name: "pci device with mapping",
			line: `"mapping=mappeddevice,pcie=0,rombar=1,x-vga=0"`,
			want: &CustomPCIDevice{
				DeviceIDs:  nil,
				Mapping:    types.StrPtr("mappeddevice"),
				MDev:       nil,
				PCIExpress: types.BoolPtr(false),
				ROMBAR:     types.BoolPtr(true),
				ROMFile:    nil,
				XVGA:       types.BoolPtr(false),
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := &CustomPCIDevice{}
			if err := r.UnmarshalJSON([]byte(tt.line)); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

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
				HostDevice: types.StrPtr("0000:81"),
			},
		},
		{
			name: "usb device with more details",
			line: `"host=81:00,usb3=0"`,
			want: &CustomUSBDevice{
				HostDevice: types.StrPtr("81:00"),
				USB3:       types.BoolPtr(false),
			},
		},
		{
			name: "usb device with mapping",
			line: `"mapping=mappeddevice,usb=0"`,
			want: &CustomUSBDevice{
				HostDevice: nil,
				Mapping:    types.StrPtr("mappeddevice"),
				USB3:       types.BoolPtr(false),
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := &CustomUSBDevice{}
			if err := r.UnmarshalJSON([]byte(tt.line)); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
