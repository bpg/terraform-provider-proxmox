//go:build acceptance || all

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package vnet_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/test"
)

func TestAccDataSourceSDNVNets(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"create multiple VNets with all attributes and read with vnets datasource", []resource.TestStep{{
			Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_sdn_zone_simple" "test" {
					id    = "vnetstst"
					nodes = ["{{.NodeName}}"]
				}

				resource "proxmox_virtual_environment_sdn_vnet" "test1" {
					id            = "dsvnet1"
					zone          = proxmox_virtual_environment_sdn_zone_simple.test.id
					alias         = "Test VNet 1"
					isolate_ports = true
					vlan_aware    = false
				}

				resource "proxmox_virtual_environment_sdn_vnet" "test2" {
					id            = "dsvnet2"
					zone          = proxmox_virtual_environment_sdn_zone_simple.test.id
					alias         = "Test VNet 2"
					isolate_ports = false
					vlan_aware    = true
				}

				resource "proxmox_virtual_environment_sdn_applier" "test_applier" {
					depends_on = [
						proxmox_virtual_environment_sdn_vnet.test1,
						proxmox_virtual_environment_sdn_vnet.test2
					]
				}

				data "proxmox_virtual_environment_sdn_vnets" "all" {
					depends_on = [
						proxmox_virtual_environment_sdn_applier.test_applier
					]
				}
			`),
			Check: resource.ComposeTestCheckFunc(
				// Check that the vnets datasource returns multiple vnets
				resource.TestCheckResourceAttrSet("data.proxmox_virtual_environment_sdn_vnets.all", "vnets.#"),
				// Verify that vnets contain expected attributes
				resource.TestCheckTypeSetElemNestedAttrs("data.proxmox_virtual_environment_sdn_vnets.all", "vnets.*", map[string]string{
					"id":            "dsvnet1",
					"zone":          "vnetstst",
					"alias":         "Test VNet 1",
					"isolate_ports": "true",
					"vlan_aware":    "false",
				}),
				resource.TestCheckTypeSetElemNestedAttrs("data.proxmox_virtual_environment_sdn_vnets.all", "vnets.*", map[string]string{
					"id":            "dsvnet2",
					"zone":          "vnetstst",
					"alias":         "Test VNet 2",
					"isolate_ports": "false",
					"vlan_aware":    "true",
				}),
			),
		}}},
		{"create VNets with VLAN zone and tags", []resource.TestStep{{
			Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_sdn_zone_vlan" "test_vlan" {
					id     = "vnetvlan"
					nodes  = ["{{.NodeName}}"]
					bridge = "vmbr0"
				}

				resource "proxmox_virtual_environment_sdn_vnet" "vlan1" {
					id         = "vlan1"
					zone       = proxmox_virtual_environment_sdn_zone_vlan.test_vlan.id
					alias      = "VLAN VNet 1"
					tag        = 100
					vlan_aware = true
				}

				resource "proxmox_virtual_environment_sdn_vnet" "vlan2" {
					id         = "vlan2"
					zone       = proxmox_virtual_environment_sdn_zone_vlan.test_vlan.id
					alias      = "VLAN VNet 2"
					tag        = 200
					vlan_aware = false
				}

				resource "proxmox_virtual_environment_sdn_applier" "test_applier" {
					depends_on = [
						proxmox_virtual_environment_sdn_vnet.vlan1,
						proxmox_virtual_environment_sdn_vnet.vlan2
					]
				}

				data "proxmox_virtual_environment_sdn_vnets" "all" {
					depends_on = [
						proxmox_virtual_environment_sdn_applier.test_applier
					]
				}
			`),
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet("data.proxmox_virtual_environment_sdn_vnets.all", "vnets.#"),
				resource.TestCheckTypeSetElemNestedAttrs("data.proxmox_virtual_environment_sdn_vnets.all", "vnets.*", map[string]string{
					"id":         "vlan1",
					"zone":       "vnetvlan",
					"alias":      "VLAN VNet 1",
					"tag":        "100",
					"vlan_aware": "true",
				}),
				resource.TestCheckTypeSetElemNestedAttrs("data.proxmox_virtual_environment_sdn_vnets.all", "vnets.*", map[string]string{
					"id":         "vlan2",
					"zone":       "vnetvlan",
					"alias":      "VLAN VNet 2",
					"tag":        "200",
					"vlan_aware": "false",
				}),
			),
		}}},
		{"create VNets with minimal attributes", []resource.TestStep{{
			Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_sdn_zone_simple" "test_minimal" {
					id    = "vnetmin"
					nodes = ["{{.NodeName}}"]
				}

				resource "proxmox_virtual_environment_sdn_vnet" "minimal1" {
					id   = "minvnet1"
					zone = proxmox_virtual_environment_sdn_zone_simple.test_minimal.id
				}

				resource "proxmox_virtual_environment_sdn_vnet" "minimal2" {
					id   = "minvnet2"
					zone = proxmox_virtual_environment_sdn_zone_simple.test_minimal.id
				}

				resource "proxmox_virtual_environment_sdn_applier" "test_applier" {
					depends_on = [
						proxmox_virtual_environment_sdn_vnet.minimal1,
						proxmox_virtual_environment_sdn_vnet.minimal2
					]
				}

				data "proxmox_virtual_environment_sdn_vnets" "all" {
					depends_on = [
						proxmox_virtual_environment_sdn_applier.test_applier
					]
				}
			`),
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet("data.proxmox_virtual_environment_sdn_vnets.all", "vnets.#"),
				resource.TestCheckTypeSetElemNestedAttrs("data.proxmox_virtual_environment_sdn_vnets.all", "vnets.*", map[string]string{
					"id":   "minvnet1",
					"zone": "vnetmin",
				}),
				resource.TestCheckTypeSetElemNestedAttrs("data.proxmox_virtual_environment_sdn_vnets.all", "vnets.*", map[string]string{
					"id":   "minvnet2",
					"zone": "vnetmin",
				}),
			),
		}}},
		{"read vnets datasource with mixed zone types", []resource.TestStep{{
			Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_sdn_zone_simple" "simple" {
					id    = "mixsimpl"
					nodes = ["{{.NodeName}}"]
				}

				resource "proxmox_virtual_environment_sdn_zone_vlan" "vlan" {
					id     = "mixvlan"
					nodes  = ["{{.NodeName}}"]
					bridge = "vmbr0"
				}

				resource "proxmox_virtual_environment_sdn_vnet" "simple_vnet" {
					id         = "mixvnet1"
					zone       = proxmox_virtual_environment_sdn_zone_simple.simple.id
					alias      = "Simple Zone VNet"
					vlan_aware = false
				}

				resource "proxmox_virtual_environment_sdn_vnet" "vlan_vnet" {
					id         = "mixvnet2"
					zone       = proxmox_virtual_environment_sdn_zone_vlan.vlan.id
					alias      = "VLAN Zone VNet"
					tag        = 300
					vlan_aware = true
				}

				resource "proxmox_virtual_environment_sdn_applier" "test_applier" {
					depends_on = [
						proxmox_virtual_environment_sdn_vnet.simple_vnet,
						proxmox_virtual_environment_sdn_vnet.vlan_vnet
					]
				}

				data "proxmox_virtual_environment_sdn_vnets" "all" {
					depends_on = [
						proxmox_virtual_environment_sdn_applier.test_applier
					]
				}
			`),
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet("data.proxmox_virtual_environment_sdn_vnets.all", "vnets.#"),
				resource.TestCheckTypeSetElemNestedAttrs("data.proxmox_virtual_environment_sdn_vnets.all", "vnets.*", map[string]string{
					"id":         "mixvnet1",
					"zone":       "mixsimpl",
					"alias":      "Simple Zone VNet",
					"vlan_aware": "false",
				}),
				resource.TestCheckTypeSetElemNestedAttrs("data.proxmox_virtual_environment_sdn_vnets.all", "vnets.*", map[string]string{
					"id":         "mixvnet2",
					"zone":       "mixvlan",
					"alias":      "VLAN Zone VNet",
					"tag":        "300",
					"vlan_aware": "true",
				}),
			),
		}}},
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
