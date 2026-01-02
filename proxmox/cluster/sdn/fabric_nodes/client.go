/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package fabric_nodes

import (
	"fmt"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
)

// Client is a client for accessing the Proxmox SDN Fabrics Node API.
type Client struct {
	api.Client

	NodeID         string
	FabricProtocol string
	FabricID       string
}

func (c *Client) basePath() string {
	return c.Client.ExpandPath("sdn/fabrics/node")
}

// ExpandPath expands a relative path to a full VM API path.
func (c *Client) ExpandPath(path string) string {
	p := fmt.Sprintf("%s/%s", c.basePath(), c.FabricID)
	if path != "" {
		p = fmt.Sprintf("%s/%s", p, path)
	}

	return p
}
