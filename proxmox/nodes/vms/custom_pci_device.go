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
type CustomPCIDevices map[string]*CustomPCIDevice

// EncodeValues converts a CustomPCIDevice struct to a URL value.
func (d *CustomPCIDevice) EncodeValues(key string, v *url.Values) error {
	var values []string

	if d.DeviceIDs == nil && d.Mapping == nil {
		return fmt.Errorf("either device ID or resource mapping must be set")
	}

	if d.DeviceIDs != nil {
		values = append(values, fmt.Sprintf("host=%s", strings.Join(*d.DeviceIDs, ";")))
	}

	if d.Mapping != nil {
		values = append(values, fmt.Sprintf("mapping=%s", *d.Mapping))
	}

	if d.MDev != nil {
		values = append(values, fmt.Sprintf("mdev=%s", *d.MDev))
	}

	if d.PCIExpress != nil {
		if *d.PCIExpress {
			values = append(values, "pcie=1")
		} else {
			values = append(values, "pcie=0")
		}
	}

	if d.ROMBAR != nil {
		if *d.ROMBAR {
			values = append(values, "rombar=1")
		} else {
			values = append(values, "rombar=0")
		}
	}

	if d.ROMFile != nil {
		values = append(values, fmt.Sprintf("romfile=%s", *d.ROMFile))
	}

	if d.XVGA != nil {
		if *d.XVGA {
			values = append(values, "x-vga=1")
		} else {
			values = append(values, "x-vga=0")
		}
	}

	v.Add(key, strings.Join(values, ","))

	return nil
}

// EncodeValues converts a CustomPCIDevices array to multiple URL values.
func (r CustomPCIDevices) EncodeValues(_ string, v *url.Values) error {
	for s, d := range r {
		if err := d.EncodeValues(s, v); err != nil {
			return fmt.Errorf("failed to encode PCI device %s: %w", s, err)
		}
	}

	return nil
}

// UnmarshalJSON converts a CustomPCIDevice string to an object.
func (d *CustomPCIDevice) UnmarshalJSON(b []byte) error {
	var s string

	if err := json.Unmarshal(b, &s); err != nil {
		return fmt.Errorf("failed to unmarshal CustomPCIDevice: %w", err)
	}

	pairs := strings.Split(s, ",")

	for _, p := range pairs {
		v := strings.Split(strings.TrimSpace(p), "=")
		if len(v) == 1 {
			dIDs := strings.Split(v[0], ";")
			d.DeviceIDs = &dIDs
		} else if len(v) == 2 {
			switch v[0] {
			case "host":
				dIDs := strings.Split(v[1], ";")
				d.DeviceIDs = &dIDs
			case "mapping":
				d.Mapping = &v[1]
			case "mdev":
				d.MDev = &v[1]
			case "pcie":
				bv := types.CustomBool(v[1] == "1")
				d.PCIExpress = &bv
			case "rombar":
				bv := types.CustomBool(v[1] == "1")
				d.ROMBAR = &bv
			case "romfile":
				d.ROMFile = &v[1]
			case "x-vga":
				bv := types.CustomBool(v[1] == "1")
				d.XVGA = &bv
			}
		}
	}

	return nil
}
