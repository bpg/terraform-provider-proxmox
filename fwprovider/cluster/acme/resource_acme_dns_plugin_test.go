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
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/test"
)

func TestAccResourceACMEDNSPlugin(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)
	pluginName := test.SafeResourceName("test-plugin")
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
					resource "proxmox_acme_dns_plugin" "test_plugin" {
						plugin = "{{.PluginName}}"
						api = "cf"
						data = {
							"CF_API_EMAIL" = "test@example.com"
							"CF_API_KEY"   = "test-api-key"
						}
					}`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_acme_dns_plugin.test_plugin", map[string]string{
						"plugin": pluginName,
						"api":    "cf",
					}),
					test.ResourceAttributesSet("proxmox_acme_dns_plugin.test_plugin", []string{
						"digest",
					}),
				),
			},
			{
				ResourceName:                         "proxmox_acme_dns_plugin.test_plugin",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        pluginName,
				ImportStateVerifyIdentifierAttribute: "plugin",
				ImportStateVerifyIgnore:              []string{"digest"}, // changes on re-read
			},
		}},
		{"plugin with validation delay", []resource.TestStep{
			{
				Config: te.RenderConfig(`
					resource "proxmox_acme_dns_plugin" "test_plugin_delay" {
						plugin = "{{.PluginName}}-delay"
						api = "cf"
						data = {
							"CF_API_EMAIL" = "test@example.com"
							"CF_API_KEY"   = "test-api-key"
						}
						validation_delay = 60
					}`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_acme_dns_plugin.test_plugin_delay", map[string]string{
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
					resource "proxmox_acme_dns_plugin" "test_plugin_disabled" {
						plugin = "{{.PluginName}}-disabled"
						api = "cf"
						data = {
							"CF_API_EMAIL" = "test@example.com"
							"CF_API_KEY"   = "test-api-key"
						}
						disable = true
					}`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_acme_dns_plugin.test_plugin_disabled", map[string]string{
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
					resource "proxmox_acme_dns_plugin" "test_plugin_update" {
						plugin = "{{.PluginName}}-update"
						api = "cf"
						data = {
							"CF_API_EMAIL" = "test@example.com"
							"CF_API_KEY"   = "test-api-key"
						}
						validation_delay = 30
					}`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_acme_dns_plugin.test_plugin_update", map[string]string{
						"validation_delay": "30",
					}),
				),
			},
			{
				Config: te.RenderConfig(`
					resource "proxmox_acme_dns_plugin" "test_plugin_update" {
						plugin = "{{.PluginName}}-update"
						api = "cf"
						data = {
							"CF_API_EMAIL" = "test@example.com"
							"CF_API_KEY"   = "test-api-key"
						}
						validation_delay = 120
					}`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_acme_dns_plugin.test_plugin_update", map[string]string{
						"validation_delay": "120",
					}),
				),
			},
		}},
		{"invalid validation delay", []resource.TestStep{
			{
				Config: te.RenderConfig(`
					resource "proxmox_acme_dns_plugin" "test_plugin_invalid" {
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
			resource.ParallelTest(t, resource.TestCase{
				ProtoV6ProviderFactories: te.AccProviders,
				Steps:                    tt.step,
			})
		})
	}
}

func TestAccResourceACMEDNSPluginWriteOnly(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)
	pluginName := test.SafeResourceName("test-plugin-wo")
	te.AddTemplateVars(map[string]interface{}{
		"PluginName": pluginName,
	})

	tests := []struct {
		name string
		step []resource.TestStep
	}{
		{"data_wo keeps secrets out of state", []resource.TestStep{
			{
				Config: te.RenderConfig(`
					resource "proxmox_acme_dns_plugin" "test_plugin_wo" {
						plugin = "{{.PluginName}}"
						api = "cf"
						data_wo = {
							"CF_API_EMAIL" = "test@example.com"
							"CF_API_KEY"   = "test-api-key"
						}
					}`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_acme_dns_plugin.test_plugin_wo", map[string]string{
						"plugin": pluginName,
						"api":    "cf",
					}),
					test.ResourceAttributesSet("proxmox_acme_dns_plugin.test_plugin_wo", []string{
						"digest",
					}),
					// Write-only data_wo must never be persisted to state, and the
					// deprecated data mirror must stay empty when data_wo is used.
					resource.TestCheckNoResourceAttr("proxmox_acme_dns_plugin.test_plugin_wo", "data_wo.%"),
					resource.TestCheckNoResourceAttr("proxmox_acme_dns_plugin.test_plugin_wo", "data.%"),
					// Behavioral proof: the write-only value actually reached the API
					// and was stored, even though it never lands in state.
					testCheckACMEPluginDataStored(te, pluginName, map[string]string{
						"CF_API_EMAIL": "test@example.com",
						"CF_API_KEY":   "test-api-key",
					}),
				),
			},
		}},
		{"data and data_wo are mutually exclusive", []resource.TestStep{
			{
				Config: te.RenderConfig(`
					resource "proxmox_acme_dns_plugin" "test_plugin_conflict" {
						plugin = "{{.PluginName}}-conflict"
						api = "cf"
						data = {
							"CF_API_KEY" = "test-api-key"
						}
						data_wo = {
							"CF_API_KEY" = "test-api-key"
						}
					}`),
				ExpectError: regexp.MustCompile(`These attributes cannot be configured together: \[data,data_wo\]`),
			},
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource.ParallelTest(t, resource.TestCase{
				// Write-only attributes require Terraform 1.11+.
				TerraformVersionChecks: []tfversion.TerraformVersionCheck{
					tfversion.SkipBelow(tfversion.Version1_11_0),
				},
				ProtoV6ProviderFactories: te.AccProviders,
				Steps:                    tt.step,
			})
		})
	}
}

// testCheckACMEPluginDataStored verifies, via a direct API read, that the plugin's
// DNS data matches want. Used to prove write-only data_wo reaches the API even
// though it is absent from Terraform state.
func testCheckACMEPluginDataStored(te *test.Environment, pluginID string, want map[string]string) resource.TestCheckFunc {
	return func(*terraform.State) error {
		plugin, err := te.ClusterClient().ACME().Plugins().Get(context.Background(), pluginID)
		if err != nil {
			return fmt.Errorf("reading ACME plugin %q: %w", pluginID, err)
		}

		if plugin.Data == nil {
			return fmt.Errorf("ACME plugin %q has no data; write-only data_wo was not stored", pluginID)
		}

		for k, v := range want {
			if got := (*plugin.Data)[k]; got != v {
				return fmt.Errorf("ACME plugin %q data[%q] = %q, want %q", pluginID, k, got, v)
			}
		}

		return nil
	}
}
