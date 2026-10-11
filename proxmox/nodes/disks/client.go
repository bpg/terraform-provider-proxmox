/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package disks

import (
	"fmt"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	zfsapi "github.com/bpg/terraform-provider-proxmox/proxmox/nodes/disks/zfs"
)

// Client provides access to the Proxmox node disks API (/nodes/{node}/disks).
type Client struct {
	api.Client
}

// ExpandPath expands a relative path to the node disks API base.
func (c *Client) ExpandPath(path string) string {
	return c.Client.ExpandPath(fmt.Sprintf("disks/%s", path))
}

// ZFS returns a client for managing ZFS pools on this node.
func (c *Client) ZFS() *zfsapi.Client {
	return &zfsapi.Client{Client: c}
}
