/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package vms

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

// CustomAgent handles QEMU agent parameters.
type CustomAgent struct {
	Enabled         *types.CustomBool `json:"enabled,omitempty"   url:"enabled,int"`
	TrimClonedDisks *types.CustomBool `json:"fstrim_cloned_disks" url:"fstrim_cloned_disks,int"`
	Type            *string           `json:"type"                url:"type"`
}

// CustomAudioDevice handles QEMU audio parameters.
type CustomAudioDevice struct {
	Device  string  `json:"device" url:"device"`
	Driver  *string `json:"driver" url:"driver"`
	Enabled bool    `json:"-"      url:"-"`
}

// CustomAudioDevices handles QEMU audio device parameters.
type CustomAudioDevices []CustomAudioDevice

// CustomBoot handles QEMU boot parameters.
type CustomBoot struct {
	Order *[]string `json:"order,omitempty" url:"order,omitempty,semicolon"`
}

// CustomCloudInitConfig handles QEMU cloud-init parameters.
type CustomCloudInitConfig struct {
	Files        *CustomCloudInitFiles     `json:"cicustom,omitempty"     url:"cicustom,omitempty"`
	IPConfig     []CustomCloudInitIPConfig `json:"ipconfig,omitempty"     url:"ipconfig,omitempty,numbered"`
	Nameserver   *string                   `json:"nameserver,omitempty"   url:"nameserver,omitempty"`
	Password     *string                   `json:"cipassword,omitempty"   url:"cipassword,omitempty"`
	SearchDomain *string                   `json:"searchdomain,omitempty" url:"searchdomain,omitempty"`
	SSHKeys      *CustomCloudInitSSHKeys   `json:"sshkeys,omitempty"      url:"sshkeys,omitempty"`
	Type         *string                   `json:"citype,omitempty"       url:"citype,omitempty"`
	Username     *string                   `json:"ciuser,omitempty"       url:"ciuser,omitempty"`
}

// CustomCloudInitFiles handles QEMU cloud-init custom files parameters.
type CustomCloudInitFiles struct {
	MetaVolume    *string `json:"meta,omitempty"    url:"meta,omitempty"`
	NetworkVolume *string `json:"network,omitempty" url:"network,omitempty"`
	UserVolume    *string `json:"user,omitempty"    url:"user,omitempty"`
	VendorVolume  *string `json:"vendor,omitempty"  url:"vendor,omitempty"`
}

// CustomCloudInitIPConfig handles QEMU cloud-init IP configuration parameters.
type CustomCloudInitIPConfig struct {
	GatewayIPv4 *string `json:"gw,omitempty"  url:"gw,omitempty"`
	GatewayIPv6 *string `json:"gw6,omitempty" url:"gw6,omitempty"`
	IPv4        *string `json:"ip,omitempty"  url:"ip,omitempty"`
	IPv6        *string `json:"ip6,omitempty" url:"ip6,omitempty"`
}

// CustomCloudInitSSHKeys handles QEMU cloud-init SSH keys parameters.
type CustomCloudInitSSHKeys []string

// CustomCPUEmulation handles QEMU CPU emulation parameters.
type CustomCPUEmulation struct {
	Flags      *[]string         `json:"flags,omitempty"        url:"flags,omitempty,semicolon"`
	Hidden     *types.CustomBool `json:"hidden,omitempty"       url:"hidden,omitempty,int"`
	HVVendorID *string           `json:"hv-vendor-id,omitempty" url:"hv-vendor-id,omitempty"`
	Type       string            `json:"cputype,omitempty"      url:"cputype,omitempty"`
}

// CustomEFIDisk handles QEMU EFI disk parameters.
type CustomEFIDisk struct {
	FileVolume      string            `json:"file"                        url:"file"`
	Format          *string           `json:"format,omitempty"            url:"format,omitempty"`
	Type            *string           `json:"efitype,omitempty"           url:"efitype,omitempty"`
	PreEnrolledKeys *types.CustomBool `json:"pre-enrolled-keys,omitempty" url:"pre-enrolled-keys,omitempty,int"`
}

// CustomNetworkDevice handles QEMU network device parameters.
type CustomNetworkDevice struct {
	Model      string            `json:"model"               url:"model"`
	Bridge     *string           `json:"bridge,omitempty"    url:"bridge,omitempty"`
	Enabled    bool              `json:"-"                   url:"-"`
	Firewall   *types.CustomBool `json:"firewall,omitempty"  url:"firewall,omitempty,int"`
	LinkDown   *types.CustomBool `json:"link_down,omitempty" url:"link_down,omitempty,int"`
	MACAddress *string           `json:"macaddr,omitempty"   url:"macaddr,omitempty"`
	Queues     *int              `json:"queues,omitempty"    url:"queues,omitempty"`
	RateLimit  *float64          `json:"rate,omitempty"      url:"rate,omitempty"`
	Tag        *int              `json:"tag,omitempty"       url:"tag,omitempty"`
	MTU        *int              `json:"mtu,omitempty"       url:"mtu,omitempty"`
	Trunks     []int             `json:"trunks,omitempty"    url:"trunks,omitempty"`
}

// CustomNetworkDevices handles QEMU network device parameters.
type CustomNetworkDevices []CustomNetworkDevice

// CustomNUMADevice handles QEMU NUMA device parameters.
type CustomNUMADevice struct {
	CPUIDs        []string  `json:"cpus"                url:"cpus,semicolon"`
	HostNodeNames *[]string `json:"hostnodes,omitempty" url:"hostnodes,omitempty,semicolon"`
	Memory        *float64  `json:"memory,omitempty"    url:"memory,omitempty"`
	Policy        *string   `json:"policy,omitempty"    url:"policy,omitempty"`
}

// CustomNUMADevices handles QEMU NUMA device parameters.
type CustomNUMADevices []CustomNUMADevice

// CustomPCIDevice handles QEMU host PCI device mapping parameters.
type CustomPCIDevice struct {
	DeviceIDs  *[]string         `json:"host,omitempty"    url:"host,omitempty,semicolon"`
	Mapping    *string           `json:"mapping,omitempty" url:"mapping,omitempty"`
	MDev       *string           `json:"mdev,omitempty"    url:"mdev,omitempty"`
	PCIExpress *types.CustomBool `json:"pcie,omitempty"    url:"pcie,omitempty,int"`
	ROMBAR     *types.CustomBool `json:"rombar,omitempty"  url:"rombar,omitempty,int"`
	ROMFile    *string           `json:"romfile,omitempty" url:"romfile,omitempty"`
	XVGA       *types.CustomBool `json:"x-vga,omitempty"   url:"x-vga,omitempty,int"`
}

// CustomPCIDevices handles QEMU host PCI device mapping parameters.
type CustomPCIDevices []CustomPCIDevice

// CustomSerialDevices handles QEMU serial device parameters.
type CustomSerialDevices []string

// CustomSharedMemory handles QEMU Inter-VM shared memory parameters.
type CustomSharedMemory struct {
	Name *string `json:"name,omitempty" url:"name,omitempty"`
	Size int     `json:"size"           url:"size"`
}

// CustomSMBIOS handles QEMU SMBIOS parameters.
type CustomSMBIOS struct {
	Base64       *types.CustomBool `json:"base64,omitempty"       url:"base64,omitempty,int"`
	Family       *string           `json:"family,omitempty"       url:"family,omitempty"`
	Manufacturer *string           `json:"manufacturer,omitempty" url:"manufacturer,omitempty"`
	Product      *string           `json:"product,omitempty"      url:"product,omitempty"`
	Serial       *string           `json:"serial,omitempty"       url:"serial,omitempty"`
	SKU          *string           `json:"sku,omitempty"          url:"sku,omitempty"`
	UUID         *string           `json:"uuid,omitempty"         url:"uuid,omitempty"`
	Version      *string           `json:"version,omitempty"      url:"version,omitempty"`
}

// CustomSpiceEnhancements handles QEMU spice enhancement parameters.
type CustomSpiceEnhancements struct {
	FolderSharing  *types.CustomBool `json:"foldersharing,omitempty"  url:"foldersharing,omitempty"`
	VideoStreaming *string           `json:"videostreaming,omitempty" url:"videostreaming,omitempty"`
}

// CustomStartupOrder handles QEMU startup order parameters.
type CustomStartupOrder struct {
	Down  *int `json:"down,omitempty"  url:"down,omitempty"`
	Order *int `json:"order,omitempty" url:"order,omitempty"`
	Up    *int `json:"up,omitempty"    url:"up,omitempty"`
}

// CustomTPMState handles QEMU TPM state parameters.
type CustomTPMState struct {
	FileVolume string  `json:"file"              url:"file"`
	Version    *string `json:"version,omitempty" url:"version,omitempty"`
}

// CustomUSBDevice handles QEMU USB device parameters.
type CustomUSBDevice struct {
	HostDevice *string           `json:"host"              url:"host"`
	Mapping    *string           `json:"mapping,omitempty" url:"mapping,omitempty"`
	USB3       *types.CustomBool `json:"usb3,omitempty"    url:"usb3,omitempty,int"`
}

// CustomUSBDevices handles QEMU USB device parameters.
type CustomUSBDevices []CustomUSBDevice

// CustomVGADevice handles QEMU VGA device parameters.
type CustomVGADevice struct {
	Memory *int    `json:"memory,omitempty" url:"memory,omitempty"`
	Type   *string `json:"type,omitempty"   url:"type,omitempty"`
}

// CustomVirtualIODevice handles QEMU VirtIO device parameters.
type CustomVirtualIODevice struct {
	AIO           *string           `json:"aio,omitempty"    url:"aio,omitempty"`
	BackupEnabled *types.CustomBool `json:"backup,omitempty" url:"backup,omitempty,int"`
	Enabled       bool              `json:"-"                url:"-"`
	FileVolume    string            `json:"file"             url:"file"`
}

