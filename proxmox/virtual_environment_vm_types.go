/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmox

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// CustomAgent handles QEMU agent parameters.
type CustomAgent struct {
	Enabled         CustomBool `json:"enabled,omitempty" url:"enabled,int"`
	TrimClonedDisks CustomBool `json:"fstrim_cloned_disks" url:"fstrim_cloned_disks,int"`
	Type            string     `json:"type" url:"type"`
}

// CustomAudioDevice handles QEMU audio parameters.
type CustomAudioDevice struct {
	Device string `json:"device" url:"device"`
	Driver string `json:"driver" url:"driver"`
}

// CustomCloudInitConfig handles QEMU cloud-init parameters.
type CustomCloudInitConfig struct {
	Files        *CustomCloudInitFiles     `json:"cicustom,omitempty" url:"cicustom,omitempty"`
	IPConfig     []CustomCloudInitIPConfig `json:"ipconfig,omitempty" url:"ipconfig,omitempty,numbered"`
	Nameserver   *string                   `json:"nameserver,omitempty" url:"nameserver,omitempty"`
	Password     *string                   `json:"cipassword,omitempty" url:"cipassword,omitempty"`
	SearchDomain *string                   `json:"searchdomain,omitempty" url:"searchdomain,omitempty"`
	SSHKeys      *CustomCloudInitSSHKeys   `json:"sshkeys,omitempty" url:"sshkeys,omitempty"`
	Type         *string                   `json:"citype,omitempty" url:"citype,omitempty"`
	Username     *string                   `json:"ciuser,omitempty" url:"ciuser,omitempty"`
}

// CustomCloudInitFiles handles QEMU cloud-init custom files parameters.
type CustomCloudInitFiles struct {
	MetaVolume    *string `json:"meta,omitempty" url:"meta,omitempty"`
	NetworkVolume *string `json:"network,omitempty" url:"network,omitempty"`
	UserVolume    *string `json:"user,omitempty" url:"user,omitempty"`
}

// CustomCloudInitIPConfig handles QEMU cloud-init IP configuration parameters.
type CustomCloudInitIPConfig struct {
	GatewayIPv4 *string `json:"gw,omitempty" url:"gw,omitempty"`
	GatewayIPv6 *string `json:"gw6,omitempty" url:"gw6,omitempty"`
	IPv4        *string `json:"ip,omitempty" url:"ip,omitempty"`
	IPv6        *string `json:"ip6,omitempty" url:"ip6,omitempty"`
}

// CustomCloudInitSSHKeys handles QEMU cloud-init SSH keys parameters.
type CustomCloudInitSSHKeys []string

// CustomEFIDisk handles QEMU EFI disk parameters.
type CustomEFIDisk struct {
	DiskSize   *int    `json:"size,omitempty" url:"size,omitempty"`
	FileVolume string  `json:"file" url:"file"`
	Format     *string `json:"format,omitempty" url:"format,omitempty"`
}

// CustomNetworkDevice handles QEMU network device parameters.
type CustomNetworkDevice struct {
	Model      string      `json:"model" url:"model"`
	Bridge     *string     `json:"bridge,omitempty" url:"bridge,omitempty"`
	Firewall   *CustomBool `json:"firewall,omitempty" url:"firewall,omitempty,int"`
	LinkDown   *CustomBool `json:"link_down,omitempty" url:"link_down,omitempty,int"`
	MACAddress *string     `json:"macaddr,omitempty" url:"macaddr,omitempty"`
	Queues     *int        `json:"queues,omitempty" url:"queues,omitempty"`
	RateLimit  *float64    `json:"rate,omitempty" url:"rate,omitempty"`
	Tag        *int        `json:"tag,omitempty" url:"tag,omitempty"`
	Trunks     []int       `json:"trunks,omitempty" url:"trunks,omitempty"`
}

// CustomNetworkDevices handles QEMU network device parameters.
type CustomNetworkDevices []CustomNetworkDevice

