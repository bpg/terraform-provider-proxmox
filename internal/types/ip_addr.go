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
var _ basetypes.StringTypable = IPAddrType{}

// IPAddrType is a type that represents an IP address.
type IPAddrType struct {
	basetypes.StringType
}

// Equal returns true if the two types are equal.
func (t IPAddrType) Equal(o attr.Type) bool {
	other, ok := o.(IPAddrType)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

// String returns a string representation of the type.
func (t IPAddrType) String() string {
	return "IPAddrType"
}

// ValueFromString converts a string value to a StringValuable.
func (t IPAddrType) ValueFromString(
	_ context.Context, in basetypes.StringValue,
) (basetypes.StringValuable, diag.Diagnostics) {
	value := IPAddrValue{
		StringValue: in,
	}

	return value, nil
}

// ValueFromTerraform converts a Terraform value to a StringValuable.
func (t IPAddrType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
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
func (t IPAddrType) ValueType(_ context.Context) attr.Value {
	return IPAddrValue{}
}

// Validate ensures the value is valid IP address.
func (t IPAddrType) Validate(_ context.Context, value tftypes.Value, valuePath path.Path) diag.Diagnostics {
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

	if ip := net.ParseIP(valueString); ip == nil {
		diags.AddAttributeError(
			valuePath,
			"Invalid IP String Value",
			"An unexpected error occurred while converting a string value that was expected to be IPv4/IPv6.\n\n"+
				"Path: "+valuePath.String()+"\n"+
				"Given Value: "+valueString,
		)

		return diags
	}

	return diags
}
