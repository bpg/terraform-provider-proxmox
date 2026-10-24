//go:build acceptance || all

//testacc:tier=medium
//testacc:resource=vm

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"slices"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
)

const (
	accVMCPUTypeResource = "proxmox_virtual_environment_vm.test_cpu_type"
	accCPUModelGenoaV2   = "EPYC-Genoa-v2"
	accCPUModelLegacy    = "qemu64"
	accCPUModelInvalid   = "not-a-qemu-cpu"
)

func TestAccResourceVMCPUType(t *testing.T) {
	t.Parallel()

	t.Run("EPYC-Genoa-v2", func(t *testing.T) {
		te := InitEnvironment(t)

		resource.ParallelTest(t, resource.TestCase{
			PreCheck: func() {
				requireAdvertisedCPUModel(t, te, accCPUModelGenoaV2)
			},
			ProtoV6ProviderFactories: te.AccProviders,
			Steps: []resource.TestStep{
				{
					Config: te.RenderConfig(cpuTypeVMConfig(accCPUModelGenoaV2)),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(accVMCPUTypeResource, "cpu.0.type", accCPUModelGenoaV2),
						checkVMCPUTypeOnPVE(te, accVMCPUTypeResource, accCPUModelGenoaV2),
					),
				},
				{
					Config: te.RenderConfig(cpuTypeVMConfig(accCPUModelLegacy)),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(accVMCPUTypeResource, "cpu.0.type", accCPUModelLegacy),
						checkVMCPUTypeOnPVE(te, accVMCPUTypeResource, accCPUModelLegacy),
					),
				},
				{
					Config: te.RenderConfig(cpuTypeVMConfig(accCPUModelGenoaV2)),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(accVMCPUTypeResource, "cpu.0.type", accCPUModelGenoaV2),
						checkVMCPUTypeOnPVE(te, accVMCPUTypeResource, accCPUModelGenoaV2),
					),
				},
				{
					ResourceName:      accVMCPUTypeResource,
					ImportState:       true,
					ImportStateVerify: true,
					ImportStateIdFunc: vmImportID(accVMCPUTypeResource),
				},
			},
		})
	})

	t.Run("invalid model", func(t *testing.T) {
		te := InitEnvironment(t)

		resource.ParallelTest(t, resource.TestCase{
			ProtoV6ProviderFactories: te.AccProviders,
			Steps: []resource.TestStep{{
				Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_vm" "test_cpu_type_invalid" {
						node_name = "{{.NodeName}}"
						started   = false
						name      = "test-cpu-type-invalid"

						cpu {
							type = "` + accCPUModelInvalid + `"
						}
					}`),
				ExpectError: regexp.MustCompile(`(?s)cputype\s+'` + regexp.QuoteMeta(accCPUModelInvalid) + `'\s+is\s+not\s+defined`),
			}},
		})
	})
}

func cpuTypeVMConfig(cpuType string) string {
	return `
		resource "proxmox_virtual_environment_vm" "test_cpu_type" {
			node_name = "{{.NodeName}}"
			started   = false
			name      = "test-cpu-type"

			cpu {
				type = "` + cpuType + `"
			}
		}`
}

func requireAdvertisedCPUModel(t *testing.T, te *Environment, model string) {
	t.Helper()

	ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
	defer cancel()

	var resp struct {
		Data []struct {
			Name string `json:"name"`
		} `json:"data"`
	}

	path := fmt.Sprintf("nodes/%s/capabilities/qemu/cpu", te.NodeName)

	err := te.Client().DoRequest(ctx, http.MethodGet, path, nil, &resp)
	if errors.Is(err, api.ErrResourceDoesNotExist) {
		t.Skipf("skipping: PVE node %q has no QEMU CPU model list; newer model acceptance requires %s", te.NodeName, model)
	}

	if err != nil {
		t.Fatalf("list QEMU CPU models on %s: %v", te.NodeName, err)
	}

	names := make([]string, len(resp.Data))
	for i, cpuModel := range resp.Data {
		names[i] = cpuModel.Name
	}

	if !slices.Contains(names, model) {
		t.Skipf("skipping: PVE node %q does not advertise %s; newer model acceptance requires that model", te.NodeName, model)
	}
}

func checkVMCPUTypeOnPVE(te *Environment, resourceName, want string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource %q not found in state", resourceName)
		}

		vmID, err := strconv.Atoi(rs.Primary.Attributes["vm_id"])
		if err != nil {
			return fmt.Errorf("failed to parse vm_id: %w", err)
		}

		vmConfig, err := te.NodeClient().VM(vmID).GetVM(context.Background())
		if err != nil {
			return fmt.Errorf("failed to get VM config: %w", err)
		}

		if vmConfig.CPUEmulation == nil {
			return fmt.Errorf("cpu is not set on PVE, want %q", want)
		}

		if vmConfig.CPUEmulation.Type != want {
			return fmt.Errorf("unexpected cpu type on PVE: got %q, want %q", vmConfig.CPUEmulation.Type, want)
		}

		return nil
	}
}

func vmImportID(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource %q not found in state", resourceName)
		}

		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["node_name"], rs.Primary.Attributes["vm_id"]), nil
	}
}
