//go:build acceptance || all

//testacc:tier=medium
//testacc:resource=vm

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package vm_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/test"
)

func TestAccDatasourceVM(t *testing.T) {
	te := test.InitEnvironment(t)

	datasourceName := "data.proxmox_virtual_environment_vm2.test"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: te.RenderConfig(`
					resource "proxmox_vm" "test_vm" {
						node_name = "{{.NodeName}}"
					}

					data "proxmox_virtual_environment_vm2" "test" {
						node_name = "{{.NodeName}}"
						id        = proxmox_vm.test_vm.id
					}
				`),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "node_name", te.NodeName),
					resource.TestCheckResourceAttr(datasourceName, "status", "stopped"),
					resource.TestCheckResourceAttrSet(datasourceName, "id"),
				),
			},
		},
	})
}

func TestAccDatasourceVMNotFound(t *testing.T) {
	te := test.InitEnvironment(t)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: te.RenderConfig(`
					data "proxmox_virtual_environment_vm2" "test" {
						node_name = "{{.NodeName}}"
						id        = 999999
					}
				`),
				ExpectError: regexp.MustCompile("VM Not Found"),
			},
		},
	})
}
