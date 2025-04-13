/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package hardwaremapping

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	proxmoxtypes "github.com/bpg/terraform-provider-proxmox/proxmox/types/hardwaremapping"
)

// Ensure the implementations satisfy the required interfaces.
var (
	_ basetypes.StringTypable  = PathType{}
	_ basetypes.StringValuable = PathValue{}
)

// ErrValueConversion indicates an error while converting a value for a hardware mapping path.
//
//nolint:gochecknoglobals
var ErrValueConversion = func(format string, attrs ...any) error {
	// bpg: this doesn't seem to be a proper use of the function.NewFuncError
	return function.NewFuncError(fmt.Sprintf(format, attrs...))
}

var (
	// PathDirValueRegEx is the regular expression for a POSIX path.
	PathDirValueRegEx = regexp.MustCompile(`^/.+$`)

	// PathPCIValueRegEx is the regular expression for a PCI hardware mapping path.
	PathPCIValueRegEx = regexp.MustCompile(`^[a-f0-9]{4,}:[a-f0-9]{2}:[a-f0-9]{2}(\.[a-f0-9])?$`)

	// PathUSBValueRegEx is the regular expression for a USB hardware mapping path.
	PathUSBValueRegEx = regexp.MustCompile(`^\d+-(\d+)(\.\d+)?$`)
)

// PathType is a type that represents a path of a hardware mapping.
type PathType struct {
	basetypes.StringType
	Type proxmoxtypes.Type
}

// PathValue is a type that represents the value of a hardware mapping path.
type PathValue struct {
	basetypes.StringValue
}

// Equal returns true if the two types are equal.
func (t PathType) Equal(o attr.Type) bool {
	other, ok := o.(PathType)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

// String returns a string representation of the type.
func (t PathType) String() string {
	return "PathType"
}

// ValueFromString converts a string value to a StringValuable.
func (t PathType) ValueFromString(_ context.Context, in basetypes.StringValue) (
	basetypes.StringValuable,
	diag.Diagnostics,
) {
	value := PathValue{
		StringValue: in,
	}

	return value, nil
}

// ValueFromTerraform converts a Terraform value to a StringValuable.
func (t PathType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	attrValue, err := t.StringType.ValueFromTerraform(ctx, in)
	if err != nil {
		return nil, errors.Join(
			ErrValueConversion("unexpected error converting Terraform value to StringValue"),
			err,
		)
	}

	stringValue, ok := attrValue.(basetypes.StringValue)
	if !ok {
		return nil, ErrValueConversion("unexpected value type of %T", attrValue)
	}

	stringValuable, diags := t.ValueFromString(ctx, stringValue)
	if diags.HasError() {
		return nil, ErrValueConversion(
			"unexpected error converting StringValue to StringValuable: %v",
			diags,
		)
	}

	return stringValuable, nil
}

// ValueType returns the underlying value type.
func (t PathType) ValueType(_ context.Context) attr.Value {
	return PathValue{}
}

// Equal returns true if the two values are equal.
func (v PathValue) Equal(o attr.Value) bool {
	other, ok := o.(PathValue)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

// IsProxmoxType checks whether the value match the given hardware mapping type.
func (v PathValue) IsProxmoxType(hmType proxmoxtypes.Type) bool {
	switch hmType {
	case proxmoxtypes.TypeDir:
		return PathDirValueRegEx.MatchString(v.ValueString())
	case proxmoxtypes.TypePCI:
		return PathPCIValueRegEx.MatchString(v.ValueString())
	case proxmoxtypes.TypeUSB:
		return PathUSBValueRegEx.MatchString(v.ValueString()) || v.ValueString() == ""
	default:
		return false
	}
}

// Type returns the type of the value.
func (v PathValue) Type(_ context.Context) attr.Type {
	return PathType{}
}

// NewPathPointerValue returns a new PathValue from a string pointer.
func NewPathPointerValue(value *string) PathValue {
	return PathValue{
		StringValue: types.StringPointerValue(value),
	}
}
