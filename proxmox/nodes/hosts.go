/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package nodes

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

// GetHosts retrieves the Hosts configuration for a node.
func (c *Client) GetHosts(ctx context.Context) (*HostsGetResponseData, error) {
	resBody := &HostsGetResponseBody{}
	err := c.DoRequest(ctx, http.MethodGet, c.ExpandPath("hosts"), nil, resBody)
	if err != nil {
		return nil, fmt.Errorf("error retrieving hosts configuration: %w", err)
	}

	if resBody.Data == nil {
		return nil, errors.New("the server did not include a data object in the response")
	}

	return resBody.Data, nil
}

// UpdateHosts updates the Hosts configuration for a node.
func (c *Client) UpdateHosts(ctx context.Context, d *HostsUpdateRequestBody) error {
	err := c.DoRequest(ctx, http.MethodPost, c.ExpandPath("hosts"), d, nil)
	if err != nil {
		return fmt.Errorf("error updating hosts configuration: %w", err)
	}
	return nil
}