// CustomNUMADevice handles QEMU NUMA device parameters.
type CustomNUMADevice struct {
	CPUIDs        []string  `json:"cpus" url:"cpus,semicolon"`
	HostNodeNames *[]string `json:"hostnodes,omitempty" url:"hostnodes,omitempty,semicolon"`
	Memory        *float64  `json:"memory,omitempty" url:"memory,omitempty"`
	Policy        *string   `json:"policy,omitempty" url:"policy,omitempty"`
}

// CustomNUMADevices handles QEMU NUMA device parameters.
type CustomNUMADevices []CustomNUMADevice

// CustomPCIDevice handles QEMU host PCI device mapping parameters.
type CustomPCIDevice struct {
	DeviceIDs  []string    `json:"host" url:"host,semicolon"`
	DevicePath *string     `json:"mdev,omitempty" url:"mdev,omitempty"`
	PCIExpress *CustomBool `json:"pcie,omitempty" url:"pcie,omitempty,int"`
	ROMBAR     *CustomBool `json:"rombar,omitempty" url:"rombar,omitempty,int"`
	ROMFile    *string     `json:"romfile,omitempty" url:"romfile,omitempty"`
	XVGA       *CustomBool `json:"x-vga,omitempty" url:"x-vga,omitempty,int"`
}

// CustomPCIDevices handles QEMU host PCI device mapping parameters.
type CustomPCIDevices []CustomPCIDevice

// CustomSerialDevices handles QEMU serial device parameters.
type CustomSerialDevices []string

// CustomSharedMemory handles QEMU Inter-VM shared memory parameters.
type CustomSharedMemory struct {
	Name *string `json:"name,omitempty" url:"name,omitempty"`
	Size int     `json:"size" url:"size"`
}

// CustomSMBIOS handles QEMU SMBIOS parameters.
type CustomSMBIOS struct {
	Base64       *CustomBool `json:"base64,omitempty" url:"base64,omitempty"`
	Family       *string     `json:"family,omitempty" url:"family,omitempty"`
	Manufacturer *string     `json:"manufacturer,omitempty" url:"manufacturer,omitempty"`
	Product      *string     `json:"product,omitempty" url:"product,omitempty"`
	Serial       *string     `json:"serial,omitempty" url:"serial,omitempty"`
	SKU          *string     `json:"sku,omitempty" url:"sku,omitempty"`
	UUID         *string     `json:"uuid,omitempty" url:"uuid,omitempty"`
	Version      *string     `json:"version,omitempty" url:"version,omitempty"`
}

// CustomSpiceEnhancements handles QEMU spice enhancement parameters.
type CustomSpiceEnhancements struct {
	FolderSharing  *CustomBool `json:"foldersharing,omitempty" url:"foldersharing,omitempty"`
	VideoStreaming *string     `json:"videostreaming,omitempty" url:"videostreaming,omitempty"`
}

// CustomStartupOrder handles QEMU startup order parameters.
type CustomStartupOrder struct {
	Down  *int `json:"down,omitempty" url:"down,omitempty"`
	Order *int `json:"order,omitempty" url:"order,omitempty"`
	Up    *int `json:"up,omitempty" url:"up,omitempty"`
}

// CustomStorageDevice handles QEMU SATA device parameters.
type CustomStorageDevice struct {
	AIO           *string     `json:"aio,omitempty" url:"aio,omitempty"`
	BackupEnabled *CustomBool `json:"backup,omitempty" url:"backup,omitempty,int"`
	Enabled       bool        `json:"-" url:"-"`
	FileVolume    string      `json:"file" url:"file"`
}

// CustomStorageDevices handles QEMU SATA device parameters.
type CustomStorageDevices []CustomStorageDevice

// CustomUSBDevice handles QEMU USB device parameters.
type CustomUSBDevice struct {
	HostDevice string      `json:"host" url:"host"`
	USB3       *CustomBool `json:"usb3,omitempty" url:"usb3,omitempty,int"`
}

