/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmox

import (
	"fmt"
	"net/url"
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

// CustomIDEDevice handles QEMU host IDE device parameters.
type CustomIDEDevice struct {
	AIO           *string     `json:"aio,omitempty" url:"aio,omitempty"`
	BackupEnabled *CustomBool `json:"backup,omitempty" url:"backup,omitempty,int"`
	FileVolume    string      `json:"file" url:"file"`
}

// CustomIDEDevices handles QEMU host IDE device parameters.
type CustomIDEDevices []CustomIDEDevice

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
	Trunks     []string    `json:"trunks,omitempty" url:"trunks,omitempty"`
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

// CustomSharedMemory handles QEMU Inter-VM shared memory parameters.
type CustomSharedMemory struct {
	Name *string `json:"name,omitempty" url:"name,omitempty"`
	Size int     `json:"size" url:"size"`
}

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

// EncodeValues converts a CustomIDEDevice struct to a URL vlaue.
func (r CustomIDEDevice) EncodeValues(key string, v *url.Values) error {
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

// EncodeValues converts a CustomIDEDevices array to multiple URL values.
func (r CustomIDEDevices) EncodeValues(key string, v *url.Values) error {
	for i, d := range r {
		d.EncodeValues(fmt.Sprintf("%s[%d]", key, i), v)
	}

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
		values = append(values, fmt.Sprintf("trunks=%s", strings.Join(r.Trunks, ";")))
	}

	v.Add(key, strings.Join(values, ","))

	return nil
}

// EncodeValues converts a CustomNetworkDevices array to multiple URL values.
func (r CustomNetworkDevices) EncodeValues(key string, v *url.Values) error {
	for i, d := range r {
		d.EncodeValues(fmt.Sprintf("%s[%d]", key, i), v)
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
		d.EncodeValues(fmt.Sprintf("%s[%d]", key, i), v)
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
		d.EncodeValues(fmt.Sprintf("%s[%d]", key, i), v)
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
