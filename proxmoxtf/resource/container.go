/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package resource

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/containers"
	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
	"github.com/bpg/terraform-provider-proxmox/proxmoxtf"
	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/resource/validator"
	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/structure"
)

const (
	dvResourceVirtualEnvironmentContainerCloneDatastoreID                  = ""
	dvResourceVirtualEnvironmentContainerCloneNodeName                     = ""
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
	dvResourceVirtualEnvironmentContainerDiskDatastoreID                   = "local"
	dvResourceVirtualEnvironmentContainerDiskSize                          = 4
	dvResourceVirtualEnvironmentContainerFeaturesNesting                   = false
	dvResourceVirtualEnvironmentContainerFeaturesKeyControl                = false
	dvResourceVirtualEnvironmentContainerFeaturesFUSE                      = false
	dvResourceVirtualEnvironmentContainerMemoryDedicated                   = 512
	dvResourceVirtualEnvironmentContainerMemorySwap                        = 0
	dvResourceVirtualEnvironmentContainerMountPointACL                     = false
	dvResourceVirtualEnvironmentContainerMountPointBackup                  = true
	dvResourceVirtualEnvironmentContainerMountPointPath                    = ""
	dvResourceVirtualEnvironmentContainerMountPointQuota                   = false
	dvResourceVirtualEnvironmentContainerMountPointReadOnly                = false
	dvResourceVirtualEnvironmentContainerMountPointReplicate               = true
	dvResourceVirtualEnvironmentContainerMountPointShared                  = false
	dvResourceVirtualEnvironmentContainerMountPointSize                    = ""
	dvResourceVirtualEnvironmentContainerNetworkInterfaceBridge            = "vmbr0"
	dvResourceVirtualEnvironmentContainerNetworkInterfaceEnabled           = true
	dvResourceVirtualEnvironmentContainerNetworkInterfaceFirewall          = false
	dvResourceVirtualEnvironmentContainerNetworkInterfaceMACAddress        = ""
	dvResourceVirtualEnvironmentContainerNetworkInterfaceRateLimit         = 0
	dvResourceVirtualEnvironmentContainerNetworkInterfaceVLANID            = 0
	dvResourceVirtualEnvironmentContainerOperatingSystemType               = "unmanaged"
	dvResourceVirtualEnvironmentContainerPoolID                            = ""
	dvResourceVirtualEnvironmentContainerStarted                           = true
	dvResourceVirtualEnvironmentContainerStartupOrder                      = -1
	dvResourceVirtualEnvironmentContainerStartupUpDelay                    = -1
	dvResourceVirtualEnvironmentContainerStartupDownDelay                  = -1
	dvResourceVirtualEnvironmentContainerStartOnBoot                       = true
	dvResourceVirtualEnvironmentContainerTemplate                          = false
	dvResourceVirtualEnvironmentContainerUnprivileged                      = false
	dvResourceVirtualEnvironmentContainerVMID                              = -1

	maxResourceVirtualEnvironmentContainerNetworkInterfaces = 8

	mkResourceVirtualEnvironmentContainerClone                             = "clone"
	mkResourceVirtualEnvironmentContainerCloneDatastoreID                  = "datastore_id"
	mkResourceVirtualEnvironmentContainerCloneNodeName                     = "node_name"
	mkResourceVirtualEnvironmentContainerCloneVMID                         = "vm_id"
	mkResourceVirtualEnvironmentContainerConsole                           = "console"
	mkResourceVirtualEnvironmentContainerConsoleEnabled                    = "enabled"
	mkResourceVirtualEnvironmentContainerConsoleMode                       = "type"
	mkResourceVirtualEnvironmentContainerConsoleTTYCount                   = "tty_count"
	mkResourceVirtualEnvironmentContainerCPU                               = "cpu"
	mkResourceVirtualEnvironmentContainerCPUArchitecture                   = "architecture"
	mkResourceVirtualEnvironmentContainerCPUCores                          = "cores"
	mkResourceVirtualEnvironmentContainerCPUUnits                          = "units"
	mkResourceVirtualEnvironmentContainerDescription                       = "description"
	mkResourceVirtualEnvironmentContainerDisk                              = "disk"
	mkResourceVirtualEnvironmentContainerDiskDatastoreID                   = "datastore_id"
	mkResourceVirtualEnvironmentContainerDiskSize                          = "size"
	mkResourceVirtualEnvironmentContainerFeatures                          = "features"
	mkResourceVirtualEnvironmentContainerFeaturesNesting                   = "nesting"
	mkResourceVirtualEnvironmentContainerFeaturesKeyControl                = "keyctl"
	mkResourceVirtualEnvironmentContainerFeaturesFUSE                      = "fuse"
	mkResourceVirtualEnvironmentContainerFeaturesMountTypes                = "mount"
	mkResourceVirtualEnvironmentContainerInitialization                    = "initialization"
	mkResourceVirtualEnvironmentContainerInitializationDNS                 = "dns"
	mkResourceVirtualEnvironmentContainerInitializationDNSDomain           = "domain"
	mkResourceVirtualEnvironmentContainerInitializationDNSServer           = "server"
	mkResourceVirtualEnvironmentContainerInitializationDNSServers          = "servers"
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
	mkResourceVirtualEnvironmentContainerMemory                            = "memory"
	mkResourceVirtualEnvironmentContainerMemoryDedicated                   = "dedicated"
	mkResourceVirtualEnvironmentContainerMemorySwap                        = "swap"
	mkResourceVirtualEnvironmentContainerMountPoint                        = "mount_point"
	mkResourceVirtualEnvironmentContainerMountPointACL                     = "acl"
	mkResourceVirtualEnvironmentContainerMountPointBackup                  = "backup"
	mkResourceVirtualEnvironmentContainerMountPointMountOptions            = "mount_options"
	mkResourceVirtualEnvironmentContainerMountPointPath                    = "path"
	mkResourceVirtualEnvironmentContainerMountPointQuota                   = "quota"
	mkResourceVirtualEnvironmentContainerMountPointReadOnly                = "read_only"
	mkResourceVirtualEnvironmentContainerMountPointReplicate               = "replicate"
	mkResourceVirtualEnvironmentContainerMountPointShared                  = "shared"
	mkResourceVirtualEnvironmentContainerMountPointSize                    = "size"
	mkResourceVirtualEnvironmentContainerMountPointVolume                  = "volume"
	mkResourceVirtualEnvironmentContainerNetworkInterface                  = "network_interface"
	mkResourceVirtualEnvironmentContainerNetworkInterfaceBridge            = "bridge"
	mkResourceVirtualEnvironmentContainerNetworkInterfaceEnabled           = "enabled"
	mkResourceVirtualEnvironmentContainerNetworkInterfaceFirewall          = "firewall"
	mkResourceVirtualEnvironmentContainerNetworkInterfaceMACAddress        = "mac_address"
	mkResourceVirtualEnvironmentContainerNetworkInterfaceName              = "name"
	mkResourceVirtualEnvironmentContainerNetworkInterfaceRateLimit         = "rate_limit"
	mkResourceVirtualEnvironmentContainerNetworkInterfaceVLANID            = "vlan_id"
	mkResourceVirtualEnvironmentContainerNetworkInterfaceMTU               = "mtu"
	mkResourceVirtualEnvironmentContainerNodeName                          = "node_name"
	mkResourceVirtualEnvironmentContainerOperatingSystem                   = "operating_system"
	mkResourceVirtualEnvironmentContainerOperatingSystemTemplateFileID     = "template_file_id"
	mkResourceVirtualEnvironmentContainerOperatingSystemType               = "type"
	mkResourceVirtualEnvironmentContainerPoolID                            = "pool_id"
	mkResourceVirtualEnvironmentContainerStarted                           = "started"
	mkResourceVirtualEnvironmentContainerStartup                           = "startup"
	mkResourceVirtualEnvironmentContainerStartupOrder                      = "order"
	mkResourceVirtualEnvironmentContainerStartupUpDelay                    = "up_delay"
	mkResourceVirtualEnvironmentContainerStartupDownDelay                  = "down_delay"
	mkResourceVirtualEnvironmentContainerStartOnBoot                       = "start_on_boot"
	mkResourceVirtualEnvironmentContainerTags                              = "tags"
	mkResourceVirtualEnvironmentContainerTemplate                          = "template"
	mkResourceVirtualEnvironmentContainerUnprivileged                      = "unprivileged"
	mkResourceVirtualEnvironmentContainerVMID                              = "vm_id"
)

