/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package storage

// OCIRegistryPullResponseBody contains the response for an OCI registry pull request.
type OCIRegistryPullResponseBody struct {
	TaskID *string `json:"data,omitempty"`
}

// OCIRegistryPullRequestBody contains the data for an OCI registry pull request.
type OCIRegistryPullRequestBody struct {
	FileName  *string `json:"filename,omitempty"  url:"filename,omitempty"`
	Reference *string `json:"reference,omitempty" url:"reference,omitempty"`
	Storage   *string `json:"storage,omitempty"   url:"storage,omitempty"`
}
