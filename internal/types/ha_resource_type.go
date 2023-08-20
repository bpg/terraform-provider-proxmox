/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package types

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/google/go-querystring/query"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/internal/validators"
)

// HAResourceType represents the type of a HA resource.
type HAResourceType int

// Ensure various interfaces are supported by the HA resource type type.
// NOTE: to my knowledge, this "global" here is required for the static type checks to work.
var (
	//nolint:gochecknoglobals
	_haResourceTypeValue HAResourceType
	_                    fmt.Stringer     = &_haResourceTypeValue
	_                    json.Marshaler   = &_haResourceTypeValue
	_                    json.Unmarshaler = &_haResourceTypeValue
	_                    query.Encoder    = &_haResourceTypeValue
)

const (
	// HAResourceTypeVM indicates that a HA resource refers to a virtual machine.
	HAResourceTypeVM HAResourceType = 0
	// HAResourceTypeContainer indicates that a HA resource refers to a container.
	HAResourceTypeContainer HAResourceType = 1
)

// ParseHAResourceType converts the string representation of a HA resource type into the corresponding
// enum value. An error is returned if the input string does not match any known type.
func ParseHAResourceType(input string) (HAResourceType, error) {
	switch input {
	case "vm":
		return HAResourceTypeVM, nil
	case "ct":
		return HAResourceTypeContainer, nil
	default:
		return _haResourceTypeValue, fmt.Errorf("illegal HA resource type '%s'", input)
	}
}

// HAResourceTypeValidator returns a new HA resource type validator.
func HAResourceTypeValidator() validator.String {
	return validators.NewParseValidator(ParseHAResourceType, "value must be a valid HA resource type")
}

// String converts a HAResourceType value into a string.
func (t HAResourceType) String() string {
	switch t {
	case HAResourceTypeVM:
		return "vm"
	case HAResourceTypeContainer:
		return "ct"
	default:
		panic(fmt.Sprintf("unknown HA resource type value: %d", t))
	}
}

// MarshalJSON marshals a HA resource type into JSON value.
func (t HAResourceType) MarshalJSON() ([]byte, error) {
	bytes, err := json.Marshal(t.String())
	if err != nil {
		return nil, fmt.Errorf("cannot marshal HA resource type: %w", err)
	}

	return bytes, nil
}

// UnmarshalJSON unmarshals a Proxmox HA resource type.
func (t *HAResourceType) UnmarshalJSON(b []byte) error {
	var rtString string

	err := json.Unmarshal(b, &rtString)
	if err != nil {
		return fmt.Errorf("cannot unmarshal HA resource type: %w", err)
	}

	resType, err := ParseHAResourceType(rtString)
	if err == nil {
		*t = resType
	}

	return err
}

// EncodeValues encodes a HA resource type field into an URL-encoded set of values.
func (t HAResourceType) EncodeValues(key string, v *url.Values) error {
	v.Add(key, t.String())
	return nil
}

// ToValue converts a HA resource type into a Terraform value.
func (t HAResourceType) ToValue() types.String {
	return types.StringValue(t.String())
}
