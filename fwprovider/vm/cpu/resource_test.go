//go:build acceptance_vm || all

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package cpu_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/test"
)

func TestAccResourceVM2CPU(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"create VM with no cpu params", []resource.TestStep{{
			Config: te.RenderConfig(`
			resource "proxmox_virtual_environment_vm2" "test_vm" {
				node_name = "{{.NodeName}}"
				name = "test-cpu"
			}`),
			Check: resource.ComposeTestCheckFunc(
				test.ResourceAttributes("proxmox_virtual_environment_vm2.test_vm", map[string]string{
					// default values that are set by PVE if not specified
					"cpu.cores":   "1",
					"cpu.sockets": "1",
					"cpu.type":    "kvm64",
				}),
			),
		}}},
		{"create VM with some cpu params", []resource.TestStep{{
			Config: te.RenderConfig(`
			resource "proxmox_virtual_environment_vm2" "test_vm" {
				node_name = "{{.NodeName}}"
				name = "test-cpu"
				cpu = {
					cores = 2
					sockets = 2
					type = "host"
					flags = ["+aes"]
				}
			}`),
			Check: resource.ComposeTestCheckFunc(
				test.ResourceAttributes("proxmox_virtual_environment_vm2.test_vm", map[string]string{
					"cpu.cores":   "2",
					"cpu.sockets": "2",
					"cpu.type":    "host",
					"cpu.flags.#": "1",
					"cpu.flags.0": `\+aes`,
				}),
			),
		}}},
		{"create VM with all cpu params and then update them", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm2" "test_vm" {
					node_name = "{{.NodeName}}"
					name = "test-cpu"
					cpu = {
						# affinity = "0-1"          only root can set affinity
						# architecture = "x86_64"   only root can set architecture
						cores = 2
						hotplugged = 2
						limit = 64
						numa = false
						sockets = 2
						type = "host"
						units = 1024
						flags = ["+aes"]
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_virtual_environment_vm2.test_vm", map[string]string{
						"cpu.cores":      "2",
						"cpu.hotplugged": "2",
						"cpu.limit":      "64",
						"cpu.numa":       "false",
						"cpu.sockets":    "2",
						"cpu.type":       "host",
						"cpu.units":      "1024",
					}),
				),
			},
			{ // now update the cpu params and check if they are updated
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm2" "test_vm" {
					node_name = "{{.NodeName}}"
					name = "test-cpu"
					cpu = {
						cores = 4
						hotplugged = 2
						limit = null     # setting to null is the same as removal
						# numa = false
						# sockets = 2    remove sockets, so it should fall back to 1 (PVE default)
						# type = "host"  remove type, so it should fall back to kvm64 (PVE default)
						units = 2048
						# flags = ["+aes"]
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_virtual_environment_vm2.test_vm", map[string]string{
						"cpu.cores":      "4",
						"cpu.hotplugged": "2",
						"cpu.sockets":    "1",     // default value, but it is a special case.
						"cpu.type":       "kvm64", // default value, but it is a special case.
						"cpu.units":      "2048",
					}),
					test.NoResourceAttributesSet("proxmox_virtual_environment_vm2.test_vm", []string{
						"cpu.limit", // other defaults are not set in the state
						"cpu.numa",
						"cpu.flags",
					}),
				),
			},
			{
				RefreshState: true,
			},
		}},
		{"clone VM with some cpu params", []resource.TestStep{{
			Config: te.RenderConfig(`
			resource "proxmox_virtual_environment_vm2" "template_vm" {
				node_name = "{{.NodeName}}"
				name = "template-cpu"
				cpu = {
					cores = 2
					sockets = 2
					type = "host"
				}
			}
			resource "proxmox_virtual_environment_vm2" "test_vm" {
				node_name = "{{.NodeName}}"
				name = "test-cpu"
				clone = {
					id = proxmox_virtual_environment_vm2.template_vm.id
				}	
			}`),
			Check: resource.ComposeTestCheckFunc(
				test.ResourceAttributes("proxmox_virtual_environment_vm2.test_vm", map[string]string{
					"cpu.cores":   "2",
					"cpu.sockets": "2",
					"cpu.type":    "host",
				}),
			),
		}}},
		{"clone VM with some cpu params and updating them in the clone", []resource.TestStep{{
			Config: te.RenderConfig(`
			resource "proxmox_virtual_environment_vm2" "template_vm" {
				node_name = "{{.NodeName}}"
				name = "template-cpu"
				cpu = {
					cores = 2
					sockets = 2
					type = "host"
				}
			}
			resource "proxmox_virtual_environment_vm2" "test_vm" {
				node_name = "{{.NodeName}}"
				name = "test-cpu"
				clone = {
					id = proxmox_virtual_environment_vm2.template_vm.id
				}
				cpu = {
					cores = 4
					units = 1024
				}
			}`),
			Check: resource.ComposeTestCheckFunc(
				test.ResourceAttributes("proxmox_virtual_environment_vm2.test_vm", map[string]string{
					"cpu.cores":   "4",
					"cpu.sockets": "2",
					"cpu.type":    "host",
					"cpu.units":   "1024",
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
