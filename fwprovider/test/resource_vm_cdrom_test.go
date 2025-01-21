//go:build acceptance_vm || all

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

func TestAccResourceVMCDROM(t *testing.T) {
	t.Parallel()

	te := InitEnvironment(t)

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"default no cdrom", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_cdrom" {
					node_name = "{{.NodeName}}"
					started   = false
					name 	  = "test-cdrom"
				}`),
				Check: NoResourceAttributesSet("proxmox_virtual_environment_vm.test_cdrom", []string{"cdrom.#"}),
			},
		}},
		{"default cdrom", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_cdrom" {
					node_name = "{{.NodeName}}"
					started   = false
					name 	  = "test-cdrom"
					cdrom {
						enabled   = true
					}
				}`),
				Check: ResourceAttributes("proxmox_virtual_environment_vm.test_cdrom", map[string]string{
					"cdrom.0.enabled": "true",
				}),
			},
			{
				RefreshState: true,
			},
		}},
		{"none cdrom", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_cdrom" {
					node_name = "{{.NodeName}}"
					started   = false
					name 	  = "test-cdrom"
					cdrom {
						enabled   = true
						file_id   = "none"
					}
				}`),
				Check: ResourceAttributes("proxmox_virtual_environment_vm.test_cdrom", map[string]string{
					"cdrom.0.enabled": "true",
					"cdrom.0.file_id": "none",
				}),
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
