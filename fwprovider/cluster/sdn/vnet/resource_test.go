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

func TestAccResourceSDNVNet(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"create and update vnet with simple zone", []resource.TestStep{{
			Config: te.RenderConfig(`
			resource "proxmox_virtual_environment_sdn_zone_simple" "test_zone" {
				id    = "testz"
				nodes = ["{{.NodeName}}"]
			}

			resource "proxmox_virtual_environment_sdn_vnet" "test_vnet" {
				id            = "testv"
				zone          = proxmox_virtual_environment_sdn_zone_simple.test_zone.id
				alias         = "Test VNet"
				isolate_ports = true
				vlan_aware    = false
			}`),
			Check: resource.ComposeTestCheckFunc(
				test.ResourceAttributes("proxmox_virtual_environment_sdn_vnet.test_vnet", map[string]string{
					"id":            "testv",
					"zone":          "testz",
					"alias":         "Test VNet",
					"isolate_ports": "true",
					"vlan_aware":    "false",
				}),
			),
		}, {
			Config: te.RenderConfig(`
			resource "proxmox_virtual_environment_sdn_zone_simple" "test_zone" {
				id    = "testz"
				nodes = ["{{.NodeName}}"]
			}

			resource "proxmox_virtual_environment_sdn_vnet" "test_vnet" {
				id            = "testv"
				zone          = proxmox_virtual_environment_sdn_zone_simple.test_zone.id
				alias         = "Updated VNet"
				isolate_ports = false
				vlan_aware    = true
			}`),
			Check: resource.ComposeTestCheckFunc(
				test.ResourceAttributes("proxmox_virtual_environment_sdn_vnet.test_vnet", map[string]string{
					"id":            "testv",
					"zone":          "testz",
					"alias":         "Updated VNet",
					"isolate_ports": "false",
					"vlan_aware":    "true",
				}),
			),
			ResourceName:      "proxmox_virtual_environment_sdn_vnet.test_vnet",
			ImportStateId:     "testv",
			ImportState:       true,
			ImportStateVerify: true,
		}}},
		{"create vnet with minimal attributes", []resource.TestStep{{
			Config: te.RenderConfig(`
			resource "proxmox_virtual_environment_sdn_zone_simple" "test_zone" {
				id    = "testz2"
				nodes = ["{{.NodeName}}"]
			}

			resource "proxmox_virtual_environment_sdn_vnet" "test_vnet" {
				id   = "testv2"
				zone = proxmox_virtual_environment_sdn_zone_simple.test_zone.id
				depends_on = [
					proxmox_virtual_environment_sdn_zone_simple.test_zone
				]
			}`),
			Check: resource.ComposeTestCheckFunc(
				test.ResourceAttributes("proxmox_virtual_environment_sdn_vnet.test_vnet", map[string]string{
					"id":   "testv2",
					"zone": "testz2",
				}),
			),
		}}},
		{"create vnet with VLAN zone and tag", []resource.TestStep{{
			Config: te.RenderConfig(`
			resource "proxmox_virtual_environment_sdn_zone_vlan" "test_zone3" {
				id     = "testz3"
				nodes  = ["{{.NodeName}}"]
				bridge = "vmbr0"
			}

			resource "proxmox_virtual_environment_sdn_vnet" "test_vnet3" {
				id            = "testv3"
				zone          = proxmox_virtual_environment_sdn_zone_vlan.test_zone3.id
				alias         = "VNet with Tag"
				isolate_ports = true
				tag           = 300
				vlan_aware    = true
				depends_on = [
					proxmox_virtual_environment_sdn_applier.finalizer
				]
			}

			resource "proxmox_virtual_environment_sdn_applier" "test_applier3" {
				depends_on = [
					proxmox_virtual_environment_sdn_zone_vlan.test_zone3,
					proxmox_virtual_environment_sdn_vnet.test_vnet3
				]
			}

			resource "proxmox_virtual_environment_sdn_applier" "finalizer" {}
			`),
			Check: resource.ComposeTestCheckFunc(
				test.ResourceAttributes("proxmox_virtual_environment_sdn_vnet.test_vnet3", map[string]string{
					"id":            "testv3",
					"zone":          "testz3",
					"alias":         "VNet with Tag",
					"isolate_ports": "true",
					"tag":           "300",
					"vlan_aware":    "true",
				}),
			),
		}, {
			Config: te.RenderConfig(`
			resource "proxmox_virtual_environment_sdn_zone_vlan" "test_zone3" {
				id     = "testz3"
				nodes  = ["{{.NodeName}}"]
				bridge = "vmbr0"
				depends_on = [
					proxmox_virtual_environment_sdn_applier.finalizer
				]
			}

			resource "proxmox_virtual_environment_sdn_vnet" "test_vnet3" {
				id            = "testv3"
				zone          = proxmox_virtual_environment_sdn_zone_vlan.test_zone3.id
				alias         = "VNet with Tag"
				isolate_ports = true
				tag           = 300
				vlan_aware    = true
				depends_on = [
					proxmox_virtual_environment_sdn_applier.finalizer
				]
			}

			resource "proxmox_virtual_environment_sdn_applier" "test_applier3" {
				depends_on = [
					proxmox_virtual_environment_sdn_zone_vlan.test_zone3,
					proxmox_virtual_environment_sdn_vnet.test_vnet3
				]
			}

			resource "proxmox_virtual_environment_sdn_applier" "finalizer" {}
			`),
			Check: resource.ComposeTestCheckFunc(
				test.ResourceAttributes("proxmox_virtual_environment_sdn_vnet.test_vnet3", map[string]string{
					"id":            "testv3",
					"zone":          "testz3",
					"alias":         "VNet with Tag",
					"isolate_ports": "true",
					"tag":           "300",
					"vlan_aware":    "true",
				}),
			),
			ResourceName:      "proxmox_virtual_environment_sdn_vnet.test_vnet3",
			ImportStateId:     "testv3",
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
