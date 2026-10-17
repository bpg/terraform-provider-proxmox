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
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccResourceVMCloneRetryOnStorageLock(t *testing.T) {
	const template = `
	resource "proxmox_virtual_environment_vm" "template" {
		node_name = "{{.NodeName}}"
		name      = "test-clone-lock-template"
		started   = false
		template  = true

		disk {
			datastore_id = "local-lvm"
			interface    = "scsi0"
			size         = 1
		}
	}
	`

	tests := []struct {
		name     string
		clone    string
		resource string
		vmIDAttr string
	}{
		{
			name: "sdk vm",
			clone: `
			resource "proxmox_virtual_environment_vm" "clone" {
				node_name = "{{.NodeName}}"
				name      = "test-clone-lock"
				started   = false

				clone {
					vm_id   = proxmox_virtual_environment_vm.template.vm_id
					full    = true
					retries = 3
				}
			}`,
			resource: "proxmox_virtual_environment_vm.clone",
			vmIDAttr: "vm_id",
		},
		{
			name: "framework cloned vm",
			clone: `
			resource "proxmox_cloned_vm" "clone" {
				node_name = "{{.NodeName}}"
				name      = "test-clone-lock"
				started   = false

				clone = {
					source_vm_id = proxmox_virtual_environment_vm.template.vm_id
					full         = true
					retries      = 3
				}
			}`,
			resource: "proxmox_cloned_vm.clone",
			vmIDAttr: "id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			te := InitEnvironment(t)

			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: te.AccProviders,
				Steps: []resource.TestStep{
					{
						Config: te.RenderConfig(template),
					},
					{
						// Hold the storage lock longer than PVE's 10s lock timeout, so the first clone task fails
						// on it and only a retried task can succeed.
						PreConfig: func() {
							te.ExecuteNodeCommands([]string{
								"(nohup flock /var/lock/pve-manager/pve-storage-local-lvm sleep 15 >/dev/null 2>&1 </dev/null &)",
								"sleep 1",
							})
						},
						Config: te.RenderConfig(template + tt.clone),
						Check: func(s *terraform.State) error {
							rs, ok := s.RootModule().Resources[tt.resource]
							if !ok {
								return fmt.Errorf("resource %s not found", tt.resource)
							}

							vmID, err := strconv.Atoi(rs.Primary.Attributes[tt.vmIDAttr])
							if err != nil {
								return fmt.Errorf("failed to parse VM ID: %w", err)
							}

							vmConfig, err := te.NodeClient().VM(vmID).GetVM(context.Background())
							if err != nil {
								return fmt.Errorf("failed to get VM %d config: %w", vmID, err)
							}

							if vmConfig.StorageDevices["scsi0"] == nil {
								return fmt.Errorf("cloned VM %d has no scsi0 disk", vmID)
							}

							return nil
						},
					},
				},
			})
		})
	}
}
