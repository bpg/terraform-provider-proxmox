/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package firewall

import (
	"fmt"

	"github.com/bpg/terraform-provider-proxmox/proxmox/firewall"
)

// API is an interface for managing cluster firewall
type API interface {
	firewall.API
	SecurityGroup
	Options
	SecurityGroup(group string) firewall.Rule
}

// Client is an interface for accessing the Proxmox cluster firewall API
type Client struct {
	firewall.Client
}

type groupClient struct {
	firewall.Client
	Group string
}

// SecurityGroup returns a client for managing a specific security group
func (c *Client) SecurityGroup(group string) firewall.Rule {
	// My head really hurts when I'm looking at this code
	// I'm not sure if this is the best way to do the required
	// interface composition and method "overrides", but it works.
	return &Client{
		Client: firewall.Client{
			Client: &groupClient{
				Client: c.Client,
				Group:  group,
			},
		},
	}
}

func (c *groupClient) ExpandPath(_ string) string {
	return fmt.Sprintf("cluster/firewall/groups/%s", c.Group)
}
