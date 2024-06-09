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
				IOThread:   types.CustomBool(true).Pointer(),
				Size:       ds8gig,
				SSD:        types.CustomBool(true).Pointer(),
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
				IOThread:   types.CustomBool(true).Pointer(),
				Size:       ds8gig,
				SSD:        types.CustomBool(true).Pointer(),
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
			name:  "nothing matches",
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
			name:  "partially matches",
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

func TestMapCustomStorageDevices(t *testing.T) {
	t.Parallel()

	type args struct {
		resp GetResponseData
	}

	tests := []struct {
		name string
		args args
		want CustomStorageDevices
	}{
		{"no storage devices", args{GetResponseData{}}, CustomStorageDevices{}},
		{
			"ide0 storage devices",
			args{GetResponseData{IDEDevice0: &CustomStorageDevice{}}},
			map[string]*CustomStorageDevice{"ide0": {}},
		},
		{
			"multiple ide storage devices",
			args{GetResponseData{
				IDEDevice1: &CustomStorageDevice{},
				IDEDevice3: &CustomStorageDevice{},
			}},
			map[string]*CustomStorageDevice{"ide1": {}, "ide3": {}},
		},
		{
			"mixed storage devices",
			args{GetResponseData{
				IDEDevice1:       &CustomStorageDevice{},
				VirtualIODevice5: &CustomStorageDevice{},
				SATADevice0:      &CustomStorageDevice{},
				IDEDevice3:       &CustomStorageDevice{},
				SCSIDevice10:     &CustomStorageDevice{},
			}},
			map[string]*CustomStorageDevice{
				"ide1":    {},
				"virtio5": {},
				"sata0":   {},
				"ide3":    {},
				"scsi10":  {},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equalf(t, tt.want, MapCustomStorageDevices(tt.args.resp), "MapCustomStorageDevices(%v)", tt.args.resp)
		})
	}
}
