/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package fabrics

import (
	"fmt"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
)

// Client is a client for accessing the Proxmox SDN Fabrics API.
type Client struct {
	api.Client

	ID       string
	Protocol string
}

func (c *Client) basePath() string {
	return c.Client.ExpandPath("sdn/fabrics/fabric")
}

// ExpandPath expands a relative path to a full VM API path.
func (c *Client) ExpandPath(path string) string {
	p := fmt.Sprintf("%s/%s", c.basePath(), c.ID)
	if path != "" {
		p = fmt.Sprintf("%s/%s", p, path)
	}

	return p
}
