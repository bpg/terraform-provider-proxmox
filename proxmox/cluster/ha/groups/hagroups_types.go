/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package groups

import (
	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

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

// HAGroupDataBase contains fields which are both received from and send to the HA group API.
type HAGroupDataBase struct {
	// A SHA1 digest of the group's configuration.
	Digest *string `json:"digest,omitempty" url:"digest,omitempty"`
	// The group's comment, if defined
	Comment *string `json:"comment,omitempty" url:"comment,omitempty"`
	// A comma-separated list of node fields. Each node field contains a node name, and may
	// include a priority, with a semicolon acting as a separator.
	Nodes string `json:"nodes" url:"nodes"`
	// A boolean (0/1) indicating that failing back to the highest priority node is disabled.
	NoFailback types.CustomBool `json:"nofailback" url:"nofailback,int"`
	// A boolean (0/1) indicating that associated resources cannot run on other nodes.
	Restricted types.CustomBool `json:"restricted" url:"restricted,int"`
}

// HAGroupGetResponseData contains the data from a HA group get response.
type HAGroupGetResponseData struct {
	// The group's data
	HAGroupDataBase

	// The group's identifier
	ID string `json:"group"`
	// The type. Always set to `group`.
	Type string `json:"type"`
}

// HAGroupCreateRequestBody contains the data which must be sent when creating a HA group.
type HAGroupCreateRequestBody struct {
	// The group's data
	HAGroupDataBase

	// The group's identifier
	ID string `url:"group"`
	// The type. Always set to `group`.
	Type string `url:"type"`
}

// HAGroupUpdateRequestBody contains the data which must be sent when updating a HA group.
type HAGroupUpdateRequestBody struct {
	// The group's data
	HAGroupDataBase

	// A list of settings to delete
	Delete string `url:"delete"`
}
