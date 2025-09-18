/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package zones

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/sdn"
)

// GetZone retrieves a single SDN zone by ID.
func (c *Client) GetZone(ctx context.Context, id string) (*ZoneData, error) {
	return c.GetZoneWithParams(ctx, id, nil)
}

// GetZoneWithParams retrieves a single SDN zone by ID with query parameters.
func (c *Client) GetZoneWithParams(ctx context.Context, id string, params *sdn.QueryParams) (*ZoneData, error) {
	resBody := &zoneResponseBody{}

	err := c.DoRequest(ctx, http.MethodGet, c.ExpandPath(id), params, resBody)
	if err != nil {
		return nil, fmt.Errorf("error reading SDN zone %s: %w", id, err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	return resBody.Data, nil
}

// GetZones lists all SDN zones.
func (c *Client) GetZones(ctx context.Context) ([]ZoneData, error) {
	return c.GetZonesWithParams(ctx, nil)
}

// GetZonesWithParams lists all SDN zones with query parameters.
func (c *Client) GetZonesWithParams(ctx context.Context, params *sdn.QueryParams) ([]ZoneData, error) {
	resBody := &zonesResponseBody{}

	err := c.DoRequest(ctx, http.MethodGet, c.ExpandPath(""), params, resBody)
	if err != nil {
		return nil, fmt.Errorf("error listing SDN zones: %w", err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	return *resBody.Data, nil
}

// CreateZone creates a new SDN zone.
func (c *Client) CreateZone(ctx context.Context, zone *Zone) error {
	err := c.DoRequest(ctx, http.MethodPost, c.ExpandPath(""), zone, nil)
	if err != nil {
		return fmt.Errorf("error creating SDN zone: %w", err)
	}

	return nil
}

// UpdateZone updates an existing SDN zone.
func (c *Client) UpdateZone(ctx context.Context, update *ZoneUpdate) error {
	/* PVE API does not allow to pass "type" in PUT requests, this doesn't makes any sense
	since other required params like port, server must still be there
	while we could spawn another struct, let's just fix it silently */
	update.Type = nil

	err := c.DoRequest(ctx, http.MethodPut, c.ExpandPath(update.ID), update, nil)
	if err != nil {
		return fmt.Errorf("error updating SDN zone: %w", err)
	}

	return nil
}

// DeleteZone deletes an SDN zone by ID.
func (c *Client) DeleteZone(ctx context.Context, id string) error {
	err := c.DoRequest(ctx, http.MethodDelete, c.ExpandPath(id), nil, nil)
	if err != nil {
		return fmt.Errorf("error deleting SDN zone: %w", err)
	}

	return nil
}
