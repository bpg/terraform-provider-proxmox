//go:build acceptance || all

//testacc:tier=light
//testacc:resource=misc

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package config_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/test"
)

func TestAccResourceNodeConfig(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: te.RenderConfig(`
					resource "proxmox_node_config" "test" {
						node_name   = "{{.NodeName}}"
						description = "test notes"
					}

					data "proxmox_node_config" "test" {
						node_name = proxmox_node_config.test.node_name
					}
				`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_node_config.test", map[string]string{
						"description": "test notes",
					}),
					test.ResourceAttributes("data.proxmox_node_config.test", map[string]string{
						"description": "test notes",
					}),
				),
			},
			{
				Config: te.RenderConfig(`
					resource "proxmox_node_config" "test" {
						node_name   = "{{.NodeName}}"
						description = "updated notes"
					}
				`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_node_config.test", map[string]string{
						"description": "updated notes",
					}),
				),
			},
			{
				Config: te.RenderConfig(`
					resource "proxmox_node_config" "test" {
						node_name = "{{.NodeName}}"
					}
				`),
				Check: resource.ComposeTestCheckFunc(
					test.NoResourceAttributesSet("proxmox_node_config.test", []string{"description"}),
				),
			},
			{
				ResourceName:      "proxmox_node_config.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
