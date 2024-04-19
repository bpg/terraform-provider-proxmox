/*
	This Source Code Form is subject to the terms of the Mozilla Public
	License, v. 2.0. If a copy of the MPL was not distributed with this
	file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/

package hardwaremapping

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/function"
)

//nolint:gochecknoglobals
var (
	// ErrIDMarshal indicates an error while marshalling a hardware mapping ID.
	ErrIDMarshal = function.NewFuncError("cannot unmarshal hardware mapping ID")

	// ErrIDParsing indicates an error while parsing a hardware mapping ID.
	ErrIDParsing = func(hmID string) error {
		return function.NewFuncError(fmt.Sprintf("%q is not a valid hardware mapping ID", hmID))
	}

	// ErrDeviceIDMarshal indicates an error while marshalling a hardware mapping device ID.
	ErrDeviceIDMarshal = function.NewFuncError("cannot marshal hardware mapping device ID")

	// ErrDeviceIDParsing indicates an error while parsing a hardware mapping device ID.
	ErrDeviceIDParsing = func(hmID string) error {
		return function.NewFuncError(
			fmt.Sprintf(
				"invalid value %q for hardware mapping device ID attribute %q: no match for regular expression %q",
				hmID,
				attrNameDeviceID,
				DeviceIDAttrValueRegEx.String(),
			),
		)
	}

	// ErrDeviceIDUnmarshal indicates an error while unmarshalling a hardware mapping device ID.
	ErrDeviceIDUnmarshal = function.NewFuncError("cannot unmarshal hardware mapping device ID")

	// ErrMapMarshal indicates an error while marshalling a hardware mapping.
	ErrMapMarshal = function.NewFuncError("cannot marshal hardware mapping")

	// ErrMapParsingFormat indicates an error the format of a hardware mapping while parsing.
	ErrMapParsingFormat = func(format string, attrs ...any) error {
		return function.NewFuncError(fmt.Sprintf(format, attrs...))
	}

	// ErrMapUnknownAttribute indicates an unknown hardware mapping attribute.
	ErrMapUnknownAttribute = func(attr string) error {
		return function.NewFuncError(fmt.Sprintf("unknown hardware mapping attribute %q", attr))
	}

	// ErrMapUnmarshal indicates an error while unmarshalling a hardware mapping.
	ErrMapUnmarshal = function.NewFuncError("cannot unmarshal hardware mapping")

	// ErrTypeIllegal indicates an error for an illegal hardware mapping type.
	ErrTypeIllegal = func(hmTypeName string) error {
		return function.NewFuncError(fmt.Sprintf("illegal hardware mapping type %q", hmTypeName))
	}

	// ErrTypeMarshal indicates an error while marshalling a hardware mapping type.
	ErrTypeMarshal = function.NewFuncError("cannot marshal hardware mapping type")

	// ErrTypeUnmarshal indicates an error while unmarshalling a hardware mapping type.
	ErrTypeUnmarshal = function.NewFuncError("cannot unmarshal hardware mapping type")
)
