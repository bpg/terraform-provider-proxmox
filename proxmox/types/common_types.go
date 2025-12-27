/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package types

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/proxmox/helpers/ptr"
)

// CustomBool allows a JSON boolean value to also be an integer.
type CustomBool bool

// CustomCommaSeparatedList allows a JSON string to also be a string array.
type CustomCommaSeparatedList []string

// CustomFloat64 allows a JSON float64 value to also be a string.
type CustomFloat64 float64

// CustomInt allows a JSON integer value to also be a string.
type CustomInt int

// CustomInt64 allows a JSON int64 value to also be a string.
type CustomInt64 int64

// CustomLineBreakSeparatedList allows a multiline JSON string to also be a string array.
type CustomLineBreakSeparatedList []string

// CustomPrivileges allows a JSON object of privileges to also be a string array.
type CustomPrivileges []string

// CustomTimestamp allows a JSON integer value to also be a unix timestamp.
type CustomTimestamp time.Time

// CustomBoolPtr creates a pointer to a CustomBool.
func CustomBoolPtr(b *bool) *CustomBool {
	if b == nil {
		return nil
	}

	return ptr.Ptr(CustomBool(*b))
}

// MarshalJSON converts a CustomBool to a JSON integer (1 or 0).
func (r CustomBool) MarshalJSON() ([]byte, error) {
	buffer := new(bytes.Buffer)

	if r {
		buffer.WriteString("1")
	} else {
		buffer.WriteString("0")
	}

	return buffer.Bytes(), nil
}

// UnmarshalJSON converts a JSON value to a CustomBool.
func (r *CustomBool) UnmarshalJSON(b []byte) error {
	s := string(b)
	*r = s == "1" || s == "true"

	return nil
}

// Pointer returns a pointers.
func (r CustomBool) Pointer() *CustomBool {
	return &r
}

// PointerBool returns a pointer to a boolean.
func (r *CustomBool) PointerBool() *bool {
	return (*bool)(r)
}

// ToValue returns a Terraform attribute value.
func (r CustomBool) ToValue() types.Bool {
	return types.BoolValue(bool(r))
}

// FromValue sets the numeric boolean based on the value of a Terraform attribute.
func (r *CustomBool) FromValue(tfValue types.Bool) {
	*r = CustomBool(tfValue.ValueBool())
}

// MarshalJSON converts a CustomCommaSeparatedList to a JSON string.
func (r *CustomCommaSeparatedList) MarshalJSON() ([]byte, error) {
	s := strings.Join(*r, ",")

	b, err := json.Marshal(s)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal CustomCommaSeparatedList: %w", err)
	}

	return b, nil
}

// UnmarshalJSON converts a JSON value to a CustomCommaSeparatedList.
func (r *CustomCommaSeparatedList) UnmarshalJSON(b []byte) error {
	var s string

	err := json.Unmarshal(b, &s)
	if err != nil {
		return fmt.Errorf("failed to unmarshal CustomCommaSeparatedList: %w", err)
	}

	*r = strings.Split(s, ",")

	return nil
}

// UnmarshalJSON converts a JSON value to a float64 value.
func (r *CustomFloat64) UnmarshalJSON(b []byte) error {
	s := string(b)

	if strings.HasPrefix(s, "\"") && strings.HasSuffix(s, "\"") {
		s = s[1 : len(s)-1]
	}

	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return fmt.Errorf("cannot parse float64 %q: %w", s, err)
	}

	*r = CustomFloat64(f)

	return nil
}

// PointerFloat64 returns a pointer to a float64.
func (r *CustomFloat64) PointerFloat64() *float64 {
	return (*float64)(r)
}

// UnmarshalJSON converts a JSON value to an integer.
func (r *CustomInt) UnmarshalJSON(b []byte) error {
	s := string(b)

	if len(s) > 0 && s[0] == '"' {
		var val string
		if err := json.Unmarshal(b, &val); err != nil {
			return fmt.Errorf("cannot unmarshal string for CustomInt: %w", err)
		}

		s = val
		if s == "" {
			return fmt.Errorf("cannot parse int from empty string")
		}
	}

	// Try plain integer first
	if i, err := strconv.ParseInt(s, 10, 0); err == nil {
		*r = CustomInt(i)
		return nil
	}

	// Fall back to float/scientific notation, then convert to int
	if f, _, err := new(big.Float).Parse(s, 10); err == nil {
		i := new(big.Int)
		f.Int(i) // Truncates toward zero, similar to int64(float64_val)

		if i.IsInt64() {
			val := i.Int64()
			if val >= math.MinInt && val <= math.MaxInt {
				*r = CustomInt(val)
				return nil
			}
		}

		return fmt.Errorf("cannot parse int %q: value out of range", s)
	}

	return fmt.Errorf("cannot parse int %q: unsupported numeric format", s)
}

