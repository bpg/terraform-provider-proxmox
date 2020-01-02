/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmox

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// VirtualEnvironmentContainerCreateRequestBody contains the data for an user create request.
type VirtualEnvironmentContainerCreateRequestBody struct {
	BandwidthLimit       *float64                                               `json:"bwlimit,omitempty" url:"bwlimit,omitempty"`
	ConsoleEnabled       *CustomBool                                            `json:"console,omitempty" url:"console,omitempty,int"`
	ConsoleMode          *string                                                `json:"cmode,omitempty" url:"cmode,omitempty"`
	CPUArchitecture      *string                                                `json:"arch,omitempty" url:"arch,omitempty"`
	CPUCores             *int                                                   `json:"cores,omitempty" url:"cores,omitempty"`
	CPULimit             *int                                                   `json:"cpulimit,omitempty" url:"cpulimit,omitempty"`
	CPUUnits             *int                                                   `json:"cpuunits,omitempty" url:"cpuunits,omitempty"`
	DatastoreID          *string                                                `json:"storage,omitempty" url:"storage,omitempty"`
	DedicatedMemory      *int                                                   `json:"memory,omitempty" url:"memory,omitempty"`
	Delete               []string                                               `json:"delete,omitempty" url:"delete,omitempty"`
	Description          *string                                                `json:"description,omitempty" url:"description,omitempty"`
	DNSDomain            *string                                                `json:"searchdomain,omitempty" url:"searchdomain,omitempty"`
	DNSServer            *string                                                `json:"nameserver,omitempty" url:"nameserver,omitempty"`
	Features             *VirtualEnvironmentContainerCustomFeatures             `json:"features,omitempty" url:"features,omitempty"`
	Force                *CustomBool                                            `json:"force,omitempty" url:"force,omitempty,int"`
	HookScript           *string                                                `json:"hookscript,omitempty" url:"hookscript,omitempty"`
	Hostname             *string                                                `json:"hostname,omitempty" url:"hostname,omitempty"`
	IgnoreUnpackErrors   *CustomBool                                            `json:"ignore-unpack-errors,omitempty" url:"force,omitempty,int"`
	Lock                 *string                                                `json:"lock,omitempty" url:"lock,omitempty,int"`
	MountPoints          VirtualEnvironmentContainerCustomMountPointArray       `json:"mp,omitempty" url:"mp,omitempty,numbered"`
	NetworkInterfaces    VirtualEnvironmentContainerCustomNetworkInterfaceArray `json:"net,omitempty" url:"net,omitempty,numbered"`
	OSTemplateFileVolume *string                                                `json:"ostemplate,omitempty" url:"ostemplate,omitempty"`
	OSType               *string                                                `json:"ostype,omitempty" url:"ostype,omitempty"`
	Password             *string                                                `json:"password,omitempty" url:"password,omitempty"`
	PoolID               *string                                                `json:"pool,omitempty" url:"pool,omitempty"`
	Protection           *CustomBool                                            `json:"protection,omitempty" url:"protection,omitempty,int"`
	Restore              *CustomBool                                            `json:"restore,omitempty" url:"restore,omitempty,int"`
	RootFS               *VirtualEnvironmentContainerCustomRootFS               `json:"rootfs,omitempty" url:"rootfs,omitempty"`
	SSHKeys              *VirtualEnvironmentContainerCustomSSHKeys              `json:"ssh-public-keys,omitempty" url:"ssh-public-keys,omitempty"`
	Start                *CustomBool                                            `json:"start,omitempty" url:"start,omitempty,int"`
	StartOnBoot          *CustomBool                                            `json:"onboot,omitempty" url:"onboot,omitempty,int"`
	StartupBehavior      *VirtualEnvironmentContainerCustomStartupBehavior      `json:"startup,omitempty" url:"startup,omitempty"`
	Swap                 *int                                                   `json:"swap,omitempty" url:"swap,omitempty"`
	Tags                 *string                                                `json:"tags,omitempty" url:"tags,omitempty"`
	Template             *CustomBool                                            `json:"template,omitempty" url:"template,omitempty,int"`
	TTY                  *int                                                   `json:"tty,omitempty" url:"tty,omitempty"`
	Unique               *CustomBool                                            `json:"unique,omitempty" url:"unique,omitempty,int"`
	Unprivileged         *CustomBool                                            `json:"unprivileged,omitempty" url:"unprivileged,omitempty,int"`
	VMID                 *int                                                   `json:"vmid,omitempty" url:"vmid,omitempty"`
}

