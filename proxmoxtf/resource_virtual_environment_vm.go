/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmoxtf

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/danitso/terraform-provider-proxmox/proxmox"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

const (
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
	dvResourceVirtualEnvironmentVMCPUArchitecture                   = "x86_64"
	dvResourceVirtualEnvironmentVMCPUCores                          = 1
	dvResourceVirtualEnvironmentVMCPUHotplugged                     = 0
	dvResourceVirtualEnvironmentVMCPUSockets                        = 1
	dvResourceVirtualEnvironmentVMCPUType                           = "qemu64"
	dvResourceVirtualEnvironmentVMCPUUnits                          = 1024
	dvResourceVirtualEnvironmentVMDescription                       = ""
	dvResourceVirtualEnvironmentVMDiskDatastoreID                   = "local-lvm"
	dvResourceVirtualEnvironmentVMDiskFileFormat                    = "qcow2"
	dvResourceVirtualEnvironmentVMDiskFileID                        = ""
	dvResourceVirtualEnvironmentVMDiskSize                          = 8
	dvResourceVirtualEnvironmentVMDiskSpeedRead                     = 0
	dvResourceVirtualEnvironmentVMDiskSpeedReadBurstable            = 0
	dvResourceVirtualEnvironmentVMDiskSpeedWrite                    = 0
	dvResourceVirtualEnvironmentVMDiskSpeedWriteBurstable           = 0
	dvResourceVirtualEnvironmentVMInitializationDatastoreID         = "local-lvm"
	dvResourceVirtualEnvironmentVMInitializationDNSDomain           = ""
	dvResourceVirtualEnvironmentVMInitializationDNSServer           = ""
	dvResourceVirtualEnvironmentVMInitializationIPConfigIPv4Address = ""
	dvResourceVirtualEnvironmentVMInitializationIPConfigIPv4Gateway = ""
	dvResourceVirtualEnvironmentVMInitializationIPConfigIPv6Address = ""
	dvResourceVirtualEnvironmentVMInitializationIPConfigIPv6Gateway = ""
	dvResourceVirtualEnvironmentVMInitializationUserAccountPassword = ""
	dvResourceVirtualEnvironmentVMInitializationUserDataFileID      = ""
	dvResourceVirtualEnvironmentVMKeyboardLayout                    = "en-us"
	dvResourceVirtualEnvironmentVMMemoryDedicated                   = 512
	dvResourceVirtualEnvironmentVMMemoryFloating                    = 0
	dvResourceVirtualEnvironmentVMMemoryShared                      = 0
	dvResourceVirtualEnvironmentVMName                              = ""
	dvResourceVirtualEnvironmentVMNetworkDeviceBridge               = "vmbr0"
	dvResourceVirtualEnvironmentVMNetworkDeviceEnabled              = true
	dvResourceVirtualEnvironmentVMNetworkDeviceMACAddress           = ""
	dvResourceVirtualEnvironmentVMNetworkDeviceModel                = "virtio"
	dvResourceVirtualEnvironmentVMNetworkDeviceRateLimit            = 0
	dvResourceVirtualEnvironmentVMNetworkDeviceVLANID               = 0
	dvResourceVirtualEnvironmentVMOperatingSystemType               = "other"
	dvResourceVirtualEnvironmentVMPoolID                            = ""
	dvResourceVirtualEnvironmentVMSerialDeviceDevice                = "socket"
	dvResourceVirtualEnvironmentVMStarted                           = true
	dvResourceVirtualEnvironmentVMTabletDevice                      = true
	dvResourceVirtualEnvironmentVMTemplate                          = false
	dvResourceVirtualEnvironmentVMVGAEnabled                        = true
	dvResourceVirtualEnvironmentVMVGAMemory                         = 16
	dvResourceVirtualEnvironmentVMVGAType                           = "std"
	dvResourceVirtualEnvironmentVMVMID                              = -1

	maxResourceVirtualEnvironmentVMAudioDevices   = 1
	maxResourceVirtualEnvironmentVMNetworkDevices = 8
	maxResourceVirtualEnvironmentVMSerialDevices  = 4

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
	mkResourceVirtualEnvironmentVMCloneDatastoreID                  = "datastore_id"
	mkResourceVirtualEnvironmentVMCloneNodeName                     = "node_name"
	mkResourceVirtualEnvironmentVMCloneVMID                         = "vm_id"
	mkResourceVirtualEnvironmentVMCPU                               = "cpu"
	mkResourceVirtualEnvironmentVMCPUArchitecture                   = "architecture"
	mkResourceVirtualEnvironmentVMCPUCores                          = "cores"
	mkResourceVirtualEnvironmentVMCPUFlags                          = "flags"
	mkResourceVirtualEnvironmentVMCPUHotplugged                     = "hotplugged"
	mkResourceVirtualEnvironmentVMCPUSockets                        = "sockets"
	mkResourceVirtualEnvironmentVMCPUType                           = "type"
	mkResourceVirtualEnvironmentVMCPUUnits                          = "units"
	mkResourceVirtualEnvironmentVMDescription                       = "description"
	mkResourceVirtualEnvironmentVMDisk                              = "disk"
	mkResourceVirtualEnvironmentVMDiskDatastoreID                   = "datastore_id"
	mkResourceVirtualEnvironmentVMDiskFileFormat                    = "file_format"
	mkResourceVirtualEnvironmentVMDiskFileID                        = "file_id"
	mkResourceVirtualEnvironmentVMDiskSize                          = "size"
	mkResourceVirtualEnvironmentVMDiskSpeed                         = "speed"
	mkResourceVirtualEnvironmentVMDiskSpeedRead                     = "read"
	mkResourceVirtualEnvironmentVMDiskSpeedReadBurstable            = "read_burstable"
	mkResourceVirtualEnvironmentVMDiskSpeedWrite                    = "write"
	mkResourceVirtualEnvironmentVMDiskSpeedWriteBurstable           = "write_burstable"
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
	mkResourceVirtualEnvironmentVMInitializationUserAccount         = "user_account"
	mkResourceVirtualEnvironmentVMInitializationUserAccountKeys     = "keys"
	mkResourceVirtualEnvironmentVMInitializationUserAccountPassword = "password"
	mkResourceVirtualEnvironmentVMInitializationUserAccountUsername = "username"
	mkResourceVirtualEnvironmentVMInitializationUserDataFileID      = "user_data_file_id"
	mkResourceVirtualEnvironmentVMIPv4Addresses                     = "ipv4_addresses"
	mkResourceVirtualEnvironmentVMIPv6Addresses                     = "ipv6_addresses"
	mkResourceVirtualEnvironmentVMKeyboardLayout                    = "keyboard_layout"
	mkResourceVirtualEnvironmentVMMACAddresses                      = "mac_addresses"
	mkResourceVirtualEnvironmentVMMemory                            = "memory"
	mkResourceVirtualEnvironmentVMMemoryDedicated                   = "dedicated"
	mkResourceVirtualEnvironmentVMMemoryFloating                    = "floating"
	mkResourceVirtualEnvironmentVMMemoryShared                      = "shared"
	mkResourceVirtualEnvironmentVMName                              = "name"
	mkResourceVirtualEnvironmentVMNetworkDevice                     = "network_device"
	mkResourceVirtualEnvironmentVMNetworkDeviceBridge               = "bridge"
	mkResourceVirtualEnvironmentVMNetworkDeviceEnabled              = "enabled"
	mkResourceVirtualEnvironmentVMNetworkDeviceMACAddress           = "mac_address"
	mkResourceVirtualEnvironmentVMNetworkDeviceModel                = "model"
	mkResourceVirtualEnvironmentVMNetworkDeviceRateLimit            = "rate_limit"
	mkResourceVirtualEnvironmentVMNetworkDeviceVLANID               = "vlan_id"
	mkResourceVirtualEnvironmentVMNetworkInterfaceNames             = "network_interface_names"
	mkResourceVirtualEnvironmentVMNodeName                          = "node_name"
	mkResourceVirtualEnvironmentVMOperatingSystem                   = "operating_system"
	mkResourceVirtualEnvironmentVMOperatingSystemType               = "type"
	mkResourceVirtualEnvironmentVMPoolID                            = "pool_id"
	mkResourceVirtualEnvironmentVMSerialDevice                      = "serial_device"
	mkResourceVirtualEnvironmentVMSerialDeviceDevice                = "device"
	mkResourceVirtualEnvironmentVMStarted                           = "started"
	mkResourceVirtualEnvironmentVMTabletDevice                      = "tablet_device"
	mkResourceVirtualEnvironmentVMTemplate                          = "template"
	mkResourceVirtualEnvironmentVMVGA                               = "vga"
	mkResourceVirtualEnvironmentVMVGAEnabled                        = "enabled"
	mkResourceVirtualEnvironmentVMVGAMemory                         = "memory"
	mkResourceVirtualEnvironmentVMVGAType                           = "type"
	mkResourceVirtualEnvironmentVMVMID                              = "vm_id"
)

