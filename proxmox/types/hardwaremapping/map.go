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
	"strconv"
	"strings"

	"github.com/google/go-querystring/query"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	// attrCountMax is the maximum number of attributes for a hardware mapping where only TypePCI can reach
	// this limit.
	attrCountMax = 6

	// attrNameIOMMUGroup is the attribute key name of the IOMMU group in a hardware mapping.
	attrNameIOMMUGroup = "iommugroup"

	// attrNameNode is the attribute key name of the node in a hardware mapping.
	attrNameNode = "node"

	// attrNameNode is the attribute key name of the path in a hardware mapping.
	attrNamePath = "path"

	// attrSeparator is the separator for the attributes in a hardware mapping PCI map.
	attrSeparator = ','

	// attrValueSeparator is the separator for the attribute key-value pairs in a hardware mapping.
	attrValueSeparator = '='

	// AttrNameDescription is the attribute key name of the description in a hardware mapping.
	// The Proxmox VE API attribute is named "description" while we name it "comment" internally since this naming is
	// generally used across the Proxmox VE web UI and API documentations. This still follows the
	// [Terraform "best practices"] as it improves the user experience by matching the field name to the naming used in
	// the human-facing interfaces.
	//
	// [Terraform "best practices"]: https://developer.hashicorp.com/terraform/plugin/best-practices/hashicorp-provider-design-principles#resource-and-attribute-schema-should-closely-match-the-underlying-api
	AttrNameDescription = "description"
)

// Ensure the hardware mapping type implements required interfaces.
var (
	_ fmt.Stringer     = &Map{}
	_ json.Marshaler   = &Map{}
	_ json.Unmarshaler = &Map{}
	_ query.Encoder    = &Map{}
)

// Map represents a hardware mapping composed of multiple attributes.
type Map struct {
	// Description is the optional "description" for a hardware mapping for both TypePCI and TypeUSB.
	Description *string

	// Description is the required "ID" for a hardware mapping for both TypePCI and TypeUSB.
	ID DeviceID

	// IOMMUGroup is the optional "IOMMU group" for a hardware mapping for TypePCI.
	// The value is not mandatory for the Proxmox VE API, but causes a TypePCI to be incomplete when not set.
	// It is not used for TypeUSB.
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

	// Node is the required "node name" for a hardware mapping for both TypePCI and TypeUSB.
	Node string

	// Path is the "path" for a hardware mapping where this field is required for TypePCI but optional for
	// TypeUSB.
	Path *string

	// SubsystemID is the optional "subsystem ID" for a hardware mapping for TypePCI.
	// The value is not mandatory for the Proxmox VE API, but causes a TypePCI to be incomplete when not set.
	// It is not used for TypeUSB.
	SubsystemID DeviceID
}

// EncodeValues encodes a cluster mapping PCI map field into a URL-encoded set of values.
func (hm Map) EncodeValues(key string, v *url.Values) error {
	v.Add(key, hm.String())
	return nil
}

// MarshalJSON marshals a hardware mapping into JSON value.
func (hm Map) MarshalJSON() ([]byte, error) {
	bytes, err := json.Marshal(hm.String())
	if err != nil {
		return nil, errors.Join(ErrMapMarshal, err)
	}

	return bytes, nil
}

// String converts a Map value into a string.
func (hm Map) String() string {
	joinKV := func(k, v string) string {
		return fmt.Sprintf("%s%s%s", k, string(attrValueSeparator), v)
	}
	attrs := make([]string, 0, attrCountMax)

	// ID is optional for directory mappings
	if hm.ID != "" {
		attrs = append(attrs, joinKV(attrNameDeviceID, hm.ID.String()))
	}

	// Node is common among all mappings
	attrs = append(attrs, joinKV(attrNameNode, hm.Node))

	if hm.Path != nil {
		attrs = append(attrs, joinKV(attrNamePath, *hm.Path))
	}

	if hm.Description != nil {
		attrs = append(attrs, joinKV(AttrNameDescription, *hm.Description))
	}

	if hm.IOMMUGroup != nil {
		attrs = append(attrs, joinKV(attrNameIOMMUGroup, strconv.FormatInt(*hm.IOMMUGroup, 10)))
	}

	if hm.SubsystemID != "" {
		attrs = append(attrs, joinKV(attrNameSubsystemID, hm.SubsystemID.String()))
	}

	return strings.Join(attrs, string(attrSeparator))
}

// ToValue converts a hardware mapping into a Terraform value.
func (hm Map) ToValue() types.String {
	return types.StringValue(hm.String())
}

// UnmarshalJSON unmarshals a hardware mapping.
func (hm *Map) UnmarshalJSON(b []byte) error {
	var hmString string

	err := json.Unmarshal(b, &hmString)
	if err != nil {
		return errors.Join(ErrMapUnmarshal, err)
	}

	resType, err := ParseMap(hmString)
	if err == nil {
		*hm = resType
	}

	return err
}

// ParseMap parses a string that represents a hardware mapping into a Map.
func ParseMap(input string) (Map, error) {
	hm := Map{}
	// Scoped function to return an error when a regular expression for an attribute did not match.
	regExNotMatchErr := func(attr, attrName string, err error) error {
		return errors.Join(
			ErrMapParsingFormat(
				fmt.Sprintf(
					"invalid format %q for hardware mapping %q attribute",
					attr,
					attrName,
				),
			), err,
		)
	}

	// Split the full PCI map string into its attributes…
	attrs := strings.Split(input, string(attrSeparator))
	// …and iterate over each attribute to parse it into the struct fields.
	for _, attr := range attrs {
		attrSplit := strings.Split(attr, string(attrValueSeparator))
		if len(attrSplit) != 2 {
			return hm, ErrMapParsingFormat(
				fmt.Sprintf(
					`invalid "key=value" format for hardware mapping attribute %q`,
					attr,
				),
			)
		}

		switch attrSplit[0] {
		case AttrNameDescription:
			hm.Description = &attrSplit[1]

		case attrNameDeviceID:
			id, err := ParseDeviceID(attrSplit[1])
			if err != nil {
				return hm, regExNotMatchErr(attrSplit[1], attrNameDeviceID, err)
			}

			hm.ID = id

		case attrNameNode:
			hm.Node = attrSplit[1]

		case attrNamePath:
			hm.Path = &attrSplit[1]

		case attrNameIOMMUGroup:
			iommuGroup, err := strconv.ParseInt(attrSplit[1], 10, 0)
			if err != nil {
				return hm, regExNotMatchErr(attrSplit[1], attrNameIOMMUGroup, err)
			}

			hm.IOMMUGroup = &iommuGroup

		case attrNameSubsystemID:
			subsystemID, err := ParseDeviceID(attrSplit[1])
			if err != nil {
				return hm, regExNotMatchErr(attrSplit[1], attrNameSubsystemID, err)
			}

			hm.SubsystemID = subsystemID

		default:
			return hm, ErrMapUnknownAttribute(attr)
		}
	}

	return hm, nil
}
