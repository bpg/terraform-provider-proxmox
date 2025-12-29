/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package zones

import (
	"fmt"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
)

// Client is a client for accessing the Proxmox SDN Zones API.
type Client struct {
	api.Client
}

func (c *Client) basePath() string {
	return c.Client.ExpandPath("sdn/zones")
}

// ExpandPath returns the API path for SDN zones.
func (c *Client) ExpandPath(path string) string {
	p := c.basePath()
	if path != "" {
		p = fmt.Sprintf("%s/%s", p, path)
	}

	return p
}