func resourceVirtualEnvironmentVM() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
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
							mkResourceVirtualEnvironmentVMAgentTrim:    dvResourceVirtualEnvironmentVMAgentEnabled,
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
							Type:         schema.TypeString,
							Description:  "The maximum amount of time to wait for data from the QEMU agent to become available",
							Optional:     true,
							Default:      dvResourceVirtualEnvironmentVMAgentTimeout,
							ValidateFunc: getTimeoutValidator(),
						},
						mkResourceVirtualEnvironmentVMAgentTrim: {
							Type:        schema.TypeBool,
							Description: "Whether to enable the FSTRIM feature in the QEMU agent",
							Optional:    true,
							Default:     dvResourceVirtualEnvironmentVMAgentTrim,
						},
						mkResourceVirtualEnvironmentVMAgentType: {
							Type:         schema.TypeString,
							Description:  "The QEMU agent interface type",
							Optional:     true,
							Default:      dvResourceVirtualEnvironmentVMAgentType,
							ValidateFunc: getQEMUAgentTypeValidator(),
						},
					},
				},
				MaxItems: 1,
				MinItems: 0,
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
							Type:         schema.TypeString,
							Description:  "The device",
							Optional:     true,
							Default:      dvResourceVirtualEnvironmentVMAudioDeviceDevice,
							ValidateFunc: resourceVirtualEnvironmentVMGetAudioDeviceValidator(),
						},
						mkResourceVirtualEnvironmentVMAudioDeviceDriver: {
							Type:         schema.TypeString,
							Description:  "The driver",
							Optional:     true,
							Default:      dvResourceVirtualEnvironmentVMAudioDeviceDriver,
							ValidateFunc: resourceVirtualEnvironmentVMGetAudioDriverValidator(),
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
				Type:         schema.TypeString,
				Description:  "The BIOS implementation",
				Optional:     true,
				Default:      dvResourceVirtualEnvironmentVMBIOS,
				ValidateFunc: getBIOSValidator(),
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
							Type:         schema.TypeString,
							Description:  "The file id",
							Optional:     true,
							Default:      dvResourceVirtualEnvironmentVMCDROMFileID,
							ValidateFunc: getFileIDValidator(),
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
							Type:         schema.TypeInt,
							Description:  "The ID of the source VM",
							Required:     true,
							ForceNew:     true,
							ValidateFunc: getVMIDValidator(),
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
							Type:         schema.TypeString,
							Description:  "The CPU architecture",
							Optional:     true,
							Default:      dvResourceVirtualEnvironmentVMCPUArchitecture,
							ValidateFunc: resourceVirtualEnvironmentVMGetCPUArchitectureValidator(),
						},
						mkResourceVirtualEnvironmentVMCPUCores: {
							Type:         schema.TypeInt,
							Description:  "The number of CPU cores",
							Optional:     true,
							Default:      dvResourceVirtualEnvironmentVMCPUCores,
							ValidateFunc: validation.IntBetween(1, 2304),
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
							Type:         schema.TypeInt,
							Description:  "The number of hotplugged vCPUs",
							Optional:     true,
							Default:      dvResourceVirtualEnvironmentVMCPUHotplugged,
							ValidateFunc: validation.IntBetween(0, 2304),
						},
						mkResourceVirtualEnvironmentVMCPUSockets: {
							Type:         schema.TypeInt,
							Description:  "The number of CPU sockets",
							Optional:     true,
							Default:      dvResourceVirtualEnvironmentVMCPUSockets,
							ValidateFunc: validation.IntBetween(1, 16),
						},
						mkResourceVirtualEnvironmentVMCPUType: {
							Type:         schema.TypeString,
							Description:  "The emulated CPU type",
							Optional:     true,
							Default:      dvResourceVirtualEnvironmentVMCPUType,
							ValidateFunc: getCPUTypeValidator(),
						},
						mkResourceVirtualEnvironmentVMCPUUnits: {
							Type:         schema.TypeInt,
							Description:  "The CPU units",
							Optional:     true,
							Default:      dvResourceVirtualEnvironmentVMCPUUnits,
							ValidateFunc: validation.IntBetween(2, 262144),
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
							mkResourceVirtualEnvironmentVMDiskFileFormat:  dvResourceVirtualEnvironmentVMDiskFileFormat,
							mkResourceVirtualEnvironmentVMDiskFileID:      dvResourceVirtualEnvironmentVMDiskFileID,
							mkResourceVirtualEnvironmentVMDiskSize:        dvResourceVirtualEnvironmentVMDiskSize,
						},
					}, nil
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mkResourceVirtualEnvironmentVMDiskDatastoreID: {
							Type:        schema.TypeString,
							Description: "The datastore id",
							Optional:    true,
							ForceNew:    true,
							Default:     dvResourceVirtualEnvironmentVMDiskDatastoreID,
						},
						mkResourceVirtualEnvironmentVMDiskFileFormat: {
							Type:         schema.TypeString,
							Description:  "The file format",
							Optional:     true,
							ForceNew:     true,
							Default:      dvResourceVirtualEnvironmentVMDiskFileFormat,
							ValidateFunc: getFileFormatValidator(),
						},
						mkResourceVirtualEnvironmentVMDiskFileID: {
							Type:         schema.TypeString,
							Description:  "The file id for a disk image",
							Optional:     true,
							ForceNew:     true,
							Default:      dvResourceVirtualEnvironmentVMDiskFileID,
							ValidateFunc: getFileIDValidator(),
						},
						mkResourceVirtualEnvironmentVMDiskSize: {
							Type:         schema.TypeInt,
							Description:  "The disk size in gigabytes",
							Optional:     true,
							ForceNew:     true,
							Default:      dvResourceVirtualEnvironmentVMDiskSize,
							ValidateFunc: validation.IntBetween(1, 8192),
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
											return len(old) > 0 && strings.ReplaceAll(old, "*", "") == ""
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
							Type:         schema.TypeString,
							Description:  "The ID of a file containing custom user data",
							Optional:     true,
							ForceNew:     true,
							Default:      dvResourceVirtualEnvironmentVMInitializationUserDataFileID,
							ValidateFunc: getFileIDValidator(),
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
			mkResourceVirtualEnvironmentVMKeyboardLayout: {
				Type:         schema.TypeString,
				Description:  "The keyboard layout",
				Optional:     true,
				Default:      dvResourceVirtualEnvironmentVMKeyboardLayout,
				ValidateFunc: getKeyboardLayoutValidator(),
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
							Type:         schema.TypeInt,
							Description:  "The dedicated memory in megabytes",
							Optional:     true,
							Default:      dvResourceVirtualEnvironmentVMMemoryDedicated,
							ValidateFunc: validation.IntBetween(64, 268435456),
						},
						mkResourceVirtualEnvironmentVMMemoryFloating: {
							Type:         schema.TypeInt,
							Description:  "The floating memory in megabytes (balloon)",
							Optional:     true,
							Default:      dvResourceVirtualEnvironmentVMMemoryFloating,
							ValidateFunc: validation.IntBetween(0, 268435456),
						},
						mkResourceVirtualEnvironmentVMMemoryShared: {
							Type:         schema.TypeInt,
							Description:  "The shared memory in megabytes",
							Optional:     true,
							Default:      dvResourceVirtualEnvironmentVMMemoryShared,
							ValidateFunc: validation.IntBetween(0, 268435456),
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
						mkResourceVirtualEnvironmentVMNetworkDeviceMACAddress: {
							Type:        schema.TypeString,
							Description: "The MAC address",
							Optional:    true,
							Default:     dvResourceVirtualEnvironmentVMNetworkDeviceMACAddress,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								return new == ""
							},
							ValidateFunc: getMACAddressValidator(),
						},
						mkResourceVirtualEnvironmentVMNetworkDeviceModel: {
							Type:         schema.TypeString,
							Description:  "The model",
							Optional:     true,
							Default:      dvResourceVirtualEnvironmentVMNetworkDeviceModel,
							ValidateFunc: getNetworkDeviceModelValidator(),
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
							Type:         schema.TypeString,
							Description:  "The type",
							Optional:     true,
							Default:      dvResourceVirtualEnvironmentVMOperatingSystemType,
							ValidateFunc: resourceVirtualEnvironmentVMGetOperatingSystemTypeValidator(),
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
							Type:         schema.TypeString,
							Description:  "The device",
							Optional:     true,
							Default:      dvResourceVirtualEnvironmentVMSerialDeviceDevice,
							ValidateFunc: resourceVirtualEnvironmentVMGetSerialDeviceValidator(),
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
			mkResourceVirtualEnvironmentVMTemplate: {
				Type:        schema.TypeBool,
				Description: "Whether to create a template",
				Optional:    true,
				ForceNew:    true,
				Default:     dvResourceVirtualEnvironmentVMTemplate,
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
							Type:         schema.TypeInt,
							Description:  "The VGA memory in megabytes (4-512 MB)",
							Optional:     true,
							Default:      dvResourceVirtualEnvironmentVMVGAMemory,
							ValidateFunc: getVGAMemoryValidator(),
						},
						mkResourceVirtualEnvironmentVMVGAType: {
							Type:         schema.TypeString,
							Description:  "The VGA type",
							Optional:     true,
							Default:      dvResourceVirtualEnvironmentVMVGAType,
							ValidateFunc: getVGATypeValidator(),
						},
					},
				},
				MaxItems: 1,
				MinItems: 0,
			},
			mkResourceVirtualEnvironmentVMVMID: {
				Type:         schema.TypeInt,
				Description:  "The VM identifier",
				Optional:     true,
				ForceNew:     true,
				Default:      dvResourceVirtualEnvironmentVMVMID,
				ValidateFunc: getVMIDValidator(),
			},
		},
		Create: resourceVirtualEnvironmentVMCreate,
		Read:   resourceVirtualEnvironmentVMRead,
		Update: resourceVirtualEnvironmentVMUpdate,
		Delete: resourceVirtualEnvironmentVMDelete,
	}
}

func resourceVirtualEnvironmentVMCreate(d *schema.ResourceData, m interface{}) error {
	clone := d.Get(mkResourceVirtualEnvironmentVMClone).([]interface{})

	if len(clone) > 0 {
		return resourceVirtualEnvironmentVMCreateClone(d, m)
	}

	return resourceVirtualEnvironmentVMCreateCustom(d, m)
}

