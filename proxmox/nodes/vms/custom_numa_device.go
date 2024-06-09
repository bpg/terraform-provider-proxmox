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

// CustomNUMADevice handles QEMU NUMA device parameters.
type CustomNUMADevice struct {
	CPUIDs        []string  `json:"cpus"                url:"cpus,semicolon"`
	HostNodeNames *[]string `json:"hostnodes,omitempty" url:"hostnodes,omitempty,semicolon"`
	Memory        *int      `json:"memory,omitempty"    url:"memory,omitempty"`
	Policy        *string   `json:"policy,omitempty"    url:"policy,omitempty"`
}

// CustomNUMADevices handles QEMU NUMA device parameters.
type CustomNUMADevices []CustomNUMADevice

// EncodeValues converts a CustomNUMADevice struct to a URL value.
func (r *CustomNUMADevice) EncodeValues(key string, v *url.Values) error {
	values := []string{
		fmt.Sprintf("cpus=%s", strings.Join(r.CPUIDs, ";")),
	}

	if r.HostNodeNames != nil {
		values = append(values, fmt.Sprintf("hostnodes=%s", strings.Join(*r.HostNodeNames, ";")))
	}

	if r.Memory != nil {
		values = append(values, fmt.Sprintf("memory=%d", *r.Memory))
	}

	if r.Policy != nil {
		values = append(values, fmt.Sprintf("policy=%s", *r.Policy))
	}

	v.Add(key, strings.Join(values, ","))

	return nil
}

// EncodeValues converts a CustomNUMADevices array to multiple URL values.
func (r CustomNUMADevices) EncodeValues(key string, v *url.Values) error {
	for i, d := range r {
		if err := d.EncodeValues(fmt.Sprintf("%s%d", key, i), v); err != nil {
			return fmt.Errorf("failed to encode NUMA device %d: %w", i, err)
		}
	}

	return nil
}

// UnmarshalJSON converts a CustomNUMADevice string to an object.
func (r *CustomNUMADevice) UnmarshalJSON(b []byte) error {
	var s string

	if err := json.Unmarshal(b, &s); err != nil {
		return fmt.Errorf("failed to unmarshal CustomNUMADevice: %w", err)
	}

	pairs := strings.Split(s, ",")

	for _, p := range pairs {
		v := strings.Split(strings.TrimSpace(p), "=")
		if len(v) == 2 {
			switch v[0] {
			case "cpus":
				r.CPUIDs = strings.Split(v[1], ";")
			case "hostnodes":
				hostnodes := strings.Split(v[1], ";")
				r.HostNodeNames = &hostnodes
			case "memory":
				memory, err := strconv.Atoi(v[1])
				if err != nil {
					return fmt.Errorf("failed to parse memory size: %w", err)
				}

				r.Memory = &memory
			case "policy":
				r.Policy = &v[1]
			}
		}
	}

	return nil
}