// CustomUSBDevices handles QEMU USB device parameters.
type CustomUSBDevices []CustomUSBDevice

// CustomVGADevice handles QEMU VGA device parameters.
type CustomVGADevice struct {
	Memory *int   `json:"memory,omitempty" url:"memory,omitempty"`
	Type   string `json:"type" url:"type"`
}

// CustomVirtualIODevice handles QEMU VirtIO device parameters.
type CustomVirtualIODevice struct {
	AIO           *string     `json:"aio,omitempty" url:"aio,omitempty"`
	BackupEnabled *CustomBool `json:"backup,omitempty" url:"backup,omitempty,int"`
	Enabled       bool        `json:"-" url:"-"`
	FileVolume    string      `json:"file" url:"file"`
}

// CustomVirtualIODevices handles QEMU VirtIO device parameters.
type CustomVirtualIODevices []CustomVirtualIODevice

// CustomWatchdogDevice handles QEMU watchdog device parameters.
type CustomWatchdogDevice struct {
	Action *string `json:"action,omitempty" url:"action,omitempty"`
	Model  string  `json:"model" url:"model"`
}

// VirtualEnvironmentVMCreateRequestBody contains the data for an virtual machine create request.
type VirtualEnvironmentVMCreateRequestBody struct {
	ACPI                 *CustomBool                  `json:"acpi,omitempty" url:"acpi,omitempty,int"`
	Agent                *CustomAgent                 `json:"agent,omitempty" url:"agent,omitempty"`
	AllowReboot          *CustomBool                  `json:"reboot,omitempty" url:"reboot,omitempty,int"`
	AudioDevice          *CustomAudioDevice           `json:"audio0,omitempty" url:"audio0,omitempty"`
	Autostart            *CustomBool                  `json:"autostart,omitempty" url:"autostart,omitempty,int"`
	BackupFile           *string                      `json:"archive,omitempty" url:"archive,omitempty"`
	BandwidthLimit       *int                         `json:"bwlimit,omitempty" url:"bwlimit,omitempty"`
	BIOS                 *string                      `json:"bios,omitempty" url:"bios,omitempty"`
	BootDisk             *string                      `json:"bootdisk,omitempty" url:"bootdisk,omitempty"`
	BootOrder            *string                      `json:"boot,omitempty" url:"boot,omitempty"`
	CDROM                *string                      `json:"cdrom,omitempty" url:"cdrom,omitempty"`
	CloudInitConfig      *CustomCloudInitConfig       `json:"cloudinit,omitempty" url:"cloudinit,omitempty"`
	CPUArchitecture      *string                      `json:"arch,omitempty" url:"arch,omitempty"`
	CPUCores             *int                         `json:"cores,omitempty" url:"cores,omitempty"`
	CPULimit             *int                         `json:"cpulimit,omitempty" url:"cpulimit,omitempty"`
	CPUSockets           *int                         `json:"sockets,omitempty" url:"sockets,omitempty"`
	CPUUnits             *int                         `json:"cpuunits,omitempty" url:"cpuunits,omitempty"`
	DedicatedMemory      *int                         `json:"memory,omitempty" url:"memory,omitempty"`
	DeletionProtection   *CustomBool                  `json:"protection,omitempty" url:"force,omitempty,int"`
	Description          *string                      `json:"description,omitempty" url:"description,omitempty"`
	EFIDisk              *CustomEFIDisk               `json:"efidisk0,omitempty" url:"efidisk0,omitempty"`
	FloatingMemory       *int                         `json:"balloon,omitempty" url:"balloon,omitempty"`
	FloatingMemoryShares *int                         `json:"shares,omitempty" url:"shares,omitempty"`
	Freeze               *CustomBool                  `json:"freeze,omitempty" url:"freeze,omitempty,int"`
	HookScript           *string                      `json:"hookscript,omitempty" url:"hookscript,omitempty"`
	Hotplug              CustomCommaSeparatedList     `json:"hotplug,omitempty" url:"hotplug,omitempty,comma"`
	Hugepages            *string                      `json:"hugepages,omitempty" url:"hugepages,omitempty"`
	IDEDevices           CustomStorageDevices         `json:"ide,omitempty" url:"ide,omitempty"`
	KeyboardLayout       *string                      `json:"keyboard,omitempty" url:"keyboard,omitempty"`
	KVMArguments         CustomLineBreakSeparatedList `json:"args,omitempty" url:"args,omitempty,space"`
	KVMEnabled           *CustomBool                  `json:"kvm,omitempty" url:"kvm,omitempty,int"`
	LocalTime            *CustomBool                  `json:"localtime,omitempty" url:"localtime,omitempty,int"`
	Lock                 *string                      `json:"lock,omitempty" url:"lock,omitempty"`
	MachineType          *string                      `json:"machine,omitempty" url:"machine,omitempty"`
	MigrateDowntime      *float64                     `json:"migrate_downtime,omitempty" url:"migrate_downtime,omitempty"`
	MigrateSpeed         *int                         `json:"migrate_speed,omitempty" url:"migrate_speed,omitempty"`
	Name                 *string                      `json:"name,omitempty" url:"name,omitempty"`
	NetworkDevices       CustomNetworkDevices         `json:"net,omitempty" url:"net,omitempty"`
	NUMADevices          CustomNUMADevices            `json:"numa_devices,omitempty" url:"numa,omitempty"`
	NUMAEnabled          *CustomBool                  `json:"numa,omitempty" url:"numa,omitempty,int"`
	OSType               *string                      `json:"ostype,omitempty" url:"ostype,omitempty"`
	Overwrite            *CustomBool                  `json:"force,omitempty" url:"force,omitempty,int"`
	PCIDevices           CustomPCIDevices             `json:"hostpci,omitempty" url:"hostpci,omitempty"`
	Revert               *string                      `json:"revert,omitempty" url:"revert,omitempty"`
	SATADevices          CustomStorageDevices         `json:"sata,omitempty" url:"sata,omitempty"`
	SCSIDevices          CustomStorageDevices         `json:"scsi,omitempty" url:"sata,omitempty"`
	SCSIHardware         *string                      `json:"scsihw,omitempty" url:"scsihw,omitempty"`
	SerialDevices        CustomSerialDevices          `json:"serial,omitempty" url:"serial,omitempty"`
	SharedMemory         *CustomSharedMemory          `json:"ivshmem,omitempty" url:"ivshmem,omitempty"`
	SkipLock             *CustomBool                  `json:"skiplock,omitempty" url:"skiplock,omitempty,int"`
	SMBIOS               *CustomSMBIOS                `json:"smbios1,omitempty" url:"smbios1,omitempty"`
	SpiceEnhancements    *CustomSpiceEnhancements     `json:"spice_enhancements,omitempty" url:"spice_enhancements,omitempty"`
	StartDate            *string                      `json:"startdate,omitempty" url:"startdate,omitempty"`
	StartOnBoot          *CustomBool                  `json:"onboot,omitempty" url:"onboot,omitempty,int"`
	StartupOrder         *CustomStartupOrder          `json:"startup,omitempty" url:"startup,omitempty"`
	TabletDeviceEnabled  *CustomBool                  `json:"tablet,omitempty" url:"tablet,omitempty,int"`
	Tags                 *string                      `json:"tags,omitempty" url:"tags,omitempty"`
	Template             *CustomBool                  `json:"template,omitempty" url:"template,omitempty,int"`
	TimeDriftFixEnabled  *CustomBool                  `json:"tdf,omitempty" url:"tdf,omitempty,int"`
	USBDevices           CustomUSBDevices             `json:"usb,omitempty" url:"usb,omitempty"`
	VGADevice            *CustomVGADevice             `json:"vga,omitempty" url:"vga,omitempty"`
	VirtualCPUCount      *int                         `json:"vcpus,omitempty" url:"vcpus,omitempty"`
	VirtualIODevices     CustomVirtualIODevices       `json:"virtio,omitempty" url:"virtio,omitempty"`
	VMGenerationID       *string                      `json:"vmgenid,omitempty" url:"vmgenid,omitempty"`
	VMStateDatastoreID   *string                      `json:"vmstatestorage,omitempty" url:"vmstatestorage,omitempty"`
	WatchdogDevice       *CustomWatchdogDevice        `json:"watchdog,omitempty" url:"watchdog,omitempty"`
}

