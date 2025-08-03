//go:build acceptance || all

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package zone_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/test"
)

func TestAccDataSourceSDNZoneSimple(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"create simple zone and read with datasource", []resource.TestStep{{
			Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_sdn_zone_simple" "test" {
					id    = "dstest1"
					nodes = ["pve"]
					mtu   = 1500
				}
				
				data "proxmox_virtual_environment_sdn_zone_simple" "test" {
					id = proxmox_virtual_environment_sdn_zone_simple.test.id
				}
			`),
			Check: resource.ComposeTestCheckFunc(
				test.ResourceAttributes("data.proxmox_virtual_environment_sdn_zone_simple.test", map[string]string{
					"id":  "dstest1",
					"mtu": "1500",
				}),
				test.ResourceAttributes("data.proxmox_virtual_environment_sdn_zone_simple.test", map[string]string{
					"nodes.#": "1",
					"nodes.0": "pve",
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

func TestAccDataSourceSDNZoneVLAN(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"create VLAN zone and read with datasource", []resource.TestStep{{
			Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_sdn_zone_vlan" "test" {
					id     = "dstest2"
					nodes  = ["pve"]
					bridge = "vmbr0"
					mtu    = 1496
				}
				
				data "proxmox_virtual_environment_sdn_zone_vlan" "test" {
					id = proxmox_virtual_environment_sdn_zone_vlan.test.id
				}
			`),
			Check: resource.ComposeTestCheckFunc(
				test.ResourceAttributes("data.proxmox_virtual_environment_sdn_zone_vlan.test", map[string]string{
					"id":      "dstest2",
					"bridge":  "vmbr0",
					"mtu":     "1496",
					"nodes.#": "1",
					"nodes.0": "pve",
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

func TestAccDataSourceSDNZoneQinQ(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"create QinQ zone and read with datasource", []resource.TestStep{{
			Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_sdn_zone_qinq" "test" {
					id                    = "dstest3"
					nodes                 = ["pve"]
					bridge                = "vmbr0"
					service_vlan          = 100
					service_vlan_protocol = "802.1ad"
					mtu                   = 1492
				}
				
				data "proxmox_virtual_environment_sdn_zone_qinq" "test" {
					id = proxmox_virtual_environment_sdn_zone_qinq.test.id
				}
			`),
			Check: resource.ComposeTestCheckFunc(
				test.ResourceAttributes("data.proxmox_virtual_environment_sdn_zone_qinq.test", map[string]string{
					"id":                    "dstest3",
					"bridge":                "vmbr0",
					"service_vlan":          "100",
					"service_vlan_protocol": "802.1ad",
					"mtu":                   "1492",
					"nodes.#":               "1",
					"nodes.0":               "pve",
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

func TestAccDataSourceSDNZoneVXLAN(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"create VXLAN zone and read with datasource", []resource.TestStep{{
			Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_sdn_zone_vxlan" "test" {
					id    = "dstest4"
					nodes = ["pve"]
					peers = ["10.0.0.1", "10.0.0.2"]
					mtu   = 1450
				}
				
				data "proxmox_virtual_environment_sdn_zone_vxlan" "test" {
					id = proxmox_virtual_environment_sdn_zone_vxlan.test.id
				}
			`),
			Check: resource.ComposeTestCheckFunc(
				test.ResourceAttributes("data.proxmox_virtual_environment_sdn_zone_vxlan.test", map[string]string{
					"id":      "dstest4",
					"mtu":     "1450",
					"nodes.#": "1",
					"nodes.0": "pve",
					"peers.#": "2",
					"peers.0": "10.0.0.1",
					"peers.1": "10.0.0.2",
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

func TestAccDataSourceSDNZones(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"create multiple zones and read with zones datasource", []resource.TestStep{{
			Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_sdn_zone_simple" "test1" {
					id    = "dstest6"
					nodes = ["pve"]
					mtu   = 1500
				}
				
				resource "proxmox_virtual_environment_sdn_zone_vlan" "test2" {
					id     = "dstest7"
					nodes  = ["pve"]
					bridge = "vmbr0"
					mtu    = 1496
				}
				
				data "proxmox_virtual_environment_sdn_zones" "all" {
					depends_on = [
						proxmox_virtual_environment_sdn_zone_simple.test1,
						proxmox_virtual_environment_sdn_zone_vlan.test2
					]
				}
				
				data "proxmox_virtual_environment_sdn_zones" "simple_only" {
					type = "simple"
					depends_on = [
						proxmox_virtual_environment_sdn_zone_simple.test1,
						proxmox_virtual_environment_sdn_zone_vlan.test2
					]
				}
			`),
			Check: resource.ComposeTestCheckFunc(
				// Check that all zones datasource returns multiple zones
				resource.TestCheckResourceAttrSet("data.proxmox_virtual_environment_sdn_zones.all", "zones.#"),

				// Check that filtered datasource works
				resource.TestCheckResourceAttr("data.proxmox_virtual_environment_sdn_zones.simple_only", "type", "simple"),
				resource.TestCheckResourceAttrSet("data.proxmox_virtual_environment_sdn_zones.simple_only", "zones.#"),
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
