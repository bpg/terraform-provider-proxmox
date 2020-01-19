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
		mkResourceVirtualEnvironmentVMClone,
		mkResourceVirtualEnvironmentVMCPU,
		mkResourceVirtualEnvironmentVMDescription,
		mkResourceVirtualEnvironmentVMDisk,
		mkResourceVirtualEnvironmentVMInitialization,
		mkResourceVirtualEnvironmentVMKeyboardLayout,
		mkResourceVirtualEnvironmentVMMemory,
		mkResourceVirtualEnvironmentVMName,
		mkResourceVirtualEnvironmentVMNetworkDevice,
		mkResourceVirtualEnvironmentVMOperatingSystem,
		mkResourceVirtualEnvironmentVMPoolID,
		mkResourceVirtualEnvironmentVMSerialDevice,
		mkResourceVirtualEnvironmentVMStarted,
		mkResourceVirtualEnvironmentVMTabletDevice,
		mkResourceVirtualEnvironmentVMTemplate,
		mkResourceVirtualEnvironmentVMVMID,
	})

	testComputedAttributes(t, s, []string{
		mkResourceVirtualEnvironmentVMIPv4Addresses,
		mkResourceVirtualEnvironmentVMIPv6Addresses,
		mkResourceVirtualEnvironmentVMMACAddresses,
		mkResourceVirtualEnvironmentVMNetworkInterfaceNames,
	})

	testSchemaValueTypes(t, s, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentVMACPI:                  schema.TypeBool,
		mkResourceVirtualEnvironmentVMAgent:                 schema.TypeList,
		mkResourceVirtualEnvironmentVMAudioDevice:           schema.TypeList,
		mkResourceVirtualEnvironmentVMBIOS:                  schema.TypeString,
		mkResourceVirtualEnvironmentVMCDROM:                 schema.TypeList,
		mkResourceVirtualEnvironmentVMCPU:                   schema.TypeList,
		mkResourceVirtualEnvironmentVMDescription:           schema.TypeString,
		mkResourceVirtualEnvironmentVMDisk:                  schema.TypeList,
		mkResourceVirtualEnvironmentVMInitialization:        schema.TypeList,
		mkResourceVirtualEnvironmentVMIPv4Addresses:         schema.TypeList,
		mkResourceVirtualEnvironmentVMIPv6Addresses:         schema.TypeList,
		mkResourceVirtualEnvironmentVMKeyboardLayout:        schema.TypeString,
		mkResourceVirtualEnvironmentVMMemory:                schema.TypeList,
		mkResourceVirtualEnvironmentVMName:                  schema.TypeString,
		mkResourceVirtualEnvironmentVMNetworkDevice:         schema.TypeList,
		mkResourceVirtualEnvironmentVMMACAddresses:          schema.TypeList,
		mkResourceVirtualEnvironmentVMNetworkInterfaceNames: schema.TypeList,
		mkResourceVirtualEnvironmentVMOperatingSystem:       schema.TypeList,
		mkResourceVirtualEnvironmentVMPoolID:                schema.TypeString,
		mkResourceVirtualEnvironmentVMSerialDevice:          schema.TypeList,
		mkResourceVirtualEnvironmentVMStarted:               schema.TypeBool,
		mkResourceVirtualEnvironmentVMTabletDevice:          schema.TypeBool,
		mkResourceVirtualEnvironmentVMTemplate:              schema.TypeBool,
		mkResourceVirtualEnvironmentVMVMID:                  schema.TypeInt,
	})

	agentSchema := testNestedSchemaExistence(t, s, mkResourceVirtualEnvironmentVMAgent)

	testOptionalArguments(t, agentSchema, []string{
		mkResourceVirtualEnvironmentVMAgentEnabled,
		mkResourceVirtualEnvironmentVMAgentTrim,
		mkResourceVirtualEnvironmentVMAgentType,
	})

	testSchemaValueTypes(t, agentSchema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentVMAgentEnabled: schema.TypeBool,
		mkResourceVirtualEnvironmentVMAgentTrim:    schema.TypeBool,
		mkResourceVirtualEnvironmentVMAgentType:    schema.TypeString,
	})

	audioDeviceSchema := testNestedSchemaExistence(t, s, mkResourceVirtualEnvironmentVMAudioDevice)

	testOptionalArguments(t, audioDeviceSchema, []string{
		mkResourceVirtualEnvironmentVMAudioDeviceDevice,
		mkResourceVirtualEnvironmentVMAudioDeviceDriver,
	})

	testSchemaValueTypes(t, audioDeviceSchema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentVMAudioDeviceDevice: schema.TypeString,
		mkResourceVirtualEnvironmentVMAudioDeviceDriver: schema.TypeString,
	})

	cdromSchema := testNestedSchemaExistence(t, s, mkResourceVirtualEnvironmentVMCDROM)

	testOptionalArguments(t, cdromSchema, []string{
		mkResourceVirtualEnvironmentVMCDROMEnabled,
		mkResourceVirtualEnvironmentVMCDROMFileID,
	})

	testSchemaValueTypes(t, cdromSchema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentVMCDROMEnabled: schema.TypeBool,
		mkResourceVirtualEnvironmentVMCDROMFileID:  schema.TypeString,
	})

	cloneSchema := testNestedSchemaExistence(t, s, mkResourceVirtualEnvironmentVMClone)

	testRequiredArguments(t, cloneSchema, []string{
		mkResourceVirtualEnvironmentVMCloneVMID,
	})

	testOptionalArguments(t, cloneSchema, []string{
		mkResourceVirtualEnvironmentVMCloneDatastoreID,
		mkResourceVirtualEnvironmentVMCloneNodeName,
	})

	testSchemaValueTypes(t, cloneSchema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentVMCloneDatastoreID: schema.TypeString,
		mkResourceVirtualEnvironmentVMCloneNodeName:    schema.TypeString,
		mkResourceVirtualEnvironmentVMCloneVMID:        schema.TypeInt,
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

	testSchemaValueTypes(t, cpuSchema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentVMCPUArchitecture: schema.TypeString,
		mkResourceVirtualEnvironmentVMCPUCores:        schema.TypeInt,
		mkResourceVirtualEnvironmentVMCPUFlags:        schema.TypeList,
		mkResourceVirtualEnvironmentVMCPUHotplugged:   schema.TypeInt,
		mkResourceVirtualEnvironmentVMCPUSockets:      schema.TypeInt,
		mkResourceVirtualEnvironmentVMCPUType:         schema.TypeString,
		mkResourceVirtualEnvironmentVMCPUUnits:        schema.TypeInt,
	})

	diskSchema := testNestedSchemaExistence(t, s, mkResourceVirtualEnvironmentVMDisk)

	testOptionalArguments(t, diskSchema, []string{
		mkResourceVirtualEnvironmentVMDiskDatastoreID,
		mkResourceVirtualEnvironmentVMDiskFileFormat,
		mkResourceVirtualEnvironmentVMDiskFileID,
		mkResourceVirtualEnvironmentVMDiskSize,
	})

	testSchemaValueTypes(t, diskSchema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentVMDiskDatastoreID: schema.TypeString,
		mkResourceVirtualEnvironmentVMDiskFileFormat:  schema.TypeString,
		mkResourceVirtualEnvironmentVMDiskFileID:      schema.TypeString,
		mkResourceVirtualEnvironmentVMDiskSize:        schema.TypeInt,
	})

	diskSpeedSchema := testNestedSchemaExistence(t, diskSchema, mkResourceVirtualEnvironmentVMDiskSpeed)

	testOptionalArguments(t, diskSpeedSchema, []string{
		mkResourceVirtualEnvironmentVMDiskSpeedRead,
		mkResourceVirtualEnvironmentVMDiskSpeedReadBurstable,
		mkResourceVirtualEnvironmentVMDiskSpeedWrite,
		mkResourceVirtualEnvironmentVMDiskSpeedWriteBurstable,
	})

	testSchemaValueTypes(t, diskSpeedSchema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentVMDiskSpeedRead:           schema.TypeInt,
		mkResourceVirtualEnvironmentVMDiskSpeedReadBurstable:  schema.TypeInt,
		mkResourceVirtualEnvironmentVMDiskSpeedWrite:          schema.TypeInt,
		mkResourceVirtualEnvironmentVMDiskSpeedWriteBurstable: schema.TypeInt,
	})

	initializationSchema := testNestedSchemaExistence(t, s, mkResourceVirtualEnvironmentVMInitialization)

	testOptionalArguments(t, initializationSchema, []string{
		mkResourceVirtualEnvironmentVMInitializationDNS,
		mkResourceVirtualEnvironmentVMInitializationIPConfig,
		mkResourceVirtualEnvironmentVMInitializationUserAccount,
	})

	testSchemaValueTypes(t, initializationSchema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentVMInitializationDNS:         schema.TypeList,
		mkResourceVirtualEnvironmentVMInitializationIPConfig:    schema.TypeList,
		mkResourceVirtualEnvironmentVMInitializationUserAccount: schema.TypeList,
	})

	initializationDNSSchema := testNestedSchemaExistence(t, initializationSchema, mkResourceVirtualEnvironmentVMInitializationDNS)

	testOptionalArguments(t, initializationDNSSchema, []string{
		mkResourceVirtualEnvironmentVMInitializationDNSDomain,
		mkResourceVirtualEnvironmentVMInitializationDNSServer,
	})

	testSchemaValueTypes(t, initializationDNSSchema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentVMInitializationDNSDomain: schema.TypeString,
		mkResourceVirtualEnvironmentVMInitializationDNSServer: schema.TypeString,
	})

	initializationIPConfigSchema := testNestedSchemaExistence(t, initializationSchema, mkResourceVirtualEnvironmentVMInitializationIPConfig)

	testOptionalArguments(t, initializationIPConfigSchema, []string{
		mkResourceVirtualEnvironmentVMInitializationIPConfigIPv4,
		mkResourceVirtualEnvironmentVMInitializationIPConfigIPv6,
	})

	testSchemaValueTypes(t, initializationIPConfigSchema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentVMInitializationIPConfigIPv4: schema.TypeList,
		mkResourceVirtualEnvironmentVMInitializationIPConfigIPv6: schema.TypeList,
	})

	initializationIPConfigIPv4Schema := testNestedSchemaExistence(t, initializationIPConfigSchema, mkResourceVirtualEnvironmentVMInitializationIPConfigIPv4)

	testOptionalArguments(t, initializationIPConfigIPv4Schema, []string{
		mkResourceVirtualEnvironmentVMInitializationIPConfigIPv4Address,
		mkResourceVirtualEnvironmentVMInitializationIPConfigIPv4Gateway,
	})

	testSchemaValueTypes(t, initializationIPConfigIPv4Schema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentVMInitializationIPConfigIPv4Address: schema.TypeString,
		mkResourceVirtualEnvironmentVMInitializationIPConfigIPv4Gateway: schema.TypeString,
	})

	initializationIPConfigIPv6Schema := testNestedSchemaExistence(t, initializationIPConfigSchema, mkResourceVirtualEnvironmentVMInitializationIPConfigIPv6)

	testOptionalArguments(t, initializationIPConfigIPv6Schema, []string{
		mkResourceVirtualEnvironmentVMInitializationIPConfigIPv6Address,
		mkResourceVirtualEnvironmentVMInitializationIPConfigIPv6Gateway,
	})

	testSchemaValueTypes(t, initializationIPConfigIPv6Schema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentVMInitializationIPConfigIPv6Address: schema.TypeString,
		mkResourceVirtualEnvironmentVMInitializationIPConfigIPv6Gateway: schema.TypeString,
	})

	initializationUserAccountSchema := testNestedSchemaExistence(t, initializationSchema, mkResourceVirtualEnvironmentVMInitializationUserAccount)

	testOptionalArguments(t, initializationUserAccountSchema, []string{
		mkResourceVirtualEnvironmentVMInitializationUserAccountKeys,
		mkResourceVirtualEnvironmentVMInitializationUserAccountPassword,
		mkResourceVirtualEnvironmentVMInitializationUserAccountUsername,
	})

	testSchemaValueTypes(t, initializationUserAccountSchema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentVMInitializationUserAccountKeys:     schema.TypeList,
		mkResourceVirtualEnvironmentVMInitializationUserAccountPassword: schema.TypeString,
		mkResourceVirtualEnvironmentVMInitializationUserAccountUsername: schema.TypeString,
	})

	memorySchema := testNestedSchemaExistence(t, s, mkResourceVirtualEnvironmentVMMemory)

	testOptionalArguments(t, memorySchema, []string{
		mkResourceVirtualEnvironmentVMMemoryDedicated,
		mkResourceVirtualEnvironmentVMMemoryFloating,
		mkResourceVirtualEnvironmentVMMemoryShared,
	})

	testSchemaValueTypes(t, memorySchema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentVMMemoryDedicated: schema.TypeInt,
		mkResourceVirtualEnvironmentVMMemoryFloating:  schema.TypeInt,
		mkResourceVirtualEnvironmentVMMemoryShared:    schema.TypeInt,
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

	testSchemaValueTypes(t, networkDeviceSchema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentVMNetworkDeviceBridge:     schema.TypeString,
		mkResourceVirtualEnvironmentVMNetworkDeviceEnabled:    schema.TypeBool,
		mkResourceVirtualEnvironmentVMNetworkDeviceMACAddress: schema.TypeString,
		mkResourceVirtualEnvironmentVMNetworkDeviceModel:      schema.TypeString,
		mkResourceVirtualEnvironmentVMNetworkDeviceRateLimit:  schema.TypeFloat,
		mkResourceVirtualEnvironmentVMNetworkDeviceVLANID:     schema.TypeInt,
	})

	operatingSystemSchema := testNestedSchemaExistence(t, s, mkResourceVirtualEnvironmentVMOperatingSystem)

	testOptionalArguments(t, operatingSystemSchema, []string{
		mkResourceVirtualEnvironmentVMOperatingSystemType,
	})

	testSchemaValueTypes(t, operatingSystemSchema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentVMOperatingSystemType: schema.TypeString,
	})

	serialDeviceSchema := testNestedSchemaExistence(t, s, mkResourceVirtualEnvironmentVMSerialDevice)

	testOptionalArguments(t, serialDeviceSchema, []string{
		mkResourceVirtualEnvironmentVMSerialDeviceDevice,
	})

	testSchemaValueTypes(t, serialDeviceSchema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentVMSerialDeviceDevice: schema.TypeString,
	})

	vgaSchema := testNestedSchemaExistence(t, s, mkResourceVirtualEnvironmentVMVGA)

	testOptionalArguments(t, vgaSchema, []string{
		mkResourceVirtualEnvironmentVMVGAEnabled,
		mkResourceVirtualEnvironmentVMVGAMemory,
		mkResourceVirtualEnvironmentVMVGAType,
	})

	testSchemaValueTypes(t, vgaSchema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentVMVGAEnabled: schema.TypeBool,
		mkResourceVirtualEnvironmentVMVGAMemory:  schema.TypeInt,
		mkResourceVirtualEnvironmentVMVGAType:    schema.TypeString,
	})
}
