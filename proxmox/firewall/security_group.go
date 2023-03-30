/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package firewall

import "context"

type SecurityGroup interface {
	CreateGroup(ctx context.Context, d *GroupCreateRequestBody) error
	ListGroups(ctx context.Context) ([]*GroupListResponseData, error)
	UpdateGroup(ctx context.Context, d *GroupUpdateRequestBody) error
	DeleteGroup(ctx context.Context, group string) error
}

// GroupCreateRequestBody contains the data for an security group create request.
type GroupCreateRequestBody struct {
	Group   string  `json:"group"             url:"group"`
	Comment *string `json:"comment,omitempty" url:"comment,omitempty"`
	Digest  *string `json:"digest,omitempty"  url:"digest,omitempty"`
}

// GroupListResponseData contains the data from a group list response.
type GroupListResponseData struct {
	Comment *string `json:"comment,omitempty" url:"comment,omitempty"`
	Group   string  `json:"group"             url:"group"`
	Digest  string  `json:"digest"            url:"digest"`
}

// GroupListResponseBody contains the data from a group get response.
type GroupListResponseBody struct {
	Data []*GroupListResponseData `json:"data,omitempty"`
}

// GroupUpdateRequestBody contains the data for a group update request.
type GroupUpdateRequestBody struct {
	Group string `json:"group"             url:"group"`

	Comment *string `json:"comment,omitempty" url:"comment,omitempty"`
	ReName  *string `json:"rename,omitempty"  url:"rename,omitempty"`
	Digest  *string `json:"digest,omitempty"  url:"digest,omitempty"`
}
