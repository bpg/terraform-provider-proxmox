//go:build acceptance || all

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package firewall_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/test"
)

func TestAccResourceNodeFirewallOptions(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: te.RenderConfig(`
						resource "proxmox_virtual_environment_node_firewall" "test" {
							node_name      = "{{.NodeName}}"
						}
					`),
			},
			{
				Config: te.RenderConfig(`
						resource "proxmox_virtual_environment_node_firewall" "test" {
							node_name                            = "{{.NodeName}}"
							enabled                              = true
							log_level_in                         = "err"
							log_level_out                        = "alert"
							log_level_forward                    = "warning"
							ndp                                  = true
							nf_conntrack_max                     = 999999999
							nf_conntrack_tcp_timeout_established = 999999999
							nftables                             = true
							nosmurfs                             = true
							smurf_log_level                      = "emerg"
							tcp_flags_log_level                  = "alert"
						}
					`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_virtual_environment_node_firewall.test", map[string]string{
						"enabled":                              "true",
						"log_level_in":                         "err",
						"log_level_out":                        "alert",
						"log_level_forward":                    "warning",
						"ndp":                                  "true",
						"nf_conntrack_max":                     "999999999",
						"nf_conntrack_tcp_timeout_established": "999999999",
						"nftables":                             "true",
						"nosmurfs":                             "true",
						"smurf_log_level":                      "emerg",
						"tcp_flags_log_level":                  "alert",
					}),
				),
			},
			{
				Config: te.RenderConfig(`
						resource "proxmox_virtual_environment_node_firewall" "test" {
							node_name           = "{{.NodeName}}"
							enabled             = true
							log_level_in        = "alert"
							log_level_out       = "alert"
							log_level_forward   = "alert"
							ndp                 = true
							nftables            = true
							nosmurfs            = true
							smurf_log_level     = "alert"
							tcp_flags_log_level = "alert"
						}
					`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_virtual_environment_node_firewall.test", map[string]string{
						"enabled":             "true",
						"log_level_in":        "alert",
						"log_level_out":       "alert",
						"log_level_forward":   "alert",
						"ndp":                 "true",
						"nftables":            "true",
						"nosmurfs":            "true",
						"smurf_log_level":     "alert",
						"tcp_flags_log_level": "alert",
					}),
				),
			},
			{
				ResourceName:      "proxmox_virtual_environment_node_firewall.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
