/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmox

import (
	"errors"
	"fmt"
	"net/url"
	"sort"
	"time"
)

// ChangeUserPassword changes a user's password.
func (c *VirtualEnvironmentClient) ChangeUserPassword(id, password string) error {
	d := VirtualEnvironmentUserChangePasswordRequestBody{
		ID:       id,
		Password: password,
	}

	return c.DoRequest(hmPUT, "access/password", d, nil)
}

// CreateUser creates an user.
func (c *VirtualEnvironmentClient) CreateUser(d *VirtualEnvironmentUserCreateRequestBody) error {
	return c.DoRequest(hmPOST, "access/users", d, nil)
}

// DeleteUser deletes an user.
func (c *VirtualEnvironmentClient) DeleteUser(id string) error {
	return c.DoRequest(hmDELETE, fmt.Sprintf("access/users/%s", url.PathEscape(id)), nil, nil)
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

	if resBody.Data.ExpirationDate != nil {
		expirationDate := CustomTimestamp(time.Time(*resBody.Data.ExpirationDate).UTC())
		resBody.Data.ExpirationDate = &expirationDate
	}

	if resBody.Data.Groups != nil {
		sort.Strings(*resBody.Data.Groups)
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

	for i := range resBody.Data {
		if resBody.Data[i].ExpirationDate != nil {
			expirationDate := CustomTimestamp(time.Time(*resBody.Data[i].ExpirationDate).UTC())
			resBody.Data[i].ExpirationDate = &expirationDate
		}

		if resBody.Data[i].Groups != nil {
			sort.Strings(*resBody.Data[i].Groups)
		}
	}

	return resBody.Data, nil
}

// UpdateUser updates an user.
func (c *VirtualEnvironmentClient) UpdateUser(id string, d *VirtualEnvironmentUserUpdateRequestBody) error {
	return c.DoRequest(hmPUT, fmt.Sprintf("access/users/%s", url.PathEscape(id)), d, nil)
}
