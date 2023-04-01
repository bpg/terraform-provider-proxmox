/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package vm

import (
	"fmt"
	"net/url"

	"github.com/bpg/terraform-provider-proxmox/proxmox/firewall"
	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
	vmfirewall "github.com/bpg/terraform-provider-proxmox/proxmox/vm/firewall"
)

type Client struct {
	types.Client
	NodeName string
	VMID     int
}

func (c *Client) AdjustPath(path string) string {
	return fmt.Sprintf("nodes/%s/qemu/%d/%s", url.PathEscape(c.NodeName), c.VMID, path)
}

func (c *Client) Firewall() firewall.API {
	return &vmfirewall.Client{
		Client: firewall.Client{Client: c},
	}
}
