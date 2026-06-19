//go:build acceptance || all

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

//testacc:tier=heavy
//testacc:resource=misc

package zfs_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/test"
)

// TestAccResourceDiskZFS tests lifecycle management of a ZFS pool on a Proxmox node.
// Requires PROXMOX_VE_ACC_ZFS_DISK to be set to a spare block device (e.g. /dev/vdb).
// The device will be fully wiped during testing.
func TestAccResourceDiskZFS(t *testing.T) {
	te := test.InitEnvironment(t)

	if te.ZfsDisk == "" {
		t.Skip("Skipping ZFS pool tests: PROXMOX_VE_ACC_ZFS_DISK is not set")
	}

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"create and read back", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_node_disk_zfs" "test" {
					node_name     = "{{.NodeName}}"
					name          = "test-zpool"
					devices       = ["{{.ZfsDisk}}"]
					raidlevel     = "single"
					cleanup_disks = true
				}`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_node_disk_zfs.test", map[string]string{
						"node_name": te.NodeName,
						"name":      "test-zpool",
						"raidlevel": "single",
						"state":     "ONLINE",
					}),
				),
			},
		}},
		{"update cleanup flags in-place", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_node_disk_zfs" "test" {
					node_name     = "{{.NodeName}}"
					name          = "test-zpool"
					devices       = ["{{.ZfsDisk}}"]
					raidlevel     = "single"
					cleanup_disks = false
				}`),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_node_disk_zfs" "test" {
					node_name      = "{{.NodeName}}"
					name           = "test-zpool"
					devices        = ["{{.ZfsDisk}}"]
					raidlevel      = "single"
					cleanup_config = true
					cleanup_disks  = true
				}`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_node_disk_zfs.test", map[string]string{
						"cleanup_config": "true",
						"cleanup_disks":  "true",
					}),
				),
				// cleanup_config and cleanup_disks have no RequiresReplace — must not trigger replace.
				ExpectNonEmptyPlan: false,
			},
		}},
		{"import", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_node_disk_zfs" "test" {
					node_name     = "{{.NodeName}}"
					name          = "test-zpool"
					devices       = ["{{.ZfsDisk}}"]
					raidlevel     = "single"
					cleanup_disks = true
				}`),
			},
			{
				ResourceName:  "proxmox_node_disk_zfs.test",
				ImportState:   true,
				ImportStateId: te.NodeName + "/test-zpool",
				// Write-only attributes (devices, raidlevel) cannot be reconstructed from the API.
				ImportStateVerifyIgnore: []string{"devices", "raidlevel", "ashift", "compression", "draid_config", "add_storage"},
			},
			{
				// After import + providing write-only values in config, plan must be empty (no replace).
				Config: te.RenderConfig(`
				resource "proxmox_node_disk_zfs" "test" {
					node_name     = "{{.NodeName}}"
					name          = "test-zpool"
					devices       = ["{{.ZfsDisk}}"]
					raidlevel     = "single"
					cleanup_disks = true
				}`),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: te.AccProviders,
				Steps:                    tt.steps,
			})
		})
	}
}
