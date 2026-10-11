/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package utils

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// GetAnyStringEnv returns the first non-empty string value from the environment variables.
func GetAnyStringEnv(ks ...string) string {
	for _, k := range ks {
		if v := os.Getenv(k); v != "" {
			return v
		}
	}

	return ""
}

// GetAnyStringMapEnv parses the first non-empty environment variable as a comma-separated list of
// `Name=Value` pairs. Returns nil when none of the variables is set. A value must not contain a
// comma. Errors reference the position of the offending pair, never its content, as the variable
// may hold credentials.
func GetAnyStringMapEnv(ks ...string) (map[string]string, error) {
	v := GetAnyStringEnv(ks...)
	if v == "" {
		return nil, nil //nolint:nilnil // a nil map means "no headers configured", which is not an error
	}

	result := map[string]string{}

	for i, pair := range strings.Split(v, ",") {
		if strings.TrimSpace(pair) == "" {
			continue
		}

		name, value, found := strings.Cut(pair, "=")
		if !found {
			return nil, fmt.Errorf("pair %d is not in the `Name=Value` format", i+1)
		}

		name = strings.TrimSpace(name)
		if name == "" {
			return nil, fmt.Errorf("pair %d has an empty name", i+1)
		}

		result[name] = strings.TrimSpace(value)
	}

	return result, nil
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

// GetAnyIntEnv returns the first non-empty integer value from the environment variables.
func GetAnyIntEnv(ks ...string) int {
	for _, k := range ks {
		if v := os.Getenv(k); v != "" {
			if i, err := strconv.Atoi(v); err == nil {
				return i
			}
		}
	}

	return 0
}
