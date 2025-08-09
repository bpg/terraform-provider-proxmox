/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package containers

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

var (
	// regexDeviceKey matches device keys like dev0, dev1, etc.
	regexDeviceKey = regexp.MustCompile(`^dev\d+$`)
	// regexNetworkKey matches network interface keys like net0, net1, etc.
	regexNetworkKey = regexp.MustCompile(`^net\d+$`)
	// regexMountPointKey matches mount point keys like mp0, mp1, etc.
	regexMountPointKey = regexp.MustCompile(`^mp\d+$`)
)

// CloneRequestBody contains the data for an container clone request.
type CloneRequestBody struct {
	BandwidthLimit *int              `json:"bwlimit,omitempty"     url:"bwlimit,omitempty"`
	Description    *string           `json:"description,omitempty" url:"description,omitempty"`
	FullCopy       *types.CustomBool `json:"full,omitempty"        url:"full,omitempty,int"`
	Hostname       *string           `json:"hostname,omitempty"    url:"hostname,omitempty"`
	PoolID         *string           `json:"pool,omitempty"        url:"pool,omitempty"`
	SnapshotName   *string           `json:"snapname,omitempty"    url:"snapname,omitempty"`
	TargetNodeName *string           `json:"target,omitempty"      url:"target,omitempty"`
	TargetStorage  *string           `json:"storage,omitempty"     url:"storage,omitempty"`
	VMIDNew        int               `json:"newid"                 url:"newid"`
}

// CreateRequestBody contains the data for a user create request.
type CreateRequestBody struct {
	BandwidthLimit       *float64                 `json:"bwlimit,omitempty"              url:"bwlimit,omitempty"`
	ConsoleEnabled       *types.CustomBool        `json:"console,omitempty"              url:"console,omitempty,int"`
	ConsoleMode          *string                  `json:"cmode,omitempty"                url:"cmode,omitempty"`
	CPUArchitecture      *string                  `json:"arch,omitempty"                 url:"arch,omitempty"`
	CPUCores             *int                     `json:"cores,omitempty"                url:"cores,omitempty"`
	CPULimit             *int                     `json:"cpulimit,omitempty"             url:"cpulimit,omitempty"`
	CPUUnits             *int                     `json:"cpuunits,omitempty"             url:"cpuunits,omitempty"`
	DatastoreID          *string                  `json:"storage,omitempty"              url:"storage,omitempty"`
	DedicatedMemory      *int                     `json:"memory,omitempty"               url:"memory,omitempty"`
	Delete               []string                 `json:"delete,omitempty"               url:"delete,omitempty"`
	Description          *string                  `json:"description,omitempty"          url:"description,omitempty"`
	DNSDomain            *string                  `json:"searchdomain,omitempty"         url:"searchdomain,omitempty"`
	DNSServer            *string                  `json:"nameserver,omitempty"           url:"nameserver,omitempty"`
	Features             *CustomFeatures          `json:"features,omitempty"             url:"features,omitempty"`
	Force                *types.CustomBool        `json:"force,omitempty"                url:"force,omitempty,int"`
	HookScript           *string                  `json:"hookscript,omitempty"           url:"hookscript,omitempty"`
	Hostname             *string                  `json:"hostname,omitempty"             url:"hostname,omitempty"`
	IgnoreUnpackErrors   *types.CustomBool        `json:"ignore-unpack-errors,omitempty" url:"force,omitempty,int"`
	Lock                 *string                  `json:"lock,omitempty"                 url:"lock,omitempty,int"`
	MountPoints          CustomMountPoints        `json:"mp,omitempty"                   url:"mp,omitempty,numbered"`
	PassthroughDevices   CustomPassthroughDevices `json:"dev,omitempty"                  url:"dev,omitempty,numbered"`
	NetworkInterfaces    CustomNetworkInterfaces  `json:"net,omitempty"                  url:"net,omitempty,numbered"`
	OSTemplateFileVolume *string                  `json:"ostemplate,omitempty"           url:"ostemplate,omitempty"`
	OSType               *string                  `json:"ostype,omitempty"               url:"ostype,omitempty"`
	Password             *string                  `json:"password,omitempty"             url:"password,omitempty"`
	PoolID               *string                  `json:"pool,omitempty"                 url:"pool,omitempty"`
	Protection           *types.CustomBool        `json:"protection,omitempty"           url:"protection,omitempty,int"`
	Restore              *types.CustomBool        `json:"restore,omitempty"              url:"restore,omitempty,int"`
	RootFS               *CustomRootFS            `json:"rootfs,omitempty"               url:"rootfs,omitempty"`
	SSHKeys              *CustomSSHKeys           `json:"ssh-public-keys,omitempty"      url:"ssh-public-keys,omitempty"`
	Start                *types.CustomBool        `json:"start,omitempty"                url:"start,omitempty,int"`
	StartOnBoot          *types.CustomBool        `json:"onboot,omitempty"               url:"onboot,omitempty,int"`
	StartupBehavior      *CustomStartupBehavior   `json:"startup,omitempty"              url:"startup,omitempty"`
	Swap                 *int                     `json:"swap,omitempty"                 url:"swap,omitempty"`
	Tags                 *string                  `json:"tags,omitempty"                 url:"tags,omitempty"`
	Template             *types.CustomBool        `json:"template,omitempty"             url:"template,omitempty,int"`
	TTY                  *int                     `json:"tty,omitempty"                  url:"tty,omitempty"`
	Unique               *types.CustomBool        `json:"unique,omitempty"               url:"unique,omitempty,int"`
	Unprivileged         *types.CustomBool        `json:"unprivileged,omitempty"         url:"unprivileged,omitempty,int"`
	VMID                 *int                     `json:"vmid,omitempty"                 url:"vmid,omitempty"`
}

