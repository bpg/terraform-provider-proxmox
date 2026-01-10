//go:build acceptance || all

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package test

import (
	"fmt"
	"math/rand"
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
		// regression test for https://github.com/bpg/terraform-provider-proxmox/issues/2353
		{"create VM without cpu.units and verify no drift", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_cpu_units" {
					node_name = "{{.NodeName}}"
					started   = false
					cpu {
						cores = 2
					}
				}`),
			},
			{
				RefreshState: true,
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
		{"purge_on_destroy and delete_unreferenced_disks_on_destroy defaults", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_destroy_params" {
					node_name = "{{.NodeName}}"
					started   = false
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_vm.test_destroy_params", map[string]string{
						"purge_on_destroy":                     "true",
						"delete_unreferenced_disks_on_destroy": "true",
					}),
				),
			},
		}},
		{"purge_on_destroy and delete_unreferenced_disks_on_destroy set to false", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_destroy_params_false" {
					node_name = "{{.NodeName}}"
					started   = false
					
					purge_on_destroy                      = false
					delete_unreferenced_disks_on_destroy = false
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_vm.test_destroy_params_false", map[string]string{
						"purge_on_destroy":                     "false",
						"delete_unreferenced_disks_on_destroy": "false",
					}),
				),
			},
		}},
		{"purge_on_destroy and delete_unreferenced_disks_on_destroy update", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_destroy_params_update" {
					node_name = "{{.NodeName}}"
					started   = false
					
					purge_on_destroy                      = true
					delete_unreferenced_disks_on_destroy = true
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_vm.test_destroy_params_update", map[string]string{
						"purge_on_destroy":                     "true",
						"delete_unreferenced_disks_on_destroy": "true",
					}),
				),
			}, {
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_destroy_params_update" {
					node_name = "{{.NodeName}}"
					started   = false
					
					purge_on_destroy                      = false
					delete_unreferenced_disks_on_destroy = false
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_vm.test_destroy_params_update", map[string]string{
						"purge_on_destroy":                     "false",
						"delete_unreferenced_disks_on_destroy": "false",
					}),
				),
			},
		}},
		{"hotplug", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_hotplug" {
					node_name = "{{.NodeName}}"
					started   = false

					hotplug = "disk,usb"
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_vm.test_hotplug", map[string]string{
						"hotplug": "disk,usb",
					}),
				),
			}, {
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_hotplug" {
					node_name = "{{.NodeName}}"
					started   = false

					hotplug = "network,disk,usb,memory"
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_vm.test_hotplug", map[string]string{
						"hotplug": "network,disk,usb,memory",
					}),
				),
			},
		}},
		{"hotplug disabled", []resource.TestStep{{
			Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_hotplug_disabled" {
					node_name = "{{.NodeName}}"
					started   = false

					hotplug = "0"
				}`),
			Check: resource.ComposeTestCheckFunc(
				ResourceAttributes("proxmox_virtual_environment_vm.test_hotplug_disabled", map[string]string{
					"hotplug": "0",
				}),
			),
		}}},
		{"hotplug order insensitive", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_hotplug_order" {
					node_name = "{{.NodeName}}"
					started   = false

					hotplug = "disk,usb,network"
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_vm.test_hotplug_order", map[string]string{
						"hotplug": "disk,usb,network",
					}),
				),
			}, {
				// change order but same values - should not cause update
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_hotplug_order" {
					node_name = "{{.NodeName}}"
					started   = false

					hotplug = "network,disk,usb"
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_vm.test_hotplug_order", map[string]string{
						"hotplug": "disk,usb,network",
					}),
				),
			},
		}},
		{"hotplug duplicate rejected", []resource.TestStep{{
			Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_hotplug_dup" {
					node_name = "{{.NodeName}}"
					started   = false

					hotplug = "cpu,cpu"
				}`),
			ExpectError: regexp.MustCompile(`duplicate hotplug feature "cpu"`),
		}}},
		{"hotplug explicit reset", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_hotplug_reset" {
					node_name = "{{.NodeName}}"
					started   = false

					hotplug = "cpu,disk"
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_vm.test_hotplug_reset", map[string]string{
						"hotplug": "cpu,disk",
					}),
				),
			}, {
				// explicitly set to different value
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_hotplug_reset" {
					node_name = "{{.NodeName}}"
					started   = false

					hotplug = "network,disk,usb"
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_vm.test_hotplug_reset", map[string]string{
						"hotplug": "network,disk,usb",
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

func TestAccResourceVMImport(t *testing.T) {
	te := InitEnvironment(t)

	// Generate dynamic VM ID to avoid conflicts
	testVMID := 100000 + rand.Intn(99999)

	te.AddTemplateVars(map[string]interface{}{
		"TestVMID": testVMID,
	})

	tests := []struct {
		name string
		step []resource.TestStep
	}{
		{"vm import", []resource.TestStep{
			{
				Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_vm" "vm_import" {
						node_name = "{{.NodeName}}"

						vm_id = {{.TestVMID}}

						started   = false
						agent {
							enabled = true
						}
						boot_order = ["virtio0", "net0"]
						cpu {
							cores = 2
						}
						memory {
							dedicated = 2048
						}

						disk {
							datastore_id = "local-lvm"
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
						}
						network_device {
							bridge = "vmbr0"
						}
					}`),
			},
			{
				ResourceName:      "proxmox_virtual_environment_vm.vm_import",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     fmt.Sprintf("%s/%d", te.NodeName, testVMID),
			},
		}},
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

