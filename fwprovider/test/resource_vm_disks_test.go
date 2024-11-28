//go:build acceptance || all

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

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
				resource "proxmox_virtual_environment_vm" "test_disk" {
					node_name = "{{.NodeName}}"
					started   = false
					name 	  = "test-disk"
					
					disk {
						// note: default qcow2 is not supported by lvm (?)
						file_format  = "raw"
						datastore_id = "local-lvm"
						interface    = "virtio0"
						size         = 8
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_vm.test_disk", map[string]string{
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
				resource "proxmox_virtual_environment_vm" "test_disk" {
					node_name = "{{.NodeName}}"
					started   = false
					name 	  = "test-disk"

					disk {
						// note: default qcow2 is not supported by lvm (?)
						file_format  = "raw"
						datastore_id = "local-lvm"
						interface    = "virtio0"
						serial	     = "-dead_beef-"
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
					ResourceAttributes("proxmox_virtual_environment_vm.test_disk", map[string]string{
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
						"disk.0.serial":                       "-dead_beef-",
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
				resource "proxmox_virtual_environment_download_file" "test_disk_image" {
					content_type = "iso"
					datastore_id = "local"
					node_name    = "{{.NodeName}}"
					url          = "{{.CloudImagesServer}}/jammy/current/jammy-server-cloudimg-amd64.img"
					overwrite_unmanaged = true
				}
				resource "proxmox_virtual_environment_vm" "test_disk" {
					node_name = "{{.NodeName}}"
					started   = false
					name 	  = "test-disk"	
					disk {
						datastore_id = "local-lvm"
						file_id      = proxmox_virtual_environment_download_file.test_disk_image.id
						interface    = "virtio0"
						iothread     = true
						discard      = "on"
						serial       = "dead_beef"
						size         = 20
					}
				}`),
			Check: resource.ComposeTestCheckFunc(
				ResourceAttributes("proxmox_virtual_environment_vm.test_disk", map[string]string{
					"disk.0.cache":             "none",
					"disk.0.datastore_id":      "local-lvm",
					"disk.0.discard":           "on",
					"disk.0.file_format":       "raw",
					"disk.0.interface":         "virtio0",
					"disk.0.iothread":          "true",
					"disk.0.path_in_datastore": `vm-\d+-disk-\d+`,
					"disk.0.serial":            "dead_beef",
					"disk.0.size":              "20",
					"disk.0.ssd":               "false",
				}),
			),
		}}},
		{"clone default disk without overrides", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_disk_template" {
					node_name = "{{.NodeName}}"
					started   = false
					name 	  = "test-disk-template"
					template  = "true"
					
					disk {
						file_format  = "raw"
						datastore_id = "local-lvm"
						interface    = "virtio0"
						size         = 8
					}
				}
				resource "proxmox_virtual_environment_vm" "test_disk" {
					node_name = "{{.NodeName}}"
					started   = false
					name 	  = "test-disk"

					clone {
						vm_id = proxmox_virtual_environment_vm.test_disk_template.id
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					// fully cloned disk, does not have any attributes in state
					resource.TestCheckNoResourceAttr("proxmox_virtual_environment_vm.test_disk", "disk.0"),
					ResourceAttributes("proxmox_virtual_environment_vm.test_disk", map[string]string{}),
				),
			},
			{
				RefreshState: true,
			},
		}},
		{"multiple disks", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_disk" {
					node_name = "{{.NodeName}}"
					started   = false
					name 	  = "test-disk"
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
				Check: ResourceAttributes("proxmox_virtual_environment_vm.test_disk", map[string]string{
					"disk.0.interface":         "virtio0",
					"disk.0.path_in_datastore": `base-\d+-disk-1`,
					"disk.1.interface":         "scsi0",
					"disk.1.path_in_datastore": `base-\d+-disk-0`,
				}),
			},
			{
				RefreshState: true,
			},
		}},
		{"adding disks", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_disk" {
					node_name = "{{.NodeName}}"
					started   = false
					name 	  = "test-disk"
					
					disk {
						file_format  = "raw"
						datastore_id = "local-lvm"
						interface    = "scsi0"
						size         = 8
					}
				}`),
				Check: ResourceAttributes("proxmox_virtual_environment_vm.test_disk", map[string]string{
					"disk.0.interface":         "scsi0",
					"disk.0.path_in_datastore": `vm-\d+-disk-0`,
				}),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_disk" {
					node_name = "{{.NodeName}}"
					started   = false
					name 	  = "test-disk"
					
					disk {
						file_format  = "raw"
						datastore_id = "local-lvm"
						interface    = "scsi0"
						size         = 8
					}

					disk {
						file_format  = "raw"
						datastore_id = "local-lvm"
						interface    = "scsi1"
						size         = 8
					}
				}`),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("proxmox_virtual_environment_vm.test_disk", plancheck.ResourceActionUpdate),
					},
				},
				Check: ResourceAttributes("proxmox_virtual_environment_vm.test_disk", map[string]string{
					"disk.0.interface":         "scsi0",
					"disk.0.path_in_datastore": `vm-\d+-disk-0`,
					"disk.1.interface":         "scsi1",
					"disk.1.path_in_datastore": `vm-\d+-disk-1`,
				}),
			},
			{
				RefreshState: true,
			},
		}},
		{"removing disks", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_disk" {
					node_name = "{{.NodeName}}"
					started   = false
					name 	  = "test-disk"
					
					disk {
						file_format  = "raw"
						datastore_id = "local-lvm"
						interface    = "scsi0"
						size         = 8
					}

					disk {
						file_format  = "raw"
						datastore_id = "local-lvm"
						interface    = "scsi1"
						size         = 8
					}
				}`),
				Check: ResourceAttributes("proxmox_virtual_environment_vm.test_disk", map[string]string{
					"disk.0.interface":         "scsi0",
					"disk.0.path_in_datastore": `vm-\d+-disk-0`,
					"disk.1.interface":         "scsi1",
					"disk.1.path_in_datastore": `vm-\d+-disk-1`,
				}),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_disk" {
					node_name = "{{.NodeName}}"
					started   = false
					name 	  = "test-disk"
					
					disk {
						file_format  = "raw"
						datastore_id = "local-lvm"
						interface    = "scsi0"
						size         = 8
					}
				}`),
				ExpectError: regexp.MustCompile(`deletion of disks not supported`),
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
					name 	  = "test-disk-template"
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
				resource "proxmox_virtual_environment_vm" "test_disk" {
					node_name = "{{.NodeName}}"
					started   = false
					name 	  = "test-disk"
		
					clone {
						vm_id = proxmox_virtual_environment_vm.test_disk_template.id
					}
		
					disk {
						interface    = "scsi0"
						//size = 10
					}
				}`),
				Check: ResourceAttributes("proxmox_virtual_environment_vm.test_disk", map[string]string{
					"disk.0.datastore_id":      "local-lvm",
					"disk.0.discard":           "on",
					"disk.0.file_format":       "raw",
					"disk.0.interface":         "scsi0",
					"disk.0.iothread":          "true",
					"disk.0.path_in_datastore": `base-\d+-disk-\d+`,
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
				resource "proxmox_virtual_environment_vm" "test_disk_template" {
					node_name = "{{.NodeName}}"
					started   = false
					name 	  = "test-disk-template"
					template  = "true"
					
					disk {
						file_format  = "raw"
						datastore_id = "local-lvm"
						interface    = "virtio0"
						size         = 8
					}
				}
				resource "proxmox_virtual_environment_vm" "test_disk" {
					node_name = "{{.NodeName}}"
					started   = false
					name 	  = "test-disk"

					clone {
						vm_id = proxmox_virtual_environment_vm.test_disk_template.id
					}

					disk {
						datastore_id = "local-lvm"
						interface    = "virtio0"
						size         = 10
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_vm.test_disk", map[string]string{
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
		{"clone with adding disk", []resource.TestStep{
			{
				SkipFunc: func() (bool, error) {
					// this test is failing because of "Attribute 'disk.1.size' expected to be set"
					return true, nil
				},
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_disk_template" {
					node_name = "{{.NodeName}}"
					started   = false
					name 	  = "test-disk-template"
					template  = "true"
					
					disk {
						file_format  = "raw"
						datastore_id = "local-lvm"
						interface    = "virtio0"
						size         = 8
					}
				}
				resource "proxmox_virtual_environment_vm" "test_disk" {
					node_name = "{{.NodeName}}"
					started   = false
					name 	  = "test-disk"

					clone {
						vm_id = proxmox_virtual_environment_vm.test_disk_template.id
					}

					disk {
						file_format  = "raw"
						datastore_id = "local-lvm"
						interface    = "scsi0"
						size         = 10
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_vm.test_disk", map[string]string{
						"disk.1.datastore_id": "local-lvm",
						"disk.1.interface":    "virtio0",
						"disk.1.size":         "8",
						"disk.0.datastore_id": "local-lvm",
						"disk.0.interface":    "scsi0",
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
