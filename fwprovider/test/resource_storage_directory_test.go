//go:build acceptance || all

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/bpg/terraform-provider-proxmox/utils"
)

func TestAccResourceStorageDirectory(t *testing.T) {
	te := InitEnvironment(t)

	storageID := fmt.Sprintf("dir-%s-%d", strings.ToLower(gofakeit.Word()), time.Now().UnixNano())
	dirPath := utils.GetAnyStringEnv("PROXMOX_VE_ACC_STORAGE_DIR_PATH")
	if dirPath == "" {
		dirPath = "/var/lib/vz"
	}

	preallocation1 := "off"
	preallocation2 := "metadata"

	te.AddTemplateVars(map[string]any{
		"StorageID":      storageID,
		"DirPath":        dirPath,
		"DirPathReplace": dirPath + "/acc-replace",
		"Preallocation1": preallocation1,
		"Preallocation2": preallocation2,
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

					preallocation = "{{.Preallocation1}}"
					shared  = true
					disable = false
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_storage_directory.test", map[string]string{
						"id":            storageID,
						"path":          dirPath,
						"preallocation": preallocation1,
						"shared":        "true",
						"disable":       "false",
						"nodes.#":       "1",
						"content.#":     "1",
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
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_storage_directory" "test" {
					id     = "{{.StorageID}}"
					path   = "{{.DirPathReplace}}"
					nodes  = ["{{.NodeName}}"]
					content = ["images"]

					preallocation = "{{.Preallocation1}}"
					shared  = true
					disable = false
				}`),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_storage_directory" "test" {
					id     = "{{.StorageID}}"
					path   = "{{.DirPath}}"
					nodes  = ["{{.NodeName}}"]
					content = ["iso"]

					preallocation = "{{.Preallocation2}}"
					shared  = false
					disable = true
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_storage_directory.test", map[string]string{
						"preallocation": preallocation2,
						"shared":        "false",
						"disable":       "true",
						"nodes.#":       "1",
						"content.#":     "1",
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

					preallocation = "{{.Preallocation2}}"
					shared  = false
					disable = false

					backups {
						max_protected_backups = 5
						keep_last             = 3
						keep_hourly           = 12
						keep_daily            = 7
						keep_weekly           = 4
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_storage_directory.test", map[string]string{
						"disable":                       "false",
						"content.#":                     "1",
						"backups.max_protected_backups": "5",
						"backups.keep_last":             "3",
						"backups.keep_hourly":           "12",
						"backups.keep_daily":            "7",
						"backups.keep_weekly":           "4",
					}),
					resource.TestCheckTypeSetElemAttr(
						"proxmox_virtual_environment_storage_directory.test",
						"content.*",
						"backup",
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

					preallocation = "{{.Preallocation2}}"
					shared  = false
					disable = false

					backups {
						max_protected_backups = 5
						keep_all              = true
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_storage_directory.test", map[string]string{
						"disable":                       "false",
						"content.#":                     "1",
						"backups.max_protected_backups": "5",
						"backups.keep_all":              "true",
					}),
					resource.TestCheckTypeSetElemAttr(
						"proxmox_virtual_environment_storage_directory.test",
						"content.*",
						"backup",
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

					preallocation = "{{.Preallocation2}}"
					shared  = false
					disable = false

					backups {
						keep_all   = true
						keep_daily = 7
					}
				}`),
				ExpectError: regexp.MustCompile(`(?s)invalid backup retention settings|keep_all.*keep_`),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_storage_directory" "test" {
					id     = "{{.StorageID}}"
					path   = "{{.DirPath}}"
					nodes  = ["{{.NodeName}}"]
					content = ["backup"]

					preallocation = "{{.Preallocation2}}"
					shared  = false
					disable = false

					backups {
						max_protected_backups = 5
						keep_all              = true
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_storage_directory.test", map[string]string{
						"disable":                       "false",
						"content.#":                     "1",
						"backups.max_protected_backups": "5",
						"backups.keep_all":              "true",
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

// TestAccResourceStorageDirectoryRemoveBackups tests removing backups block from storage directory.
// Regression test for https://github.com/bpg/terraform-provider-proxmox/issues/2463
func TestAccResourceStorageDirectoryRemoveBackups(t *testing.T) {
	te := InitEnvironment(t)

	storageID := fmt.Sprintf("dir-%s-%d", strings.ToLower(gofakeit.Word()), time.Now().UnixNano())
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
			// create with backups block
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_storage_directory" "test" {
					id      = "{{.StorageID}}"
					path    = "{{.DirPath}}"
					content = ["backup", "iso", "snippets", "vztmpl"]
					shared  = false
					disable = false

					backups {
						keep_all = true
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_storage_directory.test", map[string]string{
						"id":               storageID,
						"path":             dirPath,
						"shared":           "false",
						"disable":          "false",
						"backups.keep_all": "true",
					}),
				),
			},
			// remove backups block - should not cause "was absent, but now present" error
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_storage_directory" "test" {
					id      = "{{.StorageID}}"
					path    = "{{.DirPath}}"
					content = ["backup", "iso", "snippets", "vztmpl"]
					shared  = false
					disable = false
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_storage_directory.test", map[string]string{
						"id": storageID,
					}),
					resource.TestCheckNoResourceAttr(
						"proxmox_virtual_environment_storage_directory.test",
						"backups",
					),
				),
			},
		},
	})
}
