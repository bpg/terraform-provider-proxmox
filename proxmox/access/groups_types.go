/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package access

// GroupCreateRequestBody contains the data for an access group create request.
type GroupCreateRequestBody struct {
	Comment *string `json:"comment,omitempty" url:"comment,omitempty"`
	ID      string  `json:"groupid"           url:"groupid"`
}

// GroupGetResponseBody contains the body from an access group get response.
type GroupGetResponseBody struct {
	Data *GroupGetResponseData `json:"data,omitempty"`
}

// GroupGetResponseData contains the data from an access group get response.
type GroupGetResponseData struct {
	Comment *string  `json:"comment,omitempty"`
	Members []string `json:"members"`
}

// GroupListResponseBody contains the body from an access group list response.
type GroupListResponseBody struct {
	Data []*GroupListResponseData `json:"data,omitempty"`
}

// GroupListResponseData contains the data from an access group list response.
type GroupListResponseData struct {
	Comment *string `json:"comment,omitempty"`
	ID      string  `json:"groupid"`
}

// GroupUpdateRequestBody contains the data for an access group update request.
type GroupUpdateRequestBody struct {
	Comment *string `json:"comment,omitempty" url:"comment,omitempty"`
}
