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
	"strconv"
	"strings"

	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

// CustomNetworkDevice handles QEMU network device parameters.
type CustomNetworkDevice struct {
	Enabled    bool              `json:"-"                   url:"-"`
	Bridge     *string           `json:"bridge,omitempty"    url:"bridge,omitempty"`
	Firewall   *types.CustomBool `json:"firewall,omitempty"  url:"firewall,omitempty,int"`
	LinkDown   *types.CustomBool `json:"link_down,omitempty" url:"link_down,omitempty,int"`
	MACAddress *string           `json:"macaddr,omitempty"   url:"macaddr,omitempty"`
	MTU        *int              `json:"mtu,omitempty"       url:"mtu,omitempty"`
	Model      string            `json:"model"               url:"model"`
	Queues     *int              `json:"queues,omitempty"    url:"queues,omitempty"`
	RateLimit  *float64          `json:"rate,omitempty"      url:"rate,omitempty"`
	Tag        *int              `json:"tag,omitempty"       url:"tag,omitempty"`
	Trunks     []int             `json:"trunks,omitempty"    url:"trunks,omitempty"`
}

// CustomNetworkDevices handles QEMU network device parameters.
type CustomNetworkDevices []CustomNetworkDevice

// EncodeValues converts a CustomNetworkDevice struct to a URL value.
func (r *CustomNetworkDevice) EncodeValues(key string, v *url.Values) error {
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
