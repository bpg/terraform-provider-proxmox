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
	t.Parallel()

	te := test.InitEnvironment(t)

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"create and update zones", []resource.TestStep{{
			Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_sdn_zone_simple" "zone_simple" {
				  id  = "zoneS"
				  nodes = ["pve"]
				  mtu   = 1496
				}
			`),
		}, {
			Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_sdn_zone_simple" "zone_simple" {
				  id  = "zoneS"
				  nodes = ["pve"]
				  mtu   = 1495
				}
			`),
			ResourceName:      "proxmox_virtual_environment_sdn_zone_simple.zone_simple",
			ImportStateId:     "zoneS",
			ImportState:       true,
			ImportStateVerify: true,
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

func TestAccResourceSDNZoneVLAN(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"create and update VLAN zone", []resource.TestStep{{
			Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_sdn_zone_vlan" "zone_vlan" {
				  id    = "zoneV"
				  nodes = ["pve"]
				  mtu   = 1496
				  bridge = "vmbr0"
				}
			`),
		}, {
			Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_sdn_zone_vlan" "zone_vlan" {
				  id    = "zoneV"
				  nodes = ["pve"]
				  mtu   = 1495
				  bridge = "vmbr0"
				}
			`),
			ResourceName:      "proxmox_virtual_environment_sdn_zone_vlan.zone_vlan",
			ImportStateId:     "zoneV",
			ImportState:       true,
			ImportStateVerify: true,
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

func TestAccResourceSDNZoneQinQ(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"create and update QinQ zone", []resource.TestStep{{
			Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_sdn_zone_qinq" "zone_qinq" {
				  id    = "zoneQ"
				  nodes = ["pve"]
				  mtu   = 1496
				  bridge = "vmbr0"
				  service_vlan = 100
				  service_vlan_protocol = "802.1ad"
				}
			`),
		}, {
			Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_sdn_zone_qinq" "zone_qinq" {
				  id    = "zoneQ"
				  nodes = ["pve"]
				  mtu   = 1495
				  bridge = "vmbr0"
				  service_vlan = 200
				  service_vlan_protocol = "802.1q"
				}
			`),
			ResourceName:      "proxmox_virtual_environment_sdn_zone_qinq.zone_qinq",
			ImportStateId:     "zoneQ",
			ImportState:       true,
			ImportStateVerify: true,
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

func TestAccResourceSDNZoneVXLAN(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"create and update VXLAN zone", []resource.TestStep{{
			Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_sdn_zone_vxlan" "zone_vxlan" {
				  id    = "zoneX"
				  nodes = ["pve"]
				  mtu   = 1450
				  peers = ["10.0.0.1", "10.0.0.2"]
				}
			`),
		}, {
			Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_sdn_zone_vxlan" "zone_vxlan" {
				  id    = "zoneX"
				  nodes = ["pve"]
				  mtu   = 1440
				  peers = ["10.0.0.3", "10.0.0.4"]
				}
			`),
			ResourceName:      "proxmox_virtual_environment_sdn_zone_vxlan.zone_vxlan",
			ImportStateId:     "zoneX",
			ImportState:       true,
			ImportStateVerify: true,
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
