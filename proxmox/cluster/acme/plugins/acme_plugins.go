/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package plugins

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"sort"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
)

// List returns a list of ACME plugins.
func (c *Client) List(ctx context.Context) ([]*ACMEPluginsListResponseData, error) {
	resBody := &ACMEPluginsListResponseBody{}

	err := c.DoRequest(ctx, http.MethodGet, c.ExpandPath(""), nil, resBody)
	if err != nil {
		return nil, fmt.Errorf("error listing ACME plugins: %w", err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	sort.Slice(resBody.Data, func(i, j int) bool {
		return resBody.Data[i].Plugin < resBody.Data[j].Plugin
	})

	return resBody.Data, nil
}

// Get retrieves a single ACME plugin based on its identifier.
func (c *Client) Get(ctx context.Context, id string) (*ACMEPluginsGetResponseData, error) {
	resBody := &ACMEPluginsGetResponseBody{}

	err := c.DoRequest(ctx, http.MethodGet, c.ExpandPath(url.PathEscape(id)), nil, resBody)
	if err != nil {
		return nil, fmt.Errorf("error reading ACME plugin: %w", err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	return resBody.Data, nil
}

// Create creates a new ACME plugin.
func (c *Client) Create(ctx context.Context, data *ACMEPluginsCreateRequestBody) error {
	err := c.DoRequest(ctx, http.MethodPost, c.ExpandPath(""), data, nil)
	if err != nil {
		return fmt.Errorf("error creating ACME plugin: %w", err)
	}

	return nil
}

// Update updates an existing ACME plugin.
func (c *Client) Update(ctx context.Context, id string, data *ACMEPluginsUpdateRequestBody) error {
	err := c.DoRequest(ctx, http.MethodPut, c.ExpandPath(url.PathEscape(id)), data, nil)
	if err != nil {
		return fmt.Errorf("error updating ACME plugin: %w", err)
	}

	return nil
}

// Delete removes an ACME plugin.
func (c *Client) Delete(ctx context.Context, id string) error {
	err := c.DoRequest(ctx, http.MethodDelete, c.ExpandPath(url.PathEscape(id)), nil, nil)
	if err != nil {
		return fmt.Errorf("error deleting ACME plugin: %w", err)
	}

	return nil
}
