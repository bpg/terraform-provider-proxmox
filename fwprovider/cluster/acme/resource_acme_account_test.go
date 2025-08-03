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

func TestAccResourceACMEAccount(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)
	accountName := fmt.Sprintf("test-account-%s", gofakeit.Word())
	te.AddTemplateVars(map[string]interface{}{
		"AccountName": accountName,
	})

	tests := []struct {
		name string
		step []resource.TestStep
	}{
		{"basic account creation", []resource.TestStep{
			{
				Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_acme_account" "test_account" {
						name = "{{.AccountName}}"
						contact = "le.ge9ro@passmail.net"
						directory = "https://acme-staging-v02.api.letsencrypt.org/directory"
						tos = "https://letsencrypt.org/documents/LE-SA-v1.2-November-15-2017.pdf"
					}`, test.WithRootUser()),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_virtual_environment_acme_account.test_account", map[string]string{
						"name":      accountName,
						"directory": "https://acme-staging-v02.api.letsencrypt.org/directory",
						"tos":       "https://letsencrypt.org/documents/LE-SA-v1.2-November-15-2017.pdf",
					}),
					test.ResourceAttributesSet("proxmox_virtual_environment_acme_account.test_account", []string{
						"created_at",
						"location",
					}),
				),
			},
		}},
		{"account with EAB", []resource.TestStep{
			{
				Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_acme_account" "test_account_eab" {
						name = "{{.AccountName}}-eab"
						contact = "le.ge9ro@passmail.net"
						directory = "https://acme-staging-v02.api.letsencrypt.org/directory"
						tos = "https://letsencrypt.org/documents/LE-SA-v1.2-November-15-2017.pdf"
						eab_hmac_key = "test-hmac-key"
						eab_kid = "test-kid"
					}`, test.WithRootUser()),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_virtual_environment_acme_account.test_account_eab", map[string]string{
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
					resource "proxmox_virtual_environment_acme_account" "test_account_update" {
						name = "{{.AccountName}}-update"
						contact = "le.ge9ro@passmail.net"
						directory = "https://acme-staging-v02.api.letsencrypt.org/directory"
						tos = "https://letsencrypt.org/documents/LE-SA-v1.2-November-15-2017.pdf"
					}`, test.WithRootUser()),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_virtual_environment_acme_account.test_account_update", map[string]string{
						"name": fmt.Sprintf("%s-update", accountName),
					}),
				),
			},
			{
				Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_acme_account" "test_account_update" {
						name = "{{.AccountName}}-update"
						contact = "le.ge9ro@passmail.net"
						directory = "https://acme-staging-v02.api.letsencrypt.org/directory"
						tos = "https://letsencrypt.org/documents/LE-SA-v1.2-November-15-2017.pdf"
					}`, test.WithRootUser()),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_virtual_environment_acme_account.test_account_update", map[string]string{
						"name": fmt.Sprintf("%s-update", accountName),
					}),
				),
			},
		}},
		{"invalid directory URL", []resource.TestStep{
			{
				Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_acme_account" "test_account_invalid" {
						name = "{{.AccountName}}-invalid"
						contact = "le.ge9ro@passmail.net"
						directory = "invalid-url"
						tos = "https://letsencrypt.org/documents/LE-SA-v1.2-November-15-2017.pdf"
					}`, test.WithRootUser()),
				ExpectError: regexp.MustCompile(`must be a valid URL`),
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
