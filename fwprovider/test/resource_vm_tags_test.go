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
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccResourceVMTagsUnknownAtPlan(t *testing.T) {
	t.Parallel()

	te := InitEnvironment(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "terraform_data" "tag" {
					input = "unknown-at-plan"
				}

				resource "proxmox_virtual_environment_vm" "test_tags" {
					node_name = "{{.NodeName}}"
					started   = false
					name      = "test-tags-unknown"
					tags      = ["known-tag", terraform_data.tag.output]
				}`),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectUnknownValue("proxmox_virtual_environment_vm.test_tags", tfjsonpath.New("tags")),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_vm.test_tags", map[string]string{
						"tags.#": "2",
					}),
					resource.TestCheckTypeSetElemAttr("proxmox_virtual_environment_vm.test_tags", "tags.*", "known-tag"),
					resource.TestCheckTypeSetElemAttr("proxmox_virtual_environment_vm.test_tags", "tags.*", "unknown-at-plan"),
					checkVMTagsOnPVE(te, "proxmox_virtual_environment_vm.test_tags", "known-tag;unknown-at-plan"),
				),
			},
			{
				Config: te.RenderConfig(`
				resource "terraform_data" "tag" {
					input = "unknown-on-update"
				}

				resource "proxmox_virtual_environment_vm" "test_tags" {
					node_name = "{{.NodeName}}"
					started   = false
					name      = "test-tags-unknown"
					tags      = ["known-tag", terraform_data.tag.output]
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_vm.test_tags", map[string]string{
						"tags.#": "2",
					}),
					resource.TestCheckTypeSetElemAttr("proxmox_virtual_environment_vm.test_tags", "tags.*", "known-tag"),
					resource.TestCheckTypeSetElemAttr("proxmox_virtual_environment_vm.test_tags", "tags.*", "unknown-on-update"),
					checkVMTagsOnPVE(te, "proxmox_virtual_environment_vm.test_tags", "known-tag;unknown-on-update"),
				),
			},
			{
				Config: te.RenderConfig(`
				resource "terraform_data" "tag" {
					input = "unknown-on-update"
				}

				resource "proxmox_virtual_environment_vm" "test_tags" {
					node_name = "{{.NodeName}}"
					started   = false
					name      = "test-tags-unknown"
					tags      = ["unknown-on-update", "known-tag"]
				}`),
				PlanOnly: true,
			},
		},
	})
}

func checkVMTagsOnPVE(te *Environment, resourceName string, want string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource %q not found in state", resourceName)
		}

		vmID, err := strconv.Atoi(rs.Primary.Attributes["vm_id"])
		if err != nil {
			return fmt.Errorf("failed to parse vm_id: %w", err)
		}

		vmConfig, err := te.NodeClient().VM(vmID).GetVM(context.Background())
		if err != nil {
			return fmt.Errorf("failed to get VM config: %w", err)
		}

		got := ""
		if vmConfig.Tags != nil {
			got = *vmConfig.Tags
		}

		if got != want {
			return fmt.Errorf("unexpected tags on PVE: got %q, want %q", got, want)
		}

		return nil
	}
}
