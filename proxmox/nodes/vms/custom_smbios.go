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

// CustomSMBIOS handles QEMU SMBIOS parameters.
type CustomSMBIOS struct {
	Base64       *types.CustomBool `json:"base64,omitempty"       url:"base64,omitempty,int"`
	Family       *string           `json:"family,omitempty"       url:"family,omitempty"`
	Manufacturer *string           `json:"manufacturer,omitempty" url:"manufacturer,omitempty"`
	Product      *string           `json:"product,omitempty"      url:"product,omitempty"`
	Serial       *string           `json:"serial,omitempty"       url:"serial,omitempty"`
	SKU          *string           `json:"sku,omitempty"          url:"sku,omitempty"`
	UUID         *string           `json:"uuid,omitempty"         url:"uuid,omitempty"`
	Version      *string           `json:"version,omitempty"      url:"version,omitempty"`
}

// EncodeValues converts a CustomSMBIOS struct to a URL value.
func (r *CustomSMBIOS) EncodeValues(key string, v *url.Values) error {
	var values []string

	if r.Base64 != nil {
		if *r.Base64 {
			values = append(values, "base64=1")
		} else {
			values = append(values, "base64=0")
		}
	}

	if r.Family != nil {
		values = append(values, fmt.Sprintf("family=%s", *r.Family))
	}

	if r.Manufacturer != nil {
		values = append(values, fmt.Sprintf("manufacturer=%s", *r.Manufacturer))
	}

	if r.Product != nil {
		values = append(values, fmt.Sprintf("product=%s", *r.Product))
	}

	if r.Serial != nil {
		values = append(values, fmt.Sprintf("serial=%s", *r.Serial))
	}

	if r.SKU != nil {
		values = append(values, fmt.Sprintf("sku=%s", *r.SKU))
	}

	if r.UUID != nil {
		values = append(values, fmt.Sprintf("uuid=%s", *r.UUID))
	}

	if r.Version != nil {
		values = append(values, fmt.Sprintf("version=%s", *r.Version))
	}

	if len(values) > 0 {
		v.Add(key, strings.Join(values, ","))
	}

	return nil
}

// UnmarshalJSON converts a CustomSMBIOS string to an object.
func (r *CustomSMBIOS) UnmarshalJSON(b []byte) error {
	var s string

	if err := json.Unmarshal(b, &s); err != nil {
		return fmt.Errorf("failed to unmarshal CustomSMBIOS: %w", err)
	}

	pairs := strings.Split(s, ",")

	for _, p := range pairs {
		v := strings.SplitN(strings.TrimSpace(p), "=", 2)

		if len(v) == 2 {
			switch v[0] {
			case "base64":
				base64 := types.CustomBool(v[1] == "1")
				r.Base64 = &base64
			case "family":
				r.Family = &v[1]
			case "manufacturer":
				r.Manufacturer = &v[1]
			case "product":
				r.Product = &v[1]
			case "serial":
				r.Serial = &v[1]
			case "sku":
				r.SKU = &v[1]
			case "uuid":
				r.UUID = &v[1]
			case "version":
				r.Version = &v[1]
			}
		}
	}

	return nil
}
