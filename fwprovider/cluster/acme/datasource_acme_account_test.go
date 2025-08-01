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

func TestAccDatasourceACMEAccount(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)
	accountName := fmt.Sprintf("test-ds-account-%s", gofakeit.Word())
	te.AddTemplateVars(map[string]interface{}{
		"AccountName": accountName,
	})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_acme_account" "test_account" {
						name = "{{.AccountName}}"
						contact = "le.ge9ro@passmail.net"
						directory = "https://acme-staging-v02.api.letsencrypt.org/directory"
						tos = "https://letsencrypt.org/documents/LE-SA-v1.2-November-15-2017.pdf"
					}

					data "proxmox_virtual_environment_acme_account" "test" {
						depends_on = [proxmox_virtual_environment_acme_account.test_account]
						name = "{{.AccountName}}"
					}
				`, test.WithRootUser()),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("data.proxmox_virtual_environment_acme_account.test", map[string]string{
						"name": accountName,
					}),
					test.ResourceAttributesSet("data.proxmox_virtual_environment_acme_account.test", []string{
						"account.created_at",
						"account.status",
						"directory",
						"location",
						"tos",
					}),
					resource.TestCheckResourceAttrSet("data.proxmox_virtual_environment_acme_account.test", "account.contact.#"),
				),
			},
		},
	})
}
