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
	ExpirationDate *int64            `json:"expire,omitempty"  url:"expire,omitempty"`
	PrivSeparate   *types.CustomBool `json:"privsep,omitempty" url:"privsep,omitempty,int"`
}

// UserTokenUpdateRequestBody contains the data for a user token update request.
type UserTokenUpdateRequestBody UserTokenCreateRequestBody

// UserTokenCreateResponseBody contains the body from a user token create response.
type UserTokenCreateResponseBody struct {
	Data *UserTokenCreateResponseData `json:"data,omitempty"`
}

// UserTokenCreateResponseData contains the data from a user token create response.
type UserTokenCreateResponseData struct {
	// The full token id, format "<userid>!<tokenid>"
	FullTokenID string                   `json:"full-tokenid"`
	Info        UserTokenGetResponseData `json:"info"`
	Value       string                   `json:"value"`
}

// UserTokenGetResponseBody contains the body from a user token get response.
type UserTokenGetResponseBody struct {
	Data *UserTokenGetResponseData `json:"data,omitempty"`
}

// UserTokenGetResponseData contains the data from a user token get response.
type UserTokenGetResponseData struct {
	Comment        *string           `json:"comment,omitempty" url:"comment,omitempty"`
	PrivSeparate   *types.CustomBool `json:"privsep,omitempty" url:"privsep,omitempty,int"`
	ExpirationDate *int64            `json:"expire,omitempty"  url:"expire,omitempty"`
}

// UserTokenListResponseBody contains the body from a user token list response.
type UserTokenListResponseBody struct {
	Data []*UserTokenListResponseData `json:"data,omitempty"`
}

// UserTokenListResponseData contains the data from a user token list response.
type UserTokenListResponseData struct {
	UserTokenGetResponseData

	TokenID string `json:"tokenid"`
}
