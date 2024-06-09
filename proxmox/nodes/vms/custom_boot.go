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

// CustomBoot handles QEMU boot parameters.
type CustomBoot struct {
	Order *[]string `json:"order,omitempty" url:"order,omitempty,semicolon"`
}

// EncodeValues converts a CustomBoot struct to multiple URL values.
func (r *CustomBoot) EncodeValues(key string, v *url.Values) error {
	if r.Order != nil && len(*r.Order) > 0 {
		v.Add(key, fmt.Sprintf("order=%s", strings.Join(*r.Order, ";")))
	}

	return nil
}

// UnmarshalJSON converts a CustomBoot string to an object.
func (r *CustomBoot) UnmarshalJSON(b []byte) error {
	var s string

	if err := json.Unmarshal(b, &s); err != nil {
		return fmt.Errorf("error unmarshalling CustomBoot: %w", err)
	}

	pairs := strings.Split(s, ",")

	for _, p := range pairs {
		v := strings.Split(strings.TrimSpace(p), "=")

		if len(v) == 2 {
			if v[0] == "order" {
				o := strings.Split(strings.TrimSpace(v[1]), ";")
				r.Order = &o
			}
		}
	}

	return nil
}
