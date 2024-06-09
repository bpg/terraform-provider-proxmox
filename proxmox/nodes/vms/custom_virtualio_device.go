/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package vms

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

// CustomVirtualIODevice handles QEMU VirtIO device parameters.
type CustomVirtualIODevice struct {
	AIO           *string           `json:"aio,omitempty"    url:"aio,omitempty"`
	BackupEnabled *types.CustomBool `json:"backup,omitempty" url:"backup,omitempty,int"`
	Enabled       bool              `json:"-"                url:"-"`
	FileVolume    string            `json:"file"             url:"file"`
}

// CustomVirtualIODevices handles QEMU VirtIO device parameters.
type CustomVirtualIODevices []CustomVirtualIODevice

// EncodeValues converts a CustomVirtualIODevice struct to a URL value.
func (r CustomVirtualIODevice) EncodeValues(key string, v *url.Values) error {
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

// EncodeValues converts a CustomVirtualIODevices array to multiple URL values.
func (r CustomVirtualIODevices) EncodeValues(key string, v *url.Values) error {
	for i, d := range r {
		if d.Enabled {
			if err := d.EncodeValues(fmt.Sprintf("%s%d", key, i), v); err != nil {
				return fmt.Errorf("error encoding virtual IO device %d: %w", i, err)
			}
		}
	}

	return nil
}
