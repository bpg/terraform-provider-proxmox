/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmoxtf

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/danitso/terraform-provider-proxmox/proxmox"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

const (
	dvResourceVirtualEnvironmentVMAgentEnabled                 = false
	dvResourceVirtualEnvironmentVMAgentTrim                    = false
	dvResourceVirtualEnvironmentVMAgentType                    = "virtio"
	dvResourceVirtualEnvironmentVMCDROMEnabled                 = false
	dvResourceVirtualEnvironmentVMCDROMFileID                  = ""
	dvResourceVirtualEnvironmentVMCloudInitDNSDomain           = ""
	dvResourceVirtualEnvironmentVMCloudInitDNSServer           = ""
	dvResourceVirtualEnvironmentVMCloudInitUserAccountPassword = ""
	dvResourceVirtualEnvironmentVMCloudInitUserDataFileID      = ""
	dvResourceVirtualEnvironmentVMCPUCores                     = 1
	dvResourceVirtualEnvironmentVMCPUHotplugged                = 0
	dvResourceVirtualEnvironmentVMCPUSockets                   = 1
	dvResourceVirtualEnvironmentVMDescription                  = ""
	dvResourceVirtualEnvironmentVMDiskDatastoreID              = "local-lvm"
	dvResourceVirtualEnvironmentVMDiskFileFormat               = "qcow2"
	dvResourceVirtualEnvironmentVMDiskFileID                   = ""
	dvResourceVirtualEnvironmentVMDiskSize                     = 8
	dvResourceVirtualEnvironmentVMDiskSpeedRead                = 0
	dvResourceVirtualEnvironmentVMDiskSpeedReadBurstable       = 0
	dvResourceVirtualEnvironmentVMDiskSpeedWrite               = 0
	dvResourceVirtualEnvironmentVMDiskSpeedWriteBurstable      = 0
	dvResourceVirtualEnvironmentVMKeyboardLayout               = "en-us"
	dvResourceVirtualEnvironmentVMMemoryDedicated              = 512
	dvResourceVirtualEnvironmentVMMemoryFloating               = 0
	dvResourceVirtualEnvironmentVMMemoryShared                 = 0
	dvResourceVirtualEnvironmentVMName                         = ""
	dvResourceVirtualEnvironmentVMNetworkDeviceBridge          = "vmbr0"
	dvResourceVirtualEnvironmentVMNetworkDeviceEnabled         = true
	dvResourceVirtualEnvironmentVMNetworkDeviceMACAddress      = ""
	dvResourceVirtualEnvironmentVMNetworkDeviceModel           = "virtio"
	dvResourceVirtualEnvironmentVMNetworkDeviceRateLimit       = 0
	dvResourceVirtualEnvironmentVMNetworkDeviceVLANID          = 0
	dvResourceVirtualEnvironmentVMOSType                       = "other"
	dvResourceVirtualEnvironmentVMPoolID                       = ""
	dvResourceVirtualEnvironmentVMStarted                      = true
	dvResourceVirtualEnvironmentVMVMID                         = -1

	mkResourceVirtualEnvironmentVMAgent                        = "agent"
	mkResourceVirtualEnvironmentVMAgentEnabled                 = "enabled"
	mkResourceVirtualEnvironmentVMAgentTrim                    = "trim"
	mkResourceVirtualEnvironmentVMAgentType                    = "type"
	mkResourceVirtualEnvironmentVMCDROM                        = "cdrom"
	mkResourceVirtualEnvironmentVMCDROMEnabled                 = "enabled"
	mkResourceVirtualEnvironmentVMCDROMFileID                  = "file_id"
	mkResourceVirtualEnvironmentVMCloudInit                    = "cloud_init"
	mkResourceVirtualEnvironmentVMCloudInitDNS                 = "dns"
	mkResourceVirtualEnvironmentVMCloudInitDNSDomain           = "domain"
	mkResourceVirtualEnvironmentVMCloudInitDNSServer           = "server"
	mkResourceVirtualEnvironmentVMCloudInitIPConfig            = "ip_config"
	mkResourceVirtualEnvironmentVMCloudInitIPConfigIPv4        = "ipv4"
	mkResourceVirtualEnvironmentVMCloudInitIPConfigIPv4Address = "address"
	mkResourceVirtualEnvironmentVMCloudInitIPConfigIPv4Gateway = "gateway"
	mkResourceVirtualEnvironmentVMCloudInitIPConfigIPv6        = "ipv6"
	mkResourceVirtualEnvironmentVMCloudInitIPConfigIPv6Address = "address"
	mkResourceVirtualEnvironmentVMCloudInitIPConfigIPv6Gateway = "gateway"
	mkResourceVirtualEnvironmentVMCloudInitUserAccount         = "user_account"
	mkResourceVirtualEnvironmentVMCloudInitUserAccountKeys     = "keys"
	mkResourceVirtualEnvironmentVMCloudInitUserAccountPassword = "password"
	mkResourceVirtualEnvironmentVMCloudInitUserAccountUsername = "username"
	mkResourceVirtualEnvironmentVMCloudInitUserDataFileID      = "user_data_file_id"
	mkResourceVirtualEnvironmentVMCPU                          = "cpu"
	mkResourceVirtualEnvironmentVMCPUCores                     = "cores"
	mkResourceVirtualEnvironmentVMCPUHotplugged                = "hotplugged"
	mkResourceVirtualEnvironmentVMCPUSockets                   = "sockets"
	mkResourceVirtualEnvironmentVMDescription                  = "description"
	mkResourceVirtualEnvironmentVMDisk                         = "disk"
	mkResourceVirtualEnvironmentVMDiskDatastoreID              = "datastore_id"
	mkResourceVirtualEnvironmentVMDiskFileFormat               = "file_format"
	mkResourceVirtualEnvironmentVMDiskFileID                   = "file_id"
	mkResourceVirtualEnvironmentVMDiskSize                     = "size"
	mkResourceVirtualEnvironmentVMDiskSpeed                    = "speed"
	mkResourceVirtualEnvironmentVMDiskSpeedRead                = "read"
	mkResourceVirtualEnvironmentVMDiskSpeedReadBurstable       = "read_burstable"
	mkResourceVirtualEnvironmentVMDiskSpeedWrite               = "write"
	mkResourceVirtualEnvironmentVMDiskSpeedWriteBurstable      = "write_burstable"
	mkResourceVirtualEnvironmentVMIPv4Addresses                = "ipv4_addresses"
	mkResourceVirtualEnvironmentVMIPv6Addresses                = "ipv6_addresses"
	mkResourceVirtualEnvironmentVMKeyboardLayout               = "keyboard_layout"
	mkResourceVirtualEnvironmentVMMACAddresses                 = "mac_addresses"
	mkResourceVirtualEnvironmentVMMemory                       = "memory"
	mkResourceVirtualEnvironmentVMMemoryDedicated              = "dedicated"
	mkResourceVirtualEnvironmentVMMemoryFloating               = "floating"
	mkResourceVirtualEnvironmentVMMemoryShared                 = "shared"
	mkResourceVirtualEnvironmentVMName                         = "name"
	mkResourceVirtualEnvironmentVMNetworkDevice                = "network_device"
	mkResourceVirtualEnvironmentVMNetworkDeviceBridge          = "bridge"
	mkResourceVirtualEnvironmentVMNetworkDeviceEnabled         = "enabled"
	mkResourceVirtualEnvironmentVMNetworkDeviceMACAddress      = "mac_address"
	mkResourceVirtualEnvironmentVMNetworkDeviceModel           = "model"
	mkResourceVirtualEnvironmentVMNetworkDeviceRateLimit       = "rate_limit"
	mkResourceVirtualEnvironmentVMNetworkDeviceVLANID          = "vlan_id"
	mkResourceVirtualEnvironmentVMNetworkInterfaceNames        = "network_interface_names"
	mkResourceVirtualEnvironmentVMNodeName                     = "node_name"
	mkResourceVirtualEnvironmentVMOSType                       = "os_type"
	mkResourceVirtualEnvironmentVMPoolID                       = "pool_id"
	mkResourceVirtualEnvironmentVMStarted                      = "started"
	mkResourceVirtualEnvironmentVMVMID                         = "vm_id"
)

