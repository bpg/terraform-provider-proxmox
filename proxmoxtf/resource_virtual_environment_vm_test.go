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

	testRequiredArguments(t, s, []string{
		mkResourceVirtualEnvironmentVMNodeName,
	})

	testOptionalArguments(t, s, []string{
		mkResourceVirtualEnvironmentVMCDROM,
		mkResourceVirtualEnvironmentVMCloudInit,
		mkResourceVirtualEnvironmentVMCPU,
		mkResourceVirtualEnvironmentVMDescription,
		mkResourceVirtualEnvironmentVMDisk,
		mkResourceVirtualEnvironmentVMKeyboardLayout,
		mkResourceVirtualEnvironmentVMMemory,
		mkResourceVirtualEnvironmentVMName,
		mkResourceVirtualEnvironmentVMNetworkDevice,
		mkResourceVirtualEnvironmentVMOSType,
		mkResourceVirtualEnvironmentVMPoolID,
		mkResourceVirtualEnvironmentVMVMID,
	})

	testSchemaValueTypes(t, s, []string{
		mkResourceVirtualEnvironmentVMCDROM,
		mkResourceVirtualEnvironmentVMCloudInit,
		mkResourceVirtualEnvironmentVMCPU,
		mkResourceVirtualEnvironmentVMDescription,
		mkResourceVirtualEnvironmentVMDisk,
		mkResourceVirtualEnvironmentVMKeyboardLayout,
		mkResourceVirtualEnvironmentVMMemory,
		mkResourceVirtualEnvironmentVMName,
		mkResourceVirtualEnvironmentVMNetworkDevice,
		mkResourceVirtualEnvironmentVMOSType,
		mkResourceVirtualEnvironmentVMPoolID,
		mkResourceVirtualEnvironmentVMVMID,
	}, []schema.ValueType{
		schema.TypeList,
		schema.TypeList,
		schema.TypeList,
		schema.TypeString,
		schema.TypeList,
		schema.TypeString,
		schema.TypeList,
		schema.TypeString,
		schema.TypeList,
		schema.TypeString,
		schema.TypeString,
		schema.TypeInt,
	})

	agentSchema := testNestedSchemaExistence(t, s, mkResourceVirtualEnvironmentVMAgent)

	testOptionalArguments(t, agentSchema, []string{
		mkResourceVirtualEnvironmentVMAgentEnabled,
		mkResourceVirtualEnvironmentVMAgentTrim,
		mkResourceVirtualEnvironmentVMAgentType,
	})

	testSchemaValueTypes(t, agentSchema, []string{
		mkResourceVirtualEnvironmentVMAgentEnabled,
		mkResourceVirtualEnvironmentVMAgentTrim,
		mkResourceVirtualEnvironmentVMAgentType,
	}, []schema.ValueType{
		schema.TypeBool,
		schema.TypeBool,
		schema.TypeString,
	})

	cdromSchema := testNestedSchemaExistence(t, s, mkResourceVirtualEnvironmentVMCDROM)

	testOptionalArguments(t, cdromSchema, []string{
		mkResourceVirtualEnvironmentVMCDROMEnabled,
		mkResourceVirtualEnvironmentVMCDROMFileID,
	})

	testSchemaValueTypes(t, cdromSchema, []string{
		mkResourceVirtualEnvironmentVMCDROMEnabled,
		mkResourceVirtualEnvironmentVMCDROMFileID,
	}, []schema.ValueType{
		schema.TypeBool,
		schema.TypeString,
	})

	cloudInitSchema := testNestedSchemaExistence(t, s, mkResourceVirtualEnvironmentVMCloudInit)

	testRequiredArguments(t, cloudInitSchema, []string{
		mkResourceVirtualEnvironmentVMCloudInitUserAccount,
	})

	testOptionalArguments(t, cloudInitSchema, []string{
		mkResourceVirtualEnvironmentVMCloudInitDNS,
		mkResourceVirtualEnvironmentVMCloudInitIPConfig,
	})

	testSchemaValueTypes(t, cloudInitSchema, []string{
		mkResourceVirtualEnvironmentVMCloudInitDNS,
		mkResourceVirtualEnvironmentVMCloudInitIPConfig,
		mkResourceVirtualEnvironmentVMCloudInitUserAccount,
	}, []schema.ValueType{
		schema.TypeList,
		schema.TypeList,
		schema.TypeList,
	})

	cloudInitDNSSchema := testNestedSchemaExistence(t, cloudInitSchema, mkResourceVirtualEnvironmentVMCloudInitDNS)

	testOptionalArguments(t, cloudInitDNSSchema, []string{
		mkResourceVirtualEnvironmentVMCloudInitDNSDomain,
		mkResourceVirtualEnvironmentVMCloudInitDNSServer,
	})

	testSchemaValueTypes(t, cloudInitDNSSchema, []string{
		mkResourceVirtualEnvironmentVMCloudInitDNSDomain,
		mkResourceVirtualEnvironmentVMCloudInitDNSServer,
	}, []schema.ValueType{
		schema.TypeString,
		schema.TypeString,
	})

	cloudInitIPConfigSchema := testNestedSchemaExistence(t, cloudInitSchema, mkResourceVirtualEnvironmentVMCloudInitIPConfig)

	testOptionalArguments(t, cloudInitIPConfigSchema, []string{
		mkResourceVirtualEnvironmentVMCloudInitIPConfigIPv4,
		mkResourceVirtualEnvironmentVMCloudInitIPConfigIPv6,
	})

	testSchemaValueTypes(t, cloudInitIPConfigSchema, []string{
		mkResourceVirtualEnvironmentVMCloudInitIPConfigIPv4,
		mkResourceVirtualEnvironmentVMCloudInitIPConfigIPv6,
	}, []schema.ValueType{
		schema.TypeList,
		schema.TypeList,
	})

	cloudInitIPConfigIPv4Schema := testNestedSchemaExistence(t, cloudInitIPConfigSchema, mkResourceVirtualEnvironmentVMCloudInitIPConfigIPv4)

	testOptionalArguments(t, cloudInitIPConfigIPv4Schema, []string{
		mkResourceVirtualEnvironmentVMCloudInitIPConfigIPv4Address,
		mkResourceVirtualEnvironmentVMCloudInitIPConfigIPv4Gateway,
	})

	testSchemaValueTypes(t, cloudInitIPConfigIPv4Schema, []string{
		mkResourceVirtualEnvironmentVMCloudInitIPConfigIPv4Address,
		mkResourceVirtualEnvironmentVMCloudInitIPConfigIPv4Gateway,
	}, []schema.ValueType{
		schema.TypeString,
		schema.TypeString,
	})

	cloudInitIPConfigIPv6Schema := testNestedSchemaExistence(t, cloudInitIPConfigSchema, mkResourceVirtualEnvironmentVMCloudInitIPConfigIPv6)

	testOptionalArguments(t, cloudInitIPConfigIPv6Schema, []string{
		mkResourceVirtualEnvironmentVMCloudInitIPConfigIPv6Address,
		mkResourceVirtualEnvironmentVMCloudInitIPConfigIPv6Gateway,
	})

	testSchemaValueTypes(t, cloudInitIPConfigIPv6Schema, []string{
		mkResourceVirtualEnvironmentVMCloudInitIPConfigIPv6Address,
		mkResourceVirtualEnvironmentVMCloudInitIPConfigIPv6Gateway,
	}, []schema.ValueType{
		schema.TypeString,
		schema.TypeString,
	})

	cloudInitUserAccountSchema := testNestedSchemaExistence(t, cloudInitSchema, mkResourceVirtualEnvironmentVMCloudInitUserAccount)

	testRequiredArguments(t, cloudInitUserAccountSchema, []string{
		mkResourceVirtualEnvironmentVMCloudInitUserAccountKeys,
		mkResourceVirtualEnvironmentVMCloudInitUserAccountUsername,
	})

	testOptionalArguments(t, cloudInitUserAccountSchema, []string{
		mkResourceVirtualEnvironmentVMCloudInitUserAccountPassword,
	})

	testSchemaValueTypes(t, cloudInitUserAccountSchema, []string{
		mkResourceVirtualEnvironmentVMCloudInitUserAccountKeys,
		mkResourceVirtualEnvironmentVMCloudInitUserAccountPassword,
		mkResourceVirtualEnvironmentVMCloudInitUserAccountUsername,
	}, []schema.ValueType{
		schema.TypeList,
		schema.TypeString,
		schema.TypeString,
	})

	cpuSchema := testNestedSchemaExistence(t, s, mkResourceVirtualEnvironmentVMCPU)

	testOptionalArguments(t, cpuSchema, []string{
		mkResourceVirtualEnvironmentVMCPUCores,
		mkResourceVirtualEnvironmentVMCPUHotplugged,
		mkResourceVirtualEnvironmentVMCPUSockets,
	})

	testSchemaValueTypes(t, cpuSchema, []string{
		mkResourceVirtualEnvironmentVMCPUCores,
		mkResourceVirtualEnvironmentVMCPUHotplugged,
		mkResourceVirtualEnvironmentVMCPUSockets,
	}, []schema.ValueType{
		schema.TypeInt,
		schema.TypeInt,
		schema.TypeInt,
	})

	diskSchema := testNestedSchemaExistence(t, s, mkResourceVirtualEnvironmentVMDisk)

	testOptionalArguments(t, diskSchema, []string{
		mkResourceVirtualEnvironmentVMDiskDatastoreID,
		mkResourceVirtualEnvironmentVMDiskEnabled,
		mkResourceVirtualEnvironmentVMDiskFileFormat,
		mkResourceVirtualEnvironmentVMDiskFileID,
		mkResourceVirtualEnvironmentVMDiskSize,
	})

	testSchemaValueTypes(t, diskSchema, []string{
		mkResourceVirtualEnvironmentVMDiskDatastoreID,
		mkResourceVirtualEnvironmentVMDiskEnabled,
		mkResourceVirtualEnvironmentVMDiskFileFormat,
		mkResourceVirtualEnvironmentVMDiskFileID,
		mkResourceVirtualEnvironmentVMDiskSize,
	}, []schema.ValueType{
		schema.TypeString,
		schema.TypeBool,
		schema.TypeString,
		schema.TypeString,
		schema.TypeInt,
	})

	diskSpeedSchema := testNestedSchemaExistence(t, diskSchema, mkResourceVirtualEnvironmentVMDiskSpeed)

	testOptionalArguments(t, diskSpeedSchema, []string{
		mkResourceVirtualEnvironmentVMDiskSpeedRead,
		mkResourceVirtualEnvironmentVMDiskSpeedWrite,
	})

	testSchemaValueTypes(t, diskSpeedSchema, []string{
		mkResourceVirtualEnvironmentVMDiskSpeedRead,
		mkResourceVirtualEnvironmentVMDiskSpeedWrite,
	}, []schema.ValueType{
		schema.TypeInt,
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
		mkResourceVirtualEnvironmentVMNetworkDeviceEnabled,
		mkResourceVirtualEnvironmentVMNetworkDeviceMACAddress,
		mkResourceVirtualEnvironmentVMNetworkDeviceModel,
		mkResourceVirtualEnvironmentVMNetworkDeviceVLANIDs,
	})

	testSchemaValueTypes(t, networkDeviceSchema, []string{
		mkResourceVirtualEnvironmentVMNetworkDeviceBridge,
		mkResourceVirtualEnvironmentVMNetworkDeviceEnabled,
		mkResourceVirtualEnvironmentVMNetworkDeviceMACAddress,
		mkResourceVirtualEnvironmentVMNetworkDeviceModel,
		mkResourceVirtualEnvironmentVMNetworkDeviceVLANIDs,
	}, []schema.ValueType{
		schema.TypeString,
		schema.TypeBool,
		schema.TypeString,
		schema.TypeString,
		schema.TypeList,
	})
}
