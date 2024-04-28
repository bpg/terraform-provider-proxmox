/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package tests

import (
	"regexp"
	"strconv"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceVM2(t *testing.T) {
	t.Parallel()

	te := initTestEnvironment(t)
	vmID := gofakeit.IntRange(90000, 100000)
	te.addTemplateVars(map[string]any{
		"VMID": vmID,
	})

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"create minimal VM", []resource.TestStep{{
			Config: te.renderConfig(`
			resource "proxmox_virtual_environment_vm2" "test_vm" {
				node_name = "{{.NodeName}}"
			}`),
			Check: resource.ComposeTestCheckFunc(
				testResourceAttributes("proxmox_virtual_environment_vm2.test_vm", map[string]string{
					"node_name": te.nodeName,
				}),
				testResourceAttributesSet("proxmox_virtual_environment_vm2.test_vm", []string{
					"id",
				}),
			),
		}}},
		{"create minimal VM with ID", []resource.TestStep{{
			Config: te.renderConfig(`
			resource "proxmox_virtual_environment_vm2" "test_vm" {
				node_name = "{{.NodeName}}"
				
				id = {{.VMID}}
			}`),
			Check: resource.ComposeTestCheckFunc(
				testResourceAttributes("proxmox_virtual_environment_vm2.test_vm", map[string]string{
					"node_name": te.nodeName,
					"id":        strconv.Itoa(vmID),
				}),
			),
		}}},
		{"set an invalid VM name", []resource.TestStep{{
			Config: te.renderConfig(`
			resource "proxmox_virtual_environment_vm2" "test_vm" {
				node_name = "{{.NodeName}}"
				
				name = "not a valid DNS name"

			}`),
			ExpectError: regexp.MustCompile(`name must be a valid DNS name`),
		}}},
		{"set, update, import with primitive fields", []resource.TestStep{
			{
				Config: te.renderConfig(`
				resource "proxmox_virtual_environment_vm2" "test_vm" {
					node_name = "{{.NodeName}}"
					
					name = "test-vm"
					description = "test description"
				}`),
				Check: testResourceAttributes("proxmox_virtual_environment_vm2.test_vm", map[string]string{
					"name":        "test-vm",
					"description": "test description",
				}),
			},
			{
				Config: te.renderConfig(`
				resource "proxmox_virtual_environment_vm2" "test_vm" {
					node_name = "{{.NodeName}}"
					
					name = "test-vm"
				}`),
				Check: resource.ComposeTestCheckFunc(
					testResourceAttributes("proxmox_virtual_environment_vm2.test_vm", map[string]string{
						"name": "test-vm",
					}),
					testNoResourceAttributesSet("proxmox_virtual_environment_vm2.test_vm", []string{
						"description",
					}),
				),
			},
			{
				ResourceName:        "proxmox_virtual_environment_vm2.test_vm",
				ImportState:         true,
				ImportStateVerify:   true,
				ImportStateIdPrefix: te.nodeName + "/",
			},
		}},
		{"set, update, import with tags", []resource.TestStep{
			{
				Config: te.renderConfig(`
				resource "proxmox_virtual_environment_vm2" "test_vm" {
					node_name = "{{.NodeName}}"
					
					name = "test-tags"
					tags = ["tag2", "tag1"]
				}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckTypeSetElemAttr("proxmox_virtual_environment_vm2.test_vm", "tags.*", "tag1"),
					resource.TestCheckTypeSetElemAttr("proxmox_virtual_environment_vm2.test_vm", "tags.*", "tag2"),
				),
			},
			{
				Config: te.renderConfig(`
				resource "proxmox_virtual_environment_vm2" "test_vm" {
					node_name = "{{.NodeName}}"
					name = "test-tags"
					tags = ["tag1"]
				}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("proxmox_virtual_environment_vm2.test_vm", "tags.#", "1"),
					resource.TestCheckTypeSetElemAttr("proxmox_virtual_environment_vm2.test_vm", "tags.*", "tag1"),
				),
			},
			{
				Config: te.renderConfig(`
				resource "proxmox_virtual_environment_vm2" "test_vm" {
					node_name = "{{.NodeName}}"
					name = "test-tags"
					// no tags
				}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("proxmox_virtual_environment_vm2.test_vm", "tags.#", "1"),
					resource.TestCheckTypeSetElemAttr("proxmox_virtual_environment_vm2.test_vm", "tags.*", "tag1"),
				),
			},
			{
				Config: te.renderConfig(`
				resource "proxmox_virtual_environment_vm2" "test_vm" {
					node_name = "{{.NodeName}}"
					name = "test-tags"
					tags = []
				}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("proxmox_virtual_environment_vm2.test_vm", "tags.#", "0"),
				),
			},
		}},
		{"a VM can't have empty tags", []resource.TestStep{{
			Config: te.renderConfig(`
			resource "proxmox_virtual_environment_vm2" "test_vm" {
				node_name = "{{.NodeName}}"
		
				tags = ["", "tag1"]
			}`),
			ExpectError: regexp.MustCompile(`string length must be at least 1, got: 0`),
		}}},
		{"a VM can't have empty tags", []resource.TestStep{{
			Config: te.renderConfig(`
			resource "proxmox_virtual_environment_vm2" "test_vm" {
				node_name = "{{.NodeName}}"
		
				tags = [" ", "tag1"]
			}`),
			ExpectError: regexp.MustCompile(`must be a non-empty and non-whitespace string`),
		}}},
		{"multiline description", []resource.TestStep{{
			Config: te.renderConfig(`
				resource "proxmox_virtual_environment_vm2" "test_vm" {
					node_name = "{{.NodeName}}"
					
					description = trimspace(<<-EOT
						my
						description
						value
					EOT
					)
				}`),
			Check: resource.ComposeTestCheckFunc(
				testResourceAttributes("proxmox_virtual_environment_vm2.test_vm", map[string]string{
					"description": "my\ndescription\nvalue",
				}),
			),
		}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource.ParallelTest(t, resource.TestCase{
				ProtoV6ProviderFactories: te.accProviders,
				Steps:                    tt.steps,
			})
		})
	}
}

