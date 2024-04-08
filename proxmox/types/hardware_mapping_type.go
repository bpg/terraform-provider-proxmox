package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"

	"github.com/google/go-querystring/query"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

//nolint:gochecknoglobals
var (
	// HardwareMappingTypeErrIllegal indicates an error for an illegal hardware mapping type.
	HardwareMappingTypeErrIllegal = func(hmTypeName string) error {
		return function.NewFuncError(fmt.Sprintf("illegal hardware mapping type %q", hmTypeName))
	}

	// HardwareMappingTypeErrMarshal indicates an error while marshalling a hardware mapping type.
	HardwareMappingTypeErrMarshal = function.NewFuncError("cannot marshal hardware mapping type")

	// HardwareMappingTypeErrUnmarshal indicates an error while unmarshalling a hardware mapping type.
	HardwareMappingTypeErrUnmarshal = function.NewFuncError("cannot unmarshal hardware mapping type")
)

//nolint:gochecknoglobals
var (
	// HardwareMappingTypePCI is an identifier for a PCI hardware mapping type.
	// Do not modify this package-global variable as it acts as a safer variant compared to "iota" based constants!
	HardwareMappingTypePCI = HardwareMappingType{"pci"}

	// HardwareMappingTypeUSB is an identifier for a PCI hardware mapping type.
	// Do not modify this package-global variable as it acts as a safer variant compared to "iota" based constants!
	HardwareMappingTypeUSB = HardwareMappingType{"usb"}
)

// Ensure the hardware mapping type supports required interfaces.
var (
	_ fmt.Stringer     = new(HardwareMappingType)
	_ json.Marshaler   = new(HardwareMappingType)
	_ json.Unmarshaler = new(HardwareMappingType)
	_ query.Encoder    = new(HardwareMappingType)
)

// HardwareMappingType is the type of the hardware mapping.
type HardwareMappingType struct {
	name string
}

// EncodeValues encodes a hardware mapping type field into a URL-encoded set of values.
func (t HardwareMappingType) EncodeValues(key string, v *url.Values) error {
	v.Add(key, t.String())
	return nil
}

// MarshalJSON marshals a hardware mapping type into JSON value.
func (t HardwareMappingType) MarshalJSON() ([]byte, error) {
	bytes, err := json.Marshal(t.String())
	if err != nil {
		return nil, errors.Join(HardwareMappingTypeErrMarshal, err)
	}

	return bytes, nil
}

// String converts a HardwareMappingType value into a string.
func (t HardwareMappingType) String() string {
	return t.name
}

// ToValue converts a hardware mapping type into a Terraform value.
func (t HardwareMappingType) ToValue() types.String {
	return types.StringValue(t.String())
}

// UnmarshalJSON unmarshals a hardware mapping type.
func (t *HardwareMappingType) UnmarshalJSON(b []byte) error {
	var rtString string

	err := json.Unmarshal(b, &rtString)
	if err != nil {
		return errors.Join(HardwareMappingTypeErrUnmarshal, err)
	}

	resType, err := ParseHardwareMappingType(rtString)
	if err == nil {
		*t = resType
	}

	return err
}

// ParseHardwareMappingType converts the string representation of a hardware mapping type into the corresponding value.
// An error is returned if the input string does not match any known type.
func ParseHardwareMappingType(input string) (HardwareMappingType, error) {
	switch input {
	case HardwareMappingTypePCI.String():
		return HardwareMappingTypePCI, nil
	case HardwareMappingTypeUSB.String():
		return HardwareMappingTypeUSB, nil
	default:
		return HardwareMappingType{}, HardwareMappingTypeErrIllegal(input)
	}
}
