/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package fabric_nodes

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/sdn"
	"github.com/bpg/terraform-provider-proxmox/proxmox/helpers/ptr"
)

// GetFabricNode retrieves a single SDN Fabric Node by Fabric Node ID.
func (c *Client) GetFabricNode(ctx context.Context, nodeID string) (*FabricNodeData, error) {
	return c.GetFabricNodeWithParams(ctx, nodeID, nil)
}

// GetFabricNodeWithParams retrieves a single SDN Fabric Node by Fabric Node ID with query parameters.
func (c *Client) GetFabricNodeWithParams(ctx context.Context, nodeID string, params *sdn.QueryParams) (*FabricNodeData, error) {
	resBody := &struct {
		Data *FabricNodeData `json:"data"`
	}{}

	err := c.DoRequest(ctx, http.MethodGet, c.ExpandPath(nodeID), params, resBody)
	if err != nil {
		return nil, fmt.Errorf("error reading SDN Fabric Node %s: %w", c.NodeID, err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	return resBody.Data, nil
}

// GetFabricNodes lists all SDN Fabric Nodes.
func (c *Client) GetFabricNodes(ctx context.Context) ([]FabricNodeData, error) {
	return c.GetFabricNodesWithParams(ctx, nil)
}

// GetFabricNodesWithParams lists all SDN Fabric Nodes with query parameters.
func (c *Client) GetFabricNodesWithParams(ctx context.Context, params *sdn.QueryParams) ([]FabricNodeData, error) {
	resBody := &struct {
		Data *[]FabricNodeData `json:"data"`
	}{}

	err := c.DoRequest(ctx, http.MethodGet, c.basePath(), params, resBody)
	if err != nil {
		return nil, fmt.Errorf("error listing SDN Fabric Nodes: %w", err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	return *resBody.Data, nil
}

// CreateFabricNode creates a new SDN Fabric Node.
func (c *Client) CreateFabricNode(ctx context.Context, data *FabricNode) error {
	createRequest := FabricNodeCreate{
		FabricNode: *data,
	}

	err := c.DoRequest(ctx, http.MethodPost, c.ExpandPath(""), createRequest, nil)
	if err != nil {
		return fmt.Errorf("error creating SDN Fabric Node: %w", err)
	}

	return nil
}

// UpdateFabricNode Updates an existing Fabric Node.
func (c *Client) UpdateFabricNode(ctx context.Context, data *FabricNodeUpdate) error {
	data.Protocol = ptr.Ptr(c.FabricProtocol)
	data.FabricID = c.FabricID
	err := c.DoRequest(ctx, http.MethodPut, c.ExpandPath(data.NodeID), data, nil)
	if err != nil {
		return fmt.Errorf("error updating SDN Fabric Node: %w", err)
	}

	return nil
}

// DeleteFabricNode deletes an SDN Fabric Node by Node ID.
func (c *Client) DeleteFabricNode(ctx context.Context, nodeID string) error {
	err := c.DoRequest(ctx, http.MethodDelete, c.ExpandPath(nodeID), nil, nil)
	if err != nil {
		return fmt.Errorf("error deleting SDN Fabric Node: %w", err)
	}

	return nil
}
