/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package cdrom

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/vms"
)

// Conversion helpers.

func TestGetCDROMDeviceObjects(t *testing.T) {
	t.Parallel()

	deviceObjects := GetCDROMDeviceObjects([]any{
		map[string]any{
			MkCDROMInterface: "ide2",
			MkCDROMFileID:    "none",
		},
		map[string]any{
			MkCDROMInterface: "sata3",
			MkCDROMFileID:    "",
		},
	})

	require.Len(t, deviceObjects, 2)
	require.Equal(t, "none", deviceObjects["ide2"].FileVolume)
	require.Equal(t, "cdrom", deviceObjects["sata3"].FileVolume)
	require.NotNil(t, deviceObjects["ide2"].Media)
	require.Equal(t, "cdrom", *deviceObjects["ide2"].Media)
}

func TestGetCDROMStorageDevices(t *testing.T) {
	t.Parallel()

	cdromMedia := "cdrom"
	diskMedia := "disk"

	vmConfig := &vms.GetResponseData{
		StorageDevices: vms.CustomStorageDevices{
			"ide2": {
				FileVolume: "none",
				Media:      &cdromMedia,
			},
			"sata3": {
				FileVolume: "local:iso/example.iso",
				Media:      &cdromMedia,
			},
			"scsi0": {
				FileVolume: "local-lvm:vm-100-disk-0",
				Media:      &diskMedia,
			},
		},
	}

	cdroms := GetCDROMStorageDevices(vmConfig)

	require.Len(t, cdroms, 2)
	require.Contains(t, cdroms, "ide2")
	require.Contains(t, cdroms, "sata3")
	require.NotContains(t, cdroms, "scsi0")
}

// Schema.

func TestSchema(t *testing.T) {
	t.Parallel()

	s := Schema()
	require.Contains(t, s, MkCDROM)

	cdromSchema := s[MkCDROM].Elem.(*schema.Resource).Schema
	require.True(t, cdromSchema[MkCDROMInterface].Required)
	require.False(t, cdromSchema[MkCDROMInterface].Optional)
	require.Equal(t, DefaultFileID, cdromSchema[MkCDROMFileID].Default)
}

// Read and state building.

func TestBuildStatePreservesOrderingAndExplicitEmptyFileID(t *testing.T) {
	t.Parallel()

	currentCDROMList := []any{
		map[string]any{
			MkCDROMInterface: "sata3",
			MkCDROMFileID:    "",
		},
		map[string]any{
			MkCDROMInterface: "ide2",
			MkCDROMFileID:    "none",
		},
	}

	cdromMedia := "cdrom"
	deviceObjects := vms.CustomStorageDevices{
		"ide2": {
			FileVolume: "none",
			Media:      &cdromMedia,
		},
		"sata3": {
			FileVolume: "cdrom",
			Media:      &cdromMedia,
		},
	}

	cdromList := BuildState(currentCDROMList, deviceObjects, false)
	require.Len(t, cdromList, 2)

	first := cdromList[0].(map[string]any)
	second := cdromList[1].(map[string]any)

	require.Equal(t, "sata3", first[MkCDROMInterface])
	require.Empty(t, first[MkCDROMFileID], "explicit empty file_id should be preserved from current state")
	require.Equal(t, "ide2", second[MkCDROMInterface])
	require.Equal(t, "none", second[MkCDROMFileID])
}

