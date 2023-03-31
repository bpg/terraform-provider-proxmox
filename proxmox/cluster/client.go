/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package cluster

import (
	"fmt"

	clusterfirewall "github.com/bpg/terraform-provider-proxmox/proxmox/cluster/firewall"
	"github.com/bpg/terraform-provider-proxmox/proxmox/firewall"
	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

type Client struct {
	types.Client
}

func (c *Client) AddPrefix(path string) string {
	return fmt.Sprintf("cluster/%s", path)
}

func (c *Client) Firewall() clusterfirewall.API {
	return &clusterfirewall.Client{
		Client: firewall.Client{Client: c},
	}
}
