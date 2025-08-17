/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package disk

import (
	"context"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/proxmox"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/vms"
	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

// TestDiskOrderingDeterministic tests that disk ordering is deterministic
// and preserves the order from currentDiskList. This test addresses the issue where
// disk interfaces would randomly reorder on each terraform apply due to
// Go's non-deterministic map iteration.
func TestDiskOrderingDeterministic(t *testing.T) {
	t.Parallel()

	// Create test schema
	diskSchema := Schema()

	// Create resource data with multiple disks in a specific order that could be affected by map iteration
	currentDiskList := []interface{}{
		map[string]interface{}{
			mkDiskInterface:   "scsi1", // Intentionally put scsi1 first
			mkDiskDatastoreID: "local",
			mkDiskSize:        150,
			mkDiskSpeed:       []interface{}{},
		},
		map[string]interface{}{
			mkDiskInterface:   "scsi0", // Then scsi0 second
			mkDiskDatastoreID: "local",
			mkDiskSize:        50,
			mkDiskSpeed:       []interface{}{},
		},
	}

	// Mock VM disk data from API that matches the currentDiskList
	// Set Format to avoid API calls in the Read function
	qcow2Format := "qcow2"
	diskDeviceObjects := vms.CustomStorageDevices{
		"scsi0": &vms.CustomStorageDevice{
			FileVolume: "local:50",
			Size:       types.DiskSizeFromGigabytes(50),
			Format:     &qcow2Format,
		},
		"scsi1": &vms.CustomStorageDevice{
			FileVolume: "local:150",
			Size:       types.DiskSizeFromGigabytes(150),
			Format:     &qcow2Format,
		},
	}

	// Run the ordering logic multiple times to ensure deterministic results
	const iterations = 10

	results := make([][]interface{}, 0, iterations)

	for range iterations {
		// Create a new resource data for each iteration
		resourceData := schema.TestResourceDataRaw(t, diskSchema, map[string]interface{}{
			MkDisk: currentDiskList,
		})

		// Call the Read function which contains our fixed ordering logic
		ctx := context.Background()
		vmID := 100 // Test VM ID

		var client proxmox.Client = nil

		diags := Read(ctx, resourceData, diskDeviceObjects, vmID, client, "test-node", false)
		require.Empty(t, diags, "Read should not return any diagnostics")

		// Get the resulting disk list
		diskList := resourceData.Get(MkDisk).([]interface{})
		results = append(results, diskList)
	}

	// Verify all results are identical (deterministic ordering)
	expectedResult := results[0]
	for i, result := range results {
		require.True(t, reflect.DeepEqual(expectedResult, result),
			"Disk ordering should be deterministic across multiple calls. Iteration %d differs from first result", i)
	}

	// Verify the specific ordering - should preserve currentDiskList order: scsi1, then scsi0
	require.Len(t, expectedResult, 2, "Should have 2 disks")

	disk0 := expectedResult[0].(map[string]interface{})
	disk1 := expectedResult[1].(map[string]interface{})

	require.Equal(t, "scsi1", disk0[mkDiskInterface], "First disk should be scsi1 (as in currentDiskList)")
	require.Equal(t, 150, disk0[mkDiskSize], "First disk should have size 150")

	require.Equal(t, "scsi0", disk1[mkDiskInterface], "Second disk should be scsi0 (as in currentDiskList)")
	require.Equal(t, 50, disk1[mkDiskSize], "Second disk should have size 50")
}

// TestDiskOrderingVariousInterfaces tests deterministic ordering with various disk interfaces,
// ensuring the order from currentDiskList is preserved.
func TestDiskOrderingVariousInterfaces(t *testing.T) {
	t.Parallel()

	diskSchema := Schema()

	// Test with various interface types in random order
	currentDiskList := []interface{}{
		map[string]interface{}{
			mkDiskInterface:   "virtio2",
			mkDiskDatastoreID: "local",
			mkDiskSize:        30,
			mkDiskSpeed:       []interface{}{},
		},
		map[string]interface{}{
			mkDiskInterface:   "scsi0",
			mkDiskDatastoreID: "local",
			mkDiskSize:        10,
			mkDiskSpeed:       []interface{}{},
		},
		map[string]interface{}{
			mkDiskInterface:   "sata1",
			mkDiskDatastoreID: "local",
			mkDiskSize:        20,
			mkDiskSpeed:       []interface{}{},
		},
		map[string]interface{}{
			mkDiskInterface:   "virtio0",
			mkDiskDatastoreID: "local",
			mkDiskSize:        40,
			mkDiskSpeed:       []interface{}{},
		},
	}

	qcow2Format := "qcow2"
	diskDeviceObjects := vms.CustomStorageDevices{
		"scsi0":   &vms.CustomStorageDevice{FileVolume: "local:10", Size: types.DiskSizeFromGigabytes(10), Format: &qcow2Format},
		"sata1":   &vms.CustomStorageDevice{FileVolume: "local:20", Size: types.DiskSizeFromGigabytes(20), Format: &qcow2Format},
		"virtio2": &vms.CustomStorageDevice{FileVolume: "local:30", Size: types.DiskSizeFromGigabytes(30), Format: &qcow2Format},
		"virtio0": &vms.CustomStorageDevice{FileVolume: "local:40", Size: types.DiskSizeFromGigabytes(40), Format: &qcow2Format},
	}

	// Test multiple iterations
	const iterations = 5

	results := make([][]interface{}, 0, iterations)

	for range iterations {
		resourceData := schema.TestResourceDataRaw(t, diskSchema, map[string]interface{}{
			MkDisk: currentDiskList,
		})

		ctx := context.Background()
		vmID := 100

		var client proxmox.Client = nil

		diags := Read(ctx, resourceData, diskDeviceObjects, vmID, client, "test-node", false)
		require.Empty(t, diags)

		diskList := resourceData.Get(MkDisk).([]interface{})
		results = append(results, diskList)
	}

	// Verify deterministic ordering
	expectedResult := results[0]
	for i, result := range results {
		require.True(t, reflect.DeepEqual(expectedResult, result),
			"Disk ordering should be deterministic for various interfaces. Iteration %d differs", i)
	}

	// Verify ordering preserves currentDiskList order: virtio2, scsi0, sata1, virtio0
	require.Len(t, expectedResult, 4)

	expectedOrder := []string{"virtio2", "scsi0", "sata1", "virtio0"}
	for i, expectedInterface := range expectedOrder {
		disk := expectedResult[i].(map[string]interface{})
		require.Equal(t, expectedInterface, disk[mkDiskInterface],
			"Disk at position %d should be %s (as in currentDiskList)", i, expectedInterface)
	}
}