// CustomVirtualIODevices handles QEMU VirtIO device parameters.
type CustomVirtualIODevices []CustomVirtualIODevice

// CustomWatchdogDevice handles QEMU watchdog device parameters.
type CustomWatchdogDevice struct {
	Action *string `json:"action,omitempty" url:"action,omitempty"`
	Model  *string `json:"model"            url:"model"`
}

// CloneRequestBody contains the data for an virtual machine clone request.
type CloneRequestBody struct {
	BandwidthLimit      *int              `json:"bwlimit,omitempty"     url:"bwlimit,omitempty"`
	Description         *string           `json:"description,omitempty" url:"description,omitempty"`
	FullCopy            *types.CustomBool `json:"full,omitempty"        url:"full,omitempty,int"`
	Name                *string           `json:"name,omitempty"        url:"name,omitempty"`
	PoolID              *string           `json:"pool,omitempty"        url:"pool,omitempty"`
	SnapshotName        *string           `json:"snapname,omitempty"    url:"snapname,omitempty"`
	TargetNodeName      *string           `json:"target,omitempty"      url:"target,omitempty"`
	TargetStorage       *string           `json:"storage,omitempty"     url:"storage,omitempty"`
	TargetStorageFormat *string           `json:"format,omitempty"      url:"format,omitempty"`
	VMIDNew             int               `json:"newid"                 url:"newid"`
}

// CreateRequestBody contains the data for a virtual machine create request.
type CreateRequestBody struct {
	ACPI                 *types.CustomBool              `json:"acpi,omitempty"               url:"acpi,omitempty,int"`
	Agent                *CustomAgent                   `json:"agent,omitempty"              url:"agent,omitempty"`
	AllowReboot          *types.CustomBool              `json:"reboot,omitempty"             url:"reboot,omitempty,int"`
	AudioDevices         CustomAudioDevices             `json:"audio,omitempty"              url:"audio,omitempty"`
	Autostart            *types.CustomBool              `json:"autostart,omitempty"          url:"autostart,omitempty,int"`
	BackupFile           *string                        `json:"archive,omitempty"            url:"archive,omitempty"`
	BandwidthLimit       *int                           `json:"bwlimit,omitempty"            url:"bwlimit,omitempty"`
	BIOS                 *string                        `json:"bios,omitempty"               url:"bios,omitempty"`
	Boot                 *CustomBoot                    `json:"boot,omitempty"               url:"boot,omitempty"`
	CDROM                *string                        `json:"cdrom,omitempty"              url:"cdrom,omitempty"`
	CloudInitConfig      *CustomCloudInitConfig         `json:"cloudinit,omitempty"          url:"cloudinit,omitempty"`
	CPUArchitecture      *string                        `json:"arch,omitempty"               url:"arch,omitempty"`
	CPUCores             *int                           `json:"cores,omitempty"              url:"cores,omitempty"`
	CPUEmulation         *CustomCPUEmulation            `json:"cpu,omitempty"                url:"cpu,omitempty"`
	CPULimit             *int                           `json:"cpulimit,omitempty"           url:"cpulimit,omitempty"`
	CPUSockets           *int                           `json:"sockets,omitempty"            url:"sockets,omitempty"`
	CPUUnits             *int                           `json:"cpuunits,omitempty"           url:"cpuunits,omitempty"`
	DedicatedMemory      *int                           `json:"memory,omitempty"             url:"memory,omitempty"`
	Delete               []string                       `json:"delete,omitempty"             url:"delete,omitempty,comma"`
	DeletionProtection   *types.CustomBool              `json:"protection,omitempty"         url:"force,omitempty,int"`
	Description          *string                        `json:"description,omitempty"        url:"description,omitempty"`
	EFIDisk              *CustomEFIDisk                 `json:"efidisk0,omitempty"           url:"efidisk0,omitempty"`
	FloatingMemory       *int                           `json:"balloon,omitempty"            url:"balloon,omitempty"`
	FloatingMemoryShares *int                           `json:"shares,omitempty"             url:"shares,omitempty"`
	Freeze               *types.CustomBool              `json:"freeze,omitempty"             url:"freeze,omitempty,int"`
	HookScript           *string                        `json:"hookscript,omitempty"         url:"hookscript,omitempty"`
	Hotplug              types.CustomCommaSeparatedList `json:"hotplug,omitempty"            url:"hotplug,omitempty,comma"`
	Hugepages            *string                        `json:"hugepages,omitempty"          url:"hugepages,omitempty"`
	IDEDevices           CustomStorageDevices           `json:"ide,omitempty"                url:",omitempty"`
	KeyboardLayout       *string                        `json:"keyboard,omitempty"           url:"keyboard,omitempty"`
	KVMArguments         *string                        `json:"args,omitempty"               url:"args,omitempty,space"`
	KVMEnabled           *types.CustomBool              `json:"kvm,omitempty"                url:"kvm,omitempty,int"`
	LocalTime            *types.CustomBool              `json:"localtime,omitempty"          url:"localtime,omitempty,int"`
	Lock                 *string                        `json:"lock,omitempty"               url:"lock,omitempty"`
	Machine              *string                        `json:"machine,omitempty"            url:"machine,omitempty"`
	MigrateDowntime      *float64                       `json:"migrate_downtime,omitempty"   url:"migrate_downtime,omitempty"`
	MigrateSpeed         *int                           `json:"migrate_speed,omitempty"      url:"migrate_speed,omitempty"`
	Name                 *string                        `json:"name,omitempty"               url:"name,omitempty"`
	NetworkDevices       CustomNetworkDevices           `json:"net,omitempty"                url:"net,omitempty"`
	NUMADevices          CustomNUMADevices              `json:"numa_devices,omitempty"       url:"numa,omitempty"`
	NUMAEnabled          *types.CustomBool              `json:"numa,omitempty"               url:"numa,omitempty,int"`
	OSType               *string                        `json:"ostype,omitempty"             url:"ostype,omitempty"`
	Overwrite            *types.CustomBool              `json:"force,omitempty"              url:"force,omitempty,int"`
	PCIDevices           CustomPCIDevices               `json:"hostpci,omitempty"            url:"hostpci,omitempty"`
	PoolID               *string                        `json:"pool,omitempty"               url:"pool,omitempty"`
	Revert               *string                        `json:"revert,omitempty"             url:"revert,omitempty"`
	SATADevices          CustomStorageDevices           `json:"sata,omitempty"               url:"sata,omitempty"`
	SCSIDevices          CustomStorageDevices           `json:"scsi,omitempty"               url:"scsi,omitempty"`
	SCSIHardware         *string                        `json:"scsihw,omitempty"             url:"scsihw,omitempty"`
	SerialDevices        CustomSerialDevices            `json:"serial,omitempty"             url:"serial,omitempty"`
	SharedMemory         *CustomSharedMemory            `json:"ivshmem,omitempty"            url:"ivshmem,omitempty"`
	SkipLock             *types.CustomBool              `json:"skiplock,omitempty"           url:"skiplock,omitempty,int"`
	SMBIOS               *CustomSMBIOS                  `json:"smbios1,omitempty"            url:"smbios1,omitempty"`
	SpiceEnhancements    *CustomSpiceEnhancements       `json:"spice_enhancements,omitempty" url:"spice_enhancements,omitempty"`
	StartDate            *string                        `json:"startdate,omitempty"          url:"startdate,omitempty"`
	StartOnBoot          *types.CustomBool              `json:"onboot,omitempty"             url:"onboot,omitempty,int"`
	StartupOrder         *CustomStartupOrder            `json:"startup,omitempty"            url:"startup,omitempty"`
	TabletDeviceEnabled  *types.CustomBool              `json:"tablet,omitempty"             url:"tablet,omitempty,int"`
	Tags                 *string                        `json:"tags,omitempty"               url:"tags,omitempty"`
	Template             *types.CustomBool              `json:"template,omitempty"           url:"template,omitempty,int"`
	TimeDriftFixEnabled  *types.CustomBool              `json:"tdf,omitempty"                url:"tdf,omitempty,int"`
	TPMState             *CustomTPMState                `json:"tpmstate0,omitempty"          url:"tpmstate0,omitempty"`
	USBDevices           CustomUSBDevices               `json:"usb,omitempty"                url:"usb,omitempty"`
	VGADevice            *CustomVGADevice               `json:"vga,omitempty"                url:"vga,omitempty"`
	VirtualCPUCount      *int                           `json:"vcpus,omitempty"              url:"vcpus,omitempty"`
	VirtualIODevices     CustomStorageDevices           `json:"virtio,omitempty"             url:"virtio,omitempty"`
	VMGenerationID       *string                        `json:"vmgenid,omitempty"            url:"vmgenid,omitempty"`
	VMID                 *int                           `json:"vmid,omitempty"               url:"vmid,omitempty"`
	VMStateDatastoreID   *string                        `json:"vmstatestorage,omitempty"     url:"vmstatestorage,omitempty"`
	WatchdogDevice       *CustomWatchdogDevice          `json:"watchdog,omitempty"           url:"watchdog,omitempty"`
}

// CreateResponseBody contains the body from a create response.
type CreateResponseBody struct {
	Data *string `json:"data,omitempty"`
}

// GetQEMUNetworkInterfacesResponseBody contains the body from a QEMU get network interfaces response.
type GetQEMUNetworkInterfacesResponseBody struct {
	Data *GetQEMUNetworkInterfacesResponseData `json:"data,omitempty"`
}

// GetQEMUNetworkInterfacesResponseData contains the data from a QEMU get network interfaces response.
type GetQEMUNetworkInterfacesResponseData struct {
	Result *[]GetQEMUNetworkInterfacesResponseResult `json:"result,omitempty"`
}

