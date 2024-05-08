/*
	This Source Code Form is subject to the terms of the Mozilla Public
	License, v. 2.0. If a copy of the MPL was not distributed with this
	file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/

package mapping

import (
	"fmt"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	proxmoxtypes "github.com/bpg/terraform-provider-proxmox/proxmox/types/hardwaremapping"
)

// Client is an interface for accessing the Proxmox cluster mapping API.
type Client struct {
	api.Client
}

func (c *Client) basePath() string {
	return c.Client.ExpandPath("mapping")
}

// ExpandPath expands a relative path to a full hardware mapping API path.
func (c *Client) ExpandPath(hmType proxmoxtypes.Type, path string) string {
	ep := c.basePath()
	if hmType.String() != "" {
		ep = fmt.Sprintf("%s/%s", ep, hmType.String())
	}

	if path != "" {
		ep = fmt.Sprintf("%s/%s", ep, path)
	}

	return ep
}
