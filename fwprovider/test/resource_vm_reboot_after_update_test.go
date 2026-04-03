//go:build acceptance || all

//testacc:tier=heavy
//testacc:resource=vm

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package test

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// These tests cover two of the three reboot-after-update policy buckets:
//  1. changes that cannot be applied at all without taking the VM offline,
//     which should fail when reboot_after_update = false
//  2. the same changes succeeding when reboot_after_update = true.
//
// The third bucket — config changes that are applied but require a manual reboot,
// emitting a warning via vmFinalizePowerState — is not covered here.
func TestAccResourceVMRebootAfterUpdateTPMStatePolicy(t *testing.T) {
	t.Parallel()

	te := InitEnvironment(t)
	imageFileID := te.DownloadCloudImage()
	te.AddTemplateVars(map[string]any{"ImageFileID": imageFileID})

	var vmID string

	resourceName := "proxmox_virtual_environment_vm.test_reboot_after_update_tpm"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: te.RenderConfig(vmRebootAfterUpdateCloudImageConfig(
					"test_reboot_after_update_tpm",
					false,
					"local-lvm",
					"local-lvm",
					"",
				)),
				Check: resource.ComposeTestCheckFunc(
					captureVMID(resourceName, &vmID),
					resource.TestCheckResourceAttr(resourceName, "tpm_state.#", "0"),
					checkVMStatus(te, &vmID, "running"),
				),
			},
			{
				Config: te.RenderConfig(vmRebootAfterUpdateCloudImageConfig(
					"test_reboot_after_update_tpm",
					false,
					"local-lvm",
					"local-lvm",
					`
					tpm_state {
						datastore_id = "local-lvm"
						version      = "v2.0"
					}
					`,
				)),
				PreConfig: func() {
					ensureVMRunning(te, vmID)
				},
				ExpectError: regexp.MustCompile(`cannot add, remove, or update the TPM device`),
			},
			{
				Config: te.RenderConfig(vmRebootAfterUpdateCloudImageConfig(
					"test_reboot_after_update_tpm",
					true,
					"local-lvm",
					"local-lvm",
					`
					tpm_state {
						datastore_id = "local-lvm"
						version      = "v2.0"
					}
					`,
				)),
				PreConfig: func() {
					ensureVMRunning(te, vmID)
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "tpm_state.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "tpm_state.0.datastore_id", "local-lvm"),
					resource.TestCheckResourceAttr(resourceName, "tpm_state.0.version", "v2.0"),
					checkVMStatus(te, &vmID, "running"),
				),
			},
		},
	})
}

func TestAccResourceVMRebootAfterUpdateCloudInitMovePolicy(t *testing.T) {
	t.Parallel()

	te := InitEnvironment(t)
	imageFileID := te.DownloadCloudImage()
	te.AddTemplateVars(map[string]any{"ImageFileID": imageFileID})

	var vmID string

	resourceName := "proxmox_virtual_environment_vm.test_reboot_after_update_cloudinit"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: te.RenderConfig(vmRebootAfterUpdateCloudImageConfig(
					"test_reboot_after_update_cloudinit",
					false,
					"local-lvm",
					"local-lvm",
					"",
				)),
				Check: resource.ComposeTestCheckFunc(
					captureVMID(resourceName, &vmID),
					resource.TestCheckResourceAttr(resourceName, "initialization.0.datastore_id", "local-lvm"),
					checkVMStatus(te, &vmID, "running"),
				),
			},
			{
				Config: te.RenderConfig(vmRebootAfterUpdateCloudImageConfig(
					"test_reboot_after_update_cloudinit",
					false,
					"local-lvm",
					"local",
					"",
				)),
				PreConfig: func() {
					ensureVMRunning(te, vmID)
				},
				ExpectError: regexp.MustCompile(`cannot move the Cloud-Init drive`),
			},
			{
				Config: te.RenderConfig(vmRebootAfterUpdateCloudImageConfig(
					"test_reboot_after_update_cloudinit",
					true,
					"local-lvm",
					"local",
					"",
				)),
				PreConfig: func() {
					ensureVMRunning(te, vmID)
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "initialization.0.datastore_id", "local"),
					checkVMStatus(te, &vmID, "running"),
				),
			},
		},
	})
}

