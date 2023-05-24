/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package node

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

// GetDNS retrieves the DNS configuration for a node.
func (c *Client) GetDNS(ctx context.Context) (*DNSGetResponseData, error) {
	resBody := &DNSGetResponseBody{}

	err := c.DoRequest(ctx, http.MethodGet, c.ExpandPath("dns"), nil, resBody)
	if err != nil {
		return nil, fmt.Errorf("error retrieving DNS configuration: %w", err)
	}

	if resBody.Data == nil {
		return nil, types.ErrNoDataObjectInResponse
	}

	return resBody.Data, nil
}

// UpdateDNS updates the DNS configuration for a node.
func (c *Client) UpdateDNS(ctx context.Context, d *DNSUpdateRequestBody) error {
	err := c.DoRequest(ctx, http.MethodPut, c.ExpandPath("dns"), d, nil)
	if err != nil {
		return fmt.Errorf("error updating DNS configuration: %w", err)
	}

	return nil
}
