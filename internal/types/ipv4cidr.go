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
var _ basetypes.StringTypable = IPv4CIDRType{}

type IPv4CIDRType struct {
	basetypes.StringType
}

func (t IPv4CIDRType) Equal(o attr.Type) bool {
	other, ok := o.(IPv4CIDRType)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

func (t IPv4CIDRType) String() string {
	return "IPv4CIDRType"
}

func (t IPv4CIDRType) ValueFromString(
	ctx context.Context, in basetypes.StringValue,
) (basetypes.StringValuable, diag.Diagnostics) {
	value := IPv4CIDRValue{
		StringValue: in,
	}

	return value, nil
}

func (t IPv4CIDRType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
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

func (t IPv4CIDRType) ValueType(_ context.Context) attr.Value {
	return IPv4CIDRValue{}
}

func (t IPv4CIDRType) Validate(_ context.Context, value tftypes.Value, valuePath path.Path) diag.Diagnostics {
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
			"Invalid IPv4/CIDR String Value",
			"An unexpected error occurred while converting a string value that was expected to be IPv4/CIDR.\n\n"+
				"Path: "+valuePath.String()+"\n"+
				"Given Value: "+valueString+"\n"+
				"Error: "+err.Error(),
		)

		return diags
	}

	return diags
}
