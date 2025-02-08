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
