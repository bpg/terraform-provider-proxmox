/*
	This Source Code Form is subject to the terms of the Mozilla Public
	License, v. 2.0. If a copy of the MPL was not distributed with this
	file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/

package mapping

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

// Create creates a new hardware mapping.
func (c *Client) Create(
	ctx context.Context,
	hmType types.HardwareMappingType,
	data *HardwareMappingCreateRequestBody,
) error {
	err := c.DoRequest(ctx, http.MethodPost, c.ExpandPath(hmType, ""), data, nil)
	if err != nil {
		return fmt.Errorf("creating hardware mapping %q: %w", data.ID, err)
	}

	return nil
}

// Delete deletes a hardware mapping.
func (c *Client) Delete(ctx context.Context, hmType types.HardwareMappingType, name string) error {
	err := c.DoRequest(ctx, http.MethodDelete, c.ExpandPath(hmType, url.PathEscape(name)), nil, nil)
	if err != nil {
		return fmt.Errorf("deleting hardware mapping %q: %w", name, err)
	}

	return nil
}

// Get retrieves the configuration of a single hardware mapping.
func (c *Client) Get(
	ctx context.Context,
	hmType types.HardwareMappingType,
	name string,
) (*HardwareMappingGetResponseData, error) {
	resBody := &HardwareMappingGetResponseBody{}

	err := c.DoRequest(ctx, http.MethodGet, c.ExpandPath(hmType, url.PathEscape(name)), nil, resBody)
	if err != nil {
		return nil, fmt.Errorf("reading hardware mapping %q: %w", name, err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	return resBody.Data, nil
}

// List retrieves the list of hardware mappings.
// If "checkNode" is not empty, the "checks" list will be included in the response that might include configuration
// correctness diagnostics for the given node.
func (c *Client) List(
	ctx context.Context,
	hmType types.HardwareMappingType,
	checkNode string,
) ([]*HardwareMappingListResponseData, error) {
	options := &hardwareMappingListQuery{
		CheckNode: checkNode,
	}

	resBody := &HardwareMappingListResponseBody{}

	err := c.DoRequest(ctx, http.MethodGet, c.ExpandPath(hmType, ""), options, resBody)
	if err != nil {
		return nil, fmt.Errorf("listing hardware mapping: %w", err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	return resBody.Data, nil
}

// Update updates an existing hardware mapping.
func (c *Client) Update(
	ctx context.Context,
	hmType types.HardwareMappingType,
	name string,
	data *HardwareMappingUpdateRequestBody,
) error {
	err := c.DoRequest(ctx, http.MethodPut, c.ExpandPath(hmType, url.PathEscape(name)), data, nil)
	if err != nil {
		return fmt.Errorf("udating hardware mapping %q: %w", name, err)
	}

	return nil
}
