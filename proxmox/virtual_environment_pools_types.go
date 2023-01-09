/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmox

// VirtualEnvironmentPoolCreateRequestBody contains the data for an pool create request.
type VirtualEnvironmentPoolCreateRequestBody struct {
	Comment *string `json:"comment,omitempty" url:"comment,omitempty"`
	ID      string  `json:"groupid"           url:"poolid"`
}

// VirtualEnvironmentPoolGetResponseBody contains the body from an pool get response.
type VirtualEnvironmentPoolGetResponseBody struct {
	Data *VirtualEnvironmentPoolGetResponseData `json:"data,omitempty"`
}

// VirtualEnvironmentPoolGetResponseData contains the data from an pool get response.
type VirtualEnvironmentPoolGetResponseData struct {
	Comment *string                                    `json:"comment,omitempty"`
	Members []VirtualEnvironmentPoolGetResponseMembers `json:"members,omitempty"`
}

// VirtualEnvironmentPoolGetResponseMembers contains the members data from an pool get response.
type VirtualEnvironmentPoolGetResponseMembers struct {
	ID          string  `json:"id"`
	Node        string  `json:"node"`
	DatastoreID *string `json:"storage,omitempty"`
	Type        string  `json:"type"`
	VMID        *int    `json:"vmid"`
}

// VirtualEnvironmentPoolListResponseBody contains the body from an pool list response.
type VirtualEnvironmentPoolListResponseBody struct {
	Data []*VirtualEnvironmentPoolListResponseData `json:"data,omitempty"`
}

// VirtualEnvironmentPoolListResponseData contains the data from an pool list response.
type VirtualEnvironmentPoolListResponseData struct {
	Comment *string `json:"comment,omitempty"`
	ID      string  `json:"poolid"`
}

// VirtualEnvironmentPoolUpdateRequestBody contains the data for an pool update request.
type VirtualEnvironmentPoolUpdateRequestBody struct {
	Comment *string `json:"comment,omitempty" url:"comment,omitempty"`
}
