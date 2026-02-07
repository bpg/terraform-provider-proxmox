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
		{
			"create and update backup job",
			[]resource.TestStep{
				{
					Config: te.RenderConfig(`
resource "proxmox_virtual_environment_backup_job" "test_backup" {
  id       = "test-backup-job"
  schedule = "02:00"
  storage  = "local"
  all      = true
  mode     = "snapshot"
  compress = "zstd"

  enabled          = true
  mailto           = "admin@example.com"
  mailnotification = "failure"
  prune_backups    = "keep-daily=7,keep-last=3"
}
					`),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("proxmox_virtual_environment_backup_job.test_backup", "id", "test-backup-job"),
						resource.TestCheckResourceAttr("proxmox_virtual_environment_backup_job.test_backup", "schedule", "02:00"),
						resource.TestCheckResourceAttr("proxmox_virtual_environment_backup_job.test_backup", "storage", "local"),
						resource.TestCheckResourceAttr("proxmox_virtual_environment_backup_job.test_backup", "all", "true"),
						resource.TestCheckResourceAttr("proxmox_virtual_environment_backup_job.test_backup", "mode", "snapshot"),
						resource.TestCheckResourceAttr("proxmox_virtual_environment_backup_job.test_backup", "compress", "zstd"),
						resource.TestCheckResourceAttr("proxmox_virtual_environment_backup_job.test_backup", "enabled", "true"),
						resource.TestCheckResourceAttr("proxmox_virtual_environment_backup_job.test_backup", "mailto", "admin@example.com"),
						resource.TestCheckResourceAttr("proxmox_virtual_environment_backup_job.test_backup", "mailnotification", "failure"),
						resource.TestCheckResourceAttr("proxmox_virtual_environment_backup_job.test_backup", "prune_backups", "keep-daily=7,keep-last=3"),
					),
				},
				{
					Config: te.RenderConfig(`
resource "proxmox_virtual_environment_backup_job" "test_backup" {
  id       = "test-backup-job"
  schedule = "03:00"
  storage  = "local"
  all      = true
  mode     = "snapshot"
  compress = "lzo"

  enabled          = false
  mailto           = "backup@example.com"
  mailnotification = "always"
  prune_backups    = "keep-last=5,keep-weekly=4"
}
					`),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("proxmox_virtual_environment_backup_job.test_backup", "schedule", "03:00"),
						resource.TestCheckResourceAttr("proxmox_virtual_environment_backup_job.test_backup", "compress", "lzo"),
						resource.TestCheckResourceAttr("proxmox_virtual_environment_backup_job.test_backup", "enabled", "false"),
						resource.TestCheckResourceAttr("proxmox_virtual_environment_backup_job.test_backup", "mailto", "backup@example.com"),
						resource.TestCheckResourceAttr("proxmox_virtual_environment_backup_job.test_backup", "mailnotification", "always"),
						resource.TestCheckResourceAttr("proxmox_virtual_environment_backup_job.test_backup", "prune_backups", "keep-last=5,keep-weekly=4"),
					),
				},
			},
		},
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

func TestAccResourceBackupJobImport(t *testing.T) {
	te := test.InitEnvironment(t)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: te.RenderConfig(`
resource "proxmox_virtual_environment_backup_job" "test_import" {
  id       = "test-import-job"
  schedule = "sun 02:00"
  storage  = "local"
  all      = true
}
				`),
			},
			{
				ResourceName:      "proxmox_virtual_environment_backup_job.test_import",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccDataSourceBackupJobs(t *testing.T) {
	te := test.InitEnvironment(t)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: te.RenderConfig(`
resource "proxmox_virtual_environment_backup_job" "test_ds1" {
  id       = "test-datasource-job-1"
  schedule = "02:00"
  storage  = "local"
  all      = true
}

resource "proxmox_virtual_environment_backup_job" "test_ds2" {
  id       = "test-datasource-job-2"
  schedule = "03:00"
  storage  = "local"
  all      = true
}

data "proxmox_virtual_environment_backup_jobs" "all" {
  depends_on = [
    proxmox_virtual_environment_backup_job.test_ds1,
    proxmox_virtual_environment_backup_job.test_ds2,
  ]
}
				`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.proxmox_virtual_environment_backup_jobs.all", "jobs.#"),
					test.ResourceAttributesSet("data.proxmox_virtual_environment_backup_jobs.all", []string{"jobs.0.id", "jobs.0.schedule"}),
				),
			},
		},
	})
}
