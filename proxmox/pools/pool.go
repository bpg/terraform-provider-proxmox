/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package pools

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"sort"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
)

// CreatePool creates a pool.
func (c *Client) CreatePool(ctx context.Context, d *PoolCreateRequestBody) error {
	err := c.DoRequest(ctx, http.MethodPost, "pools", d, nil)
	if err != nil {
		return fmt.Errorf("error creating pool: %w", err)
	}

	return nil
}

// DeletePool deletes a pool.
func (c *Client) DeletePool(ctx context.Context, id string) error {
	err := c.DoRequest(ctx, http.MethodDelete, fmt.Sprintf("pools?poolid=%s", url.PathEscape(id)), nil, nil)
	if err != nil {
		return fmt.Errorf("error deleting pool: %w", err)
	}

	return nil
}

// GetPool retrieves a pool.
func (c *Client) GetPool(ctx context.Context, id string) (*PoolGetResponseData, error) {
	resBody := &PoolGetResponseBody{}

	err := c.DoRequest(ctx, http.MethodGet, fmt.Sprintf("pools?poolid=%s", url.PathEscape(id)), nil, resBody)
	if err != nil {
		return nil, fmt.Errorf("error getting pool: %w", err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	// Should always be with one single element
	if len(resBody.Data) == 0 {
		return nil, api.ErrResourceDoesNotExist
	}

	if len(resBody.Data) > 1 {
		return nil, fmt.Errorf("error getting pool: expected 1 pool, but got %d", len(resBody.Data))
	}

	sort.Slice(resBody.Data[0].Members, func(i, j int) bool {
		return resBody.Data[0].Members[i].ID < resBody.Data[0].Members[j].ID
	})

	return resBody.Data[0], nil
}

// ListPools retrieves a list of pools.
func (c *Client) ListPools(ctx context.Context) ([]*PoolListResponseData, error) {
	resBody := &PoolListResponseBody{}

	err := c.DoRequest(ctx, http.MethodGet, "pools", nil, resBody)
	if err != nil {
		return nil, fmt.Errorf("error listing pools: %w", err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	sort.Slice(resBody.Data, func(i, j int) bool {
		return resBody.Data[i].ID < resBody.Data[j].ID
	})

	return resBody.Data, nil
}

// UpdatePool updates a pool.
func (c *Client) UpdatePool(ctx context.Context, id string, d *PoolUpdateRequestBody) error {
	err := c.DoRequest(ctx, http.MethodPut, fmt.Sprintf("pools?poolid=%s", url.PathEscape(id)), d, nil)
	if err != nil {
		return fmt.Errorf("error updating pool: %w", err)
	}

	return nil
}
