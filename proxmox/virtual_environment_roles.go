/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmox

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sort"

	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

// CreateRole creates an access role.
func (c *VirtualEnvironmentClient) CreateRole(
	ctx context.Context,
	d *VirtualEnvironmentRoleCreateRequestBody,
) error {
	return c.DoRequest(ctx, http.MethodPost, "access/roles", d, nil)
}

// DeleteRole deletes an access role.
func (c *VirtualEnvironmentClient) DeleteRole(ctx context.Context, id string) error {
	return c.DoRequest(ctx, http.MethodDelete, fmt.Sprintf("access/roles/%s", url.PathEscape(id)), nil, nil)
}

// GetRole retrieves an access role.
func (c *VirtualEnvironmentClient) GetRole(
	ctx context.Context,
	id string,
) (*types.CustomPrivileges, error) {
	resBody := &VirtualEnvironmentRoleGetResponseBody{}
	err := c.DoRequest(ctx, http.MethodGet, fmt.Sprintf("access/roles/%s", url.PathEscape(id)), nil, resBody)
	if err != nil {
		return nil, err
	}

	if resBody.Data == nil {
		return nil, errors.New("the server did not include a data object in the response")
	}

	sort.Strings(*resBody.Data)

	return resBody.Data, nil
}

// ListRoles retrieves a list of access roles.
func (c *VirtualEnvironmentClient) ListRoles(
	ctx context.Context,
) ([]*VirtualEnvironmentRoleListResponseData, error) {
	resBody := &VirtualEnvironmentRoleListResponseBody{}
	err := c.DoRequest(ctx, http.MethodGet, "access/roles", nil, resBody)
	if err != nil {
		return nil, err
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
func (c *VirtualEnvironmentClient) UpdateRole(
	ctx context.Context,
	id string,
	d *VirtualEnvironmentRoleUpdateRequestBody,
) error {
	return c.DoRequest(ctx, http.MethodPut, fmt.Sprintf("access/roles/%s", url.PathEscape(id)), d, nil)
}
