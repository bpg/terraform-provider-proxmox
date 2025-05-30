package subnets

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
)

// GetSubnet retrieves a single Subnet by ID and containing Vnet's ID
func (c *Client) GetSubnet(ctx context.Context, vnetID string, id string) (*SubnetData, error) {
	resBody := &SubnetResponseBody{}

	err := c.DoRequest(ctx, http.MethodGet, c.ExpandPath(vnetID, id), nil, resBody)
	if err != nil {
		return nil, fmt.Errorf("Error reading SDN subnet %s for Vnet %s: %w", id, vnetID, err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	return resBody.Data, nil
}

// GetSubnets lists all Subnets related to a Vnet
func (c *Client) GetSubnets(ctx context.Context, vnetID string) ([]SubnetData, error) {
	resBody := &SubnetsResponseBody{}

	err := c.DoRequest(ctx, http.MethodGet, c.ExpandPath(vnetID, ""), nil, resBody)
	if err != nil {
		return nil, fmt.Errorf("Error listing Subnets for Vnet %s: %w", vnetID, err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	return *resBody.Data, nil
}

// CreateSubnet creates a new Subnet in the defined Vnet
func (c *Client) CreateSubnet(ctx context.Context, vnetID string, data *SubnetRequestData) error {
	err := c.DoRequest(ctx, http.MethodPost, c.ExpandPath(vnetID, ""), data, nil)
	if err != nil {
		return fmt.Errorf("Error creating subnet %s on VNet %s: %w", data.ID, vnetID, err)
	}

	return nil
}

// UpdateSubnet updates an existing subnet inside a defined vnet
func (c *Client) UpdateSubnet(ctx context.Context, vnetID string, data *SubnetRequestData) error {
	err := c.DoRequest(ctx, http.MethodPut, c.ExpandPath(vnetID, data.ID), data, nil)
	if err != nil {
		return fmt.Errorf("Error updating subnet %s on VNet %s: %w", data.ID, vnetID, err)
	}

	return nil
}

// DeleteSubnet deletes an existing subnet inside a defined vnet
func (c *Client) DeleteSubnet(ctx context.Context, vnetID string, id string) error {
	err := c.DoRequest(ctx, http.MethodDelete, c.ExpandPath(vnetID, id), nil, nil)
	if err != nil {
		return fmt.Errorf("Error deleting subnet %s on VNet %s: %s", id, vnetID, err)
	}

	return nil
}