// VirtualEnvironmentContainerCustomFeatures contains the values for the "features" property.
type VirtualEnvironmentContainerCustomFeatures struct {
	FUSE       *CustomBool `json:"fuse,omitempty" url:"fuse,omitempty,int"`
	KeyControl *CustomBool `json:"keyctl,omitempty" url:"keyctl,omitempty,int"`
	MountTypes *[]string   `json:"mount,omitempty" url:"mount,omitempty"`
	Nesting    *CustomBool `json:"nesting,omitempty" url:"nesting,omitempty,int"`
}

// VirtualEnvironmentContainerCustomMountPoint contains the values for the "mp[n]" properties.
type VirtualEnvironmentContainerCustomMountPoint struct {
	ACL          *CustomBool `json:"acl,omitempty" url:"acl,omitempty,int"`
	Backup       *CustomBool `json:"backup,omitempty" url:"backup,omitempty,int"`
	DiskSize     *string     `json:"size,omitempty" url:"size,omitempty"`
	Enabled      bool        `json:"-" url:"-"`
	MountOptions *[]string   `json:"mountoptions,omitempty" url:"mountoptions,omitempty"`
	MountPoint   string      `json:"mp" url:"mp"`
	Quota        *CustomBool `json:"quota,omitempty" url:"quota,omitempty,int"`
	ReadOnly     *CustomBool `json:"ro,omitempty" url:"ro,omitempty,int"`
	Replicate    *CustomBool `json:"replicate,omitempty" url:"replicate,omitempty,int"`
	Shared       *CustomBool `json:"shared,omitempty" url:"shared,omitempty,int"`
	Volume       string      `json:"volume" url:"volume"`
}

// VirtualEnvironmentContainerCustomMountPointArray is an array of VirtualEnvironmentContainerCustomMountPoint.
type VirtualEnvironmentContainerCustomMountPointArray []VirtualEnvironmentContainerCustomMountPoint

// VirtualEnvironmentContainerCustomNetworkInterface contains the values for the "net[n]" properties.
type VirtualEnvironmentContainerCustomNetworkInterface struct {
	Bridge      *string     `json:"bridge,omitempty" url:"bridge,omitempty"`
	Enabled     bool        `json:"-" url:"-"`
	Firewall    *CustomBool `json:"firewall,omitempty" url:"firewall,omitempty,int"`
	IPv4Address *string     `json:"ip,omitempty" url:"ip,omitempty"`
	IPv4Gateway *string     `json:"gw,omitempty" url:"gw,omitempty"`
	IPv6Address *string     `json:"ip6,omitempty" url:"ip6,omitempty"`
	IPv6Gateway *string     `json:"gw6,omitempty" url:"gw6,omitempty"`
	MACAddress  *string     `json:"hwaddr,omitempty" url:"hwaddr,omitempty"`
	MTU         *int        `json:"mtu,omitempty" url:"mtu,omitempty"`
	Name        string      `json:"name" url:"name"`
	RateLimit   *float64    `json:"rate,omitempty" url:"rate,omitempty"`
	Tag         *int        `json:"tag,omitempty" url:"tag,omitempty"`
	Trunks      *[]int      `json:"trunks,omitempty" url:"trunks,omitempty"`
	Type        *string     `json:"type,omitempty" url:"type,omitempty"`
}

// VirtualEnvironmentContainerCustomNetworkInterfaceArray is an array of VirtualEnvironmentContainerCustomNetworkInterface.
type VirtualEnvironmentContainerCustomNetworkInterfaceArray []VirtualEnvironmentContainerCustomNetworkInterface

