//go:build acceptance || all

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package applier_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/test"
)

func TestAccResourceSDNApplier(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"create zones, apply, destroy zones, apply again", []resource.TestStep{
			{
				Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_sdn_zone_simple" "test_zone_simple" {
						id    = "testZS"
						nodes = ["{{.NodeName}}"]
						mtu   = 1500
					}

					resource "proxmox_virtual_environment_sdn_zone_vlan" "test_zone_vlan" {
						id     = "testZV"
						nodes  = ["{{.NodeName}}"]
						mtu    = 1500
						bridge = "vmbr0"
					}

					resource "proxmox_virtual_environment_sdn_applier" "test_applier" {
						depends_on = [
							proxmox_virtual_environment_sdn_zone_simple.test_zone_simple,
							proxmox_virtual_environment_sdn_zone_vlan.test_zone_vlan
						]
					}

					data "proxmox_virtual_environment_sdn_zone_simple" "test_zone_simple" {
						id = "testZS"
						depends_on = [
							proxmox_virtual_environment_sdn_applier.test_applier
						]
					}

					data "proxmox_virtual_environment_sdn_zone_vlan" "test_zone_vlan" {
						id = "testZV"
						depends_on = [
							proxmox_virtual_environment_sdn_applier.test_applier
						]
					}
				`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributesSet("proxmox_virtual_environment_sdn_applier.test_applier", []string{
						"id",
					}),
					test.ResourceAttributes("data.proxmox_virtual_environment_sdn_zone_simple.test_zone_simple", map[string]string{
						"id":  "testZS",
						"mtu": "1500",
					}),
					test.ResourceAttributes("data.proxmox_virtual_environment_sdn_zone_vlan.test_zone_vlan", map[string]string{
						"id":     "testZV",
						"mtu":    "1500",
						"bridge": "vmbr0",
					}),
					test.NoResourceAttributesSet("data.proxmox_virtual_environment_sdn_zone_simple.test_zone_simple", []string{
						"pending",
					}),
					test.NoResourceAttributesSet("data.proxmox_virtual_environment_sdn_zone_vlan.test_zone_vlan", []string{
						"pending",
					}),
				),
			},
			{
				Config: te.RenderConfig(`
					# Zones are destroyed, applier should still work
					resource "proxmox_virtual_environment_sdn_applier" "test_applier" {
					}
				`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributesSet("proxmox_virtual_environment_sdn_applier.test_applier", []string{
						"id",
					}),
					// Zones are destroyed - no specific checks needed as Terraform handles state removal
				),
			},
		}},
		{"applier with zone updates and replace_triggered_by", []resource.TestStep{
			{
				Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_sdn_zone_simple" "update_zone" {
						id    = "updZ"
						nodes = ["{{.NodeName}}"]
						mtu   = 1500
					}

					resource "proxmox_virtual_environment_sdn_applier" "update_applier" {
						lifecycle {
							replace_triggered_by = [
								proxmox_virtual_environment_sdn_zone_simple.update_zone
							]
						}
					}
				`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributesSet("proxmox_virtual_environment_sdn_applier.update_applier", []string{
						"id",
					}),
					test.ResourceAttributes("proxmox_virtual_environment_sdn_zone_simple.update_zone", map[string]string{
						"id":  "updZ",
						"mtu": "1500",
					}),
				),
			},
			{
				Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_sdn_zone_simple" "update_zone" {
						id    = "updZ"
						nodes = ["{{.NodeName}}"]
						mtu   = 1450  # Changed MTU triggers applier replacement
						depends_on = [
							proxmox_virtual_environment_sdn_applier.finalizer
						]
					}	

					resource "proxmox_virtual_environment_sdn_applier" "update_applier" {
						lifecycle {
							replace_triggered_by = [
								proxmox_virtual_environment_sdn_zone_simple.update_zone
							]
						}
					}

					data "proxmox_virtual_environment_sdn_zone_simple" "update_zone" {
						id = "updZ"
						depends_on = [
							proxmox_virtual_environment_sdn_applier.update_applier
						]
					}

					resource "proxmox_virtual_environment_sdn_applier" "finalizer" {}
				`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributesSet("proxmox_virtual_environment_sdn_applier.update_applier", []string{
						"id",
					}),
					test.ResourceAttributes("data.proxmox_virtual_environment_sdn_zone_simple.update_zone", map[string]string{
						"id":  "updZ",
						"mtu": "1450",
					}),
					test.NoResourceAttributesSet("data.proxmox_virtual_environment_sdn_zone_simple.update_zone", []string{
						"pending",
					}),
				),
			},
		}},
		{"applier with complex zone lifecycle", []resource.TestStep{
			{
				Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_sdn_zone_simple" "lifecycle_zone_1" {
						id    = "lifeZ1"
						nodes = ["{{.NodeName}}"]
						mtu   = 1500
						depends_on = [
							proxmox_virtual_environment_sdn_applier.finalizer
						]
					}

					resource "proxmox_virtual_environment_sdn_applier" "lifecycle_applier" {
						lifecycle {
							replace_triggered_by = [
								proxmox_virtual_environment_sdn_zone_simple.lifecycle_zone_1
							]
						}
					}

					resource "proxmox_virtual_environment_sdn_applier" "finalizer" {}
				`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributesSet("proxmox_virtual_environment_sdn_applier.lifecycle_applier", []string{
						"id",
					}),
					test.ResourceAttributes("proxmox_virtual_environment_sdn_zone_simple.lifecycle_zone_1", map[string]string{
						"id":  "lifeZ1",
						"mtu": "1500",
					}),
				),
			},
			{
				Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_sdn_zone_simple" "lifecycle_zone_1" {
						id    = "lifeZ1"
						nodes = ["{{.NodeName}}"]
						mtu   = 1450
						depends_on = [
							proxmox_virtual_environment_sdn_applier.finalizer
						]
					}

					resource "proxmox_virtual_environment_sdn_zone_vlan" "lifecycle_zone_2" {
						id     = "lifeZ2"
						nodes  = ["{{.NodeName}}"]
						mtu    = 1500
						bridge = "vmbr0"
						depends_on = [
							proxmox_virtual_environment_sdn_applier.finalizer
						]
					}

					resource "proxmox_virtual_environment_sdn_applier" "lifecycle_applier" {
						lifecycle {
							replace_triggered_by = [
								proxmox_virtual_environment_sdn_zone_simple.lifecycle_zone_1,
								proxmox_virtual_environment_sdn_zone_vlan.lifecycle_zone_2
							]
						}
					}

					resource "proxmox_virtual_environment_sdn_applier" "finalizer" {}
				`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributesSet("proxmox_virtual_environment_sdn_applier.lifecycle_applier", []string{
						"id",
					}),
					test.ResourceAttributes("proxmox_virtual_environment_sdn_zone_simple.lifecycle_zone_1", map[string]string{
						"id":  "lifeZ1",
						"mtu": "1450",
					}),
					test.ResourceAttributes("proxmox_virtual_environment_sdn_zone_vlan.lifecycle_zone_2", map[string]string{
						"id":     "lifeZ2",
						"mtu":    "1500",
						"bridge": "vmbr0",
					}),
				),
			},
			{
				Config: te.RenderConfig(`
					# Remove first zone, keep second, add third
					resource "proxmox_virtual_environment_sdn_zone_vlan" "lifecycle_zone_2" {
						id     = "lifeZ2"
						nodes  = ["{{.NodeName}}"]
						mtu    = 1400  # Changed MTU
						bridge = "vmbr0"
						depends_on = [
							proxmox_virtual_environment_sdn_applier.finalizer
						]
					}

					resource "proxmox_virtual_environment_sdn_zone_simple" "lifecycle_zone_3" {
						id    = "lifeZ3"
						nodes = ["{{.NodeName}}"]
						mtu   = 1500
						depends_on = [
							proxmox_virtual_environment_sdn_applier.finalizer
						]
					}

					resource "proxmox_virtual_environment_sdn_applier" "lifecycle_applier" {
						lifecycle {
							replace_triggered_by = [
								proxmox_virtual_environment_sdn_zone_vlan.lifecycle_zone_2,
								proxmox_virtual_environment_sdn_zone_simple.lifecycle_zone_3
							]
						}
					}

					resource "proxmox_virtual_environment_sdn_applier" "finalizer" {}
				`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributesSet("proxmox_virtual_environment_sdn_applier.lifecycle_applier", []string{
						"id",
					}),
					test.ResourceAttributes("proxmox_virtual_environment_sdn_zone_vlan.lifecycle_zone_2", map[string]string{
						"id":     "lifeZ2",
						"mtu":    "1400",
						"bridge": "vmbr0",
					}),
					test.ResourceAttributes("proxmox_virtual_environment_sdn_zone_simple.lifecycle_zone_3", map[string]string{
						"id":  "lifeZ3",
						"mtu": "1500",
					}),
				),
			},
		}},
		{"applier without dependencies", []resource.TestStep{
			{
				Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_sdn_applier" "standalone_applier" {
					}
				`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributesSet("proxmox_virtual_environment_sdn_applier.standalone_applier", []string{
						"id",
					}),
				),
			},
		}},
		{"applier with mixed zone types", []resource.TestStep{
			{
				Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_sdn_zone_simple" "combo_simple" {
						id    = "comboS"
						nodes = ["{{.NodeName}}"]
						mtu   = 1500
						ipam  = "pve"
						depends_on = [
							proxmox_virtual_environment_sdn_applier.finalizer
						]
					}

					resource "proxmox_virtual_environment_sdn_zone_vlan" "combo_vlan" {
						id     = "comboV"
						nodes  = ["{{.NodeName}}"]
						bridge = "vmbr0"
						mtu    = 1400
						depends_on = [
							proxmox_virtual_environment_sdn_applier.finalizer
						]
					}

					resource "proxmox_virtual_environment_sdn_applier" "combo_applier" {
						lifecycle {
							replace_triggered_by = [
								proxmox_virtual_environment_sdn_zone_simple.combo_simple,
								proxmox_virtual_environment_sdn_zone_vlan.combo_vlan
							]
						}
					}

					resource "proxmox_virtual_environment_sdn_applier" "finalizer" {}
				`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributesSet("proxmox_virtual_environment_sdn_applier.combo_applier", []string{
						"id",
					}),
					test.ResourceAttributes("proxmox_virtual_environment_sdn_zone_simple.combo_simple", map[string]string{
						"id":   "comboS",
						"mtu":  "1500",
						"ipam": "pve",
					}),
					test.ResourceAttributes("proxmox_virtual_environment_sdn_zone_vlan.combo_vlan", map[string]string{
						"id":     "comboV",
						"bridge": "vmbr0",
						"mtu":    "1400",
					}),
				),
			},
			{
				Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_sdn_zone_simple" "combo_simple" {
						id    = "comboS"
						nodes = ["{{.NodeName}}"]
						mtu   = 1450  # Update triggers applier replacement
						ipam  = "pve"
						depends_on = [
							proxmox_virtual_environment_sdn_applier.finalizer
						]
					}

					resource "proxmox_virtual_environment_sdn_zone_vlan" "combo_vlan" {
						id     = "comboV"
						nodes  = ["{{.NodeName}}"]
						bridge = "vmbr0"
						mtu    = 1350  # Update triggers applier replacement
						depends_on = [
							proxmox_virtual_environment_sdn_applier.finalizer
						]
					}

					resource "proxmox_virtual_environment_sdn_applier" "combo_applier" {
						lifecycle {
							replace_triggered_by = [
								proxmox_virtual_environment_sdn_zone_simple.combo_simple,
								proxmox_virtual_environment_sdn_zone_vlan.combo_vlan
							]
						}
					}

					resource "proxmox_virtual_environment_sdn_applier" "finalizer" {}
				`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributesSet("proxmox_virtual_environment_sdn_applier.combo_applier", []string{
						"id",
					}),
					test.ResourceAttributes("proxmox_virtual_environment_sdn_zone_simple.combo_simple", map[string]string{
						"id":   "comboS",
						"mtu":  "1450",
						"ipam": "pve",
					}),
					test.ResourceAttributes("proxmox_virtual_environment_sdn_zone_vlan.combo_vlan", map[string]string{
						"id":     "comboV",
						"bridge": "vmbr0",
						"mtu":    "1350",
					}),
				),
			},
		}},
		{"applier multiple instances", []resource.TestStep{
			{
				Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_sdn_zone_simple" "multi_inst_zone" {
						id    = "miZ"
						nodes = ["{{.NodeName}}"]
						mtu   = 1500
						depends_on = [
							proxmox_virtual_environment_sdn_applier.finalizer
						]
					}

					resource "proxmox_virtual_environment_sdn_applier" "multi_inst_applier_1" {
						depends_on = [
							proxmox_virtual_environment_sdn_zone_simple.multi_inst_zone
						]
					}

					resource "proxmox_virtual_environment_sdn_applier" "multi_inst_applier_2" {
						depends_on = [
							proxmox_virtual_environment_sdn_applier.multi_inst_applier_1
						]
					}

					resource "proxmox_virtual_environment_sdn_applier" "finalizer" {}
				`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributesSet("proxmox_virtual_environment_sdn_applier.multi_inst_applier_1", []string{
						"id",
					}),
					test.ResourceAttributesSet("proxmox_virtual_environment_sdn_applier.multi_inst_applier_2", []string{
						"id",
					}),
					test.ResourceAttributes("proxmox_virtual_environment_sdn_zone_simple.multi_inst_zone", map[string]string{
						"id":  "miZ",
						"mtu": "1500",
					}),
				),
			},
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource.ParallelTest(t, resource.TestCase{
				ProtoV6ProviderFactories: te.AccProviders,
				Steps:                    tt.steps,
			})
		})
	}
}
