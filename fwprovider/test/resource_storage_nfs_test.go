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
	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/utils"
)

func TestAccResourceStorageNFS(t *testing.T) {
	te := InitEnvironment(t)

	nfsServer := utils.GetAnyStringEnv("PROXMOX_VE_ACC_NFS_SERVER")
	nfsExport := utils.GetAnyStringEnv("PROXMOX_VE_ACC_NFS_EXPORT")

	if nfsServer == "" {
		t.Skip("PROXMOX_VE_ACC_NFS_SERVER is not set")
	}

	require.NotEmpty(t, nfsExport, "PROXMOX_VE_ACC_NFS_EXPORT must be set when PROXMOX_VE_ACC_NFS_SERVER is set")

	storageID := fmt.Sprintf("nfs-%s-%d", strings.ToLower(gofakeit.Word()), time.Now().UnixNano())
	options1 := "vers=4"
	options2 := "vers=4,proto=tcp"
	preallocation := "off"
	preallocationPlanOnlyReplace := "metadata"

	te.AddTemplateVars(map[string]any{
		"StorageID":             storageID,
		"NFSServer":             nfsServer,
		"NFSExport":             nfsExport,
		"Options1":              options1,
		"Options2":              options2,
		"Preallocation":         preallocation,
		"SnapshotAsVolumeChain": true,
		"PreallocationReplace":  preallocationPlanOnlyReplace,
		"SnapshotReplace":       false,
	})

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_storage_nfs" "test" {
					id     = "{{.StorageID}}"
					server = "{{.NFSServer}}"
					export = "{{.NFSExport}}"

					nodes   = ["{{.NodeName}}"]
					content = ["images"]
					disable = false

					options  = "{{.Options1}}"
					preallocation           = "{{.Preallocation}}"
					snapshot_as_volume_chain = {{.SnapshotAsVolumeChain}}
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_storage_nfs.test", map[string]string{
						"id":                       storageID,
						"server":                   nfsServer,
						"export":                   nfsExport,
						"disable":                  "false",
						"options":                  options1,
						"shared":                   "true",
						"preallocation":            preallocation,
						"snapshot_as_volume_chain": "true",
						"nodes.#":                  "1",
						"content.#":                "1",
					}),
					resource.TestCheckTypeSetElemAttr(
						"proxmox_virtual_environment_storage_nfs.test",
						"nodes.*",
						te.NodeName,
					),
					resource.TestCheckTypeSetElemAttr(
						"proxmox_virtual_environment_storage_nfs.test",
						"content.*",
						"images",
					),
				),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_storage_nfs" "test" {
					id     = "{{.StorageID}}"
					server = "{{.NFSServer}}"
					export = "{{.NFSExport}}"

					nodes   = ["{{.NodeName}}"]
					content = ["images"]
					disable = false

					shared  = true
					options  = "{{.Options1}}"
					preallocation           = "{{.Preallocation}}"
					snapshot_as_volume_chain = {{.SnapshotAsVolumeChain}}
				}`),
				ExpectError: regexp.MustCompile(`(?i)(read.?only|cannot be set|shared)`),
			},
			{
				// both preallocation and snapshot_as_volume_chain are replace-only: ensure plan marks replacement
				PlanOnly: true,
				// PlanOnly can't use ConfigPlanChecks.PreApply, so use ExpectNonEmptyPlan for a stable assertion.
				ExpectNonEmptyPlan: true,
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_storage_nfs" "test" {
					id     = "{{.StorageID}}"
					server = "{{.NFSServer}}"
					export = "{{.NFSExport}}"

					nodes   = ["{{.NodeName}}"]
					content = ["images"]
					disable = false

					options  = "{{.Options1}}"
					preallocation           = "{{.PreallocationReplace}}"
					snapshot_as_volume_chain = {{.SnapshotReplace}}
				}`),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_storage_nfs" "test" {
					id     = "{{.StorageID}}"
					server = "{{.NFSServer}}"
					export = "{{.NFSExport}}"

					nodes   = ["{{.NodeName}}"]
					content = ["iso"]
					disable = true

					options  = "{{.Options2}}"
					preallocation           = "{{.Preallocation}}"
					snapshot_as_volume_chain = {{.SnapshotAsVolumeChain}}
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_storage_nfs.test", map[string]string{
						"disable":                  "true",
						"options":                  options2,
						"shared":                   "true",
						"preallocation":            preallocation,
						"snapshot_as_volume_chain": "true",
						"nodes.#":                  "1",
						"content.#":                "1",
					}),
					resource.TestCheckTypeSetElemAttr(
						"proxmox_virtual_environment_storage_nfs.test",
						"content.*",
						"iso",
					),
				),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_storage_nfs" "test" {
					id     = "{{.StorageID}}"
					server = "{{.NFSServer}}"
					export = "{{.NFSExport}}"

					nodes   = ["{{.NodeName}}"]
					content = ["backup"]
					disable = false

					options  = "{{.Options2}}"
					preallocation           = "{{.Preallocation}}"
					snapshot_as_volume_chain = {{.SnapshotAsVolumeChain}}

					backups {
						max_protected_backups = 5
						keep_daily            = 7
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_storage_nfs.test", map[string]string{
						"disable":                       "false",
						"options":                       options2,
						"shared":                        "true",
						"preallocation":                 preallocation,
						"snapshot_as_volume_chain":      "true",
						"content.#":                     "1",
						"backups.max_protected_backups": "5",
						"backups.keep_daily":            "7",
					}),
					resource.TestCheckTypeSetElemAttr(
						"proxmox_virtual_environment_storage_nfs.test",
						"content.*",
						"backup",
					),
				),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_storage_nfs" "test" {
					id     = "{{.StorageID}}"
					server = "{{.NFSServer}}"
					export = "{{.NFSExport}}"

					nodes   = ["{{.NodeName}}"]
					content = ["backup"]
					disable = false

					options  = "{{.Options2}}"
					preallocation           = "{{.Preallocation}}"
					snapshot_as_volume_chain = {{.SnapshotAsVolumeChain}}

					backups {
						max_protected_backups = 5
						keep_all              = true
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_storage_nfs.test", map[string]string{
						"disable":                       "false",
						"options":                       options2,
						"shared":                        "true",
						"preallocation":                 preallocation,
						"snapshot_as_volume_chain":      "true",
						"content.#":                     "1",
						"backups.max_protected_backups": "5",
						"backups.keep_all":              "true",
					}),
					resource.TestCheckTypeSetElemAttr(
						"proxmox_virtual_environment_storage_nfs.test",
						"content.*",
						"backup",
					),
				),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_storage_nfs" "test" {
					id     = "{{.StorageID}}"
					server = "{{.NFSServer}}"
					export = "{{.NFSExport}}"

					nodes   = ["{{.NodeName}}"]
					content = ["backup"]
					disable = false

					options  = "{{.Options2}}"
					preallocation           = "{{.Preallocation}}"
					snapshot_as_volume_chain = {{.SnapshotAsVolumeChain}}

					backups {
						keep_all   = true
						keep_daily = 7
					}
				}`),
				ExpectError: regexp.MustCompile(`(?s)invalid backup retention settings|keep_all.*keep_`),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_storage_nfs" "test" {
					id     = "{{.StorageID}}"
					server = "{{.NFSServer}}"
					export = "{{.NFSExport}}"

					nodes   = ["{{.NodeName}}"]
					content = ["backup"]
					disable = false

					options  = "{{.Options2}}"
					preallocation           = "{{.Preallocation}}"
					snapshot_as_volume_chain = {{.SnapshotAsVolumeChain}}

					backups {
						max_protected_backups = 5
						keep_all              = true
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_storage_nfs.test", map[string]string{
						"disable":                       "false",
						"options":                       options2,
						"shared":                        "true",
						"preallocation":                 preallocation,
						"snapshot_as_volume_chain":      "true",
						"content.#":                     "1",
						"backups.max_protected_backups": "5",
						"backups.keep_all":              "true",
					}),
					resource.TestCheckTypeSetElemAttr(
						"proxmox_virtual_environment_storage_nfs.test",
						"content.*",
						"backup",
					),
				),
			},
			{
				ResourceName:      "proxmox_virtual_environment_storage_nfs.test",
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
