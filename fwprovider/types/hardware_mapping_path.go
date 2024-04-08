/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package types

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

	proxmoxtypes "github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

// Ensure the implementations satisfy the required interfaces.
var (
	_ basetypes.StringTypable  = HardwareMappingPathType{}
	_ basetypes.StringValuable = HardwareMappingPathValue{}
)

// HardwareMappingPathErrValueConversion indicates an error while converting a value for a hardware mapping path.
//
//nolint:gochecknoglobals
var HardwareMappingPathErrValueConversion = func(format string, attrs ...any) error {
	return function.NewFuncError(fmt.Sprintf(format, attrs...))
}

var (
	// HardwareMappingPathPCIValueRegEx is the regular expression for a PCI hardware mapping path.
	HardwareMappingPathPCIValueRegEx = regexp.MustCompile(`^[a-f0-9]{4,}:[a-f0-9]{2}:[a-f0-9]{2}(\.[a-f0-9])?$`)

	// HardwareMappingPathUSBValueRegEx is the regular expression for a USB hardware mapping path.
	HardwareMappingPathUSBValueRegEx = regexp.MustCompile(`^\d+-(\d+)(\.\d+)?$`)
)

// HardwareMappingPathType is a type that represents a path of a hardware mapping.
type HardwareMappingPathType struct {
	basetypes.StringType
	HardwareMappingType proxmoxtypes.HardwareMappingType
}

// HardwareMappingPathValue is a type that represents the value of a hardware mapping path.
type HardwareMappingPathValue struct {
	basetypes.StringValue
}

// Equal returns true if the two types are equal.
func (t HardwareMappingPathType) Equal(o attr.Type) bool {
	other, ok := o.(HardwareMappingPathType)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

// String returns a string representation of the type.
func (t HardwareMappingPathType) String() string {
	return "HardwareMappingPathType"
}

// ValueFromString converts a string value to a StringValuable.
func (t HardwareMappingPathType) ValueFromString(_ context.Context, in basetypes.StringValue) (
	basetypes.StringValuable,
	diag.Diagnostics,
) {
	value := HardwareMappingPathValue{
		StringValue: in,
	}

	return value, nil
}

// ValueFromTerraform converts a Terraform value to a StringValuable.
func (t HardwareMappingPathType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	attrValue, err := t.StringType.ValueFromTerraform(ctx, in)
	if err != nil {
		return nil, errors.Join(
			HardwareMappingPathErrValueConversion("unexpected error converting Terraform value to StringValue"),
			err,
		)
	}

	stringValue, ok := attrValue.(basetypes.StringValue)
	if !ok {
		return nil, HardwareMappingPathErrValueConversion("unexpected value type of %T", attrValue)
	}

	stringValuable, diags := t.ValueFromString(ctx, stringValue)
	if diags.HasError() {
		return nil, HardwareMappingPathErrValueConversion(
			"unexpected error converting StringValue to StringValuable: %v",
			diags,
		)
	}

	return stringValuable, nil
}

// ValueType returns the underlying value type.
func (t HardwareMappingPathType) ValueType(_ context.Context) attr.Value {
	return HardwareMappingPathValue{}
}

// Equal returns true if the two values are equal.
func (v HardwareMappingPathValue) Equal(o attr.Value) bool {
	other, ok := o.(HardwareMappingPathValue)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

// IsProxmoxType checks whether the value match the given hardware mapping type.
func (v HardwareMappingPathValue) IsProxmoxType(hmType proxmoxtypes.HardwareMappingType) bool {
	switch hmType {
	case proxmoxtypes.HardwareMappingTypePCI:
		return HardwareMappingPathPCIValueRegEx.MatchString(v.ValueString())
	case proxmoxtypes.HardwareMappingTypeUSB:
		return HardwareMappingPathUSBValueRegEx.MatchString(v.ValueString()) || v.ValueString() == ""
	default:
		return false
	}
}

// Type returns the type of the value.
func (v HardwareMappingPathValue) Type(_ context.Context) attr.Type {
	return HardwareMappingPathType{}
}

// NewHardwareMappingPathPointerValue returns a new HardwareMappingPathValue from a string pointer.
func NewHardwareMappingPathPointerValue(value *string) HardwareMappingPathValue {
	return HardwareMappingPathValue{
		StringValue: types.StringPointerValue(value),
	}
}
