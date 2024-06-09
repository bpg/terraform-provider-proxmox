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
type CustomPCIDevices []CustomPCIDevice

// EncodeValues converts a CustomPCIDevice struct to a URL value.
func (r *CustomPCIDevice) EncodeValues(key string, v *url.Values) error {
	var values []string

	if r.DeviceIDs == nil && r.Mapping == nil {
		return fmt.Errorf("either device ID or resource mapping must be set")
	}

	if r.DeviceIDs != nil {
		values = append(values, fmt.Sprintf("host=%s", strings.Join(*r.DeviceIDs, ";")))
	}

	if r.Mapping != nil {
		values = append(values, fmt.Sprintf("mapping=%s", *r.Mapping))
	}

	if r.MDev != nil {
		values = append(values, fmt.Sprintf("mdev=%s", *r.MDev))
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
		if err := d.EncodeValues(fmt.Sprintf("%s%d", key, i), v); err != nil {
			return fmt.Errorf("failed to encode PCI device %d: %w", i, err)
		}
	}

	return nil
}

// UnmarshalJSON converts a CustomPCIDevice string to an object.
func (r *CustomPCIDevice) UnmarshalJSON(b []byte) error {
	var s string

	if err := json.Unmarshal(b, &s); err != nil {
		return fmt.Errorf("failed to unmarshal CustomPCIDevice: %w", err)
	}

	pairs := strings.Split(s, ",")

	for _, p := range pairs {
		v := strings.Split(strings.TrimSpace(p), "=")
		if len(v) == 1 {
			dIDs := strings.Split(v[0], ";")
			r.DeviceIDs = &dIDs
		} else if len(v) == 2 {
			switch v[0] {
			case "host":
				dIDs := strings.Split(v[1], ";")
				r.DeviceIDs = &dIDs
			case "mapping":
				r.Mapping = &v[1]
			case "mdev":
				r.MDev = &v[1]
			case "pcie":
				bv := types.CustomBool(v[1] == "1")
				r.PCIExpress = &bv
			case "rombar":
				bv := types.CustomBool(v[1] == "1")
				r.ROMBAR = &bv
			case "romfile":
				r.ROMFile = &v[1]
			case "x-vga":
				bv := types.CustomBool(v[1] == "1")
				r.XVGA = &bv
			}
		}
	}

	return nil
}
