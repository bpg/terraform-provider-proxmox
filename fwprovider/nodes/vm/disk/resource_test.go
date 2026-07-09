//go:build acceptance || all

//testacc:tier=medium
//testacc:resource=vm

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package disk_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/test"
)

const resourceName = "proxmox_vm.test_vm"

func TestAccResourceVMDisk(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		// Verifies that a VM without a disk block produces no drift (NullValue guard).
		{"VM without disk block produces no drift", []resource.TestStep{{
			Config: te.RenderConfig(`
			resource "proxmox_vm" "test_vm" {
				node_name = "{{.NodeName}}"
				name = "test-disk"
			}`),
			Check: test.NoResourceAttributesSet(resourceName, []string{"disk.%"}),
		}}},
		// Create a VM with a single disk using a subset of fields.
		{"create VM with single disk", []resource.TestStep{{
			Config: te.RenderConfig(`
			resource "proxmox_vm" "test_vm" {
				node_name = "{{.NodeName}}"
				name = "test-disk"
				disk = {
					"scsi0" = {
						datastore_id = "local-lvm"
						size         = 8
					}
				}
			}`),
			Check: resource.ComposeTestCheckFunc(
				test.ResourceAttributes(resourceName, map[string]string{
					"disk.%":                  "1",
					"disk.scsi0.size":         "8",
					"disk.scsi0.datastore_id": "local-lvm",
				}),
			),
		}}},
		// Create a VM with a disk and all optional fields set.
		{"create VM with disk and all optional fields", []resource.TestStep{{
			Config: te.RenderConfig(`
			resource "proxmox_vm" "test_vm" {
				node_name = "{{.NodeName}}"
				name = "test-disk"
				disk = {
					"scsi0" = {
						datastore_id = "local-lvm"
						size         = 8
						aio          = "native"
						backup       = false
						cache        = "writeback"
						discard      = "on"
						file_format  = "raw"
						iothread     = true
					}
				}
			}`),
			Check: resource.ComposeTestCheckFunc(
				test.ResourceAttributes(resourceName, map[string]string{
					"disk.%":                  "1",
					"disk.scsi0.size":         "8",
					"disk.scsi0.datastore_id": "local-lvm",
					"disk.scsi0.aio":          "native",
					"disk.scsi0.backup":       "false",
					"disk.scsi0.cache":        "writeback",
					"disk.scsi0.discard":      "on",
					"disk.scsi0.file_format":  "raw",
					"disk.scsi0.iothread":     "true",
				}),
			),
		}}},
		// Create a VM with multiple disks on different interfaces.
		{"create VM with multiple disks", []resource.TestStep{{
			Config: te.RenderConfig(`
			resource "proxmox_vm" "test_vm" {
				node_name = "{{.NodeName}}"
				name = "test-disk"
				disk = {
					"scsi0" = {
						datastore_id = "local-lvm"
						size         = 8
					}
					"scsi1" = {
						datastore_id = "local-lvm"
						size         = 4
					}
				}
			}`),
			Check: resource.ComposeTestCheckFunc(
				test.ResourceAttributes(resourceName, map[string]string{
					"disk.%":                  "2",
					"disk.scsi0.size":         "8",
					"disk.scsi0.datastore_id": "local-lvm",
					"disk.scsi1.size":         "4",
					"disk.scsi1.datastore_id": "local-lvm",
				}),
			),
		}}},
		// Update: change disk attributes and verify state reflects changes.
		{"update disk attributes", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_vm" "test_vm" {
					node_name = "{{.NodeName}}"
					name = "test-disk"
					disk = {
						"scsi0" = {
							datastore_id = "local-lvm"
							size         = 8
							cache        = "none"
							discard      = "ignore"
						}
					}
				}`),
				Check: test.ResourceAttributes(resourceName, map[string]string{
					"disk.scsi0.size":    "8",
					"disk.scsi0.cache":   "none",
					"disk.scsi0.discard": "ignore",
				}),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_vm" "test_vm" {
					node_name = "{{.NodeName}}"
					name = "test-disk"
					disk = {
						"scsi0" = {
							datastore_id = "local-lvm"
							size         = 8
							cache        = "writeback"
							discard      = "on"
						}
					}
				}`),
				Check: test.ResourceAttributes(resourceName, map[string]string{
					"disk.scsi0.cache":   "writeback",
					"disk.scsi0.discard": "on",
				}),
			},
		}},
		// MapDiff: add a slot, then remove a slot.
		{"add and remove disk slots", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_vm" "test_vm" {
					node_name = "{{.NodeName}}"
					name = "test-disk"
					disk = {
						"scsi0" = {
							datastore_id = "local-lvm"
							size         = 8
						}
					}
				}`),
				Check: test.ResourceAttributes(resourceName, map[string]string{
					"disk.%":          "1",
					"disk.scsi0.size": "8",
				}),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_vm" "test_vm" {
					node_name = "{{.NodeName}}"
					name = "test-disk"
					disk = {
						"scsi0" = {
							datastore_id = "local-lvm"
							size         = 8
						}
						"scsi1" = {
							datastore_id = "local-lvm"
							size         = 4
						}
					}
				}`),
				Check: test.ResourceAttributes(resourceName, map[string]string{
					"disk.%":          "2",
					"disk.scsi0.size": "8",
					"disk.scsi1.size": "4",
				}),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_vm" "test_vm" {
					node_name = "{{.NodeName}}"
					name = "test-disk"
					disk = {
						"scsi1" = {
							datastore_id = "local-lvm"
							size         = 4
						}
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes(resourceName, map[string]string{
						"disk.%":          "1",
						"disk.scsi1.size": "4",
					}),
				),
			},
		}},
		// Block removal: set disk block then remove it entirely.
		{"add disk then remove the block entirely", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_vm" "test_vm" {
					node_name = "{{.NodeName}}"
					name = "test-disk"
					disk = {
						"scsi0" = {
							datastore_id = "local-lvm"
							size         = 8
						}
					}
				}`),
				Check: test.ResourceAttributes(resourceName, map[string]string{
					"disk.%":          "1",
					"disk.scsi0.size": "8",
				}),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_vm" "test_vm" {
					node_name = "{{.NodeName}}"
					name = "test-disk"
				}`),
				Check: test.NoResourceAttributesSet(resourceName, []string{"disk.%"}),
			},
		}},
		// Resize: increase disk size (PVE does not allow shrinking).
		{"resize disk larger", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_vm" "test_vm" {
					node_name = "{{.NodeName}}"
					name = "test-disk"
					disk = {
						"scsi0" = {
							datastore_id = "local-lvm"
							size         = 8
						}
					}
				}`),
				Check: test.ResourceAttributes(resourceName, map[string]string{
					"disk.scsi0.size": "8",
				}),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_vm" "test_vm" {
					node_name = "{{.NodeName}}"
					name = "test-disk"
					disk = {
						"scsi0" = {
							datastore_id = "local-lvm"
							size         = 16
						}
					}
				}`),
				Check: test.ResourceAttributes(resourceName, map[string]string{
					"disk.scsi0.size": "16",
				}),
			},
		}},
		// Import: verify import round-trip preserves disk state.
		{"create disk and import", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_vm" "test_vm" {
					node_name = "{{.NodeName}}"
					name = "test-disk"
					disk = {
						"scsi0" = {
							datastore_id = "local-lvm"
							size         = 8
							discard      = "on"
							iothread     = true
						}
					}
				}`),
				Check: test.ResourceAttributes(resourceName, map[string]string{
					"disk.scsi0.size":     "8",
					"disk.scsi0.discard":  "on",
					"disk.scsi0.iothread": "true",
				}),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdPrefix:     te.NodeName + "/",
				ImportStateVerifyIgnore: []string{"disk.scsi0.import_from"},
			},
		}},
		// Disk on different interface types.
		{"create disks on ide and sata interfaces", []resource.TestStep{{
			Config: te.RenderConfig(`
			resource "proxmox_vm" "test_vm" {
				node_name = "{{.NodeName}}"
				name = "test-disk"
				disk = {
					"ide0" = {
						datastore_id = "local-lvm"
						size         = 4
					}
					"sata0" = {
						datastore_id = "local-lvm"
						size         = 4
					}
				}
			}`),
			Check: test.ResourceAttributes(resourceName, map[string]string{
				"disk.%":          "2",
				"disk.ide0.size":  "4",
				"disk.sata0.size": "4",
			}),
		}}},
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

func TestAccResourceVMDiskValidators(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"invalid disk interface key", []resource.TestStep{{
			Config: te.RenderConfig(`
			resource "proxmox_vm" "test_vm" {
				node_name = "{{.NodeName}}"
				disk = {
					"invalid0" = {
						datastore_id = "local-lvm"
						size         = 8
					}
				}
			}`),
			ExpectError: regexp.MustCompile(`one of`),
		}}},
		{"invalid aio value", []resource.TestStep{{
			Config: te.RenderConfig(`
			resource "proxmox_vm" "test_vm" {
				node_name = "{{.NodeName}}"
				disk = {
					"scsi0" = {
						datastore_id = "local-lvm"
						size         = 8
						aio          = "invalid"
					}
				}
			}`),
			ExpectError: regexp.MustCompile(`must be one of`),
		}}},
		{"invalid cache value", []resource.TestStep{{
			Config: te.RenderConfig(`
			resource "proxmox_vm" "test_vm" {
				node_name = "{{.NodeName}}"
				disk = {
					"scsi0" = {
						datastore_id = "local-lvm"
						size         = 8
						cache        = "invalid"
					}
				}
			}`),
			ExpectError: regexp.MustCompile(`must be one of`),
		}}},
		{"invalid file_format value", []resource.TestStep{{
			Config: te.RenderConfig(`
			resource "proxmox_vm" "test_vm" {
				node_name = "{{.NodeName}}"
				disk = {
					"scsi0" = {
						datastore_id = "local-lvm"
						size         = 8
						file_format  = "invalid"
					}
				}
			}`),
			ExpectError: regexp.MustCompile(`must be one of`),
		}}},
		{"scsi index out of range", []resource.TestStep{{
			Config: te.RenderConfig(`
			resource "proxmox_vm" "test_vm" {
				node_name = "{{.NodeName}}"
				disk = {
					"scsi31" = {
						datastore_id = "local-lvm"
						size         = 8
					}
				}
			}`),
			ExpectError: regexp.MustCompile(`one of`),
		}}},
		{"ide index out of range", []resource.TestStep{{
			Config: te.RenderConfig(`
			resource "proxmox_vm" "test_vm" {
				node_name = "{{.NodeName}}"
				disk = {
					"ide4" = {
						datastore_id = "local-lvm"
						size         = 8
					}
				}
			}`),
			ExpectError: regexp.MustCompile(`one of`),
		}}},
		{"sata index out of range", []resource.TestStep{{
			Config: te.RenderConfig(`
			resource "proxmox_vm" "test_vm" {
				node_name = "{{.NodeName}}"
				disk = {
					"sata6" = {
						datastore_id = "local-lvm"
						size         = 8
					}
				}
			}`),
			ExpectError: regexp.MustCompile(`one of`),
		}}},
		// Verifies the relaxed slot regex accepts scsi30 (MAX_SCSI_DISKS=31).
		{"scsi30 is a valid disk slot", []resource.TestStep{{
			Config: te.RenderConfig(`
			resource "proxmox_vm" "test_vm" {
				node_name = "{{.NodeName}}"
				disk = {
					"scsi30" = {
						datastore_id = "local-lvm"
						size         = 8
					}
				}
			}`),
			Check: test.ResourceAttributes(resourceName, map[string]string{
				"disk.%":           "1",
				"disk.scsi30.size": "8",
			}),
		}}},
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
