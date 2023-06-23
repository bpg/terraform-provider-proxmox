/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package types

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Ensure the implementation satisfies the expected interfaces.
var _ basetypes.StringValuable = IPAddrValue{}

// IPAddrValue is a type that represents an IP address value.
type IPAddrValue struct {
	basetypes.StringValue
}

// Equal returns true if the two values are equal.
func (v IPAddrValue) Equal(o attr.Value) bool {
	other, ok := o.(IPAddrValue)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

// Type returns the type of the value.
func (v IPAddrValue) Type(_ context.Context) attr.Type {
	return IPAddrType{}
}

// NewIPAddrPointerValue returns a new IPAddrValue from a string pointer.
func NewIPAddrPointerValue(value *string) IPAddrValue {
	return IPAddrValue{
		StringValue: types.StringPointerValue(value),
	}
}
