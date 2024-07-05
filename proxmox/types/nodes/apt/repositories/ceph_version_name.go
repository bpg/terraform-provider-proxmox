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

const (
	// CephStandardRepoHandlePrefix is the prefix for Ceph APT standard repositories.
	CephStandardRepoHandlePrefix = "ceph"
)

// Do not modify any of these package-global variables as they act as safer variants compared to "iota" based constants!
//
//nolint:gochecknoglobals
var (
	// CephVersionNameQuincy is the name for the "Quincy" Ceph major version.
	CephVersionNameQuincy = CephVersionName{"quincy"}

	// CephVersionNameReef is the name for the "Reef" Ceph major version.
	CephVersionNameReef = CephVersionName{"reef"}

	// CephVersionNameUnknown is the name for an unknown Ceph major version.
	CephVersionNameUnknown = CephVersionName{"unknown"}
)

// Ensure the hardware mapping type supports required interfaces.
var (
	_ fmt.Stringer     = new(CephVersionName)
	_ json.Marshaler   = new(CephVersionName)
	_ json.Unmarshaler = new(CephVersionName)
	_ query.Encoder    = new(CephVersionName)
)

// CephVersionName is the name a Ceph major version.
type CephVersionName struct {
	name string
}

// EncodeValues encodes Ceph major version name field into a URL-encoded set of values.
func (n CephVersionName) EncodeValues(key string, v *url.Values) error {
	v.Add(key, n.String())
	return nil
}

// MarshalJSON marshals a Ceph major version name into JSON value.
func (n CephVersionName) MarshalJSON() ([]byte, error) {
	bytes, err := json.Marshal(n.String())
	if err != nil {
		return nil, errors.Join(ErrCephVersionNameMarshal, err)
	}

	return bytes, nil
}

// String converts a CephVersionName value into a string.
func (n CephVersionName) String() string {
	return n.name
}

// ToValue converts a Ceph major version name into a Terraform value.
func (n CephVersionName) ToValue() types.String {
	return types.StringValue(n.String())
}

// UnmarshalJSON unmarshals a Ceph major version name.
func (n *CephVersionName) UnmarshalJSON(b []byte) error {
	var rtString string

	err := json.Unmarshal(b, &rtString)
	if err != nil {
		return errors.Join(ErrCephVersionNameUnmarshal, err)
	}

	resType, err := ParseCephVersionName(rtString)
	if err == nil {
		*n = resType
	}

	return err
}

// ParseCephVersionName converts the string representation of a Ceph major version name into the corresponding value.
// An error is returned if the input string does not match any known type.
func ParseCephVersionName(input string) (CephVersionName, error) {
	switch input {
	case CephVersionNameQuincy.String():
		return CephVersionNameQuincy, nil
	case CephVersionNameReef.String():
		return CephVersionNameReef, nil
	default:
		return CephVersionName{}, ErrCephVersionNameIllegal(input)
	}
}
