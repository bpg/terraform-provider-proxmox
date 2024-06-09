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

// CustomWatchdogDevice handles QEMU watchdog device parameters.
type CustomWatchdogDevice struct {
	Action *string `json:"action,omitempty" url:"action,omitempty"`
	Model  *string `json:"model"            url:"model"`
}

// EncodeValues converts a CustomWatchdogDevice struct to a URL value.
func (r *CustomWatchdogDevice) EncodeValues(key string, v *url.Values) error {
	values := []string{
		fmt.Sprintf("model=%+v", r.Model),
	}

	if r.Action != nil {
		values = append(values, fmt.Sprintf("action=%s", *r.Action))
	}

	v.Add(key, strings.Join(values, ","))

	return nil
}

// UnmarshalJSON converts a CustomWatchdogDevice string to an object.
func (r *CustomWatchdogDevice) UnmarshalJSON(b []byte) error {
	var s string

	if err := json.Unmarshal(b, &s); err != nil {
		return fmt.Errorf("failed to unmarshal CustomWatchdogDevice: %w", err)
	}

	if s == "" {
		return nil
	}

	pairs := strings.Split(s, ",")

	for _, p := range pairs {
		v := strings.Split(strings.TrimSpace(p), "=")

		if len(v) == 1 {
			r.Model = &v[0]
		} else if len(v) == 2 {
			switch v[0] {
			case "action":
				r.Action = &v[1]
			case "model":
				r.Model = &v[1]
			}
		}
	}

	return nil
}
