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
	"strconv"
	"strings"

	"github.com/google/go-querystring/query"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/internal/validators"
)

// NOTE: the linter believes the `HAResourceID` structure below should be tagged with `json:` due to some values of it
// being passed to a JSON marshaler in the tests. As far as I can tell this is unnecessary, so I'm silencing the lint.

// HAResourceID represents a HA resource identifier, composed of a resource type and identifier.
//
//nolint:musttag
type HAResourceID struct {
	Type HAResourceType // The type of this HA resource.
	Name string         // The name of the element the HA resource refers to.
}

// Ensure the HA resource identifier type implements various interfaces.
var (
	_ fmt.Stringer     = &HAResourceID{}
	_ json.Marshaler   = &HAResourceID{}
	_ json.Unmarshaler = &HAResourceID{}
	_ query.Encoder    = &HAResourceID{}
)

// ParseHAResourceID parses a string that represents a HA resource identifier into a value of `HAResourceID`.
func ParseHAResourceID(input string) (HAResourceID, error) {
	resID := HAResourceID{}

	inParts := strings.SplitN(input, ":", 2)
	if len(inParts) < 2 {
		return resID, fmt.Errorf("'%s' is not a valid HA resource identifier", input)
	}

	resType, err := ParseHAResourceType(inParts[0])
	if err != nil {
		return resID, fmt.Errorf("could not extract type from HA resource identifier '%s': %w", input, err)
	}

	// For types VM and Container, we know the resource "name" should be a valid integer between 100
	// and 999_999_999.
	if resType == HAResourceTypeVM || resType == HAResourceTypeContainer {
		id, err := strconv.Atoi(inParts[1])
		if err != nil {
			return resID, fmt.Errorf("invalid %s HA resource name '%s': %w", resType, inParts[1], err)
		}

		if id < 100 {
			return resID, fmt.Errorf("invalid %s HA resource name '%s': minimum value is 100", resType, inParts[1])
		}

		if id > 999_999_999 {
			return resID, fmt.Errorf("invalid %s HA resource name '%s': maximum value is 999999999", resType, inParts[1])
		}
	}

	resID.Type = resType
	resID.Name = inParts[1]

	return resID, nil
}

// HAResourceIDValidator returns a new HA resource identifier validator.
func HAResourceIDValidator() validator.String {
	return validators.NewParseValidator(ParseHAResourceID, "value must be a valid HA resource identifier")
}

// String converts a HAResourceID value into a string.
func (rid HAResourceID) String() string {
	return fmt.Sprintf("%s:%s", rid.Type, rid.Name)
}

// MarshalJSON marshals a HA resource identifier into JSON value.
func (rid HAResourceID) MarshalJSON() ([]byte, error) {
	bytes, err := json.Marshal(rid.String())
	if err != nil {
		return nil, fmt.Errorf("cannot marshal HA resource identifier: %w", err)
	}

	return bytes, nil
}

// UnmarshalJSON unmarshals a Proxmox HA resource identifier.
func (rid *HAResourceID) UnmarshalJSON(b []byte) error {
	var ridString string

	err := json.Unmarshal(b, &ridString)
	if err != nil {
		return fmt.Errorf("cannot unmarshal HA resource type: %w", err)
	}

	resType, err := ParseHAResourceID(ridString)
	if err == nil {
		*rid = resType
	}

	return err
}

// EncodeValues encodes a HA resource ID field into an URL-encoded set of values.
func (rid HAResourceID) EncodeValues(key string, v *url.Values) error {
	v.Add(key, rid.String())
	return nil
}

// ToValue converts a HA resource ID into a Terraform value.
func (rid HAResourceID) ToValue() types.String {
	return types.StringValue(rid.String())
}
