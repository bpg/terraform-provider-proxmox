/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/google/go-querystring/query"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

//nolint:gochecknoglobals
var (
	// HardwareMappingIDErrMarshal indicates an error while marshalling a hardware mapping ID.
	HardwareMappingIDErrMarshal = function.NewFuncError("cannot unmarshal hardware mapping type")

	// HardwareMappingIDErrParsing indicates an error while parsing a hardware mapping ID.
	HardwareMappingIDErrParsing = func(hmID string) error {
		return function.NewFuncError(fmt.Sprintf("%q is not a valid hardware mapping ID", hmID))
	}
)

// HardwareMappingID represents a hardware mapping ID, composed of the type and identifier.
type HardwareMappingID struct {
	// Name is the name of the hardware mapping.
	Name string

	// Type is the type of the hardware mapping.
	Type HardwareMappingType
}

// Ensure the hardware mapping ID type implements required interfaces.
var (
	_ fmt.Stringer     = &HardwareMappingID{}
	_ json.Marshaler   = &HardwareMappingID{}
	_ json.Unmarshaler = &HardwareMappingID{}
	_ query.Encoder    = &HardwareMappingID{}
)

// EncodeValues encodes a hardware mapping ID field into a URL-encoded set of values.
func (hmid HardwareMappingID) EncodeValues(key string, v *url.Values) error {
	v.Add(key, hmid.String())
	return nil
}

// MarshalJSON marshals a hardware mapping ID into JSON value.
func (hmid HardwareMappingID) MarshalJSON() ([]byte, error) {
	bytes, err := json.Marshal(hmid.String())
	if err != nil {
		return nil, errors.Join(HardwareMappingIDErrMarshal, err)
	}

	return bytes, nil
}

// String converts a HardwareMappingID value into a string.
func (hmid HardwareMappingID) String() string {
	return fmt.Sprintf("%s:%s", hmid.Type, hmid.Name)
}

// ToValue converts a hardware mapping ID into a Terraform value.
func (hmid HardwareMappingID) ToValue() types.String {
	return types.StringValue(hmid.String())
}

// UnmarshalJSON unmarshals a hardware mapping ID.
func (hmid *HardwareMappingID) UnmarshalJSON(b []byte) error {
	var hmIDString string

	err := json.Unmarshal(b, &hmIDString)
	if err != nil {
		return errors.Join(HardwareMappingTypeErrUnmarshal, err)
	}

	hmID, err := ParseHardwareMappingID(hmIDString)
	if err == nil {
		*hmid = hmID
	}

	return err
}

// ParseHardwareMappingID parses a string that represents a hardware mapping ID into a value of `HardwareMappingID`.
func ParseHardwareMappingID(input string) (HardwareMappingID, error) {
	hmID := HardwareMappingID{}

	inParts := strings.SplitN(input, ":", 2)
	if len(inParts) < 2 {
		return hmID, HardwareMappingIDErrParsing(input)
	}

	hmType, err := ParseHardwareMappingType(inParts[0])
	if err != nil {
		return hmID, errors.Join(fmt.Errorf("could not extract type from hardware mapping ID %q", input), err)
	}

	hmID.Type = hmType
	hmID.Name = inParts[1]

	return hmID, nil
}
