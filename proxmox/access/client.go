/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package access

import (
	"fmt"

	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

// Client is an interface for performing requests against the Proxmox 'access' API.
type Client struct {
	types.Client
}

// ExpandPath expands a path relative to the client's base path.
func (c *Client) ExpandPath(path string) string {
	return fmt.Sprintf("access/%s", path)
}
