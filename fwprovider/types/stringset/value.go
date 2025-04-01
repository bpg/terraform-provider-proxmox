/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package stringset

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/bpg/terraform-provider-proxmox/proxmox/helpers/ptr"
)

// Ensure the implementations satisfy the required interfaces.
var (
	_ basetypes.SetValuable = Value{}
)

// Value defines the value for string set.
type Value struct {
	basetypes.SetValue
}

// Type returns the type of the value.
func (v Value) Type(_ context.Context) attr.Type {
	return Type{
		SetType: basetypes.SetType{ElemType: basetypes.StringType{}},
	}
}

// Equal returns true if the two values are equal.
func (v Value) Equal(o attr.Value) bool {
	other, ok := o.(Value)

	if !ok {
		return false
	}

	return v.SetValue.Equal(other.SetValue)
}

// ValueList returns a string slice of set value.
func (v Value) ValueList(ctx context.Context, diags *diag.Diagnostics) []string {
	if v.IsNull() || v.IsUnknown() || len(v.Elements()) == 0 {
		return nil
	}

	elems := make([]types.String, 0, len(v.Elements()))
	d := v.ElementsAs(ctx, &elems, false)
	diags.Append(d...)

	if d.HasError() {
		return nil
	}

	var sanitizedItems []string

	for _, el := range elems {
		if el.IsNull() || el.IsUnknown() {
			continue
		}

		sanitizedItem := strings.TrimSpace(el.ValueString())
		if len(sanitizedItem) > 0 {
			sanitizedItems = append(sanitizedItems, sanitizedItem)
		}
	}

	return sanitizedItems
}

// ValueStringPointer returns a pointer to the string representation of string set value.
func (v Value) ValueStringPointer(ctx context.Context, diags *diag.Diagnostics, opts ...Option) *string {
	elems := v.ValueList(ctx, diags)
	if len(elems) == 0 {
		return nil
	}

	o := defaultOptions(opts...)

	return ptr.Ptr(strings.Join(elems, o.separator))
}

// NewValueList converts a slice of items to a new string set value.
func NewValueList(items []string, diags *diag.Diagnostics) Value {
	if len(items) == 0 {
		return Value{types.SetValueMust(types.StringType, []attr.Value{})}
	}

	elems := make([]attr.Value, len(items))

	for i, item := range items {
		elems[i] = types.StringValue(item)
	}

	setValue, d := types.SetValue(types.StringType, elems)
	diags.Append(d...)

	return Value{setValue}
}

// NewValueString converts a string of items to a new string set value.
func NewValueString(str *string, diags *diag.Diagnostics, opts ...Option) Value {
	if str == nil || *str == "" {
		return Value{types.SetValueMust(types.StringType, []attr.Value{})}
	}

	o := defaultOptions(opts...)

	items := strings.Split(*str, o.separator)

	return NewValueList(items, diags)
}
