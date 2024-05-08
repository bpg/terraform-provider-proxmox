/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package hardwaremapping

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/google/go-querystring/query"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ID represents a hardware mapping ID, composed of the type and identifier.
type ID struct {
	// Name is the name of the hardware mapping.
	Name string

	// Type is the type of the hardware mapping.
	Type Type
}

// Ensure the hardware mapping ID type implements required interfaces.
var (
	_ fmt.Stringer     = &ID{}
	_ json.Marshaler   = &ID{}
	_ json.Unmarshaler = &ID{}
	_ query.Encoder    = &ID{}
)

// EncodeValues encodes a hardware mapping ID field into a URL-encoded set of values.
func (hmid ID) EncodeValues(key string, v *url.Values) error {
	v.Add(key, hmid.String())
	return nil
}

// MarshalJSON marshals a hardware mapping ID into JSON value.
func (hmid ID) MarshalJSON() ([]byte, error) {
	bytes, err := json.Marshal(hmid.String())
	if err != nil {
		return nil, errors.Join(ErrIDMarshal, err)
	}

	return bytes, nil
}

// String converts an ID value into a string.
func (hmid ID) String() string {
	return fmt.Sprintf("%s:%s", hmid.Type, hmid.Name)
}

// ToValue converts a hardware mapping ID into a Terraform value.
func (hmid ID) ToValue() types.String {
	return types.StringValue(hmid.String())
}

// UnmarshalJSON unmarshals a hardware mapping ID.
func (hmid *ID) UnmarshalJSON(b []byte) error {
	var hmIDString string

	err := json.Unmarshal(b, &hmIDString)
	if err != nil {
		return errors.Join(ErrTypeUnmarshal, err)
	}

	hmID, err := ParseID(hmIDString)
	if err == nil {
		*hmid = hmID
	}

	return err
}

// ParseID parses a string that represents a hardware mapping ID into a value of `ID`.
func ParseID(input string) (ID, error) {
	hmID := ID{}

	inParts := strings.SplitN(input, ":", 2)
	if len(inParts) < 2 {
		return hmID, ErrIDParsing(input)
	}

	hmType, err := ParseType(inParts[0])
	if err != nil {
		return hmID, errors.Join(fmt.Errorf("could not extract type from hardware mapping ID %q", input), err)
	}

	hmID.Type = hmType
	hmID.Name = inParts[1]

	return hmID, nil
}
