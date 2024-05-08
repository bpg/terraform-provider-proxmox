/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package tests

import (
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v7"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceLinuxBridge(t *testing.T) {
	te := initTestEnvironment(t)

	iface := fmt.Sprintf("vmbr%d", gofakeit.Number(10, 9999))
	ipV4cidr1 := fmt.Sprintf("%s/24", gofakeit.IPv4Address())
	ipV4cidr2 := fmt.Sprintf("%s/24", gofakeit.IPv4Address())
	ipV6cidr := "FE80:0000:0000:0000:0202:B3FF:FE1E:8329/64"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: te.accProviders,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: te.renderConfig(fmt.Sprintf(`
				resource "proxmox_virtual_environment_network_linux_bridge" "test" {
					address = "%s"
					autostart = true
					comment = "created by terraform"
					mtu = 1499
					name = "%s"
					node_name = "{{.NodeName}}"
					vlan_aware = true
				}
				`, ipV4cidr1, iface)),
				Check: resource.ComposeTestCheckFunc(
					testResourceAttributes("proxmox_virtual_environment_network_linux_bridge.test", map[string]string{
						"address":    ipV4cidr1,
						"autostart":  "true",
						"comment":    "created by terraform",
						"mtu":        "1499",
						"name":       iface,
						"vlan_aware": "true",
					}),
					testResourceAttributesSet("proxmox_virtual_environment_network_linux_bridge.test", []string{
						"id",
					}),
				),
			},
			// Update testing
			{
				Config: te.renderConfig(fmt.Sprintf(`
				resource "proxmox_virtual_environment_network_linux_bridge" "test" {
					address = "%s"
					address6 = "%s"
					autostart = false
					comment = ""
					mtu = null
					name = "%s"
					node_name = "{{.NodeName}}"
					vlan_aware = false
				}`, ipV4cidr2, ipV6cidr, iface)),
				Check: resource.ComposeTestCheckFunc(
					testResourceAttributes("proxmox_virtual_environment_network_linux_bridge.test", map[string]string{
						"address":    ipV4cidr2,
						"address6":   ipV6cidr,
						"autostart":  "false",
						"comment":    "",
						"name":       iface,
						"vlan_aware": "false",
					}),
					testNoResourceAttributesSet("proxmox_virtual_environment_network_linux_bridge.test", []string{
						"mtu",
					}),
					testResourceAttributesSet("proxmox_virtual_environment_network_linux_bridge.test", []string{
						"id",
					}),
				),
			},
			// ImportState testing
			{
				ResourceName:      "proxmox_virtual_environment_network_linux_bridge.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"comment", // "" comments translates to null in the PVE, but nulls are not imported as empty strings.
				},
			},
		},
	})
}
