/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package resource

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/bpg/proxmox-api/cluster"
	"github.com/bpg/proxmox-api/nodes/vms"

	"github.com/bpg/terraform-provider-proxmox/proxmoxtf"
	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/resource/validator"

	"github.com/bpg/proxmox-api/types"
)

const (
	dvResourceVirtualEnvironmentVMRebootAfterCreation               = false
	dvResourceVirtualEnvironmentVMOnBoot                            = true
	dvResourceVirtualEnvironmentVMACPI                              = true
	dvResourceVirtualEnvironmentVMAgentEnabled                      = false
	dvResourceVirtualEnvironmentVMAgentTimeout                      = "15m"
	dvResourceVirtualEnvironmentVMAgentTrim                         = false
	dvResourceVirtualEnvironmentVMAgentType                         = "virtio"
	dvResourceVirtualEnvironmentVMAudioDeviceDevice                 = "intel-hda"
	dvResourceVirtualEnvironmentVMAudioDeviceDriver                 = "spice"
	dvResourceVirtualEnvironmentVMAudioDeviceEnabled                = true
	dvResourceVirtualEnvironmentVMBIOS                              = "seabios"
	dvResourceVirtualEnvironmentVMCDROMEnabled                      = false
	dvResourceVirtualEnvironmentVMCDROMFileID                       = ""
	dvResourceVirtualEnvironmentVMCloneDatastoreID                  = ""
	dvResourceVirtualEnvironmentVMCloneNodeName                     = ""
	dvResourceVirtualEnvironmentVMCloneFull                         = true
	dvResourceVirtualEnvironmentVMCloneRetries                      = 1
	dvResourceVirtualEnvironmentVMCPUArchitecture                   = "x86_64"
	dvResourceVirtualEnvironmentVMCPUCores                          = 1
	dvResourceVirtualEnvironmentVMCPUHotplugged                     = 0
	dvResourceVirtualEnvironmentVMCPUNUMA                           = false
	dvResourceVirtualEnvironmentVMCPUSockets                        = 1
	dvResourceVirtualEnvironmentVMCPUType                           = "qemu64"
	dvResourceVirtualEnvironmentVMCPUUnits                          = 1024
	dvResourceVirtualEnvironmentVMDescription                       = ""
	dvResourceVirtualEnvironmentVMDiskInterface                     = "scsi0"
	dvResourceVirtualEnvironmentVMDiskDatastoreID                   = "local-lvm"
	dvResourceVirtualEnvironmentVMDiskFileFormat                    = "qcow2"
	dvResourceVirtualEnvironmentVMDiskFileID                        = ""
	dvResourceVirtualEnvironmentVMDiskSize                          = 8
	dvResourceVirtualEnvironmentVMDiskIOThread                      = false
	dvResourceVirtualEnvironmentVMDiskSSD                           = false
	dvResourceVirtualEnvironmentVMDiskDiscard                       = ""
	dvResourceVirtualEnvironmentVMDiskSpeedRead                     = 0
	dvResourceVirtualEnvironmentVMDiskSpeedReadBurstable            = 0
	dvResourceVirtualEnvironmentVMDiskSpeedWrite                    = 0
	dvResourceVirtualEnvironmentVMDiskSpeedWriteBurstable           = 0
	dvResourceVirtualEnvironmentVMEFIDiskDatastoreID                = "local-lvm"
	dvResourceVirtualEnvironmentVMEFIDiskFileFormat                 = "qcow2"
	dvResourceVirtualEnvironmentVMEFIDiskType                       = "2m"
	dvResourceVirtualEnvironmentVMEFIDiskPreEnrolledKeys            = false
	dvResourceVirtualEnvironmentVMInitializationDatastoreID         = "local-lvm"
	dvResourceVirtualEnvironmentVMInitializationDNSDomain           = ""
	dvResourceVirtualEnvironmentVMInitializationDNSServer           = ""
	dvResourceVirtualEnvironmentVMInitializationIPConfigIPv4Address = ""
	dvResourceVirtualEnvironmentVMInitializationIPConfigIPv4Gateway = ""
	dvResourceVirtualEnvironmentVMInitializationIPConfigIPv6Address = ""
	dvResourceVirtualEnvironmentVMInitializationIPConfigIPv6Gateway = ""
	dvResourceVirtualEnvironmentVMInitializationUserAccountPassword = ""
	dvResourceVirtualEnvironmentVMInitializationUserDataFileID      = ""
	dvResourceVirtualEnvironmentVMInitializationVendorDataFileID    = ""
	dvResourceVirtualEnvironmentVMInitializationNetworkDataFileID   = ""
	dvResourceVirtualEnvironmentVMInitializationMetaDataFileID      = ""
	dvResourceVirtualEnvironmentVMInitializationType                = ""
	dvResourceVirtualEnvironmentVMKeyboardLayout                    = "en-us"
	dvResourceVirtualEnvironmentVMKVMArguments                      = ""
	dvResourceVirtualEnvironmentVMMachineType                       = ""
	dvResourceVirtualEnvironmentVMMemoryDedicated                   = 512
	dvResourceVirtualEnvironmentVMMemoryFloating                    = 0
	dvResourceVirtualEnvironmentVMMemoryShared                      = 0
	dvResourceVirtualEnvironmentVMName                              = ""
	dvResourceVirtualEnvironmentVMNetworkDeviceBridge               = "vmbr0"
	dvResourceVirtualEnvironmentVMNetworkDeviceEnabled              = true
	dvResourceVirtualEnvironmentVMNetworkDeviceFirewall             = false
	dvResourceVirtualEnvironmentVMNetworkDeviceModel                = "virtio"
	dvResourceVirtualEnvironmentVMNetworkDeviceRateLimit            = 0
	dvResourceVirtualEnvironmentVMNetworkDeviceVLANID               = 0
	dvResourceVirtualEnvironmentVMNetworkDeviceMTU                  = 0
	dvResourceVirtualEnvironmentVMOperatingSystemType               = "other"
	dvResourceVirtualEnvironmentVMPoolID                            = ""
	dvResourceVirtualEnvironmentVMSerialDeviceDevice                = "socket"
	dvResourceVirtualEnvironmentVMStarted                           = true
	dvResourceVirtualEnvironmentVMTabletDevice                      = true
	dvResourceVirtualEnvironmentVMTemplate                          = false
	dvResourceVirtualEnvironmentVMTimeoutClone                      = 1800
	dvResourceVirtualEnvironmentVMTimeoutMoveDisk                   = 1800
	dvResourceVirtualEnvironmentVMTimeoutReboot                     = 1800
	dvResourceVirtualEnvironmentVMTimeoutShutdownVM                 = 1800
	dvResourceVirtualEnvironmentVMTimeoutStartVM                    = 1800
	dvResourceVirtualEnvironmentVMTimeoutStopVM                     = 300
	dvResourceVirtualEnvironmentVMVGAEnabled                        = true
	dvResourceVirtualEnvironmentVMVGAMemory                         = 16
	dvResourceVirtualEnvironmentVMVGAType                           = "std"
	dvResourceVirtualEnvironmentVMSCSIHardware                      = "virtio-scsi-pci"

	maxResourceVirtualEnvironmentVMAudioDevices   = 1
	maxResourceVirtualEnvironmentVMNetworkDevices = 8
	maxResourceVirtualEnvironmentVMSerialDevices  = 4

	mkResourceVirtualEnvironmentVMRebootAfterCreation               = "reboot"
	mkResourceVirtualEnvironmentVMOnBoot                            = "on_boot"
	mkResourceVirtualEnvironmentVMBootOrder                         = "boot_order"
	mkResourceVirtualEnvironmentVMACPI                              = "acpi"
	mkResourceVirtualEnvironmentVMAgent                             = "agent"
	mkResourceVirtualEnvironmentVMAgentEnabled                      = "enabled"
	mkResourceVirtualEnvironmentVMAgentTimeout                      = "timeout"
	mkResourceVirtualEnvironmentVMAgentTrim                         = "trim"
	mkResourceVirtualEnvironmentVMAgentType                         = "type"
	mkResourceVirtualEnvironmentVMAudioDevice                       = "audio_device"
	mkResourceVirtualEnvironmentVMAudioDeviceDevice                 = "device"
	mkResourceVirtualEnvironmentVMAudioDeviceDriver                 = "driver"
	mkResourceVirtualEnvironmentVMAudioDeviceEnabled                = "enabled"
	mkResourceVirtualEnvironmentVMBIOS                              = "bios"
	mkResourceVirtualEnvironmentVMCDROM                             = "cdrom"
	mkResourceVirtualEnvironmentVMCDROMEnabled                      = "enabled"
	mkResourceVirtualEnvironmentVMCDROMFileID                       = "file_id"
	mkResourceVirtualEnvironmentVMClone                             = "clone"
	mkResourceVirtualEnvironmentVMCloneRetries                      = "retries"
	mkResourceVirtualEnvironmentVMCloneDatastoreID                  = "datastore_id"
	mkResourceVirtualEnvironmentVMCloneNodeName                     = "node_name"
	mkResourceVirtualEnvironmentVMCloneVMID                         = "vm_id"
	mkResourceVirtualEnvironmentVMCloneFull                         = "full"
	mkResourceVirtualEnvironmentVMCPU                               = "cpu"
	mkResourceVirtualEnvironmentVMCPUArchitecture                   = "architecture"
	mkResourceVirtualEnvironmentVMCPUCores                          = "cores"
	mkResourceVirtualEnvironmentVMCPUFlags                          = "flags"
	mkResourceVirtualEnvironmentVMCPUHotplugged                     = "hotplugged"
	mkResourceVirtualEnvironmentVMCPUNUMA                           = "numa"
	mkResourceVirtualEnvironmentVMCPUSockets                        = "sockets"
	mkResourceVirtualEnvironmentVMCPUType                           = "type"
	mkResourceVirtualEnvironmentVMCPUUnits                          = "units"
	mkResourceVirtualEnvironmentVMDescription                       = "description"
	mkResourceVirtualEnvironmentVMDisk                              = "disk"
	mkResourceVirtualEnvironmentVMDiskInterface                     = "interface"
	mkResourceVirtualEnvironmentVMDiskDatastoreID                   = "datastore_id"
	mkResourceVirtualEnvironmentVMDiskFileFormat                    = "file_format"
	mkResourceVirtualEnvironmentVMDiskFileID                        = "file_id"
	mkResourceVirtualEnvironmentVMDiskSize                          = "size"
	mkResourceVirtualEnvironmentVMDiskIOThread                      = "iothread"
	mkResourceVirtualEnvironmentVMDiskSSD                           = "ssd"
	mkResourceVirtualEnvironmentVMDiskDiscard                       = "discard"
	mkResourceVirtualEnvironmentVMDiskSpeed                         = "speed"
	mkResourceVirtualEnvironmentVMDiskSpeedRead                     = "read"
	mkResourceVirtualEnvironmentVMDiskSpeedReadBurstable            = "read_burstable"
	mkResourceVirtualEnvironmentVMDiskSpeedWrite                    = "write"
	mkResourceVirtualEnvironmentVMDiskSpeedWriteBurstable           = "write_burstable"
	mkResourceVirtualEnvironmentVMEFIDisk                           = "efi_disk"
	mkResourceVirtualEnvironmentVMEFIDiskDatastoreID                = "datastore_id"
	mkResourceVirtualEnvironmentVMEFIDiskFileFormat                 = "file_format"
	mkResourceVirtualEnvironmentVMEFIDiskType                       = "type"
	mkResourceVirtualEnvironmentVMEFIDiskPreEnrolledKeys            = "pre_enrolled_keys"
	mkResourceVirtualEnvironmentVMHostPCI                           = "hostpci"
	mkResourceVirtualEnvironmentVMHostPCIDevice                     = "device"
	mkResourceVirtualEnvironmentVMHostPCIDeviceID                   = "id"
	mkResourceVirtualEnvironmentVMHostPCIDeviceMDev                 = "mdev"
	mkResourceVirtualEnvironmentVMHostPCIDevicePCIE                 = "pcie"
	mkResourceVirtualEnvironmentVMHostPCIDeviceROMBAR               = "rombar"
	mkResourceVirtualEnvironmentVMHostPCIDeviceROMFile              = "rom_file"
	mkResourceVirtualEnvironmentVMHostPCIDeviceXVGA                 = "xvga"
	mkResourceVirtualEnvironmentVMInitialization                    = "initialization"
	mkResourceVirtualEnvironmentVMInitializationDatastoreID         = "datastore_id"
	mkResourceVirtualEnvironmentVMInitializationDNS                 = "dns"
	mkResourceVirtualEnvironmentVMInitializationDNSDomain           = "domain"
	mkResourceVirtualEnvironmentVMInitializationDNSServer           = "server"
	mkResourceVirtualEnvironmentVMInitializationIPConfig            = "ip_config"
	mkResourceVirtualEnvironmentVMInitializationIPConfigIPv4        = "ipv4"
	mkResourceVirtualEnvironmentVMInitializationIPConfigIPv4Address = "address"
	mkResourceVirtualEnvironmentVMInitializationIPConfigIPv4Gateway = "gateway"
	mkResourceVirtualEnvironmentVMInitializationIPConfigIPv6        = "ipv6"
	mkResourceVirtualEnvironmentVMInitializationIPConfigIPv6Address = "address"
	mkResourceVirtualEnvironmentVMInitializationIPConfigIPv6Gateway = "gateway"
	mkResourceVirtualEnvironmentVMInitializationType                = "type"
	mkResourceVirtualEnvironmentVMInitializationUserAccount         = "user_account"
	mkResourceVirtualEnvironmentVMInitializationUserAccountKeys     = "keys"
	mkResourceVirtualEnvironmentVMInitializationUserAccountPassword = "password"
	mkResourceVirtualEnvironmentVMInitializationUserAccountUsername = "username"
	mkResourceVirtualEnvironmentVMInitializationUserDataFileID      = "user_data_file_id"
	mkResourceVirtualEnvironmentVMInitializationVendorDataFileID    = "vendor_data_file_id"
	mkResourceVirtualEnvironmentVMInitializationNetworkDataFileID   = "network_data_file_id"
	mkResourceVirtualEnvironmentVMInitializationMetaDataFileID      = "meta_data_file_id"
	mkResourceVirtualEnvironmentVMIPv4Addresses                     = "ipv4_addresses"
	mkResourceVirtualEnvironmentVMIPv6Addresses                     = "ipv6_addresses"
	mkResourceVirtualEnvironmentVMKeyboardLayout                    = "keyboard_layout"
	mkResourceVirtualEnvironmentVMKVMArguments                      = "kvm_arguments"
	mkResourceVirtualEnvironmentVMMachine                           = "machine"
	mkResourceVirtualEnvironmentVMMACAddresses                      = "mac_addresses"
	mkResourceVirtualEnvironmentVMMemory                            = "memory"
	mkResourceVirtualEnvironmentVMMemoryDedicated                   = "dedicated"
	mkResourceVirtualEnvironmentVMMemoryFloating                    = "floating"
	mkResourceVirtualEnvironmentVMMemoryShared                      = "shared"
	mkResourceVirtualEnvironmentVMName                              = "name"
	mkResourceVirtualEnvironmentVMNetworkDevice                     = "network_device"
	mkResourceVirtualEnvironmentVMNetworkDeviceBridge               = "bridge"
	mkResourceVirtualEnvironmentVMNetworkDeviceEnabled              = "enabled"
	mkResourceVirtualEnvironmentVMNetworkDeviceFirewall             = "firewall"
	mkResourceVirtualEnvironmentVMNetworkDeviceMACAddress           = "mac_address"
	mkResourceVirtualEnvironmentVMNetworkDeviceModel                = "model"
	mkResourceVirtualEnvironmentVMNetworkDeviceRateLimit            = "rate_limit"
	mkResourceVirtualEnvironmentVMNetworkDeviceVLANID               = "vlan_id"
	mkResourceVirtualEnvironmentVMNetworkDeviceMTU                  = "mtu"
	mkResourceVirtualEnvironmentVMNetworkInterfaceNames             = "network_interface_names"
	mkResourceVirtualEnvironmentVMNodeName                          = "node_name"
	mkResourceVirtualEnvironmentVMOperatingSystem                   = "operating_system"
	mkResourceVirtualEnvironmentVMOperatingSystemType               = "type"
	mkResourceVirtualEnvironmentVMPoolID                            = "pool_id"
	mkResourceVirtualEnvironmentVMSerialDevice                      = "serial_device"
	mkResourceVirtualEnvironmentVMSerialDeviceDevice                = "device"
	mkResourceVirtualEnvironmentVMStarted                           = "started"
	mkResourceVirtualEnvironmentVMTabletDevice                      = "tablet_device"
	mkResourceVirtualEnvironmentVMTags                              = "tags"
	mkResourceVirtualEnvironmentVMTemplate                          = "template"
	mkResourceVirtualEnvironmentVMTimeoutClone                      = "timeout_clone"
	mkResourceVirtualEnvironmentVMTimeoutMoveDisk                   = "timeout_move_disk"
	mkResourceVirtualEnvironmentVMTimeoutReboot                     = "timeout_reboot"
	mkResourceVirtualEnvironmentVMTimeoutShutdownVM                 = "timeout_shutdown_vm"
	mkResourceVirtualEnvironmentVMTimeoutStartVM                    = "timeout_start_vm"
	mkResourceVirtualEnvironmentVMTimeoutStopVM                     = "timeout_stop_vm"
	mkResourceVirtualEnvironmentVMVGA                               = "vga"
	mkResourceVirtualEnvironmentVMVGAEnabled                        = "enabled"
	mkResourceVirtualEnvironmentVMVGAMemory                         = "memory"
	mkResourceVirtualEnvironmentVMVGAType                           = "type"
	mkResourceVirtualEnvironmentVMVMID                              = "vm_id"
	mkResourceVirtualEnvironmentVMSCSIHardware                      = "scsi_hardware"

	vmCreateTimeoutSeconds = 10
)

