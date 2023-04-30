/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package access

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sort"

	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

func (c *Client) rolesPath() string {
	return c.ExpandPath("roles")
}

func (c *Client) rolePath(id string) string {
	return fmt.Sprintf("%s/%s", c.rolesPath(), url.PathEscape(id))
}

// CreateRole creates an access role.
func (c *Client) CreateRole(ctx context.Context, d *RoleCreateRequestBody) error {
	err := c.DoRequest(ctx, http.MethodPost, c.rolesPath(), d, nil)
	if err != nil {
		return fmt.Errorf("error creating role: %w", err)
	}
	return nil
}

// DeleteRole deletes an access role.
func (c *Client) DeleteRole(ctx context.Context, id string) error {
	err := c.DoRequest(ctx, http.MethodDelete, c.rolePath(id), nil, nil)
	if err != nil {
		return fmt.Errorf("error deleting role: %w", err)
	}
	return nil
}

// GetRole retrieves an access role.
func (c *Client) GetRole(ctx context.Context, id string) (*types.CustomPrivileges, error) {
	resBody := &RoleGetResponseBody{}
	err := c.DoRequest(ctx, http.MethodGet, c.rolePath(id), nil, resBody)
	if err != nil {
		return nil, fmt.Errorf("error getting role: %w", err)
	}

	if resBody.Data == nil {
		return nil, errors.New("the server did not include a data object in the response")
	}

	sort.Strings(*resBody.Data)

	return resBody.Data, nil
}

// ListRoles retrieves a list of access roles.
func (c *Client) ListRoles(ctx context.Context) ([]*RoleListResponseData, error) {
	resBody := &RoleListResponseBody{}
	err := c.DoRequest(ctx, http.MethodGet, c.rolesPath(), nil, resBody)
	if err != nil {
		return nil, fmt.Errorf("error listing roles: %w", err)
	}

	if resBody.Data == nil {
		return nil, errors.New("the server did not include a data object in the response")
	}

	sort.Slice(resBody.Data, func(i, j int) bool {
		return resBody.Data[i].ID < resBody.Data[j].ID
	})

	for i := range resBody.Data {
		if resBody.Data[i].Privileges != nil {
			sort.Strings(*resBody.Data[i].Privileges)
		}
	}

	return resBody.Data, nil
}

// UpdateRole updates an access role.
func (c *Client) UpdateRole(ctx context.Context, id string, d *RoleUpdateRequestBody) error {
	err := c.DoRequest(ctx, http.MethodPut, c.rolePath(id), d, nil)
	if err != nil {
		return fmt.Errorf("error updating role: %w", err)
	}
	return nil
}
