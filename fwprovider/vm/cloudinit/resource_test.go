/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package cloudinit_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/test"
)

const resourceName = "proxmox_virtual_environment_vm2.test_vm"

func TestAccResourceVM2CloudInit(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)
	te.AddTemplateVars(map[string]interface{}{
		"UpdateVMID": te.RandomVMID(),
	})

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"create VM with cloud-init", []resource.TestStep{{
			Config: te.RenderConfig(`
			resource "proxmox_virtual_environment_vm2" "test_vm" {
				node_name = "{{.NodeName}}"
				id = {{.RandomVMID}}
				name = "test-cloudinit"
				initialization = {
					dns = {
						domain = "example.com"
			        }
				}
			}`),
			Check: resource.ComposeTestCheckFunc(
				test.ResourceAttributes("proxmox_virtual_environment_vm2.test_vm", map[string]string{
					"initialization.datastore_id": te.DatastoreID,
					"initialization.interface":    "ide2",
				}),
			),
		}}},
		{"update VM with cloud-init", []resource.TestStep{
			//{
			//	Config: te.RenderConfig(`
			//	resource "proxmox_virtual_environment_vm2" "test_vm" {
			//		node_name = "{{.NodeName}}"
			//		id = {{.UpdateVMID}}
			//		name = "test-cloudinit"
			//		initialization = {
			//			dns = {
			//				domain = "example.com"
			//			}
			//		}
			//	}`),
			//	Destroy: false,
			//},
			//{
			//	Config: te.RenderConfig(`
			//	resource "proxmox_virtual_environment_vm2" "test_vm" {
			//		node_name = "{{.NodeName}}"
			//		id = {{.UpdateVMID}}
			//		name = "test-cloudinit"
			//		initialization = {
			//			dns = {
			//				domain = "example.com"
			//				servers = [
			//					"1.1.1.1",
			//					"8.8.8.8"
			//				]
			//			}
			//		}
			//	}`),
			//	Destroy: false,
			//},
			//{
			//	Config: te.RenderConfig(`
			//	resource "proxmox_virtual_environment_vm2" "test_vm" {
			//		node_name = "{{.NodeName}}"
			//		id = {{.UpdateVMID}}
			//		name = "test-cloudinit"
			//		initialization = {
			//			dns = {
			//				domain = "another.domain.com"
			//				servers = [
			//					"8.8.8.8",
			//					"1.1.1.1"
			//				]
			//			}
			//		}
			//	}`),
			//	Destroy: false,
			//},
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm2" "test_vm" {
					node_name = "{{.NodeName}}"
					id = {{.UpdateVMID}}
					name = "test-cloudinit"
					initialization = {
						dns = {
							servers = [
								"1.1.1.1"
							]
						}
					}
				}`),
				Destroy: false,
				Check: resource.ComposeTestCheckFunc(
					test.NoResourceAttributesSet("proxmox_virtual_environment_vm2.test_vm", []string{
						"initialization.dns.domain",
					}),
					test.ResourceAttributes("proxmox_virtual_environment_vm2.test_vm", map[string]string{
						"initialization.dns.servers.#": "1",
					}),
				),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm2" "test_vm" {
					node_name = "{{.NodeName}}"
					id = {{.UpdateVMID}}
					name = "test-cloudinit"
					initialization = {
						dns = {
							//servers = []
						}
					}
				}`),
				Destroy: false,
				Check: resource.ComposeTestCheckFunc(
					test.NoResourceAttributesSet("proxmox_virtual_environment_vm2.test_vm", []string{
						"initialization.dns.servers",
					}),
				),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm2" "test_vm" {
					node_name = "{{.NodeName}}"
					id = {{.UpdateVMID}}
					name = "test-cloudinit"
					initialization = {
						dns = {}
					}
				}`),
				Destroy: false,
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm2" "test_vm" {
					node_name = "{{.NodeName}}"
					id = {{.UpdateVMID}}
					name = "test-cloudinit"
					initialization = {}
				}`),
			},
		}},
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
