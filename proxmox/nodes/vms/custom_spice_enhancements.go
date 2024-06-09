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

// CustomSpiceEnhancements handles QEMU spice enhancement parameters.
type CustomSpiceEnhancements struct {
	FolderSharing  *types.CustomBool `json:"foldersharing,omitempty"  url:"foldersharing,omitempty"`
	VideoStreaming *string           `json:"videostreaming,omitempty" url:"videostreaming,omitempty"`
}

// EncodeValues converts a CustomSpiceEnhancements struct to a URL value.
func (r *CustomSpiceEnhancements) EncodeValues(key string, v *url.Values) error {
	var values []string

	if r.FolderSharing != nil {
		if *r.FolderSharing {
			values = append(values, "foldersharing=1")
		} else {
			values = append(values, "foldersharing=0")
		}
	}

	if r.VideoStreaming != nil {
		values = append(values, fmt.Sprintf("videostreaming=%s", *r.VideoStreaming))
	}

	if len(values) > 0 {
		v.Add(key, strings.Join(values, ","))
	}

	return nil
}

// UnmarshalJSON converts JSON to a CustomSpiceEnhancements struct.
func (r *CustomSpiceEnhancements) UnmarshalJSON(b []byte) error {
	var s string

	if err := json.Unmarshal(b, &s); err != nil {
		return fmt.Errorf("failed to unmarshal CustomSpiceEnhancements: %w", err)
	}

	pairs := strings.Split(s, ",")

	for _, p := range pairs {
		v := strings.Split(strings.TrimSpace(p), "=")

		if len(v) == 2 {
			switch v[0] {
			case "foldersharing":
				v := types.CustomBool(v[1] == "1")
				r.FolderSharing = &v
			case "videostreaming":
				r.VideoStreaming = &v[1]
			}
		}
	}

	return nil
}
