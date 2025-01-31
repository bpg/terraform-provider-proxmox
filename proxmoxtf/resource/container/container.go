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
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/containers"
	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
	"github.com/bpg/terraform-provider-proxmox/proxmoxtf"
	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/resource/validators"
	resource "github.com/bpg/terraform-provider-proxmox/proxmoxtf/resource/vm"
	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/structure"
	"github.com/bpg/terraform-provider-proxmox/utils"
)

const (
	dvCloneDatastoreID                  = ""
	dvCloneNodeName                     = ""
	dvConsoleEnabled                    = true
	dvConsoleMode                       = "tty"
	dvConsoleTTYCount                   = 2
	dvInitializationDNSDomain           = ""
	dvInitializationDNSServer           = ""
	dvInitializationIPConfigIPv4Address = ""
	dvInitializationIPConfigIPv4Gateway = ""
	dvInitializationIPConfigIPv6Address = ""
	dvInitializationIPConfigIPv6Gateway = ""
	dvInitializationHostname            = ""
	dvInitializationUserAccountPassword = ""
	dvCPUArchitecture                   = "amd64"
	dvCPUCores                          = 1
	dvCPUUnits                          = 1024
	dvDescription                       = ""
	dvDiskDatastoreID                   = "local"
	dvDiskSize                          = 4
	dvFeaturesNesting                   = false
	dvFeaturesKeyControl                = false
	dvFeaturesFUSE                      = false
	dvHookScript                        = ""
	dvMemoryDedicated                   = 512
	dvMemorySwap                        = 0
	dvMountPointACL                     = false
	dvMountPointBackup                  = false
	dvMountPointPath                    = ""
	dvMountPointQuota                   = false
	dvMountPointReadOnly                = false
	dvMountPointReplicate               = true
	dvMountPointShared                  = false
	dvMountPointSize                    = ""
	dvNetworkInterfaceBridge            = "vmbr0"
	dvNetworkInterfaceEnabled           = true
	dvNetworkInterfaceFirewall          = false
	dvNetworkInterfaceMACAddress        = ""
	dvNetworkInterfaceRateLimit         = 0
	dvNetworkInterfaceVLANID            = 0
	dvNetworkInterfaceMTU               = 0
	dvOperatingSystemType               = "unmanaged"
	dvPoolID                            = ""
	dvProtection                        = false
	dvStarted                           = true
	dvStartupOrder                      = -1
	dvStartupUpDelay                    = -1
	dvStartupDownDelay                  = -1
	dvStartOnBoot                       = true
	dvTemplate                          = false
	dvTimeoutCreate                     = 1800
	dvTimeoutClone                      = 1800
	dvTimeoutUpdate                     = 1800
	dvTimeoutDelete                     = 60
	dvUnprivileged                      = false

	maxResourceVirtualEnvironmentContainerNetworkInterfaces = 8

	mkClone                             = "clone"
	mkCloneDatastoreID                  = "datastore_id"
	mkCloneNodeName                     = "node_name"
	mkCloneVMID                         = "vm_id"
	mkConsole                           = "console"
	mkConsoleEnabled                    = "enabled"
	mkConsoleMode                       = "type"
	mkConsoleTTYCount                   = "tty_count"
	mkCPU                               = "cpu"
	mkCPUArchitecture                   = "architecture"
	mkCPUCores                          = "cores"
	mkCPUUnits                          = "units"
	mkDescription                       = "description"
	mkDisk                              = "disk"
	mkDiskDatastoreID                   = "datastore_id"
	mkDiskSize                          = "size"
	mkFeatures                          = "features"
	mkFeaturesNesting                   = "nesting"
	mkFeaturesKeyControl                = "keyctl"
	mkFeaturesFUSE                      = "fuse"
	mkFeaturesMountTypes                = "mount"
	mkHookScriptFileID                  = "hook_script_file_id"
	mkInitialization                    = "initialization"
	mkInitializationDNS                 = "dns"
	mkInitializationDNSDomain           = "domain"
	mkInitializationDNSServer           = "server"
	mkInitializationDNSServers          = "servers"
	mkInitializationHostname            = "hostname"
	mkInitializationIPConfig            = "ip_config"
	mkInitializationIPConfigIPv4        = "ipv4"
	mkInitializationIPConfigIPv4Address = "address"
	mkInitializationIPConfigIPv4Gateway = "gateway"
	mkInitializationIPConfigIPv6        = "ipv6"
	mkInitializationIPConfigIPv6Address = "address"
	mkInitializationIPConfigIPv6Gateway = "gateway"
	mkInitializationUserAccount         = "user_account"
	mkInitializationUserAccountKeys     = "keys"
	mkInitializationUserAccountPassword = "password"
	mkMemory                            = "memory"
	mkMemoryDedicated                   = "dedicated"
	mkMemorySwap                        = "swap"
	mkMountPoint                        = "mount_point"
	mkMountPointACL                     = "acl"
	mkMountPointBackup                  = "backup"
	mkMountPointMountOptions            = "mount_options"
	mkMountPointPath                    = "path"
	mkMountPointQuota                   = "quota"
	mkMountPointReadOnly                = "read_only"
	mkMountPointReplicate               = "replicate"
	mkMountPointShared                  = "shared"
	mkMountPointSize                    = "size"
	mkMountPointVolume                  = "volume"
	mkDevicePassthroughDenyWrite        = "deny_write"
	mkDevicePassthrough                 = "device_passthrough" // #nosec G101
	mkDevicePassthroughPath             = "path"
	mkDevicePassthroughUID              = "uid"
	mkDevicePassthroughGID              = "gid"
	mkDevicePassthroughMode             = "mode"
	mkNetworkInterface                  = "network_interface"
	mkNetworkInterfaceBridge            = "bridge"
	mkNetworkInterfaceEnabled           = "enabled"
	mkNetworkInterfaceFirewall          = "firewall"
	mkNetworkInterfaceMACAddress        = "mac_address"
	mkNetworkInterfaceName              = "name"
	mkNetworkInterfaceRateLimit         = "rate_limit"
	mkNetworkInterfaceVLANID            = "vlan_id"
	mkNetworkInterfaceMTU               = "mtu"
	mkNodeName                          = "node_name"
	mkOperatingSystem                   = "operating_system"
	mkOperatingSystemTemplateFileID     = "template_file_id"
	mkOperatingSystemType               = "type"
	mkPoolID                            = "pool_id"
	mkProtection                        = "protection"
	mkStarted                           = "started"
	mkStartup                           = "startup"
	mkStartupOrder                      = "order"
	mkStartupUpDelay                    = "up_delay"
	mkStartupDownDelay                  = "down_delay"
	mkStartOnBoot                       = "start_on_boot"
	mkTags                              = "tags"
	mkTemplate                          = "template"
	mkTimeoutCreate                     = "timeout_create"
	mkTimeoutClone                      = "timeout_clone"
	mkTimeoutUpdate                     = "timeout_update"
	mkTimeoutDelete                     = "timeout_delete"
	mkUnprivileged                      = "unprivileged"
	mkVMID                              = "vm_id"
)

