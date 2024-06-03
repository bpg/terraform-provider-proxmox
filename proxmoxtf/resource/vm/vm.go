/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package resource

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster"
	"github.com/bpg/terraform-provider-proxmox/proxmox/helpers/ptr"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/vms"
	"github.com/bpg/terraform-provider-proxmox/proxmox/pools"
	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
	"github.com/bpg/terraform-provider-proxmox/proxmoxtf"
	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/resource/validators"
	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/resource/vm/disk"
	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/resource/vm/network"
	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/structure"
	"github.com/bpg/terraform-provider-proxmox/utils"
)

const (
	dvRebootAfterCreation = false
	dvOnBoot              = true
	dvACPI                = true
	dvAgentEnabled        = false
	dvAgentTimeout        = "15m"
	dvAgentTrim           = false
	dvAgentType           = "virtio"
	dvAudioDeviceDevice   = "intel-hda"
	dvAudioDeviceDriver   = "spice"
	dvAudioDeviceEnabled  = true
	dvBIOS                = "seabios"
	dvCDROMEnabled        = false
	dvCDROMFileID         = ""
	dvCDROMInterface      = "ide3"
	dvCloneDatastoreID    = ""
	dvCloneNodeName       = ""
	dvCloneFull           = true
	dvCloneRetries        = 1
	dvCPUArchitecture     = "x86_64"
	dvCPUCores            = 1
	dvCPUHotplugged       = 0
	dvCPULimit            = 0
	dvCPUNUMA             = false
	dvCPUSockets          = 1
	dvCPUType             = "qemu64"
	dvCPUUnits            = 1024
	dvCPUAffinity         = ""
	dvDescription         = ""

	dvEFIDiskDatastoreID                = "local-lvm"
	dvEFIDiskFileFormat                 = "raw"
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
	dvInitializationUpgrade             = true
	dvKeyboardLayout                    = "en-us"
	dvKVMArguments                      = ""
	dvMachineType                       = ""
	dvMemoryDedicated                   = 512
	dvMemoryFloating                    = 0
	dvMemoryShared                      = 0
	dvMemoryHugepages                   = ""
	dvMemoryKeepHugepages               = false
	dvMigrate                           = false
	dvName                              = ""

	dvOperatingSystemType = "other"
	dvPoolID              = ""
	dvProtection          = false
	dvSerialDeviceDevice  = "socket"
	dvSMBIOSFamily        = ""
	dvSMBIOSManufacturer  = ""
	dvSMBIOSProduct       = ""
	dvSMBIOSSKU           = ""
	dvSMBIOSSerial        = ""
	dvSMBIOSVersion       = ""
	dvStarted             = true
	dvStartupOrder        = -1
	dvStartupUpDelay      = -1
	dvStartupDownDelay    = -1
	dvTabletDevice        = true
	dvTemplate            = false
	dvTimeoutClone        = 1800
	dvTimeoutCreate       = 1800
	dvTimeoutMigrate      = 1800
	dvTimeoutReboot       = 1800
	dvTimeoutShutdownVM   = 1800
	dvTimeoutStartVM      = 1800
	dvTimeoutStopVM       = 300
	dvVGAClipboard        = ""
	dvVGAMemory           = 16
	dvVGAType             = "std"
	dvSCSIHardware        = "virtio-scsi-pci"
	dvStopOnDestroy       = false
	dvHookScript          = ""

	maxResourceVirtualEnvironmentVMAudioDevices   = 1
	maxResourceVirtualEnvironmentVMSerialDevices  = 4
	maxResourceVirtualEnvironmentVMHostPCIDevices = 8
	maxResourceVirtualEnvironmentVMHostUSBDevices = 4
	// hardcoded /usr/share/perl5/PVE/QemuServer/Memory.pm: "our $MAX_NUMA = 8".
	maxResourceVirtualEnvironmentVMNUMADevices = 8

	mkRebootAfterCreation = "reboot"
	mkOnBoot              = "on_boot"
	mkBootOrder           = "boot_order"
	mkACPI                = "acpi"
	mkAgent               = "agent"
	mkAgentEnabled        = "enabled"
	mkAgentTimeout        = "timeout"
	mkAgentTrim           = "trim"
	mkAgentType           = "type"
	mkAudioDevice         = "audio_device"
	mkAudioDeviceDevice   = "device"
	mkAudioDeviceDriver   = "driver"
	mkAudioDeviceEnabled  = "enabled"
	mkBIOS                = "bios"
	mkCDROM               = "cdrom"
	mkCDROMEnabled        = "enabled"
	mkCDROMFileID         = "file_id"
	mkCDROMInterface      = "interface"
	mkClone               = "clone"
	mkCloneRetries        = "retries"
	mkCloneDatastoreID    = "datastore_id"
	mkCloneNodeName       = "node_name"
	mkCloneVMID           = "vm_id"
	mkCloneFull           = "full"
	mkCPU                 = "cpu"
	mkCPUArchitecture     = "architecture"
	mkCPUCores            = "cores"
	mkCPUFlags            = "flags"
	mkCPUHotplugged       = "hotplugged"
	mkCPULimit            = "limit"
	mkCPUNUMA             = "numa"
	mkCPUSockets          = "sockets"
	mkCPUType             = "type"
	mkCPUUnits            = "units"
	mkCPUAffinity         = "affinity"
	mkDescription         = "description"

	mkNUMA              = "numa"
	mkNUMADevice        = "device"
	mkNUMACPUIDs        = "cpus"
	mkNUMAHostNodeNames = "hostnodes"
	mkNUMAMemory        = "memory"
	mkNUMAPolicy        = "policy"

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
	mkInitializationUpgrade             = "upgrade"

	mkKeyboardLayout      = "keyboard_layout"
	mkKVMArguments        = "kvm_arguments"
	mkMachine             = "machine"
	mkMemory              = "memory"
	mkMemoryDedicated     = "dedicated"
	mkMemoryFloating      = "floating"
	mkMemoryShared        = "shared"
	mkMemoryHugepages     = "hugepages"
	mkMemoryKeepHugepages = "keep_hugepages"
	mkMigrate             = "migrate"
	mkName                = "name"

	mkNodeName             = "node_name"
	mkOperatingSystem      = "operating_system"
	mkOperatingSystemType  = "type"
	mkPoolID               = "pool_id"
	mkProtection           = "protection"
	mkSerialDevice         = "serial_device"
	mkSerialDeviceDevice   = "device"
	mkSMBIOS               = "smbios"
	mkSMBIOSFamily         = "family"
	mkSMBIOSManufacturer   = "manufacturer"
	mkSMBIOSProduct        = "product"
	mkSMBIOSSKU            = "sku"
	mkSMBIOSSerial         = "serial"
	mkSMBIOSUUID           = "uuid"
	mkSMBIOSVersion        = "version"
	mkStarted              = "started"
	mkStartup              = "startup"
	mkStartupOrder         = "order"
	mkStartupUpDelay       = "up_delay"
	mkStartupDownDelay     = "down_delay"
	mkTabletDevice         = "tablet_device"
	mkTags                 = "tags"
	mkTemplate             = "template"
	mkTimeoutClone         = "timeout_clone"
	mkTimeoutCreate        = "timeout_create"
	mkTimeoutMigrate       = "timeout_migrate" // this is essentially an "timeout_update", needs to be refactored
	mkTimeoutReboot        = "timeout_reboot"
	mkTimeoutShutdownVM    = "timeout_shutdown_vm"
	mkTimeoutStartVM       = "timeout_start_vm"
	mkTimeoutStopVM        = "timeout_stop_vm"
	mkHostUSB              = "usb"
	mkHostUSBDevice        = "host"
	mkHostUSBDeviceMapping = "mapping"
	mkHostUSBDeviceUSB3    = "usb3"
	mkVGA                  = "vga"
	mkVGAClipboard         = "clipboard"
	mkVGAEnabled           = "enabled"
	mkVGAMemory            = "memory"
	mkVGAType              = "type"
	mkVMID                 = "vm_id"
	mkSCSIHardware         = "scsi_hardware"
	mkHookScriptFileID     = "hook_script_file_id"
	mkStopOnDestroy        = "stop_on_destroy"
)

