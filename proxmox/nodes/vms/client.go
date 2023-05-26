/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package vms

import (
	"fmt"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/firewall"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/tasks"
	vmfirewall "github.com/bpg/terraform-provider-proxmox/proxmox/nodes/vms/firewall"
)

// Client is an interface for accessing the Proxmox VM API.
type Client struct {
	api.Client
	VMID int
}

func (c *Client) basePath() string {
	return c.Client.ExpandPath("qemu")
}

// ExpandPath expands a relative path to a full VM API path.
func (c *Client) ExpandPath(path string) string {
	ep := fmt.Sprintf("%s/%d", c.basePath(), c.VMID)
	if path != "" {
		ep = fmt.Sprintf("%s/%s", ep, path)
	}

	return ep
}

// Tasks returns a client for managing VM tasks.
func (c *Client) Tasks() *tasks.Client {
	return &tasks.Client{
		Client: c.Client,
	}
}

// Firewall returns a client for managing the VM firewall.
func (c *Client) Firewall() firewall.API {
	return &vmfirewall.Client{
		Client: firewall.Client{Client: c},
	}
}
