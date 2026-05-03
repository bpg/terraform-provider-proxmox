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
	"regexp"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/test"
)

func TestAccResourceLinuxBridge(t *testing.T) {
	te := test.InitEnvironment(t)

	iface := fmt.Sprintf("vmbr%d", gofakeit.Number(10, 9999))
	ipV4cidr1 := fmt.Sprintf("%s/24", gofakeit.IPv4Address())
	ipV4cidr2 := fmt.Sprintf("%s/24", gofakeit.IPv4Address())
	ipV6cidr := "FE80:0000:0000:0000:0202:B3FF:FE1E:8329/64"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: te.RenderConfig(fmt.Sprintf(`
				resource "proxmox_network_linux_bridge" "test" {
					address = "%s"
					autostart = true
					comment = "created by terraform"
					mtu = 1499
					name = "%s"
					node_name = "{{.NodeName}}"
					timeout_reload = 60
					vlan_aware = true
				}
				`, ipV4cidr1, iface)),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_network_linux_bridge.test", map[string]string{
						"address":        ipV4cidr1,
						"autostart":      "true",
						"comment":        "created by terraform",
						"mtu":            "1499",
						"name":           iface,
						"timeout_reload": "60",
						"vlan_aware":     "true",
					}),
					test.ResourceAttributesSet("proxmox_network_linux_bridge.test", []string{
						"id",
					}),
				),
			},
			// Update testing
			{
				Config: te.RenderConfig(fmt.Sprintf(`
				resource "proxmox_network_linux_bridge" "test" {
					address = "%s"
					address6 = "%s"
					autostart = false
					comment = ""
					mtu = null
					name = "%s"
					node_name = "{{.NodeName}}"
					timeout_reload = 60
					vlan_aware = false
				}`, ipV4cidr2, ipV6cidr, iface)),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_network_linux_bridge.test", map[string]string{
						"address":        ipV4cidr2,
						"address6":       ipV6cidr,
						"autostart":      "false",
						"comment":        "",
						"name":           iface,
						"timeout_reload": "60",
						"vlan_aware":     "false",
					}),
					test.NoResourceAttributesSet("proxmox_network_linux_bridge.test", []string{
						"mtu",
					}),
					test.ResourceAttributesSet("proxmox_network_linux_bridge.test", []string{
						"id",
					}),
				),
			},
			// Update testing (remove v4 + v6)
			{
				Config: te.RenderConfig(fmt.Sprintf(`
				resource "proxmox_network_linux_bridge" "test" {
					autostart  = false
					comment    = ""
					mtu        = null
					name       = "%s"
					node_name  = "{{.NodeName}}"
					vlan_aware = false
				}`, iface)),
				Check: resource.ComposeTestCheckFunc(
					test.NoResourceAttributesSet("proxmox_network_linux_bridge.test", []string{
						"address",
						"address6",
						"mtu",
					}),
					test.ResourceAttributes("proxmox_network_linux_bridge.test", map[string]string{
						"autostart":  "false",
						"comment":    "",
						"name":       iface,
						"vlan_aware": "false",
					}),
					test.ResourceAttributesSet("proxmox_network_linux_bridge.test", []string{
						"id",
					}),
				),
			},
			// ImportState testing
			{
				ResourceName:      "proxmox_network_linux_bridge.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"comment", // "" comments translates to null in the PVE, but nulls are not imported as empty strings.
					"timeout_reload",
				},
			},
		},
	})
}

func TestAccResourceLinuxBridgeVIDs(t *testing.T) {
	te := test.InitEnvironment(t)

	iface := fmt.Sprintf("vmbr%d", gofakeit.Number(10, 9999))
	ipV4cidr := fmt.Sprintf("%s/24", gofakeit.IPv4Address())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			// Create with a hyphenated VID range.
			{
				Config: te.RenderConfig(fmt.Sprintf(`
				resource "proxmox_network_linux_bridge" "test" {
					address    = "%s"
					name       = "%s"
					node_name  = "{{.NodeName}}"
					vlan_aware = true
					vids       = "2-4094"
				}
				`, ipV4cidr, iface)),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_network_linux_bridge.test", map[string]string{
						"vids":       "2-4094",
						"vlan_aware": "true",
					}),
				),
			},
			// Update to a space-separated VID list.
			{
				Config: te.RenderConfig(fmt.Sprintf(`
				resource "proxmox_network_linux_bridge" "test" {
					address    = "%s"
					name       = "%s"
					node_name  = "{{.NodeName}}"
					vlan_aware = true
					vids       = "10 20 30"
				}
				`, ipV4cidr, iface)),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_network_linux_bridge.test", map[string]string{
						"vids": "10 20 30",
					}),
				),
			},
			// Remove vids from config — vidsPlanModifier preserves the prior state value when vlan_aware stays
			// true. Users explicitly reset by setting vids = "2-4094".
			{
				Config: te.RenderConfig(fmt.Sprintf(`
				resource "proxmox_network_linux_bridge" "test" {
					address    = "%s"
					name       = "%s"
					node_name  = "{{.NodeName}}"
					vlan_aware = true
				}
				`, ipV4cidr, iface)),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_network_linux_bridge.test", map[string]string{
						"vids": "10 20 30",
					}),
				),
			},
			// Reset vids to the implicit PVE default by setting it explicitly.
			{
				Config: te.RenderConfig(fmt.Sprintf(`
				resource "proxmox_network_linux_bridge" "test" {
					address    = "%s"
					name       = "%s"
					node_name  = "{{.NodeName}}"
					vlan_aware = true
					vids       = "2-4094"
				}
				`, ipV4cidr, iface)),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_network_linux_bridge.test", map[string]string{
						"vids": "2-4094",
					}),
				),
			},
			// ImportState testing — round-trip with vids set.
			{
				Config: te.RenderConfig(fmt.Sprintf(`
				resource "proxmox_network_linux_bridge" "test" {
					address    = "%s"
					name       = "%s"
					node_name  = "{{.NodeName}}"
					vlan_aware = true
					vids       = "100-200"
				}
				`, ipV4cidr, iface)),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_network_linux_bridge.test", map[string]string{
						"vids": "100-200",
					}),
				),
			},
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

