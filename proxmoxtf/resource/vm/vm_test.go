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

	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes"
	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/resource/vm/disk"
	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/resource/vm/network"
	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/test"
)

// TestVMInstantiation tests whether the VM instance can be instantiated.
func TestVMInstantiation(t *testing.T) {
	t.Parallel()

	r := VM()
	if r == nil {
		t.Fatalf("Cannot instantiate VM")
	}
}

// TestVMSchema tests the VM schema.
func TestVMSchema(t *testing.T) {
	t.Parallel()

	s := VM().Schema

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
		disk.MkDisk,
		mkEFIDisk,
		mkInitialization,
		mkHostPCI,
		mkHostUSB,
		mkKeyboardLayout,
		mkKVMArguments,
		mkMachine,
		mkMemory,
		mkName,
		mkNodeName,
		mkNodeNames,
		network.MkNetworkDevice,
		mkOperatingSystem,
		mkPoolID,
		mkSerialDevice,
		mkStarted,
		mkTabletDevice,
		mkTemplate,
		mkVirtiofs,
		mkVMID,
		mkSCSIHardware,
	})

	test.AssertValueTypes(t, s, map[string]schema.ValueType{
		mkACPI:            schema.TypeBool,
		mkAgent:           schema.TypeList,
		mkAudioDevice:     schema.TypeList,
		mkBIOS:            schema.TypeString,
		mkBootOrder:       schema.TypeList,
		mkCDROM:           schema.TypeList,
		mkCPU:             schema.TypeList,
		mkDescription:     schema.TypeString,
		disk.MkDisk:       schema.TypeList,
		mkEFIDisk:         schema.TypeList,
		mkHostPCI:         schema.TypeList,
		mkHostUSB:         schema.TypeList,
		mkInitialization:  schema.TypeList,
		mkKeyboardLayout:  schema.TypeString,
		mkKVMArguments:    schema.TypeString,
		mkMachine:         schema.TypeString,
		mkMemory:          schema.TypeList,
		mkName:            schema.TypeString,
		mkNodeName:        schema.TypeString,
		mkNodeNames:       schema.TypeSet,
		mkOperatingSystem: schema.TypeList,
		mkPoolID:          schema.TypeString,
		mkSerialDevice:    schema.TypeList,
		mkStarted:         schema.TypeBool,
		mkTabletDevice:    schema.TypeBool,
		mkTemplate:        schema.TypeBool,
		mkVirtiofs:        schema.TypeList,
		mkVMID:            schema.TypeInt,
		mkSCSIHardware:    schema.TypeString,
	})

	test.AssertComputedAttributes(t, s, []string{
		mkCurrentNodeName,
	})

	test.AssertExactlyOneOfArguments(t, s, map[string][]string{
		mkNodeName:  {mkNodeName, mkNodeNames},
		mkNodeNames: {mkNodeName, mkNodeNames},
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
		mkInitializationFileFormat,
		mkInitializationDNS,
		mkInitializationIPConfig,
		mkInitializationUserAccount,
	})

	test.AssertValueTypes(t, initializationSchema, map[string]schema.ValueType{
		mkInitializationDatastoreID: schema.TypeString,
		mkInitializationInterface:   schema.TypeString,
		mkInitializationFileFormat:  schema.TypeString,
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
		mkInitializationDNSServers,
	})

	test.AssertValueTypes(t, initializationDNSSchema, map[string]schema.ValueType{
		mkInitializationDNSDomain:  schema.TypeString,
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

	numaSchema := test.AssertNestedSchemaExistence(t, s, mkNUMA)

	test.AssertOptionalArguments(t, numaSchema, []string{
		mkNUMAHostNodeNames,
		mkNUMAPolicy,
	})

	test.AssertValueTypes(t, numaSchema, map[string]schema.ValueType{
		mkNUMADevice:        schema.TypeString,
		mkNUMACPUIDs:        schema.TypeString,
		mkNUMAMemory:        schema.TypeInt,
		mkNUMAHostNodeNames: schema.TypeString,
		mkNUMAPolicy:        schema.TypeString,
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

	virtiofsSchema := test.AssertNestedSchemaExistence(t, s, mkVirtiofs)

	test.AssertOptionalArguments(t, virtiofsSchema, []string{
		mkVirtiofsCache,
		mkVirtiofsDirectIO,
		mkVirtiofsExposeACL,
		mkVirtiofsExposeXAttr,
	})

	test.AssertValueTypes(t, virtiofsSchema, map[string]schema.ValueType{
		mkVirtiofsMapping:     schema.TypeString,
		mkVirtiofsCache:       schema.TypeString,
		mkVirtiofsDirectIO:    schema.TypeBool,
		mkVirtiofsExposeACL:   schema.TypeBool,
		mkVirtiofsExposeXAttr: schema.TypeBool,
	})

	vgaSchema := test.AssertNestedSchemaExistence(t, s, mkVGA)

	test.AssertOptionalArguments(t, vgaSchema, []string{
		mkVGAMemory,
		mkVGAType,
	})

	test.AssertValueTypes(t, vgaSchema, map[string]schema.ValueType{
		mkVGAMemory: schema.TypeInt,
		mkVGAType:   schema.TypeString,
	})
}

func TestHotplugContains(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		hotplug string
		feature string
		want    bool
	}{
		{"empty string disables all", "", "memory", false},
		{"zero disables all", "0", "memory", false},
		{"one enables all", "1", "memory", true},
		{"one enables cpu", "1", "cpu", true},
		{"feature present", "disk,network,memory", "memory", true},
		{"feature absent", "disk,network,usb", "memory", false},
		{"single feature match", "memory", "memory", true},
		{"single feature no match", "cpu", "memory", false},
		{"default proxmox value", "disk,network,usb", "cpu", false},
		{"cpu in list", "disk,cpu,network", "cpu", true},
		{"pve default includes disk", "network,disk,usb", "disk", true},
		{"pve default includes network", "network,disk,usb", "network", true},
		{"pve default includes usb", "network,disk,usb", "usb", true},
		{"pve default excludes cpu", "network,disk,usb", "cpu", false},
		{"pve default excludes memory", "network,disk,usb", "memory", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := hotplugContains(tt.hotplug, tt.feature)
			require.Equal(t, tt.want, got)
		})
	}
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

func TestSelectLeastUtilizedNodeFromList(t *testing.T) {
	t.Parallel()

	online := "online"
	offline := "offline"
	lowCPU := 0.10
	highCPU := 0.80
	midCPU := 0.40
	lowMemTotal := int64(100)
	lowMemUsed := int64(20)
	highMemTotal := int64(100)
	highMemUsed := int64(70)
	tieMemTotal := int64(100)
	tieMemUsed := int64(20)

	tests := []struct {
		name      string
		nodeNames []string
		nodes     []*nodes.ListResponseData
		want      string
		wantErr   string
	}{
		{
			name:      "prefers lower memory utilization",
			nodeNames: []string{"node-a", "node-b"},
			nodes: []*nodes.ListResponseData{
				{Name: "node-a", Status: &online, CPUUtilization: &midCPU, MemoryAvailable: &highMemTotal, MemoryUsed: &highMemUsed},
				{Name: "node-b", Status: &online, CPUUtilization: &highCPU, MemoryAvailable: &lowMemTotal, MemoryUsed: &lowMemUsed},
			},
			want: "node-b",
		},
		{
			name:      "uses cpu as tie breaker",
			nodeNames: []string{"node-a", "node-b"},
			nodes: []*nodes.ListResponseData{
				{Name: "node-a", Status: &online, CPUUtilization: &highCPU, MemoryAvailable: &tieMemTotal, MemoryUsed: &tieMemUsed},
				{Name: "node-b", Status: &online, CPUUtilization: &lowCPU, MemoryAvailable: &tieMemTotal, MemoryUsed: &tieMemUsed},
			},
			want: "node-b",
		},
		{
			name:      "ignores offline nodes",
			nodeNames: []string{"node-a", "node-b"},
			nodes: []*nodes.ListResponseData{
				{Name: "node-a", Status: &offline, CPUUtilization: &lowCPU, MemoryAvailable: &lowMemTotal, MemoryUsed: &lowMemUsed},
				{Name: "node-b", Status: &online, CPUUtilization: &midCPU, MemoryAvailable: &highMemTotal, MemoryUsed: &highMemUsed},
			},
			wantErr: "failed to find online nodes from node_names: node-a",
		},
		{
			name:      "fails when no configured node is found",
			nodeNames: []string{"node-a"},
			nodes:     []*nodes.ListResponseData{},
			wantErr:   "failed to find online nodes from node_names: node-a",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := selectLeastUtilizedNodeFromList(tt.nodeNames, tt.nodes)

			if tt.wantErr != "" {
				require.EqualError(t, err, tt.wantErr)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestNodePlacementHelpers(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		config      map[string]any
		wantNodes   []string
		wantMulti   bool
		wantCurrent string
		wantErr     string
	}{
		{
			name: "single node placement",
			config: map[string]any{
				mkNodeName: "node-a",
			},
			wantNodes: []string{"node-a"},
			wantMulti: false,
			wantErr:   "failed to determine current VM node from state",
		},
		{
			name: "multi node placement",
			config: map[string]any{
				mkNodeNames: []any{"node-b", "node-a"},
			},
			wantNodes: []string{"node-a", "node-b"},
			wantMulti: true,
			wantErr:   "failed to determine current VM node from state",
		},
		{
			name: "current node uses computed state",
			config: map[string]any{
				mkNodeNames:       []any{"node-b", "node-a"},
				mkCurrentNodeName: "node-b",
			},
			wantCurrent: "node-b",
			wantNodes:   []string{"node-a", "node-b"},
			wantMulti:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			d := schema.TestResourceDataRaw(t, VM().Schema, tt.config)

			require.Equal(t, tt.wantMulti, isMultiNodePlacementConfigured(d))
			require.Equal(t, tt.wantNodes, getConfiguredVMNodeNames(d))

			currentNode, err := getCurrentVMNodeName(d)
			if tt.wantErr != "" {
				require.EqualError(t, err, tt.wantErr)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.wantCurrent, currentNode)
		})
	}
}
