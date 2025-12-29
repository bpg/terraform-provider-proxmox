/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package applier

import (
	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/tasks"
)

// Client is a client for accessing the Proxmox SDN Apply API.
type Client struct {
	api.Client
}

func (c *Client) basePath() string {
	return c.Client.ExpandPath("sdn")
}

// ExpandPath returns the API path for cluster-wide SDN apply.
func (c *Client) ExpandPath(_ string) string {
	return c.basePath()
}

// Tasks returns a client for managing SDN tasks.
func (c *Client) Tasks() *tasks.Client {
	return &tasks.Client{
		Client: c.Client,
	}
}