func TestReadPreservesOrderingAndExplicitEmptyFileID(t *testing.T) {
	t.Parallel()

	cdromSchema := Schema()
	currentCDROMList := []any{
		map[string]any{
			MkCDROMInterface: "sata3",
			MkCDROMFileID:    "",
		},
		map[string]any{
			MkCDROMInterface: "ide2",
			MkCDROMFileID:    "none",
		},
	}

	resourceData := schema.TestResourceDataRaw(t, cdromSchema, map[string]any{
		MkCDROM: currentCDROMList,
	})

	cdromMedia := "cdrom"
	deviceObjects := vms.CustomStorageDevices{
		"ide2": {
			FileVolume: "none",
			Media:      &cdromMedia,
		},
		"sata3": {
			FileVolume: "cdrom",
			Media:      &cdromMedia,
		},
	}

	diags := Read(resourceData, deviceObjects, false)
	require.Empty(t, diags)

	cdromList := resourceData.Get(MkCDROM).([]any)
	require.Len(t, cdromList, 2)

	first := cdromList[0].(map[string]any)
	second := cdromList[1].(map[string]any)

	require.Equal(t, "sata3", first[MkCDROMInterface])
	require.Empty(t, first[MkCDROMFileID], "explicit empty file_id should be preserved from current state")
	require.Equal(t, "ide2", second[MkCDROMInterface])
	require.Equal(t, "none", second[MkCDROMFileID])
}

func TestBuildStateSkipsCloneWithoutCurrentState(t *testing.T) {
	t.Parallel()

	cdromMedia := "cdrom"
	deviceObjects := vms.CustomStorageDevices{
		"ide2": {
			FileVolume: "none",
			Media:      &cdromMedia,
		},
	}

	require.Nil(t, BuildState(nil, deviceObjects, true))
}

func TestReadSkipsSettingCDROMForCloneWithoutCurrentState(t *testing.T) {
	t.Parallel()

	cdromSchema := Schema()
	resourceData := schema.TestResourceDataRaw(t, cdromSchema, map[string]any{})

	cdromMedia := "cdrom"
	deviceObjects := vms.CustomStorageDevices{
		"ide2": {
			FileVolume: "none",
			Media:      &cdromMedia,
		},
	}

	diags := Read(resourceData, deviceObjects, true)
	require.Empty(t, diags)
	require.Empty(t, resourceData.Get(MkCDROM).([]any))
}

// Clone and update behavior.

func TestMergeCloneDevices(t *testing.T) {
	t.Parallel()

	ideDevices := vms.CustomStorageDevices{
		"ide1": {
			FileVolume: "existing",
		},
	}

	MergeCloneDevices(GetCDROMDeviceObjects([]any{
		map[string]any{
			MkCDROMInterface: "ide2",
			MkCDROMFileID:    "none",
		},
		map[string]any{
			MkCDROMInterface: "sata3",
			MkCDROMFileID:    "cdrom",
		},
	}), ideDevices)

	require.Len(t, ideDevices, 3)
	require.Equal(t, "existing", ideDevices["ide1"].FileVolume)
	require.Equal(t, "none", ideDevices["ide2"].FileVolume)
	require.Equal(t, "cdrom", ideDevices["sata3"].FileVolume)
}

func TestApplyDeviceObjectDiff(t *testing.T) {
	t.Parallel()

	updateBody := &vms.UpdateRequestBody{}
	del := applyDeviceObjectDiff(
		GetCDROMDeviceObjects([]any{
			map[string]any{
				MkCDROMInterface: "ide2",
				MkCDROMFileID:    "none",
			},
			map[string]any{
				MkCDROMInterface: "sata3",
				MkCDROMFileID:    "none",
			},
		}),
		GetCDROMDeviceObjects([]any{
			map[string]any{
				MkCDROMInterface: "ide2",
				MkCDROMFileID:    "cdrom",
			},
			map[string]any{
				MkCDROMInterface: "scsi5",
				MkCDROMFileID:    "none",
			},
		}),
		updateBody,
		nil,
	)

	require.Equal(t, []string{"sata3"}, del)
	require.NotNil(t, updateBody.CustomStorageDevices)
	require.Len(t, updateBody.CustomStorageDevices, 2)

	require.Equal(t, "cdrom", updateBody.CustomStorageDevices["ide2"].FileVolume)
	require.Equal(t, "none", updateBody.CustomStorageDevices["scsi5"].FileVolume)
}

// Ordering and compatibility.

