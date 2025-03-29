//go:build acceptance || all

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package cdrom_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/test"
)

const resourceName = "proxmox_virtual_environment_vm2.test_vm"

func TestAccResourceVM2CDROM(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"create VM default CDROM", []resource.TestStep{{
			Config: te.RenderConfig(`
			resource "proxmox_virtual_environment_vm2" "test_vm" {
				node_name = "{{.NodeName}}"
				name = "test-cdrom"
				cdrom = {
					"ide3" = {}
				}
			}`),
			Check: test.ResourceAttributes(resourceName, map[string]string{
				"cdrom.%":            "1",
				"cdrom.ide3.file_id": "cdrom",
			}),
		}}},
		{"create VM multiple CDROMs", []resource.TestStep{{
			Config: te.RenderConfig(`
			resource "proxmox_virtual_environment_vm2" "test_vm" {
				node_name = "{{.NodeName}}"
				name = "test-cdrom"
				cdrom = {
					"ide3" = {},
					"ide1" = {
						file_id   = "none"
					}
				}
			}`),
			Check: test.ResourceAttributes(resourceName, map[string]string{
				"cdrom.%":            "2",
				"cdrom.ide3.file_id": "cdrom",
				"cdrom.ide1.file_id": "none",
			}),
		}}},
		{"create VM with CDROM and then update it", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm2" "test_vm" {
					node_name = "{{.NodeName}}"
					name = "test-cdrom"
					cdrom = {
						"scsi2" = {
							file_id   = "none"
						},
						"ide2" = {
							file_id   = "cdrom"
						}
					}
				}`),
				Check: test.ResourceAttributes(resourceName, map[string]string{
					"cdrom.%":             "2",
					"cdrom.scsi2.file_id": "none",
					"cdrom.ide2.file_id":  "cdrom",
				}),
			},
			{ // now update the cdrom params and check if they are updated
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm2" "test_vm" {
					node_name = "{{.NodeName}}"
					name = "test-cdrom"
					cdrom = {
						"scsi2" = {
							file_id   = "cdrom"
						}
					}
				}`),
				Check: test.ResourceAttributes(resourceName, map[string]string{
					"cdrom.%":             "1",
					"cdrom.scsi2.file_id": "cdrom",
				}),
			},
			{
				RefreshState: true,
			},
		}},
		{"clone VM with CDROM", []resource.TestStep{{
			Config: te.RenderConfig(`
			resource "proxmox_virtual_environment_vm2" "template_vm" {
				node_name = "{{.NodeName}}"
				name = "template-cdrom"
				cdrom = {
					"ide3" = {
						file_id   = "cdrom"
					}
				}
			}
			resource "proxmox_virtual_environment_vm2" "test_vm" {
				node_name = "{{.NodeName}}"
				name = "test-cdrom"
				clone = {
					id = proxmox_virtual_environment_vm2.template_vm.id
				}
			}`),
			Check: test.ResourceAttributes(resourceName, map[string]string{
				"cdrom.%":            "1",
				"cdrom.ide3.file_id": "cdrom",
			}),
		}}},
		{"clone VM with some CDROM params and updating them in the clone", []resource.TestStep{{
			Config: te.RenderConfig(`
			resource "proxmox_virtual_environment_vm2" "template_vm" {
				node_name = "{{.NodeName}}"
				name = "template-cdrom"
				cdrom = {
					"ide1" = {
						file_id   = "none"
					},
					"ide2" = {
						file_id   = "cdrom"
					}
				}
			}
			resource "proxmox_virtual_environment_vm2" "test_vm" {
				node_name = "{{.NodeName}}"
				name = "test-cpu"
				clone = {
					id = proxmox_virtual_environment_vm2.template_vm.id
				}
				cdrom = {
					"ide1" = {
						file_id   = "cdrom"
					}
				}
			}`),
			Check: test.ResourceAttributes(resourceName, map[string]string{
				"cdrom.%":            "1",
				"cdrom.ide1.file_id": "cdrom",
			}),
		}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource.ParallelTest(t, resource.TestCase{
				ProtoV6ProviderFactories: te.AccProviders,
				Steps:                    tt.steps,
			})
		})
	}
}
