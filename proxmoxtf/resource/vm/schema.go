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
	dvRebootAfterCreation               = false
	dvOnBoot                            = true
	dvACPI                              = true
	dvAgentEnabled                      = false
	dvAgentTimeout                      = "15m"
	dvAgentTrim                         = false
	dvAgentType                         = "virtio"
	dvAudioDeviceDevice                 = "intel-hda"
	dvAudioDeviceDriver                 = "spice"
	dvAudioDeviceEnabled                = true
	dvBIOS                              = "seabios"
	dvCDROMEnabled                      = false
	dvCDROMFileID                       = ""
	dvCDROMInterface                    = "ide3"
	dvCloneDatastoreID                  = ""
	dvCloneNodeName                     = ""
	dvCloneFull                         = true
	dvCloneRetries                      = 1
	dvCPUArchitecture                   = "x86_64"
	dvCPUCores                          = 1
	dvCPUHotplugged                     = 0
	dvCPULimit                          = 0
	dvCPUNUMA                           = false
	dvCPUSockets                        = 1
	dvCPUType                           = "qemu64"
	dvCPUUnits                          = 1024
	dvDescription                       = ""
	dvDiskInterface                     = "scsi0"
	dvDiskDatastoreID                   = "local-lvm"
	dvDiskFileFormat                    = "qcow2"
	dvDiskFileID                        = ""
	dvDiskSize                          = 8
	dvDiskIOThread                      = false
	dvDiskSSD                           = false
	dvDiskDiscard                       = "ignore"
	dvDiskCache                         = "none"
	dvDiskSpeedRead                     = 0
	dvDiskSpeedReadBurstable            = 0
	dvDiskSpeedWrite                    = 0
	dvDiskSpeedWriteBurstable           = 0
	dvEFIDiskDatastoreID                = "local-lvm"
	dvEFIDiskFileFormat                 = "qcow2"
	dvEFIDiskType                       = "2m"
	dvEFIDiskPreEnrolledKeys            = false
	dvTPMStateDatastoreID               = "local-lvm"
	dvTPMStateVersion                   = "v2.0"
	dvInitializationDatastoreID         = "local-lvm"
	dvInitializationInterface           = ""
	dvInitializationDNSDomain           = ""
	dvInitializationDNSServer           = ""
	dvInitializationIPConfigIPv4Address = ""
	dvInitializationIPConfigIPv4Gateway = ""
	dvInitializationIPConfigIPv6Address = ""
	dvInitializationIPConfigIPv6Gateway = ""
	dvInitializationUserAccountPassword = ""
	dvInitializationUserDataFileID      = ""
	dvInitializationVendorDataFileID    = ""
	dvInitializationNetworkDataFileID   = ""
	dvInitializationMetaDataFileID      = ""
	dvInitializationType                = ""
	dvKeyboardLayout                    = "en-us"
	dvKVMArguments                      = ""
	dvMachineType                       = ""
	dvMemoryDedicated                   = 512
	dvMemoryFloating                    = 0
	dvMemoryShared                      = 0
	dvMigrate                           = false
	dvName                              = ""
	dvNetworkDeviceBridge               = "vmbr0"
	dvNetworkDeviceEnabled              = true
	dvNetworkDeviceFirewall             = false
	dvNetworkDeviceModel                = "virtio"
	dvNetworkDeviceQueues               = 0
	dvNetworkDeviceRateLimit            = 0
	dvNetworkDeviceVLANID               = 0
	dvNetworkDeviceMTU                  = 0
	dvOperatingSystemType               = "other"
	dvPoolID                            = ""
	dvSerialDeviceDevice                = "socket"
	dvSMBIOSFamily                      = ""
	dvSMBIOSManufacturer                = ""
	dvSMBIOSProduct                     = ""
	dvSMBIOSSKU                         = ""
	dvSMBIOSSerial                      = ""
	dvSMBIOSVersion                     = ""
	dvStarted                           = true
	dvStartupOrder                      = -1
	dvStartupUpDelay                    = -1
	dvStartupDownDelay                  = -1
	dvTabletDevice                      = true
	dvTemplate                          = false
	dvTimeoutClone                      = 1800
	dvTimeoutCreate                     = 1800
	dvTimeoutMoveDisk                   = 1800
	dvTimeoutMigrate                    = 1800
	dvTimeoutReboot                     = 1800
	dvTimeoutShutdownVM                 = 1800
	dvTimeoutStartVM                    = 1800
	dvTimeoutStopVM                     = 300
	dvVGAEnabled                        = true
	dvVGAMemory                         = 16
	dvVGAType                           = "std"
	dvSCSIHardware                      = "virtio-scsi-pci"
	dvStopOnDestroy                     = false
	dvHookScript                        = ""

	maxResourceVirtualEnvironmentVMAudioDevices   = 1
	maxResourceVirtualEnvironmentVMNetworkDevices = 32
	maxResourceVirtualEnvironmentVMSerialDevices  = 4
	maxResourceVirtualEnvironmentVMHostPCIDevices = 8
	maxResourceVirtualEnvironmentVMHostUSBDevices = 4

	mkRebootAfterCreation               = "reboot"
	mkOnBoot                            = "on_boot"
	mkBootOrder                         = "boot_order"
	mkACPI                              = "acpi"
	mkAgent                             = "agent"
	mkAgentEnabled                      = "enabled"
	mkAgentTimeout                      = "timeout"
	mkAgentTrim                         = "trim"
	mkAgentType                         = "type"
	mkAudioDevice                       = "audio_device"
	mkAudioDeviceDevice                 = "device"
	mkAudioDeviceDriver                 = "driver"
	mkAudioDeviceEnabled                = "enabled"
	mkBIOS                              = "bios"
	mkCDROM                             = "cdrom"
	mkCDROMEnabled                      = "enabled"
	mkCDROMFileID                       = "file_id"
	mkCDROMInterface                    = "interface"
	mkClone                             = "clone"
	mkCloneRetries                      = "retries"
	mkCloneDatastoreID                  = "datastore_id"
	mkCloneNodeName                     = "node_name"
	mkCloneVMID                         = "vm_id"
	mkCloneFull                         = "full"
	mkCPU                               = "cpu"
	mkCPUArchitecture                   = "architecture"
	mkCPUCores                          = "cores"
	mkCPUFlags                          = "flags"
	mkCPUHotplugged                     = "hotplugged"
	mkCPULimit                          = "limit"
	mkCPUNUMA                           = "numa"
	mkCPUSockets                        = "sockets"
	mkCPUType                           = "type"
	mkCPUUnits                          = "units"
	mkDescription                       = "description"
	mkDisk                              = "disk"
	mkDiskInterface                     = "interface"
	mkDiskDatastoreID                   = "datastore_id"
	mkDiskPathInDatastore               = "path_in_datastore"
	mkDiskFileFormat                    = "file_format"
	mkDiskFileID                        = "file_id"
	mkDiskSize                          = "size"
	mkDiskIOThread                      = "iothread"
	mkDiskSSD                           = "ssd"
	mkDiskDiscard                       = "discard"
	mkDiskCache                         = "cache"
	mkDiskSpeed                         = "speed"
	mkDiskSpeedRead                     = "read"
	mkDiskSpeedReadBurstable            = "read_burstable"
	mkDiskSpeedWrite                    = "write"
	mkDiskSpeedWriteBurstable           = "write_burstable"
	mkEFIDisk                           = "efi_disk"
	mkEFIDiskDatastoreID                = "datastore_id"
	mkEFIDiskFileFormat                 = "file_format"
	mkEFIDiskType                       = "type"
	mkEFIDiskPreEnrolledKeys            = "pre_enrolled_keys"
	mkTPMState                          = "tpm_state"
	mkTPMStateDatastoreID               = "datastore_id"
	mkTPMStateVersion                   = "version"
	mkHostPCI                           = "hostpci"
	mkHostPCIDevice                     = "device"
	mkHostPCIDeviceID                   = "id"
	mkHostPCIDeviceMapping              = "mapping"
	mkHostPCIDeviceMDev                 = "mdev"
	mkHostPCIDevicePCIE                 = "pcie"
	mkHostPCIDeviceROMBAR               = "rombar"
	mkHostPCIDeviceROMFile              = "rom_file"
	mkHostPCIDeviceXVGA                 = "xvga"
	mkInitialization                    = "initialization"
	mkInitializationDatastoreID         = "datastore_id"
	mkInitializationInterface           = "interface"
	mkInitializationDNS                 = "dns"
	mkInitializationDNSDomain           = "domain"
	mkInitializationDNSServer           = "server"
	mkInitializationDNSServers          = "servers"
	mkInitializationIPConfig            = "ip_config"
	mkInitializationIPConfigIPv4        = "ipv4"
	mkInitializationIPConfigIPv4Address = "address"
	mkInitializationIPConfigIPv4Gateway = "gateway"
	mkInitializationIPConfigIPv6        = "ipv6"
	mkInitializationIPConfigIPv6Address = "address"
	mkInitializationIPConfigIPv6Gateway = "gateway"
	mkInitializationType                = "type"
	mkInitializationUserAccount         = "user_account"
	mkInitializationUserAccountKeys     = "keys"
	mkInitializationUserAccountPassword = "password"
	mkInitializationUserAccountUsername = "username"
	mkInitializationUserDataFileID      = "user_data_file_id"
	mkInitializationVendorDataFileID    = "vendor_data_file_id"
	mkInitializationNetworkDataFileID   = "network_data_file_id"
	mkInitializationMetaDataFileID      = "meta_data_file_id"
	mkIPv4Addresses                     = "ipv4_addresses"
	mkIPv6Addresses                     = "ipv6_addresses"
	mkKeyboardLayout                    = "keyboard_layout"
	mkKVMArguments                      = "kvm_arguments"
	mkMachine                           = "machine"
	mkMACAddresses                      = "mac_addresses"
	mkMemory                            = "memory"
	mkMemoryDedicated                   = "dedicated"
	mkMemoryFloating                    = "floating"
	mkMemoryShared                      = "shared"
	mkMigrate                           = "migrate"
	mkName                              = "name"
	mkNetworkDevice                     = "network_device"
	mkNetworkDeviceBridge               = "bridge"
	mkNetworkDeviceEnabled              = "enabled"
	mkNetworkDeviceFirewall             = "firewall"
	mkNetworkDeviceMACAddress           = "mac_address"
	mkNetworkDeviceModel                = "model"
	mkNetworkDeviceQueues               = "queues"
	mkNetworkDeviceRateLimit            = "rate_limit"
	mkNetworkDeviceVLANID               = "vlan_id"
	mkNetworkDeviceMTU                  = "mtu"
	mkNetworkInterfaceNames             = "network_interface_names"
	mkNodeName                          = "node_name"
	mkOperatingSystem                   = "operating_system"
	mkOperatingSystemType               = "type"
	mkPoolID                            = "pool_id"
	mkSerialDevice                      = "serial_device"
	mkSerialDeviceDevice                = "device"
	mkSMBIOS                            = "smbios"
	mkSMBIOSFamily                      = "family"
	mkSMBIOSManufacturer                = "manufacturer"
	mkSMBIOSProduct                     = "product"
	mkSMBIOSSKU                         = "sku"
	mkSMBIOSSerial                      = "serial"
	mkSMBIOSUUID                        = "uuid"
	mkSMBIOSVersion                     = "version"
	mkStarted                           = "started"
	mkStartup                           = "startup"
	mkStartupOrder                      = "order"
	mkStartupUpDelay                    = "up_delay"
	mkStartupDownDelay                  = "down_delay"
	mkTabletDevice                      = "tablet_device"
	mkTags                              = "tags"
	mkTemplate                          = "template"
	mkTimeoutClone                      = "timeout_clone"
	mkTimeoutCreate                     = "timeout_create"
	mkTimeoutMoveDisk                   = "timeout_move_disk"
	mkTimeoutMigrate                    = "timeout_migrate"
	mkTimeoutReboot                     = "timeout_reboot"
	mkTimeoutShutdownVM                 = "timeout_shutdown_vm"
	mkTimeoutStartVM                    = "timeout_start_vm"
	mkTimeoutStopVM                     = "timeout_stop_vm"
	mkHostUSB                           = "usb"
	mkHostUSBDevice                     = "host"
	mkHostUSBDeviceMapping              = "mapping"
	mkHostUSBDeviceUSB3                 = "usb3"
	mkVGA                               = "vga"
	mkVGAEnabled                        = "enabled"
	mkVGAMemory                         = "memory"
	mkVGAType                           = "type"
	mkVMID                              = "vm_id"
	mkSCSIHardware                      = "scsi_hardware"
	mkHookScriptFileID                  = "hook_script_file_id"
	mkStopOnDestroy                     = "stop_on_destroy"
)

