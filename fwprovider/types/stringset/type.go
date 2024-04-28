package stringset

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
	_ basetypes.SetTypable = Type{}
)

// Type defines the type for string set.
type Type struct {
	basetypes.SetType
}

// Equal returns true if the two types are equal.
func (t Type) Equal(o attr.Type) bool {
	other, ok := o.(Type)

	if !ok {
		return false
	}

	return t.SetType.Equal(other.SetType)
}

// String returns a string representation of the type.
func (t Type) String() string {
	return "StringSetType"
}

// ValueFromSet converts the set value to a SetValuable type.
func (t Type) ValueFromSet(_ context.Context, in basetypes.SetValue) (basetypes.SetValuable, diag.Diagnostics) {
	value := Value{
		SetValue: in,
	}

	return value, nil
}

// ValueFromTerraform converts the Terraform value to a NewValue type.
func (t Type) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	attrValue, err := t.SetType.ValueFromTerraform(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("error converting Terraform value to NewValue")
	}

	setValue, ok := attrValue.(basetypes.SetValue)
	if !ok {
		return nil, fmt.Errorf("unexpected value type of %T", attrValue)
	}

	setValuable, diags := t.ValueFromSet(ctx, setValue)
	if diags.HasError() {
		return nil, fmt.Errorf("error converting NewValue to SetValuable: %v", diags)
	}

	return setValuable, nil
}

// ValueType returns the underlying value type.
func (t Type) ValueType(_ context.Context) attr.Value {
	return Value{}
}
