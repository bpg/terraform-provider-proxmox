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
		mkResourceVirtualEnvironmentVMACPI,
		mkResourceVirtualEnvironmentVMAgent,
		mkResourceVirtualEnvironmentVMAudioDevice,
		mkResourceVirtualEnvironmentVMBIOS,
		mkResourceVirtualEnvironmentVMCDROM,
		mkResourceVirtualEnvironmentVMInitialization,
		mkResourceVirtualEnvironmentVMCPU,
		mkResourceVirtualEnvironmentVMDescription,
		mkResourceVirtualEnvironmentVMDisk,
		mkResourceVirtualEnvironmentVMKeyboardLayout,
		mkResourceVirtualEnvironmentVMMemory,
		mkResourceVirtualEnvironmentVMName,
		mkResourceVirtualEnvironmentVMNetworkDevice,
		mkResourceVirtualEnvironmentVMOperatingSystem,
		mkResourceVirtualEnvironmentVMPoolID,
		mkResourceVirtualEnvironmentVMSerialDevice,
		mkResourceVirtualEnvironmentVMStarted,
		mkResourceVirtualEnvironmentVMTabletDevice,
		mkResourceVirtualEnvironmentVMVMID,
	})

	testComputedAttributes(t, s, []string{
		mkResourceVirtualEnvironmentVMIPv4Addresses,
		mkResourceVirtualEnvironmentVMIPv6Addresses,
		mkResourceVirtualEnvironmentVMMACAddresses,
		mkResourceVirtualEnvironmentVMNetworkInterfaceNames,
	})

	testSchemaValueTypes(t, s, []string{
		mkResourceVirtualEnvironmentVMACPI,
		mkResourceVirtualEnvironmentVMAgent,
		mkResourceVirtualEnvironmentVMAudioDevice,
		mkResourceVirtualEnvironmentVMBIOS,
		mkResourceVirtualEnvironmentVMCDROM,
		mkResourceVirtualEnvironmentVMInitialization,
		mkResourceVirtualEnvironmentVMCPU,
		mkResourceVirtualEnvironmentVMDescription,
		mkResourceVirtualEnvironmentVMDisk,
		mkResourceVirtualEnvironmentVMIPv4Addresses,
		mkResourceVirtualEnvironmentVMIPv6Addresses,
		mkResourceVirtualEnvironmentVMKeyboardLayout,
		mkResourceVirtualEnvironmentVMMemory,
		mkResourceVirtualEnvironmentVMName,
		mkResourceVirtualEnvironmentVMNetworkDevice,
		mkResourceVirtualEnvironmentVMMACAddresses,
		mkResourceVirtualEnvironmentVMNetworkInterfaceNames,
		mkResourceVirtualEnvironmentVMOperatingSystem,
		mkResourceVirtualEnvironmentVMPoolID,
		mkResourceVirtualEnvironmentVMSerialDevice,
		mkResourceVirtualEnvironmentVMStarted,
		mkResourceVirtualEnvironmentVMTabletDevice,
		mkResourceVirtualEnvironmentVMVMID,
	}, []schema.ValueType{
		schema.TypeBool,
		schema.TypeList,
		schema.TypeList,
		schema.TypeString,
		schema.TypeList,
		schema.TypeList,
		schema.TypeList,
		schema.TypeString,
		schema.TypeList,
		schema.TypeList,
		schema.TypeList,
		schema.TypeString,
		schema.TypeList,
		schema.TypeString,
		schema.TypeList,
		schema.TypeList,
		schema.TypeList,
		schema.TypeList,
		schema.TypeString,
		schema.TypeList,
		schema.TypeBool,
		schema.TypeBool,
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

	audioDeviceSchema := testNestedSchemaExistence(t, s, mkResourceVirtualEnvironmentVMAudioDevice)

	testOptionalArguments(t, audioDeviceSchema, []string{
		mkResourceVirtualEnvironmentVMAudioDeviceDevice,
		mkResourceVirtualEnvironmentVMAudioDeviceDriver,
	})

	testSchemaValueTypes(t, audioDeviceSchema, []string{
		mkResourceVirtualEnvironmentVMAudioDeviceDevice,
		mkResourceVirtualEnvironmentVMAudioDeviceDriver,
	}, []schema.ValueType{
		schema.TypeString,
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

	cpuSchema := testNestedSchemaExistence(t, s, mkResourceVirtualEnvironmentVMCPU)

	testOptionalArguments(t, cpuSchema, []string{
		mkResourceVirtualEnvironmentVMCPUArchitecture,
		mkResourceVirtualEnvironmentVMCPUCores,
		mkResourceVirtualEnvironmentVMCPUFlags,
		mkResourceVirtualEnvironmentVMCPUHotplugged,
		mkResourceVirtualEnvironmentVMCPUSockets,
		mkResourceVirtualEnvironmentVMCPUType,
		mkResourceVirtualEnvironmentVMCPUUnits,
	})

	testSchemaValueTypes(t, cpuSchema, []string{
		mkResourceVirtualEnvironmentVMCPUArchitecture,
		mkResourceVirtualEnvironmentVMCPUCores,
		mkResourceVirtualEnvironmentVMCPUFlags,
		mkResourceVirtualEnvironmentVMCPUHotplugged,
		mkResourceVirtualEnvironmentVMCPUSockets,
		mkResourceVirtualEnvironmentVMCPUType,
		mkResourceVirtualEnvironmentVMCPUUnits,
	}, []schema.ValueType{
		schema.TypeString,
		schema.TypeInt,
		schema.TypeList,
		schema.TypeInt,
		schema.TypeInt,
		schema.TypeString,
		schema.TypeInt,
	})

	diskSchema := testNestedSchemaExistence(t, s, mkResourceVirtualEnvironmentVMDisk)

	testOptionalArguments(t, diskSchema, []string{
		mkResourceVirtualEnvironmentVMDiskDatastoreID,
		mkResourceVirtualEnvironmentVMDiskFileFormat,
		mkResourceVirtualEnvironmentVMDiskFileID,
		mkResourceVirtualEnvironmentVMDiskSize,
	})

	testSchemaValueTypes(t, diskSchema, []string{
		mkResourceVirtualEnvironmentVMDiskDatastoreID,
		mkResourceVirtualEnvironmentVMDiskFileFormat,
		mkResourceVirtualEnvironmentVMDiskFileID,
		mkResourceVirtualEnvironmentVMDiskSize,
	}, []schema.ValueType{
		schema.TypeString,
		schema.TypeString,
		schema.TypeString,
		schema.TypeInt,
	})

	diskSpeedSchema := testNestedSchemaExistence(t, diskSchema, mkResourceVirtualEnvironmentVMDiskSpeed)

	testOptionalArguments(t, diskSpeedSchema, []string{
		mkResourceVirtualEnvironmentVMDiskSpeedRead,
		mkResourceVirtualEnvironmentVMDiskSpeedReadBurstable,
		mkResourceVirtualEnvironmentVMDiskSpeedWrite,
		mkResourceVirtualEnvironmentVMDiskSpeedWriteBurstable,
	})

	testSchemaValueTypes(t, diskSpeedSchema, []string{
		mkResourceVirtualEnvironmentVMDiskSpeedRead,
		mkResourceVirtualEnvironmentVMDiskSpeedReadBurstable,
		mkResourceVirtualEnvironmentVMDiskSpeedWrite,
		mkResourceVirtualEnvironmentVMDiskSpeedWriteBurstable,
	}, []schema.ValueType{
		schema.TypeInt,
		schema.TypeInt,
		schema.TypeInt,
		schema.TypeInt,
	})

	initializationSchema := testNestedSchemaExistence(t, s, mkResourceVirtualEnvironmentVMInitialization)

	testOptionalArguments(t, initializationSchema, []string{
		mkResourceVirtualEnvironmentVMInitializationDNS,
		mkResourceVirtualEnvironmentVMInitializationIPConfig,
		mkResourceVirtualEnvironmentVMInitializationUserAccount,
	})

	testSchemaValueTypes(t, initializationSchema, []string{
		mkResourceVirtualEnvironmentVMInitializationDNS,
		mkResourceVirtualEnvironmentVMInitializationIPConfig,
		mkResourceVirtualEnvironmentVMInitializationUserAccount,
	}, []schema.ValueType{
		schema.TypeList,
		schema.TypeList,
		schema.TypeList,
	})

	initializationDNSSchema := testNestedSchemaExistence(t, initializationSchema, mkResourceVirtualEnvironmentVMInitializationDNS)

	testOptionalArguments(t, initializationDNSSchema, []string{
		mkResourceVirtualEnvironmentVMInitializationDNSDomain,
		mkResourceVirtualEnvironmentVMInitializationDNSServer,
	})

	testSchemaValueTypes(t, initializationDNSSchema, []string{
		mkResourceVirtualEnvironmentVMInitializationDNSDomain,
		mkResourceVirtualEnvironmentVMInitializationDNSServer,
	}, []schema.ValueType{
		schema.TypeString,
		schema.TypeString,
	})

	initializationIPConfigSchema := testNestedSchemaExistence(t, initializationSchema, mkResourceVirtualEnvironmentVMInitializationIPConfig)

	testOptionalArguments(t, initializationIPConfigSchema, []string{
		mkResourceVirtualEnvironmentVMInitializationIPConfigIPv4,
		mkResourceVirtualEnvironmentVMInitializationIPConfigIPv6,
	})

	testSchemaValueTypes(t, initializationIPConfigSchema, []string{
		mkResourceVirtualEnvironmentVMInitializationIPConfigIPv4,
		mkResourceVirtualEnvironmentVMInitializationIPConfigIPv6,
	}, []schema.ValueType{
		schema.TypeList,
		schema.TypeList,
	})

	initializationIPConfigIPv4Schema := testNestedSchemaExistence(t, initializationIPConfigSchema, mkResourceVirtualEnvironmentVMInitializationIPConfigIPv4)

	testOptionalArguments(t, initializationIPConfigIPv4Schema, []string{
		mkResourceVirtualEnvironmentVMInitializationIPConfigIPv4Address,
		mkResourceVirtualEnvironmentVMInitializationIPConfigIPv4Gateway,
	})

	testSchemaValueTypes(t, initializationIPConfigIPv4Schema, []string{
		mkResourceVirtualEnvironmentVMInitializationIPConfigIPv4Address,
		mkResourceVirtualEnvironmentVMInitializationIPConfigIPv4Gateway,
	}, []schema.ValueType{
		schema.TypeString,
		schema.TypeString,
	})

	initializationIPConfigIPv6Schema := testNestedSchemaExistence(t, initializationIPConfigSchema, mkResourceVirtualEnvironmentVMInitializationIPConfigIPv6)

	testOptionalArguments(t, initializationIPConfigIPv6Schema, []string{
		mkResourceVirtualEnvironmentVMInitializationIPConfigIPv6Address,
		mkResourceVirtualEnvironmentVMInitializationIPConfigIPv6Gateway,
	})

	testSchemaValueTypes(t, initializationIPConfigIPv6Schema, []string{
		mkResourceVirtualEnvironmentVMInitializationIPConfigIPv6Address,
		mkResourceVirtualEnvironmentVMInitializationIPConfigIPv6Gateway,
	}, []schema.ValueType{
		schema.TypeString,
		schema.TypeString,
	})

	initializationUserAccountSchema := testNestedSchemaExistence(t, initializationSchema, mkResourceVirtualEnvironmentVMInitializationUserAccount)

	testOptionalArguments(t, initializationUserAccountSchema, []string{
		mkResourceVirtualEnvironmentVMInitializationUserAccountKeys,
		mkResourceVirtualEnvironmentVMInitializationUserAccountPassword,
		mkResourceVirtualEnvironmentVMInitializationUserAccountUsername,
	})

	testSchemaValueTypes(t, initializationUserAccountSchema, []string{
		mkResourceVirtualEnvironmentVMInitializationUserAccountKeys,
		mkResourceVirtualEnvironmentVMInitializationUserAccountPassword,
		mkResourceVirtualEnvironmentVMInitializationUserAccountUsername,
	}, []schema.ValueType{
		schema.TypeList,
		schema.TypeString,
		schema.TypeString,
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
		mkResourceVirtualEnvironmentVMNetworkDeviceRateLimit,
		mkResourceVirtualEnvironmentVMNetworkDeviceVLANID,
	})

	testSchemaValueTypes(t, networkDeviceSchema, []string{
		mkResourceVirtualEnvironmentVMNetworkDeviceBridge,
		mkResourceVirtualEnvironmentVMNetworkDeviceEnabled,
		mkResourceVirtualEnvironmentVMNetworkDeviceMACAddress,
		mkResourceVirtualEnvironmentVMNetworkDeviceModel,
		mkResourceVirtualEnvironmentVMNetworkDeviceRateLimit,
		mkResourceVirtualEnvironmentVMNetworkDeviceVLANID,
	}, []schema.ValueType{
		schema.TypeString,
		schema.TypeBool,
		schema.TypeString,
		schema.TypeString,
		schema.TypeFloat,
		schema.TypeInt,
	})

	operatingSystemSchema := testNestedSchemaExistence(t, s, mkResourceVirtualEnvironmentVMOperatingSystem)

	testOptionalArguments(t, operatingSystemSchema, []string{
		mkResourceVirtualEnvironmentVMOperatingSystemType,
	})

	testSchemaValueTypes(t, operatingSystemSchema, []string{
		mkResourceVirtualEnvironmentVMOperatingSystemType,
	}, []schema.ValueType{
		schema.TypeString,
	})

	serialDeviceSchema := testNestedSchemaExistence(t, s, mkResourceVirtualEnvironmentVMSerialDevice)

	testOptionalArguments(t, serialDeviceSchema, []string{
		mkResourceVirtualEnvironmentVMSerialDeviceDevice,
	})

	testSchemaValueTypes(t, serialDeviceSchema, []string{
		mkResourceVirtualEnvironmentVMSerialDeviceDevice,
	}, []schema.ValueType{
		schema.TypeString,
	})

	vgaSchema := testNestedSchemaExistence(t, s, mkResourceVirtualEnvironmentVMVGA)

	testOptionalArguments(t, vgaSchema, []string{
		mkResourceVirtualEnvironmentVMVGAEnabled,
		mkResourceVirtualEnvironmentVMVGAMemory,
		mkResourceVirtualEnvironmentVMVGAType,
	})

	testSchemaValueTypes(t, vgaSchema, []string{
		mkResourceVirtualEnvironmentVMVGAEnabled,
		mkResourceVirtualEnvironmentVMVGAMemory,
		mkResourceVirtualEnvironmentVMVGAType,
	}, []schema.ValueType{
		schema.TypeBool,
		schema.TypeInt,
		schema.TypeString,
	})
}