// VM returns a resource that manages VMs.
func VM() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			mkRebootAfterCreation: {
				Type:        schema.TypeBool,
				Description: "Whether to reboot vm after creation",
				Optional:    true,
				Default:     dvRebootAfterCreation,
			},
			mkOnBoot: {
				Type:        schema.TypeBool,
				Description: "Start VM on Node boot",
				Optional:    true,
				Default:     dvOnBoot,
			},
			mkBootOrder: {
				Type:        schema.TypeList,
				Description: "The guest will attempt to boot from devices in the order they appear here",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				DefaultFunc: func() (interface{}, error) {
					return []interface{}{}, nil
				},
			},
			mkACPI: {
				Type:        schema.TypeBool,
				Description: "Whether to enable ACPI",
				Optional:    true,
				Default:     dvACPI,
			},
			mkAgent: {
				Type:        schema.TypeList,
				Description: "The QEMU agent configuration",
				Optional:    true,
				DefaultFunc: func() (interface{}, error) {
					return []interface{}{
						map[string]interface{}{
							mkAgentEnabled: dvAgentEnabled,
							mkAgentTimeout: dvAgentTimeout,
							mkAgentTrim:    dvAgentTrim,
							mkAgentType:    dvAgentType,
						},
					}, nil
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mkAgentEnabled: {
							Type:        schema.TypeBool,
							Description: "Whether to enable the QEMU agent",
							Optional:    true,
							Default:     dvAgentEnabled,
						},
						mkAgentTimeout: {
							Type:             schema.TypeString,
							Description:      "The maximum amount of time to wait for data from the QEMU agent to become available",
							Optional:         true,
							Default:          dvAgentTimeout,
							ValidateDiagFunc: validator.Timeout(),
						},
						mkAgentTrim: {
							Type:        schema.TypeBool,
							Description: "Whether to enable the FSTRIM feature in the QEMU agent",
							Optional:    true,
							Default:     dvAgentTrim,
						},
						mkAgentType: {
							Type:             schema.TypeString,
							Description:      "The QEMU agent interface type",
							Optional:         true,
							Default:          dvAgentType,
							ValidateDiagFunc: validator.QEMUAgentType(),
						},
					},
				},
				MaxItems: 1,
				MinItems: 0,
			},
			mkKVMArguments: {
				Type:        schema.TypeString,
				Description: "The args implementation",
				Optional:    true,
				Default:     dvKVMArguments,
			},
			mkAudioDevice: {
				Type:        schema.TypeList,
				Description: "The audio devices",
				Optional:    true,
				DefaultFunc: func() (interface{}, error) {
					return []interface{}{}, nil
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mkAudioDeviceDevice: {
							Type:             schema.TypeString,
							Description:      "The device",
							Optional:         true,
							Default:          dvAudioDeviceDevice,
							ValidateDiagFunc: vmGetAudioDeviceValidator(),
						},
						mkAudioDeviceDriver: {
							Type:             schema.TypeString,
							Description:      "The driver",
							Optional:         true,
							Default:          dvAudioDeviceDriver,
							ValidateDiagFunc: vmGetAudioDriverValidator(),
						},
						mkAudioDeviceEnabled: {
							Type:        schema.TypeBool,
							Description: "Whether to enable the audio device",
							Optional:    true,
							Default:     dvAudioDeviceEnabled,
						},
					},
				},
				MaxItems: maxResourceVirtualEnvironmentVMAudioDevices,
				MinItems: 0,
			},
			mkBIOS: {
				Type:             schema.TypeString,
				Description:      "The BIOS implementation",
				Optional:         true,
				Default:          dvBIOS,
				ValidateDiagFunc: validator.BIOS(),
			},
			mkCDROM: {
				Type:        schema.TypeList,
				Description: "The CDROM drive",
				Optional:    true,
				DefaultFunc: func() (interface{}, error) {
					return []interface{}{
						map[string]interface{}{
							mkCDROMEnabled:   dvCDROMEnabled,
							mkCDROMFileID:    dvCDROMFileID,
							mkCDROMInterface: dvCDROMInterface,
						},
					}, nil
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mkCDROMEnabled: {
							Type:        schema.TypeBool,
							Description: "Whether to enable the CDROM drive",
							Optional:    true,
							Default:     dvCDROMEnabled,
						},
						mkCDROMFileID: {
							Type:             schema.TypeString,
							Description:      "The file id",
							Optional:         true,
							Default:          dvCDROMFileID,
							ValidateDiagFunc: validator.FileID(),
						},
						mkCDROMInterface: {
							Type:             schema.TypeString,
							Description:      "The CDROM interface",
							Optional:         true,
							Default:          dvCDROMInterface,
							ValidateDiagFunc: validator.IDEInterface(),
						},
					},
				},
				MaxItems: 1,
				MinItems: 0,
			},
			mkClone: {
				Type:        schema.TypeList,
				Description: "The cloning configuration",
				Optional:    true,
				DefaultFunc: func() (interface{}, error) {
					return []interface{}{}, nil
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mkCloneRetries: {
							Type:        schema.TypeInt,
							Description: "The number of Retries to create a clone",
							Optional:    true,
							ForceNew:    true,
							Default:     dvCloneRetries,
						},
						mkCloneDatastoreID: {
							Type:        schema.TypeString,
							Description: "The ID of the target datastore",
							Optional:    true,
							ForceNew:    true,
							Default:     dvCloneDatastoreID,
						},
						mkCloneNodeName: {
							Type:        schema.TypeString,
							Description: "The name of the source node",
							Optional:    true,
							ForceNew:    true,
							Default:     dvCloneNodeName,
						},
						mkCloneVMID: {
							Type:             schema.TypeInt,
							Description:      "The ID of the source VM",
							Required:         true,
							ForceNew:         true,
							ValidateDiagFunc: validator.VMID(),
						},
						mkCloneFull: {
							Type:        schema.TypeBool,
							Description: "The Clone Type, create a Full Clone (true) or a linked Clone (false)",
							Optional:    true,
							ForceNew:    true,
							Default:     dvCloneFull,
						},
					},
				},
				MaxItems: 1,
				MinItems: 0,
			},
			mkCPU: {
				Type:        schema.TypeList,
				Description: "The CPU allocation",
				Optional:    true,
				DefaultFunc: func() (interface{}, error) {
					return []interface{}{
						map[string]interface{}{
							mkCPUArchitecture: dvCPUArchitecture,
							mkCPUCores:        dvCPUCores,
							mkCPUFlags:        []interface{}{},
							mkCPUHotplugged:   dvCPUHotplugged,
							mkCPULimit:        dvCPULimit,
							mkCPUNUMA:         dvCPUNUMA,
							mkCPUSockets:      dvCPUSockets,
							mkCPUType:         dvCPUType,
							mkCPUUnits:        dvCPUUnits,
						},
					}, nil
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mkCPUArchitecture: {
							Type:             schema.TypeString,
							Description:      "The CPU architecture",
							Optional:         true,
							Default:          dvCPUArchitecture,
							ValidateDiagFunc: vmGetCPUArchitectureValidator(),
						},
						mkCPUCores: {
							Type:             schema.TypeInt,
							Description:      "The number of CPU cores",
							Optional:         true,
							Default:          dvCPUCores,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(1, 2304)),
						},
						mkCPUFlags: {
							Type:        schema.TypeList,
							Description: "The CPU flags",
							Optional:    true,
							DefaultFunc: func() (interface{}, error) {
								return []interface{}{}, nil
							},
							Elem: &schema.Schema{Type: schema.TypeString},
						},
						mkCPUHotplugged: {
							Type:             schema.TypeInt,
							Description:      "The number of hotplugged vCPUs",
							Optional:         true,
							Default:          dvCPUHotplugged,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(0, 2304)),
						},
						mkCPULimit: {
							Type:        schema.TypeInt,
							Description: "Limit of CPU usage",
							Optional:    true,
							Default:     dvCPULimit,
							ValidateDiagFunc: validation.ToDiagFunc(
								validation.IntBetween(0, 128),
							),
						},
						mkCPUNUMA: {
							Type:        schema.TypeBool,
							Description: "Enable/disable NUMA.",
							Optional:    true,
							Default:     dvCPUNUMA,
						},
						mkCPUSockets: {
							Type:             schema.TypeInt,
							Description:      "The number of CPU sockets",
							Optional:         true,
							Default:          dvCPUSockets,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(1, 16)),
						},
						mkCPUType: {
							Type:             schema.TypeString,
							Description:      "The emulated CPU type",
							Optional:         true,
							Default:          dvCPUType,
							ValidateDiagFunc: validator.CPUType(),
						},
						mkCPUUnits: {
							Type:        schema.TypeInt,
							Description: "The CPU units",
							Optional:    true,
							Default:     dvCPUUnits,
							ValidateDiagFunc: validation.ToDiagFunc(
								validation.IntBetween(2, 262144),
							),
						},
					},
				},
				MaxItems: 1,
				MinItems: 0,
			},
			mkDescription: {
				Type:        schema.TypeString,
				Description: "The description",
				Optional:    true,
				Default:     dvDescription,
			},
			mkDisk: diskSchema(),
			mkEFIDisk: {
				Type:        schema.TypeList,
				Description: "The efidisk device",
				Optional:    true,
				ForceNew:    true,
				DefaultFunc: func() (interface{}, error) {
					return []interface{}{}, nil
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mkEFIDiskDatastoreID: {
							Type:        schema.TypeString,
							Description: "The datastore id",
							Optional:    true,
							Default:     dvEFIDiskDatastoreID,
						},
						mkEFIDiskFileFormat: {
							Type:             schema.TypeString,
							Description:      "The file format",
							Optional:         true,
							ForceNew:         true,
							Computed:         true,
							ValidateDiagFunc: validator.FileFormat(),
						},
						mkEFIDiskType: {
							Type:        schema.TypeString,
							Description: "Size and type of the OVMF EFI disk",
							Optional:    true,
							ForceNew:    true,
							Default:     dvEFIDiskType,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
								"2m",
								"4m",
							}, true)),
						},
						mkEFIDiskPreEnrolledKeys: {
							Type: schema.TypeBool,
							Description: "Use an EFI vars template with distribution-specific and Microsoft Standard " +
								"keys enrolled, if used with efi type=`4m`.",
							Optional: true,
							ForceNew: true,
							Default:  dvEFIDiskPreEnrolledKeys,
						},
					},
				},
				MaxItems: 1,
				MinItems: 0,
			},
			mkTPMState: {
				Type:        schema.TypeList,
				Description: "The tpmstate device",
				Optional:    true,
				ForceNew:    true,
				DefaultFunc: func() (interface{}, error) {
					return []interface{}{}, nil
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mkTPMStateDatastoreID: {
							Type:        schema.TypeString,
							Description: "Datastore ID",
							Optional:    true,
							Default:     dvTPMStateDatastoreID,
						},
						mkTPMStateVersion: {
							Type:        schema.TypeString,
							Description: "TPM version",
							Optional:    true,
							ForceNew:    true,
							Default:     dvTPMStateVersion,
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
			mkInitialization: {
				Type:        schema.TypeList,
				Description: "The cloud-init configuration",
				Optional:    true,
				DefaultFunc: func() (interface{}, error) {
					return []interface{}{}, nil
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mkInitializationDatastoreID: {
							Type:        schema.TypeString,
							Description: "The datastore id",
							Optional:    true,
							Default:     dvInitializationDatastoreID,
						},
						mkInitializationInterface: {
							Type:             schema.TypeString,
							Description:      "The IDE interface on which the CloudInit drive will be added",
							Optional:         true,
							Default:          dvInitializationInterface,
							ValidateDiagFunc: validator.CloudInitInterface(),
							DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
								return newValue == ""
							},
						},
						mkInitializationDNS: {
							Type:        schema.TypeList,
							Description: "The DNS configuration",
							Optional:    true,
							DefaultFunc: func() (interface{}, error) {
								return []interface{}{}, nil
							},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									mkInitializationDNSDomain: {
										Type:        schema.TypeString,
										Description: "The DNS search domain",
										Optional:    true,
										Default:     dvInitializationDNSDomain,
									},
									mkInitializationDNSServer: {
										Type:        schema.TypeString,
										Description: "The DNS server",
										Deprecated: "The `server` attribute is deprecated and will be removed in a future release. " +
											"Please use the `servers` attribute instead.",
										Optional: true,
										Default:  dvInitializationDNSServer,
									},
									mkInitializationDNSServers: {
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
						mkInitializationIPConfig: {
							Type:        schema.TypeList,
							Description: "The IP configuration",
							Optional:    true,
							DefaultFunc: func() (interface{}, error) {
								return []interface{}{}, nil
							},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									mkInitializationIPConfigIPv4: {
										Type:        schema.TypeList,
										Description: "The IPv4 configuration",
										Optional:    true,
										DefaultFunc: func() (interface{}, error) {
											return []interface{}{}, nil
										},
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												mkInitializationIPConfigIPv4Address: {
													Type:        schema.TypeString,
													Description: "The IPv4 address",
													Optional:    true,
													Default:     dvInitializationIPConfigIPv4Address,
												},
												mkInitializationIPConfigIPv4Gateway: {
													Type:        schema.TypeString,
													Description: "The IPv4 gateway",
													Optional:    true,
													Default:     dvInitializationIPConfigIPv4Gateway,
												},
											},
										},
										MaxItems: 1,
										MinItems: 0,
									},
									mkInitializationIPConfigIPv6: {
										Type:        schema.TypeList,
										Description: "The IPv6 configuration",
										Optional:    true,
										DefaultFunc: func() (interface{}, error) {
											return []interface{}{}, nil
										},
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												mkInitializationIPConfigIPv6Address: {
													Type:        schema.TypeString,
													Description: "The IPv6 address",
													Optional:    true,
													Default:     dvInitializationIPConfigIPv6Address,
												},
												mkInitializationIPConfigIPv6Gateway: {
													Type:        schema.TypeString,
													Description: "The IPv6 gateway",
													Optional:    true,
													Default:     dvInitializationIPConfigIPv6Gateway,
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
						mkInitializationUserAccount: {
							Type:        schema.TypeList,
							Description: "The user account configuration",
							Optional:    true,
							ForceNew:    true,
							DefaultFunc: func() (interface{}, error) {
								return []interface{}{}, nil
							},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									mkInitializationUserAccountKeys: {
										Type:        schema.TypeList,
										Description: "The SSH keys",
										Optional:    true,
										ForceNew:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									mkInitializationUserAccountPassword: {
										Type:        schema.TypeString,
										Description: "The SSH password",
										Optional:    true,
										ForceNew:    true,
										Sensitive:   true,
										Default:     dvInitializationUserAccountPassword,
										DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
											return len(oldValue) > 0 &&
												strings.ReplaceAll(oldValue, "*", "") == ""
										},
									},
									mkInitializationUserAccountUsername: {
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
						mkInitializationUserDataFileID: {
							Type:             schema.TypeString,
							Description:      "The ID of a file containing custom user data",
							Optional:         true,
							ForceNew:         true,
							Default:          dvInitializationUserDataFileID,
							ValidateDiagFunc: validator.FileID(),
						},
						mkInitializationVendorDataFileID: {
							Type:             schema.TypeString,
							Description:      "The ID of a file containing vendor data",
							Optional:         true,
							ForceNew:         true,
							Default:          dvInitializationVendorDataFileID,
							ValidateDiagFunc: validator.FileID(),
						},
						mkInitializationNetworkDataFileID: {
							Type:             schema.TypeString,
							Description:      "The ID of a file containing network config",
							Optional:         true,
							ForceNew:         true,
							Default:          dvInitializationNetworkDataFileID,
							ValidateDiagFunc: validator.FileID(),
						},
						mkInitializationMetaDataFileID: {
							Type:             schema.TypeString,
							Description:      "The ID of a file containing meta data config",
							Optional:         true,
							ForceNew:         true,
							Default:          dvInitializationMetaDataFileID,
							ValidateDiagFunc: validator.FileID(),
						},
						mkInitializationType: {
							Type:             schema.TypeString,
							Description:      "The cloud-init configuration format",
							Optional:         true,
							ForceNew:         true,
							Default:          dvInitializationType,
							ValidateDiagFunc: validator.CloudInitType(),
						},
					},
				},
				MaxItems: 1,
				MinItems: 0,
			},
			mkIPv4Addresses: {
				Type:        schema.TypeList,
				Description: "The IPv4 addresses published by the QEMU agent",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeList,
					Elem: &schema.Schema{Type: schema.TypeString},
				},
			},
			mkIPv6Addresses: {
				Type:        schema.TypeList,
				Description: "The IPv6 addresses published by the QEMU agent",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeList,
					Elem: &schema.Schema{Type: schema.TypeString},
				},
			},
			mkHostPCI: {
				Type:        schema.TypeList,
				Description: "The Host PCI devices mapped to the VM",
				Optional:    true,
				ForceNew:    false,
				DefaultFunc: func() (interface{}, error) {
					return []interface{}{}, nil
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mkHostPCIDevice: {
							Type:        schema.TypeString,
							Description: "The PCI device name for Proxmox, in form of 'hostpciX' where X is a sequential number from 0 to 3",
							Required:    true,
						},
						mkHostPCIDeviceID: {
							Type: schema.TypeString,
							Description: "The PCI ID of the device, for example 0000:00:1f.0 (or 0000:00:1f.0;0000:00:1f.1 for multiple " +
								"device functions, or 0000:00:1f for all functions). Use either this or mapping.",
							Optional: true,
						},
						mkHostPCIDeviceMapping: {
							Type:        schema.TypeString,
							Description: "The resource mapping name of the device, for example gpu. Use either this or id.",
							Optional:    true,
						},
						mkHostPCIDeviceMDev: {
							Type:        schema.TypeString,
							Description: "The the mediated device to use",
							Optional:    true,
						},
						mkHostPCIDevicePCIE: {
							Type: schema.TypeBool,
							Description: "Tells Proxmox VE to use a PCIe or PCI port. Some guests/device combination require PCIe rather " +
								"than PCI. PCIe is only available for q35 machine types.",
							Optional: true,
						},
						mkHostPCIDeviceROMBAR: {
							Type:        schema.TypeBool,
							Description: "Makes the firmware ROM visible for the guest. Default is true",
							Optional:    true,
						},
						mkHostPCIDeviceROMFile: {
							Type:        schema.TypeString,
							Description: "A path to a ROM file for the device to use. This is a relative path under /usr/share/kvm/",
							Optional:    true,
						},
						mkHostPCIDeviceXVGA: {
							Type: schema.TypeBool,
							Description: "Marks the PCI(e) device as the primary GPU of the VM. " +
								"With this enabled the vga configuration argument will be ignored.",
							Optional: true,
						},
					},
				},
			},
			mkHostUSB: {
				Type:        schema.TypeList,
				Description: "The Host USB devices mapped to the VM",
				Optional:    true,
				ForceNew:    false,
				DefaultFunc: func() (interface{}, error) {
					return []interface{}{}, nil
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mkHostUSBDevice: {
							Type:        schema.TypeString,
							Description: "The USB device ID for Proxmox, in form of '<MANUFACTURER>:<ID>'",
							Required:    true,
						},
						mkHostUSBDeviceMapping: {
							Type:        schema.TypeString,
							Description: "The resource mapping name of the device, for example usbdisk. Use either this or id.",
							Optional:    true,
						},
						mkHostUSBDeviceUSB3: {
							Type:        schema.TypeBool,
							Description: "Makes the USB device a USB3 device for the machine. Default is false",
							Optional:    true,
						},
					},
				},
			},
			mkKeyboardLayout: {
				Type:             schema.TypeString,
				Description:      "The keyboard layout",
				Optional:         true,
				Default:          dvKeyboardLayout,
				ValidateDiagFunc: validator.KeyboardLayout(),
			},
			mkMachine: {
				Type:             schema.TypeString,
				Description:      "The VM machine type, either default `pc` or `q35`",
				Optional:         true,
				Default:          dvMachineType,
				ValidateDiagFunc: validator.MachineType(),
			},
			mkMACAddresses: {
				Type:        schema.TypeList,
				Description: "The MAC addresses for the network interfaces",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			mkMemory: {
				Type:        schema.TypeList,
				Description: "The memory allocation",
				Optional:    true,
				DefaultFunc: func() (interface{}, error) {
					return []interface{}{
						map[string]interface{}{
							mkMemoryDedicated: dvMemoryDedicated,
							mkMemoryFloating:  dvMemoryFloating,
							mkMemoryShared:    dvMemoryShared,
						},
					}, nil
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mkMemoryDedicated: {
							Type:        schema.TypeInt,
							Description: "The dedicated memory in megabytes",
							Optional:    true,
							Default:     dvMemoryDedicated,
							ValidateDiagFunc: validation.ToDiagFunc(
								validation.IntBetween(64, 268435456),
							),
						},
						mkMemoryFloating: {
							Type:        schema.TypeInt,
							Description: "The floating memory in megabytes (balloon)",
							Optional:    true,
							Default:     dvMemoryFloating,
							ValidateDiagFunc: validation.ToDiagFunc(
								validation.IntBetween(0, 268435456),
							),
						},
						mkMemoryShared: {
							Type:        schema.TypeInt,
							Description: "The shared memory in megabytes",
							Optional:    true,
							Default:     dvMemoryShared,
							ValidateDiagFunc: validation.ToDiagFunc(
								validation.IntBetween(0, 268435456),
							),
						},
					},
				},
				MaxItems: 1,
				MinItems: 0,
			},
			mkName: {
				Type:        schema.TypeString,
				Description: "The name",
				Optional:    true,
				Default:     dvName,
			},
			mkNetworkDevice: {
				Type:        schema.TypeList,
				Description: "The network devices",
				Optional:    true,
				DefaultFunc: func() (interface{}, error) {
					return make([]interface{}, 1), nil
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mkNetworkDeviceBridge: {
							Type:        schema.TypeString,
							Description: "The bridge",
							Optional:    true,
							Default:     dvNetworkDeviceBridge,
						},
						mkNetworkDeviceEnabled: {
							Type:        schema.TypeBool,
							Description: "Whether to enable the network device",
							Optional:    true,
							Default:     dvNetworkDeviceEnabled,
						},
						mkNetworkDeviceFirewall: {
							Type:        schema.TypeBool,
							Description: "Whether this interface's firewall rules should be used",
							Optional:    true,
							Default:     dvNetworkDeviceFirewall,
						},
						mkNetworkDeviceMACAddress: {
							Type:             schema.TypeString,
							Description:      "The MAC address",
							Optional:         true,
							Computed:         true,
							ValidateDiagFunc: validator.MACAddress(),
						},
						mkNetworkDeviceModel: {
							Type:             schema.TypeString,
							Description:      "The model",
							Optional:         true,
							Default:          dvNetworkDeviceModel,
							ValidateDiagFunc: validator.NetworkDeviceModel(),
						},
						mkNetworkDeviceQueues: {
							Type:             schema.TypeInt,
							Description:      "Number of packet queues to be used on the device",
							Optional:         true,
							Default:          dvNetworkDeviceQueues,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(0, 64)),
						},
						mkNetworkDeviceRateLimit: {
							Type:        schema.TypeFloat,
							Description: "The rate limit in megabytes per second",
							Optional:    true,
							Default:     dvNetworkDeviceRateLimit,
						},
						mkNetworkDeviceVLANID: {
							Type:        schema.TypeInt,
							Description: "The VLAN identifier",
							Optional:    true,
							Default:     dvNetworkDeviceVLANID,
						},
						mkNetworkDeviceMTU: {
							Type:        schema.TypeInt,
							Description: "Maximum transmission unit (MTU)",
							Optional:    true,
							Default:     dvNetworkDeviceMTU,
						},
					},
				},
				MaxItems: maxResourceVirtualEnvironmentVMNetworkDevices,
				MinItems: 0,
			},
			mkNetworkInterfaceNames: {
				Type:        schema.TypeList,
				Description: "The network interface names published by the QEMU agent",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			mkNodeName: {
				Type:        schema.TypeString,
				Description: "The node name",
				Required:    true,
			},
			mkMigrate: {
				Type:        schema.TypeBool,
				Description: "Whether to migrate the VM on node change instead of re-creating it",
				Optional:    true,
				Default:     dvMigrate,
			},
			mkOperatingSystem: {
				Type:        schema.TypeList,
				Description: "The operating system configuration",
				Optional:    true,
				DefaultFunc: func() (interface{}, error) {
					return []interface{}{
						map[string]interface{}{
							mkOperatingSystemType: dvOperatingSystemType,
						},
					}, nil
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mkOperatingSystemType: {
							Type:             schema.TypeString,
							Description:      "The type",
							Optional:         true,
							Default:          dvOperatingSystemType,
							ValidateDiagFunc: vmGetOperatingSystemTypeValidator(),
						},
					},
				},
				MaxItems: 1,
				MinItems: 0,
			},
			mkPoolID: {
				Type:        schema.TypeString,
				Description: "The ID of the pool to assign the virtual machine to",
				Optional:    true,
				Default:     dvPoolID,
			},
			mkSerialDevice: {
				Type:        schema.TypeList,
				Description: "The serial devices",
				Optional:    true,
				DefaultFunc: func() (interface{}, error) {
					return []interface{}{
						map[string]interface{}{
							mkSerialDeviceDevice: dvSerialDeviceDevice,
						},
					}, nil
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mkSerialDeviceDevice: {
							Type:             schema.TypeString,
							Description:      "The device",
							Optional:         true,
							Default:          dvSerialDeviceDevice,
							ValidateDiagFunc: vmGetSerialDeviceValidator(),
						},
					},
				},
				MaxItems: maxResourceVirtualEnvironmentVMSerialDevices,
				MinItems: 0,
			},
			mkSMBIOS: {
				Type:        schema.TypeList,
				Description: "Specifies SMBIOS (type1) settings for the VM",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mkSMBIOSFamily: {
							Type:        schema.TypeString,
							Description: "Sets SMBIOS family string",
							Optional:    true,
							Default:     dvSMBIOSFamily,
						},
						mkSMBIOSManufacturer: {
							Type:        schema.TypeString,
							Description: "Sets SMBIOS manufacturer",
							Optional:    true,
							Default:     dvSMBIOSManufacturer,
						},
						mkSMBIOSProduct: {
							Type:        schema.TypeString,
							Description: "Sets SMBIOS product ID",
							Optional:    true,
							Default:     dvSMBIOSProduct,
						},
						mkSMBIOSSerial: {
							Type:        schema.TypeString,
							Description: "Sets SMBIOS serial number",
							Optional:    true,
							Default:     dvSMBIOSSerial,
						},
						mkSMBIOSSKU: {
							Type:        schema.TypeString,
							Description: "Sets SMBIOS SKU",
							Optional:    true,
							Default:     dvSMBIOSSKU,
						},
						mkSMBIOSUUID: {
							Type:        schema.TypeString,
							Description: "Sets SMBIOS UUID",
							Optional:    true,
							Computed:    true,
						},
						mkSMBIOSVersion: {
							Type:        schema.TypeString,
							Description: "Sets SMBIOS version",
							Optional:    true,
							Default:     dvSMBIOSVersion,
						},
					},
				},
				MaxItems: 1,
				MinItems: 0,
			},
			mkStarted: {
				Type:        schema.TypeBool,
				Description: "Whether to start the virtual machine",
				Optional:    true,
				Default:     dvStarted,
				DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
					return d.Get(mkTemplate).(bool)
				},
			},
			mkStartup: {
				Type:        schema.TypeList,
				Description: "Defines startup and shutdown behavior of the VM",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mkStartupOrder: {
							Type:        schema.TypeInt,
							Description: "A non-negative number defining the general startup order",
							Optional:    true,
							Default:     dvStartupOrder,
						},
						mkStartupUpDelay: {
							Type:        schema.TypeInt,
							Description: "A non-negative number defining the delay in seconds before the next VM is started",
							Optional:    true,
							Default:     dvStartupUpDelay,
						},
						mkStartupDownDelay: {
							Type:        schema.TypeInt,
							Description: "A non-negative number defining the delay in seconds before the next VM is shut down",
							Optional:    true,
							Default:     dvStartupDownDelay,
						},
					},
				},
				MaxItems: 1,
				MinItems: 0,
			},
			mkTabletDevice: {
				Type:        schema.TypeBool,
				Description: "Whether to enable the USB tablet device",
				Optional:    true,
				Default:     dvTabletDevice,
			},
			mkTags: {
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
			mkTemplate: {
				Type:        schema.TypeBool,
				Description: "Whether to create a template",
				Optional:    true,
				ForceNew:    true,
				Default:     dvTemplate,
			},
			mkTimeoutClone: {
				Type:        schema.TypeInt,
				Description: "Clone VM timeout",
				Optional:    true,
				Default:     dvTimeoutClone,
			},
			mkTimeoutCreate: {
				Type:        schema.TypeInt,
				Description: "Create VM timeout",
				Optional:    true,
				Default:     dvTimeoutCreate,
			},
			mkTimeoutMoveDisk: {
				Type:        schema.TypeInt,
				Description: "MoveDisk timeout",
				Optional:    true,
				Default:     dvTimeoutMoveDisk,
			},
			mkTimeoutMigrate: {
				Type:        schema.TypeInt,
				Description: "Migrate VM timeout",
				Optional:    true,
				Default:     dvTimeoutMigrate,
			},
			mkTimeoutReboot: {
				Type:        schema.TypeInt,
				Description: "Reboot timeout",
				Optional:    true,
				Default:     dvTimeoutReboot,
			},
			mkTimeoutShutdownVM: {
				Type:        schema.TypeInt,
				Description: "Shutdown timeout",
				Optional:    true,
				Default:     dvTimeoutShutdownVM,
			},
			mkTimeoutStartVM: {
				Type:        schema.TypeInt,
				Description: "Start VM timeout",
				Optional:    true,
				Default:     dvTimeoutStartVM,
			},
			mkTimeoutStopVM: {
				Type:        schema.TypeInt,
				Description: "Stop VM timeout",
				Optional:    true,
				Default:     dvTimeoutStopVM,
			},
			mkVGA: {
				Type:        schema.TypeList,
				Description: "The VGA configuration",
				Optional:    true,
				DefaultFunc: func() (interface{}, error) {
					return []interface{}{
						map[string]interface{}{
							mkVGAEnabled: dvVGAEnabled,
							mkVGAMemory:  dvVGAMemory,
							mkVGAType:    dvVGAType,
						},
					}, nil
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mkVGAEnabled: {
							Type:        schema.TypeBool,
							Description: "Whether to enable the VGA device",
							Optional:    true,
							Default:     dvVGAEnabled,
						},
						mkVGAMemory: {
							Type:             schema.TypeInt,
							Description:      "The VGA memory in megabytes (4-512 MB)",
							Optional:         true,
							Default:          dvVGAMemory,
							ValidateDiagFunc: validator.VGAMemory(),
						},
						mkVGAType: {
							Type:             schema.TypeString,
							Description:      "The VGA type",
							Optional:         true,
							Default:          dvVGAType,
							ValidateDiagFunc: validator.VGAType(),
						},
					},
				},
				MaxItems: 1,
				MinItems: 0,
			},
			mkVMID: {
				Type:        schema.TypeInt,
				Description: "The VM identifier",
				Optional:    true,
				Computed:    true,
				// "ForceNew: true" handled in CustomizeDiff, making sure VMs with legacy configs with vm_id = -1
				// do not require re-creation.
				ValidateDiagFunc: validator.VMID(),
			},
			mkSCSIHardware: {
				Type:             schema.TypeString,
				Description:      "The SCSI hardware type",
				Optional:         true,
				Default:          dvSCSIHardware,
				ValidateDiagFunc: validator.SCSIHardware(),
			},
			mkHookScriptFileID: {
				Type:        schema.TypeString,
				Description: "A hook script",
				Optional:    true,
				Default:     dvHookScript,
			},
			mkStopOnDestroy: {
				Type:        schema.TypeBool,
				Description: "Whether to stop rather than shutdown on VM destroy",
				Optional:    true,
				Default:     dvStopOnDestroy,
			},
		},
		CreateContext: vmCreate,
		ReadContext:   vmRead,
		UpdateContext: vmUpdate,
		DeleteContext: vmDelete,
		CustomizeDiff: customdiff.All(
			customdiff.ComputedIf(
				mkIPv4Addresses,
				func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) bool {
					return d.HasChange(mkStarted) ||
						d.HasChange(mkNetworkDevice)
				},
			),
			customdiff.ComputedIf(
				mkIPv6Addresses,
				func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) bool {
					return d.HasChange(mkStarted) ||
						d.HasChange(mkNetworkDevice)
				},
			),
			customdiff.ComputedIf(
				mkNetworkInterfaceNames,
				func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) bool {
					return d.HasChange(mkStarted) ||
						d.HasChange(mkNetworkDevice)
				},
			),
			customdiff.ForceNewIf(
				mkVMID,
				func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) bool {
					newValue := d.Get(mkVMID)

					// 'vm_id' is ForceNew, except when changing 'vm_id' to existing correct id
					// (automatic fix from -1 to actual vm_id must not re-create VM)
					return strconv.Itoa(newValue.(int)) != d.Id()
				},
			),
			customdiff.ForceNewIf(
				mkNodeName,
				func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) bool {
					return !d.Get(mkMigrate).(bool)
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
				err = d.Set(mkNodeName, node)
				if err != nil {
					return nil, fmt.Errorf("failed setting state during import: %w", err)
				}

				return []*schema.ResourceData{d}, nil
			},
		},
	}
}
