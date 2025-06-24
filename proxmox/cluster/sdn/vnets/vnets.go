package vnets

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/sdn/zones"
)

// GetVnet retrieves a single SDN Vnet by ID.
func (c *Client) GetVnet(ctx context.Context, id string) (*VnetData, error) {
	resBody := &VnetResponseBody{}

	err := c.DoRequest(ctx, http.MethodGet, c.ExpandPath(id), nil, resBody)
	if err != nil {
		return nil, fmt.Errorf("error reading SDN Vnet %s: %w", id, err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	return resBody.Data, nil
}

// GetVnets lists all SDN Vnets.
func (c *Client) GetVnets(ctx context.Context) ([]VnetData, error) {
	resBody := &VnetsResponseBody{}

	err := c.DoRequest(ctx, http.MethodGet, c.ExpandPath(""), nil, resBody)
	if err != nil {
		return nil, fmt.Errorf("error listing SDN Vnets: %w", err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	return *resBody.Data, nil
}

// CreateVnet creates a new SDN VNET.
func (c *Client) CreateVnet(ctx context.Context, data *VnetRequestData) error {
	err := c.DoRequest(ctx, http.MethodPost, c.ExpandPath(""), data, nil)
	if err != nil {
		return fmt.Errorf("error creating SDN VNET: %w", err)
	}

	return nil
}

// UpdateVnet Updates an existing VNet.
func (c *Client) UpdateVnet(ctx context.Context, data *VnetRequestData) error {
	err := c.DoRequest(ctx, http.MethodPut, c.ExpandPath(data.ID), data, nil)
	if err != nil {
		return fmt.Errorf("error updating SDN VNET: %w", err)
	}

	return nil
}

// DeleteVnet deletes an SDN VNET by ID.
func (c *Client) DeleteVnet(ctx context.Context, id string) error {
	err := c.DoRequest(ctx, http.MethodDelete, c.ExpandPath(id), nil, nil)
	if err != nil {
		return fmt.Errorf("error deleting SDN VNET: %w", err)
	}

	return nil
}

func (c *Client) GetParentZone(ctx context.Context, zoneId string) (*zones.ZoneData, error) {
	parentZone := zones.ZoneResponseBody{}

	err := c.DoRequest(ctx, http.MethodGet, c.ParentPath(zoneId), nil, parentZone)
	if err != nil {
		return nil, fmt.Errorf("error fetching vnet's parent zone %s: %w", zoneId, err)
	}

	return parentZone.Data, nil
}
