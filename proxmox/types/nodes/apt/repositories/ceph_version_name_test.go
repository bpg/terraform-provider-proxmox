/*
	This Source Code Form is subject to the terms of the Mozilla Public
	License, v. 2.0. If a copy of the MPL was not distributed with this
	file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/

package repositories

import (
	"encoding/json"
	"testing"
)

func TestParseCephVersionName(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		input       string
		expected    CephVersionName
		expectError bool
	}{
		"quincy": {
			input:    "quincy",
			expected: CephVersionNameQuincy,
		},
		"reef": {
			input:    "reef",
			expected: CephVersionNameReef,
		},
		"squid": {
			input:    "squid",
			expected: CephVersionNameSquid,
		},
		"invalid": {
			input:       "invalid",
			expectError: true,
		},
		"empty": {
			input:       "",
			expectError: true,
		},
		"uppercase": {
			input:       "QUINCY",
			expectError: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result, err := ParseCephVersionName(test.input)

			if test.expectError {
				if err == nil {
					t.Errorf("expected error for input %q, got nil", test.input)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for input %q: %v", test.input, err)
				}

				if result != test.expected {
					t.Errorf("expected %v, got %v", test.expected, result)
				}
			}
		})
	}
}

func TestCephVersionNameString(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		version  CephVersionName
		expected string
	}{
		"quincy": {
			version:  CephVersionNameQuincy,
			expected: "quincy",
		},
		"reef": {
			version:  CephVersionNameReef,
			expected: "reef",
		},
		"squid": {
			version:  CephVersionNameSquid,
			expected: "squid",
		},
		"unknown": {
			version:  CephVersionNameUnknown,
			expected: "unknown",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result := test.version.String()
			if result != test.expected {
				t.Errorf("expected %q, got %q", test.expected, result)
			}
		})
	}
}

func TestCephVersionNameMarshalJSON(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		version  CephVersionName
		expected string
	}{
		"quincy": {
			version:  CephVersionNameQuincy,
			expected: `"quincy"`,
		},
		"reef": {
			version:  CephVersionNameReef,
			expected: `"reef"`,
		},
		"squid": {
			version:  CephVersionNameSquid,
			expected: `"squid"`,
		},
		"unknown": {
			version:  CephVersionNameUnknown,
			expected: `"unknown"`,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result, err := json.Marshal(test.version)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if string(result) != test.expected {
				t.Errorf("expected %q, got %q", test.expected, string(result))
			}
		})
	}
}

func TestCephVersionNameUnmarshalJSON(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		input       string
		expected    CephVersionName
		expectError bool
	}{
		"quincy": {
			input:    `"quincy"`,
			expected: CephVersionNameQuincy,
		},
		"reef": {
			input:    `"reef"`,
			expected: CephVersionNameReef,
		},
		"squid": {
			input:    `"squid"`,
			expected: CephVersionNameSquid,
		},
		"invalid": {
			input:       `"invalid"`,
			expectError: true,
		},
		"malformed json": {
			input:       `quincy`,
			expectError: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var result CephVersionName
			err := json.Unmarshal([]byte(test.input), &result)

			if test.expectError {
				if err == nil {
					t.Errorf("expected error for input %q, got nil", test.input)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for input %q: %v", test.input, err)
				}

				if result != test.expected {
					t.Errorf("expected %v, got %v", test.expected, result)
				}
			}
		})
	}
}
