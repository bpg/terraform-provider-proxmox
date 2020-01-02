/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmoxtf

import (
	"strconv"
	"strings"

	"github.com/danitso/terraform-provider-proxmox/proxmox"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

const (
	dvResourceVirtualEnvironmentContainerConsoleEnabled                    = true
	dvResourceVirtualEnvironmentContainerConsoleMode                       = "tty"
	dvResourceVirtualEnvironmentContainerConsoleTTYCount                   = 2
	dvResourceVirtualEnvironmentContainerInitializationDNSDomain           = ""
	dvResourceVirtualEnvironmentContainerInitializationDNSServer           = ""
	dvResourceVirtualEnvironmentContainerInitializationIPConfigIPv4Address = ""
	dvResourceVirtualEnvironmentContainerInitializationIPConfigIPv4Gateway = ""
	dvResourceVirtualEnvironmentContainerInitializationIPConfigIPv6Address = ""
	dvResourceVirtualEnvironmentContainerInitializationIPConfigIPv6Gateway = ""
	dvResourceVirtualEnvironmentContainerInitializationHostname            = ""
	dvResourceVirtualEnvironmentContainerInitializationUserAccountPassword = ""
	dvResourceVirtualEnvironmentContainerCPUArchitecture                   = "amd64"
	dvResourceVirtualEnvironmentContainerCPUCores                          = 1
	dvResourceVirtualEnvironmentContainerCPUUnits                          = 1024
	dvResourceVirtualEnvironmentContainerDescription                       = ""
	dvResourceVirtualEnvironmentContainerDiskDatastoreID                   = "local-lvm"
	dvResourceVirtualEnvironmentContainerDiskFileFormat                    = "qcow2"
	dvResourceVirtualEnvironmentContainerDiskFileID                        = ""
	dvResourceVirtualEnvironmentContainerDiskSize                          = 8
	dvResourceVirtualEnvironmentContainerDiskSpeedRead                     = 0
	dvResourceVirtualEnvironmentContainerDiskSpeedReadBurstable            = 0
	dvResourceVirtualEnvironmentContainerDiskSpeedWrite                    = 0
	dvResourceVirtualEnvironmentContainerDiskSpeedWriteBurstable           = 0
	dvResourceVirtualEnvironmentContainerMemoryDedicated                   = 512
	dvResourceVirtualEnvironmentContainerMemorySwap                        = 0
	dvResourceVirtualEnvironmentContainerNetworkInterfaceBridge            = "vmbr0"
	dvResourceVirtualEnvironmentContainerNetworkInterfaceEnabled           = true
	dvResourceVirtualEnvironmentContainerNetworkInterfaceMACAddress        = ""
	dvResourceVirtualEnvironmentContainerNetworkInterfaceRateLimit         = 0
	dvResourceVirtualEnvironmentContainerNetworkInterfaceVLANID            = 0
	dvResourceVirtualEnvironmentContainerOperatingSystemType               = "unmanaged"
	dvResourceVirtualEnvironmentContainerPoolID                            = ""
	dvResourceVirtualEnvironmentContainerStarted                           = true
	dvResourceVirtualEnvironmentContainerVMID                              = -1

	mkResourceVirtualEnvironmentContainerConsole                           = "console"
	mkResourceVirtualEnvironmentContainerConsoleEnabled                    = "enabled"
	mkResourceVirtualEnvironmentContainerConsoleMode                       = "type"
	mkResourceVirtualEnvironmentContainerConsoleTTYCount                   = "tty_count"
	mkResourceVirtualEnvironmentContainerCPU                               = "cpu"
	mkResourceVirtualEnvironmentContainerCPUArchitecture                   = "architecture"
	mkResourceVirtualEnvironmentContainerCPUCores                          = "cores"
	mkResourceVirtualEnvironmentContainerCPUUnits                          = "units"
	mkResourceVirtualEnvironmentContainerDescription                       = "description"
	mkResourceVirtualEnvironmentContainerInitialization                    = "initialization"
	mkResourceVirtualEnvironmentContainerInitializationDNS                 = "dns"
	mkResourceVirtualEnvironmentContainerInitializationDNSDomain           = "domain"
	mkResourceVirtualEnvironmentContainerInitializationDNSServer           = "server"
	mkResourceVirtualEnvironmentContainerInitializationHostname            = "hostname"
	mkResourceVirtualEnvironmentContainerInitializationIPConfig            = "ip_config"
	mkResourceVirtualEnvironmentContainerInitializationIPConfigIPv4        = "ipv4"
	mkResourceVirtualEnvironmentContainerInitializationIPConfigIPv4Address = "address"
	mkResourceVirtualEnvironmentContainerInitializationIPConfigIPv4Gateway = "gateway"
	mkResourceVirtualEnvironmentContainerInitializationIPConfigIPv6        = "ipv6"
	mkResourceVirtualEnvironmentContainerInitializationIPConfigIPv6Address = "address"
	mkResourceVirtualEnvironmentContainerInitializationIPConfigIPv6Gateway = "gateway"
	mkResourceVirtualEnvironmentContainerInitializationUserAccount         = "user_account"
	mkResourceVirtualEnvironmentContainerInitializationUserAccountKeys     = "keys"
	mkResourceVirtualEnvironmentContainerInitializationUserAccountPassword = "password"
	mkResourceVirtualEnvironmentContainerInitializationUserAccountUsername = "username"
	mkResourceVirtualEnvironmentContainerMemory                            = "memory"
	mkResourceVirtualEnvironmentContainerMemoryDedicated                   = "dedicated"
	mkResourceVirtualEnvironmentContainerMemorySwap                        = "swap"
	mkResourceVirtualEnvironmentContainerNetworkInterface                  = "network_device"
	mkResourceVirtualEnvironmentContainerNetworkInterfaceBridge            = "bridge"
	mkResourceVirtualEnvironmentContainerNetworkInterfaceEnabled           = "enabled"
	mkResourceVirtualEnvironmentContainerNetworkInterfaceMACAddress        = "mac_address"
	mkResourceVirtualEnvironmentContainerNetworkInterfaceName              = "name"
	mkResourceVirtualEnvironmentContainerNetworkInterfaceRateLimit         = "rate_limit"
	mkResourceVirtualEnvironmentContainerNetworkInterfaceVLANID            = "vlan_id"
	mkResourceVirtualEnvironmentContainerNodeName                          = "node_name"
	mkResourceVirtualEnvironmentContainerOperatingSystem                   = "operating_system"
	mkResourceVirtualEnvironmentContainerOperatingSystemTemplateFileID     = "template_file_id"
	mkResourceVirtualEnvironmentContainerOperatingSystemType               = "type"
	mkResourceVirtualEnvironmentContainerPoolID                            = "pool_id"
	mkResourceVirtualEnvironmentContainerStarted                           = "started"
	mkResourceVirtualEnvironmentContainerVMID                              = "vm_id"
)

func resourceVirtualEnvironmentContainer() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			mkResourceVirtualEnvironmentContainerConsole: &schema.Schema{
				Type:        schema.TypeList,
				Description: "The console configuration",
				Optional:    true,
				DefaultFunc: func() (interface{}, error) {
					defaultList := make([]interface{}, 1)
					defaultMap := map[string]interface{}{}

					defaultMap[mkResourceVirtualEnvironmentContainerConsoleEnabled] = dvResourceVirtualEnvironmentContainerConsoleEnabled
					defaultMap[mkResourceVirtualEnvironmentContainerConsoleMode] = dvResourceVirtualEnvironmentContainerConsoleMode
					defaultMap[mkResourceVirtualEnvironmentContainerConsoleTTYCount] = dvResourceVirtualEnvironmentContainerConsoleTTYCount

					defaultList[0] = defaultMap

					return defaultList, nil
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mkResourceVirtualEnvironmentContainerConsoleEnabled: {
							Type:        schema.TypeBool,
							Description: "Whether to enable the console device",
							Optional:    true,
							Default:     dvResourceVirtualEnvironmentContainerConsoleEnabled,
						},
						mkResourceVirtualEnvironmentContainerConsoleMode: {
							Type:         schema.TypeString,
							Description:  "The console mode",
							Optional:     true,
							Default:      dvResourceVirtualEnvironmentContainerConsoleMode,
							ValidateFunc: resourceVirtualEnvironmentContainerGetConsoleModeValidator(),
						},
						mkResourceVirtualEnvironmentContainerConsoleTTYCount: {
							Type:         schema.TypeInt,
							Description:  "The number of available TTY",
							Optional:     true,
							Default:      dvResourceVirtualEnvironmentContainerConsoleTTYCount,
							ValidateFunc: validation.IntBetween(0, 6),
						},
					},
				},
				MaxItems: 1,
				MinItems: 0,
			},
			mkResourceVirtualEnvironmentContainerCPU: &schema.Schema{
				Type:        schema.TypeList,
				Description: "The CPU allocation",
				Optional:    true,
				DefaultFunc: func() (interface{}, error) {
					defaultList := make([]interface{}, 1)
					defaultMap := map[string]interface{}{}

					defaultMap[mkResourceVirtualEnvironmentContainerCPUArchitecture] = dvResourceVirtualEnvironmentContainerCPUArchitecture
					defaultMap[mkResourceVirtualEnvironmentContainerCPUCores] = dvResourceVirtualEnvironmentContainerCPUCores
					defaultMap[mkResourceVirtualEnvironmentContainerCPUUnits] = dvResourceVirtualEnvironmentContainerCPUUnits

					defaultList[0] = defaultMap

					return defaultList, nil
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mkResourceVirtualEnvironmentContainerCPUArchitecture: {
							Type:         schema.TypeString,
							Description:  "The CPU architecture",
							Optional:     true,
							Default:      dvResourceVirtualEnvironmentContainerCPUArchitecture,
							ValidateFunc: resourceVirtualEnvironmentContainerGetCPUArchitectureValidator(),
						},
						mkResourceVirtualEnvironmentContainerCPUCores: {
							Type:         schema.TypeInt,
							Description:  "The number of CPU cores",
							Optional:     true,
							Default:      dvResourceVirtualEnvironmentContainerCPUCores,
							ValidateFunc: validation.IntBetween(1, 128),
						},
						mkResourceVirtualEnvironmentContainerCPUUnits: {
							Type:         schema.TypeInt,
							Description:  "The CPU units",
							Optional:     true,
							Default:      dvResourceVirtualEnvironmentContainerCPUUnits,
							ValidateFunc: validation.IntBetween(0, 500000),
						},
					},
				},
				MaxItems: 1,
				MinItems: 0,
			},
			mkResourceVirtualEnvironmentContainerDescription: {
				Type:        schema.TypeString,
				Description: "The description",
				Optional:    true,
				Default:     dvResourceVirtualEnvironmentContainerDescription,
			},
			mkResourceVirtualEnvironmentContainerInitialization: &schema.Schema{
				Type:        schema.TypeList,
				Description: "The initialization configuration",
				Optional:    true,
				DefaultFunc: func() (interface{}, error) {
					return []interface{}{}, nil
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mkResourceVirtualEnvironmentContainerInitializationDNS: {
							Type:        schema.TypeList,
							Description: "The DNS configuration",
							Optional:    true,
							DefaultFunc: func() (interface{}, error) {
								return []interface{}{}, nil
							},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									mkResourceVirtualEnvironmentContainerInitializationDNSDomain: {
										Type:        schema.TypeString,
										Description: "The DNS search domain",
										Optional:    true,
										Default:     dvResourceVirtualEnvironmentContainerInitializationDNSDomain,
									},
									mkResourceVirtualEnvironmentContainerInitializationDNSServer: {
										Type:        schema.TypeString,
										Description: "The DNS server",
										Optional:    true,
										Default:     dvResourceVirtualEnvironmentContainerInitializationDNSServer,
									},
								},
							},
							MaxItems: 1,
							MinItems: 0,
						},
						mkResourceVirtualEnvironmentContainerInitializationHostname: {
							Type:        schema.TypeString,
							Description: "The hostname",
							Optional:    true,
							Default:     dvResourceVirtualEnvironmentContainerInitializationHostname,
						},
						mkResourceVirtualEnvironmentContainerInitializationIPConfig: {
							Type:        schema.TypeList,
							Description: "The IP configuration",
							Optional:    true,
							DefaultFunc: func() (interface{}, error) {
								return []interface{}{}, nil
							},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									mkResourceVirtualEnvironmentContainerInitializationIPConfigIPv4: {
										Type:        schema.TypeList,
										Description: "The IPv4 configuration",
										Optional:    true,
										DefaultFunc: func() (interface{}, error) {
											return []interface{}{}, nil
										},
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												mkResourceVirtualEnvironmentContainerInitializationIPConfigIPv4Address: {
													Type:        schema.TypeString,
													Description: "The IPv4 address",
													Optional:    true,
													Default:     dvResourceVirtualEnvironmentContainerInitializationIPConfigIPv4Address,
												},
												mkResourceVirtualEnvironmentContainerInitializationIPConfigIPv4Gateway: {
													Type:        schema.TypeString,
													Description: "The IPv4 gateway",
													Optional:    true,
													Default:     dvResourceVirtualEnvironmentContainerInitializationIPConfigIPv4Gateway,
												},
											},
										},
										MaxItems: 1,
										MinItems: 0,
									},
									mkResourceVirtualEnvironmentContainerInitializationIPConfigIPv6: {
										Type:        schema.TypeList,
										Description: "The IPv6 configuration",
										Optional:    true,
										DefaultFunc: func() (interface{}, error) {
											return []interface{}{}, nil
										},
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												mkResourceVirtualEnvironmentContainerInitializationIPConfigIPv6Address: {
													Type:        schema.TypeString,
													Description: "The IPv6 address",
													Optional:    true,
													Default:     dvResourceVirtualEnvironmentContainerInitializationIPConfigIPv6Address,
												},
												mkResourceVirtualEnvironmentContainerInitializationIPConfigIPv6Gateway: {
													Type:        schema.TypeString,
													Description: "The IPv6 gateway",
													Optional:    true,
													Default:     dvResourceVirtualEnvironmentContainerInitializationIPConfigIPv6Gateway,
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
						mkResourceVirtualEnvironmentContainerInitializationUserAccount: {
							Type:        schema.TypeList,
							Description: "The user account configuration",
							Optional:    true,
							ForceNew:    true,
							DefaultFunc: func() (interface{}, error) {
								return []interface{}{}, nil
							},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									mkResourceVirtualEnvironmentContainerInitializationUserAccountKeys: {
										Type:        schema.TypeList,
										Description: "The SSH keys",
										Optional:    true,
										ForceNew:    true,
										DefaultFunc: func() (interface{}, error) {
											return []interface{}{}, nil
										},
										Elem: &schema.Schema{Type: schema.TypeString},
									},
									mkResourceVirtualEnvironmentContainerInitializationUserAccountPassword: {
										Type:        schema.TypeString,
										Description: "The SSH password",
										Optional:    true,
										ForceNew:    true,
										Sensitive:   true,
										Default:     dvResourceVirtualEnvironmentContainerInitializationUserAccountPassword,
										DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
											return strings.ReplaceAll(old, "*", "") == ""
										},
									},
								},
							},
							MaxItems: 1,
							MinItems: 0,
						},
					},
				},
				MaxItems: 1,
				MinItems: 0,
			},
			mkResourceVirtualEnvironmentContainerMemory: &schema.Schema{
				Type:        schema.TypeList,
				Description: "The memory allocation",
				Optional:    true,
				DefaultFunc: func() (interface{}, error) {
					defaultList := make([]interface{}, 1)
					defaultMap := map[string]interface{}{}

					defaultMap[mkResourceVirtualEnvironmentContainerMemoryDedicated] = dvResourceVirtualEnvironmentContainerMemoryDedicated
					defaultMap[mkResourceVirtualEnvironmentContainerMemorySwap] = dvResourceVirtualEnvironmentContainerMemorySwap

					defaultList[0] = defaultMap

					return defaultList, nil
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mkResourceVirtualEnvironmentContainerMemoryDedicated: {
							Type:         schema.TypeInt,
							Description:  "The dedicated memory in megabytes",
							Optional:     true,
							Default:      dvResourceVirtualEnvironmentContainerMemoryDedicated,
							ValidateFunc: validation.IntBetween(16, 268435456),
						},
						mkResourceVirtualEnvironmentContainerMemorySwap: {
							Type:         schema.TypeInt,
							Description:  "The swap size in megabytes",
							Optional:     true,
							Default:      dvResourceVirtualEnvironmentContainerMemorySwap,
							ValidateFunc: validation.IntBetween(0, 268435456),
						},
					},
				},
				MaxItems: 1,
				MinItems: 0,
			},
			mkResourceVirtualEnvironmentContainerNetworkInterface: &schema.Schema{
				Type:        schema.TypeList,
				Description: "The network interfaces",
				Optional:    true,
				DefaultFunc: func() (interface{}, error) {
					return make([]interface{}, 1), nil
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mkResourceVirtualEnvironmentContainerNetworkInterfaceBridge: {
							Type:        schema.TypeString,
							Description: "The bridge",
							Optional:    true,
							Default:     dvResourceVirtualEnvironmentContainerNetworkInterfaceBridge,
						},
						mkResourceVirtualEnvironmentContainerNetworkInterfaceEnabled: {
							Type:        schema.TypeBool,
							Description: "Whether to enable the network device",
							Optional:    true,
							Default:     dvResourceVirtualEnvironmentContainerNetworkInterfaceEnabled,
						},
						mkResourceVirtualEnvironmentContainerNetworkInterfaceMACAddress: {
							Type:        schema.TypeString,
							Description: "The MAC address",
							Optional:    true,
							Default:     dvResourceVirtualEnvironmentContainerNetworkInterfaceMACAddress,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								return new == ""
							},
							ValidateFunc: getMACAddressValidator(),
						},
						mkResourceVirtualEnvironmentContainerNetworkInterfaceName: {
							Type:        schema.TypeString,
							Description: "The network interface name",
							Required:    true,
						},
						mkResourceVirtualEnvironmentContainerNetworkInterfaceRateLimit: {
							Type:        schema.TypeFloat,
							Description: "The rate limit in megabytes per second",
							Optional:    true,
							Default:     dvResourceVirtualEnvironmentContainerNetworkInterfaceRateLimit,
						},
						mkResourceVirtualEnvironmentContainerNetworkInterfaceVLANID: {
							Type:        schema.TypeInt,
							Description: "The VLAN identifier",
							Optional:    true,
							Default:     dvResourceVirtualEnvironmentContainerNetworkInterfaceVLANID,
						},
					},
				},
				MaxItems: 8,
				MinItems: 0,
			},
			mkResourceVirtualEnvironmentContainerNodeName: &schema.Schema{
				Type:        schema.TypeString,
				Description: "The node name",
				Required:    true,
				ForceNew:    true,
			},
			mkResourceVirtualEnvironmentContainerOperatingSystem: &schema.Schema{
				Type:        schema.TypeList,
				Description: "The operating system configuration",
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mkResourceVirtualEnvironmentContainerOperatingSystemTemplateFileID: {
							Type:         schema.TypeString,
							Description:  "The ID of an OS template file",
							Required:     true,
							ForceNew:     true,
							ValidateFunc: getFileIDValidator(),
						},
						mkResourceVirtualEnvironmentContainerOperatingSystemType: {
							Type:         schema.TypeString,
							Description:  "The type",
							Optional:     true,
							Default:      dvResourceVirtualEnvironmentContainerOperatingSystemType,
							ValidateFunc: resourceVirtualEnvironmentContainerGetOperatingSystemTypeValidator(),
						},
					},
				},
				MaxItems: 1,
				MinItems: 0,
			},
			mkResourceVirtualEnvironmentContainerPoolID: {
				Type:        schema.TypeString,
				Description: "The ID of the pool to assign the virtual machine to",
				Optional:    true,
				ForceNew:    true,
				Default:     dvResourceVirtualEnvironmentContainerPoolID,
			},
			mkResourceVirtualEnvironmentContainerStarted: {
				Type:        schema.TypeBool,
				Description: "Whether to start the container",
				Optional:    true,
				Default:     dvResourceVirtualEnvironmentContainerStarted,
			},
			mkResourceVirtualEnvironmentContainerVMID: {
				Type:         schema.TypeInt,
				Description:  "The VM identifier",
				Optional:     true,
				ForceNew:     true,
				Default:      dvResourceVirtualEnvironmentContainerVMID,
				ValidateFunc: getVMIDValidator(),
			},
		},
		Create: resourceVirtualEnvironmentContainerCreate,
		Read:   resourceVirtualEnvironmentContainerRead,
		Update: resourceVirtualEnvironmentContainerUpdate,
		Delete: resourceVirtualEnvironmentContainerDelete,
	}
}

