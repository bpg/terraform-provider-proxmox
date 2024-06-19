/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package stringlist

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// Ensure the implementations satisfy the required interfaces.
var (
	_ basetypes.ListTypable = Type{}
)

// Type defines the type for string list.
type Type struct {
	basetypes.ListType
}

// Equal returns true if the two types are equal.
func (t Type) Equal(o attr.Type) bool {
	other, ok := o.(Type)

	if !ok {
		return false
	}

	return t.ListType.Equal(other.ListType)
}

// String returns a string representation of the type.
func (t Type) String() string {
	return "StringListType"
}

// ValueFromList converts the list value to a ListValuable type.
func (t Type) ValueFromList(_ context.Context, in basetypes.ListValue) (basetypes.ListValuable, diag.Diagnostics) {
	value := Value{
		ListValue: in,
	}

	return value, nil
}

// ValueFromTerraform converts the Terraform value to a NewValue type.
func (t Type) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	attrValue, err := t.ListType.ValueFromTerraform(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("error converting Terraform value to NewValue")
	}

	listValue, ok := attrValue.(basetypes.ListValue)
	if !ok {
		return nil, fmt.Errorf("unexpected value type of %T", attrValue)
	}

	listValuable, diags := t.ValueFromList(ctx, listValue)
	if diags.HasError() {
		return nil, fmt.Errorf("error converting NewValue to ListValuable: %v", diags)
	}

	return listValuable, nil
}

// ValueType returns the underlying value type.
func (t Type) ValueType(_ context.Context) attr.Value {
	return Value{}
}