// VirtualEnvironmentVMGetResponseBody contains the body from an virtual machine get response.
type VirtualEnvironmentVMGetResponseBody struct {
	Data *VirtualEnvironmentVMGetResponseData `json:"data,omitempty"`
}

// VirtualEnvironmentVMGetResponseData contains the data from an virtual machine get response.
type VirtualEnvironmentVMGetResponseData struct {
	ACPI                 *CustomBool                   `json:"acpi,omitempty"`
	Agent                *CustomAgent                  `json:"agent,omitempty"`
	AllowReboot          *CustomBool                   `json:"reboot,omitempty"`
	AudioDevice          *CustomAudioDevice            `json:"audio0,omitempty"`
	Autostart            *CustomBool                   `json:"autostart,omitempty"`
	BackupFile           *string                       `json:"archive,omitempty"`
	BandwidthLimit       *int                          `json:"bwlimit,omitempty"`
	BIOS                 *string                       `json:"bios,omitempty"`
	BootDisk             *string                       `json:"bootdisk,omitempty"`
	BootOrder            *string                       `json:"boot,omitempty"`
	CDROM                *string                       `json:"cdrom,omitempty"`
	CloudInitConfig      *CustomCloudInitConfig        `json:"cloudinit,omitempty"`
	CPUArchitecture      *string                       `json:"arch,omitempty"`
	CPUCores             *int                          `json:"cores,omitempty"`
	CPULimit             *int                          `json:"cpulimit,omitempty"`
	CPUSockets           *int                          `json:"sockets,omitempty"`
	CPUUnits             *int                          `json:"cpuunits,omitempty"`
	DedicatedMemory      *int                          `json:"memory,omitempty"`
	DeletionProtection   *CustomBool                   `json:"protection,omitempty"`
	Description          *string                       `json:"description,omitempty"`
	EFIDisk              *CustomEFIDisk                `json:"efidisk0,omitempty"`
	FloatingMemory       *int                          `json:"balloon,omitempty"`
	FloatingMemoryShares *int                          `json:"shares,omitempty"`
	Freeze               *CustomBool                   `json:"freeze,omitempty"`
	HookScript           *string                       `json:"hookscript,omitempty"`
	Hotplug              *CustomCommaSeparatedList     `json:"hotplug,omitempty"`
	Hugepages            *string                       `json:"hugepages,omitempty"`
	IDEDevices           *CustomStorageDevices         `json:"ide,omitempty"`
	KeyboardLayout       *string                       `json:"keyboard,omitempty"`
	KVMArguments         *CustomLineBreakSeparatedList `json:"args,omitempty"`
	KVMEnabled           *CustomBool                   `json:"kvm,omitempty"`
	LocalTime            *CustomBool                   `json:"localtime,omitempty"`
	Lock                 *string                       `json:"lock,omitempty"`
	MachineType          *string                       `json:"machine,omitempty"`
	MigrateDowntime      *float64                      `json:"migrate_downtime,omitempty"`
	MigrateSpeed         *int                          `json:"migrate_speed,omitempty"`
	Name                 *string                       `json:"name,omitempty"`
	NetworkDevices       *CustomNetworkDevices         `json:"net,omitempty"`
	NUMADevices          *CustomNUMADevices            `json:"numa_devices,omitempty"`
	NUMAEnabled          *CustomBool                   `json:"numa,omitempty"`
	OSType               *string                       `json:"ostype,omitempty"`
	Overwrite            *CustomBool                   `json:"force,omitempty"`
	PCIDevices           *CustomPCIDevices             `json:"hostpci,omitempty"`
	Revert               *string                       `json:"revert,omitempty"`
	SATADevices          *CustomStorageDevices         `json:"sata,omitempty"`
	SCSIDevices          *CustomStorageDevices         `json:"scsi,omitempty"`
	SCSIHardware         *string                       `json:"scsihw,omitempty"`
	SerialDevices        *CustomSerialDevices          `json:"serial,omitempty"`
	SharedMemory         *CustomSharedMemory           `json:"ivshmem,omitempty"`
	SkipLock             *CustomBool                   `json:"skiplock,omitempty"`
	SMBIOS               *CustomSMBIOS                 `json:"smbios1,omitempty"`
	SpiceEnhancements    *CustomSpiceEnhancements      `json:"spice_enhancements,omitempty"`
	StartDate            *string                       `json:"startdate,omitempty"`
	StartOnBoot          *CustomBool                   `json:"onboot,omitempty"`
	StartupOrder         *CustomStartupOrder           `json:"startup,omitempty"`
	TabletDeviceEnabled  *CustomBool                   `json:"tablet,omitempty"`
	Tags                 *string                       `json:"tags,omitempty"`
	Template             *CustomBool                   `json:"template,omitempty"`
	TimeDriftFixEnabled  *CustomBool                   `json:"tdf,omitempty"`
	USBDevices           *CustomUSBDevices             `json:"usb,omitempty"`
	VGADevice            *CustomVGADevice              `json:"vga,omitempty"`
	VirtualCPUCount      *int                          `json:"vcpus,omitempty"`
	VirtualIODevices     *CustomVirtualIODevices       `json:"virtio,omitempty"`
	VMGenerationID       *string                       `json:"vmgenid,omitempty"`
	VMStateDatastoreID   *string                       `json:"vmstatestorage,omitempty"`
	WatchdogDevice       *CustomWatchdogDevice         `json:"watchdog,omitempty"`
}

