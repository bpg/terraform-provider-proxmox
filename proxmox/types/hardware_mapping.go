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
	"strconv"
	"strings"

	"github.com/google/go-querystring/query"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	// hardwareMappingAttrCountMax is the maximum number of attributes for a hardware mapping where only
	// HardwareMappingTypePCI can reach this limit.
	hardwareMappingAttrCountMax = 6

	// hardwareMappingAttrNameIOMMUGroup is the attribute key name of the IOMMU group in a hardware mapping.
	hardwareMappingAttrNameIOMMUGroup = "iommugroup"

	// hardwareMappingAttrNameNode is the attribute key name of the node in a hardware mapping.
	hardwareMappingAttrNameNode = "node"

	// hardwareMappingAttrNameNode is the attribute key name of the path in a hardware mapping.
	hardwareMappingAttrNamePath = "path"

	// hardwareMappingAttrSeparator is the separator for the attributes in a hardware mapping PCI map.
	hardwareMappingAttrSeparator = ','

	// hardwareMappingAttrValueSeparator is the separator for the attribute key-value pairs in a hardware mapping.
	hardwareMappingAttrValueSeparator = '='

	// HardwareMappingAttrNameDescription is the attribute key name of the description in a hardware mapping.
	HardwareMappingAttrNameDescription = "description"
)

// Ensure the hardware mapping type implements required interfaces.
var (
	_ fmt.Stringer     = &HardwareMapping{}
	_ json.Marshaler   = &HardwareMapping{}
	_ json.Unmarshaler = &HardwareMapping{}
	_ query.Encoder    = &HardwareMapping{}
)

//nolint:gochecknoglobals
var (
	// HardwareMappingErrMarshal indicates an error while marshalling a hardware mapping.
	HardwareMappingErrMarshal = function.NewFuncError("cannot marshal hardware mapping")

	// HardwareMappingErrParsingFormat indicates an error the format of a hardware mapping while parsing.
	HardwareMappingErrParsingFormat = func(format string, attrs ...any) error {
		return function.NewFuncError(fmt.Sprintf(format, attrs...))
	}

	// HardwareMappingErrUnknownAttribute indicates an unknown hardware mapping attribute.
	HardwareMappingErrUnknownAttribute = func(attr string) error {
		return function.NewFuncError(fmt.Sprintf("unknown hardware mapping attribute %q", attr))
	}

	// HardwareMappingErrUnmarshal indicates an error while unmarshalling a hardware mapping.
	HardwareMappingErrUnmarshal = function.NewFuncError("cannot unmarshal hardware mapping")
)

// HardwareMapping represents a hardware mapping composed of multiple attributes.
type HardwareMapping struct {
	// Description is the optional "description" for a hardware mapping for both HardwareMappingTypePCI and
	// HardwareMappingTypeUSB.
	Description *string

	// Description is the required "ID" for a hardware mapping for both HardwareMappingTypePCI and HardwareMappingTypeUSB.
	ID HardwareMappingDeviceID

	// IOMMUGroup is the optional "IOMMU group" for a hardware mapping for HardwareMappingTypePCI.
	// The value is not mandatory for the Proxmox API, but causes a HardwareMappingTypePCI to be incomplete when not set.
	// It is not used for HardwareMappingTypeUSB.
	//
	// Using a pointer is required to prevent the default value of 0 to be used as a valid IOMMU group but differentiate
	// between and unset value instead.
	//
	// References:
	//   - [Proxmox VE Wiki — PCI Passthrough]
	//   - [Linux Kernel Documentations — VFIO - "Virtual Function I/O"]
	//   - [IOMMU DB]
	//
	// [Proxmox VE Wiki — PCI Passthrough]: https://pve.proxmox.com/wiki/PCI_Passthrough
	// [Linux Kernel Documentations — VFIO - "Virtual Function I/O"]: https://docs.kernel.org/driver-api/vfio.html
	// [IOMMU DB]: https://iommu.info
	IOMMUGroup *int64

	// Node is the required "node name" for a hardware mapping for both HardwareMappingTypePCI and HardwareMappingTypeUSB.
	Node string

	// Path is the "path" for a hardware mapping where this field is required for HardwareMappingTypePCI but optional for
	// HardwareMappingTypeUSB.
	Path *string

	// SubsystemID is the optional "subsystem ID" for a hardware mapping for HardwareMappingTypePCI.
	// The value is not mandatory for the Proxmox API, but causes a HardwareMappingTypePCI to be incomplete when not set.
	// It is not used for HardwareMappingTypeUSB.
	SubsystemID HardwareMappingDeviceID
}

// EncodeValues encodes a cluster mapping PCI map field into a URL-encoded set of values.
func (hm HardwareMapping) EncodeValues(key string, v *url.Values) error {
	v.Add(key, hm.String())
	return nil
}

