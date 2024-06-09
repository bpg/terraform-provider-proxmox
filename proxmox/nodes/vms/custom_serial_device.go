/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package vms

import (
	"fmt"
	"net/url"
)

// CustomSerialDevices handles QEMU serial device parameters.
type CustomSerialDevices []string

// EncodeValues converts a CustomSerialDevices array to multiple URL values.
func (r CustomSerialDevices) EncodeValues(key string, v *url.Values) error {
	for i, d := range r {
		v.Add(fmt.Sprintf("%s%d", key, i), d)
	}

	return nil
}
