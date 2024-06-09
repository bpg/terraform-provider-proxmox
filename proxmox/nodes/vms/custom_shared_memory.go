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
	"strconv"
	"strings"
)

// CustomSharedMemory handles QEMU Inter-VM shared memory parameters.
type CustomSharedMemory struct {
	Name *string `json:"name,omitempty" url:"name,omitempty"`
	Size int     `json:"size"           url:"size"`
}

// EncodeValues converts a CustomSharedMemory struct to a URL value.
func (r *CustomSharedMemory) EncodeValues(key string, v *url.Values) error {
	values := []string{
		fmt.Sprintf("size=%d", r.Size),
	}

	if r.Name != nil {
		values = append(values, fmt.Sprintf("name=%s", *r.Name))
	}

	v.Add(key, strings.Join(values, ","))

	return nil
}

// UnmarshalJSON converts a CustomSharedMemory string to an object.
func (r *CustomSharedMemory) UnmarshalJSON(b []byte) error {
	var s string

	if err := json.Unmarshal(b, &s); err != nil {
		return fmt.Errorf("failed to unmarshal CustomSharedMemory: %w", err)
	}

	pairs := strings.Split(s, ",")

	for _, p := range pairs {
		v := strings.Split(strings.TrimSpace(p), "=")

		if len(v) == 2 {
			switch v[0] {
			case "name":
				r.Name = &v[1]
			case "size":
				var err error

				r.Size, err = strconv.Atoi(v[1])
				if err != nil {
					return fmt.Errorf("failed to parse shared memory size: %w", err)
				}
			}
		}
	}

	return nil
}
