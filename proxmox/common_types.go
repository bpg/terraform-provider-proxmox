/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmox

import (
	"bytes"
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

// CustomBool allows a JSON boolean value to also be an integer.
type CustomBool bool

// CustomCommaSeparatedList allows a JSON string to also be a string array.
type CustomCommaSeparatedList []string

// CustomLineBreakSeparatedList allows a multiline JSON string to also be a string array.
type CustomLineBreakSeparatedList []string

// CustomPrivileges allows a JSON object of privileges to also be a string array.
type CustomPrivileges []string

// CustomTimestamp allows a JSON boolean value to also be a unix timestamp.
type CustomTimestamp time.Time

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

// UnmarshalJSON converts a JSON value to a boolean.
func (r *CustomBool) UnmarshalJSON(b []byte) error {
	s := string(b)
	*r = CustomBool(s == "1" || s == "true")

	return nil
}

// MarshalJSON converts a boolean to a JSON value.
func (r *CustomCommaSeparatedList) MarshalJSON() ([]byte, error) {
	s := strings.Join(*r, ",")

	return json.Marshal(s)
}

// UnmarshalJSON converts a JSON value to a boolean.
func (r *CustomCommaSeparatedList) UnmarshalJSON(b []byte) error {
	var s string

	err := json.Unmarshal(b, &s)

	if err != nil {
		return err
	}

	*r = strings.Split(s, ",")

	return nil
}

// MarshalJSON converts a boolean to a JSON value.
func (r *CustomLineBreakSeparatedList) MarshalJSON() ([]byte, error) {
	s := strings.Join(*r, "\n")

	return json.Marshal(s)
}

// UnmarshalJSON converts a JSON value to a boolean.
func (r *CustomLineBreakSeparatedList) UnmarshalJSON(b []byte) error {
	var s string

	err := json.Unmarshal(b, &s)

	if err != nil {
		return err
	}

	*r = strings.Split(s, "\n")

	return nil
}

// MarshalJSON converts a boolean to a JSON value.
func (r *CustomPrivileges) MarshalJSON() ([]byte, error) {
	var privileges map[string]CustomBool

	for _, v := range *r {
		privileges[v] = true
	}

	return json.Marshal(privileges)
}

// UnmarshalJSON converts a JSON value to a boolean.
func (r *CustomPrivileges) UnmarshalJSON(b []byte) error {
	var privileges interface{}

	err := json.Unmarshal(b, &privileges)

	if err != nil {
		return err
	}

	switch privileges.(type) {
	case string:
		s := privileges.(string)

		if s != "" {
			*r = strings.Split(s, ",")
		} else {
			*r = make(CustomPrivileges, 0)
		}
	default:
		*r = make(CustomPrivileges, 0)

		for k, v := range privileges.(map[string]interface{}) {
			if v.(float64) >= 1 {
				*r = append(*r, k)
			}
		}
	}

	return nil
}

// MarshalJSON converts a boolean to a JSON value.
func (r CustomTimestamp) MarshalJSON() ([]byte, error) {
	var timestamp time.Time

	timestamp = time.Time(r)
	buffer := bytes.NewBufferString(strconv.FormatInt(timestamp.Unix(), 10))

	return buffer.Bytes(), nil
}

// UnmarshalJSON converts a JSON value to a boolean.
func (r *CustomTimestamp) UnmarshalJSON(b []byte) error {
	s := string(b)
	i, err := strconv.ParseInt(s, 10, 64)

	if err != nil {
		return err
	}

	*r = CustomTimestamp(time.Unix(i, 0))

	return nil
}
