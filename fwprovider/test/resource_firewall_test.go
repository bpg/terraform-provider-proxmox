//go:build acceptance || all

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestAccResourceClusterFirewall(t *testing.T) {
	te := InitEnvironment(t)

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"rules1", []resource.TestStep{{
			Config: te.RenderConfig(`
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
				ResourceAttributes("proxmox_virtual_environment_firewall_rules.rules1", map[string]string{
					"rule.0.type":    "in",
					"rule.0.action":  "ACCEPT",
					"rule.0.iface":   "vmbr0",
					"rule.0.dport":   "8006",
					"rule.0.proto":   "tcp",
					"rule.0.comment": "PVE Admin Interface",
				}),
				NoResourceAttributesSet("proxmox_virtual_environment_firewall_rules.rules1", []string{
					"node_name",
				}),
			),
		}}},
		{"rule attribute removal", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_firewall_rules" "attr_removal" {
					node_name = "{{.NodeName}}"
					vm_id     = 9995
					rule {
						type    = "in"
						action  = "ACCEPT"
						comment = "Allow HTTPS"
						dest    = "192.168.1.5"
						dport   = "443"
						proto   = "tcp"
						log     = "info"
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_firewall_rules.attr_removal", map[string]string{
						"rule.0.type":    "in",
						"rule.0.action":  "ACCEPT",
						"rule.0.comment": "Allow HTTPS",
						"rule.0.dest":    "192.168.1.5",
						"rule.0.dport":   "443",
						"rule.0.proto":   "tcp",
						"rule.0.log":     "info",
					}),
				),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_firewall_rules" "attr_removal" {
					node_name = "{{.NodeName}}"
					vm_id     = 9995
					rule {
						type   = "in"
						action = "ACCEPT"
						dest   = "192.168.1.5"
						log    = "info"
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_firewall_rules.attr_removal", map[string]string{
						"rule.0.type":    "in",
						"rule.0.action":  "ACCEPT",
						"rule.0.dest":    "192.168.1.5",
						"rule.0.log":     "info",
						"rule.0.comment": "",
						"rule.0.dport":   "",
						"rule.0.proto":   "",
					}),
				),
			},
		}},
		{"interface field removal", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_firewall_rules" "iface_removal" {
					node_name = "{{.NodeName}}"
					vm_id     = 9999
					rule {
						type   = "in"
						action = "ACCEPT"
						dest   = "192.168.1.10"
						iface  = "net0"
						log    = "info"
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_firewall_rules.iface_removal", map[string]string{
						"rule.0.type":   "in",
						"rule.0.action": "ACCEPT",
						"rule.0.dest":   "192.168.1.10",
						"rule.0.iface":  "net0",
						"rule.0.log":    "info",
					}),
				),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_firewall_rules" "iface_removal" {
					node_name = "{{.NodeName}}"
					vm_id     = 9999
					rule {
						type   = "in"
						action = "ACCEPT"
						dest   = "192.168.1.10"
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_firewall_rules.iface_removal", map[string]string{
						"rule.0.type":   "in",
						"rule.0.action": "ACCEPT",
						"rule.0.dest":   "192.168.1.10",
						"rule.0.iface":  "",
						"rule.0.log":    "",
					}),
				),
			},
		}},
		{"complete attribute removal", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_firewall_rules" "complete_removal" {
					node_name = "{{.NodeName}}"
					vm_id     = 9998
					rule {
						type    = "in"
						action  = "ACCEPT"
						comment = "Full rule"
						source  = "192.168.1.0/24"
						dest    = "192.168.2.0/24"
						sport   = "1024:65535"
						dport   = "80,443"
						proto   = "tcp"
						iface   = "net0"
						log     = "debug"
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_firewall_rules.complete_removal", map[string]string{
						"rule.0.type":    "in",
						"rule.0.action":  "ACCEPT",
						"rule.0.comment": "Full rule",
						"rule.0.source":  "192.168.1.0/24",
						"rule.0.dest":    "192.168.2.0/24",
						"rule.0.sport":   "1024:65535",
						"rule.0.dport":   "80,443",
						"rule.0.proto":   "tcp",
						"rule.0.iface":   "net0",
						"rule.0.log":     "debug",
					}),
				),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_firewall_rules" "complete_removal" {
					node_name = "{{.NodeName}}"
					vm_id     = 9998
					rule {
						type   = "in"
						action = "ACCEPT"
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_firewall_rules.complete_removal", map[string]string{
						"rule.0.type":    "in",
						"rule.0.action":  "ACCEPT",
						"rule.0.comment": "",
						"rule.0.source":  "",
						"rule.0.dest":    "",
						"rule.0.sport":   "",
						"rule.0.dport":   "",
						"rule.0.proto":   "",
						"rule.0.macro":   "",
						"rule.0.iface":   "",
						"rule.0.log":     "",
					}),
				),
			},
		}},
		{"multiple rules attribute removal", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_firewall_rules" "multi_removal" {
					node_name = "{{.NodeName}}"
					vm_id     = 9997
					rule {
						type    = "in"
						action  = "ACCEPT"
						comment = "HTTP Rule"
						dest    = "192.168.1.100"
						dport   = "80"
						proto   = "tcp"
					}
					rule {
						type    = "in"
						action  = "ACCEPT"
						comment = "SSH Rule"
						dest    = "192.168.1.101"
						dport   = "22"
						proto   = "tcp"
						log     = "info"
						iface   = "net0"
					}
				}`),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_firewall_rules" "multi_removal" {
					node_name = "{{.NodeName}}"
					vm_id     = 9997
					rule {
						type   = "in"
						action = "ACCEPT"
						dest   = "192.168.1.100"
					}
					rule {
						type   = "in"
						action = "ACCEPT"
						dest   = "192.168.1.101"
						log    = "info"
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_firewall_rules.multi_removal", map[string]string{
						"rule.0.type":    "in",
						"rule.0.action":  "ACCEPT",
						"rule.0.dest":    "192.168.1.100",
						"rule.0.comment": "",
						"rule.0.dport":   "",
						"rule.0.proto":   "",
						"rule.1.type":    "in",
						"rule.1.action":  "ACCEPT",
						"rule.1.dest":    "192.168.1.101",
						"rule.1.log":     "info",
						"rule.1.comment": "",
						"rule.1.dport":   "",
						"rule.1.proto":   "",
						"rule.1.iface":   "",
					}),
				),
			},
		}},
		{"attribute re-addition", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_firewall_rules" "re_addition" {
					node_name = "{{.NodeName}}"
					vm_id     = 9996
					rule {
						type    = "in"
						action  = "ACCEPT"
						comment = "Initial rule"
						dest    = "192.168.3.5"
						dport   = "443"
						proto   = "tcp"
					}
				}`),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_firewall_rules" "re_addition" {
					node_name = "{{.NodeName}}"
					vm_id     = 9996
					rule {
						type   = "in"
						action = "ACCEPT"
						dest   = "192.168.3.5"
					}
				}`),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_firewall_rules" "re_addition" {
					node_name = "{{.NodeName}}"
					vm_id     = 9996
					rule {
						type    = "in"
						action  = "ACCEPT"
						comment = "Re-added rule"
						dest    = "192.168.3.5"
						dport   = "8080"
						proto   = "tcp"
						log     = "warning"
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_firewall_rules.re_addition", map[string]string{
						"rule.0.type":    "in",
						"rule.0.action":  "ACCEPT",
						"rule.0.comment": "Re-added rule",
						"rule.0.dest":    "192.168.3.5",
						"rule.0.dport":   "8080",
						"rule.0.proto":   "tcp",
						"rule.0.log":     "warning",
					}),
				),
			},
		}},
		{"drift detection", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_firewall_rules" "drift_detection" {
					node_name = "{{.NodeName}}"
					vm_id     = 9994
					rule {
						type    = "in"
						action  = "ACCEPT"
						comment = "Allow HTTP"
						dest    = "192.168.1.5"
						dport   = "80"
						proto   = "tcp"
						log     = "info"
					}
					rule {
						type    = "in"
						action  = "ACCEPT"
						comment = "Allow HTTPS"
						dest    = "192.168.1.5"
						dport   = "443"
						proto   = "tcp"
						log     = "info"
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_firewall_rules.drift_detection", map[string]string{
						"rule.#":         "2",
						"rule.0.type":    "in",
						"rule.0.action":  "ACCEPT",
						"rule.0.comment": "Allow HTTP",
						"rule.0.dest":    "192.168.1.5",
						"rule.0.dport":   "80",
						"rule.0.proto":   "tcp",
						"rule.0.log":     "info",
						"rule.1.type":    "in",
						"rule.1.action":  "ACCEPT",
						"rule.1.comment": "Allow HTTPS",
						"rule.1.dest":    "192.168.1.5",
						"rule.1.dport":   "443",
						"rule.1.proto":   "tcp",
						"rule.1.log":     "info",
					}),
				),
			},
			{
				PreConfig: func() {
					err := deleteFirewallRuleManually(te, te.NodeName, 9994, 1)
					if err != nil {
						t.Errorf("Failed to manually delete rule: %v", err)
					}
				},
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_firewall_rules" "drift_detection" {
					node_name = "{{.NodeName}}"
					vm_id     = 9994
					rule {
						type    = "in"
						action  = "ACCEPT"
						comment = "Allow HTTP"
						dest    = "192.168.1.5"
						dport   = "80"
						proto   = "tcp"
						log     = "info"
					}
					rule {
						type    = "in"
						action  = "ACCEPT"
						comment = "Allow HTTPS"
						dest    = "192.168.1.5"
						dport   = "443"
						proto   = "tcp"
						log     = "info"
					}
				}`),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("proxmox_virtual_environment_firewall_rules.drift_detection",
							plancheck.ResourceActionDestroyBeforeCreate),
						plancheck.ExpectNonEmptyPlan(),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_firewall_rules.drift_detection", map[string]string{
						"rule.#":         "2",
						"rule.0.comment": "Allow HTTP",
						"rule.1.comment": "Allow HTTPS",
					}),
				),
			},
		}},
		{"ipset with ipV4 and ipV6 cidrs", []resource.TestStep{{
			Config: te.RenderConfig(`
			resource "proxmox_virtual_environment_firewall_ipset" "ipset" {
				name = "test"
				cidr {
					name    = "192.168.0.0/24"
					comment = "Local IPv4"
				}
				cidr {
					name    = "2001:db8:ab21:7b00::/64"
					comment = "LAN IPv6"
				}
				cidr {
					name    = "172.10.0.0/24"
					comment = "ext IPv4"
				}
				cidr {
					name    = "2001:db8:5a93:1e00::/64"
					comment = "ext IPv6"
				}
				cidr {
					name    = "2001:0DB8:91AA:7C30::1"
					comment = "ext 2 IPv6"
				}
			}`),
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

// deleteFirewallRuleManually simulates manual deletion of a firewall rule.
func deleteFirewallRuleManually(te *Environment, nodeName string, vmID int, rulePosition int) error {
	ctx := context.Background()
	firewallClient := te.NodeClient().VM(vmID).Firewall()

	err := firewallClient.DeleteRule(ctx, rulePosition)
	if err != nil {
		return fmt.Errorf("failed to manually delete firewall rule at position %d for VM %d on node %s: %w", rulePosition, vmID, nodeName, err)
	}

	return nil
}
