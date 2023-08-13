/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package ha

import (
	"fmt"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	hagroups "github.com/bpg/terraform-provider-proxmox/proxmox/cluster/ha/groups"
)

// Client is an interface for accessing the Proxmox High Availability API.
type Client struct {
	api.Client
}

// ExpandPath expands a relative path to a full cluster API path.
func (c *Client) ExpandPath(path string) string {
	return fmt.Sprintf("cluster/ha/%s", path)
}

// Groups returns a client for managing the cluster's High Availability groups.
func (c *Client) Groups() *hagroups.Client {
	return &hagroups.Client{Client: c.Client}
}
