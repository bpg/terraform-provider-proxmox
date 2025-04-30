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
	"reflect"
	"regexp"
	"strings"

	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

var (
	//nolint:gochecknoglobals
	// regexStorageInterface is a regex pattern for matching storage interface names.
	regexStorageInterface = func(prefix string) *regexp.Regexp {
		return regexp.MustCompile(`^` + prefix + `\d+$`)
	}
	// regexPCIDevice is a regex pattern for matching PCI device names.
	regexPCIDevice = regexp.MustCompile(`^hostpci\d+$`)
	// regexVirtiofsShare is a regex pattern for matching virtiofs share names.
	regexVirtiofsShare = regexp.MustCompile(`^virtiofs\d+$`)
)

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
	CPUAffinity          *string                        `json:"affinity,omitempty"           url:"affinity,omitempty"`
	CPUArchitecture      *string                        `json:"arch,omitempty"               url:"arch,omitempty"`
	CPUCores             *int64                         `json:"cores,omitempty"              url:"cores,omitempty"`
	CPUEmulation         *CustomCPUEmulation            `json:"cpu,omitempty"                url:"cpu,omitempty"`
	CPULimit             *int64                         `json:"cpulimit,omitempty"           url:"cpulimit,omitempty"`
	CPUSockets           *int64                         `json:"sockets,omitempty"            url:"sockets,omitempty"`
	CPUUnits             *int64                         `json:"cpuunits,omitempty"           url:"cpuunits,omitempty"`
	DedicatedMemory      *int                           `json:"memory,omitempty"             url:"memory,omitempty"`
	Delete               []string                       `json:"delete,omitempty"             url:"delete,omitempty,comma"`
	DeletionProtection   *types.CustomBool              `json:"protection,omitempty"         url:"protection,omitempty,int"`
	Description          *string                        `json:"description,omitempty"        url:"description,omitempty"`
	EFIDisk              *CustomEFIDisk                 `json:"efidisk0,omitempty"           url:"efidisk0,omitempty"`
	FloatingMemory       *int                           `json:"balloon,omitempty"            url:"balloon,omitempty"`
	FloatingMemoryShares *int                           `json:"shares,omitempty"             url:"shares,omitempty"`
	Freeze               *types.CustomBool              `json:"freeze,omitempty"             url:"freeze,omitempty,int"`
	HookScript           *string                        `json:"hookscript,omitempty"         url:"hookscript,omitempty"`
	Hotplug              types.CustomCommaSeparatedList `json:"hotplug,omitempty"            url:"hotplug,omitempty,comma"`
	Hugepages            *string                        `json:"hugepages,omitempty"          url:"hugepages,omitempty"`
	KeepHugepages        *types.CustomBool              `json:"keephugepages,omitempty"      url:"keephugepages,omitempty,int"`
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
	RNGDevice            *CustomRNGDevice               `json:"rng0,omitempty"               url:"rng0,omitempty"`
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
	VirtualCPUCount      *int64                         `json:"vcpus,omitempty"              url:"vcpus,omitempty"`
	VirtiofsShares       CustomVirtiofsShares           `json:"virtiofs,omitempty"           url:"virtiofs,omitempty"`
	VMGenerationID       *string                        `json:"vmgenid,omitempty"            url:"vmgenid,omitempty"`
	VMID                 int                            `json:"vmid,omitempty"               url:"vmid,omitempty"`
	VMStateDatastoreID   *string                        `json:"vmstatestorage,omitempty"     url:"vmstatestorage,omitempty"`
	WatchdogDevice       *CustomWatchdogDevice          `json:"watchdog,omitempty"           url:"watchdog,omitempty"`
	CustomStorageDevices CustomStorageDevices           `json:"-"`
}

// AddCustomStorageDevice adds a custom storage device to the create request body.
func (b *CreateRequestBody) AddCustomStorageDevice(iface string, device CustomStorageDevice) {
	if b.CustomStorageDevices == nil {
		b.CustomStorageDevices = make(CustomStorageDevices, 1)
	}

	b.CustomStorageDevices[iface] = &device
}

// CreateResponseBody contains the body from a create response.
type CreateResponseBody struct {
	TaskID *string `json:"data,omitempty"`
}

