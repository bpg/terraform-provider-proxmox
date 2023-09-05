/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package cluster

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
)

// GetOptions retrieves the cluster options.
func (c *Client) GetOptions(ctx context.Context) (*OptionsResponseData, error) {
	resBody := &OptionsResponseBody{}

	err := c.DoRequest(ctx, http.MethodGet, c.ExpandPath("options"), nil, resBody)
	if err != nil {
		return nil, fmt.Errorf("error reading Cluster options: %w", err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	return resBody.Data, nil
}

// CreateUpdateOptions updates the cluster options.
func (c *Client) CreateUpdateOptions(ctx context.Context, data *OptionsRequestData) error {
	err := c.DoRequest(ctx, http.MethodPut, c.ExpandPath("options"), data, nil)
	if err != nil {
		return fmt.Errorf("error updating Cluster resource: %w", err)
	}

	return nil
}
