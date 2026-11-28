//go:build acceptance || all

//testacc:tier=heavy
//testacc:resource=vm

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package clonedvm_test

import (
	"context"
	"fmt"
	"maps"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/test"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/vms"
)

const templateWithTwoNets = `
	resource "proxmox_virtual_environment_vm" "template_vm" {
		node_name = "{{.NodeName}}"
		started   = false
		template  = true

		disk {
			datastore_id = "local-lvm"
			file_id      = "{{.ImageFileID}}"
			interface    = "virtio0"
			size         = 16
		}

		cpu {
			cores = 1
		}

		memory {
			dedicated = 1024
		}

		network_device {
			model  = "virtio"
			bridge = "vmbr0"
		}

		network_device {
			model  = "virtio"
			bridge = "vmbr0"
		}
	}
	`

const templateWithOneNet = `
	resource "proxmox_virtual_environment_vm" "template_vm" {
		node_name = "{{.NodeName}}"
		started   = false
		template  = true

		disk {
			datastore_id = "local-lvm"
			file_id      = "{{.ImageFileID}}"
			interface    = "virtio0"
			size         = 16
		}

		cpu {
			cores = 1
		}

		memory {
			dedicated = 1024
		}

		network_device {
			model  = "virtio"
			bridge = "vmbr0"
		}
	}
	`

