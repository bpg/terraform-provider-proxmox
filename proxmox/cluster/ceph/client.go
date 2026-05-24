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
)

// Client is an interface for accessing the Proxmox cluster-wide Ceph API.
type Client struct {
	api.Client
}

// ExpandPath expands a relative path to a full cluster Ceph API path.
func (c *Client) ExpandPath(path string) string {
	return fmt.Sprintf("cluster/ceph/%s", path)
}

// StatusResponseBody wraps the cluster Ceph status response. The full payload is
// large and version-dependent; this minimal struct exists so callers can probe
// whether Ceph is initialized without paying for full decode.
type StatusResponseBody struct {
	Data map[string]any `json:"data,omitempty"`
}

// GetStatus returns the cluster-wide Ceph status. Returns an error when Ceph is
// not installed or not initialized on the cluster, which makes it a cheap probe
// for "is Ceph usable?".
func (c *Client) GetStatus(ctx context.Context) (*StatusResponseBody, error) {
	resBody := &StatusResponseBody{}

	if err := c.DoRequest(ctx, http.MethodGet, c.ExpandPath("status"), nil, resBody); err != nil {
		return nil, fmt.Errorf("error getting Ceph cluster status: %w", err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	return resBody, nil
}
