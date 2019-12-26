/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmoxtf

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/danitso/terraform-provider-proxmox/proxmox"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

const (
	dvResourceVirtualEnvironmentVMCDROMEnabled            = false
	dvResourceVirtualEnvironmentVMCDROMFileID             = ""
	dvResourceVirtualEnvironmentVMCloudInitDNSDomain      = ""
	dvResourceVirtualEnvironmentVMCloudInitDNSServer      = ""
	dvResourceVirtualEnvironmentVMCPUCores                = 1
	dvResourceVirtualEnvironmentVMCPUSockets              = 1
	dvResourceVirtualEnvironmentVMDiskDatastoreID         = "local-lvm"
	dvResourceVirtualEnvironmentVMDiskEnabled             = true
	dvResourceVirtualEnvironmentVMDiskFileFormat          = "qcow2"
	dvResourceVirtualEnvironmentVMDiskFileID              = ""
	dvResourceVirtualEnvironmentVMDiskSize                = 8
	dvResourceVirtualEnvironmentVMKeyboardLayout          = "en-us"
	dvResourceVirtualEnvironmentVMMemoryDedicated         = 512
	dvResourceVirtualEnvironmentVMMemoryFloating          = 0
	dvResourceVirtualEnvironmentVMMemoryShared            = 0
	dvResourceVirtualEnvironmentVMName                    = ""
	dvResourceVirtualEnvironmentVMNetworkDeviceBridge     = "vmbr0"
	dvResourceVirtualEnvironmentVMNetworkDeviceEnabled    = true
	dvResourceVirtualEnvironmentVMNetworkDeviceMACAddress = ""
	dvResourceVirtualEnvironmentVMNetworkDeviceModel      = "virtio"
	dvResourceVirtualEnvironmentVMNetworkDeviceVLANID     = -1
	dvResourceVirtualEnvironmentVMOSType                  = "other"
	dvResourceVirtualEnvironmentVMVMID                    = -1

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
	mkResourceVirtualEnvironmentVMCloudInitUserAccountUsername = "username"
	mkResourceVirtualEnvironmentVMCPU                          = "cpu"
	mkResourceVirtualEnvironmentVMCPUCores                     = "cores"
	mkResourceVirtualEnvironmentVMCPUSockets                   = "sockets"
	mkResourceVirtualEnvironmentVMDisk                         = "disk"
	mkResourceVirtualEnvironmentVMDiskDatastoreID              = "datastore_id"
	mkResourceVirtualEnvironmentVMDiskEnabled                  = "enabled"
	mkResourceVirtualEnvironmentVMDiskFileFormat               = "file_format"
	mkResourceVirtualEnvironmentVMDiskFileID                   = "file_id"
	mkResourceVirtualEnvironmentVMDiskSize                     = "size"
	mkResourceVirtualEnvironmentVMKeyboardLayout               = "keyboard_layout"
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
	mkResourceVirtualEnvironmentVMNetworkDeviceVLANID          = "vlan_id"
	mkResourceVirtualEnvironmentVMNodeName                     = "node_name"
	mkResourceVirtualEnvironmentVMOSType                       = "os_type"
	mkResourceVirtualEnvironmentVMVMID                         = "vm_id"
)

