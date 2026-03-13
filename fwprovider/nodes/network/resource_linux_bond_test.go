//go:build acceptance || all

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package network_test

import (
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/test"
)

func TestAccResourceLinuxBond(t *testing.T) {
	te := test.InitEnvironment(t)

	iface := fmt.Sprintf("bond%d", gofakeit.Number(10, 9999))
	ipV4cidr1 := fmt.Sprintf("%s/24", gofakeit.IPv4Address())
	ipV4cidr2 := fmt.Sprintf("%s/24", gofakeit.IPv4Address())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: te.RenderConfig(fmt.Sprintf(`
				resource "proxmox_virtual_environment_network_linux_bond" "test" {
					address   = "%s"
					autostart = true
					comment   = "created by terraform"
					name      = "%s"
					node_name = "{{.NodeName}}"
					slaves    = ["eth0", "eth1"]
					bond_mode = "802.3ad"
					bond_xmit_hash_policy = "layer3+4"
					timeout_reload = 60
				}
				`, ipV4cidr1, iface)),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_virtual_environment_network_linux_bond.test", map[string]string{
						"address":               ipV4cidr1,
						"autostart":             "true",
						"comment":               "created by terraform",
						"name":                  iface,
						"bond_mode":             "802.3ad",
						"bond_xmit_hash_policy": "layer3+4",
						"slaves.#":              "2",
						"timeout_reload":        "60",
					}),
					test.ResourceAttributesSet("proxmox_virtual_environment_network_linux_bond.test", []string{
						"id",
					}),
				),
			},
			// Update testing
			{
				Config: te.RenderConfig(fmt.Sprintf(`
				resource "proxmox_virtual_environment_network_linux_bond" "test" {
					address   = "%s"
					autostart = false
					comment   = ""
					name      = "%s"
					node_name = "{{.NodeName}}"
					slaves    = ["eth0", "eth1"]
					bond_mode = "active-backup"
					bond_primary = "eth0"
					timeout_reload = 60
				}`, ipV4cidr2, iface)),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_virtual_environment_network_linux_bond.test", map[string]string{
						"address":      ipV4cidr2,
						"autostart":    "false",
						"comment":      "",
						"name":         iface,
						"bond_mode":    "active-backup",
						"bond_primary": "eth0",
						"slaves.#":     "2",
					}),
					test.NoResourceAttributesSet("proxmox_virtual_environment_network_linux_bond.test", []string{
						"bond_xmit_hash_policy",
					}),
					test.ResourceAttributesSet("proxmox_virtual_environment_network_linux_bond.test", []string{
						"id",
					}),
				),
			},
			// Update testing (remove address)
			{
				Config: te.RenderConfig(fmt.Sprintf(`
				resource "proxmox_virtual_environment_network_linux_bond" "test" {
					autostart = false
					comment   = ""
					name      = "%s"
					node_name = "{{.NodeName}}"
					slaves    = ["eth0", "eth1"]
					bond_mode = "balance-rr"
				}`, iface)),
				Check: resource.ComposeTestCheckFunc(
					test.NoResourceAttributesSet("proxmox_virtual_environment_network_linux_bond.test", []string{
						"address",
						"bond_primary",
						"bond_xmit_hash_policy",
					}),
					test.ResourceAttributes("proxmox_virtual_environment_network_linux_bond.test", map[string]string{
						"autostart": "false",
						"comment":   "",
						"name":      iface,
						"bond_mode": "balance-rr",
						"slaves.#":  "2",
					}),
					test.ResourceAttributesSet("proxmox_virtual_environment_network_linux_bond.test", []string{
						"id",
					}),
				),
			},
			// ImportState testing
			{
				ResourceName:      "proxmox_virtual_environment_network_linux_bond.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"comment",
					"timeout_reload",
				},
			},
		},
	})
}
