/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package vms

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnmarshalGetResponseData(t *testing.T) {
	t.Parallel()

	jsonData := fmt.Sprintf(`{
		"archive": "test",
		"ide0": "%[1]s",
		"ide1": "%[1]s",
		"ide2": "%[1]s",
		"ide3": "%[1]s",
		"virtio13": "%[1]s",
		"scsi22": "%[1]s",
		"hostpci0": "0000:81:00.2",
		"hostpci1": "host=81:00.4,pcie=0,rombar=1,x-vga=0",
		"hostpci12": "mapping=mappeddevice,pcie=0,rombar=1,x-vga=0",
		"virtiofs0":"test,cache=always,direct-io=1,expose-acl=1"
	}`, "local-lvm:vm-100-disk-0,aio=io_uring,backup=1,cache=none,discard=ignore,replicate=1,size=8G,ssd=1")

	var data GetResponseData
	err := json.Unmarshal([]byte(jsonData), &data)
	require.NoError(t, err)

	assert.Equal(t, "test", *data.BackupFile)

	assert.NotNil(t, data.StorageDevices)
	assert.Len(t, data.StorageDevices, 6)
	assert.NotNil(t, data.StorageDevices["ide0"])
	assertDevice(t, data.StorageDevices["ide0"])
	assert.NotNil(t, data.StorageDevices["ide1"])
	assertDevice(t, data.StorageDevices["ide1"])
	assert.NotNil(t, data.StorageDevices["ide2"])
	assertDevice(t, data.StorageDevices["ide2"])
	assert.NotNil(t, data.StorageDevices["ide3"])
	assertDevice(t, data.StorageDevices["ide3"])
	assert.NotNil(t, data.StorageDevices["virtio13"])
	assertDevice(t, data.StorageDevices["virtio13"])
	assert.NotNil(t, data.StorageDevices["scsi22"])
	assertDevice(t, data.StorageDevices["scsi22"])

	assert.NotNil(t, data.PCIDevices)
	assert.Len(t, data.PCIDevices, 3)
	assert.NotNil(t, data.PCIDevices["hostpci0"])
	assert.NotNil(t, data.PCIDevices["hostpci1"])
	assert.NotNil(t, data.PCIDevices["hostpci12"])

	assert.NotNil(t, data.VirtiofsShares)
	assert.Len(t, data.VirtiofsShares, 1)
	assert.Equal(t, "always", *data.VirtiofsShares["virtiofs0"].Cache)
}

func assertDevice(t *testing.T, dev *CustomStorageDevice) {
	t.Helper()

	assert.Equal(t, "io_uring", *dev.AIO)
	assert.True(t, bool(*dev.Backup))
	assert.Equal(t, "none", *dev.Cache)
	assert.Equal(t, "ignore", *dev.Discard)
	assert.Equal(t, "local-lvm:vm-100-disk-0", dev.FileVolume)
	assert.True(t, bool(*dev.Replicate))
	assert.Equal(t, "8G", dev.Size.String())
	assert.True(t, bool(*dev.SSD))
}
