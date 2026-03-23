//go:build acceptance || all

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package storage_test

import (
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/test"
)

func TestAccResourceShortNameStorageDirectory(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)

	storageID := fmt.Sprintf("dir-short-%s", gofakeit.Word())
	te.AddTemplateVars(map[string]any{
		"StorageID": storageID,
	})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: te.RenderConfig(`
					resource "proxmox_storage_directory" "test" {
						id      = "{{.StorageID}}"
						path    = "/var/lib/vz"
						content = ["images"]
					}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("proxmox_storage_directory.test", "id"),
					resource.TestCheckResourceAttr("proxmox_storage_directory.test", "id", storageID),
					resource.TestCheckResourceAttr("proxmox_storage_directory.test", "path", "/var/lib/vz"),
				),
			},
		},
	})
}
