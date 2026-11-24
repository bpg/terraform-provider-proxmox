//go:build acceptance || all

//testacc:tier=light
//testacc:resource=firewall

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package firewall_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/test"
)

func TestAccDataSourceClusterFirewallSecurityGroups(t *testing.T) {
	te := test.InitEnvironment(t)

	commented := randomGroupName()
	plain := randomGroupName()

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: te.RenderConfig(fmt.Sprintf(`
					resource "proxmox_cluster_firewall_security_group" "commented" {
						name    = %q
						comment = "Web servers"
					}

					resource "proxmox_cluster_firewall_security_group" "plain" {
						name = %q
					}

					data "proxmox_cluster_firewall_security_groups" "all" {
						depends_on = [
							proxmox_cluster_firewall_security_group.commented,
							proxmox_cluster_firewall_security_group.plain,
						]
					}`, commented, plain)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckTypeSetElemNestedAttrs("data.proxmox_cluster_firewall_security_groups.all", "groups.*",
						map[string]string{"name": commented, "comment": "Web servers"}),
					resource.TestCheckTypeSetElemNestedAttrs("data.proxmox_cluster_firewall_security_groups.all", "groups.*",
						map[string]string{"name": plain, "comment": ""}),
				),
			},
		},
	})
}
