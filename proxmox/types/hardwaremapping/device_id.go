/*
	This Source Code Form is subject to the terms of the Mozilla Public
	License, v. 2.0. If a copy of the MPL was not distributed with this
	file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/

package hardwaremapping

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"regexp"

	"github.com/google/go-querystring/query"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	// attrNameDeviceID is the attribute name of the device ID in a hardware mapping.
	attrNameDeviceID = "id"

	// attrNameSubsystemID is the attribute name of the device subsystem ID in a hardware mapping.
	attrNameSubsystemID = "subsystem-id"
)

// DeviceIDAttrValueRegEx is the regular expression for device ID attribute value in a hardware mapping.
var DeviceIDAttrValueRegEx = regexp.MustCompile(`^[0-9A-Fa-f]{4}:[0-9A-Fa-f]{4}$`)

// Ensure the hardware mapping device ID type implements required interfaces.
var (
	_ fmt.Stringer     = new(DeviceID)
	_ json.Marshaler   = new(DeviceID)
	_ json.Unmarshaler = new(DeviceID)
	_ query.Encoder    = new(DeviceID)
)

// DeviceID represents a hardware mapping device ID.
// An ID is composed of two parts, either…
//   - a Vendor ID and device ID.
//     This is the device class and subclass (two 8-bit numbers).
//   - Subsystem ID and Subsystem device ID.
//     This identifies the assembly in which the device is contained.
//     Subsystems have their vendor ID (from the same namespace as device vendors) and subsystem ID.
//
// References:
//   - [Linux Kernel Documentation — PCI drivers]
//   - [Linux Hardware Database]
//   - [Linux USB ID Repository]
//   - [man(5) — pci.ids]
//
// [Linux Kernel Documentation — PCI drivers]: https://docs.kernel.org/admin-guide/media/pci-cardlist.html
// [Linux Hardware Database]: https://linux-hardware.org
// [Linux USB ID Repository]: http://www.linux-usb.org/usb-ids.html
// [man(5) — pci.ids]: https://man.archlinux.org/man/core/pciutils/pci.ids.5.en#INTRODUCTION
type DeviceID string

// EncodeValues encodes a hardware mapping device ID field into a URL-encoded set of values.
func (did DeviceID) EncodeValues(key string, v *url.Values) error {
	v.Add(key, did.String())

	return nil
}

// MarshalJSON marshals a hardware mapping device ID into JSON value.
func (did DeviceID) MarshalJSON() ([]byte, error) {
	bytes, err := json.Marshal(did)
	if err != nil {
		return nil, errors.Join(ErrDeviceIDMarshal, err)
	}

	return bytes, nil
}

// String converts a DeviceID value into a string.
func (did DeviceID) String() string {
	return string(did)
}

// ToValue converts a hardware mapping device ID into a Terraform value.
func (did DeviceID) ToValue() types.String {
	return types.StringValue(did.String())
}

// UnmarshalJSON unmarshals a hardware mapping device ID.
func (did *DeviceID) UnmarshalJSON(b []byte) error {
	var pciMapID string

	err := json.Unmarshal(b, &pciMapID)
	if err != nil {
		return errors.Join(ErrDeviceIDUnmarshal, err)
	}

	resType, err := ParseDeviceID(pciMapID)
	if err == nil {
		*did = resType
	}

	return err
}

// ParseDeviceID parses a string that represents a hardware mapping device ID into a DeviceID.
func ParseDeviceID(input string) (DeviceID, error) {
	if !DeviceIDAttrValueRegEx.MatchString(input) {
		return "", ErrDeviceIDParsing(input)
	}

	return DeviceID(input), nil
}