func resourceVirtualEnvironmentVMCreateClone(d *schema.ResourceData, m interface{}) error {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()

	if err != nil {
		return err
	}

	clone := d.Get(mkResourceVirtualEnvironmentVMClone).([]interface{})
	cloneBlock := clone[0].(map[string]interface{})
	cloneDatastoreID := cloneBlock[mkResourceVirtualEnvironmentVMCloneDatastoreID].(string)
	cloneNodeName := cloneBlock[mkResourceVirtualEnvironmentVMCloneNodeName].(string)
	cloneVMID := cloneBlock[mkResourceVirtualEnvironmentVMCloneVMID].(int)

	description := d.Get(mkResourceVirtualEnvironmentVMDescription).(string)
	name := d.Get(mkResourceVirtualEnvironmentVMName).(string)
	nodeName := d.Get(mkResourceVirtualEnvironmentVMNodeName).(string)
	poolID := d.Get(mkResourceVirtualEnvironmentVMPoolID).(string)
	vmID := d.Get(mkResourceVirtualEnvironmentVMVMID).(int)

	if vmID == -1 {
		vmIDNew, err := veClient.GetVMID()

		if err != nil {
			return err
		}

		vmID = *vmIDNew
	}

	fullCopy := proxmox.CustomBool(true)

	cloneBody := &proxmox.VirtualEnvironmentVMCloneRequestBody{
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

	if cloneNodeName != "" && cloneNodeName != nodeName {
		cloneBody.TargetNodeName = &nodeName

		err = veClient.CloneVM(cloneNodeName, cloneVMID, cloneBody)
	} else {
		err = veClient.CloneVM(nodeName, cloneVMID, cloneBody)
	}

	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(vmID))

	// Wait for the virtual machine to be created and its configuration lock to be released.
	err = veClient.WaitForVMConfigUnlock(nodeName, vmID, 600, 5, true)

	if err != nil {
		return err
	}

	// Now that the virtual machine has been cloned, we need to perform some modifications.
	acpi := proxmox.CustomBool(d.Get(mkResourceVirtualEnvironmentVMACPI).(bool))
	agent := d.Get(mkResourceVirtualEnvironmentVMAgent).([]interface{})
	audioDevices, err := resourceVirtualEnvironmentVMGetAudioDeviceList(d, m)

	if err != nil {
		return err
	}

	bios := d.Get(mkResourceVirtualEnvironmentVMBIOS).(string)
	cdrom := d.Get(mkResourceVirtualEnvironmentVMCDROM).([]interface{})
	cpu := d.Get(mkResourceVirtualEnvironmentVMCPU).([]interface{})
	initialization := d.Get(mkResourceVirtualEnvironmentVMInitialization).([]interface{})
	keyboardLayout := d.Get(mkResourceVirtualEnvironmentVMKeyboardLayout).(string)
	memory := d.Get(mkResourceVirtualEnvironmentVMMemory).([]interface{})
	networkDevice := d.Get(mkResourceVirtualEnvironmentVMNetworkDevice).([]interface{})
	operatingSystem := d.Get(mkResourceVirtualEnvironmentVMOperatingSystem).([]interface{})
	serialDevice := d.Get(mkResourceVirtualEnvironmentVMSerialDevice).([]interface{})
	started := proxmox.CustomBool(d.Get(mkResourceVirtualEnvironmentVMStarted).(bool))
	tabletDevice := proxmox.CustomBool(d.Get(mkResourceVirtualEnvironmentVMTabletDevice).(bool))
	template := proxmox.CustomBool(d.Get(mkResourceVirtualEnvironmentVMTemplate).(bool))
	vga := d.Get(mkResourceVirtualEnvironmentVMVGA).([]interface{})

	updateBody := &proxmox.VirtualEnvironmentVMUpdateRequestBody{
		AudioDevices: audioDevices,
	}

	delete := []string{}

	if acpi != dvResourceVirtualEnvironmentVMACPI {
		updateBody.ACPI = &acpi
	}

	if len(agent) > 0 {
		agentBlock := agent[0].(map[string]interface{})

		agentEnabled := proxmox.CustomBool(agentBlock[mkResourceVirtualEnvironmentVMAgentEnabled].(bool))
		agentTrim := proxmox.CustomBool(agentBlock[mkResourceVirtualEnvironmentVMAgentTrim].(bool))
		agentType := agentBlock[mkResourceVirtualEnvironmentVMAgentType].(string)

		updateBody.Agent = &proxmox.CustomAgent{
			Enabled:         &agentEnabled,
			TrimClonedDisks: &agentTrim,
			Type:            &agentType,
		}
	}

	if bios != dvResourceVirtualEnvironmentVMBIOS {
		updateBody.BIOS = &bios
	}

	if len(cdrom) > 0 {
		cdromBlock := cdrom[0].(map[string]interface{})

		cdromEnabled := cdromBlock[mkResourceVirtualEnvironmentVMCDROMEnabled].(bool)
		cdromFileID := cdromBlock[mkResourceVirtualEnvironmentVMCDROMFileID].(string)

		if cdromFileID == "" {
			cdromFileID = "cdrom"
		}

		cdromMedia := "cdrom"

		updateBody.IDEDevices = proxmox.CustomStorageDevices{
			proxmox.CustomStorageDevice{
				Enabled: false,
			},
			proxmox.CustomStorageDevice{
				Enabled: false,
			},
			proxmox.CustomStorageDevice{
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
		cpuSockets := cpuBlock[mkResourceVirtualEnvironmentVMCPUSockets].(int)
		cpuType := cpuBlock[mkResourceVirtualEnvironmentVMCPUType].(string)
		cpuUnits := cpuBlock[mkResourceVirtualEnvironmentVMCPUUnits].(int)

		cpuFlagsConverted := make([]string, len(cpuFlags))

		for fi, flag := range cpuFlags {
			cpuFlagsConverted[fi] = flag.(string)
		}

		// Only the root account is allowed to change the CPU architecture, which makes this check necessary.
		if veClient.Username == proxmox.DefaultRootAccount || cpuArchitecture != dvResourceVirtualEnvironmentVMCPUArchitecture {
			updateBody.CPUArchitecture = &cpuArchitecture
		}

		updateBody.CPUCores = &cpuCores
		updateBody.CPUEmulation = &proxmox.CustomCPUEmulation{
			Flags: &cpuFlagsConverted,
			Type:  cpuType,
		}
		updateBody.CPUSockets = &cpuSockets
		updateBody.CPUUnits = &cpuUnits

		if cpuHotplugged > 0 {
			updateBody.VirtualCPUCount = &cpuHotplugged
		}
	}

	if len(initialization) > 0 {
		initializationBlock := initialization[0].(map[string]interface{})
		initializationDatastoreID := initializationBlock[mkResourceVirtualEnvironmentVMInitializationDatastoreID].(string)

		cdromEnabled := true
		cdromFileID := fmt.Sprintf("%s:cloudinit", initializationDatastoreID)
		cdromMedia := "cdrom"

		updateBody.IDEDevices = proxmox.CustomStorageDevices{
			proxmox.CustomStorageDevice{
				Enabled: false,
			},
			proxmox.CustomStorageDevice{
				Enabled: false,
			},
			proxmox.CustomStorageDevice{
				Enabled:    cdromEnabled,
				FileVolume: cdromFileID,
				Media:      &cdromMedia,
			},
		}

		initializationConfig, err := resourceVirtualEnvironmentVMGetCloudInitConfig(d, m)

		if err != nil {
			return err
		}

		updateBody.CloudInitConfig = initializationConfig
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

			updateBody.SharedMemory = &proxmox.CustomSharedMemory{
				Name: &memorySharedName,
				Size: memoryShared,
			}
		}
	}

	if len(networkDevice) > 0 {
		updateBody.NetworkDevices, err = resourceVirtualEnvironmentVMGetNetworkDeviceObjects(d, m)

		if err != nil {
			return err
		}

		for i := 0; i < len(updateBody.NetworkDevices); i++ {
			if !updateBody.NetworkDevices[i].Enabled {
				delete = append(delete, fmt.Sprintf("net%d", i))
			}
		}

		for i := len(updateBody.NetworkDevices); i < maxResourceVirtualEnvironmentVMNetworkDevices; i++ {
			delete = append(delete, fmt.Sprintf("net%d", i))
		}
	}

	if len(operatingSystem) > 0 {
		operatingSystemBlock := operatingSystem[0].(map[string]interface{})
		operatingSystemType := operatingSystemBlock[mkResourceVirtualEnvironmentVMOperatingSystemType].(string)

		updateBody.OSType = &operatingSystemType
	}

	if len(serialDevice) > 0 {
		updateBody.SerialDevices, err = resourceVirtualEnvironmentVMGetSerialDeviceList(d, m)

		if err != nil {
			return err
		}

		for i := len(updateBody.SerialDevices); i < maxResourceVirtualEnvironmentVMSerialDevices; i++ {
			delete = append(delete, fmt.Sprintf("serial%d", i))
		}
	}

	if started != dvResourceVirtualEnvironmentVMStarted {
		updateBody.StartOnBoot = &started
	}

	if tabletDevice != dvResourceVirtualEnvironmentVMTabletDevice {
		updateBody.TabletDeviceEnabled = &tabletDevice
	}

	if template != dvResourceVirtualEnvironmentVMTemplate {
		updateBody.Template = &template
	}

	if len(vga) > 0 {
		vgaDevice, err := resourceVirtualEnvironmentVMGetVGADeviceObject(d, m)

		if err != nil {
			return err
		}

		updateBody.VGADevice = vgaDevice
	}

	updateBody.Delete = delete

	err = veClient.UpdateVM(nodeName, vmID, updateBody)

	if err != nil {
		return err
	}

	return resourceVirtualEnvironmentVMCreateStart(d, m)
}

func resourceVirtualEnvironmentVMCreateCustom(d *schema.ResourceData, m interface{}) error {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()

	if err != nil {
		return err
	}

	resource := resourceVirtualEnvironmentVM()

	acpi := proxmox.CustomBool(d.Get(mkResourceVirtualEnvironmentVMACPI).(bool))

	agentBlock, err := getSchemaBlock(resource, d, m, []string{mkResourceVirtualEnvironmentVMAgent}, 0, true)

	if err != nil {
		return err
	}

	agentEnabled := proxmox.CustomBool(agentBlock[mkResourceVirtualEnvironmentVMAgentEnabled].(bool))
	agentTrim := proxmox.CustomBool(agentBlock[mkResourceVirtualEnvironmentVMAgentTrim].(bool))
	agentType := agentBlock[mkResourceVirtualEnvironmentVMAgentType].(string)

	audioDevices, err := resourceVirtualEnvironmentVMGetAudioDeviceList(d, m)

	if err != nil {
		return err
	}

	bios := d.Get(mkResourceVirtualEnvironmentVMBIOS).(string)

	cdromBlock, err := getSchemaBlock(resource, d, m, []string{mkResourceVirtualEnvironmentVMCDROM}, 0, true)

	if err != nil {
		return err
	}

	cdromEnabled := cdromBlock[mkResourceVirtualEnvironmentVMCDROMEnabled].(bool)
	cdromFileID := cdromBlock[mkResourceVirtualEnvironmentVMCDROMFileID].(string)

	if cdromFileID == "" {
		cdromFileID = "cdrom"
	}

	cpuBlock, err := getSchemaBlock(resource, d, m, []string{mkResourceVirtualEnvironmentVMCPU}, 0, true)

	if err != nil {
		return err
	}

	cpuArchitecture := cpuBlock[mkResourceVirtualEnvironmentVMCPUArchitecture].(string)
	cpuCores := cpuBlock[mkResourceVirtualEnvironmentVMCPUCores].(int)
	cpuFlags := cpuBlock[mkResourceVirtualEnvironmentVMCPUFlags].([]interface{})
	cpuHotplugged := cpuBlock[mkResourceVirtualEnvironmentVMCPUHotplugged].(int)
	cpuSockets := cpuBlock[mkResourceVirtualEnvironmentVMCPUSockets].(int)
	cpuType := cpuBlock[mkResourceVirtualEnvironmentVMCPUType].(string)
	cpuUnits := cpuBlock[mkResourceVirtualEnvironmentVMCPUUnits].(int)

	description := d.Get(mkResourceVirtualEnvironmentVMDescription).(string)
	diskDeviceObjects, err := resourceVirtualEnvironmentVMGetDiskDeviceObjects(d, m)

	if err != nil {
		return err
	}

	initializationConfig, err := resourceVirtualEnvironmentVMGetCloudInitConfig(d, m)

	if err != nil {
		return err
	}

	if initializationConfig != nil {
		initialization := d.Get(mkResourceVirtualEnvironmentVMInitialization).([]interface{})
		initializationBlock := initialization[0].(map[string]interface{})
		initializationDatastoreID := initializationBlock[mkResourceVirtualEnvironmentVMInitializationDatastoreID].(string)

		cdromEnabled = true
		cdromFileID = fmt.Sprintf("%s:cloudinit", initializationDatastoreID)
	}

	keyboardLayout := d.Get(mkResourceVirtualEnvironmentVMKeyboardLayout).(string)
	memoryBlock, err := getSchemaBlock(resource, d, m, []string{mkResourceVirtualEnvironmentVMMemory}, 0, true)

	if err != nil {
		return err
	}

	memoryDedicated := memoryBlock[mkResourceVirtualEnvironmentVMMemoryDedicated].(int)
	memoryFloating := memoryBlock[mkResourceVirtualEnvironmentVMMemoryFloating].(int)
	memoryShared := memoryBlock[mkResourceVirtualEnvironmentVMMemoryShared].(int)

	name := d.Get(mkResourceVirtualEnvironmentVMName).(string)

	networkDeviceObjects, err := resourceVirtualEnvironmentVMGetNetworkDeviceObjects(d, m)

	if err != nil {
		return err
	}

	nodeName := d.Get(mkResourceVirtualEnvironmentVMNodeName).(string)

	operatingSystem, err := getSchemaBlock(resource, d, m, []string{mkResourceVirtualEnvironmentVMOperatingSystem}, 0, true)

	if err != nil {
		return err
	}

	operatingSystemType := operatingSystem[mkResourceVirtualEnvironmentVMOperatingSystemType].(string)

	poolID := d.Get(mkResourceVirtualEnvironmentVMPoolID).(string)

	serialDevices, err := resourceVirtualEnvironmentVMGetSerialDeviceList(d, m)

	if err != nil {
		return err
	}

	started := proxmox.CustomBool(d.Get(mkResourceVirtualEnvironmentVMStarted).(bool))
	tabletDevice := proxmox.CustomBool(d.Get(mkResourceVirtualEnvironmentVMTabletDevice).(bool))
	template := proxmox.CustomBool(d.Get(mkResourceVirtualEnvironmentVMTemplate).(bool))

	vgaDevice, err := resourceVirtualEnvironmentVMGetVGADeviceObject(d, m)

	if err != nil {
		return err
	}

	vmID := d.Get(mkResourceVirtualEnvironmentVMVMID).(int)

	if vmID == -1 {
		vmIDNew, err := veClient.GetVMID()

		if err != nil {
			return err
		}

		vmID = *vmIDNew
	}

	var memorySharedObject *proxmox.CustomSharedMemory

	bootDisk := "scsi0"
	bootOrder := "c"

	if cdromEnabled {
		bootOrder = "cd"
	}

	cpuFlagsConverted := make([]string, len(cpuFlags))

	for fi, flag := range cpuFlags {
		cpuFlagsConverted[fi] = flag.(string)
	}

	ideDevice2Media := "cdrom"
	ideDevices := proxmox.CustomStorageDevices{
		proxmox.CustomStorageDevice{
			Enabled: false,
		},
		proxmox.CustomStorageDevice{
			Enabled: false,
		},
		proxmox.CustomStorageDevice{
			Enabled:    cdromEnabled,
			FileVolume: cdromFileID,
			Media:      &ideDevice2Media,
		},
	}

	if memoryShared > 0 {
		memorySharedName := fmt.Sprintf("vm-%d-ivshmem", vmID)
		memorySharedObject = &proxmox.CustomSharedMemory{
			Name: &memorySharedName,
			Size: memoryShared,
		}
	}

	scsiHardware := "virtio-scsi-pci"

	createBody := &proxmox.VirtualEnvironmentVMCreateRequestBody{
		ACPI: &acpi,
		Agent: &proxmox.CustomAgent{
			Enabled:         &agentEnabled,
			TrimClonedDisks: &agentTrim,
			Type:            &agentType,
		},
		AudioDevices:    audioDevices,
		BIOS:            &bios,
		BootDisk:        &bootDisk,
		BootOrder:       &bootOrder,
		CloudInitConfig: initializationConfig,
		CPUCores:        &cpuCores,
		CPUEmulation: &proxmox.CustomCPUEmulation{
			Flags: &cpuFlagsConverted,
			Type:  cpuType,
		},
		CPUSockets:          &cpuSockets,
		CPUUnits:            &cpuUnits,
		DedicatedMemory:     &memoryDedicated,
		FloatingMemory:      &memoryFloating,
		IDEDevices:          ideDevices,
		KeyboardLayout:      &keyboardLayout,
		NetworkDevices:      networkDeviceObjects,
		OSType:              &operatingSystemType,
		PoolID:              &poolID,
		SCSIDevices:         diskDeviceObjects,
		SCSIHardware:        &scsiHardware,
		SerialDevices:       serialDevices,
		SharedMemory:        memorySharedObject,
		StartOnBoot:         &started,
		TabletDeviceEnabled: &tabletDevice,
		Template:            &template,
		VGADevice:           vgaDevice,
		VMID:                &vmID,
	}

	// Only the root account is allowed to change the CPU architecture, which makes this check necessary.
	if veClient.Username == proxmox.DefaultRootAccount || cpuArchitecture != dvResourceVirtualEnvironmentVMCPUArchitecture {
		createBody.CPUArchitecture = &cpuArchitecture
	}

	if cpuHotplugged > 0 {
		createBody.VirtualCPUCount = &cpuHotplugged
	}

	if description != "" {
		createBody.Description = &description
	}

	if name != "" {
		createBody.Name = &name
	}

	err = veClient.CreateVM(nodeName, createBody)

	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(vmID))

	return resourceVirtualEnvironmentVMCreateCustomDisks(d, m)
}