func resourceVirtualEnvironmentVM() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			mkResourceVirtualEnvironmentVMAgent: &schema.Schema{
				Type:        schema.TypeList,
				Description: "The QEMU agent configuration",
				Optional:    true,
				DefaultFunc: func() (interface{}, error) {
					defaultList := make([]interface{}, 1)
					defaultMap := map[string]interface{}{}

					defaultMap[mkResourceVirtualEnvironmentVMAgentEnabled] = dvResourceVirtualEnvironmentVMAgentEnabled
					defaultMap[mkResourceVirtualEnvironmentVMAgentTrim] = dvResourceVirtualEnvironmentVMAgentTrim
					defaultMap[mkResourceVirtualEnvironmentVMAgentType] = dvResourceVirtualEnvironmentVMAgentType

					defaultList[0] = defaultMap

					return defaultList, nil
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mkResourceVirtualEnvironmentVMAgentEnabled: {
							Type:        schema.TypeBool,
							Description: "Whether to enable the QEMU agent",
							Optional:    true,
							Default:     dvResourceVirtualEnvironmentVMAgentEnabled,
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
			mkResourceVirtualEnvironmentVMCDROM: &schema.Schema{
				Type:        schema.TypeList,
				Description: "The CDROM drive",
				Optional:    true,
				DefaultFunc: func() (interface{}, error) {
					defaultList := make([]interface{}, 1)
					defaultMap := map[string]interface{}{}

					defaultMap[mkResourceVirtualEnvironmentVMCDROMEnabled] = dvResourceVirtualEnvironmentVMCDROMEnabled
					defaultMap[mkResourceVirtualEnvironmentVMCDROMFileID] = dvResourceVirtualEnvironmentVMCDROMFileID

					defaultList[0] = defaultMap

					return defaultList, nil
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
			mkResourceVirtualEnvironmentVMCloudInit: &schema.Schema{
				Type:        schema.TypeList,
				Description: "The cloud-init configuration",
				Optional:    true,
				DefaultFunc: func() (interface{}, error) {
					return []interface{}{}, nil
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mkResourceVirtualEnvironmentVMCloudInitDNS: {
							Type:        schema.TypeList,
							Description: "The DNS configuration",
							Optional:    true,
							DefaultFunc: func() (interface{}, error) {
								return []interface{}{}, nil
							},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									mkResourceVirtualEnvironmentVMCloudInitDNSDomain: {
										Type:        schema.TypeString,
										Description: "The DNS search domain",
										Optional:    true,
										Default:     dvResourceVirtualEnvironmentVMCloudInitDNSDomain,
									},
									mkResourceVirtualEnvironmentVMCloudInitDNSServer: {
										Type:        schema.TypeString,
										Description: "The DNS server",
										Optional:    true,
										Default:     dvResourceVirtualEnvironmentVMCloudInitDNSServer,
									},
								},
							},
							MaxItems: 1,
							MinItems: 0,
						},
						mkResourceVirtualEnvironmentVMCloudInitIPConfig: {
							Type:        schema.TypeList,
							Description: "The IP configuration",
							Optional:    true,
							DefaultFunc: func() (interface{}, error) {
								return []interface{}{}, nil
							},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									mkResourceVirtualEnvironmentVMCloudInitIPConfigIPv4: {
										Type:        schema.TypeList,
										Description: "The IPv4 configuration",
										Optional:    true,
										DefaultFunc: func() (interface{}, error) {
											return []interface{}{}, nil
										},
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												mkResourceVirtualEnvironmentVMCloudInitIPConfigIPv4Address: {
													Type:        schema.TypeString,
													Description: "The IPv4 address",
													Optional:    true,
													Default:     "",
												},
												mkResourceVirtualEnvironmentVMCloudInitIPConfigIPv4Gateway: {
													Type:        schema.TypeString,
													Description: "The IPv4 gateway",
													Optional:    true,
													Default:     "",
												},
											},
										},
										MaxItems: 1,
										MinItems: 0,
									},
									mkResourceVirtualEnvironmentVMCloudInitIPConfigIPv6: {
										Type:        schema.TypeList,
										Description: "The IPv6 configuration",
										Optional:    true,
										DefaultFunc: func() (interface{}, error) {
											return []interface{}{}, nil
										},
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												mkResourceVirtualEnvironmentVMCloudInitIPConfigIPv6Address: {
													Type:        schema.TypeString,
													Description: "The IPv6 address",
													Optional:    true,
													Default:     "",
												},
												mkResourceVirtualEnvironmentVMCloudInitIPConfigIPv6Gateway: {
													Type:        schema.TypeString,
													Description: "The IPv6 gateway",
													Optional:    true,
													Default:     "",
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
						mkResourceVirtualEnvironmentVMCloudInitUserAccount: {
							Type:        schema.TypeList,
							Description: "The user account configuration",
							Required:    true,
							ForceNew:    true,
							DefaultFunc: func() (interface{}, error) {
								return []interface{}{}, nil
							},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									mkResourceVirtualEnvironmentVMCloudInitUserAccountKeys: {
										Type:        schema.TypeList,
										Description: "The SSH keys",
										Required:    true,
										ForceNew:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									mkResourceVirtualEnvironmentVMCloudInitUserAccountPassword: {
										Type:        schema.TypeString,
										Description: "The SSH password",
										Optional:    true,
										ForceNew:    true,
										Default:     dvResourceVirtualEnvironmentVMCloudInitUserAccountPassword,
										DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
											return strings.ReplaceAll(old, "*", "") == ""
										},
									},
									mkResourceVirtualEnvironmentVMCloudInitUserAccountUsername: {
										Type:        schema.TypeString,
										Description: "The SSH username",
										Required:    true,
										ForceNew:    true,
									},
								},
							},
							MaxItems: 1,
							MinItems: 0,
						},
						mkResourceVirtualEnvironmentVMCloudInitUserDataFileID: {
							Type:         schema.TypeString,
							Description:  "The ID of a file containing custom user data",
							Optional:     true,
							ForceNew:     true,
							Default:      dvResourceVirtualEnvironmentVMCloudInitUserDataFileID,
							ValidateFunc: getFileIDValidator(),
						},
					},
				},
				MaxItems: 1,
				MinItems: 0,
			},
			mkResourceVirtualEnvironmentVMCPU: &schema.Schema{
				Type:        schema.TypeList,
				Description: "The CPU allocation",
				Optional:    true,
				DefaultFunc: func() (interface{}, error) {
					defaultList := make([]interface{}, 1)
					defaultMap := map[string]interface{}{}

					defaultMap[mkResourceVirtualEnvironmentVMCPUCores] = dvResourceVirtualEnvironmentVMCPUCores
					defaultMap[mkResourceVirtualEnvironmentVMCPUHotplugged] = dvResourceVirtualEnvironmentVMCPUHotplugged
					defaultMap[mkResourceVirtualEnvironmentVMCPUSockets] = dvResourceVirtualEnvironmentVMCPUSockets

					defaultList[0] = defaultMap

					return defaultList, nil
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mkResourceVirtualEnvironmentVMCPUCores: {
							Type:         schema.TypeInt,
							Description:  "The number of CPU cores",
							Optional:     true,
							Default:      dvResourceVirtualEnvironmentVMCPUCores,
							ValidateFunc: validation.IntBetween(1, 2304),
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
			mkResourceVirtualEnvironmentVMDisk: &schema.Schema{
				Type:        schema.TypeList,
				Description: "The disk devices",
				Optional:    true,
				ForceNew:    true,
				DefaultFunc: func() (interface{}, error) {
					defaultList := make([]interface{}, 1)
					defaultMap := map[string]interface{}{}

					defaultMap[mkResourceVirtualEnvironmentVMDiskDatastoreID] = dvResourceVirtualEnvironmentVMDiskDatastoreID
					defaultMap[mkResourceVirtualEnvironmentVMDiskFileFormat] = dvResourceVirtualEnvironmentVMDiskFileFormat
					defaultMap[mkResourceVirtualEnvironmentVMDiskFileID] = dvResourceVirtualEnvironmentVMDiskFileID
					defaultMap[mkResourceVirtualEnvironmentVMDiskSize] = dvResourceVirtualEnvironmentVMDiskSize

					defaultList[0] = defaultMap

					return defaultList, nil
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
								defaultList := make([]interface{}, 1)
								defaultMap := map[string]interface{}{}

								defaultMap[mkResourceVirtualEnvironmentVMDiskSpeedRead] = dvResourceVirtualEnvironmentVMDiskSpeedRead
								defaultMap[mkResourceVirtualEnvironmentVMDiskSpeedReadBurstable] = dvResourceVirtualEnvironmentVMDiskSpeedReadBurstable
								defaultMap[mkResourceVirtualEnvironmentVMDiskSpeedWrite] = dvResourceVirtualEnvironmentVMDiskSpeedWrite
								defaultMap[mkResourceVirtualEnvironmentVMDiskSpeedWriteBurstable] = dvResourceVirtualEnvironmentVMDiskSpeedWriteBurstable

								defaultList[0] = defaultMap

								return defaultList, nil
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
			mkResourceVirtualEnvironmentVMMemory: &schema.Schema{
				Type:        schema.TypeList,
				Description: "The memory allocation",
				Optional:    true,
				DefaultFunc: func() (interface{}, error) {
					defaultList := make([]interface{}, 1)
					defaultMap := map[string]interface{}{}

					defaultMap[mkResourceVirtualEnvironmentVMMemoryDedicated] = dvResourceVirtualEnvironmentVMMemoryDedicated
					defaultMap[mkResourceVirtualEnvironmentVMMemoryFloating] = dvResourceVirtualEnvironmentVMMemoryFloating
					defaultMap[mkResourceVirtualEnvironmentVMMemoryShared] = dvResourceVirtualEnvironmentVMMemoryShared

					defaultList[0] = defaultMap

					return defaultList, nil
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
			mkResourceVirtualEnvironmentVMNetworkDevice: &schema.Schema{
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
				MaxItems: 8,
				MinItems: 0,
			},
			mkResourceVirtualEnvironmentVMNetworkInterfaceNames: {
				Type:        schema.TypeList,
				Description: "The network interface names published by the QEMU agent",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			mkResourceVirtualEnvironmentVMNodeName: &schema.Schema{
				Type:        schema.TypeString,
				Description: "The node name",
				Required:    true,
				ForceNew:    true,
			},
			mkResourceVirtualEnvironmentVMOSType: {
				Type:         schema.TypeString,
				Description:  "The OS type",
				Optional:     true,
				Default:      dvResourceVirtualEnvironmentVMOSType,
				ValidateFunc: getOSTypeValidator(),
			},
			mkResourceVirtualEnvironmentVMPoolID: {
				Type:        schema.TypeString,
				Description: "The ID of the pool to assign the virtual machine to",
				Optional:    true,
				ForceNew:    true,
				Default:     dvResourceVirtualEnvironmentVMPoolID,
			},
			mkResourceVirtualEnvironmentVMStarted: {
				Type:        schema.TypeBool,
				Description: "Whether to start the virtual machine",
				Optional:    true,
				Default:     dvResourceVirtualEnvironmentVMStarted,
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
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()

	if err != nil {
		return err
	}

	resource := resourceVirtualEnvironmentVM()

	agentBlock, err := getSchemaBlock(resource, d, m, []string{mkResourceVirtualEnvironmentVMAgent}, 0, true)

	if err != nil {
		return err
	}

	agentEnabled := proxmox.CustomBool(agentBlock[mkResourceVirtualEnvironmentVMAgentEnabled].(bool))
	agentTrim := proxmox.CustomBool(agentBlock[mkResourceVirtualEnvironmentVMAgentTrim].(bool))
	agentType := agentBlock[mkResourceVirtualEnvironmentVMAgentType].(string)

	cdromBlock, err := getSchemaBlock(resource, d, m, []string{mkResourceVirtualEnvironmentVMCDROM}, 0, true)

	if err != nil {
		return err
	}

	cdromEnabled := cdromBlock[mkResourceVirtualEnvironmentVMCDROMEnabled].(bool)
	cdromFileID := cdromBlock[mkResourceVirtualEnvironmentVMCDROMFileID].(string)

	if cdromFileID == "" {
		cdromFileID = "cdrom"
	}

	cloudInitConfig, err := resourceVirtualEnvironmentVMGetCloudConfig(d, m)

	if err != nil {
		return err
	}

	if cloudInitConfig != nil {
		cdromEnabled = true
		cdromFileID = "local-lvm:cloudinit"
	}

	cpuBlock, err := getSchemaBlock(resource, d, m, []string{mkResourceVirtualEnvironmentVMCPU}, 0, true)

	if err != nil {
		return err
	}

	cpuCores := cpuBlock[mkResourceVirtualEnvironmentVMCPUCores].(int)
	cpuHotplugged := cpuBlock[mkResourceVirtualEnvironmentVMCPUHotplugged].(int)
	cpuSockets := cpuBlock[mkResourceVirtualEnvironmentVMCPUSockets].(int)

	description := d.Get(mkResourceVirtualEnvironmentVMDescription).(string)
	diskDeviceObjects, err := resourceVirtualEnvironmentVMGetDiskDeviceObjects(d, m)

	if err != nil {
		return err
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
	osType := d.Get(mkResourceVirtualEnvironmentVMOSType).(string)
	poolID := d.Get(mkResourceVirtualEnvironmentVMPoolID).(string)
	started := proxmox.CustomBool(d.Get(mkResourceVirtualEnvironmentVMStarted).(bool))
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
	tabletDeviceEnabled := proxmox.CustomBool(true)

	body := &proxmox.VirtualEnvironmentVMCreateRequestBody{
		Agent: &proxmox.CustomAgent{
			Enabled:         &agentEnabled,
			TrimClonedDisks: &agentTrim,
			Type:            &agentType,
		},
		BootDisk:            &bootDisk,
		BootOrder:           &bootOrder,
		CloudInitConfig:     cloudInitConfig,
		CPUCores:            &cpuCores,
		CPUSockets:          &cpuSockets,
		DedicatedMemory:     &memoryDedicated,
		FloatingMemory:      &memoryFloating,
		IDEDevices:          ideDevices,
		KeyboardLayout:      &keyboardLayout,
		NetworkDevices:      networkDeviceObjects,
		OSType:              &osType,
		PoolID:              &poolID,
		SCSIDevices:         diskDeviceObjects,
		SCSIHardware:        &scsiHardware,
		SerialDevices:       []string{"socket"},
		SharedMemory:        memorySharedObject,
		StartOnBoot:         &started,
		TabletDeviceEnabled: &tabletDeviceEnabled,
		VMID:                &vmID,
	}

	if cpuHotplugged > 0 {
		body.VirtualCPUCount = &cpuHotplugged
	}

	if description != "" {
		body.Description = &description
	}

	if name != "" {
		body.Name = &name
	}

	err = veClient.CreateVM(nodeName, body)

	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(vmID))

	return resourceVirtualEnvironmentVMCreateImportedDisks(d, m)
}

func resourceVirtualEnvironmentVMCreateImportedDisks(d *schema.ResourceData, m interface{}) error {
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
			fmt.Sprintf(`cp "$(grep -Pzo ': %s\s+path\s+[^\s]+' /etc/pve/storage.cfg | grep -Pzo '/[^\s]*' | tr -d '\000')%s" %s`, fileIDParts[0], filePath, filePathTmp),
			fmt.Sprintf(`qemu-img resize %s %dG`, filePathTmp, size),
			fmt.Sprintf(`qm importdisk %d %s %s -format qcow2`, vmID, filePathTmp, datastoreID),
			fmt.Sprintf(`qm set %d -scsi%d %s:vm-%d-disk-%d%s`, vmID, i, datastoreID, vmID, diskCount+importedDiskCount, diskOptions),
			fmt.Sprintf(`rm -f %s`, filePathTmp),
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

	if !started {
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

	err = veClient.WaitForState(nodeName, vmID, "running", 120, 5)

	if err != nil {
		return err
	}

	return resourceVirtualEnvironmentVMRead(d, m)
}

func resourceVirtualEnvironmentVMGetCloudConfig(d *schema.ResourceData, m interface{}) (*proxmox.CustomCloudInitConfig, error) {
	var cloudInitConfig *proxmox.CustomCloudInitConfig

	cloudInit := d.Get(mkResourceVirtualEnvironmentVMCloudInit).([]interface{})

	if len(cloudInit) > 0 {
		cloudInitBlock := cloudInit[0].(map[string]interface{})
		cloudInitConfig = &proxmox.CustomCloudInitConfig{}
		cloudInitDNS := cloudInitBlock[mkResourceVirtualEnvironmentVMCloudInitDNS].([]interface{})

		if len(cloudInitDNS) > 0 {
			cloudInitDNSBlock := cloudInitDNS[0].(map[string]interface{})
			domain := cloudInitDNSBlock[mkResourceVirtualEnvironmentVMCloudInitDNSDomain].(string)

			if domain != "" {
				cloudInitConfig.SearchDomain = &domain
			}

			server := cloudInitDNSBlock[mkResourceVirtualEnvironmentVMCloudInitDNSServer].(string)

			if server != "" {
				cloudInitConfig.Nameserver = &server
			}
		}

		cloudInitIPConfig := cloudInitBlock[mkResourceVirtualEnvironmentVMCloudInitIPConfig].([]interface{})
		cloudInitConfig.IPConfig = make([]proxmox.CustomCloudInitIPConfig, len(cloudInitIPConfig))

		for i, c := range cloudInitIPConfig {
			configBlock := c.(map[string]interface{})
			ipv4 := configBlock[mkResourceVirtualEnvironmentVMCloudInitIPConfigIPv4].([]interface{})

			if len(ipv4) > 0 {
				ipv4Block := ipv4[0].(map[string]interface{})
				ipv4Address := ipv4Block[mkResourceVirtualEnvironmentVMCloudInitIPConfigIPv4Address].(string)

				if ipv4Address != "" {
					cloudInitConfig.IPConfig[i].IPv4 = &ipv4Address
				}

				ipv4Gateway := ipv4Block[mkResourceVirtualEnvironmentVMCloudInitIPConfigIPv4Gateway].(string)

				if ipv4Gateway != "" {
					cloudInitConfig.IPConfig[i].GatewayIPv4 = &ipv4Gateway
				}
			}

			ipv6 := configBlock[mkResourceVirtualEnvironmentVMCloudInitIPConfigIPv6].([]interface{})

			if len(ipv6) > 0 {
				ipv6Block := ipv6[0].(map[string]interface{})
				ipv6Address := ipv6Block[mkResourceVirtualEnvironmentVMCloudInitIPConfigIPv6Address].(string)

				if ipv6Address != "" {
					cloudInitConfig.IPConfig[i].IPv6 = &ipv6Address
				}

				ipv6Gateway := ipv6Block[mkResourceVirtualEnvironmentVMCloudInitIPConfigIPv6Gateway].(string)

				if ipv6Gateway != "" {
					cloudInitConfig.IPConfig[i].GatewayIPv6 = &ipv6Gateway
				}
			}
		}

		cloudInitUserAccount := cloudInitBlock[mkResourceVirtualEnvironmentVMCloudInitUserAccount].([]interface{})

		if len(cloudInitUserAccount) > 0 {
			cloudInitUserAccountBlock := cloudInitUserAccount[0].(map[string]interface{})
			keys := cloudInitUserAccountBlock[mkResourceVirtualEnvironmentVMCloudInitUserAccountKeys].([]interface{})

			if len(keys) > 0 {
				sshKeys := make(proxmox.CustomCloudInitSSHKeys, len(keys))

				for i, k := range keys {
					sshKeys[i] = k.(string)
				}

				cloudInitConfig.SSHKeys = &sshKeys
			}

			password := cloudInitUserAccountBlock[mkResourceVirtualEnvironmentVMCloudInitUserAccountPassword].(string)

			if password != "" {
				cloudInitConfig.Password = &password
			}

			username := cloudInitUserAccountBlock[mkResourceVirtualEnvironmentVMCloudInitUserAccountUsername].(string)

			cloudInitConfig.Username = &username
		}

		cloudInitUserDataFileID := cloudInitBlock[mkResourceVirtualEnvironmentVMCloudInitUserDataFileID].(string)

		if cloudInitUserDataFileID != "" {
			cloudInitConfig.Files = &proxmox.CustomCloudInitFiles{
				UserVolume: &cloudInitUserDataFileID,
			}
		}
	}

	return cloudInitConfig, nil
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
		if strings.Contains(err.Error(), "HTTP 404") {
			d.SetId("")

			return nil
		}

		return err
	}

	// Compare the agent configuration to the one stored in the state.
	if vmConfig.Agent != nil {
		agent := map[string]interface{}{}

		if vmConfig.Agent.Enabled != nil {
			agent[mkResourceVirtualEnvironmentVMAgentEnabled] = bool(*vmConfig.Agent.Enabled)
		} else {
			agent[mkResourceVirtualEnvironmentVMAgentEnabled] = dvResourceVirtualEnvironmentVMAgentEnabled
		}

		if vmConfig.Agent.TrimClonedDisks != nil {
			agent[mkResourceVirtualEnvironmentVMAgentTrim] = bool(*vmConfig.Agent.TrimClonedDisks)
		} else {
			agent[mkResourceVirtualEnvironmentVMAgentTrim] = dvResourceVirtualEnvironmentVMAgentTrim
		}

		if vmConfig.Agent.Type != nil {
			agent[mkResourceVirtualEnvironmentVMAgentType] = *vmConfig.Agent.Type
		} else {
			agent[mkResourceVirtualEnvironmentVMAgentType] = dvResourceVirtualEnvironmentVMAgentType
		}

		currentAgent := d.Get(mkResourceVirtualEnvironmentVMAgent).([]interface{})

		if len(currentAgent) > 0 ||
			agent[mkResourceVirtualEnvironmentVMAgentEnabled] != dvResourceVirtualEnvironmentVMAgentEnabled ||
			agent[mkResourceVirtualEnvironmentVMAgentTrim] != dvResourceVirtualEnvironmentVMAgentTrim ||
			agent[mkResourceVirtualEnvironmentVMAgentType] != dvResourceVirtualEnvironmentVMAgentType {
			d.Set(mkResourceVirtualEnvironmentVMAgent, []interface{}{agent})
		}
	} else {
		d.Set(mkResourceVirtualEnvironmentVMAgent, make([]interface{}, 0))
	}

	// Compare the IDE devices to the CDROM and cloud-init configurations stored in the state.
	if vmConfig.IDEDevice2 != nil {
		if *vmConfig.IDEDevice2.Media == "cdrom" {
			if strings.Contains(vmConfig.IDEDevice2.FileVolume, fmt.Sprintf("vm-%d-cloudinit", vmID)) {
				d.Set(mkResourceVirtualEnvironmentVMCDROM, make([]interface{}, 0))
			} else {
				d.Set(mkResourceVirtualEnvironmentVMCloudInit, make([]interface{}, 0))

				cdrom := make([]interface{}, 1)
				cdromBlock := map[string]interface{}{}

				cdromBlock[mkResourceVirtualEnvironmentVMCDROMEnabled] = true
				cdromBlock[mkResourceVirtualEnvironmentVMCDROMFileID] = vmConfig.IDEDevice2.FileVolume

				cdrom[0] = cdromBlock

				d.Set(mkResourceVirtualEnvironmentVMCDROM, cdrom)
			}
		} else {
			d.Set(mkResourceVirtualEnvironmentVMCDROM, make([]interface{}, 0))
			d.Set(mkResourceVirtualEnvironmentVMCloudInit, make([]interface{}, 0))
		}
	} else {
		d.Set(mkResourceVirtualEnvironmentVMCDROM, make([]interface{}, 0))
		d.Set(mkResourceVirtualEnvironmentVMCloudInit, make([]interface{}, 0))
	}

	// Compare the cloud-init configuration to the one stored in the state.
	cloudInit := map[string]interface{}{}

	if vmConfig.CloudInitDNSDomain != nil || vmConfig.CloudInitDNSServer != nil {
		cloudInitDNS := map[string]interface{}{}

		if vmConfig.CloudInitDNSDomain != nil {
			cloudInitDNS[mkResourceVirtualEnvironmentVMCloudInitDNSDomain] = *vmConfig.CloudInitDNSDomain
		} else {
			cloudInitDNS[mkResourceVirtualEnvironmentVMCloudInitDNSDomain] = dvResourceVirtualEnvironmentVMCloudInitDNSDomain
		}

		if vmConfig.CloudInitDNSServer != nil {
			cloudInitDNS[mkResourceVirtualEnvironmentVMCloudInitDNSServer] = *vmConfig.CloudInitDNSServer
		} else {
			cloudInitDNS[mkResourceVirtualEnvironmentVMCloudInitDNSServer] = dvResourceVirtualEnvironmentVMCloudInitDNSServer
		}

		cloudInit[mkResourceVirtualEnvironmentVMCloudInitDNS] = []interface{}{cloudInitDNS}
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
					ipv4[mkResourceVirtualEnvironmentVMCloudInitIPConfigIPv4Address] = *ipConfig.IPv4
				} else {
					ipv4[mkResourceVirtualEnvironmentVMCloudInitIPConfigIPv4Address] = ""
				}

				if ipConfig.GatewayIPv4 != nil {
					ipv4[mkResourceVirtualEnvironmentVMCloudInitIPConfigIPv4Gateway] = *ipConfig.GatewayIPv4
				} else {
					ipv4[mkResourceVirtualEnvironmentVMCloudInitIPConfigIPv4Gateway] = ""
				}

				ipConfigItem[mkResourceVirtualEnvironmentVMCloudInitIPConfigIPv4] = []interface{}{ipv4}
			} else {
				ipConfigItem[mkResourceVirtualEnvironmentVMCloudInitIPConfigIPv4] = []interface{}{}
			}

			if ipConfig.GatewayIPv6 != nil || ipConfig.IPv6 != nil {
				ipv6 := map[string]interface{}{}

				if ipConfig.IPv4 != nil {
					ipv6[mkResourceVirtualEnvironmentVMCloudInitIPConfigIPv4Address] = *ipConfig.IPv6
				} else {
					ipv6[mkResourceVirtualEnvironmentVMCloudInitIPConfigIPv4Address] = ""
				}

				if ipConfig.GatewayIPv4 != nil {
					ipv6[mkResourceVirtualEnvironmentVMCloudInitIPConfigIPv4Gateway] = *ipConfig.GatewayIPv6
				} else {
					ipv6[mkResourceVirtualEnvironmentVMCloudInitIPConfigIPv4Gateway] = ""
				}

				ipConfigItem[mkResourceVirtualEnvironmentVMCloudInitIPConfigIPv6] = []interface{}{ipv6}
			} else {
				ipConfigItem[mkResourceVirtualEnvironmentVMCloudInitIPConfigIPv6] = []interface{}{}
			}
		} else {
			ipConfigItem[mkResourceVirtualEnvironmentVMCloudInitIPConfigIPv4] = []interface{}{}
			ipConfigItem[mkResourceVirtualEnvironmentVMCloudInitIPConfigIPv6] = []interface{}{}
		}

		ipConfigList[ipConfigIndex] = ipConfigItem
	}

	cloudInit[mkResourceVirtualEnvironmentVMCloudInitIPConfig] = ipConfigList[:ipConfigLast+1]

	if vmConfig.CloudInitPassword != nil || vmConfig.CloudInitSSHKeys != nil || vmConfig.CloudInitUsername != nil {
		cloudInitUserAccount := map[string]interface{}{}

		if vmConfig.CloudInitSSHKeys != nil {
			cloudInitUserAccount[mkResourceVirtualEnvironmentVMCloudInitUserAccountKeys] = []string(*vmConfig.CloudInitSSHKeys)
		} else {
			cloudInitUserAccount[mkResourceVirtualEnvironmentVMCloudInitUserAccountKeys] = []string{}
		}

		if vmConfig.CloudInitPassword != nil {
			cloudInitUserAccount[mkResourceVirtualEnvironmentVMCloudInitUserAccountPassword] = *vmConfig.CloudInitPassword
		} else {
			cloudInitUserAccount[mkResourceVirtualEnvironmentVMCloudInitUserAccountPassword] = dvResourceVirtualEnvironmentVMCloudInitUserAccountPassword
		}

		if vmConfig.CloudInitUsername != nil {
			cloudInitUserAccount[mkResourceVirtualEnvironmentVMCloudInitUserAccountUsername] = *vmConfig.CloudInitUsername
		} else {
			cloudInitUserAccount[mkResourceVirtualEnvironmentVMCloudInitUserAccountUsername] = ""
		}

		cloudInit[mkResourceVirtualEnvironmentVMCloudInitUserAccount] = []interface{}{cloudInitUserAccount}
	}

	if vmConfig.CloudInitFiles != nil {
		if vmConfig.CloudInitFiles.UserVolume != nil {
			cloudInit[mkResourceVirtualEnvironmentVMCloudInitUserDataFileID] = *vmConfig.CloudInitFiles.UserVolume
		} else {
			cloudInit[mkResourceVirtualEnvironmentVMCloudInitUserDataFileID] = dvResourceVirtualEnvironmentVMCloudInitUserDataFileID
		}
	} else {
		cloudInit[mkResourceVirtualEnvironmentVMCloudInitUserDataFileID] = dvResourceVirtualEnvironmentVMCloudInitUserDataFileID
	}

	if len(cloudInit) > 0 {
		d.Set(mkResourceVirtualEnvironmentVMCloudInit, []interface{}{cloudInit})
	} else {
		d.Set(mkResourceVirtualEnvironmentVMCloudInit, []interface{}{})
	}

	// Compare the CPU configuration to the one stored in the state.
	cpu := map[string]interface{}{}

	if vmConfig.CPUCores != nil {
		cpu[mkResourceVirtualEnvironmentVMCPUCores] = *vmConfig.CPUCores
	} else {
		cpu[mkResourceVirtualEnvironmentVMCPUCores] = dvResourceVirtualEnvironmentVMCPUCores
	}

	if vmConfig.VirtualCPUCount != nil {
		cpu[mkResourceVirtualEnvironmentVMCPUHotplugged] = *vmConfig.VirtualCPUCount
	} else {
		cpu[mkResourceVirtualEnvironmentVMCPUHotplugged] = dvResourceVirtualEnvironmentVMCPUHotplugged
	}

	if vmConfig.CPUSockets != nil {
		cpu[mkResourceVirtualEnvironmentVMCPUSockets] = *vmConfig.CPUSockets
	} else {
		cpu[mkResourceVirtualEnvironmentVMCPUSockets] = dvResourceVirtualEnvironmentVMCPUSockets
	}

	currentCPU := d.Get(mkResourceVirtualEnvironmentVMCPU).([]interface{})

	if len(currentCPU) > 0 ||
		cpu[mkResourceVirtualEnvironmentVMCPUCores] != dvResourceVirtualEnvironmentVMCPUCores ||
		cpu[mkResourceVirtualEnvironmentVMCPUHotplugged] != dvResourceVirtualEnvironmentVMCPUHotplugged ||
		cpu[mkResourceVirtualEnvironmentVMCPUSockets] != dvResourceVirtualEnvironmentVMCPUSockets {
		d.Set(mkResourceVirtualEnvironmentVMCPU, []interface{}{cpu})
	}

	// Compare the description and keyboard layout to the values stored in the state.
	if vmConfig.Description != nil {
		d.Set(mkResourceVirtualEnvironmentVMDescription, *vmConfig.Description)
	} else {
		d.Set(mkResourceVirtualEnvironmentVMDescription, "")
	}

	if vmConfig.KeyboardLayout != nil {
		d.Set(mkResourceVirtualEnvironmentVMKeyboardLayout, *vmConfig.KeyboardLayout)
	} else {
		d.Set(mkResourceVirtualEnvironmentVMKeyboardLayout, "")
	}

	// Compare the disks to those stored in the state.
	currentDiskList := d.Get(mkResourceVirtualEnvironmentVMDisk).([]interface{})

	diskList := make([]interface{}, 0)
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

		if len(currentDiskList) > di {
			currentDisk := currentDiskList[di].(map[string]interface{})

			disk[mkResourceVirtualEnvironmentVMDiskFileFormat] = currentDisk[mkResourceVirtualEnvironmentVMDiskFileFormat]
			disk[mkResourceVirtualEnvironmentVMDiskFileID] = currentDisk[mkResourceVirtualEnvironmentVMDiskFileID]
		}

		diskSize := 0

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

	if len(currentDiskList) > 0 || len(diskList) > 0 {
		d.Set(mkResourceVirtualEnvironmentVMDisk, diskList)
	}

	// Compare the memory configuration to the one stored in the state.
	memory := map[string]interface{}{}

	if vmConfig.DedicatedMemory != nil {
		memory[mkResourceVirtualEnvironmentVMMemoryDedicated] = *vmConfig.DedicatedMemory
	} else {
		memory[mkResourceVirtualEnvironmentVMMemoryDedicated] = dvResourceVirtualEnvironmentVMMemoryDedicated
	}

	if vmConfig.FloatingMemory != nil {
		memory[mkResourceVirtualEnvironmentVMMemoryFloating] = *vmConfig.FloatingMemory
	} else {
		memory[mkResourceVirtualEnvironmentVMMemoryFloating] = dvResourceVirtualEnvironmentVMMemoryFloating
	}

	if vmConfig.SharedMemory != nil {
		memory[mkResourceVirtualEnvironmentVMMemoryShared] = vmConfig.SharedMemory.Size
	} else {
		memory[mkResourceVirtualEnvironmentVMMemoryShared] = dvResourceVirtualEnvironmentVMMemoryShared
	}

	currentMemory := d.Get(mkResourceVirtualEnvironmentVMMemory).([]interface{})

	if len(currentMemory) > 0 ||
		memory[mkResourceVirtualEnvironmentVMMemoryDedicated] != dvResourceVirtualEnvironmentVMMemoryDedicated ||
		memory[mkResourceVirtualEnvironmentVMMemoryFloating] != dvResourceVirtualEnvironmentVMMemoryFloating ||
		memory[mkResourceVirtualEnvironmentVMMemoryShared] != dvResourceVirtualEnvironmentVMMemoryShared {
		d.Set(mkResourceVirtualEnvironmentVMMemory, []interface{}{memory})
	}

	// Compare the name to the value stored in the state.
	if vmConfig.Name != nil {
		d.Set(mkResourceVirtualEnvironmentVMName, *vmConfig.Name)
	} else {
		d.Set(mkResourceVirtualEnvironmentVMName, "")
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

	d.Set(mkResourceVirtualEnvironmentVMMACAddresses, macAddresses[0:len(currentNetworkDeviceList)])

	if len(currentNetworkDeviceList) > 0 || networkDeviceLast > -1 {
		d.Set(mkResourceVirtualEnvironmentVMNetworkDevice, networkDeviceList[:networkDeviceLast+1])
	}

	// Compare the OS type and pool ID to the values stored in the state.
	if vmConfig.OSType != nil {
		d.Set(mkResourceVirtualEnvironmentVMOSType, *vmConfig.OSType)
	} else {
		d.Set(mkResourceVirtualEnvironmentVMOSType, "")
	}

	if vmConfig.PoolID != nil {
		d.Set(mkResourceVirtualEnvironmentVMPoolID, *vmConfig.PoolID)
	}

	// Determine the state of the virtual machine in order to update the "started" argument.
	status, err := veClient.GetVMStatus(nodeName, vmID)

	if err != nil {
		return err
	}

	d.Set(mkResourceVirtualEnvironmentVMStarted, status.Status == "running")

	// Populate the attributes that rely on the QEMU agent.
	ipv4Addresses := []interface{}{}
	ipv6Addresses := []interface{}{}
	networkInterfaceNames := []interface{}{}

	if vmConfig.Agent != nil && vmConfig.Agent.Enabled != nil && *vmConfig.Agent.Enabled {
		networkInterfaces, err := veClient.WaitForNetworkInterfacesFromAgent(nodeName, vmID, 1800, 5)

		if err == nil && networkInterfaces.Result != nil {
			ipv4Addresses = make([]interface{}, len(*networkInterfaces.Result))
			ipv6Addresses = make([]interface{}, len(*networkInterfaces.Result))
			macAddresses = make([]interface{}, len(*networkInterfaces.Result))
			networkInterfaceNames = make([]interface{}, len(*networkInterfaces.Result))

			for ri, rv := range *networkInterfaces.Result {
				rvIPv4Addresses := []interface{}{}
				rvIPv6Addresses := []interface{}{}

				for _, ip := range *rv.IPAddresses {
					switch ip.Type {
					case "ipv4":
						rvIPv4Addresses = append(rvIPv4Addresses, ip.Address)
					case "ipv6":
						rvIPv6Addresses = append(rvIPv6Addresses, ip.Address)
					}
				}

				ipv4Addresses[ri] = rvIPv4Addresses
				ipv6Addresses[ri] = rvIPv6Addresses
				macAddresses[ri] = strings.ToUpper(rv.MACAddress)
				networkInterfaceNames[ri] = rv.Name
			}
		}
	}

	d.Set(mkResourceVirtualEnvironmentVMIPv4Addresses, ipv4Addresses)
	d.Set(mkResourceVirtualEnvironmentVMIPv6Addresses, ipv6Addresses)
	d.Set(mkResourceVirtualEnvironmentVMMACAddresses, macAddresses)
	d.Set(mkResourceVirtualEnvironmentVMNetworkInterfaceNames, networkInterfaceNames)

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

	body := &proxmox.VirtualEnvironmentVMUpdateRequestBody{}
	resource := resourceVirtualEnvironmentVM()

	// Retrieve the entire configuration as we need to process certain values.
	vmConfig, err := veClient.GetVM(nodeName, vmID)

	if err != nil {
		return err
	}

	// Prepare the new primitive configuration values.
	description := d.Get(mkResourceVirtualEnvironmentVMDescription).(string)
	keyboardLayout := d.Get(mkResourceVirtualEnvironmentVMKeyboardLayout).(string)
	name := d.Get(mkResourceVirtualEnvironmentVMName).(string)
	osType := d.Get(mkResourceVirtualEnvironmentVMOSType).(string)

	if description != "" {
		body.Description = &description
	}

	body.KeyboardLayout = &keyboardLayout

	if name != "" {
		body.Name = &name
	}

	body.OSType = &osType

	// Prepare the new agent configuration.
	if d.HasChange(mkResourceVirtualEnvironmentVMAgent) {
		agentBlock, err := getSchemaBlock(resource, d, m, []string{mkResourceVirtualEnvironmentVMAgent}, 0, true)

		if err != nil {
			return err
		}

		agentEnabled := proxmox.CustomBool(agentBlock[mkResourceVirtualEnvironmentVMAgentEnabled].(bool))
		agentTrim := proxmox.CustomBool(agentBlock[mkResourceVirtualEnvironmentVMAgentTrim].(bool))
		agentType := agentBlock[mkResourceVirtualEnvironmentVMAgentType].(string)

		body.Agent = &proxmox.CustomAgent{
			Enabled:         &agentEnabled,
			TrimClonedDisks: &agentTrim,
			Type:            &agentType,
		}

		rebootRequired = true
	}

	// Prepare the new CPU configuration.
	if d.HasChange(mkResourceVirtualEnvironmentVMCPU) {
		cpuBlock, err := getSchemaBlock(resource, d, m, []string{mkResourceVirtualEnvironmentVMCPU}, 0, true)

		if err != nil {
			return err
		}

		cpuCores := cpuBlock[mkResourceVirtualEnvironmentVMCPUCores].(int)
		cpuHotplugged := cpuBlock[mkResourceVirtualEnvironmentVMCPUHotplugged].(int)
		cpuSockets := cpuBlock[mkResourceVirtualEnvironmentVMCPUSockets].(int)

		if cpuCores > 0 {
			body.CPUCores = &cpuCores
		}

		if cpuHotplugged > 0 {
			body.VirtualCPUCount = &cpuHotplugged
		}

		if cpuSockets > 0 {
			body.CPUSockets = &cpuSockets
		}

		rebootRequired = true
	}

	// Prepare the new cloud-init configuration.
	if d.HasChange(mkResourceVirtualEnvironmentVMCloudInit) {
		cloudInitConfig, err := resourceVirtualEnvironmentVMGetCloudConfig(d, m)

		if err != nil {
			return err
		}

		body.CloudInitConfig = cloudInitConfig

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

		body.SCSIDevices = make(proxmox.CustomStorageDevices, len(diskDeviceObjects))

		for di, do := range diskDeviceObjects {
			if scsiDevices[di] == nil {
				return fmt.Errorf("Missing SCSI device %d (scsi%d)", di, di)
			}

			body.SCSIDevices[di] = *scsiDevices[di]
			body.SCSIDevices[di].BurstableReadSpeedMbps = do.BurstableReadSpeedMbps
			body.SCSIDevices[di].BurstableWriteSpeedMbps = do.BurstableWriteSpeedMbps
			body.SCSIDevices[di].MaxReadSpeedMbps = do.MaxReadSpeedMbps
			body.SCSIDevices[di].MaxWriteSpeedMbps = do.MaxWriteSpeedMbps
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

		body.DedicatedMemory = &memoryDedicated
		body.FloatingMemory = &memoryFloating

		if memoryShared > 0 {
			memorySharedName := fmt.Sprintf("vm-%d-ivshmem", vmID)

			body.SharedMemory = &proxmox.CustomSharedMemory{
				Name: &memorySharedName,
				Size: memoryShared,
			}
		}

		rebootRequired = true
	}

	// Prepare the new network device configuration.
	if d.HasChange(mkResourceVirtualEnvironmentVMNetworkDevice) {
		body.NetworkDevices, err = resourceVirtualEnvironmentVMGetNetworkDeviceObjects(d, m)

		if err != nil {
			return err
		}

		rebootRequired = true
	}

	// Update the configuration now that everything has been prepared.
	err = veClient.UpdateVM(nodeName, vmID, body)

	if err != nil {
		return err
	}

	// Determine if the state of the virtual machine needs to be changed.
	if d.HasChange(mkResourceVirtualEnvironmentVMStarted) {
		started := d.Get(mkResourceVirtualEnvironmentVMStarted).(bool)

		if started {
			err = veClient.StartVM(nodeName, vmID)

			if err != nil {
				return err
			}

			err = veClient.WaitForState(nodeName, vmID, "running", 120, 5)

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

			err = veClient.WaitForState(nodeName, vmID, "stopped", 30, 5)

			if err != nil {
				return err
			}

			rebootRequired = false
		}
	}

	// As a final step in the update procedure, we might need to reboot the virtual machine.
	if rebootRequired {
		rebootTimeout := 300

		err = veClient.RebootVM(nodeName, vmID, &proxmox.VirtualEnvironmentVMRebootRequestBody{
			Timeout: &rebootTimeout,
		})

		if err != nil {
			return err
		}

		// Wait for the agent to unpublish the network interfaces, if it's enabled.
		if vmConfig.Agent != nil && vmConfig.Agent.Enabled != nil && *vmConfig.Agent.Enabled {
			err = veClient.WaitForNoNetworkInterfacesFromAgent(nodeName, vmID, 300, 5)

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
	forceStop := proxmox.CustomBool(true)
	shutdownTimeout := 300

	err = veClient.ShutdownVM(nodeName, vmID, &proxmox.VirtualEnvironmentVMShutdownRequestBody{
		ForceStop: &forceStop,
		Timeout:   &shutdownTimeout,
	})

	if err != nil {
		return err
	}

	err = veClient.WaitForState(nodeName, vmID, "stopped", 30, 5)

	if err != nil {
		return err
	}

	err = veClient.DeleteVM(nodeName, vmID)

	if err != nil {
		if strings.Contains(err.Error(), "HTTP 404") {
			d.SetId("")

			return nil
		}

		return err
	}

	// Wait for the state to become unavailable as that clearly indicates the destruction of the VM.
	err = veClient.WaitForState(nodeName, vmID, "", 30, 2)

	if err == nil {
		return fmt.Errorf("Failed to delete VM \"%d\"", vmID)
	}

	d.SetId("")

	return nil
}
