/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package mds

import (
	"fmt"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/tasks"
)

// Client is an interface for accessing the Proxmox Ceph metadata server (MDS) API.
type Client struct {
	api.Client
}

// ExpandPath expands a relative path to a full Ceph MDS API path.
func (c *Client) ExpandPath(path string) string {
	return c.Client.ExpandPath(fmt.Sprintf("mds/%s", path))
}

// Tasks returns a client for managing Ceph MDS tasks.
func (c *Client) Tasks() *tasks.Client {
	return &tasks.Client{Client: c.Client}
}
