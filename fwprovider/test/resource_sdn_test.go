//go:build acceptance || all

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceSDN(t *testing.T) {
	te := InitEnvironment(t)

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"create zones, vnets and subnets", []resource.TestStep{{
			Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_sdn_zone" "zone_simple" {
				  name  = "zoneS"
				  type  = "simple"
				  nodes = "weisshorn-proxmox"
				  mtu   = 1496
				}

				resource "proxmox_virtual_environment_sdn_zone" "zone_vlan" {
				  name   = "zoneVLAN"
				  type   = "vlan"
				  nodes  = "weisshorn-proxmox"
				  mtu    = 1500
				  bridge = "vmbr0"
				}

				resource "proxmox_virtual_environment_sdn_vnet" "vnet_simple" {
				  name           = "vnetM"
				  zone           = proxmox_virtual_environment_sdn_zone.zone_simple.name
				  alias          = "vnet in zoneM"
				  isolate_ports  = "0"
				  vlanaware      = "0"
				  zonetype       = proxmox_virtual_environment_sdn_zone.zone_simple.type
				  depends_on     = [proxmox_virtual_environment_sdn_zone.zone_simple]
				}

				resource "proxmox_virtual_environment_sdn_vnet" "vnet_vlan" {
				  name       = "vnetVLAN"
				  zone       = proxmox_virtual_environment_sdn_zone.zone_vlan.name
				  alias      = "vnet in zoneVLAN"
				  tag        = 1000
				  zonetype   = proxmox_virtual_environment_sdn_zone.zone_vlan.type
				  depends_on = [proxmox_virtual_environment_sdn_zone.zone_vlan]
				}

				resource "proxmox_virtual_environment_sdn_subnet" "subnet_simple" {
				  subnet            = "10.10.0.0/24"
				  vnet              = proxmox_virtual_environment_sdn_vnet.vnet_simple.name
				  dhcp_dns_server   = "10.10.0.53"
				  dhcp_range = [
				    {
				      start_address = "10.10.0.10"
				      end_address   = "10.10.0.100"
				    }
				  ]
				  gateway    = "10.10.0.1"
				  snat       = true
				  depends_on = [proxmox_virtual_environment_sdn_vnet.vnet_simple]
				}

				resource "proxmox_virtual_environment_sdn_subnet" "subnet_simple2" {
				  subnet            = "10.40.0.0/24"
				  vnet              = proxmox_virtual_environment_sdn_vnet.vnet_simple.name
				  dhcp_dns_server   = "10.40.0.53"
				  dhcp_range = [
				    {
				      start_address = "10.40.0.10"
				      end_address   = "10.40.0.100"
				    }
				  ]
				  gateway    = "10.40.0.1"
				  snat       = true
				  depends_on = [proxmox_virtual_environment_sdn_vnet.vnet_simple]
				}

				resource "proxmox_virtual_environment_sdn_subnet" "subnet_vlan" {
				  subnet            = "10.20.0.0/24"
				  vnet              = proxmox_virtual_environment_sdn_vnet.vnet_vlan.name
				  dhcp_dns_server   = "10.20.0.53"
				  dhcp_range = [
				    {
				      start_address = "10.20.0.10"
				      end_address   = "10.20.0.100"
				    }
				  ]
				  gateway    = "10.20.0.100"
				  snat       = false
				  depends_on = [proxmox_virtual_environment_sdn_vnet.vnet_vlan]
				}
			`),
			Check: resource.ComposeTestCheckFunc(
				// Zones
				ResourceAttributes("proxmox_virtual_environment_sdn_zone.zone_simple", map[string]string{
					"name":  "zoneS",
					"type":  "simple",
					"mtu":   "1496",
					"nodes": "weisshorn-proxmox",
				}),
				ResourceAttributes("proxmox_virtual_environment_sdn_zone.zone_vlan", map[string]string{
					"name":   "zoneVLAN",
					"type":   "vlan",
					"mtu":    "1500",
					"bridge": "vmbr0",
				}),

				// VNets
				ResourceAttributes("proxmox_virtual_environment_sdn_vnet.vnet_simple", map[string]string{
					"name":          "vnetM",
					"alias":         "vnet in zoneM",
					"zone":          "zoneS",
					"isolate_ports": "false",
					"vlanaware":     "false",
					"zonetype":      "simple",
				}),
				ResourceAttributes("proxmox_virtual_environment_sdn_vnet.vnet_vlan", map[string]string{
					"name":     "vnetVLAN",
					"alias":    "vnet in zoneVLAN",
					"zone":     "zoneVLAN",
					"tag":      "1000",
					"zonetype": "vlan",
				}),

				// Subnet (only check one in detail to avoid too many long checks)
				ResourceAttributes("proxmox_virtual_environment_sdn_subnet.subnet_simple", map[string]string{
					"subnet":          "10.10.0.0/24",
					"vnet":            "vnetM",
					"gateway":         "10.10.0.1",
					"dhcp_dns_server": "10.10.0.53",
					"snat":            "true",
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
