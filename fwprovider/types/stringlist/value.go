/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package stringlist

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
	_ basetypes.ListValuable = Value{}
)

// Value defines the value for string list.
type Value struct {
	basetypes.ListValue
}

// Type returns the type of the value.
func (v Value) Type(_ context.Context) attr.Type {
	return Type{
		ListType: basetypes.ListType{ElemType: basetypes.StringType{}},
	}
}

// Equal returns true if the two values are equal.
func (v Value) Equal(o attr.Value) bool {
	other, ok := o.(Value)

	if !ok {
		return false
	}

	return v.ListValue.Equal(other.ListValue)
}

// ValueStringPointer returns a pointer to the string representation of string list value.
func (v Value) ValueStringPointer(ctx context.Context, diags *diag.Diagnostics) *string {
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

	return ptr.Ptr(strings.Join(sanitizedItems, ";"))
}

// NewValue converts a string of items to a new string list value.
func NewValue(str *string, diags *diag.Diagnostics) Value {
	if str == nil {
		return Value{types.ListValueMust(types.StringType, []attr.Value{})}
	}

	items := strings.Split(*str, ";")
	elems := make([]attr.Value, len(items))

	for i, item := range items {
		elems[i] = types.StringValue(item)
	}

	listValue, d := types.ListValue(types.StringType, elems)
	diags.Append(d...)

	return Value{listValue}
}
