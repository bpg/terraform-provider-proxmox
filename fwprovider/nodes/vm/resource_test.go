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
	"math/rand"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/test"
)

func TestAccResourceVMShort(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)
	te.AddTemplateVars(map[string]interface{}{
		"TestVMID": 100000 + rand.Intn(99999),
	})

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"create minimal VM", []resource.TestStep{{
			Config: te.RenderConfig(`
			resource "proxmox_vm" "test_vm" {
				node_name = "{{.NodeName}}"
			}`),
			Check: resource.ComposeTestCheckFunc(
				test.ResourceAttributes("proxmox_vm.test_vm", map[string]string{
					"node_name": te.NodeName,
				}),
				test.ResourceAttributesSet("proxmox_vm.test_vm", []string{
					"id",
				}),
			),
		}}},
		{"create minimal VM with ID", []resource.TestStep{{
			Config: te.RenderConfig(`
			resource "proxmox_vm" "test_vm" {
				node_name = "{{.NodeName}}"
				id = {{.TestVMID}}
			}`),
		}}},
		{"set an invalid VM name", []resource.TestStep{{
			Config: te.RenderConfig(`
			resource "proxmox_vm" "test_vm" {
				node_name = "{{.NodeName}}"
				name = "not a valid DNS name"
			}`),
			ExpectError: regexp.MustCompile(`name must be a valid DNS name`),
		}}},
		{"set a FQDN VM name", []resource.TestStep{{
			Config: te.RenderConfig(`
			resource "proxmox_vm" "test_vm" {
				node_name = "{{.NodeName}}"
				name = "vm.example.com"
			}`),
			Check: test.ResourceAttributes("proxmox_vm.test_vm", map[string]string{
				"name": "vm.example.com",
			}),
		}}},
		{"set, update, import with primitive fields", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_vm" "test_vm" {
					node_name = "{{.NodeName}}"
					name = "test-vm"
					description = "test description"
				}`),
				Check: test.ResourceAttributes("proxmox_vm.test_vm", map[string]string{
					"name":        "test-vm",
					"description": "test description",
				}),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_vm" "test_vm" {
					node_name = "{{.NodeName}}"
					name = "test-vm"
				}`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_vm.test_vm", map[string]string{
						"name": "test-vm",
					}),
					test.NoResourceAttributesSet("proxmox_vm.test_vm", []string{
						"description",
					}),
				),
			},
			{
				ResourceName:        "proxmox_vm.test_vm",
				ImportState:         true,
				ImportStateVerify:   true,
				ImportStateIdPrefix: te.NodeName + "/",
			},
		}},
		{"set, update, import with tags", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_vm" "test_vm" {
					node_name = "{{.NodeName}}"
					name = "test-tags"
					tags = ["tag2", "tag1"]
				}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckTypeSetElemAttr("proxmox_vm.test_vm", "tags.*", "tag1"),
					resource.TestCheckTypeSetElemAttr("proxmox_vm.test_vm", "tags.*", "tag2"),
				),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_vm" "test_vm" {
					node_name = "{{.NodeName}}"
					name = "test-tags"
					tags = ["tag1"]
				}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("proxmox_vm.test_vm", "tags.#", "1"),
					resource.TestCheckTypeSetElemAttr("proxmox_vm.test_vm", "tags.*", "tag1"),
				),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_vm" "test_vm" {
					node_name = "{{.NodeName}}"
					name = "test-tags"
					// no tags
				}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("proxmox_vm.test_vm", "tags.#", "1"),
					resource.TestCheckTypeSetElemAttr("proxmox_vm.test_vm", "tags.*", "tag1"),
				),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_vm" "test_vm" {
					node_name = "{{.NodeName}}"
					name = "test-tags"
					tags = []
				}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("proxmox_vm.test_vm", "tags.#", "0"),
				),
			},
		}},
		{"a VM can't have empty tags", []resource.TestStep{{
			Config: te.RenderConfig(`
			resource "proxmox_vm" "test_vm" {
				node_name = "{{.NodeName}}"
				tags = ["", "tag1"]
			}`),
			ExpectError: regexp.MustCompile(`string length must be at least 1, got: 0`),
		}}},
		{"a VM can't have empty tags", []resource.TestStep{{
			Config: te.RenderConfig(`
			resource "proxmox_vm" "test_vm" {
				node_name = "{{.NodeName}}"
				tags = [" ", "tag1"]
			}`),
			ExpectError: regexp.MustCompile(`must be a non-empty and non-whitespace string`),
		}}},
		{"multiline description", []resource.TestStep{{
			Config: te.RenderConfig(`
				resource "proxmox_vm" "test_vm" {
					node_name = "{{.NodeName}}"
					description = trimspace(<<-EOT
						my
						description
						value
					EOT
					)
				}`),
			Check: resource.ComposeTestCheckFunc(
				test.ResourceAttributes("proxmox_vm.test_vm", map[string]string{
					"description": "my\ndescription\nvalue",
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

func TestAccResourceVMv2Agent(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"set, update, remove agent", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_vm" "test_agent" {
					node_name = "{{.NodeName}}"
					name      = "test-agent"

					agent = {
						enabled = true
						type    = "virtio"
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_vm.test_agent", map[string]string{
						"agent.enabled": "true",
						"agent.type":    "virtio",
					}),
				),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_vm" "test_agent" {
					node_name = "{{.NodeName}}"
					name      = "test-agent"

					agent = {
						enabled = false
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_vm.test_agent", map[string]string{
						"agent.enabled": "false",
					}),
					test.NoResourceAttributesSet("proxmox_vm.test_agent", []string{
						"agent.type",
					}),
				),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_vm" "test_agent" {
					node_name = "{{.NodeName}}"
					name      = "test-agent"
				}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("proxmox_vm.test_agent", "agent"),
				),
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

func TestAccResourceVMv2Memory(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"set, update, remove memory", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_vm" "test_mem" {
					node_name = "{{.NodeName}}"
					name      = "test-memory"

					memory = {
						size    = 1024
						balloon = 512
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_vm.test_mem", map[string]string{
						"memory.size":    "1024",
						"memory.balloon": "512",
					}),
				),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_vm" "test_mem" {
					node_name = "{{.NodeName}}"
					name      = "test-memory"

					memory = {
						size = 2048
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_vm.test_mem", map[string]string{
						"memory.size": "2048",
					}),
					test.NoResourceAttributesSet("proxmox_vm.test_mem", []string{
						"memory.balloon",
					}),
				),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_vm" "test_mem" {
					node_name = "{{.NodeName}}"
					name      = "test-memory"
				}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("proxmox_vm.test_mem", "memory"),
				),
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

