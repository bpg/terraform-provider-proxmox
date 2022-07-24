/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmoxtf

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// TestResourceVirtualEnvironmentContainerInstantiation tests whether the ResourceVirtualEnvironmentContainer instance can be instantiated.
func TestResourceVirtualEnvironmentContainerInstantiation(t *testing.T) {
	s := resourceVirtualEnvironmentContainer()

	if s == nil {
		t.Fatalf("Cannot instantiate resourceVirtualEnvironmentContainer")
	}
}

// TestResourceVirtualEnvironmentContainerSchema tests the resourceVirtualEnvironmentContainer schema.
func TestResourceVirtualEnvironmentContainerSchema(t *testing.T) {
	s := resourceVirtualEnvironmentContainer()

	testRequiredArguments(t, s, []string{
		mkResourceVirtualEnvironmentContainerNodeName,
	})

	testOptionalArguments(t, s, []string{
		mkResourceVirtualEnvironmentContainerCPU,
		mkResourceVirtualEnvironmentContainerDescription,
		mkResourceVirtualEnvironmentContainerDisk,
		mkResourceVirtualEnvironmentContainerInitialization,
		mkResourceVirtualEnvironmentContainerMemory,
		mkResourceVirtualEnvironmentContainerOperatingSystem,
		mkResourceVirtualEnvironmentContainerPoolID,
		mkResourceVirtualEnvironmentContainerStarted,
		mkResourceVirtualEnvironmentContainerTemplate,
		mkResourceVirtualEnvironmentContainerVMID,
	})

	testValueTypes(t, s, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentContainerCPU:             schema.TypeList,
		mkResourceVirtualEnvironmentContainerDescription:     schema.TypeString,
		mkResourceVirtualEnvironmentContainerDisk:            schema.TypeList,
		mkResourceVirtualEnvironmentContainerInitialization:  schema.TypeList,
		mkResourceVirtualEnvironmentContainerMemory:          schema.TypeList,
		mkResourceVirtualEnvironmentContainerOperatingSystem: schema.TypeList,
		mkResourceVirtualEnvironmentContainerPoolID:          schema.TypeString,
		mkResourceVirtualEnvironmentContainerStarted:         schema.TypeBool,
		mkResourceVirtualEnvironmentContainerTemplate:        schema.TypeBool,
		mkResourceVirtualEnvironmentContainerVMID:            schema.TypeInt,
	})

	cloneSchema := testNestedSchemaExistence(t, s, mkResourceVirtualEnvironmentContainerClone)

	testRequiredArguments(t, cloneSchema, []string{
		mkResourceVirtualEnvironmentContainerCloneVMID,
	})

	testOptionalArguments(t, cloneSchema, []string{
		mkResourceVirtualEnvironmentContainerCloneDatastoreID,
		mkResourceVirtualEnvironmentContainerCloneNodeName,
	})

	testValueTypes(t, cloneSchema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentContainerCloneDatastoreID: schema.TypeString,
		mkResourceVirtualEnvironmentContainerCloneNodeName:    schema.TypeString,
		mkResourceVirtualEnvironmentContainerCloneVMID:        schema.TypeInt,
	})

	cpuSchema := testNestedSchemaExistence(t, s, mkResourceVirtualEnvironmentContainerCPU)

	testOptionalArguments(t, cpuSchema, []string{
		mkResourceVirtualEnvironmentContainerCPUArchitecture,
		mkResourceVirtualEnvironmentContainerCPUCores,
		mkResourceVirtualEnvironmentContainerCPUUnits,
	})

	testValueTypes(t, cpuSchema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentContainerCPUArchitecture: schema.TypeString,
		mkResourceVirtualEnvironmentContainerCPUCores:        schema.TypeInt,
		mkResourceVirtualEnvironmentContainerCPUUnits:        schema.TypeInt,
	})

	diskSchema := testNestedSchemaExistence(t, s, mkResourceVirtualEnvironmentContainerDisk)

	testOptionalArguments(t, diskSchema, []string{
		mkResourceVirtualEnvironmentContainerDiskDatastoreID,
		mkResourceVirtualEnvironmentContainerCustomRootfsSize,
	})

	testValueTypes(t, diskSchema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentContainerDiskDatastoreID:  schema.TypeString,
		mkResourceVirtualEnvironmentContainerCustomRootfsSize: schema.TypeInt,
	})

	initializationSchema := testNestedSchemaExistence(t, s, mkResourceVirtualEnvironmentContainerInitialization)

	testOptionalArguments(t, initializationSchema, []string{
		mkResourceVirtualEnvironmentContainerInitializationDNS,
		mkResourceVirtualEnvironmentContainerInitializationHostname,
		mkResourceVirtualEnvironmentContainerInitializationIPConfig,
		mkResourceVirtualEnvironmentContainerInitializationUserAccount,
	})

	testValueTypes(t, initializationSchema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentContainerInitializationDNS:         schema.TypeList,
		mkResourceVirtualEnvironmentContainerInitializationHostname:    schema.TypeString,
		mkResourceVirtualEnvironmentContainerInitializationIPConfig:    schema.TypeList,
		mkResourceVirtualEnvironmentContainerInitializationUserAccount: schema.TypeList,
	})

	initializationDNSSchema := testNestedSchemaExistence(t, initializationSchema, mkResourceVirtualEnvironmentContainerInitializationDNS)

	testOptionalArguments(t, initializationDNSSchema, []string{
		mkResourceVirtualEnvironmentContainerInitializationDNSDomain,
		mkResourceVirtualEnvironmentContainerInitializationDNSServer,
	})

	testValueTypes(t, initializationDNSSchema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentContainerInitializationDNSDomain: schema.TypeString,
		mkResourceVirtualEnvironmentContainerInitializationDNSServer: schema.TypeString,
	})

	initializationIPConfigSchema := testNestedSchemaExistence(t, initializationSchema, mkResourceVirtualEnvironmentContainerInitializationIPConfig)

	testOptionalArguments(t, initializationIPConfigSchema, []string{
		mkResourceVirtualEnvironmentContainerInitializationIPConfigIPv4,
		mkResourceVirtualEnvironmentContainerInitializationIPConfigIPv6,
	})

	testValueTypes(t, initializationIPConfigSchema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentContainerInitializationIPConfigIPv4: schema.TypeList,
		mkResourceVirtualEnvironmentContainerInitializationIPConfigIPv6: schema.TypeList,
	})

	initializationIPConfigIPv4Schema := testNestedSchemaExistence(t, initializationIPConfigSchema, mkResourceVirtualEnvironmentContainerInitializationIPConfigIPv4)

	testOptionalArguments(t, initializationIPConfigIPv4Schema, []string{
		mkResourceVirtualEnvironmentContainerInitializationIPConfigIPv4Address,
		mkResourceVirtualEnvironmentContainerInitializationIPConfigIPv4Gateway,
	})

	testValueTypes(t, initializationIPConfigIPv4Schema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentContainerInitializationIPConfigIPv4Address: schema.TypeString,
		mkResourceVirtualEnvironmentContainerInitializationIPConfigIPv4Gateway: schema.TypeString,
	})

	initializationIPConfigIPv6Schema := testNestedSchemaExistence(t, initializationIPConfigSchema, mkResourceVirtualEnvironmentContainerInitializationIPConfigIPv6)

	testOptionalArguments(t, initializationIPConfigIPv6Schema, []string{
		mkResourceVirtualEnvironmentContainerInitializationIPConfigIPv6Address,
		mkResourceVirtualEnvironmentContainerInitializationIPConfigIPv6Gateway,
	})

	testValueTypes(t, initializationIPConfigIPv6Schema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentContainerInitializationIPConfigIPv6Address: schema.TypeString,
		mkResourceVirtualEnvironmentContainerInitializationIPConfigIPv6Gateway: schema.TypeString,
	})

	initializationUserAccountSchema := testNestedSchemaExistence(t, initializationSchema, mkResourceVirtualEnvironmentContainerInitializationUserAccount)

	testOptionalArguments(t, initializationUserAccountSchema, []string{
		mkResourceVirtualEnvironmentContainerInitializationUserAccountKeys,
		mkResourceVirtualEnvironmentContainerInitializationUserAccountPassword,
	})

	testValueTypes(t, initializationUserAccountSchema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentContainerInitializationUserAccountKeys:     schema.TypeList,
		mkResourceVirtualEnvironmentContainerInitializationUserAccountPassword: schema.TypeString,
	})

	memorySchema := testNestedSchemaExistence(t, s, mkResourceVirtualEnvironmentContainerMemory)

	testOptionalArguments(t, memorySchema, []string{
		mkResourceVirtualEnvironmentContainerMemoryDedicated,
		mkResourceVirtualEnvironmentContainerMemorySwap,
	})

	testValueTypes(t, memorySchema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentContainerMemoryDedicated: schema.TypeInt,
		mkResourceVirtualEnvironmentContainerMemorySwap:      schema.TypeInt,
	})

	networkInterfaceSchema := testNestedSchemaExistence(t, s, mkResourceVirtualEnvironmentContainerNetworkInterface)

	testRequiredArguments(t, networkInterfaceSchema, []string{
		mkResourceVirtualEnvironmentContainerNetworkInterfaceName,
	})

	testOptionalArguments(t, networkInterfaceSchema, []string{
		mkResourceVirtualEnvironmentContainerNetworkInterfaceBridge,
		mkResourceVirtualEnvironmentContainerNetworkInterfaceEnabled,
		mkResourceVirtualEnvironmentContainerNetworkInterfaceMACAddress,
		mkResourceVirtualEnvironmentContainerNetworkInterfaceRateLimit,
		mkResourceVirtualEnvironmentContainerNetworkInterfaceVLANID,
	})

	testValueTypes(t, networkInterfaceSchema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentContainerNetworkInterfaceBridge:     schema.TypeString,
		mkResourceVirtualEnvironmentContainerNetworkInterfaceEnabled:    schema.TypeBool,
		mkResourceVirtualEnvironmentContainerNetworkInterfaceMACAddress: schema.TypeString,
		mkResourceVirtualEnvironmentContainerNetworkInterfaceName:       schema.TypeString,
		mkResourceVirtualEnvironmentContainerNetworkInterfaceRateLimit:  schema.TypeFloat,
		mkResourceVirtualEnvironmentContainerNetworkInterfaceVLANID:     schema.TypeInt,
	})

	operatingSystemSchema := testNestedSchemaExistence(t, s, mkResourceVirtualEnvironmentContainerOperatingSystem)

	testRequiredArguments(t, operatingSystemSchema, []string{
		mkResourceVirtualEnvironmentContainerOperatingSystemTemplateFileID,
	})

	testOptionalArguments(t, operatingSystemSchema, []string{
		mkResourceVirtualEnvironmentContainerOperatingSystemType,
	})

	testValueTypes(t, operatingSystemSchema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentContainerOperatingSystemTemplateFileID: schema.TypeString,
		mkResourceVirtualEnvironmentContainerOperatingSystemType:           schema.TypeString,
	})
}
