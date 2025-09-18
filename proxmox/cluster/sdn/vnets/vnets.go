/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package vnets

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/sdn"
)

// GetVnet retrieves a single SDN Vnet by ID.
func (c *Client) GetVnet(ctx context.Context) (*VnetData, error) {
	return c.GetVnetWithParams(ctx, nil)
}

// GetVnet retrieves a single SDN Vnet by ID.
func (c *Client) GetVnetWithParams(ctx context.Context, params *sdn.QueryParams) (*VnetData, error) {
	resBody := &vnetResponse{}

	err := c.DoRequest(ctx, http.MethodGet, c.ExpandPath(""), params, resBody)
	if err != nil {
		return nil, fmt.Errorf("error reading SDN Vnet %s: %w", c.ID, err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	return resBody.Data, nil
}

// GetVnets lists all SDN VNETs.
func (c *Client) GetVnets(ctx context.Context) ([]VnetData, error) {
	return c.GetVnetsWithParams(ctx, nil)
}

// GetVnets lists all SDN Vnets.
func (c *Client) GetVnetsWithParams(ctx context.Context, params *sdn.QueryParams) ([]VnetData, error) {
	resBody := &vnetsResponse{}

	err := c.DoRequest(ctx, http.MethodGet, c.basePath(), params, resBody)
	if err != nil {
		return nil, fmt.Errorf("error listing SDN VNETs: %w", err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	return *resBody.Data, nil
}

// CreateVnet creates a new SDN VNET.
func (c *Client) CreateVnet(ctx context.Context, data *Vnet) error {
	createRequest := VnetCreate{
		Vnet: *data,
		ID:   c.ID,
	}

	err := c.DoRequest(ctx, http.MethodPost, c.basePath(), createRequest, nil)
	if err != nil {
		return fmt.Errorf("error creating SDN VNET: %w", err)
	}

	return nil
}

// UpdateVnet Updates an existing VNet.
func (c *Client) UpdateVnet(ctx context.Context, data *VnetUpdate) error {
	err := c.DoRequest(ctx, http.MethodPut, c.ExpandPath(""), data, nil)
	if err != nil {
		return fmt.Errorf("error updating SDN VNET: %w", err)
	}

	return nil
}

// DeleteVnet deletes an SDN VNET by ID.
func (c *Client) DeleteVnet(ctx context.Context) error {
	err := c.DoRequest(ctx, http.MethodDelete, c.ExpandPath(""), nil, nil)
	if err != nil {
		return fmt.Errorf("error deleting SDN VNET: %w", err)
	}

	return nil
}
