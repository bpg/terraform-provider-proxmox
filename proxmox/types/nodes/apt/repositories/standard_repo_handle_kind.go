/*
	This Source Code Form is subject to the terms of the Mozilla Public
	License, v. 2.0. If a copy of the MPL was not distributed with this
	file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/

package repositories

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"

	"github.com/google/go-querystring/query"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Do not modify any of these package-global variables as they act as safer variants compared to "iota" based constants!
//
//nolint:gochecknoglobals
var (
	// StandardRepoHandleKindEnterprise is the name for the "Enterprise" APT standard repository handle kind.
	StandardRepoHandleKindEnterprise = StandardRepoHandleKind{"enterprise"}

	// StandardRepoHandleKindNoSubscription is the name for the "No Subscription" APT standard repository handle kind.
	StandardRepoHandleKindNoSubscription = StandardRepoHandleKind{"no-subscription"}

	// StandardRepoHandleKindTest is the name for the "Test" APT standard repository handle kind.
	StandardRepoHandleKindTest = StandardRepoHandleKind{"test"}

	// StandardRepoHandleKindUnknown is the name for an unknown APT standard repository handle kind.
	StandardRepoHandleKindUnknown = StandardRepoHandleKind{"unknown"}
)

// Ensure the hardware mapping type supports required interfaces.
var (
	_ fmt.Stringer     = new(StandardRepoHandleKind)
	_ json.Marshaler   = new(StandardRepoHandleKind)
	_ json.Unmarshaler = new(StandardRepoHandleKind)
	_ query.Encoder    = new(StandardRepoHandleKind)
)

// StandardRepoHandleKind is the kind of APT standard repository handle.
type StandardRepoHandleKind struct {
	handle string
}

// EncodeValues encodes the APT standard repository handle kind field into a URL-encoded set of values.
func (h *StandardRepoHandleKind) EncodeValues(key string, v *url.Values) error {
	v.Add(key, h.String())
	return nil
}

// MarshalJSON marshals an APT standard repository handle kind into JSON value.
func (h *StandardRepoHandleKind) MarshalJSON() ([]byte, error) {
	bytes, err := json.Marshal(h.String())
	if err != nil {
		return nil, errors.Join(ErrStandardRepoHandleKindMarshal, err)
	}

	return bytes, nil
}

// String converts a StandardRepoHandleKind value into a string.
func (h StandardRepoHandleKind) String() string {
	return h.handle
}

// ToValue converts an APT standard repository handle kind into a Terraform value.
func (h *StandardRepoHandleKind) ToValue() types.String {
	return types.StringValue(h.String())
}

// UnmarshalJSON unmarshals an APT standard repository handle kind.
func (h *StandardRepoHandleKind) UnmarshalJSON(b []byte) error {
	var rtString string

	err := json.Unmarshal(b, &rtString)
	if err != nil {
		return errors.Join(ErrStandardRepoHandleKindUnmarshal, err)
	}

	resType, err := ParseStandardRepoHandleKind(rtString)
	if err == nil {
		*h = resType
	}

	return err
}

// ParseStandardRepoHandleKind converts the string representation of an APT standard repository handle kind into the
// corresponding type.
// StandardRepoHandleKindUnknown and an error is returned if the input string does not match any known handle kind.
func ParseStandardRepoHandleKind(input string) (StandardRepoHandleKind, error) {
	switch input {
	case StandardRepoHandleKindEnterprise.String():
		return StandardRepoHandleKindEnterprise, nil
	case StandardRepoHandleKindNoSubscription.String():
		return StandardRepoHandleKindNoSubscription, nil
	case StandardRepoHandleKindTest.String():
		return StandardRepoHandleKindTest, nil
	}

	return StandardRepoHandleKindUnknown, fmt.Errorf(
		"parse APT standard repository handle kind: %w",
		ErrStandardRepoHandleKindIllegal(input),
	)
}
