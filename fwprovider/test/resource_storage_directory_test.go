//go:build acceptance || all

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/bpg/terraform-provider-proxmox/utils"
)

func TestAccResourceStorageDirectory(t *testing.T) {
	te := InitEnvironment(t)

	storageID := fmt.Sprintf("dir-%s", strings.ToLower(gofakeit.Word()))
	dirPath := utils.GetAnyStringEnv("PROXMOX_VE_ACC_STORAGE_DIR_PATH")
	if dirPath == "" {
		dirPath = "/var/lib/vz"
	}

	te.AddTemplateVars(map[string]any{
		"StorageID": storageID,
		"DirPath":   dirPath,
	})

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_storage_directory" "test" {
					id     = "{{.StorageID}}"
					path   = "{{.DirPath}}"
					nodes  = ["{{.NodeName}}"]
					content = ["images"]

					shared  = true
					disable = false
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_storage_directory.test", map[string]string{
						"id":        storageID,
						"path":      dirPath,
						"shared":    "true",
						"disable":   "false",
						"nodes.#":   "1",
						"content.#": "1",
					}),
					resource.TestCheckTypeSetElemAttr(
						"proxmox_virtual_environment_storage_directory.test",
						"nodes.*",
						te.NodeName,
					),
					resource.TestCheckTypeSetElemAttr(
						"proxmox_virtual_environment_storage_directory.test",
						"content.*",
						"images",
					),
				),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_storage_directory" "test" {
					id     = "{{.StorageID}}"
					path   = "{{.DirPath}}"
					nodes  = ["{{.NodeName}}"]
					content = ["iso"]

					shared  = false
					disable = true
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_storage_directory.test", map[string]string{
						"shared":    "false",
						"disable":   "true",
						"nodes.#":   "1",
						"content.#": "1",
					}),
					resource.TestCheckTypeSetElemAttr(
						"proxmox_virtual_environment_storage_directory.test",
						"content.*",
						"iso",
					),
				),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_storage_directory" "test" {
					id     = "{{.StorageID}}"
					path   = "{{.DirPath}}"
					nodes  = ["{{.NodeName}}"]
					content = ["backup"]

					shared  = false
					disable = false

					backups {
						max_protected_backups = 5
						keep_daily            = 7
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_storage_directory.test", map[string]string{
						"disable":                       "false",
						"content.#":                     "1",
						"backups.max_protected_backups": "5",
						"backups.keep_daily":            "7",
					}),
					resource.TestCheckTypeSetElemAttr(
						"proxmox_virtual_environment_storage_directory.test",
						"content.*",
						"backup",
					),
				),
			},
			{
				ResourceName:      "proxmox_virtual_environment_storage_directory.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     storageID,
				ImportStateVerifyIgnore: []string{
					"backups",
					"backups.keep_all",
					"backups.keep_daily",
					"backups.max_protected_backups",
				},
			},
		},
	})
}
