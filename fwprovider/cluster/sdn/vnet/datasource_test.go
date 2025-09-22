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

func TestAccDataSourceSDNVNet(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"read vnet data source with all attributes", []resource.TestStep{{
			Config: te.RenderConfig(`
			resource "proxmox_virtual_environment_sdn_zone_simple" "test_zone" {
				id    = "testdz"
				nodes = ["{{.NodeName}}"]
			}

			resource "proxmox_virtual_environment_sdn_vnet" "test_vnet" {
				id            = "testdv"
				zone          = proxmox_virtual_environment_sdn_zone_simple.test_zone.id
				alias         = "Test Data VNet"
				isolate_ports = true
				vlan_aware    = false
			}

			data "proxmox_virtual_environment_sdn_vnet" "test_vnet_data" {
				id = proxmox_virtual_environment_sdn_vnet.test_vnet.id
			}`),
			Check: resource.ComposeTestCheckFunc(
				test.ResourceAttributes("data.proxmox_virtual_environment_sdn_vnet.test_vnet_data", map[string]string{
					"id":            "testdv",
					"zone":          "testdz",
					"alias":         "Test Data VNet",
					"isolate_ports": "true",
					"vlan_aware":    "false",
				}),
			),
		}}},
		{"read vnet data source with minimal attributes", []resource.TestStep{{
			Config: te.RenderConfig(`
			resource "proxmox_virtual_environment_sdn_zone_simple" "test_zone" {
				id    = "testdz2"
				nodes = ["{{.NodeName}}"]
			}

			resource "proxmox_virtual_environment_sdn_vnet" "test_vnet" {
				id   = "testdv2"
				zone = proxmox_virtual_environment_sdn_zone_simple.test_zone.id
			}

			data "proxmox_virtual_environment_sdn_vnet" "test_vnet_data" {
				id = proxmox_virtual_environment_sdn_vnet.test_vnet.id
			}`),
			Check: resource.ComposeTestCheckFunc(
				test.ResourceAttributes("data.proxmox_virtual_environment_sdn_vnet.test_vnet_data", map[string]string{
					"id":   "testdv2",
					"zone": "testdz2",
				}),
			),
		}}},
		{"read vnet data source with VLAN zone and tag", []resource.TestStep{{
			Config: te.RenderConfig(`
			resource "proxmox_virtual_environment_sdn_zone_vlan" "test_zone3" {
				id     = "testdz3"
				nodes  = ["{{.NodeName}}"]
				bridge = "vmbr0"
			}

			resource "proxmox_virtual_environment_sdn_vnet" "test_vnet3" {
				id            = "testdv3"
				zone          = proxmox_virtual_environment_sdn_zone_vlan.test_zone3.id
				alias         = "Data VNet with Tag"
				isolate_ports = true
				tag           = 400
				vlan_aware    = true
				depends_on = [
					proxmox_virtual_environment_sdn_applier.finalizer
				]
			}

			data "proxmox_virtual_environment_sdn_vnet" "test_vnet_data3" {
				id = proxmox_virtual_environment_sdn_vnet.test_vnet3.id
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
				test.ResourceAttributes("data.proxmox_virtual_environment_sdn_vnet.test_vnet_data3", map[string]string{
					"id":            "testdv3",
					"zone":          "testdz3",
					"alias":         "Data VNet with Tag",
					"isolate_ports": "true",
					"tag":           "400",
					"vlan_aware":    "true",
				}),
			),
		}}},
		{"data source reads vnet with pending changes scenario", []resource.TestStep{{
			Config: te.RenderConfig(`
			resource "proxmox_virtual_environment_sdn_zone_simple" "test_zone_ds_pending" {
				id    = "testdzp"
				nodes = ["{{.NodeName}}"]
			}

			resource "proxmox_virtual_environment_sdn_vnet" "test_vnet_ds_pending" {
				id            = "testdvp"
				zone          = proxmox_virtual_environment_sdn_zone_simple.test_zone_ds_pending.id
				alias         = "Data Source Pending VNet"
				isolate_ports = false
				vlan_aware    = false
			}

			data "proxmox_virtual_environment_sdn_vnet" "test_vnet_ds_pending" {
				id = proxmox_virtual_environment_sdn_vnet.test_vnet_ds_pending.id
			}`),
			Check: resource.ComposeTestCheckFunc(
				test.ResourceAttributes("data.proxmox_virtual_environment_sdn_vnet.test_vnet_ds_pending", map[string]string{
					"id":            "testdvp",
					"zone":          "testdzp",
					"alias":         "Data Source Pending VNet",
					"isolate_ports": "false",
					"vlan_aware":    "false",
				}),
			),
		}, {
			Config: te.RenderConfig(`
			resource "proxmox_virtual_environment_sdn_zone_simple" "test_zone_ds_pending" {
				id    = "testdzp"
				nodes = ["{{.NodeName}}"]
			}

			resource "proxmox_virtual_environment_sdn_vnet" "test_vnet_ds_pending" {
				id            = "testdvp"
				zone          = proxmox_virtual_environment_sdn_zone_simple.test_zone_ds_pending.id
				alias         = "Updated Data Source Pending VNet"
				isolate_ports = true
				vlan_aware    = true
			}

			data "proxmox_virtual_environment_sdn_vnet" "test_vnet_ds_pending" {
				id = proxmox_virtual_environment_sdn_vnet.test_vnet_ds_pending.id
			}`),
			Check: resource.ComposeTestCheckFunc(
				test.ResourceAttributes("data.proxmox_virtual_environment_sdn_vnet.test_vnet_ds_pending", map[string]string{
					"id":            "testdvp",
					"zone":          "testdzp",
					"alias":         "Updated Data Source Pending VNet",
					"isolate_ports": "true",
					"vlan_aware":    "true",
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