// CustomFeatures contains the values for the "features" property.
type CustomFeatures struct {
	FUSE       *types.CustomBool `json:"fuse,omitempty"    url:"fuse,omitempty,int"`
	KeyControl *types.CustomBool `json:"keyctl,omitempty"  url:"keyctl,omitempty,int"`
	MountTypes *[]string         `json:"mount,omitempty"   url:"mount,omitempty"`
	Nesting    *types.CustomBool `json:"nesting,omitempty" url:"nesting,omitempty,int"`
}

// CustomMountPoint contains the values for the "mp[n]" properties.
type CustomMountPoint struct {
	ACL          *types.CustomBool `json:"acl,omitempty"          url:"acl,omitempty,int"`
	Backup       *types.CustomBool `json:"backup,omitempty"       url:"backup,omitempty,int"`
	DiskSize     *string           `json:"size,omitempty"         url:"size,omitempty"` // read-only
	Enabled      bool              `json:"-"                      url:"-"`
	MountOptions *[]string         `json:"mountoptions,omitempty" url:"mountoptions,omitempty"`
	MountPoint   string            `json:"mp"                     url:"mp"`
	Quota        *types.CustomBool `json:"quota,omitempty"        url:"quota,omitempty,int"`
	ReadOnly     *types.CustomBool `json:"ro,omitempty"           url:"ro,omitempty,int"`
	Replicate    *types.CustomBool `json:"replicate,omitempty"    url:"replicate,omitempty,int"`
	Shared       *types.CustomBool `json:"shared,omitempty"       url:"shared,omitempty,int"`
	Volume       string            `json:"volume"                 url:"volume"`
}

// CustomMountPoints is a map of CustomMountPoint per mount point ID (i.e. `mp1`).
type CustomMountPoints map[string]*CustomMountPoint

// CustomNetworkInterface contains the values for the "net[n]" properties.
type CustomNetworkInterface struct {
	Bridge      *string           `json:"bridge,omitempty"   url:"bridge,omitempty"`
	Enabled     bool              `json:"-"                  url:"-"`
	Firewall    *types.CustomBool `json:"firewall,omitempty" url:"firewall,omitempty,int"`
	IPv4Address *string           `json:"ip,omitempty"       url:"ip,omitempty"`
	IPv4Gateway *string           `json:"gw,omitempty"       url:"gw,omitempty"`
	IPv6Address *string           `json:"ip6,omitempty"      url:"ip6,omitempty"`
	IPv6Gateway *string           `json:"gw6,omitempty"      url:"gw6,omitempty"`
	MACAddress  *string           `json:"hwaddr,omitempty"   url:"hwaddr,omitempty"`
	MTU         *int              `json:"mtu,omitempty"      url:"mtu,omitempty"`
	Name        string            `json:"name"               url:"name"`
	RateLimit   *float64          `json:"rate,omitempty"     url:"rate,omitempty"`
	Tag         *int              `json:"tag,omitempty"      url:"tag,omitempty"`
	Trunks      *[]int            `json:"trunks,omitempty"   url:"trunks,omitempty"`
	Type        *string           `json:"type,omitempty"     url:"type,omitempty"`
}

