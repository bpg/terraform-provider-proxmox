/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package ceph

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	clusterceph "github.com/bpg/terraform-provider-proxmox/proxmox/cluster/ceph"
)

// GetStatus returns the node-scoped Ceph status. Returns an error when Ceph
// is not installed or not initialized on the node.
func (c *Client) GetStatus(ctx context.Context) (*clusterceph.StatusResponseData, error) {
	resBody := &StatusResponseBody{}

	if err := c.DoRequest(ctx, http.MethodGet, c.ExpandPath("status"), nil, resBody); err != nil {
		return nil, fmt.Errorf("error getting node Ceph status: %w", err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	return resBody.Data, nil
}
