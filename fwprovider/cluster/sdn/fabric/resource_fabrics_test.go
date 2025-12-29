//go:build acceptance || all

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package fabric_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/test"
)

func TestAccResourceSDNFabricOpenFabric(t *testing.T) {
	// Cannot run in parallel due to SDN applier functionality affecting global state

	te := test.InitEnvironment(t)

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"create and update fabrics", []resource.TestStep{{
			Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_sdn_fabric_openfabric" "fabric_openfabric" {
				  id  = "fabricS"
				  ip_prefix = "10.0.0.0/16
				  depends_on = [
					proxmox_virtual_environment_sdn_applier.finalizer
				  ]
				}

				resource "proxmox_virtual_environment_sdn_applier" "fabric_openfabric_applier" {
				  depends_on = [
					proxmox_virtual_environment_sdn_fabric_openfabric.fabric_openfabric
				  ]
				}

				resource "proxmox_virtual_environment_sdn_applier" "finalizer" {}
			`),
			Check: resource.ComposeTestCheckFunc(
				test.ResourceAttributes("proxmox_virtual_environment_sdn_fabric_openfabric.fabric_openfabric", map[string]string{
					"id":        "fabricS",
					"ip_prefix": "10.0.0.0/16",
					"state":     "new",
				}),
			),
		}, {
			Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_sdn_fabric_openfabric" "fabric_openfabric" {
				  id  = "fabricS"
				  ip_prefix = "10.0.1.0/16"
				  depends_on = [
					proxmox_virtual_environment_sdn_applier.finalizer
				  ]
				}

				resource "proxmox_virtual_environment_sdn_applier" "fabric_openfabric_applier" {
				  depends_on = [
					proxmox_virtual_environment_sdn_fabric_openfabric.fabric_openfabric
				  ]
				}

				resource "proxmox_virtual_environment_sdn_applier" "finalizer" {}
			`),
			Check: resource.ComposeTestCheckFunc(
				test.ResourceAttributes("proxmox_virtual_environment_sdn_fabric_openfabric.fabric_openfabric", map[string]string{
					"id":        "fabricS",
					"ip_prefix": "10.0.1.0/16",
					"state":     "changed",
				}),
				test.NoResourceAttributesSet("proxmox_virtual_environment_sdn_fabric_openfabric.fabric_openfabric", []string{
					"ip6_prefix",
				}),
			),
			ResourceName: "proxmox_virtual_environment_sdn_fabric_openfabric.fabric_openfabric",
			// ImportStateId:     "fabricS",
			// ImportState:       true,
			// ImportStateVerify: true,
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

func TestAccResourceSDNFabricOSPF(t *testing.T) {
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
			`),
			Check: resource.ComposeTestCheckFunc(
				test.ResourceAttributes("proxmox_virtual_environment_sdn_fabric_ospf.fabric_ospf", map[string]string{
					"id":    "fabricO",
					"area":  "0",
					"state": "new",
				}),
			),
		}, {
			Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_sdn_fabric_ospf" "fabric_ospf" {
				  id    = "fabricO"
				  area  = "0"
				  ip_prefix = "10.0.1.0/16"
				}
			`),
			Check: resource.ComposeTestCheckFunc(
				test.ResourceAttributes("proxmox_virtual_environment_sdn_fabric_ospf.fabric_ospf", map[string]string{
					"id":        "fabricO",
					"area":      "0",
					"ip_prefix": "10.0.1.0/16",
					"state":     "changed",
				}),
			),
			ResourceName:      "proxmox_virtual_environment_sdn_fabric_ospf.fabric_ospf",
			ImportStateId:     "fabricO",
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