// GetQEMUNetworkInterfacesResponseResult contains the result from a QEMU get network interfaces response.
type GetQEMUNetworkInterfacesResponseResult struct {
	MACAddress  string                                             `json:"hardware-address"`
	Name        string                                             `json:"name"`
	Statistics  *GetQEMUNetworkInterfacesResponseResultStatistics  `json:"statistics,omitempty"`
	IPAddresses *[]GetQEMUNetworkInterfacesResponseResultIPAddress `json:"ip-addresses,omitempty"`
}

// GetQEMUNetworkInterfacesResponseResultIPAddress contains the IP address from a QEMU get network interfaces response.
type GetQEMUNetworkInterfacesResponseResultIPAddress struct {
	Address string `json:"ip-address"`
	Prefix  int    `json:"prefix"`
	Type    string `json:"ip-address-type"`
}

// GetQEMUNetworkInterfacesResponseResultStatistics contains the statistics from a QEMU get network interfaces response.
type GetQEMUNetworkInterfacesResponseResultStatistics struct {
	RXBytes   int `json:"rx-bytes"`
	RXDropped int `json:"rx-dropped"`
	RXErrors  int `json:"rx-errs"`
	RXPackets int `json:"rx-packets"`
	TXBytes   int `json:"tx-bytes"`
	TXDropped int `json:"tx-dropped"`
	TXErrors  int `json:"tx-errs"`
	TXPackets int `json:"tx-packets"`
}

// GetResponseBody contains the body from a virtual machine get response.
type GetResponseBody struct {
	Data *GetResponseData `json:"data,omitempty"`
}

// GetResponseData contains the data from an virtual machine get response.
type GetResponseData struct {
	ACPI                 *types.CustomBool               `json:"acpi,omitempty"`
	Agent                *CustomAgent                    `json:"agent,omitempty"`
	AllowReboot          *types.CustomBool               `json:"reboot,omitempty"`
	AudioDevice          *CustomAudioDevice              `json:"audio0,omitempty"`
	Autostart            *types.CustomBool               `json:"autostart,omitempty"`
	BackupFile           *string                         `json:"archive,omitempty"`
	BandwidthLimit       *int                            `json:"bwlimit,omitempty"`
	BIOS                 *string                         `json:"bios,omitempty"`
	BootDisk             *string                         `json:"bootdisk,omitempty"`
	BootOrder            *string                         `json:"boot,omitempty"`
	CDROM                *string                         `json:"cdrom,omitempty"`
	CloudInitDNSDomain   *string                         `json:"searchdomain,omitempty"`
	CloudInitDNSServer   *string                         `json:"nameserver,omitempty"`
	CloudInitFiles       *CustomCloudInitFiles           `json:"cicustom,omitempty"`
	CloudInitPassword    *string                         `json:"cipassword,omitempty"`
	CloudInitSSHKeys     *CustomCloudInitSSHKeys         `json:"sshkeys,omitempty"`
	CloudInitType        *string                         `json:"citype,omitempty"`
	CloudInitUsername    *string                         `json:"ciuser,omitempty"`
	CPUArchitecture      *string                         `json:"arch,omitempty"`
	CPUCores             *int                            `json:"cores,omitempty"`
	CPUEmulation         *CustomCPUEmulation             `json:"cpu,omitempty"`
	CPULimit             *types.CustomInt                `json:"cpulimit,omitempty"`
	CPUSockets           *int                            `json:"sockets,omitempty"`
	CPUUnits             *int                            `json:"cpuunits,omitempty"`
	DedicatedMemory      *types.CustomInt64              `json:"memory,omitempty"`
	DeletionProtection   *types.CustomBool               `json:"protection,omitempty"`
	Description          *string                         `json:"description,omitempty"`
	EFIDisk              *CustomEFIDisk                  `json:"efidisk0,omitempty"`
	FloatingMemory       *types.CustomInt64              `json:"balloon,omitempty"`
	FloatingMemoryShares *int                            `json:"shares,omitempty"`
	Freeze               *types.CustomBool               `json:"freeze,omitempty"`
	HookScript           *string                         `json:"hookscript,omitempty"`
	Hotplug              *types.CustomCommaSeparatedList `json:"hotplug,omitempty"`
	Hugepages            *string                         `json:"hugepages,omitempty"`
	IDEDevice0           *CustomStorageDevice            `json:"ide0,omitempty"`
	IDEDevice1           *CustomStorageDevice            `json:"ide1,omitempty"`
	IDEDevice2           *CustomStorageDevice            `json:"ide2,omitempty"`
	IDEDevice3           *CustomStorageDevice            `json:"ide3,omitempty"`
	IPConfig0            *CustomCloudInitIPConfig        `json:"ipconfig0,omitempty"`
	IPConfig1            *CustomCloudInitIPConfig        `json:"ipconfig1,omitempty"`
	IPConfig2            *CustomCloudInitIPConfig        `json:"ipconfig2,omitempty"`
	IPConfig3            *CustomCloudInitIPConfig        `json:"ipconfig3,omitempty"`
	IPConfig4            *CustomCloudInitIPConfig        `json:"ipconfig4,omitempty"`
	IPConfig5            *CustomCloudInitIPConfig        `json:"ipconfig5,omitempty"`
	IPConfig6            *CustomCloudInitIPConfig        `json:"ipconfig6,omitempty"`
	IPConfig7            *CustomCloudInitIPConfig        `json:"ipconfig7,omitempty"`
	IPConfig8            *CustomCloudInitIPConfig        `json:"ipconfig8,omitempty"`
	IPConfig9            *CustomCloudInitIPConfig        `json:"ipconfig9,omitempty"`
	IPConfig10           *CustomCloudInitIPConfig        `json:"ipconfig10,omitempty"`
	IPConfig11           *CustomCloudInitIPConfig        `json:"ipconfig11,omitempty"`
	IPConfig12           *CustomCloudInitIPConfig        `json:"ipconfig12,omitempty"`
	IPConfig13           *CustomCloudInitIPConfig        `json:"ipconfig13,omitempty"`
	IPConfig14           *CustomCloudInitIPConfig        `json:"ipconfig14,omitempty"`
	IPConfig15           *CustomCloudInitIPConfig        `json:"ipconfig15,omitempty"`
	IPConfig16           *CustomCloudInitIPConfig        `json:"ipconfig16,omitempty"`
	IPConfig17           *CustomCloudInitIPConfig        `json:"ipconfig17,omitempty"`
	IPConfig18           *CustomCloudInitIPConfig        `json:"ipconfig18,omitempty"`
	IPConfig19           *CustomCloudInitIPConfig        `json:"ipconfig19,omitempty"`
	IPConfig20           *CustomCloudInitIPConfig        `json:"ipconfig20,omitempty"`
	IPConfig21           *CustomCloudInitIPConfig        `json:"ipconfig21,omitempty"`
	IPConfig22           *CustomCloudInitIPConfig        `json:"ipconfig22,omitempty"`
	IPConfig23           *CustomCloudInitIPConfig        `json:"ipconfig23,omitempty"`
	IPConfig24           *CustomCloudInitIPConfig        `json:"ipconfig24,omitempty"`
	IPConfig25           *CustomCloudInitIPConfig        `json:"ipconfig25,omitempty"`
	IPConfig26           *CustomCloudInitIPConfig        `json:"ipconfig26,omitempty"`
	IPConfig27           *CustomCloudInitIPConfig        `json:"ipconfig27,omitempty"`
	IPConfig28           *CustomCloudInitIPConfig        `json:"ipconfig28,omitempty"`
	IPConfig29           *CustomCloudInitIPConfig        `json:"ipconfig29,omitempty"`
	IPConfig30           *CustomCloudInitIPConfig        `json:"ipconfig30,omitempty"`
	IPConfig31           *CustomCloudInitIPConfig        `json:"ipconfig31,omitempty"`
	KeyboardLayout       *string                         `json:"keyboard,omitempty"`
	KVMArguments         *string                         `json:"args,omitempty"`
	KVMEnabled           *types.CustomBool               `json:"kvm,omitempty"`
	LocalTime            *types.CustomBool               `json:"localtime,omitempty"`
	Lock                 *string                         `json:"lock,omitempty"`
	Machine              *string                         `json:"machine,omitempty"`
	MigrateDowntime      *float64                        `json:"migrate_downtime,omitempty"`
	MigrateSpeed         *int                            `json:"migrate_speed,omitempty"`
	Name                 *string                         `json:"name,omitempty"`
	NetworkDevice0       *CustomNetworkDevice            `json:"net0,omitempty"`
	NetworkDevice1       *CustomNetworkDevice            `json:"net1,omitempty"`
	NetworkDevice2       *CustomNetworkDevice            `json:"net2,omitempty"`
	NetworkDevice3       *CustomNetworkDevice            `json:"net3,omitempty"`
	NetworkDevice4       *CustomNetworkDevice            `json:"net4,omitempty"`
	NetworkDevice5       *CustomNetworkDevice            `json:"net5,omitempty"`
	NetworkDevice6       *CustomNetworkDevice            `json:"net6,omitempty"`
	NetworkDevice7       *CustomNetworkDevice            `json:"net7,omitempty"`
	NetworkDevice8       *CustomNetworkDevice            `json:"net8,omitempty"`
	NetworkDevice9       *CustomNetworkDevice            `json:"net9,omitempty"`
	NetworkDevice10      *CustomNetworkDevice            `json:"net10,omitempty"`
	NetworkDevice11      *CustomNetworkDevice            `json:"net11,omitempty"`
	NetworkDevice12      *CustomNetworkDevice            `json:"net12,omitempty"`
	NetworkDevice13      *CustomNetworkDevice            `json:"net13,omitempty"`
	NetworkDevice14      *CustomNetworkDevice            `json:"net14,omitempty"`
	NetworkDevice15      *CustomNetworkDevice            `json:"net15,omitempty"`
	NetworkDevice16      *CustomNetworkDevice            `json:"net16,omitempty"`
	NetworkDevice17      *CustomNetworkDevice            `json:"net17,omitempty"`
	NetworkDevice18      *CustomNetworkDevice            `json:"net18,omitempty"`
	NetworkDevice19      *CustomNetworkDevice            `json:"net19,omitempty"`
	NetworkDevice20      *CustomNetworkDevice            `json:"net20,omitempty"`
	NetworkDevice21      *CustomNetworkDevice            `json:"net21,omitempty"`
	NetworkDevice22      *CustomNetworkDevice            `json:"net22,omitempty"`
	NetworkDevice23      *CustomNetworkDevice            `json:"net23,omitempty"`
	NetworkDevice24      *CustomNetworkDevice            `json:"net24,omitempty"`
	NetworkDevice25      *CustomNetworkDevice            `json:"net25,omitempty"`
	NetworkDevice26      *CustomNetworkDevice            `json:"net26,omitempty"`
	NetworkDevice27      *CustomNetworkDevice            `json:"net27,omitempty"`
	NetworkDevice28      *CustomNetworkDevice            `json:"net28,omitempty"`
	NetworkDevice29      *CustomNetworkDevice            `json:"net29,omitempty"`
	NetworkDevice30      *CustomNetworkDevice            `json:"net30,omitempty"`
	NetworkDevice31      *CustomNetworkDevice            `json:"net31,omitempty"`
	NUMADevices          *CustomNUMADevices              `json:"numa_devices,omitempty"`
	NUMAEnabled          *types.CustomBool               `json:"numa,omitempty"`
	OSType               *string                         `json:"ostype,omitempty"`
	Overwrite            *types.CustomBool               `json:"force,omitempty"`
	PCIDevice0           *CustomPCIDevice                `json:"hostpci0,omitempty"`
	PCIDevice1           *CustomPCIDevice                `json:"hostpci1,omitempty"`
	PCIDevice2           *CustomPCIDevice                `json:"hostpci2,omitempty"`
	PCIDevice3           *CustomPCIDevice                `json:"hostpci3,omitempty"`
	PoolID               *string                         `json:"pool,omitempty"               url:"pool,omitempty"`
	Revert               *string                         `json:"revert,omitempty"`
	SATADevice0          *CustomStorageDevice            `json:"sata0,omitempty"`
	SATADevice1          *CustomStorageDevice            `json:"sata1,omitempty"`
	SATADevice2          *CustomStorageDevice            `json:"sata2,omitempty"`
	SATADevice3          *CustomStorageDevice            `json:"sata3,omitempty"`
	SATADevice4          *CustomStorageDevice            `json:"sata4,omitempty"`
	SATADevice5          *CustomStorageDevice            `json:"sata5,omitempty"`
	SCSIDevice0          *CustomStorageDevice            `json:"scsi0,omitempty"`
	SCSIDevice1          *CustomStorageDevice            `json:"scsi1,omitempty"`
	SCSIDevice2          *CustomStorageDevice            `json:"scsi2,omitempty"`
	SCSIDevice3          *CustomStorageDevice            `json:"scsi3,omitempty"`
	SCSIDevice4          *CustomStorageDevice            `json:"scsi4,omitempty"`
	SCSIDevice5          *CustomStorageDevice            `json:"scsi5,omitempty"`
	SCSIDevice6          *CustomStorageDevice            `json:"scsi6,omitempty"`
	SCSIDevice7          *CustomStorageDevice            `json:"scsi7,omitempty"`
	SCSIDevice8          *CustomStorageDevice            `json:"scsi8,omitempty"`
	SCSIDevice9          *CustomStorageDevice            `json:"scsi9,omitempty"`
	SCSIDevice10         *CustomStorageDevice            `json:"scsi10,omitempty"`
	SCSIDevice11         *CustomStorageDevice            `json:"scsi11,omitempty"`
	SCSIDevice12         *CustomStorageDevice            `json:"scsi12,omitempty"`
	SCSIDevice13         *CustomStorageDevice            `json:"scsi13,omitempty"`
	SCSIHardware         *string                         `json:"scsihw,omitempty"`
	SerialDevice0        *string                         `json:"serial0,omitempty"`
	SerialDevice1        *string                         `json:"serial1,omitempty"`
	SerialDevice2        *string                         `json:"serial2,omitempty"`
	SerialDevice3        *string                         `json:"serial3,omitempty"`
	SharedMemory         *CustomSharedMemory             `json:"ivshmem,omitempty"`
	SkipLock             *types.CustomBool               `json:"skiplock,omitempty"`
	SMBIOS               *CustomSMBIOS                   `json:"smbios1,omitempty"`
	SpiceEnhancements    *CustomSpiceEnhancements        `json:"spice_enhancements,omitempty"`
	StartDate            *string                         `json:"startdate,omitempty"`
	StartOnBoot          *types.CustomBool               `json:"onboot,omitempty"`
	StartupOrder         *CustomStartupOrder             `json:"startup,omitempty"`
	TabletDeviceEnabled  *types.CustomBool               `json:"tablet,omitempty"`
	Tags                 *string                         `json:"tags,omitempty"`
	Template             *types.CustomBool               `json:"template,omitempty"`
	TimeDriftFixEnabled  *types.CustomBool               `json:"tdf,omitempty"`
	TPMState             *CustomTPMState                 `json:"tpmstate0,omitempty"`
	USBDevice0           *CustomUSBDevice                `json:"usb0,omitempty"`
	USBDevice1           *CustomUSBDevice                `json:"usb1,omitempty"`
	USBDevice2           *CustomUSBDevice                `json:"usb2,omitempty"`
	USBDevice3           *CustomUSBDevice                `json:"usb3,omitempty"`
	VGADevice            *CustomVGADevice                `json:"vga,omitempty"`
	VirtualCPUCount      *int                            `json:"vcpus,omitempty"`
	VirtualIODevice0     *CustomStorageDevice            `json:"virtio0,omitempty"`
	VirtualIODevice1     *CustomStorageDevice            `json:"virtio1,omitempty"`
	VirtualIODevice2     *CustomStorageDevice            `json:"virtio2,omitempty"`
	VirtualIODevice3     *CustomStorageDevice            `json:"virtio3,omitempty"`
	VirtualIODevice4     *CustomStorageDevice            `json:"virtio4,omitempty"`
	VirtualIODevice5     *CustomStorageDevice            `json:"virtio5,omitempty"`
	VirtualIODevice6     *CustomStorageDevice            `json:"virtio6,omitempty"`
	VirtualIODevice7     *CustomStorageDevice            `json:"virtio7,omitempty"`
	VirtualIODevice8     *CustomStorageDevice            `json:"virtio8,omitempty"`
	VirtualIODevice9     *CustomStorageDevice            `json:"virtio9,omitempty"`
	VirtualIODevice10    *CustomStorageDevice            `json:"virtio10,omitempty"`
	VirtualIODevice11    *CustomStorageDevice            `json:"virtio11,omitempty"`
	VirtualIODevice12    *CustomStorageDevice            `json:"virtio12,omitempty"`
	VirtualIODevice13    *CustomStorageDevice            `json:"virtio13,omitempty"`
	VirtualIODevice14    *CustomStorageDevice            `json:"virtio14,omitempty"`
	VirtualIODevice15    *CustomStorageDevice            `json:"virtio15,omitempty"`
	VMGenerationID       *string                         `json:"vmgenid,omitempty"`
	VMStateDatastoreID   *string                         `json:"vmstatestorage,omitempty"`
	WatchdogDevice       *CustomWatchdogDevice           `json:"watchdog,omitempty"`
}