func resourceVirtualEnvironmentVMCreateCustomDisks(d *schema.ResourceData, m interface{}) error {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()

	if err != nil {
		return err
	}

	nodeName := d.Get(mkResourceVirtualEnvironmentVMNodeName).(string)
	vmID, err := strconv.Atoi(d.Id())

	if err != nil {
		return err
	}

	commands := []string{}

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
	resourceSchema := resourceVirtualEnvironmentVM().Schema
	diskSchemaElem := resourceSchema[mkResourceVirtualEnvironmentVMDisk].Elem
	diskSchemaResource := diskSchemaElem.(*schema.Resource)
	diskSpeedResource := diskSchemaResource.Schema[mkResourceVirtualEnvironmentVMDiskSpeed]

	// Generate the commands required to import the specified disks.
	importedDiskCount := 0

	for i, d := range disk {
		block := d.(map[string]interface{})

		fileID, _ := block[mkResourceVirtualEnvironmentVMDiskFileID].(string)

		if fileID == "" {
			continue
		}

		datastoreID, _ := block[mkResourceVirtualEnvironmentVMDiskDatastoreID].(string)
		fileFormat, _ := block[mkResourceVirtualEnvironmentVMDiskFileFormat].(string)
		size, _ := block[mkResourceVirtualEnvironmentVMDiskSize].(int)
		speed := block[mkResourceVirtualEnvironmentVMDiskSpeed].([]interface{})

		if len(speed) == 0 {
			diskSpeedDefault, err := diskSpeedResource.DefaultValue()

			if err != nil {
				return err
			}

			speed = diskSpeedDefault.([]interface{})
		}

		speedBlock := speed[0].(map[string]interface{})
		speedLimitRead := speedBlock[mkResourceVirtualEnvironmentVMDiskSpeedRead].(int)
		speedLimitReadBurstable := speedBlock[mkResourceVirtualEnvironmentVMDiskSpeedReadBurstable].(int)
		speedLimitWrite := speedBlock[mkResourceVirtualEnvironmentVMDiskSpeedWrite].(int)
		speedLimitWriteBurstable := speedBlock[mkResourceVirtualEnvironmentVMDiskSpeedWriteBurstable].(int)

		diskOptions := ""

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

		fileIDParts := strings.Split(fileID, ":")
		filePath := ""

		if strings.HasPrefix(fileIDParts[1], "iso/") {
			filePath = fmt.Sprintf("/template/%s", fileIDParts[1])
		} else {
			filePath = fmt.Sprintf("/%s", fileIDParts[1])
		}

		filePathTmp := fmt.Sprintf("/tmp/vm-%d-disk-%d.%s", vmID, diskCount+importedDiskCount, fileFormat)

		commands = append(
			commands,
			`set -e`,
			fmt.Sprintf(`datastore_id_image="%s"`, fileIDParts[0]),
			fmt.Sprintf(`datastore_id_target="%s"`, datastoreID),
			fmt.Sprintf(`disk_count="%d"`, diskCount+importedDiskCount),
			fmt.Sprintf(`disk_index="%d"`, i),
			fmt.Sprintf(`disk_options="%s"`, diskOptions),
			fmt.Sprintf(`disk_size="%d"`, size),
			fmt.Sprintf(`file_path="%s"`, filePath),
			fmt.Sprintf(`file_path_tmp="%s"`, filePathTmp),
			fmt.Sprintf(`vm_id="%d"`, vmID),
			`getdsi() { local nr='^([A-Za-z0-9_-]+): ([A-Za-z0-9_-]+)$'; local pr='^[[:space:]]+path[[:space:]]+([^[:space:]]+)$'; local dn=""; local dt=""; while IFS='' read -r l || [[ -n "$l" ]]; do if [[ "$l" =~ $nr ]]; then dt="${BASH_REMATCH[1]}"; dn="${BASH_REMATCH[2]}"; elif [[ "$l" =~ $pr ]] && [[ "$dn" == "$1" ]]; then echo "${BASH_REMATCH[1]};${dt}"; break; fi; done < /etc/pve/storage.cfg; }`,
			`dsi_image="$(getdsi "$datastore_id_image")"`,
			`dsp_image="$(echo "$dsi_image" | cut -d ";" -f 1)"`,
			`dst_image="$(echo "$dsi_image" | cut -d ";" -f 2)"`,
			`if [[ -z "$dsp_image" ]]; then echo "Failed to determine the path for datastore '${datastore_id_image}' (${dsi_image})"; exit 1; fi`,
			`dsi_target="$(getdsi "$datastore_id_target")"`,
			`dst_target="$(echo "$dsi_target" | cut -d ";" -f 2)"`,
			`cp "${dsp_image}${file_path}" "$file_path_tmp"`,
			`qemu-img resize "$file_path_tmp" "${disk_size}G"`,
			`qm importdisk "$vm_id" "$file_path_tmp" "$datastore_id_target" -format qcow2`,
			`disk_id="${datastore_id_target}:$([[ "$dst_target" == "dir" ]] && echo "${vm_id}/" || echo "")vm-${vm_id}-disk-${disk_count}$([[ "$dst_target" == "dir" ]] && echo ".qcow2" || echo "")${disk_options}"`,
			`qm set "$vm_id" "-scsi${disk_index}" "$disk_id"`,
			`rm -f "$file_path_tmp"`,
		)

		importedDiskCount++
	}

	// Execute the commands on the node and wait for the result.
	// This is a highly experimental approach to disk imports and is not recommended by Proxmox.
	if len(commands) > 0 {
		err = veClient.ExecuteNodeCommands(nodeName, commands)

		if err != nil {
			return err
		}
	}

	return resourceVirtualEnvironmentVMCreateStart(d, m)
}

func resourceVirtualEnvironmentVMCreateStart(d *schema.ResourceData, m interface{}) error {
	started := d.Get(mkResourceVirtualEnvironmentVMStarted).(bool)
	template := d.Get(mkResourceVirtualEnvironmentVMTemplate).(bool)

	if !started || template {
		return resourceVirtualEnvironmentVMRead(d, m)
	}

	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()

	if err != nil {
		return err
	}

	nodeName := d.Get(mkResourceVirtualEnvironmentVMNodeName).(string)
	vmID, err := strconv.Atoi(d.Id())

	if err != nil {
		return err
	}

	// Start the virtual machine and wait for it to reach a running state before continuing.
	err = veClient.StartVM(nodeName, vmID)

	if err != nil {
		return err
	}

	err = veClient.WaitForVMState(nodeName, vmID, "running", 120, 5)

	if err != nil {
		return err
	}

	return resourceVirtualEnvironmentVMRead(d, m)
}

func resourceVirtualEnvironmentVMGetAudioDeviceList(d *schema.ResourceData, m interface{}) (proxmox.CustomAudioDevices, error) {
	devices := d.Get(mkResourceVirtualEnvironmentVMAudioDevice).([]interface{})
	list := make(proxmox.CustomAudioDevices, len(devices))

	for i, v := range devices {
		block := v.(map[string]interface{})

		device, _ := block[mkResourceVirtualEnvironmentVMAudioDeviceDevice].(string)
		driver, _ := block[mkResourceVirtualEnvironmentVMAudioDeviceDriver].(string)
		enabled, _ := block[mkResourceVirtualEnvironmentVMAudioDeviceEnabled].(bool)

		list[i].Device = device
		list[i].Driver = &driver
		list[i].Enabled = enabled
	}

	return list, nil
}

func resourceVirtualEnvironmentVMGetAudioDeviceValidator() schema.SchemaValidateFunc {
	return validation.StringInSlice([]string{
		"AC97",
		"ich9-intel-hda",
		"intel-hda",
	}, false)
}

func resourceVirtualEnvironmentVMGetAudioDriverValidator() schema.SchemaValidateFunc {
	return validation.StringInSlice([]string{
		"spice",
	}, false)
}

