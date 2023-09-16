/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package tests

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestLinuxVLANResource(t *testing.T) {
	t.Parallel()

	accProviders := AccMuxProviders(context.Background(), t)

	resourceName := "proxmox_virtual_environment_network_linux_vlan.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: accProviders,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + `
resource "proxmox_virtual_environment_network_linux_vlan" "test" {
	node_name = "pve"
	name = "eno0.33"
	comment = "created by terraform"
	mtu = 1499
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "eno0.33"),
					resource.TestCheckResourceAttr(resourceName, "comment", "created by terraform"),
					resource.TestCheckResourceAttr(resourceName, "vlan", "33"),
					resource.TestCheckResourceAttr(resourceName, "interface", "eno0"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update testing
			{
				Config: ProviderConfig + `
resource "proxmox_virtual_environment_network_linux_vlan" "test" {
	node_name = "pve"
	name = "eno0.33"
	address = "1.1.1.1/24"
	address6 = "FE80:0000:0000:0000:0202:B3FF:FE1E:8329/64"
	comment = "updated by terraform"
	mtu = null
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "eno0.33"),
					resource.TestCheckResourceAttr(resourceName, "address", "1.1.1.1/24"),
					resource.TestCheckResourceAttr(resourceName, "address6", "FE80:0000:0000:0000:0202:B3FF:FE1E:8329/64"),
					resource.TestCheckResourceAttr(resourceName, "comment", "updated by terraform"),
					resource.TestCheckResourceAttr(resourceName, "vlan", "33"),
					resource.TestCheckResourceAttr(resourceName, "interface", "eno0"),
					resource.TestCheckNoResourceAttr(resourceName, "mtu"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
		},
	})
}
