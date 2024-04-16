/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceVMCloneCPU(t *testing.T) {
	te := initTestEnvironment(t)

	tests := []struct {
		name string
		step []resource.TestStep
	}{
		{"copy cpu from template in full", []resource.TestStep{{
			Config: te.renderConfig(`
				resource "proxmox_virtual_environment_vm" "test_cpu_clone_template" {
					node_name = "{{.NodeName}}"
					started   = false
					name 	  = "test-cpu-clone-template"
					template  = "true"
					
					cpu {
						cores = 2
						type  = "host"
					}
				}
				resource "proxmox_virtual_environment_vm" "test_cpu_clone" {
					node_name = "{{.NodeName}}"
					started   = false
					name 	  = "test-cpu-clone"

					clone {
						vm_id = proxmox_virtual_environment_vm.test_cpu_clone_template.id
					}
				}`),
			Check: resource.ComposeTestCheckFunc(
				testResourceAttributes("proxmox_virtual_environment_vm.test_cpu_clone", map[string]string{
					"cpu.0.cores": "2",
					"cpu.0.type":  "host",
				}),
			),
		}}},
		{"merge cpu attributes", []resource.TestStep{{
			Config: te.renderConfig(`
				resource "proxmox_virtual_environment_vm" "test_cpu_clone_template" {
					node_name = "{{.NodeName}}"
					started   = false
					name 	  = "test-cpu-clone-template"
					template  = "true"
					
					cpu {
						cores = 2
					}
				}
				resource "proxmox_virtual_environment_vm" "test_cpu_clone" {
					node_name = "{{.NodeName}}"
					started   = false
					name 	  = "test-cpu-clone"

					cpu {
						type  = "host"
					}

					clone {
						vm_id = proxmox_virtual_environment_vm.test_cpu_clone_template.id
					}
				}`),
			Check: resource.ComposeTestCheckFunc(
				testResourceAttributes("proxmox_virtual_environment_vm.test_cpu_clone", map[string]string{
					"cpu.0.cores": "2",
					"cpu.0.type":  "host",
				}),
			),
		}}},
		{"overwrite cpu attributes in full", []resource.TestStep{{
			Config: te.renderConfig(`
				resource "proxmox_virtual_environment_vm" "test_cpu_clone_template" {
					node_name = "{{.NodeName}}"
					started   = false
					name 	  = "test-cpu-clone-template"
					template  = "true"
				}
				resource "proxmox_virtual_environment_vm" "test_cpu_clone" {
					node_name = "{{.NodeName}}"
					started   = false
					name 	  = "test-cpu-clone"

					cpu {
						cores = 2
						type  = "host"
					}

					clone {
						vm_id = proxmox_virtual_environment_vm.test_cpu_clone_template.id
					}
				}`),
			Check: resource.ComposeTestCheckFunc(
				testResourceAttributes("proxmox_virtual_environment_vm.test_cpu_clone", map[string]string{
					"cpu.0.cores": "2",
					"cpu.0.type":  "host",
				}),
			),
		}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource.ParallelTest(t, resource.TestCase{
				ProtoV6ProviderFactories: te.accProviders,
				Steps:                    tt.step,
			})
		})
	}
}
