/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package subnets

import (
	"fmt"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
)

// Client is a client for accessing the Proxmox SDN VNETs API.
type Client struct {
	api.Client
}

func (c *Client) basePath() string {
	return c.Client.ExpandPath("subnets")
}

// ExpandPath expands a relative path to a full VM API path.
func (c *Client) ExpandPath(subnetID string) string {
	p := c.basePath()
	if subnetID != "" {
		p = fmt.Sprintf("%s/%s", p, subnetID)
	}

	return p
}