// VirtualEnvironmentVMListResponseBody contains the body from an virtual machine list response.
type VirtualEnvironmentVMListResponseBody struct {
	Data []*VirtualEnvironmentVMListResponseData `json:"data,omitempty"`
}

// VirtualEnvironmentVMListResponseData contains the data from an virtual machine list response.
type VirtualEnvironmentVMListResponseData struct {
	ACPI *CustomBool `json:"acpi,omitempty" url:"acpi,omitempty,int"`
}

// VirtualEnvironmentVMUpdateRequestBody contains the data for an virtual machine update request.
type VirtualEnvironmentVMUpdateRequestBody VirtualEnvironmentVMCreateRequestBody

// EncodeValues converts a CustomAgent struct to a URL vlaue.
func (r CustomAgent) EncodeValues(key string, v *url.Values) error {
	enabled := 0

	if r.Enabled {
		enabled = 1
	}

	trimClonedDisks := 0

	if r.TrimClonedDisks {
		trimClonedDisks = 1
	}

	v.Add(key, fmt.Sprintf("enabled=%d,fstrim_cloned_disks=%d,type=%s", enabled, trimClonedDisks, r.Type))

	return nil
}

// EncodeValues converts a CustomAudioDevice struct to a URL vlaue.
func (r CustomAudioDevice) EncodeValues(key string, v *url.Values) error {
	v.Add(key, fmt.Sprintf("device=%s,driver=%s", r.Device, r.Driver))

	return nil
}