// Container returns a resource that manages a container.
func Container() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			mkResourceVirtualEnvironmentContainerClone: {
				Type:        schema.TypeList,
				Description: "The cloning configuration",
				Optional:    true,
				DefaultFunc: func() (interface{}, error) {
					return []interface{}{}, nil
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mkResourceVirtualEnvironmentContainerCloneDatastoreID: {
							Type:        schema.TypeString,
							Description: "The ID of the target datastore",
							Optional:    true,
							ForceNew:    true,
							Default:     dvResourceVirtualEnvironmentContainerCloneDatastoreID,
						},
						mkResourceVirtualEnvironmentContainerCloneNodeName: {
							Type:        schema.TypeString,
							Description: "The name of the source node",
							Optional:    true,
							ForceNew:    true,
							Default:     dvResourceVirtualEnvironmentContainerCloneNodeName,
						},
						mkResourceVirtualEnvironmentContainerCloneVMID: {
							Type:             schema.TypeInt,
							Description:      "The ID of the source container",
							Required:         true,
							ForceNew:         true,
							ValidateDiagFunc: validator.VMID(),
						},
					},
				},
				MaxItems: 1,
				MinItems: 0,
			},
			mkResourceVirtualEnvironmentContainerConsole: {
				Type:        schema.TypeList,
				Description: "The console configuration",
				Optional:    true,
				DefaultFunc: func() (interface{}, error) {
					return []interface{}{
						map[string]interface{}{
							mkResourceVirtualEnvironmentContainerConsoleEnabled:  dvResourceVirtualEnvironmentContainerConsoleEnabled,
							mkResourceVirtualEnvironmentContainerConsoleMode:     dvResourceVirtualEnvironmentContainerConsoleMode,
							mkResourceVirtualEnvironmentContainerConsoleTTYCount: dvResourceVirtualEnvironmentContainerConsoleTTYCount,
						},
					}, nil
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
							Type:             schema.TypeString,
							Description:      "The console mode",
							Optional:         true,
							Default:          dvResourceVirtualEnvironmentContainerConsoleMode,
							ValidateDiagFunc: containerGetConsoleModeValidator(),
						},
						mkResourceVirtualEnvironmentContainerConsoleTTYCount: {
							Type:             schema.TypeInt,
							Description:      "The number of available TTY",
							Optional:         true,
							Default:          dvResourceVirtualEnvironmentContainerConsoleTTYCount,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(0, 6)),
						},
					},
				},
				MaxItems: 1,
				MinItems: 0,
			},
			mkResourceVirtualEnvironmentContainerCPU: {
				Type:        schema.TypeList,
				Description: "The CPU allocation",
				Optional:    true,
				DefaultFunc: func() (interface{}, error) {
					return []interface{}{
						map[string]interface{}{
							mkResourceVirtualEnvironmentContainerCPUArchitecture: dvResourceVirtualEnvironmentContainerCPUArchitecture,
							mkResourceVirtualEnvironmentContainerCPUCores:        dvResourceVirtualEnvironmentContainerCPUCores,
							mkResourceVirtualEnvironmentContainerCPUUnits:        dvResourceVirtualEnvironmentContainerCPUUnits,
						},
					}, nil
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mkResourceVirtualEnvironmentContainerCPUArchitecture: {
							Type:             schema.TypeString,
							Description:      "The CPU architecture",
							Optional:         true,
							Default:          dvResourceVirtualEnvironmentContainerCPUArchitecture,
							ValidateDiagFunc: containerGetCPUArchitectureValidator(),
						},
						mkResourceVirtualEnvironmentContainerCPUCores: {
							Type:             schema.TypeInt,
							Description:      "The number of CPU cores",
							Optional:         true,
							Default:          dvResourceVirtualEnvironmentContainerCPUCores,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(1, 128)),
						},
						mkResourceVirtualEnvironmentContainerCPUUnits: {
							Type:        schema.TypeInt,
							Description: "The CPU units",
							Optional:    true,
							Default:     dvResourceVirtualEnvironmentContainerCPUUnits,
							ValidateDiagFunc: validation.ToDiagFunc(
								validation.IntBetween(0, 500000),
							),
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
				StateFunc: func(i interface{}) string {
					// PVE always adds a newline to the description, so we have to do the same,
					// also taking in account the CLRF case (Windows)
					if i.(string) != "" {
						return strings.ReplaceAll(strings.TrimSpace(i.(string)), "\r\n", "\n") + "\n"
					}
					return ""
				},
			},
			mkResourceVirtualEnvironmentContainerDisk: {
				Type:        schema.TypeList,
				Description: "The disks",
				Optional:    true,
				ForceNew:    true,
				DefaultFunc: func() (interface{}, error) {
					return []interface{}{
						map[string]interface{}{
							mkResourceVirtualEnvironmentContainerDiskDatastoreID: dvResourceVirtualEnvironmentContainerDiskDatastoreID,
							mkResourceVirtualEnvironmentContainerDiskSize:        dvResourceVirtualEnvironmentContainerDiskSize,
						},
					}, nil
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mkResourceVirtualEnvironmentContainerDiskDatastoreID: {
							Type:        schema.TypeString,
							Description: "The datastore id",
							Optional:    true,
							ForceNew:    true,
							Default:     dvResourceVirtualEnvironmentContainerDiskDatastoreID,
						},
						mkResourceVirtualEnvironmentContainerDiskSize: {
							Type:             schema.TypeInt,
							Description:      "The rootfs size in gigabytes",
							Optional:         true,
							ForceNew:         true,
							Default:          dvResourceVirtualEnvironmentContainerDiskSize,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(1)),
						},
					},
				},
				MaxItems: 1,
				MinItems: 0,
			},
			mkResourceVirtualEnvironmentContainerFeatures: {
				Type:        schema.TypeList,
				Description: "Features",
				Optional:    true,
				DefaultFunc: func() (interface{}, error) {
					return []interface{}{
						map[string]interface{}{
							mkResourceVirtualEnvironmentContainerFeaturesNesting:    dvResourceVirtualEnvironmentContainerFeaturesNesting,
							mkResourceVirtualEnvironmentContainerFeaturesKeyControl: dvResourceVirtualEnvironmentContainerFeaturesKeyControl,
							mkResourceVirtualEnvironmentContainerFeaturesFUSE:       dvResourceVirtualEnvironmentContainerFeaturesFUSE,
							mkResourceVirtualEnvironmentContainerFeaturesMountTypes: []interface{}{},
						},
					}, nil
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mkResourceVirtualEnvironmentContainerFeaturesNesting: {
							Type:        schema.TypeBool,
							Description: "Whether the container runs as nested",
							Optional:    true,
							Default:     dvResourceVirtualEnvironmentContainerFeaturesNesting,
						},
						mkResourceVirtualEnvironmentContainerFeaturesKeyControl: {
							Type:        schema.TypeBool,
							Description: "Whether the container supports `keyctl()` system call",
							Optional:    true,
							Default:     dvResourceVirtualEnvironmentContainerFeaturesKeyControl,
						},
						mkResourceVirtualEnvironmentContainerFeaturesFUSE: {
							Type:        schema.TypeBool,
							Description: "Whether the container supports FUSE mounts",
							Optional:    true,
							Default:     dvResourceVirtualEnvironmentContainerFeaturesFUSE,
						},
						mkResourceVirtualEnvironmentContainerFeaturesMountTypes: {
							Type:        schema.TypeList,
							Description: "List of allowed mount types",
							Optional:    true,
							Elem: &schema.Schema{
								Type:             schema.TypeString,
								ValidateDiagFunc: validator.MountType(),
							},
						},
					},
				},
				MaxItems: 1,
				MinItems: 0,
			},
			mkResourceVirtualEnvironmentContainerInitialization: {
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
										Deprecated: "The `server` attribute is deprecated and will be removed in a future release. " +
											"Please use the `servers` attribute instead.",
										Optional: true,
										Default:  dvResourceVirtualEnvironmentContainerInitializationDNSServer,
									},
									mkResourceVirtualEnvironmentContainerInitializationDNSServers: {
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
											return len(old) > 0 &&
												strings.ReplaceAll(old, "*", "") == ""
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
			mkResourceVirtualEnvironmentContainerMemory: {
				Type:        schema.TypeList,
				Description: "The memory allocation",
				Optional:    true,
				DefaultFunc: func() (interface{}, error) {
					return []interface{}{
						map[string]interface{}{
							mkResourceVirtualEnvironmentContainerMemoryDedicated: dvResourceVirtualEnvironmentContainerMemoryDedicated,
							mkResourceVirtualEnvironmentContainerMemorySwap:      dvResourceVirtualEnvironmentContainerMemorySwap,
						},
					}, nil
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mkResourceVirtualEnvironmentContainerMemoryDedicated: {
							Type:        schema.TypeInt,
							Description: "The dedicated memory in megabytes",
							Optional:    true,
							Default:     dvResourceVirtualEnvironmentContainerMemoryDedicated,
							ValidateDiagFunc: validation.ToDiagFunc(
								validation.IntBetween(16, 268435456),
							),
						},
						mkResourceVirtualEnvironmentContainerMemorySwap: {
							Type:        schema.TypeInt,
							Description: "The swap size in megabytes",
							Optional:    true,
							Default:     dvResourceVirtualEnvironmentContainerMemorySwap,
							ValidateDiagFunc: validation.ToDiagFunc(
								validation.IntBetween(0, 268435456),
							),
						},
					},
				},
				MaxItems: 1,
				MinItems: 0,
			},
			mkResourceVirtualEnvironmentContainerMountPoint: {
				Type:        schema.TypeList,
				Description: "A mount point",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mkResourceVirtualEnvironmentContainerMountPointACL: {
							Type:        schema.TypeBool,
							Description: "Explicitly enable or disable ACL support",
							Optional:    true,
							Default:     dvResourceVirtualEnvironmentContainerMountPointACL,
						},
						mkResourceVirtualEnvironmentContainerMountPointBackup: {
							Type:        schema.TypeBool,
							Description: "Whether to include the mount point in backups (only used for volume mount points)",
							Optional:    true,
							Default:     dvResourceVirtualEnvironmentContainerMountPointBackup,
						},
						mkResourceVirtualEnvironmentContainerMountPointMountOptions: {
							Type:        schema.TypeList,
							Description: "Extra mount options.",
							Optional:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						mkResourceVirtualEnvironmentContainerMountPointPath: {
							Type:        schema.TypeString,
							Description: "Path to the mount point as seen from inside the container",
							Required:    true,
							// StateFunc: func(i interface{}) string {
							// 	// PVE strips leading slashes from the path, so we have to do the same
							// 	return strings.TrimPrefix(i.(string), "/")
							// },
							DiffSuppressFunc: func(k, oldVal, newVal string, d *schema.ResourceData) bool {
								return "/"+oldVal == newVal
							},
						},
						mkResourceVirtualEnvironmentContainerMountPointQuota: {
							Type:        schema.TypeBool,
							Description: "Enable user quotas inside the container (not supported with volume mounts)",
							Optional:    true,
							Default:     dvResourceVirtualEnvironmentContainerMountPointQuota,
						},
						mkResourceVirtualEnvironmentContainerMountPointReadOnly: {
							Type:        schema.TypeBool,
							Description: "Read-only mount point",
							Optional:    true,
							Default:     dvResourceVirtualEnvironmentContainerMountPointReadOnly,
						},
						mkResourceVirtualEnvironmentContainerMountPointReplicate: {
							Type:        schema.TypeBool,
							Description: "Will include this volume to a storage replica job",
							Optional:    true,
							Default:     dvResourceVirtualEnvironmentContainerMountPointReplicate,
						},
						mkResourceVirtualEnvironmentContainerMountPointShared: {
							Type:        schema.TypeBool,
							Description: "Mark this non-volume mount point as available on all nodes",
							Optional:    true,
							Default:     dvResourceVirtualEnvironmentContainerMountPointShared,
						},
						mkResourceVirtualEnvironmentContainerMountPointSize: {
							Type:             schema.TypeString,
							Description:      "Volume size (only used for volume mount points)",
							Optional:         true,
							Default:          dvResourceVirtualEnvironmentContainerMountPointSize,
							ValidateDiagFunc: validator.FileSize(),
						},
						mkResourceVirtualEnvironmentContainerMountPointVolume: {
							Type:        schema.TypeString,
							Description: "Volume, device or directory to mount into the container",
							Required:    true,
							DiffSuppressFunc: func(k, oldVal, newVal string, d *schema.ResourceData) bool {
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
			mkResourceVirtualEnvironmentContainerNetworkInterface: {
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
						mkResourceVirtualEnvironmentContainerNetworkInterfaceFirewall: {
							Type:        schema.TypeBool,
							Description: "Whether this interface's firewall rules should be used.",
							Optional:    true,
							Default:     dvResourceVirtualEnvironmentContainerNetworkInterfaceFirewall,
						},
						mkResourceVirtualEnvironmentContainerNetworkInterfaceMACAddress: {
							Type:        schema.TypeString,
							Description: "The MAC address",
							Optional:    true,
							Default:     dvResourceVirtualEnvironmentContainerNetworkInterfaceMACAddress,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								return new == ""
							},
							ValidateDiagFunc: validator.MACAddress(),
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
						mkResourceVirtualEnvironmentVMNetworkDeviceMTU: {
							Type:        schema.TypeInt,
							Description: "Maximum transmission unit (MTU)",
							Optional:    true,
							Default:     dvResourceVirtualEnvironmentVMNetworkDeviceMTU,
						},
					},
				},
				MaxItems: maxResourceVirtualEnvironmentContainerNetworkInterfaces,
				MinItems: 0,
			},
			mkResourceVirtualEnvironmentContainerNodeName: {
				Type:        schema.TypeString,
				Description: "The node name",
				Required:    true,
				ForceNew:    true,
			},
			mkResourceVirtualEnvironmentContainerOperatingSystem: {
				Type:        schema.TypeList,
				Description: "The operating system configuration",
				Optional:    true,
				ForceNew:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mkResourceVirtualEnvironmentContainerOperatingSystemTemplateFileID: {
							Type:             schema.TypeString,
							Description:      "The ID of an OS template file",
							Required:         true,
							ForceNew:         true,
							ValidateDiagFunc: validator.FileID(),
						},
						mkResourceVirtualEnvironmentContainerOperatingSystemType: {
							Type:             schema.TypeString,
							Description:      "The type",
							Optional:         true,
							Default:          dvResourceVirtualEnvironmentContainerOperatingSystemType,
							ValidateDiagFunc: containerGetOperatingSystemTypeValidator(),
						},
					},
				},
				MaxItems: 1,
				MinItems: 0,
			},
			mkResourceVirtualEnvironmentContainerPoolID: {
				Type:        schema.TypeString,
				Description: "The ID of the pool to assign the container to",
				Optional:    true,
				ForceNew:    true,
				Default:     dvResourceVirtualEnvironmentContainerPoolID,
			},
			mkResourceVirtualEnvironmentContainerStarted: {
				Type:        schema.TypeBool,
				Description: "Whether to start the container",
				Optional:    true,
				Default:     dvResourceVirtualEnvironmentContainerStarted,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get(mkResourceVirtualEnvironmentContainerTemplate).(bool)
				},
			},
			mkResourceVirtualEnvironmentContainerStartup: {
				Type:        schema.TypeList,
				Description: "Defines startup and shutdown behavior of the container",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mkResourceVirtualEnvironmentContainerStartupOrder: {
							Type:        schema.TypeInt,
							Description: "A non-negative number defining the general startup order",
							Optional:    true,
							Default:     dvResourceVirtualEnvironmentContainerStartupOrder,
						},
						mkResourceVirtualEnvironmentContainerStartupUpDelay: {
							Type:        schema.TypeInt,
							Description: "A non-negative number defining the delay in seconds before the next container is started",
							Optional:    true,
							Default:     dvResourceVirtualEnvironmentContainerStartupUpDelay,
						},
						mkResourceVirtualEnvironmentContainerStartupDownDelay: {
							Type:        schema.TypeInt,
							Description: "A non-negative number defining the delay in seconds before the next container is shut down",
							Optional:    true,
							Default:     dvResourceVirtualEnvironmentContainerStartupDownDelay,
						},
					},
				},
				MaxItems: 1,
				MinItems: 0,
			},
			mkResourceVirtualEnvironmentContainerStartOnBoot: {
				Type:        schema.TypeBool,
				Description: "Automatically start container when the host system boots.",
				Optional:    true,
				ForceNew:    false,
				Default:     dvResourceVirtualEnvironmentContainerStartOnBoot,
			},
			mkResourceVirtualEnvironmentContainerTags: {
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
			mkResourceVirtualEnvironmentContainerTemplate: {
				Type:        schema.TypeBool,
				Description: "Whether to create a template",
				Optional:    true,
				ForceNew:    true,
				Default:     dvResourceVirtualEnvironmentContainerTemplate,
			},
			mkResourceVirtualEnvironmentContainerUnprivileged: {
				Type:        schema.TypeBool,
				Description: "Whether the container runs as unprivileged on the host",
				Optional:    true,
				ForceNew:    true,
				Default:     dvResourceVirtualEnvironmentContainerUnprivileged,
			},
			mkResourceVirtualEnvironmentContainerVMID: {
				Type:             schema.TypeInt,
				Description:      "The VM identifier",
				Optional:         true,
				ForceNew:         true,
				Default:          dvResourceVirtualEnvironmentContainerVMID,
				ValidateDiagFunc: validator.VMID(),
			},
		},
		CreateContext: containerCreate,
		ReadContext:   containerRead,
		UpdateContext: containerUpdate,
		DeleteContext: containerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, i interface{}) ([]*schema.ResourceData, error) {
				node, id, err := parseImportIDWithNodeName(d.Id())
				if err != nil {
					return nil, err
				}

				d.SetId(id)
				err = d.Set(mkResourceVirtualEnvironmentContainerNodeName, node)
				if err != nil {
					return nil, fmt.Errorf("failed setting state during import: %w", err)
				}

				return []*schema.ResourceData{d}, nil
			},
		},
	}
}

func containerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clone := d.Get(mkResourceVirtualEnvironmentContainerClone).([]interface{})

	if len(clone) > 0 {
		return containerCreateClone(ctx, d, m)
	}

	return containerCreateCustom(ctx, d, m)
}

func containerCreateClone(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(proxmoxtf.ProviderConfiguration)
	api, err := config.GetClient()
	if err != nil {
		return diag.FromErr(err)
	}

	clone := d.Get(mkResourceVirtualEnvironmentContainerClone).([]interface{})
	cloneBlock := clone[0].(map[string]interface{})
	cloneDatastoreID := cloneBlock[mkResourceVirtualEnvironmentContainerCloneDatastoreID].(string)
	cloneNodeName := cloneBlock[mkResourceVirtualEnvironmentContainerCloneNodeName].(string)
	cloneVMID := cloneBlock[mkResourceVirtualEnvironmentContainerCloneVMID].(int)

	description := d.Get(mkResourceVirtualEnvironmentContainerDescription).(string)

	initialization := d.Get(mkResourceVirtualEnvironmentContainerInitialization).([]interface{})
	initializationHostname := ""

	if len(initialization) > 0 {
		initializationBlock := initialization[0].(map[string]interface{})
		initializationHostname = initializationBlock[mkResourceVirtualEnvironmentContainerInitializationHostname].(string)
	}

	nodeName := d.Get(mkResourceVirtualEnvironmentContainerNodeName).(string)
	poolID := d.Get(mkResourceVirtualEnvironmentContainerPoolID).(string)
	tags := d.Get(mkResourceVirtualEnvironmentContainerTags).([]interface{})
	vmID := d.Get(mkResourceVirtualEnvironmentContainerVMID).(int)

	if vmID == -1 {
		vmIDNew, e := api.Cluster().GetVMID(ctx)
		if e != nil {
			return diag.FromErr(e)
		}

		vmID = *vmIDNew
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

		err = api.Node(cloneNodeName).Container(cloneVMID).CloneContainer(ctx, cloneBody)
	} else {
		err = api.Node(nodeName).Container(cloneVMID).CloneContainer(ctx, cloneBody)
	}

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(vmID))

	containerAPI := api.Node(nodeName).Container(vmID)

	// Wait for the container to be created and its configuration lock to be released.
	err = containerAPI.WaitForContainerLock(ctx, 600, 5, true)
	if err != nil {
		return diag.FromErr(err)
	}

	// Now that the virtual machine has been cloned, we need to perform some modifications.
	updateBody := &containers.UpdateRequestBody{}

	startOnBoot := types.CustomBool(d.Get(mkResourceVirtualEnvironmentContainerStartOnBoot).(bool))
	updateBody.StartOnBoot = &startOnBoot

	updateBody.StartupBehavior = containerGetStartupBehavior(d)

	console := d.Get(mkResourceVirtualEnvironmentContainerConsole).([]interface{})

	if len(console) > 0 {
		consoleBlock := console[0].(map[string]interface{})

		consoleEnabled := types.CustomBool(
			consoleBlock[mkResourceVirtualEnvironmentContainerConsoleEnabled].(bool),
		)
		consoleMode := consoleBlock[mkResourceVirtualEnvironmentContainerConsoleMode].(string)
		consoleTTYCount := consoleBlock[mkResourceVirtualEnvironmentContainerConsoleTTYCount].(int)

		updateBody.ConsoleEnabled = &consoleEnabled
		updateBody.ConsoleMode = &consoleMode
		updateBody.TTY = &consoleTTYCount
	}

	cpu := d.Get(mkResourceVirtualEnvironmentContainerCPU).([]interface{})

	if len(cpu) > 0 {
		cpuBlock := cpu[0].(map[string]interface{})

		cpuArchitecture := cpuBlock[mkResourceVirtualEnvironmentContainerCPUArchitecture].(string)
		cpuCores := cpuBlock[mkResourceVirtualEnvironmentContainerCPUCores].(int)
		cpuUnits := cpuBlock[mkResourceVirtualEnvironmentContainerCPUUnits].(int)

		updateBody.CPUArchitecture = &cpuArchitecture
		updateBody.CPUCores = &cpuCores
		updateBody.CPUUnits = &cpuUnits
	}

	var initializationIPConfigIPv4Address []string
	var initializationIPConfigIPv4Gateway []string
	var initializationIPConfigIPv6Address []string
	var initializationIPConfigIPv6Gateway []string

	if len(initialization) > 0 {
		initializationBlock := initialization[0].(map[string]interface{})
		initializationDNS := initializationBlock[mkResourceVirtualEnvironmentContainerInitializationDNS].([]interface{})

		if len(initializationDNS) > 0 {
			initializationDNSBlock := initializationDNS[0].(map[string]interface{})
			initializationDNSDomain := initializationDNSBlock[mkResourceVirtualEnvironmentContainerInitializationDNSDomain].(string)
			updateBody.DNSDomain = &initializationDNSDomain

			servers := initializationDNSBlock[mkResourceVirtualEnvironmentContainerInitializationDNSServers].([]interface{})
			deprecatedServer := initializationDNSBlock[mkResourceVirtualEnvironmentContainerInitializationDNSServer].(string)

			if len(servers) > 0 {
				nameserver := strings.Join(ConvertToStringSlice(servers), " ")

				updateBody.DNSServer = &nameserver
			} else {
				updateBody.DNSServer = &deprecatedServer
			}
		}

		initializationHostname := initializationBlock[mkResourceVirtualEnvironmentContainerInitializationHostname].(string)

		if initializationHostname != dvResourceVirtualEnvironmentContainerInitializationHostname {
			updateBody.Hostname = &initializationHostname
		}

		initializationIPConfig := initializationBlock[mkResourceVirtualEnvironmentContainerInitializationIPConfig].([]interface{})

		for _, c := range initializationIPConfig {
			configBlock := c.(map[string]interface{})
			ipv4 := configBlock[mkResourceVirtualEnvironmentContainerInitializationIPConfigIPv4].([]interface{})

			if len(ipv4) > 0 {
				ipv4Block := ipv4[0].(map[string]interface{})

				initializationIPConfigIPv4Address = append(
					initializationIPConfigIPv4Address,
					ipv4Block[mkResourceVirtualEnvironmentContainerInitializationIPConfigIPv4Address].(string),
				)

				initializationIPConfigIPv4Gateway = append(
					initializationIPConfigIPv4Gateway,
					ipv4Block[mkResourceVirtualEnvironmentContainerInitializationIPConfigIPv4Gateway].(string),
				)
			} else {
				initializationIPConfigIPv4Address = append(initializationIPConfigIPv4Address, "")
				initializationIPConfigIPv4Gateway = append(initializationIPConfigIPv4Gateway, "")
			}

			ipv6 := configBlock[mkResourceVirtualEnvironmentContainerInitializationIPConfigIPv6].([]interface{})

			if len(ipv6) > 0 {
				ipv6Block := ipv6[0].(map[string]interface{})

				initializationIPConfigIPv6Address = append(
					initializationIPConfigIPv6Address,
					ipv6Block[mkResourceVirtualEnvironmentContainerInitializationIPConfigIPv6Address].(string),
				)

				initializationIPConfigIPv6Gateway = append(
					initializationIPConfigIPv6Gateway,
					ipv6Block[mkResourceVirtualEnvironmentContainerInitializationIPConfigIPv6Gateway].(string),
				)
			} else {
				initializationIPConfigIPv6Address = append(initializationIPConfigIPv6Address, "")
				initializationIPConfigIPv6Gateway = append(initializationIPConfigIPv6Gateway, "")
			}
		}

		initializationUserAccount := initializationBlock[mkResourceVirtualEnvironmentContainerInitializationUserAccount].([]interface{})

		if len(initializationUserAccount) > 0 {
			initializationUserAccountBlock := initializationUserAccount[0].(map[string]interface{})
			keys := initializationUserAccountBlock[mkResourceVirtualEnvironmentContainerInitializationUserAccountKeys].([]interface{})

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

			initializationUserAccountPassword := initializationUserAccountBlock[mkResourceVirtualEnvironmentContainerInitializationUserAccountPassword].(string)

			if initializationUserAccountPassword != dvResourceVirtualEnvironmentContainerInitializationUserAccountPassword {
				updateBody.Password = &initializationUserAccountPassword
			} else {
				updateBody.Delete = append(updateBody.Delete, "password")
			}
		}
	}

	memory := d.Get(mkResourceVirtualEnvironmentContainerMemory).([]interface{})

	if len(memory) > 0 {
		memoryBlock := memory[0].(map[string]interface{})

		memoryDedicated := memoryBlock[mkResourceVirtualEnvironmentContainerMemoryDedicated].(int)
		memorySwap := memoryBlock[mkResourceVirtualEnvironmentContainerMemorySwap].(int)

		updateBody.DedicatedMemory = &memoryDedicated
		updateBody.Swap = &memorySwap
	}

	networkInterface := d.Get(mkResourceVirtualEnvironmentContainerNetworkInterface).([]interface{})

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

		bridge := networkInterfaceMap[mkResourceVirtualEnvironmentContainerNetworkInterfaceBridge].(string)
		enabled := networkInterfaceMap[mkResourceVirtualEnvironmentContainerNetworkInterfaceEnabled].(bool)
		firewall := types.CustomBool(
			networkInterfaceMap[mkResourceVirtualEnvironmentContainerNetworkInterfaceFirewall].(bool),
		)
		macAddress := networkInterfaceMap[mkResourceVirtualEnvironmentContainerNetworkInterfaceMACAddress].(string)
		name := networkInterfaceMap[mkResourceVirtualEnvironmentContainerNetworkInterfaceName].(string)
		rateLimit := networkInterfaceMap[mkResourceVirtualEnvironmentContainerNetworkInterfaceRateLimit].(float64)
		vlanID := networkInterfaceMap[mkResourceVirtualEnvironmentContainerNetworkInterfaceVLANID].(int)
		mtu, _ := networkInterfaceMap[mkResourceVirtualEnvironmentVMNetworkDeviceMTU].(int)

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

	for i := 0; i < len(updateBody.NetworkInterfaces); i++ {
		if !updateBody.NetworkInterfaces[i].Enabled {
			updateBody.Delete = append(updateBody.Delete, fmt.Sprintf("net%d", i))
		}
	}

	for i := len(updateBody.NetworkInterfaces); i < maxResourceVirtualEnvironmentContainerNetworkInterfaces; i++ {
		updateBody.Delete = append(updateBody.Delete, fmt.Sprintf("net%d", i))
	}

	operatingSystem := d.Get(mkResourceVirtualEnvironmentContainerOperatingSystem).([]interface{})

	if len(operatingSystem) > 0 {
		operatingSystemBlock := operatingSystem[0].(map[string]interface{})

		operatingSystemTemplateFileID := operatingSystemBlock[mkResourceVirtualEnvironmentContainerOperatingSystemTemplateFileID].(string)
		operatingSystemType := operatingSystemBlock[mkResourceVirtualEnvironmentContainerOperatingSystemType].(string)

		updateBody.OSTemplateFileVolume = &operatingSystemTemplateFileID
		updateBody.OSType = &operatingSystemType
	}

	if len(tags) > 0 {
		tagString := containerGetTagsString(d)
		updateBody.Tags = &tagString
	}

	template := types.CustomBool(d.Get(mkResourceVirtualEnvironmentContainerTemplate).(bool))

	//nolint:gosimple
	if template != dvResourceVirtualEnvironmentContainerTemplate {
		updateBody.Template = &template
	}

	err = containerAPI.UpdateContainer(ctx, updateBody)
	if err != nil {
		return diag.FromErr(err)
	}

	// Wait for the container's lock to be released.
	err = containerAPI.WaitForContainerLock(ctx, 600, 5, true)
	if err != nil {
		return diag.FromErr(err)
	}

	return containerCreateStart(ctx, d, m)
}

