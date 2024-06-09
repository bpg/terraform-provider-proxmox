/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package vms

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

// CustomCloudInitConfig handles QEMU cloud-init parameters.
type CustomCloudInitConfig struct {
	Files        *CustomCloudInitFiles     `json:"cicustom,omitempty"     url:"cicustom,omitempty"`
	IPConfig     []CustomCloudInitIPConfig `json:"ipconfig,omitempty"     url:"ipconfig,omitempty,numbered"`
	Nameserver   *string                   `json:"nameserver,omitempty"   url:"nameserver,omitempty"`
	Password     *string                   `json:"cipassword,omitempty"   url:"cipassword,omitempty"`
	SearchDomain *string                   `json:"searchdomain,omitempty" url:"searchdomain,omitempty"`
	SSHKeys      *CustomCloudInitSSHKeys   `json:"sshkeys,omitempty"      url:"sshkeys,omitempty"`
	Type         *string                   `json:"citype,omitempty"       url:"citype,omitempty"`
	// Can't be reliably set, it is TRUE by default in PVE
	// Upgrade      *types.CustomBool         `json:"ciupgrade,omitempty"    url:"ciupgrade,omitempty,int"`
	Username *string `json:"ciuser,omitempty" url:"ciuser,omitempty"`
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
