/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package containers

import (
	"fmt"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/firewall"
	containerfirewall "github.com/bpg/terraform-provider-proxmox/proxmox/nodes/containers/firewall"
)

// Client is an interface for accessing the Proxmox container API.
type Client struct {
	api.Client
	VMID int
}

func (c *Client) basePath() string {
	return c.Client.ExpandPath("lxc")
}

// ExpandPath expands a relative path to a full container API path.
func (c *Client) ExpandPath(path string) string {
	ep := fmt.Sprintf("%s/%d", c.basePath(), c.VMID)
	if path != "" {
		ep = fmt.Sprintf("%s/%s", ep, path)
	}

	return ep
}

// Firewall returns a client for managing the container firewall.
func (c *Client) Firewall() firewall.API {
	return &containerfirewall.Client{
		Client: firewall.Client{Client: c},
	}
}