func TestAccResourceClonedVM(t *testing.T) {
	imageFileID := test.InitEnvironment(t).DownloadCloudImage()

	t.Run("InheritAndDelete", func(t *testing.T) {
		te := test.InitEnvironment(t)
		te.AddTemplateVars(map[string]any{"ImageFileID": imageFileID})

		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: te.AccProviders,
			Steps: []resource.TestStep{
				{
					Config: te.RenderConfig(templateWithTwoNets + `
					resource "proxmox_cloned_vm" "keep_inherited" {
						node_name = "{{.NodeName}}"
						name      = "fwk-cloned-keep"
						started   = false

						clone = {
							source_vm_id = proxmox_virtual_environment_vm.template_vm.vm_id
						}

						network = {
							net0 = {
								bridge = "vmbr0"
								model  = "virtio"
							}
						}
					}

					resource "proxmox_cloned_vm" "delete_inherited" {
						node_name = "{{.NodeName}}"
						name      = "fwk-cloned-delete"
						started   = false

						clone = {
							source_vm_id = proxmox_virtual_environment_vm.template_vm.vm_id
						}

						network = {
							net0 = {
								bridge = "vmbr0"
								model  = "virtio"
							}
						}

						delete = {
							network = ["net1"]
						}
					}
					`),
					Check: resource.ComposeTestCheckFunc(
						checkNetworkSlot(te, "proxmox_cloned_vm.keep_inherited", "net1", true),
						checkNetworkSlot(te, "proxmox_cloned_vm.delete_inherited", "net1", false),
					),
				},
			},
		})
	})

	t.Run("StopManagingDoesNotDelete", func(t *testing.T) {
		te := test.InitEnvironment(t)
		te.AddTemplateVars(map[string]any{"ImageFileID": imageFileID})

		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: te.AccProviders,
			Steps: []resource.TestStep{
				{
					Config: te.RenderConfig(templateWithOneNet + `
					resource "proxmox_cloned_vm" "unmanaged" {
						node_name = "{{.NodeName}}"
						name      = "fwk-cloned-unmanage"
						started   = false

						clone = {
							source_vm_id = proxmox_virtual_environment_vm.template_vm.vm_id
						}

						network = {
							net0 = {
								bridge = "vmbr0"
								model  = "virtio"
							}
						}
					}
					`),
					Check: resource.ComposeTestCheckFunc(
						checkNetworkSlot(te, "proxmox_cloned_vm.unmanaged", "net0", true),
					),
				},
				{
					Config: te.RenderConfig(templateWithOneNet + `
					resource "proxmox_cloned_vm" "unmanaged" {
						node_name = "{{.NodeName}}"
						name      = "fwk-cloned-unmanage"
						started   = false

						clone = {
							source_vm_id = proxmox_virtual_environment_vm.template_vm.vm_id
						}
					}
					`),
					Check: resource.ComposeTestCheckFunc(
						checkNetworkSlot(te, "proxmox_cloned_vm.unmanaged", "net0", true),
					),
				},
			},
		})
	})

	t.Run("FailedCreateKeepsVMInState", func(t *testing.T) {
		te := test.InitEnvironment(t)
		te.AddTemplateVars(map[string]any{"ImageFileID": imageFileID})

		clonedVM := func(diskSize string) string {
			return te.RenderConfig(templateWithOneNet + fmt.Sprintf(`
					resource "proxmox_cloned_vm" "failed_create" {
						node_name = "{{.NodeName}}"
						name      = "fwk-cloned-failed-create"
						started   = false

						clone = {
							source_vm_id = proxmox_virtual_environment_vm.template_vm.vm_id
						}

						disk = {
							virtio1 = {
								datastore_id = "local-lvm"
								%s
							}
						}
					}
					`, diskSize))
		}

		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: te.AccProviders,
			Steps: []resource.TestStep{
				{
					// The clone succeeds and the disk step fails, so the VM must stay in the state (tainted).
					Config:      clonedVM(""),
					ExpectError: regexp.MustCompile(`Missing size_gb`),
				},
				{
					Config: clonedVM("size_gb = 1"),
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectResourceAction("proxmox_cloned_vm.failed_create", plancheck.ResourceActionReplace),
						},
					},
				},
			},
		})
	})

	t.Run("MapKeyStability", func(t *testing.T) {
		te := test.InitEnvironment(t)
		te.AddTemplateVars(map[string]any{"ImageFileID": imageFileID})

		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: te.AccProviders,
			Steps: []resource.TestStep{
				{
					Config: te.RenderConfig(templateWithTwoNets + `
					resource "proxmox_cloned_vm" "stability" {
						node_name = "{{.NodeName}}"
						name      = "fwk-cloned-stability"
						started   = false

						clone = {
							source_vm_id = proxmox_virtual_environment_vm.template_vm.vm_id
						}

						network = {
							net0 = {
								bridge = "vmbr0"
								model  = "virtio"
							}
						}
					}
					`),
					Check: resource.ComposeTestCheckFunc(
						checkNetworkSlot(te, "proxmox_cloned_vm.stability", "net0", true),
						checkNetworkSlot(te, "proxmox_cloned_vm.stability", "net1", true),
					),
				},
				{
					Config: te.RenderConfig(templateWithTwoNets + `
					resource "proxmox_cloned_vm" "stability" {
						node_name = "{{.NodeName}}"
						name      = "fwk-cloned-stability"
						started   = false

						clone = {
							source_vm_id = proxmox_virtual_environment_vm.template_vm.vm_id
						}

						network = {
							net0 = {
								bridge = "vmbr0"
								model  = "e1000"
								tag    = 100
							}
						}
					}
					`),
					Check: resource.ComposeTestCheckFunc(
						checkNetworkSlot(te, "proxmox_cloned_vm.stability", "net0", true),
						checkNetworkSlot(te, "proxmox_cloned_vm.stability", "net1", true),
					),
				},
			},
		})
	})

	t.Run("MemoryConfiguration", func(t *testing.T) {
		te := test.InitEnvironment(t)
		te.AddTemplateVars(map[string]any{"ImageFileID": imageFileID})

		baseConfig := `
		resource "proxmox_virtual_environment_vm" "template_vm" {
			node_name = "{{.NodeName}}"
			started   = false
			template  = true

			disk {
				datastore_id = "local-lvm"
				file_id      = "{{.ImageFileID}}"
				interface    = "virtio0"
				size         = 16
			}

			cpu {
				cores = 1
			}

			memory {
				dedicated = 512
			}

			network_device {
				model  = "virtio"
				bridge = "vmbr0"
			}
		}
		`

		resName := "proxmox_cloned_vm.memory_test"

		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: te.AccProviders,
			Steps: []resource.TestStep{
				{
					Config: te.RenderConfig(baseConfig + `
					resource "proxmox_cloned_vm" "memory_test" {
						node_name = "{{.NodeName}}"
						name      = "fwk-cloned-memory"
						started   = false

						clone = {
							source_vm_id = proxmox_virtual_environment_vm.template_vm.vm_id
						}

						memory = {
							size    = 2048
							balloon = 1024
							shares  = 1500
						}
					}
					`),
					Check: checkMemoryConfig(te, resName, 2048, 1024, 1500),
				},
				{
					Config: te.RenderConfig(baseConfig + `
					resource "proxmox_cloned_vm" "memory_test" {
						node_name = "{{.NodeName}}"
						name      = "fwk-cloned-memory"
						started   = false

						clone = {
							source_vm_id = proxmox_virtual_environment_vm.template_vm.vm_id
						}

						memory = {
							size    = 4096
							balloon = 2048
							shares  = 2000
						}
					}
					`),
					Check: checkMemoryConfig(te, resName, 4096, 2048, 2000),
				},
			},
		})
	})

	t.Run("Started", func(t *testing.T) {
		te := test.InitEnvironment(t)
		te.AddTemplateVars(map[string]any{"ImageFileID": imageFileID})

		baseConfig := `
		resource "proxmox_virtual_environment_vm" "template_vm" {
			node_name = "{{.NodeName}}"
			started   = false
			template  = true

			disk {
				datastore_id = "local-lvm"
				file_id      = "{{.ImageFileID}}"
				interface    = "virtio0"
				size         = 16
			}

			cpu {
				cores = 1
			}

			memory {
				dedicated = 512
			}

			network_device {
				model  = "virtio"
				bridge = "vmbr0"
			}
		}
		`

		resName := "proxmox_cloned_vm.started_test"

		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: te.AccProviders,
			Steps: []resource.TestStep{
				{
					Config: te.RenderConfig(baseConfig + `
					resource "proxmox_cloned_vm" "started_test" {
						node_name       = "{{.NodeName}}"
						name            = "fwk-cloned-started"
						stop_on_destroy = true

						clone = {
							source_vm_id = proxmox_virtual_environment_vm.template_vm.vm_id
						}
					}
					`),
					Check: resource.ComposeTestCheckFunc(
						checkVMStatus(te, resName, "running"),
						test.ResourceAttributes(resName, map[string]string{
							"started": "true",
						}),
					),
				},
				{
					Config: te.RenderConfig(baseConfig + `
					resource "proxmox_cloned_vm" "started_test" {
						node_name       = "{{.NodeName}}"
						name            = "fwk-cloned-started"
						started         = false
						stop_on_destroy = true

						clone = {
							source_vm_id = proxmox_virtual_environment_vm.template_vm.vm_id
						}
					}
					`),
					Check: resource.ComposeTestCheckFunc(
						checkVMStatus(te, resName, "stopped"),
						test.ResourceAttributes(resName, map[string]string{
							"started": "false",
						}),
					),
				},
				{
					Config: te.RenderConfig(baseConfig + `
					resource "proxmox_cloned_vm" "started_test" {
						node_name       = "{{.NodeName}}"
						name            = "fwk-cloned-started"
						started         = true
						stop_on_destroy = true

						clone = {
							source_vm_id = proxmox_virtual_environment_vm.template_vm.vm_id
						}
					}
					`),
					Check: resource.ComposeTestCheckFunc(
						checkVMStatus(te, resName, "running"),
						test.ResourceAttributes(resName, map[string]string{
							"started": "true",
						}),
					),
				},
			},
		})
	})

	t.Run("DiskManagement", func(t *testing.T) {
		te := test.InitEnvironment(t)
		te.AddTemplateVars(map[string]any{"ImageFileID": imageFileID})

		baseConfig := `
		resource "proxmox_virtual_environment_vm" "template_vm" {
			node_name = "{{.NodeName}}"
			started   = false
			template  = true

			disk {
				datastore_id = "local-lvm"
				file_id      = "{{.ImageFileID}}"
				interface    = "virtio0"
				size         = 16
			}

			disk {
				datastore_id = "local-lvm"
				interface    = "scsi0"
				size         = 8
			}

			cpu {
				cores = 1
			}

			memory {
				dedicated = 1024
			}

			network_device {
				model  = "virtio"
				bridge = "vmbr0"
			}
		}
		`

		resName := "proxmox_cloned_vm.disk_test"

		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: te.AccProviders,
			Steps: []resource.TestStep{
				{
					Config: te.RenderConfig(baseConfig + `
					resource "proxmox_cloned_vm" "disk_test" {
						node_name = "{{.NodeName}}"
						name      = "fwk-cloned-disk"
						started   = false

						clone = {
							source_vm_id = proxmox_virtual_environment_vm.template_vm.vm_id
						}

						disk = {
							virtio0 = {
								datastore_id = "local-lvm"
								size_gb      = 32
							}
						}
					}
					`),
					Check: resource.ComposeTestCheckFunc(
						checkDiskSlot(te, resName, "virtio0", true),
						checkDiskSlot(te, resName, "scsi0", true),
					),
				},
				{
					Config: te.RenderConfig(baseConfig + `
					resource "proxmox_cloned_vm" "disk_test" {
						node_name = "{{.NodeName}}"
						name      = "fwk-cloned-disk"
						started   = false

						clone = {
							source_vm_id = proxmox_virtual_environment_vm.template_vm.vm_id
						}

						disk = {
							virtio0 = {
								datastore_id = "local-lvm"
								size_gb      = 64
								discard      = "on"
								ssd          = true
							}

							scsi1 = {
								datastore_id = "local-lvm"
								size_gb      = 20
							}
						}
					}
					`),
					Check: resource.ComposeTestCheckFunc(
						checkDiskSlot(te, resName, "virtio0", true),
						checkDiskSlot(te, resName, "scsi0", true),
						checkDiskSlot(te, resName, "scsi1", true),
					),
				},
				{
					Config: te.RenderConfig(baseConfig + `
					resource "proxmox_cloned_vm" "disk_test" {
						node_name = "{{.NodeName}}"
						name      = "fwk-cloned-disk"
						started   = false

						clone = {
							source_vm_id = proxmox_virtual_environment_vm.template_vm.vm_id
						}

						disk = {
							virtio0 = {
								datastore_id = "local-lvm"
								size_gb      = 64
							}
						}

						delete = {
							disk = ["scsi0"]
						}
					}
					`),
					Check: resource.ComposeTestCheckFunc(
						checkDiskSlot(te, resName, "virtio0", true),
						checkDiskSlot(te, resName, "scsi0", false),
						checkDiskSlot(te, resName, "scsi1", true),
					),
				},
			},
		})
	})

	t.Run("DiskImportFrom", func(t *testing.T) {
		te := test.InitEnvironment(t)

		config := func(description string) string {
			return te.RenderConfig(fmt.Sprintf(`
			resource "proxmox_virtual_environment_download_file" "image" {
				content_type        = "import"
				datastore_id        = "local"
				node_name           = "{{.NodeName}}"
				url                 = "{{.CloudImagesServer}}/minimal/releases/noble/release/ubuntu-24.04-minimal-cloudimg-amd64.img"
				file_name           = "{{.TestName}}.qcow2"
				overwrite_unmanaged = true
			}

			resource "proxmox_virtual_environment_vm" "template_vm" {
				node_name = "{{.NodeName}}"
				started   = false
				template  = true

				disk {
					datastore_id = "local-lvm"
					interface    = "scsi0"
					size         = 8
				}
			}

			resource "proxmox_cloned_vm" "import_from_test" {
				node_name   = "{{.NodeName}}"
				name        = "fwk-cloned-import-from"
				description = %q
				started     = false

				clone = {
					source_vm_id = proxmox_virtual_environment_vm.template_vm.vm_id
				}

				disk = {
					scsi0 = {
						import_from = proxmox_virtual_environment_download_file.image.id
					}

					virtio1 = {
						datastore_id = "local-lvm"
						import_from  = proxmox_virtual_environment_download_file.image.id
					}

					virtio2 = {
						datastore_id = "local-lvm"
						import_from  = proxmox_virtual_environment_download_file.image.id
						size_gb      = 8
					}
				}
			}
			`, description))
		}

		resName := "proxmox_cloned_vm.import_from_test"

		var imported map[string]string

		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: te.AccProviders,
			Steps: []resource.TestStep{
				{
					Config: config("created"),
					Check: func(s *terraform.State) error {
						disks, err := rawDisks(te, s, resName)
						if err != nil {
							return err
						}

						for slot, pattern := range map[string]string{
							"scsi0":   `^local-lvm:vm-\d+-disk-\d+,.*size=8G`,
							"virtio1": `^local-lvm:vm-\d+-disk-\d+,.*size=\d+[MG]`,
							"virtio2": `^local-lvm:vm-\d+-disk-\d+,.*size=8G`,
						} {
							if !regexp.MustCompile(pattern).MatchString(disks[slot]) {
								return fmt.Errorf("disk %s is %q, expected to match %s", slot, disks[slot], pattern)
							}
						}

						if len(disks) != 3 {
							return fmt.Errorf("expected scsi0, virtio1 and virtio2 only, got %v", disks)
						}

						imported = disks

						return nil
					},
				},
				{
					Config: config("updated"),
					Check: func(s *terraform.State) error {
						disks, err := rawDisks(te, s, resName)
						if err != nil {
							return err
						}

						if !maps.Equal(disks, imported) {
							return fmt.Errorf("disks changed on an unrelated update: before %v, after %v", imported, disks)
						}

						return nil
					},
				},
			},
		})
	})

	t.Run("ResizeInheritedDiskKeepsOptions", func(t *testing.T) {
		te := test.InitEnvironment(t)

		templateConfig := `
					resource "proxmox_virtual_environment_vm" "template_vm" {
						node_name = "{{.NodeName}}"
						started   = false
						template  = true

						disk {
							datastore_id = "local-lvm"
							interface    = "scsi0"
							size         = 8
							discard      = "on"
							ssd          = true
							cache        = "writeback"
						}
					}
`

		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: te.AccProviders,
			Steps: []resource.TestStep{
				{
					Config: te.RenderConfig(templateConfig + `
					resource "proxmox_cloned_vm" "resize_inherited" {
						node_name = "{{.NodeName}}"
						name      = "fwk-cloned-resize-inherited"
						started   = false

						clone = {
							source_vm_id = proxmox_virtual_environment_vm.template_vm.vm_id
						}

						disk = {
							scsi0 = {
								size_gb = 16
							}
						}
					}
					`),
					Check: checkDiskOptions(te, "proxmox_cloned_vm.resize_inherited", "scsi0", func(disk *vms.CustomStorageDevice) error {
						if disk.Size == nil || disk.Size.InGigabytes() != 16 {
							return fmt.Errorf("expected scsi0 size 16G, got %v", disk.Size)
						}

						if disk.Discard == nil || *disk.Discard != "on" {
							return fmt.Errorf("expected scsi0 discard=on from the template, got %v", disk.Discard)
						}

						if disk.SSD == nil || !bool(*disk.SSD) {
							return fmt.Errorf("expected scsi0 ssd=1 from the template, got %v", disk.SSD)
						}

						if disk.Cache == nil || *disk.Cache != "writeback" {
							return fmt.Errorf("expected scsi0 cache=writeback from the template, got %v", disk.Cache)
						}

						return nil
					}),
				},
				{
					Config: te.RenderConfig(templateConfig + `
					resource "proxmox_cloned_vm" "resize_inherited" {
						node_name = "{{.NodeName}}"
						name      = "fwk-cloned-resize-inherited"
						started   = false

						clone = {
							source_vm_id = proxmox_virtual_environment_vm.template_vm.vm_id
						}

						disk = {
							scsi0 = {
								size_gb = 16
								cache   = "none"
								ssd     = false
							}
						}
					}
					`),
					Check: checkDiskOptions(te, "proxmox_cloned_vm.resize_inherited", "scsi0", func(disk *vms.CustomStorageDevice) error {
						if disk.Size == nil || disk.Size.InGigabytes() != 16 {
							return fmt.Errorf("expected scsi0 size 16G, got %v", disk.Size)
						}

						if disk.Discard == nil || *disk.Discard != "on" {
							return fmt.Errorf("expected scsi0 discard=on from the template, got %v", disk.Discard)
						}

						if disk.SSD == nil || bool(*disk.SSD) {
							return fmt.Errorf("expected scsi0 ssd=0 from the configuration, got %v", disk.SSD)
						}

						if disk.Cache == nil || *disk.Cache != "none" {
							return fmt.Errorf("expected scsi0 cache=none from the configuration, got %v", disk.Cache)
						}

						return nil
					}),
				},
			},
		})
	})

	t.Run("InheritedDiskOnFileStorageKeepsNoFormat", func(t *testing.T) {
		te := test.InitEnvironment(t)

		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: te.AccProviders,
			Steps: []resource.TestStep{
				{
					Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_vm" "template_vm" {
						node_name = "{{.NodeName}}"
						started   = false
						template  = true

						disk {
							datastore_id = "local"
							file_format  = "qcow2"
							interface    = "scsi0"
							size         = 1
						}
					}

					resource "proxmox_cloned_vm" "file_storage" {
						node_name = "{{.NodeName}}"
						name      = "fwk-cloned-file-storage"
						started   = false

						clone = {
							source_vm_id = proxmox_virtual_environment_vm.template_vm.vm_id
						}

						disk = {
							scsi0 = {
								size_gb = 2
							}
						}
					}
					`),
					Check: func(s *terraform.State) error {
						disks, err := rawDisks(te, s, "proxmox_cloned_vm.file_storage")
						if err != nil {
							return err
						}

						// the volume name ends in .qcow2, which the client reads back as a format
						if strings.Contains(disks["scsi0"], "format=") {
							return fmt.Errorf("expected no format on scsi0, got %q", disks["scsi0"])
						}

						if !strings.Contains(disks["scsi0"], "size=2G") {
							return fmt.Errorf("expected scsi0 size 2G, got %q", disks["scsi0"])
						}

						return nil
					},
				},
			},
		})
	})

	t.Run("Import", func(t *testing.T) {
		te := test.InitEnvironment(t)
		te.AddTemplateVars(map[string]any{"ImageFileID": imageFileID})

		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: te.AccProviders,
			Steps: []resource.TestStep{
				{
					Config: te.RenderConfig(templateWithOneNet + `
					resource "proxmox_cloned_vm" "import_test" {
						node_name = "{{.NodeName}}"
						name      = "fwk-cloned-import"
						started   = false

						clone = {
							source_vm_id = proxmox_virtual_environment_vm.template_vm.vm_id
						}

						cpu = {
							cores = 2
						}

						memory = {
							size = 2048
						}

						network = {
							net0 = {
								bridge = "vmbr0"
								model  = "virtio"
							}
						}
					}
					`),
				},
				{
					ResourceName:      "proxmox_cloned_vm.import_test",
					ImportState:       true,
					ImportStateVerify: true,
					// cloned_vm uses opt-in management, so imported state only contains
					// minimal fields - clone config and explicitly managed attributes are not preserved
					ImportStateVerifyIgnore: []string{
						"clone",
						"cpu",
						"memory",
						"network",
						"name",
					},
					ImportStateIdFunc: func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources["proxmox_cloned_vm.import_test"]
						if !ok {
							return "", fmt.Errorf("resource not found")
						}

						nodeName := rs.Primary.Attributes["node_name"]
						idStr := rs.Primary.Attributes["id"]

						return fmt.Sprintf("%s/%s", nodeName, idStr), nil
					},
				},
			},
		})
	})
}

// --- helper functions ---

func checkNetworkSlot(te *test.Environment, resourceName, slot string, expected bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource %s not found", resourceName)
		}

		idStr := rs.Primary.Attributes["id"]
		if idStr == "" {
			return fmt.Errorf("resource %s missing id in state", resourceName)
		}

		vmid, err := strconv.Atoi(idStr)
		if err != nil {
			return err
		}

		config, err := te.NodeClient().VM(vmid).GetVM(context.Background())
		if err != nil {
			return err
		}

		found := networkSlotPresent(config, slot)

		if expected && !found {
			return fmt.Errorf("expected slot %s to exist for %s", slot, resourceName)
		}

		if !expected && found {
			return fmt.Errorf("expected slot %s to be absent for %s", slot, resourceName)
		}

		return nil
	}
}

func checkDiskOptions(
	te *test.Environment,
	resourceName, slot string,
	check func(disk *vms.CustomStorageDevice) error,
) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource %s not found", resourceName)
		}

		vmid, err := strconv.Atoi(rs.Primary.Attributes["id"])
		if err != nil {
			return fmt.Errorf("parsing the id of %s: %w", resourceName, err)
		}

		config, err := te.NodeClient().VM(vmid).GetVM(context.Background())
		if err != nil {
			return err
		}

		disk := config.StorageDevices[slot]
		if disk == nil {
			return fmt.Errorf("slot %s not found for %s", slot, resourceName)
		}

		return check(disk)
	}
}

func networkSlotPresent(config *vms.GetResponseData, slot string) bool {
	return config.NetworkDevices[slot] != nil
}

func checkMemoryConfig(
	te *test.Environment,
	resourceName string,
	expectedMemory, expectedBalloon, expectedShares int,
) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource %s not found", resourceName)
		}

		idStr := rs.Primary.Attributes["id"]
		if idStr == "" {
			return fmt.Errorf("resource %s missing id in state", resourceName)
		}

		vmid, err := strconv.Atoi(idStr)
		if err != nil {
			return err
		}

		config, err := te.NodeClient().VM(vmid).GetVM(context.Background())
		if err != nil {
			return err
		}

		if config.DedicatedMemory == nil {
			return fmt.Errorf("memory size is nil for %s", resourceName)
		}

		if int64(*config.DedicatedMemory) != int64(expectedMemory) {
			return fmt.Errorf("expected memory size %d, got %d for %s", expectedMemory, *config.DedicatedMemory, resourceName)
		}

		if config.FloatingMemory == nil {
			return fmt.Errorf("balloon is nil for %s", resourceName)
		}

		if int64(*config.FloatingMemory) != int64(expectedBalloon) {
			return fmt.Errorf("expected balloon %d, got %d for %s", expectedBalloon, *config.FloatingMemory, resourceName)
		}

		if config.FloatingMemoryShares == nil {
			return fmt.Errorf("memory shares is nil for %s", resourceName)
		}

		if *config.FloatingMemoryShares != expectedShares {
			return fmt.Errorf("expected shares %d, got %d for %s", expectedShares, *config.FloatingMemoryShares, resourceName)
		}

		return nil
	}
}

func checkVMStatus(te *test.Environment, resourceName, expectedStatus string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource %s not found", resourceName)
		}

		idStr := rs.Primary.Attributes["id"]
		if idStr == "" {
			return fmt.Errorf("resource %s missing id in state", resourceName)
		}

		vmid, err := strconv.Atoi(idStr)
		if err != nil {
			return err
		}

		status, err := te.NodeClient().VM(vmid).GetVMStatus(context.Background())
		if err != nil {
			return err
		}

		if status.Status != expectedStatus {
			return fmt.Errorf("expected VM status %q, got %q for %s", expectedStatus, status.Status, resourceName)
		}

		return nil
	}
}

func checkDiskSlot(te *test.Environment, resourceName, slot string, expected bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource %s not found", resourceName)
		}

		idStr := rs.Primary.Attributes["id"]
		if idStr == "" {
			return fmt.Errorf("resource %s missing id in state", resourceName)
		}

		vmid, err := strconv.Atoi(idStr)
		if err != nil {
			return err
		}

		config, err := te.NodeClient().VM(vmid).GetVM(context.Background())
		if err != nil {
			return err
		}

		found := config.StorageDevices[slot] != nil

		if expected && !found {
			return fmt.Errorf("expected slot %s to exist for %s", slot, resourceName)
		}

		if !expected && found {
			return fmt.Errorf("expected slot %s to be absent for %s", slot, resourceName)
		}

		return nil
	}
}

func rawDisks(te *test.Environment, s *terraform.State, resourceName string) (map[string]string, error) {
	rs, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("resource %s not found", resourceName)
	}

	vmid, err := strconv.Atoi(rs.Primary.Attributes["id"])
	if err != nil {
		return nil, err
	}

	vm := te.NodeClient().VM(vmid)

	var resBody struct {
		Data map[string]any `json:"data,omitempty"`
	}

	if err := vm.DoRequest(context.Background(), http.MethodGet, vm.ExpandPath("config"), nil, &resBody); err != nil {
		return nil, err
	}

	disks := map[string]string{}

	for k, v := range resBody.Data {
		for _, prefix := range []string{"scsi", "virtio", "sata", "ide", "unused"} {
			if _, isSlot := slotIndex(k, prefix); isSlot {
				disks[k] = fmt.Sprint(v)
			}
		}
	}

	return disks, nil
}

func slotIndex(slot string, prefix string) (int, bool) {
	if !strings.HasPrefix(slot, prefix) {
		return 0, false
	}

	idx, err := strconv.Atoi(strings.TrimPrefix(slot, prefix))
	if err != nil || idx < 0 {
		return 0, false
	}

	return idx, true
}

func TestAccResourceClonedVMDiskValidation(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)

	tests := []struct {
		name        string
		disk        string
		expectError *regexp.Regexp
	}{
		{"import_from with an absolute path", `import_from = "/mnt/images/a+b.qcow2"`, nil},
		{"serial with characters PVE rejects", `serial = "a,discard=on"`, regexp.MustCompile(`(?s)serial.*must contain`)},
		{"import_from with file", "file = \"local-lvm:vm-100-disk-1\"\nimport_from = \"local:import/a.qcow2\"", regexp.MustCompile(`(?s)file.*cannot\s+be\s+specified\s+when.*import_from`)},
		{"empty datastore_id", `datastore_id = ""`, regexp.MustCompile(`(?s)datastore_id\s+string\s+length\s+must\s+be\s+at\s+least\s+1`)},
		{"import_from with a comma", `import_from = "local:import/a,b.qcow2"`, regexp.MustCompile(`(?s)import_from.*without\s+whitespace.*got:\s+local:import/a,b.qcow2`)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: te.AccProviders,
				Steps: []resource.TestStep{{
					Config: te.RenderConfig(fmt.Sprintf(`
					resource "proxmox_cloned_vm" "validation" {
						node_name = "{{.NodeName}}"
						name      = "fwk-cloned-validation"
						started   = false

						clone = {
							source_vm_id = 100
						}

						disk = {
							virtio0 = {
								datastore_id = "local-lvm"
								%s
							}
						}
					}`, tt.disk)),
					ExpectError:        tt.expectError,
					PlanOnly:           true,
					ExpectNonEmptyPlan: tt.expectError == nil,
				}},
			})
		})
	}
}
