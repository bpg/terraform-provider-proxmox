//go:build all

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

func TestAccResourceSDNZoneEVPN(t *testing.T) {
	// Cannot run in parallel due to SDN applier functionality affecting global state

	te := test.InitEnvironment(t)

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{
			name: "create evpn zone with all values then minimal then all values again",
			steps: []resource.TestStep{
				{
					Config: te.RenderConfig(`
						resource "proxmox_virtual_environment_sdn_zone_evpn" "evpn_update_test" {
						  id                         = "evpntest"
						  controller                 = "evpnctl"
						  vrf_vxlan                  = 99999
						  advertise_subnets          = true
						  exit_nodes_local_routing   = true
						  disable_arp_nd_suppression = true
						  mtu                        = 1370
						  nodes                      = ["{{.NodeName}}"]
						  exit_nodes                 = ["{{.NodeName}}"]
						  depends_on = [
						    proxmox_virtual_environment_sdn_applier.finalizer
						  ]
						}

						resource "proxmox_virtual_environment_sdn_applier" "evpn_applier" {
						  lifecycle {
							  replace_triggered_by = [
                  proxmox_virtual_environment_sdn_zone_evpn.evpn_update_test,
					      ]
						  }

						  depends_on = [
						    proxmox_virtual_environment_sdn_zone_evpn.evpn_update_test
						  ]
						}

						resource "proxmox_virtual_environment_sdn_applier" "finalizer" {}
					`),
					Check: resource.ComposeTestCheckFunc(
						test.ResourceAttributes("proxmox_virtual_environment_sdn_zone_evpn.evpn_update_test", map[string]string{
							"id":                         "evpntest",
							"controller":                 "evpnctl",
							"vrf_vxlan":                  "99999",
							"advertise_subnets":          "true",
							"exit_nodes_local_routing":   "true",
							"disable_arp_nd_suppression": "true",
							"mtu":                        "1370",
							"nodes.#":                    "1",
							"exit_nodes.#":               "1",
							"pending":                    "true",
							"state":                      "new",
						}),
					),
				},
				{
					Config: te.RenderConfig(`
						resource "proxmox_virtual_environment_sdn_zone_evpn" "evpn_update_test" {
						  id         = "evpntest"
						  controller = "evpnctl"
						  vrf_vxlan  = 99998
						  mtu        = 1450
						  nodes      = []
						  depends_on = [
						    proxmox_virtual_environment_sdn_applier.finalizer
						  ]
						}

						resource "proxmox_virtual_environment_sdn_applier" "evpn_applier" {
						  lifecycle {
							  replace_triggered_by = [
                  proxmox_virtual_environment_sdn_zone_evpn.evpn_update_test,
					      ]
						  }

						  depends_on = [
						    proxmox_virtual_environment_sdn_zone_evpn.evpn_update_test
						  ]
						}

						resource "proxmox_virtual_environment_sdn_applier" "finalizer" {}
					`),
					Check: resource.ComposeTestCheckFunc(
						test.ResourceAttributes("proxmox_virtual_environment_sdn_zone_evpn.evpn_update_test", map[string]string{
							"id":                         "evpntest",
							"controller":                 "evpnctl",
							"vrf_vxlan":                  "99998",
							"mtu":                        "1450",
							"advertise_subnets":          "false",
							"exit_nodes_local_routing":   "false",
							"disable_arp_nd_suppression": "false",
							"nodes.#":                    "0",
							"exit_nodes.#":               "0",
							"pending":                    "true",
							"state":                      "changed",
						}),
					),
				},
				{
					Config: te.RenderConfig(`
						resource "proxmox_virtual_environment_sdn_zone_evpn" "evpn_update_test" {
						  id                         = "evpntest"
						  controller                 = "evpnctl"
						  vrf_vxlan                  = 99999
						  advertise_subnets          = true
						  exit_nodes_local_routing   = true
						  disable_arp_nd_suppression = true
						  mtu                        = 1500
						  nodes                      = ["{{.NodeName}}"]
						  exit_nodes                 = ["{{.NodeName}}"]
						  depends_on = [
						    proxmox_virtual_environment_sdn_applier.finalizer
						  ]
						}

						resource "proxmox_virtual_environment_sdn_applier" "evpn_applier" {
						  lifecycle {
							  replace_triggered_by = [
                  proxmox_virtual_environment_sdn_zone_evpn.evpn_update_test,
					      ]
						  }

						  depends_on = [
						    proxmox_virtual_environment_sdn_zone_evpn.evpn_update_test
						  ]
						}

						resource "proxmox_virtual_environment_sdn_applier" "finalizer" {}
					`),
					Check: resource.ComposeTestCheckFunc(
						test.ResourceAttributes("proxmox_virtual_environment_sdn_zone_evpn.evpn_update_test", map[string]string{
							"id":                         "evpntest",
							"controller":                 "evpnctl",
							"vrf_vxlan":                  "99999",
							"advertise_subnets":          "true",
							"exit_nodes_local_routing":   "true",
							"disable_arp_nd_suppression": "true",
							"mtu":                        "1500",
							"nodes.#":                    "1",
							"exit_nodes.#":               "1",
							"pending":                    "true",
							"state":                      "changed",
						}),
					),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: te.AccProviders,
				Steps:                    tt.steps,
			})
		})
	}
}
