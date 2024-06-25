/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package utils

// ConvertToStringSlice helps convert interface slice to string slice.
func ConvertToStringSlice(interfaceSlice []interface{}) []string {
	resultSlice := make([]string, len(interfaceSlice))

	for i, val := range interfaceSlice {
		resultSlice[i] = val.(string)
	}

	return resultSlice
}
