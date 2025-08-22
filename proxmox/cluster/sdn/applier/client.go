/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package applier

import (
	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
)

// Client is a client for accessing the Proxmox SDN Apply API.
type Client struct {
	api.Client
}

// ExpandPath returns the API path for cluster-wide SDN apply.
func (c *Client) ExpandPath(path string) string {
	_ = path
	return "cluster/sdn"
}
