/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package access

import (
	"github.com/bpg/terraform-provider-proxmox/internal/types"
)

// UserChangePasswordRequestBody contains the data for a user password change request.
type UserChangePasswordRequestBody struct {
	ID       string `json:"userid"   url:"userid"`
	Password string `json:"password" url:"password"`
}

// UserCreateRequestBody contains the data for a user create request.
type UserCreateRequestBody struct {
	Comment        *string                `json:"comment,omitempty"   url:"comment,omitempty"`
	Email          *string                `json:"email,omitempty"     url:"email,omitempty"`
	Enabled        *types.CustomBool      `json:"enable,omitempty"    url:"enable,omitempty,int"`
	ExpirationDate *types.CustomTimestamp `json:"expire,omitempty"    url:"expire,omitempty,unix"`
	FirstName      *string                `json:"firstname,omitempty" url:"firstname,omitempty"`
	Groups         []string               `json:"groups,omitempty"    url:"groups,omitempty,comma"`
	ID             string                 `json:"userid"              url:"userid"`
	Keys           *string                `json:"keys,omitempty"      url:"keys,omitempty"`
	LastName       *string                `json:"lastname,omitempty"  url:"lastname,omitempty"`
	Password       string                 `json:"password"            url:"password,omitempty"`
}

// UserGetResponseBody contains the body from a user get response.
type UserGetResponseBody struct {
	Data *UserGetResponseData `json:"data,omitempty"`
}

// UserGetResponseData contains the data from an user get response.
type UserGetResponseData struct {
	Comment        *string                `json:"comment,omitempty"`
	Email          *string                `json:"email,omitempty"`
	Enabled        *types.CustomBool      `json:"enable,omitempty"`
	ExpirationDate *types.CustomTimestamp `json:"expire,omitempty"`
	FirstName      *string                `json:"firstname,omitempty"`
	Groups         *[]string              `json:"groups,omitempty"`
	Keys           *string                `json:"keys,omitempty"`
	LastName       *string                `json:"lastname,omitempty"`
}

// UserListResponseBody contains the body from a user list response.
type UserListResponseBody struct {
	Data []*UserListResponseData `json:"data,omitempty"`
}

// UserListResponseData contains the data from an user list response.
type UserListResponseData struct {
	Comment        *string                `json:"comment,omitempty"`
	Email          *string                `json:"email,omitempty"`
	Enabled        *types.CustomBool      `json:"enable,omitempty"`
	ExpirationDate *types.CustomTimestamp `json:"expire,omitempty"`
	FirstName      *string                `json:"firstname,omitempty"`
	Groups         *[]string              `json:"groups,omitempty"`
	ID             string                 `json:"userid"`
	Keys           *string                `json:"keys,omitempty"`
	LastName       *string                `json:"lastname,omitempty"`
}

// UserUpdateRequestBody contains the data for an user update request.
type UserUpdateRequestBody struct {
	Append         *types.CustomBool      `json:"append,omitempty"    url:"append,omitempty"`
	Comment        *string                `json:"comment,omitempty"   url:"comment,omitempty"`
	Email          *string                `json:"email,omitempty"     url:"email,omitempty"`
	Enabled        *types.CustomBool      `json:"enable,omitempty"    url:"enable,omitempty,int"`
	ExpirationDate *types.CustomTimestamp `json:"expire,omitempty"    url:"expire,omitempty,unix"`
	FirstName      *string                `json:"firstname,omitempty" url:"firstname,omitempty"`
	Groups         []string               `json:"groups,omitempty"    url:"groups,omitempty,comma"`
	Keys           *string                `json:"keys,omitempty"      url:"keys,omitempty"`
	LastName       *string                `json:"lastname,omitempty"  url:"lastname,omitempty"`
}
