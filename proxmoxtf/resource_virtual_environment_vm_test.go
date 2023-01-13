/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmoxtf

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
		mkResourceVirtualEnvironmentVMHostPCI,
		mkResourceVirtualEnvironmentVMKeyboardLayout,
		mkResourceVirtualEnvironmentVMKVMArguments,
		mkResourceVirtualEnvironmentVMMachine,
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

	testValueTypes(t, s, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentVMACPI:                  schema.TypeBool,
		mkResourceVirtualEnvironmentVMAgent:                 schema.TypeList,
		mkResourceVirtualEnvironmentVMAudioDevice:           schema.TypeList,
		mkResourceVirtualEnvironmentVMBIOS:                  schema.TypeString,
		mkResourceVirtualEnvironmentVMCDROM:                 schema.TypeList,
		mkResourceVirtualEnvironmentVMCPU:                   schema.TypeList,
		mkResourceVirtualEnvironmentVMDescription:           schema.TypeString,
		mkResourceVirtualEnvironmentVMDisk:                  schema.TypeList,
		mkResourceVirtualEnvironmentVMHostPCI:               schema.TypeList,
		mkResourceVirtualEnvironmentVMInitialization:        schema.TypeList,
		mkResourceVirtualEnvironmentVMIPv4Addresses:         schema.TypeList,
		mkResourceVirtualEnvironmentVMIPv6Addresses:         schema.TypeList,
		mkResourceVirtualEnvironmentVMKeyboardLayout:        schema.TypeString,
		mkResourceVirtualEnvironmentVMKVMArguments:          schema.TypeString,
		mkResourceVirtualEnvironmentVMMachine:               schema.TypeString,
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
		mkResourceVirtualEnvironmentVMAgentTimeout,
		mkResourceVirtualEnvironmentVMAgentTrim,
		mkResourceVirtualEnvironmentVMAgentType,
	})

	testValueTypes(t, agentSchema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentVMAgentEnabled: schema.TypeBool,
		mkResourceVirtualEnvironmentVMAgentTrim:    schema.TypeBool,
		mkResourceVirtualEnvironmentVMAgentType:    schema.TypeString,
	})

	audioDeviceSchema := testNestedSchemaExistence(t, s, mkResourceVirtualEnvironmentVMAudioDevice)

	testOptionalArguments(t, audioDeviceSchema, []string{
		mkResourceVirtualEnvironmentVMAudioDeviceDevice,
		mkResourceVirtualEnvironmentVMAudioDeviceDriver,
	})

	testValueTypes(t, audioDeviceSchema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentVMAudioDeviceDevice: schema.TypeString,
		mkResourceVirtualEnvironmentVMAudioDeviceDriver: schema.TypeString,
	})

	cdromSchema := testNestedSchemaExistence(t, s, mkResourceVirtualEnvironmentVMCDROM)

	testOptionalArguments(t, cdromSchema, []string{
		mkResourceVirtualEnvironmentVMCDROMEnabled,
		mkResourceVirtualEnvironmentVMCDROMFileID,
	})

	testValueTypes(t, cdromSchema, map[string]schema.ValueType{
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

	testValueTypes(t, cloneSchema, map[string]schema.ValueType{
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

	testValueTypes(t, cpuSchema, map[string]schema.ValueType{
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

	testValueTypes(t, diskSchema, map[string]schema.ValueType{
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

	testValueTypes(t, diskSpeedSchema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentVMDiskSpeedRead:           schema.TypeInt,
		mkResourceVirtualEnvironmentVMDiskSpeedReadBurstable:  schema.TypeInt,
		mkResourceVirtualEnvironmentVMDiskSpeedWrite:          schema.TypeInt,
		mkResourceVirtualEnvironmentVMDiskSpeedWriteBurstable: schema.TypeInt,
	})

	initializationSchema := testNestedSchemaExistence(t, s, mkResourceVirtualEnvironmentVMInitialization)

	testOptionalArguments(t, initializationSchema, []string{
		mkResourceVirtualEnvironmentVMInitializationDatastoreID,
		mkResourceVirtualEnvironmentVMInitializationDNS,
		mkResourceVirtualEnvironmentVMInitializationIPConfig,
		mkResourceVirtualEnvironmentVMInitializationUserAccount,
	})

	testValueTypes(t, initializationSchema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentVMInitializationDatastoreID: schema.TypeString,
		mkResourceVirtualEnvironmentVMInitializationDNS:         schema.TypeList,
		mkResourceVirtualEnvironmentVMInitializationIPConfig:    schema.TypeList,
		mkResourceVirtualEnvironmentVMInitializationUserAccount: schema.TypeList,
	})

	hostPCISchema := testNestedSchemaExistence(t, s, mkResourceVirtualEnvironmentVMHostPCI)

	testOptionalArguments(t, hostPCISchema, []string{
		mkResourceVirtualEnvironmentVMHostPCIDeviceMDev,
		mkResourceVirtualEnvironmentVMHostPCIDevicePCIE,
		mkResourceVirtualEnvironmentVMHostPCIDeviceROMBAR,
		mkResourceVirtualEnvironmentVMHostPCIDeviceROMFile,
		mkResourceVirtualEnvironmentVMHostPCIDeviceXVGA,
	})

	testValueTypes(t, hostPCISchema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentVMHostPCIDevice:        schema.TypeString,
		mkResourceVirtualEnvironmentVMHostPCIDeviceMDev:    schema.TypeString,
		mkResourceVirtualEnvironmentVMHostPCIDevicePCIE:    schema.TypeBool,
		mkResourceVirtualEnvironmentVMHostPCIDeviceROMBAR:  schema.TypeBool,
		mkResourceVirtualEnvironmentVMHostPCIDeviceROMFile: schema.TypeString,
		mkResourceVirtualEnvironmentVMHostPCIDeviceXVGA:    schema.TypeBool,
	})

	initializationDNSSchema := testNestedSchemaExistence(t, initializationSchema, mkResourceVirtualEnvironmentVMInitializationDNS)

	testOptionalArguments(t, initializationDNSSchema, []string{
		mkResourceVirtualEnvironmentVMInitializationDNSDomain,
		mkResourceVirtualEnvironmentVMInitializationDNSServer,
	})

	testValueTypes(t, initializationDNSSchema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentVMInitializationDNSDomain: schema.TypeString,
		mkResourceVirtualEnvironmentVMInitializationDNSServer: schema.TypeString,
	})

	initializationIPConfigSchema := testNestedSchemaExistence(t, initializationSchema, mkResourceVirtualEnvironmentVMInitializationIPConfig)

	testOptionalArguments(t, initializationIPConfigSchema, []string{
		mkResourceVirtualEnvironmentVMInitializationIPConfigIPv4,
		mkResourceVirtualEnvironmentVMInitializationIPConfigIPv6,
	})

	testValueTypes(t, initializationIPConfigSchema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentVMInitializationIPConfigIPv4: schema.TypeList,
		mkResourceVirtualEnvironmentVMInitializationIPConfigIPv6: schema.TypeList,
	})

	initializationIPConfigIPv4Schema := testNestedSchemaExistence(t, initializationIPConfigSchema, mkResourceVirtualEnvironmentVMInitializationIPConfigIPv4)

	testOptionalArguments(t, initializationIPConfigIPv4Schema, []string{
		mkResourceVirtualEnvironmentVMInitializationIPConfigIPv4Address,
		mkResourceVirtualEnvironmentVMInitializationIPConfigIPv4Gateway,
	})

	testValueTypes(t, initializationIPConfigIPv4Schema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentVMInitializationIPConfigIPv4Address: schema.TypeString,
		mkResourceVirtualEnvironmentVMInitializationIPConfigIPv4Gateway: schema.TypeString,
	})

	initializationIPConfigIPv6Schema := testNestedSchemaExistence(t, initializationIPConfigSchema, mkResourceVirtualEnvironmentVMInitializationIPConfigIPv6)

	testOptionalArguments(t, initializationIPConfigIPv6Schema, []string{
		mkResourceVirtualEnvironmentVMInitializationIPConfigIPv6Address,
		mkResourceVirtualEnvironmentVMInitializationIPConfigIPv6Gateway,
	})

	testValueTypes(t, initializationIPConfigIPv6Schema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentVMInitializationIPConfigIPv6Address: schema.TypeString,
		mkResourceVirtualEnvironmentVMInitializationIPConfigIPv6Gateway: schema.TypeString,
	})

	initializationUserAccountSchema := testNestedSchemaExistence(t, initializationSchema, mkResourceVirtualEnvironmentVMInitializationUserAccount)

	testOptionalArguments(t, initializationUserAccountSchema, []string{
		mkResourceVirtualEnvironmentVMInitializationUserAccountKeys,
		mkResourceVirtualEnvironmentVMInitializationUserAccountPassword,
		mkResourceVirtualEnvironmentVMInitializationUserAccountUsername,
	})

	testValueTypes(t, initializationUserAccountSchema, map[string]schema.ValueType{
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

	testValueTypes(t, memorySchema, map[string]schema.ValueType{
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
		mkResourceVirtualEnvironmentVMNetworkDeviceMTU,
	})

	testValueTypes(t, networkDeviceSchema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentVMNetworkDeviceBridge:     schema.TypeString,
		mkResourceVirtualEnvironmentVMNetworkDeviceEnabled:    schema.TypeBool,
		mkResourceVirtualEnvironmentVMNetworkDeviceMACAddress: schema.TypeString,
		mkResourceVirtualEnvironmentVMNetworkDeviceModel:      schema.TypeString,
		mkResourceVirtualEnvironmentVMNetworkDeviceRateLimit:  schema.TypeFloat,
		mkResourceVirtualEnvironmentVMNetworkDeviceVLANID:     schema.TypeInt,
		mkResourceVirtualEnvironmentVMNetworkDeviceMTU:        schema.TypeInt,
	})

	operatingSystemSchema := testNestedSchemaExistence(t, s, mkResourceVirtualEnvironmentVMOperatingSystem)

	testOptionalArguments(t, operatingSystemSchema, []string{
		mkResourceVirtualEnvironmentVMOperatingSystemType,
	})

	testValueTypes(t, operatingSystemSchema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentVMOperatingSystemType: schema.TypeString,
	})

	serialDeviceSchema := testNestedSchemaExistence(t, s, mkResourceVirtualEnvironmentVMSerialDevice)

	testOptionalArguments(t, serialDeviceSchema, []string{
		mkResourceVirtualEnvironmentVMSerialDeviceDevice,
	})

	testValueTypes(t, serialDeviceSchema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentVMSerialDeviceDevice: schema.TypeString,
	})

	vgaSchema := testNestedSchemaExistence(t, s, mkResourceVirtualEnvironmentVMVGA)

	testOptionalArguments(t, vgaSchema, []string{
		mkResourceVirtualEnvironmentVMVGAEnabled,
		mkResourceVirtualEnvironmentVMVGAMemory,
		mkResourceVirtualEnvironmentVMVGAType,
	})

	testValueTypes(t, vgaSchema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentVMVGAEnabled: schema.TypeBool,
		mkResourceVirtualEnvironmentVMVGAMemory:  schema.TypeInt,
		mkResourceVirtualEnvironmentVMVGAType:    schema.TypeString,
	})
}
