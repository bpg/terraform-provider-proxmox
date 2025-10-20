/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package resource

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

	s := Container().Schema

	test.AssertRequiredArguments(t, s, []string{
		mkNodeName,
	})

	test.AssertOptionalArguments(t, s, []string{
		mkCPU,
		mkDescription,
		mkDisk,
		mkInitialization,
		mkHookScriptFileID,
		mkMemory,
		mkDevicePassthrough,
		mkMountPoint,
		mkOperatingSystem,
		mkPoolID,
		mkProtection,
		mkStarted,
		mkTags,
		mkTemplate,
		mkUnprivileged,
		mkStartOnBoot,
		mkFeatures,
		mkVMID,
	})

	test.AssertValueTypes(t, s, map[string]schema.ValueType{
		mkCPU:               schema.TypeList,
		mkDescription:       schema.TypeString,
		mkDisk:              schema.TypeList,
		mkInitialization:    schema.TypeList,
		mkHookScriptFileID:  schema.TypeString,
		mkMemory:            schema.TypeList,
		mkDevicePassthrough: schema.TypeList,
		mkMountPoint:        schema.TypeList,
		mkOperatingSystem:   schema.TypeList,
		mkPoolID:            schema.TypeString,
		mkProtection:        schema.TypeBool,
		mkStarted:           schema.TypeBool,
		mkTags:              schema.TypeList,
		mkTemplate:          schema.TypeBool,
		mkUnprivileged:      schema.TypeBool,
		mkStartOnBoot:       schema.TypeBool,
		mkFeatures:          schema.TypeList,
		mkVMID:              schema.TypeInt,
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
		mkCPUUnits,
	})

	test.AssertValueTypes(t, cpuSchema, map[string]schema.ValueType{
		mkCPUArchitecture: schema.TypeString,
		mkCPUCores:        schema.TypeInt,
		mkCPUUnits:        schema.TypeInt,
	})

	diskSchema := test.AssertNestedSchemaExistence(t, s, mkDisk)

	test.AssertOptionalArguments(t, diskSchema, []string{
		mkDiskDatastoreID,
	})

	test.AssertValueTypes(t, diskSchema, map[string]schema.ValueType{
		mkDiskDatastoreID: schema.TypeString,
	})

	featuresSchema := test.AssertNestedSchemaExistence(t, s, mkFeatures)

	test.AssertOptionalArguments(t, featuresSchema, []string{
		mkFeaturesNesting,
		mkFeaturesKeyControl,
		mkFeaturesFUSE,
	})

	test.AssertValueTypes(t, featuresSchema, map[string]schema.ValueType{
		mkFeaturesNesting:    schema.TypeBool,
		mkFeaturesKeyControl: schema.TypeBool,
		mkFeaturesFUSE:       schema.TypeBool,
	})

	initializationSchema := test.AssertNestedSchemaExistence(
		t,
		s,
		mkInitialization,
	)

	test.AssertOptionalArguments(t, initializationSchema, []string{
		mkInitializationDNS,
		mkInitializationHostname,
		mkInitializationIPConfig,
		mkInitializationUserAccount,
	})

	test.AssertValueTypes(t, initializationSchema, map[string]schema.ValueType{
		mkInitializationDNS:         schema.TypeList,
		mkInitializationHostname:    schema.TypeString,
		mkInitializationIPConfig:    schema.TypeList,
		mkInitializationUserAccount: schema.TypeList,
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
	})

	test.AssertValueTypes(t, initializationUserAccountSchema, map[string]schema.ValueType{
		mkInitializationUserAccountKeys:     schema.TypeList,
		mkInitializationUserAccountPassword: schema.TypeString,
	})

	memorySchema := test.AssertNestedSchemaExistence(t, s, mkMemory)

	test.AssertOptionalArguments(t, memorySchema, []string{
		mkMemoryDedicated,
		mkMemorySwap,
	})

	test.AssertValueTypes(t, memorySchema, map[string]schema.ValueType{
		mkMemoryDedicated: schema.TypeInt,
		mkMemorySwap:      schema.TypeInt,
	})

	devicePassthroughSchema := test.AssertNestedSchemaExistence(t, s, mkDevicePassthrough)

	test.AssertRequiredArguments(t, devicePassthroughSchema, []string{
		mkDevicePassthroughPath,
	})

	test.AssertOptionalArguments(t, devicePassthroughSchema, []string{
		mkDevicePassthroughDenyWrite,
		mkDevicePassthroughGID,
		mkDevicePassthroughMode,
		mkDevicePassthroughUID,
	})

	test.AssertValueTypes(t, devicePassthroughSchema, map[string]schema.ValueType{
		mkDevicePassthroughDenyWrite: schema.TypeBool,
		mkDevicePassthroughGID:       schema.TypeInt,
		mkDevicePassthroughMode:      schema.TypeString,
		mkDevicePassthroughPath:      schema.TypeString,
		mkDevicePassthroughUID:       schema.TypeInt,
	})

	mountPointSchema := test.AssertNestedSchemaExistence(t, s, mkMountPoint)

	test.AssertOptionalArguments(t, mountPointSchema, []string{
		mkMountPointACL,
		mkMountPointBackup,
		mkMountPointMountOptions,
		mkMountPointQuota,
		mkMountPointReadOnly,
		mkMountPointReplicate,
		mkMountPointShared,
		mkMountPointSize,
	})

	test.AssertValueTypes(t, mountPointSchema, map[string]schema.ValueType{
		mkMountPointACL:          schema.TypeBool,
		mkMountPointBackup:       schema.TypeBool,
		mkMountPointMountOptions: schema.TypeList,
		mkMountPointPath:         schema.TypeString,
		mkMountPointQuota:        schema.TypeBool,
		mkMountPointReadOnly:     schema.TypeBool,
		mkMountPointReplicate:    schema.TypeBool,
		mkMountPointShared:       schema.TypeBool,
		mkMountPointSize:         schema.TypeString,
		mkMountPointVolume:       schema.TypeString,
	})

	networkInterfaceSchema := test.AssertNestedSchemaExistence(
		t,
		s,
		mkNetworkInterface,
	)

	test.AssertRequiredArguments(t, networkInterfaceSchema, []string{
		mkNetworkInterfaceName,
	})

	test.AssertOptionalArguments(t, networkInterfaceSchema, []string{
		mkNetworkInterfaceBridge,
		mkNetworkInterfaceEnabled,
		mkNetworkInterfaceMACAddress,
		mkNetworkInterfaceRateLimit,
		mkNetworkInterfaceVLANID,
		mkNetworkInterfaceMTU,
	})

	test.AssertValueTypes(t, networkInterfaceSchema, map[string]schema.ValueType{
		mkNetworkInterfaceBridge:     schema.TypeString,
		mkNetworkInterfaceEnabled:    schema.TypeBool,
		mkNetworkInterfaceMACAddress: schema.TypeString,
		mkNetworkInterfaceName:       schema.TypeString,
		mkNetworkInterfaceRateLimit:  schema.TypeFloat,
		mkNetworkInterfaceVLANID:     schema.TypeInt,
		mkNetworkInterfaceMTU:        schema.TypeInt,
	})

	operatingSystemSchema := test.AssertNestedSchemaExistence(
		t,
		s,
		mkOperatingSystem,
	)

	test.AssertRequiredArguments(t, operatingSystemSchema, []string{
		mkOperatingSystemTemplateFileID,
	})

	test.AssertOptionalArguments(t, operatingSystemSchema, []string{
		mkOperatingSystemType,
	})

	test.AssertValueTypes(t, operatingSystemSchema, map[string]schema.ValueType{
		mkOperatingSystemTemplateFileID: schema.TypeString,
		mkOperatingSystemType:           schema.TypeString,
	})
}

func TestInitializationDnsBlockDiffIgnore(t *testing.T) {
	container := Container()

	tests := []struct {
		domain   string
		server   string
		servers  []string
		expected bool
	}{
		{"somedomain", "", []string{}, false},
		{"somedomain", "127.0.0.1", []string{}, false},
		{"somedomain", "", []string{"127.0.0.1"}, false},
		{"", "127.0.0.1", []string{}, false},
		{"", "", []string{"127.0.0.1"}, false},
		{"", "", []string{}, true},
	}

	for _, tt := range tests {
		d := container.TestResourceData()
		dnsBlockKey := mkInitialization + ".0." + mkInitializationDNS
		m := make(map[string]any)
		m[mkInitializationDNS] = []any{
			map[string]any{
				mkInitializationDNSDomain:  tt.domain,
				mkInitializationDNSServer:  tt.server,
				mkInitializationDNSServers: tt.servers,
			},
		}
		err := d.Set(mkInitialization, []any{m})
		require.NoError(t, err)
		actual := skipDnsDiffIfEmpty(dnsBlockKey+".#", "0", "1", d)
		assert.Equal(t, tt.expected, actual)
	}
}
