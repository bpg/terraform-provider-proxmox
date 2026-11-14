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

func TestAccResourceVMCDROM(t *testing.T) {
	t.Parallel()

	te := InitEnvironment(t)

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"default no cdrom", []resource.TestStep{{
			Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_cdrom" {
					node_name = "{{.NodeName}}"
					started   = false
					name 	  = "test-cdrom"
				}`),
			Check: NoResourceAttributesSet("proxmox_virtual_environment_vm.test_cdrom", []string{"cdrom.#"}),
		}}},
		{"none cdrom", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_cdrom" {
					node_name = "{{.NodeName}}"
					started   = false
					name 	  = "test-cdrom"
					cdrom {
						file_id   = "none"
					}
				}`),
				Check: ResourceAttributes("proxmox_virtual_environment_vm.test_cdrom", map[string]string{
					"cdrom.0.file_id": "none",
				}),
			},
			{
				RefreshState: true,
			},
		}},
		{"sata cdrom", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_cdrom" {
					node_name = "{{.NodeName}}"
					started   = false
					name 	  = "test-cdrom"
					cdrom {
						file_id   = "none"
						interface = "sata3"	
					}
				}`),
				Check: ResourceAttributes("proxmox_virtual_environment_vm.test_cdrom", map[string]string{
					"cdrom.0.interface": "sata3",
				}),
			},
			{
				RefreshState: true,
			},
		}},
		{"scsi cdrom", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_cdrom" {
					node_name = "{{.NodeName}}"
					started   = false
					name 	  = "test-cdrom"
					cdrom {
						file_id   = "none"
						interface = "scsi5"	
					}
				}`),
				Check: ResourceAttributes("proxmox_virtual_environment_vm.test_cdrom", map[string]string{
					"cdrom.0.interface": "scsi5",
				}),
			},
			{
				RefreshState: true,
			},
		}},
		{"enable cdrom", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_cdrom" {
					node_name = "{{.NodeName}}"
					started   = false
					name 	  = "test-cdrom"
					cdrom {
						file_id   = "none"
					}
				}`),
				Check: ResourceAttributes("proxmox_virtual_environment_vm.test_cdrom", map[string]string{
					"cdrom.0.file_id": "none",
				}),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_cdrom" {
					node_name = "{{.NodeName}}"
					started   = false
					name 	  = "test-cdrom"
					cdrom {
						file_id   = "cdrom"
					}
				}`),
				Check: ResourceAttributes("proxmox_virtual_environment_vm.test_cdrom", map[string]string{
					"cdrom.0.file_id": "cdrom",
				}),
			},
		}},
		{"remove cdrom block on the default interface", removeCDROMSteps(te, "", "ide3")},
		{"remove cdrom block on sata", removeCDROMSteps(te, `interface = "sata3"`, "sata3")},
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

func removeCDROMSteps(te *Environment, interfaceAttr, iface string) []resource.TestStep {
	const name = "proxmox_virtual_environment_vm.test_cdrom"

	checkDrives := func(attached bool) resource.TestCheckFunc {
		return func(s *terraform.State) error {
			rs, ok := s.RootModule().Resources[name]
			if !ok {
				return fmt.Errorf("%s not found in state", name)
			}

			vmID, err := strconv.Atoi(rs.Primary.Attributes["vm_id"])
			if err != nil {
				return fmt.Errorf("failed to parse vm_id: %w", err)
			}

			vmConfig, err := te.NodeClient().VM(vmID).GetVM(context.Background())
			if err != nil {
				return fmt.Errorf("failed to get VM config: %w", err)
			}

			if attached {
				if _, exists := vmConfig.StorageDevices[iface]; !exists {
					return fmt.Errorf("%s is missing in Proxmox", iface)
				}

				return nil
			}

			for device, drive := range vmConfig.StorageDevices {
				if drive != nil {
					return fmt.Errorf("%s (%s) is still attached in Proxmox after the cdrom block was removed", device, drive.FileVolume)
				}
			}

			return nil
		}
	}

	return []resource.TestStep{
		{
			Config: te.RenderConfig(`
			resource "proxmox_virtual_environment_vm" "test_cdrom" {
				node_name = "{{.NodeName}}"
				started   = false
				name 	  = "test-cdrom"
				cdrom {
					file_id = "none"
					` + interfaceAttr + `
				}
			}`),
			Check: resource.ComposeTestCheckFunc(
				ResourceAttributes(name, map[string]string{"cdrom.0.interface": iface}),
				checkDrives(true),
			),
		},
		{
			Config: te.RenderConfig(`
			resource "proxmox_virtual_environment_vm" "test_cdrom" {
				node_name = "{{.NodeName}}"
				started   = false
				name 	  = "test-cdrom"
			}`),
			Check: resource.ComposeTestCheckFunc(
				NoResourceAttributesSet(name, []string{"cdrom.#"}),
				checkDrives(false),
			),
		},
	}
}
