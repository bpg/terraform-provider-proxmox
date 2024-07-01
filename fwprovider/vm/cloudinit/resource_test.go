/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package cloudinit_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/test"
)

func TestResource_VM2_CloudInit_Create(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"create VM with cloud-init", []resource.TestStep{{
			Config: te.RenderConfig(`
			resource "proxmox_virtual_environment_vm2" "test_vm" {
				node_name = "{{.NodeName}}"
				id = {{.RandomVMID}}
				name = "test-ci"
				initialization = {
					dns = {
						domain = "example.com"
			        }
				}
			}`),
			Check: test.ResourceAttributes("proxmox_virtual_environment_vm2.test_vm", map[string]string{
				"initialization.datastore_id": te.DatastoreID,
				"initialization.interface":    "ide2",
			}),
		}}},
		{"domain can't be empty", []resource.TestStep{{
			Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm2" "test_vm" {
					node_name = "{{.NodeName}}"
					id = {{.RandomVMID}}
					name = "test-ci"
					initialization = {
						dns = {
							domain = ""
						}
					}
				}`),
			ExpectError: regexp.MustCompile(`string length must be at least 1, got: 0`),
		}}},
		{"servers can't be empty", []resource.TestStep{{
			Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm2" "test_vm" {
					node_name = "{{.NodeName}}"
					id = {{.RandomVMID}}
					name = "test-ci"
					initialization = {
						dns = {
							servers = []
						}
					}
				}`),
			ExpectError: regexp.MustCompile(`list must contain at least 1 elements`),
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

func TestResource_VM2_CloudInit_Update(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"add servers", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm2" "test_vm" {
					node_name = "{{.NodeName}}"
					id = {{.RandomVMID}}
					name = "test-ci"
					initialization = {
						dns = {
							domain = "example.com"
						}
					}
				}`),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm2" "test_vm" {
					node_name = "{{.NodeName}}"
					name = "test-ci"
					initialization = {
						dns = {
							domain = "example.com"
							servers = [
								"1.1.1.1",
								"8.8.8.8"
							]
						}
					}
				}`),
			},
		}},
		{"change domain and servers", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm2" "test_vm" {
					node_name = "{{.NodeName}}"
					id = {{.RandomVMID}}
					name = "test-ci"
					initialization = {
						dns = {
							domain = "example.com"
							servers = [
								"1.1.1.1",
								"8.8.8.8"
							]
						}
					}
				}`),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm2" "test_vm" {
					node_name = "{{.NodeName}}"
					name = "test-ci"
					initialization = {
						dns = {
							domain = "another.domain.com"
							servers = [
								"8.8.8.8",
								"1.1.1.1"
							]
						}
					}
				}`),
			},
		}},
		{"update VM: delete dns.domain", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm2" "test_vm" {
					node_name = "{{.NodeName}}"
					id = {{.RandomVMID}}
					name = "test-ci"
					initialization = {
						dns = {
							domain = "example.com"
						}
					}
				}`),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm2" "test_vm" {
					node_name = "{{.NodeName}}"
					name = "test-ci"
					initialization = {
						dns = {}
					}
				}`),
				Check: test.NoResourceAttributesSet("proxmox_virtual_environment_vm2.test_vm", []string{
					"initialization.dns.domain",
				}),
			},
		}},
		{"delete one of the servers", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm2" "test_vm" {
					node_name = "{{.NodeName}}"
					id = {{.RandomVMID}}
					name = "test-ci"
					initialization = {
						dns = {
							servers = [
								"1.1.1.1",
								"8.8.8.8"
							]
						}
					}
				}`),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm2" "test_vm" {
					node_name = "{{.NodeName}}"
					name = "test-ci"
					initialization = {
						dns = {
							domain = "another.domain.com"
							servers = [
								"1.1.1.1"
							]
						}
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("proxmox_virtual_environment_vm2.test_vm", "initialization.dns.servers.#", "1"),
				),
			},
		}},
		{"delete servers", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm2" "test_vm" {
					node_name = "{{.NodeName}}"
					id = {{.RandomVMID}}
					name = "test-ci"
					initialization = {
						dns = {
							servers = [
								"1.1.1.1",
								"8.8.8.8"
							]
						}
					}
				}`),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm2" "test_vm" {
					node_name = "{{.NodeName}}"
					name = "test-ci"
					initialization = {
						dns = {
							// remove, or set to servers = null
						}
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("proxmox_virtual_environment_vm2.test_vm", "initialization.dns.servers.#", "0"),
				),
			},
		}},
		//	{
		//		// step 9: update the VM: remove the dns block
		//		Config: te.RenderConfig(`
		//		resource "proxmox_virtual_environment_vm2" "test_vm" {
		//			node_name = "{{.NodeName}}"
		//			name = "test-ci"
		//			initialization = {}
		//		}`),
		//	},
		//}},

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
