/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package nodes

import (
	"fmt"
	"net/url"

	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/containers"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/vms"
	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

// Client is an interface for accessing the Proxmox node API.
type Client struct {
	types.Client
	NodeName string
}

// ExpandPath expands a relative path to a full node API path.
func (c *Client) ExpandPath(path string) string {
	return fmt.Sprintf("nodes/%s/%s", url.PathEscape(c.NodeName), path)
}

// Container returns a client for managing a specific container.
func (c *Client) Container(vmID int) *containers.Client {
	return &containers.Client{
		Client: c,
		VMID:   vmID,
	}
}

// VM returns a client for managing a specific VM.
func (c *Client) VM(vmID int) *vms.Client {
	return &vms.Client{
		Client: c,
		VMID:   vmID,
	}
}