// DeleteResponseBody contains the body from a delete response.
type DeleteResponseBody struct {
	TaskID *string `json:"data,omitempty"`
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
	CloudInitUpgrade     *types.CustomBool               `json:"ciupgrade,omitempty"`
	CPUArchitecture      *string                         `json:"arch,omitempty"`
	CPUCores             *int64                          `json:"cores,omitempty"`
	CPUEmulation         *CustomCPUEmulation             `json:"cpu,omitempty"`
	CPULimit             *types.CustomInt64              `json:"cpulimit,omitempty"`
	CPUSockets           *int64                          `json:"sockets,omitempty"`
	CPUUnits             *int64                          `json:"cpuunits,omitempty"`
	CPUAffinity          *string                         `json:"affinity,omitempty"`
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
	KeepHugepages        *types.CustomBool               `json:"keephugepages,omitempty"`
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
	NUMAEnabled          *types.CustomBool               `json:"numa,omitempty"`
	NUMADevices0         *CustomNUMADevice               `json:"numa0,omitempty"`
	NUMADevices1         *CustomNUMADevice               `json:"numa1,omitempty"`
	NUMADevices2         *CustomNUMADevice               `json:"numa2,omitempty"`
	NUMADevices3         *CustomNUMADevice               `json:"numa3,omitempty"`
	NUMADevices4         *CustomNUMADevice               `json:"numa4,omitempty"`
	NUMADevices5         *CustomNUMADevice               `json:"numa5,omitempty"`
	NUMADevices6         *CustomNUMADevice               `json:"numa6,omitempty"`
	NUMADevices7         *CustomNUMADevice               `json:"numa7,omitempty"`
	OSType               *string                         `json:"ostype,omitempty"`
	Overwrite            *types.CustomBool               `json:"force,omitempty"`
	PoolID               *string                         `json:"pool,omitempty"`
	Revert               *string                         `json:"revert,omitempty"`
	RNGDevice            *CustomRNGDevice                `json:"rng0,omitempty"`
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
	VirtualCPUCount      *int64                          `json:"vcpus,omitempty"`
	VMGenerationID       *string                         `json:"vmgenid,omitempty"`
	VMStateDatastoreID   *string                         `json:"vmstatestorage,omitempty"`
	WatchdogDevice       *CustomWatchdogDevice           `json:"watchdog,omitempty"`
	StorageDevices       CustomStorageDevices            `json:"-"`
	PCIDevices           CustomPCIDevices                `json:"-"`
	VirtiofsShares       CustomVirtiofsShares            `json:"-"`
}

// GetStatusResponseBody contains the body from a VM get status response.
type GetStatusResponseBody struct {
	Data *GetStatusResponseData `json:"data,omitempty"`
}

// GetStatusResponseData contains the data from a VM get status response.
type GetStatusResponseData struct {
	AgentEnabled     *types.CustomBool `json:"agent,omitempty"`
	CPUCount         *int64            `json:"cpus,omitempty"`
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
	Name     *string           `json:"name,omitempty"`
	Tags     *string           `json:"tags,omitempty"`
	Template *types.CustomBool `json:"template,omitempty"`
	Status   *string           `json:"status,omitempty"`
	VMID     int               `json:"vmid,omitempty"`
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
type UpdateRequestBody = CreateRequestBody

// UnmarshalJSON unmarshals the data from the JSON response, populating the CustomStorageDevices field.
func (d *GetResponseData) UnmarshalJSON(b []byte) error {
	type Alias GetResponseData

	var data Alias

	// get original struct
	if err := json.Unmarshal(b, &data); err != nil {
		return fmt.Errorf("failed to unmarshal data: %w", err)
	}

	var byAttr map[string]interface{}

	// now get map by attribute name
	err := json.Unmarshal(b, &byAttr)
	if err != nil {
		return fmt.Errorf("failed to unmarshal data: %w", err)
	}

	data.StorageDevices = make(CustomStorageDevices)
	data.PCIDevices = make(CustomPCIDevices)
	data.VirtiofsShares = make(CustomVirtiofsShares)

	for key, value := range byAttr {
		for _, prefix := range StorageInterfaces {
			// the device names can overlap with other fields, for example`scsi0` and `scsihw`, so just checking
			// the prefix is not enough
			if regexStorageInterface(prefix).MatchString(key) {
				var device CustomStorageDevice
				if err := json.Unmarshal([]byte(`"`+value.(string)+`"`), &device); err != nil {
					return fmt.Errorf("failed to unmarshal %s: %w", key, err)
				}

				data.StorageDevices[key] = &device
			}
		}

		if regexPCIDevice.MatchString(key) {
			var device CustomPCIDevice
			if err := json.Unmarshal([]byte(`"`+value.(string)+`"`), &device); err != nil {
				return fmt.Errorf("failed to unmarshal %s: %w", key, err)
			}

			data.PCIDevices[key] = &device
		}

		if regexVirtiofsShare.MatchString(key) {
			var share CustomVirtiofsShare
			if err := json.Unmarshal([]byte(`"`+value.(string)+`"`), &share); err != nil {
				return fmt.Errorf("failed to unmarshal %s: %w", key, err)
			}

			data.VirtiofsShares[key] = &share
		}
	}

	*d = GetResponseData(data)

	return nil
}

// ToDelete adds a field to the delete list. The field name should be the **actual** field name in the struct.
func (b *UpdateRequestBody) ToDelete(fieldName string) error {
	if b == nil {
		return errors.New("update request body is nil")
	}

	if field, ok := reflect.TypeOf(*b).FieldByName(fieldName); ok {
		fieldTag := field.Tag.Get("url")
		name := strings.Split(fieldTag, ",")[0]
		b.Delete = append(b.Delete, name)
	} else {
		return fmt.Errorf("field %s not found in struct %s", fieldName, reflect.TypeOf(b).Name())
	}

	return nil
}

// IsEmpty checks if the update request body is empty.
func (b *UpdateRequestBody) IsEmpty() bool {
	if b == nil {
		return true
	}

	return reflect.DeepEqual(*b, UpdateRequestBody{})
}