func resourceVirtualEnvironmentContainerCreate(d *schema.ResourceData, m interface{}) error {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()

	if err != nil {
		return err
	}

	nodeName := d.Get(mkResourceVirtualEnvironmentContainerNodeName).(string)
	resource := resourceVirtualEnvironmentContainer()

	consoleBlock, err := getSchemaBlock(resource, d, m, []string{mkResourceVirtualEnvironmentContainerConsole}, 0, true)

	if err != nil {
		return err
	}

	consoleEnabled := proxmox.CustomBool(consoleBlock[mkResourceVirtualEnvironmentContainerConsoleEnabled].(bool))
	consoleMode := consoleBlock[mkResourceVirtualEnvironmentContainerConsoleMode].(string)
	consoleTTYCount := consoleBlock[mkResourceVirtualEnvironmentContainerConsoleTTYCount].(int)

	cpuBlock, err := getSchemaBlock(resource, d, m, []string{mkResourceVirtualEnvironmentContainerCPU}, 0, true)

	if err != nil {
		return err
	}

	cpuArchitecture := cpuBlock[mkResourceVirtualEnvironmentContainerCPUArchitecture].(string)
	cpuCores := cpuBlock[mkResourceVirtualEnvironmentContainerCPUCores].(int)
	cpuUnits := cpuBlock[mkResourceVirtualEnvironmentContainerCPUUnits].(int)

	description := d.Get(mkResourceVirtualEnvironmentContainerDescription).(string)

	initialization := d.Get(mkResourceVirtualEnvironmentContainerInitialization).([]interface{})
	initializationDNSDomain := dvResourceVirtualEnvironmentContainerInitializationDNSDomain
	initializationDNSServer := dvResourceVirtualEnvironmentContainerInitializationDNSServer
	initializationHostname := dvResourceVirtualEnvironmentContainerInitializationHostname
	initializationIPConfigIPv4Address := []string{}
	initializationIPConfigIPv4Gateway := []string{}
	initializationIPConfigIPv6Address := []string{}
	initializationIPConfigIPv6Gateway := []string{}
	initializationUserAccountKeys := proxmox.VirtualEnvironmentContainerCustomSSHKeys{}
	initializationUserAccountPassword := dvResourceVirtualEnvironmentContainerInitializationUserAccountPassword

	if len(initialization) > 0 {
		initializationBlock := initialization[0].(map[string]interface{})
		initializationDNS := initializationBlock[mkResourceVirtualEnvironmentContainerInitializationDNS].([]interface{})

		if len(initializationDNS) > 0 {
			initializationDNSBlock := initializationDNS[0].(map[string]interface{})
			initializationDNSDomain = initializationDNSBlock[mkResourceVirtualEnvironmentContainerInitializationDNSDomain].(string)
			initializationDNSServer = initializationDNSBlock[mkResourceVirtualEnvironmentContainerInitializationDNSServer].(string)
		}

		initializationHostname = initializationBlock[mkResourceVirtualEnvironmentContainerInitializationHostname].(string)
		initializationIPConfig := initializationBlock[mkResourceVirtualEnvironmentContainerInitializationIPConfig].([]interface{})

		for _, c := range initializationIPConfig {
			configBlock := c.(map[string]interface{})
			ipv4 := configBlock[mkResourceVirtualEnvironmentContainerInitializationIPConfigIPv4].([]interface{})

			if len(ipv4) > 0 {
				ipv4Block := ipv4[0].(map[string]interface{})

				initializationIPConfigIPv4Address = append(initializationIPConfigIPv4Address, ipv4Block[mkResourceVirtualEnvironmentContainerInitializationIPConfigIPv4Address].(string))
				initializationIPConfigIPv4Gateway = append(initializationIPConfigIPv4Gateway, ipv4Block[mkResourceVirtualEnvironmentContainerInitializationIPConfigIPv4Gateway].(string))
			} else {
				initializationIPConfigIPv4Address = append(initializationIPConfigIPv4Address, "")
				initializationIPConfigIPv4Gateway = append(initializationIPConfigIPv4Gateway, "")
			}

			ipv6 := configBlock[mkResourceVirtualEnvironmentContainerInitializationIPConfigIPv6].([]interface{})

			if len(ipv6) > 0 {
				ipv6Block := ipv6[0].(map[string]interface{})

				initializationIPConfigIPv6Address = append(initializationIPConfigIPv6Address, ipv6Block[mkResourceVirtualEnvironmentContainerInitializationIPConfigIPv6Address].(string))
				initializationIPConfigIPv6Gateway = append(initializationIPConfigIPv6Gateway, ipv6Block[mkResourceVirtualEnvironmentContainerInitializationIPConfigIPv6Gateway].(string))
			} else {
				initializationIPConfigIPv6Address = append(initializationIPConfigIPv6Address, "")
				initializationIPConfigIPv6Gateway = append(initializationIPConfigIPv6Gateway, "")
			}
		}

		initializationUserAccount := initializationBlock[mkResourceVirtualEnvironmentContainerInitializationUserAccount].([]interface{})

		if len(initializationUserAccount) > 0 {
			initializationUserAccountBlock := initializationUserAccount[0].(map[string]interface{})

			keys := initializationUserAccountBlock[mkResourceVirtualEnvironmentContainerInitializationUserAccountKeys].([]interface{})
			initializationUserAccountKeys := make(proxmox.VirtualEnvironmentContainerCustomSSHKeys, len(keys))

			for ki, kv := range keys {
				initializationUserAccountKeys[ki] = kv.(string)
			}

			initializationUserAccountPassword = initializationUserAccountBlock[mkResourceVirtualEnvironmentContainerInitializationUserAccountPassword].(string)
		}
	}

	memoryBlock, err := getSchemaBlock(resource, d, m, []string{mkResourceVirtualEnvironmentContainerMemory}, 0, true)

	if err != nil {
		return err
	}

	memoryDedicated := memoryBlock[mkResourceVirtualEnvironmentContainerMemoryDedicated].(int)
	memorySwap := memoryBlock[mkResourceVirtualEnvironmentContainerMemorySwap].(int)

	networkInterface := d.Get(mkResourceVirtualEnvironmentContainerNetworkInterface).([]interface{})
	networkInterfaceArray := make(proxmox.VirtualEnvironmentContainerCustomNetworkInterfaceArray, len(networkInterface))

	for ni, nv := range networkInterface {
		networkInterfaceMap := nv.(map[string]interface{})
		networkInterfaceObject := proxmox.VirtualEnvironmentContainerCustomNetworkInterface{}

		bridge := networkInterfaceMap[mkResourceVirtualEnvironmentContainerNetworkInterfaceBridge].(string)
		enabled := networkInterfaceMap[mkResourceVirtualEnvironmentContainerNetworkInterfaceEnabled].(bool)
		macAddress := networkInterfaceMap[mkResourceVirtualEnvironmentContainerNetworkInterfaceMACAddress].(string)
		name := networkInterfaceMap[mkResourceVirtualEnvironmentContainerNetworkInterfaceName].(string)
		rateLimit := networkInterfaceMap[mkResourceVirtualEnvironmentContainerNetworkInterfaceRateLimit].(int)
		vlanID := networkInterfaceMap[mkResourceVirtualEnvironmentContainerNetworkInterfaceVLANID].(int)

		if bridge != "" {
			networkInterfaceObject.Bridge = &bridge
		}

		networkInterfaceObject.Enabled = enabled

		if len(initializationIPConfigIPv4Address) > ni {
			if initializationIPConfigIPv4Address[ni] != "" {
				networkInterfaceObject.IPv4Address = &initializationIPConfigIPv4Address[ni]
			}

			if initializationIPConfigIPv4Gateway[ni] != "" {
				networkInterfaceObject.IPv4Gateway = &initializationIPConfigIPv4Gateway[ni]
			}

			if initializationIPConfigIPv6Address[ni] != "" {
				networkInterfaceObject.IPv6Address = &initializationIPConfigIPv6Address[ni]
			}

			if initializationIPConfigIPv6Gateway[ni] != "" {
				networkInterfaceObject.IPv6Gateway = &initializationIPConfigIPv6Gateway[ni]
			}
		}

		if macAddress != "" {
			networkInterfaceObject.MACAddress = &macAddress
		}

		if name != "" {
			networkInterfaceObject.Name = name
		}

		if rateLimit != 0 {
			networkInterfaceObject.RateLimit = &rateLimit
		}

		if vlanID != 0 {
			networkInterfaceObject.Tag = &vlanID
		}

		networkInterfaceArray[ni] = networkInterfaceObject
	}

	operatingSystem, err := getSchemaBlock(resource, d, m, []string{mkResourceVirtualEnvironmentContainerOperatingSystem}, 0, true)

	if err != nil {
		return err
	}

	operatingSystemTemplateFileID := operatingSystem[mkResourceVirtualEnvironmentContainerOperatingSystemTemplateFileID].(string)
	operatingSystemType := operatingSystem[mkResourceVirtualEnvironmentContainerOperatingSystemType].(string)

	poolID := d.Get(mkResourceVirtualEnvironmentContainerPoolID).(string)
	started := proxmox.CustomBool(d.Get(mkResourceVirtualEnvironmentContainerStarted).(bool))
	vmID := d.Get(mkResourceVirtualEnvironmentContainerVMID).(int)

	if vmID == -1 {
		vmIDNew, err := veClient.GetVMID()

		if err != nil {
			return err
		}

		vmID = *vmIDNew
	}

	// Attempt to create the resource using the retrieved values.
	body := proxmox.VirtualEnvironmentContainerCreateRequestBody{
		ConsoleEnabled:       &consoleEnabled,
		ConsoleMode:          &consoleMode,
		CPUArchitecture:      &cpuArchitecture,
		CPUCores:             &cpuCores,
		CPUUnits:             &cpuUnits,
		DedicatedMemory:      &memoryDedicated,
		NetworkInterfaces:    networkInterfaceArray,
		OSTemplateFileVolume: operatingSystemTemplateFileID,
		OSType:               &operatingSystemType,
		StartOnBoot:          &started,
		Swap:                 &memorySwap,
		TTY:                  &consoleTTYCount,
		VMID:                 vmID,
	}

	if description != "" {
		body.Description = &description
	}

	if initializationDNSDomain != "" {
		body.DNSDomain = &initializationDNSDomain
	}

	if initializationDNSServer != "" {
		body.DNSServer = &initializationDNSServer
	}

	if initializationHostname != "" {
		body.Hostname = &initializationHostname
	}

	if len(initializationUserAccountKeys) > 0 {
		body.SSHKeys = &initializationUserAccountKeys
	}

	if initializationUserAccountPassword != "" {
		body.Password = &initializationUserAccountPassword
	}

	if poolID != "" {
		body.PoolID = &poolID
	}

	err = veClient.CreateContainer(nodeName, &body)

	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(vmID))

	return resourceVirtualEnvironmentContainerCreateStart(d, m)
}

