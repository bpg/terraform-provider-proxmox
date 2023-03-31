/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package vm

import (
	fw "github.com/bpg/terraform-provider-proxmox/proxmox/firewall"
	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
	"github.com/bpg/terraform-provider-proxmox/proxmox/vm/firewall"
)

type Client struct {
	types.Client
	NodeName string
	VMID     int
}

func (c *Client) Firewall() fw.API {
	return &firewall.Client{Client: c, NodeName: c.NodeName, VMID: c.VMID}
}