// MarshalJSON marshals a hardware mapping into JSON value.
func (hm HardwareMapping) MarshalJSON() ([]byte, error) {
	bytes, err := json.Marshal(hm.String())
	if err != nil {
		return nil, errors.Join(HardwareMappingErrMarshal, err)
	}

	return bytes, nil
}

// String converts a HardwareMapping value into a string.
func (hm HardwareMapping) String() string {
	joinKV := func(k, v string) string {
		return fmt.Sprintf("%s%s%s", k, string(hardwareMappingAttrValueSeparator), v)
	}
	attrs := make([]string, 0, hardwareMappingAttrCountMax)
	attrs = append(
		attrs,
		joinKV(hardwareMappingAttrNameDeviceID, hm.ID.String()),
		joinKV(hardwareMappingAttrNameNode, hm.Node),
	)

	if hm.Path != nil {
		attrs = append(attrs, joinKV(hardwareMappingAttrNamePath, *hm.Path))
	}

	if hm.Description != nil {
		attrs = append(attrs, joinKV(HardwareMappingAttrNameDescription, *hm.Description))
	}

	if hm.IOMMUGroup != nil {
		attrs = append(attrs, joinKV(hardwareMappingAttrNameIOMMUGroup, strconv.FormatInt(*hm.IOMMUGroup, 10)))
	}

	if hm.SubsystemID != "" {
		attrs = append(attrs, joinKV(hardwareMappingDeviceSubsystemIDAttrName, hm.SubsystemID.String()))
	}

	return strings.Join(attrs, string(hardwareMappingAttrSeparator))
}

// ToValue converts a hardware mapping into a Terraform value.
func (hm HardwareMapping) ToValue() types.String {
	return types.StringValue(hm.String())
}

// UnmarshalJSON unmarshals a hardware mapping.
func (hm *HardwareMapping) UnmarshalJSON(b []byte) error {
	var hmString string

	err := json.Unmarshal(b, &hmString)
	if err != nil {
		return errors.Join(HardwareMappingErrUnmarshal, err)
	}

	resType, err := ParseHardwareMapping(hmString)
	if err == nil {
		*hm = resType
	}

	return err
}

// ParseHardwareMapping parses a string that represents a hardware mapping into a HardwareMapping.
func ParseHardwareMapping(input string) (HardwareMapping, error) {
	hm := HardwareMapping{}
	// Scoped function to return an error when a regular expression for an attribute did not match.
	regExNotMatchErr := func(attr, attrName string, err error) error {
		return errors.Join(
			HardwareMappingErrParsingFormat(
				fmt.Sprintf(
					"invalid format %q for hardware mapping %q attribute",
					attr,
					attrName,
				),
			), err,
		)
	}

	// Split the full PCI map string into its attributes…
	attrs := strings.Split(input, string(hardwareMappingAttrSeparator))
	// …and iterate over each attribute to parse it into the struct fields.
	for _, attr := range attrs {
		attrSplit := strings.Split(attr, string(hardwareMappingAttrValueSeparator))
		if len(attrSplit) != 2 {
			return hm, HardwareMappingErrParsingFormat(
				fmt.Sprintf(
					`invalid "key=value" format for hardware mapping attribute %q`,
					attr,
				),
			)
		}

		switch attrSplit[0] {
		case HardwareMappingAttrNameDescription:
			hm.Description = &attrSplit[1]

		case hardwareMappingAttrNameDeviceID:
			id, err := ParseHardwareMappingDeviceID(attrSplit[1])
			if err != nil {
				return hm, regExNotMatchErr(attrSplit[1], hardwareMappingAttrNameDeviceID, err)
			}

			hm.ID = id

		case hardwareMappingAttrNameNode:
			hm.Node = attrSplit[1]

		case hardwareMappingAttrNamePath:
			hm.Path = &attrSplit[1]

		case hardwareMappingAttrNameIOMMUGroup:
			iommuGroup, err := strconv.ParseInt(attrSplit[1], 10, 0)
			if err != nil {
				return hm, regExNotMatchErr(attrSplit[1], hardwareMappingAttrNameIOMMUGroup, err)
			}

			hm.IOMMUGroup = &iommuGroup

		case hardwareMappingDeviceSubsystemIDAttrName:
			subsystemID, err := ParseHardwareMappingDeviceID(attrSplit[1])
			if err != nil {
				return hm, regExNotMatchErr(attrSplit[1], hardwareMappingDeviceSubsystemIDAttrName, err)
			}

			hm.SubsystemID = subsystemID

		default:
			return hm, HardwareMappingErrUnknownAttribute(attr)
		}
	}

	return hm, nil
}
