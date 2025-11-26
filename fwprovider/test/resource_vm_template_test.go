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

	"github.com/bpg/terraform-provider-proxmox/utils"
)

func TestAccResourceVMTemplateConversion(t *testing.T) {
	if utils.GetAnyStringEnv("TF_ACC") == "" {
		t.Skip("Acceptance tests are disabled")
	}

	t.Parallel()

	te := InitEnvironment(t)

	tests := []struct {
		name string
		step []resource.TestStep
	}{
		{"create template from VM with imported disk", []resource.TestStep{{
			Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_download_file" "cloud_image" {
					content_type = "iso"
					datastore_id = "local"
					node_name    = "{{.NodeName}}"
					url          = "{{.CloudImagesServer}}/minimal/releases/noble/release/ubuntu-24.04-minimal-cloudimg-amd64.img"
					overwrite_unmanaged = true
				}

				resource "proxmox_virtual_environment_vm" "template_vm" {
					node_name = "{{.NodeName}}"
					started   = false
					template  = true

					disk {
						datastore_id = "local-lvm"
						file_id      = proxmox_virtual_environment_download_file.cloud_image.id
						interface    = "virtio0"
						size         = 20
					}

					cpu {
						cores = 2
					}

					memory {
						dedicated = 2048
					}
				}`),
			Check: resource.ComposeTestCheckFunc(
				ResourceAttributes("proxmox_virtual_environment_vm.template_vm", map[string]string{
					"template": "true",
				}),
			),
		}}},
		{"create template and clone with linked clone", []resource.TestStep{{
			Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_download_file" "cloud_image" {
					content_type = "iso"
					datastore_id = "local"
					node_name    = "{{.NodeName}}"
					url          = "{{.CloudImagesServer}}/minimal/releases/noble/release/ubuntu-24.04-minimal-cloudimg-amd64.img"
					overwrite_unmanaged = true
				}

				resource "proxmox_virtual_environment_vm" "template_vm" {
					node_name = "{{.NodeName}}"
					started   = false
					template  = true

					disk {
						datastore_id = "local-lvm"
						file_id      = proxmox_virtual_environment_download_file.cloud_image.id
						interface    = "virtio0"
						size         = 20
					}

					cpu {
						cores = 2
					}

					memory {
						dedicated = 2048
					}
				}

				resource "proxmox_virtual_environment_vm" "clone_vm" {
					node_name = "{{.NodeName}}"
					started   = false

					clone {
						vm_id = proxmox_virtual_environment_vm.template_vm.vm_id
					}
				}`),
			Check: resource.ComposeTestCheckFunc(
				ResourceAttributes("proxmox_virtual_environment_vm.template_vm", map[string]string{
					"template": "true",
				}),
				ResourceAttributes("proxmox_virtual_environment_vm.clone_vm", map[string]string{
					"template": "false",
				}),
			),
		}}},
		{"create template with imported disk and linked clone - verifies disk naming fix", []resource.TestStep{{
			Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_download_file" "cloud_image" {
					content_type = "iso"
					datastore_id = "local"
					node_name    = "{{.NodeName}}"
					url          = "{{.CloudImagesServer}}/minimal/releases/noble/release/ubuntu-24.04-minimal-cloudimg-amd64.img"
					overwrite_unmanaged = true
				}

				resource "proxmox_virtual_environment_vm" "template_vm" {
					node_name = "{{.NodeName}}"
					started   = false
					template  = true

					disk {
						datastore_id = "local-lvm"
						file_id      = proxmox_virtual_environment_download_file.cloud_image.id
						interface    = "virtio0"
						size         = 20
					}

					cpu {
						cores = 2
					}

					memory {
						dedicated = 2048
					}
				}

				resource "proxmox_virtual_environment_vm" "linked_clone_vm" {
					node_name = "{{.NodeName}}"
					started   = false

					clone {
						vm_id = proxmox_virtual_environment_vm.template_vm.vm_id
						full  = false
					}

					cpu {
						cores = 2
					}

					memory {
						dedicated = 2048
					}
				}`),
			Check: resource.ComposeTestCheckFunc(
				ResourceAttributes("proxmox_virtual_environment_vm.template_vm", map[string]string{
					"template": "true",
				}),
				ResourceAttributes("proxmox_virtual_environment_vm.linked_clone_vm", map[string]string{
					"template": "false",
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