// GetStatusResponseBody contains the body from a VM get status response.
type GetStatusResponseBody struct {
	Data *GetStatusResponseData `json:"data,omitempty"`
}

// GetStatusResponseData contains the data from a VM get status response.
type GetStatusResponseData struct {
	AgentEnabled     *types.CustomBool `json:"agent,omitempty"`
	CPUCount         *float64          `json:"cpus,omitempty"`
	Lock             *string           `json:"lock,omitempty"`
	MemoryAllocation *int64            `json:"maxmem,omitempty"`
	Name             *string           `json:"name,omitempty"`
	PID              *int              `json:"pid,omitempty"`
	QMPStatus        *string           `json:"qmpstatus,omitempty"`
	RootDiskSize     *int64            `json:"maxdisk,omitempty"`
	SpiceSupport     *types.CustomBool `json:"spice,omitempty"`
	Status           string            `json:"status,omitempty"`
	Tags             *string           `json:"tags,omitempty"`
	Uptime           *int              `json:"uptime,omitempty"`
	VMID             *int              `json:"vmid,omitempty"`
}

// ListResponseBody contains the body from a virtual machine list response.
type ListResponseBody struct {
	Data []*ListResponseData `json:"data,omitempty"`
}

// ListResponseData contains the data from an virtual machine list response.
type ListResponseData struct {
	Name *string `json:"name,omitempty"`
	Tags *string `json:"tags,omitempty"`
	VMID int     `json:"vmid,omitempty"`
}

// MigrateRequestBody contains the body for a VM migration request.
type MigrateRequestBody struct {
	OnlineMigration *types.CustomBool `json:"online,omitempty"           url:"online,omitempty,int"`
	TargetNode      string            `json:"target"                     url:"target"`
	TargetStorage   *string           `json:"targetstorage,omitempty"    url:"targetstorage,omitempty"`
	WithLocalDisks  *types.CustomBool `json:"with-local-disks,omitempty" url:"with-local-disks,omitempty,int"`
}

// MigrateResponseBody contains the body from a VM migrate response.
type MigrateResponseBody struct {
	Data *string `json:"data,omitempty"`
}