// CustomPassthroughDevices is a map of CustomPassthroughDevice per passthrough device ID (i.e. `dev0`).
type CustomPassthroughDevices map[string]*CustomPassthroughDevice

// CustomPassthroughDevice contains the values for the "dev[n]" properties.
type CustomPassthroughDevice struct {
	DenyWrite *types.CustomBool `json:"deny-write,omitempty" url:"deny-write,omitempty,int"`
	Path      string            `json:"path"                 url:"path"`
	UID       *int              `json:"uid,omitempty"        url:"uid,omitempty"`
	GID       *int              `json:"gid,omitempty"        url:"gid,omitempty"`
	Mode      *string           `json:"mode,omitempty"       url:"mode,omitempty"`
}

// CustomNetworkInterfaces is a map of CustomNetworkInterface per network interface ID (i.e. `net0`).
type CustomNetworkInterfaces map[string]*CustomNetworkInterface

// CustomRootFS contains the values for the "rootfs" property.
type CustomRootFS struct {
	ACL          *types.CustomBool `json:"acl,omitempty"          url:"acl,omitempty,int"`
	Size         *types.DiskSize   `json:"size,omitempty"         url:"size,omitempty"`
	MountOptions *[]string         `json:"mountoptions,omitempty" url:"mountoptions,omitempty"`
	Quota        *types.CustomBool `json:"quota,omitempty"        url:"quota,omitempty,int"`
	ReadOnly     *types.CustomBool `json:"ro,omitempty"           url:"ro,omitempty,int"`
	Replicate    *types.CustomBool `json:"replicate,omitempty"    url:"replicate,omitempty,int"`
	Shared       *types.CustomBool `json:"shared,omitempty"       url:"shared,omitempty,int"`
	Volume       string            `json:"volume"                 url:"volume"`
}

// CustomSSHKeys contains the values for the "ssh-public-keys" property.
type CustomSSHKeys []string

// CustomStartupBehavior contains the values for the "startup" property.
type CustomStartupBehavior struct {
	Down  *int `json:"down,omitempty"  url:"down,omitempty"`
	Order *int `json:"order,omitempty" url:"order,omitempty"`
	Up    *int `json:"up,omitempty"    url:"up,omitempty"`
}

// CreateResponseBody contains the body from a container create response.
type CreateResponseBody struct {
	Data *string `json:"data,omitempty"`
}

type ShutdownResponseBody = CreateResponseBody

// GetResponseBody contains the body from a user get response.
type GetResponseBody struct {
	Data *GetResponseData `json:"data,omitempty"`
}

// GetResponseData contains the data from a user get response.
type GetResponseData struct {
	ConsoleEnabled     *types.CustomBool        `json:"console,omitempty"`
	ConsoleMode        *string                  `json:"cmode,omitempty"`
	CPUArchitecture    *string                  `json:"arch,omitempty"`
	CPUCores           *int                     `json:"cores,omitempty"`
	CPULimit           *types.CustomInt         `json:"cpulimit,omitempty"`
	CPUUnits           *int                     `json:"cpuunits,omitempty"`
	DedicatedMemory    *int                     `json:"memory,omitempty"`
	Description        *string                  `json:"description,omitempty"`
	Digest             string                   `json:"digest"`
	DNSDomain          *string                  `json:"searchdomain,omitempty"`
	DNSServer          *string                  `json:"nameserver,omitempty"`
	Features           *CustomFeatures          `json:"features,omitempty"`
	HookScript         *string                  `json:"hookscript,omitempty"`
	Hostname           *string                  `json:"hostname,omitempty"`
	Lock               *types.CustomBool        `json:"lock,omitempty"`
	LXCConfiguration   *[][2]string             `json:"lxc,omitempty"`
	MountPoints        CustomMountPoints        `json:"mp,omitempty"`
	PassthroughDevices CustomPassthroughDevices `json:"dev,omitempty"`
	NetworkInterfaces  CustomNetworkInterfaces  `json:"net,omitempty"`
	OSType             *string                  `json:"ostype,omitempty"`
	Protection         *types.CustomBool        `json:"protection,omitempty"`
	RootFS             *CustomRootFS            `json:"rootfs,omitempty"`
	StartOnBoot        *types.CustomBool        `json:"onboot,omitempty"`
	StartupBehavior    *CustomStartupBehavior   `json:"startup,omitempty"`
	Swap               *int                     `json:"swap,omitempty"`
	Tags               *string                  `json:"tags,omitempty"`
	Template           *types.CustomBool        `json:"template,omitempty"`
	TTY                *int                     `json:"tty,omitempty"`
	Unprivileged       *types.CustomBool        `json:"unprivileged,omitempty"`
}

