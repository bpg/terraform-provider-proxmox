/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package rules

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"sort"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
)

// List retrieves the list of HA rules.
func (c *Client) List(ctx context.Context) ([]*HARuleListResponseData, error) {
	resBody := &HARuleListResponseBody{}

	err := c.DoRequest(ctx, http.MethodGet, c.ExpandPath(""), nil, resBody)
	if err != nil {
		return nil, fmt.Errorf("error listing HA rules: %w", err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	sort.Slice(resBody.Data, func(i, j int) bool {
		return resBody.Data[i].Rule < resBody.Data[j].Rule
	})

	return resBody.Data, nil
}

// Get retrieves a single HA rule based on its identifier.
func (c *Client) Get(ctx context.Context, ruleID string) (*HARuleGetResponseData, error) {
	resBody := &HARuleGetResponseBody{}

	err := c.DoRequest(
		ctx, http.MethodGet,
		c.ExpandPath(url.PathEscape(ruleID)), nil, resBody,
	)
	if err != nil {
		return nil, fmt.Errorf("error reading HA rule: %w", err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	return resBody.Data, nil
}

// Create creates a new HA rule.
func (c *Client) Create(ctx context.Context, data *HARuleCreateRequestBody) error {
	err := c.DoRequest(ctx, http.MethodPost, c.ExpandPath(""), data, nil)
	if err != nil {
		return fmt.Errorf("error creating HA rule: %w", err)
	}

	return nil
}

// Update updates a HA rule's configuration.
func (c *Client) Update(ctx context.Context, ruleID string, data *HARuleUpdateRequestBody) error {
	err := c.DoRequest(ctx, http.MethodPut, c.ExpandPath(url.PathEscape(ruleID)), data, nil)
	if err != nil {
		return fmt.Errorf("error updating HA rule: %w", err)
	}

	return nil
}

// Delete deletes a HA rule.
func (c *Client) Delete(ctx context.Context, ruleID string) error {
	err := c.DoRequest(ctx, http.MethodDelete, c.ExpandPath(url.PathEscape(ruleID)), nil, nil)
	if err != nil {
		return fmt.Errorf("error deleting HA rule: %w", err)
	}

	return nil
}
