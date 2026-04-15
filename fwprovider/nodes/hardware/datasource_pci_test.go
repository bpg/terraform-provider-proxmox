//go:build acceptance || all

//testacc:tier=light
//testacc:resource=hardwaremapping

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package hardware_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/test"
)

func TestAccDataSourceHardwarePCI(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"read pci devices with default blacklist", []resource.TestStep{{
			Config: te.RenderConfig(`data "proxmox_hardware_pci" "test" {
				node_name = "{{.NodeName}}"
			}`),
			Check: resource.ComposeTestCheckFunc(
				test.ResourceAttributesSet("data.proxmox_hardware_pci.test", []string{
					"node_name",
					"devices.#",
					"devices.0.id",
					"devices.0.class",
					"devices.0.vendor",
					"devices.0.device",
					"devices.0.iommu_group",
				}),
			),
		}}},
		{"read pci devices with no blacklist", []resource.TestStep{{
			Config: te.RenderConfig(`data "proxmox_hardware_pci" "test" {
				node_name = "{{.NodeName}}"
				pci_class_blacklist = []
			}`),
			Check: resource.ComposeTestCheckFunc(
				test.ResourceAttributesSet("data.proxmox_hardware_pci.test", []string{
					"devices.#",
					"devices.0.id",
				}),
			),
		}}},
		{"read pci devices with vendor filter", []resource.TestStep{{
			Config: te.RenderConfig(`data "proxmox_hardware_pci" "test" {
				node_name = "{{.NodeName}}"
				pci_class_blacklist = []
				filters = {
					vendor_id = "8086"
				}
			}`),
			Check: resource.ComposeTestCheckFunc(
				test.ResourceAttributesSet("data.proxmox_hardware_pci.test", []string{
					"devices.#",
				}),
			),
		}}},
		{"read pci devices with class filter for bridges", []resource.TestStep{{
			Config: te.RenderConfig(`data "proxmox_hardware_pci" "test" {
				node_name = "{{.NodeName}}"
				pci_class_blacklist = []
				filters = {
					class = "06"
				}
			}`),
			Check: resource.ComposeTestCheckFunc(
				test.ResourceAttributesSet("data.proxmox_hardware_pci.test", []string{
					"devices.#",
				}),
			),
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