// VirtualEnvironmentContainerCustomRootFS contains the values for the "rootfs" property.
type VirtualEnvironmentContainerCustomRootFS struct {
	ACL          *CustomBool `json:"acl,omitempty" url:"acl,omitempty,int"`
	DiskSize     *string     `json:"size,omitempty" url:"size,omitempty"`
	MountOptions *[]string   `json:"mountoptions,omitempty" url:"mountoptions,omitempty"`
	Quota        *CustomBool `json:"quota,omitempty" url:"quota,omitempty,int"`
	ReadOnly     *CustomBool `json:"ro,omitempty" url:"ro,omitempty,int"`
	Replicate    *CustomBool `json:"replicate,omitempty" url:"replicate,omitempty,int"`
	Shared       *CustomBool `json:"shared,omitempty" url:"shared,omitempty,int"`
	Volume       string      `json:"volume" url:"volume"`
}

// VirtualEnvironmentContainerCustomSSHKeys contains the values for the "ssh-public-keys" property.
type VirtualEnvironmentContainerCustomSSHKeys []string

// VirtualEnvironmentContainerCustomStartupBehavior contains the values for the "startup" property.
type VirtualEnvironmentContainerCustomStartupBehavior struct {
	Down  *int `json:"down,omitempty" url:"down,omitempty"`
	Order *int `json:"order,omitempty" url:"order,omitempty"`
	Up    *int `json:"up,omitempty" url:"up,omitempty"`
}

// VirtualEnvironmentContainerGetResponseBody contains the body from an user get response.
type VirtualEnvironmentContainerGetResponseBody struct {
	Data *VirtualEnvironmentContainerGetResponseData `json:"data,omitempty"`
}

// VirtualEnvironmentContainerGetResponseData contains the data from an user get response.
type VirtualEnvironmentContainerGetResponseData struct {
	ConsoleEnabled    *CustomBool                                        `json:"console,omitempty"`
	ConsoleMode       *string                                            `json:"cmode,omitempty"`
	CPUArchitecture   *string                                            `json:"arch,omitempty"`
	CPUCores          *int                                               `json:"cores,omitempty"`
	CPULimit          *int                                               `json:"cpulimit,omitempty"`
	CPUUnits          *int                                               `json:"cpuunits,omitempty"`
	DedicatedMemory   *int                                               `json:"memory,omitempty"`
	Description       *string                                            `json:"description,omitempty"`
	Digest            string                                             `json:"digest"`
	DNSDomain         *string                                            `json:"searchdomain,omitempty"`
	DNSServer         *string                                            `json:"nameserver,omitempty"`
	Features          *VirtualEnvironmentContainerCustomFeatures         `json:"features,omitempty"`
	HookScript        *string                                            `json:"hookscript,omitempty"`
	Hostname          *string                                            `json:"hostname,omitempty"`
	Lock              *CustomBool                                        `json:"lock,omitempty"`
	LXCConfiguration  *[]string                                          `json:"lxc,omitempty"`
	MountPoint0       VirtualEnvironmentContainerCustomMountPointArray   `json:"mp0,omitempty"`
	MountPoint1       VirtualEnvironmentContainerCustomMountPointArray   `json:"mp1,omitempty"`
	MountPoint2       VirtualEnvironmentContainerCustomMountPointArray   `json:"mp2,omitempty"`
	MountPoint3       VirtualEnvironmentContainerCustomMountPointArray   `json:"mp3,omitempty"`
	NetworkInterface0 *VirtualEnvironmentContainerCustomNetworkInterface `json:"net0,omitempty"`
	NetworkInterface1 *VirtualEnvironmentContainerCustomNetworkInterface `json:"net1,omitempty"`
	NetworkInterface2 *VirtualEnvironmentContainerCustomNetworkInterface `json:"net2,omitempty"`
	NetworkInterface3 *VirtualEnvironmentContainerCustomNetworkInterface `json:"net3,omitempty"`
	NetworkInterface4 *VirtualEnvironmentContainerCustomNetworkInterface `json:"net4,omitempty"`
	NetworkInterface5 *VirtualEnvironmentContainerCustomNetworkInterface `json:"net5,omitempty"`
	NetworkInterface6 *VirtualEnvironmentContainerCustomNetworkInterface `json:"net6,omitempty"`
	NetworkInterface7 *VirtualEnvironmentContainerCustomNetworkInterface `json:"net7,omitempty"`
	OSType            *string                                            `json:"ostype,omitempty"`
	Protection        *CustomBool                                        `json:"protection,omitempty"`
	RootFS            *VirtualEnvironmentContainerCustomRootFS           `json:"rootfs,omitempty"`
	StartOnBoot       *CustomBool                                        `json:"onboot,omitempty"`
	StartupBehavior   *VirtualEnvironmentContainerCustomStartupBehavior  `json:"startup,omitempty"`
	Swap              *int                                               `json:"swap,omitempty"`
	Tags              *string                                            `json:"tags,omitempty"`
	Template          *CustomBool                                        `json:"template,omitempty"`
	TTY               *int                                               `json:"tty,omitempty"`
	Unprivileged      *CustomBool                                        `json:"unprivileged,omitempty"`
}