func resourceVirtualEnvironmentContainerCreateStart(d *schema.ResourceData, m interface{}) error {
	started := d.Get(mkResourceVirtualEnvironmentContainerStarted).(bool)

	if !started {
		return resourceVirtualEnvironmentContainerRead(d, m)
	}

	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()

	if err != nil {
		return err
	}

	nodeName := d.Get(mkResourceVirtualEnvironmentContainerNodeName).(string)
	vmID, err := strconv.Atoi(d.Id())

	if err != nil {
		return err
	}

	// Start the container and wait for it to reach a running state before continuing.
	err = veClient.StartContainer(nodeName, vmID)

	if err != nil {
		return err
	}

	err = veClient.WaitForContainerState(nodeName, vmID, "running", 120, 5)

	if err != nil {
		return err
	}

	return resourceVirtualEnvironmentContainerRead(d, m)
}

func resourceVirtualEnvironmentContainerGetConsoleModeValidator() schema.SchemaValidateFunc {
	return validation.StringInSlice([]string{
		"console",
		"shell",
		"tty",
	}, false)
}

func resourceVirtualEnvironmentContainerGetCPUArchitectureValidator() schema.SchemaValidateFunc {
	return validation.StringInSlice([]string{
		"amd64",
		"arm64",
		"armhf",
		"i386",
	}, false)
}

