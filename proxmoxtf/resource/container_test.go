/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package resource

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/test"
)

// TestContainerInstantiation tests whether the Container instance can be instantiated.
func TestContainerInstantiation(t *testing.T) {
	t.Parallel()
	s := Container()

	if s == nil {
		t.Fatalf("Cannot instantiate Container")
	}
}

// TestContainerSchema tests the Container schema.
func TestContainerSchema(t *testing.T) {
	t.Parallel()
	s := Container()

	test.AssertRequiredArguments(t, s, []string{
		mkResourceVirtualEnvironmentContainerNodeName,
	})

	test.AssertOptionalArguments(t, s, []string{
		mkResourceVirtualEnvironmentContainerCPU,
		mkResourceVirtualEnvironmentContainerDescription,
		mkResourceVirtualEnvironmentContainerDisk,
		mkResourceVirtualEnvironmentContainerInitialization,
		mkResourceVirtualEnvironmentContainerMemory,
		mkResourceVirtualEnvironmentContainerOperatingSystem,
		mkResourceVirtualEnvironmentContainerPoolID,
		mkResourceVirtualEnvironmentContainerStarted,
		mkResourceVirtualEnvironmentContainerTags,
		mkResourceVirtualEnvironmentContainerTemplate,
		mkResourceVirtualEnvironmentContainerUnprivileged,
		mkResourceVirtualEnvironmentContainerFeatures,
		mkResourceVirtualEnvironmentContainerVMID,
	})

	test.AssertValueTypes(t, s, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentContainerCPU:             schema.TypeList,
		mkResourceVirtualEnvironmentContainerDescription:     schema.TypeString,
		mkResourceVirtualEnvironmentContainerDisk:            schema.TypeList,
		mkResourceVirtualEnvironmentContainerInitialization:  schema.TypeList,
		mkResourceVirtualEnvironmentContainerMemory:          schema.TypeList,
		mkResourceVirtualEnvironmentContainerOperatingSystem: schema.TypeList,
		mkResourceVirtualEnvironmentContainerPoolID:          schema.TypeString,
		mkResourceVirtualEnvironmentContainerStarted:         schema.TypeBool,
		mkResourceVirtualEnvironmentContainerTags:            schema.TypeList,
		mkResourceVirtualEnvironmentContainerTemplate:        schema.TypeBool,
		mkResourceVirtualEnvironmentContainerUnprivileged:    schema.TypeBool,
		mkResourceVirtualEnvironmentContainerFeatures:        schema.TypeList,
		mkResourceVirtualEnvironmentContainerVMID:            schema.TypeInt,
	})

	cloneSchema := test.AssertNestedSchemaExistence(t, s, mkResourceVirtualEnvironmentContainerClone)

	test.AssertRequiredArguments(t, cloneSchema, []string{
		mkResourceVirtualEnvironmentContainerCloneVMID,
	})

	test.AssertOptionalArguments(t, cloneSchema, []string{
		mkResourceVirtualEnvironmentContainerCloneDatastoreID,
		mkResourceVirtualEnvironmentContainerCloneNodeName,
	})

	test.AssertValueTypes(t, cloneSchema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentContainerCloneDatastoreID: schema.TypeString,
		mkResourceVirtualEnvironmentContainerCloneNodeName:    schema.TypeString,
		mkResourceVirtualEnvironmentContainerCloneVMID:        schema.TypeInt,
	})

	cpuSchema := test.AssertNestedSchemaExistence(t, s, mkResourceVirtualEnvironmentContainerCPU)

	test.AssertOptionalArguments(t, cpuSchema, []string{
		mkResourceVirtualEnvironmentContainerCPUArchitecture,
		mkResourceVirtualEnvironmentContainerCPUCores,
		mkResourceVirtualEnvironmentContainerCPUUnits,
	})

	test.AssertValueTypes(t, cpuSchema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentContainerCPUArchitecture: schema.TypeString,
		mkResourceVirtualEnvironmentContainerCPUCores:        schema.TypeInt,
		mkResourceVirtualEnvironmentContainerCPUUnits:        schema.TypeInt,
	})

	diskSchema := test.AssertNestedSchemaExistence(t, s, mkResourceVirtualEnvironmentContainerDisk)

	test.AssertOptionalArguments(t, diskSchema, []string{
		mkResourceVirtualEnvironmentContainerDiskDatastoreID,
	})

	test.AssertValueTypes(t, diskSchema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentContainerDiskDatastoreID: schema.TypeString,
	})

	featuresSchema := testNestedSchemaExistence(t, s, mkResourceVirtualEnvironmentContainerFeatures)

	testOptionalArguments(t, featuresSchema, []string{
		mkResourceVirtualEnvironmentContainerFeaturesNesting,
	})

	testValueTypes(t, featuresSchema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentContainerFeaturesNesting: schema.TypeBool,
	})

	initializationSchema := test.AssertNestedSchemaExistence(
		t,
		s,
		mkResourceVirtualEnvironmentContainerInitialization,
	)

	test.AssertOptionalArguments(t, initializationSchema, []string{
		mkResourceVirtualEnvironmentContainerInitializationDNS,
		mkResourceVirtualEnvironmentContainerInitializationHostname,
		mkResourceVirtualEnvironmentContainerInitializationIPConfig,
		mkResourceVirtualEnvironmentContainerInitializationUserAccount,
	})

	test.AssertValueTypes(t, initializationSchema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentContainerInitializationDNS:         schema.TypeList,
		mkResourceVirtualEnvironmentContainerInitializationHostname:    schema.TypeString,
		mkResourceVirtualEnvironmentContainerInitializationIPConfig:    schema.TypeList,
		mkResourceVirtualEnvironmentContainerInitializationUserAccount: schema.TypeList,
	})

	initializationDNSSchema := test.AssertNestedSchemaExistence(
		t,
		initializationSchema,
		mkResourceVirtualEnvironmentContainerInitializationDNS,
	)

	test.AssertOptionalArguments(t, initializationDNSSchema, []string{
		mkResourceVirtualEnvironmentContainerInitializationDNSDomain,
		mkResourceVirtualEnvironmentContainerInitializationDNSServer,
	})

	test.AssertValueTypes(t, initializationDNSSchema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentContainerInitializationDNSDomain: schema.TypeString,
		mkResourceVirtualEnvironmentContainerInitializationDNSServer: schema.TypeString,
	})

	initializationIPConfigSchema := test.AssertNestedSchemaExistence(
		t,
		initializationSchema,
		mkResourceVirtualEnvironmentContainerInitializationIPConfig,
	)

	test.AssertOptionalArguments(t, initializationIPConfigSchema, []string{
		mkResourceVirtualEnvironmentContainerInitializationIPConfigIPv4,
		mkResourceVirtualEnvironmentContainerInitializationIPConfigIPv6,
	})

	test.AssertValueTypes(t, initializationIPConfigSchema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentContainerInitializationIPConfigIPv4: schema.TypeList,
		mkResourceVirtualEnvironmentContainerInitializationIPConfigIPv6: schema.TypeList,
	})

	initializationIPConfigIPv4Schema := test.AssertNestedSchemaExistence(
		t,
		initializationIPConfigSchema,
		mkResourceVirtualEnvironmentContainerInitializationIPConfigIPv4,
	)

	test.AssertOptionalArguments(t, initializationIPConfigIPv4Schema, []string{
		mkResourceVirtualEnvironmentContainerInitializationIPConfigIPv4Address,
		mkResourceVirtualEnvironmentContainerInitializationIPConfigIPv4Gateway,
	})

	test.AssertValueTypes(t, initializationIPConfigIPv4Schema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentContainerInitializationIPConfigIPv4Address: schema.TypeString,
		mkResourceVirtualEnvironmentContainerInitializationIPConfigIPv4Gateway: schema.TypeString,
	})

	initializationIPConfigIPv6Schema := test.AssertNestedSchemaExistence(
		t,
		initializationIPConfigSchema,
		mkResourceVirtualEnvironmentContainerInitializationIPConfigIPv6,
	)

	test.AssertOptionalArguments(t, initializationIPConfigIPv6Schema, []string{
		mkResourceVirtualEnvironmentContainerInitializationIPConfigIPv6Address,
		mkResourceVirtualEnvironmentContainerInitializationIPConfigIPv6Gateway,
	})

	test.AssertValueTypes(t, initializationIPConfigIPv6Schema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentContainerInitializationIPConfigIPv6Address: schema.TypeString,
		mkResourceVirtualEnvironmentContainerInitializationIPConfigIPv6Gateway: schema.TypeString,
	})

	initializationUserAccountSchema := test.AssertNestedSchemaExistence(
		t,
		initializationSchema,
		mkResourceVirtualEnvironmentContainerInitializationUserAccount,
	)

	test.AssertOptionalArguments(t, initializationUserAccountSchema, []string{
		mkResourceVirtualEnvironmentContainerInitializationUserAccountKeys,
		mkResourceVirtualEnvironmentContainerInitializationUserAccountPassword,
	})

	test.AssertValueTypes(t, initializationUserAccountSchema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentContainerInitializationUserAccountKeys:     schema.TypeList,
		mkResourceVirtualEnvironmentContainerInitializationUserAccountPassword: schema.TypeString,
	})

	memorySchema := test.AssertNestedSchemaExistence(t, s, mkResourceVirtualEnvironmentContainerMemory)

	test.AssertOptionalArguments(t, memorySchema, []string{
		mkResourceVirtualEnvironmentContainerMemoryDedicated,
		mkResourceVirtualEnvironmentContainerMemorySwap,
	})

	test.AssertValueTypes(t, memorySchema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentContainerMemoryDedicated: schema.TypeInt,
		mkResourceVirtualEnvironmentContainerMemorySwap:      schema.TypeInt,
	})

	networkInterfaceSchema := test.AssertNestedSchemaExistence(
		t,
		s,
		mkResourceVirtualEnvironmentContainerNetworkInterface,
	)

	test.AssertRequiredArguments(t, networkInterfaceSchema, []string{
		mkResourceVirtualEnvironmentContainerNetworkInterfaceName,
	})

	test.AssertOptionalArguments(t, networkInterfaceSchema, []string{
		mkResourceVirtualEnvironmentContainerNetworkInterfaceBridge,
		mkResourceVirtualEnvironmentContainerNetworkInterfaceEnabled,
		mkResourceVirtualEnvironmentContainerNetworkInterfaceMACAddress,
		mkResourceVirtualEnvironmentContainerNetworkInterfaceRateLimit,
		mkResourceVirtualEnvironmentContainerNetworkInterfaceVLANID,
		mkResourceVirtualEnvironmentContainerNetworkInterfaceMTU,
	})

	test.AssertValueTypes(t, networkInterfaceSchema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentContainerNetworkInterfaceBridge:     schema.TypeString,
		mkResourceVirtualEnvironmentContainerNetworkInterfaceEnabled:    schema.TypeBool,
		mkResourceVirtualEnvironmentContainerNetworkInterfaceMACAddress: schema.TypeString,
		mkResourceVirtualEnvironmentContainerNetworkInterfaceName:       schema.TypeString,
		mkResourceVirtualEnvironmentContainerNetworkInterfaceRateLimit:  schema.TypeFloat,
		mkResourceVirtualEnvironmentContainerNetworkInterfaceVLANID:     schema.TypeInt,
		mkResourceVirtualEnvironmentContainerNetworkInterfaceMTU:        schema.TypeInt,
	})

	operatingSystemSchema := test.AssertNestedSchemaExistence(
		t,
		s,
		mkResourceVirtualEnvironmentContainerOperatingSystem,
	)

	test.AssertRequiredArguments(t, operatingSystemSchema, []string{
		mkResourceVirtualEnvironmentContainerOperatingSystemTemplateFileID,
	})

	test.AssertOptionalArguments(t, operatingSystemSchema, []string{
		mkResourceVirtualEnvironmentContainerOperatingSystemType,
	})

	test.AssertValueTypes(t, operatingSystemSchema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentContainerOperatingSystemTemplateFileID: schema.TypeString,
		mkResourceVirtualEnvironmentContainerOperatingSystemType:           schema.TypeString,
	})
}