func containerCreateCustom(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(proxmoxtf.ProviderConfiguration)
	api, err := config.GetClient()
	if err != nil {
		return diag.FromErr(err)
	}

	nodeName := d.Get(mkResourceVirtualEnvironmentContainerNodeName).(string)
	resource := Container()

	consoleBlock, err := structure.GetSchemaBlock(
		resource,
		d,
		[]string{mkResourceVirtualEnvironmentContainerConsole},
		0,
		true,
	)
	if err != nil {
		return diag.FromErr(err)
	}

	consoleEnabled := types.CustomBool(
		consoleBlock[mkResourceVirtualEnvironmentContainerConsoleEnabled].(bool),
	)
	consoleMode := consoleBlock[mkResourceVirtualEnvironmentContainerConsoleMode].(string)
	consoleTTYCount := consoleBlock[mkResourceVirtualEnvironmentContainerConsoleTTYCount].(int)

	cpuBlock, err := structure.GetSchemaBlock(
		resource,
		d,
		[]string{mkResourceVirtualEnvironmentContainerCPU},
		0,
		true,
	)
	if err != nil {
		return diag.FromErr(err)
	}

	cpuArchitecture := cpuBlock[mkResourceVirtualEnvironmentContainerCPUArchitecture].(string)
	cpuCores := cpuBlock[mkResourceVirtualEnvironmentContainerCPUCores].(int)
	cpuUnits := cpuBlock[mkResourceVirtualEnvironmentContainerCPUUnits].(int)

	description := d.Get(mkResourceVirtualEnvironmentContainerDescription).(string)

	diskBlock, err := structure.GetSchemaBlock(
		resource,
		d,
		[]string{mkResourceVirtualEnvironmentContainerDisk},
		0,
		true,
	)
	if err != nil {
		return diag.FromErr(err)
	}

	diskDatastoreID := diskBlock[mkResourceVirtualEnvironmentContainerDiskDatastoreID].(string)

	var rootFS *containers.CustomRootFS

	diskSize := diskBlock[mkResourceVirtualEnvironmentContainerDiskSize].(int)
	if diskSize != dvResourceVirtualEnvironmentContainerDiskSize && diskDatastoreID != "" {
		// This is a special case where the rootfs size is set to a non-default value at creation time.
		// see https://pve.proxmox.com/pve-docs/chapter-pct.html#_storage_backed_mount_points
		rootFS = &containers.CustomRootFS{
			Volume: fmt.Sprintf("%s:%d", diskDatastoreID, diskSize),
		}
	}

	features, err := containerGetFeatures(resource, d)
	if err != nil {
		return diag.FromErr(err)
	}

	initialization := d.Get(mkResourceVirtualEnvironmentContainerInitialization).([]interface{})
	initializationDNSDomain := dvResourceVirtualEnvironmentContainerInitializationDNSDomain
	initializationDNSServer := dvResourceVirtualEnvironmentContainerInitializationDNSServer
	initializationHostname := dvResourceVirtualEnvironmentContainerInitializationHostname
	var initializationIPConfigIPv4Address []string
	var initializationIPConfigIPv4Gateway []string
	var initializationIPConfigIPv6Address []string
	var initializationIPConfigIPv6Gateway []string

	initializationUserAccountKeys := containers.CustomSSHKeys{}
	initializationUserAccountPassword := dvResourceVirtualEnvironmentContainerInitializationUserAccountPassword

	if len(initialization) > 0 {
		initializationBlock := initialization[0].(map[string]interface{})
		initializationDNS := initializationBlock[mkResourceVirtualEnvironmentContainerInitializationDNS].([]interface{})

		if len(initializationDNS) > 0 {
			initializationDNSBlock := initializationDNS[0].(map[string]interface{})
			initializationDNSDomain = initializationDNSBlock[mkResourceVirtualEnvironmentContainerInitializationDNSDomain].(string)

			servers := initializationDNSBlock[mkResourceVirtualEnvironmentContainerInitializationDNSServers].([]interface{})
			deprecatedServer := initializationDNSBlock[mkResourceVirtualEnvironmentContainerInitializationDNSServer].(string)

			if len(servers) > 0 {
				nameserver := strings.Join(ConvertToStringSlice(servers), " ")

				initializationDNSServer = nameserver
			} else {
				initializationDNSServer = deprecatedServer
			}
		}

		initializationHostname = initializationBlock[mkResourceVirtualEnvironmentContainerInitializationHostname].(string)
		initializationIPConfig := initializationBlock[mkResourceVirtualEnvironmentContainerInitializationIPConfig].([]interface{})

		for _, c := range initializationIPConfig {
			if c == nil {
				continue
			}

			configBlock := c.(map[string]interface{})
			ipv4 := configBlock[mkResourceVirtualEnvironmentContainerInitializationIPConfigIPv4].([]interface{})

			if len(ipv4) > 0 {
				ipv4Block := ipv4[0].(map[string]interface{})

				initializationIPConfigIPv4Address = append(
					initializationIPConfigIPv4Address,
					ipv4Block[mkResourceVirtualEnvironmentContainerInitializationIPConfigIPv4Address].(string),
				)

				initializationIPConfigIPv4Gateway = append(
					initializationIPConfigIPv4Gateway,
					ipv4Block[mkResourceVirtualEnvironmentContainerInitializationIPConfigIPv4Gateway].(string),
				)
			} else {
				initializationIPConfigIPv4Address = append(initializationIPConfigIPv4Address, "")
				initializationIPConfigIPv4Gateway = append(initializationIPConfigIPv4Gateway, "")
			}

			ipv6 := configBlock[mkResourceVirtualEnvironmentContainerInitializationIPConfigIPv6].([]interface{})

			if len(ipv6) > 0 {
				ipv6Block := ipv6[0].(map[string]interface{})

				initializationIPConfigIPv6Address = append(
					initializationIPConfigIPv6Address,
					ipv6Block[mkResourceVirtualEnvironmentContainerInitializationIPConfigIPv6Address].(string),
				)

				initializationIPConfigIPv6Gateway = append(
					initializationIPConfigIPv6Gateway,
					ipv6Block[mkResourceVirtualEnvironmentContainerInitializationIPConfigIPv6Gateway].(string),
				)
			} else {
				initializationIPConfigIPv6Address = append(initializationIPConfigIPv6Address, "")
				initializationIPConfigIPv6Gateway = append(initializationIPConfigIPv6Gateway, "")
			}
		}

		initializationUserAccount := initializationBlock[mkResourceVirtualEnvironmentContainerInitializationUserAccount].([]interface{})

		if len(initializationUserAccount) > 0 {
			initializationUserAccountBlock := initializationUserAccount[0].(map[string]interface{})

			keys := initializationUserAccountBlock[mkResourceVirtualEnvironmentContainerInitializationUserAccountKeys].([]interface{})
			initializationUserAccountKeys = make(
				containers.CustomSSHKeys,
				len(keys),
			)

			for ki, kv := range keys {
				initializationUserAccountKeys[ki] = kv.(string)
			}

			initializationUserAccountPassword = initializationUserAccountBlock[mkResourceVirtualEnvironmentContainerInitializationUserAccountPassword].(string)
		}
	}

	memoryBlock, err := structure.GetSchemaBlock(
		resource,
		d,
		[]string{mkResourceVirtualEnvironmentContainerMemory},
		0,
		true,
	)
	if err != nil {
		return diag.FromErr(err)
	}

	memoryDedicated := memoryBlock[mkResourceVirtualEnvironmentContainerMemoryDedicated].(int)
	memorySwap := memoryBlock[mkResourceVirtualEnvironmentContainerMemorySwap].(int)

	mountPoint := d.Get(mkResourceVirtualEnvironmentContainerMountPoint).([]interface{})
	mountPointArray := make(containers.CustomMountPointArray, 0, len(mountPoint))

	// because of default bool values:
	//nolint:gosimple
	for _, mp := range mountPoint {
		mountPointMap := mp.(map[string]interface{})
		mountPointObject := containers.CustomMountPoint{}

		acl := types.CustomBool(mountPointMap[mkResourceVirtualEnvironmentContainerMountPointACL].(bool))
		backup := types.CustomBool(mountPointMap[mkResourceVirtualEnvironmentContainerMountPointBackup].(bool))
		mountOptions := mountPointMap[mkResourceVirtualEnvironmentContainerMountPointMountOptions].([]interface{})
		path := mountPointMap[mkResourceVirtualEnvironmentContainerMountPointPath].(string)
		quota := types.CustomBool(mountPointMap[mkResourceVirtualEnvironmentContainerMountPointQuota].(bool))
		readOnly := types.CustomBool(mountPointMap[mkResourceVirtualEnvironmentContainerMountPointReadOnly].(bool))
		replicate := types.CustomBool(mountPointMap[mkResourceVirtualEnvironmentContainerMountPointReplicate].(bool))
		shared := types.CustomBool(mountPointMap[mkResourceVirtualEnvironmentContainerMountPointShared].(bool))
		size := mountPointMap[mkResourceVirtualEnvironmentContainerMountPointSize].(string)
		volume := mountPointMap[mkResourceVirtualEnvironmentContainerMountPointVolume].(string)

		// we have to set only the values that are different from the provider's defaults,
		if acl != dvResourceVirtualEnvironmentContainerMountPointACL {
			mountPointObject.ACL = &acl
		}

		if backup != dvResourceVirtualEnvironmentContainerMountPointBackup {
			mountPointObject.Backup = &backup
		}

		if path != dvResourceVirtualEnvironmentContainerMountPointPath {
			mountPointObject.MountPoint = path
		}

		if quota != dvResourceVirtualEnvironmentContainerMountPointQuota {
			mountPointObject.Quota = &quota
		}

		if readOnly != dvResourceVirtualEnvironmentContainerMountPointReadOnly {
			mountPointObject.ReadOnly = &readOnly
		}

		if replicate != dvResourceVirtualEnvironmentContainerMountPointReplicate {
			mountPointObject.Replicate = &replicate
		}

		if shared != dvResourceVirtualEnvironmentContainerMountPointShared {
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

	networkInterface := d.Get(mkResourceVirtualEnvironmentContainerNetworkInterface).([]interface{})
	networkInterfaceArray := make(containers.CustomNetworkInterfaceArray, len(networkInterface))

	for ni, nv := range networkInterface {
		networkInterfaceMap := nv.(map[string]interface{})
		networkInterfaceObject := containers.CustomNetworkInterface{}

		bridge := networkInterfaceMap[mkResourceVirtualEnvironmentContainerNetworkInterfaceBridge].(string)
		enabled := networkInterfaceMap[mkResourceVirtualEnvironmentContainerNetworkInterfaceEnabled].(bool)
		macAddress := networkInterfaceMap[mkResourceVirtualEnvironmentContainerNetworkInterfaceMACAddress].(string)
		name := networkInterfaceMap[mkResourceVirtualEnvironmentContainerNetworkInterfaceName].(string)
		rateLimit := networkInterfaceMap[mkResourceVirtualEnvironmentContainerNetworkInterfaceRateLimit].(float64)
		vlanID := networkInterfaceMap[mkResourceVirtualEnvironmentContainerNetworkInterfaceVLANID].(int)
		mtu := networkInterfaceMap[mkResourceVirtualEnvironmentContainerNetworkInterfaceMTU].(int)

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

	operatingSystem := d.Get(mkResourceVirtualEnvironmentContainerOperatingSystem).([]interface{})

	if len(operatingSystem) == 0 {
		return diag.Errorf(
			"\"%s\": required field is not set",
			mkResourceVirtualEnvironmentContainerOperatingSystem,
		)
	}

	operatingSystemBlock := operatingSystem[0].(map[string]interface{})
	operatingSystemTemplateFileID := operatingSystemBlock[mkResourceVirtualEnvironmentContainerOperatingSystemTemplateFileID].(string)
	operatingSystemType := operatingSystemBlock[mkResourceVirtualEnvironmentContainerOperatingSystemType].(string)

	poolID := d.Get(mkResourceVirtualEnvironmentContainerPoolID).(string)
	started := types.CustomBool(d.Get(mkResourceVirtualEnvironmentContainerStarted).(bool))
	startOnBoot := types.CustomBool(d.Get(mkResourceVirtualEnvironmentContainerStartOnBoot).(bool))
	startupBehavior := containerGetStartupBehavior(d)
	tags := d.Get(mkResourceVirtualEnvironmentContainerTags).([]interface{})
	template := types.CustomBool(d.Get(mkResourceVirtualEnvironmentContainerTemplate).(bool))
	unprivileged := types.CustomBool(d.Get(mkResourceVirtualEnvironmentContainerUnprivileged).(bool))
	vmID := d.Get(mkResourceVirtualEnvironmentContainerVMID).(int)

	if vmID == -1 {
		vmIDNew, e := api.Cluster().GetVMID(ctx)
		if e != nil {
			return diag.FromErr(e)
		}

		vmID = *vmIDNew
	}

	// Attempt to create the resource using the retrieved values.
	createBody := containers.CreateRequestBody{
		ConsoleEnabled:       &consoleEnabled,
		ConsoleMode:          &consoleMode,
		CPUArchitecture:      &cpuArchitecture,
		CPUCores:             &cpuCores,
		CPUUnits:             &cpuUnits,
		DatastoreID:          &diskDatastoreID,
		DedicatedMemory:      &memoryDedicated,
		Features:             features,
		MountPoints:          mountPointArray,
		NetworkInterfaces:    networkInterfaceArray,
		OSTemplateFileVolume: &operatingSystemTemplateFileID,
		OSType:               &operatingSystemType,
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

	err = api.Node(nodeName).Container(0).CreateContainer(ctx, &createBody, 60)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(vmID))

	// Wait for the container's lock to be released.
	err = api.Node(nodeName).Container(vmID).WaitForContainerLock(ctx, 600, 5, true)
	if err != nil {
		return diag.FromErr(err)
	}

	return containerCreateStart(ctx, d, m)
}

func containerCreateStart(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	started := d.Get(mkResourceVirtualEnvironmentContainerStarted).(bool)
	template := d.Get(mkResourceVirtualEnvironmentContainerTemplate).(bool)

	if !started || template {
		return containerRead(ctx, d, m)
	}

	config := m.(proxmoxtf.ProviderConfiguration)
	api, err := config.GetClient()
	if err != nil {
		return diag.FromErr(err)
	}

	nodeName := d.Get(mkResourceVirtualEnvironmentContainerNodeName).(string)
	vmID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	containerAPI := api.Node(nodeName).Container(vmID)

	// Start the container and wait for it to reach a running state before continuing.
	err = containerAPI.StartContainer(ctx, 60)
	if err != nil {
		return diag.FromErr(err)
	}

	return containerRead(ctx, d, m)
}

func containerGetConsoleModeValidator() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.StringInSlice([]string{
		"console",
		"shell",
		"tty",
	}, false))
}

func containerGetCPUArchitectureValidator() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.StringInSlice([]string{
		"amd64",
		"arm64",
		"armhf",
		"i386",
	}, false))
}

