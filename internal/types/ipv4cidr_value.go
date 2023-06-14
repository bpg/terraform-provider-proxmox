/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package types

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Ensure the implementation satisfies the expected interfaces.
var _ basetypes.StringValuable = IPv4CIDRValue{}

type IPv4CIDRValue struct {
	basetypes.StringValue
}

func (v IPv4CIDRValue) Equal(o attr.Value) bool {
	other, ok := o.(IPv4CIDRValue)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

func (v IPv4CIDRValue) Type(_ context.Context) attr.Type {
	return IPv4CIDRType{}
}
