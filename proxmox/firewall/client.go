/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package firewall

import "github.com/bpg/terraform-provider-proxmox/proxmox/types"

// API is an interface for the Proxmox firewall API
type API interface {
	Alias
	IPSet
	Rule
	Options
}

// Client is an interface for accessing the Proxmox firewall API
type Client struct {
	types.Client
}
