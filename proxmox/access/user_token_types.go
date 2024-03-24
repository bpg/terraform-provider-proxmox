/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package access

import (
	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

// UserTokenCreateRequestBody contains the data for a user token create request.
type UserTokenCreateRequestBody struct {
	Comment        *string           `json:"comment,omitempty" url:"comment,omitempty"`
	PrivSeparate   *types.CustomBool `json:"privsep,omitempty" url:"privsep,omitempty,int"`
	ExpirationDate *int64            `json:"expire,omitempty"  url:"expire,omitempty"`
	ID             string            `json:"userid"            url:"userid"`
	TokenID        string            `json:"tokenid"           url:"tokenid"`
}

// UserTokenResponseBody contains the body from a user get token response.
type UserTokenResponseBody struct {
	Data *UserTokenGetResponseData `json:"data,omitempty"`
}

// UserTokenGetResponseData contains the data from an user token response.
type UserTokenGetResponseData struct {
	Comment        *string           `json:"comment,omitempty" url:"comment,omitempty"`
	PrivSeparate   *types.CustomBool `json:"privsep,omitempty" url:"privsep,omitempty,int"`
	ExpirationDate *int64            `json:"expire,omitempty"  url:"expire,omitempty"`
	ID             *string           `json:"userid"            url:"userid"`
	TokenID        *string           `json:"tokenid"           url:"tokenid"`
}

// UserListResponseBody contains the body from a user list response.
type UserTokenListResponseBody struct {
	Data []*UserTokenGetResponseData `json:"data,omitempty"`
}
