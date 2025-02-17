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

// CustomRNGDevice represents a random number generator device configuration.
type CustomRNGDevice struct {
	Source   string `json:"source,omitempty"    url:"source,omitempty"`
	MaxBytes *int   `json:"max_bytes,omitempty" url:"max_bytes,omitempty"`
	Period   *int   `json:"period,omitempty"    url:"period,omitempty"`
}

// EncodeValues converts a CustomRNGDevice struct to a URL value.
func (r *CustomRNGDevice) EncodeValues(key string, v *url.Values) error {
	var values []string

	if r.Source != "" {
		values = append(values, fmt.Sprintf("source=%s", r.Source))
	}

	if r.MaxBytes != nil {
		values = append(values, fmt.Sprintf("max_bytes=%d", *r.MaxBytes))
	}

	if r.Period != nil {
		values = append(values, fmt.Sprintf("period=%d", *r.Period))
	}

	if len(values) > 0 {
		v.Add(key, strings.Join(values, ","))
	}

	return nil
}

// UnmarshalJSON unmarshals a JSON object into a CustomRNGDevice struct.
func (r *CustomRNGDevice) UnmarshalJSON(b []byte) error {
	var s string

	if err := json.Unmarshal(b, &s); err != nil {
		return fmt.Errorf("failed to unmarshal CustomRNGDevice: %w", err)
	}

	pairs := strings.Split(s, ",")

	for _, p := range pairs {
		v := strings.Split(strings.TrimSpace(p), "=")
		if len(v) == 1 {
			r.Source = v[0]
		} else if len(v) == 2 {
			switch v[0] {
			case "source":
				r.Source = v[1]

			case "max_bytes":
				maxBytes, err := strconv.Atoi(v[1])
				if err != nil {
					return fmt.Errorf("failed to parse max_bytes: %w", err)
				}

				r.MaxBytes = &maxBytes

			case "period":
				period, err := strconv.Atoi(v[1])
				if err != nil {
					return fmt.Errorf("failed to parse period: %w", err)
				}

				r.Period = &period
			}
		}
	}

	return nil
}
