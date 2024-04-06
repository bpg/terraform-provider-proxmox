/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceVM(t *testing.T) {
	t.Parallel()

	te := initTestEnvironment(t)

	tests := []struct {
		name string
		step []resource.TestStep
	}{
		{"multiline description", []resource.TestStep{{
			Config: te.providerConfig + fmt.Sprintf(`
				resource "proxmox_virtual_environment_vm" "test_vm1" {
					node_name = "%s"
					started   = false
					
					description = <<-EOT
						my
						description
						value
					EOT
				}`, te.nodeName),
			Check: resource.ComposeTestCheckFunc(
				testResourceAttributes("proxmox_virtual_environment_vm.test_vm1", map[string]string{
					"description": "my\ndescription\nvalue",
				}),
			),
		}}},
		{"single line description", []resource.TestStep{{
			Config: te.providerConfig + fmt.Sprintf(`
				resource "proxmox_virtual_environment_vm" "test_vm2" {
					node_name = "%s"
					started   = false
					
					description = "my description value"
				}`, te.nodeName),
			Check: resource.ComposeTestCheckFunc(
				testResourceAttributes("proxmox_virtual_environment_vm.test_vm2", map[string]string{
					"description": "my description value",
				}),
			),
		}}},
		{"no description", []resource.TestStep{{
			Config: fmt.Sprintf(`
				resource "proxmox_virtual_environment_vm" "test_vm3" {
					node_name = "%s"
					started   = false
					
					description = ""
				}`, te.nodeName),
			Check: resource.ComposeTestCheckFunc(
				testResourceAttributes("proxmox_virtual_environment_vm.test_vm3", map[string]string{
					"description": "",
				}),
			),
		}}},
		{
			"protection", []resource.TestStep{{
				Config: fmt.Sprintf(`
				resource "proxmox_virtual_environment_vm" "test_vm4" {
					node_name = "%s"
					started   = false
					
					protection = true
				}`, te.nodeName),
				Check: resource.ComposeTestCheckFunc(
					testResourceAttributes("proxmox_virtual_environment_vm.test_vm4", map[string]string{
						"protection": "true",
					}),
				),
			}, {
				Config: fmt.Sprintf(`
				resource "proxmox_virtual_environment_vm" "test_vm4" {
					node_name = "%s"
					started   = false
					
					protection = false
				}`, te.nodeName),
				Check: resource.ComposeTestCheckFunc(
					testResourceAttributes("proxmox_virtual_environment_vm.test_vm4", map[string]string{
						"protection": "false",
					}),
				),
			}},
		},
		{
			"update cpu block", []resource.TestStep{{
				Config: fmt.Sprintf(`resource "proxmox_virtual_environment_vm" "test_vm5" {
					node_name = "%s"
					started   = false
					
					cpu {
						cores = 2
					}
				}`, te.nodeName),
				Check: resource.ComposeTestCheckFunc(
					testResourceAttributes("proxmox_virtual_environment_vm.test_vm5", map[string]string{
						"cpu.0.sockets": "1",
					}),
				),
			}, {
				Config: fmt.Sprintf(`resource "proxmox_virtual_environment_vm" "test_vm5" {
					node_name = "%s"
					started   = false
					
					cpu {
						cores = 1
					}
				}`, te.nodeName),
				Check: resource.ComposeTestCheckFunc(
					testResourceAttributes("proxmox_virtual_environment_vm.test_vm5", map[string]string{
						"cpu.0.sockets": "1",
					}),
				),
			}},
		},
		{
			"update memory block", []resource.TestStep{{
				Config: fmt.Sprintf(`resource "proxmox_virtual_environment_vm" "test_vm6" {
					node_name = "%s"
					started   = false
					
					memory {
						dedicated = 2048
					}
				}`, te.nodeName),
				Check: resource.ComposeTestCheckFunc(
					testResourceAttributes("proxmox_virtual_environment_vm.test_vm6", map[string]string{
						"memory.0.dedicated": "2048",
					}),
				),
			}, {
				Config: fmt.Sprintf(`resource "proxmox_virtual_environment_vm" "test_vm6" {
					node_name = "%s"
					started   = false
					
					memory {
						dedicated = 1024
					}
				}`, te.nodeName),
				Check: resource.ComposeTestCheckFunc(
					testResourceAttributes("proxmox_virtual_environment_vm.test_vm6", map[string]string{
						"memory.0.dedicated": "1024",
					}),
				),
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: te.accProviders,
				Steps:                    tt.step,
			})
		})
	}
}

func TestAccResourceVMInitialization(t *testing.T) {
	te := initTestEnvironment(t)

	tests := []struct {
		name string
		step []resource.TestStep
	}{
		{"initialization works with cloud-init config provided over SCSI interface", []resource.TestStep{{
			Config: te.providerConfig + fmt.Sprintf(`
				resource "proxmox_virtual_environment_file" "cloud_config" {
					content_type = "snippets"
					datastore_id = "local"
					node_name    = "%[1]s"
					source_raw {
						data = <<EOF
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
					node_name = "%[1]s"
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
					node_name    = "%[1]s"
					url = "https://cloud-images.ubuntu.com/jammy/current/jammy-server-cloudimg-amd64.img"
					overwrite_unmanaged = true
				}`, te.nodeName),
		}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: te.accProviders,
				Steps:                    tt.step,
			})
		})
	}
}

