/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package account

import (
	"fmt"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/tasks"
)

// Client is an interface for accessing the Proxmox ACME management API.
type Client struct {
	api.Client
}

// ExpandPath expands a relative path to the Proxmox ACME management API path.
func (c *Client) ExpandPath(path string) string {
	return fmt.Sprintf("cluster/acme/account/%s", path)
}

// Tasks returns a client for managing ACME account tasks.
func (c *Client) Tasks() *tasks.Client {
	return &tasks.Client{
		Client: c.Client,
	}
}
