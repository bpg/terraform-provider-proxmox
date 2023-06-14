/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package types

import "github.com/bpg/terraform-provider-proxmox/internal/types"

// StrPtr returns a pointer to a string.
func StrPtr(s string) *string {
	return &s
}

// BoolPtr returns a pointer to a bool.
func BoolPtr(s bool) *types.CustomBool {
	customBool := types.CustomBool(s)
	return &customBool
}
