//go:build acceptance || all

//testacc:tier=heavy
//testacc:resource=vm

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

// TestAccResourceVMRebootAfterCreationWithAgent verifies that a VM with
// guest-agent enabled and reboot=true can be created successfully.
// This reproduces the issue where the provider attempts to reboot the VM
// immediately after starting it, before the guest agent is ready, causing
// 'qmp command guest-ping failed - got timeout'.
func TestAccResourceVMRebootAfterCreationWithAgent(t *testing.T) {
	t.Parallel()

	te := InitEnvironment(t)
	imageFileID := te.DownloadCloudImage()
	te.AddTemplateVars(map[string]any{"ImageFileID": imageFileID})

	var vmID string

	resourceName := "proxmox_virtual_environment_vm.test_reboot_creation_agent"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_file" "cloud_config_reboot" {
						content_type = "snippets"
						datastore_id = "local"
						node_name    = "{{.NodeName}}"
						overwrite    = true
						source_raw {
							data = <<-EOF
							#cloud-config
							runcmd:
							  - apt-get update
							  - apt-get install -y qemu-guest-agent
							  - systemctl enable qemu-guest-agent
							  - systemctl start qemu-guest-agent
							EOF
							file_name = "cloud-config-reboot-agent.yaml"
						}
					}

					resource "proxmox_virtual_environment_vm" "test_reboot_creation_agent" {
						node_name       = "{{.NodeName}}"
						name            = "test-reboot-creation-agent"
						started         = true
						stop_on_destroy = true
						reboot          = true

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
							file_format  = "raw"
							file_id      = "{{.ImageFileID}}"
							interface    = "scsi0"
							discard      = "on"
							size         = 20
						}

						initialization {
							interface    = "scsi1"
							datastore_id = "local-lvm"
							ip_config {
								ipv4 {
									address = "dhcp"
								}
							}
							user_data_file_id = proxmox_virtual_environment_file.cloud_config_reboot.id
						}

						network_device {
							bridge = "vmbr0"
						}
					}
				`),
				Check: resource.ComposeTestCheckFunc(
					captureVMID(resourceName, &vmID),
					checkVMStatus(te, &vmID, "running"),
				),
			},
		},
	})
}