func TestAccResourceVMNetwork(t *testing.T) {
	te := initTestEnvironment(t)

	tests := []struct {
		name string
		step []resource.TestStep
	}{
		{"network interfaces", []resource.TestStep{{
			Config: te.providerConfig + fmt.Sprintf(`
				resource "proxmox_virtual_environment_file" "cloud_config" {
					content_type = "snippets"
					datastore_id = "local"
					node_name    = "%[1]s"
					source_raw {
						data = <<EOF
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
					node_name = "%[1]s"
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
					node_name    = "%[1]s"
					url = "https://cloud-images.ubuntu.com/jammy/current/jammy-server-cloudimg-amd64.img"
					overwrite_unmanaged = true
				}`, te.nodeName),
			Check: resource.ComposeTestCheckFunc(
				testResourceAttributes("proxmox_virtual_environment_vm.test_vm_network1", map[string]string{
					"ipv4_addresses.#":        "2",
					"mac_addresses.#":         "2",
					"network_device.0.bridge": "vmbr0",
					"network_device.0.trunks": "10;20;30",
				}),
			),
		}}},
		{"network device disconnected", []resource.TestStep{{
			Config: fmt.Sprintf(`
				resource "proxmox_virtual_environment_vm" "test_vm_network2" {
					node_name = "%s"
					started   = false
					
					network_device {
						bridge = "vmbr0"
					}
				}`, te.nodeName),
			Check: resource.ComposeTestCheckFunc(
				testResourceAttributes("proxmox_virtual_environment_vm.test_vm_network2", map[string]string{
					"network_device.0.bridge":       "vmbr0",
					"network_device.0.disconnected": "false",
				}),
			),
		}, {
			Config: fmt.Sprintf(`
				resource "proxmox_virtual_environment_vm" "test_vm_network2" {
					node_name = "%s"
					started   = false
					
					network_device {
						bridge = "vmbr0"
						disconnected = true
					}
				}`, te.nodeName),
			Check: resource.ComposeTestCheckFunc(
				testResourceAttributes("proxmox_virtual_environment_vm.test_vm_network2", map[string]string{
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
				ProtoV6ProviderFactories: te.accProviders,
				Steps:                    tt.step,
			})
		})
	}
}

func TestAccResourceVMDisks(t *testing.T) {
	t.Parallel()

	te := initTestEnvironment(t)

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"create disk with default parameters, then update it", []resource.TestStep{
			{
				Config: te.providerConfig + fmt.Sprintf(`
				resource "proxmox_virtual_environment_vm" "test_disk1" {
					node_name = "%s"
					started   = false
					name 	  = "test-disk1"
					
					disk {
						// note: default qcow2 is not supported by lvm (?)
						file_format  = "raw"
						datastore_id = "local-lvm"
						interface    = "virtio0"
						size         = 8
					}
				}`, te.nodeName),
				Check: resource.ComposeTestCheckFunc(
					testResourceAttributes("proxmox_virtual_environment_vm.test_disk1", map[string]string{
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
				Config: te.providerConfig + fmt.Sprintf(`
				resource "proxmox_virtual_environment_vm" "test_disk1" {
					node_name = "%s"
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
				}`, te.nodeName),
				Check: resource.ComposeTestCheckFunc(
					testResourceAttributes("proxmox_virtual_environment_vm.test_disk1", map[string]string{
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
			Config: te.providerConfig + fmt.Sprintf(`
				resource "proxmox_virtual_environment_download_file" "test_disk2_image" {
					content_type = "iso"
					datastore_id = "local"
					node_name    = "%[1]s"
					url          = "https://cloud-images.ubuntu.com/jammy/current/jammy-server-cloudimg-amd64.img"
					overwrite_unmanaged = true
				}
				resource "proxmox_virtual_environment_vm" "test_disk2" {
					node_name = "%[1]s"
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
				}`, te.nodeName),
			Check: resource.ComposeTestCheckFunc(
				testResourceAttributes("proxmox_virtual_environment_vm.test_disk2", map[string]string{
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
				Config: te.providerConfig + fmt.Sprintf(`
				resource "proxmox_virtual_environment_vm" "test_disk3_template" {
					node_name = "%[1]s"
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
					node_name = "%[1]s"
					started   = false
					name 	  = "test-disk3"

					clone {
						vm_id = proxmox_virtual_environment_vm.test_disk3_template.id
					}
				}
				`, te.nodeName),
				Check: resource.ComposeTestCheckFunc(
					// fully cloned disk, does not have any attributes in state
					resource.TestCheckNoResourceAttr("proxmox_virtual_environment_vm.test_disk3", "disk.0"),
					testResourceAttributes("proxmox_virtual_environment_vm.test_disk3", map[string]string{}),
				),
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
				Config: te.providerConfig + fmt.Sprintf(`
				resource "proxmox_virtual_environment_vm" "test_disk3_template" {
					node_name = "%[1]s"
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
					node_name = "%[1]s"
					started   = false
					name 	  = "test-disk3"
		
					clone {
						vm_id = proxmox_virtual_environment_vm.test_disk3_template.id
					}
		
					disk {
						interface    = "scsi0"
						//size = 10
					}
				}
				`, te.nodeName),
				Check: resource.ComposeTestCheckFunc(
					testResourceAttributes("proxmox_virtual_environment_vm.test_disk3", map[string]string{
						"disk.0.datastore_id":      "local-lvm",
						"disk.0.discard":           "on",
						"disk.0.file_format":       "raw",
						"disk.0.interface":         "scsi0",
						"disk.0.iothread":          "true",
						"disk.0.path_in_datastore": `vm-\d+-disk-\d+`,
						"disk.0.size":              "8",
						"disk.0.ssd":               "true",
					}),
				),
			},
			{
				RefreshState: true,
				Destroy:      false,
			},
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: te.accProviders,
				Steps:                    tt.steps,
			})
		})
	}
}
