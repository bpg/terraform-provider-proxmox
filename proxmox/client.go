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

// API is the interface for the Proxmox Virtual Environment API.
type API interface {
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

	// RESTAPI returns a lower-lever REST API client.
	RESTAPI() api.Client

	// SSHAPI returns a lower-lever SSH client.
	SSHAPI() ssh.Client
}

// Client is an implementation of the Proxmox Virtual Environment API.
type Client struct {
	a api.Client
	s ssh.Client
}

// NewClient creates a new API client.
func NewClient(a api.Client, s ssh.Client) *Client {
	return &Client{a: a, s: s}
}

// Access returns a client for managing access control.
func (c *Client) Access() *access.Client {
	return &access.Client{Client: c.a}
}

// Cluster returns a client for managing the cluster.
func (c *Client) Cluster() *cluster.Client {
	return &cluster.Client{Client: c.a}
}

// Node returns a client for managing resources on a specific node.
func (c *Client) Node(nodeName string) *nodes.Client {
	return &nodes.Client{Client: c.a, NodeName: nodeName}
}

// Pool returns a client for managing resource pools.
func (c *Client) Pool() *pools.Client {
	return &pools.Client{Client: c.a}
}

// Storage returns a client for managing storage.
func (c *Client) Storage() *storage.Client {
	return &storage.Client{Client: c.a}
}

// Version returns a client for getting the version of the Proxmox Virtual Environment API.
func (c *Client) Version() *version.Client {
	return &version.Client{Client: c.a}
}

// RESTAPI returns a lower-lever REST API client.
func (c *Client) RESTAPI() api.Client {
	return c.a
}

// SSHAPI returns a lower-lever SSH client.s.
func (c *Client) SSHAPI() ssh.Client {
	return c.s
}