// VM returns a resource that manages VMs.
func VM() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			mkResourceVirtualEnvironmentVMRebootAfterCreation: {
				Type:        schema.TypeBool,
				Description: "Whether to reboot vm after creation",
				Optional:    true,
				Default:     dvResourceVirtualEnvironmentVMRebootAfterCreation,
			},
			mkResourceVirtualEnvironmentVMOnBoot: {
				Type:        schema.TypeBool,
				Description: "Start VM on Node boot",
				Optional:    true,
				Default:     dvResourceVirtualEnvironmentVMOnBoot,
			},
			mkResourceVirtualEnvironmentVMBootOrder: {
				Type:        schema.TypeList,
				Description: "The guest will attempt to boot from devices in the order they appear here",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				DefaultFunc: func() (interface{}, error) {
					return []interface{}{}, nil
				},
			},
			mkResourceVirtualEnvironmentVMACPI: {
				Type:        schema.TypeBool,
				Description: "Whether to enable ACPI",
				Optional:    true,
				Default:     dvResourceVirtualEnvironmentVMACPI,
			},
			mkResourceVirtualEnvironmentVMAgent: {
				Type:        schema.TypeList,
				Description: "The QEMU agent configuration",
				Optional:    true,
				DefaultFunc: func() (interface{}, error) {
					return []interface{}{
						map[string]interface{}{
							mkResourceVirtualEnvironmentVMAgentEnabled: dvResourceVirtualEnvironmentVMAgentEnabled,
							mkResourceVirtualEnvironmentVMAgentTimeout: dvResourceVirtualEnvironmentVMAgentTimeout,
							mkResourceVirtualEnvironmentVMAgentTrim:    dvResourceVirtualEnvironmentVMAgentTrim,
							mkResourceVirtualEnvironmentVMAgentType:    dvResourceVirtualEnvironmentVMAgentType,
						},
					}, nil
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mkResourceVirtualEnvironmentVMAgentEnabled: {
							Type:        schema.TypeBool,
							Description: "Whether to enable the QEMU agent",
							Optional:    true,
							Default:     dvResourceVirtualEnvironmentVMAgentEnabled,
						},
						mkResourceVirtualEnvironmentVMAgentTimeout: {
							Type:             schema.TypeString,
							Description:      "The maximum amount of time to wait for data from the QEMU agent to become available",
							Optional:         true,
							Default:          dvResourceVirtualEnvironmentVMAgentTimeout,
							ValidateDiagFunc: getTimeoutValidator(),
						},
						mkResourceVirtualEnvironmentVMAgentTrim: {
							Type:        schema.TypeBool,
							Description: "Whether to enable the FSTRIM feature in the QEMU agent",
							Optional:    true,
							Default:     dvResourceVirtualEnvironmentVMAgentTrim,
						},
						mkResourceVirtualEnvironmentVMAgentType: {
							Type:             schema.TypeString,
							Description:      "The QEMU agent interface type",
							Optional:         true,
							Default:          dvResourceVirtualEnvironmentVMAgentType,
							ValidateDiagFunc: getQEMUAgentTypeValidator(),
						},
					},
				},
				MaxItems: 1,
				MinItems: 0,
			},
			mkResourceVirtualEnvironmentVMKVMArguments: {
				Type:        schema.TypeString,
				Description: "The args implementation",
				Optional:    true,
				Default:     dvResourceVirtualEnvironmentVMKVMArguments,
			},
			mkResourceVirtualEnvironmentVMAudioDevice: {
				Type:        schema.TypeList,
				Description: "The audio devices",
				Optional:    true,
				DefaultFunc: func() (interface{}, error) {
					return []interface{}{}, nil
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mkResourceVirtualEnvironmentVMAudioDeviceDevice: {
							Type:             schema.TypeString,
							Description:      "The device",
							Optional:         true,
							Default:          dvResourceVirtualEnvironmentVMAudioDeviceDevice,
							ValidateDiagFunc: vmGetAudioDeviceValidator(),
						},
						mkResourceVirtualEnvironmentVMAudioDeviceDriver: {
							Type:             schema.TypeString,
							Description:      "The driver",
							Optional:         true,
							Default:          dvResourceVirtualEnvironmentVMAudioDeviceDriver,
							ValidateDiagFunc: vmGetAudioDriverValidator(),
						},
						mkResourceVirtualEnvironmentVMAudioDeviceEnabled: {
							Type:        schema.TypeBool,
							Description: "Whether to enable the audio device",
							Optional:    true,
							Default:     dvResourceVirtualEnvironmentVMAudioDeviceEnabled,
						},
					},
				},
				MaxItems: maxResourceVirtualEnvironmentVMAudioDevices,
				MinItems: 0,
			},
			mkResourceVirtualEnvironmentVMBIOS: {
				Type:             schema.TypeString,
				Description:      "The BIOS implementation",
				Optional:         true,
				Default:          dvResourceVirtualEnvironmentVMBIOS,
				ValidateDiagFunc: getBIOSValidator(),
			},
			mkResourceVirtualEnvironmentVMCDROM: {
				Type:        schema.TypeList,
				Description: "The CDROM drive",
				Optional:    true,
				DefaultFunc: func() (interface{}, error) {
					return []interface{}{
						map[string]interface{}{
							mkResourceVirtualEnvironmentVMCDROMEnabled: dvResourceVirtualEnvironmentVMCDROMEnabled,
							mkResourceVirtualEnvironmentVMCDROMFileID:  dvResourceVirtualEnvironmentVMCDROMFileID,
						},
					}, nil
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mkResourceVirtualEnvironmentVMCDROMEnabled: {
							Type:        schema.TypeBool,
							Description: "Whether to enable the CDROM drive",
							Optional:    true,
							Default:     dvResourceVirtualEnvironmentVMCDROMEnabled,
						},
						mkResourceVirtualEnvironmentVMCDROMFileID: {
							Type:             schema.TypeString,
							Description:      "The file id",
							Optional:         true,
							Default:          dvResourceVirtualEnvironmentVMCDROMFileID,
							ValidateDiagFunc: getFileIDValidator(),
						},
					},
				},
				MaxItems: 1,
				MinItems: 0,
			},
			mkResourceVirtualEnvironmentVMClone: {
				Type:        schema.TypeList,
				Description: "The cloning configuration",
				Optional:    true,
				DefaultFunc: func() (interface{}, error) {
					return []interface{}{}, nil
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mkResourceVirtualEnvironmentVMCloneRetries: {
							Type:        schema.TypeInt,
							Description: "The number of Retries to create a clone",
							Optional:    true,
							ForceNew:    true,
							Default:     dvResourceVirtualEnvironmentVMCloneRetries,
						},
						mkResourceVirtualEnvironmentVMCloneDatastoreID: {
							Type:        schema.TypeString,
							Description: "The ID of the target datastore",
							Optional:    true,
							ForceNew:    true,
							Default:     dvResourceVirtualEnvironmentVMCloneDatastoreID,
						},
						mkResourceVirtualEnvironmentVMCloneNodeName: {
							Type:        schema.TypeString,
							Description: "The name of the source node",
							Optional:    true,
							ForceNew:    true,
							Default:     dvResourceVirtualEnvironmentVMCloneNodeName,
						},
						mkResourceVirtualEnvironmentVMCloneVMID: {
							Type:             schema.TypeInt,
							Description:      "The ID of the source VM",
							Required:         true,
							ForceNew:         true,
							ValidateDiagFunc: validator.VMID(),
						},
						mkResourceVirtualEnvironmentVMCloneFull: {
							Type:        schema.TypeBool,
							Description: "The Clone Type, create a Full Clone (true) or a linked Clone (false)",
							Optional:    true,
							ForceNew:    true,
							Default:     dvResourceVirtualEnvironmentVMCloneFull,
						},
					},
				},
				MaxItems: 1,
				MinItems: 0,
			},
			mkResourceVirtualEnvironmentVMCPU: {
				Type:        schema.TypeList,
				Description: "The CPU allocation",
				Optional:    true,
				DefaultFunc: func() (interface{}, error) {
					return []interface{}{
						map[string]interface{}{
							mkResourceVirtualEnvironmentVMCPUArchitecture: dvResourceVirtualEnvironmentVMCPUArchitecture,
							mkResourceVirtualEnvironmentVMCPUCores:        dvResourceVirtualEnvironmentVMCPUCores,
							mkResourceVirtualEnvironmentVMCPUFlags:        []interface{}{},
							mkResourceVirtualEnvironmentVMCPUNUMA:         dvResourceVirtualEnvironmentVMCPUNUMA,
							mkResourceVirtualEnvironmentVMCPUHotplugged:   dvResourceVirtualEnvironmentVMCPUHotplugged,
							mkResourceVirtualEnvironmentVMCPUSockets:      dvResourceVirtualEnvironmentVMCPUSockets,
							mkResourceVirtualEnvironmentVMCPUType:         dvResourceVirtualEnvironmentVMCPUType,
							mkResourceVirtualEnvironmentVMCPUUnits:        dvResourceVirtualEnvironmentVMCPUUnits,
						},
					}, nil
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mkResourceVirtualEnvironmentVMCPUArchitecture: {
							Type:             schema.TypeString,
							Description:      "The CPU architecture",
							Optional:         true,
							Default:          dvResourceVirtualEnvironmentVMCPUArchitecture,
							ValidateDiagFunc: vmGetCPUArchitectureValidator(),
						},
						mkResourceVirtualEnvironmentVMCPUCores: {
							Type:             schema.TypeInt,
							Description:      "The number of CPU cores",
							Optional:         true,
							Default:          dvResourceVirtualEnvironmentVMCPUCores,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(1, 2304)),
						},
						mkResourceVirtualEnvironmentVMCPUFlags: {
							Type:        schema.TypeList,
							Description: "The CPU flags",
							Optional:    true,
							DefaultFunc: func() (interface{}, error) {
								return []interface{}{}, nil
							},
							Elem: &schema.Schema{Type: schema.TypeString},
						},
						mkResourceVirtualEnvironmentVMCPUHotplugged: {
							Type:             schema.TypeInt,
							Description:      "The number of hotplugged vCPUs",
							Optional:         true,
							Default:          dvResourceVirtualEnvironmentVMCPUHotplugged,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(0, 2304)),
						},
						mkResourceVirtualEnvironmentVMCPUNUMA: {
							Type:        schema.TypeBool,
							Description: "Enable/disable NUMA.",
							Optional:    true,
							Default:     dvResourceVirtualEnvironmentVMCPUNUMA,
						},
						mkResourceVirtualEnvironmentVMCPUSockets: {
							Type:             schema.TypeInt,
							Description:      "The number of CPU sockets",
							Optional:         true,
							Default:          dvResourceVirtualEnvironmentVMCPUSockets,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(1, 16)),
						},
						mkResourceVirtualEnvironmentVMCPUType: {
							Type:             schema.TypeString,
							Description:      "The emulated CPU type",
							Optional:         true,
							Default:          dvResourceVirtualEnvironmentVMCPUType,
							ValidateDiagFunc: getCPUTypeValidator(),
						},
						mkResourceVirtualEnvironmentVMCPUUnits: {
							Type:        schema.TypeInt,
							Description: "The CPU units",
							Optional:    true,
							Default:     dvResourceVirtualEnvironmentVMCPUUnits,
							ValidateDiagFunc: validation.ToDiagFunc(
								validation.IntBetween(2, 262144),
							),
						},
					},
				},
				MaxItems: 1,
				MinItems: 0,
			},
			mkResourceVirtualEnvironmentVMDescription: {
				Type:        schema.TypeString,
				Description: "The description",
				Optional:    true,
				Default:     dvResourceVirtualEnvironmentVMDescription,
			},
			mkResourceVirtualEnvironmentVMDisk: {
				Type:        schema.TypeList,
				Description: "The disk devices",
				Optional:    true,
				ForceNew:    true,
				DefaultFunc: func() (interface{}, error) {
					return []interface{}{
						map[string]interface{}{
							mkResourceVirtualEnvironmentVMDiskDatastoreID: dvResourceVirtualEnvironmentVMDiskDatastoreID,
							mkResourceVirtualEnvironmentVMDiskFileID:      dvResourceVirtualEnvironmentVMDiskFileID,
							mkResourceVirtualEnvironmentVMDiskInterface:   dvResourceVirtualEnvironmentVMDiskInterface,
							mkResourceVirtualEnvironmentVMDiskSize:        dvResourceVirtualEnvironmentVMDiskSize,
							mkResourceVirtualEnvironmentVMDiskIOThread:    dvResourceVirtualEnvironmentVMDiskIOThread,
							mkResourceVirtualEnvironmentVMDiskSSD:         dvResourceVirtualEnvironmentVMDiskSSD,
							mkResourceVirtualEnvironmentVMDiskDiscard:     dvResourceVirtualEnvironmentVMDiskDiscard,
						},
					}, nil
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mkResourceVirtualEnvironmentVMDiskInterface: {
							Type:        schema.TypeString,
							Description: "The datastore name",
							Required:    true,
						},
						mkResourceVirtualEnvironmentVMDiskDatastoreID: {
							Type:        schema.TypeString,
							Description: "The datastore id",
							Optional:    true,
							Default:     dvResourceVirtualEnvironmentVMDiskDatastoreID,
						},
						mkResourceVirtualEnvironmentVMDiskFileFormat: {
							Type:             schema.TypeString,
							Description:      "The file format",
							Optional:         true,
							ForceNew:         true,
							Computed:         true,
							ValidateDiagFunc: getFileFormatValidator(),
						},
						mkResourceVirtualEnvironmentVMDiskFileID: {
							Type:             schema.TypeString,
							Description:      "The file id for a disk image",
							Optional:         true,
							ForceNew:         true,
							Default:          dvResourceVirtualEnvironmentVMDiskFileID,
							ValidateDiagFunc: getFileIDValidator(),
						},
						mkResourceVirtualEnvironmentVMDiskSize: {
							Type:             schema.TypeInt,
							Description:      "The disk size in gigabytes",
							Optional:         true,
							Default:          dvResourceVirtualEnvironmentVMDiskSize,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(1)),
						},
						mkResourceVirtualEnvironmentVMDiskIOThread: {
							Type:        schema.TypeBool,
							Description: "Whether to use iothreads for this disk drive",
							Optional:    true,
							Default:     dvResourceVirtualEnvironmentVMDiskIOThread,
						},
						mkResourceVirtualEnvironmentVMDiskSSD: {
							Type:        schema.TypeBool,
							Description: "Whether to use ssd for this disk drive",
							Optional:    true,
							Default:     dvResourceVirtualEnvironmentVMDiskSSD,
						},
						mkResourceVirtualEnvironmentVMDiskDiscard: {
							Type:        schema.TypeString,
							Description: "Whether to pass discard/trim requests to the underlying storage.",
							Optional:    true,
							Default:     dvResourceVirtualEnvironmentVMDiskDiscard,
						},
						mkResourceVirtualEnvironmentVMDiskSpeed: {
							Type:        schema.TypeList,
							Description: "The speed limits",
							Optional:    true,
							DefaultFunc: func() (interface{}, error) {
								return []interface{}{
									map[string]interface{}{
										mkResourceVirtualEnvironmentVMDiskSpeedRead:           dvResourceVirtualEnvironmentVMDiskSpeedRead,
										mkResourceVirtualEnvironmentVMDiskSpeedReadBurstable:  dvResourceVirtualEnvironmentVMDiskSpeedReadBurstable,
										mkResourceVirtualEnvironmentVMDiskSpeedWrite:          dvResourceVirtualEnvironmentVMDiskSpeedWrite,
										mkResourceVirtualEnvironmentVMDiskSpeedWriteBurstable: dvResourceVirtualEnvironmentVMDiskSpeedWriteBurstable,
									},
								}, nil
							},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									mkResourceVirtualEnvironmentVMDiskSpeedRead: {
										Type:        schema.TypeInt,
										Description: "The maximum read speed in megabytes per second",
										Optional:    true,
										Default:     dvResourceVirtualEnvironmentVMDiskSpeedRead,
									},
									mkResourceVirtualEnvironmentVMDiskSpeedReadBurstable: {
										Type:        schema.TypeInt,
										Description: "The maximum burstable read speed in megabytes per second",
										Optional:    true,
										Default:     dvResourceVirtualEnvironmentVMDiskSpeedReadBurstable,
									},
									mkResourceVirtualEnvironmentVMDiskSpeedWrite: {
										Type:        schema.TypeInt,
										Description: "The maximum write speed in megabytes per second",
										Optional:    true,
										Default:     dvResourceVirtualEnvironmentVMDiskSpeedWrite,
									},
									mkResourceVirtualEnvironmentVMDiskSpeedWriteBurstable: {
										Type:        schema.TypeInt,
										Description: "The maximum burstable write speed in megabytes per second",
										Optional:    true,
										Default:     dvResourceVirtualEnvironmentVMDiskSpeedWriteBurstable,
									},
								},
							},
							MaxItems: 1,
							MinItems: 0,
						},
					},
				},
				MaxItems: 14,
				MinItems: 0,
			},
			mkResourceVirtualEnvironmentVMEFIDisk: {
				Type:        schema.TypeList,
				Description: "The efidisk device",
				Optional:    true,
				ForceNew:    true,
				DefaultFunc: func() (interface{}, error) {
					return []interface{}{
						map[string]interface{}{
							mkResourceVirtualEnvironmentVMEFIDiskDatastoreID:     dvResourceVirtualEnvironmentVMEFIDiskDatastoreID,
							mkResourceVirtualEnvironmentVMEFIDiskType:            dvResourceVirtualEnvironmentVMEFIDiskType,
							mkResourceVirtualEnvironmentVMEFIDiskFileFormat:      dvResourceVirtualEnvironmentVMEFIDiskFileFormat,
							mkResourceVirtualEnvironmentVMEFIDiskPreEnrolledKeys: dvResourceVirtualEnvironmentVMEFIDiskPreEnrolledKeys,
						},
					}, nil
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mkResourceVirtualEnvironmentVMEFIDiskDatastoreID: {
							Type:        schema.TypeString,
							Description: "The datastore id",
							Optional:    true,
							Default:     dvResourceVirtualEnvironmentVMEFIDiskDatastoreID,
						},
						mkResourceVirtualEnvironmentVMEFIDiskFileFormat: {
							Type:             schema.TypeString,
							Description:      "The file format",
							Optional:         true,
							ForceNew:         true,
							Computed:         true,
							ValidateDiagFunc: getFileFormatValidator(),
						},
						mkResourceVirtualEnvironmentVMEFIDiskType: {
							Type:        schema.TypeString,
							Description: "Size and type of the OVMF EFI disk",
							Optional:    true,
							ForceNew:    true,
							Default:     dvResourceVirtualEnvironmentVMEFIDiskType,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
								"2m",
								"4m",
							}, true)),
						},
						mkResourceVirtualEnvironmentVMEFIDiskPreEnrolledKeys: {
							Type: schema.TypeBool,
							Description: "Use an EFI vars template with distribution-specific and Microsoft Standard " +
								"keys enrolled, if used with efi type=`4m`.",
							Optional: true,
							ForceNew: true,
							Default:  dvResourceVirtualEnvironmentVMEFIDiskPreEnrolledKeys,
						},
					},
				},
				MaxItems: 1,
				MinItems: 0,
			},
			mkResourceVirtualEnvironmentVMInitialization: {
				Type:        schema.TypeList,
				Description: "The cloud-init configuration",
				Optional:    true,
				DefaultFunc: func() (interface{}, error) {
					return []interface{}{}, nil
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mkResourceVirtualEnvironmentVMInitializationDatastoreID: {
							Type:        schema.TypeString,
							Description: "The datastore id",
							Optional:    true,
							ForceNew:    true,
							Default:     dvResourceVirtualEnvironmentVMInitializationDatastoreID,
						},
						mkResourceVirtualEnvironmentVMInitializationDNS: {
							Type:        schema.TypeList,
							Description: "The DNS configuration",
							Optional:    true,
							DefaultFunc: func() (interface{}, error) {
								return []interface{}{}, nil
							},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									mkResourceVirtualEnvironmentVMInitializationDNSDomain: {
										Type:        schema.TypeString,
										Description: "The DNS search domain",
										Optional:    true,
										Default:     dvResourceVirtualEnvironmentVMInitializationDNSDomain,
									},
									mkResourceVirtualEnvironmentVMInitializationDNSServer: {
										Type:        schema.TypeString,
										Description: "The DNS server",
										Optional:    true,
										Default:     dvResourceVirtualEnvironmentVMInitializationDNSServer,
									},
								},
							},
							MaxItems: 1,
							MinItems: 0,
						},
						mkResourceVirtualEnvironmentVMInitializationIPConfig: {
							Type:        schema.TypeList,
							Description: "The IP configuration",
							Optional:    true,
							DefaultFunc: func() (interface{}, error) {
								return []interface{}{}, nil
							},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									mkResourceVirtualEnvironmentVMInitializationIPConfigIPv4: {
										Type:        schema.TypeList,
										Description: "The IPv4 configuration",
										Optional:    true,
										DefaultFunc: func() (interface{}, error) {
											return []interface{}{}, nil
										},
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												mkResourceVirtualEnvironmentVMInitializationIPConfigIPv4Address: {
													Type:        schema.TypeString,
													Description: "The IPv4 address",
													Optional:    true,
													Default:     dvResourceVirtualEnvironmentVMInitializationIPConfigIPv4Address,
												},
												mkResourceVirtualEnvironmentVMInitializationIPConfigIPv4Gateway: {
													Type:        schema.TypeString,
													Description: "The IPv4 gateway",
													Optional:    true,
													Default:     dvResourceVirtualEnvironmentVMInitializationIPConfigIPv4Gateway,
												},
											},
										},
										MaxItems: 1,
										MinItems: 0,
									},
									mkResourceVirtualEnvironmentVMInitializationIPConfigIPv6: {
										Type:        schema.TypeList,
										Description: "The IPv6 configuration",
										Optional:    true,
										DefaultFunc: func() (interface{}, error) {
											return []interface{}{}, nil
										},
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												mkResourceVirtualEnvironmentVMInitializationIPConfigIPv6Address: {
													Type:        schema.TypeString,
													Description: "The IPv6 address",
													Optional:    true,
													Default:     dvResourceVirtualEnvironmentVMInitializationIPConfigIPv6Address,
												},
												mkResourceVirtualEnvironmentVMInitializationIPConfigIPv6Gateway: {
													Type:        schema.TypeString,
													Description: "The IPv6 gateway",
													Optional:    true,
													Default:     dvResourceVirtualEnvironmentVMInitializationIPConfigIPv6Gateway,
												},
											},
										},
										MaxItems: 1,
										MinItems: 0,
									},
								},
							},
							MaxItems: 8,
							MinItems: 0,
						},
						mkResourceVirtualEnvironmentVMInitializationUserAccount: {
							Type:        schema.TypeList,
							Description: "The user account configuration",
							Optional:    true,
							ForceNew:    true,
							DefaultFunc: func() (interface{}, error) {
								return []interface{}{}, nil
							},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									mkResourceVirtualEnvironmentVMInitializationUserAccountKeys: {
										Type:        schema.TypeList,
										Description: "The SSH keys",
										Optional:    true,
										ForceNew:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									mkResourceVirtualEnvironmentVMInitializationUserAccountPassword: {
										Type:        schema.TypeString,
										Description: "The SSH password",
										Optional:    true,
										ForceNew:    true,
										Sensitive:   true,
										Default:     dvResourceVirtualEnvironmentVMInitializationUserAccountPassword,
										DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
											return len(old) > 0 &&
												strings.ReplaceAll(old, "*", "") == ""
										},
									},
									mkResourceVirtualEnvironmentVMInitializationUserAccountUsername: {
										Type:        schema.TypeString,
										Description: "The SSH username",
										Optional:    true,
										ForceNew:    true,
									},
								},
							},
							MaxItems: 1,
							MinItems: 0,
						},
						mkResourceVirtualEnvironmentVMInitializationUserDataFileID: {
							Type:             schema.TypeString,
							Description:      "The ID of a file containing custom user data",
							Optional:         true,
							ForceNew:         true,
							Default:          dvResourceVirtualEnvironmentVMInitializationUserDataFileID,
							ValidateDiagFunc: getFileIDValidator(),
						},
						mkResourceVirtualEnvironmentVMInitializationVendorDataFileID: {
							Type:             schema.TypeString,
							Description:      "The ID of a file containing vendor data",
							Optional:         true,
							ForceNew:         true,
							Default:          dvResourceVirtualEnvironmentVMInitializationVendorDataFileID,
							ValidateDiagFunc: getFileIDValidator(),
						},
						mkResourceVirtualEnvironmentVMInitializationNetworkDataFileID: {
							Type:             schema.TypeString,
							Description:      "The ID of a file containing network config",
							Optional:         true,
							ForceNew:         true,
							Default:          dvResourceVirtualEnvironmentVMInitializationNetworkDataFileID,
							ValidateDiagFunc: getFileIDValidator(),
						},
						mkResourceVirtualEnvironmentVMInitializationMetaDataFileID: {
							Type:             schema.TypeString,
							Description:      "The ID of a file containing meta data config",
							Optional:         true,
							ForceNew:         true,
							Default:          dvResourceVirtualEnvironmentVMInitializationMetaDataFileID,
							ValidateDiagFunc: getFileIDValidator(),
						},
						mkResourceVirtualEnvironmentVMInitializationType: {
							Type:             schema.TypeString,
							Description:      "The cloud-init configuration format",
							Optional:         true,
							ForceNew:         true,
							Default:          dvResourceVirtualEnvironmentVMInitializationType,
							ValidateDiagFunc: getCloudInitTypeValidator(),
						},
					},
				},
				MaxItems: 1,
				MinItems: 0,
			},
			mkResourceVirtualEnvironmentVMIPv4Addresses: {
				Type:        schema.TypeList,
				Description: "The IPv4 addresses published by the QEMU agent",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeList,
					Elem: &schema.Schema{Type: schema.TypeString},
				},
			},
			mkResourceVirtualEnvironmentVMIPv6Addresses: {
				Type:        schema.TypeList,
				Description: "The IPv6 addresses published by the QEMU agent",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeList,
					Elem: &schema.Schema{Type: schema.TypeString},
				},
			},
			mkResourceVirtualEnvironmentVMHostPCI: {
				Type:        schema.TypeList,
				Description: "The Host PCI devices mapped to the VM",
				Optional:    true,
				ForceNew:    true,
				DefaultFunc: func() (interface{}, error) {
					return []interface{}{}, nil
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mkResourceVirtualEnvironmentVMHostPCIDevice: {
							Type:        schema.TypeString,
							Description: "The PCI device name for Proxmox, in form of 'hostpciX' where X is a sequential number from 0 to 3",
							Required:    true,
						},
						mkResourceVirtualEnvironmentVMHostPCIDeviceID: {
							Type:        schema.TypeString,
							Description: "The PCI ID of the device, for example 0000:00:1f.0 (or 0000:00:1f.0;0000:00:1f.1 for multiple device functions, or 0000:00:1f for all functions)",
							Required:    true,
						},
						mkResourceVirtualEnvironmentVMHostPCIDeviceMDev: {
							Type:        schema.TypeString,
							Description: "The the mediated device to use",
							Optional:    true,
						},
						mkResourceVirtualEnvironmentVMHostPCIDevicePCIE: {
							Type:        schema.TypeBool,
							Description: "Tells Proxmox VE to use a PCIe or PCI port. Some guests/device combination require PCIe rather than PCI. PCIe is only available for q35 machine types.",
							Optional:    true,
						},
						mkResourceVirtualEnvironmentVMHostPCIDeviceROMBAR: {
							Type:        schema.TypeBool,
							Description: "Makes the firmware ROM visible for the guest. Default is true",
							Optional:    true,
						},
						mkResourceVirtualEnvironmentVMHostPCIDeviceROMFile: {
							Type:        schema.TypeString,
							Description: "A path to a ROM file for the device to use. This is a relative path under /usr/share/kvm/",
							Optional:    true,
						},
						mkResourceVirtualEnvironmentVMHostPCIDeviceXVGA: {
							Type:        schema.TypeBool,
							Description: "Marks the PCI(e) device as the primary GPU of the VM. With this enabled the vga configuration argument will be ignored.",
							Optional:    true,
						},
					},
				},
			},
			mkResourceVirtualEnvironmentVMKeyboardLayout: {
				Type:             schema.TypeString,
				Description:      "The keyboard layout",
				Optional:         true,
				Default:          dvResourceVirtualEnvironmentVMKeyboardLayout,
				ValidateDiagFunc: getKeyboardLayoutValidator(),
			},
			mkResourceVirtualEnvironmentVMMachine: {
				Type:        schema.TypeString,
				Description: "The VM machine type, either default i440fx or q35",
				Optional:    true,
				Default:     dvResourceVirtualEnvironmentVMMachineType,
			},
			mkResourceVirtualEnvironmentVMMACAddresses: {
				Type:        schema.TypeList,
				Description: "The MAC addresses for the network interfaces",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			mkResourceVirtualEnvironmentVMMemory: {
				Type:        schema.TypeList,
				Description: "The memory allocation",
				Optional:    true,
				DefaultFunc: func() (interface{}, error) {
					return []interface{}{
						map[string]interface{}{
							mkResourceVirtualEnvironmentVMMemoryDedicated: dvResourceVirtualEnvironmentVMMemoryDedicated,
							mkResourceVirtualEnvironmentVMMemoryFloating:  dvResourceVirtualEnvironmentVMMemoryFloating,
							mkResourceVirtualEnvironmentVMMemoryShared:    dvResourceVirtualEnvironmentVMMemoryShared,
						},
					}, nil
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mkResourceVirtualEnvironmentVMMemoryDedicated: {
							Type:        schema.TypeInt,
							Description: "The dedicated memory in megabytes",
							Optional:    true,
							Default:     dvResourceVirtualEnvironmentVMMemoryDedicated,
							ValidateDiagFunc: validation.ToDiagFunc(
								validation.IntBetween(64, 268435456),
							),
						},
						mkResourceVirtualEnvironmentVMMemoryFloating: {
							Type:        schema.TypeInt,
							Description: "The floating memory in megabytes (balloon)",
							Optional:    true,
							Default:     dvResourceVirtualEnvironmentVMMemoryFloating,
							ValidateDiagFunc: validation.ToDiagFunc(
								validation.IntBetween(0, 268435456),
							),
						},
						mkResourceVirtualEnvironmentVMMemoryShared: {
							Type:        schema.TypeInt,
							Description: "The shared memory in megabytes",
							Optional:    true,
							Default:     dvResourceVirtualEnvironmentVMMemoryShared,
							ValidateDiagFunc: validation.ToDiagFunc(
								validation.IntBetween(0, 268435456),
							),
						},
					},
				},
				MaxItems: 1,
				MinItems: 0,
			},
			mkResourceVirtualEnvironmentVMName: {
				Type:        schema.TypeString,
				Description: "The name",
				Optional:    true,
				Default:     dvResourceVirtualEnvironmentVMName,
			},
			mkResourceVirtualEnvironmentVMNetworkDevice: {
				Type:        schema.TypeList,
				Description: "The network devices",
				Optional:    true,
				DefaultFunc: func() (interface{}, error) {
					return make([]interface{}, 1), nil
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mkResourceVirtualEnvironmentVMNetworkDeviceBridge: {
							Type:        schema.TypeString,
							Description: "The bridge",
							Optional:    true,
							Default:     dvResourceVirtualEnvironmentVMNetworkDeviceBridge,
						},
						mkResourceVirtualEnvironmentVMNetworkDeviceEnabled: {
							Type:        schema.TypeBool,
							Description: "Whether to enable the network device",
							Optional:    true,
							Default:     dvResourceVirtualEnvironmentVMNetworkDeviceEnabled,
						},
						mkResourceVirtualEnvironmentVMNetworkDeviceFirewall: {
							Type:        schema.TypeBool,
							Description: "Whether this interface's firewall rules should be used",
							Optional:    true,
							Default:     dvResourceVirtualEnvironmentVMNetworkDeviceFirewall,
						},
						mkResourceVirtualEnvironmentVMNetworkDeviceMACAddress: {
							Type:             schema.TypeString,
							Description:      "The MAC address",
							Optional:         true,
							Computed:         true,
							ValidateDiagFunc: getMACAddressValidator(),
						},
						mkResourceVirtualEnvironmentVMNetworkDeviceModel: {
							Type:             schema.TypeString,
							Description:      "The model",
							Optional:         true,
							Default:          dvResourceVirtualEnvironmentVMNetworkDeviceModel,
							ValidateDiagFunc: getNetworkDeviceModelValidator(),
						},
						mkResourceVirtualEnvironmentVMNetworkDeviceRateLimit: {
							Type:        schema.TypeFloat,
							Description: "The rate limit in megabytes per second",
							Optional:    true,
							Default:     dvResourceVirtualEnvironmentVMNetworkDeviceRateLimit,
						},
						mkResourceVirtualEnvironmentVMNetworkDeviceVLANID: {
							Type:        schema.TypeInt,
							Description: "The VLAN identifier",
							Optional:    true,
							Default:     dvResourceVirtualEnvironmentVMNetworkDeviceVLANID,
						},
						mkResourceVirtualEnvironmentVMNetworkDeviceMTU: {
							Type:        schema.TypeInt,
							Description: "Maximum transmission unit (MTU)",
							Optional:    true,
							Default:     dvResourceVirtualEnvironmentVMNetworkDeviceMTU,
						},
					},
				},
				MaxItems: maxResourceVirtualEnvironmentVMNetworkDevices,
				MinItems: 0,
			},
			mkResourceVirtualEnvironmentVMNetworkInterfaceNames: {
				Type:        schema.TypeList,
				Description: "The network interface names published by the QEMU agent",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			mkResourceVirtualEnvironmentVMNodeName: {
				Type:        schema.TypeString,
				Description: "The node name",
				Required:    true,
				ForceNew:    true,
			},
			mkResourceVirtualEnvironmentVMOperatingSystem: {
				Type:        schema.TypeList,
				Description: "The operating system configuration",
				Optional:    true,
				DefaultFunc: func() (interface{}, error) {
					return []interface{}{
						map[string]interface{}{
							mkResourceVirtualEnvironmentVMOperatingSystemType: dvResourceVirtualEnvironmentVMOperatingSystemType,
						},
					}, nil
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mkResourceVirtualEnvironmentVMOperatingSystemType: {
							Type:             schema.TypeString,
							Description:      "The type",
							Optional:         true,
							Default:          dvResourceVirtualEnvironmentVMOperatingSystemType,
							ValidateDiagFunc: vmGetOperatingSystemTypeValidator(),
						},
					},
				},
				MaxItems: 1,
				MinItems: 0,
			},
			mkResourceVirtualEnvironmentVMPoolID: {
				Type:        schema.TypeString,
				Description: "The ID of the pool to assign the virtual machine to",
				Optional:    true,
				ForceNew:    true,
				Default:     dvResourceVirtualEnvironmentVMPoolID,
			},
			mkResourceVirtualEnvironmentVMSerialDevice: {
				Type:        schema.TypeList,
				Description: "The serial devices",
				Optional:    true,
				DefaultFunc: func() (interface{}, error) {
					return []interface{}{
						map[string]interface{}{
							mkResourceVirtualEnvironmentVMSerialDeviceDevice: dvResourceVirtualEnvironmentVMSerialDeviceDevice,
						},
					}, nil
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mkResourceVirtualEnvironmentVMSerialDeviceDevice: {
							Type:             schema.TypeString,
							Description:      "The device",
							Optional:         true,
							Default:          dvResourceVirtualEnvironmentVMSerialDeviceDevice,
							ValidateDiagFunc: vmGetSerialDeviceValidator(),
						},
					},
				},
				MaxItems: maxResourceVirtualEnvironmentVMSerialDevices,
				MinItems: 0,
			},
			mkResourceVirtualEnvironmentVMStarted: {
				Type:        schema.TypeBool,
				Description: "Whether to start the virtual machine",
				Optional:    true,
				Default:     dvResourceVirtualEnvironmentVMStarted,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get(mkResourceVirtualEnvironmentVMTemplate).(bool)
				},
			},
			mkResourceVirtualEnvironmentVMTabletDevice: {
				Type:        schema.TypeBool,
				Description: "Whether to enable the USB tablet device",
				Optional:    true,
				Default:     dvResourceVirtualEnvironmentVMTabletDevice,
			},
			mkResourceVirtualEnvironmentVMTags: {
				Type:        schema.TypeList,
				Description: "Tags of the virtual machine. This is only meta information.",
				Optional:    true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				DiffSuppressFunc:      suppressIfListsAreEqualIgnoringOrder,
				DiffSuppressOnRefresh: true,
			},
			mkResourceVirtualEnvironmentVMTemplate: {
				Type:        schema.TypeBool,
				Description: "Whether to create a template",
				Optional:    true,
				ForceNew:    true,
				Default:     dvResourceVirtualEnvironmentVMTemplate,
			},
			mkResourceVirtualEnvironmentVMTimeoutClone: {
				Type:        schema.TypeInt,
				Description: "Clone VM timeout",
				Optional:    true,
				Default:     dvResourceVirtualEnvironmentVMTimeoutClone,
			},
			mkResourceVirtualEnvironmentVMTimeoutMoveDisk: {
				Type:        schema.TypeInt,
				Description: "MoveDisk timeout",
				Optional:    true,
				Default:     dvResourceVirtualEnvironmentVMTimeoutMoveDisk,
			},
			mkResourceVirtualEnvironmentVMTimeoutReboot: {
				Type:        schema.TypeInt,
				Description: "Reboot timeout",
				Optional:    true,
				Default:     dvResourceVirtualEnvironmentVMTimeoutReboot,
			},
			mkResourceVirtualEnvironmentVMTimeoutShutdownVM: {
				Type:        schema.TypeInt,
				Description: "Shutdown timeout",
				Optional:    true,
				Default:     dvResourceVirtualEnvironmentVMTimeoutShutdownVM,
			},
			mkResourceVirtualEnvironmentVMTimeoutStartVM: {
				Type:        schema.TypeInt,
				Description: "Start VM timeout",
				Optional:    true,
				Default:     dvResourceVirtualEnvironmentVMTimeoutStartVM,
			},
			mkResourceVirtualEnvironmentVMTimeoutStopVM: {
				Type:        schema.TypeInt,
				Description: "Stop VM timeout",
				Optional:    true,
				Default:     dvResourceVirtualEnvironmentVMTimeoutStopVM,
			},
			mkResourceVirtualEnvironmentVMVGA: {
				Type:        schema.TypeList,
				Description: "The VGA configuration",
				Optional:    true,
				DefaultFunc: func() (interface{}, error) {
					return []interface{}{
						map[string]interface{}{
							mkResourceVirtualEnvironmentVMVGAEnabled: dvResourceVirtualEnvironmentVMVGAEnabled,
							mkResourceVirtualEnvironmentVMVGAMemory:  dvResourceVirtualEnvironmentVMVGAMemory,
							mkResourceVirtualEnvironmentVMVGAType:    dvResourceVirtualEnvironmentVMVGAType,
						},
					}, nil
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mkResourceVirtualEnvironmentVMVGAEnabled: {
							Type:        schema.TypeBool,
							Description: "Whether to enable the VGA device",
							Optional:    true,
							Default:     dvResourceVirtualEnvironmentVMVGAEnabled,
						},
						mkResourceVirtualEnvironmentVMVGAMemory: {
							Type:             schema.TypeInt,
							Description:      "The VGA memory in megabytes (4-512 MB)",
							Optional:         true,
							Default:          dvResourceVirtualEnvironmentVMVGAMemory,
							ValidateDiagFunc: getVGAMemoryValidator(),
						},
						mkResourceVirtualEnvironmentVMVGAType: {
							Type:             schema.TypeString,
							Description:      "The VGA type",
							Optional:         true,
							Default:          dvResourceVirtualEnvironmentVMVGAType,
							ValidateDiagFunc: getVGATypeValidator(),
						},
					},
				},
				MaxItems: 1,
				MinItems: 0,
			},
			mkResourceVirtualEnvironmentVMVMID: {
				Type:        schema.TypeInt,
				Description: "The VM identifier",
				Optional:    true,
				Computed:    true,
				// "ForceNew: true" handled in CustomizeDiff, making sure VMs with legacy configs with vm_id = -1
				// do not require re-creation.
				ValidateDiagFunc: validator.VMID(),
			},
			mkResourceVirtualEnvironmentVMSCSIHardware: {
				Type:             schema.TypeString,
				Description:      "The SCSI hardware type",
				Optional:         true,
				Default:          dvResourceVirtualEnvironmentVMSCSIHardware,
				ValidateDiagFunc: getSCSIHardwareValidator(),
			},
		},
		CreateContext: vmCreate,
		ReadContext:   vmRead,
		UpdateContext: vmUpdate,
		DeleteContext: vmDelete,
		CustomizeDiff: customdiff.All(
			customdiff.ComputedIf(
				mkResourceVirtualEnvironmentVMIPv4Addresses,
				func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) bool {
					return d.HasChange(mkResourceVirtualEnvironmentVMStarted) ||
						d.HasChange(mkResourceVirtualEnvironmentVMNetworkDevice)
				},
			),
			customdiff.ComputedIf(
				mkResourceVirtualEnvironmentVMIPv6Addresses,
				func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) bool {
					return d.HasChange(mkResourceVirtualEnvironmentVMStarted) ||
						d.HasChange(mkResourceVirtualEnvironmentVMNetworkDevice)
				},
			),
			customdiff.ComputedIf(
				mkResourceVirtualEnvironmentVMNetworkInterfaceNames,
				func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) bool {
					return d.HasChange(mkResourceVirtualEnvironmentVMStarted) ||
						d.HasChange(mkResourceVirtualEnvironmentVMNetworkDevice)
				},
			),
			customdiff.ForceNewIf(
				mkResourceVirtualEnvironmentVMVMID,
				func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) bool {
					newValue := d.Get(mkResourceVirtualEnvironmentVMVMID)

					// 'vm_id' is ForceNew, except when changing 'vm_id' to existing correct id
					// (automatic fix from -1 to actual vm_id must not re-create VM)
					return strconv.Itoa(newValue.(int)) != d.Id()
				},
			),
		),
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, i interface{}) ([]*schema.ResourceData, error) {
				node, id, err := parseImportIDWithNodeName(d.Id())
				if err != nil {
					return nil, err
				}

				d.SetId(id)
				err = d.Set(mkResourceVirtualEnvironmentVMNodeName, node)
				if err != nil {
					return nil, fmt.Errorf("failed setting state during import: %w", err)
				}

				return []*schema.ResourceData{d}, nil
			},
		},
	}
}

func vmCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clone := d.Get(mkResourceVirtualEnvironmentVMClone).([]interface{})

	if len(clone) > 0 {
		return vmCreateClone(ctx, d, m)
	}

	return vmCreateCustom(ctx, d, m)
}

func vmCreateClone(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(proxmoxtf.ProviderConfiguration)

	api, e := config.GetClient()
	if e != nil {
		return diag.FromErr(e)
	}

	clone := d.Get(mkResourceVirtualEnvironmentVMClone).([]interface{})
	cloneBlock := clone[0].(map[string]interface{})
	cloneRetries := cloneBlock[mkResourceVirtualEnvironmentVMCloneRetries].(int)
	cloneDatastoreID := cloneBlock[mkResourceVirtualEnvironmentVMCloneDatastoreID].(string)
	cloneNodeName := cloneBlock[mkResourceVirtualEnvironmentVMCloneNodeName].(string)
	cloneVMID := cloneBlock[mkResourceVirtualEnvironmentVMCloneVMID].(int)
	cloneFull := cloneBlock[mkResourceVirtualEnvironmentVMCloneFull].(bool)

	description := d.Get(mkResourceVirtualEnvironmentVMDescription).(string)
	name := d.Get(mkResourceVirtualEnvironmentVMName).(string)
	tags := d.Get(mkResourceVirtualEnvironmentVMTags).([]interface{})
	nodeName := d.Get(mkResourceVirtualEnvironmentVMNodeName).(string)
	poolID := d.Get(mkResourceVirtualEnvironmentVMPoolID).(string)
	vmIDUntyped, hasVMID := d.GetOk(mkResourceVirtualEnvironmentVMVMID)
	vmID := vmIDUntyped.(int)

	if !hasVMID {
		vmIDNew, err := api.Cluster().GetVMID(ctx)
		if err != nil {
			return diag.FromErr(err)
		}

		vmID = *vmIDNew
		err = d.Set(mkResourceVirtualEnvironmentVMVMID, vmID)

		if err != nil {
			return diag.FromErr(err)
		}
	}

	fullCopy := types.CustomBool(cloneFull)

	cloneBody := &vms.CloneRequestBody{
		FullCopy: &fullCopy,
		VMIDNew:  vmID,
	}

	if cloneDatastoreID != "" {
		cloneBody.TargetStorage = &cloneDatastoreID
	}

	if description != "" {
		cloneBody.Description = &description
	}

	if name != "" {
		cloneBody.Name = &name
	}

	if poolID != "" {
		cloneBody.PoolID = &poolID
	}

	cloneTimeout := d.Get(mkResourceVirtualEnvironmentVMTimeoutClone).(int)

	if cloneNodeName != "" && cloneNodeName != nodeName {
		// Check if any used datastores of the source VM are not shared
		vmConfig, err := api.Node(cloneNodeName).VM(cloneVMID).GetVM(ctx)
		if err != nil {
			return diag.FromErr(err)
		}

		datastores := getDiskDatastores(vmConfig, d)

		onlySharedDatastores := true
		for _, datastore := range datastores {
			datastoreStatus, err2 := api.Node(cloneNodeName).GetDatastoreStatus(ctx, datastore)
			if err2 != nil {
				return diag.FromErr(err2)
			}

			if datastoreStatus.Shared != nil && !*datastoreStatus.Shared {
				onlySharedDatastores = false
				break
			}
		}

		if onlySharedDatastores {
			// If the source and the target node are not the same, only clone directly to the target node if
			//  all used datastores in the source VM are shared. Directly cloning to non-shared storage
			//  on a different node is currently not supported by proxmox.
			cloneBody.TargetNodeName = &nodeName
			err = api.Node(cloneNodeName).VM(cloneVMID).CloneVM(
				ctx,
				cloneRetries,
				cloneBody,
				cloneTimeout,
			)
			if err != nil {
				return diag.FromErr(err)
			}
		} else {
			// If the source and the target node are not the same and any used datastore in the source VM is
			//  not shared, clone to the source node and then migrate to the target node. This is a workaround
			//  for missing functionality in the proxmox api as recommended per
			//  https://forum.proxmox.com/threads/500-cant-clone-to-non-shared-storage-local.49078/#post-229727

			// Temporarily clone to local node
			err = api.Node(cloneNodeName).VM(cloneVMID).CloneVM(ctx, cloneRetries, cloneBody, cloneTimeout)
			if err != nil {
				return diag.FromErr(err)
			}

			// Wait for the virtual machine to be created and its configuration lock to be released before migrating.

			err = api.Node(cloneNodeName).VM(vmID).WaitForVMConfigUnlock(ctx, 600, 5, true)
			if err != nil {
				return diag.FromErr(err)
			}

			// Migrate to target node
			withLocalDisks := types.CustomBool(true)
			migrateBody := &vms.MigrateRequestBody{
				TargetNode:     nodeName,
				WithLocalDisks: &withLocalDisks,
			}

			if cloneDatastoreID != "" {
				migrateBody.TargetStorage = &cloneDatastoreID
			}

			err = api.Node(cloneNodeName).VM(vmID).MigrateVM(ctx, migrateBody, cloneTimeout)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	} else {
		e = api.Node(nodeName).VM(cloneVMID).CloneVM(ctx, cloneRetries, cloneBody, cloneTimeout)
	}

	if e != nil {
		return diag.FromErr(e)
	}

	d.SetId(strconv.Itoa(vmID))

	vmAPI := api.Node(nodeName).VM(vmID)

	// Wait for the virtual machine to be created and its configuration lock to be released.
	e = vmAPI.WaitForVMConfigUnlock(ctx, 600, 5, true)
	if e != nil {
		return diag.FromErr(e)
	}

	// Now that the virtual machine has been cloned, we need to perform some modifications.
	acpi := types.CustomBool(d.Get(mkResourceVirtualEnvironmentVMACPI).(bool))
	agent := d.Get(mkResourceVirtualEnvironmentVMAgent).([]interface{})
	audioDevices := vmGetAudioDeviceList(d)

	bios := d.Get(mkResourceVirtualEnvironmentVMBIOS).(string)
	kvmArguments := d.Get(mkResourceVirtualEnvironmentVMKVMArguments).(string)
	scsiHardware := d.Get(mkResourceVirtualEnvironmentVMSCSIHardware).(string)
	cdrom := d.Get(mkResourceVirtualEnvironmentVMCDROM).([]interface{})
	cpu := d.Get(mkResourceVirtualEnvironmentVMCPU).([]interface{})
	initialization := d.Get(mkResourceVirtualEnvironmentVMInitialization).([]interface{})
	hostPCI := d.Get(mkResourceVirtualEnvironmentVMHostPCI).([]interface{})
	keyboardLayout := d.Get(mkResourceVirtualEnvironmentVMKeyboardLayout).(string)
	memory := d.Get(mkResourceVirtualEnvironmentVMMemory).([]interface{})
	networkDevice := d.Get(mkResourceVirtualEnvironmentVMNetworkDevice).([]interface{})
	operatingSystem := d.Get(mkResourceVirtualEnvironmentVMOperatingSystem).([]interface{})
	serialDevice := d.Get(mkResourceVirtualEnvironmentVMSerialDevice).([]interface{})
	onBoot := types.CustomBool(d.Get(mkResourceVirtualEnvironmentVMOnBoot).(bool))
	tabletDevice := types.CustomBool(d.Get(mkResourceVirtualEnvironmentVMTabletDevice).(bool))
	template := types.CustomBool(d.Get(mkResourceVirtualEnvironmentVMTemplate).(bool))
	vga := d.Get(mkResourceVirtualEnvironmentVMVGA).([]interface{})

	updateBody := &vms.UpdateRequestBody{
		AudioDevices: audioDevices,
	}

	ideDevices := vms.CustomStorageDevices{}

	var del []string

	//nolint:gosimple
	if acpi != dvResourceVirtualEnvironmentVMACPI {
		updateBody.ACPI = &acpi
	}

	if len(agent) > 0 {
		agentBlock := agent[0].(map[string]interface{})

		agentEnabled := types.CustomBool(
			agentBlock[mkResourceVirtualEnvironmentVMAgentEnabled].(bool),
		)
		agentTrim := types.CustomBool(agentBlock[mkResourceVirtualEnvironmentVMAgentTrim].(bool))
		agentType := agentBlock[mkResourceVirtualEnvironmentVMAgentType].(string)

		updateBody.Agent = &vms.CustomAgent{
			Enabled:         &agentEnabled,
			TrimClonedDisks: &agentTrim,
			Type:            &agentType,
		}
	}

	if kvmArguments != dvResourceVirtualEnvironmentVMKVMArguments {
		updateBody.KVMArguments = &kvmArguments
	}

	if bios != dvResourceVirtualEnvironmentVMBIOS {
		updateBody.BIOS = &bios
	}

	if scsiHardware != dvResourceVirtualEnvironmentVMSCSIHardware {
		updateBody.SCSIHardware = &scsiHardware
	}

	if len(cdrom) > 0 || len(initialization) > 0 {
		ideDevices = vms.CustomStorageDevices{
			"ide0": vms.CustomStorageDevice{
				Enabled: false,
			},
			"ide1": vms.CustomStorageDevice{
				Enabled: false,
			},
			"ide2": vms.CustomStorageDevice{
				Enabled: false,
			},
			"ide3": vms.CustomStorageDevice{
				Enabled: false,
			},
		}
	}

	if len(cdrom) > 0 {
		cdromBlock := cdrom[0].(map[string]interface{})

		cdromEnabled := cdromBlock[mkResourceVirtualEnvironmentVMCDROMEnabled].(bool)
		cdromFileID := cdromBlock[mkResourceVirtualEnvironmentVMCDROMFileID].(string)

		if cdromFileID == "" {
			cdromFileID = "cdrom"
		}

		cdromMedia := "cdrom"

		ideDevices = vms.CustomStorageDevices{
			"ide3": vms.CustomStorageDevice{
				Enabled:    cdromEnabled,
				FileVolume: cdromFileID,
				Media:      &cdromMedia,
			},
		}
	}

	if len(cpu) > 0 {
		cpuBlock := cpu[0].(map[string]interface{})

		cpuArchitecture := cpuBlock[mkResourceVirtualEnvironmentVMCPUArchitecture].(string)
		cpuCores := cpuBlock[mkResourceVirtualEnvironmentVMCPUCores].(int)
		cpuFlags := cpuBlock[mkResourceVirtualEnvironmentVMCPUFlags].([]interface{})
		cpuHotplugged := cpuBlock[mkResourceVirtualEnvironmentVMCPUHotplugged].(int)
		cpuNUMA := types.CustomBool(cpuBlock[mkResourceVirtualEnvironmentVMCPUNUMA].(bool))
		cpuSockets := cpuBlock[mkResourceVirtualEnvironmentVMCPUSockets].(int)
		cpuType := cpuBlock[mkResourceVirtualEnvironmentVMCPUType].(string)
		cpuUnits := cpuBlock[mkResourceVirtualEnvironmentVMCPUUnits].(int)

		cpuFlagsConverted := make([]string, len(cpuFlags))

		for fi, flag := range cpuFlags {
			cpuFlagsConverted[fi] = flag.(string)
		}

		// Only the root account is allowed to change the CPU architecture, which makes this check necessary.
		if api.API().IsRootTicket() ||
			cpuArchitecture != dvResourceVirtualEnvironmentVMCPUArchitecture {
			updateBody.CPUArchitecture = &cpuArchitecture
		}

		updateBody.CPUCores = &cpuCores
		updateBody.CPUEmulation = &vms.CustomCPUEmulation{
			Flags: &cpuFlagsConverted,
			Type:  cpuType,
		}
		updateBody.NUMAEnabled = &cpuNUMA
		updateBody.CPUSockets = &cpuSockets
		updateBody.CPUUnits = &cpuUnits

		if cpuHotplugged > 0 {
			updateBody.VirtualCPUCount = &cpuHotplugged
		}
	}

	if len(initialization) > 0 {
		initializationBlock := initialization[0].(map[string]interface{})
		initializationDatastoreID := initializationBlock[mkResourceVirtualEnvironmentVMInitializationDatastoreID].(string)

		const cdromCloudInitEnabled = true
		cdromCloudInitFileID := fmt.Sprintf("%s:cloudinit", initializationDatastoreID)
		cdromCloudInitMedia := "cdrom"

		ideDevices = vms.CustomStorageDevices{
			"ide2": vms.CustomStorageDevice{
				Enabled:    cdromCloudInitEnabled,
				FileVolume: cdromCloudInitFileID,
				Media:      &cdromCloudInitMedia,
			},
		}
		ciUpdateBody := &vms.UpdateRequestBody{}
		ciUpdateBody.Delete = append(ciUpdateBody.Delete, "ide2")

		e = vmAPI.UpdateVM(ctx, ciUpdateBody)
		if e != nil {
			return diag.FromErr(e)
		}

		updateBody.CloudInitConfig = vmGetCloudInitConfig(d)
	}

	if len(hostPCI) > 0 {
		updateBody.PCIDevices = vmGetHostPCIDeviceObjects(d)
	}

	if len(cdrom) > 0 || len(initialization) > 0 {
		updateBody.IDEDevices = ideDevices
	}

	if keyboardLayout != dvResourceVirtualEnvironmentVMKeyboardLayout {
		updateBody.KeyboardLayout = &keyboardLayout
	}

	if len(memory) > 0 {
		memoryBlock := memory[0].(map[string]interface{})

		memoryDedicated := memoryBlock[mkResourceVirtualEnvironmentVMMemoryDedicated].(int)
		memoryFloating := memoryBlock[mkResourceVirtualEnvironmentVMMemoryFloating].(int)
		memoryShared := memoryBlock[mkResourceVirtualEnvironmentVMMemoryShared].(int)

		updateBody.DedicatedMemory = &memoryDedicated
		updateBody.FloatingMemory = &memoryFloating

		if memoryShared > 0 {
			memorySharedName := fmt.Sprintf("vm-%d-ivshmem", vmID)

			updateBody.SharedMemory = &vms.CustomSharedMemory{
				Name: &memorySharedName,
				Size: memoryShared,
			}
		}
	}

	if len(networkDevice) > 0 {
		updateBody.NetworkDevices = vmGetNetworkDeviceObjects(d)

		for i := 0; i < len(updateBody.NetworkDevices); i++ {
			if !updateBody.NetworkDevices[i].Enabled {
				del = append(del, fmt.Sprintf("net%d", i))
			}
		}

		for i := len(updateBody.NetworkDevices); i < maxResourceVirtualEnvironmentVMNetworkDevices; i++ {
			del = append(del, fmt.Sprintf("net%d", i))
		}
	}

	if len(operatingSystem) > 0 {
		operatingSystemBlock := operatingSystem[0].(map[string]interface{})
		operatingSystemType := operatingSystemBlock[mkResourceVirtualEnvironmentVMOperatingSystemType].(string)

		updateBody.OSType = &operatingSystemType
	}

	if len(serialDevice) > 0 {
		updateBody.SerialDevices = vmGetSerialDeviceList(d)

		for i := len(updateBody.SerialDevices); i < maxResourceVirtualEnvironmentVMSerialDevices; i++ {
			del = append(del, fmt.Sprintf("serial%d", i))
		}
	}

	updateBody.StartOnBoot = &onBoot

	//nolint:gosimple
	if tabletDevice != dvResourceVirtualEnvironmentVMTabletDevice {
		updateBody.TabletDeviceEnabled = &tabletDevice
	}

	if len(tags) > 0 {
		tagString := vmGetTagsString(d)
		updateBody.Tags = &tagString
	}

	//nolint:gosimple
	if template != dvResourceVirtualEnvironmentVMTemplate {
		updateBody.Template = &template
	}

	if len(vga) > 0 {
		vgaDevice, err := vmGetVGADeviceObject(d)
		if err != nil {
			return diag.FromErr(err)
		}

		updateBody.VGADevice = vgaDevice
	}

	updateBody.Delete = del

	e = vmAPI.UpdateVM(ctx, updateBody)
	if e != nil {
		return diag.FromErr(e)
	}

	disk := d.Get(mkResourceVirtualEnvironmentVMDisk).([]interface{})
	efiDisk := d.Get(mkResourceVirtualEnvironmentVMEFIDisk).([]interface{})

	vmConfig, e := vmAPI.GetVM(ctx)
	if e != nil {
		if strings.Contains(e.Error(), "HTTP 404") ||
			(strings.Contains(e.Error(), "HTTP 500") && strings.Contains(e.Error(), "does not exist")) {
			d.SetId("")

			return nil
		}

		return diag.FromErr(e)
	}

	allDiskInfo := getDiskInfo(vmConfig, d) // from the cloned VM

	diskDeviceObjects, e := vmGetDiskDeviceObjects(d, nil) // from the resource config
	if e != nil {
		return diag.FromErr(e)
	}

	for i := range disk {
		diskBlock := disk[i].(map[string]interface{})
		diskInterface := diskBlock[mkResourceVirtualEnvironmentVMDiskInterface].(string)
		dataStoreID := diskBlock[mkResourceVirtualEnvironmentVMDiskDatastoreID].(string)
		diskSize := diskBlock[mkResourceVirtualEnvironmentVMDiskSize].(int)
		prefix := diskDigitPrefix(diskInterface)

		currentDiskInfo := allDiskInfo[diskInterface]
		configuredDiskInfo := diskDeviceObjects[prefix][diskInterface]

		if currentDiskInfo == nil {
			diskUpdateBody := &vms.UpdateRequestBody{}

			switch prefix {
			case "virtio":
				if diskUpdateBody.VirtualIODevices == nil {
					diskUpdateBody.VirtualIODevices = vms.CustomStorageDevices{}
				}

				diskUpdateBody.VirtualIODevices[diskInterface] = configuredDiskInfo
			case "sata":
				if diskUpdateBody.SATADevices == nil {
					diskUpdateBody.SATADevices = vms.CustomStorageDevices{}
				}

				diskUpdateBody.SATADevices[diskInterface] = configuredDiskInfo
			case "scsi":
				if diskUpdateBody.SCSIDevices == nil {
					diskUpdateBody.SCSIDevices = vms.CustomStorageDevices{}
				}

				diskUpdateBody.SCSIDevices[diskInterface] = configuredDiskInfo
			}

			e = vmAPI.UpdateVM(ctx, diskUpdateBody)
			if e != nil {
				return diag.FromErr(e)
			}

			continue
		}

		if diskSize < currentDiskInfo.Size.InGigabytes() {
			return diag.Errorf(
				"disk resize fails requests size (%dG) is lower than current size (%s)",
				diskSize,
				*currentDiskInfo.Size,
			)
		}

		deleteOriginalDisk := types.CustomBool(true)

		diskMoveBody := &vms.MoveDiskRequestBody{
			DeleteOriginalDisk: &deleteOriginalDisk,
			Disk:               diskInterface,
			TargetStorage:      dataStoreID,
		}

		diskResizeBody := &vms.ResizeDiskRequestBody{
			Disk: diskInterface,
			Size: types.DiskSizeFromGigabytes(diskSize),
		}

		moveDisk := false

		if dataStoreID != "" {
			moveDisk = true

			if allDiskInfo[diskInterface] != nil {
				fileIDParts := strings.Split(allDiskInfo[diskInterface].FileVolume, ":")
				moveDisk = dataStoreID != fileIDParts[0]
			}
		}

		if moveDisk {
			moveDiskTimeout := d.Get(mkResourceVirtualEnvironmentVMTimeoutMoveDisk).(int)

			e = vmAPI.MoveVMDisk(ctx, diskMoveBody, moveDiskTimeout)
			if e != nil {
				return diag.FromErr(e)
			}
		}

		if diskSize > currentDiskInfo.Size.InGigabytes() {
			e = vmAPI.ResizeVMDisk(ctx, diskResizeBody)
			if e != nil {
				return diag.FromErr(e)
			}
		}
	}

	efiDiskInfo := vmGetEfiDisk(d, nil) // from the resource config

	for i := range efiDisk {
		diskBlock := efiDisk[i].(map[string]interface{})
		diskInterface := "efidisk0"
		dataStoreID := diskBlock[mkResourceVirtualEnvironmentVMEFIDiskDatastoreID].(string)
		efiType := diskBlock[mkResourceVirtualEnvironmentVMEFIDiskType].(string)

		currentDiskInfo := vmConfig.EFIDisk
		configuredDiskInfo := efiDiskInfo

		if currentDiskInfo == nil {
			diskUpdateBody := &vms.UpdateRequestBody{}

			diskUpdateBody.EFIDisk = configuredDiskInfo

			e = vmAPI.UpdateVM(ctx, diskUpdateBody)
			if e != nil {
				return diag.FromErr(e)
			}

			continue
		}

		if &efiType != currentDiskInfo.Type {
			return diag.Errorf(
				"resizing of efidisks is not supported.",
			)
		}

		deleteOriginalDisk := types.CustomBool(true)

		diskMoveBody := &vms.MoveDiskRequestBody{
			DeleteOriginalDisk: &deleteOriginalDisk,
			Disk:               diskInterface,
			TargetStorage:      dataStoreID,
		}

		moveDisk := false

		if dataStoreID != "" {
			moveDisk = true

			if allDiskInfo[diskInterface] != nil {
				fileIDParts := strings.Split(allDiskInfo[diskInterface].FileVolume, ":")
				moveDisk = dataStoreID != fileIDParts[0]
			}
		}

		if moveDisk {
			moveDiskTimeout := d.Get(mkResourceVirtualEnvironmentVMTimeoutMoveDisk).(int)

			e = vmAPI.MoveVMDisk(ctx, diskMoveBody, moveDiskTimeout)
			if e != nil {
				return diag.FromErr(e)
			}
		}
	}

	return vmCreateStart(ctx, d, m)
}

func vmCreateCustom(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(proxmoxtf.ProviderConfiguration)
	api, err := config.GetClient()
	if err != nil {
		return diag.FromErr(err)
	}

	resource := VM()

	acpi := types.CustomBool(d.Get(mkResourceVirtualEnvironmentVMACPI).(bool))

	agentBlock, err := getSchemaBlock(
		resource,
		d,
		[]string{mkResourceVirtualEnvironmentVMAgent},
		0,
		true,
	)
	if err != nil {
		return diag.FromErr(err)
	}

	agentEnabled := types.CustomBool(
		agentBlock[mkResourceVirtualEnvironmentVMAgentEnabled].(bool),
	)
	agentTrim := types.CustomBool(agentBlock[mkResourceVirtualEnvironmentVMAgentTrim].(bool))
	agentType := agentBlock[mkResourceVirtualEnvironmentVMAgentType].(string)

	kvmArguments := d.Get(mkResourceVirtualEnvironmentVMKVMArguments).(string)

	audioDevices := vmGetAudioDeviceList(d)

	bios := d.Get(mkResourceVirtualEnvironmentVMBIOS).(string)

	cdromBlock, err := getSchemaBlock(
		resource,
		d,
		[]string{mkResourceVirtualEnvironmentVMCDROM},
		0,
		true,
	)
	if err != nil {
		return diag.FromErr(err)
	}

	cdromEnabled := cdromBlock[mkResourceVirtualEnvironmentVMCDROMEnabled].(bool)
	cdromFileID := cdromBlock[mkResourceVirtualEnvironmentVMCDROMFileID].(string)

	cdromCloudInitEnabled := false
	cdromCloudInitFileID := ""

	if cdromFileID == "" {
		cdromFileID = "cdrom"
	}

	cpuBlock, err := getSchemaBlock(
		resource,
		d,
		[]string{mkResourceVirtualEnvironmentVMCPU},
		0,
		true,
	)
	if err != nil {
		return diag.FromErr(err)
	}

	cpuArchitecture := cpuBlock[mkResourceVirtualEnvironmentVMCPUArchitecture].(string)
	cpuCores := cpuBlock[mkResourceVirtualEnvironmentVMCPUCores].(int)
	cpuFlags := cpuBlock[mkResourceVirtualEnvironmentVMCPUFlags].([]interface{})
	cpuHotplugged := cpuBlock[mkResourceVirtualEnvironmentVMCPUHotplugged].(int)
	cpuSockets := cpuBlock[mkResourceVirtualEnvironmentVMCPUSockets].(int)
	cpuNUMA := types.CustomBool(cpuBlock[mkResourceVirtualEnvironmentVMCPUNUMA].(bool))
	cpuType := cpuBlock[mkResourceVirtualEnvironmentVMCPUType].(string)
	cpuUnits := cpuBlock[mkResourceVirtualEnvironmentVMCPUUnits].(int)

	description := d.Get(mkResourceVirtualEnvironmentVMDescription).(string)
	diskDeviceObjects, err := vmGetDiskDeviceObjects(d, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	var efiDisk *vms.CustomEFIDisk

	efiDiskBlock := d.Get(mkResourceVirtualEnvironmentVMEFIDisk).([]interface{})
	if len(efiDiskBlock) > 0 {
		block := efiDiskBlock[0].(map[string]interface{})

		datastoreID, _ := block[mkResourceVirtualEnvironmentVMEFIDiskDatastoreID].(string)
		fileFormat, _ := block[mkResourceVirtualEnvironmentVMEFIDiskFileFormat].(string)
		efiType, _ := block[mkResourceVirtualEnvironmentVMEFIDiskType].(string)
		preEnrolledKeys := types.CustomBool(block[mkResourceVirtualEnvironmentVMEFIDiskPreEnrolledKeys].(bool))

		if fileFormat == "" {
			fileFormat = dvResourceVirtualEnvironmentVMEFIDiskFileFormat
		}

		efiDisk = &vms.CustomEFIDisk{
			Type:            &efiType,
			FileVolume:      fmt.Sprintf("%s:1", datastoreID),
			Format:          &fileFormat,
			PreEnrolledKeys: &preEnrolledKeys,
		}
	}

	virtioDeviceObjects := diskDeviceObjects["virtio"]
	scsiDeviceObjects := diskDeviceObjects["scsi"]
	// ideDeviceObjects := getOrderedDiskDeviceList(diskDeviceObjects, "ide")
	sataDeviceObjects := diskDeviceObjects["sata"]

	initializationConfig := vmGetCloudInitConfig(d)

	if initializationConfig != nil {
		initialization := d.Get(mkResourceVirtualEnvironmentVMInitialization).([]interface{})
		initializationBlock := initialization[0].(map[string]interface{})
		initializationDatastoreID := initializationBlock[mkResourceVirtualEnvironmentVMInitializationDatastoreID].(string)

		cdromCloudInitEnabled = true
		cdromCloudInitFileID = fmt.Sprintf("%s:cloudinit", initializationDatastoreID)
	}

	pciDeviceObjects := vmGetHostPCIDeviceObjects(d)

	keyboardLayout := d.Get(mkResourceVirtualEnvironmentVMKeyboardLayout).(string)
	memoryBlock, err := getSchemaBlock(
		resource,
		d,
		[]string{mkResourceVirtualEnvironmentVMMemory},
		0,
		true,
	)
	if err != nil {
		return diag.FromErr(err)
	}

	memoryDedicated := memoryBlock[mkResourceVirtualEnvironmentVMMemoryDedicated].(int)
	memoryFloating := memoryBlock[mkResourceVirtualEnvironmentVMMemoryFloating].(int)
	memoryShared := memoryBlock[mkResourceVirtualEnvironmentVMMemoryShared].(int)

	machine := d.Get(mkResourceVirtualEnvironmentVMMachine).(string)
	name := d.Get(mkResourceVirtualEnvironmentVMName).(string)
	tags := d.Get(mkResourceVirtualEnvironmentVMTags).([]interface{})

	networkDeviceObjects := vmGetNetworkDeviceObjects(d)

	nodeName := d.Get(mkResourceVirtualEnvironmentVMNodeName).(string)

	operatingSystem, err := getSchemaBlock(
		resource,
		d,
		[]string{mkResourceVirtualEnvironmentVMOperatingSystem},
		0,
		true,
	)
	if err != nil {
		return diag.FromErr(err)
	}

	operatingSystemType := operatingSystem[mkResourceVirtualEnvironmentVMOperatingSystemType].(string)

	poolID := d.Get(mkResourceVirtualEnvironmentVMPoolID).(string)

	serialDevices := vmGetSerialDeviceList(d)

	onBoot := types.CustomBool(d.Get(mkResourceVirtualEnvironmentVMOnBoot).(bool))
	tabletDevice := types.CustomBool(d.Get(mkResourceVirtualEnvironmentVMTabletDevice).(bool))
	template := types.CustomBool(d.Get(mkResourceVirtualEnvironmentVMTemplate).(bool))

	vgaDevice, err := vmGetVGADeviceObject(d)
	if err != nil {
		return diag.FromErr(err)
	}

	vmIDUntyped, hasVMID := d.GetOk(mkResourceVirtualEnvironmentVMVMID)
	vmID := vmIDUntyped.(int)

	if !hasVMID {
		vmIDNew, e := api.Cluster().GetVMID(ctx)
		if e != nil {
			return diag.FromErr(e)
		}

		vmID = *vmIDNew
		e = d.Set(mkResourceVirtualEnvironmentVMVMID, vmID)

		if e != nil {
			return diag.FromErr(e)
		}
	}

	var memorySharedObject *vms.CustomSharedMemory

	var bootOrderConverted []string
	if cdromEnabled {
		bootOrderConverted = []string{"ide3"}
	}

	bootOrder := d.Get(mkResourceVirtualEnvironmentVMBootOrder).([]interface{})
	//nolint:nestif
	if len(bootOrder) == 0 {
		if sataDeviceObjects != nil {
			bootOrderConverted = append(bootOrderConverted, "sata0")
		}

		if scsiDeviceObjects != nil {
			bootOrderConverted = append(bootOrderConverted, "scsi0")
		}

		if virtioDeviceObjects != nil {
			bootOrderConverted = append(bootOrderConverted, "virtio0")
		}

		if networkDeviceObjects != nil {
			bootOrderConverted = append(bootOrderConverted, "net0")
		}
	} else {
		bootOrderConverted = make([]string, len(bootOrder))
		for i, device := range bootOrder {
			bootOrderConverted[i] = device.(string)
		}
	}

	cpuFlagsConverted := make([]string, len(cpuFlags))
	for fi, flag := range cpuFlags {
		cpuFlagsConverted[fi] = flag.(string)
	}

	ideDevice2Media := "cdrom"
	ideDevices := vms.CustomStorageDevices{
		"ide2": vms.CustomStorageDevice{
			Enabled:    cdromCloudInitEnabled,
			FileVolume: cdromCloudInitFileID,
			Media:      &ideDevice2Media,
		},
		"ide3": vms.CustomStorageDevice{
			Enabled:    cdromEnabled,
			FileVolume: cdromFileID,
			Media:      &ideDevice2Media,
		},
	}

	if memoryShared > 0 {
		memorySharedName := fmt.Sprintf("vm-%d-ivshmem", vmID)
		memorySharedObject = &vms.CustomSharedMemory{
			Name: &memorySharedName,
			Size: memoryShared,
		}
	}

	scsiHardware := d.Get(mkResourceVirtualEnvironmentVMSCSIHardware).(string)

	createBody := &vms.CreateRequestBody{
		ACPI: &acpi,
		Agent: &vms.CustomAgent{
			Enabled:         &agentEnabled,
			TrimClonedDisks: &agentTrim,
			Type:            &agentType,
		},
		AudioDevices: audioDevices,
		BIOS:         &bios,
		Boot: &vms.CustomBoot{
			Order: &bootOrderConverted,
		},
		CloudInitConfig: initializationConfig,
		CPUCores:        &cpuCores,
		CPUEmulation: &vms.CustomCPUEmulation{
			Flags: &cpuFlagsConverted,
			Type:  cpuType,
		},
		CPUSockets:          &cpuSockets,
		CPUUnits:            &cpuUnits,
		DedicatedMemory:     &memoryDedicated,
		EFIDisk:             efiDisk,
		FloatingMemory:      &memoryFloating,
		IDEDevices:          ideDevices,
		KeyboardLayout:      &keyboardLayout,
		NetworkDevices:      networkDeviceObjects,
		NUMAEnabled:         &cpuNUMA,
		OSType:              &operatingSystemType,
		PCIDevices:          pciDeviceObjects,
		SCSIHardware:        &scsiHardware,
		SerialDevices:       serialDevices,
		SharedMemory:        memorySharedObject,
		StartOnBoot:         &onBoot,
		TabletDeviceEnabled: &tabletDevice,
		Template:            &template,
		VGADevice:           vgaDevice,
		VMID:                &vmID,
	}

	if sataDeviceObjects != nil {
		createBody.SATADevices = sataDeviceObjects
	}

	if scsiDeviceObjects != nil {
		createBody.SCSIDevices = scsiDeviceObjects
	}

	if virtioDeviceObjects != nil {
		createBody.VirtualIODevices = virtioDeviceObjects
	}

	// Only the root account is allowed to change the CPU architecture, which makes this check necessary.
	if api.API().IsRootTicket() ||
		cpuArchitecture != dvResourceVirtualEnvironmentVMCPUArchitecture {
		createBody.CPUArchitecture = &cpuArchitecture
	}

	if cpuHotplugged > 0 {
		createBody.VirtualCPUCount = &cpuHotplugged
	}

	if description != "" {
		createBody.Description = &description
	}

	if len(tags) > 0 {
		tagsString := vmGetTagsString(d)
		createBody.Tags = &tagsString
	}

	if kvmArguments != "" {
		createBody.KVMArguments = &kvmArguments
	}

	if machine != "" {
		createBody.Machine = &machine
	}

	if name != "" {
		createBody.Name = &name
	}

	if poolID != "" {
		createBody.PoolID = &poolID
	}

	err = api.Node(nodeName).VM(0).CreateVM(ctx, createBody, vmCreateTimeoutSeconds)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(vmID))

	return vmCreateCustomDisks(ctx, d, m)
}

func vmCreateCustomDisks(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	vmID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	var commands []string

	// Determine the ID of the next disk.
	disk := d.Get(mkResourceVirtualEnvironmentVMDisk).([]interface{})
	diskCount := 0

	for _, d := range disk {
		block := d.(map[string]interface{})
		fileID, _ := block[mkResourceVirtualEnvironmentVMDiskFileID].(string)

		if fileID == "" {
			diskCount++
		}
	}

	// Retrieve some information about the disk schema.
	resourceSchema := VM().Schema
	diskSchemaElem := resourceSchema[mkResourceVirtualEnvironmentVMDisk].Elem
	diskSchemaResource := diskSchemaElem.(*schema.Resource)
	diskSpeedResource := diskSchemaResource.Schema[mkResourceVirtualEnvironmentVMDiskSpeed]

	// Generate the commands required to import the specified disks.
	importedDiskCount := 0

	for _, d := range disk {
		block := d.(map[string]interface{})

		fileID, _ := block[mkResourceVirtualEnvironmentVMDiskFileID].(string)

		if fileID == "" {
			continue
		}

		datastoreID, _ := block[mkResourceVirtualEnvironmentVMDiskDatastoreID].(string)
		fileFormat, _ := block[mkResourceVirtualEnvironmentVMDiskFileFormat].(string)
		size, _ := block[mkResourceVirtualEnvironmentVMDiskSize].(int)
		speed := block[mkResourceVirtualEnvironmentVMDiskSpeed].([]interface{})
		diskInterface, _ := block[mkResourceVirtualEnvironmentVMDiskInterface].(string)
		ioThread := types.CustomBool(block[mkResourceVirtualEnvironmentVMDiskIOThread].(bool))
		ssd := types.CustomBool(block[mkResourceVirtualEnvironmentVMDiskSSD].(bool))
		discard, _ := block[mkResourceVirtualEnvironmentVMDiskDiscard].(string)

		if fileFormat == "" {
			fileFormat = dvResourceVirtualEnvironmentVMDiskFileFormat
		}

		if len(speed) == 0 {
			diskSpeedDefault, err := diskSpeedResource.DefaultValue()
			if err != nil {
				return diag.FromErr(err)
			}
			speed = diskSpeedDefault.([]interface{})
		}

		speedBlock := speed[0].(map[string]interface{})
		speedLimitRead := speedBlock[mkResourceVirtualEnvironmentVMDiskSpeedRead].(int)
		speedLimitReadBurstable := speedBlock[mkResourceVirtualEnvironmentVMDiskSpeedReadBurstable].(int)
		speedLimitWrite := speedBlock[mkResourceVirtualEnvironmentVMDiskSpeedWrite].(int)
		speedLimitWriteBurstable := speedBlock[mkResourceVirtualEnvironmentVMDiskSpeedWriteBurstable].(int)

		diskOptions := ""

		if ioThread {
			diskOptions += ",iothread=1"
		}

		if ssd {
			diskOptions += ",ssd=1"
		}

		if discard != "" {
			diskOptions += fmt.Sprintf(",discard=%s", discard)
		}

		if speedLimitRead > 0 {
			diskOptions += fmt.Sprintf(",mbps_rd=%d", speedLimitRead)
		}

		if speedLimitReadBurstable > 0 {
			diskOptions += fmt.Sprintf(",mbps_rd_max=%d", speedLimitReadBurstable)
		}

		if speedLimitWrite > 0 {
			diskOptions += fmt.Sprintf(",mbps_wr=%d", speedLimitWrite)
		}

		if speedLimitWriteBurstable > 0 {
			diskOptions += fmt.Sprintf(",mbps_wr_max=%d", speedLimitWriteBurstable)
		}

		filePathTmp := fmt.Sprintf(
			"/tmp/vm-%d-disk-%d.%s",
			vmID,
			diskCount+importedDiskCount,
			fileFormat,
		)

		//nolint:lll
		commands = append(
			commands,
			`set -e`,
			fmt.Sprintf(`file_id="%s"`, fileID),
			fmt.Sprintf(`file_format="%s"`, fileFormat),
			fmt.Sprintf(`datastore_id_target="%s"`, datastoreID),
			fmt.Sprintf(`disk_options="%s"`, diskOptions),
			fmt.Sprintf(`disk_size="%d"`, size),
			fmt.Sprintf(`disk_interface="%s"`, diskInterface),
			fmt.Sprintf(`file_path_tmp="%s"`, filePathTmp),
			fmt.Sprintf(`vm_id="%d"`, vmID),
			`source_image=$(pvesm path "$file_id")`,
			`cp "$source_image" "$file_path_tmp"`,
			`qemu-img resize -f "$file_format" "$file_path_tmp" "${disk_size}G"`,
			`imported_disk="$(qm importdisk "$vm_id" "$file_path_tmp" "$datastore_id_target" -format $file_format | grep "unused0" | cut -d ":" -f 3 | cut -d "'" -f 1)"`,
			`disk_id="${datastore_id_target}:$imported_disk${disk_options}"`,
			`qm set "$vm_id" "-${disk_interface}" "$disk_id"`,
			`rm -f "$file_path_tmp"`,
		)

		importedDiskCount++
	}

	// Execute the commands on the node and wait for the result.
	// This is a highly experimental approach to disk imports and is not recommended by Proxmox.
	if len(commands) > 0 {
		config := m.(proxmoxtf.ProviderConfiguration)

		api, err := config.GetClient()
		if err != nil {
			return diag.FromErr(err)
		}

		nodeName := d.Get(mkResourceVirtualEnvironmentVMNodeName).(string)

		err = api.SSH().ExecuteNodeCommands(ctx, nodeName, commands)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return vmCreateStart(ctx, d, m)
}

func vmCreateStart(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	started := d.Get(mkResourceVirtualEnvironmentVMStarted).(bool)
	template := d.Get(mkResourceVirtualEnvironmentVMTemplate).(bool)
	reboot := d.Get(mkResourceVirtualEnvironmentVMRebootAfterCreation).(bool)

	if !started || template {
		return vmRead(ctx, d, m)
	}

	config := m.(proxmoxtf.ProviderConfiguration)
	api, err := config.GetClient()
	if err != nil {
		return diag.FromErr(err)
	}

	nodeName := d.Get(mkResourceVirtualEnvironmentVMNodeName).(string)
	vmID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	vmAPI := api.Node(nodeName).VM(vmID)

	// Start the virtual machine and wait for it to reach a running state before continuing.
	startVMTimeout := d.Get(mkResourceVirtualEnvironmentVMTimeoutStartVM).(int)
	err = vmAPI.StartVM(ctx, startVMTimeout)
	if err != nil {
		return diag.FromErr(err)
	}

	if reboot {
		rebootTimeout := d.Get(mkResourceVirtualEnvironmentVMTimeoutReboot).(int)

		err := vmAPI.RebootVM(
			ctx,
			&vms.RebootRequestBody{
				Timeout: &rebootTimeout,
			},
			rebootTimeout+30,
		)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return vmRead(ctx, d, m)
}

func vmGetAudioDeviceList(d *schema.ResourceData) vms.CustomAudioDevices {
	devices := d.Get(mkResourceVirtualEnvironmentVMAudioDevice).([]interface{})
	list := make(vms.CustomAudioDevices, len(devices))

	for i, v := range devices {
		block := v.(map[string]interface{})

		device, _ := block[mkResourceVirtualEnvironmentVMAudioDeviceDevice].(string)
		driver, _ := block[mkResourceVirtualEnvironmentVMAudioDeviceDriver].(string)
		enabled, _ := block[mkResourceVirtualEnvironmentVMAudioDeviceEnabled].(bool)

		list[i].Device = device
		list[i].Driver = &driver
		list[i].Enabled = enabled
	}

	return list
}

func vmGetAudioDeviceValidator() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.StringInSlice([]string{
		"AC97",
		"ich9-intel-hda",
		"intel-hda",
	}, false))
}

func vmGetAudioDriverValidator() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.StringInSlice([]string{
		"spice",
	}, false))
}

