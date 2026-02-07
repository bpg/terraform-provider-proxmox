/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package disk

import (
	"context"
	"reflect"
	"slices"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/proxmox"
	"github.com/bpg/terraform-provider-proxmox/proxmox/helpers/ptr"
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
	currentDiskList := []any{
		map[string]any{
			mkDiskInterface:   "scsi1", // Intentionally put scsi1 first
			mkDiskDatastoreID: "local",
			mkDiskSize:        150,
			mkDiskSpeed:       []any{},
		},
		map[string]any{
			mkDiskInterface:   "scsi0", // Then scsi0 second
			mkDiskDatastoreID: "local",
			mkDiskSize:        50,
			mkDiskSpeed:       []any{},
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

	results := make([][]any, 0, iterations)

	for range iterations {
		// Create a new resource data for each iteration
		resourceData := schema.TestResourceDataRaw(t, diskSchema, map[string]any{
			MkDisk: currentDiskList,
		})

		// Call the Read function which contains our fixed ordering logic
		ctx := context.Background()
		vmID := 100 // Test VM ID

		var client proxmox.Client

		diags := Read(ctx, resourceData, diskDeviceObjects, vmID, client, "test-node", false)
		require.Empty(t, diags, "Read should not return any diagnostics")

		// Get the resulting disk list
		diskList := resourceData.Get(MkDisk).([]any)
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

	disk0 := expectedResult[0].(map[string]any)
	disk1 := expectedResult[1].(map[string]any)

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
	currentDiskList := []any{
		map[string]any{
			mkDiskInterface:   "virtio2",
			mkDiskDatastoreID: "local",
			mkDiskSize:        30,
			mkDiskSpeed:       []any{},
		},
		map[string]any{
			mkDiskInterface:   "scsi0",
			mkDiskDatastoreID: "local",
			mkDiskSize:        10,
			mkDiskSpeed:       []any{},
		},
		map[string]any{
			mkDiskInterface:   "sata1",
			mkDiskDatastoreID: "local",
			mkDiskSize:        20,
			mkDiskSpeed:       []any{},
		},
		map[string]any{
			mkDiskInterface:   "virtio0",
			mkDiskDatastoreID: "local",
			mkDiskSize:        40,
			mkDiskSpeed:       []any{},
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

	results := make([][]any, 0, iterations)

	for range iterations {
		resourceData := schema.TestResourceDataRaw(t, diskSchema, map[string]any{
			MkDisk: currentDiskList,
		})

		ctx := context.Background()
		vmID := 100

		var client proxmox.Client

		diags := Read(ctx, resourceData, diskDeviceObjects, vmID, client, "test-node", false)
		require.Empty(t, diags)

		diskList := resourceData.Get(MkDisk).([]any)
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
		disk := expectedResult[i].(map[string]any)
		require.Equal(t, expectedInterface, disk[mkDiskInterface],
			"Disk at position %d should be %s (as in currentDiskList)", i, expectedInterface)
	}
}

// TestDiskDevicesEqual tests the disk Equals method to ensure proper comparison.
// noinspection:GoDfaNilDereference // verifying nil receiver handling
func TestDiskDevicesEqual(t *testing.T) {
	t.Parallel()

	// Test nil cases
	var nilDisk *vms.CustomStorageDevice
	require.False(t, nilDisk.Equals(nil))
	require.False(t, nilDisk.Equals(&vms.CustomStorageDevice{}))
	require.False(t, (&vms.CustomStorageDevice{}).Equals(nil))

	// Create identical disks
	aio1 := "io_uring"
	aio2 := "io_uring"
	cache1 := "writeback"
	cache2 := "writeback"
	size1 := types.DiskSizeFromGigabytes(10)
	size2 := types.DiskSizeFromGigabytes(10)
	datastore1 := "local"
	datastore2 := "local"
	importFrom1 := "local:import/test.qcow2"
	importFrom2 := "local:import/test.qcow2"

	disk1 := &vms.CustomStorageDevice{
		AIO:         &aio1,
		Cache:       &cache1,
		Size:        size1,
		DatastoreID: &datastore1,
		ImportFrom:  &importFrom1,
	}

	disk2 := &vms.CustomStorageDevice{
		AIO:         &aio2,
		Cache:       &cache2,
		Size:        size2,
		DatastoreID: &datastore2,
		ImportFrom:  &importFrom2,
	}

	// Test identical disks
	require.True(t, disk1.Equals(disk2))

	// Test different AIO
	aio2Changed := "native"
	disk2Changed := &vms.CustomStorageDevice{
		AIO:         &aio2Changed,
		Cache:       &cache2,
		Size:        size2,
		DatastoreID: &datastore2,
	}
	require.False(t, disk1.Equals(disk2Changed))

	// Test different size
	size2Changed := types.DiskSizeFromGigabytes(20)
	disk2SizeChanged := &vms.CustomStorageDevice{
		AIO:         &aio2,
		Cache:       &cache2,
		Size:        size2Changed,
		DatastoreID: &datastore2,
	}
	require.False(t, disk1.Equals(disk2SizeChanged))

	// Test different ImportFrom values
	importFrom2Changed := "local:import/test2.qcow2"
	disk2ImportFromChanged := &vms.CustomStorageDevice{
		AIO:         &aio2,
		Cache:       &cache2,
		Size:        size2,
		ImportFrom:  &importFrom2Changed,
		DatastoreID: &datastore2,
	}
	require.False(t, disk1.Equals(disk2ImportFromChanged))
}

// TestDiskUpdateSkipsUnchangedDisks tests that the Update function only updates changed disks.
func TestDiskUpdateSkipsUnchangedDisks(t *testing.T) {
	t.Parallel()

	// Mock resource data
	diskSchema := Schema()

	var err error

	resourceData := schema.TestResourceDataRaw(t, diskSchema, map[string]any{
		MkDisk: []any{
			map[string]any{
				mkDiskInterface:   "scsi0",
				mkDiskDatastoreID: "local",
				mkDiskSize:        10,
				mkDiskImportFrom:  "local:iso/disk.qcow2",
				mkDiskSpeed:       []any{},
			},
			map[string]any{
				mkDiskInterface:   "scsi1",
				mkDiskDatastoreID: "local",
				mkDiskSize:        20,
				mkDiskSpeed:       []any{},
			},
		},
	})

	// Mark that the disk configuration has changed (terraform detected a change)
	resourceData.MarkNewResource()

	// Create current disks (what Proxmox currently has)
	importFrom := "local:iso/disk.qcow2"
	datastoreID := "local"
	currentDisks := vms.CustomStorageDevices{
		"scsi0": &vms.CustomStorageDevice{
			Size:        types.DiskSizeFromGigabytes(10),
			DatastoreID: &datastoreID,
			ImportFrom:  &importFrom,
		},
		"scsi1": &vms.CustomStorageDevice{
			Size:        types.DiskSizeFromGigabytes(5), // This is different (current=5, plan=20)
			DatastoreID: &datastoreID,
		},
	}

	// Create plan disks (what terraform wants)
	planDisks := vms.CustomStorageDevices{
		"scsi0": &vms.CustomStorageDevice{
			Size:        types.DiskSizeFromGigabytes(10), // Same as current
			DatastoreID: &datastoreID,
			ImportFrom:  &importFrom,
		},
		"scsi1": &vms.CustomStorageDevice{
			Size:        types.DiskSizeFromGigabytes(20), // Different from current (5 -> 20)
			DatastoreID: &datastoreID,
		},
	}

	// Mock update body to capture what gets sent to the API
	updateBody := &vms.UpdateRequestBody{}

	// Mock client (not used in this test, but required by function signature)
	var client proxmox.Client = nil

	ctx := context.Background()
	vmID := 100
	nodeName := "test-node"

	// Force HasChange to return true by setting old and new values
	err = resourceData.Set(MkDisk, []any{
		map[string]any{
			mkDiskInterface:   "scsi1",
			mkDiskDatastoreID: "local",
			mkDiskSize:        5, // Old size
			mkDiskSpeed:       []any{},
		},
	})
	require.NoError(t, err)

	err = resourceData.Set(MkDisk, []any{
		map[string]any{
			mkDiskInterface:   "scsi1",
			mkDiskDatastoreID: "local",
			mkDiskSize:        20, // New size
			mkDiskSpeed:       []any{},
		},
	})
	require.NoError(t, err)

	// Call the Update function
	_, err = Update(ctx, client, nodeName, vmID, resourceData, planDisks, currentDisks, updateBody)
	require.NoError(t, err)

	// The update body should only contain scsi1, not scsi0
	// This prevents the "can't unplug bootdisk 'scsi0'" error
	require.Contains(t, updateBody.CustomStorageDevices, "scsi1", "Update body should contain the changed disk scsi1")
	require.NotContains(t, updateBody.CustomStorageDevices, "scsi0", "Update body should not contain the unchanged disk scsi0")
}

// TestImportFromDiskNotReimportedOnSizeChange tests issue #2385:
// when a disk with import_from is resized in Proxmox GUI, terraform should NOT
// attempt to re-import the disk (which would fail with "cannot shrink" error).
func TestImportFromDiskNotReimportedOnSizeChange(t *testing.T) {
	t.Parallel()

	diskSchema := Schema()

	// Terraform config has a disk with import_from and size=20GB
	resourceData := schema.TestResourceDataRaw(t, diskSchema, map[string]any{
		MkDisk: []any{
			map[string]any{
				mkDiskInterface:   "scsi0",
				mkDiskDatastoreID: "nfs-v3-tmp-01",
				mkDiskSize:        20,
				mkDiskImportFrom:  "nfs-v3-tmp-01:iso/ubuntu-22.04-cloud.qcow2",
				mkDiskSpeed:       []any{},
			},
		},
	})

	resourceData.MarkNewResource()

	// Proxmox currently has the disk at 30GB (resized via GUI)
	// Note: PVE does NOT return import_from for existing disks
	datastoreID := "nfs-v3-tmp-01"
	currentDisks := vms.CustomStorageDevices{
		"scsi0": &vms.CustomStorageDevice{
			Size:        types.DiskSizeFromGigabytes(30), // larger than plan (resized in GUI)
			DatastoreID: &datastoreID,
			// ImportFrom is nil - PVE doesn't return this for existing disks
		},
	}

	// Plan has the original size from terraform config
	importFrom := "nfs-v3-tmp-01:iso/ubuntu-22.04-cloud.qcow2"
	planDisks := vms.CustomStorageDevices{
		"scsi0": &vms.CustomStorageDevice{
			Size:        types.DiskSizeFromGigabytes(20), // smaller than current
			DatastoreID: &datastoreID,
			ImportFrom:  &importFrom, // preserved from terraform config
		},
	}

	updateBody := &vms.UpdateRequestBody{}
	var client proxmox.Client

	ctx := context.Background()
	vmID := 1002
	nodeName := "test-node"

	// Force HasChange to return true
	err := resourceData.Set(MkDisk, []any{
		map[string]any{
			mkDiskInterface:   "scsi0",
			mkDiskDatastoreID: "nfs-v3-tmp-01",
			mkDiskSize:        30, // current size
			mkDiskSpeed:       []any{},
		},
	})
	require.NoError(t, err)

	err = resourceData.Set(MkDisk, []any{
		map[string]any{
			mkDiskInterface:   "scsi0",
			mkDiskDatastoreID: "nfs-v3-tmp-01",
			mkDiskSize:        20, // plan size
			mkDiskSpeed:       []any{},
		},
	})
	require.NoError(t, err)

	_, err = Update(ctx, client, nodeName, vmID, resourceData, planDisks, currentDisks, updateBody)
	require.NoError(t, err)

	// the update body should NOT contain ImportFrom - we're updating existing disk, not re-importing
	if updateBody.CustomStorageDevices != nil && updateBody.CustomStorageDevices["scsi0"] != nil {
		require.Nil(t, updateBody.CustomStorageDevices["scsi0"].ImportFrom,
			"should not set ImportFrom when updating existing disk")
	}
}

func TestDiskDeletionDetectionInGetDiskDeviceObjects(t *testing.T) {
	t.Parallel()

	diskSchema := Schema()

	// Create a simple resource schema for testing
	resource := &schema.Resource{
		Schema: diskSchema,
	}

	// Create test configuration with multiple disks in specific order
	oldDiskList := []any{
		map[string]any{
			mkDiskInterface:   "scsi0",
			mkDiskDatastoreID: "local",
			mkDiskSize:        32,
			mkDiskAIO:         "io_uring",
			mkDiskCache:       "none",
			mkDiskDiscard:     "ignore",
			mkDiskBackup:      true,
			mkDiskIOThread:    false,
			mkDiskReplicate:   true,
			mkDiskSSD:         false,
			mkDiskSerial:      "",
			mkDiskSpeed:       []any{},
		},
		map[string]any{
			mkDiskInterface:   "scsi1",
			mkDiskDatastoreID: "local",
			mkDiskSize:        20,
			mkDiskAIO:         "io_uring",
			mkDiskCache:       "none",
			mkDiskDiscard:     "ignore",
			mkDiskBackup:      true,
			mkDiskIOThread:    false,
			mkDiskReplicate:   true,
			mkDiskSSD:         false,
			mkDiskSerial:      "",
			mkDiskSpeed:       []any{},
		},
		map[string]any{
			mkDiskInterface:   "scsi2",
			mkDiskDatastoreID: "local",
			mkDiskSize:        1,
			mkDiskAIO:         "io_uring",
			mkDiskCache:       "none",
			mkDiskDiscard:     "ignore",
			mkDiskBackup:      true,
			mkDiskIOThread:    false,
			mkDiskReplicate:   true,
			mkDiskSSD:         false,
			mkDiskSerial:      "",
			mkDiskSpeed:       []any{},
		},
		map[string]any{
			mkDiskInterface:   "scsi3",
			mkDiskDatastoreID: "local",
			mkDiskSize:        50,
			mkDiskAIO:         "io_uring",
			mkDiskCache:       "none",
			mkDiskDiscard:     "ignore",
			mkDiskBackup:      true,
			mkDiskIOThread:    false,
			mkDiskReplicate:   true,
			mkDiskSSD:         false,
			mkDiskSerial:      "",
			mkDiskSpeed:       []any{},
		},
		map[string]any{
			mkDiskInterface:   "scsi7",
			mkDiskDatastoreID: "local",
			mkDiskSize:        4,
			mkDiskAIO:         "io_uring",
			mkDiskCache:       "none",
			mkDiskDiscard:     "ignore",
			mkDiskBackup:      true,
			mkDiskIOThread:    false,
			mkDiskReplicate:   true,
			mkDiskSSD:         false,
			mkDiskSerial:      "",
			mkDiskSpeed:       []any{},
		},
	}

	// New configuration with scsi7 removed
	newDiskList := []any{
		map[string]any{
			mkDiskInterface:   "scsi0",
			mkDiskDatastoreID: "local",
			mkDiskSize:        32,
			mkDiskAIO:         "io_uring",
			mkDiskCache:       "none",
			mkDiskDiscard:     "ignore",
			mkDiskBackup:      true,
			mkDiskIOThread:    false,
			mkDiskReplicate:   true,
			mkDiskSSD:         false,
			mkDiskSerial:      "",
			mkDiskSpeed:       []any{},
		},
		map[string]any{
			mkDiskInterface:   "scsi1",
			mkDiskDatastoreID: "local",
			mkDiskSize:        20,
			mkDiskAIO:         "io_uring",
			mkDiskCache:       "none",
			mkDiskDiscard:     "ignore",
			mkDiskBackup:      true,
			mkDiskIOThread:    false,
			mkDiskReplicate:   true,
			mkDiskSSD:         false,
			mkDiskSerial:      "",
			mkDiskSpeed:       []any{},
		},
		map[string]any{
			mkDiskInterface:   "scsi2",
			mkDiskDatastoreID: "local",
			mkDiskSize:        1,
			mkDiskAIO:         "io_uring",
			mkDiskCache:       "none",
			mkDiskDiscard:     "ignore",
			mkDiskBackup:      true,
			mkDiskIOThread:    false,
			mkDiskReplicate:   true,
			mkDiskSSD:         false,
			mkDiskSerial:      "",
			mkDiskSpeed:       []any{},
		},
		map[string]any{
			mkDiskInterface:   "scsi3",
			mkDiskDatastoreID: "local",
			mkDiskSize:        50,
			mkDiskAIO:         "io_uring",
			mkDiskCache:       "none",
			mkDiskDiscard:     "ignore",
			mkDiskBackup:      true,
			mkDiskIOThread:    false,
			mkDiskReplicate:   true,
			mkDiskSSD:         false,
			mkDiskSerial:      "",
			mkDiskSpeed:       []any{},
		},
	}

	// Create resource data
	resourceData := schema.TestResourceDataRaw(t, diskSchema, map[string]any{
		MkDisk: oldDiskList,
	})

	// Get old disk device objects
	oldDiskDevices, err := GetDiskDeviceObjects(resourceData, resource, oldDiskList)
	require.NoError(t, err)
	require.Len(t, oldDiskDevices, 5)

	// Get new disk device objects
	newDiskDevices, err := GetDiskDeviceObjects(resourceData, resource, newDiskList)
	require.NoError(t, err)
	require.Len(t, newDiskDevices, 4)

	// Verify that the removed interface (scsi7) is detected
	require.Contains(t, oldDiskDevices, "scsi7")
	require.NotContains(t, newDiskDevices, "scsi7")

	// Verify that all other interfaces remain in correct positions
	for _, iface := range []string{"scsi0", "scsi1", "scsi2", "scsi3"} {
		require.Contains(t, oldDiskDevices, iface)
		require.Contains(t, newDiskDevices, iface)

		// Verify disk properties remain unchanged for existing disks
		oldDisk := oldDiskDevices[iface]
		newDisk := newDiskDevices[iface]
		require.Equal(t, oldDisk.Size, newDisk.Size, "Disk size should remain unchanged for interface %s", iface)
		require.Equal(t, oldDisk.DatastoreID, newDisk.DatastoreID, "Datastore ID should remain unchanged for interface %s", iface)
	}

	// Simulate the deletion detection logic that should happen in vmUpdate
	// This is what should identify the disk for deletion
	var deletedInterfaces []string

	for oldIface := range oldDiskDevices {
		if _, present := newDiskDevices[oldIface]; !present {
			deletedInterfaces = append(deletedInterfaces, oldIface)
		}
	}

	// Verify exactly one disk is detected as deleted
	require.Len(t, deletedInterfaces, 1, "Exactly one disk should be detected as deleted")
	require.Equal(t, "scsi7", deletedInterfaces[0], "scsi7 should be the deleted disk")
}

func TestDiskDeletionWithBootDiskProtection(t *testing.T) {
	t.Parallel()

	// Mock current disk configuration (what Proxmox currently has)
	currentDisks := vms.CustomStorageDevices{
		"scsi0": &vms.CustomStorageDevice{
			Size:        types.DiskSizeFromGigabytes(32),
			DatastoreID: ptr.Ptr("local"),
		},
		"scsi1": &vms.CustomStorageDevice{
			Size:        types.DiskSizeFromGigabytes(20),
			DatastoreID: ptr.Ptr("local"),
		},
		"scsi7": &vms.CustomStorageDevice{
			Size:        types.DiskSizeFromGigabytes(4),
			DatastoreID: ptr.Ptr("local"),
		},
	}

	// Test case 1: Try to delete boot disk (should be protected)
	planDisksWithBootDiskRemoved := vms.CustomStorageDevices{
		"scsi1": &vms.CustomStorageDevice{
			Size:        types.DiskSizeFromGigabytes(20),
			DatastoreID: ptr.Ptr("local"),
		},
		"scsi7": &vms.CustomStorageDevice{
			Size:        types.DiskSizeFromGigabytes(4),
			DatastoreID: ptr.Ptr("local"),
		},
	}

	bootDevices := []string{"scsi0", "net0"} // scsi0 is in boot order

	// Check for deleted disks and boot protection
	bootDiskDeletionAttempted := false

	for currentInterface := range currentDisks {
		if _, present := planDisksWithBootDiskRemoved[currentInterface]; !present {
			// Check if this is a boot disk
			if slices.Contains(bootDevices, currentInterface) {
				bootDiskDeletionAttempted = true
			}
		}
	}

	require.True(t, bootDiskDeletionAttempted, "Boot disk deletion should be detected and would be blocked")

	// Test case 2: Delete non-boot disk (should be allowed)
	planDisksWithNonBootDiskRemoved := vms.CustomStorageDevices{
		"scsi0": &vms.CustomStorageDevice{
			Size:        types.DiskSizeFromGigabytes(32),
			DatastoreID: ptr.Ptr("local"),
		},
		"scsi1": &vms.CustomStorageDevice{
			Size:        types.DiskSizeFromGigabytes(20),
			DatastoreID: ptr.Ptr("local"),
		},
	}

	var deletedInterfaces []string
	bootDiskInDeletion := false

	for currentInterface := range currentDisks {
		if _, present := planDisksWithNonBootDiskRemoved[currentInterface]; !present {
			deletedInterfaces = append(deletedInterfaces, currentInterface)

			// Check if this is a boot disk
			if slices.Contains(bootDevices, currentInterface) {
				bootDiskInDeletion = true
			}
		}
	}

	require.Len(t, deletedInterfaces, 1, "One disk should be marked for deletion")
	require.Equal(t, "scsi7", deletedInterfaces[0], "scsi7 should be the deleted disk")
	require.False(t, bootDiskInDeletion, "No boot disk should be in deletion list")
}

func TestOriginalBugScenario(t *testing.T) {
	t.Parallel()

	// This represents the original VM configuration from the bug report
	originalDisks := vms.CustomStorageDevices{
		"scsi0": &vms.CustomStorageDevice{
			Size:        types.DiskSizeFromGigabytes(32),
			DatastoreID: ptr.Ptr("local-lvm"),
		},
		"scsi1": &vms.CustomStorageDevice{
			Size:        types.DiskSizeFromGigabytes(20),
			DatastoreID: ptr.Ptr("local-lvm"),
		},
		"scsi2": &vms.CustomStorageDevice{
			Size:        types.DiskSizeFromGigabytes(1),
			DatastoreID: ptr.Ptr("local-lvm"),
		},
		"scsi3": &vms.CustomStorageDevice{
			Size:        types.DiskSizeFromGigabytes(50),
			DatastoreID: ptr.Ptr("local-lvm"),
		},
		"scsi4": &vms.CustomStorageDevice{
			Size:        types.DiskSizeFromGigabytes(1),
			DatastoreID: ptr.Ptr("local-lvm"),
		},
		"scsi5": &vms.CustomStorageDevice{
			Size:        types.DiskSizeFromGigabytes(50),
			DatastoreID: ptr.Ptr("local-lvm"),
		},
		"scsi6": &vms.CustomStorageDevice{
			Size:        types.DiskSizeFromGigabytes(1),
			DatastoreID: ptr.Ptr("local-lvm"),
		},
		"scsi7": &vms.CustomStorageDevice{
			Size:        types.DiskSizeFromGigabytes(4),
			DatastoreID: ptr.Ptr("local-lvm"),
		},
	}

	// This represents the new configuration with scsi7 removed
	newDisks := vms.CustomStorageDevices{
		"scsi0": &vms.CustomStorageDevice{
			Size:        types.DiskSizeFromGigabytes(32),
			DatastoreID: ptr.Ptr("local-lvm"),
		},
		"scsi1": &vms.CustomStorageDevice{
			Size:        types.DiskSizeFromGigabytes(20),
			DatastoreID: ptr.Ptr("local-lvm"),
		},
		"scsi2": &vms.CustomStorageDevice{
			Size:        types.DiskSizeFromGigabytes(1),
			DatastoreID: ptr.Ptr("local-lvm"),
		},
		"scsi3": &vms.CustomStorageDevice{
			Size:        types.DiskSizeFromGigabytes(50),
			DatastoreID: ptr.Ptr("local-lvm"),
		},
		"scsi4": &vms.CustomStorageDevice{
			Size:        types.DiskSizeFromGigabytes(1),
			DatastoreID: ptr.Ptr("local-lvm"),
		},
		"scsi5": &vms.CustomStorageDevice{
			Size:        types.DiskSizeFromGigabytes(50),
			DatastoreID: ptr.Ptr("local-lvm"),
		},
		"scsi6": &vms.CustomStorageDevice{
			Size:        types.DiskSizeFromGigabytes(1),
			DatastoreID: ptr.Ptr("local-lvm"),
		},
	}

	// Simulate the disk deletion detection logic
	var deletesToAdd []string

	for oldIface := range originalDisks {
		if _, present := newDisks[oldIface]; !present {
			deletesToAdd = append(deletesToAdd, oldIface)
		}
	}

	// Verify behavior
	require.Len(t, deletesToAdd, 1, "Exactly one disk should be marked for deletion")
	require.Equal(t, "scsi7", deletesToAdd[0], "scsi7 should be the only disk marked for deletion")

	// Verify all other disks remain unchanged and in correct positions
	expectedInterfaces := []string{"scsi0", "scsi1", "scsi2", "scsi3", "scsi4", "scsi5", "scsi6"}
	for _, iface := range expectedInterfaces {
		require.Contains(t, originalDisks, iface, "Original disks should contain %s", iface)
		require.Contains(t, newDisks, iface, "New disks should contain %s", iface)

		originalDisk := originalDisks[iface]
		newDisk := newDisks[iface]

		// Verify no disk properties changed (no re-ordering/reassignment)
		require.Equal(t, originalDisk.Size, newDisk.Size, "Disk %s size should remain unchanged", iface)
		require.Equal(t, originalDisk.DatastoreID, newDisk.DatastoreID, "Disk %s datastore should remain unchanged", iface)
	}
}

// TestDiskSpeedSettingsPerDisk verifies that each disk gets its own speed settings
// and not the speed settings from the first disk (fixes issue #2467).
func TestDiskSpeedSettingsPerDisk(t *testing.T) {
	t.Parallel()

	diskSchema := Schema()
	resource := &schema.Resource{Schema: diskSchema}

	// create two disks with different speed settings
	diskList := []any{
		map[string]any{
			mkDiskInterface:   "scsi0",
			mkDiskDatastoreID: "local",
			mkDiskSize:        50,
			mkDiskAIO:         "io_uring",
			mkDiskBackup:      true,
			mkDiskCache:       "none",
			mkDiskDiscard:     "ignore",
			mkDiskIOThread:    false,
			mkDiskReplicate:   true,
			mkDiskSerial:      "",
			mkDiskSSD:         false,
			mkDiskSpeed: []any{
				map[string]any{
					mkDiskIopsRead:            100,
					mkDiskIopsWrite:           200,
					mkDiskIopsReadBurstable:   1000,
					mkDiskIopsWriteBurstable:  2000,
					mkDiskSpeedRead:           10,
					mkDiskSpeedWrite:          20,
					mkDiskSpeedReadBurstable:  100,
					mkDiskSpeedWriteBurstable: 200,
				},
			},
		},
		map[string]any{
			mkDiskInterface:   "scsi1",
			mkDiskDatastoreID: "local",
			mkDiskSize:        100,
			mkDiskAIO:         "io_uring",
			mkDiskBackup:      true,
			mkDiskCache:       "none",
			mkDiskDiscard:     "ignore",
			mkDiskIOThread:    false,
			mkDiskReplicate:   true,
			mkDiskSerial:      "",
			mkDiskSSD:         false,
			mkDiskSpeed: []any{
				map[string]any{
					mkDiskIopsRead:            300,
					mkDiskIopsWrite:           400,
					mkDiskIopsReadBurstable:   3000,
					mkDiskIopsWriteBurstable:  4000,
					mkDiskSpeedRead:           30,
					mkDiskSpeedWrite:          40,
					mkDiskSpeedReadBurstable:  300,
					mkDiskSpeedWriteBurstable: 400,
				},
			},
		},
		map[string]any{
			mkDiskInterface:   "scsi2",
			mkDiskDatastoreID: "local",
			mkDiskSize:        200,
			mkDiskAIO:         "io_uring",
			mkDiskBackup:      true,
			mkDiskCache:       "none",
			mkDiskDiscard:     "ignore",
			mkDiskIOThread:    false,
			mkDiskReplicate:   true,
			mkDiskSerial:      "",
			mkDiskSSD:         false,
			mkDiskSpeed:       []any{}, // no speed limits
		},
	}

	resourceData := schema.TestResourceDataRaw(t, diskSchema, map[string]any{
		MkDisk: diskList,
	})

	diskDevices, err := GetDiskDeviceObjects(resourceData, resource, diskList)
	require.NoError(t, err)
	require.Len(t, diskDevices, 3)

	// verify scsi0 has its own speed settings
	scsi0 := diskDevices["scsi0"]
	require.NotNil(t, scsi0)
	require.NotNil(t, scsi0.IopsRead, "scsi0 should have IopsRead")
	require.Equal(t, 100, *scsi0.IopsRead)
	require.NotNil(t, scsi0.IopsWrite, "scsi0 should have IopsWrite")
	require.Equal(t, 200, *scsi0.IopsWrite)
	require.NotNil(t, scsi0.MaxIopsRead, "scsi0 should have MaxIopsRead")
	require.Equal(t, 1000, *scsi0.MaxIopsRead)
	require.NotNil(t, scsi0.MaxIopsWrite, "scsi0 should have MaxIopsWrite")
	require.Equal(t, 2000, *scsi0.MaxIopsWrite)
	require.NotNil(t, scsi0.MaxReadSpeedMbps, "scsi0 should have MaxReadSpeedMbps")
	require.Equal(t, 10, *scsi0.MaxReadSpeedMbps)
	require.NotNil(t, scsi0.MaxWriteSpeedMbps, "scsi0 should have MaxWriteSpeedMbps")
	require.Equal(t, 20, *scsi0.MaxWriteSpeedMbps)
	require.NotNil(t, scsi0.BurstableReadSpeedMbps, "scsi0 should have BurstableReadSpeedMbps")
	require.Equal(t, 100, *scsi0.BurstableReadSpeedMbps)
	require.NotNil(t, scsi0.BurstableWriteSpeedMbps, "scsi0 should have BurstableWriteSpeedMbps")
	require.Equal(t, 200, *scsi0.BurstableWriteSpeedMbps)

	// verify scsi1 has DIFFERENT speed settings (not scsi0's)
	scsi1 := diskDevices["scsi1"]
	require.NotNil(t, scsi1)
	require.NotNil(t, scsi1.IopsRead, "scsi1 should have IopsRead")
	require.Equal(t, 300, *scsi1.IopsRead, "scsi1 IopsRead should be 300, not 100 from scsi0")
	require.NotNil(t, scsi1.IopsWrite, "scsi1 should have IopsWrite")
	require.Equal(t, 400, *scsi1.IopsWrite, "scsi1 IopsWrite should be 400, not 200 from scsi0")
	require.NotNil(t, scsi1.MaxIopsRead, "scsi1 should have MaxIopsRead")
	require.Equal(t, 3000, *scsi1.MaxIopsRead, "scsi1 MaxIopsRead should be 3000, not 1000 from scsi0")
	require.NotNil(t, scsi1.MaxIopsWrite, "scsi1 should have MaxIopsWrite")
	require.Equal(t, 4000, *scsi1.MaxIopsWrite, "scsi1 MaxIopsWrite should be 4000, not 2000 from scsi0")
	require.NotNil(t, scsi1.MaxReadSpeedMbps, "scsi1 should have MaxReadSpeedMbps")
	require.Equal(t, 30, *scsi1.MaxReadSpeedMbps, "scsi1 MaxReadSpeedMbps should be 30, not 10 from scsi0")
	require.NotNil(t, scsi1.MaxWriteSpeedMbps, "scsi1 should have MaxWriteSpeedMbps")
	require.Equal(t, 40, *scsi1.MaxWriteSpeedMbps, "scsi1 MaxWriteSpeedMbps should be 40, not 20 from scsi0")
	require.NotNil(t, scsi1.BurstableReadSpeedMbps, "scsi1 should have BurstableReadSpeedMbps")
	require.Equal(t, 300, *scsi1.BurstableReadSpeedMbps, "scsi1 BurstableReadSpeedMbps should be 300, not 100 from scsi0")
	require.NotNil(t, scsi1.BurstableWriteSpeedMbps, "scsi1 should have BurstableWriteSpeedMbps")
	require.Equal(t, 400, *scsi1.BurstableWriteSpeedMbps, "scsi1 BurstableWriteSpeedMbps should be 400, not 200 from scsi0")

	// verify scsi2 has NO speed settings (empty speed block)
	scsi2 := diskDevices["scsi2"]
	require.NotNil(t, scsi2)
	require.Nil(t, scsi2.IopsRead, "scsi2 should NOT have IopsRead (empty speed block)")
	require.Nil(t, scsi2.IopsWrite, "scsi2 should NOT have IopsWrite (empty speed block)")
	require.Nil(t, scsi2.MaxIopsRead, "scsi2 should NOT have MaxIopsRead (empty speed block)")
	require.Nil(t, scsi2.MaxIopsWrite, "scsi2 should NOT have MaxIopsWrite (empty speed block)")
	require.Nil(t, scsi2.MaxReadSpeedMbps, "scsi2 should NOT have MaxReadSpeedMbps (empty speed block)")
	require.Nil(t, scsi2.MaxWriteSpeedMbps, "scsi2 should NOT have MaxWriteSpeedMbps (empty speed block)")
	require.Nil(t, scsi2.BurstableReadSpeedMbps, "scsi2 should NOT have BurstableReadSpeedMbps (empty speed block)")
	require.Nil(t, scsi2.BurstableWriteSpeedMbps, "scsi2 should NOT have BurstableWriteSpeedMbps (empty speed block)")
}
