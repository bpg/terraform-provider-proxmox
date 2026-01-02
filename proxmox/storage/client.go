/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package storage

import (
	"fmt"
	"net/url"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
)

// Client provides access to the Proxmox VE cluster-level storage configuration API.
type Client struct {
	api.Client
}

func (c *Client) basePath() string {
	return c.Client.ExpandPath("storage")
}

// ExpandPath expands a relative path to a full storage API path.
func (c *Client) ExpandPath(path string) string {
	ep := c.basePath()
	if path == "" {
		return ep
	}

	return fmt.Sprintf("%s/%s", ep, url.PathEscape(path))
}
