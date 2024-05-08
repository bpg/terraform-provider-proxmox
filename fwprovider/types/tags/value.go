package tags

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

// Value defines the value for tags.
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

// ValueStringPointer returns a pointer to the string representation of tags set value.
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

	var sanitizedTags []string

	for _, el := range elems {
		if el.IsNull() || el.IsUnknown() {
			continue
		}

		sanitizedTag := strings.TrimSpace(el.ValueString())
		if len(sanitizedTag) > 0 {
			sanitizedTags = append(sanitizedTags, sanitizedTag)
		}
	}

	return ptr.Ptr(strings.Join(sanitizedTags, ";"))
}

// SetValue converts a string of tags to a tags set value.
func SetValue(tagsStr *string, diags *diag.Diagnostics) Value {
	if tagsStr == nil {
		return Value{types.SetNull(types.StringType)}
	}

	tags := strings.Split(*tagsStr, ";")
	elems := make([]attr.Value, len(tags))

	for i, tag := range tags {
		elems[i] = types.StringValue(tag)
	}

	setValue, d := types.SetValue(types.StringType, elems)
	diags.Append(d...)

	return Value{setValue}
}