// TestAccResourceLinuxBridgeVIDsToggleVLANAware exercises the plan modifier
// branch that nulls vids when vlan_aware flips from true to false, and the
// CheckDelete plumbing that removes bridge_vids from PVE on that transition.
func TestAccResourceLinuxBridgeVIDsToggleVLANAware(t *testing.T) {
	te := test.InitEnvironment(t)

	iface := fmt.Sprintf("vmbr%d", gofakeit.Number(10, 9999))
	ipV4cidr := fmt.Sprintf("%s/24", gofakeit.IPv4Address())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			// Create with vlan_aware = true and an explicit vids list.
			{
				Config: te.RenderConfig(fmt.Sprintf(`
				resource "proxmox_network_linux_bridge" "test" {
					address    = "%s"
					name       = "%s"
					node_name  = "{{.NodeName}}"
					vlan_aware = true
					vids       = "100 200"
				}
				`, ipV4cidr, iface)),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_network_linux_bridge.test", map[string]string{
						"vlan_aware": "true",
						"vids":       "100 200",
					}),
				),
			},
			// Flip vlan_aware to false and drop vids from config — the plan modifier
			// must null vids (rather than pin the prior "100 200" via state) and the
			// Update path must send bridge_vids in the delete list.
			{
				Config: te.RenderConfig(fmt.Sprintf(`
				resource "proxmox_network_linux_bridge" "test" {
					address    = "%s"
					name       = "%s"
					node_name  = "{{.NodeName}}"
					vlan_aware = false
				}
				`, ipV4cidr, iface)),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_network_linux_bridge.test", map[string]string{
						"vlan_aware": "false",
					}),
					test.NoResourceAttributesSet("proxmox_network_linux_bridge.test", []string{
						"vids",
					}),
				),
			},
		},
	})
}

// Regression test for #2851: `ports` referencing an unknown value must not
// break `tofu validate` / `terraform plan`.
func TestAccResourceLinuxBridgeUnknownPorts(t *testing.T) {
	te := test.InitEnvironment(t)

	iface := fmt.Sprintf("vmbr%d", gofakeit.Number(10, 9999))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: te.RenderConfig(fmt.Sprintf(`
				resource "terraform_data" "ports_source" {
					input = ["enp1s0"]
				}

				resource "proxmox_network_linux_bridge" "test" {
					name      = "%s"
					node_name = "{{.NodeName}}"
					ports     = terraform_data.ports_source.output
				}
				`, iface)),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccResourceLinuxBridgeVIDsValidation(t *testing.T) {
	te := test.InitEnvironment(t)

	iface := fmt.Sprintf("vmbr%d", gofakeit.Number(10, 9999))
	ipV4cidr := fmt.Sprintf("%s/24", gofakeit.IPv4Address())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			// `vids = ""` rejected by per-attribute LengthAtLeast(1) validator.
			{
				Config: te.RenderConfig(fmt.Sprintf(`
				resource "proxmox_network_linux_bridge" "test" {
					address    = "%s"
					name       = "%s"
					node_name  = "{{.NodeName}}"
					vlan_aware = true
					vids       = ""
				}
				`, ipV4cidr, iface)),
				ExpectError: regexp.MustCompile(`(?s)at least 1`),
				PlanOnly:    true,
			},
			// vids set + vlan_aware omitted → ValidateConfig still errors. PVE
			// defaults vlan_aware to false, so the misconfiguration is real;
			// only an unresolved (unknown) vlan_aware should be skipped.
			{
				Config: te.RenderConfig(fmt.Sprintf(`
				resource "proxmox_network_linux_bridge" "test" {
					address   = "%s"
					name      = "%s"
					node_name = "{{.NodeName}}"
					vids      = "1 2 3"
				}
				`, ipV4cidr, iface)),
				ExpectError: regexp.MustCompile(`(?s)requires.*vlan_aware`),
				PlanOnly:    true,
			},
			// vids set + vlan_aware explicitly false → ValidateConfig errors.
			// Covers the migration scenario: a user toggling vlan_aware off
			// without first removing vids.
			{
				Config: te.RenderConfig(fmt.Sprintf(`
				resource "proxmox_network_linux_bridge" "test" {
					address    = "%s"
					name       = "%s"
					node_name  = "{{.NodeName}}"
					vlan_aware = false
					vids       = "1 2 3"
				}
				`, ipV4cidr, iface)),
				ExpectError: regexp.MustCompile(`(?s)requires.*vlan_aware`),
				PlanOnly:    true,
			},
			// Comma-separated list rejected by per-attribute RegexMatches validator.
			// Catches the most common malformed input (users reaching for CSV).
			{
				Config: te.RenderConfig(fmt.Sprintf(`
				resource "proxmox_network_linux_bridge" "test" {
					address    = "%s"
					name       = "%s"
					node_name  = "{{.NodeName}}"
					vlan_aware = true
					vids       = "1,2,3"
				}
				`, ipV4cidr, iface)),
				ExpectError: regexp.MustCompile(`(?s)space-separated list of VLAN IDs`),
				PlanOnly:    true,
			},
		},
	})
}