// GetStatusResponseBody contains the body from a container get status response.
type GetStatusResponseBody struct {
	Data *GetStatusResponseData `json:"data,omitempty"`
}

// GetStatusResponseData contains the data from a container get status response.
type GetStatusResponseData struct {
	CPUCount         *float64         `json:"cpus,omitempty"`
	Lock             *string          `json:"lock,omitempty"`
	MemoryAllocation *int64           `json:"maxmem,omitempty"`
	Name             *string          `json:"name,omitempty"`
	RootDiskSize     *interface{}     `json:"maxdisk,omitempty"`
	Status           string           `json:"status,omitempty"`
	SwapAllocation   *int64           `json:"maxswap,omitempty"`
	Tags             *string          `json:"tags,omitempty"`
	Uptime           *int             `json:"uptime,omitempty"`
	VMID             *types.CustomInt `json:"vmid,omitempty"`
}

// ListResponseBody contains the body from a container list response.
type ListResponseBody struct {
	Data []*ListResponseData `json:"data,omitempty"`
}

// ListResponseData contains the data from an container list response.
type ListResponseData struct {
	Name     *string           `json:"name,omitempty"`
	Tags     *string           `json:"tags,omitempty"`
	Template *types.CustomBool `json:"template,omitempty"`
	Status   *string           `json:"status,omitempty"`
	VMID     int               `json:"vmid,omitempty"`
}

// GetNetworkInterfaceResponseBody contains the body from a container get network interface response.
type GetNetworkInterfaceResponseBody struct {
	Data []GetNetworkInterfacesData `json:"data,omitempty"`
}

// GetNetworkInterfacesData contains the data from a container get network interfaces response.
type GetNetworkInterfacesData struct {
	MACAddress  string `json:"hardware-address"`
	Name        string `json:"name"`
	IPAddresses *[]struct {
		Address string          `json:"ip-address"`
		Prefix  types.CustomInt `json:"prefix"`
		Type    string          `json:"ip-address-type"`
	} `json:"ip-addresses,omitempty"`
}

// StartResponseBody contains the body from a container start response.
type StartResponseBody struct {
	Data *string `json:"data,omitempty"`
}

// RebootRequestBody contains the body for a container reboot request.
type RebootRequestBody struct {
	Timeout *int `json:"timeout,omitempty" url:"timeout,omitempty"`
}

// ShutdownRequestBody contains the body for a container shutdown request.
type ShutdownRequestBody struct {
	ForceStop *types.CustomBool `json:"forceStop,omitempty" url:"forceStop,omitempty,int"`
	Timeout   *int              `json:"timeout,omitempty"   url:"timeout,omitempty"`
}

// UpdateRequestBody contains the data for an user update request.
type UpdateRequestBody CreateRequestBody