// EncodeValues converts a CustomCloudInitConfig struct to multiple URL vlaues.
func (r CustomCloudInitConfig) EncodeValues(key string, v *url.Values) error {
	if r.Files != nil {
		volumes := []string{}

		if r.Files.MetaVolume != nil {
			volumes = append(volumes, fmt.Sprintf("meta=%s", *r.Files.MetaVolume))
		}

		if r.Files.NetworkVolume != nil {
			volumes = append(volumes, fmt.Sprintf("network=%s", *r.Files.NetworkVolume))
		}

		if r.Files.UserVolume != nil {
			volumes = append(volumes, fmt.Sprintf("user=%s", *r.Files.UserVolume))
		}

		if len(volumes) > 0 {
			v.Add("cicustom", strings.Join(volumes, ","))
		}
	}

	for i, c := range r.IPConfig {
		config := []string{}

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
		keys := []string{}

		for _, pk := range *r.SSHKeys {
			keys = append(keys, pk)
		}

		v.Add("sshkeys", strings.Join(keys, "\n"))
	}

	if r.Type != nil {
		v.Add("citype", *r.Type)
	}

	if r.Username != nil {
		v.Add("ciuser", *r.Username)
	}

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

	if r.DiskSize != nil {
		values = append(values, fmt.Sprintf("size=%d", *r.DiskSize))
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
		d.EncodeValues(fmt.Sprintf("%s%d", key, i), v)
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
		d.EncodeValues(fmt.Sprintf("%s%d", key, i), v)
	}

	return nil
}