// MoveDiskRequestBody contains the body for a VM move disk request.
type MoveDiskRequestBody struct {
	BandwidthLimit      *int              `json:"bwlimit,omitempty" url:"bwlimit,omitempty"`
	DeleteOriginalDisk  *types.CustomBool `json:"delete,omitempty"  url:"delete,omitempty,int"`
	Digest              *string           `json:"digest,omitempty"  url:"digest,omitempty"`
	Disk                string            `json:"disk"              url:"disk"`
	TargetStorage       string            `json:"storage"           url:"storage"`
	TargetStorageFormat *string           `json:"format,omitempty"  url:"format,omitempty"`
}

// MoveDiskResponseBody contains the body from a VM move disk response.
type MoveDiskResponseBody struct {
	Data *string `json:"data,omitempty"`
}

// RebootRequestBody contains the body for a VM reboot request.
type RebootRequestBody struct {
	Timeout *int `json:"timeout,omitempty" url:"timeout,omitempty"`
}

// RebootResponseBody contains the body from a VM reboot response.
type RebootResponseBody struct {
	Data *string `json:"data,omitempty"`
}

// ResizeDiskRequestBody contains the body for a VM resize disk request.
type ResizeDiskRequestBody struct {
	Digest   *string           `json:"digest,omitempty"   url:"digest,omitempty"`
	Disk     string            `json:"disk"               url:"disk"`
	Size     types.DiskSize    `json:"size"               url:"size"`
	SkipLock *types.CustomBool `json:"skiplock,omitempty" url:"skiplock,omitempty,int"`
}

// ShutdownRequestBody contains the body for a VM shutdown request.
type ShutdownRequestBody struct {
	ForceStop  *types.CustomBool `json:"forceStop,omitempty"  url:"forceStop,omitempty,int"`
	KeepActive *types.CustomBool `json:"keepActive,omitempty" url:"keepActive,omitempty,int"`
	SkipLock   *types.CustomBool `json:"skipLock,omitempty"   url:"skipLock,omitempty,int"`
	Timeout    *int              `json:"timeout,omitempty"    url:"timeout,omitempty"`
}

// ShutdownResponseBody contains the body from a VM shutdown response.
type ShutdownResponseBody struct {
	Data *string `json:"data,omitempty"`
}

// StartRequestBody contains the body for a VM start request.
type StartRequestBody struct {
	ForceCPU         *string           `json:"force-cpu,omitempty"         url:"force-cpu,omitempty"`
	Machine          *string           `json:"machine,omitempty"           url:"machine,omitempty"`
	MigrateFrom      *string           `json:"migratefrom,omitempty"       url:"migratefrom,omitempty"`
	MigrationNetwork *string           `json:"migration_network,omitempty" url:"migration_network,omitempty"`
	MigrationType    *string           `json:"migration_type,omitempty"    url:"migration_type,omitempty"`
	SkipLock         *types.CustomBool `json:"skipLock,omitempty"          url:"skipLock,omitempty,int"`
	StateURI         *string           `json:"stateuri,omitempty"          url:"stateuri,omitempty"`
	TargetStorage    *string           `json:"targetstorage,omitempty"     url:"targetstorage,omitempty"`
	TimeoutSeconds   *int              `json:"timeout,omitempty"           url:"timeout,omitempty"`
}

// StartResponseBody contains the body from a VM start response.
type StartResponseBody struct {
	Data *string `json:"data,omitempty"`
}

// StopResponseBody contains the body from a VM stop response.
type StopResponseBody struct {
	Data *string `json:"data,omitempty"`
}

// UpdateAsyncResponseBody contains the body from a VM async update response.
type UpdateAsyncResponseBody struct {
	Data *string `json:"data,omitempty"`
}

// UpdateRequestBody contains the data for an virtual machine update request.
type UpdateRequestBody CreateRequestBody

// EncodeValues converts a CustomAgent struct to a URL vlaue.
func (r CustomAgent) EncodeValues(key string, v *url.Values) error {
	var values []string

	if r.Enabled != nil {
		if *r.Enabled {
			values = append(values, "enabled=1")
		} else {
			values = append(values, "enabled=0")
		}
	}

	if r.TrimClonedDisks != nil {
		if *r.TrimClonedDisks {
			values = append(values, "fstrim_cloned_disks=1")
		} else {
			values = append(values, "fstrim_cloned_disks=0")
		}
	}

	if r.Type != nil {
		values = append(values, fmt.Sprintf("type=%s", *r.Type))
	}

	if len(values) > 0 {
		v.Add(key, strings.Join(values, ","))
	}

	return nil
}

// EncodeValues converts a CustomAudioDevice struct to a URL vlaue.
func (r CustomAudioDevice) EncodeValues(key string, v *url.Values) error {
	values := []string{fmt.Sprintf("device=%s", r.Device)}

	if r.Driver != nil {
		values = append(values, fmt.Sprintf("driver=%s", *r.Driver))
	}

	v.Add(key, strings.Join(values, ","))

	return nil
}

// EncodeValues converts a CustomAudioDevices array to multiple URL values.
func (r CustomAudioDevices) EncodeValues(key string, v *url.Values) error {
	for i, d := range r {
		if d.Enabled {
			if err := d.EncodeValues(fmt.Sprintf("%s%d", key, i), v); err != nil {
				return fmt.Errorf("unable to encode audio device %d: %w", i, err)
			}
		}
	}

	return nil
}

// EncodeValues converts a CustomBoot struct to multiple URL values.
func (r CustomBoot) EncodeValues(key string, v *url.Values) error {
	if r.Order != nil && len(*r.Order) > 0 {
		v.Add(key, fmt.Sprintf("order=%s", strings.Join(*r.Order, ";")))
	}

	return nil
}

// EncodeValues converts a CustomCloudInitConfig struct to multiple URL values.
func (r CustomCloudInitConfig) EncodeValues(_ string, v *url.Values) error {
	//nolint:nestif
	if r.Files != nil {
		var volumes []string

		if r.Files.MetaVolume != nil {
			volumes = append(volumes, fmt.Sprintf("meta=%s", *r.Files.MetaVolume))
		}

		if r.Files.NetworkVolume != nil {
			volumes = append(volumes, fmt.Sprintf("network=%s", *r.Files.NetworkVolume))
		}

		if r.Files.UserVolume != nil {
			volumes = append(volumes, fmt.Sprintf("user=%s", *r.Files.UserVolume))
		}

		if r.Files.VendorVolume != nil {
			volumes = append(volumes, fmt.Sprintf("vendor=%s", *r.Files.VendorVolume))
		}

		if len(volumes) > 0 {
			v.Add("cicustom", strings.Join(volumes, ","))
		}
	}

	for i, c := range r.IPConfig {
		var config []string

		if c.GatewayIPv4 != nil {
			config = append(config, fmt.Sprintf("gw=%s", *c.GatewayIPv4))
		}

		if c.GatewayIPv6 != nil {
			config = append(config, fmt.Sprintf("gw6=%s", *c.GatewayIPv6))
		}

		if c.IPv4 != nil {
			config = append(config, fmt.Sprintf("ip=%s", *c.IPv4))
		}

		if c.IPv6 != nil {
			config = append(config, fmt.Sprintf("ip6=%s", *c.IPv6))
		}

		if len(config) > 0 {
			v.Add(fmt.Sprintf("ipconfig%d", i), strings.Join(config, ","))
		}
	}

	if r.Nameserver != nil {
		v.Add("nameserver", *r.Nameserver)
	}

	if r.Password != nil {
		v.Add("cipassword", *r.Password)
	}

	if r.SearchDomain != nil {
		v.Add("searchdomain", *r.SearchDomain)
	}

	if r.SSHKeys != nil {
		v.Add(
			"sshkeys",
			strings.ReplaceAll(url.QueryEscape(strings.Join(*r.SSHKeys, "\n")), "+", "%20"),
		)
	}

	if r.Type != nil {
		v.Add("citype", *r.Type)
	}

	if r.Username != nil {
		v.Add("ciuser", *r.Username)
	}

	return nil
}

// EncodeValues converts a CustomCPUEmulation struct to a URL vlaue.
func (r CustomCPUEmulation) EncodeValues(key string, v *url.Values) error {
	values := []string{
		fmt.Sprintf("cputype=%s", r.Type),
	}

	if r.Flags != nil && len(*r.Flags) > 0 {
		values = append(values, fmt.Sprintf("flags=%s", strings.Join(*r.Flags, ";")))
	}

	if r.Hidden != nil {
		if *r.Hidden {
			values = append(values, "hidden=1")
		} else {
			values = append(values, "hidden=0")
		}
	}

	if r.HVVendorID != nil {
		values = append(values, fmt.Sprintf("hv-vendor-id=%s", *r.HVVendorID))
	}

	v.Add(key, strings.Join(values, ","))

	return nil
}

// EncodeValues converts a CustomEFIDisk struct to a URL vlaue.
func (r CustomEFIDisk) EncodeValues(key string, v *url.Values) error {
	values := []string{
		fmt.Sprintf("file=%s", r.FileVolume),
	}

	if r.Format != nil {
		values = append(values, fmt.Sprintf("format=%s", *r.Format))
	}

	if r.Type != nil {
		values = append(values, fmt.Sprintf("efitype=%s", *r.Type))
	}

	if r.PreEnrolledKeys != nil {
		if *r.PreEnrolledKeys {
			values = append(values, "pre-enrolled-keys=1")
		} else {
			values = append(values, "pre-enrolled-keys=0")
		}
	}

	v.Add(key, strings.Join(values, ","))

	return nil
}

