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
var _ basetypes.StringTypable = IPv4Type{}

type IPv4Type struct {
	basetypes.StringType
}

func (t IPv4Type) Equal(o attr.Type) bool {
	other, ok := o.(IPv4Type)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

func (t IPv4Type) String() string {
	return "IPv4Type"
}

func (t IPv4Type) ValueFromString(
	ctx context.Context, in basetypes.StringValue,
) (basetypes.StringValuable, diag.Diagnostics) {
	value := IPv4Value{
		StringValue: in,
	}

	return value, nil
}

func (t IPv4Type) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
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

func (t IPv4Type) ValueType(_ context.Context) attr.Value {
	return IPv4Value{}
}

func (t IPv4Type) Validate(_ context.Context, value tftypes.Value, valuePath path.Path) diag.Diagnostics {
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

	if ip := net.ParseIP(valueString); ip == nil || ip.To4() == nil {
		diags.AddAttributeError(
			valuePath,
			"Invalid IPv4 String Value",
			"An unexpected error occurred while converting a string value that was expected to be IPv4.\n\n"+
				"Path: "+valuePath.String()+"\n"+
				"Given Value: "+valueString,
		)

		return diags
	}

	return diags
}
