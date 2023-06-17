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

func TestInterfaceLinuxBridgeResource(t *testing.T) {
	resourceName := "proxmox_virtual_environment_network_linux_bridge.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + `
resource "proxmox_virtual_environment_network_linux_bridge" "test" {
	node_name = "pve"
	iface = "vmbr99"
	address = "3.3.3.3/24"
	comment = "created by terraform"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "iface", "vmbr99"),
					resource.TestCheckResourceAttr(resourceName, "address", "3.3.3.3/24"),
					resource.TestCheckResourceAttr(resourceName, "comment", "created by terraform"),
					resource.TestCheckResourceAttr(resourceName, "bridge_vlan_aware", "true"),
				),
			},
		},
	})
}