// EncodeValues converts a CustomPCIDevice struct to a URL vlaue.
func (r CustomPCIDevice) EncodeValues(key string, v *url.Values) error {
	values := []string{
		fmt.Sprintf("host=%s", strings.Join(r.DeviceIDs, ";")),
	}

	if r.DevicePath != nil {
		values = append(values, fmt.Sprintf("mdev=%s", *r.DevicePath))
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
		d.EncodeValues(fmt.Sprintf("%s%d", key, i), v)
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
	values := []string{}

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
	values := []string{}

	if r.FolderSharing != nil {
		if *r.FolderSharing {
			values = append(values, fmt.Sprintf("foldersharing=1"))
		} else {
			values = append(values, fmt.Sprintf("foldersharing=0"))
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
	values := []string{}

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

// EncodeValues converts a CustomStorageDevice struct to a URL vlaue.
func (r CustomStorageDevice) EncodeValues(key string, v *url.Values) error {
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

// EncodeValues converts a CustomStorageDevices array to multiple URL values.
func (r CustomStorageDevices) EncodeValues(key string, v *url.Values) error {
	for i, d := range r {
		if d.Enabled {
			d.EncodeValues(fmt.Sprintf("%s%d", key, i), v)
		}
	}

	return nil
}

// EncodeValues converts a CustomUSBDevice struct to a URL vlaue.
func (r CustomUSBDevice) EncodeValues(key string, v *url.Values) error {
	values := []string{
		fmt.Sprintf("host=%s", r.HostDevice),
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
		d.EncodeValues(fmt.Sprintf("%s%d", key, i), v)
	}

	return nil
}

// EncodeValues converts a CustomVGADevice struct to a URL vlaue.
func (r CustomVGADevice) EncodeValues(key string, v *url.Values) error {
	values := []string{
		fmt.Sprintf("type=%s", r.Type),
	}

	if r.Memory != nil {
		values = append(values, fmt.Sprintf("memory=%d", *r.Memory))
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
			d.EncodeValues(fmt.Sprintf("%s%d", key, i), v)
		}
	}

	return nil
}

// EncodeValues converts a CustomWatchdogDevice struct to a URL vlaue.
func (r CustomWatchdogDevice) EncodeValues(key string, v *url.Values) error {
	values := []string{
		fmt.Sprintf("model=%s", r.Model),
	}

	if r.Action != nil {
		values = append(values, fmt.Sprintf("action=%s", *r.Action))
	}

	v.Add(key, strings.Join(values, ","))

	return nil
}
