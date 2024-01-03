/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package storage

import (
	"fmt"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/tasks"
)

// Client is an interface for accessing the Proxmox node storage API.
type Client struct {
	api.Client
	StorageName string
}

func (c *Client) basePath() string {
	return c.Client.ExpandPath("storage")
}

// ExpandPath expands a relative path to a full node storage API path.
func (c *Client) ExpandPath(path string) string {
	ep := fmt.Sprintf("%s/%s", c.basePath(), c.StorageName)
	if path != "" {
		ep = fmt.Sprintf("%s/%s", ep, path)
	}

	return ep
}

// Tasks returns a client for managing node storage tasks.
func (c *Client) Tasks() *tasks.Client {
	return &tasks.Client{
		Client: c.Client,
	}
}
