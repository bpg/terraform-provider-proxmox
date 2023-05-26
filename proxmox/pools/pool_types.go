/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package pools

// PoolCreateRequestBody contains the data for a pool create request.
type PoolCreateRequestBody struct {
	Comment *string `json:"comment,omitempty" url:"comment,omitempty"`
	ID      string  `json:"groupid"           url:"poolid"`
}

// PoolGetResponseBody contains the body from a pool get response.
type PoolGetResponseBody struct {
	Data *PoolGetResponseData `json:"data,omitempty"`
}

// PoolGetResponseData contains the data from a pool get response.
type PoolGetResponseData struct {
	Comment *string                                    `json:"comment,omitempty"`
	Members []VirtualEnvironmentPoolGetResponseMembers `json:"members,omitempty"`
}

// VirtualEnvironmentPoolGetResponseMembers contains the members data from a pool get response.
type VirtualEnvironmentPoolGetResponseMembers struct {
	ID          string  `json:"id"`
	Node        string  `json:"node"`
	DatastoreID *string `json:"storage,omitempty"`
	Type        string  `json:"type"`
	VMID        *int    `json:"vmid"`
}

// PoolListResponseBody contains the body from a pool list response.
type PoolListResponseBody struct {
	Data []*PoolListResponseData `json:"data,omitempty"`
}

// PoolListResponseData contains the data from a pool list response.
type PoolListResponseData struct {
	Comment *string `json:"comment,omitempty"`
	ID      string  `json:"poolid"`
}

// PoolUpdateRequestBody contains the data for an pool update request.
type PoolUpdateRequestBody struct {
	Comment *string `json:"comment,omitempty" url:"comment,omitempty"`
}
