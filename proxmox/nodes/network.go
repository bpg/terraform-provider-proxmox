/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package nodes

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"sort"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
)

// ListNetworkDevices retrieves a list of network devices for a specific nodes.
func (c *Client) ListNetworkDevices(ctx context.Context) ([]*NetworkDeviceListResponseData, error) {
	resBody := &NetworkDeviceListResponseBody{}

	err := c.DoRequest(ctx, http.MethodGet, c.ExpandPath("network"), nil, resBody)
	if err != nil {
		return nil, fmt.Errorf("failed to get network devices for node \"%s\": %w", c.NodeName, err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	sort.Slice(resBody.Data, func(i, j int) bool {
		return resBody.Data[i].Priority < resBody.Data[j].Priority
	})

	return resBody.Data, nil
}

// CreateNetworkDevice creates a network device for a specific node.
func (c *Client) CreateNetworkDevice(ctx context.Context, d *NetworkDeviceCreateUpdateRequestBody) error {
	err := c.DoRequest(ctx, http.MethodPost, c.ExpandPath("network"), d, nil)
	if err != nil {
		return fmt.Errorf("failed to create network device for node \"%s\": %w", c.NodeName, err)
	}

	return nil
}

// ReloadNetworkConfiguration reloads the network configuration for a specific node.
func (c *Client) ReloadNetworkConfiguration(ctx context.Context) error {
	err := c.DoRequest(ctx, http.MethodPut, c.ExpandPath("network"), nil, nil)
	if err != nil {
		return fmt.Errorf("failed to reload network configuration for node \"%s\": %w", c.NodeName, err)
	}

	return nil
}

// RevertNetworkConfiguration reverts the network configuration changes for a specific node.
func (c *Client) RevertNetworkConfiguration(ctx context.Context) error {
	err := c.DoRequest(ctx, http.MethodDelete, c.ExpandPath("network"), nil, nil)
	if err != nil {
		return fmt.Errorf("failed to revert network configuration for node \"%s\": %w", c.NodeName, err)
	}

	return nil
}

// UpdateNetworkDevice updates a network device for a specific node.
func (c *Client) UpdateNetworkDevice(ctx context.Context, iface string, d *NetworkDeviceCreateUpdateRequestBody) error {
	err := c.DoRequest(
		ctx,
		http.MethodPut,
		c.ExpandPath(fmt.Sprintf("network/%s", url.PathEscape(iface))),
		d,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to update network device for node \"%s\": %w", c.NodeName, err)
	}

	return nil
}

// DeleteNetworkDevice deletes a network device configuration for a specific node.
func (c *Client) DeleteNetworkDevice(ctx context.Context, iface string) error {
	err := c.DoRequest(
		ctx,
		http.MethodDelete,
		c.ExpandPath(fmt.Sprintf("network/%s", url.PathEscape(iface))),
		nil,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to delete network device for node \"%s\": %w", c.NodeName, err)
	}

	return nil
}
