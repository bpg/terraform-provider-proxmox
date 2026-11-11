//go:build acceptance || all

//testacc:tier=medium
//testacc:resource=vm

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package test

import (
	"context"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/vms"
	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

func TestAccResourceVMSMBIOSClone(t *testing.T) {
	t.Parallel()

	te := InitEnvironment(t)

	var cloneID int

	config := te.RenderConfig(`
	resource "proxmox_virtual_environment_vm" "template" {
		node_name = "{{.NodeName}}"
		started   = false
		template  = true
		name      = "test-smbios-clone-template"

		smbios {
			manufacturer = "template-manufacturer"
		}
	}

	resource "proxmox_virtual_environment_vm" "clone" {
		node_name = "{{.NodeName}}"
		started   = false
		name      = "test-smbios-clone"

		clone {
			vm_id = proxmox_virtual_environment_vm.template.vm_id
		}

		smbios {
			sku = "test-sku"
		}
	}

	resource "proxmox_virtual_environment_vm" "clone_inherit" {
		node_name = "{{.NodeName}}"
		started   = false
		name      = "test-smbios-clone-inherit"

		clone {
			vm_id = proxmox_virtual_environment_vm.template.vm_id
		}
	}`)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_vm.clone", map[string]string{
						"smbios.0.sku": "^test-sku$",
					}),
					checkVMSMBIOSOnPVE(te, "proxmox_virtual_environment_vm.clone", map[string]string{"sku": "test-sku"}),
					resource.TestCheckResourceAttr("proxmox_virtual_environment_vm.clone_inherit", "smbios.#", "0"),
					func(s *terraform.State) error {
						var err error

						cloneID, err = strconv.Atoi(s.RootModule().Resources["proxmox_virtual_environment_vm.clone"].Primary.Attributes["vm_id"])

						return err
					},
				),
			},
			{
				PreConfig: func() {
					vm := te.NodeClient().VM(cloneID)
					cfg, err := vm.GetVM(context.Background())
					require.NoError(t, err)
					require.NotNil(t, cfg.SMBIOS)

					err = vm.UpdateVM(context.Background(), &vms.UpdateRequestBody{
						SMBIOS: &vms.CustomSMBIOS{Base64: types.CustomBool(true).Pointer(), UUID: cfg.SMBIOS.UUID},
					})
					require.NoError(t, err)
				},
				Config:             config,
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config: config,
				Check:  checkVMSMBIOSOnPVE(te, "proxmox_virtual_environment_vm.clone", map[string]string{"sku": "test-sku"}),
			},
		},
	})
}
