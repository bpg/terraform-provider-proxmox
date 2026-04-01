//go:build acceptance || all

//testacc:tier=light
//testacc:resource=backup

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package backup_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/test"
)

func TestAccResourceBackupJob(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"create and update backup job", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_backup_job" "test" {
					id       = "acc-test-bj"
					schedule = "*-*-* 02:00"
					storage  = "local"
					all      = true
					mode     = "snapshot"
					compress = "zstd"
				}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("proxmox_backup_job.test", "id", "acc-test-bj"),
					resource.TestCheckResourceAttr("proxmox_backup_job.test", "schedule", "*-*-* 02:00"),
					resource.TestCheckResourceAttr("proxmox_backup_job.test", "storage", "local"),
					resource.TestCheckResourceAttr("proxmox_backup_job.test", "all", "true"),
					resource.TestCheckResourceAttr("proxmox_backup_job.test", "mode", "snapshot"),
					resource.TestCheckResourceAttr("proxmox_backup_job.test", "compress", "zstd"),
				),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_backup_job" "test" {
					id       = "acc-test-bj"
					schedule = "*-*-* 03:00"
					storage  = "local"
					all      = true
					mode     = "stop"
					compress = "lzo"
					enabled  = false
				}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("proxmox_backup_job.test", "id", "acc-test-bj"),
					resource.TestCheckResourceAttr("proxmox_backup_job.test", "schedule", "*-*-* 03:00"),
					resource.TestCheckResourceAttr("proxmox_backup_job.test", "mode", "stop"),
					resource.TestCheckResourceAttr("proxmox_backup_job.test", "compress", "lzo"),
					resource.TestCheckResourceAttr("proxmox_backup_job.test", "enabled", "false"),
				),
			},
			{
				ResourceName:      "proxmox_backup_job.test",
				ImportStateId:     "acc-test-bj",
				ImportState:       true,
				ImportStateVerify: true,
			},
		}},
		{"create with minimal attributes", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_backup_job" "test_minimal" {
					id       = "acc-test-min"
					schedule = "sun 01:00"
					storage  = "local"
					all      = true
				}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("proxmox_backup_job.test_minimal", "id", "acc-test-min"),
					resource.TestCheckResourceAttr("proxmox_backup_job.test_minimal", "storage", "local"),
					resource.TestCheckResourceAttr("proxmox_backup_job.test_minimal", "all", "true"),
					resource.TestCheckResourceAttrSet("proxmox_backup_job.test_minimal", "enabled"),
				),
			},
		}},
		{"field deletion", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_backup_job" "test_delete" {
					id       = "acc-test-del"
					schedule = "*-*-* 04:00"
					storage  = "local"
					all      = true
					mode     = "snapshot"
					compress = "zstd"
					mailto   = ["test@example.com"]
				}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("proxmox_backup_job.test_delete", "id", "acc-test-del"),
					resource.TestCheckResourceAttr("proxmox_backup_job.test_delete", "mode", "snapshot"),
					resource.TestCheckResourceAttr("proxmox_backup_job.test_delete", "compress", "zstd"),
					resource.TestCheckResourceAttr("proxmox_backup_job.test_delete", "mailto.#", "1"),
					resource.TestCheckResourceAttr("proxmox_backup_job.test_delete", "mailto.0", "test@example.com"),
				),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_backup_job" "test_delete" {
					id       = "acc-test-del"
					schedule = "*-*-* 04:00"
					storage  = "local"
					all      = true
				}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("proxmox_backup_job.test_delete", "id", "acc-test-del"),
					resource.TestCheckResourceAttr("proxmox_backup_job.test_delete", "all", "true"),
					resource.TestCheckResourceAttr("proxmox_backup_job.test_delete", "mailto.#", "0"),
				),
			},
		}},
		{"backup specific VMs by ID", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_backup_job" "test_vmid" {
					id       = "acc-test-vmid"
					schedule = "*-*-* 06:00"
					storage  = "local"
					vmid     = ["100", "101", "102"]
					mode     = "snapshot"
				}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("proxmox_backup_job.test_vmid", "id", "acc-test-vmid"),
					resource.TestCheckResourceAttr("proxmox_backup_job.test_vmid", "vmid.#", "3"),
					resource.TestCheckResourceAttr("proxmox_backup_job.test_vmid", "vmid.0", "100"),
					resource.TestCheckResourceAttr("proxmox_backup_job.test_vmid", "vmid.1", "101"),
					resource.TestCheckResourceAttr("proxmox_backup_job.test_vmid", "vmid.2", "102"),
					resource.TestCheckResourceAttr("proxmox_backup_job.test_vmid", "mode", "snapshot"),
				),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_backup_job" "test_vmid" {
					id       = "acc-test-vmid"
					schedule = "*-*-* 06:00"
					storage  = "local"
					vmid     = ["100", "200"]
					mode     = "stop"
				}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("proxmox_backup_job.test_vmid", "vmid.#", "2"),
					resource.TestCheckResourceAttr("proxmox_backup_job.test_vmid", "vmid.0", "100"),
					resource.TestCheckResourceAttr("proxmox_backup_job.test_vmid", "vmid.1", "200"),
					resource.TestCheckResourceAttr("proxmox_backup_job.test_vmid", "mode", "stop"),
				),
			},
			{
				ResourceName:      "proxmox_backup_job.test_vmid",
				ImportStateId:     "acc-test-vmid",
				ImportState:       true,
				ImportStateVerify: true,
			},
		}},
		{"backup with retention policy", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_backup_job" "test_retention" {
					id            = "acc-test-ret"
					schedule      = "*-*-* 07:00"
					storage       = "local"
					all           = true
					prune_backups = {
						keep-daily = "7"
						keep-last  = "3"
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("proxmox_backup_job.test_retention", "id", "acc-test-ret"),
					resource.TestCheckResourceAttr("proxmox_backup_job.test_retention", "prune_backups.%", "2"),
					resource.TestCheckResourceAttr("proxmox_backup_job.test_retention", "prune_backups.keep-daily", "7"),
					resource.TestCheckResourceAttr("proxmox_backup_job.test_retention", "prune_backups.keep-last", "3"),
				),
			},
		}},
		{"backup with fleecing", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_backup_job" "test_fleecing" {
					id       = "acc-test-flc"
					schedule = "*-*-* 09:00"
					storage  = "local"
					all      = true
					fleecing = {
						enabled = true
						storage = "local"
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("proxmox_backup_job.test_fleecing", "id", "acc-test-flc"),
					resource.TestCheckResourceAttr("proxmox_backup_job.test_fleecing", "fleecing.enabled", "true"),
					resource.TestCheckResourceAttr("proxmox_backup_job.test_fleecing", "fleecing.storage", "local"),
				),
			},
		}},
		{"backup with performance settings", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_backup_job" "test_perf" {
					id       = "acc-test-perf"
					schedule = "*-*-* 10:00"
					storage  = "local"
					all      = true
					performance = {
						max_workers = 2
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("proxmox_backup_job.test_perf", "id", "acc-test-perf"),
					resource.TestCheckResourceAttr("proxmox_backup_job.test_perf", "performance.max_workers", "2"),
				),
			},
		}},
		{"backup with multiple mailto addresses", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_backup_job" "test_mailto" {
					id       = "acc-test-mail"
					schedule = "*-*-* 11:00"
					storage  = "local"
					all      = true
					mailto   = ["admin@example.com", "ops@example.com"]
				}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("proxmox_backup_job.test_mailto", "id", "acc-test-mail"),
					resource.TestCheckResourceAttr("proxmox_backup_job.test_mailto", "mailto.#", "2"),
					resource.TestCheckResourceAttr("proxmox_backup_job.test_mailto", "mailto.0", "admin@example.com"),
					resource.TestCheckResourceAttr("proxmox_backup_job.test_mailto", "mailto.1", "ops@example.com"),
				),
			},
			{
				ResourceName:      "proxmox_backup_job.test_mailto",
				ImportStateId:     "acc-test-mail",
				ImportState:       true,
				ImportStateVerify: true,
			},
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource.ParallelTest(t, resource.TestCase{
				ProtoV6ProviderFactories: te.AccProviders,
				Steps:                    tt.steps,
			})
		})
	}
}

func TestAccDataSourceBackupJobs(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_backup_job" "test_ds" {
					id       = "acc-test-ds"
					schedule = "*-*-* 05:00"
					storage  = "local"
					all      = true
				}

				data "proxmox_backup_jobs" "all" {
					depends_on = [proxmox_backup_job.test_ds]
				}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.proxmox_backup_jobs.all", "jobs.#"),
				),
			},
		},
	})
}