func resourceVirtualEnvironmentContainerGetOperatingSystemTypeValidator() schema.SchemaValidateFunc {
	return validation.StringInSlice([]string{
		"alpine",
		"archlinux",
		"centos",
		"debian",
		"fedora",
		"gentoo",
		"opensuse",
		"ubuntu",
		"unmanaged",
	}, false)
}

func resourceVirtualEnvironmentContainerRead(d *schema.ResourceData, m interface{}) error {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()

	if err != nil {
		return err
	}

	nodeName := d.Get(mkResourceVirtualEnvironmentContainerNodeName).(string)
	vmID, err := strconv.Atoi(d.Id())

	if err != nil {
		return err
	}

	// Retrieve the entire configuration in order to compare it to the state.
	_, err = veClient.GetContainer(nodeName, vmID)

	if err != nil {
		if strings.Contains(err.Error(), "HTTP 404") {
			d.SetId("")

			return nil
		}

		return err
	}

	return nil
}

func resourceVirtualEnvironmentContainerUpdate(d *schema.ResourceData, m interface{}) error {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()

	if err != nil {
		return err
	}

	nodeName := d.Get(mkResourceVirtualEnvironmentContainerNodeName).(string)
	vmID, err := strconv.Atoi(d.Id())

	if err != nil {
		return err
	}

	// Retrieve the entire configuration as we need to process certain values.
	_, err = veClient.GetContainer(nodeName, vmID)

	if err != nil {
		return err
	}

	return resourceVirtualEnvironmentContainerRead(d, m)
}

func resourceVirtualEnvironmentContainerDelete(d *schema.ResourceData, m interface{}) error {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()

	if err != nil {
		return err
	}

	nodeName := d.Get(mkResourceVirtualEnvironmentContainerNodeName).(string)
	vmID, err := strconv.Atoi(d.Id())

	if err != nil {
		return err
	}

	err = veClient.DeleteContainer(nodeName, vmID)

	if err != nil {
		if strings.Contains(err.Error(), "HTTP 404") {
			d.SetId("")

			return nil
		}

		return err
	}

	d.SetId("")

	return nil
}
