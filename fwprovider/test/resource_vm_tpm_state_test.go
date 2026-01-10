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
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestAccResourceVMTpmState(t *testing.T) {
	t.Parallel()

	te := InitEnvironment(t)

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"add and remove TPM state without replacement", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_tpm" {
					node_name = "{{.NodeName}}"
					started   = false
					name 	  = "test-tpm"

					disk {
						datastore_id = "local-lvm"
						interface    = "virtio0"
						size         = 8
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("proxmox_virtual_environment_vm.test_tpm", "tpm_state.0"),
				),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_tpm" {
					node_name = "{{.NodeName}}"
					started   = false
					name 	  = "test-tpm"

					disk {
						datastore_id = "local-lvm"
						interface    = "virtio0"
						size         = 8
					}

					tpm_state {
						datastore_id = "local-lvm"
						version      = "v2.0"
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_vm.test_tpm", map[string]string{
						"tpm_state.#":              "1",
						"tpm_state.0.datastore_id": "local-lvm",
						"tpm_state.0.version":      "v2.0",
					}),
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("proxmox_virtual_environment_vm.test_tpm", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_tpm" {
					node_name = "{{.NodeName}}"
					started   = false
					name 	  = "test-tpm"

					disk {
						datastore_id = "local-lvm"
						interface    = "virtio0"
						size         = 8
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("proxmox_virtual_environment_vm.test_tpm", "tpm_state.0"),
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("proxmox_virtual_environment_vm.test_tpm", plancheck.ResourceActionUpdate),
					},
				},
			},
		}},
		{"changing TPM version forces replacement", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_tpm_version" {
					node_name = "{{.NodeName}}"
					started   = false
					name 	  = "test-tpm-version"

					disk {
						datastore_id = "local-lvm"
						interface    = "virtio0"
						size         = 8
					}

					tpm_state {
						datastore_id = "local-lvm"
						version      = "v2.0"
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_vm.test_tpm_version", map[string]string{
						"tpm_state.#":              "1",
						"tpm_state.0.datastore_id": "local-lvm",
						"tpm_state.0.version":      "v2.0",
					}),
				),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_tpm_version" {
					node_name = "{{.NodeName}}"
					started   = false
					name 	  = "test-tpm-version"

					disk {
						datastore_id = "local-lvm"
						interface    = "virtio0"
						size         = 8
					}

					tpm_state {
						datastore_id = "local-lvm"
						version      = "v1.2"
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_vm.test_tpm_version", map[string]string{
						"tpm_state.#":              "1",
						"tpm_state.0.datastore_id": "local-lvm",
						"tpm_state.0.version":      "v1.2",
					}),
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("proxmox_virtual_environment_vm.test_tpm_version", plancheck.ResourceActionReplace),
					},
				},
			},
		}},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: te.AccProviders,
				Steps:                    tt.steps,
			})
		})
	}
}
