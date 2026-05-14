//go:build acceptance || all

//testacc:tier=light
//testacc:resource=misc

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

	"github.com/bpg/terraform-provider-proxmox/utils"
)

func TestAccResourceGroupImport(t *testing.T) {
	if utils.GetAnyStringEnv("TF_ACC") == "" {
		t.Skip("Acceptance tests are disabled")
	}

	te := InitEnvironment(t)
	groupID := SafeResourceName("test-group-import")

	te.AddTemplateVars(map[string]interface{}{
		"GroupID": groupID,
	})

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_group" "test" {
						group_id = "{{.GroupID}}"
						comment  = "Test group for import"
						acl {
							path      = "/"
							propagate = true
							role_id   = "NoAccess"
						}
					}
				`),
				Check: ResourceAttributes("proxmox_virtual_environment_group.test", map[string]string{
					"group_id": groupID,
					"id":       groupID,
					"comment":  "Test group for import",
				}),
			},
			{
				ResourceName:      "proxmox_virtual_environment_group.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     groupID,
				// `acl` is no longer auto-populated on import (#2866):
				// imported groups start with empty acl state. Users opt in
				// by adding acl{} blocks (or migrate to `proxmox_acl`).
				ImportStateVerifyIgnore: []string{"acl"},
			},
		},
	})
}

// TestAccResourceGroupACLNotManagedHere verifies that when a group has no
// inline acl{} blocks and ACLs for that group are managed by separate
// proxmox_acl resources, refresh does not produce a spurious diff and apply
// does not destructively revoke the live ACL. Reproducer for #2866.
func TestAccResourceGroupACLNotManagedHere(t *testing.T) {
	if utils.GetAnyStringEnv("TF_ACC") == "" {
		t.Skip("Acceptance tests are disabled")
	}

	te := InitEnvironment(t)
	groupID := SafeResourceName("test-group-acl-external")

	te.AddTemplateVars(map[string]interface{}{
		"GroupID": groupID,
	})

	config := te.RenderConfig(`
		resource "proxmox_virtual_environment_group" "test" {
			group_id = "{{.GroupID}}"
			comment  = "ACLs managed via proxmox_acl"
		}
		resource "proxmox_acl" "test" {
			group_id  = proxmox_virtual_environment_group.test.group_id
			path      = "/"
			role_id   = "PVEAuditor"
			propagate = true
		}
	`)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_group.test", map[string]string{
						"group_id": groupID,
						"acl.#":    "0",
					}),
					ResourceAttributes("proxmox_acl.test", map[string]string{
						"group_id": groupID,
						"path":     "/",
						"role_id":  "PVEAuditor",
					}),
				),
			},
			{
				// Refresh + plan with the same config — must produce no diff.
				// Without the fix, groupRead unconditionally repopulates `acl`
				// from a cluster-wide ACL fetch, producing a spurious diff.
				Config: config,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}
