/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceVM(t *testing.T) {
	t.Parallel()

	te := InitEnvironment(t)

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
		{
			"protection", []resource.TestStep{{
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
			}},
		},
		{
			"update cpu block", []resource.TestStep{{
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
			}},
		},
		{
			"update memory block", []resource.TestStep{{
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
			}},
		},
		{
			"update vga block", []resource.TestStep{{
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
			}},
		},
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
					url = "https://cloud-images.ubuntu.com/jammy/current/jammy-server-cloudimg-amd64.img"
					overwrite_unmanaged = true
				}`),
		}}},
		{"native cloud-init: do not upgrade packages", []resource.TestStep{{
			Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_vm_cloudinit3" {
					node_name = "{{.NodeName}}"
					started   = false
					initialization {
						upgrade = false
					}
				}`),
			Check: resource.ComposeTestCheckFunc(
				ResourceAttributes("proxmox_virtual_environment_vm.test_vm_cloudinit3", map[string]string{
					"initialization.0.upgrade": "false",
				}),
			),
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
					url = "https://cloud-images.ubuntu.com/jammy/current/jammy-server-cloudimg-amd64.img"
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
		{"network device disconnected", []resource.TestStep{{
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

func TestAccResourceVMDisks(t *testing.T) {
	t.Parallel()

	te := InitEnvironment(t)

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"create disk with default parameters, then update it", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_disk1" {
					node_name = "{{.NodeName}}"
					started   = false
					name 	  = "test-disk1"
					
					disk {
						// note: default qcow2 is not supported by lvm (?)
						file_format  = "raw"
						datastore_id = "local-lvm"
						interface    = "virtio0"
						size         = 8
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_vm.test_disk1", map[string]string{
						"disk.0.aio":               "io_uring",
						"disk.0.backup":            "true",
						"disk.0.cache":             "none",
						"disk.0.discard":           "ignore",
						"disk.0.file_id":           "",
						"disk.0.datastore_id":      "local-lvm",
						"disk.0.file_format":       "raw",
						"disk.0.interface":         "virtio0",
						"disk.0.iothread":          "false",
						"disk.0.path_in_datastore": `vm-\d+-disk-\d+`,
						"disk.0.replicate":         "true",
						"disk.0.size":              "8",
						"disk.0.ssd":               "false",
					}),
				),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_disk1" {
					node_name = "{{.NodeName}}"
					started   = false
					name 	  = "test-disk1"

					disk {
						// note: default qcow2 is not supported by lvm (?)
						file_format  = "raw"
						datastore_id = "local-lvm"
						interface    = "virtio0"
						size         = 8
						replicate    = false
						aio          = "native"
						speed {
						  iops_read = 100
						  iops_read_burstable = 1000
						  iops_write = 400
						  iops_write_burstable = 800
						}
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_vm.test_disk1", map[string]string{
						"disk.0.aio":                          "native",
						"disk.0.backup":                       "true",
						"disk.0.cache":                        "none",
						"disk.0.discard":                      "ignore",
						"disk.0.file_id":                      "",
						"disk.0.datastore_id":                 "local-lvm",
						"disk.0.file_format":                  "raw",
						"disk.0.interface":                    "virtio0",
						"disk.0.iothread":                     "false",
						"disk.0.path_in_datastore":            `vm-\d+-disk-\d+`,
						"disk.0.replicate":                    "false",
						"disk.0.size":                         "8",
						"disk.0.ssd":                          "false",
						"disk.0.speed.0.iops_read":            "100",
						"disk.0.speed.0.iops_read_burstable":  "1000",
						"disk.0.speed.0.iops_write":           "400",
						"disk.0.speed.0.iops_write_burstable": "800",
					}),
				),
			},
		}},
		{"create disk from an image", []resource.TestStep{{
			Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_download_file" "test_disk2_image" {
					content_type = "iso"
					datastore_id = "local"
					node_name    = "{{.NodeName}}"
					url          = "https://cloud-images.ubuntu.com/jammy/current/jammy-server-cloudimg-amd64.img"
					overwrite_unmanaged = true
				}
				resource "proxmox_virtual_environment_vm" "test_disk2" {
					node_name = "{{.NodeName}}"
					started   = false
					name 	  = "test-disk2"	
					disk {
						datastore_id = "local-lvm"
						file_id      = proxmox_virtual_environment_download_file.test_disk2_image.id
						interface    = "virtio0"
						iothread     = true
						discard      = "on"
						size         = 20
					}
				}`),
			Check: resource.ComposeTestCheckFunc(
				ResourceAttributes("proxmox_virtual_environment_vm.test_disk2", map[string]string{
					"disk.0.cache":             "none",
					"disk.0.datastore_id":      "local-lvm",
					"disk.0.discard":           "on",
					"disk.0.file_format":       "raw",
					"disk.0.interface":         "virtio0",
					"disk.0.iothread":          "true",
					"disk.0.path_in_datastore": `vm-\d+-disk-\d+`,
					"disk.0.size":              "20",
					"disk.0.ssd":               "false",
				}),
			),
		}}},
		{"clone default disk without overrides", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_disk3_template" {
					node_name = "{{.NodeName}}"
					started   = false
					name 	  = "test-disk3-template"
					template  = "true"
					
					disk {
						file_format  = "raw"
						datastore_id = "local-lvm"
						interface    = "virtio0"
						size         = 8
					}
				}
				resource "proxmox_virtual_environment_vm" "test_disk3" {
					node_name = "{{.NodeName}}"
					started   = false
					name 	  = "test-disk3"

					clone {
						vm_id = proxmox_virtual_environment_vm.test_disk3_template.id
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					// fully cloned disk, does not have any attributes in state
					resource.TestCheckNoResourceAttr("proxmox_virtual_environment_vm.test_disk3", "disk.0"),
					ResourceAttributes("proxmox_virtual_environment_vm.test_disk3", map[string]string{}),
				),
			},
			{
				RefreshState: true,
			},
		}},
		{"multiple disks", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_disk4" {
					node_name = "{{.NodeName}}"
					started   = false
					name 	  = "test-disk4"
					template  = "true"
					
					disk {
						file_format  = "raw"
						datastore_id = "local-lvm"
						interface    = "virtio0"
						size         = 8
					}
					disk {
						file_format  = "raw"
						datastore_id = "local-lvm"
						interface    = "scsi0"
						size         = 8
					}
				}`),
				Check: ResourceAttributes("proxmox_virtual_environment_vm.test_disk4", map[string]string{
					"disk.0.interface":         "virtio0",
					"disk.0.path_in_datastore": `vm-\d+-disk-1`,
					"disk.1.interface":         "scsi0",
					"disk.1.path_in_datastore": `vm-\d+-disk-0`,
				}),
			},
			{
				RefreshState: true,
			},
		}},

		{"cdrom", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_cdrom" {
					node_name = "{{.NodeName}}"
					started   = false
					name 	  = "test-cdrom"
					cdrom {
						enabled   = true
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_vm.test_cdrom", map[string]string{
						"cdrom.0.enabled": "true",
					}),
				),
			},
			{
				RefreshState: true,
			},
		}},
		{"efi disk", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_efi_disk" {
					node_name = "{{.NodeName}}"
					started   = false
					name 	  = "test-efi-disk"

					efi_disk {
						datastore_id = "local-lvm"
						type = "4m"
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_vm.test_efi_disk", map[string]string{
						"efi_disk.0.datastore_id": "local-lvm",
						"efi_disk.0.type":         "4m",
					}),
				),
			},
			{
				RefreshState: true,
			},
		}},
		{"ide disks", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_disks" {
					node_name = "{{.NodeName}}"
					started   = false
					name 	  = "test-disks-ide"
					
					disk {
						file_format  = "raw"
						datastore_id = "local-lvm"
						interface    = "ide0"
						size         = 8
					}
				}`),
				Check: ResourceAttributes("proxmox_virtual_environment_vm.test_disks", map[string]string{
					"disk.0.interface":         "ide0",
					"disk.0.path_in_datastore": `vm-\d+-disk-0`,
				}),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_disks" {
					node_name = "{{.NodeName}}"
					started   = false
					name 	  = "test-disks-ide"
					
					disk {
						file_format  = "raw"
						datastore_id = "local-lvm"
						interface    = "ide0"
						size         = 8
					}
					disk {
						file_format  = "raw"
						datastore_id = "local-lvm"
						interface    = "ide1"
						size         = 8
					}
				}`),
				Check: ResourceAttributes("proxmox_virtual_environment_vm.test_disks", map[string]string{
					"disk.#": "2",
				}),
			},
			{
				RefreshState: true,
			},
		}},
		{"clone disk with overrides", []resource.TestStep{
			{
				SkipFunc: func() (bool, error) {
					// this test is failing because of https://github.com/bpg/terraform-provider-proxmox/issues/873
					return true, nil
				},
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_disk3_template" {
					node_name = "{{.NodeName}}"
					started   = false
					name 	  = "test-disk3-template"
					template  = "true"
		
					disk {
						file_format  = "raw"
						datastore_id = "local-lvm"
						interface    = "scsi0"
						size         = 8
						discard      = "on"
						iothread     = true
						ssd          = true
					}
				}
				resource "proxmox_virtual_environment_vm" "test_disk3" {
					node_name = "{{.NodeName}}"
					started   = false
					name 	  = "test-disk3"
		
					clone {
						vm_id = proxmox_virtual_environment_vm.test_disk3_template.id
					}
		
					disk {
						interface    = "scsi0"
						//size = 10
					}
				}`),
				Check: ResourceAttributes("proxmox_virtual_environment_vm.test_disk3", map[string]string{
					"disk.0.datastore_id":      "local-lvm",
					"disk.0.discard":           "on",
					"disk.0.file_format":       "raw",
					"disk.0.interface":         "scsi0",
					"disk.0.iothread":          "true",
					"disk.0.path_in_datastore": `vm-\d+-disk-\d+`,
					"disk.0.size":              "8",
					"disk.0.ssd":               "true",
				}),
			},
			{
				RefreshState: true,
				Destroy:      false,
			},
		}},
		{"clone with disk resize", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_disk3_template" {
					node_name = "{{.NodeName}}"
					started   = false
					name 	  = "test-disk3-template"
					template  = "true"
					
					disk {
						file_format  = "raw"
						datastore_id = "local-lvm"
						interface    = "virtio0"
						size         = 8
					}
				}
				resource "proxmox_virtual_environment_vm" "test_disk3" {
					node_name = "{{.NodeName}}"
					started   = false
					name 	  = "test-disk3"

					clone {
						vm_id = proxmox_virtual_environment_vm.test_disk3_template.id
					}

					disk {
						datastore_id = "local-lvm"
						interface    = "virtio0"
						size         = 10
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_vm.test_disk3", map[string]string{
						"disk.0.datastore_id": "local-lvm",
						"disk.0.interface":    "virtio0",
						"disk.0.size":         "10",
					}),
				),
			},
			{
				RefreshState: true,
			},
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: te.AccProviders,
				Steps:                    tt.steps,
			})
		})
	}
}