// VirtualEnvironmentContainerGetStatusResponseBody contains the body from a container get status response.
type VirtualEnvironmentContainerGetStatusResponseBody struct {
	Data *VirtualEnvironmentContainerGetStatusResponseData `json:"data,omitempty"`
}

// VirtualEnvironmentContainerGetStatusResponseData contains the data from a container get status response.
type VirtualEnvironmentContainerGetStatusResponseData struct {
	CPUCount         *float64     `json:"cpus,omitempty"`
	Lock             *string      `json:"lock,omitempty"`
	MemoryAllocation *int         `json:"maxmem,omitempty"`
	Name             *string      `json:"name,omitempty"`
	RootDiskSize     *interface{} `json:"maxdisk,omitempty"`
	Status           string       `json:"status,omitempty"`
	SwapAllocation   *int         `json:"maxswap,omitempty"`
	Tags             *string      `json:"tags,omitempty"`
	Uptime           *int         `json:"uptime,omitempty"`
	VMID             string       `json:"vmid,omitempty"`
}

// VirtualEnvironmentContainerRebootRequestBody contains the body for a container reboot request.
type VirtualEnvironmentContainerRebootRequestBody struct {
	Timeout *int `json:"timeout,omitempty" url:"timeout,omitempty"`
}

// VirtualEnvironmentContainerShutdownRequestBody contains the body for a container shutdown request.
type VirtualEnvironmentContainerShutdownRequestBody struct {
	ForceStop *CustomBool `json:"forceStop,omitempty,int" url:"forceStop,omitempty,int"`
	Timeout   *int        `json:"timeout,omitempty" url:"timeout,omitempty"`
}

// VirtualEnvironmentContainerUpdateRequestBody contains the data for an user update request.
type VirtualEnvironmentContainerUpdateRequestBody VirtualEnvironmentContainerCreateRequestBody

// EncodeValues converts a VirtualEnvironmentContainerCustomFeatures struct to a URL vlaue.
func (r VirtualEnvironmentContainerCustomFeatures) EncodeValues(key string, v *url.Values) error {
	values := []string{}

	if r.FUSE != nil {
		if *r.FUSE {
			values = append(values, "fuse=1")
		} else {
			values = append(values, "fuse=0")
		}
	}

	if r.KeyControl != nil {
		if *r.KeyControl {
			values = append(values, "keyctl=1")
		} else {
			values = append(values, "keyctl=0")
		}
	}

	if r.MountTypes != nil {
		if len(*r.MountTypes) > 0 {
			values = append(values, fmt.Sprintf("mount=%s", strings.Join(*r.MountTypes, ";")))
		}
	}

	if r.Nesting != nil {
		if *r.Nesting {
			values = append(values, "nesting=1")
		} else {
			values = append(values, "nesting=0")
		}
	}

	if len(values) > 0 {
		v.Add(key, strings.Join(values, ","))
	}

	return nil
}

