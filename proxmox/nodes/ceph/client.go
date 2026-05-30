/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package ceph

import (
	"fmt"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/ceph/pool"
)

// Client is an interface for accessing the Proxmox node-scoped Ceph API.
type Client struct {
	api.Client
}

// ExpandPath expands a relative path to a full node Ceph API path.
func (c *Client) ExpandPath(path string) string {
	return c.Client.ExpandPath(fmt.Sprintf("ceph/%s", path))
}

// Pool returns a client for managing Ceph pools.
func (c *Client) Pool() *pool.Client {
	return &pool.Client{Client: c}
}
