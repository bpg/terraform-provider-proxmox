package zones

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
)

// GetZone retrieves a single SDN zone by ID.
func (c *Client) GetZone(ctx context.Context, id string) (*ZoneData, error) {
	resBody := &ZoneResponseBody{}

	err := c.DoRequest(ctx, http.MethodGet, c.ExpandPath(id), nil, resBody)
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
	resBody := &ZonesResponseBody{}

	err := c.DoRequest(ctx, http.MethodGet, c.ExpandPath(""), nil, resBody)
	if err != nil {
		return nil, fmt.Errorf("error listing SDN zones: %w", err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	return *resBody.Data, nil
}

// CreateZone creates a new SDN zone.
func (c *Client) CreateZone(ctx context.Context, data *ZoneRequestData) error {
	err := c.DoRequest(ctx, http.MethodPost, c.ExpandPath(""), data, nil)
	if err != nil {
		return fmt.Errorf("error creating SDN zone: %w", err)
	}

	return nil
}

// UpdateZone updates an existing SDN zone.
func (c *Client) UpdateZone(ctx context.Context, data *ZoneRequestData) error {
	// PVE API does not allow to pass "type" in PUT requests, this doesn't makes any sense
	// since other required params like port, server must still be there
	// while we could spawn another struct, let's just fix it silently
	data.Type = nil
	err := c.DoRequest(ctx, http.MethodPut, c.ExpandPath(data.ID), data, nil)
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