// EncodeValues converts a VirtualEnvironmentContainerCustomMountPoint struct to a URL vlaue.
func (r VirtualEnvironmentContainerCustomMountPoint) EncodeValues(key string, v *url.Values) error {
	values := []string{}

	if r.ACL != nil {
		if *r.ACL {
			values = append(values, "acl=%d")
		} else {
			values = append(values, "acl=0")
		}
	}

	if r.Backup != nil {
		if *r.Backup {
			values = append(values, "backup=1")
		} else {
			values = append(values, "backup=0")
		}
	}

	if r.DiskSize != nil {
		values = append(values, fmt.Sprintf("size=%s", *r.DiskSize))
	}

	if r.MountOptions != nil {
		if len(*r.MountOptions) > 0 {
			values = append(values, fmt.Sprintf("mount=%s", strings.Join(*r.MountOptions, ";")))
		}
	}

	values = append(values, fmt.Sprintf("mp=%s", r.MountPoint))

	if r.Quota != nil {
		if *r.Quota {
			values = append(values, "quota=1")
		} else {
			values = append(values, "quota=0")
		}
	}

	if r.ReadOnly != nil {
		if *r.ReadOnly {
			values = append(values, "ro=1")
		} else {
			values = append(values, "ro=0")
		}
	}

	if r.Replicate != nil {
		if *r.ReadOnly {
			values = append(values, "replicate=1")
		} else {
			values = append(values, "replicate=0")
		}
	}

	if r.Shared != nil {
		if *r.Shared {
			values = append(values, "shared=1")
		} else {
			values = append(values, "shared=0")
		}
	}

	values = append(values, fmt.Sprintf("volume=%s", r.Volume))

	if len(values) > 0 {
		v.Add(key, strings.Join(values, ","))
	}

	return nil
}

// EncodeValues converts a VirtualEnvironmentContainerCustomMountPointArray array to multiple URL values.
func (r VirtualEnvironmentContainerCustomMountPointArray) EncodeValues(key string, v *url.Values) error {
	for i, d := range r {
		d.EncodeValues(fmt.Sprintf("%s%d", key, i), v)
	}

	return nil
}

// EncodeValues converts a VirtualEnvironmentContainerCustomNetworkInterface struct to a URL vlaue.
func (r VirtualEnvironmentContainerCustomNetworkInterface) EncodeValues(key string, v *url.Values) error {
	values := []string{}

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

	if r.IPv4Address != nil {
		values = append(values, fmt.Sprintf("ip=%s", *r.IPv4Address))
	}

	if r.IPv4Gateway != nil {
		values = append(values, fmt.Sprintf("gw=%s", *r.IPv4Gateway))
	}

	if r.IPv6Address != nil {
		values = append(values, fmt.Sprintf("ip6=%s", *r.IPv6Address))
	}

	if r.IPv6Gateway != nil {
		values = append(values, fmt.Sprintf("gw6=%s", *r.IPv6Gateway))
	}

	if r.MACAddress != nil {
		values = append(values, fmt.Sprintf("hwaddr=%s", *r.MACAddress))
	}

	if r.MTU != nil {
		values = append(values, fmt.Sprintf("mtu=%d", *r.MTU))
	}

	values = append(values, fmt.Sprintf("name=%s", r.Name))

	if r.RateLimit != nil {
		values = append(values, fmt.Sprintf("rate=%.2f", *r.RateLimit))
	}

	if r.Tag != nil {
		values = append(values, fmt.Sprintf("tag=%d", *r.Tag))
	}

	if r.Trunks != nil && len(*r.Trunks) > 0 {
		sTrunks := make([]string, len(*r.Trunks))

		for i, v := range *r.Trunks {
			sTrunks[i] = strconv.Itoa(v)
		}

		values = append(values, fmt.Sprintf("trunks=%s", strings.Join(sTrunks, ";")))
	}

	if r.Type != nil {
		values = append(values, fmt.Sprintf("type=%s", *r.Type))
	}

	if len(values) > 0 {
		v.Add(key, strings.Join(values, ","))
	}

	return nil
}

// EncodeValues converts a VirtualEnvironmentContainerCustomNetworkInterfaceArray array to multiple URL values.
func (r VirtualEnvironmentContainerCustomNetworkInterfaceArray) EncodeValues(key string, v *url.Values) error {
	for i, d := range r {
		d.EncodeValues(fmt.Sprintf("%s%d", key, i), v)
	}

	return nil
}

