/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package container

import (
	"fmt"
	"net/url"

	containerfirewall "github.com/bpg/terraform-provider-proxmox/proxmox/container/firewall"
	"github.com/bpg/terraform-provider-proxmox/proxmox/firewall"
	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

type Client struct {
	types.Client
	NodeName string
	VMID     int
}

func (c *Client) ExpandPath(path string) string {
	return fmt.Sprintf("nodes/%s/lxc/%d/%s", url.PathEscape(c.NodeName), c.VMID, path)
}

func (c *Client) Firewall() firewall.API {
	return &containerfirewall.Client{
		Client: firewall.Client{Client: c},
	}
}