// EncodeValues converts a CustomNetworkDevice struct to a URL vlaue.
func (r CustomNetworkDevice) EncodeValues(key string, v *url.Values) error {
	values := []string{
		fmt.Sprintf("model=%s", r.Model),
	}

	if r.Bridge != nil {
		values = append(values, fmt.Sprintf("bridge=%s", *r.Bridge))
	}

	if r.Firewall != nil {
		if *r.Firewall {
			values = append(values, "firewall=1")
		} else {
			values = append(values, "firewall=0")
		}
	}

	if r.LinkDown != nil {
		if *r.LinkDown {
			values = append(values, "link_down=1")
		} else {
			values = append(values, "link_down=0")
		}
	}

	if r.MACAddress != nil {
		values = append(values, fmt.Sprintf("macaddr=%s", *r.MACAddress))
	}

	if r.Queues != nil {
		values = append(values, fmt.Sprintf("queues=%d", *r.Queues))
	}

	if r.RateLimit != nil {
		values = append(values, fmt.Sprintf("rate=%f", *r.RateLimit))
	}

	if r.Tag != nil {
		values = append(values, fmt.Sprintf("tag=%d", *r.Tag))
	}

	if r.MTU != nil {
		values = append(values, fmt.Sprintf("mtu=%d", *r.MTU))
	}

	if len(r.Trunks) > 0 {
		trunks := make([]string, len(r.Trunks))

		for i, v := range r.Trunks {
			trunks[i] = strconv.Itoa(v)
		}

		values = append(values, fmt.Sprintf("trunks=%s", strings.Join(trunks, ";")))
	}

	v.Add(key, strings.Join(values, ","))

	return nil
}

// EncodeValues converts a CustomNetworkDevices array to multiple URL values.
func (r CustomNetworkDevices) EncodeValues(key string, v *url.Values) error {
	for i, d := range r {
		if d.Enabled {
			if err := d.EncodeValues(fmt.Sprintf("%s%d", key, i), v); err != nil {
				return fmt.Errorf("failed to encode network device %d: %w", i, err)
			}
		}
	}

	return nil
}

// EncodeValues converts a CustomNUMADevice struct to a URL vlaue.
func (r CustomNUMADevice) EncodeValues(key string, v *url.Values) error {
	values := []string{
		fmt.Sprintf("cpus=%s", strings.Join(r.CPUIDs, ";")),
	}

	if r.HostNodeNames != nil {
		values = append(values, fmt.Sprintf("hostnodes=%s", strings.Join(*r.HostNodeNames, ";")))
	}

	if r.Memory != nil {
		values = append(values, fmt.Sprintf("memory=%f", *r.Memory))
	}

	if r.Policy != nil {
		values = append(values, fmt.Sprintf("policy=%s", *r.Policy))
	}

	v.Add(key, strings.Join(values, ","))

	return nil
}

// EncodeValues converts a CustomNUMADevices array to multiple URL values.
func (r CustomNUMADevices) EncodeValues(key string, v *url.Values) error {
	for i, d := range r {
		if err := d.EncodeValues(fmt.Sprintf("%s%d", key, i), v); err != nil {
			return fmt.Errorf("failed to encode NUMA device %d: %w", i, err)
		}
	}

	return nil
}

// EncodeValues converts a CustomPCIDevice struct to a URL vlaue.
func (r CustomPCIDevice) EncodeValues(key string, v *url.Values) error {
	values := []string{}

	if r.DeviceIDs == nil && r.Mapping == nil {
		return fmt.Errorf("either device ID or resource mapping must be set")
	}

	if r.DeviceIDs != nil {
		values = append(values, fmt.Sprintf("host=%s", strings.Join(*r.DeviceIDs, ";")))
	}

	if r.Mapping != nil {
		values = append(values, fmt.Sprintf("mapping=%s", *r.Mapping))
	}

	if r.MDev != nil {
		values = append(values, fmt.Sprintf("mdev=%s", *r.MDev))
	}

	if r.PCIExpress != nil {
		if *r.PCIExpress {
			values = append(values, "pcie=1")
		} else {
			values = append(values, "pcie=0")
		}
	}

	if r.ROMBAR != nil {
		if *r.ROMBAR {
			values = append(values, "rombar=1")
		} else {
			values = append(values, "rombar=0")
		}
	}

	if r.ROMFile != nil {
		values = append(values, fmt.Sprintf("romfile=%s", *r.ROMFile))
	}

	if r.XVGA != nil {
		if *r.XVGA {
			values = append(values, "x-vga=1")
		} else {
			values = append(values, "x-vga=0")
		}
	}

	v.Add(key, strings.Join(values, ","))

	return nil
}

// EncodeValues converts a CustomPCIDevices array to multiple URL values.
func (r CustomPCIDevices) EncodeValues(key string, v *url.Values) error {
	for i, d := range r {
		if err := d.EncodeValues(fmt.Sprintf("%s%d", key, i), v); err != nil {
			return fmt.Errorf("failed to encode PCI device %d: %w", i, err)
		}
	}

	return nil
}

// EncodeValues converts a CustomSerialDevices array to multiple URL values.
func (r CustomSerialDevices) EncodeValues(key string, v *url.Values) error {
	for i, d := range r {
		v.Add(fmt.Sprintf("%s%d", key, i), d)
	}

	return nil
}

// EncodeValues converts a CustomSharedMemory struct to a URL vlaue.
func (r CustomSharedMemory) EncodeValues(key string, v *url.Values) error {
	values := []string{
		fmt.Sprintf("size=%d", r.Size),
	}

	if r.Name != nil {
		values = append(values, fmt.Sprintf("name=%s", *r.Name))
	}

	v.Add(key, strings.Join(values, ","))

	return nil
}

// EncodeValues converts a CustomSMBIOS struct to a URL vlaue.
func (r CustomSMBIOS) EncodeValues(key string, v *url.Values) error {
	var values []string

	if r.Base64 != nil {
		if *r.Base64 {
			values = append(values, "base64=1")
		} else {
			values = append(values, "base64=0")
		}
	}

	if r.Family != nil {
		values = append(values, fmt.Sprintf("family=%s", *r.Family))
	}

	if r.Manufacturer != nil {
		values = append(values, fmt.Sprintf("manufacturer=%s", *r.Manufacturer))
	}

	if r.Product != nil {
		values = append(values, fmt.Sprintf("product=%s", *r.Product))
	}

	if r.Serial != nil {
		values = append(values, fmt.Sprintf("serial=%s", *r.Serial))
	}

	if r.SKU != nil {
		values = append(values, fmt.Sprintf("sku=%s", *r.SKU))
	}

	if r.UUID != nil {
		values = append(values, fmt.Sprintf("uuid=%s", *r.UUID))
	}

	if r.Version != nil {
		values = append(values, fmt.Sprintf("version=%s", *r.Version))
	}

	if len(values) > 0 {
		v.Add(key, strings.Join(values, ","))
	}

	return nil
}

// EncodeValues converts a CustomSpiceEnhancements struct to a URL vlaue.
func (r CustomSpiceEnhancements) EncodeValues(key string, v *url.Values) error {
	var values []string

	if r.FolderSharing != nil {
		if *r.FolderSharing {
			values = append(values, "foldersharing=1")
		} else {
			values = append(values, "foldersharing=0")
		}
	}

	if r.VideoStreaming != nil {
		values = append(values, fmt.Sprintf("videostreaming=%s", *r.VideoStreaming))
	}

	if len(values) > 0 {
		v.Add(key, strings.Join(values, ","))
	}

	return nil
}

// EncodeValues converts a CustomStartupOrder struct to a URL vlaue.
func (r CustomStartupOrder) EncodeValues(key string, v *url.Values) error {
	var values []string

	if r.Order != nil {
		values = append(values, fmt.Sprintf("order=%d", *r.Order))
	}

	if r.Up != nil {
		values = append(values, fmt.Sprintf("up=%d", *r.Up))
	}

	if r.Down != nil {
		values = append(values, fmt.Sprintf("down=%d", *r.Down))
	}

	if len(values) > 0 {
		v.Add(key, strings.Join(values, ","))
	}

	return nil
}

// EncodeValues converts a CustomTPMState struct to a URL vlaue.
func (r CustomTPMState) EncodeValues(key string, v *url.Values) error {
	values := []string{
		fmt.Sprintf("file=%s", r.FileVolume),
	}

	if r.Version != nil {
		values = append(values, fmt.Sprintf("version=%s", *r.Version))
	}

	v.Add(key, strings.Join(values, ","))

	return nil
}

// EncodeValues converts a CustomUSBDevice struct to a URL vlaue.
func (r CustomUSBDevice) EncodeValues(key string, v *url.Values) error {
	if r.HostDevice == nil && r.Mapping == nil {
		return fmt.Errorf("either device ID or resource mapping must be set")
	}

	values := []string{
		fmt.Sprintf("host=%s", *(r.HostDevice)),
	}

	if r.Mapping != nil {
		values = append(values, fmt.Sprintf("mapping=%s", *r.Mapping))
	}

	if r.USB3 != nil {
		if *r.USB3 {
			values = append(values, "usb3=1")
		} else {
			values = append(values, "usb3=0")
		}
	}

	v.Add(key, strings.Join(values, ","))

	return nil
}

// EncodeValues converts a CustomUSBDevices array to multiple URL values.
func (r CustomUSBDevices) EncodeValues(key string, v *url.Values) error {
	for i, d := range r {
		if err := d.EncodeValues(fmt.Sprintf("%s%d", key, i), v); err != nil {
			return fmt.Errorf("error encoding USB device %d: %w", i, err)
		}
	}

	return nil
}

// EncodeValues converts a CustomVGADevice struct to a URL vlaue.
func (r CustomVGADevice) EncodeValues(key string, v *url.Values) error {
	var values []string

	if r.Memory != nil {
		values = append(values, fmt.Sprintf("memory=%d", *r.Memory))
	}

	if r.Type != nil {
		values = append(values, fmt.Sprintf("type=%s", *r.Type))
	}

	v.Add(key, strings.Join(values, ","))

	return nil
}

