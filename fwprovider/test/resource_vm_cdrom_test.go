//go:build acceptance || all

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package test

import (
	"fmt"
	"math/rand"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccResourceVMCDROM(t *testing.T) {
	t.Parallel()

	te := InitEnvironment(t)

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"default no cdrom", []resource.TestStep{{
			Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_cdrom" {
					node_name = "{{.NodeName}}"
					started   = false
					name 	  = "test-cdrom"
				}`),
			Check: NoResourceAttributesSet("proxmox_virtual_environment_vm.test_cdrom", []string{"cdrom.#"}),
		}}},
		{"none cdrom", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_cdrom" {
					node_name = "{{.NodeName}}"
					started   = false
					name 	  = "test-cdrom"
					cdrom {
						file_id   = "none"
						interface = "ide2"
					}
				}`),
				Check: ResourceAttributes("proxmox_virtual_environment_vm.test_cdrom", map[string]string{
					"cdrom.0.file_id":   "none",
					"cdrom.0.interface": "ide2",
				}),
			},
			{
				RefreshState: true,
			},
		}},
		{"sata cdrom", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_cdrom" {
					node_name = "{{.NodeName}}"
					started   = false
					name 	  = "test-cdrom"
					cdrom {
						file_id   = "none"
						interface = "sata3"	
					}
				}`),
				Check: ResourceAttributes("proxmox_virtual_environment_vm.test_cdrom", map[string]string{
					"cdrom.0.interface": "sata3",
				}),
			},
			{
				RefreshState: true,
			},
		}},
		{"scsi cdrom", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_cdrom" {
					node_name = "{{.NodeName}}"
					started   = false
					name 	  = "test-cdrom"
					cdrom {
						file_id   = "none"
						interface = "scsi5"	
					}
				}`),
				Check: ResourceAttributes("proxmox_virtual_environment_vm.test_cdrom", map[string]string{
					"cdrom.0.interface": "scsi5",
				}),
			},
			{
				RefreshState: true,
			},
		}},
		{"enable cdrom", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_cdrom" {
					node_name = "{{.NodeName}}"
					started   = false
					name 	  = "test-cdrom"
					cdrom {
						file_id   = "none"
						interface = "ide2"
					}
				}`),
				Check: ResourceAttributes("proxmox_virtual_environment_vm.test_cdrom", map[string]string{
					"cdrom.0.file_id":   "none",
					"cdrom.0.interface": "ide2",
				}),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_cdrom" {
					node_name = "{{.NodeName}}"
					started   = false
					name 	  = "test-cdrom"
					cdrom {
						file_id   = "cdrom"
						interface = "ide2"
					}
				}`),
				Check: ResourceAttributes("proxmox_virtual_environment_vm.test_cdrom", map[string]string{
					"cdrom.0.file_id":   "cdrom",
					"cdrom.0.interface": "ide2",
				}),
			},
		}},
		{"multiple cdroms", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_cdrom" {
					node_name = "{{.NodeName}}"
					started   = false
					name 	  = "test-cdrom"
					cdrom {
						file_id   = "none"
						interface = "ide2"
					}
					cdrom {
						file_id   = "none"
						interface = "sata3"
					}
				}`),
				Check: ResourceAttributes("proxmox_virtual_environment_vm.test_cdrom", map[string]string{
					"cdrom.#":           "2",
					"cdrom.0.interface": "ide2",
					"cdrom.0.file_id":   "none",
					"cdrom.1.interface": "sata3",
					"cdrom.1.file_id":   "none",
				}),
			},
			{
				RefreshState: true,
			},
		}},
		{"remove one cdrom from multiple", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_cdrom" {
					node_name = "{{.NodeName}}"
					started   = false
					name 	  = "test-cdrom"
					cdrom {
						file_id   = "none"
						interface = "ide2"
					}
					cdrom {
						file_id   = "none"
						interface = "sata3"
					}
				}`),
				Check: ResourceAttributes("proxmox_virtual_environment_vm.test_cdrom", map[string]string{
					"cdrom.#":           "2",
					"cdrom.0.interface": "ide2",
					"cdrom.1.interface": "sata3",
				}),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_cdrom" {
					node_name = "{{.NodeName}}"
					started   = false
					name 	  = "test-cdrom"
					cdrom {
						file_id   = "none"
						interface = "ide2"
					}
				}`),
				Check: ResourceAttributes("proxmox_virtual_environment_vm.test_cdrom", map[string]string{
					"cdrom.#":           "1",
					"cdrom.0.interface": "ide2",
					"cdrom.0.file_id":   "none",
				}),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_cdrom" {
					node_name = "{{.NodeName}}"
					started   = false
					name 	  = "test-cdrom"
					cdrom {
						file_id   = "none"
						interface = "ide2"
					}
				}`),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		}},
		{"clone with multiple cdroms", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "template_cdrom" {
					node_name = "{{.NodeName}}"
					started   = false
					name      = "test-cdrom-template"
					template  = true

					cdrom {
						file_id   = "none"
						interface = "ide2"
					}
					cdrom {
						file_id   = "none"
						interface = "sata3"
					}
				}
				resource "proxmox_virtual_environment_vm" "clone_cdrom" {
					node_name = "{{.NodeName}}"
					started   = false
					name      = "test-cdrom-clone"

					clone {
						vm_id = proxmox_virtual_environment_vm.template_cdrom.vm_id
					}

					cdrom {
						file_id   = "none"
						interface = "ide2"
					}
					cdrom {
						file_id   = "none"
						interface = "sata3"
					}
				}`),
				Check: ResourceAttributes("proxmox_virtual_environment_vm.clone_cdrom", map[string]string{
					"cdrom.#":           "2",
					"cdrom.0.interface": "ide2",
					"cdrom.0.file_id":   "none",
					"cdrom.1.interface": "sata3",
					"cdrom.1.file_id":   "none",
				}),
			},
			{
				RefreshState: true,
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "template_cdrom" {
					node_name = "{{.NodeName}}"
					started   = false
					name      = "test-cdrom-template"
					template  = true

					cdrom {
						file_id   = "none"
						interface = "ide2"
					}
					cdrom {
						file_id   = "none"
						interface = "sata3"
					}
				}
				resource "proxmox_virtual_environment_vm" "clone_cdrom" {
					node_name = "{{.NodeName}}"
					started   = false
					name      = "test-cdrom-clone"

					clone {
						vm_id = proxmox_virtual_environment_vm.template_cdrom.vm_id
					}

					cdrom {
						file_id   = "none"
						interface = "ide2"
					}
					cdrom {
						file_id   = "none"
						interface = "sata3"
					}
				}`),
				Check: ResourceAttributes("proxmox_virtual_environment_vm.clone_cdrom", map[string]string{
					"cdrom.#":           "2",
					"cdrom.0.interface": "ide2",
					"cdrom.1.interface": "sata3",
				}),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: te.AccProviders,
				Steps:                    tt.steps,
			})
		})
	}
}

