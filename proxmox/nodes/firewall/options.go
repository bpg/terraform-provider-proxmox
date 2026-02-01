/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package firewall

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
)

// Options is an interface for managing node firewall options.
type Options interface {
	SetNodeOptions(ctx context.Context, d *OptionsPutRequestBody) error
	GetNodeOptions(ctx context.Context) (*OptionsGetResponseData, error)
}

// SetNodeOptions sets the node firewall options.
func (c *Client) SetNodeOptions(ctx context.Context, d *OptionsPutRequestBody) error {
	err := c.DoRequest(ctx, http.MethodPut, c.ExpandPath("firewall/options"), d, nil)
	if err != nil {
		return fmt.Errorf("error setting node firewall options: %w", err)
	}

	return nil
}

// GetNodeOptions retrieves the node firewall options.
func (c *Client) GetNodeOptions(ctx context.Context) (*OptionsGetResponseData, error) {
	resBody := &OptionsGetResponseBody{}

	err := c.DoRequest(ctx, http.MethodGet, c.ExpandPath("firewall/options"), nil, resBody)
	if err != nil {
		return nil, fmt.Errorf("error retrieving node firewall options: %w", err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	return resBody.Data, nil
}
