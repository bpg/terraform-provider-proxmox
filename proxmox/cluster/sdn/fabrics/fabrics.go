/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package fabrics

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/sdn"
	"github.com/bpg/terraform-provider-proxmox/proxmox/helpers/ptr"
)

const (
	ProtocolOpenFabric = "openfabric"
	ProtocolOSPF       = "ospf"
)

// GetFabric retrieves a single SDN Fabric by ID.
func (c *Client) GetFabric(ctx context.Context, id string) (*FabricData, error) {
	return c.GetFabricWithParams(ctx, id, nil)
}

// GetFabricWithParams retrieves a single SDN Fabric by ID with query parameters.
func (c *Client) GetFabricWithParams(ctx context.Context, id string, params *sdn.QueryParams) (*FabricData, error) {
	resBody := &struct {
		Data *FabricData `json:"data"`
	}{}

	err := c.DoRequest(ctx, http.MethodGet, c.ExpandPath(id), params, resBody)
	if err != nil {
		return nil, fmt.Errorf("error reading SDN Fabric %s: %w", id, err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	return resBody.Data, nil
}

// GetFabrics lists all SDN Fabrics.
func (c *Client) GetFabrics(ctx context.Context) ([]FabricData, error) {
	return c.GetFabricsWithParams(ctx, nil)
}

// GetFabricsWithParams lists all SDN Fabrics with query parameters.
func (c *Client) GetFabricsWithParams(ctx context.Context, params *sdn.QueryParams) ([]FabricData, error) {
	resBody := &struct {
		Data *[]FabricData `json:"data"`
	}{}

	err := c.DoRequest(ctx, http.MethodGet, c.basePath(), params, resBody)
	if err != nil {
		return nil, fmt.Errorf("error listing SDN Fabrics: %w", err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	return *resBody.Data, nil
}

// CreateFabric creates a new SDN Fabric.
func (c *Client) CreateFabric(ctx context.Context, data *Fabric) error {
	createRequest := FabricCreate{
		Fabric: *data,
	}

	err := c.DoRequest(ctx, http.MethodPost, c.basePath(), createRequest, nil)
	if err != nil {
		return fmt.Errorf("error creating SDN Fabric: %w", err)
	}

	return nil
}

// UpdateFabric Updates an existing Fabric.
func (c *Client) UpdateFabric(ctx context.Context, data *FabricUpdate) error {
	data.Protocol = ptr.Ptr(c.Protocol)
	err := c.DoRequest(ctx, http.MethodPut, c.ExpandPath(data.ID), data, nil)
	if err != nil {
		return fmt.Errorf("error updating SDN Fabric: %w", err)
	}

	return nil
}

// DeleteFabric deletes an SDN Fabric by ID.
func (c *Client) DeleteFabric(ctx context.Context, id string) error {
	err := c.DoRequest(ctx, http.MethodDelete, c.ExpandPath(id), nil, nil)
	if err != nil {
		return fmt.Errorf("error deleting SDN Fabric: %w", err)
	}

	return nil
}
