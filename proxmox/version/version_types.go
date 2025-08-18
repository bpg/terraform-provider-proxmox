/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package version

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-version"
)

// ResponseBody contains the body from a version response.
type ResponseBody struct {
	Data *ResponseData `json:"data,omitempty"`
}

// ResponseData contains the data from a version response.
type ResponseData struct {
	Console      string         `json:"console"`
	Release      string         `json:"release"`
	RepositoryID string         `json:"repoid"`
	Version      ProxmoxVersion `json:"version"`
}

type ProxmoxVersion struct {
	version.Version
}

func (v *ProxmoxVersion) UnmarshalJSON(data []byte) error {
	// Unmarshal the version string into a go-version Version object, remove wrapping quotes if any
	ver, err := version.NewVersion(strings.Trim(string(data), "\""))
	if err != nil {
		return fmt.Errorf("failed to parse version %q: %w", string(data), err)
	}

	v.Version = *ver

	return nil
}
