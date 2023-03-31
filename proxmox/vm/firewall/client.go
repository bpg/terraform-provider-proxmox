/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package firewall

import (
	"fmt"
	"net/url"

	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

type Client struct {
	types.Client
	NodeName string
	VMID     int
}

func (c *Client) Prefix() string {
	return fmt.Sprintf("nodes/%s/qemu/%d", url.PathEscape(c.NodeName), c.VMID)
}