func TestAccResourceVMRebootAfterUpdateTemplatePolicy(t *testing.T) {
	t.Parallel()

	te := InitEnvironment(t)
	imageFileID := te.DownloadCloudImage()
	te.AddTemplateVars(map[string]any{"ImageFileID": imageFileID})

	var vmID string

	resourceName := "proxmox_virtual_environment_vm.test_reboot_after_update_template"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: te.RenderConfig(vmRebootAfterUpdateCloudImageConfig(
					"test_reboot_after_update_template",
					false,
					"local-lvm",
					"local-lvm",
					"",
				)),
				Check: resource.ComposeTestCheckFunc(
					captureVMID(resourceName, &vmID),
					resource.TestCheckResourceAttr(resourceName, "template", "false"),
					checkVMStatus(te, &vmID, "running"),
				),
			},
			{
				Config: te.RenderConfig(vmRebootAfterUpdateCloudImageConfig(
					"test_reboot_after_update_template",
					false,
					"local-lvm",
					"local-lvm",
					"template = true",
				)),
				PreConfig: func() {
					ensureVMRunning(te, vmID)
				},
				ExpectError: regexp.MustCompile(`cannot convert the VM to a template`),
			},
			{
				Config: te.RenderConfig(vmRebootAfterUpdateCloudImageConfig(
					"test_reboot_after_update_template",
					true,
					"local-lvm",
					"local-lvm",
					"template = true",
				)),
				PreConfig: func() {
					ensureVMRunning(te, vmID)
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "template", "true"),
					checkVMStatus(te, &vmID, "stopped"),
				),
			},
		},
	})
}

func TestAccResourceVMRebootAfterUpdateDiskMovePolicy(t *testing.T) {
	t.Parallel()

	te := InitEnvironment(t)
	imageFileID := te.DownloadCloudImage()
	te.AddTemplateVars(map[string]any{"ImageFileID": imageFileID})

	var vmID string

	resourceName := "proxmox_virtual_environment_vm.test_reboot_after_update_disk_move"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: te.RenderConfig(vmRebootAfterUpdateCloudImageConfig(
					"test_reboot_after_update_disk_move",
					false,
					"local-lvm",
					"local-lvm",
					"",
				)),
				Check: resource.ComposeTestCheckFunc(
					captureVMID(resourceName, &vmID),
					resource.TestCheckResourceAttr(resourceName, "disk.0.datastore_id", "local-lvm"),
					checkVMStatus(te, &vmID, "running"),
				),
			},
			{
				Config: te.RenderConfig(vmRebootAfterUpdateCloudImageConfig(
					"test_reboot_after_update_disk_move",
					false,
					"local",
					"local-lvm",
					"",
				)),
				PreConfig: func() {
					ensureVMRunning(te, vmID)
				},
				ExpectError: regexp.MustCompile(`cannot move disks between datastores`),
			},
			{
				Config: te.RenderConfig(vmRebootAfterUpdateCloudImageConfig(
					"test_reboot_after_update_disk_move",
					true,
					"local",
					"local-lvm",
					"",
				)),
				PreConfig: func() {
					ensureVMRunning(te, vmID)
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "disk.0.datastore_id", "local"),
					checkVMStatus(te, &vmID, "running"),
				),
			},
		},
	})
}

