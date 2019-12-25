/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmox

// VirtualEnvironmentGroupCreateRequestBody contains the data for an access group create request.
type VirtualEnvironmentGroupCreateRequestBody struct {
	Comment *string `json:"comment,omitempty" url:"comment,omitempty"`
	ID      string  `json:"groupid" url:"groupid"`
}

// VirtualEnvironmentGroupGetResponseBody contains the body from an access group get response.
type VirtualEnvironmentGroupGetResponseBody struct {
	Data *VirtualEnvironmentGroupGetResponseData `json:"data,omitempty"`
}

// VirtualEnvironmentGroupGetResponseData contains the data from an access group get response.
type VirtualEnvironmentGroupGetResponseData struct {
	Comment *string  `json:"comment,omitempty"`
	Members []string `json:"members"`
}

// VirtualEnvironmentGroupListResponseBody contains the body from an access group list response.
type VirtualEnvironmentGroupListResponseBody struct {
	Data []*VirtualEnvironmentGroupListResponseData `json:"data,omitempty"`
}

// VirtualEnvironmentGroupListResponseData contains the data from an access group list response.
type VirtualEnvironmentGroupListResponseData struct {
	Comment *string `json:"comment,omitempty"`
	ID      string  `json:"groupid"`
}

// VirtualEnvironmentGroupUpdateRequestBody contains the data for an access group update request.
type VirtualEnvironmentGroupUpdateRequestBody struct {
	Comment *string `json:"comment,omitempty" url:"comment,omitempty"`
}
