/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmox

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"sort"
)

// CreatePool creates a pool.
func (c *VirtualEnvironmentClient) CreatePool(
	ctx context.Context,
	d *VirtualEnvironmentPoolCreateRequestBody,
) error {
	return c.DoRequest(ctx, hmPOST, "pools", d, nil)
}

// DeletePool deletes a pool.
func (c *VirtualEnvironmentClient) DeletePool(ctx context.Context, id string) error {
	return c.DoRequest(ctx, hmDELETE, fmt.Sprintf("pools/%s", url.PathEscape(id)), nil, nil)
}

// GetPool retrieves a pool.
func (c *VirtualEnvironmentClient) GetPool(
	ctx context.Context,
	id string,
) (*VirtualEnvironmentPoolGetResponseData, error) {
	resBody := &VirtualEnvironmentPoolGetResponseBody{}
	err := c.DoRequest(ctx, hmGET, fmt.Sprintf("pools/%s", url.PathEscape(id)), nil, resBody)
	if err != nil {
		return nil, err
	}

	if resBody.Data == nil {
		return nil, errors.New("the server did not include a data object in the response")
	}

	sort.Slice(resBody.Data.Members, func(i, j int) bool {
		return resBody.Data.Members[i].ID < resBody.Data.Members[j].ID
	})

	return resBody.Data, nil
}

// ListPools retrieves a list of pools.
func (c *VirtualEnvironmentClient) ListPools(
	ctx context.Context,
) ([]*VirtualEnvironmentPoolListResponseData, error) {
	resBody := &VirtualEnvironmentPoolListResponseBody{}
	err := c.DoRequest(ctx, hmGET, "pools", nil, resBody)
	if err != nil {
		return nil, err
	}

	if resBody.Data == nil {
		return nil, errors.New("the server did not include a data object in the response")
	}

	sort.Slice(resBody.Data, func(i, j int) bool {
		return resBody.Data[i].ID < resBody.Data[j].ID
	})

	return resBody.Data, nil
}

// UpdatePool updates a pool.
func (c *VirtualEnvironmentClient) UpdatePool(
	ctx context.Context,
	id string,
	d *VirtualEnvironmentPoolUpdateRequestBody,
) error {
	return c.DoRequest(ctx, hmPUT, fmt.Sprintf("pools/%s", url.PathEscape(id)), d, nil)
}
