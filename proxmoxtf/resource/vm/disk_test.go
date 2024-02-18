package vm

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/vms"
	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

func TestMapStorageDevices(t *testing.T) {
	t.Parallel()

	devices := &vms.GetResponseData{
		VirtualIODevice0: &vms.CustomStorageDevice{
			Interface: types.StrPtr("virtio0"),
		},
		VirtualIODevice1: &vms.CustomStorageDevice{
			Interface: types.StrPtr("virtio1"),
			Size:      types.DiskSizeFromGigabytes(10),
		},
	}

	expected := map[string]*vms.CustomStorageDevice{
		"virtio0": {
			Interface: types.StrPtr("virtio0"),
			Size:      new(types.DiskSize),
		},
		"virtio1": {
			Interface: types.StrPtr("virtio1"),
			Size:      types.DiskSizeFromGigabytes(10),
		},
	}

	result := mapStorageDevices(devices)

	assert.Equal(t, expected, result)
}

func TestPopulateFileID(t *testing.T) {
	t.Parallel()

	devicesMap := map[string]*vms.CustomStorageDevice{
		"virtio0": {},
		"virtio1": {},
	}

	disk := []map[string]interface{}{
		{
			mkDiskInterface: "virtio0",
			mkDiskFileID:    "local:100/vm-100-disk-1.qcow2",
		},
		{
			mkDiskInterface: "virtio1",
			mkDiskFileID:    "local:100/vm-100-disk-2.qcow2",
		},
	}

	d := VM().TestResourceData()
	err := d.Set("disk", disk)
	require.NoError(t, err)

	expected := vms.CustomStorageDevices{
		"virtio0": {
			FileID: types.StrPtr("local:100/vm-100-disk-1.qcow2"),
		},
		"virtio1": {
			FileID: types.StrPtr("local:100/vm-100-disk-2.qcow2"),
		},
	}

	result := populateFileIDs(devicesMap, d)

	assert.Equal(t, expected, result)
}
