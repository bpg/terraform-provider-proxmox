/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package utils

// ConvertToStringSlice helps convert interface slice to string slice.
func ConvertToStringSlice(interfaceSlice []interface{}) []string {
	resultSlice := make([]string, len(interfaceSlice))

	for _, val := range interfaceSlice {
		resultSlice = append(resultSlice, val.(string))
	}

	return resultSlice
}