func resourceVirtualEnvironmentVMGetCloudInitConfig(d *schema.ResourceData, m interface{}) (*proxmox.CustomCloudInitConfig, error) {
	var initializationConfig *proxmox.CustomCloudInitConfig

	initialization := d.Get(mkResourceVirtualEnvironmentVMInitialization).([]interface{})

	if len(initialization) > 0 {
		initializationBlock := initialization[0].(map[string]interface{})
		initializationConfig = &proxmox.CustomCloudInitConfig{}
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
		initializationConfig.IPConfig = make([]proxmox.CustomCloudInitIPConfig, len(initializationIPConfig))

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
				sshKeys := make(proxmox.CustomCloudInitSSHKeys, len(keys))

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
			initializationConfig.Files = &proxmox.CustomCloudInitFiles{
				UserVolume: &initializationUserDataFileID,
			}
		}
	}

	return initializationConfig, nil
}

func resourceVirtualEnvironmentVMGetCPUArchitectureValidator() schema.SchemaValidateFunc {
	return validation.StringInSlice([]string{
		"aarch64",
		"x86_64",
	}, false)
}

func resourceVirtualEnvironmentVMGetDiskDeviceObjects(d *schema.ResourceData, m interface{}) (proxmox.CustomStorageDevices, error) {
	diskDevice := d.Get(mkResourceVirtualEnvironmentVMDisk).([]interface{})
	diskDeviceObjects := make(proxmox.CustomStorageDevices, len(diskDevice))
	resource := resourceVirtualEnvironmentVM()

	for i, diskEntry := range diskDevice {
		diskDevice := proxmox.CustomStorageDevice{
			Enabled: true,
		}

		block := diskEntry.(map[string]interface{})
		datastoreID, _ := block[mkResourceVirtualEnvironmentVMDiskDatastoreID].(string)
		fileID, _ := block[mkResourceVirtualEnvironmentVMDiskFileID].(string)
		size, _ := block[mkResourceVirtualEnvironmentVMDiskSize].(int)

		speedBlock, err := getSchemaBlock(resource, d, m, []string{mkResourceVirtualEnvironmentVMDisk, mkResourceVirtualEnvironmentVMDiskSpeed}, 0, false)

		if err != nil {
			return diskDeviceObjects, err
		}

		if fileID != "" {
			diskDevice.Enabled = false
		} else {
			diskDevice.FileVolume = fmt.Sprintf("%s:%d", datastoreID, size)
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

		diskDeviceObjects[i] = diskDevice
	}

	return diskDeviceObjects, nil
}

func resourceVirtualEnvironmentVMGetNetworkDeviceObjects(d *schema.ResourceData, m interface{}) (proxmox.CustomNetworkDevices, error) {
	networkDevice := d.Get(mkResourceVirtualEnvironmentVMNetworkDevice).([]interface{})
	networkDeviceObjects := make(proxmox.CustomNetworkDevices, len(networkDevice))

	for i, networkDeviceEntry := range networkDevice {
		block := networkDeviceEntry.(map[string]interface{})

		bridge, _ := block[mkResourceVirtualEnvironmentVMNetworkDeviceBridge].(string)
		enabled, _ := block[mkResourceVirtualEnvironmentVMNetworkDeviceEnabled].(bool)
		macAddress, _ := block[mkResourceVirtualEnvironmentVMNetworkDeviceMACAddress].(string)
		model, _ := block[mkResourceVirtualEnvironmentVMNetworkDeviceModel].(string)
		rateLimit, _ := block[mkResourceVirtualEnvironmentVMNetworkDeviceRateLimit].(float64)
		vlanID, _ := block[mkResourceVirtualEnvironmentVMNetworkDeviceVLANID].(int)

		device := proxmox.CustomNetworkDevice{
			Enabled: enabled,
			Model:   model,
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

		networkDeviceObjects[i] = device
	}

	return networkDeviceObjects, nil
}

func resourceVirtualEnvironmentVMGetOperatingSystemTypeValidator() schema.SchemaValidateFunc {
	return validation.StringInSlice([]string{
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
	}, false)
}

func resourceVirtualEnvironmentVMGetSerialDeviceList(d *schema.ResourceData, m interface{}) (proxmox.CustomSerialDevices, error) {
	device := d.Get(mkResourceVirtualEnvironmentVMSerialDevice).([]interface{})
	list := make(proxmox.CustomSerialDevices, len(device))

	for i, v := range device {
		block := v.(map[string]interface{})

		device, _ := block[mkResourceVirtualEnvironmentVMSerialDeviceDevice].(string)

		list[i] = device
	}

	return list, nil
}

func resourceVirtualEnvironmentVMGetSerialDeviceValidator() schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {
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
	}
}

func resourceVirtualEnvironmentVMGetVGADeviceObject(d *schema.ResourceData, m interface{}) (*proxmox.CustomVGADevice, error) {
	resource := resourceVirtualEnvironmentVM()

	vgaBlock, err := getSchemaBlock(resource, d, m, []string{mkResourceVirtualEnvironmentVMVGA}, 0, true)

	if err != nil {
		return nil, err
	}

	vgaEnabled := proxmox.CustomBool(vgaBlock[mkResourceVirtualEnvironmentVMAgentEnabled].(bool))
	vgaMemory := vgaBlock[mkResourceVirtualEnvironmentVMVGAMemory].(int)
	vgaType := vgaBlock[mkResourceVirtualEnvironmentVMVGAType].(string)

	vgaDevice := &proxmox.CustomVGADevice{}

	if vgaEnabled {
		if vgaMemory > 0 {
			vgaDevice.Memory = &vgaMemory
		}

		vgaDevice.Type = &vgaType
	} else {
		vgaType = "none"

		vgaDevice = &proxmox.CustomVGADevice{
			Type: &vgaType,
		}
	}

	return vgaDevice, nil
}

func resourceVirtualEnvironmentVMRead(d *schema.ResourceData, m interface{}) error {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()

	if err != nil {
		return err
	}

	nodeName := d.Get(mkResourceVirtualEnvironmentVMNodeName).(string)
	vmID, err := strconv.Atoi(d.Id())

	if err != nil {
		return err
	}

	// Retrieve the entire configuration in order to compare it to the state.
	vmConfig, err := veClient.GetVM(nodeName, vmID)

	if err != nil {
		if strings.Contains(err.Error(), "HTTP 404") ||
			(strings.Contains(err.Error(), "HTTP 500") && strings.Contains(err.Error(), "does not exist")) {
			d.SetId("")

			return nil
		}

		return err
	}

	vmStatus, err := veClient.GetVMStatus(nodeName, vmID)

	if err != nil {
		return err
	}

	return resourceVirtualEnvironmentVMReadCustom(d, m, vmID, vmConfig, vmStatus)
}

func resourceVirtualEnvironmentVMReadCustom(d *schema.ResourceData, m interface{}, vmID int, vmConfig *proxmox.VirtualEnvironmentVMGetResponseData, vmStatus *proxmox.VirtualEnvironmentVMGetStatusResponseData) error {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()

	if err != nil {
		return err
	}

	err = resourceVirtualEnvironmentVMReadPrimitiveValues(d, m, vmID, vmConfig, vmStatus)

	if err != nil {
		return err
	}

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
				agent[mkResourceVirtualEnvironmentVMAgentTrim] = bool(*vmConfig.Agent.TrimClonedDisks)
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
					d.Set(mkResourceVirtualEnvironmentVMAgent, []interface{}{agent})
				}
			} else if len(currentAgent) > 0 ||
				agent[mkResourceVirtualEnvironmentVMAgentEnabled] != dvResourceVirtualEnvironmentVMAgentEnabled ||
				agent[mkResourceVirtualEnvironmentVMAgentTimeout] != dvResourceVirtualEnvironmentVMAgentTimeout ||
				agent[mkResourceVirtualEnvironmentVMAgentTrim] != dvResourceVirtualEnvironmentVMAgentTrim ||
				agent[mkResourceVirtualEnvironmentVMAgentType] != dvResourceVirtualEnvironmentVMAgentType {
				d.Set(mkResourceVirtualEnvironmentVMAgent, []interface{}{agent})
			}
		} else if len(clone) > 0 {
			if len(currentAgent) > 0 {
				d.Set(mkResourceVirtualEnvironmentVMAgent, []interface{}{})
			}
		} else {
			d.Set(mkResourceVirtualEnvironmentVMAgent, []interface{}{})
		}
	}

	// Compare the audio devices to those stored in the state.
	currentAudioDevice := d.Get(mkResourceVirtualEnvironmentVMAudioDevice).([]interface{})

	audioDevices := make([]interface{}, 1)
	audioDevicesArray := []*proxmox.CustomAudioDevice{
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
		d.Set(mkResourceVirtualEnvironmentVMAudioDevice, audioDevices[:audioDevicesCount])
	}

	// Compare the IDE devices to the CDROM and cloud-init configurations stored in the state.
	if vmConfig.IDEDevice2 != nil {
		if *vmConfig.IDEDevice2.Media == "cdrom" {
			if strings.Contains(vmConfig.IDEDevice2.FileVolume, fmt.Sprintf("vm-%d-cloudinit", vmID)) {
				d.Set(mkResourceVirtualEnvironmentVMCDROM, []interface{}{})
			} else {
				d.Set(mkResourceVirtualEnvironmentVMInitialization, []interface{}{})

				cdrom := make([]interface{}, 1)
				cdromBlock := map[string]interface{}{}

				cdromBlock[mkResourceVirtualEnvironmentVMCDROMEnabled] = true
				cdromBlock[mkResourceVirtualEnvironmentVMCDROMFileID] = vmConfig.IDEDevice2.FileVolume

				cdrom[0] = cdromBlock

				currentCDROM := d.Get(mkResourceVirtualEnvironmentVMCDROM).([]interface{})

				if len(clone) == 0 || len(currentCDROM) > 0 {
					d.Set(mkResourceVirtualEnvironmentVMCDROM, cdrom)
				}
			}
		} else {
			d.Set(mkResourceVirtualEnvironmentVMCDROM, []interface{}{})
			d.Set(mkResourceVirtualEnvironmentVMInitialization, []interface{}{})
		}
	} else {
		d.Set(mkResourceVirtualEnvironmentVMCDROM, []interface{}{})
		d.Set(mkResourceVirtualEnvironmentVMInitialization, []interface{}{})
	}

	// Compare the CPU configuration to the one stored in the state.
	cpu := map[string]interface{}{}

	if vmConfig.CPUArchitecture != nil {
		cpu[mkResourceVirtualEnvironmentVMCPUArchitecture] = *vmConfig.CPUArchitecture
	} else {
		// Default value of "arch" is "" according to the API documentation.
		// However, assume the provider's default value as a workaround when the root account is not being used.
		if veClient.Username != proxmox.DefaultRootAccount {
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
			d.Set(mkResourceVirtualEnvironmentVMCPU, []interface{}{cpu})
		}
	} else if len(currentCPU) > 0 ||
		cpu[mkResourceVirtualEnvironmentVMCPUArchitecture] != dvResourceVirtualEnvironmentVMCPUArchitecture ||
		cpu[mkResourceVirtualEnvironmentVMCPUCores] != dvResourceVirtualEnvironmentVMCPUCores ||
		len(cpu[mkResourceVirtualEnvironmentVMCPUFlags].([]interface{})) > 0 ||
		cpu[mkResourceVirtualEnvironmentVMCPUHotplugged] != dvResourceVirtualEnvironmentVMCPUHotplugged ||
		cpu[mkResourceVirtualEnvironmentVMCPUSockets] != dvResourceVirtualEnvironmentVMCPUSockets ||
		cpu[mkResourceVirtualEnvironmentVMCPUType] != dvResourceVirtualEnvironmentVMCPUType ||
		cpu[mkResourceVirtualEnvironmentVMCPUUnits] != dvResourceVirtualEnvironmentVMCPUUnits {
		d.Set(mkResourceVirtualEnvironmentVMCPU, []interface{}{cpu})
	}

	// Compare the disks to those stored in the state.
	currentDisk := d.Get(mkResourceVirtualEnvironmentVMDisk).([]interface{})

	diskList := []interface{}{}
	diskObjects := []*proxmox.CustomStorageDevice{
		vmConfig.SCSIDevice0,
		vmConfig.SCSIDevice1,
		vmConfig.SCSIDevice2,
		vmConfig.SCSIDevice3,
		vmConfig.SCSIDevice4,
		vmConfig.SCSIDevice5,
		vmConfig.SCSIDevice6,
		vmConfig.SCSIDevice7,
		vmConfig.SCSIDevice8,
		vmConfig.SCSIDevice9,
		vmConfig.SCSIDevice10,
		vmConfig.SCSIDevice11,
		vmConfig.SCSIDevice12,
		vmConfig.SCSIDevice13,
	}

	for di, dd := range diskObjects {
		disk := map[string]interface{}{}

		if dd == nil {
			continue
		}

		fileIDParts := strings.Split(dd.FileVolume, ":")

		disk[mkResourceVirtualEnvironmentVMDiskDatastoreID] = fileIDParts[0]

		if len(currentDisk) > di {
			currentDiskEntry := currentDisk[di].(map[string]interface{})

			disk[mkResourceVirtualEnvironmentVMDiskFileFormat] = currentDiskEntry[mkResourceVirtualEnvironmentVMDiskFileFormat]
			disk[mkResourceVirtualEnvironmentVMDiskFileID] = currentDiskEntry[mkResourceVirtualEnvironmentVMDiskFileID]
		}

		diskSize := 0

		var err error

		if dd.Size != nil {
			if strings.HasSuffix(*dd.Size, "T") {
				diskSize, err = strconv.Atoi(strings.TrimSuffix(*dd.Size, "T"))

				if err != nil {
					return err
				}

				diskSize = int(math.Ceil(float64(diskSize) * 1024))
			} else if strings.HasSuffix(*dd.Size, "G") {
				diskSize, err = strconv.Atoi(strings.TrimSuffix(*dd.Size, "G"))

				if err != nil {
					return err
				}
			} else if strings.HasSuffix(*dd.Size, "M") {
				diskSize, err = strconv.Atoi(strings.TrimSuffix(*dd.Size, "M"))

				if err != nil {
					return err
				}

				diskSize = int(math.Ceil(float64(diskSize) / 1024))
			} else {
				return fmt.Errorf("Cannot parse storage size \"%s\"", *dd.Size)
			}
		}

		disk[mkResourceVirtualEnvironmentVMDiskSize] = diskSize

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

		diskList = append(diskList, disk)
	}

	if len(clone) > 0 {
		if len(currentDisk) > 0 {
			d.Set(mkResourceVirtualEnvironmentVMDisk, diskList)
		}
	} else if len(currentDisk) > 0 || len(diskList) > 0 {
		d.Set(mkResourceVirtualEnvironmentVMDisk, diskList)
	}

	// Compare the initialization configuration to the one stored in the state.
	initialization := map[string]interface{}{}

	if vmConfig.IDEDevice2 != nil {
		if *vmConfig.IDEDevice2.Media == "cdrom" {
			if strings.Contains(vmConfig.IDEDevice2.FileVolume, fmt.Sprintf("vm-%d-cloudinit", vmID)) {
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

		initialization[mkResourceVirtualEnvironmentVMInitializationDNS] = []interface{}{initializationDNS}
	}

	ipConfigLast := -1
	ipConfigObjects := []*proxmox.CustomCloudInitIPConfig{
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

				ipConfigItem[mkResourceVirtualEnvironmentVMInitializationIPConfigIPv4] = []interface{}{ipv4}
			} else {
				ipConfigItem[mkResourceVirtualEnvironmentVMInitializationIPConfigIPv4] = []interface{}{}
			}

			if ipConfig.GatewayIPv6 != nil || ipConfig.IPv6 != nil {
				ipv6 := map[string]interface{}{}

				if ipConfig.IPv4 != nil {
					ipv6[mkResourceVirtualEnvironmentVMInitializationIPConfigIPv4Address] = *ipConfig.IPv6
				} else {
					ipv6[mkResourceVirtualEnvironmentVMInitializationIPConfigIPv4Address] = ""
				}

				if ipConfig.GatewayIPv4 != nil {
					ipv6[mkResourceVirtualEnvironmentVMInitializationIPConfigIPv4Gateway] = *ipConfig.GatewayIPv6
				} else {
					ipv6[mkResourceVirtualEnvironmentVMInitializationIPConfigIPv4Gateway] = ""
				}

				ipConfigItem[mkResourceVirtualEnvironmentVMInitializationIPConfigIPv6] = []interface{}{ipv6}
			} else {
				ipConfigItem[mkResourceVirtualEnvironmentVMInitializationIPConfigIPv6] = []interface{}{}
			}
		} else {
			ipConfigItem[mkResourceVirtualEnvironmentVMInitializationIPConfigIPv4] = []interface{}{}
			ipConfigItem[mkResourceVirtualEnvironmentVMInitializationIPConfigIPv6] = []interface{}{}
		}

		ipConfigList[ipConfigIndex] = ipConfigItem
	}

	initialization[mkResourceVirtualEnvironmentVMInitializationIPConfig] = ipConfigList[:ipConfigLast+1]

	if vmConfig.CloudInitPassword != nil || vmConfig.CloudInitSSHKeys != nil || vmConfig.CloudInitUsername != nil {
		initializationUserAccount := map[string]interface{}{}

		if vmConfig.CloudInitSSHKeys != nil {
			initializationUserAccount[mkResourceVirtualEnvironmentVMInitializationUserAccountKeys] = []string(*vmConfig.CloudInitSSHKeys)
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

		initialization[mkResourceVirtualEnvironmentVMInitializationUserAccount] = []interface{}{initializationUserAccount}
	}

	if vmConfig.CloudInitFiles != nil {
		if vmConfig.CloudInitFiles.UserVolume != nil {
			initialization[mkResourceVirtualEnvironmentVMInitializationUserDataFileID] = *vmConfig.CloudInitFiles.UserVolume
		} else {
			initialization[mkResourceVirtualEnvironmentVMInitializationUserDataFileID] = ""
		}
	} else {
		initialization[mkResourceVirtualEnvironmentVMInitializationUserDataFileID] = ""
	}

	currentInitialization := d.Get(mkResourceVirtualEnvironmentVMInitialization).([]interface{})

	if len(clone) > 0 {
		if len(currentInitialization) > 0 {
			if len(initialization) > 0 {
				d.Set(mkResourceVirtualEnvironmentVMInitialization, []interface{}{initialization})
			} else {
				d.Set(mkResourceVirtualEnvironmentVMInitialization, []interface{}{})
			}
		}
	} else if len(initialization) > 0 {
		d.Set(mkResourceVirtualEnvironmentVMInitialization, []interface{}{initialization})
	} else {
		d.Set(mkResourceVirtualEnvironmentVMInitialization, []interface{}{})
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
			d.Set(mkResourceVirtualEnvironmentVMMemory, []interface{}{memory})
		}
	} else if len(currentMemory) > 0 ||
		memory[mkResourceVirtualEnvironmentVMMemoryDedicated] != dvResourceVirtualEnvironmentVMMemoryDedicated ||
		memory[mkResourceVirtualEnvironmentVMMemoryFloating] != dvResourceVirtualEnvironmentVMMemoryFloating ||
		memory[mkResourceVirtualEnvironmentVMMemoryShared] != dvResourceVirtualEnvironmentVMMemoryShared {
		d.Set(mkResourceVirtualEnvironmentVMMemory, []interface{}{memory})
	}

	// Compare the network devices to those stored in the state.
	currentNetworkDeviceList := d.Get(mkResourceVirtualEnvironmentVMNetworkDevice).([]interface{})

	macAddresses := make([]interface{}, 8)
	networkDeviceLast := -1
	networkDeviceList := make([]interface{}, 8)
	networkDeviceObjects := []*proxmox.CustomNetworkDevice{
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
		} else {
			macAddresses[ni] = ""
			networkDevice[mkResourceVirtualEnvironmentVMNetworkDeviceEnabled] = false
		}

		networkDeviceList[ni] = networkDevice
	}

	if len(clone) > 0 {
		if len(currentNetworkDeviceList) > 0 {
			d.Set(mkResourceVirtualEnvironmentVMMACAddresses, macAddresses[0:len(currentNetworkDeviceList)])
			d.Set(mkResourceVirtualEnvironmentVMNetworkDevice, networkDeviceList[:networkDeviceLast+1])
		}
	} else {
		d.Set(mkResourceVirtualEnvironmentVMMACAddresses, macAddresses[0:len(currentNetworkDeviceList)])

		if len(currentNetworkDeviceList) > 0 || networkDeviceLast > -1 {
			d.Set(mkResourceVirtualEnvironmentVMNetworkDevice, networkDeviceList[:networkDeviceLast+1])
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
			d.Set(mkResourceVirtualEnvironmentVMOperatingSystem, []interface{}{operatingSystem})
		}
	} else if len(currentOperatingSystem) > 0 ||
		operatingSystem[mkResourceVirtualEnvironmentVMOperatingSystemType] != dvResourceVirtualEnvironmentVMOperatingSystemType {
		d.Set(mkResourceVirtualEnvironmentVMOperatingSystem, []interface{}{operatingSystem})
	} else {
		d.Set(mkResourceVirtualEnvironmentVMOperatingSystem, []interface{}{})
	}

	// Compare the pool ID to the value stored in the state.
	currentPoolID := d.Get(mkResourceVirtualEnvironmentVMPoolID).(string)

	if len(clone) == 0 || currentPoolID != dvResourceVirtualEnvironmentVMPoolID {
		if vmConfig.PoolID != nil {
			d.Set(mkResourceVirtualEnvironmentVMPoolID, *vmConfig.PoolID)
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
		d.Set(mkResourceVirtualEnvironmentVMSerialDevice, serialDevices[:serialDevicesCount])
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
			d.Set(mkResourceVirtualEnvironmentVMVGA, []interface{}{vga})
		}
	} else if len(currentVGA) > 0 ||
		vga[mkResourceVirtualEnvironmentVMVGAEnabled] != dvResourceVirtualEnvironmentVMVGAEnabled ||
		vga[mkResourceVirtualEnvironmentVMVGAMemory] != dvResourceVirtualEnvironmentVMVGAMemory ||
		vga[mkResourceVirtualEnvironmentVMVGAType] != dvResourceVirtualEnvironmentVMVGAType {
		d.Set(mkResourceVirtualEnvironmentVMVGA, []interface{}{vga})
	} else {
		d.Set(mkResourceVirtualEnvironmentVMVGA, []interface{}{})
	}

	return resourceVirtualEnvironmentVMReadNetworkValues(d, m, vmID, vmConfig)
}