func TestReadDeterministicOrdering(t *testing.T) {
	t.Parallel()

	cdromSchema := Schema()
	currentCDROMList := []any{
		map[string]any{
			MkCDROMInterface: "scsi3",
			MkCDROMFileID:    "none",
		},
		map[string]any{
			MkCDROMInterface: "ide2",
			MkCDROMFileID:    "cdrom",
		},
	}

	cdromMedia := "cdrom"
	deviceObjects := vms.CustomStorageDevices{
		"ide2": {
			FileVolume: "cdrom",
			Media:      &cdromMedia,
		},
		"scsi3": {
			FileVolume: "none",
			Media:      &cdromMedia,
		},
	}

	const iterations = 5
	results := make([][]any, 0, iterations)

	for range iterations {
		resourceData := schema.TestResourceDataRaw(t, cdromSchema, map[string]any{
			MkCDROM: currentCDROMList,
		})

		diags := Read(resourceData, deviceObjects, false)
		require.Empty(t, diags)

		results = append(results, resourceData.Get(MkCDROM).([]any))
	}

	expected := results[0]
	for i, result := range results {
		require.True(t, reflect.DeepEqual(expected, result), "iteration %d differed from expected ordering", i)
	}
}

func TestBuildStateDefaultsFileIDToCDROM(t *testing.T) {
	t.Parallel()

	cdromMedia := "cdrom"
	deviceObjects := vms.CustomStorageDevices{
		"ide2": {
			FileVolume: DefaultFileID,
			Media:      &cdromMedia,
		},
	}

	cdromList := BuildState([]any{
		map[string]any{
			MkCDROMInterface: "ide2",
			MkCDROMFileID:    DefaultFileID,
		},
	}, deviceObjects, false)

	require.Len(t, cdromList, 1)
	require.Equal(t, DefaultFileID, cdromList[0].(map[string]any)[MkCDROMFileID])
}

func TestOrderedInterfaces(t *testing.T) {
	t.Parallel()

	deviceObjects := vms.CustomStorageDevices{
		"ide2":  {},
		"sata3": {},
		"scsi5": {},
	}

	ordered := OrderedInterfaces([]any{
		map[string]any{
			MkCDROMInterface: "sata3",
		},
		map[string]any{
			MkCDROMInterface: "ide2",
		},
	}, deviceObjects)

	require.Equal(t, []string{"sata3", "ide2", "scsi5"}, ordered)
}

// Validation.

func TestValidateInterfacesForMachine(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		machineType string
		cdrom       []any
		wantErr     string
	}{
		{
			name:        "non q35 allows ide3",
			machineType: "pc",
			cdrom:       []any{map[string]any{MkCDROMInterface: "ide3"}},
		},
		{
			name:        "q35 allows ide0 and ide2",
			machineType: "q35",
			cdrom: []any{
				map[string]any{MkCDROMInterface: "ide0"},
				map[string]any{MkCDROMInterface: "ide2"},
			},
		},
		{
			name:        "q35 allows sata and scsi",
			machineType: "pc-q35-9.0+pve1",
			cdrom: []any{
				map[string]any{MkCDROMInterface: "sata3"},
				map[string]any{MkCDROMInterface: "scsi5"},
			},
		},
		{
			name:        "q35 rejects ide1",
			machineType: "q35,viommu=virtio",
			cdrom:       []any{map[string]any{MkCDROMInterface: "ide1"}},
			wantErr:     `cdrom interface "ide1" is invalid for q35 machine type: only ide0 and ide2 are supported on the IDE bus`,
		},
		{
			name:        "q35 rejects ide3",
			machineType: "pc-q35-8.1",
			cdrom:       []any{map[string]any{MkCDROMInterface: "ide3"}},
			wantErr:     `cdrom interface "ide3" is invalid for q35 machine type: only ide0 and ide2 are supported on the IDE bus`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := ValidateInterfacesForMachine(tt.machineType, tt.cdrom)

			if tt.wantErr == "" {
				require.NoError(t, err)
				return
			}

			require.EqualError(t, err, tt.wantErr)
		})
	}
}
