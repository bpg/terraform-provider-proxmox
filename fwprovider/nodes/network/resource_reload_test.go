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
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/test"
)

// checkInterfaceActive asserts whether PVE reports an interface as active on the node.
//
// This is the behavioural assertion behind `reload`: a staged-only change is present in the
// interface list (PVE serves it from /etc/network/interfaces.new) but is not up in the kernel,
// so `active` stays false until something reloads the node network configuration. Asserting on
// Terraform state alone would pass with or without the flag being honoured.
func checkInterfaceActive(te *test.Environment, iface string, want bool) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		ifaces, err := te.NodeClient().ListNetworkInterfaces(context.Background())
		if err != nil {
			return fmt.Errorf("listing network interfaces: %w", err)
		}

		for _, i := range ifaces {
			if i.Iface != iface {
				continue
			}

			got := i.Active != nil && bool(*i.Active)
			if got != want {
				return fmt.Errorf("interface %q: active = %t, want %t", iface, got, want)
			}

			return nil
		}

		return fmt.Errorf("interface %q not found on node", iface)
	}
}

func TestAccResourceLinuxBridgeReload(t *testing.T) {
	te := test.InitEnvironment(t)

	iface := fmt.Sprintf("vmbr%d", gofakeit.Number(10, 9999))
	ipV4cidr := fmt.Sprintf("%s/24", gofakeit.IPv4Address())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			// reload = false stages the bridge without activating it.
			{
				Config: te.RenderConfig(fmt.Sprintf(`
				resource "proxmox_network_linux_bridge" "test" {
					address   = "%s"
					name      = "%s"
					node_name = "{{.NodeName}}"
					reload    = false
				}`, ipV4cidr, iface)),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_network_linux_bridge.test", map[string]string{
						"reload": "false",
					}),
					checkInterfaceActive(te, iface, false),
				),
			},
			// Flipping to reload = true applies the staged change.
			{
				Config: te.RenderConfig(fmt.Sprintf(`
				resource "proxmox_network_linux_bridge" "test" {
					address   = "%s"
					name      = "%s"
					node_name = "{{.NodeName}}"
					reload    = true
				}`, ipV4cidr, iface)),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_network_linux_bridge.test", map[string]string{
						"reload": "true",
					}),
					checkInterfaceActive(te, iface, true),
				),
			},
			// Omitting the attribute keeps the computed default of true — no diff, no behaviour change.
			{
				Config: te.RenderConfig(fmt.Sprintf(`
				resource "proxmox_network_linux_bridge" "test" {
					address   = "%s"
					name      = "%s"
					node_name = "{{.NodeName}}"
				}`, ipV4cidr, iface)),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_network_linux_bridge.test", map[string]string{
						"reload": "true",
					}),
					checkInterfaceActive(te, iface, true),
				),
			},
			// The default must survive an import round-trip.
			{
				ResourceName:      "proxmox_network_linux_bridge.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"timeout_reload",
				},
			},
		},
	})
}

func TestAccResourceLinuxVLANReload(t *testing.T) {
	te := test.InitEnvironment(t)

	parent := os.Getenv("PROXMOX_VE_ACC_IFACE_NAME")
	if parent == "" {
		parent = "ens18"
	}

	vlan := gofakeit.Number(10, 4094)
	iface := fmt.Sprintf("%s.%d", parent, vlan)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			// reload = false stages the VLAN interface without activating it.
			{
				Config: te.RenderConfig(fmt.Sprintf(`
				resource "proxmox_network_linux_vlan" "test" {
					name      = "%s"
					node_name = "{{.NodeName}}"
					reload    = false
				}`, iface)),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_network_linux_vlan.test", map[string]string{
						"reload": "false",
					}),
					checkInterfaceActive(te, iface, false),
				),
			},
			// Flipping to reload = true applies the staged change.
			{
				Config: te.RenderConfig(fmt.Sprintf(`
				resource "proxmox_network_linux_vlan" "test" {
					name      = "%s"
					node_name = "{{.NodeName}}"
					reload    = true
				}`, iface)),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_network_linux_vlan.test", map[string]string{
						"reload": "true",
					}),
					checkInterfaceActive(te, iface, true),
				),
			},
		},
	})
}

func TestAccResourceLinuxBondReload(t *testing.T) {
	te := test.InitEnvironment(t)

	slave1 := os.Getenv("PROXMOX_VE_ACC_BOND_SLAVE1")
	slave2 := os.Getenv("PROXMOX_VE_ACC_BOND_SLAVE2")

	if slave1 == "" || slave2 == "" {
		t.Skip("skipping: PROXMOX_VE_ACC_BOND_SLAVE1 and PROXMOX_VE_ACC_BOND_SLAVE2 must be set to eth-type interfaces")
	}

	iface := fmt.Sprintf("bond%d", gofakeit.Number(10, 9999))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			// reload = false stages the bond without activating it.
			{
				Config: te.RenderConfig(fmt.Sprintf(`
				resource "proxmox_network_linux_bond" "test" {
					name      = "%s"
					node_name = "{{.NodeName}}"
					slaves    = ["%s", "%s"]
					reload    = false
				}`, iface, slave1, slave2)),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_network_linux_bond.test", map[string]string{
						"reload": "false",
					}),
					checkInterfaceActive(te, iface, false),
				),
			},
			// Flipping to reload = true applies the staged change.
			{
				Config: te.RenderConfig(fmt.Sprintf(`
				resource "proxmox_network_linux_bond" "test" {
					name      = "%s"
					node_name = "{{.NodeName}}"
					slaves    = ["%s", "%s"]
					reload    = true
				}`, iface, slave1, slave2)),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_network_linux_bond.test", map[string]string{
						"reload": "true",
					}),
					checkInterfaceActive(te, iface, true),
				),
			},
		},
	})
}
