//go:build acceptance || all

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package zone_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/test"
)

func TestAccResourceSDNZoneVXLAN(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"create and update VXLAN zone", []resource.TestStep{{
			Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_sdn_zone_vxlan" "zone_vxlan" {
				  id    = "zoneX"
				  nodes = ["pve"]
				  mtu   = 1450
				  peers = ["10.0.0.1", "10.0.0.2"]
				}
			`),
		}, {
			Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_sdn_zone_vxlan" "zone_vxlan" {
				  id    = "zoneX"
				  nodes = ["pve"]
				  mtu   = 1440
				  peers = ["10.0.0.3", "10.0.0.4"]
				}
			`),
			ResourceName:      "proxmox_virtual_environment_sdn_zone_vxlan.zone_vxlan",
			ImportStateId:     "zoneX",
			ImportState:       true,
			ImportStateVerify: true,
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
