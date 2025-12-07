//go:build acceptance || all

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestAccResourceVMHotplug(t *testing.T) {
	te := InitEnvironment(t)

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"add disk to running VM", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_download_file" "ubuntu_cloud_image" {
					content_type = "iso"
					datastore_id = "local"
					node_name    = "{{.NodeName}}"
					url = "{{.CloudImagesServer}}/minimal/releases/noble/release/ubuntu-24.04-minimal-cloudimg-amd64.img"
					overwrite_unmanaged = true
				}

				resource "proxmox_virtual_environment_vm" "test_hotplug" {
					node_name       = "{{.NodeName}}"
					started         = true
					stop_on_destroy = true
					name            = "test-hotplug-disk"
					
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
						size         = 20
					}
					initialization {
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
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_vm.test_hotplug", map[string]string{
						"disk.#": "1",
					}),
				),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_download_file" "ubuntu_cloud_image" {
					content_type = "iso"
					datastore_id = "local"
					node_name    = "{{.NodeName}}"
					url = "{{.CloudImagesServer}}/minimal/releases/noble/release/ubuntu-24.04-minimal-cloudimg-amd64.img"
					overwrite_unmanaged = true
				}

				resource "proxmox_virtual_environment_vm" "test_hotplug" {
					node_name       = "{{.NodeName}}"
					started         = true
					stop_on_destroy = true
					name            = "test-hotplug-disk"
					
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
						size         = 20
					}
					disk {
						datastore_id = "local-lvm"
						interface    = "scsi1"
						size         = 4
					}
					initialization {
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
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_vm.test_hotplug", map[string]string{
						"disk.#": "2",
					}),
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("proxmox_virtual_environment_vm.test_hotplug", plancheck.ResourceActionUpdate),
					},
				},
			},
		}},
		{"add network device to running VM", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_download_file" "ubuntu_cloud_image" {
					content_type = "iso"
					datastore_id = "local"
					node_name    = "{{.NodeName}}"
					url = "{{.CloudImagesServer}}/minimal/releases/noble/release/ubuntu-24.04-minimal-cloudimg-amd64.img"
					overwrite_unmanaged = true
				}

				resource "proxmox_virtual_environment_vm" "test_hotplug" {
					node_name = "{{.NodeName}}"
					started   = true
					name      = "test-hotplug-network"
					
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
						size         = 20
					}
					initialization {
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
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_vm.test_hotplug", map[string]string{
						"network_device.#": "1",
					}),
				),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_download_file" "ubuntu_cloud_image" {
					content_type = "iso"
					datastore_id = "local"
					node_name    = "{{.NodeName}}"
					url = "{{.CloudImagesServer}}/minimal/releases/noble/release/ubuntu-24.04-minimal-cloudimg-amd64.img"
					overwrite_unmanaged = true
				}

				resource "proxmox_virtual_environment_vm" "test_hotplug" {
					node_name = "{{.NodeName}}"
					started   = true
					name      = "test-hotplug-network"
					
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
						size         = 20
					}
					initialization {
						ip_config {
							ipv4 {
								address = "dhcp"
							}
						}
					}
					network_device {
						bridge = "vmbr0"
					}
					network_device {
						bridge = "vmbr0"
						model  = "virtio"
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_vm.test_hotplug", map[string]string{
						"network_device.#": "2",
					}),
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("proxmox_virtual_environment_vm.test_hotplug", plancheck.ResourceActionUpdate),
					},
				},
			},
		}},
		{"increase memory on running VM", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_download_file" "ubuntu_cloud_image" {
					content_type = "iso"
					datastore_id = "local"
					node_name    = "{{.NodeName}}"
					url = "{{.CloudImagesServer}}/minimal/releases/noble/release/ubuntu-24.04-minimal-cloudimg-amd64.img"
					overwrite_unmanaged = true
				}

				resource "proxmox_virtual_environment_vm" "test_hotplug" {
					node_name = "{{.NodeName}}"
					started   = true
					name      = "test-hotplug-memory"
					
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
						size         = 20
					}
					initialization {
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
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_vm.test_hotplug", map[string]string{
						"memory.0.dedicated": "2048",
					}),
				),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_download_file" "ubuntu_cloud_image" {
					content_type = "iso"
					datastore_id = "local"
					node_name    = "{{.NodeName}}"
					url = "{{.CloudImagesServer}}/minimal/releases/noble/release/ubuntu-24.04-minimal-cloudimg-amd64.img"
					overwrite_unmanaged = true
				}

				resource "proxmox_virtual_environment_vm" "test_hotplug" {
					node_name = "{{.NodeName}}"
					started   = true
					name      = "test-hotplug-memory"
					
					cpu {
						cores = 2
					}
					memory {
						dedicated = 4096
					}
					disk {
						datastore_id = "local-lvm"
						file_id      = proxmox_virtual_environment_download_file.ubuntu_cloud_image.id
						interface    = "scsi0"
						size         = 20
					}
					initialization {
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
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_vm.test_hotplug", map[string]string{
						"memory.0.dedicated": "4096",
					}),
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("proxmox_virtual_environment_vm.test_hotplug", plancheck.ResourceActionUpdate),
					},
				},
			},
		}},
		{"increase CPU cores on running VM", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_download_file" "ubuntu_cloud_image" {
					content_type = "iso"
					datastore_id = "local"
					node_name    = "{{.NodeName}}"
					url = "{{.CloudImagesServer}}/minimal/releases/noble/release/ubuntu-24.04-minimal-cloudimg-amd64.img"
					overwrite_unmanaged = true
				}

				resource "proxmox_virtual_environment_vm" "test_hotplug" {
					node_name = "{{.NodeName}}"
					started   = true
					name      = "test-hotplug-cpu"
					
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
						size         = 20
					}
					initialization {
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
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_vm.test_hotplug", map[string]string{
						"cpu.0.cores": "2",
					}),
				),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_download_file" "ubuntu_cloud_image" {
					content_type = "iso"
					datastore_id = "local"
					node_name    = "{{.NodeName}}"
					url = "{{.CloudImagesServer}}/minimal/releases/noble/release/ubuntu-24.04-minimal-cloudimg-amd64.img"
					overwrite_unmanaged = true
				}

				resource "proxmox_virtual_environment_vm" "test_hotplug" {
					node_name = "{{.NodeName}}"
					started   = true
					name      = "test-hotplug-cpu"
					
					cpu {
						cores = 4
					}
					memory {
						dedicated = 2048
					}
					disk {
						datastore_id = "local-lvm"
						file_id      = proxmox_virtual_environment_download_file.ubuntu_cloud_image.id
						interface    = "scsi0"
						size         = 20
					}
					initialization {
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
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_vm.test_hotplug", map[string]string{
						"cpu.0.cores": "4",
					}),
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("proxmox_virtual_environment_vm.test_hotplug", plancheck.ResourceActionUpdate),
					},
				},
			},
		}},
		{"change disk properties on running VM", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_download_file" "ubuntu_cloud_image" {
					content_type = "iso"
					datastore_id = "local"
					node_name    = "{{.NodeName}}"
					url = "{{.CloudImagesServer}}/minimal/releases/noble/release/ubuntu-24.04-minimal-cloudimg-amd64.img"
					overwrite_unmanaged = true
				}

				resource "proxmox_virtual_environment_vm" "test_hotplug" {
					node_name = "{{.NodeName}}"
					started   = true
					name      = "test-hotplug-disk-props"
					
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
						size         = 20
						cache        = "none"
					}
					initialization {
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
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_vm.test_hotplug", map[string]string{
						"disk.0.cache": "none",
					}),
				),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_download_file" "ubuntu_cloud_image" {
					content_type = "iso"
					datastore_id = "local"
					node_name    = "{{.NodeName}}"
					url = "{{.CloudImagesServer}}/minimal/releases/noble/release/ubuntu-24.04-minimal-cloudimg-amd64.img"
					overwrite_unmanaged = true
				}

				resource "proxmox_virtual_environment_vm" "test_hotplug" {
					node_name = "{{.NodeName}}"
					started   = true
					name      = "test-hotplug-disk-props"
					
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
						size         = 20
						cache        = "writeback"
						discard      = "on"
					}
					initialization {
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
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_vm.test_hotplug", map[string]string{
						"disk.0.cache":   "writeback",
						"disk.0.discard": "on",
					}),
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("proxmox_virtual_environment_vm.test_hotplug", plancheck.ResourceActionUpdate),
					},
				},
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
