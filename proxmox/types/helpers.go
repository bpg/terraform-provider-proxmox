/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package types

func StrPtr(s string) *string {
	return &s
}

func BoolPtr(s bool) *CustomBool {
	customBool := CustomBool(s)
	return &customBool
}