func TestAccResourceVMRebootAfterUpdateDiskResizePolicy(t *testing.T) {
	t.Parallel()

	te := InitEnvironment(t)
	imageFileID := te.DownloadCloudImage()
	te.AddTemplateVars(map[string]any{"ImageFileID": imageFileID})

	var vmID string

	resourceName := "proxmox_virtual_environment_vm.test_reboot_after_update_disk_resize"

	diskResizeConfig := func(rebootAfterUpdate bool, size int, aio string) string {
		return te.RenderConfig(fmt.Sprintf(`
			resource "proxmox_virtual_environment_vm" "test_reboot_after_update_disk_resize" {
				node_name           = "{{.NodeName}}"
				name                = "test-reboot-after-update-disk-resize"
				started             = true
				stop_on_destroy     = true
				reboot_after_update = %t

				cpu {
					cores = 2
				}

				memory {
					dedicated = 2048
				}

				disk {
					datastore_id = "local-lvm"
					file_format  = "raw"
					file_id      = "{{.ImageFileID}}"
					interface    = "scsi0"
					discard      = "on"
					size         = %d
					aio          = "%s"
				}

				initialization {
					interface    = "scsi1"
					datastore_id = "local-lvm"
					ip_config {
						ipv4 {
							address = "dhcp"
						}
					}
				}

				network_device {
					bridge = "vmbr0"
				}
			}
		`, rebootAfterUpdate, size, aio))
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				// Step 1: Create running VM with reboot_after_update=false
				Config: diskResizeConfig(false, 20, "io_uring"),
				Check: resource.ComposeTestCheckFunc(
					captureVMID(resourceName, &vmID),
					resource.TestCheckResourceAttr(resourceName, "disk.0.size", "20"),
					resource.TestCheckResourceAttr(resourceName, "disk.0.aio", "io_uring"),
					checkVMStatus(te, &vmID, "running"),
				),
			},
			{
				// Step 2: Change aio (triggers rebootRequired via disk.Update) AND resize disk
				// with reboot_after_update=false → expect error
				Config: diskResizeConfig(false, 25, "threads"),
				PreConfig: func() {
					ensureVMRunning(te, vmID)
				},
				ExpectError: regexp.MustCompile(`cannot resize disks`),
			},
			{
				// Step 3: Same change with reboot_after_update=true → expect success
				Config: diskResizeConfig(true, 25, "threads"),
				PreConfig: func() {
					ensureVMRunning(te, vmID)
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "disk.0.size", "25"),
					resource.TestCheckResourceAttr(resourceName, "disk.0.aio", "threads"),
					checkVMStatus(te, &vmID, "running"),
				),
			},
		},
	})
}

func vmRebootAfterUpdateCloudImageConfig(
	resourceName string,
	rebootAfterUpdate bool,
	diskDatastoreID string,
	initializationDatastoreID string,
	extra string,
) string {
	vmName := strings.ReplaceAll(resourceName, "_", "-")

	return fmt.Sprintf(`
		resource "proxmox_virtual_environment_vm" "%s" {
			node_name           = "{{.NodeName}}"
			name                = "%s"
			started             = true
			stop_on_destroy     = true
			reboot_after_update = %t

			cpu {
				cores = 2
			}

			memory {
				dedicated = 2048
			}

			disk {
				datastore_id = "%s"
				file_format  = "raw"
				file_id      = "{{.ImageFileID}}"
				interface    = "scsi0"
				discard      = "on"
				size         = 20
			}

			initialization {
				interface = "scsi1"
				datastore_id = "%s"
				ip_config {
					ipv4 {
						address = "dhcp"
					}
				}
			}

			network_device {
				bridge = "vmbr0"
			}

			%s
		}
	`, resourceName, vmName, rebootAfterUpdate, diskDatastoreID, initializationDatastoreID, extra)
}

func captureVMID(resourceName string, vmID *string) resource.TestCheckFunc {
	return resource.TestCheckResourceAttrWith(resourceName, "vm_id", func(v string) error {
		*vmID = v
		return nil
	})
}

func checkVMStatus(te *Environment, vmID *string, expectedStatus string) resource.TestCheckFunc {
	return func(*terraform.State) error {
		if vmID == nil || *vmID == "" {
			return fmt.Errorf("vm_id was not captured from state")
		}

		id, err := strconv.Atoi(*vmID)
		if err != nil {
			return fmt.Errorf("invalid vm_id %q: %w", *vmID, err)
		}

		status, err := te.NodeClient().VM(id).GetVMStatus(context.Background())
		if err != nil {
			return fmt.Errorf("failed to get VM status: %w", err)
		}

		if status == nil || status.Status != expectedStatus {
			return fmt.Errorf("expected VM %s to be %s, got %v", *vmID, expectedStatus, status)
		}

		return nil
	}
}

func ensureVMRunning(te *Environment, vmID string) {
	te.t.Helper()

	id, err := strconv.Atoi(vmID)
	if err != nil {
		te.t.Fatalf("invalid vm_id %q: %v", vmID, err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	vm := te.NodeClient().VM(id)

	status, err := vm.GetVMStatus(ctx)
	if err != nil {
		te.t.Fatalf("failed to get VM status for %s: %v", vmID, err)
	}

	if status != nil && status.Status == "running" {
		return
	}

	if result := vm.StartVM(ctx, 120); result.Err() != nil {
		te.t.Fatalf("failed to start VM %s: %v", vmID, result.Err())
	}

	if err := vm.WaitForVMStatus(ctx, "running"); err != nil {
		te.t.Fatalf("failed waiting for VM %s to start: %v", vmID, err)
	}
}
