/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package resource

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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
		mkEnvironmentVariables,
		mkIDMap,
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
		mkCPU:                  schema.TypeList,
		mkDescription:          schema.TypeString,
		mkDisk:                 schema.TypeList,
		mkEnvironmentVariables: schema.TypeMap,
		mkIDMap:                schema.TypeList,
		mkInitialization:       schema.TypeList,
		mkHookScriptFileID:     schema.TypeString,
		mkMemory:               schema.TypeList,
		mkDevicePassthrough:    schema.TypeList,
		mkMountPoint:           schema.TypeList,
		mkOperatingSystem:      schema.TypeList,
		mkPoolID:               schema.TypeString,
		mkProtection:           schema.TypeBool,
		mkStarted:              schema.TypeBool,
		mkTags:                 schema.TypeList,
		mkTemplate:             schema.TypeBool,
		mkUnprivileged:         schema.TypeBool,
		mkStartOnBoot:          schema.TypeBool,
		mkFeatures:             schema.TypeList,
		mkVMID:                 schema.TypeInt,
	})

	cloneSchema := test.AssertNestedSchemaExistence(t, s, mkClone)

	test.AssertRequiredArguments(t, cloneSchema, []string{
		mkCloneVMID,
	})

	test.AssertOptionalArguments(t, cloneSchema, []string{
		mkCloneDatastoreID,
		mkCloneFull,
		mkCloneNodeName,
	})

	test.AssertValueTypes(t, cloneSchema, map[string]schema.ValueType{
		mkCloneDatastoreID: schema.TypeString,
		mkCloneFull:        schema.TypeBool,
		mkCloneNodeName:    schema.TypeString,
		mkCloneVMID:        schema.TypeInt,
	})

	cpuSchema := test.AssertNestedSchemaExistence(t, s, mkCPU)

	test.AssertOptionalArguments(t, cpuSchema, []string{
		mkCPUArchitecture,
		mkCPUCores,
		mkCPULimit,
		mkCPUUnits,
	})

	test.AssertValueTypes(t, cpuSchema, map[string]schema.ValueType{
		mkCPUArchitecture: schema.TypeString,
		mkCPUCores:        schema.TypeInt,
		mkCPULimit:        schema.TypeFloat,
		mkCPUUnits:        schema.TypeInt,
	})

	diskSchema := test.AssertNestedSchemaExistence(t, s, mkDisk)

	test.AssertOptionalArguments(t, diskSchema, []string{
		mkDiskDatastoreID,
	})

	test.AssertValueTypes(t, diskSchema, map[string]schema.ValueType{
		mkDiskDatastoreID:     schema.TypeString,
		mkDiskPathInDatastore: schema.TypeString,
	})

	featuresSchema := test.AssertNestedSchemaExistence(t, s, mkFeatures)

	test.AssertOptionalArguments(t, featuresSchema, []string{
		mkFeaturesNesting,
		mkFeaturesKeyControl,
		mkFeaturesFUSE,
		mkFeaturesMakeDeviceNode,
	})

	test.AssertValueTypes(t, featuresSchema, map[string]schema.ValueType{
		mkFeaturesNesting:        schema.TypeBool,
		mkFeaturesKeyControl:     schema.TypeBool,
		mkFeaturesFUSE:           schema.TypeBool,
		mkFeaturesMakeDeviceNode: schema.TypeBool,
	})

	initializationSchema := test.AssertNestedSchemaExistence(
		t,
		s,
		mkInitialization,
	)

	test.AssertOptionalArguments(t, initializationSchema, []string{
		mkInitializationDNS,
		mkInitializationEntrypoint,
		mkInitializationHostname,
		mkInitializationIPConfig,
		mkInitializationUserAccount,
	})

	test.AssertValueTypes(t, initializationSchema, map[string]schema.ValueType{
		mkInitializationDNS:         schema.TypeList,
		mkInitializationEntrypoint:  schema.TypeString,
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

	idmapSchema := test.AssertNestedSchemaExistence(t, s, mkIDMap)

	test.AssertRequiredArguments(t, idmapSchema, []string{
		mkIDMapType,
		mkIDMapContainerID,
		mkIDMapHostID,
		mkIDMapSize,
	})

	test.AssertValueTypes(t, idmapSchema, map[string]schema.ValueType{
		mkIDMapType:        schema.TypeString,
		mkIDMapContainerID: schema.TypeInt,
		mkIDMapHostID:      schema.TypeInt,
		mkIDMapSize:        schema.TypeInt,
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
		mkMountPointACL:             schema.TypeBool,
		mkMountPointBackup:          schema.TypeBool,
		mkMountPointMountOptions:    schema.TypeList,
		mkMountPointPath:            schema.TypeString,
		mkMountPointQuota:           schema.TypeBool,
		mkMountPointReadOnly:        schema.TypeBool,
		mkMountPointReplicate:       schema.TypeBool,
		mkMountPointShared:          schema.TypeBool,
		mkMountPointSize:            schema.TypeString,
		mkMountPointVolume:          schema.TypeString,
		mkMountPointPathInDatastore: schema.TypeString,
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

// mountPointStateAttrs builds the flatmap state attributes for a list of mount points.
func mountPointStateAttrs(mountPoints ...map[string]string) map[string]string {
	attrs := map[string]string{
		mkVMID:                     "100",
		mkNodeName:                 "pve",
		mkMountPoint + ".#":        strconv.Itoa(len(mountPoints)),
		mkNetworkInterface + ".#":  "0",
		mkDisk + ".#":              "0",
		mkInitialization + ".#":    "0",
		mkOperatingSystem + ".#":   "0",
		mkDevicePassthrough + ".#": "0",
		mkIDMap + ".#":             "0",
		// Unrelated ForceNew attributes: pinned to their defaults so they do not contribute
		// RequiresNew to the diff under test.
		mkTemplate:     "false",
		mkUnprivileged: "false",
		mkPoolID:       "",
	}

	for i, mp := range mountPoints {
		prefix := fmt.Sprintf("%s.%d.", mkMountPoint, i)
		attrs[prefix+mkMountPointMountOptions+".#"] = "0"
		attrs[prefix+mkMountPointACL] = "false"
		attrs[prefix+mkMountPointBackup] = "false"
		attrs[prefix+mkMountPointQuota] = "false"
		attrs[prefix+mkMountPointReadOnly] = "false"
		attrs[prefix+mkMountPointReplicate] = "true"
		attrs[prefix+mkMountPointShared] = "false"
		attrs[prefix+mkMountPointSize] = ""

		for k, v := range mp {
			attrs[prefix+k] = v
		}
	}

	return attrs
}

// TestContainerMountPointForceNew verifies that mount point changes only force container
// recreation when the volume of an *existing* mount point is repointed or resized. Adding,
// removing, or reconfiguring mount points must be handled in place by containerUpdate.
//
// See https://github.com/bpg/terraform-provider-proxmox/issues/3028.
func TestContainerMountPointForceNew(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		state             []map[string]string
		config            []any
		expectRequiresNew bool
	}{
		{
			name: "adding a mount point is an in-place update",
			state: []map[string]string{
				{mkMountPointVolume: "local-lvm:vm-100-disk-1", mkMountPointPath: "/mnt/local1"},
			},
			config: []any{
				map[string]any{mkMountPointVolume: "local-lvm", mkMountPointPath: "/mnt/local1"},
				map[string]any{mkMountPointVolume: "/mnt/pve/share", mkMountPointPath: "/mnt/share"},
			},
			expectRequiresNew: false,
		},
		{
			name:  "adding the first mount point is an in-place update",
			state: nil,
			config: []any{
				map[string]any{mkMountPointVolume: "/mnt/pve/share", mkMountPointPath: "/mnt/share"},
			},
			expectRequiresNew: false,
		},
		{
			name: "removing a mount point is an in-place update",
			state: []map[string]string{
				{mkMountPointVolume: "local-lvm:vm-100-disk-1", mkMountPointPath: "/mnt/local1"},
				{mkMountPointVolume: "/mnt/pve/share", mkMountPointPath: "/mnt/share"},
			},
			config: []any{
				map[string]any{mkMountPointVolume: "local-lvm", mkMountPointPath: "/mnt/local1"},
			},
			expectRequiresNew: false,
		},
		{
			name: "toggling read_only on an existing mount point is an in-place update",
			state: []map[string]string{
				{mkMountPointVolume: "/mnt/pve/share", mkMountPointPath: "/mnt/share"},
			},
			config: []any{
				map[string]any{
					mkMountPointVolume:   "/mnt/pve/share",
					mkMountPointPath:     "/mnt/share",
					mkMountPointReadOnly: true,
				},
			},
			expectRequiresNew: false,
		},
		{
			name: "repointing the volume of an existing mount point forces recreation",
			state: []map[string]string{
				{mkMountPointVolume: "/mnt/pve/share", mkMountPointPath: "/mnt/share"},
			},
			config: []any{
				map[string]any{mkMountPointVolume: "/mnt/pve/other", mkMountPointPath: "/mnt/share"},
			},
			expectRequiresNew: true,
		},
		{
			name: "resizing an existing mount point forces recreation",
			state: []map[string]string{
				{
					mkMountPointVolume: "local-lvm:vm-100-disk-1",
					mkMountPointPath:   "/mnt/local1",
					mkMountPointSize:   "4G",
				},
			},
			config: []any{
				map[string]any{
					mkMountPointVolume: "local-lvm",
					mkMountPointPath:   "/mnt/local1",
					mkMountPointSize:   "8G",
				},
			},
			expectRequiresNew: true,
		},
		{
			// PVE resolves a datastore-only volume to a concrete volume ID; the resulting
			// state-vs-config drift must not be mistaken for a repoint.
			name: "resolved volume id is not treated as a repoint",
			state: []map[string]string{
				{mkMountPointVolume: "local-lvm:vm-100-disk-1", mkMountPointPath: "/mnt/local1"},
			},
			config: []any{
				map[string]any{mkMountPointVolume: "local-lvm", mkMountPointPath: "/mnt/local1"},
			},
			expectRequiresNew: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			container := Container()

			state := &terraform.InstanceState{
				ID:         "100",
				Attributes: mountPointStateAttrs(tt.state...),
			}

			config := terraform.NewResourceConfigRaw(map[string]any{
				mkNodeName:   "pve",
				mkVMID:       100,
				mkMountPoint: tt.config,
			})

			diff, err := container.Diff(t.Context(), state, config, nil)
			require.NoError(t, err)
			require.NotNil(t, diff, "expected a non-empty diff")

			assert.Equal(t, tt.expectRequiresNew, diff.RequiresNew(),
				"unexpected RequiresNew; diff attributes: %v", diff.Attributes)
		})
	}
}

func TestInitializationEntrypointValidation(t *testing.T) {
	t.Parallel()

	// Same character set Proxmox rejects (see EnvironmentVariablesValidator).
	re := regexp.MustCompile(`^[^\x00-\x08\x10-\x1F\x7F]+$`)

	validCases := []string{
		"/sbin/init",
		"/usr/bin/my-init --option",
		"my-init",
		"/path/with spaces/init",
		"init\twith-tab",
	}

	for _, v := range validCases {
		assert.True(t, re.MatchString(v), "expected %q to be valid", v)
	}

	invalidCases := []string{
		"\x00null",
		"\x08backspace",
		"\x10dle",
		"\x1Fescape",
		"\x7Fdel",
	}

	for _, v := range invalidCases {
		assert.False(t, re.MatchString(v), "expected %q to be invalid", v)
	}
}

func TestInitializationDnsBlockDiffIgnore(t *testing.T) {
	t.Parallel()

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