func TestAccResourceVMCDROMImportMultiple(t *testing.T) {
	t.Parallel()

	te := InitEnvironment(t)
	testVMID := 100000 + rand.Intn(99999) //nolint:gosec

	te.AddTemplateVars(map[string]any{
		"TestVMID": testVMID,
	})

	config := te.RenderConfig(`
		resource "proxmox_virtual_environment_vm" "test_cdrom_import" {
			node_name = "{{.NodeName}}"
			started   = false
			vm_id     = {{.TestVMID}}
			name      = "test-cdrom-import"

			cdrom {
				file_id   = "none"
				interface = "ide2"
			}
			cdrom {
				file_id   = "none"
				interface = "sata3"
			}
		}`)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: ResourceAttributes("proxmox_virtual_environment_vm.test_cdrom_import", map[string]string{
					"cdrom.#":           "2",
					"cdrom.0.interface": "ide2",
					"cdrom.1.interface": "sata3",
				}),
			},
			{
				ResourceName:      "proxmox_virtual_environment_vm.test_cdrom_import",
				ImportState:       true,
				ImportStateVerify: false,
				ImportStateId:     fmt.Sprintf("%s/%d", te.NodeName, testVMID),
				ImportStateCheck: func(states []*terraform.InstanceState) error {
					if len(states) != 1 {
						return fmt.Errorf("expected 1 instance state, got %d", len(states))
					}

					attrs := states[0].Attributes
					if attrs["cdrom.#"] != "2" {
						return fmt.Errorf("expected cdrom.# = 2 after import, got %s", attrs["cdrom.#"])
					}

					importedCDROMs := map[string]string{
						attrs["cdrom.0.interface"]: attrs["cdrom.0.file_id"],
						attrs["cdrom.1.interface"]: attrs["cdrom.1.file_id"],
					}

					if importedCDROMs["ide2"] != "none" {
						return fmt.Errorf("expected ide2 cdrom file_id none after import, got %q", importedCDROMs["ide2"])
					}

					if importedCDROMs["sata3"] != "none" {
						return fmt.Errorf("expected sata3 cdrom file_id none after import, got %q", importedCDROMs["sata3"])
					}

					return nil
				},
			},
			{
				Config: config,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func TestAccResourceVMCDROMInterfaceMove(t *testing.T) {
	t.Parallel()

	te := InitEnvironment(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_cdrom_move" {
					node_name = "{{.NodeName}}"
					started   = false
					name      = "test-cdrom-move"

					cdrom {
						file_id   = "none"
						interface = "ide2"
					}
					cdrom {
						file_id   = "none"
						interface = "sata3"
					}
				}`),
				Check: ResourceAttributes("proxmox_virtual_environment_vm.test_cdrom_move", map[string]string{
					"cdrom.#":           "2",
					"cdrom.0.interface": "ide2",
					"cdrom.1.interface": "sata3",
				}),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_cdrom_move" {
					node_name = "{{.NodeName}}"
					started   = false
					name      = "test-cdrom-move"

					cdrom {
						file_id   = "none"
						interface = "scsi5"
					}
					cdrom {
						file_id   = "none"
						interface = "sata3"
					}
				}`),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("proxmox_virtual_environment_vm.test_cdrom_move", plancheck.ResourceActionUpdate),
					},
				},
				Check: ResourceAttributes("proxmox_virtual_environment_vm.test_cdrom_move", map[string]string{
					"cdrom.#":           "2",
					"cdrom.0.interface": "sata3",
					"cdrom.1.interface": "scsi5",
				}),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_cdrom_move" {
					node_name = "{{.NodeName}}"
					started   = false
					name      = "test-cdrom-move"

					cdrom {
						file_id   = "none"
						interface = "scsi5"
					}
					cdrom {
						file_id   = "none"
						interface = "sata3"
					}
				}`),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func TestAccResourceVMCDROMCloneMultiplePhysical(t *testing.T) {
	t.Parallel()

	te := InitEnvironment(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "template_cdrom_physical" {
					node_name = "{{.NodeName}}"
					started   = false
					name      = "test-cdrom-template-physical"
					template  = true

					cdrom {
						file_id   = "cdrom"
						interface = "ide2"
					}
					cdrom {
						file_id   = "cdrom"
						interface = "sata3"
					}
				}
				resource "proxmox_virtual_environment_vm" "clone_cdrom_physical" {
					node_name = "{{.NodeName}}"
					started   = false
					name      = "test-cdrom-clone-physical"

					clone {
						vm_id = proxmox_virtual_environment_vm.template_cdrom_physical.vm_id
					}

					cdrom {
						file_id   = "cdrom"
						interface = "ide2"
					}
					cdrom {
						file_id   = "cdrom"
						interface = "sata3"
					}
				}`),
				Check: ResourceAttributes("proxmox_virtual_environment_vm.clone_cdrom_physical", map[string]string{
					"cdrom.#":           "2",
					"cdrom.0.interface": "ide2",
					"cdrom.0.file_id":   "cdrom",
					"cdrom.1.interface": "sata3",
					"cdrom.1.file_id":   "cdrom",
				}),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "template_cdrom_physical" {
					node_name = "{{.NodeName}}"
					started   = false
					name      = "test-cdrom-template-physical"
					template  = true

					cdrom {
						file_id   = "cdrom"
						interface = "ide2"
					}
					cdrom {
						file_id   = "cdrom"
						interface = "sata3"
					}
				}
				resource "proxmox_virtual_environment_vm" "clone_cdrom_physical" {
					node_name = "{{.NodeName}}"
					started   = false
					name      = "test-cdrom-clone-physical"

					clone {
						vm_id = proxmox_virtual_environment_vm.template_cdrom_physical.vm_id
					}

					cdrom {
						file_id   = "cdrom"
						interface = "ide2"
					}
					cdrom {
						file_id   = "cdrom"
						interface = "sata3"
					}
				}`),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func TestAccResourceVMCDROMImportSubsetManaged(t *testing.T) {
	t.Parallel()

	te := InitEnvironment(t)
	testVMID := 100000 + rand.Intn(99999) //nolint:gosec

	te.AddTemplateVars(map[string]any{
		"TestVMID": testVMID,
	})

	createConfig := te.RenderConfig(`
		resource "proxmox_virtual_environment_vm" "test_cdrom_subset" {
			node_name = "{{.NodeName}}"
			started   = false
			vm_id     = {{.TestVMID}}
			name      = "test-cdrom-subset"

			cdrom {
				file_id   = "none"
				interface = "ide2"
			}
			cdrom {
				file_id   = "none"
				interface = "sata3"
			}
		}`)

	managedSubsetConfig := te.RenderConfig(`
		resource "proxmox_virtual_environment_vm" "test_cdrom_subset" {
			node_name = "{{.NodeName}}"
			started   = false
			vm_id     = {{.TestVMID}}
			name      = "test-cdrom-subset"

			cdrom {
				file_id   = "none"
				interface = "ide2"
			}
		}`)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: createConfig,
				Check: ResourceAttributes("proxmox_virtual_environment_vm.test_cdrom_subset", map[string]string{
					"cdrom.#":           "2",
					"cdrom.0.interface": "ide2",
					"cdrom.1.interface": "sata3",
				}),
			},
			{
				ResourceName:      "proxmox_virtual_environment_vm.test_cdrom_subset",
				ImportState:       true,
				ImportStateVerify: false,
				ImportStateId:     fmt.Sprintf("%s/%d", te.NodeName, testVMID),
				ImportStateCheck: func(states []*terraform.InstanceState) error {
					if len(states) != 1 {
						return fmt.Errorf("expected 1 instance state, got %d", len(states))
					}

					attrs := states[0].Attributes
					if attrs["cdrom.#"] != "2" {
						return fmt.Errorf("expected imported state to include both discovered cdroms, got %s", attrs["cdrom.#"])
					}

					return nil
				},
			},
			{
				Config: managedSubsetConfig,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: ResourceAttributes("proxmox_virtual_environment_vm.test_cdrom_subset", map[string]string{
					"cdrom.#":           "1",
					"cdrom.0.interface": "ide2",
					"cdrom.0.file_id":   "none",
				}),
			},
		},
	})
}

func TestAccResourceVMCDROMMixedUpdate(t *testing.T) {
	t.Parallel()

	te := InitEnvironment(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_cdrom_mixed_update" {
					node_name = "{{.NodeName}}"
					started   = false
					name      = "test-cdrom-mixed-update"

					cdrom {
						file_id   = "none"
						interface = "ide2"
					}
					cdrom {
						file_id   = "none"
						interface = "sata3"
					}
				}`),
				Check: ResourceAttributes("proxmox_virtual_environment_vm.test_cdrom_mixed_update", map[string]string{
					"cdrom.#":           "2",
					"cdrom.0.interface": "ide2",
					"cdrom.0.file_id":   "none",
					"cdrom.1.interface": "sata3",
					"cdrom.1.file_id":   "none",
				}),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_cdrom_mixed_update" {
					node_name = "{{.NodeName}}"
					started   = false
					name      = "test-cdrom-mixed-update"

					cdrom {
						file_id   = "cdrom"
						interface = "ide2"
					}
					cdrom {
						file_id   = "none"
						interface = "scsi5"
					}
				}`),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("proxmox_virtual_environment_vm.test_cdrom_mixed_update", plancheck.ResourceActionUpdate),
					},
				},
				Check: ResourceAttributes("proxmox_virtual_environment_vm.test_cdrom_mixed_update", map[string]string{
					"cdrom.#":           "2",
					"cdrom.0.interface": "ide2",
					"cdrom.0.file_id":   "cdrom",
					"cdrom.1.interface": "scsi5",
					"cdrom.1.file_id":   "none",
				}),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_cdrom_mixed_update" {
					node_name = "{{.NodeName}}"
					started   = false
					name      = "test-cdrom-mixed-update"

					cdrom {
						file_id   = "cdrom"
						interface = "ide2"
					}
					cdrom {
						file_id   = "none"
						interface = "scsi5"
					}
				}`),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func TestAccResourceVMCDROMRejectsInvalidQ35IDEInterface(t *testing.T) {
	t.Parallel()

	te := InitEnvironment(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_cdrom_q35_invalid_ide" {
					node_name = "{{.NodeName}}"
					started   = false
					name      = "test-cdrom-q35-invalid-ide"
					machine   = "q35"

					cdrom {
						file_id   = "none"
						interface = "ide3"
					}
				}`),
				ExpectError: regexp.MustCompile(`cdrom interface "ide3" is invalid for q35 machine type: only ide0 and ide2 are supported on the IDE bus`),
			},
		},
	})
}