// EncodeValues converts a VirtualEnvironmentContainerCustomRootFS struct to a URL vlaue.
func (r VirtualEnvironmentContainerCustomRootFS) EncodeValues(key string, v *url.Values) error {
	values := []string{}

	if r.ACL != nil {
		if *r.ACL {
			values = append(values, "acl=%d")
		} else {
			values = append(values, "acl=0")
		}
	}

	if r.DiskSize != nil {
		values = append(values, fmt.Sprintf("size=%s", *r.DiskSize))
	}

	if r.MountOptions != nil {
		if len(*r.MountOptions) > 0 {
			values = append(values, fmt.Sprintf("mount=%s", strings.Join(*r.MountOptions, ";")))
		}
	}

	if r.Quota != nil {
		if *r.Quota {
			values = append(values, "quota=1")
		} else {
			values = append(values, "quota=0")
		}
	}

	if r.ReadOnly != nil {
		if *r.ReadOnly {
			values = append(values, "ro=1")
		} else {
			values = append(values, "ro=0")
		}
	}

	if r.Replicate != nil {
		if *r.ReadOnly {
			values = append(values, "replicate=1")
		} else {
			values = append(values, "replicate=0")
		}
	}

	if r.Shared != nil {
		if *r.Shared {
			values = append(values, "shared=1")
		} else {
			values = append(values, "shared=0")
		}
	}

	values = append(values, fmt.Sprintf("volume=%s", r.Volume))

	if len(values) > 0 {
		v.Add(key, strings.Join(values, ","))
	}

	return nil
}

// EncodeValues converts a VirtualEnvironmentContainerCustomSSHKeys array to a URL vlaue.
func (r VirtualEnvironmentContainerCustomSSHKeys) EncodeValues(key string, v *url.Values) error {
	v.Add(key, strings.Join(r, "\n"))

	return nil
}

// EncodeValues converts a VirtualEnvironmentContainerCustomStartupBehavior struct to a URL vlaue.
func (r VirtualEnvironmentContainerCustomStartupBehavior) EncodeValues(key string, v *url.Values) error {
	values := []string{}

	if r.Down != nil {
		values = append(values, fmt.Sprintf("down=%d", *r.Down))
	}

	if r.Order != nil {
		values = append(values, fmt.Sprintf("order=%d", *r.Order))
	}

	if r.Up != nil {
		values = append(values, fmt.Sprintf("up=%d", *r.Up))
	}

	if len(values) > 0 {
		v.Add(key, strings.Join(values, ","))
	}

	return nil
}

// UnmarshalJSON converts a VirtualEnvironmentContainerCustomFeatures string to an object.
func (r *VirtualEnvironmentContainerCustomFeatures) UnmarshalJSON(b []byte) error {
	var s string

	err := json.Unmarshal(b, &s)

	if err != nil {
		return err
	}

	pairs := strings.Split(s, ",")

	for _, p := range pairs {
		v := strings.Split(strings.TrimSpace(p), "=")

		if len(v) == 2 {
			switch v[0] {
			case "fuse":
				bv := CustomBool(v[1] == "1")
				r.FUSE = &bv
			case "keyctl":
				bv := CustomBool(v[1] == "1")
				r.KeyControl = &bv
			case "mount":
				if v[1] != "" {
					a := strings.Split(v[1], ";")
					r.MountTypes = &a
				} else {
					a := []string{}
					r.MountTypes = &a
				}
			case "nesting":
				bv := CustomBool(v[1] == "1")
				r.Nesting = &bv
			}
		}
	}

	return nil
}

// UnmarshalJSON converts a VirtualEnvironmentContainerCustomMountPoint string to an object.
func (r *VirtualEnvironmentContainerCustomMountPoint) UnmarshalJSON(b []byte) error {
	var s string

	err := json.Unmarshal(b, &s)

	if err != nil {
		return err
	}

	pairs := strings.Split(s, ",")

	for _, p := range pairs {
		v := strings.Split(strings.TrimSpace(p), "=")

		if len(v) == 1 {
			r.Volume = v[0]
		} else if len(v) == 2 {
			switch v[0] {
			case "acl":
				bv := CustomBool(v[1] == "1")
				r.ACL = &bv
			case "backup":
				bv := CustomBool(v[1] == "1")
				r.Backup = &bv
			case "mountoptions":
				if v[1] != "" {
					a := strings.Split(v[1], ";")
					r.MountOptions = &a
				} else {
					a := []string{}
					r.MountOptions = &a
				}
			case "mp":
				r.MountPoint = v[1]
			case "quota":
				bv := CustomBool(v[1] == "1")
				r.Quota = &bv
			case "ro":
				bv := CustomBool(v[1] == "1")
				r.ReadOnly = &bv
			case "replicate":
				bv := CustomBool(v[1] == "1")
				r.Replicate = &bv
			case "shared":
				bv := CustomBool(v[1] == "1")
				r.Shared = &bv
			case "size":
				r.DiskSize = &v[1]
			}
		}
	}

	return nil
}

