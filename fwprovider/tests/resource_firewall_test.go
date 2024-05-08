/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceClusterFirewall(t *testing.T) {
	te := initTestEnvironment(t)

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"rules1", []resource.TestStep{{
			Config: te.renderConfig(`
			resource "proxmox_virtual_environment_firewall_rules" "rules1" {
				rule {
					type   = "in"
					action = "ACCEPT"
					iface  = "vmbr0"
					dport = "8006"
					proto = "tcp"
					comment = "PVE Admin Interface"
				}
			}`),
			Check: resource.ComposeTestCheckFunc(
				testResourceAttributes("proxmox_virtual_environment_firewall_rules.rules1", map[string]string{
					"rule.0.type":    "in",
					"rule.0.action":  "ACCEPT",
					"rule.0.iface":   "vmbr0",
					"rule.0.dport":   "8006",
					"rule.0.proto":   "tcp",
					"rule.0.comment": "PVE Admin Interface",
				}),
				testNoResourceAttributesSet("proxmox_virtual_environment_firewall_rules.rules1", []string{
					"node_name",
				}),
			),
		}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource.ParallelTest(t, resource.TestCase{
				ProtoV6ProviderFactories: te.accProviders,
				Steps:                    tt.steps,
			})
		})
	}
}
