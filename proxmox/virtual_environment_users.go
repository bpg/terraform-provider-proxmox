/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmox

import (
	"errors"
	"fmt"
	"net/url"
	"sort"
)

// VirtualEnvironmentUserCreateRequestBody contains the data for an user create request.
type VirtualEnvironmentUserCreateRequestBody struct {
	Comment        *string          `json:"comment,omitempty"`
	Email          *string          `json:"email,omitempty"`
	Enabled        *CustomBool      `json:"enable,omitempty"`
	ExpirationDate *CustomTimestamp `json:"expire,omitempty"`
	FirstName      *string          `json:"firstname,omitempty"`
	Groups         *[]string        `json:"groups,omitempty"`
	ID             string           `json:"userid"`
	Keys           *string          `json:"keys,omitempty"`
	LastName       *string          `json:"lastname,omitempty"`
	Password       string           `json:"password"`
}

// VirtualEnvironmentUserGetResponseBody contains the body from an user get response.
type VirtualEnvironmentUserGetResponseBody struct {
	Data *VirtualEnvironmentUserGetResponseData `json:"data,omitempty"`
}

// VirtualEnvironmentUserGetResponseData contains the data from an user get response.
type VirtualEnvironmentUserGetResponseData struct {
	Comment        *string          `json:"comment,omitempty"`
	Email          *string          `json:"email,omitempty"`
	Enabled        *CustomBool      `json:"enable,omitempty"`
	ExpirationDate *CustomTimestamp `json:"expire,omitempty"`
	FirstName      *string          `json:"firstname,omitempty"`
	Groups         *[]string        `json:"groups,omitempty"`
	Keys           *string          `json:"keys,omitempty"`
	LastName       *string          `json:"lastname,omitempty"`
}

// VirtualEnvironmentUserListResponseBody contains the body from an user list response.
type VirtualEnvironmentUserListResponseBody struct {
	Data []*VirtualEnvironmentUserListResponseData `json:"data,omitempty"`
}

// VirtualEnvironmentUserListResponseData contains the data from an user list response.
type VirtualEnvironmentUserListResponseData struct {
	Comment        *string          `json:"comment,omitempty"`
	Email          *string          `json:"email,omitempty"`
	Enabled        *CustomBool      `json:"enable,omitempty"`
	ExpirationDate *CustomTimestamp `json:"expire,omitempty"`
	FirstName      *string          `json:"firstname,omitempty"`
	Groups         *[]string        `json:"groups,omitempty"`
	ID             string           `json:"userid"`
	Keys           *string          `json:"keys,omitempty"`
	LastName       *string          `json:"lastname,omitempty"`
}

// VirtualEnvironmentUserUpdateRequestBody contains the data for an user update request.
type VirtualEnvironmentUserUpdateRequestBody struct {
	Append         *CustomBool      `json:"append,omitempty"`
	Comment        *string          `json:"comment,omitempty"`
	Email          *string          `json:"email,omitempty"`
	Enabled        *CustomBool      `json:"enable,omitempty"`
	ExpirationDate *CustomTimestamp `json:"expire,omitempty"`
	FirstName      *string          `json:"firstname,omitempty"`
	Groups         *[]string        `json:"groups,omitempty"`
	Keys           *string          `json:"keys,omitempty"`
	LastName       *string          `json:"lastname,omitempty"`
}

// CreateUser creates an user.
func (c *VirtualEnvironmentClient) CreateUser(d *VirtualEnvironmentUserCreateRequestBody) error {
	return c.DoRequest(hmPOST, "access/users", d, nil)
}

// DeleteUser deletes an user.
func (c *VirtualEnvironmentClient) DeleteUser(id string) error {
	return c.DoRequest(hmDELETE, fmt.Sprintf("access/users/%s", id), nil, nil)
}

// GetUser retrieves an user.
func (c *VirtualEnvironmentClient) GetUser(id string) (*VirtualEnvironmentUserGetResponseData, error) {
	resBody := &VirtualEnvironmentUserGetResponseBody{}
	err := c.DoRequest(hmGET, fmt.Sprintf("access/users/%s", url.PathEscape(id)), nil, resBody)

	if err != nil {
		return nil, err
	}

	if resBody.Data == nil {
		return nil, errors.New("The server did not include a data object in the response")
	}

	return resBody.Data, nil
}

// ListUsers retrieves a list of users.
func (c *VirtualEnvironmentClient) ListUsers() ([]*VirtualEnvironmentUserListResponseData, error) {
	resBody := &VirtualEnvironmentUserListResponseBody{}
	err := c.DoRequest(hmGET, "access/users", nil, resBody)

	if err != nil {
		return nil, err
	}

	if resBody.Data == nil {
		return nil, errors.New("The server did not include a data object in the response")
	}

	sort.Slice(resBody.Data, func(i, j int) bool {
		return resBody.Data[i].ID < resBody.Data[j].ID
	})

	return resBody.Data, nil
}

// UpdateUser updates an user.
func (c *VirtualEnvironmentClient) UpdateUser(id string, d *VirtualEnvironmentUserUpdateRequestBody) error {
	return c.DoRequest(hmPUT, fmt.Sprintf("access/users/%s", id), d, nil)
}
