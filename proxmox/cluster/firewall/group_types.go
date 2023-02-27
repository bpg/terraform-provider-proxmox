/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package firewall

// GroupCreateRequestBody contains the data for an security group create request.
type GroupCreateRequestBody struct {
	Group   string  `json:"group"             url:"group"`
	Comment *string `json:"comment,omitempty" url:"comment,omitempty"`
	Digest  *string `json:"digest,omitempty"  url:"digest,omitempty"`
}

// GroupGetResponseBody contains the body from a group get response.
type GroupGetResponseBody struct {
	Data *GroupGetResponseData `json:"data,omitempty"`
}

// GroupGetResponseData contains the data from a group get response.
type GroupGetResponseData struct {
	Comment *string `json:"comment,omitempty" url:"comment,omitempty"`
	Group   string  `json:"group"             url:"group"`
	Digest  string  `json:"digest"            url:"digest"`
}

// GroupListResponseBody contains the data from a group get response.
type GroupListResponseBody struct {
	Data []*GroupGetResponseData `json:"data,omitempty"`
}

// GroupUpdateRequestBody contains the data for a group update request.
type GroupUpdateRequestBody struct {
	Comment *string `json:"comment,omitempty" url:"comment,omitempty"`
	Group   string  `json:"group"             url:"group"`
	ReName  string  `json:"rename"            url:"rename"`
	Digest  *string `json:"digest,omitempty"  url:"digest,omitempty"`
}
