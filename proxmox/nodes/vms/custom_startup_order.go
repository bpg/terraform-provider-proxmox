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

// CustomStartupOrder handles QEMU startup order parameters.
type CustomStartupOrder struct {
	Down  *int `json:"down,omitempty"  url:"down,omitempty"`
	Order *int `json:"order,omitempty" url:"order,omitempty"`
	Up    *int `json:"up,omitempty"    url:"up,omitempty"`
}

// EncodeValues converts a CustomStartupOrder struct to a URL value.
func (r *CustomStartupOrder) EncodeValues(key string, v *url.Values) error {
	var values []string

	if r.Order != nil {
		values = append(values, fmt.Sprintf("order=%d", *r.Order))
	}

	if r.Up != nil {
		values = append(values, fmt.Sprintf("up=%d", *r.Up))
	}

	if r.Down != nil {
		values = append(values, fmt.Sprintf("down=%d", *r.Down))
	}

	if len(values) > 0 {
		v.Add(key, strings.Join(values, ","))
	}

	return nil
}

// UnmarshalJSON converts a CustomStartupOrder string to an object.
func (r *CustomStartupOrder) UnmarshalJSON(b []byte) error {
	var s string

	if err := json.Unmarshal(b, &s); err != nil {
		return fmt.Errorf("failed to unmarshal CustomStartupOrder: %w", err)
	}

	pairs := strings.Split(s, ",")

	for _, p := range pairs {
		v := strings.Split(strings.TrimSpace(p), "=")

		if len(v) == 2 {
			switch v[0] {
			case "order":
				order, err := strconv.Atoi(v[1])
				if err != nil {
					return fmt.Errorf("failed to parse int: %w", err)
				}

				r.Order = &order
			case "up":
				up, err := strconv.Atoi(v[1])
				if err != nil {
					return fmt.Errorf("failed to parse int: %w", err)
				}

				r.Up = &up
			case "down":
				down, err := strconv.Atoi(v[1])
				if err != nil {
					return fmt.Errorf("failed to parse int: %w", err)
				}

				r.Down = &down
			}
		}
	}

	return nil
}
