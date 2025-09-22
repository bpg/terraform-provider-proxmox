/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package subnets

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/sdn"
	"github.com/bpg/terraform-provider-proxmox/proxmox/helpers/ptr"
)

// GetSubnet retrieves a single Subnet by ID and containing Vnet's ID.
func (c *Client) GetSubnet(ctx context.Context, subnetID string) (*SubnetData, error) {
	return c.GetSubnetWithParams(ctx, subnetID, nil)
}

// GetSubnetWithParams retrieves a single Subnet by ID and containing Vnet's ID with query parameters.
func (c *Client) GetSubnetWithParams(ctx context.Context, subnetID string, params *sdn.QueryParams) (*SubnetData, error) {
	resBody := &subnetResponse{}

	err := c.DoRequest(ctx, http.MethodGet, c.ExpandPath(subnetID), params, resBody)
	if err != nil {
		return nil, fmt.Errorf("error reading SDN subnet %s: %w", subnetID, err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	return resBody.Data, nil
}

// GetSubnets lists all SDN Subnets.
func (c *Client) GetSubnets(ctx context.Context) ([]SubnetData, error) {
	return c.GetSubnetsWithParams(ctx, nil)
}

// GetSubnetsWithParams lists all SDN Subnets with query parameters.
func (c *Client) GetSubnetsWithParams(ctx context.Context, params *sdn.QueryParams) ([]SubnetData, error) {
	resBody := &subnetsResponse{}

	err := c.DoRequest(ctx, http.MethodGet, c.ExpandPath(""), params, resBody)
	if err != nil {
		return nil, fmt.Errorf("error listing SDN Subnets: %w", err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	return *resBody.Data, nil
}

// CreateSubnet creates a new Subnet in the defined Vnet.
func (c *Client) CreateSubnet(ctx context.Context, subnet *Subnet) error {
	subnet.Type = ptr.Ptr("subnet")

	err := c.DoRequest(ctx, http.MethodPost, c.ExpandPath(""), subnet, nil)
	if err != nil {
		return fmt.Errorf("error creating subnet %s: %w", subnet.ID, err)
	}

	return nil
}

// UpdateSubnet updates an existing subnet inside a defined vnet.
func (c *Client) UpdateSubnet(ctx context.Context, udate *SubnetUpdate) error {
	err := c.DoRequest(ctx, http.MethodPut, c.ExpandPath(udate.ID), udate, nil)
	if err != nil {
		return fmt.Errorf("error updating subnet %s: %w", udate.ID, err)
	}

	return nil
}

// DeleteSubnet deletes an existing subnet inside a defined vnet.
func (c *Client) DeleteSubnet(ctx context.Context, id string) error {
	err := c.DoRequest(ctx, http.MethodDelete, c.ExpandPath(id), nil, nil)
	if err != nil {
		return fmt.Errorf("error deleting subnet %s: %w", id, err)
	}

	return nil
}
