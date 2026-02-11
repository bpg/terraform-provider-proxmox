/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package replications

import (
	"fmt"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
)

// Client is a client for accessing the Proxmox Replication API.
type Client struct {
	api.Client

	ID string
}

func (c *Client) basePath() string {
	return c.Client.ExpandPath("replication")
}

// ExpandPath expands a relative path to a full Replication path.
func (c *Client) ExpandPath() string {
	p := fmt.Sprintf("%s/%s", c.basePath(), c.ID)
	return p
}
