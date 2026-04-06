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

	datasourceName := "data.proxmox_vm.test"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: te.RenderConfig(`
					resource "proxmox_vm" "test_vm" {
						node_name = "{{.NodeName}}"
						name      = "test-datasource-vm"
						tags      = ["tag1", "tag2"]
					}

					data "proxmox_vm" "test" {
						node_name = "{{.NodeName}}"
						id        = proxmox_vm.test_vm.id
					}
				`),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceName, "id"),
					resource.TestCheckResourceAttr(datasourceName, "node_name", te.NodeName),
					resource.TestCheckResourceAttr(datasourceName, "name", "test-datasource-vm"),
					resource.TestCheckResourceAttr(datasourceName, "status", "stopped"),
					resource.TestCheckResourceAttr(datasourceName, "template", "false"),
					resource.TestCheckResourceAttr(datasourceName, "tags.#", "2"),
					resource.TestCheckResourceAttr(datasourceName, "description", ""),
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
					data "proxmox_vm" "test" {
						node_name = "{{.NodeName}}"
						id        = 999999
					}
				`),
				ExpectError: regexp.MustCompile("VM Not Found"),
			},
		},
	})
}
