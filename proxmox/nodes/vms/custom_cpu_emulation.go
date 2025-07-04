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
	"strings"

	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

// CustomCPUEmulation handles QEMU CPU emulation parameters.
type CustomCPUEmulation struct {
	Flags      *[]string         `json:"flags,omitempty"        url:"flags,omitempty,semicolon"`
	Hidden     *types.CustomBool `json:"hidden,omitempty"       url:"hidden,omitempty,int"`
	HVVendorID *string           `json:"hv-vendor-id,omitempty" url:"hv-vendor-id,omitempty"`
	Type       string            `json:"cputype,omitempty"      url:"cputype,omitempty"`
}

// EncodeValues converts a CustomCPUEmulation struct to a URL value.
func (r *CustomCPUEmulation) EncodeValues(key string, v *url.Values) error {
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
