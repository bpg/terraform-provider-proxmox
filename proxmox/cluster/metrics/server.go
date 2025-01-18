/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package metrics

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
)

// GetServer retrieves the metrics server data.
func (c *Client) GetServer(ctx context.Context, id string) (*ServerData, error) {
	resBody := &ServerResponseBody{}

	err := c.DoRequest(ctx, http.MethodGet, c.ExpandPath(id), nil, resBody)
	if err != nil {
		return nil, fmt.Errorf("error reading Cluster options: %w", err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	return resBody.Data, nil
}

// GetServers lists the metrics servers.
func (c *Client) GetServers(ctx context.Context) (*[]ServerData, error) {
	resBody := &ServersResponseBody{}

	err := c.DoRequest(ctx, http.MethodGet, c.ExpandPath(""), nil, resBody)
	if err != nil {
		return nil, fmt.Errorf("error reading Cluster options: %w", err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	return resBody.Data, nil
}

// UpdateServer updates the metrics server.
func (c *Client) UpdateServer(ctx context.Context, data *ServerData) error {
	if data.ID == nil {
		return errors.New("ID must be provided in data")
	}

	err := c.DoRequest(ctx, http.MethodPut, c.ExpandPath(*data.ID), data, nil)
	if err != nil {
		return fmt.Errorf("error updating metrics Server resource: %w", err)
	}

	return nil
}

// CreateServer creates the metrics server.
func (c *Client) CreateServer(ctx context.Context, data *ServerData) error {
	if data.ID == nil {
		return errors.New("ID must be provided in data")
	}

	err := c.DoRequest(ctx, http.MethodPost, c.ExpandPath(*data.ID), data, nil)
	if err != nil {
		return fmt.Errorf("error creating metrics Server resource: %w", err)
	}

	return nil
}

// DeleteServer deletes the metrics server.
func (c *Client) DeleteServer(ctx context.Context, id string) error {
	err := c.DoRequest(ctx, http.MethodDelete, c.ExpandPath(id), nil, nil)
	if err != nil {
		return fmt.Errorf("error updating Cluster resource: %w", err)
	}

	return nil
}
