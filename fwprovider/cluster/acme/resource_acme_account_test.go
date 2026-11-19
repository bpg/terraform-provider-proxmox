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
	"fmt"
	"regexp"
	"slices"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/assert"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/test"
)

func TestAccResourceACMEAccount(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)
	accountName := test.SafeResourceName("test-account")
	te.AddTemplateVars(map[string]interface{}{
		"AccountName": accountName,
	})

	const (
		directory   = "https://acme-staging-v02.api.letsencrypt.org/directory"
		tos         = "https://letsencrypt.org/documents/LE-SA-v1.2-November-15-2017.pdf"
		tosChanged  = "https://letsencrypt.org/documents/LE-SA-v1.5-February-24-2025.pdf"
		replaceTmpl = `
			resource "proxmox_acme_account" "test_account_replace" {
				name = "%s"
				contact = "le.ge9ro@passmail.net"
				directory = "%s"
				tos = "%s"
			}`
		replaceRes = "proxmox_acme_account.test_account_replace"
	)

	replaceConfig := func(suffix, directory, tos string) string {
		return te.RenderConfig(fmt.Sprintf(replaceTmpl, accountName+suffix, directory, tos), test.WithRootPassword())
	}

	expectReplace := []plancheck.PlanCheck{
		plancheck.ExpectResourceAction(replaceRes, plancheck.ResourceActionReplace),
	}

	var tosLocation, nameLocation string

	saveLocation := func(dst *string) resource.TestCheckFunc {
		return resource.TestCheckResourceAttrWith(replaceRes, "location", func(v string) error {
			*dst = v
			return nil
		})
	}

	newLocation := func(old *string) resource.TestCheckFunc {
		return resource.TestCheckResourceAttrWith(replaceRes, "location", func(v string) error {
			if v == "" || v == *old {
				return fmt.Errorf("expected a new account location, got %q (was %q)", v, *old)
			}

			return nil
		})
	}

	tests := []struct {
		name string
		step []resource.TestStep
	}{
		{"tos change replaces account", []resource.TestStep{
			{
				Config: replaceConfig("-tos", directory, tos),
				Check:  saveLocation(&tosLocation),
			},
			{
				Config:           replaceConfig("-tos", directory, tosChanged),
				ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: expectReplace},
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes(replaceRes, map[string]string{"tos": "^" + regexp.QuoteMeta(tosChanged) + "$"}),
					newLocation(&tosLocation),
				),
			},
		}},
		{"name change replaces account", []resource.TestStep{
			{
				Config: replaceConfig("-name", directory, tos),
				Check:  saveLocation(&nameLocation),
			},
			{
				Config:           replaceConfig("-renamed", directory, tos),
				ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: expectReplace},
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes(replaceRes, map[string]string{"name": "^" + accountName + "-renamed$"}),
					newLocation(&nameLocation),
					func(*terraform.State) error {
						out := te.ExecuteNodeCommands([]string{"ls /etc/pve/priv/acme"})
						assert.False(t, slices.Contains(strings.Fields(out), accountName+"-name"), "old account still present")
						assert.True(t, slices.Contains(strings.Fields(out), accountName+"-renamed"), "renamed account missing")

						return nil
					},
				),
			},
		}},
		{"directory change replaces account", []resource.TestStep{
			{
				Config: replaceConfig("-dir", directory, tos),
			},
			{
				Config:             replaceConfig("-dir", "https://acme-v02.api.letsencrypt.org/directory", tos),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
				ConfigPlanChecks:   resource.ConfigPlanChecks{PostApplyPostRefresh: expectReplace},
			},
		}},
		{"basic account creation", []resource.TestStep{
			{
				Config: te.RenderConfig(`
					resource "proxmox_acme_account" "test_account" {
						name = "{{.AccountName}}"
						contact = "le.ge9ro@passmail.net"
						directory = "https://acme-staging-v02.api.letsencrypt.org/directory"
						tos = "https://letsencrypt.org/documents/LE-SA-v1.2-November-15-2017.pdf"
					}`, test.WithRootPassword()),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_acme_account.test_account", map[string]string{
						"name":      accountName,
						"directory": "https://acme-staging-v02.api.letsencrypt.org/directory",
						"tos":       "https://letsencrypt.org/documents/LE-SA-v1.2-November-15-2017.pdf",
					}),
					test.ResourceAttributesSet("proxmox_acme_account.test_account", []string{
						"created_at",
						"location",
					}),
				),
			},
			{
				ResourceName:                         "proxmox_acme_account.test_account",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        accountName,
				ImportStateVerifyIdentifierAttribute: "name",
				ImportStateVerifyIgnore:              []string{"contact"}, // not returned by API on fresh read
			},
		}},
		{"account with EAB", []resource.TestStep{
			{
				Config: te.RenderConfig(`
					resource "proxmox_acme_account" "test_account_eab" {
						name = "{{.AccountName}}-eab"
						contact = "le.ge9ro@passmail.net"
						directory = "https://acme-staging-v02.api.letsencrypt.org/directory"
						tos = "https://letsencrypt.org/documents/LE-SA-v1.2-November-15-2017.pdf"
						eab_hmac_key = "test-hmac-key"
						eab_kid = "test-kid"
					}`, test.WithRootPassword()),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_acme_account.test_account_eab", map[string]string{
						"name":         fmt.Sprintf("%s-eab", accountName),
						"directory":    "https://acme-staging-v02.api.letsencrypt.org/directory",
						"tos":          "https://letsencrypt.org/documents/LE-SA-v1.2-November-15-2017.pdf",
						"eab_hmac_key": "test-hmac-key",
						"eab_kid":      "test-kid",
					}),
				),
			},
		}},
		{"update account", []resource.TestStep{
			{
				Config: te.RenderConfig(`
					resource "proxmox_acme_account" "test_account_update" {
						name = "{{.AccountName}}-update"
						contact = "le.ge9ro@passmail.net"
						directory = "https://acme-staging-v02.api.letsencrypt.org/directory"
						tos = "https://letsencrypt.org/documents/LE-SA-v1.2-November-15-2017.pdf"
					}`, test.WithRootPassword()),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_acme_account.test_account_update", map[string]string{
						"name": fmt.Sprintf("%s-update", accountName),
					}),
				),
			},
			{
				Config: te.RenderConfig(`
					resource "proxmox_acme_account" "test_account_update" {
						name = "{{.AccountName}}-update"
						contact = "le.ge9ro@passmail.net"
						directory = "https://acme-staging-v02.api.letsencrypt.org/directory"
						tos = "https://letsencrypt.org/documents/LE-SA-v1.2-November-15-2017.pdf"
					}`, test.WithRootPassword()),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_acme_account.test_account_update", map[string]string{
						"name": fmt.Sprintf("%s-update", accountName),
					}),
				),
			},
		}},
		{"invalid directory URL", []resource.TestStep{
			{
				Config: te.RenderConfig(`
					resource "proxmox_acme_account" "test_account_invalid" {
						name = "{{.AccountName}}-invalid"
						contact = "le.ge9ro@passmail.net"
						directory = "invalid-url"
						tos = "https://letsencrypt.org/documents/LE-SA-v1.2-November-15-2017.pdf"
					}`, test.WithRootPassword()),
				ExpectError: regexp.MustCompile(`must be a valid URL`),
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
