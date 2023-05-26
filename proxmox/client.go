/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package proxmox

import (
	"github.com/bpg/terraform-provider-proxmox/proxmox/access"
	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes"
	"github.com/bpg/terraform-provider-proxmox/proxmox/pools"
	"github.com/bpg/terraform-provider-proxmox/proxmox/ssh"
	"github.com/bpg/terraform-provider-proxmox/proxmox/storage"
	"github.com/bpg/terraform-provider-proxmox/proxmox/version"
)

// Client defines a client interface for the Proxmox Virtual Environment API.
type Client interface {
	// Access returns a client for managing access control.
	Access() *access.Client

	// Cluster returns a client for managing the cluster.
	Cluster() *cluster.Client

	// Node returns a client for managing resources on a specific node.
	Node(nodeName string) *nodes.Client

	// Pool returns a client for managing resource pools.
	Pool() *pools.Client

	// Storage returns a client for managing storage.
	Storage() *storage.Client

	// Version returns a client for getting the version of the Proxmox Virtual Environment API.
	Version() *version.Client

	// API returns a lower-lever REST API client.
	API() api.Client

	// SSH returns a lower-lever SSH client.
	SSH() ssh.Client
}

type client struct {
	a api.Client
	s ssh.Client
}

// NewClient creates a new API client.
func NewClient(a api.Client, s ssh.Client) Client {
	return &client{a: a, s: s}
}

// Access returns a client for managing access control.
func (c *client) Access() *access.Client {
	return &access.Client{Client: c.a}
}

// Cluster returns a client for managing the cluster.
func (c *client) Cluster() *cluster.Client {
	return &cluster.Client{Client: c.a}
}

// Node returns a client for managing resources on a specific node.
func (c *client) Node(nodeName string) *nodes.Client {
	return &nodes.Client{Client: c.a, NodeName: nodeName}
}

// Pool returns a client for managing resource pools.
func (c *client) Pool() *pools.Client {
	return &pools.Client{Client: c.a}
}

// Storage returns a client for managing storage.
func (c *client) Storage() *storage.Client {
	return &storage.Client{Client: c.a}
}

// Version returns a client for getting the version of the Proxmox Virtual Environment API.
func (c *client) Version() *version.Client {
	return &version.Client{Client: c.a}
}

// API returns a lower-lever REST API client.
func (c *client) API() api.Client {
	return c.a
}

// SSH returns a lower-lever SSH client.s.
func (c *client) SSH() ssh.Client {
	return c.s
}
