//go:build acceptance || all

//testacc:tier=light
//testacc:resource=misc

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package cpumodel_test

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/test"
	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/cpumodels"
)

const res = "proxmox_custom_cpu_model.test"

func TestAccResourceCustomCPUModel(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"minimal model and import with or without prefix", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_custom_cpu_model" "test" {
					id             = "acc-cpu-min"
					reported_model = "kvm64"
				}`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes(res, exactly(map[string]string{
						"id":             "acc-cpu-min",
						"cputype":        "custom-acc-cpu-min",
						"reported_model": "kvm64",
						"hidden":         "false",
					})),
					test.NoResourceAttributesSet(res, []string{
						"flags", "phys_bits", "guest_phys_bits", "hv_vendor_id", "level",
					}),
					checkModelOnPVE(te, "acc-cpu-min", func(m *cpumodels.GetResponseData) error {
						return expectEqual("reported-model", m.ReportedModel, "kvm64")
					}),
				),
			},
			{
				ResourceName:      res,
				ImportStateId:     "acc-cpu-min",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// PVE accepts the name with the prefix too; import must normalize it.
				ResourceName:      res,
				ImportStateId:     "custom-acc-cpu-min",
				ImportState:       true,
				ImportStateVerify: true,
			},
		}},
		{"all attributes, update, and field deletion", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_custom_cpu_model" "test" {
					id              = "acc-cpu-all"
					reported_model  = "qemu64"
					flags           = ["+aes", "-pcid"]
					phys_bits       = "host"
					guest_phys_bits = 40
					hidden          = true
					hv_vendor_id    = "proxmox"
					level           = 30
				}`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes(res, exactly(map[string]string{
						"flags.#":         "2",
						"flags.0":         "+aes",
						"flags.1":         "-pcid",
						"phys_bits":       "host",
						"guest_phys_bits": "40",
						"hidden":          "true",
						"hv_vendor_id":    "proxmox",
						"level":           "30",
					})),
					checkModelOnPVE(te, "acc-cpu-all", func(m *cpumodels.GetResponseData) error {
						return errors.Join(
							expectEqual("flags", m.Flags, "+aes;-pcid"),
							expectEqual("phys-bits", m.PhysBits, "host"),
							expectEqual("hv-vendor-id", m.HVVendorID, "proxmox"),
						)
					}),
				),
			},
			{
				ResourceName:      res,
				ImportStateId:     "acc-cpu-all",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_custom_cpu_model" "test" {
					id              = "acc-cpu-all"
					reported_model  = "kvm64"
					flags           = ["-pcid"]
					phys_bits       = "48"
					guest_phys_bits = 48
					hidden          = false
					hv_vendor_id    = "homelab"
					level           = 13
				}`),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(res, plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes(res, exactly(map[string]string{
						"reported_model":  "kvm64",
						"flags.#":         "1",
						"flags.0":         "-pcid",
						"phys_bits":       "48",
						"guest_phys_bits": "48",
						"hidden":          "false",
						"hv_vendor_id":    "homelab",
						"level":           "13",
					})),
					checkModelOnPVE(te, "acc-cpu-all", func(m *cpumodels.GetResponseData) error {
						return errors.Join(
							expectEqual("reported-model", m.ReportedModel, "kvm64"),
							expectEqual("flags", m.Flags, "-pcid"),
							expectEqual("phys-bits", m.PhysBits, "48"),
						)
					}),
				),
			},
			{
				// Dropping optional attributes must clear them on PVE, not leave them stale.
				Config: te.RenderConfig(`
				resource "proxmox_custom_cpu_model" "test" {
					id             = "acc-cpu-all"
					reported_model = "kvm64"
				}`),
				Check: resource.ComposeTestCheckFunc(
					test.NoResourceAttributesSet(res, []string{
						"flags", "phys_bits", "guest_phys_bits", "hv_vendor_id", "level",
					}),
					test.ResourceAttributes(res, exactly(map[string]string{"hidden": "false"})),
					checkModelOnPVE(te, "acc-cpu-all", func(m *cpumodels.GetResponseData) error {
						if m.Flags != nil || m.PhysBits != nil || m.GuestPhysBits != nil || m.HVVendorID != nil || m.Level != nil {
							return fmt.Errorf("optional fields still set on PVE: %+v", m)
						}

						return nil
					}),
				),
			},
		}},
		{"rename forces replacement", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_custom_cpu_model" "test" {
					id             = "acc-cpu-old"
					reported_model = "kvm64"
				}`),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_custom_cpu_model" "test" {
					id             = "acc-cpu-new"
					reported_model = "kvm64"
				}`),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(res, plancheck.ResourceActionReplace),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes(res, exactly(map[string]string{"cputype": "custom-acc-cpu-new"})),
					checkModelGoneFromPVE(te, "acc-cpu-old"),
				),
			},
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource.ParallelTest(t, resource.TestCase{
				ProtoV6ProviderFactories: te.AccProviders,
				CheckDestroy:             checkNoAccModelsOnPVE(te),
				Steps:                    tt.steps,
			})
		})
	}
}

func TestUnitResourceCustomCPUModel_Validators(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)

	tests := []struct {
		name   string
		config string
		error  string
	}{
		{"name with custom- prefix", `
			resource "proxmox_custom_cpu_model" "test" {
				id             = "custom-acc-cpu"
				reported_model = "kvm64"
			}`, "Invalid Custom CPU Model Name"},
		{"name starting with a digit", `
			resource "proxmox_custom_cpu_model" "test" {
				id             = "1cpu"
				reported_model = "kvm64"
			}`, "must start with a letter"},
		{"single-character name", `
			resource "proxmox_custom_cpu_model" "test" {
				id             = "a"
				reported_model = "kvm64"
			}`, "must start with a letter"},
		{"missing reported_model", `
			resource "proxmox_custom_cpu_model" "test" {
				id = "acc-cpu-bad"
			}`, `"reported_model" is required`},
		{"empty flags list", `
			resource "proxmox_custom_cpu_model" "test" {
				id             = "acc-cpu-bad"
				reported_model = "kvm64"
				flags          = []
			}`, "at least 1"},
		{"flag without sign", `
			resource "proxmox_custom_cpu_model" "test" {
				id             = "acc-cpu-bad"
				reported_model = "kvm64"
				flags          = ["aes"]
			}`, "followed by a flag name"},
		{"phys_bits out of range", `
			resource "proxmox_custom_cpu_model" "test" {
				id             = "acc-cpu-bad"
				reported_model = "kvm64"
				phys_bits      = "65"
			}`, "from 8 to 64"},
		{"guest_phys_bits out of range", `
			resource "proxmox_custom_cpu_model" "test" {
				id              = "acc-cpu-bad"
				reported_model  = "kvm64"
				guest_phys_bits = 31
			}`, "between 32 and 64"},
		{"hv_vendor_id too long", `
			resource "proxmox_custom_cpu_model" "test" {
				id             = "acc-cpu-bad"
				reported_model = "kvm64"
				hv_vendor_id   = "thirteenchars"
			}`, "1 to 12 alphanumeric"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: te.AccProviders,
				Steps: []resource.TestStep{{
					Config:      tt.config,
					PlanOnly:    true,
					ExpectError: regexp.MustCompile(tt.error),
				}},
			})
		})
	}
}

func checkModelOnPVE(
	te *test.Environment,
	name string,
	check func(*cpumodels.GetResponseData) error,
) resource.TestCheckFunc {
	return func(*terraform.State) error {
		m, err := te.ClusterClient().CustomCPUModels().Get(context.Background(), name)
		if err != nil {
			return err
		}

		return check(m)
	}
}

func checkModelGoneFromPVE(te *test.Environment, name string) resource.TestCheckFunc {
	return func(*terraform.State) error {
		_, err := te.ClusterClient().CustomCPUModels().Get(context.Background(), name)
		if errors.Is(err, api.ErrResourceDoesNotExist) {
			return nil
		}

		if err != nil {
			return err
		}

		return fmt.Errorf("custom CPU model %q still exists on PVE", name)
	}
}

func checkNoAccModelsOnPVE(te *test.Environment) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "proxmox_custom_cpu_model" {
				continue
			}

			if err := checkModelGoneFromPVE(te, rs.Primary.ID)(s); err != nil {
				return err
			}
		}

		return nil
	}
}

func expectEqual(field string, got *string, want string) error {
	if got == nil {
		return fmt.Errorf("%s on PVE: got nothing, want %q", field, want)
	}

	if *got != want {
		return fmt.Errorf("%s on PVE: got %q, want %q", field, *got, want)
	}

	return nil
}

// exactly anchors and escapes expected values for test.ResourceAttributes, which matches them as regular expressions.
func exactly(attrs map[string]string) map[string]string {
	anchored := make(map[string]string, len(attrs))

	for k, v := range attrs {
		anchored[k] = "^" + regexp.QuoteMeta(v) + "$"
	}

	return anchored
}