// EncodeValues converts a ContainerCustomFeatures struct to a URL value.
func (r *CustomFeatures) EncodeValues(key string, v *url.Values) error {
	var values []string

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

// EncodeValues converts a CustomPassthroughDevice struct to a URL value.
func (r *CustomPassthroughDevice) EncodeValues(key string, v *url.Values) error {
	var values []string

	if r.DenyWrite != nil {
		if *r.DenyWrite {
			values = append(values, "deny-write=1")
		} else {
			values = append(values, "deny-write=0")
		}
	}

	if r.Path != "" {
		values = append(values, fmt.Sprintf("path=%s", r.Path))
	}

	if r.UID != nil {
		values = append(values, fmt.Sprintf("uid=%d", *r.UID))
	}

	if r.GID != nil {
		values = append(values, fmt.Sprintf("gid=%d", *r.GID))
	}

	if r.Mode != nil && *r.Mode != "" {
		values = append(values, fmt.Sprintf("mode=%s", *r.Mode))
	}

	if len(values) > 0 {
		v.Add(key, strings.Join(values, ","))
	}

	return nil
}

// EncodeValues converts a CustomPassthroughDevices array to multiple URL values.
func (r CustomPassthroughDevices) EncodeValues(
	_ string,
	v *url.Values,
) error {
	for key, d := range r {
		if err := d.EncodeValues(key, v); err != nil {
			return fmt.Errorf("failed to encode CustomPassthroughDevices: %w", err)
		}
	}

	return nil
}

// EncodeValues converts a CustomMountPoint struct to a URL value.
func (r *CustomMountPoint) EncodeValues(key string, v *url.Values) error {
	var values []string

	if r.ACL != nil {
		if *r.ACL {
			values = append(values, "acl=1")
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
			values = append(values, fmt.Sprintf("mountoptions=%s", strings.Join(*r.MountOptions, ";")))
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
		if *r.Replicate {
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

// EncodeValues converts a CustomMountPoints array to multiple URL values.
func (r CustomMountPoints) EncodeValues(
	_ string,
	v *url.Values,
) error {
	for key, d := range r {
		if err := d.EncodeValues(key, v); err != nil {
			return fmt.Errorf("failed to encode CustomMountPoints: %w", err)
		}
	}

	return nil
}

// EncodeValues converts a CustomNetworkInterface struct to a URL value.
func (r *CustomNetworkInterface) EncodeValues(
	key string,
	v *url.Values,
) error {
	var values []string

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

// EncodeValues converts a CustomNetworkInterfaces array to multiple URL values.
func (r CustomNetworkInterfaces) EncodeValues(
	_ string,
	v *url.Values,
) error {
	for key, d := range r {
		if err := d.EncodeValues(key, v); err != nil {
			return fmt.Errorf("failed to encode CustomNetworkInterfaces: %w", err)
		}
	}

	return nil
}

// EncodeValues converts a CustomRootFS struct to a URL value.
func (r *CustomRootFS) EncodeValues(key string, v *url.Values) error {
	var values []string

	if r.ACL != nil {
		if *r.ACL {
			values = append(values, "acl=1")
		} else {
			values = append(values, "acl=0")
		}
	}

	if r.Size != nil {
		values = append(values, fmt.Sprintf("size=%d", *r.Size))
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
		if *r.Replicate {
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

// EncodeValues converts a CustomSSHKeys array to a URL value.
func (r CustomSSHKeys) EncodeValues(key string, v *url.Values) error {
	v.Add(key, strings.Join(r, "\n"))

	return nil
}

// EncodeValues converts a CustomStartupBehavior struct to a URL value.
func (r *CustomStartupBehavior) EncodeValues(
	key string,
	v *url.Values,
) error {
	var values []string

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

// UnmarshalJSON converts a ContainerCustomFeatures string to an object.
func (r *CustomFeatures) UnmarshalJSON(b []byte) error {
	var s string

	err := json.Unmarshal(b, &s)
	if err != nil {
		return fmt.Errorf("unable to unmarshal ContainerCustomFeatures: %w", err)
	}

	pairs := strings.Split(s, ",")

	for _, p := range pairs {
		v := strings.Split(strings.TrimSpace(p), "=")

		if len(v) == 2 {
			switch v[0] {
			case "fuse":
				bv := types.CustomBool(v[1] == "1")
				r.FUSE = &bv
			case "keyctl":
				bv := types.CustomBool(v[1] == "1")
				r.KeyControl = &bv
			case "mount":
				if v[1] != "" {
					a := strings.Split(v[1], ";")
					r.MountTypes = &a
				} else {
					var a []string

					r.MountTypes = &a
				}
			case "nesting":
				bv := types.CustomBool(v[1] == "1")
				r.Nesting = &bv
			}
		}
	}

	return nil
}

// UnmarshalJSON converts a CustomPassthroughDevice string to an object.
func (r *CustomPassthroughDevice) UnmarshalJSON(b []byte) error {
	var s string

	err := json.Unmarshal(b, &s)
	if err != nil {
		return fmt.Errorf("unable to unmarshal CustomPassthroughDevice: %w", err)
	}

	pairs := strings.Split(s, ",")

	var path string

	for _, p := range pairs {
		v := strings.Split(strings.TrimSpace(p), "=")

		if len(v) == 1 {
			path = v[0]
		} else if len(v) == 2 {
			switch v[0] {
			case "deny-write":
				bv := types.CustomBool(v[1] == "1")
				r.DenyWrite = &bv
			case "path":
				path = v[1]
			case "uid":
				iv, err := strconv.Atoi(v[1])
				if err != nil {
					return fmt.Errorf("unable to unmarshal 'uid': %w", err)
				}

				r.UID = &iv
			case "gid":
				iv, err := strconv.Atoi(v[1])
				if err != nil {
					return fmt.Errorf("unable to unmarshal 'gid': %w", err)
				}

				r.GID = &iv
			case "mode":
				r.Mode = &v[1]
			}
		}
	}

	r.Path = path

	return nil
}

// UnmarshalJSON converts a CustomMountPoint string to an object.
func (r *CustomMountPoint) UnmarshalJSON(b []byte) error {
	var s string

	err := json.Unmarshal(b, &s)
	if err != nil {
		return fmt.Errorf("unable to unmarshal CustomMountPoint: %w", err)
	}

	pairs := strings.Split(s, ",")

	for _, p := range pairs {
		v := strings.Split(strings.TrimSpace(p), "=")

		if len(v) == 1 {
			r.Volume = v[0]
		} else if len(v) == 2 {
			switch v[0] {
			case "acl":
				bv := types.CustomBool(v[1] == "1")
				r.ACL = &bv
			case "backup":
				bv := types.CustomBool(v[1] == "1")
				r.Backup = &bv
			case "mountoptions":
				if v[1] != "" {
					a := strings.Split(v[1], ";")
					r.MountOptions = &a
				} else {
					var a []string

					r.MountOptions = &a
				}
			case "mp":
				r.MountPoint = v[1]
			case "quota":
				bv := types.CustomBool(v[1] == "1")
				r.Quota = &bv
			case "ro":
				bv := types.CustomBool(v[1] == "1")
				r.ReadOnly = &bv
			case "replicate":
				bv := types.CustomBool(v[1] == "1")
				r.Replicate = &bv
			case "shared":
				bv := types.CustomBool(v[1] == "1")
				r.Shared = &bv
			case "size":
				r.DiskSize = &v[1]
			}
		}
	}

	return nil
}

// UnmarshalJSON converts a CustomNetworkInterface string to an object.
func (r *CustomNetworkInterface) UnmarshalJSON(b []byte) error {
	var s string

	er := json.Unmarshal(b, &s)
	if er != nil {
		return fmt.Errorf("unable to unmarshal CustomNetworkInterface: %w", er)
	}

	pairs := strings.Split(s, ",")

	for _, p := range pairs {
		v := strings.Split(strings.TrimSpace(p), "=")

		//nolint:nestif
		if len(v) == 1 {
			r.Name = v[0]
		} else if len(v) == 2 {
			switch v[0] {
			case "bridge":
				r.Bridge = &v[1]
			case "firewall":
				bv := types.CustomBool(v[1] == "1")
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
					return fmt.Errorf("unable to unmarshal 'mtu': %w", err)
				}

				r.MTU = &iv
			case "name":
				r.Name = v[1]
			case "rate":
				fv, err := strconv.ParseFloat(v[1], 64)
				if err != nil {
					return fmt.Errorf("unable to unmarshal 'rate': %w", err)
				}

				r.RateLimit = &fv
			case "tag":
				iv, err := strconv.Atoi(v[1])
				if err != nil {
					return fmt.Errorf("unable to unmarshal 'tag': %w", err)
				}

				r.Tag = &iv
			case "trunks":
				var err error

				if v[1] != "" {
					trunks := strings.Split(v[1], ";")
					a := make([]int, len(trunks))

					for ti, tv := range trunks {
						a[ti], err = strconv.Atoi(tv)
						if err != nil {
							return fmt.Errorf("unable to unmarshal 'trunks': %w", err)
						}
					}

					r.Trunks = &a
				} else {
					var a []int

					r.Trunks = &a
				}
			case "type":
				r.Type = &v[1]
			}
		}
	}

	return nil
}

// UnmarshalJSON converts a CustomRootFS string to an object.
func (r *CustomRootFS) UnmarshalJSON(b []byte) error {
	var s string

	err := json.Unmarshal(b, &s)
	if err != nil {
		return fmt.Errorf("unable to unmarshal CustomRootFS: %w", err)
	}

	pairs := strings.Split(s, ",")

	for _, p := range pairs {
		v := strings.Split(strings.TrimSpace(p), "=")

		if len(v) == 1 {
			r.Volume = v[0]
		} else if len(v) == 2 {
			switch v[0] {
			case "acl":
				bv := types.CustomBool(v[1] == "1")
				r.ACL = &bv
			case "mountoptions":
				if v[1] != "" {
					a := strings.Split(v[1], ";")
					r.MountOptions = &a
				} else {
					var a []string

					r.MountOptions = &a
				}
			case "quota":
				bv := types.CustomBool(v[1] == "1")
				r.Quota = &bv
			case "ro":
				bv := types.CustomBool(v[1] == "1")
				r.ReadOnly = &bv
			case "replicate":
				bv := types.CustomBool(v[1] == "1")
				r.Replicate = &bv
			case "shared":
				bv := types.CustomBool(v[1] == "1")
				r.Shared = &bv
			case "size":
				r.Size = new(types.DiskSize)

				err := r.Size.UnmarshalJSON([]byte(v[1]))
				if err != nil {
					return fmt.Errorf("failed to unmarshal disk size: %w", err)
				}
			}
		}
	}

	return nil
}

// UnmarshalJSON converts a CustomStartupBehavior string to an object.
func (r *CustomStartupBehavior) UnmarshalJSON(b []byte) error {
	var s string

	err := json.Unmarshal(b, &s)
	if err != nil {
		return fmt.Errorf("unable to unmarshal CustomStartupBehavior: %w", err)
	}

	pairs := strings.Split(s, ",")

	for _, p := range pairs {
		v := strings.Split(strings.TrimSpace(p), "=")

		if len(v) == 2 {
			switch v[0] {
			case "down":
				iv, err := strconv.Atoi(v[1])
				if err != nil {
					return fmt.Errorf("unable to unmarshal 'down': %w", err)
				}

				r.Down = &iv
			case "order":
				iv, err := strconv.Atoi(v[1])
				if err != nil {
					return fmt.Errorf("unable to unmarshal 'order': %w", err)
				}

				r.Order = &iv
			case "up":
				iv, err := strconv.Atoi(v[1])
				if err != nil {
					return fmt.Errorf("unable to unmarshal 'up': %w", err)
				}

				r.Up = &iv
			}
		}
	}

	return nil
}

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

	data.PassthroughDevices = make(CustomPassthroughDevices)
	data.NetworkInterfaces = make(CustomNetworkInterfaces)
	data.MountPoints = make(CustomMountPoints)

	for key, value := range byAttr {
		valueStr, ok := value.(string)
		if !ok {
			continue // Skip non-string values
		}

		jsonValue := []byte(`"` + valueStr + `"`)

		switch {
		case regexDeviceKey.MatchString(key):
			var device CustomPassthroughDevice
			if e := json.Unmarshal(jsonValue, &device); e != nil {
				return fmt.Errorf("failed to unmarshal %s with value %q: %w", key, valueStr, e)
			}

			data.PassthroughDevices[key] = &device

		case regexNetworkKey.MatchString(key):
			var net CustomNetworkInterface
			if e := json.Unmarshal(jsonValue, &net); e != nil {
				return fmt.Errorf("failed to unmarshal %s with value %q: %w", key, valueStr, e)
			}

			data.NetworkInterfaces[key] = &net

		case regexMountPointKey.MatchString(key):
			var mp CustomMountPoint
			if e := json.Unmarshal(jsonValue, &mp); e != nil {
				return fmt.Errorf("failed to unmarshal %s with value %q: %w", key, valueStr, e)
			}

			data.MountPoints[key] = &mp
		}
	}

	*d = GetResponseData(data)

	return nil
}
