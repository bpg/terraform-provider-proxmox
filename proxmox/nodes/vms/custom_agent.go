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

// CustomAgent handles QEMU agent parameters.
type CustomAgent struct {
	Enabled         *types.CustomBool `json:"enabled,omitempty"   url:"enabled,int"`
	TrimClonedDisks *types.CustomBool `json:"fstrim_cloned_disks" url:"fstrim_cloned_disks,int"`
	Type            *string           `json:"type"                url:"type"`
}

// EncodeValues converts a CustomAgent struct to a URL value.
func (r *CustomAgent) EncodeValues(key string, v *url.Values) error {
	var values []string

	if r.Enabled != nil {
		if *r.Enabled {
			values = append(values, "enabled=1")
		} else {
			values = append(values, "enabled=0")
		}
	}

	if r.TrimClonedDisks != nil {
		if *r.TrimClonedDisks {
			values = append(values, "fstrim_cloned_disks=1")
		} else {
			values = append(values, "fstrim_cloned_disks=0")
		}
	}

	if r.Type != nil {
		values = append(values, fmt.Sprintf("type=%s", *r.Type))
	}

	if len(values) > 0 {
		v.Add(key, strings.Join(values, ","))
	}

	return nil
}

// UnmarshalJSON converts a CustomAgent string to an object.
func (r *CustomAgent) UnmarshalJSON(b []byte) error {
	var s string

	if err := json.Unmarshal(b, &s); err != nil {
		return fmt.Errorf("error unmarshalling CustomAgent: %w", err)
	}

	pairs := strings.Split(s, ",")

	for _, p := range pairs {
		v := strings.Split(strings.TrimSpace(p), "=")

		if len(v) == 1 {
			enabled := types.CustomBool(v[0] == "1")
			r.Enabled = &enabled
		} else if len(v) == 2 {
			switch v[0] {
			case "enabled":
				enabled := types.CustomBool(v[1] == "1")
				r.Enabled = &enabled
			case "fstrim_cloned_disks":
				fstrim := types.CustomBool(v[1] == "1")
				r.TrimClonedDisks = &fstrim
			case "type":
				r.Type = &v[1]
			}
		}
	}

	return nil
}
