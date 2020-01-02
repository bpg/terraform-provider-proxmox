/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmoxtf

import (
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
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
		mkResourceVirtualEnvironmentContainerOperatingSystem,
	})

	testOptionalArguments(t, s, []string{
		mkResourceVirtualEnvironmentContainerCPU,
		mkResourceVirtualEnvironmentContainerDescription,
		mkResourceVirtualEnvironmentContainerInitialization,
		mkResourceVirtualEnvironmentContainerMemory,
		mkResourceVirtualEnvironmentContainerPoolID,
		mkResourceVirtualEnvironmentContainerStarted,
		mkResourceVirtualEnvironmentContainerVMID,
	})

	testSchemaValueTypes(t, s, []string{
		mkResourceVirtualEnvironmentContainerCPU,
		mkResourceVirtualEnvironmentContainerDescription,
		mkResourceVirtualEnvironmentContainerInitialization,
		mkResourceVirtualEnvironmentContainerMemory,
		mkResourceVirtualEnvironmentContainerOperatingSystem,
		mkResourceVirtualEnvironmentContainerPoolID,
		mkResourceVirtualEnvironmentContainerStarted,
		mkResourceVirtualEnvironmentContainerVMID,
	}, []schema.ValueType{
		schema.TypeList,
		schema.TypeString,
		schema.TypeList,
		schema.TypeList,
		schema.TypeList,
		schema.TypeString,
		schema.TypeBool,
		schema.TypeInt,
	})

	cpuSchema := testNestedSchemaExistence(t, s, mkResourceVirtualEnvironmentContainerCPU)

	testOptionalArguments(t, cpuSchema, []string{
		mkResourceVirtualEnvironmentContainerCPUArchitecture,
		mkResourceVirtualEnvironmentContainerCPUCores,
		mkResourceVirtualEnvironmentContainerCPUUnits,
	})

	testSchemaValueTypes(t, cpuSchema, []string{
		mkResourceVirtualEnvironmentContainerCPUArchitecture,
		mkResourceVirtualEnvironmentContainerCPUCores,
		mkResourceVirtualEnvironmentContainerCPUUnits,
	}, []schema.ValueType{
		schema.TypeString,
		schema.TypeInt,
		schema.TypeInt,
	})

	initializationSchema := testNestedSchemaExistence(t, s, mkResourceVirtualEnvironmentContainerInitialization)

	testOptionalArguments(t, initializationSchema, []string{
		mkResourceVirtualEnvironmentContainerInitializationDNS,
		mkResourceVirtualEnvironmentContainerInitializationHostname,
		mkResourceVirtualEnvironmentContainerInitializationIPConfig,
		mkResourceVirtualEnvironmentContainerInitializationUserAccount,
	})

	testSchemaValueTypes(t, initializationSchema, []string{
		mkResourceVirtualEnvironmentContainerInitializationDNS,
		mkResourceVirtualEnvironmentContainerInitializationHostname,
		mkResourceVirtualEnvironmentContainerInitializationIPConfig,
		mkResourceVirtualEnvironmentContainerInitializationUserAccount,
	}, []schema.ValueType{
		schema.TypeList,
		schema.TypeString,
		schema.TypeList,
		schema.TypeList,
	})

	initializationDNSSchema := testNestedSchemaExistence(t, initializationSchema, mkResourceVirtualEnvironmentContainerInitializationDNS)

	testOptionalArguments(t, initializationDNSSchema, []string{
		mkResourceVirtualEnvironmentContainerInitializationDNSDomain,
		mkResourceVirtualEnvironmentContainerInitializationDNSServer,
	})

	testSchemaValueTypes(t, initializationDNSSchema, []string{
		mkResourceVirtualEnvironmentContainerInitializationDNSDomain,
		mkResourceVirtualEnvironmentContainerInitializationDNSServer,
	}, []schema.ValueType{
		schema.TypeString,
		schema.TypeString,
	})

	initializationIPConfigSchema := testNestedSchemaExistence(t, initializationSchema, mkResourceVirtualEnvironmentContainerInitializationIPConfig)

	testOptionalArguments(t, initializationIPConfigSchema, []string{
		mkResourceVirtualEnvironmentContainerInitializationIPConfigIPv4,
		mkResourceVirtualEnvironmentContainerInitializationIPConfigIPv6,
	})

	testSchemaValueTypes(t, initializationIPConfigSchema, []string{
		mkResourceVirtualEnvironmentContainerInitializationIPConfigIPv4,
		mkResourceVirtualEnvironmentContainerInitializationIPConfigIPv6,
	}, []schema.ValueType{
		schema.TypeList,
		schema.TypeList,
	})

	initializationIPConfigIPv4Schema := testNestedSchemaExistence(t, initializationIPConfigSchema, mkResourceVirtualEnvironmentContainerInitializationIPConfigIPv4)

	testOptionalArguments(t, initializationIPConfigIPv4Schema, []string{
		mkResourceVirtualEnvironmentContainerInitializationIPConfigIPv4Address,
		mkResourceVirtualEnvironmentContainerInitializationIPConfigIPv4Gateway,
	})

	testSchemaValueTypes(t, initializationIPConfigIPv4Schema, []string{
		mkResourceVirtualEnvironmentContainerInitializationIPConfigIPv4Address,
		mkResourceVirtualEnvironmentContainerInitializationIPConfigIPv4Gateway,
	}, []schema.ValueType{
		schema.TypeString,
		schema.TypeString,
	})

	initializationIPConfigIPv6Schema := testNestedSchemaExistence(t, initializationIPConfigSchema, mkResourceVirtualEnvironmentContainerInitializationIPConfigIPv6)

	testOptionalArguments(t, initializationIPConfigIPv6Schema, []string{
		mkResourceVirtualEnvironmentContainerInitializationIPConfigIPv6Address,
		mkResourceVirtualEnvironmentContainerInitializationIPConfigIPv6Gateway,
	})

	testSchemaValueTypes(t, initializationIPConfigIPv6Schema, []string{
		mkResourceVirtualEnvironmentContainerInitializationIPConfigIPv6Address,
		mkResourceVirtualEnvironmentContainerInitializationIPConfigIPv6Gateway,
	}, []schema.ValueType{
		schema.TypeString,
		schema.TypeString,
	})

	initializationUserAccountSchema := testNestedSchemaExistence(t, initializationSchema, mkResourceVirtualEnvironmentContainerInitializationUserAccount)

	testOptionalArguments(t, initializationUserAccountSchema, []string{
		mkResourceVirtualEnvironmentContainerInitializationUserAccountKeys,
		mkResourceVirtualEnvironmentContainerInitializationUserAccountPassword,
	})

	testSchemaValueTypes(t, initializationUserAccountSchema, []string{
		mkResourceVirtualEnvironmentContainerInitializationUserAccountKeys,
		mkResourceVirtualEnvironmentContainerInitializationUserAccountPassword,
	}, []schema.ValueType{
		schema.TypeList,
		schema.TypeString,
	})

	memorySchema := testNestedSchemaExistence(t, s, mkResourceVirtualEnvironmentContainerMemory)

	testOptionalArguments(t, memorySchema, []string{
		mkResourceVirtualEnvironmentContainerMemoryDedicated,
		mkResourceVirtualEnvironmentContainerMemorySwap,
	})

	testSchemaValueTypes(t, memorySchema, []string{
		mkResourceVirtualEnvironmentContainerMemoryDedicated,
		mkResourceVirtualEnvironmentContainerMemorySwap,
	}, []schema.ValueType{
		schema.TypeInt,
		schema.TypeInt,
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

	testSchemaValueTypes(t, networkInterfaceSchema, []string{
		mkResourceVirtualEnvironmentContainerNetworkInterfaceBridge,
		mkResourceVirtualEnvironmentContainerNetworkInterfaceEnabled,
		mkResourceVirtualEnvironmentContainerNetworkInterfaceMACAddress,
		mkResourceVirtualEnvironmentContainerNetworkInterfaceName,
		mkResourceVirtualEnvironmentContainerNetworkInterfaceRateLimit,
		mkResourceVirtualEnvironmentContainerNetworkInterfaceVLANID,
	}, []schema.ValueType{
		schema.TypeString,
		schema.TypeBool,
		schema.TypeString,
		schema.TypeString,
		schema.TypeFloat,
		schema.TypeInt,
	})

	operatingSystemSchema := testNestedSchemaExistence(t, s, mkResourceVirtualEnvironmentContainerOperatingSystem)

	testRequiredArguments(t, operatingSystemSchema, []string{
		mkResourceVirtualEnvironmentContainerOperatingSystemTemplateFileID,
	})

	testOptionalArguments(t, operatingSystemSchema, []string{
		mkResourceVirtualEnvironmentContainerOperatingSystemType,
	})

	testSchemaValueTypes(t, operatingSystemSchema, []string{
		mkResourceVirtualEnvironmentContainerOperatingSystemTemplateFileID,
		mkResourceVirtualEnvironmentContainerOperatingSystemType,
	}, []schema.ValueType{
		schema.TypeString,
		schema.TypeString,
	})
}
