package vm

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/resource/validator"
	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/structure"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
	dvResourceVirtualEnvironmentVMCDROMInterface                    = "ide3"
	dvResourceVirtualEnvironmentVMCloneDatastoreID                  = ""
	dvResourceVirtualEnvironmentVMCloneNodeName                     = ""
	dvResourceVirtualEnvironmentVMCloneFull                         = true
	dvResourceVirtualEnvironmentVMCloneRetries                      = 1
	dvResourceVirtualEnvironmentVMCPUArchitecture                   = "x86_64"
	dvResourceVirtualEnvironmentVMCPUCores                          = 1
	dvResourceVirtualEnvironmentVMCPUHotplugged                     = 0
	dvResourceVirtualEnvironmentVMCPULimit                          = 0
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
	dvResourceVirtualEnvironmentVMDiskDiscard                       = "ignore"
	dvResourceVirtualEnvironmentVMDiskCache                         = "none"
	dvResourceVirtualEnvironmentVMDiskSpeedRead                     = 0
	dvResourceVirtualEnvironmentVMDiskSpeedReadBurstable            = 0
	dvResourceVirtualEnvironmentVMDiskSpeedWrite                    = 0
	dvResourceVirtualEnvironmentVMDiskSpeedWriteBurstable           = 0
	dvResourceVirtualEnvironmentVMEFIDiskDatastoreID                = "local-lvm"
	dvResourceVirtualEnvironmentVMEFIDiskFileFormat                 = "qcow2"
	dvResourceVirtualEnvironmentVMEFIDiskType                       = "2m"
	dvResourceVirtualEnvironmentVMEFIDiskPreEnrolledKeys            = false
	dvResourceVirtualEnvironmentVMTPMStateDatastoreID               = "local-lvm"
	dvResourceVirtualEnvironmentVMTPMStateVersion                   = "v2.0"
	dvResourceVirtualEnvironmentVMInitializationDatastoreID         = "local-lvm"
	dvResourceVirtualEnvironmentVMInitializationInterface           = ""
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
	dvResourceVirtualEnvironmentVMMigrate                           = false
	dvResourceVirtualEnvironmentVMName                              = ""
	dvResourceVirtualEnvironmentVMNetworkDeviceBridge               = "vmbr0"
	dvResourceVirtualEnvironmentVMNetworkDeviceEnabled              = true
	dvResourceVirtualEnvironmentVMNetworkDeviceFirewall             = false
	dvResourceVirtualEnvironmentVMNetworkDeviceModel                = "virtio"
	dvResourceVirtualEnvironmentVMNetworkDeviceQueues               = 0
	dvResourceVirtualEnvironmentVMNetworkDeviceRateLimit            = 0
	dvResourceVirtualEnvironmentVMNetworkDeviceVLANID               = 0
	dvResourceVirtualEnvironmentVMNetworkDeviceMTU                  = 0
	dvResourceVirtualEnvironmentVMOperatingSystemType               = "other"
	dvResourceVirtualEnvironmentVMPoolID                            = ""
	dvResourceVirtualEnvironmentVMSerialDeviceDevice                = "socket"
	dvResourceVirtualEnvironmentVMSMBIOSFamily                      = ""
	dvResourceVirtualEnvironmentVMSMBIOSManufacturer                = ""
	dvResourceVirtualEnvironmentVMSMBIOSProduct                     = ""
	dvResourceVirtualEnvironmentVMSMBIOSSKU                         = ""
	dvResourceVirtualEnvironmentVMSMBIOSSerial                      = ""
	dvResourceVirtualEnvironmentVMSMBIOSVersion                     = ""
	dvResourceVirtualEnvironmentVMStarted                           = true
	dvResourceVirtualEnvironmentVMStartupOrder                      = -1
	dvResourceVirtualEnvironmentVMStartupUpDelay                    = -1
	dvResourceVirtualEnvironmentVMStartupDownDelay                  = -1
	dvResourceVirtualEnvironmentVMTabletDevice                      = true
	dvResourceVirtualEnvironmentVMTemplate                          = false
	dvResourceVirtualEnvironmentVMTimeoutClone                      = 1800
	dvResourceVirtualEnvironmentVMTimeoutCreate                     = 1800
	dvResourceVirtualEnvironmentVMTimeoutMoveDisk                   = 1800
	dvResourceVirtualEnvironmentVMTimeoutMigrate                    = 1800
	dvResourceVirtualEnvironmentVMTimeoutReboot                     = 1800
	dvResourceVirtualEnvironmentVMTimeoutShutdownVM                 = 1800
	dvResourceVirtualEnvironmentVMTimeoutStartVM                    = 1800
	dvResourceVirtualEnvironmentVMTimeoutStopVM                     = 300
	dvResourceVirtualEnvironmentVMVGAEnabled                        = true
	dvResourceVirtualEnvironmentVMVGAMemory                         = 16
	dvResourceVirtualEnvironmentVMVGAType                           = "std"
	dvResourceVirtualEnvironmentVMSCSIHardware                      = "virtio-scsi-pci"
	dvResourceVirtualEnvironmentVMStopOnDestroy                     = false
	dvResourceVirtualEnvironmentVMHookScript                        = ""

	maxResourceVirtualEnvironmentVMAudioDevices   = 1
	maxResourceVirtualEnvironmentVMNetworkDevices = 32
	maxResourceVirtualEnvironmentVMSerialDevices  = 4
	maxResourceVirtualEnvironmentVMHostPCIDevices = 8
	maxResourceVirtualEnvironmentVMHostUSBDevices = 4

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
	mkResourceVirtualEnvironmentVMCDROMInterface                    = "interface"
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
	mkResourceVirtualEnvironmentVMCPULimit                          = "limit"
	mkResourceVirtualEnvironmentVMCPUNUMA                           = "numa"
	mkResourceVirtualEnvironmentVMCPUSockets                        = "sockets"
	mkResourceVirtualEnvironmentVMCPUType                           = "type"
	mkResourceVirtualEnvironmentVMCPUUnits                          = "units"
	mkResourceVirtualEnvironmentVMDescription                       = "description"
	mkResourceVirtualEnvironmentVMDisk                              = "disk"
	mkResourceVirtualEnvironmentVMDiskInterface                     = "interface"
	mkResourceVirtualEnvironmentVMDiskDatastoreID                   = "datastore_id"
	mkResourceVirtualEnvironmentVMDiskPathInDatastore               = "path_in_datastore"
	mkResourceVirtualEnvironmentVMDiskFileFormat                    = "file_format"
	mkResourceVirtualEnvironmentVMDiskFileID                        = "file_id"
	mkResourceVirtualEnvironmentVMDiskSize                          = "size"
	mkResourceVirtualEnvironmentVMDiskIOThread                      = "iothread"
	mkResourceVirtualEnvironmentVMDiskSSD                           = "ssd"
	mkResourceVirtualEnvironmentVMDiskDiscard                       = "discard"
	mkResourceVirtualEnvironmentVMDiskCache                         = "cache"
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
	mkResourceVirtualEnvironmentVMTPMState                          = "tpm_state"
	mkResourceVirtualEnvironmentVMTPMStateDatastoreID               = "datastore_id"
	mkResourceVirtualEnvironmentVMTPMStateVersion                   = "version"
	mkResourceVirtualEnvironmentVMHostPCI                           = "hostpci"
	mkResourceVirtualEnvironmentVMHostPCIDevice                     = "device"
	mkResourceVirtualEnvironmentVMHostPCIDeviceID                   = "id"
	mkResourceVirtualEnvironmentVMHostPCIDeviceMapping              = "mapping"
	mkResourceVirtualEnvironmentVMHostPCIDeviceMDev                 = "mdev"
	mkResourceVirtualEnvironmentVMHostPCIDevicePCIE                 = "pcie"
	mkResourceVirtualEnvironmentVMHostPCIDeviceROMBAR               = "rombar"
	mkResourceVirtualEnvironmentVMHostPCIDeviceROMFile              = "rom_file"
	mkResourceVirtualEnvironmentVMHostPCIDeviceXVGA                 = "xvga"
	mkResourceVirtualEnvironmentVMInitialization                    = "initialization"
	mkResourceVirtualEnvironmentVMInitializationDatastoreID         = "datastore_id"
	mkResourceVirtualEnvironmentVMInitializationInterface           = "interface"
	mkResourceVirtualEnvironmentVMInitializationDNS                 = "dns"
	mkResourceVirtualEnvironmentVMInitializationDNSDomain           = "domain"
	mkResourceVirtualEnvironmentVMInitializationDNSServer           = "server"
	mkResourceVirtualEnvironmentVMInitializationDNSServers          = "servers"
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
	mkResourceVirtualEnvironmentVMMigrate                           = "migrate"
	mkResourceVirtualEnvironmentVMName                              = "name"
	mkResourceVirtualEnvironmentVMNetworkDevice                     = "network_device"
	mkResourceVirtualEnvironmentVMNetworkDeviceBridge               = "bridge"
	mkResourceVirtualEnvironmentVMNetworkDeviceEnabled              = "enabled"
	mkResourceVirtualEnvironmentVMNetworkDeviceFirewall             = "firewall"
	mkResourceVirtualEnvironmentVMNetworkDeviceMACAddress           = "mac_address"
	mkResourceVirtualEnvironmentVMNetworkDeviceModel                = "model"
	mkResourceVirtualEnvironmentVMNetworkDeviceQueues               = "queues"
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
	mkResourceVirtualEnvironmentVMSMBIOS                            = "smbios"
	mkResourceVirtualEnvironmentVMSMBIOSFamily                      = "family"
	mkResourceVirtualEnvironmentVMSMBIOSManufacturer                = "manufacturer"
	mkResourceVirtualEnvironmentVMSMBIOSProduct                     = "product"
	mkResourceVirtualEnvironmentVMSMBIOSSKU                         = "sku"
	mkResourceVirtualEnvironmentVMSMBIOSSerial                      = "serial"
	mkResourceVirtualEnvironmentVMSMBIOSUUID                        = "uuid"
	mkResourceVirtualEnvironmentVMSMBIOSVersion                     = "version"
	mkResourceVirtualEnvironmentVMStarted                           = "started"
	mkResourceVirtualEnvironmentVMStartup                           = "startup"
	mkResourceVirtualEnvironmentVMStartupOrder                      = "order"
	mkResourceVirtualEnvironmentVMStartupUpDelay                    = "up_delay"
	mkResourceVirtualEnvironmentVMStartupDownDelay                  = "down_delay"
	mkResourceVirtualEnvironmentVMTabletDevice                      = "tablet_device"
	mkResourceVirtualEnvironmentVMTags                              = "tags"
	mkResourceVirtualEnvironmentVMTemplate                          = "template"
	mkResourceVirtualEnvironmentVMTimeoutClone                      = "timeout_clone"
	mkResourceVirtualEnvironmentVMTimeoutCreate                     = "timeout_create"
	mkResourceVirtualEnvironmentVMTimeoutMoveDisk                   = "timeout_move_disk"
	mkResourceVirtualEnvironmentVMTimeoutMigrate                    = "timeout_migrate"
	mkResourceVirtualEnvironmentVMTimeoutReboot                     = "timeout_reboot"
	mkResourceVirtualEnvironmentVMTimeoutShutdownVM                 = "timeout_shutdown_vm"
	mkResourceVirtualEnvironmentVMTimeoutStartVM                    = "timeout_start_vm"
	mkResourceVirtualEnvironmentVMTimeoutStopVM                     = "timeout_stop_vm"
	mkResourceVirtualEnvironmentVMHostUSB                           = "usb"
	mkResourceVirtualEnvironmentVMHostUSBDevice                     = "host"
	mkResourceVirtualEnvironmentVMHostUSBDeviceMapping              = "mapping"
	mkResourceVirtualEnvironmentVMHostUSBDeviceUSB3                 = "usb3"
	mkResourceVirtualEnvironmentVMVGA                               = "vga"
	mkResourceVirtualEnvironmentVMVGAEnabled                        = "enabled"
	mkResourceVirtualEnvironmentVMVGAMemory                         = "memory"
	mkResourceVirtualEnvironmentVMVGAType                           = "type"
	mkResourceVirtualEnvironmentVMVMID                              = "vm_id"
	mkResourceVirtualEnvironmentVMSCSIHardware                      = "scsi_hardware"
	mkResourceVirtualEnvironmentVMHookScriptFileID                  = "hook_script_file_id"
	mkResourceVirtualEnvironmentVMStopOnDestroy                     = "stop_on_destroy"
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
							ValidateDiagFunc: validator.Timeout(),
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
							ValidateDiagFunc: validator.QEMUAgentType(),
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
				ValidateDiagFunc: validator.BIOS(),
			},
			mkResourceVirtualEnvironmentVMCDROM: {
				Type:        schema.TypeList,
				Description: "The CDROM drive",
				Optional:    true,
				DefaultFunc: func() (interface{}, error) {
					return []interface{}{
						map[string]interface{}{
							mkResourceVirtualEnvironmentVMCDROMEnabled:   dvResourceVirtualEnvironmentVMCDROMEnabled,
							mkResourceVirtualEnvironmentVMCDROMFileID:    dvResourceVirtualEnvironmentVMCDROMFileID,
							mkResourceVirtualEnvironmentVMCDROMInterface: dvResourceVirtualEnvironmentVMCDROMInterface,
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
							ValidateDiagFunc: validator.FileID(),
						},
						mkResourceVirtualEnvironmentVMCDROMInterface: {
							Type:             schema.TypeString,
							Description:      "The CDROM interface",
							Optional:         true,
							Default:          dvResourceVirtualEnvironmentVMCDROMInterface,
							ValidateDiagFunc: validator.IDEInterface(),
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
							mkResourceVirtualEnvironmentVMCPUHotplugged:   dvResourceVirtualEnvironmentVMCPUHotplugged,
							mkResourceVirtualEnvironmentVMCPULimit:        dvResourceVirtualEnvironmentVMCPULimit,
							mkResourceVirtualEnvironmentVMCPUNUMA:         dvResourceVirtualEnvironmentVMCPUNUMA,
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
						mkResourceVirtualEnvironmentVMCPULimit: {
							Type:        schema.TypeInt,
							Description: "Limit of CPU usage",
							Optional:    true,
							Default:     dvResourceVirtualEnvironmentVMCPULimit,
							ValidateDiagFunc: validation.ToDiagFunc(
								validation.IntBetween(0, 128),
							),
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
							ValidateDiagFunc: validator.CPUType(),
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
							mkResourceVirtualEnvironmentVMDiskDatastoreID:     dvResourceVirtualEnvironmentVMDiskDatastoreID,
							mkResourceVirtualEnvironmentVMDiskPathInDatastore: nil,
							mkResourceVirtualEnvironmentVMDiskFileID:          dvResourceVirtualEnvironmentVMDiskFileID,
							mkResourceVirtualEnvironmentVMDiskInterface:       dvResourceVirtualEnvironmentVMDiskInterface,
							mkResourceVirtualEnvironmentVMDiskSize:            dvResourceVirtualEnvironmentVMDiskSize,
							mkResourceVirtualEnvironmentVMDiskIOThread:        dvResourceVirtualEnvironmentVMDiskIOThread,
							mkResourceVirtualEnvironmentVMDiskSSD:             dvResourceVirtualEnvironmentVMDiskSSD,
							mkResourceVirtualEnvironmentVMDiskDiscard:         dvResourceVirtualEnvironmentVMDiskDiscard,
							mkResourceVirtualEnvironmentVMDiskCache:           dvResourceVirtualEnvironmentVMDiskCache,
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
						mkResourceVirtualEnvironmentVMDiskPathInDatastore: {
							Type:        schema.TypeString,
							Description: "The in-datastore path to disk image",
							Computed:    true,
							Optional:    true,
							Default:     nil,
						},
						mkResourceVirtualEnvironmentVMDiskFileFormat: {
							Type:             schema.TypeString,
							Description:      "The file format",
							Optional:         true,
							ForceNew:         true,
							Computed:         true,
							ValidateDiagFunc: validator.FileFormat(),
						},
						mkResourceVirtualEnvironmentVMDiskFileID: {
							Type:             schema.TypeString,
							Description:      "The file id for a disk image",
							Optional:         true,
							ForceNew:         true,
							Default:          dvResourceVirtualEnvironmentVMDiskFileID,
							ValidateDiagFunc: validator.FileID(),
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
						mkResourceVirtualEnvironmentVMDiskCache: {
							Type:        schema.TypeString,
							Description: "The drive’s cache mode",
							Optional:    true,
							Default:     dvResourceVirtualEnvironmentVMDiskCache,
							ValidateDiagFunc: validation.ToDiagFunc(
								validation.StringInSlice([]string{
									"none",
									"writethrough",
									"writeback",
									"unsafe",
									"directsync",
								}, false),
							),
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
					return []interface{}{}, nil
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
							ValidateDiagFunc: validator.FileFormat(),
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
			mkResourceVirtualEnvironmentVMTPMState: {
				Type:        schema.TypeList,
				Description: "The tpmstate device",
				Optional:    true,
				ForceNew:    true,
				DefaultFunc: func() (interface{}, error) {
					return []interface{}{}, nil
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mkResourceVirtualEnvironmentVMTPMStateDatastoreID: {
							Type:        schema.TypeString,
							Description: "Datastore ID",
							Optional:    true,
							Default:     dvResourceVirtualEnvironmentVMTPMStateDatastoreID,
						},
						mkResourceVirtualEnvironmentVMTPMStateVersion: {
							Type:        schema.TypeString,
							Description: "TPM version",
							Optional:    true,
							ForceNew:    true,
							Default:     dvResourceVirtualEnvironmentVMTPMStateVersion,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
								"v1.2",
								"v2.0",
							}, true)),
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
							Default:     dvResourceVirtualEnvironmentVMInitializationDatastoreID,
						},
						mkResourceVirtualEnvironmentVMInitializationInterface: {
							Type:             schema.TypeString,
							Description:      "The IDE interface on which the CloudInit drive will be added",
							Optional:         true,
							Default:          dvResourceVirtualEnvironmentVMInitializationInterface,
							ValidateDiagFunc: validator.CloudInitInterface(),
							DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
								return newValue == ""
							},
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
										Deprecated: "The `server` attribute is deprecated and will be removed in a future release. " +
											"Please use the `servers` attribute instead.",
										Optional: true,
										Default:  dvResourceVirtualEnvironmentVMInitializationDNSServer,
									},
									mkResourceVirtualEnvironmentVMInitializationDNSServers: {
										Type:        schema.TypeList,
										Description: "The list of DNS servers",
										Optional:    true,
										Elem:        &schema.Schema{Type: schema.TypeString, ValidateFunc: validation.IsIPAddress},
										MinItems:    0,
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
							ValidateDiagFunc: validator.FileID(),
						},
						mkResourceVirtualEnvironmentVMInitializationVendorDataFileID: {
							Type:             schema.TypeString,
							Description:      "The ID of a file containing vendor data",
							Optional:         true,
							ForceNew:         true,
							Default:          dvResourceVirtualEnvironmentVMInitializationVendorDataFileID,
							ValidateDiagFunc: validator.FileID(),
						},
						mkResourceVirtualEnvironmentVMInitializationNetworkDataFileID: {
							Type:             schema.TypeString,
							Description:      "The ID of a file containing network config",
							Optional:         true,
							ForceNew:         true,
							Default:          dvResourceVirtualEnvironmentVMInitializationNetworkDataFileID,
							ValidateDiagFunc: validator.FileID(),
						},
						mkResourceVirtualEnvironmentVMInitializationMetaDataFileID: {
							Type:             schema.TypeString,
							Description:      "The ID of a file containing meta data config",
							Optional:         true,
							ForceNew:         true,
							Default:          dvResourceVirtualEnvironmentVMInitializationMetaDataFileID,
							ValidateDiagFunc: validator.FileID(),
						},
						mkResourceVirtualEnvironmentVMInitializationType: {
							Type:             schema.TypeString,
							Description:      "The cloud-init configuration format",
							Optional:         true,
							ForceNew:         true,
							Default:          dvResourceVirtualEnvironmentVMInitializationType,
							ValidateDiagFunc: validator.CloudInitType(),
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
				ForceNew:    false,
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
							Type: schema.TypeString,
							Description: "The PCI ID of the device, for example 0000:00:1f.0 (or 0000:00:1f.0;0000:00:1f.1 for multiple " +
								"device functions, or 0000:00:1f for all functions). Use either this or mapping.",
							Optional: true,
						},
						mkResourceVirtualEnvironmentVMHostPCIDeviceMapping: {
							Type:        schema.TypeString,
							Description: "The resource mapping name of the device, for example gpu. Use either this or id.",
							Optional:    true,
						},
						mkResourceVirtualEnvironmentVMHostPCIDeviceMDev: {
							Type:        schema.TypeString,
							Description: "The the mediated device to use",
							Optional:    true,
						},
						mkResourceVirtualEnvironmentVMHostPCIDevicePCIE: {
							Type: schema.TypeBool,
							Description: "Tells Proxmox VE to use a PCIe or PCI port. Some guests/device combination require PCIe rather " +
								"than PCI. PCIe is only available for q35 machine types.",
							Optional: true,
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
			mkResourceVirtualEnvironmentVMHostUSB: {
				Type:        schema.TypeList,
				Description: "The Host USB devices mapped to the VM",
				Optional:    true,
				ForceNew:    false,
				DefaultFunc: func() (interface{}, error) {
					return []interface{}{}, nil
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mkResourceVirtualEnvironmentVMHostUSBDevice: {
							Type:        schema.TypeString,
							Description: "The USB device ID for Proxmox, in form of '<MANUFACTURER>:<ID>'",
							Required:    true,
						},
						mkResourceVirtualEnvironmentVMHostUSBDeviceMapping: {
							Type:        schema.TypeString,
							Description: "The resource mapping name of the device, for example usbdisk. Use either this or id.",
							Optional:    true,
						},
						mkResourceVirtualEnvironmentVMHostUSBDeviceUSB3: {
							Type:        schema.TypeBool,
							Description: "Makes the USB device a USB3 device for the machine. Default is false",
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
				ValidateDiagFunc: validator.KeyboardLayout(),
			},
			mkResourceVirtualEnvironmentVMMachine: {
				Type:             schema.TypeString,
				Description:      "The VM machine type, either default `pc` or `q35`",
				Optional:         true,
				Default:          dvResourceVirtualEnvironmentVMMachineType,
				ValidateDiagFunc: validator.MachineType(),
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
							ValidateDiagFunc: validator.MACAddress(),
						},
						mkResourceVirtualEnvironmentVMNetworkDeviceModel: {
							Type:             schema.TypeString,
							Description:      "The model",
							Optional:         true,
							Default:          dvResourceVirtualEnvironmentVMNetworkDeviceModel,
							ValidateDiagFunc: validator.NetworkDeviceModel(),
						},
						mkResourceVirtualEnvironmentVMNetworkDeviceQueues: {
							Type:             schema.TypeInt,
							Description:      "Number of packet queues to be used on the device",
							Optional:         true,
							Default:          dvResourceVirtualEnvironmentVMNetworkDeviceQueues,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(0, 64)),
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
			},
			mkResourceVirtualEnvironmentVMMigrate: {
				Type:        schema.TypeBool,
				Description: "Whether to migrate the VM on node change instead of re-creating it",
				Optional:    true,
				Default:     dvResourceVirtualEnvironmentVMMigrate,
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
			mkResourceVirtualEnvironmentVMSMBIOS: {
				Type:        schema.TypeList,
				Description: "Specifies SMBIOS (type1) settings for the VM",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mkResourceVirtualEnvironmentVMSMBIOSFamily: {
							Type:        schema.TypeString,
							Description: "Sets SMBIOS family string",
							Optional:    true,
							Default:     dvResourceVirtualEnvironmentVMSMBIOSFamily,
						},
						mkResourceVirtualEnvironmentVMSMBIOSManufacturer: {
							Type:        schema.TypeString,
							Description: "Sets SMBIOS manufacturer",
							Optional:    true,
							Default:     dvResourceVirtualEnvironmentVMSMBIOSManufacturer,
						},
						mkResourceVirtualEnvironmentVMSMBIOSProduct: {
							Type:        schema.TypeString,
							Description: "Sets SMBIOS product ID",
							Optional:    true,
							Default:     dvResourceVirtualEnvironmentVMSMBIOSProduct,
						},
						mkResourceVirtualEnvironmentVMSMBIOSSerial: {
							Type:        schema.TypeString,
							Description: "Sets SMBIOS serial number",
							Optional:    true,
							Default:     dvResourceVirtualEnvironmentVMSMBIOSSerial,
						},
						mkResourceVirtualEnvironmentVMSMBIOSSKU: {
							Type:        schema.TypeString,
							Description: "Sets SMBIOS SKU",
							Optional:    true,
							Default:     dvResourceVirtualEnvironmentVMSMBIOSSKU,
						},
						mkResourceVirtualEnvironmentVMSMBIOSUUID: {
							Type:        schema.TypeString,
							Description: "Sets SMBIOS UUID",
							Optional:    true,
							Computed:    true,
						},
						mkResourceVirtualEnvironmentVMSMBIOSVersion: {
							Type:        schema.TypeString,
							Description: "Sets SMBIOS version",
							Optional:    true,
							Default:     dvResourceVirtualEnvironmentVMSMBIOSVersion,
						},
					},
				},
				MaxItems: 1,
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
			mkResourceVirtualEnvironmentVMStartup: {
				Type:        schema.TypeList,
				Description: "Defines startup and shutdown behavior of the VM",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mkResourceVirtualEnvironmentVMStartupOrder: {
							Type:        schema.TypeInt,
							Description: "A non-negative number defining the general startup order",
							Optional:    true,
							Default:     dvResourceVirtualEnvironmentVMStartupOrder,
						},
						mkResourceVirtualEnvironmentVMStartupUpDelay: {
							Type:        schema.TypeInt,
							Description: "A non-negative number defining the delay in seconds before the next VM is started",
							Optional:    true,
							Default:     dvResourceVirtualEnvironmentVMStartupUpDelay,
						},
						mkResourceVirtualEnvironmentVMStartupDownDelay: {
							Type:        schema.TypeInt,
							Description: "A non-negative number defining the delay in seconds before the next VM is shut down",
							Optional:    true,
							Default:     dvResourceVirtualEnvironmentVMStartupDownDelay,
						},
					},
				},
				MaxItems: 1,
				MinItems: 0,
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
				DiffSuppressFunc:      structure.SuppressIfListsAreEqualIgnoringOrder,
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
			mkResourceVirtualEnvironmentVMTimeoutCreate: {
				Type:        schema.TypeInt,
				Description: "Create VM timeout",
				Optional:    true,
				Default:     dvResourceVirtualEnvironmentVMTimeoutCreate,
			},
			mkResourceVirtualEnvironmentVMTimeoutMoveDisk: {
				Type:        schema.TypeInt,
				Description: "MoveDisk timeout",
				Optional:    true,
				Default:     dvResourceVirtualEnvironmentVMTimeoutMoveDisk,
			},
			mkResourceVirtualEnvironmentVMTimeoutMigrate: {
				Type:        schema.TypeInt,
				Description: "Migrate VM timeout",
				Optional:    true,
				Default:     dvResourceVirtualEnvironmentVMTimeoutMigrate,
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
							ValidateDiagFunc: validator.VGAMemory(),
						},
						mkResourceVirtualEnvironmentVMVGAType: {
							Type:             schema.TypeString,
							Description:      "The VGA type",
							Optional:         true,
							Default:          dvResourceVirtualEnvironmentVMVGAType,
							ValidateDiagFunc: validator.VGAType(),
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
				ValidateDiagFunc: validator.SCSIHardware(),
			},
			mkResourceVirtualEnvironmentVMHookScriptFileID: {
				Type:        schema.TypeString,
				Description: "A hook script",
				Optional:    true,
				Default:     dvResourceVirtualEnvironmentVMHookScript,
			},
			mkResourceVirtualEnvironmentVMStopOnDestroy: {
				Type:        schema.TypeBool,
				Description: "Whether to stop rather than shutdown on VM destroy",
				Optional:    true,
				Default:     dvResourceVirtualEnvironmentVMStopOnDestroy,
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
			customdiff.ForceNewIf(
				mkResourceVirtualEnvironmentVMNodeName,
				func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) bool {
					return !d.Get(mkResourceVirtualEnvironmentVMMigrate).(bool)
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
