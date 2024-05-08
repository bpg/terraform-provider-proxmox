/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package vms

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/proxmox/helpers/ptr"
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
				Cache:      ptr.Ptr("writeback"),
				Discard:    ptr.Ptr("on"),
				Enabled:    true,
				FileVolume: "local-lvm:vm-2041-disk-0",
				IOThread:   types.CustomBoolPtr(true),
				Size:       ds8gig,
				SSD:        types.CustomBoolPtr(true),
			},
		},
		{
			name: "raw volume type",
			line: `"nfs:2041/vm-2041-disk-0.raw,discard=ignore,ssd=1,iothread=1,size=8G"`,
			want: &CustomStorageDevice{
				Discard:    ptr.Ptr("ignore"),
				Enabled:    true,
				FileVolume: "nfs:2041/vm-2041-disk-0.raw",
				Format:     ptr.Ptr("raw"),
				IOThread:   types.CustomBoolPtr(true),
				Size:       ds8gig,
				SSD:        types.CustomBoolPtr(true),
			},
		},
	}

	for _, tt := range tests {
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

func TestCustomStorageDevice_IsCloudInitDrive(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		device CustomStorageDevice
		want   bool
	}{
		{
			name: "simple volume",
			device: CustomStorageDevice{
				FileVolume: "local-lvm:vm-131-disk-0",
			},
			want: false,
		}, {
			name: "on directory storage",
			device: CustomStorageDevice{
				Media:      ptr.Ptr("cdrom"),
				FileVolume: "local:131/vm-131-cloudinit.qcow2",
			},
			want: true,
		}, {
			name: "on block storage",
			device: CustomStorageDevice{
				Media:      ptr.Ptr("cdrom"),
				FileVolume: "local-lvm:vm-131-cloudinit",
			},
			want: true,
		}, {
			name: "wrong VM ID",
			device: CustomStorageDevice{
				Media:      ptr.Ptr("cdrom"),
				FileVolume: "local-lvm:vm-123-cloudinit",
			},
			want: false,
		}, {
			name: "not a cdrom",
			device: CustomStorageDevice{
				FileVolume: "local-lvm:vm-123-cloudinit",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := tt.device.IsCloudInitDrive(131)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCustomStorageDevice_StorageInterface(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		device CustomStorageDevice
		want   string
	}{
		{
			name: "virtio0",
			device: CustomStorageDevice{
				Interface: ptr.Ptr("virtio0"),
			},
			want: "virtio",
		}, {
			name: "scsi13",
			device: CustomStorageDevice{
				Interface: ptr.Ptr("scsi13"),
			},
			want: "scsi",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := tt.device.StorageInterface()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCustomStorageDevices_ByStorageInterface(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		iface   string
		devices CustomStorageDevices
		want    CustomStorageDevices
	}{
		{
			name:    "empty",
			iface:   "virtio",
			devices: CustomStorageDevices{},
			want:    CustomStorageDevices{},
		},
		{
			name:  "not in the list",
			iface: "sata",
			devices: CustomStorageDevices{
				"virtio0": &CustomStorageDevice{
					Interface: ptr.Ptr("virtio0"),
				},
				"scsi13": &CustomStorageDevice{
					Interface: ptr.Ptr("scsi13"),
				},
			},
			want: CustomStorageDevices{},
		},
		{
			name:  "not in the list",
			iface: "virtio",
			devices: CustomStorageDevices{
				"virtio0": &CustomStorageDevice{
					Interface: ptr.Ptr("virtio0"),
				},
				"scsi13": &CustomStorageDevice{
					Interface: ptr.Ptr("scsi13"),
				},
				"virtio1": &CustomStorageDevice{
					Interface: ptr.Ptr("virtio1"),
				},
			},
			want: CustomStorageDevices{
				"virtio0": &CustomStorageDevice{
					Interface: ptr.Ptr("virtio0"),
				},
				"virtio1": &CustomStorageDevice{
					Interface: ptr.Ptr("virtio1"),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := tt.devices.ByStorageInterface(tt.iface)
			assert.Equal(t, tt.want, got)
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
				DeviceIDs: &[]string{"0000:81:00.2"},
			},
		},
		{
			name: "pci device with more details",
			line: `"host=81:00.4,pcie=0,rombar=1,x-vga=0"`,
			want: &CustomPCIDevice{
				DeviceIDs:  &[]string{"81:00.4"},
				MDev:       nil,
				PCIExpress: types.CustomBoolPtr(false),
				ROMBAR:     types.CustomBoolPtr(true),
				ROMFile:    nil,
				XVGA:       types.CustomBoolPtr(false),
			},
		},
		{
			name: "pci device with mapping",
			line: `"mapping=mappeddevice,pcie=0,rombar=1,x-vga=0"`,
			want: &CustomPCIDevice{
				DeviceIDs:  nil,
				Mapping:    ptr.Ptr("mappeddevice"),
				MDev:       nil,
				PCIExpress: types.CustomBoolPtr(false),
				ROMBAR:     types.CustomBoolPtr(true),
				ROMFile:    nil,
				XVGA:       types.CustomBoolPtr(false),
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

func TestCustomNUMADevice_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		line    string
		want    *CustomNUMADevice
		wantErr bool
	}{
		{
			name: "numa device all options",
			line: `"cpus=1-2;3-4,hostnodes=1-2,memory=1024,policy=preferred"`,
			want: &CustomNUMADevice{
				CPUIDs:        []string{"1-2", "3-4"},
				HostNodeNames: &[]string{"1-2"},
				Memory:        ptr.Ptr(1024),
				Policy:        ptr.Ptr("preferred"),
			},
		},
		{
			name: "numa device cpus/memory only",
			line: `"cpus=1-2,memory=1024"`,
			want: &CustomNUMADevice{
				CPUIDs: []string{"1-2"},
				Memory: ptr.Ptr(1024),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := &CustomNUMADevice{}
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
				HostDevice: ptr.Ptr("0000:81"),
			},
		},
		{
			name: "usb device with more details",
			line: `"host=81:00,usb3=0"`,
			want: &CustomUSBDevice{
				HostDevice: ptr.Ptr("81:00"),
				USB3:       types.CustomBoolPtr(false),
			},
		},
		{
			name: "usb device with mapping",
			line: `"mapping=mappeddevice,usb=0"`,
			want: &CustomUSBDevice{
				HostDevice: nil,
				Mapping:    ptr.Ptr("mappeddevice"),
				USB3:       types.CustomBoolPtr(false),
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