// UnmarshalJSON converts a JSON value to an integer64.
func (r *CustomInt64) UnmarshalJSON(b []byte) error {
	s := string(b)

	if len(s) > 0 && s[0] == '"' {
		var val string
		if err := json.Unmarshal(b, &val); err != nil {
			return fmt.Errorf("cannot unmarshal string for CustomInt64: %w", err)
		}

		s = val
		if s == "" {
			return fmt.Errorf("cannot parse int64 from empty string")
		}
	}

	// First, try parsing as a plain integer
	if i, err := strconv.ParseInt(s, 10, 64); err == nil {
		*r = CustomInt64(i)
		return nil
	}

	// If that fails, numbers may be provided in scientific notation or as floats.
	// Parse as a high-precision float and convert to int64 (truncating towards zero).
	if f, _, err := new(big.Float).Parse(s, 10); err == nil {
		i := new(big.Int)
		f.Int(i) // Truncates toward zero, similar to int64(float64_val)

		if i.IsInt64() {
			*r = CustomInt64(i.Int64())
			return nil
		}

		return fmt.Errorf("cannot parse int64 %q: value out of range", s)
	}

	return fmt.Errorf("cannot parse int64 %q: unsupported numeric format", s)
}

// PointerInt64 returns a pointer to an int64.
func (r *CustomInt64) PointerInt64() *int64 {
	return (*int64)(r)
}

// MarshalJSON converts a CustomLineBreakSeparatedList to a JSON string.
func (r *CustomLineBreakSeparatedList) MarshalJSON() ([]byte, error) {
	s := strings.Join(*r, "\n")

	b, err := json.Marshal(s)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal CustomLineBreakSeparatedList: %w", err)
	}

	return b, nil
}

// UnmarshalJSON converts a JSON value to a CustomLineBreakSeparatedList.
func (r *CustomLineBreakSeparatedList) UnmarshalJSON(b []byte) error {
	var s string

	err := json.Unmarshal(b, &s)
	if err != nil {
		return fmt.Errorf("failed to unmarshal CustomLineBreakSeparatedList: %w", err)
	}

	*r = strings.Split(s, "\n")

	return nil
}

// MarshalJSON converts a CustomPrivileges to a JSON object.
func (r *CustomPrivileges) MarshalJSON() ([]byte, error) {
	privileges := map[string]CustomBool{}

	for _, v := range *r {
		privileges[v] = true
	}

	b, err := json.Marshal(privileges)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal CustomPrivileges: %w", err)
	}

	return b, nil
}

// UnmarshalJSON converts a JSON value to a CustomPrivileges.
func (r *CustomPrivileges) UnmarshalJSON(b []byte) error {
	var privileges any

	err := json.Unmarshal(b, &privileges)
	if err != nil {
		return fmt.Errorf("failed to unmarshal CustomPrivileges: %w", err)
	}

	switch s := privileges.(type) {
	case string:
		if s != "" {
			*r = strings.Split(s, ",")
		} else {
			*r = CustomPrivileges{}
		}
	default:
		*r = CustomPrivileges{}

		for k, v := range privileges.(map[string]any) {
			if v.(float64) >= 1 {
				*r = append(*r, k)
			}
		}
	}

	return nil
}

// MarshalJSON converts a timestamp to a JSON value.
func (r CustomTimestamp) MarshalJSON() ([]byte, error) {
	timestamp := time.Time(r)
	buffer := bytes.NewBufferString(strconv.FormatInt(timestamp.Unix(), 10))

	return buffer.Bytes(), nil
}

// UnmarshalJSON converts a JSON value to a timestamp.
func (r *CustomTimestamp) UnmarshalJSON(b []byte) error {
	s := string(b)

	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse timestamp: %w", err)
	}

	*r = CustomTimestamp(time.Unix(i, 0).UTC())

	return nil
}