// VM returns a resource that manages VMs.
func VM() *schema.Resource {
	s := map[string]*schema.Schema{
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
						ValidateDiagFunc: TimeoutValidator(),
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
						ValidateDiagFunc: QEMUAgentTypeValidator(),
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
						ValidateDiagFunc: AudioDeviceValidator(),
					},
					mkAudioDeviceDriver: {
						Type:             schema.TypeString,
						Description:      "The driver",
						Optional:         true,
						Default:          dvAudioDeviceDriver,
						ValidateDiagFunc: AudioDriverValidator(),
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
			ValidateDiagFunc: BIOSValidator(),
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
						ValidateDiagFunc: validators.FileID(),
					},
					mkCDROMInterface: {
						Type:             schema.TypeString,
						Description:      "The CDROM interface",
						Optional:         true,
						Default:          dvCDROMInterface,
						ValidateDiagFunc: IDEInterfaceValidator(),
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
						ValidateDiagFunc: VMIDValidator(),
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
						mkCPUAffinity:     dvCPUAffinity,
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
						ValidateDiagFunc: CPUArchitectureValidator(),
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
						ValidateDiagFunc: CPUTypeValidator(),
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
					mkCPUAffinity: {
						Type:             schema.TypeString,
						Description:      "The CPU affinity",
						Optional:         true,
						Default:          dvCPUAffinity,
						ValidateDiagFunc: CPUAffinityValidator(),
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
			StateFunc: func(i interface{}) string {
				// PVE always adds a newline to the description, so we have to do the same,
				// also taking in account the CLRF case (Windows)
				// Unlike container, VM description does not have trailing "\n"
				if i.(string) != "" {
					return strings.ReplaceAll(strings.TrimSpace(i.(string)), "\r\n", "\n")
				}

				return ""
			},
		},
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
						ValidateDiagFunc: validators.FileFormat(),
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
						ValidateDiagFunc: CloudInitInterfaceValidator(),
						DiffSuppressFunc: func(_, _, newValue string, _ *schema.ResourceData) bool {
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
									DiffSuppressFunc: func(_, oldVal, _ string, _ *schema.ResourceData) bool {
										return len(oldVal) > 0 &&
											strings.ReplaceAll(oldVal, "*", "") == ""
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
						ValidateDiagFunc: validators.FileID(),
					},
					mkInitializationVendorDataFileID: {
						Type:             schema.TypeString,
						Description:      "The ID of a file containing vendor data",
						Optional:         true,
						ForceNew:         true,
						Default:          dvInitializationVendorDataFileID,
						ValidateDiagFunc: validators.FileID(),
					},
					mkInitializationNetworkDataFileID: {
						Type:             schema.TypeString,
						Description:      "The ID of a file containing network config",
						Optional:         true,
						ForceNew:         true,
						Default:          dvInitializationNetworkDataFileID,
						ValidateDiagFunc: validators.FileID(),
					},
					mkInitializationMetaDataFileID: {
						Type:             schema.TypeString,
						Description:      "The ID of a file containing meta data config",
						Optional:         true,
						ForceNew:         true,
						Default:          dvInitializationMetaDataFileID,
						ValidateDiagFunc: validators.FileID(),
					},
					mkInitializationType: {
						Type:             schema.TypeString,
						Description:      "The cloud-init configuration format",
						Optional:         true,
						ForceNew:         true,
						Default:          dvInitializationType,
						ValidateDiagFunc: CloudInitTypeValidator(),
					},
					mkInitializationUpgrade: {
						Type:        schema.TypeBool,
						Description: "Whether to do an automatic package upgrade after the first boot",
						Optional:    true,
						Computed:    true,
						Deprecated:  "The `upgrade` attribute is deprecated and will be removed in a future release.",
					},
				},
			},
			MaxItems: 1,
			MinItems: 0,
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
						Description: "Marks the PCI(e) device as the primary GPU of the VM. With this enabled, " +
							"the vga configuration argument will be ignored.",
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
						Optional:    true,
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
			ValidateDiagFunc: KeyboardLayoutValidator(),
		},
		mkMachine: {
			Type:             schema.TypeString,
			Description:      "The VM machine type, either default `pc` or `q35`",
			Optional:         true,
			Default:          dvMachineType,
			ValidateDiagFunc: MachineTypeValidator(),
		},
		mkMemory: {
			Type:        schema.TypeList,
			Description: "The memory allocation",
			Optional:    true,
			DefaultFunc: func() (interface{}, error) {
				return []interface{}{
					map[string]interface{}{
						mkMemoryDedicated:     dvMemoryDedicated,
						mkMemoryFloating:      dvMemoryFloating,
						mkMemoryShared:        dvMemoryShared,
						mkMemoryHugepages:     dvMemoryHugepages,
						mkMemoryKeepHugepages: dvMemoryKeepHugepages,
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
					mkMemoryHugepages: {
						Type:         schema.TypeString,
						Description:  "Enable/disable hugepages memory",
						Optional:     true,
						Default:      dvMemoryHugepages,
						RequiredWith: []string{"cpu.0.numa"},
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
							"1024",
							"2",
							"any",
						}, true)),
					},
					mkMemoryKeepHugepages: {
						Type:         schema.TypeBool,
						Description:  "Hugepages will not be deleted after VM shutdown and can be used for subsequent starts",
						Optional:     true,
						Default:      dvMemoryKeepHugepages,
						RequiredWith: []string{"cpu.0.numa"},
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
		mkNodeName: {
			Type:        schema.TypeString,
			Description: "The node name",
			Required:    true,
		},
		mkNUMA: {
			Type:        schema.TypeList,
			Description: "The NUMA topology",
			Optional:    true,
			ForceNew:    false,
			DefaultFunc: func() (interface{}, error) {
				return []interface{}{}, nil
			},
			DiffSuppressFunc:      structure.SuppressIfListsOfMapsAreEqualIgnoringOrderByKey(mkNUMADevice),
			DiffSuppressOnRefresh: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					mkNUMADevice: {
						Type:         schema.TypeString,
						Description:  "Numa node device ID",
						Optional:     false,
						Required:     true,
						RequiredWith: []string{"cpu.0.numa"},
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringMatch(
							regexp.MustCompile(`^numa\d+$`),
							"numa node device ID must be in the format 'numaX' where X is a number",
						)),
					},
					mkNUMACPUIDs: {
						Type:             schema.TypeString,
						Description:      "CPUs accessing this NUMA node",
						Optional:         false,
						Required:         true,
						RequiredWith:     []string{"cpu.0.numa"},
						ValidateDiagFunc: RangeSemicolonValidator(),
					},
					mkNUMAMemory: {
						Type:         schema.TypeInt,
						Description:  "Amount of memory this NUMA node provides",
						Optional:     false,
						Required:     true,
						RequiredWith: []string{"cpu.0.numa"},
						ValidateDiagFunc: validation.ToDiagFunc(
							validation.IntBetween(64, 268435456),
						),
					},
					mkNUMAHostNodeNames: {
						Type:             schema.TypeString,
						Description:      "Host NUMA nodes to use",
						Optional:         true,
						RequiredWith:     []string{"cpu.0.numa"},
						ValidateDiagFunc: RangeSemicolonValidator(),
					},
					mkNUMAPolicy: {
						Type:        schema.TypeString,
						Description: "NUMA policy",
						Optional:    true,
						Default:     "preferred",
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
							"bind",
							"interleave",
							"preferred",
						}, true)),
					},
				},
			},
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
						ValidateDiagFunc: OperatingSystemTypeValidator(),
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
		mkProtection: {
			Type:        schema.TypeBool,
			Description: "Sets the protection flag of the VM. This will disable the remove VM and remove disk operations",
			Optional:    true,
			Default:     dvProtection,
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
						ValidateDiagFunc: SerialDeviceValidator(),
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
			DiffSuppressFunc: func(_, _, _ string, d *schema.ResourceData) bool {
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
		"timeout_move_disk": {
			Type:        schema.TypeInt,
			Description: "MoveDisk timeout",
			Optional:    true,
			Default:     1800,
			Deprecated: "This field is deprecated and will be removed in a future release. " +
				"An overall operation timeout (timeout_create / timeout_clone / timeout_migrate) is used instead.",
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
						mkVGAClipboard: dvVGAClipboard,
						mkVGAMemory:    dvVGAMemory,
						mkVGAType:      dvVGAType,
					},
				}, nil
			},
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					mkVGAClipboard: {
						Type:        schema.TypeString,
						Description: "Enable clipboard support. Set to `vnc` to enable clipboard support for VNC.",
						Optional:    true,
						Default:     dvVGAClipboard,
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
							"",
							"vnc",
						}, true)),
					},
					mkVGAEnabled: {
						Type: schema.TypeBool,
						Deprecated: "The `enabled` attribute is deprecated and will be removed in a future release. " +
							"Use type `none` instead.",
						Description: "Whether to enable the VGA device",
						Optional:    true,
					},
					mkVGAMemory: {
						Type:             schema.TypeInt,
						Description:      "The VGA memory in megabytes (4-512 MB)",
						Optional:         true,
						Default:          dvVGAMemory,
						ValidateDiagFunc: VGAMemoryValidator(),
					},
					mkVGAType: {
						Type:             schema.TypeString,
						Description:      "The VGA type",
						Optional:         true,
						Default:          dvVGAType,
						ValidateDiagFunc: VGATypeValidator(),
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
			ValidateDiagFunc: VMIDValidator(),
		},
		mkSCSIHardware: {
			Type:             schema.TypeString,
			Description:      "The SCSI hardware type",
			Optional:         true,
			Default:          dvSCSIHardware,
			ValidateDiagFunc: SCSIHardwareValidator(),
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
	}

	structure.MergeSchema(s, disk.Schema())
	structure.MergeSchema(s, network.Schema())

	return &schema.Resource{
		Schema:        s,
		CreateContext: vmCreate,
		ReadContext:   vmRead,
		UpdateContext: vmUpdate,
		DeleteContext: vmDelete,
		CustomizeDiff: customdiff.All(
			customdiff.All(network.CustomizeDiff()...),
			customdiff.ForceNewIf(
				mkVMID,
				func(_ context.Context, d *schema.ResourceDiff, _ interface{}) bool {
					newValue := d.Get(mkVMID)

					// 'vm_id' is ForceNew, except when changing 'vm_id' to existing correct id
					// (automatic fix from -1 to actual vm_id must not re-create VM)
					return strconv.Itoa(newValue.(int)) != d.Id()
				},
			),
			customdiff.ForceNewIf(
				mkNodeName,
				func(_ context.Context, d *schema.ResourceDiff, _ interface{}) bool {
					return !d.Get(mkMigrate).(bool)
				},
			),
		),
		Importer: &schema.ResourceImporter{
			StateContext: func(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
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

// ConvertToStringSlice helps convert interface slice to string slice.
func ConvertToStringSlice(interfaceSlice []interface{}) []string {
	resultSlice := []string{}
	for _, val := range interfaceSlice {
		resultSlice = append(resultSlice, val.(string))
	}

	return resultSlice
}

func vmCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clone := d.Get(mkClone).([]interface{})

	if len(clone) > 0 {
		return vmCreateClone(ctx, d, m)
	}

	return vmCreateCustom(ctx, d, m)
}

// Check for an existing CloudInit IDE drive. If no such drive is found, return the specified `defaultValue`.
func findExistingCloudInitDrive(vmConfig *vms.GetResponseData, vmID int, defaultValue string) string {
	ideDevices := []*vms.CustomStorageDevice{
		vmConfig.IDEDevice0,
		vmConfig.IDEDevice1,
		vmConfig.IDEDevice2,
		vmConfig.IDEDevice3,
	}
	for i, device := range ideDevices {
		if device != nil && device.Enabled && device.IsCloudInitDrive(vmID) {
			return fmt.Sprintf("ide%d", i)
		}
	}

	sataDevices := []*vms.CustomStorageDevice{
		vmConfig.SATADevice0,
		vmConfig.SATADevice1,
		vmConfig.SATADevice2,
		vmConfig.SATADevice3,
		vmConfig.SATADevice4,
		vmConfig.SATADevice5,
	}
	for i, device := range sataDevices {
		if device != nil && device.Enabled && device.IsCloudInitDrive(vmID) {
			return fmt.Sprintf("sata%d", i)
		}
	}

	scsiDevices := []*vms.CustomStorageDevice{
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
	for i, device := range scsiDevices {
		if device != nil && device.Enabled && device.IsCloudInitDrive(vmID) {
			return fmt.Sprintf("scsi%d", i)
		}
	}

	return defaultValue
}

// Return a pointer to the storage device configuration based on a name. The device name is assumed to be a
// valid ide, sata, or scsi interface name.
func getStorageDevice(vmConfig *vms.GetResponseData, deviceName string) *vms.CustomStorageDevice {
	switch deviceName {
	case "ide0":
		return vmConfig.IDEDevice0
	case "ide1":
		return vmConfig.IDEDevice1
	case "ide2":
		return vmConfig.IDEDevice2
	case "ide3":
		return vmConfig.IDEDevice3

	case "sata0":
		return vmConfig.SATADevice0
	case "sata1":
		return vmConfig.SATADevice1
	case "sata2":
		return vmConfig.SATADevice2
	case "sata3":
		return vmConfig.SATADevice3
	case "sata4":
		return vmConfig.SATADevice4
	case "sata5":
		return vmConfig.SATADevice5

	case "scsi0":
		return vmConfig.SCSIDevice0
	case "scsi1":
		return vmConfig.SCSIDevice1
	case "scsi2":
		return vmConfig.SCSIDevice2
	case "scsi3":
		return vmConfig.SCSIDevice3
	case "scsi4":
		return vmConfig.SCSIDevice4
	case "scsi5":
		return vmConfig.SCSIDevice5
	case "scsi6":
		return vmConfig.SCSIDevice6
	case "scsi7":
		return vmConfig.SCSIDevice7
	case "scsi8":
		return vmConfig.SCSIDevice8
	case "scsi9":
		return vmConfig.SCSIDevice9
	case "scsi10":
		return vmConfig.SCSIDevice10
	case "scsi11":
		return vmConfig.SCSIDevice11
	case "scsi12":
		return vmConfig.SCSIDevice12
	case "scsi13":
		return vmConfig.SCSIDevice13

	default:
		return nil
	}
}

// Delete IDE interfaces that can then be used for CloudInit. The first interface will always
// be deleted. The second will be deleted only if it isn't empty and isn't the same as the
// first.
func deleteIdeDrives(ctx context.Context, vmAPI *vms.Client, itf1 string, itf2 string) diag.Diagnostics {
	ddUpdateBody := &vms.UpdateRequestBody{}
	ddUpdateBody.Delete = append(ddUpdateBody.Delete, itf1)
	tflog.Debug(ctx, fmt.Sprintf("Deleting IDE interface '%s'", itf1))

	if itf2 != "" && itf2 != itf1 {
		ddUpdateBody.Delete = append(ddUpdateBody.Delete, itf2)
		tflog.Debug(ctx, fmt.Sprintf("Deleting IDE interface '%s'", itf2))
	}

	e := vmAPI.UpdateVM(ctx, ddUpdateBody)
	if e != nil {
		return diag.FromErr(e)
	}

	return nil
}

// Start the VM, then wait for it to actually start; it may not be started immediately if running in HA mode.
func vmStart(ctx context.Context, vmAPI *vms.Client, d *schema.ResourceData) diag.Diagnostics {
	tflog.Debug(ctx, "Starting VM")

	startTimeoutSec := d.Get(mkTimeoutStartVM).(int)

	ctx, cancel := context.WithTimeout(ctx, time.Duration(startTimeoutSec)*time.Second)
	defer cancel()

	var diags diag.Diagnostics

	log, e := vmAPI.StartVM(ctx, startTimeoutSec)
	if e != nil {
		return diag.FromErr(e)
	}

	if len(log) > 0 {
		lines := "\n\t| " + strings.Join(log, "\n\t| ")
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  fmt.Sprintf("the VM startup task finished with a warning, task log:\n%s", lines),
		})
	}

	return append(diags, diag.FromErr(vmAPI.WaitForVMStatus(ctx, "running"))...)
}

// Shutdown the VM, then wait for it to actually shut down (it may not be shut down immediately if
// running in HA mode).
func vmShutdown(ctx context.Context, vmAPI *vms.Client, d *schema.ResourceData) diag.Diagnostics {
	tflog.Debug(ctx, "Shutting down VM")

	forceStop := types.CustomBool(true)
	shutdownTimeoutSec := d.Get(mkTimeoutShutdownVM).(int)

	ctx, cancel := context.WithTimeout(ctx, time.Duration(shutdownTimeoutSec)*time.Second)
	defer cancel()

	e := vmAPI.ShutdownVM(ctx, &vms.ShutdownRequestBody{
		ForceStop: &forceStop,
		Timeout:   &shutdownTimeoutSec,
	})
	if e != nil {
		return diag.FromErr(e)
	}

	return diag.FromErr(vmAPI.WaitForVMStatus(ctx, "stopped"))
}

// Forcefully stop the VM, then wait for it to actually stop.
func vmStop(ctx context.Context, vmAPI *vms.Client, d *schema.ResourceData) diag.Diagnostics {
	tflog.Debug(ctx, "Stopping VM")

	stopTimeout := d.Get(mkTimeoutStopVM).(int)

	ctx, cancel := context.WithTimeout(ctx, time.Duration(stopTimeout)*time.Second)
	defer cancel()

	e := vmAPI.StopVM(ctx)
	if e != nil {
		return diag.FromErr(e)
	}

	return diag.FromErr(vmAPI.WaitForVMStatus(ctx, "stopped"))
}

func vmCreateClone(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cloneTimeoutSec := d.Get(mkTimeoutClone).(int)

	ctx, cancel := context.WithTimeout(ctx, time.Duration(cloneTimeoutSec)*time.Second)
	defer cancel()

	config := m.(proxmoxtf.ProviderConfiguration)

	client, e := config.GetClient()
	if e != nil {
		return diag.FromErr(e)
	}

	clone := d.Get(mkClone).([]interface{})
	cloneBlock := clone[0].(map[string]interface{})
	cloneRetries := cloneBlock[mkCloneRetries].(int)
	cloneDatastoreID := cloneBlock[mkCloneDatastoreID].(string)
	cloneNodeName := cloneBlock[mkCloneNodeName].(string)
	cloneVMID := cloneBlock[mkCloneVMID].(int)
	cloneFull := cloneBlock[mkCloneFull].(bool)

	description := d.Get(mkDescription).(string)
	name := d.Get(mkName).(string)
	tags := d.Get(mkTags).([]interface{})
	nodeName := d.Get(mkNodeName).(string)
	poolID := d.Get(mkPoolID).(string)
	vmIDUntyped, hasVMID := d.GetOk(mkVMID)
	vmID := vmIDUntyped.(int)

	if !hasVMID {
		vmIDNew, err := client.Cluster().GetVMID(ctx)
		if err != nil {
			return diag.FromErr(err)
		}

		vmID = *vmIDNew

		err = d.Set(mkVMID, vmID)
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

	if cloneNodeName != "" && cloneNodeName != nodeName {
		// Check if any used datastores of the source VM are not shared
		vmConfig, err := client.Node(cloneNodeName).VM(cloneVMID).GetVM(ctx)
		if err != nil {
			return diag.FromErr(err)
		}

		datastores := getDiskDatastores(vmConfig, d)

		onlySharedDatastores := true

		for _, datastore := range datastores {
			datastoreStatus, err2 := client.Node(cloneNodeName).Storage(datastore).GetDatastoreStatus(ctx)
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

			err = client.Node(cloneNodeName).VM(cloneVMID).CloneVM(ctx, cloneRetries, cloneBody)
			if err != nil {
				return diag.FromErr(err)
			}
		} else { //nolint:wsl
			// If the source and the target node are not the same and any used datastore in the source VM is
			//  not shared, clone to the source node and then migrate to the target node. This is a workaround
			//  for missing functionality in the proxmox api as recommended per
			//  https://forum.proxmox.com/threads/500-cant-clone-to-non-shared-storage-local.49078/#post-229727

			// Temporarily clone to local node
			err = client.Node(cloneNodeName).VM(cloneVMID).CloneVM(ctx, cloneRetries, cloneBody)
			if err != nil {
				return diag.FromErr(err)
			}

			// Wait for the virtual machine to be created and its configuration lock to be released before migrating.

			err = client.Node(cloneNodeName).VM(vmID).WaitForVMConfigUnlock(ctx, true)
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

			err = client.Node(cloneNodeName).VM(vmID).MigrateVM(ctx, migrateBody)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	} else {
		e = client.Node(nodeName).VM(cloneVMID).CloneVM(ctx, cloneRetries, cloneBody)
	}

	if e != nil {
		return diag.FromErr(e)
	}

	d.SetId(strconv.Itoa(vmID))

	vmAPI := client.Node(nodeName).VM(vmID)

	// Wait for the virtual machine to be created and its configuration lock to be released.
	e = vmAPI.WaitForVMConfigUnlock(ctx, true)
	if e != nil {
		return diag.FromErr(e)
	}

	// Now that the virtual machine has been cloned, we need to perform some modifications.
	acpi := types.CustomBool(d.Get(mkACPI).(bool))
	agent := d.Get(mkAgent).([]interface{})
	audioDevices := vmGetAudioDeviceList(d)

	bios := d.Get(mkBIOS).(string)
	kvmArguments := d.Get(mkKVMArguments).(string)
	scsiHardware := d.Get(mkSCSIHardware).(string)
	cdrom := d.Get(mkCDROM).([]interface{})
	cpu := d.Get(mkCPU).([]interface{})
	initialization := d.Get(mkInitialization).([]interface{})
	hostPCI := d.Get(mkHostPCI).([]interface{})
	hostUSB := d.Get(mkHostUSB).([]interface{})
	keyboardLayout := d.Get(mkKeyboardLayout).(string)
	memory := d.Get(mkMemory).([]interface{})
	numa := d.Get(mkNUMA).([]interface{})
	operatingSystem := d.Get(mkOperatingSystem).([]interface{})
	serialDevice := d.Get(mkSerialDevice).([]interface{})
	onBoot := types.CustomBool(d.Get(mkOnBoot).(bool))
	tabletDevice := types.CustomBool(d.Get(mkTabletDevice).(bool))
	protection := types.CustomBool(d.Get(mkProtection).(bool))
	template := types.CustomBool(d.Get(mkTemplate).(bool))
	vga := d.Get(mkVGA).([]interface{})

	updateBody := &vms.UpdateRequestBody{
		AudioDevices: audioDevices,
	}

	ideDevices := vms.CustomStorageDevices{}

	var del []string

	//nolint:gosimple
	if acpi != dvACPI {
		updateBody.ACPI = &acpi
	}

	if len(agent) > 0 && agent[0] != nil {
		agentBlock := agent[0].(map[string]interface{})

		agentEnabled := types.CustomBool(
			agentBlock[mkAgentEnabled].(bool),
		)
		agentTrim := types.CustomBool(agentBlock[mkAgentTrim].(bool))
		agentType := agentBlock[mkAgentType].(string)

		updateBody.Agent = &vms.CustomAgent{
			Enabled:         &agentEnabled,
			TrimClonedDisks: &agentTrim,
			Type:            &agentType,
		}
	}

	if kvmArguments != dvKVMArguments {
		updateBody.KVMArguments = &kvmArguments
	}

	if bios != dvBIOS {
		updateBody.BIOS = &bios
	}

	if scsiHardware != dvSCSIHardware {
		updateBody.SCSIHardware = &scsiHardware
	}

	if len(cdrom) > 0 || len(initialization) > 0 {
		ideDevices = vms.CustomStorageDevices{
			"ide0": &vms.CustomStorageDevice{
				Enabled: false,
			},
			"ide1": &vms.CustomStorageDevice{
				Enabled: false,
			},
			"ide2": &vms.CustomStorageDevice{
				Enabled: false,
			},
			"ide3": &vms.CustomStorageDevice{
				Enabled: false,
			},
		}
	}

	if len(cdrom) > 0 && cdrom[0] != nil {
		cdromBlock := cdrom[0].(map[string]interface{})

		cdromEnabled := cdromBlock[mkCDROMEnabled].(bool)
		cdromFileID := cdromBlock[mkCDROMFileID].(string)
		cdromInterface := cdromBlock[mkCDROMInterface].(string)

		if cdromFileID == "" {
			cdromFileID = "cdrom"
		}

		cdromMedia := "cdrom"

		ideDevices[cdromInterface] = &vms.CustomStorageDevice{
			Enabled:    cdromEnabled,
			FileVolume: cdromFileID,
			Media:      &cdromMedia,
		}
	}

	if len(cpu) > 0 && cpu[0] != nil {
		cpuBlock := cpu[0].(map[string]interface{})

		cpuArchitecture := cpuBlock[mkCPUArchitecture].(string)
		cpuCores := cpuBlock[mkCPUCores].(int)
		cpuFlags := cpuBlock[mkCPUFlags].([]interface{})
		cpuHotplugged := cpuBlock[mkCPUHotplugged].(int)
		cpuLimit := cpuBlock[mkCPULimit].(int)
		cpuNUMA := types.CustomBool(cpuBlock[mkCPUNUMA].(bool))
		cpuSockets := cpuBlock[mkCPUSockets].(int)
		cpuType := cpuBlock[mkCPUType].(string)
		cpuUnits := cpuBlock[mkCPUUnits].(int)
		cpuAffinity := cpuBlock[mkCPUAffinity].(string)

		cpuFlagsConverted := make([]string, len(cpuFlags))

		for fi, flag := range cpuFlags {
			cpuFlagsConverted[fi] = flag.(string)
		}

		// Only the root account is allowed to change the CPU architecture, which makes this check necessary.
		if client.API().IsRootTicket() ||
			cpuArchitecture != dvCPUArchitecture {
			updateBody.CPUArchitecture = &cpuArchitecture
		}

		updateBody.CPUCores = ptr.Ptr(int64(cpuCores))
		updateBody.CPUEmulation = &vms.CustomCPUEmulation{
			Flags: &cpuFlagsConverted,
			Type:  cpuType,
		}
		updateBody.NUMAEnabled = &cpuNUMA
		updateBody.CPUSockets = ptr.Ptr(int64(cpuSockets))
		updateBody.CPUUnits = ptr.Ptr(int64(cpuUnits))

		if cpuAffinity != "" {
			updateBody.CPUAffinity = &cpuAffinity
		}

		if cpuHotplugged > 0 {
			updateBody.VirtualCPUCount = ptr.Ptr(int64(cpuHotplugged))
		}

		if cpuLimit > 0 {
			updateBody.CPULimit = ptr.Ptr(int64(cpuLimit))
		}
	}

	vmConfig, err := vmAPI.GetVM(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(initialization) > 0 && initialization[0] != nil {
		tflog.Trace(ctx, "Preparing the CloudInit configuration")

		initializationBlock := initialization[0].(map[string]interface{})
		initializationDatastoreID := initializationBlock[mkInitializationDatastoreID].(string)
		initializationInterface := initializationBlock[mkInitializationInterface].(string)

		existingInterface := findExistingCloudInitDrive(vmConfig, vmID, "ide2")
		if initializationInterface == "" {
			initializationInterface = existingInterface
		}

		tflog.Trace(ctx, fmt.Sprintf("CloudInit IDE interface is '%s'", initializationInterface))

		const cdromCloudInitEnabled = true

		cdromCloudInitFileID := fmt.Sprintf("%s:cloudinit", initializationDatastoreID)
		cdromCloudInitMedia := "cdrom"
		ideDevices[initializationInterface] = &vms.CustomStorageDevice{
			Enabled:    cdromCloudInitEnabled,
			FileVolume: cdromCloudInitFileID,
			Media:      &cdromCloudInitMedia,
		}

		if err := deleteIdeDrives(ctx, vmAPI, initializationInterface, existingInterface); err != nil {
			return err
		}

		updateBody.CloudInitConfig = vmGetCloudInitConfig(d)
	}

	if len(hostPCI) > 0 {
		updateBody.PCIDevices = vmGetHostPCIDeviceObjects(d)
	}

	if len(numa) > 0 {
		updateBody.NUMADevices = vmGetNumaDeviceObjects(d)
	}

	if len(hostUSB) > 0 {
		updateBody.USBDevices = vmGetHostUSBDeviceObjects(d)
	}

	if len(cdrom) > 0 || len(initialization) > 0 {
		updateBody.IDEDevices = ideDevices
	}

	if keyboardLayout != dvKeyboardLayout {
		updateBody.KeyboardLayout = &keyboardLayout
	}

	if len(memory) > 0 && memory[0] != nil {
		memoryBlock := memory[0].(map[string]interface{})

		memoryDedicated := memoryBlock[mkMemoryDedicated].(int)
		memoryFloating := memoryBlock[mkMemoryFloating].(int)
		memoryShared := memoryBlock[mkMemoryShared].(int)
		hugepages := memoryBlock[mkMemoryHugepages].(string)
		keepHugepages := types.CustomBool(memoryBlock[mkMemoryKeepHugepages].(bool))

		updateBody.DedicatedMemory = &memoryDedicated
		updateBody.FloatingMemory = &memoryFloating

		if memoryShared > 0 {
			memorySharedName := fmt.Sprintf("vm-%d-ivshmem", vmID)

			updateBody.SharedMemory = &vms.CustomSharedMemory{
				Name: &memorySharedName,
				Size: memoryShared,
			}
		}

		if hugepages != "" {
			updateBody.Hugepages = &hugepages
		}

		if keepHugepages {
			updateBody.KeepHugepages = &keepHugepages
		}
	}

	networkDevice := d.Get(network.MkNetworkDevice).([]interface{})
	if len(networkDevice) > 0 {
		updateBody.NetworkDevices, err = network.GetNetworkDeviceObjects(d)
		if err != nil {
			return diag.FromErr(err)
		}

		for i, ni := range updateBody.NetworkDevices {
			if !ni.Enabled {
				del = append(del, fmt.Sprintf("net%d", i))
			}
		}

		for i := len(updateBody.NetworkDevices); i < network.MaxNetworkDevices; i++ {
			del = append(del, fmt.Sprintf("net%d", i))
		}
	}

	if len(operatingSystem) > 0 && operatingSystem[0] != nil {
		operatingSystemBlock := operatingSystem[0].(map[string]interface{})
		operatingSystemType := operatingSystemBlock[mkOperatingSystemType].(string)

		updateBody.OSType = &operatingSystemType
	}

	if len(serialDevice) > 0 {
		updateBody.SerialDevices = vmGetSerialDeviceList(d)

		for i := len(updateBody.SerialDevices); i < maxResourceVirtualEnvironmentVMSerialDevices; i++ {
			del = append(del, fmt.Sprintf("serial%d", i))
		}
	}

	updateBody.StartOnBoot = &onBoot

	updateBody.SMBIOS = vmGetSMBIOS(d)

	updateBody.StartupOrder = vmGetStartupOrder(d)

	//nolint:gosimple
	if tabletDevice != dvTabletDevice {
		updateBody.TabletDeviceEnabled = &tabletDevice
	}

	//nolint:gosimple
	if protection != dvProtection {
		updateBody.DeletionProtection = &protection
	}

	if len(tags) > 0 {
		tagString := vmGetTagsString(d)
		updateBody.Tags = &tagString
	}

	//nolint:gosimple
	if template != dvTemplate {
		updateBody.Template = &template
	}

	if len(vga) > 0 {
		vgaDevice, err := vmGetVGADeviceObject(d)
		if err != nil {
			return diag.FromErr(err)
		}

		updateBody.VGADevice = vgaDevice
	}

	hookScript := d.Get(mkHookScriptFileID).(string)
	currentHookScript := vmConfig.HookScript

	if len(hookScript) > 0 {
		updateBody.HookScript = &hookScript
	} else if currentHookScript != nil {
		del = append(del, "hookscript")
	}

	updateBody.Delete = del

	e = vmAPI.UpdateVM(ctx, updateBody)
	if e != nil {
		return diag.FromErr(e)
	}

	vmConfig, e = vmAPI.GetVM(ctx)
	if e != nil {
		if errors.Is(e, api.ErrResourceDoesNotExist) {
			d.SetId("")
			return nil
		}

		return diag.FromErr(e)
	}

	/////////////////

	allDiskInfo := disk.GetInfo(vmConfig, d) // from the cloned VM

	planDisks, e := disk.GetDiskDeviceObjects(d, VM(), nil) // from the resource config
	if e != nil {
		return diag.FromErr(e)
	}

	e = disk.CreateClone(ctx, d, planDisks, allDiskInfo, vmAPI)
	if e != nil {
		return diag.FromErr(e)
	}

	efiDisk := d.Get(mkEFIDisk).([]interface{})
	efiDiskInfo := vmGetEfiDisk(d, nil) // from the resource config

	for i := range efiDisk {
		diskBlock := efiDisk[i].(map[string]interface{})
		diskInterface := "efidisk0"
		dataStoreID := diskBlock[mkEFIDiskDatastoreID].(string)
		efiType := diskBlock[mkEFIDiskType].(string)

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

		if efiType != *currentDiskInfo.Type {
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
			e = vmAPI.MoveVMDisk(ctx, diskMoveBody)
			if e != nil {
				return diag.FromErr(e)
			}
		}
	}

	tpmState := d.Get(mkTPMState).([]interface{})
	tpmStateInfo := vmGetTPMState(d, nil) // from the resource config

	for i := range tpmState {
		diskBlock := tpmState[i].(map[string]interface{})
		diskInterface := "tpmstate0"
		dataStoreID := diskBlock[mkTPMStateDatastoreID].(string)

		currentTPMState := vmConfig.TPMState
		configuredTPMStateInfo := tpmStateInfo

		if currentTPMState == nil {
			diskUpdateBody := &vms.UpdateRequestBody{}

			diskUpdateBody.TPMState = configuredTPMStateInfo

			e = vmAPI.UpdateVM(ctx, diskUpdateBody)
			if e != nil {
				return diag.FromErr(e)
			}

			continue
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
			e = vmAPI.MoveVMDisk(ctx, diskMoveBody)
			if e != nil {
				return diag.FromErr(e)
			}
		}
	}

	return vmCreateStart(ctx, d, m)
}

func vmCreateCustom(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	createTimeoutSec := d.Get(mkTimeoutCreate).(int)

	ctx, cancel := context.WithTimeout(ctx, time.Duration(createTimeoutSec)*time.Second)
	defer cancel()

	config := m.(proxmoxtf.ProviderConfiguration)

	client, err := config.GetClient()
	if err != nil {
		return diag.FromErr(err)
	}

	resource := VM()

	acpi := types.CustomBool(d.Get(mkACPI).(bool))

	agentBlock, err := structure.GetSchemaBlock(
		resource,
		d,
		[]string{mkAgent},
		0,
		true,
	)
	if err != nil {
		return diag.FromErr(err)
	}

	agentEnabled := types.CustomBool(
		agentBlock[mkAgentEnabled].(bool),
	)
	agentTrim := types.CustomBool(agentBlock[mkAgentTrim].(bool))
	agentType := agentBlock[mkAgentType].(string)

	kvmArguments := d.Get(mkKVMArguments).(string)

	audioDevices := vmGetAudioDeviceList(d)

	bios := d.Get(mkBIOS).(string)

	cdromBlock, err := structure.GetSchemaBlock(
		resource,
		d,
		[]string{mkCDROM},
		0,
		true,
	)
	if err != nil {
		return diag.FromErr(err)
	}

	cdromEnabled := cdromBlock[mkCDROMEnabled].(bool)
	cdromFileID := cdromBlock[mkCDROMFileID].(string)
	cdromInterface := cdromBlock[mkCDROMInterface].(string)

	cdromCloudInitEnabled := false
	cdromCloudInitFileID := ""
	cdromCloudInitInterface := ""

	if cdromFileID == "" {
		cdromFileID = "cdrom"
	}

	cpuBlock, err := structure.GetSchemaBlock(
		resource,
		d,
		[]string{mkCPU},
		0,
		true,
	)
	if err != nil {
		return diag.FromErr(err)
	}

	cpuArchitecture := cpuBlock[mkCPUArchitecture].(string)
	cpuCores := cpuBlock[mkCPUCores].(int)
	cpuFlags := cpuBlock[mkCPUFlags].([]interface{})
	cpuHotplugged := cpuBlock[mkCPUHotplugged].(int)
	cpuLimit := cpuBlock[mkCPULimit].(int)
	cpuSockets := cpuBlock[mkCPUSockets].(int)
	cpuNUMA := types.CustomBool(cpuBlock[mkCPUNUMA].(bool))
	cpuType := cpuBlock[mkCPUType].(string)
	cpuUnits := cpuBlock[mkCPUUnits].(int)
	cpuAffinity := cpuBlock[mkCPUAffinity].(string)

	description := d.Get(mkDescription).(string)

	var efiDisk *vms.CustomEFIDisk

	efiDiskBlock := d.Get(mkEFIDisk).([]interface{})
	if len(efiDiskBlock) > 0 && efiDiskBlock[0] != nil {
		block := efiDiskBlock[0].(map[string]interface{})

		datastoreID, _ := block[mkEFIDiskDatastoreID].(string)
		fileFormat, _ := block[mkEFIDiskFileFormat].(string)
		efiType, _ := block[mkEFIDiskType].(string)
		preEnrolledKeys := types.CustomBool(block[mkEFIDiskPreEnrolledKeys].(bool))

		if fileFormat == "" {
			fileFormat = dvEFIDiskFileFormat
		}

		efiDisk = &vms.CustomEFIDisk{
			Type:            &efiType,
			FileVolume:      fmt.Sprintf("%s:1", datastoreID),
			Format:          &fileFormat,
			PreEnrolledKeys: &preEnrolledKeys,
		}
	}

	var tpmState *vms.CustomTPMState

	tpmStateBlock := d.Get(mkTPMState).([]interface{})
	if len(tpmStateBlock) > 0 && tpmStateBlock[0] != nil {
		block := tpmStateBlock[0].(map[string]interface{})

		datastoreID, _ := block[mkTPMStateDatastoreID].(string)
		version, _ := block[mkTPMStateVersion].(string)

		if version == "" {
			version = dvTPMStateVersion
		}

		tpmState = &vms.CustomTPMState{
			FileVolume: fmt.Sprintf("%s:1", datastoreID),
			Version:    &version,
		}
	}

	initializationConfig := vmGetCloudInitConfig(d)
	initializationAttr := d.Get(mkInitialization)

	if initializationConfig != nil && initializationAttr != nil {
		initialization := initializationAttr.([]interface{})

		initializationBlock := initialization[0].(map[string]interface{})
		initializationDatastoreID := initializationBlock[mkInitializationDatastoreID].(string)

		cdromCloudInitEnabled = true
		cdromCloudInitFileID = fmt.Sprintf("%s:cloudinit", initializationDatastoreID)

		cdromCloudInitInterface = initializationBlock[mkInitializationInterface].(string)
		if cdromCloudInitInterface == "" {
			cdromCloudInitInterface = "ide2"
		}
	}

	pciDeviceObjects := vmGetHostPCIDeviceObjects(d)

	numaDeviceObjects := vmGetNumaDeviceObjects(d)

	usbDeviceObjects := vmGetHostUSBDeviceObjects(d)

	keyboardLayout := d.Get(mkKeyboardLayout).(string)

	memoryBlock, err := structure.GetSchemaBlock(
		resource,
		d,
		[]string{mkMemory},
		0,
		true,
	)
	if err != nil {
		return diag.FromErr(err)
	}

	memoryDedicated := memoryBlock[mkMemoryDedicated].(int)
	memoryFloating := memoryBlock[mkMemoryFloating].(int)
	memoryShared := memoryBlock[mkMemoryShared].(int)
	memoryHugepages := memoryBlock[mkMemoryHugepages].(string)
	memoryKeepHugepages := types.CustomBool(memoryBlock[mkMemoryKeepHugepages].(bool))

	machine := d.Get(mkMachine).(string)
	name := d.Get(mkName).(string)
	tags := d.Get(mkTags).([]interface{})

	networkDeviceObjects, err := network.GetNetworkDeviceObjects(d)
	if err != nil {
		return diag.FromErr(err)
	}

	nodeName := d.Get(mkNodeName).(string)

	operatingSystem, err := structure.GetSchemaBlock(
		resource,
		d,
		[]string{mkOperatingSystem},
		0,
		true,
	)
	if err != nil {
		return diag.FromErr(err)
	}

	operatingSystemType := operatingSystem[mkOperatingSystemType].(string)

	poolID := d.Get(mkPoolID).(string)
	protection := types.CustomBool(d.Get(mkProtection).(bool))

	serialDevices := vmGetSerialDeviceList(d)

	smbios := vmGetSMBIOS(d)

	startupOrder := vmGetStartupOrder(d)

	onBoot := types.CustomBool(d.Get(mkOnBoot).(bool))
	tabletDevice := types.CustomBool(d.Get(mkTabletDevice).(bool))
	template := types.CustomBool(d.Get(mkTemplate).(bool))

	vgaDevice, err := vmGetVGADeviceObject(d)
	if err != nil {
		return diag.FromErr(err)
	}

	vmIDUntyped, hasVMID := d.GetOk(mkVMID)
	vmID := vmIDUntyped.(int)

	if !hasVMID {
		vmIDNew, e := client.Cluster().GetVMID(ctx)
		if e != nil {
			return diag.FromErr(e)
		}

		vmID = *vmIDNew
		e = d.Set(mkVMID, vmID)

		if e != nil {
			return diag.FromErr(e)
		}
	}

	diskDeviceObjects, err := disk.GetDiskDeviceObjects(d, resource, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	virtioDeviceObjects := diskDeviceObjects["virtio"]
	scsiDeviceObjects := diskDeviceObjects["scsi"]
	ideDeviceObjects := diskDeviceObjects["ide"]
	sataDeviceObjects := diskDeviceObjects["sata"]

	var bootOrderConverted []string
	if cdromEnabled {
		bootOrderConverted = []string{cdromInterface}
	}

	bootOrder := d.Get(mkBootOrder).([]interface{})

	if len(bootOrder) == 0 {
		if ideDeviceObjects != nil {
			bootOrderConverted = append(bootOrderConverted, "ide0")
		}

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
	ideDevices := vms.CustomStorageDevices{}

	if cdromCloudInitInterface != "" {
		ideDevices[cdromCloudInitInterface] = &vms.CustomStorageDevice{
			Enabled:    cdromCloudInitEnabled,
			FileVolume: cdromCloudInitFileID,
			Media:      &ideDevice2Media,
		}
	}

	if cdromInterface != "" {
		ideDevices[cdromInterface] = &vms.CustomStorageDevice{
			Enabled:    cdromEnabled,
			FileVolume: cdromFileID,
			Media:      &ideDevice2Media,
		}
	}

	var memorySharedObject *vms.CustomSharedMemory

	if memoryShared > 0 {
		memorySharedName := fmt.Sprintf("vm-%d-ivshmem", vmID)
		memorySharedObject = &vms.CustomSharedMemory{
			Name: &memorySharedName,
			Size: memoryShared,
		}
	}

	scsiHardware := d.Get(mkSCSIHardware).(string)

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
		CPUCores:        ptr.Ptr(int64(cpuCores)),
		CPUEmulation: &vms.CustomCPUEmulation{
			Flags: &cpuFlagsConverted,
			Type:  cpuType,
		},
		CPUSockets:          ptr.Ptr(int64(cpuSockets)),
		CPUUnits:            ptr.Ptr(int64(cpuUnits)),
		DedicatedMemory:     &memoryDedicated,
		DeletionProtection:  &protection,
		EFIDisk:             efiDisk,
		TPMState:            tpmState,
		FloatingMemory:      &memoryFloating,
		IDEDevices:          ideDevices,
		KeyboardLayout:      &keyboardLayout,
		NetworkDevices:      networkDeviceObjects,
		NUMAEnabled:         &cpuNUMA,
		NUMADevices:         numaDeviceObjects,
		OSType:              &operatingSystemType,
		PCIDevices:          pciDeviceObjects,
		SCSIHardware:        &scsiHardware,
		SerialDevices:       serialDevices,
		SharedMemory:        memorySharedObject,
		StartOnBoot:         &onBoot,
		SMBIOS:              smbios,
		StartupOrder:        startupOrder,
		TabletDeviceEnabled: &tabletDevice,
		Template:            &template,
		USBDevices:          usbDeviceObjects,
		VGADevice:           vgaDevice,
		VMID:                vmID,
	}

	if ideDeviceObjects != nil {
		createBody.IDEDevices = ideDeviceObjects
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
	if client.API().IsRootTicket() ||
		cpuArchitecture != dvCPUArchitecture {
		createBody.CPUArchitecture = &cpuArchitecture
	}

	if cpuHotplugged > 0 {
		createBody.VirtualCPUCount = ptr.Ptr(int64(cpuHotplugged))
	}

	if cpuLimit > 0 {
		createBody.CPULimit = ptr.Ptr(int64(cpuLimit))
	}

	if cpuAffinity != "" {
		createBody.CPUAffinity = &cpuAffinity
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

	if memoryHugepages != "" {
		createBody.Hugepages = &memoryHugepages
	}

	if memoryKeepHugepages {
		createBody.KeepHugepages = &memoryKeepHugepages
	}

	if name != "" {
		createBody.Name = &name
	}

	if poolID != "" {
		createBody.PoolID = &poolID
	}

	hookScript := d.Get(mkHookScriptFileID).(string)
	if len(hookScript) > 0 {
		createBody.HookScript = &hookScript
	}

	err = client.Node(nodeName).VM(0).CreateVM(ctx, createBody)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(vmID))

	diags := disk.CreateCustomDisks(ctx, client, nodeName, vmID, diskDeviceObjects)
	if diags.HasError() {
		return diags
	}

	return vmCreateStart(ctx, d, m)
}

func vmCreateStart(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	started := d.Get(mkStarted).(bool)
	template := d.Get(mkTemplate).(bool)
	reboot := d.Get(mkRebootAfterCreation).(bool)

	if !started || template {
		return vmRead(ctx, d, m)
	}

	config := m.(proxmoxtf.ProviderConfiguration)

	client, err := config.GetClient()
	if err != nil {
		return diag.FromErr(err)
	}

	nodeName := d.Get(mkNodeName).(string)

	vmID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	vmAPI := client.Node(nodeName).VM(vmID)

	// Start the virtual machine and wait for it to reach a running state before continuing.
	if diags := vmStart(ctx, vmAPI, d); diags != nil {
		return diags
	}

	if reboot {
		rebootTimeoutSec := d.Get(mkTimeoutReboot).(int)

		err := vmAPI.RebootVM(
			ctx,
			&vms.RebootRequestBody{
				Timeout: &rebootTimeoutSec,
			},
		)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return vmRead(ctx, d, m)
}

func vmGetAudioDeviceList(d *schema.ResourceData) vms.CustomAudioDevices {
	devices := d.Get(mkAudioDevice).([]interface{})
	list := make(vms.CustomAudioDevices, len(devices))

	for i, v := range devices {
		block := v.(map[string]interface{})

		device, _ := block[mkAudioDeviceDevice].(string)
		driver, _ := block[mkAudioDeviceDriver].(string)
		enabled, _ := block[mkAudioDeviceEnabled].(bool)

		list[i].Device = device
		list[i].Driver = &driver
		list[i].Enabled = enabled
	}

	return list
}

func vmGetCloudInitConfig(d *schema.ResourceData) *vms.CustomCloudInitConfig {
	initialization := d.Get(mkInitialization).([]interface{})

	if len(initialization) == 0 || initialization[0] == nil {
		return nil
	}

	var initializationConfig *vms.CustomCloudInitConfig

	initializationBlock := initialization[0].(map[string]interface{})
	initializationConfig = &vms.CustomCloudInitConfig{}
	initializationDNS := initializationBlock[mkInitializationDNS].([]interface{})

	if len(initializationDNS) > 0 && initializationDNS[0] != nil {
		initializationDNSBlock := initializationDNS[0].(map[string]interface{})
		domain := initializationDNSBlock[mkInitializationDNSDomain].(string)

		if domain != "" {
			initializationConfig.SearchDomain = &domain
		}

		servers := initializationDNSBlock[mkInitializationDNSServers].([]interface{})
		deprecatedServer := initializationDNSBlock[mkInitializationDNSServer].(string)

		if len(servers) > 0 {
			nameserver := strings.Join(ConvertToStringSlice(servers), " ")

			initializationConfig.Nameserver = &nameserver
		} else if deprecatedServer != "" {
			initializationConfig.Nameserver = &deprecatedServer
		}
	}

	initializationIPConfig := initializationBlock[mkInitializationIPConfig].([]interface{})
	initializationConfig.IPConfig = make([]vms.CustomCloudInitIPConfig, len(initializationIPConfig))

	for i, c := range initializationIPConfig {
		configBlock := c.(map[string]interface{})
		ipv4 := configBlock[mkInitializationIPConfigIPv4].([]interface{})

		if len(ipv4) > 0 && ipv4[0] != nil {
			ipv4Block := ipv4[0].(map[string]interface{})
			ipv4Address := ipv4Block[mkInitializationIPConfigIPv4Address].(string)

			if ipv4Address != "" {
				initializationConfig.IPConfig[i].IPv4 = &ipv4Address
			}

			ipv4Gateway := ipv4Block[mkInitializationIPConfigIPv4Gateway].(string)

			if ipv4Gateway != "" {
				initializationConfig.IPConfig[i].GatewayIPv4 = &ipv4Gateway
			}
		}

		ipv6 := configBlock[mkInitializationIPConfigIPv6].([]interface{})

		if len(ipv6) > 0 && ipv6[0] != nil {
			ipv6Block := ipv6[0].(map[string]interface{})
			ipv6Address := ipv6Block[mkInitializationIPConfigIPv6Address].(string)

			if ipv6Address != "" {
				initializationConfig.IPConfig[i].IPv6 = &ipv6Address
			}

			ipv6Gateway := ipv6Block[mkInitializationIPConfigIPv6Gateway].(string)

			if ipv6Gateway != "" {
				initializationConfig.IPConfig[i].GatewayIPv6 = &ipv6Gateway
			}
		}
	}

	initializationUserAccount := initializationBlock[mkInitializationUserAccount].([]interface{})

	if len(initializationUserAccount) > 0 && initializationUserAccount[0] != nil {
		initializationUserAccountBlock := initializationUserAccount[0].(map[string]interface{})
		keys := initializationUserAccountBlock[mkInitializationUserAccountKeys].([]interface{})

		if len(keys) > 0 {
			sshKeys := make(vms.CustomCloudInitSSHKeys, len(keys))

			for i, k := range keys {
				if k != nil {
					sshKeys[i] = k.(string)
				}
			}

			initializationConfig.SSHKeys = &sshKeys
		}

		password := initializationUserAccountBlock[mkInitializationUserAccountPassword].(string)
		if password != "" {
			initializationConfig.Password = &password
		}

		username := initializationUserAccountBlock[mkInitializationUserAccountUsername].(string)
		initializationConfig.Username = &username
	}

	initializationUserDataFileID := initializationBlock[mkInitializationUserDataFileID].(string)
	if initializationUserDataFileID != "" {
		initializationConfig.Files = &vms.CustomCloudInitFiles{
			UserVolume: &initializationUserDataFileID,
		}
	}

	initializationVendorDataFileID := initializationBlock[mkInitializationVendorDataFileID].(string)
	if initializationVendorDataFileID != "" {
		if initializationConfig.Files == nil {
			initializationConfig.Files = &vms.CustomCloudInitFiles{}
		}

		initializationConfig.Files.VendorVolume = &initializationVendorDataFileID
	}

	initializationNetworkDataFileID := initializationBlock[mkInitializationNetworkDataFileID].(string)
	if initializationNetworkDataFileID != "" {
		if initializationConfig.Files == nil {
			initializationConfig.Files = &vms.CustomCloudInitFiles{}
		}

		initializationConfig.Files.NetworkVolume = &initializationNetworkDataFileID
	}

	initializationMetaDataFileID := initializationBlock[mkInitializationMetaDataFileID].(string)
	if initializationMetaDataFileID != "" {
		if initializationConfig.Files == nil {
			initializationConfig.Files = &vms.CustomCloudInitFiles{}
		}

		initializationConfig.Files.MetaVolume = &initializationMetaDataFileID
	}

	initializationType := initializationBlock[mkInitializationType].(string)
	if initializationType != "" {
		initializationConfig.Type = &initializationType
	}

	return initializationConfig
}

func vmGetEfiDisk(d *schema.ResourceData, disk []interface{}) *vms.CustomEFIDisk {
	var efiDisk []interface{}

	if disk != nil {
		efiDisk = disk
	} else {
		efiDisk = d.Get(mkEFIDisk).([]interface{})
	}

	var efiDiskConfig *vms.CustomEFIDisk

	if len(efiDisk) > 0 && efiDisk[0] != nil {
		efiDiskConfig = &vms.CustomEFIDisk{}

		block := efiDisk[0].(map[string]interface{})
		datastoreID, _ := block[mkEFIDiskDatastoreID].(string)
		fileFormat, _ := block[mkEFIDiskFileFormat].(string)
		efiType, _ := block[mkEFIDiskType].(string)
		preEnrolledKeys := types.CustomBool(block[mkEFIDiskPreEnrolledKeys].(bool))

		// use the special syntax STORAGE_ID:SIZE_IN_GiB to allocate a new volume.
		// NB SIZE_IN_GiB is ignored, see docs for more info.
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
			Enabled:     true,
			FileVolume:  efiDisk.FileVolume,
			Format:      efiDisk.Format,
			Interface:   &diskInterface,
			DatastoreID: &id,
		}

		if efiDisk.Type != nil {
			ds, err := types.ParseDiskSize(*efiDisk.Type)
			if err != nil {
				return nil, fmt.Errorf("invalid efi disk type: %s", err.Error())
			}

			storageDevice.Size = &ds
		}
	}

	return storageDevice, nil
}

func vmGetTPMState(d *schema.ResourceData, disk []interface{}) *vms.CustomTPMState {
	var tpmState []interface{}

	if disk != nil {
		tpmState = disk
	} else {
		tpmState = d.Get(mkTPMState).([]interface{})
	}

	var tpmStateConfig *vms.CustomTPMState

	if len(tpmState) > 0 && tpmState[0] != nil {
		tpmStateConfig = &vms.CustomTPMState{}

		block := tpmState[0].(map[string]interface{})
		datastoreID, _ := block[mkTPMStateDatastoreID].(string)
		version, _ := block[mkTPMStateVersion].(string)

		// use the special syntax STORAGE_ID:SIZE_IN_GiB to allocate a new volume.
		// NB SIZE_IN_GiB is ignored, see docs for more info.
		tpmStateConfig.FileVolume = fmt.Sprintf("%s:1", datastoreID)
		tpmStateConfig.Version = &version
	}

	return tpmStateConfig
}

func vmGetTPMStateAsStorageDevice(d *schema.ResourceData, disk []interface{}) *vms.CustomStorageDevice {
	tpmState := vmGetTPMState(d, disk)

	var storageDevice *vms.CustomStorageDevice

	if tpmState != nil {
		id := "0"
		baseDiskInterface := "tpmstate"
		diskInterface := fmt.Sprint(baseDiskInterface, id)

		storageDevice = &vms.CustomStorageDevice{
			Enabled:     true,
			FileVolume:  tpmState.FileVolume,
			Interface:   &diskInterface,
			DatastoreID: &id,
		}
	}

	return storageDevice
}

func vmGetHostPCIDeviceObjects(d *schema.ResourceData) vms.CustomPCIDevices {
	pciDevice := d.Get(mkHostPCI).([]interface{})
	pciDeviceObjects := make(vms.CustomPCIDevices, len(pciDevice))

	for i, pciDeviceEntry := range pciDevice {
		block := pciDeviceEntry.(map[string]interface{})

		ids, _ := block[mkHostPCIDeviceID].(string)
		mdev, _ := block[mkHostPCIDeviceMDev].(string)
		pcie := types.CustomBool(block[mkHostPCIDevicePCIE].(bool))
		rombar := types.CustomBool(
			block[mkHostPCIDeviceROMBAR].(bool),
		)
		romfile, _ := block[mkHostPCIDeviceROMFile].(string)
		xvga := types.CustomBool(block[mkHostPCIDeviceXVGA].(bool))
		mapping, _ := block[mkHostPCIDeviceMapping].(string)

		device := vms.CustomPCIDevice{
			PCIExpress: &pcie,
			ROMBAR:     &rombar,
			XVGA:       &xvga,
		}

		if ids != "" {
			dIDs := strings.Split(ids, ";")
			device.DeviceIDs = &dIDs
		}

		if mdev != "" {
			device.MDev = &mdev
		}

		if romfile != "" {
			device.ROMFile = &romfile
		}

		if mapping != "" {
			device.Mapping = &mapping
		}

		pciDeviceObjects[i] = device
	}

	return pciDeviceObjects
}

func vmGetNumaDeviceObjects(d *schema.ResourceData) vms.CustomNUMADevices {
	numaNode := d.Get(mkNUMA).([]interface{})
	numaNodeObjects := make(vms.CustomNUMADevices, len(numaNode))

	for i, numaNodeEntry := range numaNode {
		block := numaNodeEntry.(map[string]interface{})

		deviceName := block[mkNUMADevice].(string)
		ids := block[mkNUMACPUIDs].(string)
		hostNodes, _ := block[mkNUMAHostNodeNames].(string)
		memory, _ := block[mkNUMAMemory].(int)
		policy, _ := block[mkNUMAPolicy].(string)

		device := vms.CustomNUMADevice{
			Memory: &memory,
			Policy: &policy,
		}

		if ids != "" {
			dIDs := strings.Split(ids, ";")
			device.CPUIDs = dIDs
		}

		if hostNodes != "" {
			dHostNodes := strings.Split(hostNodes, ";")
			device.HostNodeNames = &dHostNodes
		}

		if strings.HasPrefix(deviceName, "numa") {
			deviceID, err := strconv.Atoi(deviceName[4:])
			if err == nil {
				numaNodeObjects[deviceID] = device

				continue
			}
		}

		numaNodeObjects[i] = device
	}

	return numaNodeObjects
}

func vmGetHostUSBDeviceObjects(d *schema.ResourceData) vms.CustomUSBDevices {
	usbDevice := d.Get(mkHostUSB).([]interface{})
	usbDeviceObjects := make(vms.CustomUSBDevices, len(usbDevice))

	for i, usbDeviceEntry := range usbDevice {
		block := usbDeviceEntry.(map[string]interface{})

		host, _ := block[mkHostUSBDevice].(string)
		usb3 := types.CustomBool(block[mkHostUSBDeviceUSB3].(bool))
		mapping, _ := block[mkHostUSBDeviceMapping].(string)

		device := vms.CustomUSBDevice{
			USB3: &usb3,
		}

		if host != "" {
			device.HostDevice = &host
		}

		if mapping != "" {
			device.Mapping = &mapping
		}

		usbDeviceObjects[i] = device
	}

	return usbDeviceObjects
}

func vmGetSerialDeviceList(d *schema.ResourceData) vms.CustomSerialDevices {
	device := d.Get(mkSerialDevice).([]interface{})
	list := make(vms.CustomSerialDevices, len(device))

	for i, v := range device {
		block := v.(map[string]interface{})

		device, _ := block[mkSerialDeviceDevice].(string)

		list[i] = device
	}

	return list
}

func vmGetSMBIOS(d *schema.ResourceData) *vms.CustomSMBIOS {
	smbiosSections := d.Get(mkSMBIOS).([]interface{})

	if len(smbiosSections) > 0 && smbiosSections[0] != nil {
		smbiosBlock := smbiosSections[0].(map[string]interface{})
		b64 := types.CustomBool(true)
		family, _ := smbiosBlock[mkSMBIOSFamily].(string)
		manufacturer, _ := smbiosBlock[mkSMBIOSManufacturer].(string)
		product, _ := smbiosBlock[mkSMBIOSProduct].(string)
		serial, _ := smbiosBlock[mkSMBIOSSerial].(string)
		sku, _ := smbiosBlock[mkSMBIOS].(string)
		version, _ := smbiosBlock[mkSMBIOSVersion].(string)
		uid, _ := smbiosBlock[mkSMBIOSUUID].(string)

		smbios := vms.CustomSMBIOS{
			Base64: &b64,
		}

		if family != "" {
			v := base64.StdEncoding.EncodeToString([]byte(family))
			smbios.Family = &v
		}

		if manufacturer != "" {
			v := base64.StdEncoding.EncodeToString([]byte(manufacturer))
			smbios.Manufacturer = &v
		}

		if product != "" {
			v := base64.StdEncoding.EncodeToString([]byte(product))
			smbios.Product = &v
		}

		if serial != "" {
			v := base64.StdEncoding.EncodeToString([]byte(serial))
			smbios.Serial = &v
		}

		if sku != "" {
			v := base64.StdEncoding.EncodeToString([]byte(sku))
			smbios.SKU = &v
		}

		if version != "" {
			v := base64.StdEncoding.EncodeToString([]byte(version))
			smbios.Version = &v
		}

		if uid != "" {
			smbios.UUID = &uid
		}

		if smbios.UUID == nil || *smbios.UUID == "" {
			smbios.UUID = ptr.Ptr(uuid.New().String())
		}

		return &smbios
	}

	return nil
}

func vmGetStartupOrder(d *schema.ResourceData) *vms.CustomStartupOrder {
	startup := d.Get(mkStartup).([]interface{})

	if len(startup) > 0 && startup[0] != nil {
		startupBlock := startup[0].(map[string]interface{})
		startupOrder := startupBlock[mkStartupOrder].(int)
		startupUpDelay := startupBlock[mkStartupUpDelay].(int)
		startupDownDelay := startupBlock[mkStartupDownDelay].(int)

		order := vms.CustomStartupOrder{}

		if startupUpDelay >= 0 {
			order.Up = &startupUpDelay
		}

		if startupDownDelay >= 0 {
			order.Down = &startupDownDelay
		}

		if startupOrder >= 0 {
			order.Order = &startupOrder
		}

		return &order
	}

	return nil
}

func vmGetTagsString(d *schema.ResourceData) string {
	var sanitizedTags []string

	tags := d.Get(mkTags).([]interface{})
	for _, tag := range tags {
		sanitizedTag := strings.TrimSpace(tag.(string))
		if len(sanitizedTag) > 0 {
			sanitizedTags = append(sanitizedTags, sanitizedTag)
		}
	}

	sort.Strings(sanitizedTags)

	return strings.Join(sanitizedTags, ";")
}

func vmGetVGADeviceObject(d *schema.ResourceData) (*vms.CustomVGADevice, error) {
	resource := VM()

	vgaBlock, err := structure.GetSchemaBlock(
		resource,
		d,
		[]string{mkVGA},
		0,
		true,
	)
	if err != nil {
		return nil, fmt.Errorf("error getting VGA block: %w", err)
	}

	vgaClipboard := vgaBlock[mkVGAClipboard].(string)
	vgaMemory := vgaBlock[mkVGAMemory].(int)
	vgaType := vgaBlock[mkVGAType].(string)

	vgaDevice := &vms.CustomVGADevice{}

	if vgaClipboard != "" {
		vgaDevice.Clipboard = &vgaClipboard
	}

	if vgaMemory > 0 {
		vgaDevice.Memory = ptr.Ptr(int64(vgaMemory))
	}

	if vgaType != "" {
		vgaDevice.Type = &vgaType
	}

	return vgaDevice, nil
}

func vmRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(proxmoxtf.ProviderConfiguration)

	client, err := config.GetClient()
	if err != nil {
		return diag.FromErr(err)
	}

	vmID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	vmNodeName, err := client.Cluster().GetVMNodeName(ctx, vmID)
	if err != nil {
		if errors.Is(err, cluster.ErrVMDoesNotExist) {
			d.SetId("")

			return nil
		}

		return diag.FromErr(err)
	}

	if vmNodeName != d.Get(mkNodeName) {
		err = d.Set(mkNodeName, vmNodeName)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	nodeName := d.Get(mkNodeName).(string)

	vmAPI := client.Node(nodeName).VM(vmID)

	// Retrieve the entire configuration in order to compare it to the state.
	vmConfig, err := vmAPI.GetVM(ctx)
	if err != nil {
		if errors.Is(err, api.ErrResourceDoesNotExist) {
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

	client, e := config.GetClient()
	if e != nil {
		return diag.FromErr(e)
	}

	diags := vmReadPrimitiveValues(d, vmConfig, vmStatus)
	if diags.HasError() {
		return diags
	}

	// Fix terraform.tfstate, by replacing '-1' (the old default value) with actual vm_id value
	if storedVMID := d.Get(mkVMID).(int); storedVMID == -1 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary: fmt.Sprintf("VM %s has stored legacy vm_id %d, setting vm_id to its correct value %d.",
				d.Id(), storedVMID, vmID),
		})

		err := d.Set(mkVMID, vmID)
		diags = append(diags, diag.FromErr(err)...)
	}

	nodeName := d.Get(mkNodeName).(string)
	clone := d.Get(mkClone).([]interface{})

	// Compare the agent configuration to the one stored in the state.
	currentAgent := d.Get(mkAgent).([]interface{})

	//nolint:gocritic
	if len(clone) == 0 || len(currentAgent) > 0 {
		if vmConfig.Agent != nil {
			agent := map[string]interface{}{}

			if vmConfig.Agent.Enabled != nil {
				agent[mkAgentEnabled] = bool(*vmConfig.Agent.Enabled)
			} else {
				agent[mkAgentEnabled] = false
			}

			if vmConfig.Agent.TrimClonedDisks != nil {
				agent[mkAgentTrim] = bool(
					*vmConfig.Agent.TrimClonedDisks,
				)
			} else {
				agent[mkAgentTrim] = false
			}

			if len(currentAgent) > 0 && currentAgent[0] != nil {
				currentAgentBlock := currentAgent[0].(map[string]interface{})
				currentAgentTimeout := currentAgentBlock[mkAgentTimeout].(string)

				if currentAgentTimeout != "" {
					agent[mkAgentTimeout] = currentAgentTimeout
				} else {
					agent[mkAgentTimeout] = dvAgentTimeout
				}
			} else {
				agent[mkAgentTimeout] = dvAgentTimeout
			}

			if vmConfig.Agent.Type != nil {
				agent[mkAgentType] = *vmConfig.Agent.Type
			} else {
				agent[mkAgentType] = ""
			}

			if len(clone) > 0 {
				if len(currentAgent) > 0 {
					err := d.Set(mkAgent, []interface{}{agent})
					diags = append(diags, diag.FromErr(err)...)
				}
			} else if len(currentAgent) > 0 ||
				agent[mkAgentEnabled] != dvAgentEnabled ||
				agent[mkAgentTimeout] != dvAgentTimeout ||
				agent[mkAgentTrim] != dvAgentTrim ||
				agent[mkAgentType] != dvAgentType {
				err := d.Set(mkAgent, []interface{}{agent})
				diags = append(diags, diag.FromErr(err)...)
			}
		} else if len(clone) > 0 {
			if len(currentAgent) > 0 {
				err := d.Set(mkAgent, []interface{}{})
				diags = append(diags, diag.FromErr(err)...)
			}
		} else {
			err := d.Set(mkAgent, []interface{}{})
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	// Compare the audio devices to those stored in the state.
	currentAudioDevice := d.Get(mkAudioDevice).([]interface{})

	audioDevices := make([]interface{}, 1)
	audioDevicesArray := []*vms.CustomAudioDevice{
		vmConfig.AudioDevice,
	}
	audioDevicesCount := 0

	for adi, ad := range audioDevicesArray {
		m := map[string]interface{}{}

		if ad != nil {
			m[mkAudioDeviceDevice] = ad.Device

			if ad.Driver != nil {
				m[mkAudioDeviceDriver] = *ad.Driver
			} else {
				m[mkAudioDeviceDriver] = ""
			}

			m[mkAudioDeviceEnabled] = true

			audioDevicesCount = adi + 1
		} else {
			m[mkAudioDeviceDevice] = ""
			m[mkAudioDeviceDriver] = ""
			m[mkAudioDeviceEnabled] = false
		}

		audioDevices[adi] = m
	}

	if len(clone) == 0 || len(currentAudioDevice) > 0 {
		err := d.Set(mkAudioDevice, audioDevices[:audioDevicesCount])
		diags = append(diags, diag.FromErr(err)...)
	}

	// Compare the IDE devices to the CD-ROM configurations stored in the state.
	currentInterface := dvCDROMInterface

	currentCDROM := d.Get(mkCDROM).([]interface{})
	if len(currentCDROM) > 0 && currentCDROM[0] != nil {
		currentBlock := currentCDROM[0].(map[string]interface{})
		currentInterface = currentBlock[mkCDROMInterface].(string)
	}

	cdromIDEDevice := getStorageDevice(vmConfig, currentInterface)

	if cdromIDEDevice != nil {
		cdrom := make([]interface{}, 1)
		cdromBlock := map[string]interface{}{}

		if len(clone) == 0 || len(currentCDROM) > 0 {
			cdromBlock[mkCDROMEnabled] = cdromIDEDevice.Enabled
			cdromBlock[mkCDROMFileID] = cdromIDEDevice.FileVolume
			cdromBlock[mkCDROMInterface] = currentInterface

			if len(currentCDROM) > 0 && currentCDROM[0] != nil {
				currentBlock := currentCDROM[0].(map[string]interface{})

				if currentBlock[mkCDROMFileID] == "" {
					cdromBlock[mkCDROMFileID] = ""
				}

				if currentBlock[mkCDROMEnabled] == false {
					cdromBlock[mkCDROMEnabled] = false
				}
			}

			cdrom[0] = cdromBlock

			err := d.Set(mkCDROM, cdrom)
			diags = append(diags, diag.FromErr(err)...)
		}
	} else {
		err := d.Set(mkCDROM, []interface{}{})
		diags = append(diags, diag.FromErr(err)...)
	}

	// Compare the CPU configuration to the one stored in the state.
	cpu := map[string]interface{}{}

	if vmConfig.CPUArchitecture != nil {
		cpu[mkCPUArchitecture] = *vmConfig.CPUArchitecture
	} else {
		// Default value of "arch" is "" according to the API documentation.
		// However, assume the provider's default value as a workaround when the root account is not being used.
		if !client.API().IsRootTicket() {
			cpu[mkCPUArchitecture] = dvCPUArchitecture
		} else {
			cpu[mkCPUArchitecture] = ""
		}
	}

	if vmConfig.CPUCores != nil {
		cpu[mkCPUCores] = int(*vmConfig.CPUCores)
	} else {
		// Default value of "cores" is "1" according to the API documentation.
		cpu[mkCPUCores] = 1
	}

	if vmConfig.VirtualCPUCount != nil {
		cpu[mkCPUHotplugged] = int(*vmConfig.VirtualCPUCount)
	} else {
		// Default value of "vcpus" is "1" according to the API documentation.
		cpu[mkCPUHotplugged] = 0
	}

	if vmConfig.CPULimit != nil {
		cpu[mkCPULimit] = int(*vmConfig.CPULimit)
	} else {
		// Default value of "cpulimit" is "0" according to the API documentation.
		cpu[mkCPULimit] = 0
	}

	if vmConfig.NUMAEnabled != nil {
		cpu[mkCPUNUMA] = *vmConfig.NUMAEnabled
	} else {
		// Default value of "numa" is "false" according to the API documentation.
		cpu[mkCPUNUMA] = false
	}

	currentNUMAList := d.Get(mkNUMA).([]interface{})
	numaMap := map[string]interface{}{}

	numaDevices := getNUMAInfo(vmConfig, d)
	for ni, np := range numaDevices {
		if np == nil || np.CPUIDs == nil || np.HostNodeNames == nil {
			continue
		}

		numaNode := map[string]interface{}{}
		numaNode[mkNUMADevice] = ni

		if len(np.CPUIDs) > 0 {
			numaNode[mkNUMACPUIDs] = strings.Join(np.CPUIDs, ";")
		}

		numaNode[mkNUMAHostNodeNames] = strings.Join(*np.HostNodeNames, ";")
		numaNode[mkNUMAMemory] = np.Memory
		numaNode[mkNUMAPolicy] = np.Policy

		numaMap[ni] = numaNode
	}

	if len(clone) == 0 || len(currentNUMAList) > 0 {
		var numaList []interface{}

		if len(currentNUMAList) > 0 {
			devices := utils.ListResourcesAttributeValue(currentNUMAList, mkNUMADevice)
			numaList = utils.OrderedListFromMapByKeyValues(numaMap, devices)
		} else {
			numaList = utils.OrderedListFromMap(numaMap)
		}

		err := d.Set(mkNUMA, numaList)
		diags = append(diags, diag.FromErr(err)...)
	}

	if vmConfig.CPUSockets != nil {
		cpu[mkCPUSockets] = int(*vmConfig.CPUSockets)
	} else {
		// Default value of "sockets" is "1" according to the API documentation.
		cpu[mkCPUSockets] = 1
	}

	if vmConfig.CPUEmulation != nil {
		if vmConfig.CPUEmulation.Flags != nil {
			convertedFlags := make([]interface{}, len(*vmConfig.CPUEmulation.Flags))

			for fi, fv := range *vmConfig.CPUEmulation.Flags {
				convertedFlags[fi] = fv
			}

			cpu[mkCPUFlags] = convertedFlags
		} else {
			cpu[mkCPUFlags] = []interface{}{}
		}

		cpu[mkCPUType] = vmConfig.CPUEmulation.Type
	} else {
		cpu[mkCPUFlags] = []interface{}{}
		// Default value of "cputype" is "qemu64" according to the QEMU documentation.
		cpu[mkCPUType] = "qemu64"
	}

	if vmConfig.CPUUnits != nil {
		cpu[mkCPUUnits] = int(*vmConfig.CPUUnits)
	} else {
		// Default value of "cpuunits" is "1024" according to the API documentation.
		cpu[mkCPUUnits] = 1024
	}

	if vmConfig.CPUAffinity != nil {
		cpu[mkCPUAffinity] = *vmConfig.CPUAffinity
	} else {
		cpu[mkCPUAffinity] = ""
	}

	currentCPU := d.Get(mkCPU).([]interface{})

	if len(clone) > 0 {
		if len(currentCPU) > 0 {
			err := d.Set(mkCPU, []interface{}{cpu})
			diags = append(diags, diag.FromErr(err)...)
		}
	} else if len(currentCPU) > 0 ||
		cpu[mkCPUArchitecture] != dvCPUArchitecture ||
		cpu[mkCPUCores] != dvCPUCores ||
		len(cpu[mkCPUFlags].([]interface{})) > 0 ||
		cpu[mkCPUHotplugged] != dvCPUHotplugged ||
		cpu[mkCPULimit] != dvCPULimit ||
		cpu[mkCPUSockets] != dvCPUSockets ||
		cpu[mkCPUType] != dvCPUType ||
		cpu[mkCPUUnits] != dvCPUUnits {
		err := d.Set(mkCPU, []interface{}{cpu})
		diags = append(diags, diag.FromErr(err)...)
	}

	allDiskInfo := disk.GetInfo(vmConfig, d)

	diags = append(diags, disk.Read(ctx, d, allDiskInfo, vmID, client, nodeName, len(clone) > 0)...)

	if vmConfig.EFIDisk != nil {
		efiDisk := map[string]interface{}{}

		fileIDParts := strings.Split(vmConfig.EFIDisk.FileVolume, ":")

		efiDisk[mkEFIDiskDatastoreID] = fileIDParts[0]

		if vmConfig.EFIDisk.Format != nil {
			efiDisk[mkEFIDiskFileFormat] = *vmConfig.EFIDisk.Format
		} else {
			// disk format may not be returned by config API if it is default for the storage, and that may be different
			// from the default qcow2, so we need to read it from the storage API to make sure we have the correct value
			volume, err := client.Node(nodeName).Storage(fileIDParts[0]).GetDatastoreFile(ctx, vmConfig.EFIDisk.FileVolume)
			if err != nil {
				diags = append(diags, diag.FromErr(e)...)
			} else {
				efiDisk[mkEFIDiskFileFormat] = volume.FileFormat
			}
		}

		if vmConfig.EFIDisk.Type != nil {
			efiDisk[mkEFIDiskType] = *vmConfig.EFIDisk.Type
		} else {
			efiDisk[mkEFIDiskType] = dvEFIDiskType
		}

		if vmConfig.EFIDisk.PreEnrolledKeys != nil {
			efiDisk[mkEFIDiskPreEnrolledKeys] = *vmConfig.EFIDisk.PreEnrolledKeys
		} else {
			efiDisk[mkEFIDiskPreEnrolledKeys] = false
		}

		currentEfiDisk := d.Get(mkEFIDisk).([]interface{})

		if len(clone) > 0 {
			if len(currentEfiDisk) > 0 {
				err := d.Set(mkEFIDisk, []interface{}{efiDisk})
				diags = append(diags, diag.FromErr(err)...)
			}
		} else if len(currentEfiDisk) > 0 ||
			efiDisk[mkEFIDiskDatastoreID] != dvEFIDiskDatastoreID ||
			efiDisk[mkEFIDiskType] != dvEFIDiskType ||
			efiDisk[mkEFIDiskPreEnrolledKeys] != dvEFIDiskPreEnrolledKeys ||
			efiDisk[mkEFIDiskFileFormat] != dvEFIDiskFileFormat {
			err := d.Set(mkEFIDisk, []interface{}{efiDisk})
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	if vmConfig.TPMState != nil {
		tpmState := map[string]interface{}{}

		fileIDParts := strings.Split(vmConfig.TPMState.FileVolume, ":")

		tpmState[mkTPMStateDatastoreID] = fileIDParts[0]
		tpmState[mkTPMStateVersion] = dvTPMStateVersion

		currentTPMState := d.Get(mkTPMState).([]interface{})

		if len(clone) > 0 {
			if len(currentTPMState) > 0 {
				err := d.Set(mkTPMState, []interface{}{tpmState})
				diags = append(diags, diag.FromErr(err)...)
			}
		} else if len(currentTPMState) > 0 ||
			tpmState[mkTPMStateDatastoreID] != dvTPMStateDatastoreID ||
			tpmState[mkTPMStateVersion] != dvTPMStateVersion {
			err := d.Set(mkTPMState, []interface{}{tpmState})
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	currentPCIList := d.Get(mkHostPCI).([]interface{})
	pciMap := map[string]interface{}{}

	pciDevices := getPCIInfo(vmConfig, d)
	for pi, pp := range pciDevices {
		if (pp == nil) || (pp.DeviceIDs == nil && pp.Mapping == nil) {
			continue
		}

		pci := map[string]interface{}{}

		pci[mkHostPCIDevice] = pi
		if pp.DeviceIDs != nil {
			pci[mkHostPCIDeviceID] = strings.Join(*pp.DeviceIDs, ";")
		} else {
			pci[mkHostPCIDeviceID] = ""
		}

		if pp.MDev != nil {
			pci[mkHostPCIDeviceMDev] = *pp.MDev
		} else {
			pci[mkHostPCIDeviceMDev] = ""
		}

		if pp.PCIExpress != nil {
			pci[mkHostPCIDevicePCIE] = *pp.PCIExpress
		} else {
			pci[mkHostPCIDevicePCIE] = false
		}

		if pp.ROMBAR != nil {
			pci[mkHostPCIDeviceROMBAR] = *pp.ROMBAR
		} else {
			pci[mkHostPCIDeviceROMBAR] = true
		}

		if pp.ROMFile != nil {
			pci[mkHostPCIDeviceROMFile] = *pp.ROMFile
		} else {
			pci[mkHostPCIDeviceROMFile] = ""
		}

		if pp.XVGA != nil {
			pci[mkHostPCIDeviceXVGA] = *pp.XVGA
		} else {
			pci[mkHostPCIDeviceXVGA] = false
		}

		if pp.Mapping != nil {
			pci[mkHostPCIDeviceMapping] = *pp.Mapping
		} else {
			pci[mkHostPCIDeviceMapping] = ""
		}

		pciMap[pi] = pci
	}

	if len(clone) == 0 || len(currentPCIList) > 0 {
		orderedPCIList := utils.OrderedListFromMap(pciMap)
		err := d.Set(mkHostPCI, orderedPCIList)
		diags = append(diags, diag.FromErr(err)...)
	}

	currentUSBList := d.Get(mkHostUSB).([]interface{})
	usbMap := map[string]interface{}{}

	usbDevices := getUSBInfo(vmConfig, d)
	for pi, pp := range usbDevices {
		if (pp == nil) || (pp.HostDevice == nil && pp.Mapping == nil) {
			continue
		}

		usb := map[string]interface{}{}

		if pp.HostDevice != nil {
			usb[mkHostUSBDevice] = *pp.HostDevice
		} else {
			usb[mkHostUSBDevice] = ""
		}

		if pp.USB3 != nil {
			usb[mkHostUSBDeviceUSB3] = *pp.USB3
		} else {
			usb[mkHostUSBDeviceUSB3] = false
		}

		if pp.Mapping != nil {
			usb[mkHostUSBDeviceMapping] = *pp.Mapping
		} else {
			usb[mkHostUSBDeviceMapping] = ""
		}

		usbMap[pi] = usb
	}

	if len(clone) == 0 || len(currentUSBList) > 0 {
		// NOTE: reordering of devices by PVE may cause an issue here
		orderedUSBList := utils.OrderedListFromMap(usbMap)
		err := d.Set(mkHostUSB, orderedUSBList)
		diags = append(diags, diag.FromErr(err)...)
	}

	// Compare the initialization configuration to the one stored in the state.
	initialization := map[string]interface{}{}

	initializationInterface := findExistingCloudInitDrive(vmConfig, vmID, "")
	if initializationInterface != "" {
		initializationDevice := getStorageDevice(vmConfig, initializationInterface)
		fileVolumeParts := strings.Split(initializationDevice.FileVolume, ":")

		initialization[mkInitializationInterface] = initializationInterface
		initialization[mkInitializationDatastoreID] = fileVolumeParts[0]
	}

	if vmConfig.CloudInitDNSDomain != nil || vmConfig.CloudInitDNSServer != nil {
		initializationDNS := map[string]interface{}{}

		if vmConfig.CloudInitDNSDomain != nil {
			initializationDNS[mkInitializationDNSDomain] = *vmConfig.CloudInitDNSDomain
		} else {
			initializationDNS[mkInitializationDNSDomain] = ""
		}

		// check what we have in the plan
		currentInitializationDNSBlock := map[string]interface{}{}
		currentInitialization := d.Get(mkInitialization).([]interface{})

		if len(currentInitialization) > 0 && currentInitialization[0] != nil {
			currentInitializationBlock := currentInitialization[0].(map[string]interface{})
			currentInitializationDNS := currentInitializationBlock[mkInitializationDNS].([]interface{})

			if len(currentInitializationDNS) > 0 && currentInitializationDNS[0] != nil {
				currentInitializationDNSBlock = currentInitializationDNS[0].(map[string]interface{})
			}
		}

		currentInitializationDNSServer, ok := currentInitializationDNSBlock[mkInitializationDNSServer]
		if vmConfig.CloudInitDNSServer != nil {
			if ok && currentInitializationDNSServer != "" {
				// the template is using deprecated attribute mkInitializationDNSServer
				initializationDNS[mkInitializationDNSServer] = *vmConfig.CloudInitDNSServer
			} else {
				dnsServer := strings.Split(*vmConfig.CloudInitDNSServer, " ")
				initializationDNS[mkInitializationDNSServers] = dnsServer
			}
		} else {
			initializationDNS[mkInitializationDNSServer] = ""
			initializationDNS[mkInitializationDNSServers] = []string{}
		}

		initialization[mkInitializationDNS] = []interface{}{
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
		vmConfig.IPConfig7,
		vmConfig.IPConfig8,
		vmConfig.IPConfig9,
		vmConfig.IPConfig10,
		vmConfig.IPConfig11,
		vmConfig.IPConfig12,
		vmConfig.IPConfig13,
		vmConfig.IPConfig14,
		vmConfig.IPConfig15,
		vmConfig.IPConfig16,
		vmConfig.IPConfig17,
		vmConfig.IPConfig18,
		vmConfig.IPConfig19,
		vmConfig.IPConfig20,
		vmConfig.IPConfig21,
		vmConfig.IPConfig22,
		vmConfig.IPConfig23,
		vmConfig.IPConfig24,
		vmConfig.IPConfig25,
		vmConfig.IPConfig26,
		vmConfig.IPConfig27,
		vmConfig.IPConfig28,
		vmConfig.IPConfig29,
		vmConfig.IPConfig30,
		vmConfig.IPConfig31,
	}
	ipConfigList := make([]interface{}, len(ipConfigObjects))

	for ipConfigIndex, ipConfig := range ipConfigObjects {
		ipConfigItem := map[string]interface{}{}

		if ipConfig != nil {
			ipConfigLast = ipConfigIndex

			if ipConfig.GatewayIPv4 != nil || ipConfig.IPv4 != nil {
				ipv4 := map[string]interface{}{}

				if ipConfig.IPv4 != nil {
					ipv4[mkInitializationIPConfigIPv4Address] = *ipConfig.IPv4
				} else {
					ipv4[mkInitializationIPConfigIPv4Address] = ""
				}

				if ipConfig.GatewayIPv4 != nil {
					ipv4[mkInitializationIPConfigIPv4Gateway] = *ipConfig.GatewayIPv4
				} else {
					ipv4[mkInitializationIPConfigIPv4Gateway] = ""
				}

				ipConfigItem[mkInitializationIPConfigIPv4] = []interface{}{
					ipv4,
				}
			} else {
				ipConfigItem[mkInitializationIPConfigIPv4] = []interface{}{}
			}

			if ipConfig.GatewayIPv6 != nil || ipConfig.IPv6 != nil {
				ipv6 := map[string]interface{}{}

				if ipConfig.IPv6 != nil {
					ipv6[mkInitializationIPConfigIPv6Address] = *ipConfig.IPv6
				} else {
					ipv6[mkInitializationIPConfigIPv6Address] = ""
				}

				if ipConfig.GatewayIPv6 != nil {
					ipv6[mkInitializationIPConfigIPv6Gateway] = *ipConfig.GatewayIPv6
				} else {
					ipv6[mkInitializationIPConfigIPv6Gateway] = ""
				}

				ipConfigItem[mkInitializationIPConfigIPv6] = []interface{}{
					ipv6,
				}
			} else {
				ipConfigItem[mkInitializationIPConfigIPv6] = []interface{}{}
			}
		} else {
			ipConfigItem[mkInitializationIPConfigIPv4] = []interface{}{}
			ipConfigItem[mkInitializationIPConfigIPv6] = []interface{}{}
		}

		ipConfigList[ipConfigIndex] = ipConfigItem
	}

	if ipConfigLast >= 0 {
		initialization[mkInitializationIPConfig] = ipConfigList[:ipConfigLast+1]
	}

	if vmConfig.CloudInitPassword != nil || vmConfig.CloudInitSSHKeys != nil ||
		vmConfig.CloudInitUsername != nil {
		initializationUserAccount := map[string]interface{}{}

		if vmConfig.CloudInitSSHKeys != nil {
			initializationUserAccount[mkInitializationUserAccountKeys] = []string(
				*vmConfig.CloudInitSSHKeys,
			)
		} else {
			initializationUserAccount[mkInitializationUserAccountKeys] = []string{}
		}

		if vmConfig.CloudInitPassword != nil {
			initializationUserAccount[mkInitializationUserAccountPassword] = *vmConfig.CloudInitPassword
		} else {
			initializationUserAccount[mkInitializationUserAccountPassword] = ""
		}

		if vmConfig.CloudInitUsername != nil {
			initializationUserAccount[mkInitializationUserAccountUsername] = *vmConfig.CloudInitUsername
		} else {
			initializationUserAccount[mkInitializationUserAccountUsername] = ""
		}

		initialization[mkInitializationUserAccount] = []interface{}{
			initializationUserAccount,
		}
	}

	if vmConfig.CloudInitFiles != nil {
		if vmConfig.CloudInitFiles.UserVolume != nil {
			initialization[mkInitializationUserDataFileID] = *vmConfig.CloudInitFiles.UserVolume
		} else {
			initialization[mkInitializationUserDataFileID] = ""
		}

		if vmConfig.CloudInitFiles.VendorVolume != nil {
			initialization[mkInitializationVendorDataFileID] = *vmConfig.CloudInitFiles.VendorVolume
		} else {
			initialization[mkInitializationVendorDataFileID] = ""
		}

		if vmConfig.CloudInitFiles.NetworkVolume != nil {
			initialization[mkInitializationNetworkDataFileID] = *vmConfig.CloudInitFiles.NetworkVolume
		} else {
			initialization[mkInitializationNetworkDataFileID] = ""
		}

		if vmConfig.CloudInitFiles.MetaVolume != nil {
			initialization[mkInitializationMetaDataFileID] = *vmConfig.CloudInitFiles.MetaVolume
		} else {
			initialization[mkInitializationMetaDataFileID] = ""
		}
	} else if len(initialization) > 0 {
		initialization[mkInitializationUserDataFileID] = ""
		initialization[mkInitializationVendorDataFileID] = ""
		initialization[mkInitializationNetworkDataFileID] = ""
		initialization[mkInitializationMetaDataFileID] = ""
	}

	if vmConfig.CloudInitType != nil {
		initialization[mkInitializationType] = *vmConfig.CloudInitType
	} else if len(initialization) > 0 {
		initialization[mkInitializationType] = ""
	}

	currentInitialization := d.Get(mkInitialization).([]interface{})

	//nolint:gocritic
	if len(clone) > 0 {
		if len(currentInitialization) > 0 {
			if len(initialization) > 0 {
				err := d.Set(mkInitialization, []interface{}{initialization})
				diags = append(diags, diag.FromErr(err)...)
			} else {
				err := d.Set(mkInitialization, []interface{}{})
				diags = append(diags, diag.FromErr(err)...)
			}
		}
	} else if len(initialization) > 0 {
		err := d.Set(mkInitialization, []interface{}{initialization})
		diags = append(diags, diag.FromErr(err)...)
	} else {
		err := d.Set(mkInitialization, []interface{}{})
		diags = append(diags, diag.FromErr(err)...)
	}

	// Compare the operating system configuration to the one stored in the state.
	kvmArguments := map[string]interface{}{}

	if vmConfig.KVMArguments != nil {
		kvmArguments[mkKVMArguments] = *vmConfig.KVMArguments
	} else {
		kvmArguments[mkKVMArguments] = ""
	}

	// Compare the memory configuration to the one stored in the state.
	memory := map[string]interface{}{}

	if vmConfig.DedicatedMemory != nil {
		memory[mkMemoryDedicated] = int(*vmConfig.DedicatedMemory)
	} else {
		memory[mkMemoryDedicated] = 0
	}

	if vmConfig.FloatingMemory != nil {
		memory[mkMemoryFloating] = int(*vmConfig.FloatingMemory)
	} else {
		memory[mkMemoryFloating] = 0
	}

	if vmConfig.SharedMemory != nil {
		memory[mkMemoryShared] = vmConfig.SharedMemory.Size
	} else {
		memory[mkMemoryShared] = 0
	}

	if vmConfig.Hugepages != nil {
		memory[mkMemoryHugepages] = *vmConfig.Hugepages
	} else {
		memory[mkMemoryHugepages] = ""
	}

	if vmConfig.KeepHugepages != nil {
		memory[mkMemoryKeepHugepages] = *vmConfig.KeepHugepages
	} else {
		memory[mkMemoryKeepHugepages] = false
	}

	currentMemory := d.Get(mkMemory).([]interface{})

	if len(clone) > 0 {
		if len(currentMemory) > 0 {
			err := d.Set(mkMemory, []interface{}{memory})
			diags = append(diags, diag.FromErr(err)...)
		}
	} else if len(currentMemory) > 0 ||
		memory[mkMemoryDedicated] != dvMemoryDedicated ||
		memory[mkMemoryFloating] != dvMemoryFloating ||
		memory[mkMemoryShared] != dvMemoryShared ||
		memory[mkMemoryHugepages] != dvMemoryHugepages ||
		memory[mkMemoryKeepHugepages] != dvMemoryKeepHugepages {
		err := d.Set(mkMemory, []interface{}{memory})
		diags = append(diags, diag.FromErr(err)...)
	}

	diags = append(diags, network.ReadNetworkDeviceObjects(d, vmConfig)...)

	// Compare the operating system configuration to the one stored in the state.
	operatingSystem := map[string]interface{}{}

	if vmConfig.OSType != nil {
		operatingSystem[mkOperatingSystemType] = *vmConfig.OSType
	} else {
		operatingSystem[mkOperatingSystemType] = ""
	}

	currentOperatingSystem := d.Get(mkOperatingSystem).([]interface{})

	switch {
	case len(clone) > 0:
		if len(currentOperatingSystem) > 0 {
			err := d.Set(
				mkOperatingSystem,
				[]interface{}{operatingSystem},
			)
			diags = append(diags, diag.FromErr(err)...)
		}
	case len(currentOperatingSystem) > 0 || operatingSystem[mkOperatingSystemType] != dvOperatingSystemType:
		err := d.Set(mkOperatingSystem, []interface{}{operatingSystem})
		diags = append(diags, diag.FromErr(err)...)
	default:
		err := d.Set(mkOperatingSystem, []interface{}{})
		diags = append(diags, diag.FromErr(err)...)
	}

	// Compare the pool ID to the value stored in the state.
	currentPoolID := d.Get(mkPoolID).(string)

	if len(clone) == 0 || currentPoolID != dvPoolID {
		if vmConfig.PoolID != nil {
			err := d.Set(mkPoolID, *vmConfig.PoolID)
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
			m[mkSerialDeviceDevice] = *sd
			serialDevicesCount = sdi + 1
		} else {
			m[mkSerialDeviceDevice] = ""
		}

		serialDevices[sdi] = m
	}

	currentSerialDevice := d.Get(mkSerialDevice).([]interface{})

	if len(clone) == 0 || len(currentSerialDevice) > 0 {
		err := d.Set(mkSerialDevice, serialDevices[:serialDevicesCount])
		diags = append(diags, diag.FromErr(err)...)
	}

	// Compare the SMBIOS to the one stored in the state.
	var smbios map[string]interface{}

	if vmConfig.SMBIOS != nil {
		smbios = map[string]interface{}{}

		if vmConfig.SMBIOS.Family != nil {
			b, err := base64.StdEncoding.DecodeString(*vmConfig.SMBIOS.Family)
			diags = append(diags, diag.FromErr(err)...)
			smbios[mkSMBIOSFamily] = string(b)
		} else {
			smbios[mkSMBIOSFamily] = dvSMBIOSFamily
		}

		if vmConfig.SMBIOS.Manufacturer != nil {
			b, err := base64.StdEncoding.DecodeString(*vmConfig.SMBIOS.Manufacturer)
			diags = append(diags, diag.FromErr(err)...)
			smbios[mkSMBIOSManufacturer] = string(b)
		} else {
			smbios[mkSMBIOSManufacturer] = dvSMBIOSManufacturer
		}

		if vmConfig.SMBIOS.Product != nil {
			b, err := base64.StdEncoding.DecodeString(*vmConfig.SMBIOS.Product)
			diags = append(diags, diag.FromErr(err)...)
			smbios[mkSMBIOSProduct] = string(b)
		} else {
			smbios[mkSMBIOSProduct] = dvSMBIOSProduct
		}

		if vmConfig.SMBIOS.Serial != nil {
			b, err := base64.StdEncoding.DecodeString(*vmConfig.SMBIOS.Serial)
			diags = append(diags, diag.FromErr(err)...)
			smbios[mkSMBIOSSerial] = string(b)
		} else {
			smbios[mkSMBIOSSerial] = dvSMBIOSSerial
		}

		if vmConfig.SMBIOS.SKU != nil {
			b, err := base64.StdEncoding.DecodeString(*vmConfig.SMBIOS.SKU)
			diags = append(diags, diag.FromErr(err)...)
			smbios[mkSMBIOSSKU] = string(b)
		} else {
			smbios[mkSMBIOSSKU] = dvSMBIOSSKU
		}

		if vmConfig.SMBIOS.Version != nil {
			b, err := base64.StdEncoding.DecodeString(*vmConfig.SMBIOS.Version)
			diags = append(diags, diag.FromErr(err)...)
			smbios[mkSMBIOSVersion] = string(b)
		} else {
			smbios[mkSMBIOSVersion] = dvSMBIOSVersion
		}

		if vmConfig.SMBIOS.UUID != nil {
			smbios[mkSMBIOSUUID] = *vmConfig.SMBIOS.UUID
		} else {
			smbios[mkSMBIOSUUID] = nil
		}
	}

	currentSMBIOS := d.Get(mkSMBIOS).([]interface{})

	switch {
	case len(clone) > 0:
		if len(currentSMBIOS) > 0 {
			err := d.Set(mkSMBIOS, currentSMBIOS)
			diags = append(diags, diag.FromErr(err)...)
		}
	case len(smbios) == 0:
		err := d.Set(mkSMBIOS, []interface{}{})
		diags = append(diags, diag.FromErr(err)...)
	default:
		if len(currentSMBIOS) > 0 ||
			smbios[mkSMBIOSFamily] != dvSMBIOSFamily ||
			smbios[mkSMBIOSManufacturer] != dvSMBIOSManufacturer ||
			smbios[mkSMBIOSProduct] != dvSMBIOSProduct ||
			smbios[mkSMBIOSSerial] != dvSMBIOSSerial ||
			smbios[mkSMBIOSSKU] != dvSMBIOSSKU ||
			smbios[mkSMBIOSVersion] != dvSMBIOSVersion {
			err := d.Set(mkSMBIOS, []interface{}{smbios})
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	// Compare the startup order to the one stored in the state.
	var startup map[string]interface{}

	if vmConfig.StartupOrder != nil {
		startup = map[string]interface{}{}

		if vmConfig.StartupOrder.Order != nil {
			startup[mkStartupOrder] = *vmConfig.StartupOrder.Order
		} else {
			startup[mkStartupOrder] = dvStartupOrder
		}

		if vmConfig.StartupOrder.Up != nil {
			startup[mkStartupUpDelay] = *vmConfig.StartupOrder.Up
		} else {
			startup[mkStartupUpDelay] = dvStartupUpDelay
		}

		if vmConfig.StartupOrder.Down != nil {
			startup[mkStartupDownDelay] = *vmConfig.StartupOrder.Down
		} else {
			startup[mkStartupDownDelay] = dvStartupDownDelay
		}
	}

	currentStartup := d.Get(mkStartup).([]interface{})

	switch {
	case len(clone) > 0:
		if len(currentStartup) > 0 {
			err := d.Set(mkStartup, []interface{}{startup})
			diags = append(diags, diag.FromErr(err)...)
		}
	case len(startup) == 0:
		err := d.Set(mkStartup, []interface{}{})
		diags = append(diags, diag.FromErr(err)...)
	default:
		if len(currentStartup) > 0 ||
			startup[mkStartupOrder] != mkStartupOrder ||
			startup[mkStartupUpDelay] != dvStartupUpDelay ||
			startup[mkStartupDownDelay] != dvStartupDownDelay {
			err := d.Set(mkStartup, []interface{}{startup})
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	// Compare the VGA configuration to the one stored in the state.
	vga := map[string]interface{}{}

	if vmConfig.VGADevice != nil {
		if vmConfig.VGADevice.Clipboard != nil {
			vga[mkVGAClipboard] = *vmConfig.VGADevice.Clipboard
		} else {
			vga[mkVGAClipboard] = dvVGAClipboard
		}

		if vmConfig.VGADevice.Memory != nil {
			vga[mkVGAMemory] = int(*vmConfig.VGADevice.Memory)
		} else {
			vga[mkVGAMemory] = dvVGAMemory
		}

		if vmConfig.VGADevice.Type != nil {
			vga[mkVGAType] = *vmConfig.VGADevice.Type
		}
	} else {
		vga[mkVGAClipboard] = ""
		vga[mkVGAMemory] = 0
		vga[mkVGAType] = ""
	}

	currentVGA := d.Get(mkVGA).([]interface{})

	switch {
	case len(clone) > 0 && len(currentVGA) > 0:
		err := d.Set(mkVGA, []interface{}{vga})
		diags = append(diags, diag.FromErr(err)...)
	case len(currentVGA) > 0 ||
		vga[mkVGAClipboard] != dvVGAClipboard ||
		vga[mkVGAMemory] != dvVGAMemory ||
		vga[mkVGAType] != dvVGAType:
		err := d.Set(mkVGA, []interface{}{vga})
		diags = append(diags, diag.FromErr(err)...)
	default:
		err := d.Set(mkVGA, []interface{}{})
		diags = append(diags, diag.FromErr(err)...)
	}

	// Compare SCSI hardware type
	scsiHardware := d.Get(mkSCSIHardware).(string)

	if len(clone) == 0 || scsiHardware != dvSCSIHardware {
		if vmConfig.SCSIHardware != nil {
			err := d.Set(mkSCSIHardware, *vmConfig.SCSIHardware)
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	vmAPI := client.Node(nodeName).VM(vmID)
	started := d.Get(mkStarted).(bool)

	agentTimeout, e := getAgentTimeout(d)
	if e != nil {
		return diag.FromErr(e)
	}

	diags = append(
		diags,
		network.ReadNetworkValues(ctx, d, vmAPI, started, vmConfig, agentTimeout)...)

	// during import these core attributes might not be set, so set them explicitly here
	d.SetId(strconv.Itoa(vmID))
	e = d.Set(mkVMID, vmID)
	diags = append(diags, diag.FromErr(e)...)
	e = d.Set(mkNodeName, nodeName)
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

	clone := d.Get(mkClone).([]interface{})
	currentACPI := d.Get(mkACPI).(bool)

	//nolint:gosimple
	if len(clone) == 0 || currentACPI != dvACPI {
		if vmConfig.ACPI != nil {
			err = d.Set(mkACPI, bool(*vmConfig.ACPI))
		} else {
			// Default value of "acpi" is "1" according to the API documentation.
			err = d.Set(mkACPI, true)
		}

		diags = append(diags, diag.FromErr(err)...)
	}

	currentKVMArguments := d.Get(mkKVMArguments).(string)

	if len(clone) == 0 || currentKVMArguments != dvKVMArguments {
		// PVE API returns "args" as " " if it is set to empty.
		if vmConfig.KVMArguments != nil && len(strings.TrimSpace(*vmConfig.KVMArguments)) > 0 {
			err = d.Set(mkKVMArguments, *vmConfig.KVMArguments)
		} else {
			// Default value of "args" is "" according to the API documentation.
			err = d.Set(mkKVMArguments, "")
		}

		diags = append(diags, diag.FromErr(err)...)
	}

	currentBIOS := d.Get(mkBIOS).(string)

	if len(clone) == 0 || currentBIOS != dvBIOS {
		if vmConfig.BIOS != nil {
			err = d.Set(mkBIOS, *vmConfig.BIOS)
		} else {
			// Default value of "bios" is "seabios" according to the API documentation.
			err = d.Set(mkBIOS, "seabios")
		}

		diags = append(diags, diag.FromErr(err)...)
	}

	currentDescription := d.Get(mkDescription).(string)

	if len(clone) == 0 || currentDescription != dvDescription {
		if vmConfig.Description != nil {
			err = d.Set(mkDescription, *vmConfig.Description)
		} else {
			// Default value of "description" is "" according to the API documentation.
			err = d.Set(mkDescription, "")
		}

		diags = append(diags, diag.FromErr(err)...)
	}

	currentTags := d.Get(mkTags).([]interface{})

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

		err = d.Set(mkTags, tags)
		diags = append(diags, diag.FromErr(err)...)
	}

	currentKeyboardLayout := d.Get(mkKeyboardLayout).(string)

	if len(clone) == 0 || currentKeyboardLayout != dvKeyboardLayout {
		if vmConfig.KeyboardLayout != nil {
			err = d.Set(mkKeyboardLayout, *vmConfig.KeyboardLayout)
		} else {
			// Default value of "keyboard" is "" according to the API documentation.
			err = d.Set(mkKeyboardLayout, "")
		}

		diags = append(diags, diag.FromErr(err)...)
	}

	currentMachine := d.Get(mkMachine).(string)

	if len(clone) == 0 || currentMachine != dvMachineType {
		if vmConfig.Machine != nil {
			err = d.Set(mkMachine, *vmConfig.Machine)
		} else {
			err = d.Set(mkMachine, "")
		}

		diags = append(diags, diag.FromErr(err)...)
	}

	currentName := d.Get(mkName).(string)

	if len(clone) == 0 || currentName != dvName {
		if vmConfig.Name != nil {
			err = d.Set(mkName, *vmConfig.Name)
		} else {
			// Default value of "name" is "" according to the API documentation.
			err = d.Set(mkName, "")
		}

		diags = append(diags, diag.FromErr(err)...)
	}

	currentProtection := d.Get(mkProtection).(bool)

	//nolint:gosimple
	if len(clone) == 0 || currentProtection != dvProtection {
		if vmConfig.DeletionProtection != nil {
			err = d.Set(
				mkProtection,
				bool(*vmConfig.DeletionProtection),
			)
		} else {
			// Default value of "protection" is "0" according to the API documentation.
			err = d.Set(mkProtection, false)
		}

		diags = append(diags, diag.FromErr(err)...)
	}

	if !d.Get(mkTemplate).(bool) {
		err = d.Set(mkStarted, vmStatus.Status == "running")
		diags = append(diags, diag.FromErr(err)...)
	}

	currentTabletDevice := d.Get(mkTabletDevice).(bool)

	//nolint:gosimple
	if len(clone) == 0 || currentTabletDevice != dvTabletDevice {
		if vmConfig.TabletDeviceEnabled != nil {
			err = d.Set(
				mkTabletDevice,
				bool(*vmConfig.TabletDeviceEnabled),
			)
		} else {
			// Default value of "tablet" is "1" according to the API documentation.
			err = d.Set(mkTabletDevice, true)
		}

		diags = append(diags, diag.FromErr(err)...)
	}

	currentTemplate := d.Get(mkTemplate).(bool)

	//nolint:gosimple
	if len(clone) == 0 || currentTemplate != dvTemplate {
		if vmConfig.Template != nil {
			err = d.Set(mkTemplate, bool(*vmConfig.Template))
		} else {
			// Default value of "template" is "0" according to the API documentation.
			err = d.Set(mkTemplate, false)
		}

		diags = append(diags, diag.FromErr(err)...)
	}

	return diags
}

// vmUpdatePool moves the VM to the pool it is supposed to be in if the pool ID changed.
func vmUpdatePool(
	ctx context.Context,
	d *schema.ResourceData,
	api *pools.Client,
	vmID int,
) error {
	oldPoolValue, newPoolValue := d.GetChange(mkPoolID)
	if cmp.Equal(newPoolValue, oldPoolValue) {
		return nil
	}

	oldPool := oldPoolValue.(string)
	newPool := newPoolValue.(string)
	vmList := (types.CustomCommaSeparatedList)([]string{strconv.Itoa(vmID)})

	tflog.Debug(ctx, fmt.Sprintf("Moving VM %d from pool '%s' to pool '%s'", vmID, oldPool, newPool))

	if oldPool != "" {
		trueValue := types.CustomBool(true)
		poolUpdate := &pools.PoolUpdateRequestBody{
			VMs:    &vmList,
			Delete: &trueValue,
		}

		err := api.UpdatePool(ctx, oldPool, poolUpdate)
		if err != nil {
			return fmt.Errorf("while removing VM %d from pool %s: %w", vmID, oldPool, err)
		}
	}

	if newPool != "" {
		poolUpdate := &pools.PoolUpdateRequestBody{VMs: &vmList}

		err := api.UpdatePool(ctx, newPool, poolUpdate)
		if err != nil {
			return fmt.Errorf("while adding VM %d to pool %s: %w", vmID, newPool, err)
		}
	}

	return nil
}

func vmUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(proxmoxtf.ProviderConfiguration)

	client, e := config.GetClient()
	if e != nil {
		return diag.FromErr(e)
	}

	nodeName := d.Get(mkNodeName).(string)
	rebootRequired := false

	vmID, e := strconv.Atoi(d.Id())
	if e != nil {
		return diag.FromErr(e)
	}

	e = vmUpdatePool(ctx, d, client.Pool(), vmID)
	if e != nil {
		return diag.FromErr(e)
	}

	// If the node name has changed we need to migrate the VM to the new node before we do anything else.
	if d.HasChange(mkNodeName) {
		migrateTimeoutSec := d.Get(mkTimeoutMigrate).(int)

		ctx, cancel := context.WithTimeout(ctx, time.Duration(migrateTimeoutSec)*time.Second)
		defer cancel()

		oldNodeNameValue, _ := d.GetChange(mkNodeName)
		oldNodeName := oldNodeNameValue.(string)
		vmAPI := client.Node(oldNodeName).VM(vmID)

		trueValue := types.CustomBool(true)
		migrateBody := &vms.MigrateRequestBody{
			TargetNode:      nodeName,
			WithLocalDisks:  &trueValue,
			OnlineMigration: &trueValue,
		}

		err := vmAPI.MigrateVM(ctx, migrateBody)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	vmAPI := client.Node(nodeName).VM(vmID)

	updateBody := &vms.UpdateRequestBody{
		IDEDevices: vms.CustomStorageDevices{
			"ide0": &vms.CustomStorageDevice{
				Enabled: false,
			},
			"ide1": &vms.CustomStorageDevice{
				Enabled: false,
			},
			"ide2": &vms.CustomStorageDevice{
				Enabled: false,
			},
			"ide3": &vms.CustomStorageDevice{
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
	if d.HasChange(mkACPI) {
		acpi := types.CustomBool(d.Get(mkACPI).(bool))
		updateBody.ACPI = &acpi
		rebootRequired = true
	}

	if d.HasChange(mkKVMArguments) {
		kvmArguments := d.Get(mkKVMArguments).(string)
		updateBody.KVMArguments = &kvmArguments
		rebootRequired = true
	}

	if d.HasChange(mkBIOS) {
		bios := d.Get(mkBIOS).(string)
		updateBody.BIOS = &bios
		rebootRequired = true
	}

	if d.HasChange(mkDescription) {
		description := d.Get(mkDescription).(string)
		updateBody.Description = &description
	}

	if d.HasChange(mkOnBoot) {
		startOnBoot := types.CustomBool(d.Get(mkOnBoot).(bool))
		updateBody.StartOnBoot = &startOnBoot
	}

	if d.HasChange(mkTags) {
		tagString := vmGetTagsString(d)
		updateBody.Tags = &tagString
	}

	if d.HasChange(mkKeyboardLayout) {
		keyboardLayout := d.Get(mkKeyboardLayout).(string)
		updateBody.KeyboardLayout = &keyboardLayout
		rebootRequired = true
	}

	if d.HasChange(mkMachine) {
		machine := d.Get(mkMachine).(string)
		updateBody.Machine = &machine
		rebootRequired = true
	}

	name := d.Get(mkName).(string)

	if name == "" {
		del = append(del, "name")
	} else {
		updateBody.Name = &name
	}

	if d.HasChange(mkProtection) {
		protection := types.CustomBool(d.Get(mkProtection).(bool))
		updateBody.DeletionProtection = &protection
	}

	if d.HasChange(mkTabletDevice) {
		tabletDevice := types.CustomBool(d.Get(mkTabletDevice).(bool))
		updateBody.TabletDeviceEnabled = &tabletDevice
		rebootRequired = true
	}

	template := types.CustomBool(d.Get(mkTemplate).(bool))

	if d.HasChange(mkTemplate) {
		updateBody.Template = &template
		rebootRequired = true
	}

	// Prepare the new agent configuration.
	if d.HasChange(mkAgent) {
		agentBlock, err := structure.GetSchemaBlock(
			resource,
			d,
			[]string{mkAgent},
			0,
			true,
		)
		if err != nil {
			return diag.FromErr(err)
		}

		agentEnabled := types.CustomBool(
			agentBlock[mkAgentEnabled].(bool),
		)
		agentTrim := types.CustomBool(agentBlock[mkAgentTrim].(bool))
		agentType := agentBlock[mkAgentType].(string)

		updateBody.Agent = &vms.CustomAgent{
			Enabled:         &agentEnabled,
			TrimClonedDisks: &agentTrim,
			Type:            &agentType,
		}

		rebootRequired = true
	}

	// Prepare the new audio devices.
	if d.HasChange(mkAudioDevice) {
		updateBody.AudioDevices = vmGetAudioDeviceList(d)

		for i, ad := range updateBody.AudioDevices {
			if !ad.Enabled {
				del = append(del, fmt.Sprintf("audio%d", i))
			}
		}

		for i := len(updateBody.AudioDevices); i < maxResourceVirtualEnvironmentVMAudioDevices; i++ {
			del = append(del, fmt.Sprintf("audio%d", i))
		}

		rebootRequired = true
	}

	// Prepare the new boot configuration.
	if d.HasChange(mkBootOrder) {
		bootOrder := d.Get(mkBootOrder).([]interface{})
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

	if d.HasChange(mkCDROM) {
		cdromBlock, err := structure.GetSchemaBlock(
			resource,
			d,
			[]string{mkCDROM},
			0,
			true,
		)
		if err != nil {
			return diag.FromErr(err)
		}

		cdromEnabled := cdromBlock[mkCDROMEnabled].(bool)
		cdromFileID := cdromBlock[mkCDROMFileID].(string)
		cdromInterface := cdromBlock[mkCDROMInterface].(string)

		old, _ := d.GetChange(mkCDROM)

		if len(old.([]interface{})) > 0 && old.([]interface{})[0] != nil {
			oldList := old.([]interface{})[0]
			oldBlock := oldList.(map[string]interface{})

			// If the interface is not set, use the default, for backward compatibility.
			oldInterface, ok := oldBlock[mkCDROMInterface].(string)
			if !ok || oldInterface == "" {
				oldInterface = dvCDROMInterface
			}

			if oldInterface != cdromInterface {
				del = append(del, oldInterface)
			}
		}

		if !cdromEnabled && cdromFileID == "" {
			del = append(del, cdromInterface)
		}

		if cdromFileID == "" {
			cdromFileID = "cdrom"
		}

		cdromMedia := "cdrom"

		updateBody.IDEDevices[cdromInterface] = &vms.CustomStorageDevice{
			Enabled:    cdromEnabled,
			FileVolume: cdromFileID,
			Media:      &cdromMedia,
		}
	}

	// Prepare the new CPU configuration.

	if d.HasChange(mkCPU) {
		cpuBlock, err := structure.GetSchemaBlock(
			resource,
			d,
			[]string{mkCPU},
			0,
			true,
		)
		if err != nil {
			return diag.FromErr(err)
		}

		cpuArchitecture := cpuBlock[mkCPUArchitecture].(string)
		cpuCores := cpuBlock[mkCPUCores].(int)
		cpuFlags := cpuBlock[mkCPUFlags].([]interface{})
		cpuHotplugged := cpuBlock[mkCPUHotplugged].(int)
		cpuLimit := cpuBlock[mkCPULimit].(int)
		cpuNUMA := types.CustomBool(cpuBlock[mkCPUNUMA].(bool))
		cpuSockets := cpuBlock[mkCPUSockets].(int)
		cpuType := cpuBlock[mkCPUType].(string)
		cpuUnits := cpuBlock[mkCPUUnits].(int)
		cpuAffinity := cpuBlock[mkCPUAffinity].(string)

		// Only the root account is allowed to change the CPU architecture, which makes this check necessary.
		if client.API().IsRootTicket() ||
			cpuArchitecture != dvCPUArchitecture {
			updateBody.CPUArchitecture = &cpuArchitecture
		}

		updateBody.CPUCores = ptr.Ptr(int64(cpuCores))
		updateBody.CPUSockets = ptr.Ptr(int64(cpuSockets))
		updateBody.CPUUnits = ptr.Ptr(int64(cpuUnits))
		updateBody.NUMAEnabled = &cpuNUMA

		// CPU affinity is a special case, only root can change it.
		// we can't even have it in the delete list, as PVE will return an error for non-root.
		// Hence, checking explicitly if it has changed.
		if d.HasChange(mkCPU + ".0." + mkCPUAffinity) {
			if cpuAffinity != "" {
				updateBody.CPUAffinity = &cpuAffinity
			} else {
				del = append(del, "affinity")
			}
		}

		if cpuHotplugged > 0 {
			updateBody.VirtualCPUCount = ptr.Ptr(int64(cpuHotplugged))
		} else {
			del = append(del, "vcpus")
		}

		if cpuLimit > 0 {
			updateBody.CPULimit = ptr.Ptr(int64(cpuLimit))
		} else {
			del = append(del, "cpulimit")
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
	allDiskInfo := disk.GetInfo(vmConfig, d)

	planDisks, err := disk.GetDiskDeviceObjects(d, resource, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	rr, err := disk.Update(ctx, client, nodeName, vmID, d, planDisks, allDiskInfo, updateBody)
	if err != nil {
		return diag.FromErr(err)
	}

	rebootRequired = rebootRequired || rr

	// Prepare the new efi disk configuration.
	if d.HasChange(mkEFIDisk) {
		efiDisk := vmGetEfiDisk(d, nil)

		updateBody.EFIDisk = efiDisk

		rebootRequired = true
	}

	// Prepare the new tpm state configuration.
	if d.HasChange(mkTPMState) {
		tpmState := vmGetTPMState(d, nil)

		updateBody.TPMState = tpmState

		rebootRequired = true
	}

	// Prepare the new cloud-init configuration.
	stoppedBeforeUpdate := false

	if d.HasChange(mkInitialization) {
		initializationConfig := vmGetCloudInitConfig(d)

		updateBody.CloudInitConfig = initializationConfig

		initialization := d.Get(mkInitialization).([]interface{})

		if updateBody.CloudInitConfig != nil && len(initialization) > 0 && initialization[0] != nil {
			var fileVolume string

			initializationBlock := initialization[0].(map[string]interface{})
			initializationDatastoreID := initializationBlock[mkInitializationDatastoreID].(string)
			initializationInterface := initializationBlock[mkInitializationInterface].(string)
			cdromMedia := "cdrom"

			existingInterface := findExistingCloudInitDrive(vmConfig, vmID, "")
			if initializationInterface == "" && existingInterface == "" {
				initializationInterface = "ide2"
			} else if initializationInterface == "" {
				initializationInterface = existingInterface
			}

			mustMove := existingInterface != "" && initializationInterface != existingInterface
			if mustMove {
				tflog.Debug(ctx, fmt.Sprintf("CloudInit must be moved from %s to %s", existingInterface, initializationInterface))
			}

			oldInit, _ := d.GetChange(mkInitialization)
			oldInitBlock := oldInit.([]interface{})[0].(map[string]interface{})
			prevDatastoreID := oldInitBlock[mkInitializationDatastoreID].(string)

			mustChangeDatastore := prevDatastoreID != initializationDatastoreID
			if mustChangeDatastore {
				tflog.Debug(ctx, fmt.Sprintf("CloudInit must be moved from datastore %s to datastore %s",
					prevDatastoreID, initializationDatastoreID))
			}

			if mustMove || mustChangeDatastore || existingInterface == "" {
				// CloudInit must be moved, either from a device to another or from a datastore
				// to another (or both). This requires the VM to be stopped.
				if err := vmShutdown(ctx, vmAPI, d); err != nil {
					return err
				}

				if err := deleteIdeDrives(ctx, vmAPI, initializationInterface, existingInterface); err != nil {
					return err
				}

				stoppedBeforeUpdate = true
				fileVolume = fmt.Sprintf("%s:cloudinit", initializationDatastoreID)
			} else {
				ideDevice := getStorageDevice(vmConfig, existingInterface)
				fileVolume = ideDevice.FileVolume
			}

			updateBody.IDEDevices[initializationInterface] = &vms.CustomStorageDevice{
				Enabled:    true,
				FileVolume: fileVolume,
				Media:      &cdromMedia,
			}
		}

		rebootRequired = true
	}

	// Prepare the new hostpci devices configuration.
	if d.HasChange(mkHostPCI) {
		updateBody.PCIDevices = vmGetHostPCIDeviceObjects(d)

		for i := len(updateBody.PCIDevices); i < maxResourceVirtualEnvironmentVMHostPCIDevices; i++ {
			del = append(del, fmt.Sprintf("hostpci%d", i))
		}

		rebootRequired = true
	}

	// Prepare the new numa devices configuration.
	if d.HasChange(mkNUMA) {
		updateBody.NUMADevices = vmGetNumaDeviceObjects(d)

		for i := len(updateBody.NUMADevices); i < maxResourceVirtualEnvironmentVMNUMADevices; i++ {
			del = append(del, fmt.Sprintf("numa%d", i))
		}

		rebootRequired = true
	}

	// Prepare the new usb devices configuration.
	if d.HasChange(mkHostUSB) {
		updateBody.USBDevices = vmGetHostUSBDeviceObjects(d)

		for i := len(updateBody.USBDevices); i < maxResourceVirtualEnvironmentVMHostUSBDevices; i++ {
			del = append(del, fmt.Sprintf("usb%d", i))
		}

		rebootRequired = true
	}

	// Prepare the new memory configuration.
	if d.HasChange(mkMemory) {
		memoryBlock, er := structure.GetSchemaBlock(
			resource,
			d,
			[]string{mkMemory},
			0,
			true,
		)
		if er != nil {
			return diag.FromErr(er)
		}

		memoryDedicated := memoryBlock[mkMemoryDedicated].(int)
		memoryFloating := memoryBlock[mkMemoryFloating].(int)
		memoryShared := memoryBlock[mkMemoryShared].(int)
		memoryHugepages := memoryBlock[mkMemoryHugepages].(string)
		memoryKeepHugepages := types.CustomBool(memoryBlock[mkMemoryKeepHugepages].(bool))

		updateBody.DedicatedMemory = &memoryDedicated
		updateBody.FloatingMemory = &memoryFloating

		if memoryShared > 0 {
			memorySharedName := fmt.Sprintf("vm-%d-ivshmem", vmID)

			updateBody.SharedMemory = &vms.CustomSharedMemory{
				Name: &memorySharedName,
				Size: memoryShared,
			}
		}

		if d.HasChange(mkMemory + ".0." + mkMemoryHugepages) {
			if memoryHugepages != "" {
				updateBody.Hugepages = &memoryHugepages
			} else {
				del = append(del, "hugepages")
			}
		}

		if d.HasChange(mkMemory + ".0." + mkMemoryKeepHugepages) {
			if memoryHugepages != "" {
				updateBody.KeepHugepages = &memoryKeepHugepages
			} else {
				del = append(del, "keephugepages")
			}
		}

		rebootRequired = true
	}

	// Prepare the new network device configuration.

	if d.HasChange(network.MkNetworkDevice) {
		updateBody.NetworkDevices, err = network.GetNetworkDeviceObjects(d)
		if err != nil {
			return diag.FromErr(err)
		}

		for i, nd := range updateBody.NetworkDevices {
			if !nd.Enabled {
				del = append(del, fmt.Sprintf("net%d", i))
			}
		}

		for i := len(updateBody.NetworkDevices); i < network.MaxNetworkDevices; i++ {
			del = append(del, fmt.Sprintf("net%d", i))
		}

		rebootRequired = true
	}

	// Prepare the new operating system configuration.
	if d.HasChange(mkOperatingSystem) {
		operatingSystem, err := structure.GetSchemaBlock(
			resource,
			d,
			[]string{mkOperatingSystem},
			0,
			true,
		)
		if err != nil {
			return diag.FromErr(err)
		}

		operatingSystemType := operatingSystem[mkOperatingSystemType].(string)

		updateBody.OSType = &operatingSystemType

		rebootRequired = true
	}

	// Prepare the new serial devices.
	if d.HasChange(mkSerialDevice) {
		updateBody.SerialDevices = vmGetSerialDeviceList(d)

		for i := len(updateBody.SerialDevices); i < maxResourceVirtualEnvironmentVMSerialDevices; i++ {
			del = append(del, fmt.Sprintf("serial%d", i))
		}

		rebootRequired = true
	}

	if d.HasChange(mkSMBIOS) {
		updateBody.SMBIOS = vmGetSMBIOS(d)
		if updateBody.SMBIOS == nil {
			del = append(del, "smbios1")
		}
	}

	if d.HasChange(mkStartup) {
		updateBody.StartupOrder = vmGetStartupOrder(d)
		if updateBody.StartupOrder == nil {
			del = append(del, "startup")
		}
	}

	// Prepare the new VGA configuration.
	if d.HasChange(mkVGA) {
		updateBody.VGADevice, e = vmGetVGADeviceObject(d)
		if e != nil {
			return diag.FromErr(e)
		}

		rebootRequired = true
	}

	// Prepare the new SCSI hardware type
	if d.HasChange(mkSCSIHardware) {
		scsiHardware := d.Get(mkSCSIHardware).(string)
		updateBody.SCSIHardware = &scsiHardware

		rebootRequired = true
	}

	if d.HasChanges(mkHookScriptFileID) {
		hookScript := d.Get(mkHookScriptFileID).(string)
		if len(hookScript) > 0 {
			updateBody.HookScript = &hookScript
		} else {
			del = append(del, "hookscript")
		}
	}

	// Update the configuration now that everything has been prepared.
	updateBody.Delete = del

	e = vmAPI.UpdateVM(ctx, updateBody)
	if e != nil {
		return diag.FromErr(e)
	}

	// Determine if the state of the virtual machine state needs to be changed.
	//nolint: nestif
	if (d.HasChange(mkStarted) || stoppedBeforeUpdate) && !bool(template) {
		started := d.Get(mkStarted).(bool)
		if started {
			if diags := vmStart(ctx, vmAPI, d); diags != nil {
				return diags
			}
		} else {
			if e := vmShutdown(ctx, vmAPI, d); e != nil {
				return e
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

	client, err := config.GetClient()
	if err != nil {
		return diag.FromErr(err)
	}

	nodeName := d.Get(mkNodeName).(string)
	started := d.Get(mkStarted).(bool)
	template := d.Get(mkTemplate).(bool)

	vmID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	vmAPI := client.Node(nodeName).VM(vmID)

	// Determine if any of the disks are changing location and/or size, and initiate the necessary actions.
	//nolint: nestif
	if d.HasChange(disk.MkDisk) {
		diskOld, diskNew := d.GetChange(disk.MkDisk)

		resource := VM()

		diskOldEntries, err := disk.GetDiskDeviceObjects(
			d,
			resource,
			diskOld.([]interface{}),
		)
		if err != nil {
			return diag.FromErr(err)
		}

		diskNewEntries, err := disk.GetDiskDeviceObjects(
			d,
			resource,
			diskNew.([]interface{}),
		)
		if err != nil {
			return diag.FromErr(err)
		}

		// Add efidisk if it has changes
		if d.HasChange(mkEFIDisk) {
			diskOld, diskNew := d.GetChange(mkEFIDisk)

			oldEfiDisk, e := vmGetEfiDiskAsStorageDevice(d, diskOld.([]interface{}))
			if e != nil {
				return diag.FromErr(e)
			}

			newEfiDisk, e := vmGetEfiDiskAsStorageDevice(d, diskNew.([]interface{}))
			if e != nil {
				return diag.FromErr(e)
			}

			if oldEfiDisk != nil {
				baseDiskInterface := disk.DigitPrefix(*oldEfiDisk.Interface)
				diskOldEntries[baseDiskInterface][*oldEfiDisk.Interface] = oldEfiDisk
			}

			if newEfiDisk != nil {
				baseDiskInterface := disk.DigitPrefix(*newEfiDisk.Interface)
				diskNewEntries[baseDiskInterface][*newEfiDisk.Interface] = newEfiDisk
			}

			if oldEfiDisk != nil && newEfiDisk != nil && oldEfiDisk.Size != newEfiDisk.Size {
				return diag.Errorf(
					"resizing of efidisk is not supported.",
				)
			}
		}

		// Add tpm state if it has changes
		if d.HasChange(mkTPMState) {
			diskOld, diskNew := d.GetChange(mkTPMState)

			oldTPMState := vmGetTPMStateAsStorageDevice(d, diskOld.([]interface{}))
			newTPMState := vmGetTPMStateAsStorageDevice(d, diskNew.([]interface{}))

			if oldTPMState != nil {
				baseDiskInterface := disk.DigitPrefix(*oldTPMState.Interface)
				diskOldEntries[baseDiskInterface][*oldTPMState.Interface] = oldTPMState
			}

			if newTPMState != nil {
				baseDiskInterface := disk.DigitPrefix(*newTPMState.Interface)
				diskNewEntries[baseDiskInterface][*newTPMState.Interface] = newTPMState
			}

			if oldTPMState != nil && newTPMState != nil && oldTPMState.Size != newTPMState.Size {
				return diag.Errorf(
					"resizing of tpm state is not supported.",
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
						"deletion of disks not supported. Please delete disk by hand. Old interface was %q",
						*oldDisk.Interface,
					)
				}

				if *oldDisk.DatastoreID != *diskNewEntries[prefix][oldKey].DatastoreID {
					if oldDisk.IsOwnedBy(vmID) {
						deleteOriginalDisk := types.CustomBool(true)

						diskMoveBodies = append(
							diskMoveBodies,
							&vms.MoveDiskRequestBody{
								DeleteOriginalDisk: &deleteOriginalDisk,
								Disk:               *oldDisk.Interface,
								TargetStorage:      *diskNewEntries[prefix][oldKey].DatastoreID,
							},
						)

						// Cannot be done while VM is running.
						shutdownForDisksRequired = true
					} else {
						return diag.Errorf(
							"Cannot move %s:%s to datastore %s in VM %d configuration, it is not owned by this VM!",
							*oldDisk.DatastoreID,
							*oldDisk.PathInDatastore(),
							*diskNewEntries[prefix][oldKey].DatastoreID,
							vmID,
						)
					}
				}

				if *oldDisk.Size < *diskNewEntries[prefix][oldKey].Size {
					if oldDisk.IsOwnedBy(vmID) {
						diskResizeBodies = append(
							diskResizeBodies,
							&vms.ResizeDiskRequestBody{
								Disk: *oldDisk.Interface,
								Size: *diskNewEntries[prefix][oldKey].Size,
							},
						)
					} else {
						return diag.Errorf(
							"Cannot resize %s:%s in VM %d configuration, it is not owned by this VM!",
							*oldDisk.DatastoreID,
							*oldDisk.PathInDatastore(),
							vmID,
						)
					}
				}
			}
		}

		if shutdownForDisksRequired && !template {
			if e := vmShutdown(ctx, vmAPI, d); e != nil {
				return e
			}
		}

		for _, reqBody := range diskMoveBodies {
			err = vmAPI.MoveVMDisk(ctx, reqBody)
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
			if diags := vmStart(ctx, vmAPI, d); diags != nil {
				return diags
			}

			// This concludes an equivalent of a reboot, avoid doing another.
			reboot = false
		}
	}

	// Perform a regular reboot in case it's necessary and haven't already been done.
	if reboot {
		vmStatus, err := vmAPI.GetVMStatus(ctx)
		if err != nil {
			return diag.FromErr(err)
		}

		if vmStatus.Status != "stopped" {
			rebootTimeoutSec := d.Get(mkTimeoutReboot).(int)

			err := vmAPI.RebootVM(
				ctx,
				&vms.RebootRequestBody{
					Timeout: &rebootTimeoutSec,
				},
			)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return vmRead(ctx, d, m)
}

func vmDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	timeout := d.Get(mkTimeoutStopVM).(int)
	shutdownTimeout := d.Get(mkTimeoutShutdownVM).(int)

	if shutdownTimeout > timeout {
		timeout = shutdownTimeout
	}

	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	config := m.(proxmoxtf.ProviderConfiguration)

	client, err := config.GetClient()
	if err != nil {
		return diag.FromErr(err)
	}

	nodeName := d.Get(mkNodeName).(string)

	vmID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	vmAPI := client.Node(nodeName).VM(vmID)

	// Stop or shut down the virtual machine before deleting it.
	status, err := vmAPI.GetVMStatus(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	stop := d.Get(mkStopOnDestroy).(bool)

	//nolint: nestif
	if status.Status != "stopped" {
		if stop {
			if e := vmStop(ctx, vmAPI, d); e != nil {
				return e
			}
		} else {
			if e := vmShutdown(ctx, vmAPI, d); e != nil {
				return e
			}
		}
	}

	err = vmAPI.DeleteVM(ctx)
	if err != nil {
		if errors.Is(err, api.ErrResourceDoesNotExist) {
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	// Wait for the state to become unavailable as that clearly indicates the destruction of the VM.
	err = vmAPI.WaitForVMStatus(ctx, "")
	if err == nil {
		return diag.Errorf("failed to delete VM \"%d\"", vmID)
	}

	d.SetId("")

	return nil
}

// getDiskDatastores returns a list of the used datastores in a VM.
func getDiskDatastores(vm *vms.GetResponseData, d *schema.ResourceData) []string {
	storageDevices := disk.GetInfo(vm, d)
	datastoresSet := map[string]int{}

	for _, diskInfo := range storageDevices {
		// Ignore empty storage devices and storage devices (like ide) which may not have any media mounted
		if diskInfo == nil || diskInfo.FileVolume == "none" {
			continue
		}

		fileIDParts := strings.Split(diskInfo.FileVolume, ":")
		datastoresSet[fileIDParts[0]] = 1
	}

	if vm.EFIDisk != nil {
		fileIDParts := strings.Split(vm.EFIDisk.FileVolume, ":")
		datastoresSet[fileIDParts[0]] = 1
	}

	if vm.TPMState != nil {
		fileIDParts := strings.Split(vm.TPMState.FileVolume, ":")
		datastoresSet[fileIDParts[0]] = 1
	}

	datastores := []string{}
	for datastore := range datastoresSet {
		datastores = append(datastores, datastore)
	}

	return datastores
}

func getNUMAInfo(resp *vms.GetResponseData, _ *schema.ResourceData) map[string]*vms.CustomNUMADevice {
	numaDevices := map[string]*vms.CustomNUMADevice{}

	numaDevices["numa0"] = resp.NUMADevices0
	numaDevices["numa1"] = resp.NUMADevices1
	numaDevices["numa2"] = resp.NUMADevices2
	numaDevices["numa3"] = resp.NUMADevices3
	numaDevices["numa4"] = resp.NUMADevices4
	numaDevices["numa5"] = resp.NUMADevices5
	numaDevices["numa6"] = resp.NUMADevices6
	numaDevices["numa7"] = resp.NUMADevices7

	return numaDevices
}

func getPCIInfo(resp *vms.GetResponseData, _ *schema.ResourceData) map[string]*vms.CustomPCIDevice {
	pciDevices := map[string]*vms.CustomPCIDevice{}

	pciDevices["hostpci0"] = resp.PCIDevice0
	pciDevices["hostpci1"] = resp.PCIDevice1
	pciDevices["hostpci2"] = resp.PCIDevice2
	pciDevices["hostpci3"] = resp.PCIDevice3

	return pciDevices
}

func getUSBInfo(resp *vms.GetResponseData, _ *schema.ResourceData) map[string]*vms.CustomUSBDevice {
	usbDevices := map[string]*vms.CustomUSBDevice{}

	usbDevices["usb0"] = resp.USBDevice0
	usbDevices["usb1"] = resp.USBDevice1
	usbDevices["usb2"] = resp.USBDevice2
	usbDevices["usb3"] = resp.USBDevice3

	return usbDevices
}

func parseImportIDWithNodeName(id string) (string, string, error) {
	nodeName, id, found := strings.Cut(id, "/")

	if !found {
		return "", "", fmt.Errorf("unexpected format of ID (%s), expected node/id", id)
	}

	return nodeName, id, nil
}

func getAgentTimeout(d *schema.ResourceData) (time.Duration, error) {
	resource := VM()

	agentBlock, err := structure.GetSchemaBlock(
		resource,
		d,
		[]string{mkAgent},
		0,
		true,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to get agent block: %w", err)
	}

	agentTimeout, err := time.ParseDuration(
		agentBlock[mkAgentTimeout].(string),
	)
	if err != nil {
		return 0, fmt.Errorf("failed to parse agent timeout: %w", err)
	}

	return agentTimeout, nil
}