func TestAccResourceVMv2NetworkDevice(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"create, update, remove network devices", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_vm" "test_net" {
					node_name = "{{.NodeName}}"
					name      = "test-network-device"

					network_device = [
						{
							model  = "virtio"
							bridge = "vmbr0"
						},
						{
							model  = "e1000"
							bridge = "vmbr0"
						}
					]
				}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("proxmox_vm.test_net", "network_device.#", "2"),
					resource.TestCheckResourceAttr("proxmox_vm.test_net", "network_device.0.model", "virtio"),
					resource.TestCheckResourceAttr("proxmox_vm.test_net", "network_device.0.bridge", "vmbr0"),
					resource.TestCheckResourceAttr("proxmox_vm.test_net", "network_device.1.model", "e1000"),
					// mac_address is Optional+Computed — PVE generates it when not set
					test.ResourceAttributesSet("proxmox_vm.test_net", []string{
						"network_device.0.mac_address",
						"network_device.1.mac_address",
					}),
				),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_vm" "test_net" {
					node_name = "{{.NodeName}}"
					name      = "test-network-device"

					network_device = [
						{
							model  = "e1000e"
							bridge = "vmbr0"
						}
					]
				}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("proxmox_vm.test_net", "network_device.#", "1"),
					resource.TestCheckResourceAttr("proxmox_vm.test_net", "network_device.0.model", "e1000e"),
				),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_vm" "test_net" {
					node_name = "{{.NodeName}}"
					name      = "test-network-device"
				}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("proxmox_vm.test_net", "network_device"),
				),
			},
		}},
		{"network device with VLAN and firewall", []resource.TestStep{{
			Config: te.RenderConfig(`
			resource "proxmox_vm" "test_net_vlan" {
				node_name = "{{.NodeName}}"
				name      = "test-network-device-vlan"

				network_device = [{
					model    = "virtio"
					bridge   = "vmbr0"
					vlan_id  = 100
					firewall = true
				}]
			}`),
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("proxmox_vm.test_net_vlan", "network_device.0.model", "virtio"),
				resource.TestCheckResourceAttr("proxmox_vm.test_net_vlan", "network_device.0.vlan_id", "100"),
				resource.TestCheckResourceAttr("proxmox_vm.test_net_vlan", "network_device.0.firewall", "true"),
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

func TestAccResourceVMv2Started(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"create VM with started=false", []resource.TestStep{{
			Config: te.RenderConfig(`
			resource "proxmox_vm" "test_started" {
				node_name = "{{.NodeName}}"
				name      = "test-started-false"
				started   = false
			}`),
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("proxmox_vm.test_started", "started", "false"),
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

func TestAccResourceVMv2Clone(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"clone from template, override network device", []resource.TestStep{{
			// Template has one virtio device; clone specifies its own e1000 device.
			// The plan overlay must prevent the template's network config from bleeding
			// into state and causing "inconsistent result" errors.
			Config: te.RenderConfig(`
			resource "proxmox_vm" "template" {
				node_name = "{{.NodeName}}"
				name      = "test-clone-template"
				template  = true
				started   = false

				network_device = [{
					model  = "virtio"
					bridge = "vmbr0"
				}]
			}
			resource "proxmox_vm" "clone" {
				node_name = "{{.NodeName}}"
				name      = "test-clone-vm"
				started   = false

				clone = {
					vm_id = proxmox_vm.template.id
					full  = true
				}

				network_device = [{
					model  = "e1000"
					bridge = "vmbr0"
				}]
			}`),
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("proxmox_vm.clone", "network_device.#", "1"),
				resource.TestCheckResourceAttr("proxmox_vm.clone", "network_device.0.model", "e1000"),
				// clone block is preserved in state from the plan (cannot be read from API)
				resource.TestCheckResourceAttr("proxmox_vm.clone", "clone.full", "true"),
			),
		}}},
		{"clone from template, inherit cpu block", []resource.TestStep{{
			// Template has cpu.sockets=2; clone sets only cores. The plan overlay must
			// keep cpu.sockets null in state (user didn't set it) rather than leaking the
			// template's value and triggering "inconsistent result after apply".
			Config: te.RenderConfig(`
			resource "proxmox_vm" "template" {
				node_name = "{{.NodeName}}"
				name      = "test-clone-cpu-template"
				template  = true
				started   = false

				cpu = {
					sockets = 2
					cores   = 1
				}
			}
			resource "proxmox_vm" "clone" {
				node_name = "{{.NodeName}}"
				name      = "test-clone-cpu-vm"
				started   = false

				clone = {
					vm_id = proxmox_vm.template.id
					full  = true
				}

				cpu = {
					cores = 2
				}
			}`),
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("proxmox_vm.clone", "cpu.cores", "2"),
				// sockets was not set in the clone plan — must stay null, not leak from template
				resource.TestCheckNoResourceAttr("proxmox_vm.clone", "cpu.sockets"),
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

func TestAccResourceVMv2Initialization(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"create VM with user_account and DNS", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_vm" "test_init" {
					node_name = "{{.NodeName}}"
					name      = "test-init-cloudinit"

					initialization = {
						user_account = {
							username = "ubuntu"
							password = "s3cr3t"
						}
						dns = {
							domain  = "example.com"
							servers = ["1.1.1.1", "8.8.8.8"]
						}
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_vm.test_init", map[string]string{
						"initialization.user_account.username": "ubuntu",
						"initialization.dns.domain":            "example.com",
					}),
					// password is write-only — must NOT appear in state
					resource.TestCheckNoResourceAttr("proxmox_vm.test_init", "initialization.user_account.password"),
					resource.TestCheckResourceAttr("proxmox_vm.test_init", "initialization.dns.servers.#", "2"),
					resource.TestCheckResourceAttr("proxmox_vm.test_init", "initialization.dns.servers.0", "1.1.1.1"),
				),
			},
			{
				// Update: change username; password absence should not trigger replacement
				Config: te.RenderConfig(`
				resource "proxmox_vm" "test_init" {
					node_name = "{{.NodeName}}"
					name      = "test-init-cloudinit"

					initialization = {
						user_account = {
							username = "admin"
						}
						dns = {
							domain  = "example.com"
							servers = ["1.1.1.1"]
						}
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_vm.test_init", map[string]string{
						"initialization.user_account.username": "admin",
						"initialization.dns.servers.#":         "1",
					}),
				),
			},
		}},
		{"create VM with ip_config", []resource.TestStep{{
			Config: te.RenderConfig(`
			resource "proxmox_vm" "test_init_ip" {
				node_name = "{{.NodeName}}"
				name      = "test-init-ipconfig"

				initialization = {
					ip_config = [{
						ipv4_address = "dhcp"
					}]
				}
			}`),
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("proxmox_vm.test_init_ip", "initialization.ip_config.#", "1"),
				resource.TestCheckResourceAttr("proxmox_vm.test_init_ip", "initialization.ip_config.0.ipv4_address", "dhcp"),
			),
		}}},
		{"remove initialization block", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_vm" "test_init_remove" {
					node_name = "{{.NodeName}}"
					name      = "test-init-remove"

					initialization = {
						dns = {
							domain = "remove.me"
						}
					}
				}`),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_vm" "test_init_remove" {
					node_name = "{{.NodeName}}"
					name      = "test-init-remove"
				}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("proxmox_vm.test_init_remove", "initialization"),
				),
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
