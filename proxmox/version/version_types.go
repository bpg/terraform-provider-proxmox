/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package version

// ResponseBody contains the body from a version response.
type ResponseBody struct {
	Data *ResponseData `json:"data,omitempty"`
}

// ResponseData contains the data from a version response.
type ResponseData struct {
	Console      string `json:"console"`
	Release      string `json:"release"`
	RepositoryID string `json:"repoid"`
	Version      string `json:"version"`
}