func TestAccResourceVM2CPU(t *testing.T) {
	t.Parallel()

	te := initTestEnvironment(t)

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"create VM with no cpu params", []resource.TestStep{{
			Config: te.renderConfig(`
			resource "proxmox_virtual_environment_vm2" "test_vm" {
				node_name = "{{.NodeName}}"
				name = "test-cpu"
			}`),
			Check: resource.ComposeTestCheckFunc(
				testResourceAttributes("proxmox_virtual_environment_vm2.test_vm", map[string]string{
					// default values that are set by PVE if not specified
					"cpu.cores":   "1",
					"cpu.sockets": "1",
					"cpu.type":    "kvm64",
				}),
			),
		}}},
		{"create VM with some cpu params", []resource.TestStep{{
			Config: te.renderConfig(`
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
				testResourceAttributes("proxmox_virtual_environment_vm2.test_vm", map[string]string{
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
				Config: te.renderConfig(`
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
					testResourceAttributes("proxmox_virtual_environment_vm2.test_vm", map[string]string{
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
				Config: te.renderConfig(`
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
					testResourceAttributes("proxmox_virtual_environment_vm2.test_vm", map[string]string{
						"cpu.cores":      "4",
						"cpu.hotplugged": "2",
						"cpu.sockets":    "1",     // default value, but it is a special case.
						"cpu.type":       "kvm64", // default value, but it is a special case.
						"cpu.units":      "2048",
					}),
					testNoResourceAttributesSet("proxmox_virtual_environment_vm2.test_vm", []string{
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
			Config: te.renderConfig(`
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
				testResourceAttributes("proxmox_virtual_environment_vm2.test_vm", map[string]string{
					"cpu.cores":   "2",
					"cpu.sockets": "2",
					"cpu.type":    "host",
				}),
			),
		}}},
		{"clone VM with some cpu params and updating them in the clone", []resource.TestStep{{
			Config: te.renderConfig(`
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
				testResourceAttributes("proxmox_virtual_environment_vm2.test_vm", map[string]string{
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
				ProtoV6ProviderFactories: te.accProviders,
				Steps:                    tt.steps,
			})
		})
	}
}

func TestAccResourceVM2Clone(t *testing.T) {
	t.Parallel()

	te := initTestEnvironment(t)
	vmID := gofakeit.IntRange(90000, 100000)
	te.addTemplateVars(map[string]any{
		"VMID": vmID,
	})

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"create a clone from template", []resource.TestStep{{
			Config: te.renderConfig(`
			resource "proxmox_virtual_environment_vm2" "test_vm" {
				node_name = "{{.NodeName}}"
				name = "template"
				description = "template description"
				template = true
			}
			resource "proxmox_virtual_environment_vm2" "test_vm_clone" {
				node_name = "{{.NodeName}}"
				name = "clone"
				clone = {
					id = proxmox_virtual_environment_vm2.test_vm.id	
				}
			}`),
			Check: resource.ComposeTestCheckFunc(
				testResourceAttributes("proxmox_virtual_environment_vm2.test_vm", map[string]string{
					"template": "true",
				}),
				testResourceAttributes("proxmox_virtual_environment_vm2.test_vm_clone", map[string]string{
					// name is overwritten
					"name": "clone",
				}),
				testNoResourceAttributesSet("proxmox_virtual_environment_vm2.test_vm_clone", []string{
					// description is not copied
					"description",
				}),
			),
		}}},
		{"tags are copied to the clone", []resource.TestStep{{
			Config: te.renderConfig(`
			resource "proxmox_virtual_environment_vm2" "test_vm" {
				node_name = "{{.NodeName}}"
				template = true
				tags = ["tag1", "tag2"]
			}
			resource "proxmox_virtual_environment_vm2" "test_vm_clone" {
				node_name = "{{.NodeName}}"
				clone = {
					id = proxmox_virtual_environment_vm2.test_vm.id
				}
			}`),
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckTypeSetElemAttr("proxmox_virtual_environment_vm2.test_vm_clone", "tags.*", "tag1"),
				resource.TestCheckTypeSetElemAttr("proxmox_virtual_environment_vm2.test_vm_clone", "tags.*", "tag2"),
			),
		}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource.ParallelTest(t, resource.TestCase{
				ProtoV6ProviderFactories: te.accProviders,
				Steps:                    tt.steps,
			})
		})
	}
}
