/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package cluster

import (
	"fmt"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	clusterfirewall "github.com/bpg/terraform-provider-proxmox/proxmox/cluster/firewall"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/ha"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/mapping"
	"github.com/bpg/terraform-provider-proxmox/proxmox/firewall"
)

// Client is an interface for accessing the Proxmox cluster API.
type Client struct {
	api.Client
}

// ExpandPath expands a relative path to a full cluster API path.
func (c *Client) ExpandPath(path string) string {
	return fmt.Sprintf("cluster/%s", path)
}

// Firewall returns a client for managing the cluster firewall.
func (c *Client) Firewall() clusterfirewall.API {
	return &clusterfirewall.Client{
		Client: firewall.Client{Client: c},
	}
}

// HA returns a client for managing the cluster's High Availability features.
func (c *Client) HA() *ha.Client {
	return &ha.Client{Client: c}
}

// HardwareMapping returns a client for managing the cluster's hardware mapping features.
func (c *Client) HardwareMapping() *mapping.Client {
	return &mapping.Client{Client: c}
}
