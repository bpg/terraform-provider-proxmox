/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package account

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"sort"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
)

// List returns a list of ACME accounts.
func (c *Client) List(ctx context.Context) ([]*ACMEAccountListResponseData, error) {
	resBody := &ACMEAccountListResponseBody{}

	err := c.DoRequest(ctx, http.MethodGet, c.ExpandPath(""), nil, resBody)
	if err != nil {
		return nil, fmt.Errorf("error listing ACME accounts: %w", err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	sort.Slice(resBody.Data, func(i, j int) bool {
		return resBody.Data[i].Name < resBody.Data[j].Name
	})

	return resBody.Data, nil
}

// Get retrieves a single ACME account based on its identifier.
func (c *Client) Get(ctx context.Context, name string) (*ACMEAccountGetResponseData, error) {
	resBody := &ACMEAccountGetResponseBody{}

	err := c.DoRequest(ctx, http.MethodGet, c.ExpandPath(url.PathEscape(name)), nil, resBody)
	if err != nil {
		return nil, fmt.Errorf("error reading ACME account: %w", err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	return resBody.Data, nil
}

// Create creates a new ACME account.
func (c *Client) Create(ctx context.Context, data *ACMEAccountCreateRequestBody) error {
	err := c.DoRequest(ctx, http.MethodPost, c.ExpandPath(""), data, nil)
	if err != nil {
		return fmt.Errorf("error creating ACME account: %w", err)
	}

	return nil
}

// Update updates an existing ACME account.
func (c *Client) Update(ctx context.Context, accountName string, data *ACMEAccountUpdateRequestBody) error {
	err := c.DoRequest(ctx, http.MethodPut, c.ExpandPath(url.PathEscape(accountName)), data, nil)
	if err != nil {
		return fmt.Errorf("error updating ACME account: %w", err)
	}

	return nil
}

// Delete removes an ACME account.
func (c *Client) Delete(ctx context.Context, accountName string) error {
	err := c.DoRequest(ctx, http.MethodDelete, c.ExpandPath(url.PathEscape(accountName)), nil, nil)
	if err != nil {
		return fmt.Errorf("error deleting ACME account: %w", err)
	}

	return nil
}
