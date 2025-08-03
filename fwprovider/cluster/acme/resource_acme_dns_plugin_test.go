//go:build acceptance || all

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package acme_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/test"
)

func TestAccResourceACMEDNSPlugin(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)
	pluginName := fmt.Sprintf("test-plugin-%s", gofakeit.Word())
	te.AddTemplateVars(map[string]interface{}{
		"PluginName": pluginName,
	})

	tests := []struct {
		name string
		step []resource.TestStep
	}{
		{"basic plugin creation", []resource.TestStep{
			{
				Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_acme_dns_plugin" "test_plugin" {
						plugin = "{{.PluginName}}"
						api = "cf"
						data = {
							"CF_API_EMAIL" = "test@example.com"
							"CF_API_KEY"   = "test-api-key"
						}
					}`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_virtual_environment_acme_dns_plugin.test_plugin", map[string]string{
						"plugin": pluginName,
						"api":    "cf",
					}),
					test.ResourceAttributesSet("proxmox_virtual_environment_acme_dns_plugin.test_plugin", []string{
						"digest",
					}),
				),
			},
		}},
		{"plugin with validation delay", []resource.TestStep{
			{
				Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_acme_dns_plugin" "test_plugin_delay" {
						plugin = "{{.PluginName}}-delay"
						api = "cf"
						data = {
							"CF_API_EMAIL" = "test@example.com"
							"CF_API_KEY"   = "test-api-key"
						}
						validation_delay = 60
					}`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_virtual_environment_acme_dns_plugin.test_plugin_delay", map[string]string{
						"plugin":           fmt.Sprintf("%s-delay", pluginName),
						"api":              "cf",
						"validation_delay": "60",
					}),
				),
			},
		}},
		{"plugin with disable flag", []resource.TestStep{
			{
				Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_acme_dns_plugin" "test_plugin_disabled" {
						plugin = "{{.PluginName}}-disabled"
						api = "cf"
						data = {
							"CF_API_EMAIL" = "test@example.com"
							"CF_API_KEY"   = "test-api-key"
						}
						disable = true
					}`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_virtual_environment_acme_dns_plugin.test_plugin_disabled", map[string]string{
						"plugin":  fmt.Sprintf("%s-disabled", pluginName),
						"api":     "cf",
						"disable": "true",
					}),
				),
			},
		}},
		{"update plugin", []resource.TestStep{
			{
				Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_acme_dns_plugin" "test_plugin_update" {
						plugin = "{{.PluginName}}-update"
						api = "cf"
						data = {
							"CF_API_EMAIL" = "test@example.com"
							"CF_API_KEY"   = "test-api-key"
						}
						validation_delay = 30
					}`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_virtual_environment_acme_dns_plugin.test_plugin_update", map[string]string{
						"validation_delay": "30",
					}),
				),
			},
			{
				Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_acme_dns_plugin" "test_plugin_update" {
						plugin = "{{.PluginName}}-update"
						api = "cf"
						data = {
							"CF_API_EMAIL" = "test@example.com"
							"CF_API_KEY"   = "test-api-key"
						}
						validation_delay = 120
					}`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_virtual_environment_acme_dns_plugin.test_plugin_update", map[string]string{
						"validation_delay": "120",
					}),
				),
			},
		}},
		{"invalid validation delay", []resource.TestStep{
			{
				Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_acme_dns_plugin" "test_plugin_invalid" {
						plugin = "{{.PluginName}}-invalid"
						api = "cf"
						data = {
							"CF_API_EMAIL" = "test@example.com"
							"CF_API_KEY"   = "test-api-key"
						}
						validation_delay = 200000
					}`),
				ExpectError: regexp.MustCompile(`Attribute validation_delay value must be between 0 and 172800`),
			},
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: te.AccProviders,
				Steps:                    tt.step,
			})
		})
	}
}