// Container returns a resource that manages a container.
func Container() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			mkClone: {
				Type:        schema.TypeList,
				Description: "The cloning configuration",
				Optional:    true,
				DefaultFunc: func() (interface{}, error) {
					return []interface{}{}, nil
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
							Description:      "The ID of the source container",
							Required:         true,
							ForceNew:         true,
							ValidateDiagFunc: resource.VMIDValidator(),
						},
					},
				},
				MaxItems: 1,
				MinItems: 0,
			},
			mkConsole: {
				Type:        schema.TypeList,
				Description: "The console configuration",
				Optional:    true,
				DefaultFunc: func() (interface{}, error) {
					return []interface{}{
						map[string]interface{}{
							mkConsoleEnabled:  dvConsoleEnabled,
							mkConsoleMode:     dvConsoleMode,
							mkConsoleTTYCount: dvConsoleTTYCount,
						},
					}, nil
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mkConsoleEnabled: {
							Type:        schema.TypeBool,
							Description: "Whether to enable the console device",
							Optional:    true,
							Default:     dvConsoleEnabled,
						},
						mkConsoleMode: {
							Type:             schema.TypeString,
							Description:      "The console mode",
							Optional:         true,
							Default:          dvConsoleMode,
							ValidateDiagFunc: ConsoleModeValidator(),
						},
						mkConsoleTTYCount: {
							Type:             schema.TypeInt,
							Description:      "The number of available TTY",
							Optional:         true,
							Default:          dvConsoleTTYCount,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(0, 6)),
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
							ValidateDiagFunc: CPUArchitectureValidator(),
						},
						mkCPUCores: {
							Type:             schema.TypeInt,
							Description:      "The number of CPU cores",
							Optional:         true,
							Default:          dvCPUCores,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(1, 128)),
						},
						mkCPUUnits: {
							Type:        schema.TypeInt,
							Description: "The CPU units",
							Optional:    true,
							Default:     dvCPUUnits,
							ValidateDiagFunc: validation.ToDiagFunc(
								validation.IntBetween(0, 500000),
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
				StateFunc: func(i interface{}) string {
					// PVE always adds a newline to the description, so we have to do the same,
					// also taking in account the CLRF case (Windows)
					if i.(string) != "" {
						return strings.ReplaceAll(strings.TrimSpace(i.(string)), "\r\n", "\n") + "\n"
					}
					return ""
				},
			},
			mkDisk: {
				Type:        schema.TypeList,
				Description: "The disks",
				Optional:    true,
				ForceNew:    true,
				DefaultFunc: func() (interface{}, error) {
					return []interface{}{
						map[string]interface{}{
							mkDiskDatastoreID: dvDiskDatastoreID,
							mkDiskSize:        dvDiskSize,
						},
					}, nil
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mkDiskDatastoreID: {
							Type:        schema.TypeString,
							Description: "The datastore id",
							Optional:    true,
							ForceNew:    true,
							Default:     dvDiskDatastoreID,
						},
						mkDiskSize: {
							Type:             schema.TypeInt,
							Description:      "The rootfs size in gigabytes",
							Optional:         true,
							ForceNew:         true,
							Default:          dvDiskSize,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(0)),
						},
					},
				},
				MaxItems: 1,
				MinItems: 0,
			},
			mkFeatures: {
				Type:        schema.TypeList,
				Description: "Features",
				Optional:    true,
				DefaultFunc: func() (interface{}, error) {
					return []interface{}{
						map[string]interface{}{
							mkFeaturesNesting:    dvFeaturesNesting,
							mkFeaturesKeyControl: dvFeaturesKeyControl,
							mkFeaturesFUSE:       dvFeaturesFUSE,
							mkFeaturesMountTypes: []interface{}{},
						},
					}, nil
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mkFeaturesNesting: {
							Type:        schema.TypeBool,
							Description: "Whether the container runs as nested",
							Optional:    true,
							Default:     dvFeaturesNesting,
						},
						mkFeaturesKeyControl: {
							Type:        schema.TypeBool,
							Description: "Whether the container supports `keyctl()` system call",
							Optional:    true,
							Default:     dvFeaturesKeyControl,
						},
						mkFeaturesFUSE: {
							Type:        schema.TypeBool,
							Description: "Whether the container supports FUSE mounts",
							Optional:    true,
							Default:     dvFeaturesFUSE,
						},
						mkFeaturesMountTypes: {
							Type:        schema.TypeList,
							Description: "List of allowed mount types",
							Optional:    true,
							Elem: &schema.Schema{
								Type:             schema.TypeString,
								ValidateDiagFunc: MountTypeValidator(),
							},
						},
					},
				},
				MaxItems: 1,
				MinItems: 0,
			},
			mkHookScriptFileID: {
				Type:        schema.TypeString,
				Description: "A hook script",
				Optional:    true,
				Default:     dvHookScript,
			},
			mkInitialization: {
				Type:        schema.TypeList,
				Description: "The initialization configuration",
				Optional:    true,
				DefaultFunc: func() (interface{}, error) {
					return []interface{}{}, nil
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
						mkInitializationHostname: {
							Type:        schema.TypeString,
							Description: "The hostname",
							Optional:    true,
							Default:     dvInitializationHostname,
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
										DefaultFunc: func() (interface{}, error) {
											return []interface{}{}, nil
										},
										Elem: &schema.Schema{Type: schema.TypeString},
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
			mkMemory: {
				Type:        schema.TypeList,
				Description: "The memory allocation",
				Optional:    true,
				DefaultFunc: func() (interface{}, error) {
					return []interface{}{
						map[string]interface{}{
							mkMemoryDedicated: dvMemoryDedicated,
							mkMemorySwap:      dvMemorySwap,
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
								validation.IntBetween(16, 268435456),
							),
						},
						mkMemorySwap: {
							Type:        schema.TypeInt,
							Description: "The swap size in megabytes",
							Optional:    true,
							Default:     dvMemorySwap,
							ValidateDiagFunc: validation.ToDiagFunc(
								validation.IntBetween(0, 268435456),
							),
						},
					},
				},
				MaxItems: 1,
				MinItems: 0,
			},
			mkMountPoint: {
				Type:        schema.TypeList,
				Description: "A mount point",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mkMountPointACL: {
							Type:        schema.TypeBool,
							Description: "Explicitly enable or disable ACL support",
							Optional:    true,
							Default:     dvMountPointACL,
						},
						mkMountPointBackup: {
							Type:        schema.TypeBool,
							Description: "Whether to include the mount point in backups (only used for volume mount points)",
							Optional:    true,
							Default:     dvMountPointBackup,
						},
						mkMountPointMountOptions: {
							Type:        schema.TypeList,
							Description: "Extra mount options.",
							Optional:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						mkMountPointPath: {
							Type:        schema.TypeString,
							Description: "Path to the mount point as seen from inside the container",
							Required:    true,
							// StateFunc: func(i interface{}) string {
							// 	// PVE strips leading slashes from the path, so we have to do the same
							// 	return strings.TrimPrefix(i.(string), "/")
							// },
							DiffSuppressFunc: func(_, oldVal, newVal string, _ *schema.ResourceData) bool {
								return "/"+oldVal == newVal
							},
						},
						mkMountPointQuota: {
							Type:        schema.TypeBool,
							Description: "Enable user quotas inside the container (not supported with volume mounts)",
							Optional:    true,
							Default:     dvMountPointQuota,
						},
						mkMountPointReadOnly: {
							Type:        schema.TypeBool,
							Description: "Read-only mount point",
							Optional:    true,
							Default:     dvMountPointReadOnly,
						},
						mkMountPointReplicate: {
							Type:        schema.TypeBool,
							Description: "Will include this volume to a storage replica job",
							Optional:    true,
							Default:     dvMountPointReplicate,
						},
						mkMountPointShared: {
							Type:        schema.TypeBool,
							Description: "Mark this non-volume mount point as available on all nodes",
							Optional:    true,
							Default:     dvMountPointShared,
						},
						mkMountPointSize: {
							Type:             schema.TypeString,
							Description:      "Volume size (only used for volume mount points)",
							Optional:         true,
							Default:          dvMountPointSize,
							ValidateDiagFunc: validators.FileSize(),
						},
						mkMountPointVolume: {
							Type:        schema.TypeString,
							Description: "Volume, device or directory to mount into the container",
							Required:    true,
							DiffSuppressFunc: func(_, oldVal, newVal string, _ *schema.ResourceData) bool {
								// For *new* volume mounts PVE returns an actual volume ID which is saved in the stare,
								// so on reapply the provider will try override it:"
								//   "local-lvm" -> "local-lvm:vm-101-disk-1"
								//   "local-lvm:8" -> "local-lvm:vm-101-disk-1"
								// There is also an option to mount an existing volume, so
								//   "local-lvm:vm-101-disk-1" -> "local-lvm:vm-101-disk-1"
								// which is a valid case.
								return oldVal == newVal || strings.HasPrefix(oldVal, strings.Split(newVal, ":")[0]+":")
							},
						},
					},
				},
				MaxItems: 8,
				MinItems: 0,
			},
			mkDevicePassthrough: {
				Type:        schema.TypeList,
				Description: "Device to pass through to the container",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mkDevicePassthroughDenyWrite: {
							Type:        schema.TypeBool,
							Description: "Deny the container to write to the device",
							Optional:    true,
							Default:     false,
						},
						mkDevicePassthroughGID: {
							Type:             schema.TypeInt,
							Description:      "Group ID to be assigned to the device node",
							Optional:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(0)),
						},
						mkDevicePassthroughMode: {
							Type:        schema.TypeString,
							Description: "Access mode to be set on the device node (e.g. 0666)",
							Optional:    true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringMatch(
								regexp.MustCompile(`0[0-7]{3}`), "Octal access mode",
							)),
						},
						mkDevicePassthroughPath: {
							Type:             schema.TypeString,
							Description:      "Device to pass through to the container",
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
						},
						mkDevicePassthroughUID: {
							Type:             schema.TypeInt,
							Description:      "Device UID in the container",
							Optional:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(0)),
						},
					},
				},
				MaxItems: 8,
				MinItems: 0,
			},
			mkNetworkInterface: {
				Type:        schema.TypeList,
				Description: "The network interfaces",
				Optional:    true,
				DefaultFunc: func() (interface{}, error) {
					return make([]interface{}, 1), nil
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mkNetworkInterfaceBridge: {
							Type:        schema.TypeString,
							Description: "The bridge",
							Optional:    true,
							Default:     dvNetworkInterfaceBridge,
						},
						mkNetworkInterfaceEnabled: {
							Type:        schema.TypeBool,
							Description: "Whether to enable the network device",
							Optional:    true,
							Default:     dvNetworkInterfaceEnabled,
						},
						mkNetworkInterfaceFirewall: {
							Type:        schema.TypeBool,
							Description: "Whether this interface's firewall rules should be used.",
							Optional:    true,
							Default:     dvNetworkInterfaceFirewall,
						},
						mkNetworkInterfaceMACAddress: {
							Type:        schema.TypeString,
							Description: "The MAC address",
							Optional:    true,
							Default:     dvNetworkInterfaceMACAddress,
							DiffSuppressFunc: func(_, _, newVal string, _ *schema.ResourceData) bool {
								return newVal == ""
							},
							ValidateDiagFunc: validators.MACAddress(),
						},
						mkNetworkInterfaceName: {
							Type:        schema.TypeString,
							Description: "The network interface name",
							Required:    true,
						},
						mkNetworkInterfaceRateLimit: {
							Type:        schema.TypeFloat,
							Description: "The rate limit in megabytes per second",
							Optional:    true,
							Default:     dvNetworkInterfaceRateLimit,
						},
						mkNetworkInterfaceVLANID: {
							Type:        schema.TypeInt,
							Description: "The VLAN identifier",
							Optional:    true,
							Default:     dvNetworkInterfaceVLANID,
						},
						mkNetworkInterfaceMTU: {
							Type:        schema.TypeInt,
							Description: "Maximum transmission unit (MTU)",
							Optional:    true,
							Default:     dvNetworkInterfaceMTU,
						},
					},
				},
				MaxItems: maxResourceVirtualEnvironmentContainerNetworkInterfaces,
				MinItems: 0,
			},
			mkNodeName: {
				Type:        schema.TypeString,
				Description: "The node name",
				Required:    true,
				ForceNew:    true,
			},
			mkOperatingSystem: {
				Type:        schema.TypeList,
				Description: "The operating system configuration",
				Optional:    true,
				ForceNew:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mkOperatingSystemTemplateFileID: {
							Type:             schema.TypeString,
							Description:      "The ID of an OS template file",
							Required:         true,
							ForceNew:         true,
							ValidateDiagFunc: validators.FileID(),
						},
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
				Description: "The ID of the pool to assign the container to",
				Optional:    true,
				ForceNew:    true,
				Default:     dvPoolID,
			},
			mkProtection: {
				Type: schema.TypeBool,
				Description: "Whether to set the protection flag of the container. " +
					"This will prevent the container itself and its disk for remove/update operations.",
				Optional: true,
				ForceNew: false,
				Default:  dvProtection,
			},
			mkStarted: {
				Type:        schema.TypeBool,
				Description: "Whether to start the container",
				Optional:    true,
				Default:     dvStarted,
				DiffSuppressFunc: func(_, _, _ string, d *schema.ResourceData) bool {
					return d.Get(mkTemplate).(bool)
				},
			},
			mkStartup: {
				Type:        schema.TypeList,
				Description: "Defines startup and shutdown behavior of the container",
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
							Description: "A non-negative number defining the delay in seconds before the next container is started",
							Optional:    true,
							Default:     dvStartupUpDelay,
						},
						mkStartupDownDelay: {
							Type:        schema.TypeInt,
							Description: "A non-negative number defining the delay in seconds before the next container is shut down",
							Optional:    true,
							Default:     dvStartupDownDelay,
						},
					},
				},
				MaxItems: 1,
				MinItems: 0,
			},
			mkStartOnBoot: {
				Type:        schema.TypeBool,
				Description: "Automatically start container when the host system boots.",
				Optional:    true,
				ForceNew:    false,
				Default:     dvStartOnBoot,
			},
			mkTags: {
				Type:        schema.TypeList,
				Description: "Tags of the container. This is only meta information.",
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
			mkTimeoutCreate: {
				Type:        schema.TypeInt,
				Description: "Create container timeout",
				Optional:    true,
				Default:     dvTimeoutCreate,
			},
			mkTimeoutClone: {
				Type:        schema.TypeInt,
				Description: "Clone container timeout",
				Optional:    true,
				Default:     dvTimeoutClone,
			},
			mkTimeoutUpdate: {
				Type:        schema.TypeInt,
				Description: "Update container timeout",
				Optional:    true,
				Default:     dvTimeoutUpdate,
			},
			mkTimeoutDelete: {
				Type:        schema.TypeInt,
				Description: "Delete container timeout",
				Optional:    true,
				Default:     dvTimeoutDelete,
			},
			"timeout_start": {
				Type:        schema.TypeInt,
				Description: "Start container timeout",
				Optional:    true,
				Default:     300,
				Deprecated: "This field is deprecated and will be removed in a future release. " +
					"An overall operation timeout (`timeout_create` / `timeout_clone`) is used instead.",
			},
			mkUnprivileged: {
				Type:        schema.TypeBool,
				Description: "Whether the container runs as unprivileged on the host",
				Optional:    true,
				ForceNew:    true,
				Default:     dvUnprivileged,
			},
			mkVMID: {
				Type:             schema.TypeInt,
				Description:      "The VM identifier",
				Optional:         true,
				Computed:         true,
				ValidateDiagFunc: resource.VMIDValidator(),
			},
		},
		CreateContext: containerCreate,
		ReadContext:   containerRead,
		UpdateContext: containerUpdate,
		DeleteContext: containerDelete,
		CustomizeDiff: customdiff.All(
			customdiff.ForceNewIf(
				mkVMID,
				func(_ context.Context, d *schema.ResourceDiff, _ interface{}) bool {
					newValue := d.Get(mkVMID)

					// 'vm_id' is ForceNew, except when changing 'vm_id' to existing correct id
					// (automatic fix from -1 to actual vm_id must not re-create VM)
					return strconv.Itoa(newValue.(int)) != d.Id()
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

func containerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clone := d.Get(mkClone).([]interface{})

	if len(clone) > 0 {
		return containerCreateClone(ctx, d, m)
	}

	return containerCreateCustom(ctx, d, m)
}

func containerCreateClone(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cloneTimeoutSec := d.Get(mkTimeoutClone).(int)

	ctx, cancel := context.WithTimeout(ctx, time.Duration(cloneTimeoutSec)*time.Second)
	defer cancel()

	config := m.(proxmoxtf.ProviderConfiguration)

	client, err := config.GetClient()
	if err != nil {
		return diag.FromErr(err)
	}

	clone := d.Get(mkClone).([]interface{})
	cloneBlock := clone[0].(map[string]interface{})
	cloneDatastoreID := cloneBlock[mkCloneDatastoreID].(string)
	cloneNodeName := cloneBlock[mkCloneNodeName].(string)
	cloneVMID := cloneBlock[mkCloneVMID].(int)

	description := d.Get(mkDescription).(string)

	initialization := d.Get(mkInitialization).([]interface{})
	initializationHostname := ""

	if len(initialization) > 0 && initialization[0] != nil {
		initializationBlock := initialization[0].(map[string]interface{})
		initializationHostname = initializationBlock[mkInitializationHostname].(string)
	}

	nodeName := d.Get(mkNodeName).(string)
	poolID := d.Get(mkPoolID).(string)
	tags := d.Get(mkTags).([]interface{})
	vmIDUntyped, hasVMID := d.GetOk(mkVMID)
	vmID := vmIDUntyped.(int)

	if !hasVMID {
		vmIDNew, err := config.GetIDGenerator().NextID(ctx)
		if err != nil {
			return diag.FromErr(err)
		}

		vmID = vmIDNew

		err = d.Set(mkVMID, vmID)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	fullCopy := types.CustomBool(true)

	cloneBody := &containers.CloneRequestBody{
		FullCopy: &fullCopy,
		VMIDNew:  vmID,
	}

	if cloneDatastoreID != "" {
		cloneBody.TargetStorage = &cloneDatastoreID
	}

	if description != "" {
		cloneBody.Description = &description
	}

	if initializationHostname != "" {
		cloneBody.Hostname = &initializationHostname
	}

	if poolID != "" {
		cloneBody.PoolID = &poolID
	}

	if cloneNodeName != "" && cloneNodeName != nodeName {
		cloneBody.TargetNodeName = &nodeName

		err = client.Node(cloneNodeName).Container(cloneVMID).CloneContainer(ctx, cloneBody)
	} else {
		err = client.Node(nodeName).Container(cloneVMID).CloneContainer(ctx, cloneBody)
	}

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(vmID))

	containerAPI := client.Node(nodeName).Container(vmID)

	// Wait for the container to be created and its configuration lock to be released.
	err = containerAPI.WaitForContainerConfigUnlock(ctx, true)
	if err != nil {
		return diag.FromErr(err)
	}

	// Now that the virtual machine has been cloned, we need to perform some modifications.
	updateBody := &containers.UpdateRequestBody{}

	startOnBoot := types.CustomBool(d.Get(mkStartOnBoot).(bool))
	updateBody.StartOnBoot = &startOnBoot

	protection := types.CustomBool(d.Get(mkProtection).(bool))
	updateBody.Protection = &protection

	updateBody.StartupBehavior = containerGetStartupBehavior(d)

	console := d.Get(mkConsole).([]interface{})

	if len(console) > 0 && console[0] != nil {
		consoleBlock := console[0].(map[string]interface{})

		consoleEnabled := types.CustomBool(
			consoleBlock[mkConsoleEnabled].(bool),
		)
		consoleMode := consoleBlock[mkConsoleMode].(string)
		consoleTTYCount := consoleBlock[mkConsoleTTYCount].(int)

		updateBody.ConsoleEnabled = &consoleEnabled
		updateBody.ConsoleMode = &consoleMode
		updateBody.TTY = &consoleTTYCount
	}

	cpu := d.Get(mkCPU).([]interface{})

	if len(cpu) > 0 && cpu[0] != nil {
		cpuBlock := cpu[0].(map[string]interface{})

		cpuArchitecture := cpuBlock[mkCPUArchitecture].(string)
		cpuCores := cpuBlock[mkCPUCores].(int)
		cpuUnits := cpuBlock[mkCPUUnits].(int)

		updateBody.CPUArchitecture = &cpuArchitecture
		updateBody.CPUCores = &cpuCores
		updateBody.CPUUnits = &cpuUnits
	}

	hookScript := d.Get(mkHookScriptFileID).(string)

	if hookScript != "" {
		updateBody.HookScript = &hookScript
	}

	var initializationIPConfigIPv4Address []string

	var initializationIPConfigIPv4Gateway []string

	var initializationIPConfigIPv6Address []string

	var initializationIPConfigIPv6Gateway []string

	if len(initialization) > 0 && initialization[0] != nil {
		initializationBlock := initialization[0].(map[string]interface{})
		initializationDNS := initializationBlock[mkInitializationDNS].([]interface{})

		if len(initializationDNS) > 0 && initializationDNS[0] != nil {
			initializationDNSBlock := initializationDNS[0].(map[string]interface{})
			initializationDNSDomain := initializationDNSBlock[mkInitializationDNSDomain].(string)
			updateBody.DNSDomain = &initializationDNSDomain

			servers := initializationDNSBlock[mkInitializationDNSServers].([]interface{})
			deprecatedServer := initializationDNSBlock[mkInitializationDNSServer].(string)

			if len(servers) > 0 {
				nameserver := strings.Join(utils.ConvertToStringSlice(servers), " ")

				updateBody.DNSServer = &nameserver
			} else {
				updateBody.DNSServer = &deprecatedServer
			}
		}

		initializationHostname := initializationBlock[mkInitializationHostname].(string)

		if initializationHostname != dvInitializationHostname {
			updateBody.Hostname = &initializationHostname
		}

		initializationIPConfig := initializationBlock[mkInitializationIPConfig].([]interface{})

		for _, c := range initializationIPConfig {
			configBlock := c.(map[string]interface{})
			ipv4 := configBlock[mkInitializationIPConfigIPv4].([]interface{})

			if len(ipv4) > 0 && ipv4[0] != nil {
				ipv4Block := ipv4[0].(map[string]interface{})

				initializationIPConfigIPv4Address = append(
					initializationIPConfigIPv4Address,
					ipv4Block[mkInitializationIPConfigIPv4Address].(string),
				)

				initializationIPConfigIPv4Gateway = append(
					initializationIPConfigIPv4Gateway,
					ipv4Block[mkInitializationIPConfigIPv4Gateway].(string),
				)
			} else {
				initializationIPConfigIPv4Address = append(initializationIPConfigIPv4Address, "")
				initializationIPConfigIPv4Gateway = append(initializationIPConfigIPv4Gateway, "")
			}

			ipv6 := configBlock[mkInitializationIPConfigIPv6].([]interface{})

			if len(ipv6) > 0 && ipv6[0] != nil {
				ipv6Block := ipv6[0].(map[string]interface{})

				initializationIPConfigIPv6Address = append(
					initializationIPConfigIPv6Address,
					ipv6Block[mkInitializationIPConfigIPv6Address].(string),
				)

				initializationIPConfigIPv6Gateway = append(
					initializationIPConfigIPv6Gateway,
					ipv6Block[mkInitializationIPConfigIPv6Gateway].(string),
				)
			} else {
				initializationIPConfigIPv6Address = append(initializationIPConfigIPv6Address, "")
				initializationIPConfigIPv6Gateway = append(initializationIPConfigIPv6Gateway, "")
			}
		}

		initializationUserAccount := initializationBlock[mkInitializationUserAccount].([]interface{})

		if len(initializationUserAccount) > 0 && initializationUserAccount[0] != nil {
			initializationUserAccountBlock := initializationUserAccount[0].(map[string]interface{})
			keys := initializationUserAccountBlock[mkInitializationUserAccountKeys].([]interface{})

			if len(keys) > 0 {
				initializationUserAccountKeys := make(
					containers.CustomSSHKeys,
					len(keys),
				)

				for ki, kv := range keys {
					initializationUserAccountKeys[ki] = kv.(string)
				}

				updateBody.SSHKeys = &initializationUserAccountKeys
			} else {
				updateBody.Delete = append(updateBody.Delete, "ssh-public-keys")
			}

			initializationUserAccountPassword := initializationUserAccountBlock[mkInitializationUserAccountPassword].(string)

			if initializationUserAccountPassword != dvInitializationUserAccountPassword {
				updateBody.Password = &initializationUserAccountPassword
			} else {
				updateBody.Delete = append(updateBody.Delete, "password")
			}
		}
	}

	memory := d.Get(mkMemory).([]interface{})

	if len(memory) > 0 && memory[0] != nil {
		memoryBlock := memory[0].(map[string]interface{})

		memoryDedicated := memoryBlock[mkMemoryDedicated].(int)
		memorySwap := memoryBlock[mkMemorySwap].(int)

		updateBody.DedicatedMemory = &memoryDedicated
		updateBody.Swap = &memorySwap
	}

	devicePassthrough := d.Get(mkDevicePassthrough).([]interface{})

	devicePassthroughArray := make(
		containers.CustomDevicePassthroughArray,
		len(devicePassthrough),
	)

	for di, dv := range devicePassthrough {
		devicePassthroughMap := dv.(map[string]interface{})
		devicePassthroughObject := containers.CustomDevicePassthrough{}

		denyWrite := types.CustomBool(
			devicePassthroughMap[mkDevicePassthroughDenyWrite].(bool),
		)
		gid := devicePassthroughMap[mkDevicePassthroughGID].(int)
		mode := devicePassthroughMap[mkDevicePassthroughMode].(string)
		path := devicePassthroughMap[mkDevicePassthroughPath].(string)
		uid := devicePassthroughMap[mkDevicePassthroughUID].(int)

		devicePassthroughObject.DenyWrite = &denyWrite
		devicePassthroughObject.GID = &gid
		devicePassthroughObject.Mode = &mode
		devicePassthroughObject.Path = path
		devicePassthroughObject.UID = &uid

		devicePassthroughArray[di] = devicePassthroughObject
	}

	updateBody.DevicePassthrough = devicePassthroughArray

	networkInterface := d.Get(mkNetworkInterface).([]interface{})

	if len(networkInterface) == 0 {
		networkInterface, err = containerGetExistingNetworkInterface(ctx, containerAPI)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	networkInterfaceArray := make(
		containers.CustomNetworkInterfaceArray,
		len(networkInterface),
	)

	for ni, nv := range networkInterface {
		networkInterfaceMap := nv.(map[string]interface{})
		networkInterfaceObject := containers.CustomNetworkInterface{}

		bridge := networkInterfaceMap[mkNetworkInterfaceBridge].(string)
		enabled := networkInterfaceMap[mkNetworkInterfaceEnabled].(bool)
		firewall := types.CustomBool(
			networkInterfaceMap[mkNetworkInterfaceFirewall].(bool),
		)
		macAddress := networkInterfaceMap[mkNetworkInterfaceMACAddress].(string)
		name := networkInterfaceMap[mkNetworkInterfaceName].(string)
		rateLimit := networkInterfaceMap[mkNetworkInterfaceRateLimit].(float64)
		vlanID := networkInterfaceMap[mkNetworkInterfaceVLANID].(int)
		mtu, _ := networkInterfaceMap[mkNetworkInterfaceMTU].(int)

		if bridge != "" {
			networkInterfaceObject.Bridge = &bridge
		}

		networkInterfaceObject.Enabled = enabled
		networkInterfaceObject.Firewall = &firewall

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

		networkInterfaceObject.Name = name

		if rateLimit != 0 {
			networkInterfaceObject.RateLimit = &rateLimit
		}

		if vlanID != 0 {
			networkInterfaceObject.Tag = &vlanID
		}

		if mtu != 0 {
			networkInterfaceObject.MTU = &mtu
		}

		networkInterfaceArray[ni] = networkInterfaceObject
	}

	updateBody.NetworkInterfaces = networkInterfaceArray

	for i, ni := range updateBody.NetworkInterfaces {
		if !ni.Enabled {
			updateBody.Delete = append(updateBody.Delete, fmt.Sprintf("net%d", i))
		}
	}

	for i := len(updateBody.NetworkInterfaces); i < maxResourceVirtualEnvironmentContainerNetworkInterfaces; i++ {
		updateBody.Delete = append(updateBody.Delete, fmt.Sprintf("net%d", i))
	}

	operatingSystem := d.Get(mkOperatingSystem).([]interface{})

	if len(operatingSystem) > 0 && operatingSystem[0] != nil {
		operatingSystemBlock := operatingSystem[0].(map[string]interface{})

		operatingSystemTemplateFileID := operatingSystemBlock[mkOperatingSystemTemplateFileID].(string)
		operatingSystemType := operatingSystemBlock[mkOperatingSystemType].(string)

		updateBody.OSTemplateFileVolume = &operatingSystemTemplateFileID
		updateBody.OSType = &operatingSystemType
	}

	if len(tags) > 0 {
		tagString := containerGetTagsString(d)
		updateBody.Tags = &tagString
	}

	template := types.CustomBool(d.Get(mkTemplate).(bool))

	//nolint:gosimple
	if template != dvTemplate {
		updateBody.Template = &template
	}

	err = containerAPI.UpdateContainer(ctx, updateBody)
	if err != nil {
		return diag.FromErr(err)
	}

	// Wait for the container's lock to be released.
	err = containerAPI.WaitForContainerConfigUnlock(ctx, true)
	if err != nil {
		return diag.FromErr(err)
	}

	return containerCreateStart(ctx, d, m)
}

func containerCreateCustom(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	createTimeoutSec := d.Get(mkTimeoutCreate).(int)

	ctx, cancel := context.WithTimeout(ctx, time.Duration(createTimeoutSec)*time.Second)
	defer cancel()

	config := m.(proxmoxtf.ProviderConfiguration)

	client, err := config.GetClient()
	if err != nil {
		return diag.FromErr(err)
	}

	nodeName := d.Get(mkNodeName).(string)
	container := Container()

	consoleBlock, err := structure.GetSchemaBlock(
		container,
		d,
		[]string{mkConsole},
		0,
		true,
	)
	if err != nil {
		return diag.FromErr(err)
	}

	consoleEnabled := types.CustomBool(
		consoleBlock[mkConsoleEnabled].(bool),
	)
	consoleMode := consoleBlock[mkConsoleMode].(string)
	consoleTTYCount := consoleBlock[mkConsoleTTYCount].(int)

	cpuBlock, err := structure.GetSchemaBlock(
		container,
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
	cpuUnits := cpuBlock[mkCPUUnits].(int)

	description := d.Get(mkDescription).(string)

	diskBlock, err := structure.GetSchemaBlock(
		container,
		d,
		[]string{mkDisk},
		0,
		true,
	)
	if err != nil {
		return diag.FromErr(err)
	}

	diskDatastoreID := diskBlock[mkDiskDatastoreID].(string)

	features, err := containerGetFeatures(container, d)
	if err != nil {
		return diag.FromErr(err)
	}

	hookScript := d.Get(mkHookScriptFileID).(string)

	initialization := d.Get(mkInitialization).([]interface{})
	initializationDNSDomain := dvInitializationDNSDomain
	initializationDNSServer := dvInitializationDNSServer
	initializationHostname := dvInitializationHostname

	var initializationIPConfigIPv4Address []string

	var initializationIPConfigIPv4Gateway []string

	var initializationIPConfigIPv6Address []string

	var initializationIPConfigIPv6Gateway []string

	initializationUserAccountKeys := containers.CustomSSHKeys{}
	initializationUserAccountPassword := dvInitializationUserAccountPassword

	if len(initialization) > 0 && initialization[0] != nil {
		initializationBlock := initialization[0].(map[string]interface{})
		initializationDNS := initializationBlock[mkInitializationDNS].([]interface{})

		if len(initializationDNS) > 0 && initializationDNS[0] != nil {
			initializationDNSBlock := initializationDNS[0].(map[string]interface{})
			initializationDNSDomain = initializationDNSBlock[mkInitializationDNSDomain].(string)

			servers := initializationDNSBlock[mkInitializationDNSServers].([]interface{})
			deprecatedServer := initializationDNSBlock[mkInitializationDNSServer].(string)

			if len(servers) > 0 {
				nameserver := strings.Join(utils.ConvertToStringSlice(servers), " ")

				initializationDNSServer = nameserver
			} else {
				initializationDNSServer = deprecatedServer
			}
		}

		initializationHostname = initializationBlock[mkInitializationHostname].(string)
		initializationIPConfig := initializationBlock[mkInitializationIPConfig].([]interface{})

		for _, c := range initializationIPConfig {
			if c == nil {
				continue
			}

			configBlock := c.(map[string]interface{})
			ipv4 := configBlock[mkInitializationIPConfigIPv4].([]interface{})

			if len(ipv4) > 0 && ipv4[0] != nil {
				ipv4Block := ipv4[0].(map[string]interface{})

				initializationIPConfigIPv4Address = append(
					initializationIPConfigIPv4Address,
					ipv4Block[mkInitializationIPConfigIPv4Address].(string),
				)

				initializationIPConfigIPv4Gateway = append(
					initializationIPConfigIPv4Gateway,
					ipv4Block[mkInitializationIPConfigIPv4Gateway].(string),
				)
			} else {
				initializationIPConfigIPv4Address = append(initializationIPConfigIPv4Address, "")
				initializationIPConfigIPv4Gateway = append(initializationIPConfigIPv4Gateway, "")
			}

			ipv6 := configBlock[mkInitializationIPConfigIPv6].([]interface{})

			if len(ipv6) > 0 && ipv6[0] != nil {
				ipv6Block := ipv6[0].(map[string]interface{})

				initializationIPConfigIPv6Address = append(
					initializationIPConfigIPv6Address,
					ipv6Block[mkInitializationIPConfigIPv6Address].(string),
				)

				initializationIPConfigIPv6Gateway = append(
					initializationIPConfigIPv6Gateway,
					ipv6Block[mkInitializationIPConfigIPv6Gateway].(string),
				)
			} else {
				initializationIPConfigIPv6Address = append(initializationIPConfigIPv6Address, "")
				initializationIPConfigIPv6Gateway = append(initializationIPConfigIPv6Gateway, "")
			}
		}

		initializationUserAccount := initializationBlock[mkInitializationUserAccount].([]interface{})

		if len(initializationUserAccount) > 0 && initializationUserAccount[0] != nil {
			initializationUserAccountBlock := initializationUserAccount[0].(map[string]interface{})

			keys := initializationUserAccountBlock[mkInitializationUserAccountKeys].([]interface{})
			initializationUserAccountKeys = make(
				containers.CustomSSHKeys,
				len(keys),
			)

			for ki, kv := range keys {
				initializationUserAccountKeys[ki] = kv.(string)
			}

			initializationUserAccountPassword = initializationUserAccountBlock[mkInitializationUserAccountPassword].(string)
		}
	}

	memoryBlock, err := structure.GetSchemaBlock(
		container,
		d,
		[]string{mkMemory},
		0,
		true,
	)
	if err != nil {
		return diag.FromErr(err)
	}

	memoryDedicated := memoryBlock[mkMemoryDedicated].(int)
	memorySwap := memoryBlock[mkMemorySwap].(int)

	mountPoint := d.Get(mkMountPoint).([]interface{})
	mountPointArray := make(containers.CustomMountPointArray, 0, len(mountPoint))

	// because of default bool values:
	//nolint:gosimple
	for _, mp := range mountPoint {
		mountPointMap := mp.(map[string]interface{})
		mountPointObject := containers.CustomMountPoint{}

		acl := types.CustomBool(mountPointMap[mkMountPointACL].(bool))
		backup := types.CustomBool(mountPointMap[mkMountPointBackup].(bool))
		mountOptions := mountPointMap[mkMountPointMountOptions].([]interface{})
		path := mountPointMap[mkMountPointPath].(string)
		quota := types.CustomBool(mountPointMap[mkMountPointQuota].(bool))
		readOnly := types.CustomBool(mountPointMap[mkMountPointReadOnly].(bool))
		replicate := types.CustomBool(mountPointMap[mkMountPointReplicate].(bool))
		shared := types.CustomBool(mountPointMap[mkMountPointShared].(bool))
		size := mountPointMap[mkMountPointSize].(string)
		volume := mountPointMap[mkMountPointVolume].(string)

		// we have to set only the values that are different from the provider's defaults,
		if acl != dvMountPointACL {
			mountPointObject.ACL = &acl
		}

		if backup != dvMountPointBackup {
			mountPointObject.Backup = &backup
		}

		if path != dvMountPointPath {
			mountPointObject.MountPoint = path
		}

		if quota != dvMountPointQuota {
			mountPointObject.Quota = &quota
		}

		if readOnly != dvMountPointReadOnly {
			mountPointObject.ReadOnly = &readOnly
		}

		if replicate != dvMountPointReplicate {
			mountPointObject.Replicate = &replicate
		}

		if shared != dvMountPointShared {
			mountPointObject.Shared = &shared
		}

		if len(size) > 0 {
			var ds types.DiskSize

			ds, err = types.ParseDiskSize(size)
			if err != nil {
				return diag.Errorf("invalid disk size: %s", err.Error())
			}

			mountPointObject.Volume = fmt.Sprintf("%s:%d", volume, ds.InGigabytes())
		} else {
			mountPointObject.Volume = volume
		}

		if len(mountOptions) > 0 {
			mountOptionsArray := make([]string, 0, len(mountPoint))

			for _, option := range mountOptions {
				mountOptionsArray = append(mountOptionsArray, option.(string))
			}

			mountPointObject.MountOptions = &mountOptionsArray
		}

		mountPointArray = append(mountPointArray, mountPointObject)
	}

	var rootFS *containers.CustomRootFS

	diskSize := diskBlock[mkDiskSize].(int)
	if diskDatastoreID != "" && (diskSize != dvDiskSize || len(mountPointArray) > 0) {
		// This is a special case where the rootfs size is set to a non-default value at creation time.
		// see https://pve.proxmox.com/pve-docs/chapter-pct.html#_storage_backed_mount_points
		rootFS = &containers.CustomRootFS{
			Volume: fmt.Sprintf("%s:%d", diskDatastoreID, diskSize),
		}
	}

	networkInterface := d.Get(mkNetworkInterface).([]interface{})
	networkInterfaceArray := make(containers.CustomNetworkInterfaceArray, len(networkInterface))

	for ni, nv := range networkInterface {
		networkInterfaceMap := nv.(map[string]interface{})
		networkInterfaceObject := containers.CustomNetworkInterface{}

		bridge := networkInterfaceMap[mkNetworkInterfaceBridge].(string)
		enabled := networkInterfaceMap[mkNetworkInterfaceEnabled].(bool)
		macAddress := networkInterfaceMap[mkNetworkInterfaceMACAddress].(string)
		name := networkInterfaceMap[mkNetworkInterfaceName].(string)
		rateLimit := networkInterfaceMap[mkNetworkInterfaceRateLimit].(float64)
		vlanID := networkInterfaceMap[mkNetworkInterfaceVLANID].(int)
		mtu := networkInterfaceMap[mkNetworkInterfaceMTU].(int)
		firewall := networkInterfaceMap[mkNetworkInterfaceFirewall].(bool)

		if bridge != "" {
			networkInterfaceObject.Bridge = &bridge
		}

		networkInterfaceObject.Enabled = enabled
		networkInterfaceObject.Name = name

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

		if firewall {
			networkInterfaceObject.Firewall = types.CustomBool(firewall).Pointer()
		}

		if macAddress != "" {
			networkInterfaceObject.MACAddress = &macAddress
		}

		if rateLimit != 0 {
			networkInterfaceObject.RateLimit = &rateLimit
		}

		if vlanID != 0 {
			networkInterfaceObject.Tag = &vlanID
		}

		if mtu != 0 {
			networkInterfaceObject.MTU = &mtu
		}

		networkInterfaceArray[ni] = networkInterfaceObject
	}

	devicePassthrough := d.Get(mkDevicePassthrough).([]interface{})

	devicePassthroughArray := make(
		containers.CustomDevicePassthroughArray,
		len(devicePassthrough),
	)

	for di, dv := range devicePassthrough {
		devicePassthroughMap := dv.(map[string]interface{})
		devicePassthroughObject := containers.CustomDevicePassthrough{}

		denyWrite := types.CustomBool(
			devicePassthroughMap[mkDevicePassthroughDenyWrite].(bool),
		)
		gid := devicePassthroughMap[mkDevicePassthroughGID].(int)
		mode := devicePassthroughMap[mkDevicePassthroughMode].(string)
		path := devicePassthroughMap[mkDevicePassthroughPath].(string)
		uid := devicePassthroughMap[mkDevicePassthroughUID].(int)

		devicePassthroughObject.DenyWrite = &denyWrite
		devicePassthroughObject.GID = &gid
		devicePassthroughObject.Mode = &mode
		devicePassthroughObject.Path = path
		devicePassthroughObject.UID = &uid

		devicePassthroughArray[di] = devicePassthroughObject
	}

	operatingSystem := d.Get(mkOperatingSystem).([]interface{})

	if len(operatingSystem) == 0 || operatingSystem[0] == nil {
		return diag.Errorf(
			"\"%s\": required field is not set",
			mkOperatingSystem,
		)
	}

	operatingSystemBlock := operatingSystem[0].(map[string]interface{})
	operatingSystemTemplateFileID := operatingSystemBlock[mkOperatingSystemTemplateFileID].(string)
	operatingSystemType := operatingSystemBlock[mkOperatingSystemType].(string)

	poolID := d.Get(mkPoolID).(string)
	protection := types.CustomBool(d.Get(mkProtection).(bool))
	started := types.CustomBool(d.Get(mkStarted).(bool))
	startOnBoot := types.CustomBool(d.Get(mkStartOnBoot).(bool))
	startupBehavior := containerGetStartupBehavior(d)
	tags := d.Get(mkTags).([]interface{})
	template := types.CustomBool(d.Get(mkTemplate).(bool))
	unprivileged := types.CustomBool(d.Get(mkUnprivileged).(bool))
	vmIDUntyped, hasVMID := d.GetOk(mkVMID)
	vmID := vmIDUntyped.(int)

	if !hasVMID {
		vmIDNew, err := config.GetIDGenerator().NextID(ctx)
		if err != nil {
			return diag.FromErr(err)
		}

		vmID = vmIDNew

		err = d.Set(mkVMID, vmID)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// Attempt to create the container using the retrieved values.
	createBody := containers.CreateRequestBody{
		ConsoleEnabled:       &consoleEnabled,
		ConsoleMode:          &consoleMode,
		CPUArchitecture:      &cpuArchitecture,
		CPUCores:             &cpuCores,
		CPUUnits:             &cpuUnits,
		DatastoreID:          &diskDatastoreID,
		DedicatedMemory:      &memoryDedicated,
		DevicePassthrough:    devicePassthroughArray,
		Features:             features,
		MountPoints:          mountPointArray,
		NetworkInterfaces:    networkInterfaceArray,
		OSTemplateFileVolume: &operatingSystemTemplateFileID,
		OSType:               &operatingSystemType,
		Protection:           &protection,
		RootFS:               rootFS,
		Start:                &started,
		StartOnBoot:          &startOnBoot,
		StartupBehavior:      startupBehavior,
		Swap:                 &memorySwap,
		Template:             &template,
		TTY:                  &consoleTTYCount,
		Unprivileged:         &unprivileged,
		VMID:                 &vmID,
	}

	if description != "" {
		createBody.Description = &description
	}

	if hookScript != "" {
		createBody.HookScript = &hookScript
	}

	if initializationDNSDomain != "" {
		createBody.DNSDomain = &initializationDNSDomain
	}

	if initializationDNSServer != "" {
		createBody.DNSServer = &initializationDNSServer
	}

	if initializationHostname != "" {
		createBody.Hostname = &initializationHostname
	}

	if len(initializationUserAccountKeys) > 0 {
		createBody.SSHKeys = &initializationUserAccountKeys
	}

	if initializationUserAccountPassword != "" {
		createBody.Password = &initializationUserAccountPassword
	}

	if poolID != "" {
		createBody.PoolID = &poolID
	}

	if len(tags) > 0 {
		tagsString := containerGetTagsString(d)
		createBody.Tags = &tagsString
	}

	err = client.Node(nodeName).Container(0).CreateContainer(ctx, &createBody)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(vmID))

	// Wait for the container's lock to be released.
	err = client.Node(nodeName).Container(vmID).WaitForContainerConfigUnlock(ctx, true)
	if err != nil {
		return diag.FromErr(err)
	}

	return containerCreateStart(ctx, d, m)
}

func containerCreateStart(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	started := d.Get(mkStarted).(bool)
	template := d.Get(mkTemplate).(bool)

	if !started || template {
		return containerRead(ctx, d, m)
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

	containerAPI := client.Node(nodeName).Container(vmID)

	// Start the container and wait for it to reach a running state before continuing.
	err = containerAPI.StartContainer(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	return containerRead(ctx, d, m)
}

func containerGetExistingNetworkInterface(
	ctx context.Context,
	containerAPI *containers.Client,
) ([]interface{}, error) {
	containerInfo, err := containerAPI.GetContainer(ctx)
	if err != nil {
		return []interface{}{}, fmt.Errorf("error getting container information: %w", err)
	}

	//nolint:prealloc
	var networkInterfaces []interface{}

	networkInterfaceArray := []*containers.CustomNetworkInterface{
		containerInfo.NetworkInterface0,
		containerInfo.NetworkInterface1,
		containerInfo.NetworkInterface2,
		containerInfo.NetworkInterface3,
		containerInfo.NetworkInterface4,
		containerInfo.NetworkInterface5,
		containerInfo.NetworkInterface6,
		containerInfo.NetworkInterface7,
	}

	for _, nv := range networkInterfaceArray {
		if nv == nil {
			continue
		}

		networkInterface := map[string]interface{}{}

		networkInterface[mkNetworkInterfaceEnabled] = true
		networkInterface[mkNetworkInterfaceName] = nv.Name

		if nv.Bridge != nil {
			networkInterface[mkNetworkInterfaceBridge] = *nv.Bridge
		} else {
			networkInterface[mkNetworkInterfaceBridge] = ""
		}

		if nv.Firewall != nil && *nv.Firewall {
			networkInterface[mkNetworkInterfaceFirewall] = true
		} else {
			networkInterface[mkNetworkInterfaceFirewall] = false
		}

		if nv.MACAddress != nil {
			networkInterface[mkNetworkInterfaceMACAddress] = *nv.MACAddress
		} else {
			networkInterface[mkNetworkInterfaceMACAddress] = ""
		}

		if nv.RateLimit != nil {
			networkInterface[mkNetworkInterfaceRateLimit] = *nv.RateLimit
		} else {
			networkInterface[mkNetworkInterfaceRateLimit] = float64(0)
		}

		if nv.Tag != nil {
			networkInterface[mkNetworkInterfaceVLANID] = *nv.Tag
		} else {
			networkInterface[mkNetworkInterfaceVLANID] = 0
		}

		if nv.MTU != nil {
			networkInterface[mkNetworkInterfaceMTU] = *nv.MTU
		} else {
			networkInterface[mkNetworkInterfaceMTU] = 0
		}

		networkInterfaces = append(networkInterfaces, networkInterface)
	}

	return networkInterfaces, nil
}

func containerGetTagsString(d *schema.ResourceData) string {
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

func containerGetStartupBehavior(d *schema.ResourceData) *containers.CustomStartupBehavior {
	startup := d.Get(mkStartup).([]interface{})
	if len(startup) > 0 && startup[0] != nil {
		startupBlock := startup[0].(map[string]interface{})
		startupOrder := startupBlock[mkStartupOrder].(int)
		startupUpDelay := startupBlock[mkStartupUpDelay].(int)
		startupDownDelay := startupBlock[mkStartupDownDelay].(int)

		order := containers.CustomStartupBehavior{}

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

func containerGetFeatures(resource *schema.Resource, d *schema.ResourceData) (*containers.CustomFeatures, error) {
	featuresBlock, err := structure.GetSchemaBlock(
		resource,
		d,
		[]string{mkFeatures},
		0,
		true,
	)
	if err != nil {
		return nil, fmt.Errorf("error getting container features from schema: %w", err)
	}

	nesting := types.CustomBool(featuresBlock[mkFeaturesNesting].(bool))
	keyctl := types.CustomBool(featuresBlock[mkFeaturesKeyControl].(bool))
	fuse := types.CustomBool(featuresBlock[mkFeaturesFUSE].(bool))
	mountTypes := featuresBlock[mkFeaturesMountTypes].([]interface{})

	var mountTypesConverted []string
	if mountTypes != nil {
		mountTypesConverted = make([]string, len(mountTypes))
		for i, mountType := range mountTypes {
			mountTypesConverted[i] = mountType.(string)
		}
	} else {
		mountTypesConverted = []string{}
	}

	features := containers.CustomFeatures{
		MountTypes: &mountTypesConverted,
	}

	if nesting {
		features.Nesting = &nesting
	}

	if keyctl {
		features.KeyControl = &keyctl
	}

	if fuse {
		features.FUSE = &fuse
	}

	return &features, nil
}

func containerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	config := m.(proxmoxtf.ProviderConfiguration)

	client, e := config.GetClient()
	if e != nil {
		return diag.FromErr(e)
	}

	nodeName := d.Get(mkNodeName).(string)

	vmID, e := strconv.Atoi(d.Id())
	if e != nil {
		return diag.FromErr(e)
	}

	containerAPI := client.Node(nodeName).Container(vmID)

	// Retrieve the entire configuration in order to compare it to the state.
	containerConfig, e := containerAPI.GetContainer(ctx)
	if e != nil {
		if errors.Is(e, api.ErrResourceDoesNotExist) {
			d.SetId("")

			return nil
		}

		return diag.FromErr(e)
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

	clone := d.Get(mkClone).([]interface{})

	// Compare the primitive values to those stored in the state.
	currentDescription := d.Get(mkDescription).(string)

	if len(clone) == 0 || currentDescription != dvDescription {
		if containerConfig.Description != nil {
			e = d.Set(mkDescription, *containerConfig.Description)
		} else {
			e = d.Set(mkDescription, "")
		}

		diags = append(diags, diag.FromErr(e)...)
	}

	// Compare the console configuration to the one stored in the state.
	console := map[string]interface{}{}

	if containerConfig.ConsoleEnabled != nil {
		console[mkConsoleEnabled] = *containerConfig.ConsoleEnabled
	} else {
		// Default value of "console" is "1" according to the API documentation.
		console[mkConsoleEnabled] = true
	}

	if containerConfig.ConsoleMode != nil {
		console[mkConsoleMode] = *containerConfig.ConsoleMode
	} else {
		// Default value of "cmode" is "tty" according to the API documentation.
		console[mkConsoleMode] = "tty"
	}

	if containerConfig.TTY != nil {
		console[mkConsoleTTYCount] = *containerConfig.TTY
	} else {
		// Default value of "tty" is "2" according to the API documentation.
		console[mkConsoleTTYCount] = 2
	}

	currentConsole := d.Get(mkConsole).([]interface{})

	if len(clone) > 0 {
		if len(currentConsole) > 0 {
			err := d.Set(mkConsole, []interface{}{console})
			diags = append(diags, diag.FromErr(err)...)
		}
	} else if len(currentConsole) > 0 ||
		console[mkConsoleEnabled] != types.CustomBool(dvConsoleEnabled) ||
		console[mkConsoleMode] != dvConsoleMode ||
		console[mkConsoleTTYCount] != dvConsoleTTYCount {
		err := d.Set(mkConsole, []interface{}{console})
		diags = append(diags, diag.FromErr(err)...)
	}

	// Compare the CPU configuration to the one stored in the state.
	cpu := map[string]interface{}{}

	if containerConfig.CPUArchitecture != nil {
		cpu[mkCPUArchitecture] = *containerConfig.CPUArchitecture
	} else {
		// Default value of "arch" is "amd64" according to the API documentation.
		cpu[mkCPUArchitecture] = "amd64"
	}

	if containerConfig.CPUCores != nil {
		cpu[mkCPUCores] = *containerConfig.CPUCores
	} else {
		// Default value of "cores" is "1" according to the API documentation.
		cpu[mkCPUCores] = 1
	}

	if containerConfig.CPUUnits != nil {
		cpu[mkCPUUnits] = *containerConfig.CPUUnits
	} else {
		// Default value of "cpuunits" is "1024" according to the API documentation.
		cpu[mkCPUUnits] = 1024
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
		cpu[mkCPUUnits] != dvCPUUnits {
		err := d.Set(mkCPU, []interface{}{cpu})
		diags = append(diags, diag.FromErr(err)...)
	}

	// Compare the disk configuration to the one stored in the state.
	disk := map[string]interface{}{}

	if containerConfig.RootFS != nil {
		volumeParts := strings.Split(containerConfig.RootFS.Volume, ":")
		disk[mkDiskDatastoreID] = volumeParts[0]
		disk[mkDiskSize] = containerConfig.RootFS.Size.InGigabytes()
	} else {
		// Default value of "storage" is "local" according to the API documentation.
		disk[mkDiskDatastoreID] = "local"
		disk[mkDiskSize] = dvDiskSize
	}

	currentDisk := d.Get(mkDisk).([]interface{})

	if len(clone) > 0 {
		if len(currentDisk) > 0 && currentDisk[0] != nil {
			// do not override the rootfs size if it was not changed during the clone operation
			if currentDisk[0].(map[string]interface{})[mkDiskSize] == dvDiskSize {
				disk[mkDiskSize] = dvDiskSize
			}

			err := d.Set(mkDisk, []interface{}{disk})
			diags = append(diags, diag.FromErr(err)...)
		}
	} else if len(currentDisk) > 0 ||
		disk[mkDiskDatastoreID] != dvDiskDatastoreID ||
		disk[mkDiskSize] != dvDiskSize {
		err := d.Set(mkDisk, []interface{}{disk})
		diags = append(diags, diag.FromErr(err)...)
	}

	// Compare the memory configuration to the one stored in the state.
	memory := map[string]interface{}{}

	if containerConfig.DedicatedMemory != nil {
		memory[mkMemoryDedicated] = *containerConfig.DedicatedMemory
	} else {
		memory[mkMemoryDedicated] = 0
	}

	if containerConfig.Swap != nil {
		memory[mkMemorySwap] = *containerConfig.Swap
	} else {
		memory[mkMemorySwap] = 0
	}

	currentMemory := d.Get(mkMemory).([]interface{})

	if len(clone) > 0 {
		if len(currentMemory) > 0 {
			err := d.Set(mkMemory, []interface{}{memory})
			diags = append(diags, diag.FromErr(err)...)
		}
	} else if len(currentMemory) > 0 ||
		memory[mkMemoryDedicated] != dvMemoryDedicated ||
		memory[mkMemorySwap] != dvMemorySwap {
		err := d.Set(mkMemory, []interface{}{memory})
		diags = append(diags, diag.FromErr(err)...)
	}

	// Compare the initialization and network interface configuration to the one stored in the state.
	initialization := map[string]interface{}{}

	if containerConfig.DNSDomain != nil || containerConfig.DNSServer != nil {
		initializationDNS := map[string]interface{}{}

		if containerConfig.DNSDomain != nil {
			initializationDNS[mkInitializationDNSDomain] = *containerConfig.DNSDomain
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
		if containerConfig.DNSServer != nil {
			if ok && currentInitializationDNSServer != "" {
				initializationDNS[mkInitializationDNSServer] = *containerConfig.DNSServer
			} else {
				dnsServer := strings.Split(*containerConfig.DNSServer, " ")
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

	if containerConfig.Hostname != nil {
		initialization[mkInitializationHostname] = *containerConfig.Hostname
	} else {
		initialization[mkInitializationHostname] = ""
	}

	devicePassthroughArray := []*containers.CustomDevicePassthrough{
		containerConfig.DevicePassthrough0,
		containerConfig.DevicePassthrough1,
		containerConfig.DevicePassthrough2,
		containerConfig.DevicePassthrough3,
		containerConfig.DevicePassthrough4,
		containerConfig.DevicePassthrough5,
		containerConfig.DevicePassthrough6,
		containerConfig.DevicePassthrough7,
	}

	devicePassthroughList := make([]interface{}, 0, len(devicePassthroughArray))

	for _, dp := range devicePassthroughArray {
		if dp == nil {
			continue
		}

		devicePassthrough := map[string]interface{}{}

		if dp.DenyWrite != nil {
			devicePassthrough[mkDevicePassthroughDenyWrite] = *dp.DenyWrite
		} else {
			devicePassthrough[mkDevicePassthroughDenyWrite] = false
		}

		if dp.GID != nil {
			devicePassthrough[mkDevicePassthroughGID] = *dp.GID
		} else {
			devicePassthrough[mkDevicePassthroughGID] = 0
		}

		if dp.Mode != nil {
			devicePassthrough[mkDevicePassthroughMode] = *dp.Mode
		} else {
			devicePassthrough[mkDevicePassthroughMode] = ""
		}

		devicePassthrough[mkDevicePassthroughPath] = dp.Path

		if dp.UID != nil {
			devicePassthrough[mkDevicePassthroughUID] = *dp.UID
		} else {
			devicePassthrough[mkDevicePassthroughUID] = 0
		}

		devicePassthroughList = append(devicePassthroughList, devicePassthrough)
	}

	currentDevicePassthrough := d.Get(mkDevicePassthrough).([]interface{})

	if len(clone) > 0 {
		if len(currentDevicePassthrough) > 0 {
			err := d.Set(mkDevicePassthrough, devicePassthroughList)
			diags = append(diags, diag.FromErr(err)...)
		}
	} else if len(devicePassthroughList) > 0 {
		err := d.Set(mkDevicePassthrough, devicePassthroughList)
		diags = append(diags, diag.FromErr(err)...)
	}

	mountPointArray := []*containers.CustomMountPoint{
		containerConfig.MountPoint0,
		containerConfig.MountPoint1,
		containerConfig.MountPoint2,
		containerConfig.MountPoint3,
		containerConfig.MountPoint4,
		containerConfig.MountPoint5,
		containerConfig.MountPoint6,
		containerConfig.MountPoint7,
	}

	mountPointList := make([]interface{}, 0, len(mountPointArray))

	for _, mp := range mountPointArray {
		if mp == nil {
			continue
		}

		mountPoint := map[string]interface{}{}

		if mp.ACL != nil {
			mountPoint[mkMountPointACL] = *mp.ACL
		} else {
			mountPoint[mkMountPointACL] = false
		}

		if mp.Backup != nil {
			mountPoint[mkMountPointBackup] = *mp.Backup
		} else {
			mountPoint[mkMountPointBackup] = dvMountPointBackup
		}

		if mp.MountOptions != nil {
			mountPoint[mkMountPointMountOptions] = *mp.MountOptions
		} else {
			mountPoint[mkMountPointMountOptions] = []string{}
		}

		mountPoint[mkMountPointPath] = mp.MountPoint

		if mp.Quota != nil {
			mountPoint[mkMountPointQuota] = *mp.Quota
		} else {
			mountPoint[mkMountPointQuota] = false
		}

		if mp.ReadOnly != nil {
			mountPoint[mkMountPointReadOnly] = *mp.ReadOnly
		} else {
			mountPoint[mkMountPointReadOnly] = false
		}

		if mp.Replicate != nil {
			mountPoint[mkMountPointReplicate] = *mp.Replicate
		} else {
			mountPoint[mkMountPointReplicate] = true
		}

		if mp.Shared != nil {
			mountPoint[mkMountPointShared] = *mp.Shared
		} else {
			mountPoint[mkMountPointShared] = false
		}

		if mp.DiskSize != nil {
			mountPoint[mkMountPointSize] = *mp.DiskSize
		} else {
			mountPoint[mkMountPointSize] = ""
		}

		mountPoint[mkMountPointVolume] = mp.Volume

		mountPointList = append(mountPointList, mountPoint)
	}

	currentMountPoint := d.Get(mkMountPoint).([]interface{})

	if len(clone) > 0 {
		if len(currentMountPoint) > 0 {
			err := d.Set(mkMountPoint, mountPointList)
			diags = append(diags, diag.FromErr(err)...)
		}
	} else if len(mountPointList) > 0 {
		err := d.Set(mkMountPoint, mountPointList)
		diags = append(diags, diag.FromErr(err)...)
	}

	var ipConfigList []interface{}

	networkInterfaceArray := []*containers.CustomNetworkInterface{
		containerConfig.NetworkInterface0,
		containerConfig.NetworkInterface1,
		containerConfig.NetworkInterface2,
		containerConfig.NetworkInterface3,
		containerConfig.NetworkInterface4,
		containerConfig.NetworkInterface5,
		containerConfig.NetworkInterface6,
		containerConfig.NetworkInterface7,
	}

	//nolint:prealloc
	var networkInterfaceList []interface{}

	for _, nv := range networkInterfaceArray {
		if nv == nil {
			continue
		}

		if nv.IPv4Address != nil || nv.IPv4Gateway != nil || nv.IPv6Address != nil ||
			nv.IPv6Gateway != nil {
			ipConfig := map[string]interface{}{}

			if nv.IPv4Address != nil || nv.IPv4Gateway != nil {
				ip := map[string]interface{}{}

				if nv.IPv4Address != nil {
					ip[mkInitializationIPConfigIPv4Address] = *nv.IPv4Address
				} else {
					ip[mkInitializationIPConfigIPv4Address] = ""
				}

				if nv.IPv4Gateway != nil {
					ip[mkInitializationIPConfigIPv4Gateway] = *nv.IPv4Gateway
				} else {
					ip[mkInitializationIPConfigIPv4Gateway] = ""
				}

				ipConfig[mkInitializationIPConfigIPv4] = []interface{}{
					ip,
				}
			} else {
				ipConfig[mkInitializationIPConfigIPv4] = []interface{}{}
			}

			if nv.IPv6Address != nil || nv.IPv6Gateway != nil {
				ip := map[string]interface{}{}

				if nv.IPv6Address != nil {
					ip[mkInitializationIPConfigIPv6Address] = *nv.IPv6Address
				} else {
					ip[mkInitializationIPConfigIPv6Address] = ""
				}

				if nv.IPv6Gateway != nil {
					ip[mkInitializationIPConfigIPv6Gateway] = *nv.IPv6Gateway
				} else {
					ip[mkInitializationIPConfigIPv6Gateway] = ""
				}

				ipConfig[mkInitializationIPConfigIPv6] = []interface{}{
					ip,
				}
			} else {
				ipConfig[mkInitializationIPConfigIPv6] = []interface{}{}
			}

			ipConfigList = append(ipConfigList, ipConfig)
		}

		networkInterface := map[string]interface{}{}

		networkInterface[mkNetworkInterfaceEnabled] = true
		networkInterface[mkNetworkInterfaceName] = nv.Name

		if nv.Bridge != nil {
			networkInterface[mkNetworkInterfaceBridge] = *nv.Bridge
		} else {
			networkInterface[mkNetworkInterfaceBridge] = ""
		}

		if nv.Firewall != nil && *nv.Firewall {
			networkInterface[mkNetworkInterfaceFirewall] = true
		} else {
			networkInterface[mkNetworkInterfaceFirewall] = false
		}

		if nv.MACAddress != nil {
			networkInterface[mkNetworkInterfaceMACAddress] = *nv.MACAddress
		} else {
			networkInterface[mkNetworkInterfaceMACAddress] = ""
		}

		if nv.RateLimit != nil {
			networkInterface[mkNetworkInterfaceRateLimit] = *nv.RateLimit
		} else {
			networkInterface[mkNetworkInterfaceRateLimit] = 0
		}

		if nv.Tag != nil {
			networkInterface[mkNetworkInterfaceVLANID] = *nv.Tag
		} else {
			networkInterface[mkNetworkInterfaceVLANID] = 0
		}

		if nv.MTU != nil {
			networkInterface[mkNetworkInterfaceMTU] = *nv.MTU
		} else {
			networkInterface[mkNetworkInterfaceMTU] = 0
		}

		networkInterfaceList = append(networkInterfaceList, networkInterface)
	}

	initialization[mkInitializationIPConfig] = ipConfigList

	currentInitialization := d.Get(mkInitialization).([]interface{})

	if len(currentInitialization) > 0 && currentInitialization[0] != nil {
		currentInitializationMap := currentInitialization[0].(map[string]interface{})

		initialization[mkInitializationUserAccount] = currentInitializationMap[mkInitializationUserAccount].([]interface{})
	}

	if len(clone) > 0 {
		if len(currentInitialization) > 0 && currentInitialization[0] != nil {
			currentInitializationBlock := currentInitialization[0].(map[string]interface{})
			currentInitializationDNS := currentInitializationBlock[mkInitializationDNS].([]interface{})

			if len(currentInitializationDNS) == 0 {
				initialization[mkInitializationDNS] = []interface{}{}
			}

			currentInitializationIPConfig := currentInitializationBlock[mkInitializationIPConfig].([]interface{})

			if len(currentInitializationIPConfig) == 0 {
				initialization[mkInitializationIPConfig] = []interface{}{}
			}

			currentInitializationUserAccount := currentInitializationBlock[mkInitializationUserAccount].([]interface{})

			if len(currentInitializationUserAccount) == 0 {
				initialization[mkInitializationUserAccount] = []interface{}{}
			}

			if len(initialization) > 0 {
				e = d.Set(
					mkInitialization,
					[]interface{}{initialization},
				)
			} else {
				e = d.Set(mkInitialization, []interface{}{})
			}

			diags = append(diags, diag.FromErr(e)...)
		}

		currentNetworkInterface := d.Get(mkNetworkInterface).([]interface{})

		if len(currentNetworkInterface) > 0 {
			err := d.Set(
				mkNetworkInterface,
				networkInterfaceList,
			)
			diags = append(diags, diag.FromErr(err)...)
		}
	} else {
		if len(initialization) > 0 {
			e = d.Set(mkInitialization, []interface{}{initialization})
		} else {
			e = d.Set(mkInitialization, []interface{}{})
		}

		diags = append(diags, diag.FromErr(e)...)

		err := d.Set(mkNetworkInterface, networkInterfaceList)
		diags = append(diags, diag.FromErr(err)...)
	}

	// Compare the startup behavior to the one stored in the state.
	var startup map[string]interface{}

	if containerConfig.StartupBehavior != nil {
		startup = map[string]interface{}{}

		if containerConfig.StartupBehavior.Order != nil {
			startup[mkStartupOrder] = *containerConfig.StartupBehavior.Order
		} else {
			startup[mkStartupOrder] = dvStartupOrder
		}

		if containerConfig.StartupBehavior.Up != nil {
			startup[mkStartupUpDelay] = *containerConfig.StartupBehavior.Up
		} else {
			startup[mkStartupUpDelay] = dvStartupUpDelay
		}

		if containerConfig.StartupBehavior.Down != nil {
			startup[mkStartupDownDelay] = *containerConfig.StartupBehavior.Down
		} else {
			startup[mkStartupDownDelay] = dvStartupDownDelay
		}
	}

	currentStartup := d.Get(mkStartup).([]interface{})

	switch {
	case len(clone) > 0 && len(currentStartup) > 0:
		err := d.Set(mkStartup, []interface{}{startup})
		diags = append(diags, diag.FromErr(err)...)
	case len(startup) == 0:
		err := d.Set(mkStartup, []interface{}{})
		diags = append(diags, diag.FromErr(err)...)
	case len(currentStartup) > 0 ||
		startup[mkStartupOrder] != mkStartupOrder ||
		startup[mkStartupUpDelay] != dvStartupUpDelay ||
		startup[mkStartupDownDelay] != dvStartupDownDelay:
		err := d.Set(mkStartup, []interface{}{startup})
		diags = append(diags, diag.FromErr(err)...)
	}

	// Compare the operating system configuration to the one stored in the state.
	operatingSystem := map[string]interface{}{}

	if containerConfig.OSType != nil {
		operatingSystem[mkOperatingSystemType] = *containerConfig.OSType
	} else {
		// Default value of "ostype" is "" according to the API documentation.
		operatingSystem[mkOperatingSystemType] = ""
	}

	currentOperatingSystem := d.Get(mkOperatingSystem).([]interface{})

	if len(currentOperatingSystem) > 0 && currentOperatingSystem[0] != nil {
		currentOperatingSystemMap := currentOperatingSystem[0].(map[string]interface{})

		operatingSystem[mkOperatingSystemTemplateFileID] = currentOperatingSystemMap[mkOperatingSystemTemplateFileID]
	}

	if len(clone) > 0 {
		if len(currentOperatingSystem) > 0 {
			err := d.Set(
				mkOperatingSystem,
				[]interface{}{operatingSystem},
			)
			diags = append(diags, diag.FromErr(err)...)
		}
	} else if len(currentOperatingSystem) > 0 ||
		operatingSystem[mkOperatingSystemType] != dvOperatingSystemType {
		err := d.Set(mkOperatingSystem, []interface{}{operatingSystem})
		diags = append(diags, diag.FromErr(err)...)
	}

	currentProtection := types.CustomBool(d.Get(mkProtection).(bool))

	//nolint:gosimple
	if len(clone) == 0 || currentProtection != dvProtection {
		if containerConfig.Protection != nil {
			e = d.Set(
				mkProtection,
				bool(*containerConfig.Protection),
			)
		} else {
			e = d.Set(mkProtection, false)
		}

		diags = append(diags, diag.FromErr(e)...)
	}

	currentTags := d.Get(mkTags).([]interface{})

	if len(clone) == 0 || len(currentTags) > 0 {
		var tags []string

		if containerConfig.Tags != nil {
			for _, tag := range strings.Split(*containerConfig.Tags, ";") {
				t := strings.TrimSpace(tag)
				if len(t) > 0 {
					tags = append(tags, t)
				}
			}

			sort.Strings(tags)
		}

		e = d.Set(mkTags, tags)
		diags = append(diags, diag.FromErr(e)...)
	}

	currentTemplate := d.Get(mkTemplate).(bool)

	//nolint:gosimple
	if len(clone) == 0 || currentTemplate != dvTemplate {
		if containerConfig.Template != nil {
			e = d.Set(
				mkTemplate,
				bool(*containerConfig.Template),
			)
		} else {
			e = d.Set(mkTemplate, false)
		}

		diags = append(diags, diag.FromErr(e)...)
	}

	// Determine the state of the container in order to update the "started" argument.
	status, e := containerAPI.GetContainerStatus(ctx)
	if e != nil {
		return diag.FromErr(e)
	}

	e = d.Set(mkStarted, status.Status == "running")
	diags = append(diags, diag.FromErr(e)...)

	return diags
}

func containerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	updateTimeoutSec := d.Get(mkTimeoutUpdate).(int)

	ctx, cancel := context.WithTimeout(ctx, time.Duration(updateTimeoutSec)*time.Second)
	defer cancel()

	config := m.(proxmoxtf.ProviderConfiguration)

	client, e := config.GetClient()
	if e != nil {
		return diag.FromErr(e)
	}

	nodeName := d.Get(mkNodeName).(string)

	vmID, e := strconv.Atoi(d.Id())
	if e != nil {
		return diag.FromErr(e)
	}

	containerAPI := client.Node(nodeName).Container(vmID)

	// Prepare the new request object.
	updateBody := containers.UpdateRequestBody{
		Delete: []string{},
	}

	rebootRequired := false
	container := Container()

	// Retrieve the clone argument as the update logic varies for clones.
	clone := d.Get(mkClone).([]interface{})

	// Prepare the new primitive values.
	description := d.Get(mkDescription).(string)
	updateBody.Description = &description

	template := types.CustomBool(d.Get(mkTemplate).(bool))

	if d.HasChange(mkTemplate) {
		updateBody.Template = &template
	}

	// Prepare the new console configuration.
	if d.HasChange(mkConsole) {
		consoleBlock, err := structure.GetSchemaBlock(
			container,
			d,
			[]string{mkConsole},
			0,
			true,
		)
		if err != nil {
			return diag.FromErr(err)
		}

		consoleEnabled := types.CustomBool(
			consoleBlock[mkConsoleEnabled].(bool),
		)
		consoleMode := consoleBlock[mkConsoleMode].(string)
		consoleTTYCount := consoleBlock[mkConsoleTTYCount].(int)

		updateBody.ConsoleEnabled = &consoleEnabled
		updateBody.ConsoleMode = &consoleMode
		updateBody.TTY = &consoleTTYCount

		rebootRequired = true
	}

	// Prepare the new CPU configuration.
	if d.HasChange(mkCPU) {
		cpuBlock, err := structure.GetSchemaBlock(
			container,
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
		cpuUnits := cpuBlock[mkCPUUnits].(int)

		updateBody.CPUArchitecture = &cpuArchitecture
		updateBody.CPUCores = &cpuCores
		updateBody.CPUUnits = &cpuUnits

		rebootRequired = true
	}

	if d.HasChange(mkFeatures) {
		features, err := containerGetFeatures(container, d)
		if err != nil {
			return diag.FromErr(err)
		}

		updateBody.Features = features
	}

	if d.HasChange(mkHookScriptFileID) {
		hookScript := d.Get(mkHookScriptFileID).(string)
		if hookScript != "" {
			updateBody.HookScript = &hookScript
		} else {
			updateBody.Delete = append(updateBody.Delete, "hookscript")
		}
	}

	// Prepare the new initialization configuration.
	initialization := d.Get(mkInitialization).([]interface{})
	initializationDNSDomain := dvInitializationDNSDomain
	initializationDNSServer := dvInitializationDNSServer
	initializationHostname := dvInitializationHostname

	var initializationIPConfigIPv4Address []string

	var initializationIPConfigIPv4Gateway []string

	var initializationIPConfigIPv6Address []string

	var initializationIPConfigIPv6Gateway []string

	if len(initialization) > 0 && initialization[0] != nil {
		initializationBlock := initialization[0].(map[string]interface{})
		initializationDNS := initializationBlock[mkInitializationDNS].([]interface{})

		if len(initializationDNS) > 0 {
			initializationDNSBlock := initializationDNS[0].(map[string]interface{})
			initializationDNSDomain = initializationDNSBlock[mkInitializationDNSDomain].(string)

			servers := initializationDNSBlock[mkInitializationDNSServers].([]interface{})
			deprecatedServer := initializationDNSBlock[mkInitializationDNSServer].(string)

			if len(servers) > 0 {
				initializationDNSServer = strings.Join(utils.ConvertToStringSlice(servers), " ")
			} else {
				initializationDNSServer = deprecatedServer
			}
		}

		initializationHostname = initializationBlock[mkInitializationHostname].(string)
		initializationIPConfig := initializationBlock[mkInitializationIPConfig].([]interface{})

		for _, c := range initializationIPConfig {
			if c == nil {
				continue
			}

			configBlock := c.(map[string]interface{})
			ipv4 := configBlock[mkInitializationIPConfigIPv4].([]interface{})

			if len(ipv4) > 0 && ipv4[0] != nil {
				ipv4Block := ipv4[0].(map[string]interface{})

				initializationIPConfigIPv4Address = append(
					initializationIPConfigIPv4Address,
					ipv4Block[mkInitializationIPConfigIPv4Address].(string),
				)

				initializationIPConfigIPv4Gateway = append(
					initializationIPConfigIPv4Gateway,
					ipv4Block[mkInitializationIPConfigIPv4Gateway].(string),
				)
			} else {
				initializationIPConfigIPv4Address = append(initializationIPConfigIPv4Address, "")
				initializationIPConfigIPv4Gateway = append(initializationIPConfigIPv4Gateway, "")
			}

			ipv6 := configBlock[mkInitializationIPConfigIPv6].([]interface{})

			if len(ipv6) > 0 && ipv6[0] != nil {
				ipv6Block := ipv6[0].(map[string]interface{})

				initializationIPConfigIPv6Address = append(
					initializationIPConfigIPv6Address,
					ipv6Block[mkInitializationIPConfigIPv6Address].(string),
				)

				initializationIPConfigIPv6Gateway = append(
					initializationIPConfigIPv6Gateway,
					ipv6Block[mkInitializationIPConfigIPv6Gateway].(string),
				)
			} else {
				initializationIPConfigIPv6Address = append(initializationIPConfigIPv6Address, "")
				initializationIPConfigIPv6Gateway = append(initializationIPConfigIPv6Gateway, "")
			}
		}
	}

	if d.HasChange(mkInitialization) {
		updateBody.DNSDomain = &initializationDNSDomain
		updateBody.DNSServer = &initializationDNSServer
		updateBody.Hostname = &initializationHostname

		rebootRequired = true
	}

	// Prepare the new memory configuration.
	if d.HasChange(mkMemory) {
		memoryBlock, err := structure.GetSchemaBlock(
			container,
			d,
			[]string{mkMemory},
			0,
			true,
		)
		if err != nil {
			return diag.FromErr(err)
		}

		memoryDedicated := memoryBlock[mkMemoryDedicated].(int)
		memorySwap := memoryBlock[mkMemorySwap].(int)

		updateBody.DedicatedMemory = &memoryDedicated
		updateBody.Swap = &memorySwap

		rebootRequired = true
	}

	// Prepare the new device passthrough configuration.
	if d.HasChange(mkDevicePassthrough) {
		_, newDevicePassthrough := d.GetChange(mkDevicePassthrough)

		devicePassthrough := newDevicePassthrough.([]interface{})
		devicePassthroughArray := make(
			containers.CustomDevicePassthroughArray,
			len(devicePassthrough),
		)

		for i, dp := range devicePassthrough {
			devicePassthroughMap := dp.(map[string]interface{})
			devicePassthroughObject := containers.CustomDevicePassthrough{}

			denyWrite := types.CustomBool(devicePassthroughMap[mkDevicePassthroughDenyWrite].(bool))
			gid := devicePassthroughMap[mkDevicePassthroughGID].(int)
			mode := devicePassthroughMap[mkDevicePassthroughMode].(string)
			path := devicePassthroughMap[mkDevicePassthroughPath].(string)
			uid := devicePassthroughMap[mkDevicePassthroughUID].(int)

			devicePassthroughObject.DenyWrite = &denyWrite
			devicePassthroughObject.GID = &gid
			devicePassthroughObject.Mode = &mode
			devicePassthroughObject.Path = path
			devicePassthroughObject.UID = &uid

			devicePassthroughArray[i] = devicePassthroughObject
		}

		updateBody.DevicePassthrough = devicePassthroughArray

		rebootRequired = true
	}

	// Prepare the new mount point configuration.
	if d.HasChange(mkMountPoint) {
		_, newMountPoints := d.GetChange(mkMountPoint)

		mountPoints := newMountPoints.([]interface{})
		mountPointArray := make(
			containers.CustomMountPointArray,
			len(mountPoints),
		)

		for i, mp := range mountPoints {
			mountPointMap := mp.(map[string]interface{})
			mountPointObject := containers.CustomMountPoint{}

			acl := types.CustomBool(mountPointMap[mkMountPointACL].(bool))
			backup := types.CustomBool(mountPointMap[mkMountPointBackup].(bool))
			mountOptions := mountPointMap[mkMountPointMountOptions].([]interface{})
			path := mountPointMap[mkMountPointPath].(string)
			quota := types.CustomBool(mountPointMap[mkMountPointQuota].(bool))
			readOnly := types.CustomBool(mountPointMap[mkMountPointReadOnly].(bool))
			replicate := types.CustomBool(mountPointMap[mkMountPointReplicate].(bool))
			shared := types.CustomBool(mountPointMap[mkMountPointShared].(bool))
			size := mountPointMap[mkMountPointSize].(string)
			volume := mountPointMap[mkMountPointVolume].(string)

			mountPointObject.ACL = &acl
			mountPointObject.Backup = &backup
			mountPointObject.MountPoint = path
			mountPointObject.Quota = &quota
			mountPointObject.ReadOnly = &readOnly
			mountPointObject.Replicate = &replicate
			mountPointObject.Shared = &shared

			// this is a totally hackish way to determine if the mount point is new or not during the container update.
			// an attached storage-backed MP has volume in the format "storage:disk file", i.e. `local-lvm:vm-123-disk-1`
			// while a new storage-backed MP has just plain volume name, i.e. `local-lvm`
			// device or directory MPs won't have a colon in the volume name either, and we don't need to do the special
			// handling for them.
			createNewMP := !strings.Contains(volume, ":")

			if len(size) > 0 && createNewMP {
				var ds types.DiskSize

				ds, err := types.ParseDiskSize(size)
				if err != nil {
					return diag.Errorf("invalid disk size: %s", err.Error())
				}

				mountPointObject.Volume = fmt.Sprintf("%s:%d", volume, ds.InGigabytes())
			} else {
				mountPointObject.Volume = volume
			}

			if len(mountOptions) > 0 {
				mountOptionsArray := make([]string, 0, len(mountPoints))

				for _, option := range mountOptions {
					mountOptionsArray = append(mountOptionsArray, option.(string))
				}

				mountPointObject.MountOptions = &mountOptionsArray
			}

			mountPointArray[i] = mountPointObject
		}

		updateBody.MountPoints = mountPointArray

		rebootRequired = true
	}

	// Prepare the new network interface configuration.
	networkInterface := d.Get(mkNetworkInterface).([]interface{})

	if len(networkInterface) == 0 && len(clone) > 0 {
		networkInterface, e = containerGetExistingNetworkInterface(ctx, containerAPI)
		if e != nil {
			return diag.FromErr(e)
		}
	}

	if d.HasChange(mkInitialization) ||
		d.HasChange(mkNetworkInterface) {
		networkInterfaceArray := make(
			containers.CustomNetworkInterfaceArray,
			len(networkInterface),
		)

		for ni, nv := range networkInterface {
			networkInterfaceMap := nv.(map[string]interface{})
			networkInterfaceObject := containers.CustomNetworkInterface{}

			bridge := networkInterfaceMap[mkNetworkInterfaceBridge].(string)
			enabled := networkInterfaceMap[mkNetworkInterfaceEnabled].(bool)
			firewall := types.CustomBool(
				networkInterfaceMap[mkNetworkInterfaceFirewall].(bool),
			)
			macAddress := networkInterfaceMap[mkNetworkInterfaceMACAddress].(string)
			name := networkInterfaceMap[mkNetworkInterfaceName].(string)
			rateLimit := networkInterfaceMap[mkNetworkInterfaceRateLimit].(float64)
			vlanID := networkInterfaceMap[mkNetworkInterfaceVLANID].(int)
			mtu := networkInterfaceMap[mkNetworkInterfaceMTU].(int)

			if bridge != "" {
				networkInterfaceObject.Bridge = &bridge
			}

			networkInterfaceObject.Enabled = enabled
			networkInterfaceObject.Firewall = &firewall

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

			networkInterfaceObject.Name = name

			if rateLimit != 0 {
				networkInterfaceObject.RateLimit = &rateLimit
			}

			if vlanID != 0 {
				networkInterfaceObject.Tag = &vlanID
			}

			if mtu != 0 {
				networkInterfaceObject.MTU = &mtu
			}

			networkInterfaceArray[ni] = networkInterfaceObject
		}

		updateBody.NetworkInterfaces = networkInterfaceArray

		for i, ni := range updateBody.NetworkInterfaces {
			if !ni.Enabled {
				updateBody.Delete = append(updateBody.Delete, fmt.Sprintf("net%d", i))
			}
		}

		for i := len(updateBody.NetworkInterfaces); i < maxResourceVirtualEnvironmentContainerNetworkInterfaces; i++ {
			updateBody.Delete = append(updateBody.Delete, fmt.Sprintf("net%d", i))
		}

		rebootRequired = true
	}

	if d.HasChange(mkStartup) {
		updateBody.StartupBehavior = containerGetStartupBehavior(d)
		if updateBody.StartupBehavior == nil {
			updateBody.Delete = append(updateBody.Delete, "startup")
		}
	}

	// Prepare the new operating system configuration.
	if d.HasChange(mkOperatingSystem) {
		operatingSystem, err := structure.GetSchemaBlock(
			container,
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

	if d.HasChange(mkProtection) {
		protection := types.CustomBool(d.Get(mkProtection).(bool))
		updateBody.Protection = &protection
	}

	if d.HasChange(mkTags) {
		tagString := containerGetTagsString(d)
		updateBody.Tags = &tagString
	}

	// Update the configuration now that everything has been prepared.
	e = containerAPI.UpdateContainer(ctx, &updateBody)
	if e != nil {
		return diag.FromErr(e)
	}

	// Determine if the state of the container needs to be changed.
	started := d.Get(mkStarted).(bool)

	if d.HasChange(mkStarted) && !bool(template) {
		if started {
			e = containerAPI.StartContainer(ctx)
			if e != nil {
				return diag.FromErr(e)
			}
		} else {
			forceStop := types.CustomBool(true)
			shutdownTimeout := 300

			e = containerAPI.ShutdownContainer(ctx, &containers.ShutdownRequestBody{
				ForceStop: &forceStop,
				Timeout:   &shutdownTimeout,
			})
			if e != nil {
				return diag.FromErr(e)
			}

			e = containerAPI.WaitForContainerStatus(ctx, "stopped")
			if e != nil {
				return diag.FromErr(e)
			}

			rebootRequired = false
		}
	}

	// As a final step in the update procedure, we might need to reboot the container.
	if !bool(template) && started && rebootRequired {
		rebootTimeout := 300

		e = containerAPI.RebootContainer(
			ctx,
			&containers.RebootRequestBody{
				Timeout: &rebootTimeout,
			},
		)
		if e != nil {
			return diag.FromErr(e)
		}
	}

	return containerRead(ctx, d, m)
}

func containerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	deleteTimeoutSec := d.Get(mkTimeoutDelete).(int)

	ctx, cancel := context.WithTimeout(ctx, time.Duration(deleteTimeoutSec)*time.Second)
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

	containerAPI := client.Node(nodeName).Container(vmID)

	// Shut down the container before deleting it.
	status, err := containerAPI.GetContainerStatus(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	if status.Status != "stopped" {
		forceStop := types.CustomBool(true)

		err = containerAPI.ShutdownContainer(
			ctx,
			&containers.ShutdownRequestBody{
				ForceStop: &forceStop,
				Timeout:   &deleteTimeoutSec,
			},
		)
		if err != nil {
			return diag.FromErr(err)
		}

		err = containerAPI.WaitForContainerStatus(ctx, "stopped")
		if err != nil {
			return diag.FromErr(err)
		}
	}

	err = containerAPI.DeleteContainer(ctx)
	if err != nil {
		if errors.Is(err, api.ErrResourceDoesNotExist) {
			d.SetId("")

			return nil
		}

		return diag.FromErr(err)
	}

	// Wait for the state to become unavailable as that clearly indicates the destruction of the container.
	err = containerAPI.WaitForContainerStatus(ctx, "")
	if err == nil {
		return diag.Errorf("failed to delete container \"%d\"", vmID)
	}

	d.SetId("")

	return nil
}

func parseImportIDWithNodeName(id string) (string, string, error) {
	nodeName, id, found := strings.Cut(id, "/")

	if !found {
		return "", "", fmt.Errorf("unexpected format of ID (%s), expected node/id", id)
	}

	return nodeName, id, nil
}
