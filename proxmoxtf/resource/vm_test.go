/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package resource

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/test"
)

// TestVMInstantiation tests whether the VM instance can be instantiated.
func TestVMInstantiation(t *testing.T) {
	t.Parallel()

	s := VM()
	if s == nil {
		t.Fatalf("Cannot instantiate VM")
	}
}

// TestVMSchema tests the VM schema.
func TestVMSchema(t *testing.T) {
	t.Parallel()

	s := VM()

	test.AssertRequiredArguments(t, s, []string{
		mkResourceVirtualEnvironmentVMNodeName,
	})

	test.AssertOptionalArguments(t, s, []string{
		mkResourceVirtualEnvironmentVMACPI,
		mkResourceVirtualEnvironmentVMAgent,
		mkResourceVirtualEnvironmentVMAudioDevice,
		mkResourceVirtualEnvironmentVMBIOS,
		mkResourceVirtualEnvironmentVMBootOrder,
		mkResourceVirtualEnvironmentVMCDROM,
		mkResourceVirtualEnvironmentVMClone,
		mkResourceVirtualEnvironmentVMCPU,
		mkResourceVirtualEnvironmentVMDescription,
		mkResourceVirtualEnvironmentVMDisk,
		mkResourceVirtualEnvironmentVMEFIDisk,
		mkResourceVirtualEnvironmentVMInitialization,
		mkResourceVirtualEnvironmentVMHostPCI,
		mkResourceVirtualEnvironmentVMHostUSB,
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
		mkResourceVirtualEnvironmentVMSCSIHardware,
	})

	test.AssertComputedAttributes(t, s, []string{
		mkResourceVirtualEnvironmentVMIPv4Addresses,
		mkResourceVirtualEnvironmentVMIPv6Addresses,
		mkResourceVirtualEnvironmentVMMACAddresses,
		mkResourceVirtualEnvironmentVMNetworkInterfaceNames,
	})

	test.AssertValueTypes(t, s, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentVMACPI:                  schema.TypeBool,
		mkResourceVirtualEnvironmentVMAgent:                 schema.TypeList,
		mkResourceVirtualEnvironmentVMAudioDevice:           schema.TypeList,
		mkResourceVirtualEnvironmentVMBIOS:                  schema.TypeString,
		mkResourceVirtualEnvironmentVMBootOrder:             schema.TypeList,
		mkResourceVirtualEnvironmentVMCDROM:                 schema.TypeList,
		mkResourceVirtualEnvironmentVMCPU:                   schema.TypeList,
		mkResourceVirtualEnvironmentVMDescription:           schema.TypeString,
		mkResourceVirtualEnvironmentVMDisk:                  schema.TypeList,
		mkResourceVirtualEnvironmentVMEFIDisk:               schema.TypeList,
		mkResourceVirtualEnvironmentVMHostPCI:               schema.TypeList,
		mkResourceVirtualEnvironmentVMHostUSB:               schema.TypeList,
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
		mkResourceVirtualEnvironmentVMSCSIHardware:          schema.TypeString,
	})

	agentSchema := test.AssertNestedSchemaExistence(t, s, mkResourceVirtualEnvironmentVMAgent)

	test.AssertOptionalArguments(t, agentSchema, []string{
		mkResourceVirtualEnvironmentVMAgentEnabled,
		mkResourceVirtualEnvironmentVMAgentTimeout,
		mkResourceVirtualEnvironmentVMAgentTrim,
		mkResourceVirtualEnvironmentVMAgentType,
	})

	test.AssertValueTypes(t, agentSchema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentVMAgentEnabled: schema.TypeBool,
		mkResourceVirtualEnvironmentVMAgentTrim:    schema.TypeBool,
		mkResourceVirtualEnvironmentVMAgentType:    schema.TypeString,
	})

	audioDeviceSchema := test.AssertNestedSchemaExistence(t, s, mkResourceVirtualEnvironmentVMAudioDevice)

	test.AssertOptionalArguments(t, audioDeviceSchema, []string{
		mkResourceVirtualEnvironmentVMAudioDeviceDevice,
		mkResourceVirtualEnvironmentVMAudioDeviceDriver,
	})

	test.AssertValueTypes(t, audioDeviceSchema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentVMAudioDeviceDevice: schema.TypeString,
		mkResourceVirtualEnvironmentVMAudioDeviceDriver: schema.TypeString,
	})

	cdromSchema := test.AssertNestedSchemaExistence(t, s, mkResourceVirtualEnvironmentVMCDROM)

	test.AssertOptionalArguments(t, cdromSchema, []string{
		mkResourceVirtualEnvironmentVMCDROMEnabled,
		mkResourceVirtualEnvironmentVMCDROMFileID,
	})

	test.AssertValueTypes(t, cdromSchema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentVMCDROMEnabled: schema.TypeBool,
		mkResourceVirtualEnvironmentVMCDROMFileID:  schema.TypeString,
	})

	cloneSchema := test.AssertNestedSchemaExistence(t, s, mkResourceVirtualEnvironmentVMClone)

	test.AssertRequiredArguments(t, cloneSchema, []string{
		mkResourceVirtualEnvironmentVMCloneVMID,
	})

	test.AssertOptionalArguments(t, cloneSchema, []string{
		mkResourceVirtualEnvironmentVMCloneDatastoreID,
		mkResourceVirtualEnvironmentVMCloneNodeName,
	})

	test.AssertValueTypes(t, cloneSchema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentVMCloneDatastoreID: schema.TypeString,
		mkResourceVirtualEnvironmentVMCloneNodeName:    schema.TypeString,
		mkResourceVirtualEnvironmentVMCloneVMID:        schema.TypeInt,
	})

	cpuSchema := test.AssertNestedSchemaExistence(t, s, mkResourceVirtualEnvironmentVMCPU)

	test.AssertOptionalArguments(t, cpuSchema, []string{
		mkResourceVirtualEnvironmentVMCPUArchitecture,
		mkResourceVirtualEnvironmentVMCPUCores,
		mkResourceVirtualEnvironmentVMCPUFlags,
		mkResourceVirtualEnvironmentVMCPUHotplugged,
		mkResourceVirtualEnvironmentVMCPUNUMA,
		mkResourceVirtualEnvironmentVMCPUSockets,
		mkResourceVirtualEnvironmentVMCPUType,
		mkResourceVirtualEnvironmentVMCPUUnits,
	})

	test.AssertValueTypes(t, cpuSchema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentVMCPUArchitecture: schema.TypeString,
		mkResourceVirtualEnvironmentVMCPUCores:        schema.TypeInt,
		mkResourceVirtualEnvironmentVMCPUFlags:        schema.TypeList,
		mkResourceVirtualEnvironmentVMCPUHotplugged:   schema.TypeInt,
		mkResourceVirtualEnvironmentVMCPUNUMA:         schema.TypeBool,
		mkResourceVirtualEnvironmentVMCPUSockets:      schema.TypeInt,
		mkResourceVirtualEnvironmentVMCPUType:         schema.TypeString,
		mkResourceVirtualEnvironmentVMCPUUnits:        schema.TypeInt,
	})

	diskSchema := test.AssertNestedSchemaExistence(t, s, mkResourceVirtualEnvironmentVMDisk)

	test.AssertOptionalArguments(t, diskSchema, []string{
		mkResourceVirtualEnvironmentVMDiskDatastoreID,
		mkResourceVirtualEnvironmentVMDiskPathInDatastore,
		mkResourceVirtualEnvironmentVMDiskFileFormat,
		mkResourceVirtualEnvironmentVMDiskFileID,
		mkResourceVirtualEnvironmentVMDiskSize,
	})

	test.AssertValueTypes(t, diskSchema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentVMDiskDatastoreID:     schema.TypeString,
		mkResourceVirtualEnvironmentVMDiskPathInDatastore: schema.TypeString,
		mkResourceVirtualEnvironmentVMDiskFileFormat:      schema.TypeString,
		mkResourceVirtualEnvironmentVMDiskFileID:          schema.TypeString,
		mkResourceVirtualEnvironmentVMDiskSize:            schema.TypeInt,
	})

	diskSpeedSchema := test.AssertNestedSchemaExistence(
		t,
		diskSchema,
		mkResourceVirtualEnvironmentVMDiskSpeed,
	)

	test.AssertOptionalArguments(t, diskSpeedSchema, []string{
		mkResourceVirtualEnvironmentVMDiskSpeedRead,
		mkResourceVirtualEnvironmentVMDiskSpeedReadBurstable,
		mkResourceVirtualEnvironmentVMDiskSpeedWrite,
		mkResourceVirtualEnvironmentVMDiskSpeedWriteBurstable,
	})

	test.AssertValueTypes(t, diskSpeedSchema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentVMDiskSpeedRead:           schema.TypeInt,
		mkResourceVirtualEnvironmentVMDiskSpeedReadBurstable:  schema.TypeInt,
		mkResourceVirtualEnvironmentVMDiskSpeedWrite:          schema.TypeInt,
		mkResourceVirtualEnvironmentVMDiskSpeedWriteBurstable: schema.TypeInt,
	})

	efiDiskSchema := test.AssertNestedSchemaExistence(t, s, mkResourceVirtualEnvironmentVMEFIDisk)

	test.AssertOptionalArguments(t, efiDiskSchema, []string{
		mkResourceVirtualEnvironmentVMEFIDiskDatastoreID,
		mkResourceVirtualEnvironmentVMEFIDiskFileFormat,
		mkResourceVirtualEnvironmentVMEFIDiskType,
	})

	test.AssertValueTypes(t, efiDiskSchema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentVMEFIDiskDatastoreID: schema.TypeString,
		mkResourceVirtualEnvironmentVMEFIDiskFileFormat:  schema.TypeString,
		mkResourceVirtualEnvironmentVMEFIDiskType:        schema.TypeString,
	})

	initializationSchema := test.AssertNestedSchemaExistence(
		t,
		s,
		mkResourceVirtualEnvironmentVMInitialization,
	)

	test.AssertOptionalArguments(t, initializationSchema, []string{
		mkResourceVirtualEnvironmentVMInitializationDatastoreID,
		mkResourceVirtualEnvironmentVMInitializationInterface,
		mkResourceVirtualEnvironmentVMInitializationDNS,
		mkResourceVirtualEnvironmentVMInitializationIPConfig,
		mkResourceVirtualEnvironmentVMInitializationUserAccount,
	})

	test.AssertValueTypes(t, initializationSchema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentVMInitializationDatastoreID: schema.TypeString,
		mkResourceVirtualEnvironmentVMInitializationInterface:   schema.TypeString,
		mkResourceVirtualEnvironmentVMInitializationDNS:         schema.TypeList,
		mkResourceVirtualEnvironmentVMInitializationIPConfig:    schema.TypeList,
		mkResourceVirtualEnvironmentVMInitializationUserAccount: schema.TypeList,
	})

	hostPCISchema := test.AssertNestedSchemaExistence(t, s, mkResourceVirtualEnvironmentVMHostPCI)

	test.AssertOptionalArguments(t, hostPCISchema, []string{
		mkResourceVirtualEnvironmentVMHostPCIDeviceMDev,
		mkResourceVirtualEnvironmentVMHostPCIDevicePCIE,
		mkResourceVirtualEnvironmentVMHostPCIDeviceROMBAR,
		mkResourceVirtualEnvironmentVMHostPCIDeviceROMFile,
		mkResourceVirtualEnvironmentVMHostPCIDeviceXVGA,
	})

	test.AssertValueTypes(t, hostPCISchema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentVMHostPCIDevice:        schema.TypeString,
		mkResourceVirtualEnvironmentVMHostPCIDeviceMDev:    schema.TypeString,
		mkResourceVirtualEnvironmentVMHostPCIDevicePCIE:    schema.TypeBool,
		mkResourceVirtualEnvironmentVMHostPCIDeviceROMBAR:  schema.TypeBool,
		mkResourceVirtualEnvironmentVMHostPCIDeviceROMFile: schema.TypeString,
		mkResourceVirtualEnvironmentVMHostPCIDeviceXVGA:    schema.TypeBool,
	})

	hostUSBSchema := test.AssertNestedSchemaExistence(t, s, mkResourceVirtualEnvironmentVMHostUSB)

	test.AssertOptionalArguments(t, hostUSBSchema, []string{
		mkResourceVirtualEnvironmentVMHostUSBDeviceMapping,
	})

	test.AssertValueTypes(t, hostUSBSchema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentVMHostUSBDevice:     schema.TypeString,
		mkResourceVirtualEnvironmentVMHostUSBDeviceUSB3: schema.TypeBool,
	})

	initializationDNSSchema := test.AssertNestedSchemaExistence(
		t,
		initializationSchema,
		mkResourceVirtualEnvironmentVMInitializationDNS,
	)

	test.AssertOptionalArguments(t, initializationDNSSchema, []string{
		mkResourceVirtualEnvironmentVMInitializationDNSDomain,
		mkResourceVirtualEnvironmentVMInitializationDNSServer,
	})

	test.AssertValueTypes(t, initializationDNSSchema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentVMInitializationDNSDomain: schema.TypeString,
		mkResourceVirtualEnvironmentVMInitializationDNSServer: schema.TypeString,
	})

	initializationIPConfigSchema := test.AssertNestedSchemaExistence(
		t,
		initializationSchema,
		mkResourceVirtualEnvironmentVMInitializationIPConfig,
	)

	test.AssertOptionalArguments(t, initializationIPConfigSchema, []string{
		mkResourceVirtualEnvironmentVMInitializationIPConfigIPv4,
		mkResourceVirtualEnvironmentVMInitializationIPConfigIPv6,
	})

	test.AssertValueTypes(t, initializationIPConfigSchema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentVMInitializationIPConfigIPv4: schema.TypeList,
		mkResourceVirtualEnvironmentVMInitializationIPConfigIPv6: schema.TypeList,
	})

	initializationIPConfigIPv4Schema := test.AssertNestedSchemaExistence(
		t,
		initializationIPConfigSchema,
		mkResourceVirtualEnvironmentVMInitializationIPConfigIPv4,
	)

	test.AssertOptionalArguments(t, initializationIPConfigIPv4Schema, []string{
		mkResourceVirtualEnvironmentVMInitializationIPConfigIPv4Address,
		mkResourceVirtualEnvironmentVMInitializationIPConfigIPv4Gateway,
	})

	test.AssertValueTypes(t, initializationIPConfigIPv4Schema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentVMInitializationIPConfigIPv4Address: schema.TypeString,
		mkResourceVirtualEnvironmentVMInitializationIPConfigIPv4Gateway: schema.TypeString,
	})

	initializationIPConfigIPv6Schema := test.AssertNestedSchemaExistence(
		t,
		initializationIPConfigSchema,
		mkResourceVirtualEnvironmentVMInitializationIPConfigIPv6,
	)

	test.AssertOptionalArguments(t, initializationIPConfigIPv6Schema, []string{
		mkResourceVirtualEnvironmentVMInitializationIPConfigIPv6Address,
		mkResourceVirtualEnvironmentVMInitializationIPConfigIPv6Gateway,
	})

	test.AssertValueTypes(t, initializationIPConfigIPv6Schema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentVMInitializationIPConfigIPv6Address: schema.TypeString,
		mkResourceVirtualEnvironmentVMInitializationIPConfigIPv6Gateway: schema.TypeString,
	})

	initializationUserAccountSchema := test.AssertNestedSchemaExistence(
		t,
		initializationSchema,
		mkResourceVirtualEnvironmentVMInitializationUserAccount,
	)

	test.AssertOptionalArguments(t, initializationUserAccountSchema, []string{
		mkResourceVirtualEnvironmentVMInitializationUserAccountKeys,
		mkResourceVirtualEnvironmentVMInitializationUserAccountPassword,
		mkResourceVirtualEnvironmentVMInitializationUserAccountUsername,
	})

	test.AssertValueTypes(t, initializationUserAccountSchema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentVMInitializationUserAccountKeys:     schema.TypeList,
		mkResourceVirtualEnvironmentVMInitializationUserAccountPassword: schema.TypeString,
		mkResourceVirtualEnvironmentVMInitializationUserAccountUsername: schema.TypeString,
	})

	memorySchema := test.AssertNestedSchemaExistence(t, s, mkResourceVirtualEnvironmentVMMemory)

	test.AssertOptionalArguments(t, memorySchema, []string{
		mkResourceVirtualEnvironmentVMMemoryDedicated,
		mkResourceVirtualEnvironmentVMMemoryFloating,
		mkResourceVirtualEnvironmentVMMemoryShared,
	})

	test.AssertValueTypes(t, memorySchema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentVMMemoryDedicated: schema.TypeInt,
		mkResourceVirtualEnvironmentVMMemoryFloating:  schema.TypeInt,
		mkResourceVirtualEnvironmentVMMemoryShared:    schema.TypeInt,
	})

	networkDeviceSchema := test.AssertNestedSchemaExistence(
		t,
		s,
		mkResourceVirtualEnvironmentVMNetworkDevice,
	)

	test.AssertOptionalArguments(t, networkDeviceSchema, []string{
		mkResourceVirtualEnvironmentVMNetworkDeviceBridge,
		mkResourceVirtualEnvironmentVMNetworkDeviceEnabled,
		mkResourceVirtualEnvironmentVMNetworkDeviceMACAddress,
		mkResourceVirtualEnvironmentVMNetworkDeviceModel,
		mkResourceVirtualEnvironmentVMNetworkDeviceRateLimit,
		mkResourceVirtualEnvironmentVMNetworkDeviceVLANID,
		mkResourceVirtualEnvironmentVMNetworkDeviceMTU,
	})

	test.AssertValueTypes(t, networkDeviceSchema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentVMNetworkDeviceBridge:     schema.TypeString,
		mkResourceVirtualEnvironmentVMNetworkDeviceEnabled:    schema.TypeBool,
		mkResourceVirtualEnvironmentVMNetworkDeviceMACAddress: schema.TypeString,
		mkResourceVirtualEnvironmentVMNetworkDeviceModel:      schema.TypeString,
		mkResourceVirtualEnvironmentVMNetworkDeviceRateLimit:  schema.TypeFloat,
		mkResourceVirtualEnvironmentVMNetworkDeviceVLANID:     schema.TypeInt,
		mkResourceVirtualEnvironmentVMNetworkDeviceMTU:        schema.TypeInt,
	})

	operatingSystemSchema := test.AssertNestedSchemaExistence(
		t,
		s,
		mkResourceVirtualEnvironmentVMOperatingSystem,
	)

	test.AssertOptionalArguments(t, operatingSystemSchema, []string{
		mkResourceVirtualEnvironmentVMOperatingSystemType,
	})

	test.AssertValueTypes(t, operatingSystemSchema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentVMOperatingSystemType: schema.TypeString,
	})

	serialDeviceSchema := test.AssertNestedSchemaExistence(
		t,
		s,
		mkResourceVirtualEnvironmentVMSerialDevice,
	)

	test.AssertOptionalArguments(t, serialDeviceSchema, []string{
		mkResourceVirtualEnvironmentVMSerialDeviceDevice,
	})

	test.AssertValueTypes(t, serialDeviceSchema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentVMSerialDeviceDevice: schema.TypeString,
	})

	vgaSchema := test.AssertNestedSchemaExistence(t, s, mkResourceVirtualEnvironmentVMVGA)

	test.AssertOptionalArguments(t, vgaSchema, []string{
		mkResourceVirtualEnvironmentVMVGAEnabled,
		mkResourceVirtualEnvironmentVMVGAMemory,
		mkResourceVirtualEnvironmentVMVGAType,
	})

	test.AssertValueTypes(t, vgaSchema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentVMVGAEnabled: schema.TypeBool,
		mkResourceVirtualEnvironmentVMVGAMemory:  schema.TypeInt,
		mkResourceVirtualEnvironmentVMVGAType:    schema.TypeString,
	})
}

func Test_parseImportIDWIthNodeName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		value            string
		valid            bool
		expectedNodeName string
		expectedID       string
	}{
		{"empty", "", false, "", ""},
		{"missing slash", "invalid", false, "", ""},
		{"valid", "host/id", true, "host", "id"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			nodeName, id, err := parseImportIDWithNodeName(tt.value)

			if !tt.valid {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.expectedNodeName, nodeName)
			require.Equal(t, tt.expectedID, id)
		})
	}
}
