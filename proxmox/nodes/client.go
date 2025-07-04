/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package nodes

import (
	"fmt"
	"net/url"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/apt"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/containers"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/storage"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/tasks"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/vms"
)

// Client is an interface for accessing the Proxmox node API.
type Client struct {
	api.Client

	NodeName string
}

// ExpandPath expands a relative path to a full node API path.
func (c *Client) ExpandPath(path string) string {
	return fmt.Sprintf("nodes/%s/%s", url.PathEscape(c.NodeName), path)
}

// APT returns a client for managing APT related settings.
func (c *Client) APT() *apt.Client {
	return &apt.Client{
		Client: c,
	}
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

// Storage returns a client for managing a specific storage.
func (c *Client) Storage(storageName string) *storage.Client {
	return &storage.Client{
		Client:      c,
		StorageName: storageName,
	}
}

// Tasks returns a client for managing VM tasks.
func (c *Client) Tasks() *tasks.Client {
	return &tasks.Client{
		Client: c,
	}
}
