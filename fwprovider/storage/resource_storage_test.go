//go:build acceptance || all

//testacc:tier=medium
//testacc:resource=storage

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package storage_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/test"
)

func TestAccResourceStorageDirectoryShortName(t *testing.T) {
	te := test.InitEnvironment(t)

	storageID := test.SafeResourceName("dir-short")
	te.AddTemplateVars(map[string]any{
		"StorageID": storageID,
	})

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			// Step 1: Create with short name and single content type
			{
				Config: te.RenderConfig(`
					resource "proxmox_storage_directory" "test" {
						id      = "{{.StorageID}}"
						path    = "/var/lib/vz"
						content = ["images"]
						nodes   = ["{{.NodeName}}"]
					}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("proxmox_storage_directory.test", "id", storageID),
					resource.TestCheckResourceAttr("proxmox_storage_directory.test", "path", "/var/lib/vz"),
					resource.TestCheckResourceAttr("proxmox_storage_directory.test", "content.#", "1"),
					resource.TestCheckResourceAttr("proxmox_storage_directory.test", "nodes.#", "1"),
				),
			},
			// Step 2: Update — change content to multiple types, toggle disable
			{
				Config: te.RenderConfig(`
					resource "proxmox_storage_directory" "test" {
						id      = "{{.StorageID}}"
						path    = "/var/lib/vz"
						content = ["images", "iso"]
						nodes   = ["{{.NodeName}}"]
						disable = true
					}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("proxmox_storage_directory.test", "content.#", "2"),
					resource.TestCheckResourceAttr("proxmox_storage_directory.test", "disable", "true"),
				),
			},
			// Step 3: Update — re-enable, change preallocation
			{
				Config: te.RenderConfig(`
					resource "proxmox_storage_directory" "test" {
						id            = "{{.StorageID}}"
						path          = "/var/lib/vz"
						content       = ["images", "iso"]
						nodes         = ["{{.NodeName}}"]
						disable       = false
						preallocation = "metadata"
					}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("proxmox_storage_directory.test", "disable", "false"),
					resource.TestCheckResourceAttr("proxmox_storage_directory.test", "preallocation", "metadata"),
				),
			},
			// Step 4: Import round-trip
			{
				ResourceName:      "proxmox_storage_directory.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     storageID,
			},
		},
	})
}
