//go:build acceptance || all

//testacc:tier=heavy
//testacc:resource=network

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package network_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/test"
)

const (
	accTestLinuxBondName = "proxmox_network_linux_bond.test"
)

func TestAccResourceLinuxBond(t *testing.T) {
	te := test.InitEnvironment(t)

	iface := fmt.Sprintf("bond%d", gofakeit.Number(10, 9999))

	// Bond slave interfaces must be eth-type (or bond-type) interfaces not already in use.
	// Set PROXMOX_VE_ACC_BOND_SLAVE1 and PROXMOX_VE_ACC_BOND_SLAVE2 to available interfaces.
	slave1 := os.Getenv("PROXMOX_VE_ACC_BOND_SLAVE1")
	slave2 := os.Getenv("PROXMOX_VE_ACC_BOND_SLAVE2")

	if slave1 == "" || slave2 == "" {
		t.Skip("skipping: PROXMOX_VE_ACC_BOND_SLAVE1 and PROXMOX_VE_ACC_BOND_SLAVE2 must be set to eth-type interfaces")
	}

	ipV4cidr1 := fmt.Sprintf("%s/24", gofakeit.IPv4Address())
	ipV4cidr2 := fmt.Sprintf("%s/24", gofakeit.IPv4Address())

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			// Step 1: Create with 802.3ad mode, hash policy, address, comment
			{
				Config: te.RenderConfig(fmt.Sprintf(`
				resource "proxmox_network_linux_bond" "test" {
					address               = "%s"
					autostart             = true
					comment               = "created by terraform"
					name                  = "%s"
					node_name             = "{{.NodeName}}"
					slaves                = ["%s", "%s"]
					bond_mode             = "802.3ad"
					bond_xmit_hash_policy = "layer3+4"
					timeout_reload        = 60
				}
				`, ipV4cidr1, iface, slave1, slave2)),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes(accTestLinuxBondName, map[string]string{
						"address":               ipV4cidr1,
						"autostart":             "true",
						"comment":               "created by terraform",
						"name":                  iface,
						"bond_mode":             "802.3ad",
						"bond_xmit_hash_policy": `layer3\+4`,
						"slaves.#":              "2",
						"timeout_reload":        "60",
					}),
					test.ResourceAttributesSet(accTestLinuxBondName, []string{
						"id",
					}),
				),
			},
			// Step 2: Update to active-backup mode with bond_primary, remove hash policy
			{
				Config: te.RenderConfig(fmt.Sprintf(`
				resource "proxmox_network_linux_bond" "test" {
					address        = "%s"
					autostart      = false
					comment        = "updated comment"
					name           = "%s"
					node_name      = "{{.NodeName}}"
					slaves         = ["%s", "%s"]
					bond_mode      = "active-backup"
					bond_primary   = "%s"
					timeout_reload = 60
				}`, ipV4cidr2, iface, slave1, slave2, slave1)),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes(accTestLinuxBondName, map[string]string{
						"address":      ipV4cidr2,
						"autostart":    "false",
						"comment":      "updated comment",
						"name":         iface,
						"bond_mode":    "active-backup",
						"bond_primary": slave1,
						"slaves.#":     "2",
					}),
					test.NoResourceAttributesSet(accTestLinuxBondName, []string{
						"bond_xmit_hash_policy",
					}),
				),
			},
			// Step 3: Remove address, comment, bond_primary; omit bond_mode to test Computed default
			{
				Config: te.RenderConfig(fmt.Sprintf(`
				resource "proxmox_network_linux_bond" "test" {
					autostart      = false
					name           = "%s"
					node_name      = "{{.NodeName}}"
					slaves         = ["%s", "%s"]
					timeout_reload = 60
				}`, iface, slave1, slave2)),
				Check: resource.ComposeTestCheckFunc(
					test.NoResourceAttributesSet(accTestLinuxBondName, []string{
						"address",
						"comment",
						"bond_primary",
						"bond_xmit_hash_policy",
					}),
					test.ResourceAttributes(accTestLinuxBondName, map[string]string{
						"autostart": "false",
						"name":      iface,
						"slaves.#":  "2",
					}),
					test.ResourceAttributesSet(accTestLinuxBondName, []string{
						"id",
						"bond_mode",
					}),
				),
			},
			// Step 4: ImportState
			{
				ResourceName:      accTestLinuxBondName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"timeout_reload",
				},
			},
		},
	})
}