func containerGetExistingNetworkInterface(
	ctx context.Context,
	containerAPI *containers.Client,
) ([]interface{}, error) {
	containerInfo, err := containerAPI.GetContainer(ctx)
	if err != nil {
		return []interface{}{}, err
	}

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

		if nv.Bridge != nil {
			networkInterface[mkResourceVirtualEnvironmentContainerNetworkInterfaceBridge] = *nv.Bridge
		} else {
			networkInterface[mkResourceVirtualEnvironmentContainerNetworkInterfaceBridge] = ""
		}

		networkInterface[mkResourceVirtualEnvironmentContainerNetworkInterfaceEnabled] = true
		if nv.Firewall != nil {
			networkInterface[mkResourceVirtualEnvironmentContainerNetworkInterfaceFirewall] = *nv.Firewall
		} else {
			networkInterface[mkResourceVirtualEnvironmentContainerNetworkInterfaceFirewall] = false
		}

		if nv.MACAddress != nil {
			networkInterface[mkResourceVirtualEnvironmentContainerNetworkInterfaceMACAddress] = *nv.MACAddress
		} else {
			networkInterface[mkResourceVirtualEnvironmentContainerNetworkInterfaceMACAddress] = ""
		}

		networkInterface[mkResourceVirtualEnvironmentContainerNetworkInterfaceName] = nv.Name

		if nv.RateLimit != nil {
			networkInterface[mkResourceVirtualEnvironmentContainerNetworkInterfaceRateLimit] = *nv.RateLimit
		} else {
			networkInterface[mkResourceVirtualEnvironmentContainerNetworkInterfaceRateLimit] = float64(0)
		}

		if nv.Tag != nil {
			networkInterface[mkResourceVirtualEnvironmentContainerNetworkInterfaceVLANID] = *nv.Tag
		} else {
			networkInterface[mkResourceVirtualEnvironmentContainerNetworkInterfaceVLANID] = 0
		}

		if nv.MTU != nil {
			networkInterface[mkResourceVirtualEnvironmentContainerNetworkInterfaceMTU] = *nv.MTU
		} else {
			networkInterface[mkResourceVirtualEnvironmentContainerNetworkInterfaceMTU] = 0
		}

		networkInterfaces = append(networkInterfaces, networkInterface)
	}

	return networkInterfaces, nil
}

func containerGetOperatingSystemTypeValidator() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.StringInSlice([]string{
		"alpine",
		"archlinux",
		"centos",
		"debian",
		"fedora",
		"gentoo",
		"nixos",
		"opensuse",
		"ubuntu",
		"unmanaged",
	}, false))
}

func containerGetTagsString(d *schema.ResourceData) string {
	var sanitizedTags []string

	tags := d.Get(mkResourceVirtualEnvironmentContainerTags).([]interface{})
	for i := 0; i < len(tags); i++ {
		tag := strings.TrimSpace(tags[i].(string))
		if len(tag) > 0 {
			sanitizedTags = append(sanitizedTags, tag)
		}
	}

	sort.Strings(sanitizedTags)

	return strings.Join(sanitizedTags, ";")
}

