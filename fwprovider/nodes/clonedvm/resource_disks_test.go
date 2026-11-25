/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package clonedvm

import (
	"net/url"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/vms"
)

func nullDiskModel() DiskModel {
	return DiskModel{
		File:        types.StringNull(),
		DatastoreID: types.StringNull(),
		SizeGB:      types.Int64Null(),
		Format:      types.StringNull(),
		AIO:         types.StringNull(),
		Backup:      types.BoolNull(),
		Discard:     types.StringNull(),
		Cache:       types.StringNull(),
		IOThread:    types.BoolNull(),
		Replicate:   types.BoolNull(),
		Serial:      types.StringNull(),
		SSD:         types.BoolNull(),
		ImportFrom:  types.StringNull(),
		Media:       types.StringNull(),
	}
}

func encodeDisk(t *testing.T, body *vms.UpdateRequestBody, slot string) string {
	t.Helper()

	device := body.CustomStorageDevices[slot]
	require.NotNil(t, device, "no %s in the update body", slot)

	values := url.Values{}
	require.NoError(t, device.EncodeValues(slot, &values))

	return values.Get(slot)
}

func TestApplyDisksKeepsInheritedOptions(t *testing.T) {
	t.Parallel()

	current := &vms.CustomStorageDevice{}
	require.NoError(t, current.UnmarshalJSON([]byte(
		`"local-lvm:vm-102-disk-0,aio=io_uring,backup=1,cache=writeback,discard=on,iothread=0,replicate=1,size=8G,ssd=1"`,
	)))

	currentConfig := &vms.GetResponseData{
		StorageDevices: vms.CustomStorageDevices{"scsi0": current},
	}

	cfg := nullDiskModel()
	cfg.SizeGB = types.Int64Value(16)

	body := &vms.UpdateRequestBody{}
	diags := diag.Diagnostics{}

	resizes := applyDisks(map[string]DiskModel{"scsi0": cfg}, currentConfig, body, &diags)
	require.False(t, diags.HasError(), diags)

	require.Len(t, resizes, 1)
	assert.Equal(t, "scsi0", resizes[0].Disk)
	assert.Equal(t, int64(16), resizes[0].Size.InGigabytes())

	encoded := encodeDisk(t, body, "scsi0")
	for _, opt := range []string{"local-lvm:vm-102-disk-0", "aio=io_uring", "backup=1", "cache=writeback", "discard=on", "ssd=1"} {
		assert.Contains(t, encoded, opt)
	}
}

func TestApplyDisksConfiguredOptionOverridesInherited(t *testing.T) {
	t.Parallel()

	current := &vms.CustomStorageDevice{}
	require.NoError(t, current.UnmarshalJSON([]byte(
		`"local-lvm:vm-102-disk-0,cache=writeback,discard=on,size=8G"`,
	)))

	currentConfig := &vms.GetResponseData{
		StorageDevices: vms.CustomStorageDevices{"scsi0": current},
	}

	cfg := nullDiskModel()
	cfg.Cache = types.StringValue("none")

	body := &vms.UpdateRequestBody{}
	diags := diag.Diagnostics{}

	resizes := applyDisks(map[string]DiskModel{"scsi0": cfg}, currentConfig, body, &diags)
	require.False(t, diags.HasError(), diags)
	assert.Empty(t, resizes)

	encoded := encodeDisk(t, body, "scsi0")
	assert.Contains(t, encoded, "cache=none")
	assert.NotContains(t, encoded, "cache=writeback")
	assert.Contains(t, encoded, "discard=on")
}

func TestApplyDisksSSDOnInheritedDisk(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		slot     string
		line     string
		ssd      bool
		contains string
		excludes string
	}{
		{
			name:     "ssd=false overrides the template on scsi",
			slot:     "scsi0",
			line:     `"local-lvm:vm-102-disk-0,discard=on,size=8G,ssd=1"`,
			ssd:      false,
			contains: "ssd=0",
			excludes: "ssd=1",
		},
		{
			name:     "ssd=true is set on an inherited sata disk",
			slot:     "sata0",
			line:     `"local-lvm:vm-102-disk-0,size=8G,ssd=0"`,
			ssd:      true,
			contains: "ssd=1",
			excludes: "ssd=0",
		},
		{
			name:     "ssd is not sent on an inherited virtio disk",
			slot:     "virtio0",
			line:     `"local-lvm:vm-102-disk-0,discard=on,size=8G"`,
			ssd:      true,
			contains: "discard=on",
			excludes: "ssd=",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			current := &vms.CustomStorageDevice{}
			require.NoError(t, current.UnmarshalJSON([]byte(tt.line)))

			currentConfig := &vms.GetResponseData{
				StorageDevices: vms.CustomStorageDevices{tt.slot: current},
			}

			cfg := nullDiskModel()
			cfg.SSD = types.BoolValue(tt.ssd)

			body := &vms.UpdateRequestBody{}
			diags := diag.Diagnostics{}

			applyDisks(map[string]DiskModel{tt.slot: cfg}, currentConfig, body, &diags)
			require.False(t, diags.HasError(), diags)

			encoded := encodeDisk(t, body, tt.slot)
			assert.Contains(t, encoded, tt.contains)
			assert.NotContains(t, encoded, tt.excludes)
		})
	}
}

func TestApplyDisksDoesNotSendDerivedFormat(t *testing.T) {
	t.Parallel()

	newCurrentConfig := func(t *testing.T) *vms.GetResponseData {
		t.Helper()

		current := &vms.CustomStorageDevice{}
		require.NoError(t, current.UnmarshalJSON([]byte(
			`"local:102/vm-102-disk-0.qcow2,aio=io_uring,size=1G,ssd=0"`,
		)))
		require.NotNil(t, current.Format, "format should be derived from the file extension")

		return &vms.GetResponseData{
			StorageDevices: vms.CustomStorageDevices{"scsi0": current},
		}
	}

	t.Run("format not configured", func(t *testing.T) {
		t.Parallel()

		cfg := nullDiskModel()
		cfg.SizeGB = types.Int64Value(2)

		body := &vms.UpdateRequestBody{}
		diags := diag.Diagnostics{}

		applyDisks(map[string]DiskModel{"scsi0": cfg}, newCurrentConfig(t), body, &diags)
		require.False(t, diags.HasError(), diags)

		encoded := encodeDisk(t, body, "scsi0")
		assert.NotContains(t, encoded, "format=")
		assert.Contains(t, encoded, "aio=io_uring")
	})

	t.Run("format configured", func(t *testing.T) {
		t.Parallel()

		cfg := nullDiskModel()
		cfg.Format = types.StringValue("raw")

		body := &vms.UpdateRequestBody{}
		diags := diag.Diagnostics{}

		applyDisks(map[string]DiskModel{"scsi0": cfg}, newCurrentConfig(t), body, &diags)
		require.False(t, diags.HasError(), diags)

		assert.Contains(t, encodeDisk(t, body, "scsi0"), "format=raw")
	})
}
