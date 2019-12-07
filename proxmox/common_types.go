/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmox

import "bytes"

// CustomBool allows a JSON boolean value to also be an integer
type CustomBool bool

// UnmarshalJSON converts a JSON value to a boolean.
func (r *CustomBool) UnmarshalJSON(b []byte) error {
	s := string(b)
	*r = CustomBool(s == "1" || s == "true")

	return nil
}

// MarshalJSON converts a boolean to a JSON value.
func (r CustomBool) MarshalJSON() ([]byte, error) {
	buffer := new(bytes.Buffer)

	if r {
		buffer.WriteString("1")
	} else {
		buffer.WriteString("0")
	}

	return buffer.Bytes(), nil
}
