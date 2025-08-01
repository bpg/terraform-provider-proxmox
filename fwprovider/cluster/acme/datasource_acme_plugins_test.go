//go:build acceptance || all

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package acme_test

import (
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/test"
)

func TestAccDatasourceACMEPlugins(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)
	pluginName1 := fmt.Sprintf("test-ds-plugins1-%s", gofakeit.Word())
	pluginName2 := fmt.Sprintf("test-ds-plugins2-%s", gofakeit.Word())
	te.AddTemplateVars(map[string]interface{}{
		"PluginName1": pluginName1,
		"PluginName2": pluginName2,
	})

	// First create some plugins to test against
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_acme_dns_plugin" "test_plugin1" {
						plugin = "{{.PluginName1}}"
						api = "cf"
						data = {
							"CF_API_EMAIL" = "test1@example.com"
							"CF_API_KEY"   = "test-api-key-1"
						}
					}
					
					resource "proxmox_virtual_environment_acme_dns_plugin" "test_plugin2" {
						plugin = "{{.PluginName2}}"
						api = "cf"
						data = {
							"CF_API_EMAIL" = "test2@example.com"
							"CF_API_KEY"   = "test-api-key-2"
						}
					}
				`),
			},
		},
	})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_acme_dns_plugin" "test_plugin1" {
						plugin = "{{.PluginName1}}"
						api = "cf"
						data = {
							"CF_API_EMAIL" = "le.ge9ro@passmail.net"
							"CF_API_KEY"   = "test-api-key-1"
						}
					}
					
					resource "proxmox_virtual_environment_acme_dns_plugin" "test_plugin2" {
						plugin = "{{.PluginName2}}"
						api = "cf"
						data = {
							"CF_API_EMAIL" = "le.ge9ro@passmail.net"
							"CF_API_KEY"   = "test-api-key-2"
						}
					}

					data "proxmox_virtual_environment_acme_plugins" "test" {
						depends_on = [
							proxmox_virtual_environment_acme_dns_plugin.test_plugin1,
							proxmox_virtual_environment_acme_dns_plugin.test_plugin2
						]
					}
				`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.proxmox_virtual_environment_acme_plugins.test", "plugins.#"),
				),
			},
		},
	})
}
