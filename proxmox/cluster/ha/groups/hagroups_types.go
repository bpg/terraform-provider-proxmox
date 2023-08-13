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

// HAGroupGetResponseBody contains the body from a HA group get response.
type HAGroupGetResponseBody struct {
	Data *HAGroupGetResponseData `json:"data,omitempty"`
}

// HAGroupGetResponseData contains the data from a HA group get response.
type HAGroupGetResponseData struct {
	// The group's identifier
	ID string `json:"group"`
	// A digest of the group's configuration
	Digest string `json:"digest"`
	// The group's comment, if defined
	Comment *string `json:"comment,omitempty"`
	// A comma-separated list of node fields. Each node field contains a node name, and may
	// include a priority, with a semicolon acting as a separator.
	Nodes string `json:"nodes"`
	// A boolean (0/1) indicating that failing back to the highest priority node is disabled.
	NoFailback int `json:"nofailback"`
	// A boolean (0/1) indicating that associated resources cannot run on other nodes.
	Restricted int `json:"restricted"`
	// The type. Always set to `group`.
	Type string `json:"type"`
}
