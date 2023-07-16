/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package types

import (
	"context"
	"fmt"
	"net"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// Ensure the implementation satisfies the expected interfaces.
var _ basetypes.StringTypable = IPCIDRType{}

// IPCIDRType is a type that represents an IP address in CIDR notation.
type IPCIDRType struct {
	basetypes.StringType
}

// Equal returns true if the two types are equal.
func (t IPCIDRType) Equal(o attr.Type) bool {
	other, ok := o.(IPCIDRType)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

// String returns a string representation of the type.
func (t IPCIDRType) String() string {
	return "IPCIDRType"
}

// ValueFromString converts a string value to a StringValuable.
func (t IPCIDRType) ValueFromString(
	_ context.Context, in basetypes.StringValue,
) (basetypes.StringValuable, diag.Diagnostics) {
	value := IPCIDRValue{
		StringValue: in,
	}

	return value, nil
}

// ValueFromTerraform converts a Terraform value to a StringValuable.
func (t IPCIDRType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	attrValue, err := t.StringType.ValueFromTerraform(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("unexpected error converting Terraform value to StringValue: %w", err)
	}

	stringValue, ok := attrValue.(basetypes.StringValue)
	if !ok {
		return nil, fmt.Errorf("unexpected value type of %T", attrValue)
	}

	stringValuable, diags := t.ValueFromString(ctx, stringValue)
	if diags.HasError() {
		return nil, fmt.Errorf("unexpected error converting StringValue to StringValuable: %v", diags)
	}

	return stringValuable, nil
}

// ValueType returns the underlying value type.
func (t IPCIDRType) ValueType(_ context.Context) attr.Value {
	return IPCIDRValue{}
}

// Validate ensures the value is valid IP address in CIDR notation.
func (t IPCIDRType) Validate(_ context.Context, value tftypes.Value, valuePath path.Path) diag.Diagnostics {
	if value.IsNull() || !value.IsKnown() {
		return nil
	}

	var diags diag.Diagnostics

	var valueString string

	if err := value.As(&valueString); err != nil {
		diags.AddAttributeError(
			valuePath,
			"Invalid Terraform Value",
			"An unexpected error occurred while attempting to convert a Terraform value to a string. "+
				"This generally is an issue with the provider schema implementation. "+
				"Please contact the provider developers.\n\n"+
				"Path: "+valuePath.String()+"\n"+
				"Error: "+err.Error(),
		)

		return diags
	}

	if _, _, err := net.ParseCIDR(valueString); err != nil {
		diags.AddAttributeError(
			valuePath,
			"Invalid IP/CIDR String Value",
			"An unexpected error occurred while converting a string value that was expected to be IP/CIDR.\n\n"+
				"Path: "+valuePath.String()+"\n"+
				"Given Value: "+valueString+"\n"+
				"Error: "+err.Error(),
		)

		return diags
	}

	return diags
}
