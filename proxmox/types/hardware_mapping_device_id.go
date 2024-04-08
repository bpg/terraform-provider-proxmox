/*
	This Source Code Form is subject to the terms of the Mozilla Public
	License, v. 2.0. If a copy of the MPL was not distributed with this
	file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/

package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"regexp"

	"github.com/google/go-querystring/query"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	// hardwareMappingAttrNameDeviceID is the attribute name of the device ID in a hardware mapping.
	hardwareMappingAttrNameDeviceID = "id"

	// hardwareMappingDeviceSubsystemIDAttrName is the attribute name of the device subsystem ID in a hardware mapping.
	hardwareMappingDeviceSubsystemIDAttrName = "subsystem-id"
)

//nolint:gochecknoglobals
var (
	// HardwareMappingDeviceIDErrMarshal indicates an error while marshalling a hardware mapping device ID.
	HardwareMappingDeviceIDErrMarshal = function.NewFuncError("cannot marshal hardware mapping device ID")

	// HardwareMappingDeviceIDErrParsing indicates an error while parsing a hardware mapping device ID.
	HardwareMappingDeviceIDErrParsing = func(hmID string) error {
		return function.NewFuncError(
			fmt.Sprintf(
				"invalid value %q for hardware mapping device ID attribute %q: no match for regular expression %q",
				hmID,
				hardwareMappingAttrNameDeviceID,
				HardwareMappingDeviceIDAttrValueRegEx.String(),
			),
		)
	}

	// HardwareMappingDeviceIDErrUnmarshal indicates an error while unmarshalling a hardware mapping device ID.
	HardwareMappingDeviceIDErrUnmarshal = function.NewFuncError("cannot unmarshal hardware mapping device ID")

	// HardwareMappingDeviceIDAttrValueRegEx is the regular expression for device ID attribute value in a hardware mapping.
	HardwareMappingDeviceIDAttrValueRegEx = regexp.MustCompile(`^[0-9A-Fa-f]{4}:[0-9A-Fa-f]{4}$`)
)

// Ensure the hardware mapping device ID type implements required interfaces.
var (
	_ fmt.Stringer     = new(HardwareMappingDeviceID)
	_ json.Marshaler   = new(HardwareMappingDeviceID)
	_ json.Unmarshaler = new(HardwareMappingDeviceID)
	_ query.Encoder    = new(HardwareMappingDeviceID)
)

// HardwareMappingDeviceID represents a hardware mapping device ID.
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
type HardwareMappingDeviceID string

// EncodeValues encodes a hardware mapping device ID field into a URL-encoded set of values.
func (did HardwareMappingDeviceID) EncodeValues(key string, v *url.Values) error {
	v.Add(key, did.String())

	return nil
}

// MarshalJSON marshals a hardware mapping device ID into JSON value.
func (did HardwareMappingDeviceID) MarshalJSON() ([]byte, error) {
	bytes, err := json.Marshal(did)
	if err != nil {
		return nil, errors.Join(HardwareMappingDeviceIDErrMarshal, err)
	}

	return bytes, nil
}

// String converts a HardwareMappingDeviceID value into a string.
func (did HardwareMappingDeviceID) String() string {
	return string(did)
}

// ToValue converts a hardware mapping device ID into a Terraform value.
func (did HardwareMappingDeviceID) ToValue() types.String {
	return types.StringValue(did.String())
}

// UnmarshalJSON unmarshals a hardware mapping device ID.
func (did *HardwareMappingDeviceID) UnmarshalJSON(b []byte) error {
	var pciMapID string

	err := json.Unmarshal(b, &pciMapID)
	if err != nil {
		return errors.Join(HardwareMappingDeviceIDErrUnmarshal, err)
	}

	resType, err := ParseHardwareMappingDeviceID(pciMapID)
	if err == nil {
		*did = resType
	}

	return err
}

// ParseHardwareMappingDeviceID parses a string that represents a hardware mapping device ID into a
// HardwareMappingDeviceID.
func ParseHardwareMappingDeviceID(input string) (HardwareMappingDeviceID, error) {
	if !HardwareMappingDeviceIDAttrValueRegEx.MatchString(input) {
		return "", HardwareMappingDeviceIDErrParsing(input)
	}

	return HardwareMappingDeviceID(input), nil
}
