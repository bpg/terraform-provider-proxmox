/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package tffwk

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// OptStringFromModel returns either a string pointer or a `nil`, depending on whether the input attribute value
// was null/an empty string or some actual string.
func OptStringFromModel(inval types.String) *string {
	value := inval.ValueString()
	if value == "" {
		return nil
	}

	return &value
}

// OptStringToModel converts a string pointer into a Terraform value. If the input value is not `nil` but
// contains an empty string, the result will still be a null.
func OptStringToModel(inval *string) types.String {
	var outval types.String

	if inval == nil || *inval == "" {
		outval = types.StringNull()
	} else {
		outval = types.StringValue(*inval)
	}

	return outval
}

// BoolintFromModel converts a boolean attribute value into the integers `0` and `1`.
func BoolintFromModel(inval types.Bool) int {
	if inval.ValueBool() {
		return 1
	}

	return 0
}

// BoolintToModel converts a Proxmox-provided boolean, encoded as either `0` or `1`, into a boolean attribute value.
func BoolintToModel(inval int) types.Bool {
	return types.BoolValue(inval != 0)
}