func vmGetCloudInitConfig(d *schema.ResourceData) *vms.CustomCloudInitConfig {
	var initializationConfig *vms.CustomCloudInitConfig

	initialization := d.Get(mkResourceVirtualEnvironmentVMInitialization).([]interface{})

	if len(initialization) > 0 {
		initializationBlock := initialization[0].(map[string]interface{})
		initializationConfig = &vms.CustomCloudInitConfig{}
		initializationDNS := initializationBlock[mkResourceVirtualEnvironmentVMInitializationDNS].([]interface{})

		if len(initializationDNS) > 0 {
			initializationDNSBlock := initializationDNS[0].(map[string]interface{})
			domain := initializationDNSBlock[mkResourceVirtualEnvironmentVMInitializationDNSDomain].(string)

			if domain != "" {
				initializationConfig.SearchDomain = &domain
			}

			server := initializationDNSBlock[mkResourceVirtualEnvironmentVMInitializationDNSServer].(string)

			if server != "" {
				initializationConfig.Nameserver = &server
			}
		}

		initializationIPConfig := initializationBlock[mkResourceVirtualEnvironmentVMInitializationIPConfig].([]interface{})
		initializationConfig.IPConfig = make(
			[]vms.CustomCloudInitIPConfig,
			len(initializationIPConfig),
		)

		for i, c := range initializationIPConfig {
			configBlock := c.(map[string]interface{})
			ipv4 := configBlock[mkResourceVirtualEnvironmentVMInitializationIPConfigIPv4].([]interface{})

			if len(ipv4) > 0 {
				ipv4Block := ipv4[0].(map[string]interface{})
				ipv4Address := ipv4Block[mkResourceVirtualEnvironmentVMInitializationIPConfigIPv4Address].(string)

				if ipv4Address != "" {
					initializationConfig.IPConfig[i].IPv4 = &ipv4Address
				}

				ipv4Gateway := ipv4Block[mkResourceVirtualEnvironmentVMInitializationIPConfigIPv4Gateway].(string)

				if ipv4Gateway != "" {
					initializationConfig.IPConfig[i].GatewayIPv4 = &ipv4Gateway
				}
			}

			ipv6 := configBlock[mkResourceVirtualEnvironmentVMInitializationIPConfigIPv6].([]interface{})

			if len(ipv6) > 0 {
				ipv6Block := ipv6[0].(map[string]interface{})
				ipv6Address := ipv6Block[mkResourceVirtualEnvironmentVMInitializationIPConfigIPv6Address].(string)

				if ipv6Address != "" {
					initializationConfig.IPConfig[i].IPv6 = &ipv6Address
				}

				ipv6Gateway := ipv6Block[mkResourceVirtualEnvironmentVMInitializationIPConfigIPv6Gateway].(string)

				if ipv6Gateway != "" {
					initializationConfig.IPConfig[i].GatewayIPv6 = &ipv6Gateway
				}
			}
		}

		initializationUserAccount := initializationBlock[mkResourceVirtualEnvironmentVMInitializationUserAccount].([]interface{})

		if len(initializationUserAccount) > 0 {
			initializationUserAccountBlock := initializationUserAccount[0].(map[string]interface{})
			keys := initializationUserAccountBlock[mkResourceVirtualEnvironmentVMInitializationUserAccountKeys].([]interface{})

			if len(keys) > 0 {
				sshKeys := make(vms.CustomCloudInitSSHKeys, len(keys))

				for i, k := range keys {
					sshKeys[i] = k.(string)
				}

				initializationConfig.SSHKeys = &sshKeys
			}

			password := initializationUserAccountBlock[mkResourceVirtualEnvironmentVMInitializationUserAccountPassword].(string)

			if password != "" {
				initializationConfig.Password = &password
			}

			username := initializationUserAccountBlock[mkResourceVirtualEnvironmentVMInitializationUserAccountUsername].(string)

			initializationConfig.Username = &username
		}

		initializationUserDataFileID := initializationBlock[mkResourceVirtualEnvironmentVMInitializationUserDataFileID].(string)

		if initializationUserDataFileID != "" {
			initializationConfig.Files = &vms.CustomCloudInitFiles{
				UserVolume: &initializationUserDataFileID,
			}
		}

		initializationVendorDataFileID := initializationBlock[mkResourceVirtualEnvironmentVMInitializationVendorDataFileID].(string)

		if initializationVendorDataFileID != "" {
			if initializationConfig.Files == nil {
				initializationConfig.Files = &vms.CustomCloudInitFiles{}
			}

			initializationConfig.Files.VendorVolume = &initializationVendorDataFileID
		}

		initializationNetworkDataFileID := initializationBlock[mkResourceVirtualEnvironmentVMInitializationNetworkDataFileID].(string)

		if initializationNetworkDataFileID != "" {
			if initializationConfig.Files == nil {
				initializationConfig.Files = &vms.CustomCloudInitFiles{}
			}

			initializationConfig.Files.NetworkVolume = &initializationNetworkDataFileID
		}

		//nolint:lll
		initializationMetaDataFileID := initializationBlock[mkResourceVirtualEnvironmentVMInitializationMetaDataFileID].(string)

		if initializationMetaDataFileID != "" {
			if initializationConfig.Files == nil {
				initializationConfig.Files = &vms.CustomCloudInitFiles{}
			}

			initializationConfig.Files.MetaVolume = &initializationMetaDataFileID
		}

		initializationType := initializationBlock[mkResourceVirtualEnvironmentVMInitializationType].(string)

		if initializationType != "" {
			initializationConfig.Type = &initializationType
		}
	}

	return initializationConfig
}

