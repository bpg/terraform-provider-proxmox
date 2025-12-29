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

func TestAccDataSourceSDNFabricOpenFabric(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"create openfabric fabric and read with datasource", []resource.TestStep{{
			Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_sdn_fabric_openfabric" "test" {
					id    = "dstest1"
					ip_prefix = "10.0.0.0/16"
				}
				
				data "proxmox_virtual_environment_sdn_fabric_openfabric" "test" {
					id = proxmox_virtual_environment_sdn_fabric_openfabric.test.id
				}
			`),
			Check: resource.ComposeTestCheckFunc(
				test.ResourceAttributes("data.proxmox_virtual_environment_sdn_fabric_openfabric.test", map[string]string{
					"id":        "dstest1",
					"ip_prefix": "10.0.0.0/16",
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
		{"create OSPF fabric and read with datasource", []resource.TestStep{{
			Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_sdn_fabric_ospf" "test" {
					id     = "dstest2"
					area  = "0"
					ip_prefix = "10.0.0.0/16"
				}
				
				data "proxmox_virtual_environment_sdn_fabric_ospf" "test" {
					id = proxmox_virtual_environment_sdn_fabric_ospf.test.id
				}
			`),
			Check: resource.ComposeTestCheckFunc(
				test.ResourceAttributes("data.proxmox_virtual_environment_sdn_fabric_ospf.test", map[string]string{
					"id":        "dstest2",
					"area":      "0",
					"ip_prefix": "10.0.0.0/16",
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
