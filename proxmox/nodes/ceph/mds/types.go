/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package mds

// ListResponseBody contains the body from a Ceph MDS list response.
type ListResponseBody struct {
	Data []*ListResponseData `json:"data,omitempty"`
}

// ListResponseData contains the data for one MDS in a list response.
type ListResponseData struct {
	Name string  `json:"name"`
	Host *string `json:"host,omitempty"`
}

// CreateResponseBody contains the body from a Ceph MDS create response.
type CreateResponseBody struct {
	Data *string `json:"data,omitempty"`
}

// DeleteResponseBody contains the body from a Ceph MDS delete response.
type DeleteResponseBody struct {
	Data *string `json:"data,omitempty"`
}