// EncodeValues converts a CustomVirtualIODevice struct to a URL vlaue.
func (r CustomVirtualIODevice) EncodeValues(key string, v *url.Values) error {
	values := []string{
		fmt.Sprintf("file=%s", r.FileVolume),
	}

	if r.AIO != nil {
		values = append(values, fmt.Sprintf("aio=%s", *r.AIO))
	}

	if r.BackupEnabled != nil {
		if *r.BackupEnabled {
			values = append(values, "backup=1")
		} else {
			values = append(values, "backup=0")
		}
	}

	v.Add(key, strings.Join(values, ","))

	return nil
}

// EncodeValues converts a CustomVirtualIODevices array to multiple URL values.
func (r CustomVirtualIODevices) EncodeValues(key string, v *url.Values) error {
	for i, d := range r {
		if d.Enabled {
			if err := d.EncodeValues(fmt.Sprintf("%s%d", key, i), v); err != nil {
				return fmt.Errorf("error encoding virtual IO device %d: %w", i, err)
			}
		}
	}

	return nil
}

// EncodeValues converts a CustomWatchdogDevice struct to a URL vlaue.
func (r CustomWatchdogDevice) EncodeValues(key string, v *url.Values) error {
	values := []string{
		fmt.Sprintf("model=%+v", r.Model),
	}

	if r.Action != nil {
		values = append(values, fmt.Sprintf("action=%s", *r.Action))
	}

	v.Add(key, strings.Join(values, ","))

	return nil
}

// UnmarshalJSON converts a CustomAgent string to an object.
func (r *CustomAgent) UnmarshalJSON(b []byte) error {
	var s string

	if err := json.Unmarshal(b, &s); err != nil {
		return fmt.Errorf("error unmarshalling CustomAgent: %w", err)
	}

	pairs := strings.Split(s, ",")

	for _, p := range pairs {
		v := strings.Split(strings.TrimSpace(p), "=")

		if len(v) == 1 {
			enabled := types.CustomBool(v[0] == "1")
			r.Enabled = &enabled
		} else if len(v) == 2 {
			switch v[0] {
			case "enabled":
				enabled := types.CustomBool(v[1] == "1")
				r.Enabled = &enabled
			case "fstrim_cloned_disks":
				fstrim := types.CustomBool(v[1] == "1")
				r.TrimClonedDisks = &fstrim
			case "type":
				r.Type = &v[1]
			}
		}
	}

	return nil
}

// UnmarshalJSON converts a CustomAgent string to an object.
func (r *CustomAudioDevice) UnmarshalJSON(b []byte) error {
	var s string

	if err := json.Unmarshal(b, &s); err != nil {
		return fmt.Errorf("error unmarshalling CustomAudioDevice: %w", err)
	}

	pairs := strings.Split(s, ",")

	for _, p := range pairs {
		v := strings.Split(strings.TrimSpace(p), "=")

		if len(v) == 2 {
			switch v[0] {
			case "device":
				r.Device = v[1]
			case "driver":
				r.Driver = &v[1]
			}
		}
	}

	return nil
}

// UnmarshalJSON converts a CustomBoot string to an object.
func (r *CustomBoot) UnmarshalJSON(b []byte) error {
	var s string

	if err := json.Unmarshal(b, &s); err != nil {
		return fmt.Errorf("error unmarshalling CustomBoot: %w", err)
	}

	pairs := strings.Split(s, ",")

	for _, p := range pairs {
		v := strings.Split(strings.TrimSpace(p), "=")

		if len(v) == 2 {
			if v[0] == "order" {
				o := strings.Split(strings.TrimSpace(v[1]), ";")
				r.Order = &o
			}
		}
	}

	return nil
}

// UnmarshalJSON converts a CustomCloudInitFiles string to an object.
func (r *CustomCloudInitFiles) UnmarshalJSON(b []byte) error {
	var s string

	if err := json.Unmarshal(b, &s); err != nil {
		return fmt.Errorf("error unmarshalling CustomCloudInitFiles: %w", err)
	}

	pairs := strings.Split(s, ",")

	for _, p := range pairs {
		v := strings.Split(strings.TrimSpace(p), "=")

		if len(v) == 2 {
			switch v[0] {
			case "meta":
				r.MetaVolume = &v[1]
			case "network":
				r.NetworkVolume = &v[1]
			case "user":
				r.UserVolume = &v[1]
			case "vendor":
				r.VendorVolume = &v[1]
			}
		}
	}

	return nil
}

// UnmarshalJSON converts a CustomCloudInitIPConfig string to an object.
func (r *CustomCloudInitIPConfig) UnmarshalJSON(b []byte) error {
	var s string

	if err := json.Unmarshal(b, &s); err != nil {
		return fmt.Errorf("error unmarshalling CustomCloudInitIPConfig: %w", err)
	}

	pairs := strings.Split(s, ",")

	for _, p := range pairs {
		v := strings.Split(strings.TrimSpace(p), "=")

		if len(v) == 2 {
			switch v[0] {
			case "gw":
				r.GatewayIPv4 = &v[1]
			case "gw6":
				r.GatewayIPv6 = &v[1]
			case "ip":
				r.IPv4 = &v[1]
			case "ip6":
				r.IPv6 = &v[1]
			}
		}
	}

	return nil
}

// UnmarshalJSON converts a CustomCloudInitFiles string to an object.
func (r *CustomCloudInitSSHKeys) UnmarshalJSON(b []byte) error {
	var s string

	if err := json.Unmarshal(b, &s); err != nil {
		return fmt.Errorf("error unmarshalling CustomCloudInitSSHKeys: %w", err)
	}

	s, err := url.QueryUnescape(s)
	if err != nil {
		return fmt.Errorf("error unescaping CustomCloudInitSSHKeys: %w", err)
	}

	if s != "" {
		*r = strings.Split(strings.TrimSpace(s), "\n")
	} else {
		*r = []string{}
	}

	return nil
}

// UnmarshalJSON converts a CustomCPUEmulation string to an object.
func (r *CustomCPUEmulation) UnmarshalJSON(b []byte) error {
	var s string

	if err := json.Unmarshal(b, &s); err != nil {
		return fmt.Errorf("error unmarshalling CustomCPUEmulation: %w", err)
	}

	if s == "" {
		return errors.New("unexpected empty string")
	}

	pairs := strings.Split(s, ",")

	for _, p := range pairs {
		v := strings.Split(strings.TrimSpace(p), "=")

		if len(v) == 1 {
			r.Type = v[0]
		} else if len(v) == 2 {
			switch v[0] {
			case "cputype":
				r.Type = v[1]
			case "flags":
				if v[1] != "" {
					f := strings.Split(v[1], ";")
					r.Flags = &f
				} else {
					var f []string
					r.Flags = &f
				}
			case "hidden":
				bv := types.CustomBool(v[1] == "1")
				r.Hidden = &bv
			case "hv-vendor-id":
				r.HVVendorID = &v[1]
			}
		}
	}

	return nil
}

// UnmarshalJSON converts a CustomEFIDisk string to an object.
func (r *CustomEFIDisk) UnmarshalJSON(b []byte) error {
	var s string

	if err := json.Unmarshal(b, &s); err != nil {
		return fmt.Errorf("failed to unmarshal CustomEFIDisk: %w", err)
	}

	pairs := strings.Split(s, ",")

	for i, p := range pairs {
		v := strings.Split(strings.TrimSpace(p), "=")

		if len(v) == 1 && i == 0 {
			r.FileVolume = v[0]
		}

		if len(v) == 2 {
			switch v[0] {
			case "file":
				r.FileVolume = v[1]
			case "format":
				r.Format = &v[1]
			case "efitype":
				t := strings.ToLower(v[1])
				r.Type = &t
			case "pre-enrolled-keys":
				bv := types.CustomBool(v[1] == "1")
				r.PreEnrolledKeys = &bv
			}
		}
	}

	return nil
}

// UnmarshalJSON converts a CustomNetworkDevice string to an object.
func (r *CustomNetworkDevice) UnmarshalJSON(b []byte) error {
	var s string

	if err := json.Unmarshal(b, &s); err != nil {
		return fmt.Errorf("failed to unmarshal CustomNetworkDevice: %w", err)
	}

	pairs := strings.Split(s, ",")

	for _, p := range pairs {
		v := strings.Split(strings.TrimSpace(p), "=")

		//nolint:nestif
		if len(v) == 2 {
			switch v[0] {
			case "bridge":
				r.Bridge = &v[1]
			case "firewall":
				bv := types.CustomBool(v[1] == "1")
				r.Firewall = &bv
			case "link_down":
				bv := types.CustomBool(v[1] == "1")
				r.LinkDown = &bv
			case "macaddr":
				r.MACAddress = &v[1]
			case "model":
				r.Model = v[1]
			case "queues":
				iv, err := strconv.Atoi(v[1])
				if err != nil {
					return fmt.Errorf("failed to parse queues: %w", err)
				}

				r.Queues = &iv
			case "rate":
				fv, err := strconv.ParseFloat(v[1], 64)
				if err != nil {
					return fmt.Errorf("failed to parse rate: %w", err)
				}

				r.RateLimit = &fv

			case "mtu":
				iv, err := strconv.Atoi(v[1])
				if err != nil {
					return fmt.Errorf("failed to parse mtu: %w", err)
				}

				r.MTU = &iv

			case "tag":
				iv, err := strconv.Atoi(v[1])
				if err != nil {
					return fmt.Errorf("failed to parse tag: %w", err)
				}

				r.Tag = &iv
			case "trunks":
				trunks := strings.Split(v[1], ";")
				r.Trunks = make([]int, len(trunks))

				for i, trunk := range trunks {
					iv, err := strconv.Atoi(trunk)
					if err != nil {
						return fmt.Errorf("failed to parse trunk %d: %w", i, err)
					}

					r.Trunks[i] = iv
				}
			default:
				r.MACAddress = &v[1]
				r.Model = v[0]
			}
		}
	}

	r.Enabled = true

	return nil
}

