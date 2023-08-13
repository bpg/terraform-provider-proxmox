/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package groups

// HAGroupListResponseBody contains the body from a HA group list response.
type HAGroupListResponseBody struct {
	Data []*HAGroupListResponseData `json:"data,omitempty"`
}

// HAGroupListResponseData contains the data from a HA group list response.
type HAGroupListResponseData struct {
	ID string `json:"group"`
}
