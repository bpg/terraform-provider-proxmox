/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package tests

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceVM(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		step resource.TestStep
	}{
		{"multiline description", resource.TestStep{
			Config: `
				resource "proxmox_virtual_environment_vm" "test_vm1" {
					node_name = "pve"
					started   = false
					
					description = <<-EOT
						my
						description
						value
					EOT
				}`,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("proxmox_virtual_environment_vm.test_vm1", "description", "my\ndescription\nvalue"),
			),
		}},
		{"single line description", resource.TestStep{
			Config: `
				resource "proxmox_virtual_environment_vm" "test_vm2" {
					node_name = "pve"
					started   = false
					
					description = "my description value"
				}`,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("proxmox_virtual_environment_vm.test_vm2", "description", "my description value"),
			),
		}},
		{"no description", resource.TestStep{
			Config: `
				resource "proxmox_virtual_environment_vm" "test_vm3" {
					node_name = "pve"
					started   = false
					
					description = ""
				}`,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("proxmox_virtual_environment_vm.test_vm3", "description", ""),
			),
		}},
	}

	accProviders := testAccMuxProviders(context.Background(), t)

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: accProviders,
				Steps:                    []resource.TestStep{tt.step},
			})
		})
	}
}

//nolint:paralleltest
func TestAccResourceVMNetwork(t *testing.T) {
	tests := []struct {
		name string
		step resource.TestStep
	}{
		{"network interfaces mac", resource.TestStep{
			Config: `
				resource "proxmox_virtual_environment_file" "cloud_config" {
					content_type = "snippets"
					datastore_id = "local"
					node_name    = "pve"
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
					node_name = "pve"
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
					}
				}

				resource "proxmox_virtual_environment_download_file" "ubuntu_cloud_image" {
					content_type = "iso"
					datastore_id = "local"
					node_name    = "pve"
					url = "https://cloud-images.ubuntu.com/jammy/current/jammy-server-cloudimg-amd64.img"
				}`,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("proxmox_virtual_environment_vm.test_vm_network1", "ipv4_addresses.#", "2"),
				resource.TestCheckResourceAttr("proxmox_virtual_environment_vm.test_vm_network1", "mac_addresses.#", "2"),
			),
		}},
	}

	accProviders := testAccMuxProviders(context.Background(), t)

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: accProviders,
				Steps:                    []resource.TestStep{tt.step},
			})
		})
	}
}

func TestAccResourceVMDisks(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"create disk with default parameters", []resource.TestStep{{
			Config: `
				resource "proxmox_virtual_environment_vm" "test_disk1" {
					node_name = "pve"
					started   = false
					name 	  = "test-disk1"
					
					disk {
						// note: default qcow2 is not supported by lvm (?)
						file_format  = "raw"
						datastore_id = "local-lvm"
						interface    = "virtio0"
						size         = 8
					}
				}`,
			Check: resource.ComposeTestCheckFunc(
				testResourceAttributes("proxmox_virtual_environment_vm.test_disk1", map[string]string{
					// those are empty by default, but we can't check for that
					// "disk.0.cache":          "",
					// "disk.0.discard":      	"",
					// "disk.0.file_id":        "",
					"disk.0.datastore_id":      "local-lvm",
					"disk.0.file_format":       "raw",
					"disk.0.interface":         "virtio0",
					"disk.0.iothread":          "false",
					"disk.0.path_in_datastore": `vm-\d+-disk-\d+`,
					"disk.0.size":              "8",
					"disk.0.ssd":               "false",
				}),
			),
		}}},
		{"create disk from an image", []resource.TestStep{{
			Config: `
				resource "proxmox_virtual_environment_download_file" "test_disk2_image" {
					content_type = "iso"
					datastore_id = "local"
					node_name    = "pve"
					url          = "https://cloud-images.ubuntu.com/jammy/current/jammy-server-cloudimg-amd64.img"
				}
				resource "proxmox_virtual_environment_vm" "test_disk2" {
					node_name = "pve"
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
				}`,
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
				Config: `
				resource "proxmox_virtual_environment_vm" "test_disk3_template" {
					node_name = "pve"
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
					node_name = "pve"
					started   = false
					name 	  = "test-disk3"

					clone {
						vm_id = proxmox_virtual_environment_vm.test_disk3_template.id
					}
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					// fully cloned disk, does not have any attributes in state
					resource.TestCheckNoResourceAttr("proxmox_virtual_environment_vm.test_disk3", "disk.0"),
				),
			},
			{
				RefreshState: true,
			},
		}},
		// this test is failing because of https://github.com/bpg/terraform-provider-proxmox/issues/360
		// {"clone disk with new size", []resource.TestStep{
		//	{
		//		Config: `
		//		resource "proxmox_virtual_environment_vm" "test_disk3_template" {
		//			node_name = "pve"
		//			started   = false
		//			name 	  = "test-disk3-template"
		//			template  = "true"
		//
		//			disk {
		//				file_format  = "raw"
		//				datastore_id = "local-lvm"
		//				interface    = "scsi0"
		//				size         = 8
		//				discard      = "on"
		//				iothread     = true
		//			}
		//		}
		//		resource "proxmox_virtual_environment_vm" "test_disk3" {
		//			node_name = "pve"
		//			started   = false
		//			name 	  = "test-disk3"
		//
		//			clone {
		//				vm_id = proxmox_virtual_environment_vm.test_disk3_template.id
		//			}
		//
		//			disk {
		//				interface    = "scsi0"
		//				size = 10
		//                //ssd = true
		//			}
		//		}
		//		`,
		//		Check: resource.ComposeTestCheckFunc(
		//			testResourceAttributes("proxmox_virtual_environment_vm.test_disk3", map[string]string{
		//				"disk.0.datastore_id":      "local-lvm",
		//				"disk.0.discard":           "on",
		//				"disk.0.file_format":       "raw",
		//				"disk.0.interface":         "scsi0",
		//				"disk.0.iothread":          "true",
		//				"disk.0.path_in_datastore": `vm-\d+-disk-\d+`,
		//				"disk.0.size":              "10",
		//				"disk.0.ssd":               "false",
		//			}),
		//		),
		//	},
		//{
		//	RefreshState: true,
		//	Destroy:      false,
		// },
		// }},
	}

	accProviders := testAccMuxProviders(context.Background(), t)

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: accProviders,
				Steps:                    tt.steps,
			})
		})
	}
}

func testResourceAttributes(res string, attrs map[string]string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for k, v := range attrs {
			if err := resource.TestCheckResourceAttrWith(res, k, func(got string) error {
				match, err := regexp.Match(v, []byte(got)) //nolint:mirror
				if err != nil {
					return fmt.Errorf("error matching '%s': %w", v, err)
				}
				if !match {
					return fmt.Errorf("expected '%s' to match '%s'", got, v)
				}
				return nil
			})(s); err != nil {
				return err
			}
		}

		return nil
	}
}
