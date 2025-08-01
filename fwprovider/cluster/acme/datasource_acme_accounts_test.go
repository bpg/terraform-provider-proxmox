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

func TestAccDatasourceACMEAccounts(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)
	accountName1 := fmt.Sprintf("test-ds-accounts1-%s", gofakeit.Word())
	accountName2 := fmt.Sprintf("test-ds-accounts2-%s", gofakeit.Word())
	te.AddTemplateVars(map[string]interface{}{
		"AccountName1": accountName1,
		"AccountName2": accountName2,
	})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_acme_account" "test_account1" {
						name = "{{.AccountName1}}"
						contact = "le.ge9ro@passmail.net"
						directory = "https://acme-staging-v02.api.letsencrypt.org/directory"
						tos = "https://letsencrypt.org/documents/LE-SA-v1.2-November-15-2017.pdf"
					}
					
					resource "proxmox_virtual_environment_acme_account" "test_account2" {
						name = "{{.AccountName2}}"
						contact = "le.ge9ro@passmail.net"
						directory = "https://acme-staging-v02.api.letsencrypt.org/directory"
						tos = "https://letsencrypt.org/documents/LE-SA-v1.2-November-15-2017.pdf"
					}

					data "proxmox_virtual_environment_acme_accounts" "test" {
						depends_on = [
							proxmox_virtual_environment_acme_account.test_account1,
							proxmox_virtual_environment_acme_account.test_account2
						]
					}
				`, test.WithRootUser()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.proxmox_virtual_environment_acme_accounts.test", "accounts.#"),
				),
			},
		},
	})
}
