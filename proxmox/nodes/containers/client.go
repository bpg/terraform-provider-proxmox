/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package containers

import (
	"fmt"

	"github.com/bpg/terraform-provider-proxmox/proxmox/firewall"
	containerfirewall "github.com/bpg/terraform-provider-proxmox/proxmox/nodes/containers/firewall"
	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

type Client struct {
	types.Client
	VMID int
}

func (c *Client) basePath() string {
	return c.Client.ExpandPath("lxc")
}

func (c *Client) ExpandPath(path string) string {
	ep := fmt.Sprintf("%s/%d", c.basePath(), c.VMID)
	if path != "" {
		ep = fmt.Sprintf("%s/%s", ep, path)
	}
	return ep
}

func (c *Client) Firewall() firewall.API {
	return &containerfirewall.Client{
		Client: firewall.Client{Client: c},
	}
}
