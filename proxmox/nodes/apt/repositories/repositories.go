/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package repositories

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
)

// Client is an interface for accessing the Proxmox node APT repositories API.
type Client struct {
	api.Client
}

// basePath returns the expanded APT repositories API base path.
func (c *Client) basePath() string {
	return c.Client.ExpandPath("repositories")
}

// Add adds an APT standard repository entry.
func (c *Client) Add(ctx context.Context, data *AddRequestBody) error {
	err := c.DoRequest(ctx, http.MethodPut, c.ExpandPath(), data, nil)
	if err != nil {
		return fmt.Errorf("adding APT standard repository: %w", err)
	}

	return nil
}

// ExpandPath expands a relative path to a full APT repositories API path.
func (c *Client) ExpandPath() string {
	return c.basePath()
}

// Get retrieves all APT repositories.
func (c *Client) Get(ctx context.Context) (*GetResponseData, error) {
	resBody := &GetResponseBody{}

	err := c.DoRequest(ctx, http.MethodGet, c.ExpandPath(), nil, resBody)
	if err != nil {
		return nil, fmt.Errorf("reading APT repositories: %w", err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	return resBody.Data, nil
}

// Modify modifies the activation status of an APT repository.
func (c *Client) Modify(ctx context.Context, data *ModifyRequestBody) error {
	err := c.DoRequest(ctx, http.MethodPost, c.ExpandPath(), data, nil)
	if err != nil {
		return fmt.Errorf(
			`modifying APT repository in file %s at index %d to activation state %v: %w`,
			data.Path,
			data.Index,
			data.Enabled,
			err,
		)
	}

	return nil
}
