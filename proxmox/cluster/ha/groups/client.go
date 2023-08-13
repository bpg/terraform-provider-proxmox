/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package groups

import (
	"fmt"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
)

// Client is an interface for accessing the Proxmox High Availability groups API.
type Client struct {
	api.Client
}

// ExpandPath expands a relative path to the HA groups management API path.
func (c *Client) ExpandPath(path string) string {
	return fmt.Sprintf("cluster/ha/groups/%s", path)
}