func TestAccResourceVMInitialization(t *testing.T) {
	te := InitEnvironment(t)

	tests := []struct {
		name string
		step []resource.TestStep
	}{
		{"custom cloud-init drive file format", []resource.TestStep{{
			Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_vm_cloudinit" {
					node_name = "{{.NodeName}}"
					started = false
					cpu {
						cores = 1
					}
					memory {
						dedicated = 1024
					}

					initialization {
						datastore_id = "local"
						file_format = "raw"
					}
				}`),
			Check: ResourceAttributes("proxmox_virtual_environment_vm.test_vm_cloudinit", map[string]string{
				"initialization.0.datastore_id": "local",
				"initialization.0.file_format":  "raw",
			}),
		}, {
			Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_vm_cloudinit" {
					node_name = "{{.NodeName}}"
					started = false
					cpu {
						cores = 1
					}
					memory {
						dedicated = 1024
					}

					initialization {
						datastore_id = "local"
						file_format  = "qcow2"
					}
				}`),
			Check: ResourceAttributes("proxmox_virtual_environment_vm.test_vm_cloudinit", map[string]string{
				"initialization.0.datastore_id": "local",
				"initialization.0.file_format":  "qcow2",
			}),
		}}},
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
					url = "{{.CloudImagesServer}}/minimal/releases/noble/release/ubuntu-24.04-minimal-cloudimg-amd64.img"
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
					url = "{{.CloudImagesServer}}/minimal/releases/noble/release/ubuntu-24.04-minimal-cloudimg-amd64.img"
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
		{"wait for IPv4 address", []resource.TestStep{{
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
				
				resource "proxmox_virtual_environment_vm" "test_vm_wait_ipv4" {
					node_name = "{{.NodeName}}"
					started   = true
					agent {
						enabled = true
						wait_for_ip {
							ipv4 = true
						}
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
					}
				}

				resource "proxmox_virtual_environment_download_file" "ubuntu_cloud_image" {
					content_type = "iso"
					datastore_id = "local"
					node_name    = "{{.NodeName}}"
					url = "{{.CloudImagesServer}}/minimal/releases/noble/release/ubuntu-24.04-minimal-cloudimg-amd64.img"
					overwrite_unmanaged = true
				}`),
			Check: resource.ComposeTestCheckFunc(
				ResourceAttributes("proxmox_virtual_environment_vm.test_vm_wait_ipv4", map[string]string{
					"ipv4_addresses.#":           "2",
					"agent.0.wait_for_ip.0.ipv4": "true",
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
		{"remove network device", []resource.TestStep{
			{
				SkipFunc: func() (bool, error) {
					// backward incompatibility with the current implementation of clone
					// see https://github.com/bpg/terraform-provider-proxmox/pull/2260
					return true, nil
				},
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_vm" {
					node_name = "{{.NodeName}}"
					started   = false

					network_device {
						bridge = "vmbr0"
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_vm.test_vm", map[string]string{
						"network_device.#":        "1",
						"network_device.0.bridge": "vmbr0",
					}),
				),
			},
			{
				SkipFunc: func() (bool, error) {
					// backward incompatibility with the current implementation of clone
					// see https://github.com/bpg/terraform-provider-proxmox/pull/2260
					return true, nil
				},
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_vm" {
					node_name = "{{.NodeName}}"
					started   = false
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_vm.test_vm", map[string]string{
						"network_device.#": "0",
					}),
				),
			},
		}},
		{"multiple network devices removal", []resource.TestStep{
			{
				SkipFunc: func() (bool, error) {
					// backward incompatibility with the current implementation of clone
					// see https://github.com/bpg/terraform-provider-proxmox/pull/2260
					return true, nil
				},
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_vm" {
					node_name = "{{.NodeName}}"
					started   = false

					network_device {
						bridge = "vmbr0"
						model  = "virtio"
					}

					network_device {
						bridge = "vmbr1"
						model  = "virtio"
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_vm.test_vm", map[string]string{
						"network_device.#":        "2",
						"network_device.0.bridge": "vmbr0",
						"network_device.1.bridge": "vmbr1",
					}),
				),
			},
			{
				SkipFunc: func() (bool, error) {
					// backward incompatibility with the current implementation of clone
					// see https://github.com/bpg/terraform-provider-proxmox/pull/2260
					return true, nil
				},
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_vm" {
					node_name = "{{.NodeName}}"
					started   = false

					# Only keep the first network device
					network_device {
						bridge = "vmbr0"
						model  = "virtio"
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_vm.test_vm", map[string]string{
						"network_device.#":        "1",
						"network_device.0.bridge": "vmbr0",
					}),
				),
			},
			{
				SkipFunc: func() (bool, error) {
					// backward incompatibility with the current implementation of clone
					// see https://github.com/bpg/terraform-provider-proxmox/pull/2260
					return true, nil
				},
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_vm" {
					node_name = "{{.NodeName}}"
					started   = false
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_vm.test_vm", map[string]string{
						"network_device.#": "0",
					}),
				),
			},
		}},
		{"network device state consistency", []resource.TestStep{
			{
				SkipFunc: func() (bool, error) {
					// backward incompatibility with the current implementation of clone
					// see https://github.com/bpg/terraform-provider-proxmox/pull/2260
					return true, nil
				},
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_vm" {
					node_name = "{{.NodeName}}"
					started   = false

					network_device {
						bridge = "vmbr0"
						model  = "virtio"
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_vm.test_vm", map[string]string{
						"network_device.#":        "1",
						"network_device.0.bridge": "vmbr0",
						"network_device.0.model":  "virtio",
					}),
				),
			},
			{
				SkipFunc: func() (bool, error) {
					// backward incompatibility with the current implementation of clone
					// see https://github.com/bpg/terraform-provider-proxmox/pull/2260
					return true, nil
				},
				// This step tests that the state is read correctly after network device removal
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_vm" {
					node_name = "{{.NodeName}}"
					started   = false
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_vm.test_vm", map[string]string{
						"network_device.#": "0",
					}),
				),
			},
			{
				SkipFunc: func() (bool, error) {
					// backward incompatibility with the current implementation of clone
					// see https://github.com/bpg/terraform-provider-proxmox/pull/2260
					return true, nil
				},
				// This step tests that we can add network devices back after removal
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_vm" {
					node_name = "{{.NodeName}}"
					started   = false

					network_device {
						bridge = "vmbr0"
						model  = "virtio"
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_vm.test_vm", map[string]string{
						"network_device.#":        "1",
						"network_device.0.bridge": "vmbr0",
						"network_device.0.model":  "virtio",
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
		{"clone with network device", []resource.TestStep{{
			Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "template" {
					node_name = "{{.NodeName}}"
					started   = false
					template = true
					network_device {
						bridge = "vmbr0"
					}
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
					"network_device.#":        "1",
					"network_device.0.bridge": "vmbr0",
				}),
			),
		}}},
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
		{"clone with network devices", []resource.TestStep{{
			Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "template" {
					node_name = "{{.NodeName}}"
					started   = false
					network_device {
						bridge = "vmbr0"
					}
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
					"network_device.#":        "1",
					"network_device.0.bridge": "vmbr0",
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
		{"clone hotplug inherited", []resource.TestStep{{
			Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "template_hotplug" {
					node_name = "{{.NodeName}}"
					started   = false
					template  = true
					hotplug   = "cpu,disk"
				}
				resource "proxmox_virtual_environment_vm" "clone_hotplug_inherit" {
					node_name = "{{.NodeName}}"
					started   = false
					clone {
						vm_id = proxmox_virtual_environment_vm.template_hotplug.vm_id
					}
				}`),
			Check: resource.ComposeTestCheckFunc(
				ResourceAttributes("proxmox_virtual_environment_vm.clone_hotplug_inherit", map[string]string{
					"hotplug": "cpu,disk",
				}),
			),
		}}},
		{"clone hotplug override", []resource.TestStep{{
			Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "template_hotplug2" {
					node_name = "{{.NodeName}}"
					started   = false
					template  = true
					hotplug   = "cpu,disk"
				}
				resource "proxmox_virtual_environment_vm" "clone_hotplug_override" {
					node_name = "{{.NodeName}}"
					started   = false
					clone {
						vm_id = proxmox_virtual_environment_vm.template_hotplug2.vm_id
					}
					hotplug = "network,usb"
				}`),
			Check: resource.ComposeTestCheckFunc(
				ResourceAttributes("proxmox_virtual_environment_vm.clone_hotplug_override", map[string]string{
					"hotplug": "network,usb",
				}),
			),
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

func TestAccResourceVMVirtioSCSISingleWithAgent(t *testing.T) {
	te := InitEnvironment(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
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

				resource "proxmox_virtual_environment_vm" "test_vm_scsi_single" {
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
						interface    = "scsi0"
						iothread     = true
						discard      = "on"
						size         = 20
					}
					scsi_hardware = "virtio-scsi-single"
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
					node_name    = "{{.NodeName}}"
					url = "{{.CloudImagesServer}}/minimal/releases/noble/release/ubuntu-24.04-minimal-cloudimg-amd64.img"
					overwrite_unmanaged = true
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_vm.test_vm_scsi_single", map[string]string{
						"scsi_hardware":             "virtio-scsi-single",
						"agent.0.enabled":           "true",
						"ipv4_addresses.#":          "2",
						"network_interface_names.#": "2",
					}),
				),
			},
		},
	})
}