// UnmarshalJSON converts a VirtualEnvironmentContainerCustomNetworkInterface string to an object.
func (r *VirtualEnvironmentContainerCustomNetworkInterface) UnmarshalJSON(b []byte) error {
	var s string

	err := json.Unmarshal(b, &s)

	if err != nil {
		return err
	}

	pairs := strings.Split(s, ",")

	for _, p := range pairs {
		v := strings.Split(strings.TrimSpace(p), "=")

		if len(v) == 1 {
			r.Name = v[0]
		} else if len(v) == 2 {
			switch v[0] {
			case "bridge":
				r.Bridge = &v[1]
			case "firewall":
				bv := CustomBool(v[1] == "1")
				r.Firewall = &bv
			case "gw":
				r.IPv4Gateway = &v[1]
			case "gw6":
				r.IPv6Gateway = &v[1]
			case "ip":
				r.IPv4Address = &v[1]
			case "ip6":
				r.IPv6Address = &v[1]
			case "hwaddr":
				r.MACAddress = &v[1]
			case "mtu":
				iv, err := strconv.Atoi(v[1])

				if err != nil {
					return err
				}

				r.MTU = &iv
			case "name":
				r.Name = v[1]
			case "rate":
				fv, err := strconv.ParseFloat(v[1], 64)

				if err != nil {
					return err
				}

				r.RateLimit = &fv
			case "tag":
				iv, err := strconv.Atoi(v[1])

				if err != nil {
					return err
				}

				r.Tag = &iv
			case "trunks":
				if v[1] != "" {
					trunks := strings.Split(v[1], ";")
					a := make([]int, len(trunks))

					for ti, tv := range trunks {
						a[ti], err = strconv.Atoi(tv)

						if err != nil {
							return err
						}
					}

					r.Trunks = &a
				} else {
					a := []int{}
					r.Trunks = &a
				}
			case "type":
				r.Type = &v[1]
			}
		}
	}

	return nil
}

// UnmarshalJSON converts a VirtualEnvironmentContainerCustomRootFS string to an object.
func (r *VirtualEnvironmentContainerCustomRootFS) UnmarshalJSON(b []byte) error {
	var s string

	err := json.Unmarshal(b, &s)

	if err != nil {
		return err
	}

	pairs := strings.Split(s, ",")

	for _, p := range pairs {
		v := strings.Split(strings.TrimSpace(p), "=")

		if len(v) == 1 {
			r.Volume = v[0]
		} else if len(v) == 2 {
			switch v[0] {
			case "acl":
				bv := CustomBool(v[1] == "1")
				r.ACL = &bv
			case "mountoptions":
				if v[1] != "" {
					a := strings.Split(v[1], ";")
					r.MountOptions = &a
				} else {
					a := []string{}
					r.MountOptions = &a
				}
			case "quota":
				bv := CustomBool(v[1] == "1")
				r.Quota = &bv
			case "ro":
				bv := CustomBool(v[1] == "1")
				r.ReadOnly = &bv
			case "replicate":
				bv := CustomBool(v[1] == "1")
				r.Replicate = &bv
			case "shared":
				bv := CustomBool(v[1] == "1")
				r.Shared = &bv
			case "size":
				r.DiskSize = &v[1]
			}
		}
	}

	return nil
}

// UnmarshalJSON converts a VirtualEnvironmentContainerCustomStartupBehavior string to an object.
func (r *VirtualEnvironmentContainerCustomStartupBehavior) UnmarshalJSON(b []byte) error {
	var s string

	err := json.Unmarshal(b, &s)

	if err != nil {
		return err
	}

	pairs := strings.Split(s, ",")

	for _, p := range pairs {
		v := strings.Split(strings.TrimSpace(p), "=")

		if len(v) == 2 {
			switch v[0] {
			case "down":
				iv, err := strconv.Atoi(v[1])

				if err != nil {
					return err
				}

				r.Down = &iv
			case "order":
				iv, err := strconv.Atoi(v[1])

				if err != nil {
					return err
				}

				r.Order = &iv
			case "up":
				iv, err := strconv.Atoi(v[1])

				if err != nil {
					return err
				}

				r.Up = &iv
			}
		}
	}

	return nil
}