func vmGetCPUArchitectureValidator() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.StringInSlice([]string{
		"aarch64",
		"x86_64",
	}, false))
}

func vmGetDiskDeviceObjects(
	d *schema.ResourceData,
	disks []interface{},
) (map[string]map[string]vms.CustomStorageDevice, error) {
	var diskDevice []interface{}

	if disks != nil {
		diskDevice = disks
	} else {
		diskDevice = d.Get(mkResourceVirtualEnvironmentVMDisk).([]interface{})
	}

	diskDeviceObjects := map[string]map[string]vms.CustomStorageDevice{}
	resource := VM()

	for _, diskEntry := range diskDevice {
		diskDevice := vms.CustomStorageDevice{
			Enabled: true,
		}

		block := diskEntry.(map[string]interface{})
		datastoreID, _ := block[mkResourceVirtualEnvironmentVMDiskDatastoreID].(string)
		fileFormat, _ := block[mkResourceVirtualEnvironmentVMDiskFileFormat].(string)
		fileID, _ := block[mkResourceVirtualEnvironmentVMDiskFileID].(string)
		size, _ := block[mkResourceVirtualEnvironmentVMDiskSize].(int)
		diskInterface, _ := block[mkResourceVirtualEnvironmentVMDiskInterface].(string)
		ioThread := types.CustomBool(block[mkResourceVirtualEnvironmentVMDiskIOThread].(bool))
		ssd := types.CustomBool(block[mkResourceVirtualEnvironmentVMDiskSSD].(bool))
		discard := block[mkResourceVirtualEnvironmentVMDiskDiscard].(string)

		speedBlock, err := getSchemaBlock(
			resource,
			d,
			[]string{mkResourceVirtualEnvironmentVMDisk, mkResourceVirtualEnvironmentVMDiskSpeed},
			0,
			false,
		)
		if err != nil {
			return diskDeviceObjects, err
		}

		if fileFormat == "" {
			fileFormat = dvResourceVirtualEnvironmentVMDiskFileFormat
		}
		if fileID != "" {
			diskDevice.Enabled = false
		} else {
			diskDevice.FileVolume = fmt.Sprintf("%s:%d", datastoreID, size)
		}

		diskDevice.ID = &datastoreID
		diskDevice.Interface = &diskInterface
		diskDevice.Format = &fileFormat
		diskDevice.FileID = &fileID
		diskSize := types.DiskSizeFromGigabytes(size)
		diskDevice.Size = &diskSize
		diskDevice.SizeInt = &size
		diskDevice.IOThread = &ioThread
		diskDevice.Discard = &discard

		if !strings.HasPrefix(diskInterface, "virtio") {
			diskDevice.SSD = &ssd
		}

		if len(speedBlock) > 0 {
			speedLimitRead := speedBlock[mkResourceVirtualEnvironmentVMDiskSpeedRead].(int)
			speedLimitReadBurstable := speedBlock[mkResourceVirtualEnvironmentVMDiskSpeedReadBurstable].(int)
			speedLimitWrite := speedBlock[mkResourceVirtualEnvironmentVMDiskSpeedWrite].(int)
			speedLimitWriteBurstable := speedBlock[mkResourceVirtualEnvironmentVMDiskSpeedWriteBurstable].(int)

			if speedLimitRead > 0 {
				diskDevice.MaxReadSpeedMbps = &speedLimitRead
			}

			if speedLimitReadBurstable > 0 {
				diskDevice.BurstableReadSpeedMbps = &speedLimitReadBurstable
			}

			if speedLimitWrite > 0 {
				diskDevice.MaxWriteSpeedMbps = &speedLimitWrite
			}

			if speedLimitWriteBurstable > 0 {
				diskDevice.BurstableWriteSpeedMbps = &speedLimitWriteBurstable
			}
		}

		baseDiskInterface := diskDigitPrefix(diskInterface)

		if baseDiskInterface != "virtio" && baseDiskInterface != "scsi" &&
			baseDiskInterface != "sata" {
			errorMsg := fmt.Sprintf(
				"Defined disk interface not supported. Interface was %s, but only virtio, sata and scsi are supported",
				diskInterface,
			)
			return diskDeviceObjects, errors.New(errorMsg)
		}

		if _, present := diskDeviceObjects[baseDiskInterface]; !present {
			diskDeviceObjects[baseDiskInterface] = map[string]vms.CustomStorageDevice{}
		}

		diskDeviceObjects[baseDiskInterface][diskInterface] = diskDevice
	}

	return diskDeviceObjects, nil
}

func vmGetEfiDisk(d *schema.ResourceData, disk []interface{}) *vms.CustomEFIDisk {
	var efiDisk []interface{}

	if disk != nil {
		efiDisk = disk
	} else {
		efiDisk = d.Get(mkResourceVirtualEnvironmentVMEFIDisk).([]interface{})
	}

	var efiDiskConfig *vms.CustomEFIDisk

	if len(efiDisk) > 0 {
		efiDiskConfig = &vms.CustomEFIDisk{}

		block := efiDisk[0].(map[string]interface{})
		datastoreID, _ := block[mkResourceVirtualEnvironmentVMEFIDiskDatastoreID].(string)
		fileFormat, _ := block[mkResourceVirtualEnvironmentVMEFIDiskFileFormat].(string)
		efiType, _ := block[mkResourceVirtualEnvironmentVMEFIDiskType].(string)
		preEnrolledKeys := types.CustomBool(block[mkResourceVirtualEnvironmentVMEFIDiskPreEnrolledKeys].(bool))

		// special case for efi disk, the size is ignored, see docs for more info
		efiDiskConfig.FileVolume = fmt.Sprintf("%s:1", datastoreID)
		efiDiskConfig.Format = &fileFormat
		efiDiskConfig.Type = &efiType
		efiDiskConfig.PreEnrolledKeys = &preEnrolledKeys
	}

	return efiDiskConfig
}

func vmGetEfiDiskAsStorageDevice(d *schema.ResourceData, disk []interface{}) (*vms.CustomStorageDevice, error) {
	efiDisk := vmGetEfiDisk(d, disk)

	var storageDevice *vms.CustomStorageDevice

	if efiDisk != nil {
		id := "0"
		baseDiskInterface := "efidisk"
		diskInterface := fmt.Sprint(baseDiskInterface, id)

		storageDevice = &vms.CustomStorageDevice{
			Enabled:    true,
			FileVolume: efiDisk.FileVolume,
			Format:     efiDisk.Format,
			Interface:  &diskInterface,
			ID:         &id,
		}

		if efiDisk.Type != nil {
			ds, err := types.ParseDiskSize(*efiDisk.Type)
			if err != nil {
				return nil, fmt.Errorf("invalid efi disk type: %s", err.Error())
			}

			sizeInt := ds.InMegabytes()
			storageDevice.Size = &ds
			storageDevice.SizeInt = &sizeInt
		}
	}

	return storageDevice, nil
}

func vmGetHostPCIDeviceObjects(d *schema.ResourceData) vms.CustomPCIDevices {
	pciDevice := d.Get(mkResourceVirtualEnvironmentVMHostPCI).([]interface{})
	pciDeviceObjects := make(vms.CustomPCIDevices, len(pciDevice))

	for i, pciDeviceEntry := range pciDevice {
		block := pciDeviceEntry.(map[string]interface{})

		ids, _ := block[mkResourceVirtualEnvironmentVMHostPCIDeviceID].(string)
		mdev, _ := block[mkResourceVirtualEnvironmentVMHostPCIDeviceMDev].(string)
		pcie := types.CustomBool(block[mkResourceVirtualEnvironmentVMHostPCIDevicePCIE].(bool))
		rombar := types.CustomBool(
			block[mkResourceVirtualEnvironmentVMHostPCIDeviceROMBAR].(bool),
		)
		romfile, _ := block[mkResourceVirtualEnvironmentVMHostPCIDeviceROMFile].(string)
		xvga := types.CustomBool(block[mkResourceVirtualEnvironmentVMHostPCIDeviceXVGA].(bool))

		device := vms.CustomPCIDevice{
			DeviceIDs:  strings.Split(ids, ";"),
			PCIExpress: &pcie,
			ROMBAR:     &rombar,
			XVGA:       &xvga,
		}
		if ids != "" {
			device.DeviceIDs = strings.Split(ids, ";")
		}

		if mdev != "" {
			device.MDev = &mdev
		}

		if romfile != "" {
			device.ROMFile = &romfile
		}

		pciDeviceObjects[i] = device
	}

	return pciDeviceObjects
}

func vmGetNetworkDeviceObjects(d *schema.ResourceData) vms.CustomNetworkDevices {
	networkDevice := d.Get(mkResourceVirtualEnvironmentVMNetworkDevice).([]interface{})
	networkDeviceObjects := make(vms.CustomNetworkDevices, len(networkDevice))

	for i, networkDeviceEntry := range networkDevice {
		block := networkDeviceEntry.(map[string]interface{})

		bridge := block[mkResourceVirtualEnvironmentVMNetworkDeviceBridge].(string)
		enabled := block[mkResourceVirtualEnvironmentVMNetworkDeviceEnabled].(bool)
		firewall := types.CustomBool(block[mkResourceVirtualEnvironmentVMNetworkDeviceFirewall].(bool))
		macAddress := block[mkResourceVirtualEnvironmentVMNetworkDeviceMACAddress].(string)
		model := block[mkResourceVirtualEnvironmentVMNetworkDeviceModel].(string)
		rateLimit := block[mkResourceVirtualEnvironmentVMNetworkDeviceRateLimit].(float64)
		vlanID := block[mkResourceVirtualEnvironmentVMNetworkDeviceVLANID].(int)
		mtu := block[mkResourceVirtualEnvironmentVMNetworkDeviceMTU].(int)

		device := vms.CustomNetworkDevice{
			Enabled:  enabled,
			Firewall: &firewall,
			Model:    model,
		}

		if bridge != "" {
			device.Bridge = &bridge
		}

		if macAddress != "" {
			device.MACAddress = &macAddress
		}

		if rateLimit != 0 {
			device.RateLimit = &rateLimit
		}

		if vlanID != 0 {
			device.Tag = &vlanID
		}

		if mtu != 0 {
			device.MTU = &mtu
		}

		networkDeviceObjects[i] = device
	}

	return networkDeviceObjects
}

func vmGetOperatingSystemTypeValidator() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.StringInSlice([]string{
		"l24",
		"l26",
		"other",
		"solaris",
		"w2k",
		"w2k3",
		"w2k8",
		"win7",
		"win8",
		"win10",
		"wvista",
		"wxp",
	}, false))
}

func vmGetSerialDeviceList(d *schema.ResourceData) vms.CustomSerialDevices {
	device := d.Get(mkResourceVirtualEnvironmentVMSerialDevice).([]interface{})
	list := make(vms.CustomSerialDevices, len(device))

	for i, v := range device {
		block := v.(map[string]interface{})

		device, _ := block[mkResourceVirtualEnvironmentVMSerialDeviceDevice].(string)

		list[i] = device
	}

	return list
}

func vmGetTagsString(d *schema.ResourceData) string {
	tags := d.Get(mkResourceVirtualEnvironmentVMTags).([]interface{})
	var sanitizedTags []string
	for i := 0; i < len(tags); i++ {
		tag := strings.TrimSpace(tags[i].(string))
		if len(tag) > 0 {
			sanitizedTags = append(sanitizedTags, tag)
		}
	}
	sort.Strings(sanitizedTags)
	return strings.Join(sanitizedTags, ";")
}

func vmGetSerialDeviceValidator() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(func(i interface{}, k string) (s []string, es []error) {
		v, ok := i.(string)

		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be string", k))
			return
		}

		if !strings.HasPrefix(v, "/dev/") && v != "socket" {
			es = append(es, fmt.Errorf("expected %s to be '/dev/*' or 'socket'", k))
			return
		}

		return
	})
}

func vmGetVGADeviceObject(d *schema.ResourceData) (*vms.CustomVGADevice, error) {
	resource := VM()

	vgaBlock, err := getSchemaBlock(
		resource,
		d,
		[]string{mkResourceVirtualEnvironmentVMVGA},
		0,
		true,
	)
	if err != nil {
		return nil, err
	}

	vgaEnabled := types.CustomBool(vgaBlock[mkResourceVirtualEnvironmentVMAgentEnabled].(bool))
	vgaMemory := vgaBlock[mkResourceVirtualEnvironmentVMVGAMemory].(int)
	vgaType := vgaBlock[mkResourceVirtualEnvironmentVMVGAType].(string)

	vgaDevice := &vms.CustomVGADevice{}

	if vgaEnabled {
		if vgaMemory > 0 {
			vgaDevice.Memory = &vgaMemory
		}

		vgaDevice.Type = &vgaType
	} else {
		vgaType = "none"

		vgaDevice = &vms.CustomVGADevice{
			Type: &vgaType,
		}
	}

	return vgaDevice, nil
}

func vmRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(proxmoxtf.ProviderConfiguration)
	api, err := config.GetClient()
	if err != nil {
		return diag.FromErr(err)
	}

	vmID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	vmNodeName, err := api.Cluster().GetVMNodeName(ctx, vmID)
	if err != nil {
		if errors.Is(err, cluster.ErrVMDoesNotExist) {
			d.SetId("")

			return nil
		}

		return diag.FromErr(err)
	}

	if vmNodeName != d.Get(mkResourceVirtualEnvironmentVMNodeName) {
		err = d.Set(mkResourceVirtualEnvironmentVMNodeName, vmNodeName)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	nodeName := d.Get(mkResourceVirtualEnvironmentVMNodeName).(string)

	vmAPI := api.Node(nodeName).VM(vmID)

	// Retrieve the entire configuration in order to compare it to the state.
	vmConfig, err := vmAPI.GetVM(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "HTTP 404") ||
			(strings.Contains(err.Error(), "HTTP 500") && strings.Contains(err.Error(), "does not exist")) {
			d.SetId("")

			return nil
		}

		return diag.FromErr(err)
	}

	vmStatus, err := vmAPI.GetVMStatus(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	return vmReadCustom(ctx, d, m, vmID, vmConfig, vmStatus)
}

