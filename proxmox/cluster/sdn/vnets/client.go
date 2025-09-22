/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package vnets

import (
	"fmt"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/tasks"
)

// Client is a client for accessing the Proxmox SDN VNETs API.
type Client struct {
	api.Client

	ID string
}

func (c *Client) basePath() string {
	return c.Client.ExpandPath("sdn/vnets")
}

// ExpandPath expands a relative path to a full VM API path.
func (c *Client) ExpandPath(path string) string {
	p := fmt.Sprintf("%s/%s", c.basePath(), c.ID)
	if path != "" {
		p = fmt.Sprintf("%s/%s", p, path)
	}

	return p
}

// Tasks returns a client for managing VNET tasks.
func (c *Client) Tasks() *tasks.Client {
	return &tasks.Client{
		Client: c.Client,
	}
}
