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

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
)

func (c *Client) userTokensPath(id string) string {
	return fmt.Sprintf("%s/%s/token", c.usersPath(), url.PathEscape(id))
}

func (c *Client) userTokenPath(userid, id string) string {
	return fmt.Sprintf("%s/%s", c.userTokensPath(userid), url.PathEscape(id))
}

// CreateUserToken creates a user token.
func (c *Client) CreateUserToken(ctx context.Context, userid string, id string, d *UserTokenCreateRequestBody) error {
	err := c.DoRequest(ctx, http.MethodPost, c.userTokenPath(userid, id), d, nil)
	if err != nil {
		return fmt.Errorf("error creating user token: %w", err)
	}

	return nil
}

// DeleteUserToken deletes an user token.
func (c *Client) DeleteUserToken(ctx context.Context, userid string, id string) error {
	err := c.DoRequest(ctx, http.MethodDelete, c.userTokenPath(userid, id), nil, nil)
	if err != nil {
		return fmt.Errorf("error deleting user token: %w", err)
	}

	return nil
}

// GetUserToken retrieves a user token.
func (c *Client) GetUserToken(ctx context.Context, userid string, id string) (*UserTokenGetResponseData, error) {
	resBody := &UserTokenResponseBody{}

	err := c.DoRequest(ctx, http.MethodGet, c.userTokenPath(userid, id), nil, resBody)
	if err != nil {
		return nil, fmt.Errorf("error retrieving user: %w", err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	return resBody.Data, nil
}

// ListUserTokens retrieves a list of user tokens.
func (c *Client) ListUserTokens(ctx context.Context, userid string) ([]*UserTokenGetResponseData, error) {
	resBody := &UserTokenListResponseBody{}

	err := c.DoRequest(ctx, http.MethodGet, c.userTokensPath(userid), nil, resBody)
	if err != nil {
		return nil, fmt.Errorf("error listing users: %w", err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	return resBody.Data, nil
}

// UpdateUserToken updates the user token.
func (c *Client) UpdateUserToken(ctx context.Context, userid string, id string, d *UserTokenCreateRequestBody) error {
	err := c.DoRequest(ctx, http.MethodPut, c.userTokenPath(userid, id), d, nil)
	if err != nil {
		return fmt.Errorf("error updating user token: %w", err)
	}

	return nil
}