func vmReadCustom(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
	vmID int,
	vmConfig *vms.GetResponseData,
	vmStatus *vms.GetStatusResponseData,
) diag.Diagnostics {
	config := m.(proxmoxtf.ProviderConfiguration)
	api, err := config.GetClient()
	if err != nil {
		return diag.FromErr(err)
	}

	diags := vmReadPrimitiveValues(d, vmConfig, vmStatus)
	if diags.HasError() {
		return diags
	}

	// Fix terraform.tfstate, by replacing '-1' (the old default value) with actual vm_id value
	if storedVMID := d.Get(mkResourceVirtualEnvironmentVMVMID).(int); storedVMID == -1 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary: fmt.Sprintf("VM %s has stored legacy vm_id %d, setting vm_id to its correct value %d.",
				d.Id(), storedVMID, vmID),
		})

		err = d.Set(mkResourceVirtualEnvironmentVMVMID, vmID)
		diags = append(diags, diag.FromErr(err)...)
	}

	nodeName := d.Get(mkResourceVirtualEnvironmentVMNodeName).(string)
	clone := d.Get(mkResourceVirtualEnvironmentVMClone).([]interface{})

	// Compare the agent configuration to the one stored in the state.
	currentAgent := d.Get(mkResourceVirtualEnvironmentVMAgent).([]interface{})

	if len(clone) == 0 || len(currentAgent) > 0 {
		if vmConfig.Agent != nil {
			agent := map[string]interface{}{}

			if vmConfig.Agent.Enabled != nil {
				agent[mkResourceVirtualEnvironmentVMAgentEnabled] = bool(*vmConfig.Agent.Enabled)
			} else {
				agent[mkResourceVirtualEnvironmentVMAgentEnabled] = false
			}

			if vmConfig.Agent.TrimClonedDisks != nil {
				agent[mkResourceVirtualEnvironmentVMAgentTrim] = bool(
					*vmConfig.Agent.TrimClonedDisks,
				)
			} else {
				agent[mkResourceVirtualEnvironmentVMAgentTrim] = false
			}

			if len(currentAgent) > 0 {
				currentAgentBlock := currentAgent[0].(map[string]interface{})
				currentAgentTimeout := currentAgentBlock[mkResourceVirtualEnvironmentVMAgentTimeout].(string)

				if currentAgentTimeout != "" {
					agent[mkResourceVirtualEnvironmentVMAgentTimeout] = currentAgentTimeout
				} else {
					agent[mkResourceVirtualEnvironmentVMAgentTimeout] = dvResourceVirtualEnvironmentVMAgentTimeout
				}
			} else {
				agent[mkResourceVirtualEnvironmentVMAgentTimeout] = dvResourceVirtualEnvironmentVMAgentTimeout
			}

			if vmConfig.Agent.Type != nil {
				agent[mkResourceVirtualEnvironmentVMAgentType] = *vmConfig.Agent.Type
			} else {
				agent[mkResourceVirtualEnvironmentVMAgentType] = ""
			}

			if len(clone) > 0 {
				if len(currentAgent) > 0 {
					err := d.Set(mkResourceVirtualEnvironmentVMAgent, []interface{}{agent})
					diags = append(diags, diag.FromErr(err)...)
				}
			} else if len(currentAgent) > 0 ||
				agent[mkResourceVirtualEnvironmentVMAgentEnabled] != dvResourceVirtualEnvironmentVMAgentEnabled ||
				agent[mkResourceVirtualEnvironmentVMAgentTimeout] != dvResourceVirtualEnvironmentVMAgentTimeout ||
				agent[mkResourceVirtualEnvironmentVMAgentTrim] != dvResourceVirtualEnvironmentVMAgentTrim ||
				agent[mkResourceVirtualEnvironmentVMAgentType] != dvResourceVirtualEnvironmentVMAgentType {
				err := d.Set(mkResourceVirtualEnvironmentVMAgent, []interface{}{agent})
				diags = append(diags, diag.FromErr(err)...)
			}
		} else if len(clone) > 0 {
			if len(currentAgent) > 0 {
				err := d.Set(mkResourceVirtualEnvironmentVMAgent, []interface{}{})
				diags = append(diags, diag.FromErr(err)...)
			}
		} else {
			err := d.Set(mkResourceVirtualEnvironmentVMAgent, []interface{}{})
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	// Compare the audio devices to those stored in the state.
	currentAudioDevice := d.Get(mkResourceVirtualEnvironmentVMAudioDevice).([]interface{})

	audioDevices := make([]interface{}, 1)
	audioDevicesArray := []*vms.CustomAudioDevice{
		vmConfig.AudioDevice,
	}
	audioDevicesCount := 0

	for adi, ad := range audioDevicesArray {
		m := map[string]interface{}{}

		if ad != nil {
			m[mkResourceVirtualEnvironmentVMAudioDeviceDevice] = ad.Device

			if ad.Driver != nil {
				m[mkResourceVirtualEnvironmentVMAudioDeviceDriver] = *ad.Driver
			} else {
				m[mkResourceVirtualEnvironmentVMAudioDeviceDriver] = ""
			}

			m[mkResourceVirtualEnvironmentVMAudioDeviceEnabled] = true

			audioDevicesCount = adi + 1
		} else {
			m[mkResourceVirtualEnvironmentVMAudioDeviceDevice] = ""
			m[mkResourceVirtualEnvironmentVMAudioDeviceDriver] = ""
			m[mkResourceVirtualEnvironmentVMAudioDeviceEnabled] = false
		}

		audioDevices[adi] = m
	}

	if len(clone) == 0 || len(currentAudioDevice) > 0 {
		err := d.Set(mkResourceVirtualEnvironmentVMAudioDevice, audioDevices[:audioDevicesCount])
		diags = append(diags, diag.FromErr(err)...)
	}

	// Compare the IDE devices to the CD-ROM and cloud-init configurations stored in the state.
	if vmConfig.IDEDevice3 != nil {
		cdrom := make([]interface{}, 1)
		cdromBlock := map[string]interface{}{}
		currentCDROM := d.Get(mkResourceVirtualEnvironmentVMCDROM).([]interface{})

		if len(clone) == 0 || len(currentCDROM) > 0 {
			cdromBlock[mkResourceVirtualEnvironmentVMCDROMEnabled] = vmConfig.IDEDevice3.Enabled
			cdromBlock[mkResourceVirtualEnvironmentVMCDROMFileID] = vmConfig.IDEDevice3.FileVolume

			if len(currentCDROM) > 0 {
				isCurrentCDROMFileId := currentCDROM[0].(map[string]interface{})

				if isCurrentCDROMFileId[mkResourceVirtualEnvironmentVMCDROMFileID] == "" {
					cdromBlock[mkResourceVirtualEnvironmentVMCDROMFileID] = ""
				}

				if isCurrentCDROMFileId[mkResourceVirtualEnvironmentVMCDROMEnabled] == false {
					cdromBlock[mkResourceVirtualEnvironmentVMCDROMEnabled] = false
				}
			}

			cdrom[0] = cdromBlock

			err := d.Set(mkResourceVirtualEnvironmentVMCDROM, cdrom)
			diags = append(diags, diag.FromErr(err)...)
		}

	} else {
		err := d.Set(mkResourceVirtualEnvironmentVMCDROM, []interface{}{})
		diags = append(diags, diag.FromErr(err)...)
	}

	// Compare the CPU configuration to the one stored in the state.
	cpu := map[string]interface{}{}

	if vmConfig.CPUArchitecture != nil {
		cpu[mkResourceVirtualEnvironmentVMCPUArchitecture] = *vmConfig.CPUArchitecture
	} else {
		// Default value of "arch" is "" according to the API documentation.
		// However, assume the provider's default value as a workaround when the root account is not being used.
		if !api.API().IsRootTicket() {
			cpu[mkResourceVirtualEnvironmentVMCPUArchitecture] = dvResourceVirtualEnvironmentVMCPUArchitecture
		} else {
			cpu[mkResourceVirtualEnvironmentVMCPUArchitecture] = ""
		}
	}

	if vmConfig.CPUCores != nil {
		cpu[mkResourceVirtualEnvironmentVMCPUCores] = *vmConfig.CPUCores
	} else {
		// Default value of "cores" is "1" according to the API documentation.
		cpu[mkResourceVirtualEnvironmentVMCPUCores] = 1
	}

	if vmConfig.VirtualCPUCount != nil {
		cpu[mkResourceVirtualEnvironmentVMCPUHotplugged] = *vmConfig.VirtualCPUCount
	} else {
		// Default value of "vcpus" is "1" according to the API documentation.
		cpu[mkResourceVirtualEnvironmentVMCPUHotplugged] = 0
	}

	if vmConfig.NUMAEnabled != nil {
		cpu[mkResourceVirtualEnvironmentVMCPUNUMA] = *vmConfig.NUMAEnabled
	} else {
		// Default value of "numa" is "false" according to the API documentation.
		cpu[mkResourceVirtualEnvironmentVMCPUNUMA] = false
	}

	if vmConfig.CPUSockets != nil {
		cpu[mkResourceVirtualEnvironmentVMCPUSockets] = *vmConfig.CPUSockets
	} else {
		// Default value of "sockets" is "1" according to the API documentation.
		cpu[mkResourceVirtualEnvironmentVMCPUSockets] = 1
	}

	if vmConfig.CPUEmulation != nil {
		if vmConfig.CPUEmulation.Flags != nil {
			convertedFlags := make([]interface{}, len(*vmConfig.CPUEmulation.Flags))

			for fi, fv := range *vmConfig.CPUEmulation.Flags {
				convertedFlags[fi] = fv
			}

			cpu[mkResourceVirtualEnvironmentVMCPUFlags] = convertedFlags
		} else {
			cpu[mkResourceVirtualEnvironmentVMCPUFlags] = []interface{}{}
		}

		cpu[mkResourceVirtualEnvironmentVMCPUType] = vmConfig.CPUEmulation.Type
	} else {
		cpu[mkResourceVirtualEnvironmentVMCPUFlags] = []interface{}{}
		// Default value of "cputype" is "qemu64" according to the QEMU documentation.
		cpu[mkResourceVirtualEnvironmentVMCPUType] = "qemu64"
	}

	if vmConfig.CPUUnits != nil {
		cpu[mkResourceVirtualEnvironmentVMCPUUnits] = *vmConfig.CPUUnits
	} else {
		// Default value of "cpuunits" is "1024" according to the API documentation.
		cpu[mkResourceVirtualEnvironmentVMCPUUnits] = 1024
	}

	currentCPU := d.Get(mkResourceVirtualEnvironmentVMCPU).([]interface{})

	if len(clone) > 0 {
		if len(currentCPU) > 0 {
			err := d.Set(mkResourceVirtualEnvironmentVMCPU, []interface{}{cpu})
			diags = append(diags, diag.FromErr(err)...)
		}
	} else if len(currentCPU) > 0 ||
		cpu[mkResourceVirtualEnvironmentVMCPUArchitecture] != dvResourceVirtualEnvironmentVMCPUArchitecture ||
		cpu[mkResourceVirtualEnvironmentVMCPUCores] != dvResourceVirtualEnvironmentVMCPUCores ||
		len(cpu[mkResourceVirtualEnvironmentVMCPUFlags].([]interface{})) > 0 ||
		cpu[mkResourceVirtualEnvironmentVMCPUHotplugged] != dvResourceVirtualEnvironmentVMCPUHotplugged ||
		cpu[mkResourceVirtualEnvironmentVMCPUSockets] != dvResourceVirtualEnvironmentVMCPUSockets ||
		cpu[mkResourceVirtualEnvironmentVMCPUType] != dvResourceVirtualEnvironmentVMCPUType ||
		cpu[mkResourceVirtualEnvironmentVMCPUUnits] != dvResourceVirtualEnvironmentVMCPUUnits {
		err := d.Set(mkResourceVirtualEnvironmentVMCPU, []interface{}{cpu})
		diags = append(diags, diag.FromErr(err)...)
	}

	currentDiskList := d.Get(mkResourceVirtualEnvironmentVMDisk).([]interface{})
	diskMap := map[string]interface{}{}
	diskObjects := getDiskInfo(vmConfig, d)
	var orderedDiskList []interface{}

	for di, dd := range diskObjects {
		disk := map[string]interface{}{}

		if dd == nil || strings.HasPrefix(di, "ide") {
			continue
		}

		fileIDParts := strings.Split(dd.FileVolume, ":")

		disk[mkResourceVirtualEnvironmentVMDiskDatastoreID] = fileIDParts[0]

		if dd.Format == nil {
			disk[mkResourceVirtualEnvironmentVMDiskFileFormat] = dvResourceVirtualEnvironmentVMDiskFileFormat
			// disk format may not be returned by config API if it is default for the storage, and that may be different
			// from the default qcow2, so we need to read it from the storage API to make sure we have the correct value
			files, err := api.Node(nodeName).ListDatastoreFiles(ctx, fileIDParts[0])
			if err != nil {
				diags = append(diags, diag.FromErr(err)...)
				continue
			}

			for _, v := range files {
				if v.VolumeID == dd.FileVolume {
					disk[mkResourceVirtualEnvironmentVMDiskFileFormat] = v.FileFormat
					break
				}
			}
		} else {
			disk[mkResourceVirtualEnvironmentVMDiskFileFormat] = dd.Format
		}

		if dd.FileID != nil {
			disk[mkResourceVirtualEnvironmentVMDiskFileID] = dd.FileID
		}

		disk[mkResourceVirtualEnvironmentVMDiskInterface] = di
		disk[mkResourceVirtualEnvironmentVMDiskSize] = dd.Size.InGigabytes()

		if dd.BurstableReadSpeedMbps != nil ||
			dd.BurstableWriteSpeedMbps != nil ||
			dd.MaxReadSpeedMbps != nil ||
			dd.MaxWriteSpeedMbps != nil {
			speed := map[string]interface{}{}

			if dd.MaxReadSpeedMbps != nil {
				speed[mkResourceVirtualEnvironmentVMDiskSpeedRead] = *dd.MaxReadSpeedMbps
			} else {
				speed[mkResourceVirtualEnvironmentVMDiskSpeedRead] = 0
			}

			if dd.BurstableReadSpeedMbps != nil {
				speed[mkResourceVirtualEnvironmentVMDiskSpeedReadBurstable] = *dd.BurstableReadSpeedMbps
			} else {
				speed[mkResourceVirtualEnvironmentVMDiskSpeedReadBurstable] = 0
			}

			if dd.MaxWriteSpeedMbps != nil {
				speed[mkResourceVirtualEnvironmentVMDiskSpeedWrite] = *dd.MaxWriteSpeedMbps
			} else {
				speed[mkResourceVirtualEnvironmentVMDiskSpeedWrite] = 0
			}

			if dd.BurstableWriteSpeedMbps != nil {
				speed[mkResourceVirtualEnvironmentVMDiskSpeedWriteBurstable] = *dd.BurstableWriteSpeedMbps
			} else {
				speed[mkResourceVirtualEnvironmentVMDiskSpeedWriteBurstable] = 0
			}

			disk[mkResourceVirtualEnvironmentVMDiskSpeed] = []interface{}{speed}
		} else {
			disk[mkResourceVirtualEnvironmentVMDiskSpeed] = []interface{}{}
		}

		if dd.IOThread != nil {
			disk[mkResourceVirtualEnvironmentVMDiskIOThread] = *dd.IOThread
		} else {
			disk[mkResourceVirtualEnvironmentVMDiskIOThread] = false
		}

		if dd.SSD != nil {
			disk[mkResourceVirtualEnvironmentVMDiskSSD] = *dd.SSD
		} else {
			disk[mkResourceVirtualEnvironmentVMDiskSSD] = false
		}

		if dd.Discard != nil {
			disk[mkResourceVirtualEnvironmentVMDiskDiscard] = *dd.Discard
		} else {
			disk[mkResourceVirtualEnvironmentVMDiskDiscard] = ""
		}

		diskMap[di] = disk
	}

	var keyList []string

	for key := range diskMap {
		keyList = append(keyList, key)
	}

	sort.Strings(keyList)

	for _, k := range keyList {
		orderedDiskList = append(orderedDiskList, diskMap[k])
	}

	if len(clone) > 0 {
		if len(currentDiskList) > 0 {
			err := d.Set(mkResourceVirtualEnvironmentVMDisk, orderedDiskList)
			diags = append(diags, diag.FromErr(err)...)
		}
	} else if len(currentDiskList) > 0 {
		err := d.Set(mkResourceVirtualEnvironmentVMDisk, orderedDiskList)
		diags = append(diags, diag.FromErr(err)...)
	}

	//nolint:nestif
	if vmConfig.EFIDisk != nil {
		efiDisk := map[string]interface{}{}

		fileIDParts := strings.Split(vmConfig.EFIDisk.FileVolume, ":")

		efiDisk[mkResourceVirtualEnvironmentVMEFIDiskDatastoreID] = fileIDParts[0]

		if vmConfig.EFIDisk.Format != nil {
			efiDisk[mkResourceVirtualEnvironmentVMEFIDiskFileFormat] = *vmConfig.EFIDisk.Format
		} else {
			// disk format may not be returned by config API if it is default for the storage, and that may be different
			// from the default qcow2, so we need to read it from the storage API to make sure we have the correct value
			files, err := api.Node(nodeName).ListDatastoreFiles(ctx, fileIDParts[0])
			if err != nil {
				diags = append(diags, diag.FromErr(err)...)
			} else {
				efiDisk[mkResourceVirtualEnvironmentVMEFIDiskFileFormat] = ""
				for _, v := range files {
					if v.VolumeID == vmConfig.EFIDisk.FileVolume {
						efiDisk[mkResourceVirtualEnvironmentVMEFIDiskFileFormat] = v.FileFormat
						break
					}
				}
			}
		}

		if vmConfig.EFIDisk.Type != nil {
			efiDisk[mkResourceVirtualEnvironmentVMEFIDiskType] = *vmConfig.EFIDisk.Type
		} else {
			efiDisk[mkResourceVirtualEnvironmentVMEFIDiskType] = dvResourceVirtualEnvironmentVMEFIDiskType
		}

		if vmConfig.EFIDisk.PreEnrolledKeys != nil {
			efiDisk[mkResourceVirtualEnvironmentVMEFIDiskPreEnrolledKeys] = *vmConfig.EFIDisk.PreEnrolledKeys
		} else {
			efiDisk[mkResourceVirtualEnvironmentVMEFIDiskPreEnrolledKeys] = false
		}

		currentEfiDisk := d.Get(mkResourceVirtualEnvironmentVMEFIDisk).([]interface{})

		if len(clone) > 0 && len(currentEfiDisk) > 0 {
			err := d.Set(mkResourceVirtualEnvironmentVMEFIDisk, []interface{}{efiDisk})
			diags = append(diags, diag.FromErr(err)...)
		} else if len(currentEfiDisk) > 0 ||
			efiDisk[mkResourceVirtualEnvironmentVMEFIDiskDatastoreID] != dvResourceVirtualEnvironmentVMEFIDiskDatastoreID ||
			efiDisk[mkResourceVirtualEnvironmentVMEFIDiskType] != dvResourceVirtualEnvironmentVMEFIDiskType ||
			efiDisk[mkResourceVirtualEnvironmentVMEFIDiskPreEnrolledKeys] != dvResourceVirtualEnvironmentVMEFIDiskPreEnrolledKeys || //nolint:lll
			efiDisk[mkResourceVirtualEnvironmentVMEFIDiskFileFormat] != dvResourceVirtualEnvironmentVMEFIDiskFileFormat {
			err := d.Set(mkResourceVirtualEnvironmentVMEFIDisk, []interface{}{efiDisk})
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	currentPCIList := d.Get(mkResourceVirtualEnvironmentVMHostPCI).([]interface{})
	pciMap := map[string]interface{}{}
	var orderedPCIList []interface{}

	pciDevices := getPCIInfo(vmConfig, d)
	for pi, pp := range pciDevices {
		if (pp == nil) || (pp.DeviceIDs == nil) {
			continue
		}

		pci := map[string]interface{}{}

		pci[mkResourceVirtualEnvironmentVMHostPCIDevice] = pi
		pci[mkResourceVirtualEnvironmentVMHostPCIDeviceID] = strings.Join(pp.DeviceIDs, ";")

		if pp.MDev != nil {
			pci[mkResourceVirtualEnvironmentVMHostPCIDeviceMDev] = *pp.MDev
		} else {
			pci[mkResourceVirtualEnvironmentVMHostPCIDeviceMDev] = ""
		}

		if pp.PCIExpress != nil {
			pci[mkResourceVirtualEnvironmentVMHostPCIDevicePCIE] = *pp.PCIExpress
		} else {
			pci[mkResourceVirtualEnvironmentVMHostPCIDevicePCIE] = false
		}

		if pp.ROMBAR != nil {
			pci[mkResourceVirtualEnvironmentVMHostPCIDeviceROMBAR] = *pp.ROMBAR
		} else {
			pci[mkResourceVirtualEnvironmentVMHostPCIDeviceROMBAR] = false
		}

		if pp.ROMFile != nil {
			pci[mkResourceVirtualEnvironmentVMHostPCIDeviceROMFile] = *pp.ROMFile
		} else {
			pci[mkResourceVirtualEnvironmentVMHostPCIDeviceROMFile] = ""
		}

		if pp.XVGA != nil {
			pci[mkResourceVirtualEnvironmentVMHostPCIDeviceXVGA] = *pp.XVGA
		} else {
			pci[mkResourceVirtualEnvironmentVMHostPCIDeviceXVGA] = false
		}

		pciMap[pi] = pci
	}

	keyList = []string{}
	for key := range pciMap {
		keyList = append(keyList, key)
	}
	sort.Strings(keyList)

	for _, k := range keyList {
		orderedPCIList = append(orderedPCIList, pciMap[k])
	}

	if len(clone) > 0 {
		if len(currentPCIList) > 0 {
			err := d.Set(mkResourceVirtualEnvironmentVMHostPCI, orderedPCIList)
			diags = append(diags, diag.FromErr(err)...)
		}
	} else if len(currentPCIList) > 0 {
		// todo: reordering of devices by PVE may cause an issue here
		err := d.Set(mkResourceVirtualEnvironmentVMHostPCI, orderedPCIList)
		diags = append(diags, diag.FromErr(err)...)
	}

	// Compare the initialization configuration to the one stored in the state.
	initialization := map[string]interface{}{}

	if vmConfig.IDEDevice2 != nil {
		if *vmConfig.IDEDevice2.Media == "cdrom" {
			if strings.Contains(
				vmConfig.IDEDevice2.FileVolume,
				fmt.Sprintf("vm-%d-cloudinit", vmID),
			) {
				fileVolumeParts := strings.Split(vmConfig.IDEDevice2.FileVolume, ":")
				initialization[mkResourceVirtualEnvironmentVMInitializationDatastoreID] = fileVolumeParts[0]
			}
		}
	}

	if vmConfig.CloudInitDNSDomain != nil || vmConfig.CloudInitDNSServer != nil {
		initializationDNS := map[string]interface{}{}

		if vmConfig.CloudInitDNSDomain != nil {
			initializationDNS[mkResourceVirtualEnvironmentVMInitializationDNSDomain] = *vmConfig.CloudInitDNSDomain
		} else {
			initializationDNS[mkResourceVirtualEnvironmentVMInitializationDNSDomain] = ""
		}

		if vmConfig.CloudInitDNSServer != nil {
			initializationDNS[mkResourceVirtualEnvironmentVMInitializationDNSServer] = *vmConfig.CloudInitDNSServer
		} else {
			initializationDNS[mkResourceVirtualEnvironmentVMInitializationDNSServer] = ""
		}

		initialization[mkResourceVirtualEnvironmentVMInitializationDNS] = []interface{}{
			initializationDNS,
		}
	}

	ipConfigLast := -1
	ipConfigObjects := []*vms.CustomCloudInitIPConfig{
		vmConfig.IPConfig0,
		vmConfig.IPConfig1,
		vmConfig.IPConfig2,
		vmConfig.IPConfig3,
		vmConfig.IPConfig4,
		vmConfig.IPConfig5,
		vmConfig.IPConfig6,
		vmConfig.IPConfig7,
	}
	ipConfigList := make([]interface{}, len(ipConfigObjects))

	for ipConfigIndex, ipConfig := range ipConfigObjects {
		ipConfigItem := map[string]interface{}{}

		if ipConfig != nil {
			ipConfigLast = ipConfigIndex

			if ipConfig.GatewayIPv4 != nil || ipConfig.IPv4 != nil {
				ipv4 := map[string]interface{}{}

				if ipConfig.IPv4 != nil {
					ipv4[mkResourceVirtualEnvironmentVMInitializationIPConfigIPv4Address] = *ipConfig.IPv4
				} else {
					ipv4[mkResourceVirtualEnvironmentVMInitializationIPConfigIPv4Address] = ""
				}

				if ipConfig.GatewayIPv4 != nil {
					ipv4[mkResourceVirtualEnvironmentVMInitializationIPConfigIPv4Gateway] = *ipConfig.GatewayIPv4
				} else {
					ipv4[mkResourceVirtualEnvironmentVMInitializationIPConfigIPv4Gateway] = ""
				}

				ipConfigItem[mkResourceVirtualEnvironmentVMInitializationIPConfigIPv4] = []interface{}{
					ipv4,
				}
			} else {
				ipConfigItem[mkResourceVirtualEnvironmentVMInitializationIPConfigIPv4] = []interface{}{}
			}

			if ipConfig.GatewayIPv6 != nil || ipConfig.IPv6 != nil {
				ipv6 := map[string]interface{}{}

				if ipConfig.IPv6 != nil {
					ipv6[mkResourceVirtualEnvironmentVMInitializationIPConfigIPv6Address] = *ipConfig.IPv6
				} else {
					ipv6[mkResourceVirtualEnvironmentVMInitializationIPConfigIPv6Address] = ""
				}

				if ipConfig.GatewayIPv6 != nil {
					ipv6[mkResourceVirtualEnvironmentVMInitializationIPConfigIPv6Gateway] = *ipConfig.GatewayIPv6
				} else {
					ipv6[mkResourceVirtualEnvironmentVMInitializationIPConfigIPv6Gateway] = ""
				}

				ipConfigItem[mkResourceVirtualEnvironmentVMInitializationIPConfigIPv6] = []interface{}{
					ipv6,
				}
			} else {
				ipConfigItem[mkResourceVirtualEnvironmentVMInitializationIPConfigIPv6] = []interface{}{}
			}
		} else {
			ipConfigItem[mkResourceVirtualEnvironmentVMInitializationIPConfigIPv4] = []interface{}{}
			ipConfigItem[mkResourceVirtualEnvironmentVMInitializationIPConfigIPv6] = []interface{}{}
		}

		ipConfigList[ipConfigIndex] = ipConfigItem
	}

	if ipConfigLast >= 0 {
		initialization[mkResourceVirtualEnvironmentVMInitializationIPConfig] = ipConfigList[:ipConfigLast+1]
	}

	//nolint:nestif
	if vmConfig.CloudInitPassword != nil || vmConfig.CloudInitSSHKeys != nil ||
		vmConfig.CloudInitUsername != nil {
		initializationUserAccount := map[string]interface{}{}

		if vmConfig.CloudInitSSHKeys != nil {
			initializationUserAccount[mkResourceVirtualEnvironmentVMInitializationUserAccountKeys] = []string(
				*vmConfig.CloudInitSSHKeys,
			)
		} else {
			initializationUserAccount[mkResourceVirtualEnvironmentVMInitializationUserAccountKeys] = []string{}
		}

		if vmConfig.CloudInitPassword != nil {
			initializationUserAccount[mkResourceVirtualEnvironmentVMInitializationUserAccountPassword] = *vmConfig.CloudInitPassword
		} else {
			initializationUserAccount[mkResourceVirtualEnvironmentVMInitializationUserAccountPassword] = ""
		}

		if vmConfig.CloudInitUsername != nil {
			initializationUserAccount[mkResourceVirtualEnvironmentVMInitializationUserAccountUsername] = *vmConfig.CloudInitUsername
		} else {
			initializationUserAccount[mkResourceVirtualEnvironmentVMInitializationUserAccountUsername] = ""
		}

		initialization[mkResourceVirtualEnvironmentVMInitializationUserAccount] = []interface{}{
			initializationUserAccount,
		}
	}

	if vmConfig.CloudInitFiles != nil {
		if vmConfig.CloudInitFiles.UserVolume != nil {
			initialization[mkResourceVirtualEnvironmentVMInitializationUserDataFileID] = *vmConfig.CloudInitFiles.UserVolume
		} else {
			initialization[mkResourceVirtualEnvironmentVMInitializationUserDataFileID] = ""
		}

		if vmConfig.CloudInitFiles.VendorVolume != nil {
			initialization[mkResourceVirtualEnvironmentVMInitializationVendorDataFileID] = *vmConfig.CloudInitFiles.VendorVolume
		} else {
			initialization[mkResourceVirtualEnvironmentVMInitializationVendorDataFileID] = ""
		}

		if vmConfig.CloudInitFiles.NetworkVolume != nil {
			initialization[mkResourceVirtualEnvironmentVMInitializationNetworkDataFileID] = *vmConfig.CloudInitFiles.NetworkVolume
		} else {
			initialization[mkResourceVirtualEnvironmentVMInitializationNetworkDataFileID] = ""
		}

		if vmConfig.CloudInitFiles.MetaVolume != nil {
			initialization[mkResourceVirtualEnvironmentVMInitializationMetaDataFileID] = *vmConfig.CloudInitFiles.MetaVolume
		} else {
			initialization[mkResourceVirtualEnvironmentVMInitializationMetaDataFileID] = ""
		}
	} else if len(initialization) > 0 {
		initialization[mkResourceVirtualEnvironmentVMInitializationUserDataFileID] = ""
		initialization[mkResourceVirtualEnvironmentVMInitializationVendorDataFileID] = ""
		initialization[mkResourceVirtualEnvironmentVMInitializationNetworkDataFileID] = ""
		initialization[mkResourceVirtualEnvironmentVMInitializationMetaDataFileID] = ""
	}

	if vmConfig.CloudInitType != nil {
		initialization[mkResourceVirtualEnvironmentVMInitializationType] = *vmConfig.CloudInitType
	} else if len(initialization) > 0 {
		initialization[mkResourceVirtualEnvironmentVMInitializationType] = ""
	}

	currentInitialization := d.Get(mkResourceVirtualEnvironmentVMInitialization).([]interface{})

	if len(clone) > 0 {
		if len(currentInitialization) > 0 {
			if len(initialization) > 0 {
				err := d.Set(
					mkResourceVirtualEnvironmentVMInitialization,
					[]interface{}{initialization},
				)
				diags = append(diags, diag.FromErr(err)...)
			} else {
				err := d.Set(mkResourceVirtualEnvironmentVMInitialization, []interface{}{})
				diags = append(diags, diag.FromErr(err)...)
			}
		}
	} else if len(initialization) > 0 {
		err := d.Set(mkResourceVirtualEnvironmentVMInitialization, []interface{}{initialization})
		diags = append(diags, diag.FromErr(err)...)
	} else {
		err := d.Set(mkResourceVirtualEnvironmentVMInitialization, []interface{}{})
		diags = append(diags, diag.FromErr(err)...)
	}

	// Compare the operating system configuration to the one stored in the state.
	kvmArguments := map[string]interface{}{}

	if vmConfig.KVMArguments != nil {
		kvmArguments[mkResourceVirtualEnvironmentVMKVMArguments] = *vmConfig.KVMArguments
	} else {
		kvmArguments[mkResourceVirtualEnvironmentVMKVMArguments] = ""
	}

	// Compare the memory configuration to the one stored in the state.
	memory := map[string]interface{}{}

	if vmConfig.DedicatedMemory != nil {
		memory[mkResourceVirtualEnvironmentVMMemoryDedicated] = *vmConfig.DedicatedMemory
	} else {
		memory[mkResourceVirtualEnvironmentVMMemoryDedicated] = 0
	}

	if vmConfig.FloatingMemory != nil {
		memory[mkResourceVirtualEnvironmentVMMemoryFloating] = *vmConfig.FloatingMemory
	} else {
		memory[mkResourceVirtualEnvironmentVMMemoryFloating] = 0
	}

	if vmConfig.SharedMemory != nil {
		memory[mkResourceVirtualEnvironmentVMMemoryShared] = vmConfig.SharedMemory.Size
	} else {
		memory[mkResourceVirtualEnvironmentVMMemoryShared] = 0
	}

	currentMemory := d.Get(mkResourceVirtualEnvironmentVMMemory).([]interface{})

	if len(clone) > 0 {
		if len(currentMemory) > 0 {
			err := d.Set(mkResourceVirtualEnvironmentVMMemory, []interface{}{memory})
			diags = append(diags, diag.FromErr(err)...)
		}
	} else if len(currentMemory) > 0 ||
		memory[mkResourceVirtualEnvironmentVMMemoryDedicated] != dvResourceVirtualEnvironmentVMMemoryDedicated ||
		memory[mkResourceVirtualEnvironmentVMMemoryFloating] != dvResourceVirtualEnvironmentVMMemoryFloating ||
		memory[mkResourceVirtualEnvironmentVMMemoryShared] != dvResourceVirtualEnvironmentVMMemoryShared {
		err := d.Set(mkResourceVirtualEnvironmentVMMemory, []interface{}{memory})
		diags = append(diags, diag.FromErr(err)...)
	}

	// Compare the network devices to those stored in the state.
	currentNetworkDeviceList := d.Get(mkResourceVirtualEnvironmentVMNetworkDevice).([]interface{})

	macAddresses := make([]interface{}, 8)
	networkDeviceLast := -1
	networkDeviceList := make([]interface{}, 8)
	networkDeviceObjects := []*vms.CustomNetworkDevice{
		vmConfig.NetworkDevice0,
		vmConfig.NetworkDevice1,
		vmConfig.NetworkDevice2,
		vmConfig.NetworkDevice3,
		vmConfig.NetworkDevice4,
		vmConfig.NetworkDevice5,
		vmConfig.NetworkDevice6,
		vmConfig.NetworkDevice7,
	}

	for ni, nd := range networkDeviceObjects {
		networkDevice := map[string]interface{}{}

		if nd != nil {
			networkDeviceLast = ni

			if nd.Bridge != nil {
				networkDevice[mkResourceVirtualEnvironmentVMNetworkDeviceBridge] = *nd.Bridge
			} else {
				networkDevice[mkResourceVirtualEnvironmentVMNetworkDeviceBridge] = ""
			}

			networkDevice[mkResourceVirtualEnvironmentVMNetworkDeviceEnabled] = nd.Enabled

			if nd.Firewall != nil {
				networkDevice[mkResourceVirtualEnvironmentVMNetworkDeviceFirewall] = *nd.Firewall
			} else {
				networkDevice[mkResourceVirtualEnvironmentVMNetworkDeviceFirewall] = false
			}

			if nd.MACAddress != nil {
				macAddresses[ni] = *nd.MACAddress
			} else {
				macAddresses[ni] = ""
			}

			networkDevice[mkResourceVirtualEnvironmentVMNetworkDeviceMACAddress] = macAddresses[ni]
			networkDevice[mkResourceVirtualEnvironmentVMNetworkDeviceModel] = nd.Model

			if nd.RateLimit != nil {
				networkDevice[mkResourceVirtualEnvironmentVMNetworkDeviceRateLimit] = *nd.RateLimit
			} else {
				networkDevice[mkResourceVirtualEnvironmentVMNetworkDeviceRateLimit] = 0
			}

			if nd.Tag != nil {
				networkDevice[mkResourceVirtualEnvironmentVMNetworkDeviceVLANID] = nd.Tag
			} else {
				networkDevice[mkResourceVirtualEnvironmentVMNetworkDeviceVLANID] = 0
			}
			if nd.MTU != nil {
				networkDevice[mkResourceVirtualEnvironmentVMNetworkDeviceMTU] = nd.MTU
			} else {
				networkDevice[mkResourceVirtualEnvironmentVMNetworkDeviceMTU] = 0
			}
		} else {
			macAddresses[ni] = ""
			networkDevice[mkResourceVirtualEnvironmentVMNetworkDeviceEnabled] = false
		}

		networkDeviceList[ni] = networkDevice
	}

	if len(clone) > 0 {
		if len(currentNetworkDeviceList) > 0 {
			err := d.Set(
				mkResourceVirtualEnvironmentVMMACAddresses,
				macAddresses[0:len(currentNetworkDeviceList)],
			)
			diags = append(diags, diag.FromErr(err)...)
			err = d.Set(
				mkResourceVirtualEnvironmentVMNetworkDevice,
				networkDeviceList[:networkDeviceLast+1],
			)
			diags = append(diags, diag.FromErr(err)...)
		}
	} else {
		err := d.Set(mkResourceVirtualEnvironmentVMMACAddresses, macAddresses[0:len(currentNetworkDeviceList)])
		diags = append(diags, diag.FromErr(err)...)

		if len(currentNetworkDeviceList) > 0 || networkDeviceLast > -1 {
			err := d.Set(mkResourceVirtualEnvironmentVMNetworkDevice, networkDeviceList[:networkDeviceLast+1])
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	// Compare the operating system configuration to the one stored in the state.
	operatingSystem := map[string]interface{}{}

	if vmConfig.OSType != nil {
		operatingSystem[mkResourceVirtualEnvironmentVMOperatingSystemType] = *vmConfig.OSType
	} else {
		operatingSystem[mkResourceVirtualEnvironmentVMOperatingSystemType] = ""
	}

	currentOperatingSystem := d.Get(mkResourceVirtualEnvironmentVMOperatingSystem).([]interface{})

	if len(clone) > 0 {
		if len(currentOperatingSystem) > 0 {
			err := d.Set(
				mkResourceVirtualEnvironmentVMOperatingSystem,
				[]interface{}{operatingSystem},
			)
			diags = append(diags, diag.FromErr(err)...)
		}
	} else if len(currentOperatingSystem) > 0 ||
		operatingSystem[mkResourceVirtualEnvironmentVMOperatingSystemType] != dvResourceVirtualEnvironmentVMOperatingSystemType {
		err := d.Set(mkResourceVirtualEnvironmentVMOperatingSystem, []interface{}{operatingSystem})
		diags = append(diags, diag.FromErr(err)...)
	} else {
		err := d.Set(mkResourceVirtualEnvironmentVMOperatingSystem, []interface{}{})
		diags = append(diags, diag.FromErr(err)...)
	}

	// Compare the pool ID to the value stored in the state.
	currentPoolID := d.Get(mkResourceVirtualEnvironmentVMPoolID).(string)

	if len(clone) == 0 || currentPoolID != dvResourceVirtualEnvironmentVMPoolID {
		if vmConfig.PoolID != nil {
			err := d.Set(mkResourceVirtualEnvironmentVMPoolID, *vmConfig.PoolID)
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	// Compare the serial devices to those stored in the state.
	serialDevices := make([]interface{}, 4)
	serialDevicesArray := []*string{
		vmConfig.SerialDevice0,
		vmConfig.SerialDevice1,
		vmConfig.SerialDevice2,
		vmConfig.SerialDevice3,
	}
	serialDevicesCount := 0

	for sdi, sd := range serialDevicesArray {
		m := map[string]interface{}{}

		if sd != nil {
			m[mkResourceVirtualEnvironmentVMSerialDeviceDevice] = *sd
			serialDevicesCount = sdi + 1
		} else {
			m[mkResourceVirtualEnvironmentVMSerialDeviceDevice] = ""
		}

		serialDevices[sdi] = m
	}

	currentSerialDevice := d.Get(mkResourceVirtualEnvironmentVMSerialDevice).([]interface{})

	if len(clone) == 0 || len(currentSerialDevice) > 0 {
		err := d.Set(mkResourceVirtualEnvironmentVMSerialDevice, serialDevices[:serialDevicesCount])
		diags = append(diags, diag.FromErr(err)...)
	}

	// Compare the VGA configuration to the one stored in the state.
	vga := map[string]interface{}{}

	if vmConfig.VGADevice != nil {
		vgaEnabled := true

		if vmConfig.VGADevice.Type != nil {
			vgaEnabled = *vmConfig.VGADevice.Type != "none"
		}

		vga[mkResourceVirtualEnvironmentVMVGAEnabled] = vgaEnabled

		if vmConfig.VGADevice.Memory != nil {
			vga[mkResourceVirtualEnvironmentVMVGAMemory] = *vmConfig.VGADevice.Memory
		} else {
			vga[mkResourceVirtualEnvironmentVMVGAMemory] = 0
		}

		if vgaEnabled {
			if vmConfig.VGADevice.Type != nil {
				vga[mkResourceVirtualEnvironmentVMVGAType] = *vmConfig.VGADevice.Type
			} else {
				vga[mkResourceVirtualEnvironmentVMVGAType] = ""
			}
		}
	} else {
		vga[mkResourceVirtualEnvironmentVMVGAEnabled] = true
		vga[mkResourceVirtualEnvironmentVMVGAMemory] = 0
		vga[mkResourceVirtualEnvironmentVMVGAType] = ""
	}

	currentVGA := d.Get(mkResourceVirtualEnvironmentVMVGA).([]interface{})

	if len(clone) > 0 {
		if len(currentVGA) > 0 {
			err := d.Set(mkResourceVirtualEnvironmentVMVGA, []interface{}{vga})
			diags = append(diags, diag.FromErr(err)...)
		}
	} else if len(currentVGA) > 0 ||
		vga[mkResourceVirtualEnvironmentVMVGAEnabled] != dvResourceVirtualEnvironmentVMVGAEnabled ||
		vga[mkResourceVirtualEnvironmentVMVGAMemory] != dvResourceVirtualEnvironmentVMVGAMemory ||
		vga[mkResourceVirtualEnvironmentVMVGAType] != dvResourceVirtualEnvironmentVMVGAType {
		err := d.Set(mkResourceVirtualEnvironmentVMVGA, []interface{}{vga})
		diags = append(diags, diag.FromErr(err)...)
	} else {
		err := d.Set(mkResourceVirtualEnvironmentVMVGA, []interface{}{})
		diags = append(diags, diag.FromErr(err)...)
	}

	// Compare SCSI hardware type
	scsiHardware := d.Get(mkResourceVirtualEnvironmentVMSCSIHardware).(string)

	if len(clone) == 0 || scsiHardware != dvResourceVirtualEnvironmentVMSCSIHardware {
		if vmConfig.SCSIHardware != nil {
			err := d.Set(mkResourceVirtualEnvironmentVMSCSIHardware, *vmConfig.SCSIHardware)
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	diags = append(
		diags,
		vmReadNetworkValues(ctx, d, m, vmID, vmConfig)...)

	return diags
}

func vmReadNetworkValues(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
	vmID int,
	vmConfig *vms.GetResponseData,
) diag.Diagnostics {
	var diags diag.Diagnostics

	config := m.(proxmoxtf.ProviderConfiguration)

	api, e := config.GetClient()
	if e != nil {
		return diag.FromErr(e)
	}

	nodeName := d.Get(mkResourceVirtualEnvironmentVMNodeName).(string)

	vmAPI := api.Node(nodeName).VM(vmID)

	started := d.Get(mkResourceVirtualEnvironmentVMStarted).(bool)

	var ipv4Addresses []interface{}
	var ipv6Addresses []interface{}
	var networkInterfaceNames []interface{}

	if started {
		if vmConfig.Agent != nil && vmConfig.Agent.Enabled != nil && *vmConfig.Agent.Enabled {
			resource := VM()
			agentBlock, err := getSchemaBlock(
				resource,
				d,
				[]string{mkResourceVirtualEnvironmentVMAgent},
				0,
				true,
			)
			if err != nil {
				return diag.FromErr(err)
			}

			agentTimeout, err := time.ParseDuration(
				agentBlock[mkResourceVirtualEnvironmentVMAgentTimeout].(string),
			)
			if err != nil {
				return diag.FromErr(err)
			}

			var macAddresses []interface{}

			networkInterfaces, err := vmAPI.WaitForNetworkInterfacesFromVMAgent(ctx, int(agentTimeout.Seconds()), 5, true)
			if err == nil && networkInterfaces.Result != nil {
				ipv4Addresses = make([]interface{}, len(*networkInterfaces.Result))
				ipv6Addresses = make([]interface{}, len(*networkInterfaces.Result))
				macAddresses = make([]interface{}, len(*networkInterfaces.Result))
				networkInterfaceNames = make([]interface{}, len(*networkInterfaces.Result))

				for ri, rv := range *networkInterfaces.Result {
					var rvIPv4Addresses []interface{}
					var rvIPv6Addresses []interface{}

					if rv.IPAddresses != nil {
						for _, ip := range *rv.IPAddresses {
							switch ip.Type {
							case "ipv4":
								rvIPv4Addresses = append(rvIPv4Addresses, ip.Address)
							case "ipv6":
								rvIPv6Addresses = append(rvIPv6Addresses, ip.Address)
							}
						}
					}

					ipv4Addresses[ri] = rvIPv4Addresses
					ipv6Addresses[ri] = rvIPv6Addresses
					macAddresses[ri] = strings.ToUpper(rv.MACAddress)
					networkInterfaceNames[ri] = rv.Name
				}
			}

			err = d.Set(mkResourceVirtualEnvironmentVMMACAddresses, macAddresses)
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	e = d.Set(mkResourceVirtualEnvironmentVMIPv4Addresses, ipv4Addresses)
	diags = append(diags, diag.FromErr(e)...)
	e = d.Set(mkResourceVirtualEnvironmentVMIPv6Addresses, ipv6Addresses)
	diags = append(diags, diag.FromErr(e)...)
	e = d.Set(mkResourceVirtualEnvironmentVMNetworkInterfaceNames, networkInterfaceNames)
	diags = append(diags, diag.FromErr(e)...)

	return diags
}

func vmReadPrimitiveValues(
	d *schema.ResourceData,
	vmConfig *vms.GetResponseData,
	vmStatus *vms.GetStatusResponseData,
) diag.Diagnostics {
	var diags diag.Diagnostics
	var err error

	clone := d.Get(mkResourceVirtualEnvironmentVMClone).([]interface{})
	currentACPI := d.Get(mkResourceVirtualEnvironmentVMACPI).(bool)

	//nolint:gosimple
	if len(clone) == 0 || currentACPI != dvResourceVirtualEnvironmentVMACPI {
		if vmConfig.ACPI != nil {
			err = d.Set(mkResourceVirtualEnvironmentVMACPI, bool(*vmConfig.ACPI))
		} else {
			// Default value of "acpi" is "1" according to the API documentation.
			err = d.Set(mkResourceVirtualEnvironmentVMACPI, true)
		}
		diags = append(diags, diag.FromErr(err)...)
	}

	currentkvmArguments := d.Get(mkResourceVirtualEnvironmentVMKVMArguments).(string)

	if len(clone) == 0 || currentkvmArguments != dvResourceVirtualEnvironmentVMKVMArguments {
		// PVE API returns "args" as " " if it is set to empty.
		if vmConfig.KVMArguments != nil && len(strings.TrimSpace(*vmConfig.KVMArguments)) > 0 {
			err = d.Set(mkResourceVirtualEnvironmentVMKVMArguments, *vmConfig.KVMArguments)
		} else {
			// Default value of "args" is "" according to the API documentation.
			err = d.Set(mkResourceVirtualEnvironmentVMKVMArguments, "")
		}

		diags = append(diags, diag.FromErr(err)...)
	}

	currentBIOS := d.Get(mkResourceVirtualEnvironmentVMBIOS).(string)

	if len(clone) == 0 || currentBIOS != dvResourceVirtualEnvironmentVMBIOS {
		if vmConfig.BIOS != nil {
			err = d.Set(mkResourceVirtualEnvironmentVMBIOS, *vmConfig.BIOS)
		} else {
			// Default value of "bios" is "seabios" according to the API documentation.
			err = d.Set(mkResourceVirtualEnvironmentVMBIOS, "seabios")
		}
		diags = append(diags, diag.FromErr(err)...)
	}

	currentDescription := d.Get(mkResourceVirtualEnvironmentVMDescription).(string)

	if len(clone) == 0 || currentDescription != dvResourceVirtualEnvironmentVMDescription {
		if vmConfig.Description != nil {
			err = d.Set(mkResourceVirtualEnvironmentVMDescription, *vmConfig.Description)
		} else {
			// Default value of "description" is "" according to the API documentation.
			err = d.Set(mkResourceVirtualEnvironmentVMDescription, "")
		}
		diags = append(diags, diag.FromErr(err)...)
	}

	currentTags := d.Get(mkResourceVirtualEnvironmentVMTags).([]interface{})

	if len(clone) == 0 || len(currentTags) > 0 {
		var tags []string
		if vmConfig.Tags != nil {
			for _, tag := range strings.Split(*vmConfig.Tags, ";") {
				t := strings.TrimSpace(tag)
				if len(t) > 0 {
					tags = append(tags, t)
				}
			}
			sort.Strings(tags)
		}
		err = d.Set(mkResourceVirtualEnvironmentVMTags, tags)
		diags = append(diags, diag.FromErr(err)...)
	}

	currentKeyboardLayout := d.Get(mkResourceVirtualEnvironmentVMKeyboardLayout).(string)

	if len(clone) == 0 || currentKeyboardLayout != dvResourceVirtualEnvironmentVMKeyboardLayout {
		if vmConfig.KeyboardLayout != nil {
			err = d.Set(mkResourceVirtualEnvironmentVMKeyboardLayout, *vmConfig.KeyboardLayout)
		} else {
			// Default value of "keyboard" is "" according to the API documentation.
			err = d.Set(mkResourceVirtualEnvironmentVMKeyboardLayout, "")
		}
		diags = append(diags, diag.FromErr(err)...)
	}

	currentMachine := d.Get(mkResourceVirtualEnvironmentVMMachine).(string)

	if len(clone) == 0 || currentMachine != dvResourceVirtualEnvironmentVMMachineType {
		if vmConfig.Machine != nil {
			err = d.Set(mkResourceVirtualEnvironmentVMMachine, *vmConfig.Machine)
		} else {
			err = d.Set(mkResourceVirtualEnvironmentVMMachine, "")
		}
		diags = append(diags, diag.FromErr(err)...)
	}

	currentName := d.Get(mkResourceVirtualEnvironmentVMName).(string)

	if len(clone) == 0 || currentName != dvResourceVirtualEnvironmentVMName {
		if vmConfig.Name != nil {
			err = d.Set(mkResourceVirtualEnvironmentVMName, *vmConfig.Name)
		} else {
			// Default value of "name" is "" according to the API documentation.
			err = d.Set(mkResourceVirtualEnvironmentVMName, "")
		}
		diags = append(diags, diag.FromErr(err)...)
	}

	if !d.Get(mkResourceVirtualEnvironmentVMTemplate).(bool) {
		err = d.Set(mkResourceVirtualEnvironmentVMStarted, vmStatus.Status == "running")
		diags = append(diags, diag.FromErr(err)...)
	}

	currentTabletDevice := d.Get(mkResourceVirtualEnvironmentVMTabletDevice).(bool)

	//nolint:gosimple
	if len(clone) == 0 || currentTabletDevice != dvResourceVirtualEnvironmentVMTabletDevice {
		if vmConfig.TabletDeviceEnabled != nil {
			err = d.Set(
				mkResourceVirtualEnvironmentVMTabletDevice,
				bool(*vmConfig.TabletDeviceEnabled),
			)
		} else {
			// Default value of "tablet" is "1" according to the API documentation.
			err = d.Set(mkResourceVirtualEnvironmentVMTabletDevice, true)
		}
		diags = append(diags, diag.FromErr(err)...)
	}

	currentTemplate := d.Get(mkResourceVirtualEnvironmentVMTemplate).(bool)

	//nolint:gosimple
	if len(clone) == 0 || currentTemplate != dvResourceVirtualEnvironmentVMTemplate {
		if vmConfig.Template != nil {
			err = d.Set(mkResourceVirtualEnvironmentVMTemplate, bool(*vmConfig.Template))
		} else {
			// Default value of "template" is "0" according to the API documentation.
			err = d.Set(mkResourceVirtualEnvironmentVMTemplate, false)
		}
		diags = append(diags, diag.FromErr(err)...)
	}

	return diags
}

func vmUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(proxmoxtf.ProviderConfiguration)

	api, e := config.GetClient()
	if e != nil {
		return diag.FromErr(e)
	}

	nodeName := d.Get(mkResourceVirtualEnvironmentVMNodeName).(string)
	rebootRequired := false

	vmID, e := strconv.Atoi(d.Id())
	if e != nil {
		return diag.FromErr(e)
	}

	vmAPI := api.Node(nodeName).VM(vmID)

	updateBody := &vms.UpdateRequestBody{
		IDEDevices: vms.CustomStorageDevices{
			"ide0": vms.CustomStorageDevice{
				Enabled: false,
			},
			"ide1": vms.CustomStorageDevice{
				Enabled: false,
			},
			"ide2": vms.CustomStorageDevice{
				Enabled: false,
			},
			"ide3": vms.CustomStorageDevice{
				Enabled: false,
			},
		},
	}

	var del []string

	resource := VM()

	// Retrieve the entire configuration as we need to process certain values.
	vmConfig, e := vmAPI.GetVM(ctx)
	if e != nil {
		return diag.FromErr(e)
	}

	// Prepare the new primitive configuration values.
	if d.HasChange(mkResourceVirtualEnvironmentVMACPI) {
		acpi := types.CustomBool(d.Get(mkResourceVirtualEnvironmentVMACPI).(bool))
		updateBody.ACPI = &acpi
		rebootRequired = true
	}

	if d.HasChange(mkResourceVirtualEnvironmentVMKVMArguments) {
		kvmArguments := d.Get(mkResourceVirtualEnvironmentVMKVMArguments).(string)
		updateBody.KVMArguments = &kvmArguments
		rebootRequired = true
	}

	if d.HasChange(mkResourceVirtualEnvironmentVMBIOS) {
		bios := d.Get(mkResourceVirtualEnvironmentVMBIOS).(string)
		updateBody.BIOS = &bios
		rebootRequired = true
	}

	if d.HasChange(mkResourceVirtualEnvironmentVMDescription) {
		description := d.Get(mkResourceVirtualEnvironmentVMDescription).(string)
		updateBody.Description = &description
	}

	if d.HasChange(mkResourceVirtualEnvironmentVMOnBoot) {
		startOnBoot := types.CustomBool(d.Get(mkResourceVirtualEnvironmentVMOnBoot).(bool))
		updateBody.StartOnBoot = &startOnBoot
	}

	if d.HasChange(mkResourceVirtualEnvironmentVMTags) {
		tagString := vmGetTagsString(d)
		updateBody.Tags = &tagString
	}

	if d.HasChange(mkResourceVirtualEnvironmentVMKeyboardLayout) {
		keyboardLayout := d.Get(mkResourceVirtualEnvironmentVMKeyboardLayout).(string)
		updateBody.KeyboardLayout = &keyboardLayout
		rebootRequired = true
	}

	if d.HasChange(mkResourceVirtualEnvironmentVMMachine) {
		machine := d.Get(mkResourceVirtualEnvironmentVMMachine).(string)
		updateBody.Machine = &machine
		rebootRequired = true
	}

	name := d.Get(mkResourceVirtualEnvironmentVMName).(string)

	if name == "" {
		del = append(del, "name")
	} else {
		updateBody.Name = &name
	}

	if d.HasChange(mkResourceVirtualEnvironmentVMTabletDevice) {
		tabletDevice := types.CustomBool(d.Get(mkResourceVirtualEnvironmentVMTabletDevice).(bool))
		updateBody.TabletDeviceEnabled = &tabletDevice
		rebootRequired = true
	}

	template := types.CustomBool(d.Get(mkResourceVirtualEnvironmentVMTemplate).(bool))

	if d.HasChange(mkResourceVirtualEnvironmentVMTemplate) {
		updateBody.Template = &template
		rebootRequired = true
	}

	// Prepare the new agent configuration.
	if d.HasChange(mkResourceVirtualEnvironmentVMAgent) {
		agentBlock, err := getSchemaBlock(
			resource,
			d,
			[]string{mkResourceVirtualEnvironmentVMAgent},
			0,
			true,
		)
		if err != nil {
			return diag.FromErr(err)
		}

		agentEnabled := types.CustomBool(
			agentBlock[mkResourceVirtualEnvironmentVMAgentEnabled].(bool),
		)
		agentTrim := types.CustomBool(agentBlock[mkResourceVirtualEnvironmentVMAgentTrim].(bool))
		agentType := agentBlock[mkResourceVirtualEnvironmentVMAgentType].(string)

		updateBody.Agent = &vms.CustomAgent{
			Enabled:         &agentEnabled,
			TrimClonedDisks: &agentTrim,
			Type:            &agentType,
		}

		rebootRequired = true
	}

	// Prepare the new audio devices.
	if d.HasChange(mkResourceVirtualEnvironmentVMAudioDevice) {
		updateBody.AudioDevices = vmGetAudioDeviceList(d)

		for i := 0; i < len(updateBody.AudioDevices); i++ {
			if !updateBody.AudioDevices[i].Enabled {
				del = append(del, fmt.Sprintf("audio%d", i))
			}
		}

		for i := len(updateBody.AudioDevices); i < maxResourceVirtualEnvironmentVMAudioDevices; i++ {
			del = append(del, fmt.Sprintf("audio%d", i))
		}

		rebootRequired = true
	}

	// Prepare the new boot configuration.
	if d.HasChange(mkResourceVirtualEnvironmentVMBootOrder) {
		bootOrder := d.Get(mkResourceVirtualEnvironmentVMBootOrder).([]interface{})
		bootOrderConverted := make([]string, len(bootOrder))

		for i, device := range bootOrder {
			bootOrderConverted[i] = device.(string)
		}

		updateBody.Boot = &vms.CustomBoot{
			Order: &bootOrderConverted,
		}
		rebootRequired = true
	}

	// Prepare the new CD-ROM configuration.
	if d.HasChange(mkResourceVirtualEnvironmentVMCDROM) {
		cdromBlock, err := getSchemaBlock(
			resource,
			d,
			[]string{mkResourceVirtualEnvironmentVMCDROM},
			0,
			true,
		)
		if err != nil {
			return diag.FromErr(err)
		}

		cdromEnabled := cdromBlock[mkResourceVirtualEnvironmentVMCDROMEnabled].(bool)
		cdromFileID := cdromBlock[mkResourceVirtualEnvironmentVMCDROMFileID].(string)

		if !cdromEnabled && cdromFileID == "" {
			del = append(del, "ide3")
		}

		if cdromFileID == "" {
			cdromFileID = "cdrom"
		}

		cdromMedia := "cdrom"

		updateBody.IDEDevices["ide3"] = vms.CustomStorageDevice{
			Enabled:    cdromEnabled,
			FileVolume: cdromFileID,
			Media:      &cdromMedia,
		}
	}

	// Prepare the new CPU configuration.
	if d.HasChange(mkResourceVirtualEnvironmentVMCPU) {
		cpuBlock, err := getSchemaBlock(
			resource,
			d,
			[]string{mkResourceVirtualEnvironmentVMCPU},
			0,
			true,
		)
		if err != nil {
			return diag.FromErr(err)
		}

		cpuArchitecture := cpuBlock[mkResourceVirtualEnvironmentVMCPUArchitecture].(string)
		cpuCores := cpuBlock[mkResourceVirtualEnvironmentVMCPUCores].(int)
		cpuFlags := cpuBlock[mkResourceVirtualEnvironmentVMCPUFlags].([]interface{})
		cpuHotplugged := cpuBlock[mkResourceVirtualEnvironmentVMCPUHotplugged].(int)
		cpuNUMA := types.CustomBool(cpuBlock[mkResourceVirtualEnvironmentVMCPUNUMA].(bool))
		cpuSockets := cpuBlock[mkResourceVirtualEnvironmentVMCPUSockets].(int)
		cpuType := cpuBlock[mkResourceVirtualEnvironmentVMCPUType].(string)
		cpuUnits := cpuBlock[mkResourceVirtualEnvironmentVMCPUUnits].(int)

		// Only the root account is allowed to change the CPU architecture, which makes this check necessary.
		if api.API().IsRootTicket() ||
			cpuArchitecture != dvResourceVirtualEnvironmentVMCPUArchitecture {
			updateBody.CPUArchitecture = &cpuArchitecture
		}

		updateBody.CPUCores = &cpuCores
		updateBody.CPUSockets = &cpuSockets
		updateBody.CPUUnits = &cpuUnits
		updateBody.NUMAEnabled = &cpuNUMA

		if cpuHotplugged > 0 {
			updateBody.VirtualCPUCount = &cpuHotplugged
		} else {
			del = append(del, "vcpus")
		}

		cpuFlagsConverted := make([]string, len(cpuFlags))

		for fi, flag := range cpuFlags {
			cpuFlagsConverted[fi] = flag.(string)
		}

		updateBody.CPUEmulation = &vms.CustomCPUEmulation{
			Flags: &cpuFlagsConverted,
			Type:  cpuType,
		}

		rebootRequired = true
	}

	// Prepare the new disk device configuration.
	if d.HasChange(mkResourceVirtualEnvironmentVMDisk) {
		diskDeviceObjects, err := vmGetDiskDeviceObjects(d, nil)
		if err != nil {
			return diag.FromErr(err)
		}

		diskDeviceInfo := getDiskInfo(vmConfig, d)

		for prefix, diskMap := range diskDeviceObjects {
			if diskMap == nil {
				continue
			}

			for key, value := range diskMap {
				if diskDeviceInfo[key] == nil {
					return diag.Errorf("missing %s device %s", prefix, key)
				}

				tmp := *diskDeviceInfo[key]
				tmp.BurstableReadSpeedMbps = value.BurstableReadSpeedMbps
				tmp.BurstableWriteSpeedMbps = value.BurstableWriteSpeedMbps
				tmp.MaxReadSpeedMbps = value.MaxReadSpeedMbps
				tmp.MaxWriteSpeedMbps = value.MaxWriteSpeedMbps

				switch prefix {
				case "virtio":
					{
						if updateBody.VirtualIODevices == nil {
							updateBody.VirtualIODevices = vms.CustomStorageDevices{}
						}

						updateBody.VirtualIODevices[key] = tmp
					}
				case "sata":
					{
						if updateBody.SATADevices == nil {
							updateBody.SATADevices = vms.CustomStorageDevices{}
						}

						updateBody.SATADevices[key] = tmp
					}
				case "scsi":
					{
						if updateBody.SCSIDevices == nil {
							updateBody.SCSIDevices = vms.CustomStorageDevices{}
						}

						updateBody.SCSIDevices[key] = tmp
					}
				case "ide":
					{
						// Investigate whether to support IDE mapping.
					}
				default:
					return diag.Errorf("device prefix %s not supported", prefix)
				}
			}
		}
	}

	// Prepare the new efi disk configuration.
	if d.HasChange(mkResourceVirtualEnvironmentVMEFIDisk) {
		efiDisk := vmGetEfiDisk(d, nil)

		updateBody.EFIDisk = efiDisk

		rebootRequired = true
	}

	// Prepare the new cloud-init configuration.
	if d.HasChange(mkResourceVirtualEnvironmentVMInitialization) {
		initializationConfig := vmGetCloudInitConfig(d)

		updateBody.CloudInitConfig = initializationConfig

		if updateBody.CloudInitConfig != nil {
			initialization := d.Get(mkResourceVirtualEnvironmentVMInitialization).([]interface{})
			initializationBlock := initialization[0].(map[string]interface{})
			initializationDatastoreID := initializationBlock[mkResourceVirtualEnvironmentVMInitializationDatastoreID].(string)

			cdromMedia := "cdrom"

			updateBody.IDEDevices["ide2"] = vms.CustomStorageDevice{
				Enabled:    true,
				FileVolume: fmt.Sprintf("%s:cloudinit", initializationDatastoreID),
				Media:      &cdromMedia,
			}

			if vmConfig.IDEDevice2 != nil &&
				strings.Contains(
					vmConfig.IDEDevice2.FileVolume,
					fmt.Sprintf("vm-%d-cloudinit", vmID),
				) {
				tmp := updateBody.IDEDevices["ide2"]
				tmp.Enabled = true
				updateBody.IDEDevices["ide2"] = tmp
			}
		}

		rebootRequired = true
	}

	// Prepare the new hostpci devices configuration.
	if d.HasChange(mkResourceVirtualEnvironmentVMHostPCI) {
		updateBody.PCIDevices = vmGetHostPCIDeviceObjects(d)
		rebootRequired = true
	}

	// Prepare the new memory configuration.
	if d.HasChange(mkResourceVirtualEnvironmentVMMemory) {
		memoryBlock, err := getSchemaBlock(
			resource,
			d,
			[]string{mkResourceVirtualEnvironmentVMMemory},
			0,
			true,
		)
		if err != nil {
			return diag.FromErr(err)
		}

		memoryDedicated := memoryBlock[mkResourceVirtualEnvironmentVMMemoryDedicated].(int)
		memoryFloating := memoryBlock[mkResourceVirtualEnvironmentVMMemoryFloating].(int)
		memoryShared := memoryBlock[mkResourceVirtualEnvironmentVMMemoryShared].(int)

		updateBody.DedicatedMemory = &memoryDedicated
		updateBody.FloatingMemory = &memoryFloating

		if memoryShared > 0 {
			memorySharedName := fmt.Sprintf("vm-%d-ivshmem", vmID)

			updateBody.SharedMemory = &vms.CustomSharedMemory{
				Name: &memorySharedName,
				Size: memoryShared,
			}
		}

		rebootRequired = true
	}

	// Prepare the new network device configuration.
	if d.HasChange(mkResourceVirtualEnvironmentVMNetworkDevice) {
		updateBody.NetworkDevices = vmGetNetworkDeviceObjects(d)

		for i := 0; i < len(updateBody.NetworkDevices); i++ {
			if !updateBody.NetworkDevices[i].Enabled {
				del = append(del, fmt.Sprintf("net%d", i))
			}
		}

		for i := len(updateBody.NetworkDevices); i < maxResourceVirtualEnvironmentVMNetworkDevices; i++ {
			del = append(del, fmt.Sprintf("net%d", i))
		}

		rebootRequired = true
	}

	// Prepare the new operating system configuration.
	if d.HasChange(mkResourceVirtualEnvironmentVMOperatingSystem) {
		operatingSystem, err := getSchemaBlock(
			resource,
			d,
			[]string{mkResourceVirtualEnvironmentVMOperatingSystem},
			0,
			true,
		)
		if err != nil {
			return diag.FromErr(err)
		}

		operatingSystemType := operatingSystem[mkResourceVirtualEnvironmentVMOperatingSystemType].(string)

		updateBody.OSType = &operatingSystemType

		rebootRequired = true
	}

	// Prepare the new serial devices.
	if d.HasChange(mkResourceVirtualEnvironmentVMSerialDevice) {
		updateBody.SerialDevices = vmGetSerialDeviceList(d)

		for i := len(updateBody.SerialDevices); i < maxResourceVirtualEnvironmentVMSerialDevices; i++ {
			del = append(del, fmt.Sprintf("serial%d", i))
		}

		rebootRequired = true
	}

	// Prepare the new VGA configuration.
	if d.HasChange(mkResourceVirtualEnvironmentVMVGA) {
		updateBody.VGADevice, e = vmGetVGADeviceObject(d)
		if e != nil {
			return diag.FromErr(e)
		}

		rebootRequired = true
	}

	// Prepare the new SCSI hardware type
	if d.HasChange(mkResourceVirtualEnvironmentVMSCSIHardware) {
		scsiHardware := d.Get(mkResourceVirtualEnvironmentVMSCSIHardware).(string)
		updateBody.SCSIHardware = &scsiHardware

		rebootRequired = true
	}

	// Update the configuration now that everything has been prepared.
	updateBody.Delete = del

	e = vmAPI.UpdateVM(ctx, updateBody)
	if e != nil {
		return diag.FromErr(e)
	}

	// Determine if the state of the virtual machine state needs to be changed.
	started := d.Get(mkResourceVirtualEnvironmentVMStarted).(bool)

	if d.HasChange(mkResourceVirtualEnvironmentVMStarted) && !bool(template) {
		if started {
			startVMTimeout := d.Get(mkResourceVirtualEnvironmentVMTimeoutStartVM).(int)

			e = vmAPI.StartVM(ctx, startVMTimeout)
			if e != nil {
				return diag.FromErr(e)
			}
		} else {
			forceStop := types.CustomBool(true)
			shutdownTimeout := d.Get(mkResourceVirtualEnvironmentVMTimeoutShutdownVM).(int)

			e = vmAPI.ShutdownVM(ctx, &vms.ShutdownRequestBody{
				ForceStop: &forceStop,
				Timeout:   &shutdownTimeout,
			}, shutdownTimeout+30)
			if e != nil {
				return diag.FromErr(e)
			}

			rebootRequired = false
		}
	}

	// Change the disk locations and/or sizes, if necessary.
	return vmUpdateDiskLocationAndSize(
		ctx,
		d,
		m,
		!bool(template) && rebootRequired,
	)
}

func vmUpdateDiskLocationAndSize(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
	reboot bool,
) diag.Diagnostics {
	config := m.(proxmoxtf.ProviderConfiguration)
	api, err := config.GetClient()
	if err != nil {
		return diag.FromErr(err)
	}

	nodeName := d.Get(mkResourceVirtualEnvironmentVMNodeName).(string)
	started := d.Get(mkResourceVirtualEnvironmentVMStarted).(bool)
	template := d.Get(mkResourceVirtualEnvironmentVMTemplate).(bool)
	vmID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	vmAPI := api.Node(nodeName).VM(vmID)

	// Determine if any of the disks are changing location and/or size, and initiate the necessary actions.
	if d.HasChange(mkResourceVirtualEnvironmentVMDisk) {
		diskOld, diskNew := d.GetChange(mkResourceVirtualEnvironmentVMDisk)

		diskOldEntries, err := vmGetDiskDeviceObjects(
			d,
			diskOld.([]interface{}),
		)
		if err != nil {
			return diag.FromErr(err)
		}

		diskNewEntries, err := vmGetDiskDeviceObjects(
			d,
			diskNew.([]interface{}),
		)
		if err != nil {
			return diag.FromErr(err)
		}

		// Add efidisk if it has changes
		if d.HasChange(mkResourceVirtualEnvironmentVMEFIDisk) {
			diskOld, diskNew := d.GetChange(mkResourceVirtualEnvironmentVMEFIDisk)

			oldEfiDisk, e := vmGetEfiDiskAsStorageDevice(d, diskOld.([]interface{}))
			if e != nil {
				return diag.FromErr(e)
			}

			newEfiDisk, e := vmGetEfiDiskAsStorageDevice(d, diskNew.([]interface{}))
			if e != nil {
				return diag.FromErr(e)
			}

			if oldEfiDisk != nil {
				baseDiskInterface := diskDigitPrefix(*oldEfiDisk.Interface)
				diskOldEntries[baseDiskInterface][*oldEfiDisk.Interface] = *oldEfiDisk
			}

			if newEfiDisk != nil {
				baseDiskInterface := diskDigitPrefix(*newEfiDisk.Interface)
				diskNewEntries[baseDiskInterface][*newEfiDisk.Interface] = *newEfiDisk
			}

			if oldEfiDisk != nil && newEfiDisk != nil && oldEfiDisk.Size != newEfiDisk.Size {
				return diag.Errorf(
					"resizing of efidisks is not supported.",
				)
			}
		}

		var diskMoveBodies []*vms.MoveDiskRequestBody

		var diskResizeBodies []*vms.ResizeDiskRequestBody

		shutdownForDisksRequired := false

		for prefix, diskMap := range diskOldEntries {
			for oldKey, oldDisk := range diskMap {
				if _, present := diskNewEntries[prefix][oldKey]; !present {
					return diag.Errorf(
						"deletion of disks not supported. Please delete disk by hand. Old Interface was %s",
						*oldDisk.Interface,
					)
				}

				if *oldDisk.ID != *diskNewEntries[prefix][oldKey].ID {
					deleteOriginalDisk := types.CustomBool(true)

					diskMoveBodies = append(
						diskMoveBodies,
						&vms.MoveDiskRequestBody{
							DeleteOriginalDisk: &deleteOriginalDisk,
							Disk:               *oldDisk.Interface,
							TargetStorage:      *diskNewEntries[prefix][oldKey].ID,
						},
					)

					// Cannot be done while VM is running.
					shutdownForDisksRequired = true
				}

				if *oldDisk.SizeInt < *diskNewEntries[prefix][oldKey].SizeInt {
					diskResizeBodies = append(
						diskResizeBodies,
						&vms.ResizeDiskRequestBody{
							Disk: *oldDisk.Interface,
							Size: *diskNewEntries[prefix][oldKey].Size,
						},
					)
				}
			}
		}

		if shutdownForDisksRequired && !template {
			forceStop := types.CustomBool(true)
			shutdownTimeout := d.Get(mkResourceVirtualEnvironmentVMTimeoutShutdownVM).(int)

			err = vmAPI.ShutdownVM(
				ctx,
				&vms.ShutdownRequestBody{
					ForceStop: &forceStop,
					Timeout:   &shutdownTimeout,
				},
				shutdownTimeout+30,
			)
			if err != nil {
				return diag.FromErr(err)
			}
		}

		for _, reqBody := range diskMoveBodies {
			moveDiskTimeout := d.Get(mkResourceVirtualEnvironmentVMTimeoutMoveDisk).(int)
			err = vmAPI.MoveVMDisk(ctx, reqBody, moveDiskTimeout)
			if err != nil {
				return diag.FromErr(err)
			}
		}

		for _, reqBody := range diskResizeBodies {
			err = vmAPI.ResizeVMDisk(ctx, reqBody)
			if err != nil {
				return diag.FromErr(err)
			}
		}

		if shutdownForDisksRequired && started && !template {
			startVMTimeout := d.Get(mkResourceVirtualEnvironmentVMTimeoutStartVM).(int)
			err = vmAPI.StartVM(ctx, startVMTimeout)
			if err != nil {
				return diag.FromErr(err)
			}

			// This concludes an equivalent of a reboot, avoid doing another.
			reboot = false
		}
	}

	// Perform a regular reboot in case it's necessary and haven't already been done.
	if reboot {
		rebootTimeout := d.Get(mkResourceVirtualEnvironmentVMTimeoutReboot).(int)

		err := vmAPI.RebootVM(
			ctx,
			&vms.RebootRequestBody{
				Timeout: &rebootTimeout,
			},
			rebootTimeout+30,
		)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return vmRead(ctx, d, m)
}

func vmDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(proxmoxtf.ProviderConfiguration)
	api, err := config.GetClient()
	if err != nil {
		return diag.FromErr(err)
	}

	nodeName := d.Get(mkResourceVirtualEnvironmentVMNodeName).(string)
	vmID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	vmAPI := api.Node(nodeName).VM(vmID)

	// Shut down the virtual machine before deleting it.
	status, err := vmAPI.GetVMStatus(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	if status.Status != "stopped" {
		forceStop := types.CustomBool(true)
		shutdownTimeout := d.Get(mkResourceVirtualEnvironmentVMTimeoutShutdownVM).(int)

		err = vmAPI.ShutdownVM(
			ctx,
			&vms.ShutdownRequestBody{
				ForceStop: &forceStop,
				Timeout:   &shutdownTimeout,
			},
			shutdownTimeout+30,
		)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	err = vmAPI.DeleteVM(ctx)

	if err != nil {
		if strings.Contains(err.Error(), "HTTP 404") ||
			(strings.Contains(err.Error(), "HTTP 500") && strings.Contains(err.Error(), "does not exist")) {
			d.SetId("")

			return nil
		}
		return diag.FromErr(err)
	}

	// Wait for the state to become unavailable as that clearly indicates the destruction of the VM.
	err = vmAPI.WaitForVMState(ctx, "", 60, 2)
	if err == nil {
		return diag.Errorf("failed to delete VM \"%d\"", vmID)
	}

	d.SetId("")

	return nil
}
