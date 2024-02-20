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
		mkNodeName,
	})

	test.AssertOptionalArguments(t, s, []string{
		mkACPI,
		mkAgent,
		mkAudioDevice,
		mkBIOS,
		mkBootOrder,
		mkCDROM,
		mkClone,
		mkCPU,
		mkDescription,
		mkDisk,
		mkEFIDisk,
		mkInitialization,
		mkHostPCI,
		mkHostUSB,
		mkKeyboardLayout,
		mkKVMArguments,
		mkMachine,
		mkMemory,
		mkName,
		mkNetworkDevice,
		mkOperatingSystem,
		mkPoolID,
		mkSerialDevice,
		mkStarted,
		mkTabletDevice,
		mkTemplate,
		mkVMID,
		mkSCSIHardware,
	})

	test.AssertComputedAttributes(t, s, []string{
		mkIPv4Addresses,
		mkIPv6Addresses,
		mkMACAddresses,
		mkNetworkInterfaceNames,
	})

	test.AssertValueTypes(t, s, map[string]schema.ValueType{
		mkACPI:                  schema.TypeBool,
		mkAgent:                 schema.TypeList,
		mkAudioDevice:           schema.TypeList,
		mkBIOS:                  schema.TypeString,
		mkBootOrder:             schema.TypeList,
		mkCDROM:                 schema.TypeList,
		mkCPU:                   schema.TypeList,
		mkDescription:           schema.TypeString,
		mkDisk:                  schema.TypeList,
		mkEFIDisk:               schema.TypeList,
		mkHostPCI:               schema.TypeList,
		mkHostUSB:               schema.TypeList,
		mkInitialization:        schema.TypeList,
		mkIPv4Addresses:         schema.TypeList,
		mkIPv6Addresses:         schema.TypeList,
		mkKeyboardLayout:        schema.TypeString,
		mkKVMArguments:          schema.TypeString,
		mkMachine:               schema.TypeString,
		mkMemory:                schema.TypeList,
		mkName:                  schema.TypeString,
		mkNetworkDevice:         schema.TypeList,
		mkMACAddresses:          schema.TypeList,
		mkNetworkInterfaceNames: schema.TypeList,
		mkOperatingSystem:       schema.TypeList,
		mkPoolID:                schema.TypeString,
		mkSerialDevice:          schema.TypeList,
		mkStarted:               schema.TypeBool,
		mkTabletDevice:          schema.TypeBool,
		mkTemplate:              schema.TypeBool,
		mkVMID:                  schema.TypeInt,
		mkSCSIHardware:          schema.TypeString,
	})

	agentSchema := test.AssertNestedSchemaExistence(t, s, mkAgent)

	test.AssertOptionalArguments(t, agentSchema, []string{
		mkAgentEnabled,
		mkAgentTimeout,
		mkAgentTrim,
		mkAgentType,
	})

	test.AssertValueTypes(t, agentSchema, map[string]schema.ValueType{
		mkAgentEnabled: schema.TypeBool,
		mkAgentTrim:    schema.TypeBool,
		mkAgentType:    schema.TypeString,
	})

	audioDeviceSchema := test.AssertNestedSchemaExistence(t, s, mkAudioDevice)

	test.AssertOptionalArguments(t, audioDeviceSchema, []string{
		mkAudioDeviceDevice,
		mkAudioDeviceDriver,
	})

	test.AssertValueTypes(t, audioDeviceSchema, map[string]schema.ValueType{
		mkAudioDeviceDevice: schema.TypeString,
		mkAudioDeviceDriver: schema.TypeString,
	})

	cdromSchema := test.AssertNestedSchemaExistence(t, s, mkCDROM)

	test.AssertOptionalArguments(t, cdromSchema, []string{
		mkCDROMEnabled,
		mkCDROMFileID,
	})

	test.AssertValueTypes(t, cdromSchema, map[string]schema.ValueType{
		mkCDROMEnabled: schema.TypeBool,
		mkCDROMFileID:  schema.TypeString,
	})

	cloneSchema := test.AssertNestedSchemaExistence(t, s, mkClone)

	test.AssertRequiredArguments(t, cloneSchema, []string{
		mkCloneVMID,
	})

	test.AssertOptionalArguments(t, cloneSchema, []string{
		mkCloneDatastoreID,
		mkCloneNodeName,
	})

	test.AssertValueTypes(t, cloneSchema, map[string]schema.ValueType{
		mkCloneDatastoreID: schema.TypeString,
		mkCloneNodeName:    schema.TypeString,
		mkCloneVMID:        schema.TypeInt,
	})

	cpuSchema := test.AssertNestedSchemaExistence(t, s, mkCPU)

	test.AssertOptionalArguments(t, cpuSchema, []string{
		mkCPUArchitecture,
		mkCPUCores,
		mkCPUFlags,
		mkCPUHotplugged,
		mkCPUNUMA,
		mkCPUSockets,
		mkCPUType,
		mkCPUUnits,
	})

	test.AssertValueTypes(t, cpuSchema, map[string]schema.ValueType{
		mkCPUArchitecture: schema.TypeString,
		mkCPUCores:        schema.TypeInt,
		mkCPUFlags:        schema.TypeList,
		mkCPUHotplugged:   schema.TypeInt,
		mkCPUNUMA:         schema.TypeBool,
		mkCPUSockets:      schema.TypeInt,
		mkCPUType:         schema.TypeString,
		mkCPUUnits:        schema.TypeInt,
	})

	diskSchema := test.AssertNestedSchemaExistence(t, s, mkDisk)

	test.AssertOptionalArguments(t, diskSchema, []string{
		mkDiskDatastoreID,
		mkDiskPathInDatastore,
		mkDiskFileFormat,
		mkDiskFileID,
		mkDiskSize,
	})

	test.AssertValueTypes(t, diskSchema, map[string]schema.ValueType{
		mkDiskDatastoreID:     schema.TypeString,
		mkDiskPathInDatastore: schema.TypeString,
		mkDiskFileFormat:      schema.TypeString,
		mkDiskFileID:          schema.TypeString,
		mkDiskSize:            schema.TypeInt,
	})

	diskSpeedSchema := test.AssertNestedSchemaExistence(
		t,
		diskSchema,
		mkDiskSpeed,
	)

	test.AssertOptionalArguments(t, diskSpeedSchema, []string{
		mkDiskSpeedRead,
		mkDiskSpeedReadBurstable,
		mkDiskSpeedWrite,
		mkDiskSpeedWriteBurstable,
	})

	test.AssertValueTypes(t, diskSpeedSchema, map[string]schema.ValueType{
		mkDiskSpeedRead:           schema.TypeInt,
		mkDiskSpeedReadBurstable:  schema.TypeInt,
		mkDiskSpeedWrite:          schema.TypeInt,
		mkDiskSpeedWriteBurstable: schema.TypeInt,
	})

	efiDiskSchema := test.AssertNestedSchemaExistence(t, s, mkEFIDisk)

	test.AssertOptionalArguments(t, efiDiskSchema, []string{
		mkEFIDiskDatastoreID,
		mkEFIDiskFileFormat,
		mkEFIDiskType,
	})

	test.AssertValueTypes(t, efiDiskSchema, map[string]schema.ValueType{
		mkEFIDiskDatastoreID: schema.TypeString,
		mkEFIDiskFileFormat:  schema.TypeString,
		mkEFIDiskType:        schema.TypeString,
	})

	initializationSchema := test.AssertNestedSchemaExistence(
		t,
		s,
		mkInitialization,
	)

	test.AssertOptionalArguments(t, initializationSchema, []string{
		mkInitializationDatastoreID,
		mkInitializationInterface,
		mkInitializationDNS,
		mkInitializationIPConfig,
		mkInitializationUserAccount,
	})

	test.AssertValueTypes(t, initializationSchema, map[string]schema.ValueType{
		mkInitializationDatastoreID: schema.TypeString,
		mkInitializationInterface:   schema.TypeString,
		mkInitializationDNS:         schema.TypeList,
		mkInitializationIPConfig:    schema.TypeList,
		mkInitializationUserAccount: schema.TypeList,
	})

	hostPCISchema := test.AssertNestedSchemaExistence(t, s, mkHostPCI)

	test.AssertOptionalArguments(t, hostPCISchema, []string{
		mkHostPCIDeviceMDev,
		mkHostPCIDevicePCIE,
		mkHostPCIDeviceROMBAR,
		mkHostPCIDeviceROMFile,
		mkHostPCIDeviceXVGA,
	})

	test.AssertValueTypes(t, hostPCISchema, map[string]schema.ValueType{
		mkHostPCIDevice:        schema.TypeString,
		mkHostPCIDeviceMDev:    schema.TypeString,
		mkHostPCIDevicePCIE:    schema.TypeBool,
		mkHostPCIDeviceROMBAR:  schema.TypeBool,
		mkHostPCIDeviceROMFile: schema.TypeString,
		mkHostPCIDeviceXVGA:    schema.TypeBool,
	})

	hostUSBSchema := test.AssertNestedSchemaExistence(t, s, mkHostUSB)

	test.AssertOptionalArguments(t, hostUSBSchema, []string{
		mkHostUSBDeviceMapping,
	})

	test.AssertValueTypes(t, hostUSBSchema, map[string]schema.ValueType{
		mkHostUSBDevice:     schema.TypeString,
		mkHostUSBDeviceUSB3: schema.TypeBool,
	})

	initializationDNSSchema := test.AssertNestedSchemaExistence(
		t,
		initializationSchema,
		mkInitializationDNS,
	)

	test.AssertOptionalArguments(t, initializationDNSSchema, []string{
		mkInitializationDNSDomain,
		mkInitializationDNSServer,
		mkInitializationDNSServers,
	})

	test.AssertValueTypes(t, initializationDNSSchema, map[string]schema.ValueType{
		mkInitializationDNSDomain:  schema.TypeString,
		mkInitializationDNSServer:  schema.TypeString,
		mkInitializationDNSServers: schema.TypeList,
	})

	initializationIPConfigSchema := test.AssertNestedSchemaExistence(
		t,
		initializationSchema,
		mkInitializationIPConfig,
	)

	test.AssertOptionalArguments(t, initializationIPConfigSchema, []string{
		mkInitializationIPConfigIPv4,
		mkInitializationIPConfigIPv6,
	})

	test.AssertValueTypes(t, initializationIPConfigSchema, map[string]schema.ValueType{
		mkInitializationIPConfigIPv4: schema.TypeList,
		mkInitializationIPConfigIPv6: schema.TypeList,
	})

	initializationIPConfigIPv4Schema := test.AssertNestedSchemaExistence(
		t,
		initializationIPConfigSchema,
		mkInitializationIPConfigIPv4,
	)

	test.AssertOptionalArguments(t, initializationIPConfigIPv4Schema, []string{
		mkInitializationIPConfigIPv4Address,
		mkInitializationIPConfigIPv4Gateway,
	})

	test.AssertValueTypes(t, initializationIPConfigIPv4Schema, map[string]schema.ValueType{
		mkInitializationIPConfigIPv4Address: schema.TypeString,
		mkInitializationIPConfigIPv4Gateway: schema.TypeString,
	})

	initializationIPConfigIPv6Schema := test.AssertNestedSchemaExistence(
		t,
		initializationIPConfigSchema,
		mkInitializationIPConfigIPv6,
	)

	test.AssertOptionalArguments(t, initializationIPConfigIPv6Schema, []string{
		mkInitializationIPConfigIPv6Address,
		mkInitializationIPConfigIPv6Gateway,
	})

	test.AssertValueTypes(t, initializationIPConfigIPv6Schema, map[string]schema.ValueType{
		mkInitializationIPConfigIPv6Address: schema.TypeString,
		mkInitializationIPConfigIPv6Gateway: schema.TypeString,
	})

	initializationUserAccountSchema := test.AssertNestedSchemaExistence(
		t,
		initializationSchema,
		mkInitializationUserAccount,
	)

	test.AssertOptionalArguments(t, initializationUserAccountSchema, []string{
		mkInitializationUserAccountKeys,
		mkInitializationUserAccountPassword,
		mkInitializationUserAccountUsername,
	})

	test.AssertValueTypes(t, initializationUserAccountSchema, map[string]schema.ValueType{
		mkInitializationUserAccountKeys:     schema.TypeList,
		mkInitializationUserAccountPassword: schema.TypeString,
		mkInitializationUserAccountUsername: schema.TypeString,
	})

	memorySchema := test.AssertNestedSchemaExistence(t, s, mkMemory)

	test.AssertOptionalArguments(t, memorySchema, []string{
		mkMemoryDedicated,
		mkMemoryFloating,
		mkMemoryShared,
	})

	test.AssertValueTypes(t, memorySchema, map[string]schema.ValueType{
		mkMemoryDedicated: schema.TypeInt,
		mkMemoryFloating:  schema.TypeInt,
		mkMemoryShared:    schema.TypeInt,
	})

	networkDeviceSchema := test.AssertNestedSchemaExistence(
		t,
		s,
		mkNetworkDevice,
	)

	test.AssertOptionalArguments(t, networkDeviceSchema, []string{
		mkNetworkDeviceBridge,
		mkNetworkDeviceEnabled,
		mkNetworkDeviceMACAddress,
		mkNetworkDeviceModel,
		mkNetworkDeviceRateLimit,
		mkNetworkDeviceVLANID,
		mkNetworkDeviceMTU,
	})

	test.AssertValueTypes(t, networkDeviceSchema, map[string]schema.ValueType{
		mkNetworkDeviceBridge:     schema.TypeString,
		mkNetworkDeviceEnabled:    schema.TypeBool,
		mkNetworkDeviceMACAddress: schema.TypeString,
		mkNetworkDeviceModel:      schema.TypeString,
		mkNetworkDeviceRateLimit:  schema.TypeFloat,
		mkNetworkDeviceVLANID:     schema.TypeInt,
		mkNetworkDeviceMTU:        schema.TypeInt,
	})

	operatingSystemSchema := test.AssertNestedSchemaExistence(
		t,
		s,
		mkOperatingSystem,
	)

	test.AssertOptionalArguments(t, operatingSystemSchema, []string{
		mkOperatingSystemType,
	})

	test.AssertValueTypes(t, operatingSystemSchema, map[string]schema.ValueType{
		mkOperatingSystemType: schema.TypeString,
	})

	serialDeviceSchema := test.AssertNestedSchemaExistence(
		t,
		s,
		mkSerialDevice,
	)

	test.AssertOptionalArguments(t, serialDeviceSchema, []string{
		mkSerialDeviceDevice,
	})

	test.AssertValueTypes(t, serialDeviceSchema, map[string]schema.ValueType{
		mkSerialDeviceDevice: schema.TypeString,
	})

	vgaSchema := test.AssertNestedSchemaExistence(t, s, mkVGA)

	test.AssertOptionalArguments(t, vgaSchema, []string{
		mkVGAEnabled,
		mkVGAMemory,
		mkVGAType,
	})

	test.AssertValueTypes(t, vgaSchema, map[string]schema.ValueType{
		mkVGAEnabled: schema.TypeBool,
		mkVGAMemory:  schema.TypeInt,
		mkVGAType:    schema.TypeString,
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
