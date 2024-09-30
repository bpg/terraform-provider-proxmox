/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package config

import "github.com/bpg/terraform-provider-proxmox/proxmox"

// DataSource is the global configuration for all datasources.
type DataSource struct {
	Client proxmox.Client
}
