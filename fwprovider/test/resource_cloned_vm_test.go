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
)

// TestAccResourceClonedVMDiskSize tests the disk size attribute with string units.
func TestAccResourceClonedVMDiskSize(t *testing.T) {
	t.Parallel()

	te := InitEnvironment(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				// Create a template VM, then clone it with disk size override using string units
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "template" {
					node_name = "{{.NodeName}}"
					started   = false
					template  = true

					disk {
						datastore_id = "local-lvm"
						interface    = "scsi0"
						size         = 8
					}
				}

				resource "proxmox_cloned_vm" "clone" {
					node_name = "{{.NodeName}}"

					clone {
						source_vm_id = proxmox_virtual_environment_vm.template.vm_id
					}

					disk = {
						scsi0 = {
							datastore_id = "local-lvm"
							size         = "10G"
						}
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_cloned_vm.clone", map[string]string{
						"disk.scsi0.size": "10G",
					}),
				),
			},
		},
	})
}

// TestAccResourceClonedVMDiskSizeResize tests resizing a cloned VM disk using string units.
func TestAccResourceClonedVMDiskSizeResize(t *testing.T) {
	t.Parallel()

	te := InitEnvironment(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				// Create template and clone with initial size
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "template" {
					node_name = "{{.NodeName}}"
					started   = false
					template  = true

					disk {
						datastore_id = "local-lvm"
						interface    = "scsi0"
						size         = 8
					}
				}

				resource "proxmox_cloned_vm" "clone" {
					node_name = "{{.NodeName}}"

					clone {
						source_vm_id = proxmox_virtual_environment_vm.template.vm_id
					}

					disk = {
						scsi0 = {
							datastore_id = "local-lvm"
							size         = "10G"
						}
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_cloned_vm.clone", map[string]string{
						"disk.scsi0.size": "10G",
					}),
				),
			},
			{
				// Resize disk to larger size
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "template" {
					node_name = "{{.NodeName}}"
					started   = false
					template  = true

					disk {
						datastore_id = "local-lvm"
						interface    = "scsi0"
						size         = 8
					}
				}

				resource "proxmox_cloned_vm" "clone" {
					node_name = "{{.NodeName}}"

					clone {
						source_vm_id = proxmox_virtual_environment_vm.template.vm_id
					}

					disk = {
						scsi0 = {
							datastore_id = "local-lvm"
							size         = "15G"
						}
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_cloned_vm.clone", map[string]string{
						"disk.scsi0.size": "15G",
					}),
				),
			},
		},
	})
}
