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

	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

func (c *Client) groupsPath() string {
	return c.ExpandPath("groups")
}

func (c *Client) groupPath(id string) string {
	return fmt.Sprintf("%s/%s", c.groupsPath(), url.PathEscape(id))
}

// CreateGroup creates an access group.
func (c *Client) CreateGroup(ctx context.Context, d *GroupCreateRequestBody) error {
	err := c.DoRequest(ctx, http.MethodPost, c.groupsPath(), d, nil)
	if err != nil {
		return fmt.Errorf("failed to create access group: %w", err)
	}

	return nil
}

// DeleteGroup deletes an access group.
func (c *Client) DeleteGroup(ctx context.Context, id string) error {
	err := c.DoRequest(ctx, http.MethodDelete, c.groupPath(id), nil, nil)
	if err != nil {
		return fmt.Errorf("failed to delete access group: %w", err)
	}

	return nil
}

// GetGroup retrieves an access group.
func (c *Client) GetGroup(ctx context.Context, id string) (*GroupGetResponseData, error) {
	resBody := &GroupGetResponseBody{}

	err := c.DoRequest(ctx, http.MethodGet, c.groupPath(id), nil, resBody)
	if err != nil {
		return nil, fmt.Errorf("failed to get access group: %w", err)
	}

	if resBody.Data == nil {
		return nil, types.ErrNoDataObjectInResponse
	}

	sort.Strings(resBody.Data.Members)

	return resBody.Data, nil
}

// ListGroups retrieves a list of access groups.
func (c *Client) ListGroups(ctx context.Context) ([]*GroupListResponseData, error) {
	resBody := &GroupListResponseBody{}

	err := c.DoRequest(ctx, http.MethodGet, c.groupsPath(), nil, resBody)
	if err != nil {
		return nil, fmt.Errorf("failed to list access groups: %w", err)
	}

	if resBody.Data == nil {
		return nil, types.ErrNoDataObjectInResponse
	}

	sort.Slice(resBody.Data, func(i, j int) bool {
		return resBody.Data[i].ID < resBody.Data[j].ID
	})

	return resBody.Data, nil
}

// UpdateGroup updates an access group.
func (c *Client) UpdateGroup(ctx context.Context, id string, d *GroupUpdateRequestBody) error {
	err := c.DoRequest(ctx, http.MethodPut, c.groupPath(id), d, nil)
	if err != nil {
		return fmt.Errorf("failed to update access group: %w", err)
	}

	return nil
}
