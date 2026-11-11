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
	"encoding/base64"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccResourceVMSMBIOS(t *testing.T) {
	t.Parallel()

	te := InitEnvironment(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_smbios" {
					node_name = "{{.NodeName}}"
					started   = false
					name      = "test-smbios"

					smbios {
						manufacturer = "test-manufacturer"
						product      = "test-product"
						sku          = "test-sku"
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_vm.test_smbios", map[string]string{
						"smbios.0.manufacturer": "^test-manufacturer$",
						"smbios.0.product":      "^test-product$",
						"smbios.0.sku":          "^test-sku$",
					}),
					checkVMSMBIOSOnPVE(te, "proxmox_virtual_environment_vm.test_smbios", map[string]string{
						"manufacturer": "test-manufacturer",
						"product":      "test-product",
						"sku":          "test-sku",
					}),
				),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_smbios" {
					node_name = "{{.NodeName}}"
					started   = false
					name      = "test-smbios"

					smbios {
						manufacturer = "test-manufacturer"
						product      = "test-product"
						sku          = "test-sku-updated"
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_vm.test_smbios", map[string]string{
						"smbios.0.sku": "^test-sku-updated$",
					}),
					checkVMSMBIOSOnPVE(te, "proxmox_virtual_environment_vm.test_smbios", map[string]string{
						"manufacturer": "test-manufacturer",
						"product":      "test-product",
						"sku":          "test-sku-updated",
					}),
				),
			},
		},
	})
}

func checkVMSMBIOSOnPVE(te *Environment, resourceName string, want map[string]string) resource.TestCheckFunc {
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

		if vmConfig.SMBIOS == nil {
			return fmt.Errorf("smbios1 is not set on PVE")
		}

		fields := map[string]*string{
			"manufacturer": vmConfig.SMBIOS.Manufacturer,
			"product":      vmConfig.SMBIOS.Product,
			"sku":          vmConfig.SMBIOS.SKU,
		}

		for name, wantValue := range want {
			ptr := fields[name]
			if ptr == nil {
				return fmt.Errorf("smbios %s is not set on PVE, want %q", name, wantValue)
			}

			decoded, err := base64.StdEncoding.DecodeString(*ptr)
			if err != nil {
				return fmt.Errorf("smbios %s on PVE is not base64: %w", name, err)
			}

			if string(decoded) != wantValue {
				return fmt.Errorf("unexpected smbios %s on PVE: got %q, want %q", name, string(decoded), wantValue)
			}
		}

		return nil
	}
}
