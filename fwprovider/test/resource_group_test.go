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
			},
		},
	})
}
