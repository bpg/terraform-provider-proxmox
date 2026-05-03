/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package controllers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/sdn"
)

// GetController retrieves a single SDN controller by ID.
func (c *Client) GetController(ctx context.Context, id string) (*ControllerData, error) {
	return c.GetControllerWithParams(ctx, id, nil)
}

// GetControllerWithParams retrieves a single SDN controller by ID with query parameters.
func (c *Client) GetControllerWithParams(ctx context.Context, id string, params *sdn.QueryParams) (*ControllerData, error) {
	resBody := &struct {
		Data *ControllerData `json:"data"`
	}{}

	err := c.DoRequest(ctx, http.MethodGet, c.ExpandPath(id), params, resBody)
	if err != nil {
		return nil, fmt.Errorf("error reading SDN controller %s: %w", id, err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	return resBody.Data, nil
}

// GetControllers lists all SDN controllers.
func (c *Client) GetControllers(ctx context.Context) ([]ControllerData, error) {
	return c.GetControllersWithParams(ctx, nil)
}

// GetControllersWithParams lists all SDN controllers with query parameters.
func (c *Client) GetControllersWithParams(ctx context.Context, params *sdn.QueryParams) ([]ControllerData, error) {
	resBody := &struct {
		Data *[]ControllerData `json:"data"`
	}{}

	err := c.DoRequest(ctx, http.MethodGet, c.ExpandPath(""), params, resBody)
	if err != nil {
		return nil, fmt.Errorf("error listing SDN controllers: %w", err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	return *resBody.Data, nil
}

// CreateController creates a new SDN controller.
func (c *Client) CreateController(ctx context.Context, controller *Controller) error {
	err := c.DoRequest(ctx, http.MethodPost, c.ExpandPath(""), controller, nil)
	if err != nil {
		return fmt.Errorf("error creating SDN controller: %w", err)
	}

	return nil
}

// UpdateController updates an existing SDN controller.
func (c *Client) UpdateController(ctx context.Context, update *ControllerUpdate) error {
	/* PVE API does not allow to pass "type" in PUT requests, this doesn't makes any sense
	since other required params like port, server must still be there
	while we could spawn another struct, let's just fix it silently */
	update.Type = nil

	err := c.DoRequest(ctx, http.MethodPut, c.ExpandPath(update.ID), update, nil)
	if err != nil {
		return fmt.Errorf("error updating SDN controller: %w", err)
	}

	return nil
}

// DeleteController deletes an SDN controller by ID.
func (c *Client) DeleteController(ctx context.Context, id string) error {
	err := c.DoRequest(ctx, http.MethodDelete, c.ExpandPath(id), nil, nil)
	if err != nil {
		return fmt.Errorf("error deleting SDN controller: %w", err)
	}

	return nil
}
