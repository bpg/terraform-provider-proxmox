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
var _ basetypes.StringValuable = IPCIDRValue{}

// IPCIDRValue is a type that represents an IP address in CIDR notation.
type IPCIDRValue struct {
	basetypes.StringValue
}

// Equal returns true if the two values are equal.
func (v IPCIDRValue) Equal(o attr.Value) bool {
	other, ok := o.(IPCIDRValue)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

// Type returns the type of the value.
func (v IPCIDRValue) Type(_ context.Context) attr.Type {
	return IPCIDRType{}
}

// NewIPCIDRPointerValue returns a new IPCIDRValue from a string pointer.
func NewIPCIDRPointerValue(value *string) IPCIDRValue {
	return IPCIDRValue{
		StringValue: types.StringPointerValue(value),
	}
}
