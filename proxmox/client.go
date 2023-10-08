/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package proxmox

import (
	"os"

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

	// API returns a lower-level REST API client.
	API() api.Client

	// SSH returns a lower-level SSH client.
	SSH() ssh.Client

	// TempDir returns (possibly overridden) os.TempDir().
	TempDir() string
}

type client struct {
	apiClient      api.Client
	sshClient      ssh.Client
	tmpDirOverride string
}

// NewClient creates a new API client.
func NewClient(apiClient api.Client, sshClient ssh.Client, tmpDirOverride string) Client {
	return &client{apiClient: apiClient, sshClient: sshClient, tmpDirOverride: tmpDirOverride}
}

// Access returns a client for managing access control.
func (c *client) Access() *access.Client {
	return &access.Client{Client: c.apiClient}
}

// Cluster returns a client for managing the cluster.
func (c *client) Cluster() *cluster.Client {
	return &cluster.Client{Client: c.apiClient}
}

// Node returns a client for managing resources on a specific node.
func (c *client) Node(nodeName string) *nodes.Client {
	return &nodes.Client{Client: c.apiClient, NodeName: nodeName}
}

// Pool returns a client for managing resource pools.
func (c *client) Pool() *pools.Client {
	return &pools.Client{Client: c.apiClient}
}

// Storage returns a client for managing storage.
func (c *client) Storage() *storage.Client {
	return &storage.Client{Client: c.apiClient}
}

// Version returns a client for getting the version of the Proxmox Virtual Environment API.
func (c *client) Version() *version.Client {
	return &version.Client{Client: c.apiClient}
}

// API returns a lower-lever REST API client.
func (c *client) API() api.Client {
	return c.apiClient
}

// SSH returns a lower-lever SSH client.
func (c *client) SSH() ssh.Client {
	return c.sshClient
}

// TempDir returns (possibly overridden) os.TempDir().
func (c *client) TempDir() string {
	if c.tmpDirOverride != "" {
		return c.tmpDirOverride
	}

	return os.TempDir()
}
