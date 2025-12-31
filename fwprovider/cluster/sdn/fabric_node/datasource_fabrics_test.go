//go:build acceptance || all

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package fabric_node_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/test"
)

func TestAccDataSourceSDNFabricNodeOpenFabric(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"create openfabric fabric node and read with datasource", []resource.TestStep{{
			Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_sdn_fabric_openfabric" "test" {
				  id        = "dstest1"
				  ip_prefix = "10.0.0.0/16"
				}
				
				resource "proxmox_virtual_environment_sdn_fabric_node_openfabric" "test_node" {
				  fabric_id = proxmox_virtual_environment_sdn_fabric_openfabric.test.id
				  node_id   = "pve"
				  ip        = "10.0.0.1"
				}

				data "proxmox_virtual_environment_sdn_fabric_node_openfabric" "test_node" {
				  fabric_id = proxmox_virtual_environment_sdn_fabric_openfabric.test.id
				  node_id   = "pve"
				}
			`),
			Check: resource.ComposeTestCheckFunc(
				test.ResourceAttributes("data.proxmox_virtual_environment_sdn_fabric_node_openfabric.test_node", map[string]string{
					"fabric_id": "dstest1",
					"node_id":   "pve",
					"ip":        "10.0.0.1",
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

func TestAccDataSourceSDNFabricOSPF(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"create OSPF fabric node and read with datasource", []resource.TestStep{{
			Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_sdn_fabric_ospf" "test" {
				  id     = "dstest2"
				  area  = "0"
				  ip_prefix = "10.0.0.0/16"
				}
				
				resource "proxmox_virtual_environment_sdn_fabric_node_ospf" "test_node" {
				  fabric_id = proxmox_virtual_environment_sdn_fabric_ospf.test.id
				  node_id   = "pve"
				  ip        = "10.0.0.1"
				}
				
				data "proxmox_virtual_environment_sdn_fabric_node_ospf" "test_node" {
				  fabric_id       = proxmox_virtual_environment_sdn_fabric_ospf.test.id
				  node_id         = "pve"
				}
			`),
			Check: resource.ComposeTestCheckFunc(
				test.ResourceAttributes("data.proxmox_virtual_environment_sdn_fabric_node_ospf.test_node", map[string]string){
					"fabric_id": "dstest2",
					"node_id":   "pve",
					"ip":        "10.0.0.1",
				}
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
