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

// CustomEFIDisk handles QEMU EFI disk parameters.
type CustomEFIDisk struct {
	FileVolume      string            `json:"file"                        url:"file"`
	Format          *string           `json:"format,omitempty"            url:"format,omitempty"`
	Type            *string           `json:"efitype,omitempty"           url:"efitype,omitempty"`
	PreEnrolledKeys *types.CustomBool `json:"pre-enrolled-keys,omitempty" url:"pre-enrolled-keys,omitempty,int"`
}

// EncodeValues converts a CustomEFIDisk struct to a URL value.
func (r *CustomEFIDisk) EncodeValues(key string, v *url.Values) error {
	values := []string{
		fmt.Sprintf("file=%s", r.FileVolume),
	}

	if r.Format != nil {
		values = append(values, fmt.Sprintf("format=%s", *r.Format))
	}

	if r.Type != nil {
		values = append(values, fmt.Sprintf("efitype=%s", *r.Type))
	}

	if r.PreEnrolledKeys != nil {
		if *r.PreEnrolledKeys {
			values = append(values, "pre-enrolled-keys=1")
		} else {
			values = append(values, "pre-enrolled-keys=0")
		}
	}

	v.Add(key, strings.Join(values, ","))

	return nil
}

// UnmarshalJSON converts a CustomEFIDisk string to an object.
func (r *CustomEFIDisk) UnmarshalJSON(b []byte) error {
	var s string

	if err := json.Unmarshal(b, &s); err != nil {
		return fmt.Errorf("failed to unmarshal CustomEFIDisk: %w", err)
	}

	pairs := strings.Split(s, ",")

	for i, p := range pairs {
		v := strings.Split(strings.TrimSpace(p), "=")

		if len(v) == 1 && i == 0 {
			r.FileVolume = v[0]
		}

		if len(v) == 2 {
			switch v[0] {
			case "file":
				r.FileVolume = v[1]
			case "format":
				r.Format = &v[1]
			case "efitype":
				t := strings.ToLower(v[1])
				r.Type = &t
			case "pre-enrolled-keys":
				bv := types.CustomBool(v[1] == "1")
				r.PreEnrolledKeys = &bv
			}
		}
	}

	return nil
}
