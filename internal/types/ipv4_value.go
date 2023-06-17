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
var _ basetypes.StringValuable = IPv4Value{}

type IPv4Value struct {
	basetypes.StringValue
}

func (v IPv4Value) Equal(o attr.Value) bool {
	other, ok := o.(IPv4Value)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

func (v IPv4Value) Type(ctx context.Context) attr.Type {
	return IPv4Type{}
}

func NewIPv4PointerValue(value *string) IPv4Value {
	return IPv4Value{
		StringValue: types.StringPointerValue(value),
	}
}
