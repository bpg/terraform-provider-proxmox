//go:build acceptance || all

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/bpg/terraform-provider-proxmox/utils"
)

func TestAccResourceVM(t *testing.T) {
	t.Parallel()

	te := InitEnvironment(t)
	dirName := fmt.Sprintf("dir_%s", gofakeit.Word())
	te.AddTemplateVars(map[string]interface{}{
		"DirName": dirName,
	})

	tests := []struct {
		name string
		step []resource.TestStep
	}{
		{"multiline description", []resource.TestStep{{
			Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_vm1" {
					node_name = "{{.NodeName}}"
					started   = false
					
					description = <<-EOT
						my
						description
						value
					EOT
				}`),
			Check: resource.ComposeTestCheckFunc(
				ResourceAttributes("proxmox_virtual_environment_vm.test_vm1", map[string]string{
					"description": "my\ndescription\nvalue",
				}),
			),
		}}},
		{"single line description", []resource.TestStep{{
			Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_vm2" {
					node_name = "{{.NodeName}}"
					started   = false
					
					description = "my description value"
				}`),
			Check: resource.ComposeTestCheckFunc(
				ResourceAttributes("proxmox_virtual_environment_vm.test_vm2", map[string]string{
					"description": "my description value",
				}),
			),
		}}},
		{"no description", []resource.TestStep{{
			Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_vm3" {
					node_name = "{{.NodeName}}"
					started   = false
					
					description = ""
				}`),
			Check: resource.ComposeTestCheckFunc(
				ResourceAttributes("proxmox_virtual_environment_vm.test_vm3", map[string]string{
					"description": "",
				}),
			),
		}}},
		{"empty node_name", []resource.TestStep{{
			Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_empty_node_name" {
					node_name = ""
					started   = false	
				}`),
			ExpectError: regexp.MustCompile(`expected "node_name" to not be an empty string, got `),
		}}},
		{"protection", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_vm4" {
					node_name = "{{.NodeName}}"
					started   = false
					
					protection = true
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_vm.test_vm4", map[string]string{
						"protection": "true",
					}),
				),
			}, {
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_vm4" {
					node_name = "{{.NodeName}}"
					started   = false
					
					protection = false
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_vm.test_vm4", map[string]string{
						"protection": "false",
					}),
				),
			},
		}},
		{"update cpu block", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_vm5" {
					node_name = "{{.NodeName}}"
					started   = false
					
					cpu {
						cores = 2
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_vm.test_vm5", map[string]string{
						"cpu.0.sockets": "1",
					}),
				),
			}, {
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_vm5" {
					node_name = "{{.NodeName}}"
					started   = false
					
					cpu {
						cores = 1
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_vm.test_vm5", map[string]string{
						"cpu.0.sockets": "1",
					}),
				),
			},
		}},
		{"set cpu.architecture as non root is not supported", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_cpu_arch" {
					node_name = "{{.NodeName}}"
					started   = false
					cpu {
						architecture = "x86_64"
					}
				}`, WithAPIToken()),
				ExpectError: regexp.MustCompile(`can only be set by the root account`),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "template" {
					node_name = "{{.NodeName}}"
					started   = false
					cpu {
						architecture = "x86_64"
					}
				}`, WithRootUser()),
				Destroy: false,
			},
		}},
		{"update memory block", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_vm6" {
					node_name = "{{.NodeName}}"
					started   = false
					
					memory {
						dedicated = 2048
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_vm.test_vm6", map[string]string{
						"memory.0.dedicated": "2048",
					}),
				),
			}, {
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_vm6" {
					node_name = "{{.NodeName}}"
					started   = false
					
					memory {
						dedicated = 1024
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_vm.test_vm6", map[string]string{
						"memory.0.dedicated": "1024",
					}),
				),
			},
		}},
		{"create vga block", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_vm" {
					node_name = "{{.NodeName}}"
					started   = false
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_vm.test_vm", map[string]string{
						"vga.#": "0",
					}),
				),
			}, {
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_vm" {
					node_name = "{{.NodeName}}"
					started   = false
					vga {
						type = "virtio-gl"
						clipboard = "vnc"
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_vm.test_vm", map[string]string{
						"vga.#": "1",
					}),
				),
			},
		}},
		{"update vga block", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_vm" {
					node_name = "{{.NodeName}}"
					started   = false
					
					vga {
						type = "none"
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_vm.test_vm", map[string]string{
						"vga.0.type":      "none",
						"vga.0.clipboard": "",
					}),
				),
			}, {
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_vm" {
					node_name = "{{.NodeName}}"
					started   = false
					
					vga {
						type = "virtio-gl"
						clipboard = "vnc"
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_vm.test_vm", map[string]string{
						"vga.0.type":      "virtio-gl",
						"vga.0.clipboard": "vnc",
					}),
				),
			}, {
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_vm" {
					node_name = "{{.NodeName}}"
					started   = false
					
					vga {
						type = "virtio-gl"
						clipboard = ""
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_vm.test_vm", map[string]string{
						"vga.0.type":      "virtio-gl",
						"vga.0.clipboard": "",
					}),
				),
			},
		}},
		{"update watchdog block", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_vm" {
					node_name = "{{.NodeName}}"
					started   = false
					
					watchdog {
						enabled = "true"
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_vm.test_vm", map[string]string{
						"watchdog.0.model":  "i6300esb",
						"watchdog.0.action": "none",
					}),
				),
			}, {
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_vm" {
					node_name = "{{.NodeName}}"
					started   = false
			
					watchdog {
						enabled = "true"
						model   = "ib700"
						action  = "reset"
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_vm.test_vm", map[string]string{
						"watchdog.0.model":  "ib700",
						"watchdog.0.action": "reset",
					}),
				),
			}, {
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_vm" {
					node_name = "{{.NodeName}}"
					started   = false
			
					watchdog {
						enabled = "false"
						model   = "ib700"
						action  = "reset"
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_vm.test_vm", map[string]string{
						"watchdog.0.enabled": "false",
						"watchdog.0.model":   "ib700",
						"watchdog.0.action":  "reset",
					}),
				),
			},
		}},
		{"update rng block", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_vm" {
					node_name = "{{.NodeName}}"
					started   = false
					
					rng {
						source = "/dev/urandom"
					}
				}`, WithRootUser()),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_vm.test_vm", map[string]string{
						"rng.0.source":    "/dev/urandom",
						"rng.0.max_bytes": "1024",
						"rng.0.period":    "1000",
					}),
				),
			}, {
				Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_vm" "test_vm" {
						node_name = "{{.NodeName}}"
						started   = false
				
						rng {
							source = "/dev/urandom"
							max_bytes = 2048
							period = 500
						}
					}`, WithRootUser()),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_vm.test_vm", map[string]string{
						"rng.0.source":    "/dev/urandom",
						"rng.0.max_bytes": "2048",
						"rng.0.period":    "500",
					}),
				),
			}, {
				Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_vm" "test_vm" {
						node_name = "{{.NodeName}}"
						started   = false
				
						rng {
							source = "/dev/random"
							max_bytes = 512
							period = 200
						}
					}`, WithRootUser()),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_vm.test_vm", map[string]string{
						"rng.0.source":    "/dev/random",
						"rng.0.max_bytes": "512",
						"rng.0.period":    "200",
					}),
				),
			},
		}},
		{"create virtiofs block", []resource.TestStep{
			{
				Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_hardware_mapping_dir" "test" {
						name      = "{{.DirName}}"

						map = [{
							node = "{{.NodeName}}"
							path = "/mnt"
						}]
					}`, WithRootUser()),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_hardware_mapping_dir.test", map[string]string{
						"name":       dirName,
						"map.0.node": te.NodeName,
						"map.0.path": "/mnt",
					}),
				),
			},
			{
				Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_vm" "test_vm" {
						node_name = "{{.NodeName}}"
						started   = false

						virtiofs {
							mapping = "{{.DirName}}"
							cache = "always"
							direct_io = true
							expose_acl = false
							expose_xattr = false
						}
					}`, WithRootUser()),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_vm.test_vm", map[string]string{
						"virtiofs.0.mapping":      dirName,
						"virtiofs.0.cache":        "always",
						"virtiofs.0.direct_io":    "true",
						"virtiofs.0.expose_acl":   "false",
						"virtiofs.0.expose_xattr": "false",
					}),
				),
			},
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: te.AccProviders,
				Steps:                    tt.step,
			})
		})
	}
}

func TestAccResourceVMInitialization(t *testing.T) {
	te := InitEnvironment(t)

	tests := []struct {
		name string
		step []resource.TestStep
	}{
		{"custom cloud-init: use SCSI interface", []resource.TestStep{{
			Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_file" "cloud_config" {
					content_type = "snippets"
					datastore_id = "local"
					node_name = "{{.NodeName}}"
					source_raw {
						data = <<-EOF
						#cloud-config
						runcmd:
						  - apt update
						  - apt install -y qemu-guest-agent
						  - systemctl enable qemu-guest-agent
						  - systemctl start qemu-guest-agent
						EOF
						file_name = "cloud-config.yaml"
					}
				}

				resource "proxmox_virtual_environment_vm" "test_vm_cloudinit1" {
					node_name = "{{.NodeName}}"
					started   = true
					agent {
						enabled = true
					}
					cpu {
						cores = 2
					}
					memory {
						dedicated = 2048
					}
					disk {
						datastore_id = "local-lvm"
						file_id      = proxmox_virtual_environment_download_file.ubuntu_cloud_image.id
						interface    = "virtio0"
						iothread     = true
						discard      = "on"
						size         = 20
					}

					initialization {
						interface = "scsi1"
						
						ip_config {
							ipv4 {
								address = "dhcp"
							}
						}
						user_data_file_id = proxmox_virtual_environment_file.cloud_config.id
					}
					network_device {
						bridge = "vmbr0"
					}
				}

				resource "proxmox_virtual_environment_download_file" "ubuntu_cloud_image" {
					content_type = "iso"
					datastore_id = "local"
					node_name = "{{.NodeName}}"
					url = "{{.CloudImagesServer}}/jammy/current/jammy-server-cloudimg-amd64.img"
					overwrite_unmanaged = true
				}`),
		}}},
		{"native cloud-init: username should not change", []resource.TestStep{{
			Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_vm_cloudinit4" {
					node_name = "{{.NodeName}}"
					started   = false
					initialization {
						user_account {
							keys = ["ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIOQCHPhOV9XsJa3uq4bmKymklNy6ktgBB/+2umizgnnY"]
						}
					}
				}`),
			Check: NoResourceAttributesSet("proxmox_virtual_environment_vm.test_vm_cloudinit4", []string{
				"initialization.0.username",
				"initialization.0.password",
			}),
		}}},
		{"native cloud-init: username should not change after update", []resource.TestStep{{
			Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_vm_cloudinit4" {
					node_name = "{{.NodeName}}"
					started   = false
					initialization {
						user_account {
							username = "ubuntu"
							password = "password"
						}
					}
				}`),
			Check: ResourceAttributes("proxmox_virtual_environment_vm.test_vm_cloudinit4", map[string]string{
				"initialization.0.user_account.0.username": "ubuntu",
				// override by PVE, set when reading back from the API
				// have to escape the asterisks because of regex match
				"initialization.0.user_account.0.password": `\*\*\*\*\*\*\*\*\*\*`,
			}),
		}, {
			Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_vm_cloudinit4" {
					node_name = "{{.NodeName}}"
					started   = false
					initialization {
						user_account {
							username = "ubuntu"
							password = "password"
						}
						dns {
							servers = ["172.16.0.15", "172.16.0.16"]
							domain = "example.com"
						}
					}
				}`),
			Check: ResourceAttributes("proxmox_virtual_environment_vm.test_vm_cloudinit4", map[string]string{
				"initialization.0.user_account.0.username": "ubuntu",
				"initialization.0.user_account.0.password": `\*\*\*\*\*\*\*\*\*\*`,
			}),
		}}},
		{"native cloud-init: username update should not cause replacement", []resource.TestStep{{
			Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_vm" {
					node_name = "{{.NodeName}}"
					started   = false
					initialization {
						user_account {
							username = "ubuntu"
							password = "password"
						}
					}
				}`),
		}, {
			Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_vm" {
					node_name = "{{.NodeName}}"
					started   = false
					initialization {
						user_account {
							username = "ubuntu-updated"
							password = "password"
						}
					}
				}`),
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction("proxmox_virtual_environment_vm.test_vm", plancheck.ResourceActionUpdate),
				},
			},
		}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: te.AccProviders,
				Steps:                    tt.step,
			})
		})
	}
}

func TestAccResourceVMNetwork(t *testing.T) {
	te := InitEnvironment(t)

	tests := []struct {
		name string
		step []resource.TestStep
	}{
		{"network interfaces", []resource.TestStep{{
			Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_file" "cloud_config" {
					content_type = "snippets"
					datastore_id = "local"
					node_name = "{{.NodeName}}"
					source_raw {
						data = <<-EOF
						#cloud-config
						runcmd:
						  - apt update
						  - apt install -y qemu-guest-agent
						  - systemctl enable qemu-guest-agent
						  - systemctl start qemu-guest-agent
						EOF
						file_name = "cloud-config.yaml"
					}
				}
				
				resource "proxmox_virtual_environment_vm" "test_vm_network1" {
					node_name = "{{.NodeName}}"
					started   = true
					agent {
						enabled = true
					}
					cpu {
						cores = 2
					}
					memory {
						dedicated = 2048
					}
					disk {
						datastore_id = "local-lvm"
						file_id      = proxmox_virtual_environment_download_file.ubuntu_cloud_image.id
						interface    = "virtio0"
						iothread     = true
						discard      = "on"
						size         = 20
					}
					initialization {
						ip_config {
							ipv4 {
								address = "dhcp"
							}
						}
						user_data_file_id = proxmox_virtual_environment_file.cloud_config.id
					}
					network_device {
						bridge = "vmbr0"
						trunks = "10;20;30"
					}
				}

				resource "proxmox_virtual_environment_download_file" "ubuntu_cloud_image" {
					content_type = "iso"
					datastore_id = "local"
					node_name    = "{{.NodeName}}"
					url = "{{.CloudImagesServer}}/jammy/current/jammy-server-cloudimg-amd64.img"
					overwrite_unmanaged = true
				}`),
			Check: resource.ComposeTestCheckFunc(
				ResourceAttributes("proxmox_virtual_environment_vm.test_vm_network1", map[string]string{
					"ipv4_addresses.#":        "2",
					"mac_addresses.#":         "2",
					"network_device.0.bridge": "vmbr0",
					"network_device.0.trunks": "10;20;30",
				}),
			),
		}}},
		{"network device disconnected", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_vm_network2" {
					node_name = "{{.NodeName}}"
					started   = false
					
					network_device {
						bridge = "vmbr0"
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_vm.test_vm_network2", map[string]string{
						"network_device.0.bridge":       "vmbr0",
						"network_device.0.disconnected": "false",
					}),
				),
			}, {
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_vm_network2" {
					node_name = "{{.NodeName}}"
					started   = false
					
					network_device {
						bridge = "vmbr0"
						disconnected = true
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_vm.test_vm_network2", map[string]string{
						"network_device.0.bridge":       "vmbr0",
						"network_device.0.disconnected": "true",
					}),
				),
			},
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: te.AccProviders,
				Steps:                    tt.step,
			})
		})
	}
}