func containerGetStartupBehavior(d *schema.ResourceData) *containers.CustomStartupBehavior {
	startup := d.Get(mkResourceVirtualEnvironmentContainerStartup).([]interface{})
	if len(startup) > 0 {
		startupBlock := startup[0].(map[string]interface{})
		startupOrder := startupBlock[mkResourceVirtualEnvironmentContainerStartupOrder].(int)
		startupUpDelay := startupBlock[mkResourceVirtualEnvironmentContainerStartupUpDelay].(int)
		startupDownDelay := startupBlock[mkResourceVirtualEnvironmentContainerStartupDownDelay].(int)

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
		[]string{mkResourceVirtualEnvironmentContainerFeatures},
		0,
		true,
	)
	if err != nil {
		return nil, fmt.Errorf("error getting container features from schema: %w", err)
	}

	nesting := types.CustomBool(featuresBlock[mkResourceVirtualEnvironmentContainerFeaturesNesting].(bool))
	keyctl := types.CustomBool(featuresBlock[mkResourceVirtualEnvironmentContainerFeaturesKeyControl].(bool))
	fuse := types.CustomBool(featuresBlock[mkResourceVirtualEnvironmentContainerFeaturesFUSE].(bool))
	mountTypes := featuresBlock[mkResourceVirtualEnvironmentContainerFeaturesMountTypes].([]interface{})

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

	api, e := config.GetClient()
	if e != nil {
		return diag.FromErr(e)
	}

	nodeName := d.Get(mkResourceVirtualEnvironmentContainerNodeName).(string)

	vmID, e := strconv.Atoi(d.Id())
	if e != nil {
		return diag.FromErr(e)
	}

	containerAPI := api.Node(nodeName).Container(vmID)

	// Retrieve the entire configuration in order to compare it to the state.
	containerConfig, e := containerAPI.GetContainer(ctx)
	if e != nil {
		if strings.Contains(e.Error(), "HTTP 404") ||
			(strings.Contains(e.Error(), "HTTP 500") && strings.Contains(e.Error(), "does not exist")) {
			d.SetId("")

			return nil
		}

		return diag.FromErr(e)
	}

	clone := d.Get(mkResourceVirtualEnvironmentContainerClone).([]interface{})

	// Compare the primitive values to those stored in the state.
	currentDescription := d.Get(mkResourceVirtualEnvironmentContainerDescription).(string)

	if len(clone) == 0 || currentDescription != dvResourceVirtualEnvironmentContainerDescription {
		if containerConfig.Description != nil {
			e = d.Set(mkResourceVirtualEnvironmentContainerDescription, *containerConfig.Description)
		} else {
			e = d.Set(mkResourceVirtualEnvironmentContainerDescription, "")
		}

		diags = append(diags, diag.FromErr(e)...)
	}

	// Compare the console configuration to the one stored in the state.
	console := map[string]interface{}{}

	if containerConfig.ConsoleEnabled != nil {
		console[mkResourceVirtualEnvironmentContainerConsoleEnabled] = *containerConfig.ConsoleEnabled
	} else {
		// Default value of "console" is "1" according to the API documentation.
		console[mkResourceVirtualEnvironmentContainerConsoleEnabled] = true
	}

	if containerConfig.ConsoleMode != nil {
		console[mkResourceVirtualEnvironmentContainerConsoleMode] = *containerConfig.ConsoleMode
	} else {
		// Default value of "cmode" is "tty" according to the API documentation.
		console[mkResourceVirtualEnvironmentContainerConsoleMode] = "tty"
	}

	if containerConfig.TTY != nil {
		console[mkResourceVirtualEnvironmentContainerConsoleTTYCount] = *containerConfig.TTY
	} else {
		// Default value of "tty" is "2" according to the API documentation.
		console[mkResourceVirtualEnvironmentContainerConsoleTTYCount] = 2
	}

	currentConsole := d.Get(mkResourceVirtualEnvironmentContainerConsole).([]interface{})

	if len(clone) > 0 {
		if len(currentConsole) > 0 {
			err := d.Set(mkResourceVirtualEnvironmentContainerConsole, []interface{}{console})
			diags = append(diags, diag.FromErr(err)...)
		}
	} else if len(currentConsole) > 0 ||
		//nolint:lll
		console[mkResourceVirtualEnvironmentContainerConsoleEnabled] != types.CustomBool(dvResourceVirtualEnvironmentContainerConsoleEnabled) ||
		console[mkResourceVirtualEnvironmentContainerConsoleMode] != dvResourceVirtualEnvironmentContainerConsoleMode ||
		console[mkResourceVirtualEnvironmentContainerConsoleTTYCount] != dvResourceVirtualEnvironmentContainerConsoleTTYCount {
		err := d.Set(mkResourceVirtualEnvironmentContainerConsole, []interface{}{console})
		diags = append(diags, diag.FromErr(err)...)
	}

	// Compare the CPU configuration to the one stored in the state.
	cpu := map[string]interface{}{}

	if containerConfig.CPUArchitecture != nil {
		cpu[mkResourceVirtualEnvironmentContainerCPUArchitecture] = *containerConfig.CPUArchitecture
	} else {
		// Default value of "arch" is "amd64" according to the API documentation.
		cpu[mkResourceVirtualEnvironmentContainerCPUArchitecture] = "amd64"
	}

	if containerConfig.CPUCores != nil {
		cpu[mkResourceVirtualEnvironmentContainerCPUCores] = *containerConfig.CPUCores
	} else {
		// Default value of "cores" is "1" according to the API documentation.
		cpu[mkResourceVirtualEnvironmentContainerCPUCores] = 1
	}

	if containerConfig.CPUUnits != nil {
		cpu[mkResourceVirtualEnvironmentContainerCPUUnits] = *containerConfig.CPUUnits
	} else {
		// Default value of "cpuunits" is "1024" according to the API documentation.
		cpu[mkResourceVirtualEnvironmentContainerCPUUnits] = 1024
	}

	currentCPU := d.Get(mkResourceVirtualEnvironmentContainerCPU).([]interface{})

	if len(clone) > 0 {
		if len(currentCPU) > 0 {
			err := d.Set(mkResourceVirtualEnvironmentContainerCPU, []interface{}{cpu})
			diags = append(diags, diag.FromErr(err)...)
		}
	} else if len(currentCPU) > 0 ||
		cpu[mkResourceVirtualEnvironmentContainerCPUArchitecture] != dvResourceVirtualEnvironmentContainerCPUArchitecture ||
		cpu[mkResourceVirtualEnvironmentContainerCPUCores] != dvResourceVirtualEnvironmentContainerCPUCores ||
		cpu[mkResourceVirtualEnvironmentContainerCPUUnits] != dvResourceVirtualEnvironmentContainerCPUUnits {
		err := d.Set(mkResourceVirtualEnvironmentContainerCPU, []interface{}{cpu})
		diags = append(diags, diag.FromErr(err)...)
	}

	// Compare the disk configuration to the one stored in the state.
	disk := map[string]interface{}{}

	if containerConfig.RootFS != nil {
		volumeParts := strings.Split(containerConfig.RootFS.Volume, ":")
		disk[mkResourceVirtualEnvironmentContainerDiskDatastoreID] = volumeParts[0]
		disk[mkResourceVirtualEnvironmentContainerDiskSize] = containerConfig.RootFS.Size.InGigabytes()
	} else {
		// Default value of "storage" is "local" according to the API documentation.
		disk[mkResourceVirtualEnvironmentContainerDiskDatastoreID] = "local"
		disk[mkResourceVirtualEnvironmentContainerDiskSize] = dvResourceVirtualEnvironmentContainerDiskSize
	}

	currentDisk := d.Get(mkResourceVirtualEnvironmentContainerDisk).([]interface{})

	if len(clone) > 0 {
		if len(currentDisk) > 0 {
			// do not override the rootfs size if it was not changed during the clone operation
			//nolint:lll
			if currentDisk[0].(map[string]interface{})[mkResourceVirtualEnvironmentContainerDiskSize] == dvResourceVirtualEnvironmentContainerDiskSize {
				disk[mkResourceVirtualEnvironmentContainerDiskSize] = dvResourceVirtualEnvironmentContainerDiskSize
			}
			err := d.Set(mkResourceVirtualEnvironmentContainerDisk, []interface{}{disk})
			diags = append(diags, diag.FromErr(err)...)
		}
	} else if len(currentDisk) > 0 ||
		disk[mkResourceVirtualEnvironmentContainerDiskDatastoreID] != dvResourceVirtualEnvironmentContainerDiskDatastoreID ||
		disk[mkResourceVirtualEnvironmentContainerDiskSize] != dvResourceVirtualEnvironmentContainerDiskSize {
		err := d.Set(mkResourceVirtualEnvironmentContainerDisk, []interface{}{disk})
		diags = append(diags, diag.FromErr(err)...)
	}

	// Compare the memory configuration to the one stored in the state.
	memory := map[string]interface{}{}

	if containerConfig.DedicatedMemory != nil {
		memory[mkResourceVirtualEnvironmentContainerMemoryDedicated] = *containerConfig.DedicatedMemory
	} else {
		memory[mkResourceVirtualEnvironmentContainerMemoryDedicated] = 0
	}

	if containerConfig.Swap != nil {
		memory[mkResourceVirtualEnvironmentContainerMemorySwap] = *containerConfig.Swap
	} else {
		memory[mkResourceVirtualEnvironmentContainerMemorySwap] = 0
	}

	currentMemory := d.Get(mkResourceVirtualEnvironmentContainerMemory).([]interface{})

	if len(clone) > 0 {
		if len(currentMemory) > 0 {
			err := d.Set(mkResourceVirtualEnvironmentContainerMemory, []interface{}{memory})
			diags = append(diags, diag.FromErr(err)...)
		}
	} else if len(currentMemory) > 0 ||
		memory[mkResourceVirtualEnvironmentContainerMemoryDedicated] != dvResourceVirtualEnvironmentContainerMemoryDedicated ||
		memory[mkResourceVirtualEnvironmentContainerMemorySwap] != dvResourceVirtualEnvironmentContainerMemorySwap {
		err := d.Set(mkResourceVirtualEnvironmentContainerMemory, []interface{}{memory})
		diags = append(diags, diag.FromErr(err)...)
	}

	// Compare the initialization and network interface configuration to the one stored in the state.
	initialization := map[string]interface{}{}

	if containerConfig.DNSDomain != nil || containerConfig.DNSServer != nil {
		initializationDNS := map[string]interface{}{}

		if containerConfig.DNSDomain != nil {
			initializationDNS[mkResourceVirtualEnvironmentContainerInitializationDNSDomain] = *containerConfig.DNSDomain
		} else {
			initializationDNS[mkResourceVirtualEnvironmentContainerInitializationDNSDomain] = ""
		}

		// check what we have in the plan
		currentInitializationDNSBlock := map[string]interface{}{}
		currentInitialization := d.Get(mkResourceVirtualEnvironmentContainerInitialization).([]interface{})

		if len(currentInitialization) > 0 {
			currentInitializationBlock := currentInitialization[0].(map[string]interface{})
			//nolint:lll
			currentInitializationDNS := currentInitializationBlock[mkResourceVirtualEnvironmentContainerInitializationDNS].([]interface{})
			if len(currentInitializationDNS) > 0 {
				currentInitializationDNSBlock = currentInitializationDNS[0].(map[string]interface{})
			}
		}

		//nolint:lll
		currentInitializationDNSServer, ok := currentInitializationDNSBlock[mkResourceVirtualEnvironmentContainerInitializationDNSServer]
		if containerConfig.DNSServer != nil {
			if ok && currentInitializationDNSServer != "" {
				initializationDNS[mkResourceVirtualEnvironmentContainerInitializationDNSServer] = *containerConfig.DNSServer
			} else {
				dnsServer := strings.Split(*containerConfig.DNSServer, " ")
				initializationDNS[mkResourceVirtualEnvironmentContainerInitializationDNSServers] = dnsServer
			}
		} else {
			initializationDNS[mkResourceVirtualEnvironmentContainerInitializationDNSServer] = ""
			initializationDNS[mkResourceVirtualEnvironmentContainerInitializationDNSServers] = []string{}
		}

		initialization[mkResourceVirtualEnvironmentContainerInitializationDNS] = []interface{}{
			initializationDNS,
		}
	}

	if containerConfig.Hostname != nil {
		initialization[mkResourceVirtualEnvironmentContainerInitializationHostname] = *containerConfig.Hostname
	} else {
		initialization[mkResourceVirtualEnvironmentContainerInitializationHostname] = ""
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
			mountPoint[mkResourceVirtualEnvironmentContainerMountPointACL] = *mp.ACL
		} else {
			mountPoint[mkResourceVirtualEnvironmentContainerMountPointACL] = false
		}

		if mp.Backup != nil {
			mountPoint[mkResourceVirtualEnvironmentContainerMountPointBackup] = *mp.Backup
		} else {
			mountPoint[mkResourceVirtualEnvironmentContainerMountPointBackup] = true
		}

		if mp.MountOptions != nil {
			mountPoint[mkResourceVirtualEnvironmentContainerMountPointMountOptions] = *mp.MountOptions
		} else {
			mountPoint[mkResourceVirtualEnvironmentContainerMountPointMountOptions] = []string{}
		}

		mountPoint[mkResourceVirtualEnvironmentContainerMountPointPath] = mp.MountPoint

		if mp.Quota != nil {
			mountPoint[mkResourceVirtualEnvironmentContainerMountPointQuota] = *mp.Quota
		} else {
			mountPoint[mkResourceVirtualEnvironmentContainerMountPointQuota] = false
		}

		if mp.ReadOnly != nil {
			mountPoint[mkResourceVirtualEnvironmentContainerMountPointReadOnly] = *mp.ReadOnly
		} else {
			mountPoint[mkResourceVirtualEnvironmentContainerMountPointReadOnly] = false
		}

		if mp.Replicate != nil {
			mountPoint[mkResourceVirtualEnvironmentContainerMountPointReplicate] = *mp.Replicate
		} else {
			mountPoint[mkResourceVirtualEnvironmentContainerMountPointReplicate] = true
		}

		if mp.Shared != nil {
			mountPoint[mkResourceVirtualEnvironmentContainerMountPointShared] = *mp.Shared
		} else {
			mountPoint[mkResourceVirtualEnvironmentContainerMountPointShared] = false
		}

		if mp.DiskSize != nil {
			mountPoint[mkResourceVirtualEnvironmentContainerMountPointSize] = *mp.DiskSize
		} else {
			mountPoint[mkResourceVirtualEnvironmentContainerMountPointSize] = ""
		}

		mountPoint[mkResourceVirtualEnvironmentContainerMountPointVolume] = mp.Volume

		mountPointList = append(mountPointList, mountPoint)
	}

	currentMountPoint := d.Get(mkResourceVirtualEnvironmentContainerMountPoint).([]interface{})

	if len(clone) > 0 {
		if len(currentMountPoint) > 0 {
			err := d.Set(mkResourceVirtualEnvironmentContainerMountPoint, mountPointList)
			diags = append(diags, diag.FromErr(err)...)
		}
	} else if len(mountPointList) > 0 {
		err := d.Set(mkResourceVirtualEnvironmentContainerMountPoint, mountPointList)
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

	var networkInterfaceList []interface{}

	for _, nv := range networkInterfaceArray {
		if nv == nil {
			continue
		}

		//nolint:nestif
		if nv.IPv4Address != nil || nv.IPv4Gateway != nil || nv.IPv6Address != nil ||
			nv.IPv6Gateway != nil {
			ipConfig := map[string]interface{}{}

			if nv.IPv4Address != nil || nv.IPv4Gateway != nil {
				ip := map[string]interface{}{}

				if nv.IPv4Address != nil {
					ip[mkResourceVirtualEnvironmentContainerInitializationIPConfigIPv4Address] = *nv.IPv4Address
				} else {
					ip[mkResourceVirtualEnvironmentContainerInitializationIPConfigIPv4Address] = ""
				}

				if nv.IPv4Gateway != nil {
					ip[mkResourceVirtualEnvironmentContainerInitializationIPConfigIPv4Gateway] = *nv.IPv4Gateway
				} else {
					ip[mkResourceVirtualEnvironmentContainerInitializationIPConfigIPv4Gateway] = ""
				}

				ipConfig[mkResourceVirtualEnvironmentContainerInitializationIPConfigIPv4] = []interface{}{
					ip,
				}
			} else {
				ipConfig[mkResourceVirtualEnvironmentContainerInitializationIPConfigIPv4] = []interface{}{}
			}

			if nv.IPv6Address != nil || nv.IPv6Gateway != nil {
				ip := map[string]interface{}{}

				if nv.IPv6Address != nil {
					ip[mkResourceVirtualEnvironmentContainerInitializationIPConfigIPv6Address] = *nv.IPv6Address
				} else {
					ip[mkResourceVirtualEnvironmentContainerInitializationIPConfigIPv6Address] = ""
				}

				if nv.IPv6Gateway != nil {
					ip[mkResourceVirtualEnvironmentContainerInitializationIPConfigIPv6Gateway] = *nv.IPv6Gateway
				} else {
					ip[mkResourceVirtualEnvironmentContainerInitializationIPConfigIPv6Gateway] = ""
				}

				ipConfig[mkResourceVirtualEnvironmentContainerInitializationIPConfigIPv6] = []interface{}{
					ip,
				}
			} else {
				ipConfig[mkResourceVirtualEnvironmentContainerInitializationIPConfigIPv6] = []interface{}{}
			}

			ipConfigList = append(ipConfigList, ipConfig)
		}

		networkInterface := map[string]interface{}{}

		if nv.Bridge != nil {
			networkInterface[mkResourceVirtualEnvironmentContainerNetworkInterfaceBridge] = *nv.Bridge
		} else {
			networkInterface[mkResourceVirtualEnvironmentContainerNetworkInterfaceBridge] = ""
		}

		networkInterface[mkResourceVirtualEnvironmentContainerNetworkInterfaceEnabled] = true

		if nv.Firewall != nil {
			networkInterface[mkResourceVirtualEnvironmentContainerNetworkInterfaceFirewall] = *nv.Firewall
		} else {
			networkInterface[mkResourceVirtualEnvironmentContainerNetworkInterfaceFirewall] = false
		}

		if nv.MACAddress != nil {
			networkInterface[mkResourceVirtualEnvironmentContainerNetworkInterfaceMACAddress] = *nv.MACAddress
		} else {
			networkInterface[mkResourceVirtualEnvironmentContainerNetworkInterfaceMACAddress] = ""
		}

		networkInterface[mkResourceVirtualEnvironmentContainerNetworkInterfaceName] = nv.Name

		if nv.RateLimit != nil {
			networkInterface[mkResourceVirtualEnvironmentContainerNetworkInterfaceRateLimit] = *nv.RateLimit
		} else {
			networkInterface[mkResourceVirtualEnvironmentContainerNetworkInterfaceRateLimit] = 0
		}

		if nv.Tag != nil {
			networkInterface[mkResourceVirtualEnvironmentContainerNetworkInterfaceVLANID] = *nv.Tag
		} else {
			networkInterface[mkResourceVirtualEnvironmentContainerNetworkInterfaceVLANID] = 0
		}

		if nv.MTU != nil {
			networkInterface[mkResourceVirtualEnvironmentContainerNetworkInterfaceMTU] = *nv.MTU
		} else {
			networkInterface[mkResourceVirtualEnvironmentContainerNetworkInterfaceMTU] = 0
		}

		networkInterfaceList = append(networkInterfaceList, networkInterface)
	}

	initialization[mkResourceVirtualEnvironmentContainerInitializationIPConfig] = ipConfigList

	currentInitialization := d.Get(mkResourceVirtualEnvironmentContainerInitialization).([]interface{})

	if len(currentInitialization) > 0 {
		currentInitializationMap := currentInitialization[0].(map[string]interface{})

		initialization[mkResourceVirtualEnvironmentContainerInitializationUserAccount] = currentInitializationMap[mkResourceVirtualEnvironmentContainerInitializationUserAccount].([]interface{})
	}

	if len(clone) > 0 {
		if len(currentInitialization) > 0 {
			currentInitializationBlock := currentInitialization[0].(map[string]interface{})
			currentInitializationDNS := currentInitializationBlock[mkResourceVirtualEnvironmentContainerInitializationDNS].([]interface{})

			if len(currentInitializationDNS) == 0 {
				initialization[mkResourceVirtualEnvironmentContainerInitializationDNS] = []interface{}{}
			}

			currentInitializationIPConfig := currentInitializationBlock[mkResourceVirtualEnvironmentContainerInitializationIPConfig].([]interface{})

			if len(currentInitializationIPConfig) == 0 {
				initialization[mkResourceVirtualEnvironmentContainerInitializationIPConfig] = []interface{}{}
			}

			currentInitializationUserAccount := currentInitializationBlock[mkResourceVirtualEnvironmentContainerInitializationUserAccount].([]interface{})

			if len(currentInitializationUserAccount) == 0 {
				initialization[mkResourceVirtualEnvironmentContainerInitializationUserAccount] = []interface{}{}
			}

			if len(initialization) > 0 {
				e = d.Set(
					mkResourceVirtualEnvironmentContainerInitialization,
					[]interface{}{initialization},
				)
			} else {
				e = d.Set(mkResourceVirtualEnvironmentContainerInitialization, []interface{}{})
			}

			diags = append(diags, diag.FromErr(e)...)
		}

		currentNetworkInterface := d.Get(mkResourceVirtualEnvironmentContainerNetworkInterface).([]interface{})

		if len(currentNetworkInterface) > 0 {
			err := d.Set(
				mkResourceVirtualEnvironmentContainerNetworkInterface,
				networkInterfaceList,
			)
			diags = append(diags, diag.FromErr(err)...)
		}
	} else {
		if len(initialization) > 0 {
			e = d.Set(mkResourceVirtualEnvironmentContainerInitialization, []interface{}{initialization})
		} else {
			e = d.Set(mkResourceVirtualEnvironmentContainerInitialization, []interface{}{})
		}
		diags = append(diags, diag.FromErr(e)...)

		err := d.Set(mkResourceVirtualEnvironmentContainerNetworkInterface, networkInterfaceList)
		diags = append(diags, diag.FromErr(err)...)
	}

	// Compare the startup behavior to the one stored in the state.
	var startup map[string]interface{}

	//nolint:nestif
	if containerConfig.StartupBehavior != nil {
		startup = map[string]interface{}{}

		if containerConfig.StartupBehavior.Order != nil {
			startup[mkResourceVirtualEnvironmentContainerStartupOrder] = *containerConfig.StartupBehavior.Order
		} else {
			startup[mkResourceVirtualEnvironmentContainerStartupOrder] = dvResourceVirtualEnvironmentContainerStartupOrder
		}

		if containerConfig.StartupBehavior.Up != nil {
			startup[mkResourceVirtualEnvironmentContainerStartupUpDelay] = *containerConfig.StartupBehavior.Up
		} else {
			startup[mkResourceVirtualEnvironmentContainerStartupUpDelay] = dvResourceVirtualEnvironmentContainerStartupUpDelay
		}

		if containerConfig.StartupBehavior.Down != nil {
			startup[mkResourceVirtualEnvironmentContainerStartupDownDelay] = *containerConfig.StartupBehavior.Down
		} else {
			//nolint:lll
			startup[mkResourceVirtualEnvironmentContainerStartupDownDelay] = dvResourceVirtualEnvironmentContainerStartupDownDelay
		}
	}

	currentStartup := d.Get(mkResourceVirtualEnvironmentContainerStartup).([]interface{})

	//nolint:gocritic
	if len(clone) > 0 {
		if len(currentStartup) > 0 {
			err := d.Set(mkResourceVirtualEnvironmentContainerStartup, []interface{}{startup})
			diags = append(diags, diag.FromErr(err)...)
		}
	} else if len(startup) == 0 {
		err := d.Set(mkResourceVirtualEnvironmentContainerStartup, []interface{}{})
		diags = append(diags, diag.FromErr(err)...)
	} else if len(currentStartup) > 0 ||
		startup[mkResourceVirtualEnvironmentContainerStartupOrder] != mkResourceVirtualEnvironmentContainerStartupOrder ||
		startup[mkResourceVirtualEnvironmentContainerStartupUpDelay] != dvResourceVirtualEnvironmentContainerStartupUpDelay ||
		//nolint:lll
		startup[mkResourceVirtualEnvironmentContainerStartupDownDelay] != dvResourceVirtualEnvironmentContainerStartupDownDelay {
		err := d.Set(mkResourceVirtualEnvironmentContainerStartup, []interface{}{startup})
		diags = append(diags, diag.FromErr(err)...)
	}

	// Compare the operating system configuration to the one stored in the state.
	operatingSystem := map[string]interface{}{}

	if containerConfig.OSType != nil {
		operatingSystem[mkResourceVirtualEnvironmentContainerOperatingSystemType] = *containerConfig.OSType
	} else {
		// Default value of "ostype" is "" according to the API documentation.
		operatingSystem[mkResourceVirtualEnvironmentContainerOperatingSystemType] = ""
	}

	currentOperatingSystem := d.Get(mkResourceVirtualEnvironmentContainerOperatingSystem).([]interface{})

	if len(currentOperatingSystem) > 0 {
		currentOperatingSystemMap := currentOperatingSystem[0].(map[string]interface{})

		operatingSystem[mkResourceVirtualEnvironmentContainerOperatingSystemTemplateFileID] = currentOperatingSystemMap[mkResourceVirtualEnvironmentContainerOperatingSystemTemplateFileID]
	}

	if len(clone) > 0 {
		if len(currentOperatingSystem) > 0 {
			err := d.Set(
				mkResourceVirtualEnvironmentContainerOperatingSystem,
				[]interface{}{operatingSystem},
			)
			diags = append(diags, diag.FromErr(err)...)
		}
	} else if len(currentOperatingSystem) > 0 ||
		operatingSystem[mkResourceVirtualEnvironmentContainerOperatingSystemType] != dvResourceVirtualEnvironmentContainerOperatingSystemType {
		err := d.Set(mkResourceVirtualEnvironmentContainerOperatingSystem, []interface{}{operatingSystem})
		diags = append(diags, diag.FromErr(err)...)
	}

	currentTags := d.Get(mkResourceVirtualEnvironmentContainerTags).([]interface{})

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

		e = d.Set(mkResourceVirtualEnvironmentContainerTags, tags)
		diags = append(diags, diag.FromErr(e)...)
	}

	currentTemplate := d.Get(mkResourceVirtualEnvironmentContainerTemplate).(bool)

	//nolint:gosimple
	if len(clone) == 0 || currentTemplate != dvResourceVirtualEnvironmentContainerTemplate {
		if containerConfig.Template != nil {
			e = d.Set(
				mkResourceVirtualEnvironmentContainerTemplate,
				bool(*containerConfig.Template),
			)
		} else {
			e = d.Set(mkResourceVirtualEnvironmentContainerTemplate, false)
		}

		diags = append(diags, diag.FromErr(e)...)
	}

	// Determine the state of the container in order to update the "started" argument.
	status, e := containerAPI.GetContainerStatus(ctx)
	if e != nil {
		return diag.FromErr(e)
	}

	e = d.Set(mkResourceVirtualEnvironmentContainerStarted, status.Status == "running")
	diags = append(diags, diag.FromErr(e)...)

	return diags
}

func containerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(proxmoxtf.ProviderConfiguration)

	api, e := config.GetClient()
	if e != nil {
		return diag.FromErr(e)
	}

	nodeName := d.Get(mkResourceVirtualEnvironmentContainerNodeName).(string)

	vmID, e := strconv.Atoi(d.Id())
	if e != nil {
		return diag.FromErr(e)
	}

	containerAPI := api.Node(nodeName).Container(vmID)

	// Prepare the new request object.
	updateBody := containers.UpdateRequestBody{
		Delete: []string{},
	}

	rebootRequired := false
	resource := Container()

	// Retrieve the clone argument as the update logic varies for clones.
	clone := d.Get(mkResourceVirtualEnvironmentVMClone).([]interface{})

	// Prepare the new primitive values.
	description := d.Get(mkResourceVirtualEnvironmentContainerDescription).(string)
	updateBody.Description = &description

	template := types.CustomBool(d.Get(mkResourceVirtualEnvironmentContainerTemplate).(bool))

	if d.HasChange(mkResourceVirtualEnvironmentContainerTemplate) {
		updateBody.Template = &template
	}

	// Prepare the new console configuration.
	if d.HasChange(mkResourceVirtualEnvironmentContainerConsole) {
		consoleBlock, err := structure.GetSchemaBlock(
			resource,
			d,
			[]string{mkResourceVirtualEnvironmentContainerConsole},
			0,
			true,
		)
		if err != nil {
			return diag.FromErr(err)
		}

		consoleEnabled := types.CustomBool(
			consoleBlock[mkResourceVirtualEnvironmentContainerConsoleEnabled].(bool),
		)
		consoleMode := consoleBlock[mkResourceVirtualEnvironmentContainerConsoleMode].(string)
		consoleTTYCount := consoleBlock[mkResourceVirtualEnvironmentContainerConsoleTTYCount].(int)

		updateBody.ConsoleEnabled = &consoleEnabled
		updateBody.ConsoleMode = &consoleMode
		updateBody.TTY = &consoleTTYCount

		rebootRequired = true
	}

	// Prepare the new CPU configuration.
	if d.HasChange(mkResourceVirtualEnvironmentContainerCPU) {
		cpuBlock, err := structure.GetSchemaBlock(
			resource,
			d,
			[]string{mkResourceVirtualEnvironmentContainerCPU},
			0,
			true,
		)
		if err != nil {
			return diag.FromErr(err)
		}

		cpuArchitecture := cpuBlock[mkResourceVirtualEnvironmentContainerCPUArchitecture].(string)
		cpuCores := cpuBlock[mkResourceVirtualEnvironmentContainerCPUCores].(int)
		cpuUnits := cpuBlock[mkResourceVirtualEnvironmentContainerCPUUnits].(int)

		updateBody.CPUArchitecture = &cpuArchitecture
		updateBody.CPUCores = &cpuCores
		updateBody.CPUUnits = &cpuUnits

		rebootRequired = true
	}

	if d.HasChange(mkResourceVirtualEnvironmentContainerFeatures) {
		features, err := containerGetFeatures(resource, d)
		if err != nil {
			return diag.FromErr(err)
		}

		updateBody.Features = features
	}

	// Prepare the new initialization configuration.
	initialization := d.Get(mkResourceVirtualEnvironmentContainerInitialization).([]interface{})
	initializationDNSDomain := dvResourceVirtualEnvironmentContainerInitializationDNSDomain
	initializationDNSServer := dvResourceVirtualEnvironmentContainerInitializationDNSServer
	initializationHostname := dvResourceVirtualEnvironmentContainerInitializationHostname
	var initializationIPConfigIPv4Address []string
	var initializationIPConfigIPv4Gateway []string
	var initializationIPConfigIPv6Address []string
	var initializationIPConfigIPv6Gateway []string

	if len(initialization) > 0 {
		initializationBlock := initialization[0].(map[string]interface{})
		initializationDNS := initializationBlock[mkResourceVirtualEnvironmentContainerInitializationDNS].([]interface{})

		if len(initializationDNS) > 0 {
			initializationDNSBlock := initializationDNS[0].(map[string]interface{})
			initializationDNSDomain = initializationDNSBlock[mkResourceVirtualEnvironmentContainerInitializationDNSDomain].(string)

			servers := initializationDNSBlock[mkResourceVirtualEnvironmentContainerInitializationDNSServers].([]interface{})
			deprecatedServer := initializationDNSBlock[mkResourceVirtualEnvironmentContainerInitializationDNSServer].(string)

			if len(servers) > 0 {
				initializationDNSServer = strings.Join(ConvertToStringSlice(servers), " ")
			} else {
				initializationDNSServer = deprecatedServer
			}
		}

		initializationHostname = initializationBlock[mkResourceVirtualEnvironmentContainerInitializationHostname].(string)
		initializationIPConfig := initializationBlock[mkResourceVirtualEnvironmentContainerInitializationIPConfig].([]interface{})

		for _, c := range initializationIPConfig {
			if c == nil {
				continue
			}

			configBlock := c.(map[string]interface{})
			ipv4 := configBlock[mkResourceVirtualEnvironmentContainerInitializationIPConfigIPv4].([]interface{})

			if len(ipv4) > 0 {
				ipv4Block := ipv4[0].(map[string]interface{})

				initializationIPConfigIPv4Address = append(
					initializationIPConfigIPv4Address,
					ipv4Block[mkResourceVirtualEnvironmentContainerInitializationIPConfigIPv4Address].(string),
				)

				initializationIPConfigIPv4Gateway = append(
					initializationIPConfigIPv4Gateway,
					ipv4Block[mkResourceVirtualEnvironmentContainerInitializationIPConfigIPv4Gateway].(string),
				)
			} else {
				initializationIPConfigIPv4Address = append(initializationIPConfigIPv4Address, "")
				initializationIPConfigIPv4Gateway = append(initializationIPConfigIPv4Gateway, "")
			}

			ipv6 := configBlock[mkResourceVirtualEnvironmentContainerInitializationIPConfigIPv6].([]interface{})

			if len(ipv6) > 0 {
				ipv6Block := ipv6[0].(map[string]interface{})

				initializationIPConfigIPv6Address = append(
					initializationIPConfigIPv6Address,
					ipv6Block[mkResourceVirtualEnvironmentContainerInitializationIPConfigIPv6Address].(string),
				)

				initializationIPConfigIPv6Gateway = append(
					initializationIPConfigIPv6Gateway,
					ipv6Block[mkResourceVirtualEnvironmentContainerInitializationIPConfigIPv6Gateway].(string),
				)
			} else {
				initializationIPConfigIPv6Address = append(initializationIPConfigIPv6Address, "")
				initializationIPConfigIPv6Gateway = append(initializationIPConfigIPv6Gateway, "")
			}
		}
	}

	if d.HasChange(mkResourceVirtualEnvironmentContainerInitialization) {
		updateBody.DNSDomain = &initializationDNSDomain
		updateBody.DNSServer = &initializationDNSServer
		updateBody.Hostname = &initializationHostname

		rebootRequired = true
	}

	// Prepare the new memory configuration.
	if d.HasChange(mkResourceVirtualEnvironmentContainerMemory) {
		memoryBlock, err := structure.GetSchemaBlock(
			resource,
			d,
			[]string{mkResourceVirtualEnvironmentContainerMemory},
			0,
			true,
		)
		if err != nil {
			return diag.FromErr(err)
		}

		memoryDedicated := memoryBlock[mkResourceVirtualEnvironmentContainerMemoryDedicated].(int)
		memorySwap := memoryBlock[mkResourceVirtualEnvironmentContainerMemorySwap].(int)

		updateBody.DedicatedMemory = &memoryDedicated
		updateBody.Swap = &memorySwap

		rebootRequired = true
	}

	// Prepare the new mount point configuration.
	if d.HasChange(mkResourceVirtualEnvironmentContainerMountPoint) {
		mountPoint := d.Get(mkResourceVirtualEnvironmentContainerMountPoint).([]interface{})
		mountPointArray := make(
			containers.CustomMountPointArray,
			len(mountPoint),
		)

		for i, mp := range mountPoint {
			mountPointMap := mp.(map[string]interface{})
			mountPointObject := containers.CustomMountPoint{}

			acl := types.CustomBool(mountPointMap[mkResourceVirtualEnvironmentContainerMountPointACL].(bool))
			backup := types.CustomBool(mountPointMap[mkResourceVirtualEnvironmentContainerMountPointBackup].(bool))
			mountOptions := mountPointMap[mkResourceVirtualEnvironmentContainerMountPointMountOptions].([]interface{})
			path := mountPointMap[mkResourceVirtualEnvironmentContainerMountPointPath].(string)
			quota := types.CustomBool(mountPointMap[mkResourceVirtualEnvironmentContainerMountPointQuota].(bool))
			readOnly := types.CustomBool(mountPointMap[mkResourceVirtualEnvironmentContainerMountPointReadOnly].(bool))
			replicate := types.CustomBool(mountPointMap[mkResourceVirtualEnvironmentContainerMountPointReplicate].(bool))
			shared := types.CustomBool(mountPointMap[mkResourceVirtualEnvironmentContainerMountPointShared].(bool))
			volume := mountPointMap[mkResourceVirtualEnvironmentContainerMountPointVolume].(string)

			mountPointObject.ACL = &acl
			mountPointObject.Backup = &backup
			mountPointObject.MountPoint = path
			mountPointObject.Quota = &quota
			mountPointObject.ReadOnly = &readOnly
			mountPointObject.Replicate = &replicate
			mountPointObject.Shared = &shared
			mountPointObject.Volume = volume

			if len(mountOptions) > 0 {
				mountOptionsArray := make([]string, 0, len(mountPoint))

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
	networkInterface := d.Get(mkResourceVirtualEnvironmentContainerNetworkInterface).([]interface{})

	if len(networkInterface) == 0 && len(clone) > 0 {
		networkInterface, e = containerGetExistingNetworkInterface(ctx, containerAPI)
		if e != nil {
			return diag.FromErr(e)
		}
	}

	//nolint:nestif
	if d.HasChange(mkResourceVirtualEnvironmentContainerInitialization) ||
		d.HasChange(mkResourceVirtualEnvironmentContainerNetworkInterface) {
		networkInterfaceArray := make(
			containers.CustomNetworkInterfaceArray,
			len(networkInterface),
		)

		for ni, nv := range networkInterface {
			networkInterfaceMap := nv.(map[string]interface{})
			networkInterfaceObject := containers.CustomNetworkInterface{}

			bridge := networkInterfaceMap[mkResourceVirtualEnvironmentContainerNetworkInterfaceBridge].(string)
			enabled := networkInterfaceMap[mkResourceVirtualEnvironmentContainerNetworkInterfaceEnabled].(bool)
			firewall := types.CustomBool(
				networkInterfaceMap[mkResourceVirtualEnvironmentContainerNetworkInterfaceFirewall].(bool),
			)
			macAddress := networkInterfaceMap[mkResourceVirtualEnvironmentContainerNetworkInterfaceMACAddress].(string)
			name := networkInterfaceMap[mkResourceVirtualEnvironmentContainerNetworkInterfaceName].(string)
			rateLimit := networkInterfaceMap[mkResourceVirtualEnvironmentContainerNetworkInterfaceRateLimit].(float64)
			vlanID := networkInterfaceMap[mkResourceVirtualEnvironmentContainerNetworkInterfaceVLANID].(int)
			mtu := networkInterfaceMap[mkResourceVirtualEnvironmentContainerNetworkInterfaceMTU].(int)

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

		for i := 0; i < len(updateBody.NetworkInterfaces); i++ {
			if !updateBody.NetworkInterfaces[i].Enabled {
				updateBody.Delete = append(updateBody.Delete, fmt.Sprintf("net%d", i))
			}
		}

		for i := len(updateBody.NetworkInterfaces); i < maxResourceVirtualEnvironmentContainerNetworkInterfaces; i++ {
			updateBody.Delete = append(updateBody.Delete, fmt.Sprintf("net%d", i))
		}

		rebootRequired = true
	}

	if d.HasChange(mkResourceVirtualEnvironmentContainerStartup) {
		updateBody.StartupBehavior = containerGetStartupBehavior(d)
		if updateBody.StartupBehavior == nil {
			updateBody.Delete = append(updateBody.Delete, "startup")
		}
	}

	// Prepare the new operating system configuration.
	if d.HasChange(mkResourceVirtualEnvironmentContainerOperatingSystem) {
		operatingSystem, err := structure.GetSchemaBlock(
			resource,
			d,
			[]string{mkResourceVirtualEnvironmentContainerOperatingSystem},
			0,
			true,
		)
		if err != nil {
			return diag.FromErr(err)
		}

		operatingSystemType := operatingSystem[mkResourceVirtualEnvironmentContainerOperatingSystemType].(string)

		updateBody.OSType = &operatingSystemType

		rebootRequired = true
	}

	if d.HasChange(mkResourceVirtualEnvironmentContainerTags) {
		tagString := containerGetTagsString(d)
		updateBody.Tags = &tagString
	}

	// Update the configuration now that everything has been prepared.
	e = containerAPI.UpdateContainer(ctx, &updateBody)
	if e != nil {
		return diag.FromErr(e)
	}

	// Determine if the state of the container needs to be changed.
	started := d.Get(mkResourceVirtualEnvironmentContainerStarted).(bool)

	if d.HasChange(mkResourceVirtualEnvironmentContainerStarted) && !bool(template) {
		if started {
			e = containerAPI.StartContainer(ctx, 60)
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

			e = containerAPI.WaitForContainerStatus(ctx, "stopped", 300, 5)
			if e != nil {
				return diag.FromErr(e)
			}

			rebootRequired = false
		}
	}

	// As a final step in the update procedure, we might need to reboot the container.
	if !bool(template) && rebootRequired {
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
	config := m.(proxmoxtf.ProviderConfiguration)
	api, err := config.GetClient()
	if err != nil {
		return diag.FromErr(err)
	}

	nodeName := d.Get(mkResourceVirtualEnvironmentContainerNodeName).(string)
	vmID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	containerAPI := api.Node(nodeName).Container(vmID)

	// Shut down the container before deleting it.
	status, err := containerAPI.GetContainerStatus(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	if status.Status != "stopped" {
		forceStop := types.CustomBool(true)
		shutdownTimeout := 300

		err = containerAPI.ShutdownContainer(
			ctx,
			&containers.ShutdownRequestBody{
				ForceStop: &forceStop,
				Timeout:   &shutdownTimeout,
			},
		)
		if err != nil {
			return diag.FromErr(err)
		}

		err = containerAPI.WaitForContainerStatus(ctx, "stopped", 30, 5)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	err = containerAPI.DeleteContainer(ctx)

	if err != nil {
		if strings.Contains(err.Error(), "HTTP 404") {
			d.SetId("")

			return nil
		}
		return diag.FromErr(err)
	}

	// Wait for the state to become unavailable as that clearly indicates the destruction of the container.
	err = containerAPI.WaitForContainerStatus(ctx, "", 60, 2)
	if err == nil {
		return diag.Errorf("failed to delete container \"%d\"", vmID)
	}

	d.SetId("")

	return nil
}
