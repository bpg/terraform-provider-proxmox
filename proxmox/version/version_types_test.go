/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package version

import (
	"encoding/json"
	"testing"

	"github.com/hashicorp/go-version"
)

func TestResponseBody_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		jsonData       string
		expectedError  bool
		expectedResult *ResponseBody
	}{
		{
			name: "valid response with version 9.0.4",
			jsonData: `{
				"data": {
					"repoid": "39d8a4de7dfb2c40",
					"release": "9.0",
					"version": "9.0.4"
				}
			}`,
			expectedError: false,
			expectedResult: &ResponseBody{
				Data: &ResponseData{
					Console:      "",
					Release:      "9.0",
					RepositoryID: "39d8a4de7dfb2c40",
					Version: ProxmoxVersion{
						Version: *mustParseVersion("9.0.4"),
					},
				},
			},
		},
		{
			name: "valid response with semantic version",
			jsonData: `{
				"data": {
					"repoid": "test123",
					"release": "8.2",
					"version": "8.2.1",
					"console": "proxmox"
				}
			}`,
			expectedError: false,
			expectedResult: &ResponseBody{
				Data: &ResponseData{
					Console:      "proxmox",
					Release:      "8.2",
					RepositoryID: "test123",
					Version: ProxmoxVersion{
						Version: *mustParseVersion("8.2.1"),
					},
				},
			},
		},
		{
			name: "invalid version format",
			jsonData: `{
				"data": {
					"repoid": "test123",
					"release": "9.0",
					"version": "invalid-version"
				}
			}`,
			expectedError: true,
		},
		{
			name: "missing data field",
			jsonData: `{
				"repoid": "test123",
				"release": "9.0",
				"version": "9.0.4"
			}`,
			expectedError: false,
			expectedResult: &ResponseBody{
				Data: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var result ResponseBody

			err := json.Unmarshal([]byte(tt.jsonData), &result)

			if tt.expectedError {
				if err == nil {
					t.Errorf("expected error but got none")
				}

				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if tt.expectedResult.Data == nil {
				if result.Data != nil {
					t.Errorf("expected nil data but got %+v", result.Data)
				}

				return
			}

			if result.Data == nil {
				t.Errorf("expected data but got nil")
				return
			}

			if result.Data.Console != tt.expectedResult.Data.Console {
				t.Errorf("console mismatch: expected %q, got %q", tt.expectedResult.Data.Console, result.Data.Console)
			}

			if result.Data.Release != tt.expectedResult.Data.Release {
				t.Errorf("release mismatch: expected %q, got %q", tt.expectedResult.Data.Release, result.Data.Release)
			}

			if result.Data.RepositoryID != tt.expectedResult.Data.RepositoryID {
				t.Errorf("repository ID mismatch: expected %q, got %q", tt.expectedResult.Data.RepositoryID, result.Data.RepositoryID)
			}

			if result.Data.Version.String() != tt.expectedResult.Data.Version.String() {
				t.Errorf("version mismatch: expected %q, got %q", tt.expectedResult.Data.Version.String(), result.Data.Version.String())
			}
		})
	}
}

func TestProxmoxVersion_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		jsonData       []byte
		expectedError  bool
		expectedResult string
	}{
		{
			name:           "valid version with quotes",
			jsonData:       []byte(`"9.0.4"`),
			expectedError:  false,
			expectedResult: "9.0.4",
		},
		{
			name:           "valid version without quotes",
			jsonData:       []byte(`9.0.4`),
			expectedError:  true, // This should fail because it's invalid JSON for a string
			expectedResult: "",
		},
		{
			name:           "semantic version",
			jsonData:       []byte(`"8.2.1"`),
			expectedError:  false,
			expectedResult: "8.2.1",
		},
		{
			name:          "invalid version format",
			jsonData:      []byte(`"not-a-version"`),
			expectedError: true,
		},
		{
			name:          "empty version",
			jsonData:      []byte(`""`),
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var pv ProxmoxVersion

			err := json.Unmarshal(tt.jsonData, &pv)

			if tt.expectedError {
				if err == nil {
					t.Errorf("expected error but got none")
				}

				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if pv.String() != tt.expectedResult {
				t.Errorf("version mismatch: expected %q, got %q", tt.expectedResult, pv.String())
			}
		})
	}
}

// Helper function to create version objects for testing.
func mustParseVersion(v string) *version.Version {
	parsed, err := version.NewVersion(v)
	if err != nil {
		panic(err)
	}

	return parsed
}
