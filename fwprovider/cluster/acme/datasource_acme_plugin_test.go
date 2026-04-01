//go:build acceptance || all

//testacc:tier=light
//testacc:resource=acme

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package acme_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/test"
)

func TestAccDatasourceACMEPlugin(t *testing.T) {
	te := test.InitEnvironment(t)
	pluginName := test.SafeResourceName("test-ds-plugin")
	te.AddTemplateVars(map[string]interface{}{
		"PluginName": pluginName,
	})

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: te.RenderConfig(`
					resource "proxmox_acme_dns_plugin" "test_plugin" {
						plugin = "{{.PluginName}}"
						api = "cf"
						data = {
							"CF_Account_ID" = "Account_ID"
							"CF_Token" = "Token"
							"CF_Zone_ID" = "Zone_ID"
						}
					}

					data "proxmox_acme_plugin" "test" {
						depends_on = [proxmox_acme_dns_plugin.test_plugin]
						plugin = "{{.PluginName}}"
					}
				`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("data.proxmox_acme_plugin.test", map[string]string{
						"plugin": pluginName,
					}),
					test.ResourceAttributesSet("data.proxmox_acme_plugin.test", []string{
						"api",
						"digest",
						"validation_delay",
					}),
					resource.TestCheckResourceAttrSet("data.proxmox_acme_plugin.test", "data.%"),
				),
			},
		},
	})
}
