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
	proxmoxtypes "github.com/bpg/terraform-provider-proxmox/proxmox/types/hardwaremapping"
)

// encodedDescription returns a percent-encoded copy of a description so it survives the round-trip through
// the Proxmox VE API, which double-encodes stored non-ASCII bytes on read.
func encodedDescription(description *string) *string {
	if description == nil {
		return nil
	}

	encoded := proxmoxtypes.EncodeText(*description)

	return &encoded
}

// decodeDescription reverses encodedDescription.
func decodeDescription(description *string) *string {
	if description == nil {
		return nil
	}

	decoded := proxmoxtypes.DecodeText(*description)

	return &decoded
}

// Create creates a new hardware mapping.
func (c *Client) Create(ctx context.Context, hmType proxmoxtypes.Type, data *CreateRequestBody) error {
	body := *data
	body.Description = encodedDescription(data.Description)

	err := c.DoRequest(ctx, http.MethodPost, c.ExpandPath(hmType, ""), &body, nil)
	if err != nil {
		return fmt.Errorf("creating hardware mapping %q: %w", data.ID, err)
	}

	return nil
}

// Delete deletes a hardware mapping.
func (c *Client) Delete(ctx context.Context, hmType proxmoxtypes.Type, name string) error {
	err := c.DoRequest(ctx, http.MethodDelete, c.ExpandPath(hmType, url.PathEscape(name)), nil, nil)
	if err != nil {
		return fmt.Errorf("deleting hardware mapping %q: %w", name, err)
	}

	return nil
}

// Get retrieves the configuration of a single hardware mapping.
func (c *Client) Get(ctx context.Context, hmType proxmoxtypes.Type, name string) (*GetResponseData, error) {
	resBody := &GetResponseBody{}

	err := c.DoRequest(ctx, http.MethodGet, c.ExpandPath(hmType, url.PathEscape(name)), nil, resBody)
	if err != nil {
		return nil, fmt.Errorf("reading hardware mapping %q: %w", name, err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	resBody.Data.Description = decodeDescription(resBody.Data.Description)

	return resBody.Data, nil
}

// List retrieves the list of hardware mappings.
// If "checkNode" is not empty, the "checks" list will be included in the response that might include configuration
// correctness diagnostics for the given node.
func (c *Client) List(ctx context.Context, hmType proxmoxtypes.Type, checkNode string) ([]*ListResponseData, error) {
	options := &listQuery{
		CheckNode: checkNode,
	}

	resBody := &ListResponseBody{}

	err := c.DoRequest(ctx, http.MethodGet, c.ExpandPath(hmType, ""), options, resBody)
	if err != nil {
		return nil, fmt.Errorf("listing hardware mapping: %w", err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	for _, item := range resBody.Data {
		item.Description = decodeDescription(item.Description)
	}

	return resBody.Data, nil
}

// Update updates an existing hardware mapping.
func (c *Client) Update(ctx context.Context, hmType proxmoxtypes.Type, name string, data *UpdateRequestBody) error {
	body := *data
	body.Description = encodedDescription(data.Description)

	err := c.DoRequest(ctx, http.MethodPut, c.ExpandPath(hmType, url.PathEscape(name)), &body, nil)
	if err != nil {
		return fmt.Errorf("udating hardware mapping %q: %w", name, err)
	}

	return nil
}