func resourceVirtualEnvironmentVM() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			mkResourceVirtualEnvironmentVMCDROM: &schema.Schema{
				Type:        schema.TypeList,
				Description: "The CDROM drive",
				Optional:    true,
				DefaultFunc: func() (interface{}, error) {
					defaultList := make([]interface{}, 1)
					defaultMap := make(map[string]interface{})

					defaultMap[mkResourceVirtualEnvironmentVMCDROMEnabled] = dvResourceVirtualEnvironmentVMCDROMEnabled
					defaultMap[mkResourceVirtualEnvironmentVMCDROMFileID] = dvResourceVirtualEnvironmentVMCDROMFileID

					defaultList[0] = defaultMap

					return defaultList, nil
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mkResourceVirtualEnvironmentVMCDROMEnabled: {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to enable the CDROM drive",
							Default:     dvResourceVirtualEnvironmentVMCDROMEnabled,
						},
						mkResourceVirtualEnvironmentVMCDROMFileID: {
							Type:         schema.TypeString,
							Optional:     true,
							Description:  "The file id",
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
					return make([]interface{}, 0), nil
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mkResourceVirtualEnvironmentVMCloudInitDNS: {
							Type:        schema.TypeList,
							Description: "The DNS configuration",
							Optional:    true,
							DefaultFunc: func() (interface{}, error) {
								return make([]interface{}, 0), nil
							},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									mkResourceVirtualEnvironmentVMCloudInitDNSDomain: {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The DNS search domain",
										Default:     dvResourceVirtualEnvironmentVMCloudInitDNSDomain,
									},
									mkResourceVirtualEnvironmentVMCloudInitDNSServer: {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The DNS server",
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
								return make([]interface{}, 0), nil
							},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									mkResourceVirtualEnvironmentVMCloudInitIPConfigIPv4: {
										Type:        schema.TypeList,
										Description: "The IPv4 configuration",
										Optional:    true,
										DefaultFunc: func() (interface{}, error) {
											return make([]interface{}, 0), nil
										},
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												mkResourceVirtualEnvironmentVMCloudInitIPConfigIPv4Address: {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The IPv4 address",
													Default:     "",
												},
												mkResourceVirtualEnvironmentVMCloudInitIPConfigIPv4Gateway: {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The IPv4 gateway",
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
											return make([]interface{}, 0), nil
										},
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												mkResourceVirtualEnvironmentVMCloudInitIPConfigIPv6Address: {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The IPv6 address",
													Default:     "",
												},
												mkResourceVirtualEnvironmentVMCloudInitIPConfigIPv6Gateway: {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The IPv6 gateway",
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
							DefaultFunc: func() (interface{}, error) {
								return make([]interface{}, 0), nil
							},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									mkResourceVirtualEnvironmentVMCloudInitUserAccountKeys: {
										Type:        schema.TypeList,
										Required:    true,
										Description: "The SSH keys",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									mkResourceVirtualEnvironmentVMCloudInitUserAccountUsername: {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The SSH username",
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
			mkResourceVirtualEnvironmentVMCPU: &schema.Schema{
				Type:        schema.TypeList,
				Description: "The CPU allocation",
				Optional:    true,
				DefaultFunc: func() (interface{}, error) {
					defaultList := make([]interface{}, 1)
					defaultMap := make(map[string]interface{})

					defaultMap[mkResourceVirtualEnvironmentVMCPUCores] = dvResourceVirtualEnvironmentVMCPUCores
					defaultMap[mkResourceVirtualEnvironmentVMCPUSockets] = dvResourceVirtualEnvironmentVMCPUSockets

					defaultList[0] = defaultMap

					return defaultList, nil
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mkResourceVirtualEnvironmentVMCPUCores: {
							Type:         schema.TypeInt,
							Optional:     true,
							Description:  "The number of CPU cores",
							Default:      dvResourceVirtualEnvironmentVMCPUCores,
							ValidateFunc: validation.IntBetween(1, 2304),
						},
						mkResourceVirtualEnvironmentVMCPUSockets: {
							Type:         schema.TypeInt,
							Optional:     true,
							Description:  "The number of CPU sockets",
							Default:      dvResourceVirtualEnvironmentVMCPUSockets,
							ValidateFunc: validation.IntBetween(1, 16),
						},
					},
				},
				MaxItems: 1,
				MinItems: 0,
			},
			mkResourceVirtualEnvironmentVMDisk: &schema.Schema{
				Type:        schema.TypeList,
				Description: "The disk devices",
				Optional:    true,
				DefaultFunc: func() (interface{}, error) {
					defaultList := make([]interface{}, 1)
					defaultMap := make(map[string]interface{})

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
							Optional:    true,
							Description: "The datastore id",
							Default:     dvResourceVirtualEnvironmentVMDiskDatastoreID,
						},
						mkResourceVirtualEnvironmentVMDiskEnabled: {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to enable the disk",
							Default:     dvResourceVirtualEnvironmentVMDiskEnabled,
						},
						mkResourceVirtualEnvironmentVMDiskFileFormat: {
							Type:         schema.TypeString,
							Optional:     true,
							Description:  "The file format",
							Default:      dvResourceVirtualEnvironmentVMDiskFileFormat,
							ValidateFunc: getFileFormatValidator(),
						},
						mkResourceVirtualEnvironmentVMDiskFileID: {
							Type:         schema.TypeString,
							Optional:     true,
							Description:  "The file id for a disk image",
							Default:      dvResourceVirtualEnvironmentVMDiskFileID,
							ValidateFunc: getFileIDValidator(),
						},
						mkResourceVirtualEnvironmentVMDiskSize: {
							Type:         schema.TypeInt,
							Optional:     true,
							Description:  "The disk size in gigabytes",
							Default:      dvResourceVirtualEnvironmentVMDiskSize,
							ValidateFunc: validation.IntBetween(1, 8192),
						},
					},
				},
				MaxItems: 14,
				MinItems: 0,
			},
			mkResourceVirtualEnvironmentVMKeyboardLayout: {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The keyboard layout",
				Default:      dvResourceVirtualEnvironmentVMKeyboardLayout,
				ValidateFunc: getKeyboardLayoutValidator(),
			},
			mkResourceVirtualEnvironmentVMMemory: &schema.Schema{
				Type:        schema.TypeList,
				Description: "The memory allocation",
				Optional:    true,
				DefaultFunc: func() (interface{}, error) {
					defaultList := make([]interface{}, 1)
					defaultMap := make(map[string]interface{})

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
							Optional:     true,
							Description:  "The dedicated memory in megabytes",
							Default:      dvResourceVirtualEnvironmentVMMemoryDedicated,
							ValidateFunc: validation.IntBetween(64, 268435456),
						},
						mkResourceVirtualEnvironmentVMMemoryFloating: {
							Type:         schema.TypeInt,
							Optional:     true,
							Description:  "The floating memory in megabytes (balloon)",
							Default:      dvResourceVirtualEnvironmentVMMemoryFloating,
							ValidateFunc: validation.IntBetween(0, 268435456),
						},
						mkResourceVirtualEnvironmentVMMemoryShared: {
							Type:         schema.TypeInt,
							Optional:     true,
							Description:  "The shared memory in megabytes",
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
				Optional:    true,
				Description: "The name",
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
							Optional:    true,
							Description: "The bridge",
							Default:     dvResourceVirtualEnvironmentVMNetworkDeviceBridge,
						},
						mkResourceVirtualEnvironmentVMNetworkDeviceEnabled: {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to enable the network device",
							Default:     dvResourceVirtualEnvironmentVMNetworkDeviceEnabled,
						},
						mkResourceVirtualEnvironmentVMNetworkDeviceMACAddress: {
							Type:         schema.TypeString,
							Optional:     true,
							Description:  "The MAC address",
							Default:      dvResourceVirtualEnvironmentVMNetworkDeviceMACAddress,
							ValidateFunc: getMACAddressValidator(),
						},
						mkResourceVirtualEnvironmentVMNetworkDeviceModel: {
							Type:         schema.TypeString,
							Optional:     true,
							Description:  "The model",
							Default:      dvResourceVirtualEnvironmentVMNetworkDeviceModel,
							ValidateFunc: getNetworkDeviceModelValidator(),
						},
						mkResourceVirtualEnvironmentVMNetworkDeviceVLANID: {
							Type:         schema.TypeInt,
							Optional:     true,
							Description:  "The VLAN identifier",
							Default:      dvResourceVirtualEnvironmentVMNetworkDeviceVLANID,
							ValidateFunc: getVLANIDValidator(),
						},
					},
				},
				MaxItems: 8,
				MinItems: 0,
			},
			mkResourceVirtualEnvironmentVMNodeName: &schema.Schema{
				Type:        schema.TypeString,
				Description: "The node name",
				Required:    true,
			},
			mkResourceVirtualEnvironmentVMOSType: {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The OS type",
				Default:      dvResourceVirtualEnvironmentVMOSType,
				ValidateFunc: getOSTypeValidator(),
			},
			mkResourceVirtualEnvironmentVMVMID: {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				Description:  "The VM identifier",
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

	schema := resourceVirtualEnvironmentVM().Schema
	cdrom := d.Get(mkResourceVirtualEnvironmentVMCDROM).([]interface{})

	if len(cdrom) == 0 {
		cdromDefault, err := schema[mkResourceVirtualEnvironmentVMCDROM].DefaultValue()

		if err != nil {
			return err
		}

		cdrom = cdromDefault.([]interface{})
	}

	cdromBlock := cdrom[0].(map[string]interface{})
	cdromEnabled := cdromBlock[mkResourceVirtualEnvironmentVMCDROMEnabled].(bool)
	cdromFileID := cdromBlock[mkResourceVirtualEnvironmentVMCDROMFileID].(string)

	if cdromFileID == "" {
		cdromFileID = "cdrom"
	}

	var cloudInitConfig *proxmox.CustomCloudInitConfig

	cloudInit := d.Get(mkResourceVirtualEnvironmentVMCloudInit).([]interface{})

	if len(cloudInit) > 0 {
		cdromEnabled = true
		cdromFileID = "local-lvm:cloudinit"

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

			username := cloudInitUserAccountBlock[mkResourceVirtualEnvironmentVMCloudInitUserAccountUsername].(string)

			cloudInitConfig.Username = &username
		}
	}

	cpu := d.Get(mkResourceVirtualEnvironmentVMCPU).([]interface{})

	if len(cpu) == 0 {
		cpuDefault, err := schema[mkResourceVirtualEnvironmentVMCPU].DefaultValue()

		if err != nil {
			return err
		}

		cpu = cpuDefault.([]interface{})
	}

	cpuBlock := cpu[0].(map[string]interface{})
	cpuCores := cpuBlock[mkResourceVirtualEnvironmentVMCPUCores].(int)
	cpuSockets := cpuBlock[mkResourceVirtualEnvironmentVMCPUSockets].(int)

	keyboardLayout := d.Get(mkResourceVirtualEnvironmentVMKeyboardLayout).(string)
	memory := d.Get(mkResourceVirtualEnvironmentVMMemory).([]interface{})

	if len(memory) == 0 {
		memoryDefault, err := schema[mkResourceVirtualEnvironmentVMMemory].DefaultValue()

		if err != nil {
			return err
		}

		memory = memoryDefault.([]interface{})
	}

	memoryBlock := memory[0].(map[string]interface{})
	memoryDedicated := memoryBlock[mkResourceVirtualEnvironmentVMMemoryDedicated].(int)
	memoryFloating := memoryBlock[mkResourceVirtualEnvironmentVMMemoryFloating].(int)
	memoryShared := memoryBlock[mkResourceVirtualEnvironmentVMMemoryShared].(int)

	name := d.Get(mkResourceVirtualEnvironmentVMName).(string)

	networkDevice := d.Get(mkResourceVirtualEnvironmentVMNetworkDevice).([]interface{})
	networkDeviceObjects := make(proxmox.CustomNetworkDevices, len(networkDevice))

	for i, d := range networkDevice {
		block := d.(map[string]interface{})

		bridge, _ := block[mkResourceVirtualEnvironmentVMNetworkDeviceBridge].(string)
		enabled, _ := block[mkResourceVirtualEnvironmentVMNetworkDeviceEnabled].(bool)
		macAddress, _ := block[mkResourceVirtualEnvironmentVMNetworkDeviceMACAddress].(string)
		model, _ := block[mkResourceVirtualEnvironmentVMNetworkDeviceModel].(string)
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

		if vlanID != -1 {
			device.Trunks = []int{vlanID}
		}

		networkDeviceObjects[i] = device
	}

	nodeName := d.Get(mkResourceVirtualEnvironmentVMNodeName).(string)
	osType := d.Get(mkResourceVirtualEnvironmentVMOSType).(string)
	vmID := d.Get(mkResourceVirtualEnvironmentVMVMID).(int)

	if vmID == -1 {
		vmIDNew, err := veClient.GetVMID()

		if err != nil {
			return err
		}

		vmID = *vmIDNew
	}

	var memorySharedObject *proxmox.CustomSharedMemory

	agentEnabled := proxmox.CustomBool(true)
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
		Agent:               &proxmox.CustomAgent{Enabled: &agentEnabled},
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
		SCSIHardware:        &scsiHardware,
		SerialDevices:       []string{"socket"},
		SharedMemory:        memorySharedObject,
		TabletDeviceEnabled: &tabletDeviceEnabled,
		VGADevice:           &proxmox.CustomVGADevice{Type: "serial0"},
		VMID:                &vmID,
	}

	if name != "" {
		body.Name = &name
	}

	err = veClient.CreateVM(nodeName, body)

	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(vmID))

	return resourceVirtualEnvironmentVMRead(d, m)
}

func resourceVirtualEnvironmentVMRead(d *schema.ResourceData, m interface{}) error {
	/*
		config := m.(providerConfiguration)
		veClient, err := config.GetVEClient()

		if err != nil {
			return err
		}
	*/

	return nil
}

func resourceVirtualEnvironmentVMUpdate(d *schema.ResourceData, m interface{}) error {
	/*
		config := m.(providerConfiguration)
		veClient, err := config.GetVEClient()

		if err != nil {
			return err
		}
	*/

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

	err = veClient.DeleteVM(nodeName, vmID)

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