// UnmarshalJSON converts a CustomPCIDevice string to an object.
func (r *CustomPCIDevice) UnmarshalJSON(b []byte) error {
	var s string

	if err := json.Unmarshal(b, &s); err != nil {
		return fmt.Errorf("failed to unmarshal CustomPCIDevice: %w", err)
	}

	pairs := strings.Split(s, ",")

	for _, p := range pairs {
		v := strings.Split(strings.TrimSpace(p), "=")
		if len(v) == 1 {
			dIds := strings.Split(v[0], ";")
			r.DeviceIDs = &dIds
		} else if len(v) == 2 {
			switch v[0] {
			case "host":
				dIds := strings.Split(v[1], ";")
				r.DeviceIDs = &dIds
			case "mapping":
				r.Mapping = &v[1]
			case "mdev":
				r.MDev = &v[1]
			case "pcie":
				bv := types.CustomBool(v[1] == "1")
				r.PCIExpress = &bv
			case "rombar":
				bv := types.CustomBool(v[1] == "1")
				r.ROMBAR = &bv
			case "romfile":
				r.ROMFile = &v[1]
			case "x-vga":
				bv := types.CustomBool(v[1] == "1")
				r.XVGA = &bv
			}
		}
	}

	return nil
}

// UnmarshalJSON converts a CustomTPMState string to an object.
func (r *CustomTPMState) UnmarshalJSON(b []byte) error {
	var s string

	if err := json.Unmarshal(b, &s); err != nil {
		return fmt.Errorf("failed to unmarshal CustomTPMState: %w", err)
	}

	pairs := strings.Split(s, ",")

	for _, p := range pairs {
		v := strings.Split(strings.TrimSpace(p), "=")
		if len(v) == 1 {
			r.FileVolume = v[0]
		} else if len(v) == 2 {
			switch v[0] {
			case "file":
				r.FileVolume = v[1]
			case "version":
				r.Version = &v[1]
			}
		}
	}

	return nil
}

// UnmarshalJSON converts a CustomUSBDevice string to an object.
func (r *CustomUSBDevice) UnmarshalJSON(b []byte) error {
	var s string

	if err := json.Unmarshal(b, &s); err != nil {
		return fmt.Errorf("failed to unmarshal CustomUSBDevice: %w", err)
	}

	pairs := strings.Split(s, ",")

	for _, p := range pairs {
		v := strings.Split(strings.TrimSpace(p), "=")
		if len(v) == 1 {
			r.HostDevice = &v[1]
		} else if len(v) == 2 {
			switch v[0] {
			case "host":
				r.HostDevice = &v[1]
			case "mapping":
				r.Mapping = &v[1]
			case "usb3":
				bv := types.CustomBool(v[1] == "1")
				r.USB3 = &bv
			}
		}
	}

	return nil
}

// UnmarshalJSON converts a CustomSharedMemory string to an object.
func (r *CustomSharedMemory) UnmarshalJSON(b []byte) error {
	var s string

	if err := json.Unmarshal(b, &s); err != nil {
		return fmt.Errorf("failed to unmarshal CustomSharedMemory: %w", err)
	}

	pairs := strings.Split(s, ",")

	for _, p := range pairs {
		v := strings.Split(strings.TrimSpace(p), "=")

		if len(v) == 2 {
			switch v[0] {
			case "name":
				r.Name = &v[1]
			case "size":
				var err error

				r.Size, err = strconv.Atoi(v[1])
				if err != nil {
					return fmt.Errorf("failed to parse shared memory size: %w", err)
				}
			}
		}
	}

	return nil
}

// UnmarshalJSON converts a CustomSMBIOS string to an object.
func (r *CustomSMBIOS) UnmarshalJSON(b []byte) error {
	var s string

	if err := json.Unmarshal(b, &s); err != nil {
		return fmt.Errorf("failed to unmarshal CustomSMBIOS: %w", err)
	}

	pairs := strings.Split(s, ",")

	for _, p := range pairs {
		v := strings.SplitN(strings.TrimSpace(p), "=", 2)

		if len(v) == 2 {
			switch v[0] {
			case "base64":
				base64 := types.CustomBool(v[1] == "1")
				r.Base64 = &base64
			case "family":
				r.Family = &v[1]
			case "manufacturer":
				r.Manufacturer = &v[1]
			case "product":
				r.Product = &v[1]
			case "serial":
				r.Serial = &v[1]
			case "sku":
				r.SKU = &v[1]
			case "uuid":
				r.UUID = &v[1]
			case "version":
				r.Version = &v[1]
			}
		}
	}

	return nil
}

// UnmarshalJSON converts a CustomStartupOrder string to an object.
func (r *CustomStartupOrder) UnmarshalJSON(b []byte) error {
	var s string

	if err := json.Unmarshal(b, &s); err != nil {
		return fmt.Errorf("failed to unmarshal CustomStartupOrder: %w", err)
	}

	pairs := strings.Split(s, ",")

	for _, p := range pairs {
		v := strings.Split(strings.TrimSpace(p), "=")

		if len(v) == 2 {
			switch v[0] {
			case "order":
				order, err := strconv.Atoi(v[1])
				if err != nil {
					return fmt.Errorf("failed to parse int: %w", err)
				}

				r.Order = &order
			case "up":
				up, err := strconv.Atoi(v[1])
				if err != nil {
					return fmt.Errorf("failed to parse int: %w", err)
				}

				r.Up = &up
			case "down":
				down, err := strconv.Atoi(v[1])
				if err != nil {
					return fmt.Errorf("failed to parse int: %w", err)
				}

				r.Down = &down
			}
		}
	}

	return nil
}

// UnmarshalJSON converts a CustomStorageDevice string to an object.
func (d *CustomStorageDevice) UnmarshalJSON(b []byte) error {
	var s string

	if err := json.Unmarshal(b, &s); err != nil {
		return fmt.Errorf("failed to unmarshal CustomStorageDevice: %w", err)
	}

	pairs := strings.Split(s, ",")

	for _, p := range pairs {
		v := strings.Split(strings.TrimSpace(p), "=")

		//nolint:nestif
		if len(v) == 1 {
			d.FileVolume = v[0]

			ext := filepath.Ext(v[0])
			if ext != "" {
				format := string([]byte(ext)[1:])
				d.Format = &format
			}
		} else if len(v) == 2 {
			switch v[0] {
			case "aio":
				d.AIO = &v[1]

			case "backup":
				bv := types.CustomBool(v[1] == "1")
				d.BackupEnabled = &bv

			case "file":
				d.FileVolume = v[1]

			case "mbps_rd":
				iv, err := strconv.Atoi(v[1])
				if err != nil {
					return fmt.Errorf("failed to convert mbps_rd to int: %w", err)
				}

				d.MaxReadSpeedMbps = &iv
			case "mbps_rd_max":
				iv, err := strconv.Atoi(v[1])
				if err != nil {
					return fmt.Errorf("failed to convert mbps_rd_max to int: %w", err)
				}

				d.BurstableReadSpeedMbps = &iv
			case "mbps_wr":
				iv, err := strconv.Atoi(v[1])
				if err != nil {
					return fmt.Errorf("failed to convert mbps_wr to int: %w", err)
				}

				d.MaxWriteSpeedMbps = &iv
			case "mbps_wr_max":
				iv, err := strconv.Atoi(v[1])
				if err != nil {
					return fmt.Errorf("failed to convert mbps_wr_max to int: %w", err)
				}

				d.BurstableWriteSpeedMbps = &iv
			case "media":
				d.Media = &v[1]

			case "size":
				d.Size = new(types.DiskSize)
				err := d.Size.UnmarshalJSON([]byte(v[1]))
				if err != nil {
					return fmt.Errorf("failed to unmarshal disk size: %w", err)
				}

			case "format":
				d.Format = &v[1]

			case "iothread":
				bv := types.CustomBool(v[1] == "1")
				d.IOThread = &bv

			case "ssd":
				bv := types.CustomBool(v[1] == "1")
				d.SSD = &bv

			case "discard":
				d.Discard = &v[1]

			case "cache":
				d.Cache = &v[1]
			}
		}
	}

	d.Enabled = true

	return nil
}

// UnmarshalJSON converts a CustomVGADevice string to an object.
func (r *CustomVGADevice) UnmarshalJSON(b []byte) error {
	var s string

	if err := json.Unmarshal(b, &s); err != nil {
		return fmt.Errorf("failed to unmarshal CustomVGADevice: %w", err)
	}

	if s == "" {
		return nil
	}

	pairs := strings.Split(s, ",")

	for _, p := range pairs {
		v := strings.Split(strings.TrimSpace(p), "=")

		if len(v) == 1 {
			r.Type = &v[0]
		} else if len(v) == 2 {
			switch v[0] {
			case "memory":
				m, err := strconv.Atoi(v[1])
				if err != nil {
					return fmt.Errorf("failed to convert memory to int: %w", err)
				}

				r.Memory = &m
			case "type":
				r.Type = &v[1]
			}
		}
	}

	return nil
}

// UnmarshalJSON converts a CustomWatchdogDevice string to an object.
func (r *CustomWatchdogDevice) UnmarshalJSON(b []byte) error {
	var s string

	if err := json.Unmarshal(b, &s); err != nil {
		return fmt.Errorf("failed to unmarshal CustomWatchdogDevice: %w", err)
	}

	if s == "" {
		return nil
	}

	pairs := strings.Split(s, ",")

	for _, p := range pairs {
		v := strings.Split(strings.TrimSpace(p), "=")

		if len(v) == 1 {
			r.Model = &v[0]
		} else if len(v) == 2 {
			switch v[0] {
			case "action":
				r.Action = &v[1]
			case "model":
				r.Model = &v[1]
			}
		}
	}

	return nil
}