func resourceVirtualEnvironmentVMReadNetworkValues(d *schema.ResourceData, m interface{}, vmID int, vmConfig *proxmox.VirtualEnvironmentVMGetResponseData) error {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()

	if err != nil {
		return err
	}

	nodeName := d.Get(mkResourceVirtualEnvironmentVMNodeName).(string)
	started := d.Get(mkResourceVirtualEnvironmentVMStarted).(bool)

	ipv4Addresses := []interface{}{}
	ipv6Addresses := []interface{}{}
	networkInterfaceNames := []interface{}{}

	if started {
		if vmConfig.Agent != nil && vmConfig.Agent.Enabled != nil && *vmConfig.Agent.Enabled {
			resource := resourceVirtualEnvironmentVM()
			agentBlock, err := getSchemaBlock(resource, d, m, []string{mkResourceVirtualEnvironmentVMAgent}, 0, true)

			if err != nil {
				return err
			}

			agentTimeout, err := time.ParseDuration(agentBlock[mkResourceVirtualEnvironmentVMAgentTimeout].(string))

			if err != nil {
				return err
			}

			macAddresses := []interface{}{}
			networkInterfaces, err := veClient.WaitForNetworkInterfacesFromVMAgent(nodeName, vmID, int(agentTimeout.Seconds()), 5, true)

			if err == nil && networkInterfaces.Result != nil {
				ipv4Addresses = make([]interface{}, len(*networkInterfaces.Result))
				ipv6Addresses = make([]interface{}, len(*networkInterfaces.Result))
				macAddresses = make([]interface{}, len(*networkInterfaces.Result))
				networkInterfaceNames = make([]interface{}, len(*networkInterfaces.Result))

				for ri, rv := range *networkInterfaces.Result {
					rvIPv4Addresses := []interface{}{}
					rvIPv6Addresses := []interface{}{}

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

			d.Set(mkResourceVirtualEnvironmentVMMACAddresses, macAddresses)
		}
	}

	d.Set(mkResourceVirtualEnvironmentVMIPv4Addresses, ipv4Addresses)
	d.Set(mkResourceVirtualEnvironmentVMIPv6Addresses, ipv6Addresses)
	d.Set(mkResourceVirtualEnvironmentVMNetworkInterfaceNames, networkInterfaceNames)

	return nil
}

func resourceVirtualEnvironmentVMReadPrimitiveValues(d *schema.ResourceData, m interface{}, vmID int, vmConfig *proxmox.VirtualEnvironmentVMGetResponseData, vmStatus *proxmox.VirtualEnvironmentVMGetStatusResponseData) error {
	clone := d.Get(mkResourceVirtualEnvironmentVMClone).([]interface{})
	currentACPI := d.Get(mkResourceVirtualEnvironmentVMACPI).(bool)

	if len(clone) == 0 || currentACPI != dvResourceVirtualEnvironmentVMACPI {
		if vmConfig.ACPI != nil {
			d.Set(mkResourceVirtualEnvironmentVMACPI, bool(*vmConfig.ACPI))
		} else {
			// Default value of "acpi" is "1" according to the API documentation.
			d.Set(mkResourceVirtualEnvironmentVMACPI, true)
		}
	}

	currentBIOS := d.Get(mkResourceVirtualEnvironmentVMBIOS).(string)

	if len(clone) == 0 || currentBIOS != dvResourceVirtualEnvironmentVMBIOS {
		if vmConfig.BIOS != nil {
			d.Set(mkResourceVirtualEnvironmentVMBIOS, *vmConfig.BIOS)
		} else {
			// Default value of "bios" is "seabios" according to the API documentation.
			d.Set(mkResourceVirtualEnvironmentVMBIOS, "seabios")
		}
	}

	currentDescription := d.Get(mkResourceVirtualEnvironmentVMDescription).(string)

	if len(clone) == 0 || currentDescription != dvResourceVirtualEnvironmentVMDescription {
		if vmConfig.Description != nil {
			d.Set(mkResourceVirtualEnvironmentVMDescription, *vmConfig.Description)
		} else {
			// Default value of "description" is "" according to the API documentation.
			d.Set(mkResourceVirtualEnvironmentVMDescription, "")
		}
	}

	currentKeyboardLayout := d.Get(mkResourceVirtualEnvironmentVMKeyboardLayout).(string)

	if len(clone) == 0 || currentKeyboardLayout != dvResourceVirtualEnvironmentVMKeyboardLayout {
		if vmConfig.KeyboardLayout != nil {
			d.Set(mkResourceVirtualEnvironmentVMKeyboardLayout, *vmConfig.KeyboardLayout)
		} else {
			// Default value of "keyboard" is "" according to the API documentation.
			d.Set(mkResourceVirtualEnvironmentVMKeyboardLayout, "")
		}
	}

	currentName := d.Get(mkResourceVirtualEnvironmentVMName).(string)

	if len(clone) == 0 || currentName != dvResourceVirtualEnvironmentVMName {
		if vmConfig.Name != nil {
			d.Set(mkResourceVirtualEnvironmentVMName, *vmConfig.Name)
		} else {
			// Default value of "name" is "" according to the API documentation.
			d.Set(mkResourceVirtualEnvironmentVMName, "")
		}
	}

	d.Set(mkResourceVirtualEnvironmentVMStarted, vmStatus.Status == "running")

	currentTabletDevice := d.Get(mkResourceVirtualEnvironmentVMTabletDevice).(bool)

	if len(clone) == 0 || currentTabletDevice != dvResourceVirtualEnvironmentVMTabletDevice {
		if vmConfig.TabletDeviceEnabled != nil {
			d.Set(mkResourceVirtualEnvironmentVMTabletDevice, bool(*vmConfig.TabletDeviceEnabled))
		} else {
			// Default value of "tablet" is "1" according to the API documentation.
			d.Set(mkResourceVirtualEnvironmentVMTabletDevice, true)
		}
	}

	currentTemplate := d.Get(mkResourceVirtualEnvironmentVMTemplate).(bool)

	if len(clone) == 0 || currentTemplate != dvResourceVirtualEnvironmentVMTemplate {
		if vmConfig.Template != nil {
			d.Set(mkResourceVirtualEnvironmentVMTemplate, bool(*vmConfig.Template))
		} else {
			// Default value of "template" is "0" according to the API documentation.
			d.Set(mkResourceVirtualEnvironmentVMTemplate, false)
		}
	}

	return nil
}

func resourceVirtualEnvironmentVMUpdate(d *schema.ResourceData, m interface{}) error {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()

	if err != nil {
		return err
	}

	nodeName := d.Get(mkResourceVirtualEnvironmentVMNodeName).(string)
	rebootRequired := false

	vmID, err := strconv.Atoi(d.Id())

	if err != nil {
		return err
	}

	updateBody := &proxmox.VirtualEnvironmentVMUpdateRequestBody{
		IDEDevices: proxmox.CustomStorageDevices{
			proxmox.CustomStorageDevice{
				Enabled: false,
			},
			proxmox.CustomStorageDevice{
				Enabled: false,
			},
			proxmox.CustomStorageDevice{
				Enabled: false,
			},
		},
	}

	delete := []string{}
	resource := resourceVirtualEnvironmentVM()

	// Retrieve the entire configuration as we need to process certain values.
	vmConfig, err := veClient.GetVM(nodeName, vmID)

	if err != nil {
		return err
	}

	// Prepare the new primitive configuration values.
	if d.HasChange(mkResourceVirtualEnvironmentVMACPI) {
		acpi := proxmox.CustomBool(d.Get(mkResourceVirtualEnvironmentVMACPI).(bool))
		updateBody.ACPI = &acpi
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

	if d.HasChange(mkResourceVirtualEnvironmentVMKeyboardLayout) {
		keyboardLayout := d.Get(mkResourceVirtualEnvironmentVMKeyboardLayout).(string)
		updateBody.KeyboardLayout = &keyboardLayout
		rebootRequired = true
	}

	name := d.Get(mkResourceVirtualEnvironmentVMName).(string)
	updateBody.Name = &name

	if d.HasChange(mkResourceVirtualEnvironmentVMTabletDevice) {
		tabletDevice := proxmox.CustomBool(d.Get(mkResourceVirtualEnvironmentVMTabletDevice).(bool))
		updateBody.TabletDeviceEnabled = &tabletDevice
		rebootRequired = true
	}

	template := proxmox.CustomBool(d.Get(mkResourceVirtualEnvironmentVMTemplate).(bool))

	if d.HasChange(mkResourceVirtualEnvironmentVMTemplate) {
		updateBody.Template = &template
		rebootRequired = true
	}

	// Prepare the new agent configuration.
	if d.HasChange(mkResourceVirtualEnvironmentVMAgent) {
		agentBlock, err := getSchemaBlock(resource, d, m, []string{mkResourceVirtualEnvironmentVMAgent}, 0, true)

		if err != nil {
			return err
		}

		agentEnabled := proxmox.CustomBool(agentBlock[mkResourceVirtualEnvironmentVMAgentEnabled].(bool))
		agentTrim := proxmox.CustomBool(agentBlock[mkResourceVirtualEnvironmentVMAgentTrim].(bool))
		agentType := agentBlock[mkResourceVirtualEnvironmentVMAgentType].(string)

		updateBody.Agent = &proxmox.CustomAgent{
			Enabled:         &agentEnabled,
			TrimClonedDisks: &agentTrim,
			Type:            &agentType,
		}

		rebootRequired = true
	}

	// Prepare the new audio devices.
	if d.HasChange(mkResourceVirtualEnvironmentVMAudioDevice) {
		updateBody.AudioDevices, err = resourceVirtualEnvironmentVMGetAudioDeviceList(d, m)

		if err != nil {
			return err
		}

		for i := 0; i < len(updateBody.AudioDevices); i++ {
			if !updateBody.AudioDevices[i].Enabled {
				delete = append(delete, fmt.Sprintf("audio%d", i))
			}
		}

		for i := len(updateBody.AudioDevices); i < maxResourceVirtualEnvironmentVMAudioDevices; i++ {
			delete = append(delete, fmt.Sprintf("audio%d", i))
		}

		rebootRequired = true
	}

	// Prepare the new CDROM configuration.
	if d.HasChange(mkResourceVirtualEnvironmentVMCDROM) {
		cdromBlock, err := getSchemaBlock(resource, d, m, []string{mkResourceVirtualEnvironmentVMCDROM}, 0, true)

		if err != nil {
			return err
		}

		cdromEnabled := cdromBlock[mkResourceVirtualEnvironmentVMCDROMEnabled].(bool)
		cdromFileID := cdromBlock[mkResourceVirtualEnvironmentVMCDROMFileID].(string)

		if cdromFileID == "" {
			cdromFileID = "cdrom"
		}

		cdromMedia := "cdrom"

		updateBody.IDEDevices[2] = proxmox.CustomStorageDevice{
			Enabled:    cdromEnabled,
			FileVolume: cdromFileID,
			Media:      &cdromMedia,
		}
	}

	// Prepare the new CPU configuration.
	if d.HasChange(mkResourceVirtualEnvironmentVMCPU) {
		cpuBlock, err := getSchemaBlock(resource, d, m, []string{mkResourceVirtualEnvironmentVMCPU}, 0, true)

		if err != nil {
			return err
		}

		cpuArchitecture := cpuBlock[mkResourceVirtualEnvironmentVMCPUArchitecture].(string)
		cpuCores := cpuBlock[mkResourceVirtualEnvironmentVMCPUCores].(int)
		cpuFlags := cpuBlock[mkResourceVirtualEnvironmentVMCPUFlags].([]interface{})
		cpuHotplugged := cpuBlock[mkResourceVirtualEnvironmentVMCPUHotplugged].(int)
		cpuSockets := cpuBlock[mkResourceVirtualEnvironmentVMCPUSockets].(int)
		cpuType := cpuBlock[mkResourceVirtualEnvironmentVMCPUType].(string)
		cpuUnits := cpuBlock[mkResourceVirtualEnvironmentVMCPUUnits].(int)

		// Only the root account is allowed to change the CPU architecture, which makes this check necessary.
		if veClient.Username == proxmox.DefaultRootAccount || cpuArchitecture != dvResourceVirtualEnvironmentVMCPUArchitecture {
			updateBody.CPUArchitecture = &cpuArchitecture
		}

		updateBody.CPUCores = &cpuCores
		updateBody.CPUSockets = &cpuSockets
		updateBody.CPUUnits = &cpuUnits

		if cpuHotplugged > 0 {
			updateBody.VirtualCPUCount = &cpuHotplugged
		} else {
			delete = append(delete, "vcpus")
		}

		cpuFlagsConverted := make([]string, len(cpuFlags))

		for fi, flag := range cpuFlags {
			cpuFlagsConverted[fi] = flag.(string)
		}

		updateBody.CPUEmulation = &proxmox.CustomCPUEmulation{
			Flags: &cpuFlagsConverted,
			Type:  cpuType,
		}

		rebootRequired = true
	}

	// Prepare the new disk device configuration.
	if d.HasChange(mkResourceVirtualEnvironmentVMDisk) {
		diskDeviceObjects, err := resourceVirtualEnvironmentVMGetDiskDeviceObjects(d, m)

		if err != nil {
			return err
		}

		scsiDevices := []*proxmox.CustomStorageDevice{
			vmConfig.SCSIDevice0,
			vmConfig.SCSIDevice1,
			vmConfig.SCSIDevice2,
			vmConfig.SCSIDevice3,
			vmConfig.SCSIDevice4,
			vmConfig.SCSIDevice5,
			vmConfig.SCSIDevice6,
			vmConfig.SCSIDevice7,
			vmConfig.SCSIDevice8,
			vmConfig.SCSIDevice9,
			vmConfig.SCSIDevice10,
			vmConfig.SCSIDevice11,
			vmConfig.SCSIDevice12,
			vmConfig.SCSIDevice13,
		}

		updateBody.SCSIDevices = make(proxmox.CustomStorageDevices, len(diskDeviceObjects))

		for di, do := range diskDeviceObjects {
			if scsiDevices[di] == nil {
				return fmt.Errorf("Missing SCSI device %d (scsi%d)", di, di)
			}

			updateBody.SCSIDevices[di] = *scsiDevices[di]
			updateBody.SCSIDevices[di].BurstableReadSpeedMbps = do.BurstableReadSpeedMbps
			updateBody.SCSIDevices[di].BurstableWriteSpeedMbps = do.BurstableWriteSpeedMbps
			updateBody.SCSIDevices[di].MaxReadSpeedMbps = do.MaxReadSpeedMbps
			updateBody.SCSIDevices[di].MaxWriteSpeedMbps = do.MaxWriteSpeedMbps
		}

		rebootRequired = true
	}

	// Prepare the new cloud-init configuration.
	if d.HasChange(mkResourceVirtualEnvironmentVMInitialization) {
		initializationConfig, err := resourceVirtualEnvironmentVMGetCloudInitConfig(d, m)

		if err != nil {
			return err
		}

		updateBody.CloudInitConfig = initializationConfig

		if updateBody.CloudInitConfig != nil {
			initialization := d.Get(mkResourceVirtualEnvironmentVMInitialization).([]interface{})
			initializationBlock := initialization[0].(map[string]interface{})
			initializationDatastoreID := initializationBlock[mkResourceVirtualEnvironmentVMInitializationDatastoreID].(string)

			cdromMedia := "cdrom"

			updateBody.IDEDevices[2] = proxmox.CustomStorageDevice{
				Enabled:    true,
				FileVolume: fmt.Sprintf("%s:cloudinit", initializationDatastoreID),
				Media:      &cdromMedia,
			}

			if vmConfig.IDEDevice2 != nil {
				if strings.Contains(vmConfig.IDEDevice2.FileVolume, fmt.Sprintf("vm-%d-cloudinit", vmID)) {
					updateBody.IDEDevices[2].Enabled = false
				}
			}
		}

		rebootRequired = true
	}

	// Prepare the new memory configuration.
	if d.HasChange(mkResourceVirtualEnvironmentVMMemory) {
		memoryBlock, err := getSchemaBlock(resource, d, m, []string{mkResourceVirtualEnvironmentVMMemory}, 0, true)

		if err != nil {
			return err
		}

		memoryDedicated := memoryBlock[mkResourceVirtualEnvironmentVMMemoryDedicated].(int)
		memoryFloating := memoryBlock[mkResourceVirtualEnvironmentVMMemoryFloating].(int)
		memoryShared := memoryBlock[mkResourceVirtualEnvironmentVMMemoryShared].(int)

		updateBody.DedicatedMemory = &memoryDedicated
		updateBody.FloatingMemory = &memoryFloating

		if memoryShared > 0 {
			memorySharedName := fmt.Sprintf("vm-%d-ivshmem", vmID)

			updateBody.SharedMemory = &proxmox.CustomSharedMemory{
				Name: &memorySharedName,
				Size: memoryShared,
			}
		}

		rebootRequired = true
	}

	// Prepare the new network device configuration.
	if d.HasChange(mkResourceVirtualEnvironmentVMNetworkDevice) {
		updateBody.NetworkDevices, err = resourceVirtualEnvironmentVMGetNetworkDeviceObjects(d, m)

		if err != nil {
			return err
		}

		for i := 0; i < len(updateBody.NetworkDevices); i++ {
			if !updateBody.NetworkDevices[i].Enabled {
				delete = append(delete, fmt.Sprintf("net%d", i))
			}
		}

		for i := len(updateBody.NetworkDevices); i < maxResourceVirtualEnvironmentVMNetworkDevices; i++ {
			delete = append(delete, fmt.Sprintf("net%d", i))
		}

		rebootRequired = true
	}

	// Prepare the new operating system configuration.
	if d.HasChange(mkResourceVirtualEnvironmentVMOperatingSystem) {
		operatingSystem, err := getSchemaBlock(resource, d, m, []string{mkResourceVirtualEnvironmentVMOperatingSystem}, 0, true)

		if err != nil {
			return err
		}

		operatingSystemType := operatingSystem[mkResourceVirtualEnvironmentVMOperatingSystemType].(string)

		updateBody.OSType = &operatingSystemType

		rebootRequired = true
	}

	// Prepare the new serial devices.
	if d.HasChange(mkResourceVirtualEnvironmentVMSerialDevice) {
		updateBody.SerialDevices, err = resourceVirtualEnvironmentVMGetSerialDeviceList(d, m)

		if err != nil {
			return err
		}

		for i := len(updateBody.SerialDevices); i < maxResourceVirtualEnvironmentVMSerialDevices; i++ {
			delete = append(delete, fmt.Sprintf("serial%d", i))
		}

		rebootRequired = true
	}

	// Prepare the new VGA configuration.
	if d.HasChange(mkResourceVirtualEnvironmentVMVGA) {
		updateBody.VGADevice, err = resourceVirtualEnvironmentVMGetVGADeviceObject(d, m)

		if err != nil {
			return err
		}

		rebootRequired = true
	}

	// Update the configuration now that everything has been prepared.
	updateBody.Delete = delete

	err = veClient.UpdateVM(nodeName, vmID, updateBody)

	if err != nil {
		return err
	}

	// Determine if the state of the virtual machine state needs to be changed.
	started := d.Get(mkResourceVirtualEnvironmentVMStarted).(bool)

	if d.HasChange(mkResourceVirtualEnvironmentVMStarted) && !bool(template) {
		if started {
			err = veClient.StartVM(nodeName, vmID)

			if err != nil {
				return err
			}

			err = veClient.WaitForVMState(nodeName, vmID, "running", 120, 5)

			if err != nil {
				return err
			}
		} else {
			forceStop := proxmox.CustomBool(true)
			shutdownTimeout := 300

			err = veClient.ShutdownVM(nodeName, vmID, &proxmox.VirtualEnvironmentVMShutdownRequestBody{
				ForceStop: &forceStop,
				Timeout:   &shutdownTimeout,
			})

			if err != nil {
				return err
			}

			err = veClient.WaitForVMState(nodeName, vmID, "stopped", 30, 5)

			if err != nil {
				return err
			}

			rebootRequired = false
		}
	}

	// Reboot the virtual machine, if required.
	if !bool(template) && rebootRequired {
		rebootTimeout := 300

		err = veClient.RebootVM(nodeName, vmID, &proxmox.VirtualEnvironmentVMRebootRequestBody{
			Timeout: &rebootTimeout,
		})

		if err != nil {
			return err
		}

		// Wait for the agent to unpublish the network interfaces, if it's enabled.
		if vmConfig.Agent != nil && vmConfig.Agent.Enabled != nil && *vmConfig.Agent.Enabled {
			err = veClient.WaitForNoNetworkInterfacesFromVMAgent(nodeName, vmID, 300, 5)

			if err != nil {
				return err
			}
		}
	}

	return resourceVirtualEnvironmentVMRead(d, m)
}

func resourceVirtualEnvironmentVMDelete(d *schema.ResourceData, m interface{}) error {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()

	if err != nil {
		return err
	}

	nodeName := d.Get(mkResourceVirtualEnvironmentVMNodeName).(string)
	vmID, err := strconv.Atoi(d.Id())

	if err != nil {
		return err
	}

	// Shut down the virtual machine before deleting it.
	status, err := veClient.GetVMStatus(nodeName, vmID)

	if err != nil {
		return err
	}

	if status.Status != "stopped" {
		forceStop := proxmox.CustomBool(true)
		shutdownTimeout := 300

		err = veClient.ShutdownVM(nodeName, vmID, &proxmox.VirtualEnvironmentVMShutdownRequestBody{
			ForceStop: &forceStop,
			Timeout:   &shutdownTimeout,
		})

		if err != nil {
			return err
		}

		err = veClient.WaitForVMState(nodeName, vmID, "stopped", 30, 5)

		if err != nil {
			return err
		}
	}

	err = veClient.DeleteVM(nodeName, vmID)

	if err != nil {
		if strings.Contains(err.Error(), "HTTP 404") ||
			(strings.Contains(err.Error(), "HTTP 500") && strings.Contains(err.Error(), "does not exist")) {
			d.SetId("")

			return nil
		}

		return err
	}

	// Wait for the state to become unavailable as that clearly indicates the destruction of the VM.
	err = veClient.WaitForVMState(nodeName, vmID, "", 60, 2)

	if err == nil {
		return fmt.Errorf("Failed to delete VM \"%d\"", vmID)
	}

	d.SetId("")

	return nil
}
