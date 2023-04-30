/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package proxmox

// VersionResponseBody contains the body from a version response.
type VersionResponseBody struct {
	Data *VersionResponseData `json:"data,omitempty"`
}

// VersionResponseData contains the data from a version response.
type VersionResponseData struct {
	Keyboard     string `json:"keyboard"`
	Release      string `json:"release"`
	RepositoryID string `json:"repoid"`
	Version      string `json:"version"`
}
