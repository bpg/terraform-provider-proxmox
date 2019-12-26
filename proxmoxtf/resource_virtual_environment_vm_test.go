/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmoxtf

import (
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
)

// TestResourceVirtualEnvironmentVMInstantiation tests whether the ResourceVirtualEnvironmentVM instance can be instantiated.
func TestResourceVirtualEnvironmentVMInstantiation(t *testing.T) {
	s := resourceVirtualEnvironmentVM()

	if s == nil {
		t.Fatalf("Cannot instantiate resourceVirtualEnvironmentVM")
	}
}

// TestResourceVirtualEnvironmentVMSchema tests the resourceVirtualEnvironmentVM schema.
func TestResourceVirtualEnvironmentVMSchema(t *testing.T) {
	s := resourceVirtualEnvironmentVM()

	testOptionalArguments(t, s, []string{
		mkResourceVirtualEnvironmentVMCPU,
		mkResourceVirtualEnvironmentVMDisk,
		mkResourceVirtualEnvironmentVMMemory,
		mkResourceVirtualEnvironmentVMName,
		mkResourceVirtualEnvironmentVMNetworkDevice,
		mkResourceVirtualEnvironmentVMVMID,
	})

	testSchemaValueTypes(t, s, []string{
		mkResourceVirtualEnvironmentVMCPU,
		mkResourceVirtualEnvironmentVMDisk,
		mkResourceVirtualEnvironmentVMMemory,
		mkResourceVirtualEnvironmentVMName,
		mkResourceVirtualEnvironmentVMNetworkDevice,
		mkResourceVirtualEnvironmentVMVMID,
	}, []schema.ValueType{
		schema.TypeList,
		schema.TypeList,
		schema.TypeList,
		schema.TypeString,
		schema.TypeList,
		schema.TypeInt,
	})

	cpuSchema := testNestedSchemaExistence(t, s, mkResourceVirtualEnvironmentVMCPU)

	testOptionalArguments(t, cpuSchema, []string{
		mkResourceVirtualEnvironmentVMCPUCores,
		mkResourceVirtualEnvironmentVMCPUSockets,
	})

	testSchemaValueTypes(t, cpuSchema, []string{
		mkResourceVirtualEnvironmentVMCPUCores,
		mkResourceVirtualEnvironmentVMCPUSockets,
	}, []schema.ValueType{
		schema.TypeInt,
		schema.TypeInt,
	})

	diskSchema := testNestedSchemaExistence(t, s, mkResourceVirtualEnvironmentVMDisk)

	testOptionalArguments(t, diskSchema, []string{
		mkResourceVirtualEnvironmentVMDiskDatastoreID,
		mkResourceVirtualEnvironmentVMDiskFileID,
		mkResourceVirtualEnvironmentVMDiskSize,
	})

	testSchemaValueTypes(t, diskSchema, []string{
		mkResourceVirtualEnvironmentVMDiskDatastoreID,
		mkResourceVirtualEnvironmentVMDiskFileID,
		mkResourceVirtualEnvironmentVMDiskSize,
	}, []schema.ValueType{
		schema.TypeString,
		schema.TypeString,
		schema.TypeInt,
	})

	memorySchema := testNestedSchemaExistence(t, s, mkResourceVirtualEnvironmentVMMemory)

	testOptionalArguments(t, memorySchema, []string{
		mkResourceVirtualEnvironmentVMMemoryDedicated,
		mkResourceVirtualEnvironmentVMMemoryFloating,
		mkResourceVirtualEnvironmentVMMemoryShared,
	})

	testSchemaValueTypes(t, memorySchema, []string{
		mkResourceVirtualEnvironmentVMMemoryDedicated,
		mkResourceVirtualEnvironmentVMMemoryFloating,
		mkResourceVirtualEnvironmentVMMemoryShared,
	}, []schema.ValueType{
		schema.TypeInt,
		schema.TypeInt,
		schema.TypeInt,
	})

	networkDeviceSchema := testNestedSchemaExistence(t, s, mkResourceVirtualEnvironmentVMNetworkDevice)

	testOptionalArguments(t, networkDeviceSchema, []string{
		mkResourceVirtualEnvironmentVMNetworkDeviceBridge,
		mkResourceVirtualEnvironmentVMNetworkDeviceMACAddress,
		mkResourceVirtualEnvironmentVMNetworkDeviceModel,
		mkResourceVirtualEnvironmentVMNetworkDeviceVLANID,
	})

	testSchemaValueTypes(t, networkDeviceSchema, []string{
		mkResourceVirtualEnvironmentVMNetworkDeviceBridge,
		mkResourceVirtualEnvironmentVMNetworkDeviceMACAddress,
		mkResourceVirtualEnvironmentVMNetworkDeviceModel,
		mkResourceVirtualEnvironmentVMNetworkDeviceVLANID,
	}, []schema.ValueType{
		schema.TypeString,
		schema.TypeString,
		schema.TypeString,
		schema.TypeInt,
	})
}
