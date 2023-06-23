/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func Test_LinuxVLANResource(t *testing.T) {
	t.Parallel()

	resourceName := "proxmox_virtual_environment_network_linux_vlan.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + `
resource "proxmox_virtual_environment_network_linux_vlan" "test" {
	node_name = "pve"
	name = "ens18.33"
	comment = "created by terraform"
	mtu = 1499
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "ens18.33"),
					resource.TestCheckResourceAttr(resourceName, "comment", "created by terraform"),
					resource.TestCheckResourceAttr(resourceName, "vlan", "33"),
					resource.TestCheckResourceAttr(resourceName, "interface", "ens18"),
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
	name = "ens18.33"
	address = "1.1.1.1/24"
	address6 = "FE80:0000:0000:0000:0202:B3FF:FE1E:8329/64"
	comment = "updated by terraform"
	mtu = null
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "ens18.33"),
					resource.TestCheckResourceAttr(resourceName, "address", "1.1.1.1/24"),
					resource.TestCheckResourceAttr(resourceName, "address6", "FE80:0000:0000:0000:0202:B3FF:FE1E:8329/64"),
					resource.TestCheckResourceAttr(resourceName, "comment", "updated by terraform"),
					resource.TestCheckResourceAttr(resourceName, "vlan", "33"),
					resource.TestCheckResourceAttr(resourceName, "interface", "ens18"),
					resource.TestCheckNoResourceAttr(resourceName, "mtu"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
		},
	})
}
