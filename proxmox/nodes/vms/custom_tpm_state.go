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
)

// CustomTPMState handles QEMU TPM state parameters.
type CustomTPMState struct {
	FileVolume string  `json:"file"              url:"file"`
	Version    *string `json:"version,omitempty" url:"version,omitempty"`
}

// EncodeValues converts a CustomTPMState struct to a URL value.
func (r *CustomTPMState) EncodeValues(key string, v *url.Values) error {
	values := []string{
		fmt.Sprintf("file=%s", r.FileVolume),
	}

	if r.Version != nil {
		values = append(values, fmt.Sprintf("version=%s", *r.Version))
	}

	v.Add(key, strings.Join(values, ","))

	return nil
}

// UnmarshalJSON converts a CustomTPMState string to an object.
func (r *CustomTPMState) UnmarshalJSON(b []byte) error {
	var s string

	if err := json.Unmarshal(b, &s); err != nil {
		return fmt.Errorf("failed to unmarshal CustomTPMState: %w", err)
	}

	pairs := strings.Split(s, ",")

	for _, p := range pairs {
		v := strings.Split(strings.TrimSpace(p), "=")
		if len(v) == 1 {
			r.FileVolume = v[0]
		} else if len(v) == 2 {
			switch v[0] {
			case "file":
				r.FileVolume = v[1]
			case "version":
				r.Version = &v[1]
			}
		}
	}

	return nil
}
