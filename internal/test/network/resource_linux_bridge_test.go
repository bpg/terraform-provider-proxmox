/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package network

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/bpg/terraform-provider-proxmox/internal/test"
)

func TestLinuxBridgeResource(t *testing.T) {
	t.Parallel()

	accProviders := test.AccMuxProviders(context.Background(), t)

	resourceName := "proxmox_virtual_environment_network_linux_bridge.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: accProviders,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: test.ProviderConfig + `
resource "proxmox_virtual_environment_network_linux_bridge" "test" {
	node_name = "pve"
	name = "vmbr99"
	address = "3.3.3.3/24"
	comment = "created by terraform"
	mtu = 1499
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "vmbr99"),
					resource.TestCheckResourceAttr(resourceName, "address", "3.3.3.3/24"),
					resource.TestCheckResourceAttr(resourceName, "comment", "created by terraform"),
					resource.TestCheckResourceAttr(resourceName, "vlan_aware", "true"),
					resource.TestCheckResourceAttr(resourceName, "mtu", "1499"),
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
				Config: test.ProviderConfig + `
resource "proxmox_virtual_environment_network_linux_bridge" "test" {	
	node_name = "pve"
	name = "vmbr99"
	address = "1.1.1.1/24"
	address6 = "FE80:0000:0000:0000:0202:B3FF:FE1E:8329/64"
	comment = "updated by terraform"
	vlan_aware = false
	mtu = null
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "vmbr99"),
					resource.TestCheckResourceAttr(resourceName, "address", "1.1.1.1/24"),
					resource.TestCheckResourceAttr(resourceName, "address6", "FE80:0000:0000:0000:0202:B3FF:FE1E:8329/64"),
					resource.TestCheckResourceAttr(resourceName, "comment", "updated by terraform"),
					resource.TestCheckResourceAttr(resourceName, "vlan_aware", "false"),
					resource.TestCheckNoResourceAttr(resourceName, "mtu"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
		},
	})
}
