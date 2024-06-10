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

// CustomVGADevice handles QEMU VGA device parameters.
type CustomVGADevice struct {
	Clipboard *string `json:"clipboard,omitempty" url:"memory,omitempty"`
	Memory    *int64  `json:"memory,omitempty"    url:"memory,omitempty"`
	Type      *string `json:"type,omitempty"      url:"type,omitempty"`
}

// EncodeValues converts a CustomVGADevice struct to a URL value.
func (r *CustomVGADevice) EncodeValues(key string, v *url.Values) error {
	var values []string

	if r.Clipboard != nil {
		values = append(values, fmt.Sprintf("clipboard=%s", *r.Clipboard))
	}

	if r.Memory != nil {
		values = append(values, fmt.Sprintf("memory=%d", *r.Memory))
	}

	if r.Type != nil {
		values = append(values, fmt.Sprintf("type=%s", *r.Type))
	}

	v.Add(key, strings.Join(values, ","))

	return nil
}

// UnmarshalJSON converts a CustomVGADevice string to an object.
func (r *CustomVGADevice) UnmarshalJSON(b []byte) error {
	var s string

	if err := json.Unmarshal(b, &s); err != nil {
		return fmt.Errorf("failed to unmarshal CustomVGADevice: %w", err)
	}

	if s == "" {
		return nil
	}

	pairs := strings.Split(s, ",")

	for _, p := range pairs {
		v := strings.Split(strings.TrimSpace(p), "=")

		if len(v) == 1 {
			r.Type = &v[0]
		} else if len(v) == 2 {
			switch v[0] {
			case "clipboard":
				r.Clipboard = &v[1]

			case "memory":
				m, err := strconv.ParseInt(v[1], 10, 64)
				if err != nil {
					return fmt.Errorf("failed to convert memory to int: %w", err)
				}

				r.Memory = &m
			case "type":
				r.Type = &v[1]
			}
		}
	}

	return nil
}
