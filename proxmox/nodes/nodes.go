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
	"sort"
	"strings"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
)

// GetIP retrieves the IP address of a node.
func (c *Client) GetIP(ctx context.Context) (string, error) {
	networkDevices, err := c.ListNetworkDevices(ctx)
	if err != nil {
		return "", err
	}

	nodeAddress := ""

	for _, d := range networkDevices {
		if d.Address != nil {
			nodeAddress = *d.Address
			break
		}
	}

	if nodeAddress == "" {
		return "", fmt.Errorf("failed to determine the IP address of node \"%s\"", c.NodeName)
	}

	nodeAddressParts := strings.Split(nodeAddress, "/")

	return nodeAddressParts[0], nil
}

// GetTime retrieves the time information for a node.
func (c *Client) GetTime(ctx context.Context) (*GetTimeResponseData, error) {
	resBody := &GetTimeResponseBody{}

	err := c.DoRequest(ctx, http.MethodGet, c.ExpandPath("time"), nil, resBody)
	if err != nil {
		return nil, fmt.Errorf("failed to get time information for node \"%s\": %w", c.NodeName, err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	return resBody.Data, nil
}

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

// ListNodes retrieves a list of nodes.
func (c *Client) ListNodes(ctx context.Context) ([]*ListResponseData, error) {
	resBody := &ListResponseBody{}

	err := c.DoRequest(ctx, http.MethodGet, "nodes", nil, resBody)
	if err != nil {
		return nil, fmt.Errorf("failed to get nodes: %w", err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	sort.Slice(resBody.Data, func(i, j int) bool {
		return resBody.Data[i].Name < resBody.Data[j].Name
	})

	return resBody.Data, nil
}

// UpdateTime updates the time on a node.
func (c *Client) UpdateTime(ctx context.Context, d *UpdateTimeRequestBody) error {
	err := c.DoRequest(ctx, http.MethodPut, c.ExpandPath("time"), d, nil)
	if err != nil {
		return fmt.Errorf("failed to update node time: %w", err)
	}

	return nil
}
