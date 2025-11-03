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

func TestAccResourceSDNZoneSimple(t *testing.T) {
	// Cannot run in parallel due to SDN applier functionality affecting global state

	te := test.InitEnvironment(t)

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"create and update zones", []resource.TestStep{{
			Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_sdn_zone_simple" "zone_simple" {
				  id  = "zoneS"
				  nodes = ["{{.NodeName}}"]
				  ipam  = "pve"
				  mtu   = 1496
				  depends_on = [
					proxmox_virtual_environment_sdn_applier.finalizer
				  ]
				}

				resource "proxmox_virtual_environment_sdn_applier" "zone_simple_applier" {
				  depends_on = [
					proxmox_virtual_environment_sdn_zone_simple.zone_simple
				  ]
				}

				resource "proxmox_virtual_environment_sdn_applier" "finalizer" {}
			`),
			Check: resource.ComposeTestCheckFunc(
				test.ResourceAttributes("proxmox_virtual_environment_sdn_zone_simple.zone_simple", map[string]string{
					"id":      "zoneS",
					"mtu":     "1496",
					"ipam":    "pve",
					"pending": "true",
					"state":   "new",
				}),
			),
		}, {
			Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_sdn_zone_simple" "zone_simple" {
				  id  = "zoneS"
				  nodes = ["{{.NodeName}}"]
				  mtu   = 1495
				  depends_on = [
					proxmox_virtual_environment_sdn_applier.finalizer
				  ]
				}

				resource "proxmox_virtual_environment_sdn_applier" "zone_simple_applier" {
				  depends_on = [
					proxmox_virtual_environment_sdn_zone_simple.zone_simple
				  ]
				}

				resource "proxmox_virtual_environment_sdn_applier" "finalizer" {}
			`),
			Check: resource.ComposeTestCheckFunc(
				test.ResourceAttributes("proxmox_virtual_environment_sdn_zone_simple.zone_simple", map[string]string{
					"id":      "zoneS",
					"mtu":     "1495",
					"pending": "true",
					"state":   "changed",
				}),
				test.NoResourceAttributesSet("proxmox_virtual_environment_sdn_zone_simple.zone_simple", []string{
					"ipam",
				}),
			),
			ResourceName: "proxmox_virtual_environment_sdn_zone_simple.zone_simple",
			// ImportStateId:     "zoneS",
			// ImportState:       true,
			// ImportStateVerify: true,
		}}},
		{"create zones with empty nodes", []resource.TestStep{{
			Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_sdn_zone_simple" "zone_simple2" {
				  id  = "zoneSE"
				  nodes = []
				  mtu   = 1496
				  depends_on = [
					proxmox_virtual_environment_sdn_applier.finalizer
				  ]
				}

				resource "proxmox_virtual_environment_sdn_applier" "zone_simple2_applier" {
				  depends_on = [
					proxmox_virtual_environment_sdn_zone_simple.zone_simple2
				  ]
				}
				  
				resource "proxmox_virtual_environment_sdn_applier" "finalizer" {}

			`),
			Check: resource.ComposeTestCheckFunc(
				test.ResourceAttributes("proxmox_virtual_environment_sdn_zone_simple.zone_simple2", map[string]string{
					"id":      "zoneSE",
					"mtu":     "1496",
					"pending": "true",
					"nodes.#": "0",
					"state":   "new",
				}),
			),
		}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: te.AccProviders,
				Steps:                    tt.steps,
			})
		})
	}
}

func TestAccResourceSDNZoneVLAN(t *testing.T) {
	// Cannot run in parallel due to SDN applier functionality affecting global state

	te := test.InitEnvironment(t)

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"create and update VLAN zone", []resource.TestStep{{
			Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_sdn_zone_vlan" "zone_vlan" {
				  id    = "zoneV"
				  nodes = ["{{.NodeName}}"]
				  mtu   = 1496
				  bridge = "vmbr0"
				}
			`),
			Check: resource.ComposeTestCheckFunc(
				test.ResourceAttributes("proxmox_virtual_environment_sdn_zone_vlan.zone_vlan", map[string]string{
					"id":      "zoneV",
					"mtu":     "1496",
					"bridge":  "vmbr0",
					"pending": "true",
					"state":   "new",
				}),
			),
		}, {
			Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_sdn_zone_vlan" "zone_vlan" {
				  id    = "zoneV"
				  nodes = ["{{.NodeName}}"]
				  mtu   = 1495
				  bridge = "vmbr0"
				}
			`),
			Check: resource.ComposeTestCheckFunc(
				test.ResourceAttributes("proxmox_virtual_environment_sdn_zone_vlan.zone_vlan", map[string]string{
					"id":      "zoneV",
					"mtu":     "1495",
					"bridge":  "vmbr0",
					"pending": "true",
					"state":   "new",
				}),
			),
			ResourceName:      "proxmox_virtual_environment_sdn_zone_vlan.zone_vlan",
			ImportStateId:     "zoneV",
			ImportState:       true,
			ImportStateVerify: true,
		}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: te.AccProviders,
				Steps:                    tt.steps,
			})
		})
	}
}

func TestAccResourceSDNZoneVLAN_NoNodes(t *testing.T) {
	// Cannot run in parallel due to SDN applier functionality affecting global state

	te := test.InitEnvironment(t)

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{
			name: "create VLAN zone without nodes",
			steps: []resource.TestStep{{
				Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_sdn_zone_vlan" "zone_vlan_no_nodes" {
					  id     = "zoneVNo"
					  bridge = "vmbr0"
					  mtu    = 1496
					}
				`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_virtual_environment_sdn_zone_vlan.zone_vlan_no_nodes", map[string]string{
						"id":      "zoneVNo",
						"bridge":  "vmbr0",
						"mtu":     "1496",
						"pending": "true",
						"state":   "new",
					}),
				),
			}},
		},
		{
			name: "create VLAN zone with empty nodes list",
			steps: []resource.TestStep{{
				Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_sdn_zone_vlan" "zone_vlan_empty_nodes" {
					  id     = "zoneVEm"
					  nodes  = []
					  bridge = "vmbr0"
					  mtu    = 1496
					}
				`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_virtual_environment_sdn_zone_vlan.zone_vlan_empty_nodes", map[string]string{
						"id":      "zoneVEm",
						"bridge":  "vmbr0",
						"mtu":     "1496",
						"pending": "true",
						"state":   "new",
					}),
				),
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: te.AccProviders,
				Steps:                    tt.steps,
			})
		})
	}
}

func TestAccResourceSDNZoneQinQ(t *testing.T) {
	// Cannot run in parallel due to SDN applier functionality affecting global state

	te := test.InitEnvironment(t)

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"create and update QinQ zone", []resource.TestStep{{
			Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_sdn_zone_qinq" "zone_qinq" {
				  id    = "zoneQ"
				  nodes = ["{{.NodeName}}"]
				  mtu   = 1496
				  bridge = "vmbr0"
				  service_vlan = 100
				  service_vlan_protocol = "802.1ad"
				}
			`),
			Check: resource.ComposeTestCheckFunc(
				test.ResourceAttributes("proxmox_virtual_environment_sdn_zone_qinq.zone_qinq", map[string]string{
					"id":                    "zoneQ",
					"mtu":                   "1496",
					"bridge":                "vmbr0",
					"service_vlan":          "100",
					"service_vlan_protocol": "802.1ad",
					"pending":               "true",
					"state":                 "new",
				}),
			),
		}, {
			Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_sdn_zone_qinq" "zone_qinq" {
				  id    = "zoneQ"
				  nodes = ["{{.NodeName}}"]
				  mtu   = 1495
				  bridge = "vmbr0"
				  service_vlan = 200
				  service_vlan_protocol = "802.1q"
				}
			`),
			Check: resource.ComposeTestCheckFunc(
				test.ResourceAttributes("proxmox_virtual_environment_sdn_zone_qinq.zone_qinq", map[string]string{
					"id":                    "zoneQ",
					"mtu":                   "1495",
					"bridge":                "vmbr0",
					"service_vlan":          "200",
					"service_vlan_protocol": "802.1q",
					"pending":               "true",
					"state":                 "new",
				}),
			),
			ResourceName:      "proxmox_virtual_environment_sdn_zone_qinq.zone_qinq",
			ImportStateId:     "zoneQ",
			ImportState:       true,
			ImportStateVerify: true,
		}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: te.AccProviders,
				Steps:                    tt.steps,
			})
		})
	}
}

func TestAccResourceSDNZoneVXLAN(t *testing.T) {
	// Cannot run in parallel due to SDN applier functionality affecting global state

	te := test.InitEnvironment(t)

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"create and update VXLAN zone", []resource.TestStep{{
			Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_sdn_zone_vxlan" "zone_vxlan" {
				  id    = "zoneX"
				  nodes = ["{{.NodeName}}"]
				  mtu   = 1450
				  peers = ["10.0.0.1", "10.0.0.2"]
				}
			`),
			Check: resource.ComposeTestCheckFunc(
				test.ResourceAttributes("proxmox_virtual_environment_sdn_zone_vxlan.zone_vxlan", map[string]string{
					"id":      "zoneX",
					"mtu":     "1450",
					"pending": "true",
					"state":   "new",
				}),
			),
		}, {
			Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_sdn_zone_vxlan" "zone_vxlan" {
				  id    = "zoneX"
				  nodes = ["{{.NodeName}}"]
				  mtu   = 1440
				  peers = ["10.0.0.3", "10.0.0.4"]
				}
			`),
			Check: resource.ComposeTestCheckFunc(
				test.ResourceAttributes("proxmox_virtual_environment_sdn_zone_vxlan.zone_vxlan", map[string]string{
					"id":      "zoneX",
					"mtu":     "1440",
					"pending": "true",
					"state":   "new",
				}),
			),
			ResourceName:      "proxmox_virtual_environment_sdn_zone_vxlan.zone_vxlan",
			ImportStateId:     "zoneX",
			ImportState:       true,
			ImportStateVerify: true,
		}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: te.AccProviders,
				Steps:                    tt.steps,
			})
		})
	}
}
