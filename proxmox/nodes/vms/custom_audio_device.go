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

// CustomAudioDevice handles QEMU audio parameters.
type CustomAudioDevice struct {
	Device  string  `json:"device" url:"device"`
	Driver  *string `json:"driver" url:"driver"`
	Enabled bool    `json:"-"      url:"-"`
}

// CustomAudioDevices handles QEMU audio device parameters.
type CustomAudioDevices []CustomAudioDevice

// EncodeValues converts a CustomAudioDevice struct to a URL value.
func (r *CustomAudioDevice) EncodeValues(key string, v *url.Values) error {
	values := []string{fmt.Sprintf("device=%s", r.Device)}

	if r.Driver != nil {
		values = append(values, fmt.Sprintf("driver=%s", *r.Driver))
	}

	v.Add(key, strings.Join(values, ","))

	return nil
}

// UnmarshalJSON converts a CustomAudioDevice string to an object.
func (r *CustomAudioDevice) UnmarshalJSON(b []byte) error {
	var s string

	if err := json.Unmarshal(b, &s); err != nil {
		return fmt.Errorf("error unmarshalling CustomAudioDevice: %w", err)
	}

	pairs := strings.Split(s, ",")

	for _, p := range pairs {
		v := strings.Split(strings.TrimSpace(p), "=")

		if len(v) == 2 {
			switch v[0] {
			case "device":
				r.Device = v[1]
			case "driver":
				r.Driver = &v[1]
			}
		}
	}

	return nil
}

// EncodeValues converts a CustomAudioDevices array to multiple URL values.
func (r CustomAudioDevices) EncodeValues(key string, v *url.Values) error {
	for i, d := range r {
		if d.Enabled {
			if err := d.EncodeValues(fmt.Sprintf("%s%d", key, i), v); err != nil {
				return fmt.Errorf("unable to encode audio device %d: %w", i, err)
			}
		}
	}

	return nil
}
