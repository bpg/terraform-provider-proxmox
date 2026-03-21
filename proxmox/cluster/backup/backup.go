/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package backup

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"sort"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
)

// List retrieves all backup jobs, sorted by ID.
func (c *Client) List(ctx context.Context) ([]*GetResponseData, error) {
	resBody := &ListResponseBody{}

	err := c.DoRequest(ctx, http.MethodGet, c.ExpandPath(""), nil, resBody)
	if err != nil {
		return nil, fmt.Errorf("error listing backup jobs: %w", err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	sort.Slice(resBody.Data, func(i, j int) bool {
		return resBody.Data[i].ID < resBody.Data[j].ID
	})

	return resBody.Data, nil
}

// Get retrieves a single backup job by ID.
func (c *Client) Get(ctx context.Context, id string) (*GetResponseData, error) {
	resBody := &GetResponseBody{}

	err := c.DoRequest(ctx, http.MethodGet, c.ExpandPath(url.PathEscape(id)), nil, resBody)
	if err != nil {
		return nil, fmt.Errorf("error reading backup job: %w", err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	return resBody.Data, nil
}

// Create creates a new backup job.
func (c *Client) Create(ctx context.Context, data *CreateRequestBody) error {
	err := c.DoRequest(ctx, http.MethodPost, c.ExpandPath(""), data, nil)
	if err != nil {
		return fmt.Errorf("error creating backup job: %w", err)
	}

	return nil
}

// Update updates an existing backup job.
func (c *Client) Update(ctx context.Context, id string, data *UpdateRequestBody) error {
	err := c.DoRequest(ctx, http.MethodPut, c.ExpandPath(url.PathEscape(id)), data, nil)
	if err != nil {
		return fmt.Errorf("error updating backup job: %w", err)
	}

	return nil
}

// Delete deletes a backup job.
func (c *Client) Delete(ctx context.Context, id string) error {
	err := c.DoRequest(ctx, http.MethodDelete, c.ExpandPath(url.PathEscape(id)), nil, nil)
	if err != nil {
		return fmt.Errorf("error deleting backup job: %w", err)
	}

	return nil
}
