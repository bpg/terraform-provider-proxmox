/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package groups

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"sort"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
)

// List retrieves the list of HA groups.
func (c *Client) List(ctx context.Context) ([]*HAGroupListResponseData, error) {
	resBody := &HAGroupListResponseBody{}

	err := c.DoRequest(ctx, http.MethodGet, c.ExpandPath(""), nil, resBody)
	if err != nil {
		return nil, fmt.Errorf("error listing HA groups: %w", err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	sort.Slice(resBody.Data, func(i, j int) bool {
		return resBody.Data[i].ID < resBody.Data[j].ID
	})

	return resBody.Data, nil
}

// Get retrieves a single HA group based on its identifier.
func (c *Client) Get(ctx context.Context, groupID string) (*HAGroupGetResponseData, error) {
	resBody := &HAGroupGetResponseBody{}

	err := c.DoRequest(
		ctx, http.MethodGet,
		c.ExpandPath(url.PathEscape(groupID)), nil, resBody,
	)
	if err != nil {
		return nil, fmt.Errorf("error listing pools: %w", err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	return resBody.Data, nil
}