func TestAccResourceVMClone(t *testing.T) {
	if utils.GetAnyStringEnv("TF_ACC") == "" {
		t.Skip("Acceptance tests are disabled")
	}

	te := InitEnvironment(t)

	tests := []struct {
		name string
		step []resource.TestStep
	}{
		{"clone cpu.architecture as root", []resource.TestStep{{
			Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "template" {
					node_name = "{{.NodeName}}"
					started   = false
					template  = true
					cpu {
						architecture = "x86_64"
					}
				}
				resource "proxmox_virtual_environment_vm" "clone" {
					node_name = "{{.NodeName}}"
					started   = false
					clone {
						vm_id = proxmox_virtual_environment_vm.template.vm_id
					}
				}`, WithRootUser()),
		}}},
		{"clone machine type", []resource.TestStep{{
			Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "template" {
					node_name = "{{.NodeName}}"
					started   = false
					template  = true
					machine   = "q35"
				}
				resource "proxmox_virtual_environment_vm" "clone" {
					node_name = "{{.NodeName}}"
					started   = false
					clone {
						vm_id = proxmox_virtual_environment_vm.template.vm_id
					}
					machine = "pc"
				}`),
			Check: ResourceAttributes("proxmox_virtual_environment_vm.clone", map[string]string{
				"machine": "pc",
			}),
		}}},
		{"clone no vga block", []resource.TestStep{{
			Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "template" {
					node_name = "{{.NodeName}}"
					started   = false
				}
				resource "proxmox_virtual_environment_vm" "clone" {
					node_name = "{{.NodeName}}"
					started   = false
					clone {
						vm_id = proxmox_virtual_environment_vm.template.vm_id
					}
				}`),
			Check: resource.ComposeTestCheckFunc(
				ResourceAttributes("proxmox_virtual_environment_vm.clone", map[string]string{
					"vga.#": "0",
				}),
			),
		}}},
		{"clone initialization datastore does not exist", []resource.TestStep{{
			Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "template" {
					node_name = "{{.NodeName}}"
					started   = false
				}
				resource "proxmox_virtual_environment_vm" "clone" {
					node_name = "{{.NodeName}}"
					started   = false
					clone {
						vm_id = proxmox_virtual_environment_vm.template.vm_id
					}
					initialization {
						datastore_id = "doesnotexist"
						ip_config {
							ipv4 {
								address = "172.16.2.57/32"
								gateway = "172.16.2.10"
							}
						}
					}
				}`),
			ExpectError: regexp.MustCompile(`storage 'doesnotexist' does not exist`),
		}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: te.AccProviders,
				Steps:                    tt.step,
			})
		})
	}
}
