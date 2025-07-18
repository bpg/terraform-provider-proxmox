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

func TestAccResourceSDNZoneVLAN(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"create and update VLAN zone", []resource.TestStep{{
			Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_sdn_zone_vlan" "zone_vlan" {
				  id    = "zoneV"
				  nodes = ["pve"]
				  mtu   = 1496
				  bridge = "vmbr0"
				}
			`),
		}, {
			Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_sdn_zone_vlan" "zone_vlan" {
				  id    = "zoneV"
				  nodes = ["pve"]
				  mtu   = 1495
				  bridge = "vmbr0"
				}
			`),
			ResourceName:      "proxmox_virtual_environment_sdn_zone_vlan.zone_vlan",
			ImportStateId:     "zoneV",
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
