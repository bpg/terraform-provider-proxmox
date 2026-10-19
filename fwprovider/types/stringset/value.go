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

	proxmoxtypes "github.com/bpg/terraform-provider-proxmox/proxmox/types"
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

	return new(strings.Join(elems, o.separator))
}

// ValueCommaSeparatedList returns the set as a PVE comma-separated list, or nil when it is null, unknown or empty.
func (v Value) ValueCommaSeparatedList(ctx context.Context, diags *diag.Diagnostics) *proxmoxtypes.CustomCommaSeparatedList {
	elems := v.ValueList(ctx, diags)
	if len(elems) == 0 {
		return nil
	}

	return new(proxmoxtypes.CustomCommaSeparatedList(elems))
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
// Handles special PVE "deleted" value as empty set value.
func NewValueString(str *string, diags *diag.Diagnostics, opts ...Option) Value {
	if str == nil || *str == "" || *str == "deleted" {
		return Value{types.SetValueMust(types.StringType, []attr.Value{})}
	}

	o := defaultOptions(opts...)

	items := strings.Split(*str, o.separator)

	return NewValueList(items, diags)
}

// NewValueCommaSeparatedList converts a PVE comma-separated list to a string set value.
// A missing list and the PVE pending "deleted" marker both map to null.
func NewValueCommaSeparatedList(list *proxmoxtypes.CustomCommaSeparatedList, diags *diag.Diagnostics) Value {
	if list == nil || (len(*list) == 1 && (*list)[0] == "deleted") {
		return NullValue()
	}

	return NewValueList(*list, diags)
}

// NewValueCommaSeparatedListOrEmpty is NewValueCommaSeparatedList for computed-only datasource attributes, which
// must not be null after Read: a missing list maps to an empty set.
func NewValueCommaSeparatedListOrEmpty(list *proxmoxtypes.CustomCommaSeparatedList, diags *diag.Diagnostics) Value {
	if list == nil {
		return NewValueList(nil, diags)
	}

	return NewValueList(*list, diags)
}

// NullValue returns a properly typed null Value.
func NullValue() Value {
	return Value{types.SetNull(types.StringType)}
}
