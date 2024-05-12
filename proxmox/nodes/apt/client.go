/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package apt

import (
	"fmt"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/apt/repositories"
)

// Client is an interface for accessing the Proxmox cluster API.
type Client struct {
	api.Client
}

// ExpandPath expands a relative path to a full cluster API path.
func (c *Client) ExpandPath(path string) string {
	return c.Client.ExpandPath(fmt.Sprintf("apt/%s", path))
}

// Repositories returns a client for managing APT repositories.
func (c *Client) Repositories() *repositories.Client {
	return &repositories.Client{Client: c}
}
