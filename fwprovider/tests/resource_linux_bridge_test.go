/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v6"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	accTestLinuxBridgeName = "proxmox_virtual_environment_network_linux_bridge.test"
)

func TestAccResourceLinuxBridge(t *testing.T) {
	t.Parallel()

	accProviders := testAccMuxProviders(context.Background(), t)

	iface := fmt.Sprintf("vmbr%d", gofakeit.Number(10, 9999))
	ipV4cidr1 := fmt.Sprintf("%s/24", gofakeit.IPv4Address())
	ipV4cidr2 := fmt.Sprintf("%s/24", gofakeit.IPv4Address())
	ipV6cidr := "FE80:0000:0000:0000:0202:B3FF:FE1E:8329/64"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: accProviders,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccResourceLinuxBridgeCreatedConfig(iface, ipV4cidr1),
				Check:  testAccResourceLinuxBridgeCreatedCheck(iface, ipV4cidr1),
			},
			// Update testing
			{
				Config: testAccResourceLinuxBridgeUpdatedConfig(iface, ipV4cidr2, ipV6cidr),
				Check:  testAccResourceLinuxBridgeUpdatedCheck(iface, ipV4cidr2, ipV6cidr),
			},
			// ImportState testing
			{
				ResourceName:      accTestLinuxBridgeName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccResourceLinuxBridgeCreatedConfig(name string, ipV4cidr string) string {
	return fmt.Sprintf(`
	resource "proxmox_virtual_environment_network_linux_bridge" "test" {
		node_name = "%s"
		name = "%s"
		address = "%s"
		comment = "created by terraform"
		vlan_aware = true
		autostart = true
		mtu = 1499
	}
	`, accTestNodeName, name, ipV4cidr)
}

func testAccResourceLinuxBridgeCreatedCheck(name string, ipV4cidr string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(accTestLinuxBridgeName, "name", name),
		resource.TestCheckResourceAttr(accTestLinuxBridgeName, "address", ipV4cidr),
		resource.TestCheckResourceAttr(accTestLinuxBridgeName, "comment", "created by terraform"),
		resource.TestCheckResourceAttr(accTestLinuxBridgeName, "vlan_aware", "true"),
		resource.TestCheckResourceAttr(accTestLinuxBridgeName, "autostart", "true"),
		resource.TestCheckResourceAttr(accTestLinuxBridgeName, "mtu", "1499"),
		resource.TestCheckResourceAttrSet(accTestLinuxBridgeName, "id"),
	)
}

func testAccResourceLinuxBridgeUpdatedConfig(name string, ipV4cidr string, ipV6cidr string) string {
	return fmt.Sprintf(`
	resource "proxmox_virtual_environment_network_linux_bridge" "test" {
		node_name = "%s"
		name = "%s"
		address = "%s"
		address6 = "%s"
		comment = "updated by terraform"
		vlan_aware = false
		autostart = false
		mtu = null
	}
	`, accTestNodeName, name, ipV4cidr, ipV6cidr)
}

func testAccResourceLinuxBridgeUpdatedCheck(name string, ipV4cidr string, ipV6cidr string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr(accTestLinuxBridgeName, "name", name),
			resource.TestCheckResourceAttr(accTestLinuxBridgeName, "address", ipV4cidr),
			resource.TestCheckResourceAttr(accTestLinuxBridgeName, "address6", ipV6cidr),
			resource.TestCheckResourceAttr(accTestLinuxBridgeName, "comment", "updated by terraform"),
			resource.TestCheckResourceAttr(accTestLinuxBridgeName, "vlan_aware", "false"),
			resource.TestCheckResourceAttr(accTestLinuxBridgeName, "autostart", "false"),
			resource.TestCheckNoResourceAttr(accTestLinuxBridgeName, "mtu"),
			resource.TestCheckResourceAttrSet(accTestLinuxBridgeName, "id"),
		),
	)
}
