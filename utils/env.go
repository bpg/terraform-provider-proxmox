/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package utils

import "os"

// GetAnyStringEnv returns the first non-empty string value from the environment variables.
func GetAnyStringEnv(ks ...string) string {
	for _, k := range ks {
		if v := os.Getenv(k); v != "" {
			return v
		}
	}

	return ""
}

// GetAnyBoolEnv returns the first non-empty boolean value from the environment variables.
func GetAnyBoolEnv(ks ...string) bool {
	val := ""

	for _, k := range ks {
		if v := os.Getenv(k); v != "" {
			val = v
			break
		}
	}

	return val == "true" || val == "1"
}
