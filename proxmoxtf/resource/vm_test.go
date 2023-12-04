/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package resource

import (
	"context"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/proxmox/ssh"
	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/test"
	"github.com/bpg/terraform-provider-proxmox/utils"
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

type nodeResolver struct {
	node ssh.ProxmoxNode
}

func (c *nodeResolver) Resolve(_ context.Context, _ string) (ssh.ProxmoxNode, error) {
	return c.node, nil
}

func TestVmGetDiskImagePath(t *testing.T) {
	ctx := context.TODO()
	testDiskImageId := utils.GetAnyStringEnv("PROXMOXTF_EXISTENT_TEST_DISKIMAGE_ID")
	testDiskImagePath := utils.GetAnyStringEnv("PROXMOXTF_EXISTENT_TEST_DISKIMAGE_PATH")

	if testDiskImageId == "" {
		t.Skip("SKIPPING TestVmGetDiskImagePath: Environment variable not set: 'PROXMOXTF_EXISTENT_TEST_DISKIMAGE_ID'")
	}
	if testDiskImagePath == "" {
		t.Skip("SKIPPING TestVmGetDiskImagePath: Environment variable not set: 'PROXMOXTF_EXISTENT_TEST_DISKIMAGE_PATH'")
	}

	endpoint := utils.GetAnyStringEnv("PROXMOX_VE_ENDPOINT")
	u, err := url.ParseRequestURI(endpoint)
	require.NoError(t, err)

	sshUsername := strings.Split(utils.GetAnyStringEnv("PROXMOX_VE_USERNAME"), "@")[0]
	sshAgentSocket := utils.GetAnyStringEnv("SSH_AUTH_SOCK", "PROXMOX_VE_SSH_AUTH_SOCK", "PM_VE_SSH_AUTH_SOCK")

	sshClient, err := ssh.NewClient(
		sshUsername, "", true, sshAgentSocket,
		&nodeResolver{
			node: ssh.ProxmoxNode{
				Address: u.Hostname(),
				Port:    22,
			},
		},
	)
	require.NoError(t, err)

	err, output := vmGetDiskImagePath(ctx, sshClient, u.Hostname(), testDiskImageId)
	t.Logf("TestVmGetDiskImagePath: stdout: %v", output)
	t.Logf("TestVmGetDiskImagePath: stderr: %v", err)

	require.Equal(t, testDiskImagePath, output)
}

func TestVmCopyDiskImageTmp(t *testing.T) {
	ctx := context.TODO()
	testDiskImageTmpPath := utils.GetAnyStringEnv("PROXMOXTF_EXISTENT_TEST_DISKIMAGE_PATH_TMP")
	testDiskImagePath := utils.GetAnyStringEnv("PROXMOXTF_EXISTENT_TEST_DISKIMAGE_PATH")

	if testDiskImageTmpPath == "" {
		t.Skip("SKIPPING TestVmCopyDiskImageTmp: Environment variable not set: 'PROXMOXTF_EXISTENT_TEST_DISKIMAGE_PATH_TMP'")
	}
	if testDiskImagePath == "" {
		t.Skip("SKIPPING TestVmCopyDiskImageTmp: Environment variable not set: 'PROXMOXTF_EXISTENT_TEST_DISKIMAGE_PATH'")
	}

	endpoint := utils.GetAnyStringEnv("PROXMOX_VE_ENDPOINT")
	u, err := url.ParseRequestURI(endpoint)
	require.NoError(t, err)

	sshUsername := strings.Split(utils.GetAnyStringEnv("PROXMOX_VE_USERNAME"), "@")[0]
	sshAgentSocket := utils.GetAnyStringEnv("SSH_AUTH_SOCK", "PROXMOX_VE_SSH_AUTH_SOCK", "PM_VE_SSH_AUTH_SOCK")

	sshClient, err := ssh.NewClient(
		sshUsername, "", true, sshAgentSocket,
		&nodeResolver{
			node: ssh.ProxmoxNode{
				Address: u.Hostname(),
				Port:    22,
			},
		},
	)
	require.NoError(t, err)

	err, output := vmCopyDiskImageTmp(ctx, sshClient, u.Hostname(), testDiskImagePath, testDiskImageTmpPath)
	t.Logf("TestVmCopyDiskImageTmp: stdout: %v", output)
	t.Logf("TestVmCopyDiskImageTmp: err: %v", err)

	require.Nil(t, err)
	require.Equal(t, "", output)

	// Proactively go and check the file exists after cp
	sshClient, err = ssh.NewClient(
		sshUsername, "", true, sshAgentSocket,
		&nodeResolver{
			node: ssh.ProxmoxNode{
				Address: u.Hostname(),
				Port:    22,
			},
		},
	)
	require.NoError(t, err)

	comErr, stdOut, stdErr := sshClient.ExecuteNodeCommand(
		ctx,
		u.Host,
		fmt.Sprintf(`ls "%v"`, testDiskImageTmpPath),
		[]string{},
	)
	require.NoError(t, comErr)
	require.Equal(t, fmt.Sprintf("%v\n", testDiskImageTmpPath), stdOut)
	require.Empty(t, stdErr)
}

func TestVmQemuResizeDiskImage(t *testing.T) {
	ctx := context.TODO()
	testDiskImageTmpPath := utils.GetAnyStringEnv("PROXMOXTF_EXISTENT_TEST_DISKIMAGE_PATH_TMP")
	testDiskImageTmpNewSize := utils.GetAnyStringEnv("PROXMOXTF_EXISTENT_TEST_DISKIMAGE_TMP_NEWSIZE")
	testDiskImageTmpFormat := utils.GetAnyStringEnv("PROXMOXTF_EXISTENT_TEST_DISKIMAGE_TMP_FORMAT")

	if testDiskImageTmpPath == "" {
		t.Skip("SKIPPING TestVmQemuResizeDiskImage: Environment variable not set: 'PROXMOXTF_EXISTENT_TEST_DISKIMAGE_PATH_TMP'")
	}
	if testDiskImageTmpNewSize == "" {
		t.Skip("SKIPPING TestVmQemuResizeDiskImage: Environment variable not set: 'PROXMOXTF_EXISTENT_TEST_DISKIMAGE_TMP_NEWSIZE'")
	}
	if testDiskImageTmpFormat == "" {
		t.Skip("SKIPPING TestVmQemuResizeDiskImage: Environment variable not set: 'PROXMOXTF_EXISTENT_TEST_DISKIMAGE_TMP_FORMAT'")
	}

	endpoint := utils.GetAnyStringEnv("PROXMOX_VE_ENDPOINT")
	u, err := url.ParseRequestURI(endpoint)
	require.NoError(t, err)

	sshUsername := strings.Split(utils.GetAnyStringEnv("PROXMOX_VE_USERNAME"), "@")[0]
	sshAgentSocket := utils.GetAnyStringEnv("SSH_AUTH_SOCK", "PROXMOX_VE_SSH_AUTH_SOCK", "PM_VE_SSH_AUTH_SOCK")

	sshClient, err := ssh.NewClient(
		sshUsername, "", true, sshAgentSocket,
		&nodeResolver{
			node: ssh.ProxmoxNode{
				Address: u.Hostname(),
				Port:    22,
			},
		},
	)
	require.NoError(t, err)

	err, output := vmQemuResizeDiskImage(ctx, sshClient, u.Hostname(), testDiskImageTmpFormat, testDiskImageTmpPath, testDiskImageTmpNewSize)
	t.Logf("TestVmQemuResizeDiskImage: stdout: %v", output)
	t.Logf("TestVmQemuResizeDiskImage: err: %v", err)

	require.Nil(t, err)
	require.Equal(t, "Image resized.", output)
}

func TestVmQemuImportDiskImage(t *testing.T) {
	ctx := context.TODO()
	testDiskImageTmpPath := utils.GetAnyStringEnv("PROXMOXTF_EXISTENT_TEST_DISKIMAGE_PATH_TMP")
	testVmId := utils.GetAnyStringEnv("PROXMOXTF_EXISTENT_TEST_VM_ID")
	testVmDatastoreId := utils.GetAnyStringEnv("PROXMOXTF_EXISTENT_TEST_VM_DATASTORE_ID")
	testDiskImageTmpFormat := utils.GetAnyStringEnv("PROXMOXTF_EXISTENT_TEST_DISKIMAGE_TMP_FORMAT")

	if testDiskImageTmpPath == "" {
		t.Skip("SKIPPING TestVmQemuImportDiskImage: Environment variable not set: 'PROXMOXTF_EXISTENT_TEST_DISKIMAGE_PATH_TMP'")
	}
	if testVmId == "" {
		t.Skip("SKIPPING TestVmQemuImportDiskImage: Environment variable not set: 'PROXMOXTF_EXISTENT_TEST_VM_ID'")
	}
	if testVmDatastoreId == "" {
		t.Skip("SKIPPING TestVmQemuResizeDiskImage: Environment variable not set: 'PROXMOXTF_EXISTENT_TEST_VM_DATASTORE_ID'")
	}
	if testDiskImageTmpFormat == "" {
		t.Skip("SKIPPING TestVmQemuResizeDiskImage: Environment variable not set: 'PROXMOXTF_EXISTENT_TEST_DISKIMAGE_TMP_FORMAT'")
	}

	endpoint := utils.GetAnyStringEnv("PROXMOX_VE_ENDPOINT")
	u, err := url.ParseRequestURI(endpoint)
	require.NoError(t, err)

	sshUsername := strings.Split(utils.GetAnyStringEnv("PROXMOX_VE_USERNAME"), "@")[0]
	sshAgentSocket := utils.GetAnyStringEnv("SSH_AUTH_SOCK", "PROXMOX_VE_SSH_AUTH_SOCK", "PM_VE_SSH_AUTH_SOCK")

	sshClient, err := ssh.NewClient(
		sshUsername, "", true, sshAgentSocket,
		&nodeResolver{
			node: ssh.ProxmoxNode{
				Address: u.Hostname(),
				Port:    22,
			},
		},
	)
	require.NoError(t, err)

	err, output := vmQemuImportDiskImage(ctx, sshClient, u.Hostname(), testVmId, testDiskImageTmpPath, testVmDatastoreId, testDiskImageTmpFormat)
	t.Logf("TestVmQemuImportDiskImage: stdout: %v", output)
	t.Logf("TestVmQemuImportDiskImage: err: %v", err)

	regexPattern := fmt.Sprintf(`^vm-%v-disk-\d*$`, testVmId)
	re := regexp.MustCompile(regexPattern)
	matches := re.FindStringSubmatch(output)

	require.Nil(t, err)
	require.Equal(t, matches[0], output)
}

func TestVmQemuVmSetDiskImageInterface(t *testing.T) {
	ctx := context.TODO()
	testVmDiskInterface := utils.GetAnyStringEnv("PROXMOXTF_EXISTENT_TEST_VM_DISK_INTERFACE_0")
	testVmId := utils.GetAnyStringEnv("PROXMOXTF_EXISTENT_TEST_VM_ID")
	testVmDatastoreId := utils.GetAnyStringEnv("PROXMOXTF_EXISTENT_TEST_VMLVM_DATASTORE_ID")
	testDiskImageId := utils.GetAnyStringEnv("PROXMOXTF_EXISTENT_TEST_VM_IMPORTED_DISKIMAGE_ID")

	if testVmDiskInterface == "" {
		t.Skip("SKIPPING TestVmQemuImportDiskImage: Environment variable not set: 'PROXMOXTF_EXISTENT_TEST_VM_DISK_INTERFACE_0'")
	}
	if testVmId == "" {
		t.Skip("SKIPPING TestVmQemuImportDiskImage: Environment variable not set: 'PROXMOXTF_EXISTENT_TEST_VM_ID'")
	}
	if testVmDatastoreId == "" {
		t.Skip("SKIPPING TestVmQemuResizeDiskImage: Environment variable not set: 'PROXMOXTF_EXISTENT_TEST_VMLVM_DATASTORE_ID'")
	}
	if testDiskImageId == "" {
		t.Skip("SKIPPING TestVmQemuResizeDiskImage: Environment variable not set: 'PROXMOXTF_EXISTENT_TEST_VM_IMPORTED_DISKIMAGE_ID'")
	}

	endpoint := utils.GetAnyStringEnv("PROXMOX_VE_ENDPOINT")
	u, err := url.ParseRequestURI(endpoint)
	require.NoError(t, err)

	sshUsername := strings.Split(utils.GetAnyStringEnv("PROXMOX_VE_USERNAME"), "@")[0]
	sshAgentSocket := utils.GetAnyStringEnv("SSH_AUTH_SOCK", "PROXMOX_VE_SSH_AUTH_SOCK", "PM_VE_SSH_AUTH_SOCK")

	sshClient, err := ssh.NewClient(
		sshUsername, "", true, sshAgentSocket,
		&nodeResolver{
			node: ssh.ProxmoxNode{
				Address: u.Hostname(),
				Port:    22,
			},
		},
	)
	require.NoError(t, err)

	err, output := vmQemuVmSetDiskImageInterface(ctx, sshClient, u.Hostname(), testVmId, testVmDatastoreId, testDiskImageId, testVmDiskInterface, []string{})
	t.Logf("TestVmQemuVmSetDiskImageInterface: stdout: %v", output)
	t.Logf("TestVmQemuVmSetDiskImageInterface: err: %v", err)

	regexPattern := fmt.Sprintf("update VM %v: -%v %v:%v%v\n", testVmId, testVmDiskInterface, testVmDatastoreId, testDiskImageId, "")
	re := regexp.MustCompile(regexPattern)
	matches := re.FindStringSubmatch(output)

	require.Nil(t, err)
	require.Equal(t, matches[0], output)
}

func TestVmQemuVmRemoveTmpDiskImageFile(t *testing.T) {
	ctx := context.TODO()
	remoteTmpDir := utils.GetAnyStringEnv("PROXMOXTF_TEST_REMOTE_TMP_DIR")

	if remoteTmpDir == "" {
		remoteTmpDir = "/tmp"
	}

	endpoint := utils.GetAnyStringEnv("PROXMOX_VE_ENDPOINT")
	u, err := url.ParseRequestURI(endpoint)
	require.NoError(t, err)

	sshUsername := strings.Split(utils.GetAnyStringEnv("PROXMOX_VE_USERNAME"), "@")[0]
	sshAgentSocket := utils.GetAnyStringEnv("SSH_AUTH_SOCK", "PROXMOX_VE_SSH_AUTH_SOCK", "PM_VE_SSH_AUTH_SOCK")

	sshClient, err := ssh.NewClient(
		sshUsername, "", true, sshAgentSocket,
		&nodeResolver{
			node: ssh.ProxmoxNode{
				Address: u.Hostname(),
				Port:    22,
			},
		},
	)
	require.NoError(t, err)

	// Create the tmp file to delete
	createdTmpFile := fmt.Sprintf(`%v/%v`, remoteTmpDir, "random-file.tmp")
	createTmpFileCommand := fmt.Sprintf(`echo 'Test' > %v`, createdTmpFile)
	tmpFileErr, tmpFilesout, tmpFilesErr := sshClient.ExecuteNodeCommand(ctx, u.Hostname(), createTmpFileCommand, []string{})

	if tmpFileErr != nil {
		t.Errorf("TestVmQemuVmRemoveTmpDiskImageFile: Error creating tmp file on server for test: file: %v - stdout: %v - stderr: %v", tmpFileErr, tmpFilesout, tmpFilesErr)
	}

	// Run the test
	sshClient, err = ssh.NewClient(
		sshUsername, "", true, sshAgentSocket,
		&nodeResolver{
			node: ssh.ProxmoxNode{
				Address: u.Hostname(),
				Port:    22,
			},
		},
	)
	require.NoError(t, err)

	err, output := vmQemuVmRemoveTmpDiskImageFile(ctx, sshClient, u.Hostname(), createdTmpFile)
	t.Logf("TestVmQemuVmRemoveTmpDiskImageFile: stdout: %v", output)
	t.Logf("TestVmQemuVmRemoveTmpDiskImageFile: err: %v", err)

	require.Nil(t, err)
	require.Equal(t, "", output)
}
