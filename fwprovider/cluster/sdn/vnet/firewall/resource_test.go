//go:build acceptance || all

//testacc:tier=light
//testacc:resource=sdn

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package firewall_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/test"
)

func TestAccResourceSDNVNetFirewallRules(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"manage ordered forward rules", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_sdn_zone_simple" "firewall" {
					id    = "fwzmain"
					nodes = ["{{.NodeName}}"]
					depends_on = [
						proxmox_sdn_applier.finalizer
					]
				}

				resource "proxmox_sdn_vnet" "firewall" {
					id   = "fwvmain"
					zone = proxmox_sdn_zone_simple.firewall.id
				}

				resource "proxmox_sdn_applier" "firewall" {
					depends_on = [
						proxmox_sdn_vnet.firewall
					]
				}

				resource "proxmox_sdn_applier" "finalizer" {}

				resource "proxmox_sdn_vnet_firewall_rules" "firewall" {
					vnet = proxmox_sdn_vnet.firewall.id
					depends_on = [
						proxmox_sdn_applier.firewall
					]
					rules = [
						{
							action  = "ACCEPT"
							comment = "Allow application HTTPS"
							dest    = "10.96.20.0/24"
							dport   = "443"
							log     = "info"
							proto   = "tcp"
							source  = "10.96.10.0/24"
						},
						{
							action  = "DROP"
							comment = "Deny remaining application traffic"
							dest    = "10.96.20.0/24"
							enabled = false
							source  = "10.96.10.0/24"
						},
					]
				}`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_sdn_vnet_firewall_rules.firewall", map[string]string{
						"id":              "fwvmain",
						"vnet":            "fwvmain",
						"rules.#":         "2",
						"rules.0.action":  "ACCEPT",
						"rules.0.comment": "Allow application HTTPS",
						"rules.0.dest":    "10.96.20.0/24",
						"rules.0.dport":   "443",
						"rules.0.enabled": "true",
						"rules.0.log":     "info",
						"rules.0.proto":   "tcp",
						"rules.0.source":  "10.96.10.0/24",
						"rules.1.action":  "DROP",
						"rules.1.comment": "Deny remaining application traffic",
						"rules.1.enabled": "false",
					}),
				),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_sdn_zone_simple" "firewall" {
					id    = "fwzmain"
					nodes = ["{{.NodeName}}"]
					depends_on = [
						proxmox_sdn_applier.finalizer
					]
				}

				resource "proxmox_sdn_vnet" "firewall" {
					id   = "fwvmain"
					zone = proxmox_sdn_zone_simple.firewall.id
				}

				resource "proxmox_sdn_applier" "firewall" {
					depends_on = [
						proxmox_sdn_vnet.firewall
					]
				}

				resource "proxmox_sdn_applier" "finalizer" {}

				resource "proxmox_sdn_vnet_firewall_rules" "firewall" {
					vnet = proxmox_sdn_vnet.firewall.id
					depends_on = [
						proxmox_sdn_applier.firewall
					]
					rules = [
						{
							action  = "ACCEPT"
							comment = "Allow DNS"
							dest    = "10.96.53.53"
							dport   = "53"
							proto   = "udp"
							source  = "10.96.10.0/24"
						},
						{
							action  = "ACCEPT"
							comment = "Allow production HTTPS"
							dest    = "10.96.20.0/24"
							dport   = "443"
							log     = "notice"
							proto   = "tcp"
							source  = "10.96.10.0/24"
						},
						{
							action = "DROP"
							dest   = "10.96.20.0/24"
							source = "10.96.10.0/24"
						},
					]
				}`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_sdn_vnet_firewall_rules.firewall", map[string]string{
						"id":              "fwvmain",
						"vnet":            "fwvmain",
						"rules.#":         "3",
						"rules.0.action":  "ACCEPT",
						"rules.0.comment": "Allow DNS",
						"rules.0.dport":   "53",
						"rules.0.proto":   "udp",
						"rules.1.action":  "ACCEPT",
						"rules.1.comment": "Allow production HTTPS",
						"rules.1.log":     "notice",
						"rules.2.action":  "DROP",
						"rules.2.enabled": "true",
					}),
					test.NoResourceAttributesSet("proxmox_sdn_vnet_firewall_rules.firewall", []string{
						"rules.2.comment",
						"rules.2.log",
					}),
				),
			},
			{
				PreConfig: func() {
					if err := deleteVNetFirewallRuleManually(te, "fwvmain", 1); err != nil {
						t.Errorf("Failed to manually delete VNet firewall rule: %v", err)
					}
				},
				Config: te.RenderConfig(`
				resource "proxmox_sdn_zone_simple" "firewall" {
					id    = "fwzmain"
					nodes = ["{{.NodeName}}"]
					depends_on = [
						proxmox_sdn_applier.finalizer
					]
				}

				resource "proxmox_sdn_vnet" "firewall" {
					id   = "fwvmain"
					zone = proxmox_sdn_zone_simple.firewall.id
				}

				resource "proxmox_sdn_applier" "firewall" {
					depends_on = [
						proxmox_sdn_vnet.firewall
					]
				}

				resource "proxmox_sdn_applier" "finalizer" {}

				resource "proxmox_sdn_vnet_firewall_rules" "firewall" {
					vnet = proxmox_sdn_vnet.firewall.id
					depends_on = [
						proxmox_sdn_applier.firewall
					]
					rules = [
						{
							action  = "ACCEPT"
							comment = "Allow DNS"
							dest    = "10.96.53.53"
							dport   = "53"
							proto   = "udp"
							source  = "10.96.10.0/24"
						},
						{
							action  = "ACCEPT"
							comment = "Allow production HTTPS"
							dest    = "10.96.20.0/24"
							dport   = "443"
							log     = "notice"
							proto   = "tcp"
							source  = "10.96.10.0/24"
						},
						{
							action = "DROP"
							dest   = "10.96.20.0/24"
							source = "10.96.10.0/24"
						},
					]
				}`),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(
							"proxmox_sdn_vnet_firewall_rules.firewall",
							plancheck.ResourceActionUpdate,
						),
						plancheck.ExpectNonEmptyPlan(),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_sdn_vnet_firewall_rules.firewall", map[string]string{
						"rules.#":         "3",
						"rules.0.comment": "Allow DNS",
						"rules.1.comment": "Allow production HTTPS",
						"rules.2.action":  "DROP",
					}),
				),
			},
			{
				ResourceName:      "proxmox_sdn_vnet_firewall_rules.firewall",
				ImportStateId:     "fwvmain",
				ImportState:       true,
				ImportStateVerify: true,
			},
		}},
		{"manage a security group reference", []resource.TestStep{{
			Config: te.RenderConfig(`
			resource "proxmox_sdn_zone_simple" "firewall_group" {
				id    = "fwzgrp"
				nodes = ["{{.NodeName}}"]
				depends_on = [
					proxmox_sdn_applier.finalizer_group
				]
			}

			resource "proxmox_sdn_vnet" "firewall_group" {
				id   = "fwvgrp"
				zone = proxmox_sdn_zone_simple.firewall_group.id
			}

			resource "proxmox_sdn_applier" "firewall_group" {
				depends_on = [
					proxmox_sdn_vnet.firewall_group
				]
			}

			resource "proxmox_sdn_applier" "finalizer_group" {}

			resource "proxmox_virtual_environment_cluster_firewall_security_group" "firewall_group" {
				name    = "fwvgroup"
				comment = "VNet firewall acceptance group"
				rule {
					action = "ACCEPT"
					type   = "in"
					proto  = "icmp"
				}
			}

			resource "proxmox_sdn_vnet_firewall_rules" "firewall_group" {
				vnet = proxmox_sdn_vnet.firewall_group.id
				depends_on = [
					proxmox_sdn_applier.firewall_group
				]
				rules = [{
					security_group = proxmox_virtual_environment_cluster_firewall_security_group.firewall_group.name
					comment        = "Apply shared monitoring policy"
				}]
			}`),
			Check: resource.ComposeTestCheckFunc(
				test.ResourceAttributes("proxmox_sdn_vnet_firewall_rules.firewall_group", map[string]string{
					"id":                     "fwvgrp",
					"vnet":                   "fwvgrp",
					"rules.#":                "1",
					"rules.0.security_group": "fwvgroup",
					"rules.0.comment":        "Apply shared monitoring policy",
					"rules.0.enabled":        "true",
				}),
			),
		}, {
			ResourceName:      "proxmox_sdn_vnet_firewall_rules.firewall_group",
			ImportStateId:     "fwvgrp",
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

func deleteVNetFirewallRuleManually(te *test.Environment, vnet string, position int) error {
	if err := te.ClusterClient().SDNVnets(vnet).Firewall().DeleteRule(context.Background(), position); err != nil {
		return fmt.Errorf("failed to manually delete VNet firewall rule %d from %q: %w", position, vnet, err)
	}

	return nil
}

func TestAccResourceSDNVNetFirewallRulesValidation(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)

	tests := []struct {
		name        string
		rule        string
		expectError *regexp.Regexp
	}{
		{
			name:        "missing action and security group",
			rule:        `{ comment = "invalid" }`,
			expectError: regexp.MustCompile(`Exactly one of .*action.*security_group`),
		},
		{
			name:        "action conflicts with security group",
			rule:        `{ action = "ACCEPT", security_group = "fwvgroup" }`,
			expectError: regexp.MustCompile(`Exactly one of .*action.*security_group`),
		},
		{
			name:        "security group conflicts with packet match",
			rule:        `{ security_group = "fwvgroup", source = "10.96.10.0/24" }`,
			expectError: regexp.MustCompile(`Packet match attributes cannot be used with security_group`),
		},
		{
			name:        "reject is invalid on the forward chain",
			rule:        `{ action = "REJECT" }`,
			expectError: regexp.MustCompile(`value must be one of: \["ACCEPT" "DROP"\]`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: te.AccProviders,
				Steps: []resource.TestStep{{
					PlanOnly: true,
					Config: te.RenderConfig(`
					resource "proxmox_sdn_vnet_firewall_rules" "validation" {
						vnet  = "fwvalid"
						rules = [` + tt.rule + `]
					}`),
					ExpectError: tt.expectError,
				}},
			})
		})
	}
}
