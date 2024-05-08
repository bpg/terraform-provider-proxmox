package hardwaremapping

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"

	"github.com/google/go-querystring/query"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

//nolint:gochecknoglobals
var (
	// TypePCI is an identifier for a PCI hardware mapping type.
	// Do not modify this package-global variable as it acts as a safer variant compared to "iota" based constants!
	TypePCI = Type{"pci"}

	// TypeUSB is an identifier for a PCI hardware mapping type.
	// Do not modify this package-global variable as it acts as a safer variant compared to "iota" based constants!
	TypeUSB = Type{"usb"}
)

// Ensure the hardware mapping type supports required interfaces.
var (
	_ fmt.Stringer     = new(Type)
	_ json.Marshaler   = new(Type)
	_ json.Unmarshaler = new(Type)
	_ query.Encoder    = new(Type)
)

// Type is the type of the hardware mapping.
type Type struct {
	name string
}

// EncodeValues encodes a hardware mapping type field into a URL-encoded set of values.
func (t Type) EncodeValues(key string, v *url.Values) error {
	v.Add(key, t.String())
	return nil
}

// MarshalJSON marshals a hardware mapping type into JSON value.
func (t Type) MarshalJSON() ([]byte, error) {
	bytes, err := json.Marshal(t.String())
	if err != nil {
		return nil, errors.Join(ErrTypeMarshal, err)
	}

	return bytes, nil
}

// String converts a Type value into a string.
func (t Type) String() string {
	return t.name
}

// ToValue converts a hardware mapping type into a Terraform value.
func (t Type) ToValue() types.String {
	return types.StringValue(t.String())
}

// UnmarshalJSON unmarshals a hardware mapping type.
func (t *Type) UnmarshalJSON(b []byte) error {
	var rtString string

	err := json.Unmarshal(b, &rtString)
	if err != nil {
		return errors.Join(ErrTypeUnmarshal, err)
	}

	resType, err := ParseType(rtString)
	if err == nil {
		*t = resType
	}

	return err
}

// ParseType converts the string representation of a hardware mapping type into the corresponding value.
// An error is returned if the input string does not match any known type.
func ParseType(input string) (Type, error) {
	switch input {
	case TypePCI.String():
		return TypePCI, nil
	case TypeUSB.String():
		return TypeUSB, nil
	default:
		return Type{}, ErrTypeIllegal(input)
	}
}
