/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package access

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"time"

	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

func (c *Client) usersPath() string {
	return c.ExpandPath("users")
}

func (c *Client) userPath(id string) string {
	return fmt.Sprintf("%s/%s", c.usersPath(), url.PathEscape(id))
}

// ChangeUserPassword changes a user's password.
func (c *Client) ChangeUserPassword(ctx context.Context, id, password string) error {
	d := UserChangePasswordRequestBody{
		ID:       id,
		Password: password,
	}

	err := c.DoRequest(ctx, http.MethodPut, c.ExpandPath("password"), d, nil)
	if err != nil {
		return fmt.Errorf("error changing user password: %w", err)
	}

	return nil
}

// CreateUser creates a user.
func (c *Client) CreateUser(ctx context.Context, d *UserCreateRequestBody) error {
	err := c.DoRequest(ctx, http.MethodPost, c.usersPath(), d, nil)
	if err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}

	return nil
}

// DeleteUser deletes an  user.
func (c *Client) DeleteUser(ctx context.Context, id string) error {
	err := c.DoRequest(ctx, http.MethodDelete, c.userPath(id), nil, nil)
	if err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}

	return nil
}

// GetUser retrieves a user.
func (c *Client) GetUser(ctx context.Context, id string) (*UserGetResponseData, error) {
	resBody := &UserGetResponseBody{}

	err := c.DoRequest(ctx, http.MethodGet, c.userPath(id), nil, resBody)
	if err != nil {
		return nil, fmt.Errorf("error retrieving user: %w", err)
	}

	if resBody.Data == nil {
		return nil, types.ErrNoDataObjectInResponse
	}

	if resBody.Data.ExpirationDate != nil {
		expirationDate := types.CustomTimestamp(time.Time(*resBody.Data.ExpirationDate).UTC())
		resBody.Data.ExpirationDate = &expirationDate
	}

	if resBody.Data.Groups != nil {
		sort.Strings(*resBody.Data.Groups)
	}

	return resBody.Data, nil
}

// ListUsers retrieves a list of users.
func (c *Client) ListUsers(ctx context.Context) ([]*UserListResponseData, error) {
	resBody := &UserListResponseBody{}

	err := c.DoRequest(ctx, http.MethodGet, c.usersPath(), nil, resBody)
	if err != nil {
		return nil, fmt.Errorf("error listing users: %w", err)
	}

	if resBody.Data == nil {
		return nil, types.ErrNoDataObjectInResponse
	}

	sort.Slice(resBody.Data, func(i, j int) bool {
		return resBody.Data[i].ID < resBody.Data[j].ID
	})

	for i := range resBody.Data {
		if resBody.Data[i].ExpirationDate != nil {
			expirationDate := types.CustomTimestamp(time.Time(*resBody.Data[i].ExpirationDate).UTC())
			resBody.Data[i].ExpirationDate = &expirationDate
		}

		if resBody.Data[i].Groups != nil {
			sort.Strings(*resBody.Data[i].Groups)
		}
	}

	return resBody.Data, nil
}

// UpdateUser updates a user.
func (c *Client) UpdateUser(ctx context.Context, id string, d *UserUpdateRequestBody) error {
	err := c.DoRequest(ctx, http.MethodPut, c.userPath(id), d, nil)
	if err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}

	return nil
}
