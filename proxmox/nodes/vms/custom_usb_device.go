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

	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

// CustomUSBDevice handles QEMU USB device parameters.
type CustomUSBDevice struct {
	HostDevice *string           `json:"host"              url:"host"`
	Mapping    *string           `json:"mapping,omitempty" url:"mapping,omitempty"`
	USB3       *types.CustomBool `json:"usb3,omitempty"    url:"usb3,omitempty,int"`
}

// CustomUSBDevices handles QEMU USB device parameters.
type CustomUSBDevices []CustomUSBDevice

// EncodeValues converts a CustomUSBDevice struct to a URL value.
func (r *CustomUSBDevice) EncodeValues(key string, v *url.Values) error {
	if r.HostDevice == nil && r.Mapping == nil {
		return fmt.Errorf("either device ID or resource mapping must be set")
	}

	var values []string
	if r.HostDevice != nil {
		values = append(values, fmt.Sprintf("host=%s", *(r.HostDevice)))
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
			r.HostDevice = &v[0]
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
