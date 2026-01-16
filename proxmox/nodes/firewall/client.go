/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package firewall

import (
	"github.com/bpg/terraform-provider-proxmox/proxmox/firewall"
)

// API is an interface for managing node firewall.
type API interface {
	firewall.API
	Options
}

// Client is an interface for accessing the Proxmox node firewall API.
type Client struct {
	firewall.Client
}
