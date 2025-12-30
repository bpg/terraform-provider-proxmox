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

func TestAccResourceSDNFabricNodeOpenFabric(t *testing.T) {
	// Cannot run in parallel due to SDN applier functionality affecting global state

	te := test.InitEnvironment(t)

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"create and update fabric nodes", []resource.TestStep{{
			Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_sdn_fabric_openfabric" "fabric_openfabric" {
				  id  = "fabricS"
				  ip_prefix = "10.0.0.0/16"
				  depends_on = [
					proxmox_virtual_environment_sdn_applier.finalizer
				  ]
				}
				
				resource "proxmox_virtual_environment_sdn_fabric_node_openfabric" "fabric_node_openfabric" {
				  fabric_id = proxmox_virtual_environment_sdn_fabric_openfabric.fabric_openfabric.id
				  node_id   = "pve"
				  ip 	    = "10.0.0.1"
				  depends_on = [
					proxmox_virtual_environment_sdn_applier.finalizer
				  ]
				}

				resource "proxmox_virtual_environment_sdn_applier" "fabric_openfabric_applier" {
				  depends_on = [
					proxmox_virtual_environment_sdn_fabric_openfabric.fabric_openfabric
					proxmox_virtual_environment_sdn_fabric_node_openfabric.fabric_node_openfabric
				  ]
				}

				resource "proxmox_virtual_environment_sdn_applier" "finalizer" {}
			`),
			Check: resource.ComposeTestCheckFunc(
				test.ResourceAttributes("proxmox_virtual_environment_sdn_fabric_node_openfabric.fabric_node_openfabric", map[string]string{
					"fabric_id": "fabricS",
					"node_id":   "pve",
					"ip":        "10.0.0.1",
					"state":     "new",
				}),
			),
		}, {
			Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_sdn_fabric_openfabric" "fabric_openfabric" {
				  id  = "fabricS"
				  ip_prefix = "10.0.0.0/16"
				  depends_on = [
					proxmox_virtual_environment_sdn_applier.finalizer
				  ]
				}

				resource "proxmox_virtual_environment_sdn_fabric_node_openfabric" "fabric_node_openfabric" {
				  fabric_id = proxmox_virtual_environment_sdn_fabric_openfabric.fabric_openfabric.id
				  node_id   = "pve"
				  ip 	    = "10.0.0.2"
				  depends_on = [
					proxmox_virtual_environment_sdn_applier.finalizer
				  ]
				}

				resource "proxmox_virtual_environment_sdn_applier" "fabric_openfabric_applier" {
				  depends_on = [
					proxmox_virtual_environment_sdn_fabric_openfabric.fabric_openfabric,
					proxmox_virtual_environment_sdn_fabric_node_openfabric.fabric_node_openfabric,
				  ]
				}

				resource "proxmox_virtual_environment_sdn_applier" "finalizer" {}
			`),
			Check: resource.ComposeTestCheckFunc(
				test.ResourceAttributes("proxmox_virtual_environment_sdn_fabric_node_openfabric.fabric_node_openfabric", map[string]string{
					"fabric_id": "fabricS",
					"node_id":   "pve",
					"ip":        "10.0.0.2",
					"state":     "changed",
				}),
				test.NoResourceAttributesSet("proxmox_virtual_environment_sdn_fabric_node_openfabric.fabric_node_openfabric", []string{
					"ip6",
				}),
			),
			ResourceName:      "proxmox_virtual_environment_sdn_fabric_node_openfabric.fabric_node_openfabric",
			ImportStateId:     "fabricS/pve",
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

func TestAccResourceSDNFabricNodeOSPF(t *testing.T) {
	// Cannot run in parallel due to SDN applier functionality affecting global state

	te := test.InitEnvironment(t)

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"create and update OSPF fabric", []resource.TestStep{{
			Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_sdn_fabric_ospf" "fabric_ospf" {
				  id    = "fabricO"
				  area  = "0"
				  ip_prefix = "10.0.0.0/16"
				}
				resource "proxmox_virtual_environment_sdn_fabric_node_ospf" "fabric_node_ospf" {
				  fabric_id = proxmox_virtual_environment_sdn_fabric_ospf.fabric_ospf.id
				  node_id   = "pve"
				  ip 	    = "10.0.0.1"
				}
			`),
			Check: resource.ComposeTestCheckFunc(
				test.ResourceAttributes("proxmox_virtual_environment_sdn_fabric_node_ospf.fabric_node_ospf", map[string]string{
					"fabric_id": "fabricO",
					"node_id":   "pve",
					"ip":        "10.0.0.1",
					"state":     "new",
				}),
			),
		}, {
			Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_sdn_fabric_ospf" "fabric_ospf" {
				  id    = "fabricO"
				  area  = "0"
				  ip_prefix = "10.0.0.0/16"
				}
				resource "proxmox_virtual_environment_sdn_fabric_node_ospf" "fabric_node_ospf" {
				  fabric_id = proxmox_virtual_environment_sdn_fabric_ospf.fabric_ospf.id
				  node_id   = "pve"
				  ip 	    = "10.0.0.2"
				}
			`),
			Check: resource.ComposeTestCheckFunc(
				test.ResourceAttributes("proxmox_virtual_environment_sdn_fabric_node_ospf.fabric_node_ospf", map[string]string{
					"fabric_id": "fabricO",
					"node_id":   "pve",
					"ip":        "10.0.0.2",
					"state":     "changed",
				}),
			),
			ResourceName:      "proxmox_virtual_environment_sdn_fabric_node_ospf.fabric_node_ospf",
			ImportStateId:     "fabricO/pve",
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
