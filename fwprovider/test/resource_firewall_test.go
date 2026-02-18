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
	"math/rand"
	"regexp"
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
					dport = "18006"
					proto = "tcp"
					comment = "PVE Admin Interface"
				}
			}`),
			Check: resource.ComposeTestCheckFunc(
				ResourceAttributes("proxmox_virtual_environment_firewall_rules.rules1", map[string]string{
					"rule.0.type":    "in",
					"rule.0.action":  "ACCEPT",
					"rule.0.iface":   "vmbr0",
					"rule.0.dport":   "18006",
					"rule.0.proto":   "tcp",
					"rule.0.comment": "PVE Admin Interface",
				}),
				NoResourceAttributesSet("proxmox_virtual_environment_firewall_rules.rules1", []string{
					"node_name",
				}),
			),
		}}},
		{"appending new rules", []resource.TestStep{{
			Config: te.RenderConfig(`
			resource "proxmox_virtual_environment_firewall_rules" "appending_new_rules" {
				rule {
					type    = "in"
					action  = "ACCEPT"
					dport   = "80"
					proto   = "tcp"
					comment = "Allow HTTP"
				}
			}`),
			Check: resource.ComposeTestCheckFunc(
				ResourceAttributes("proxmox_virtual_environment_firewall_rules.appending_new_rules", map[string]string{
					"rule.#":         "1",
					"rule.0.type":    "in",
					"rule.0.action":  "ACCEPT",
					"rule.0.dport":   "80",
					"rule.0.proto":   "tcp",
					"rule.0.comment": "Allow HTTP",
					"rule.0.pos":     "0",
				}),
			),
		}, {
			Config: te.RenderConfig(`
			resource "proxmox_virtual_environment_firewall_rules" "appending_new_rules" {
				rule {
					type    = "in"
					action  = "ACCEPT"
					dport   = "80"
					proto   = "tcp"
					comment = "Allow HTTP"
				}
				rule {
					type    = "in"
					action  = "ACCEPT"
					dport   = "443"
					proto   = "tcp"
					comment = "Allow HTTPS"
				}
			}`),
			Check: resource.ComposeTestCheckFunc(
				ResourceAttributes("proxmox_virtual_environment_firewall_rules.appending_new_rules", map[string]string{
					"rule.#":         "2",
					"rule.0.type":    "in",
					"rule.0.action":  "ACCEPT",
					"rule.0.dport":   "80",
					"rule.0.proto":   "tcp",
					"rule.0.comment": "Allow HTTP",
					"rule.0.pos":     "0",
					"rule.1.type":    "in",
					"rule.1.action":  "ACCEPT",
					"rule.1.dport":   "443",
					"rule.1.proto":   "tcp",
					"rule.1.comment": "Allow HTTPS",
					"rule.1.pos":     "1",
				}),
			),
		}}},
		{"prepending new rules", []resource.TestStep{{
			Config: te.RenderConfig(`
			resource "proxmox_virtual_environment_firewall_rules" "prepending_new_rules" {
				rule {
					type    = "in"
					action  = "ACCEPT"
					dport   = "80"
					proto   = "tcp"
					comment = "Allow HTTP"
				}
			}`),
			Check: resource.ComposeTestCheckFunc(
				ResourceAttributes("proxmox_virtual_environment_firewall_rules.prepending_new_rules", map[string]string{
					"rule.#":         "1",
					"rule.0.type":    "in",
					"rule.0.action":  "ACCEPT",
					"rule.0.dport":   "80",
					"rule.0.proto":   "tcp",
					"rule.0.comment": "Allow HTTP",
				}),
			),
		}, {
			Config: te.RenderConfig(`
			resource "proxmox_virtual_environment_firewall_rules" "prepending_new_rules" {
				rule {
					type    = "in"
					action  = "ACCEPT"
					dport   = "443"
					proto   = "tcp"
					comment = "Allow HTTPS"
				}
				rule {
					type    = "in"
					action  = "ACCEPT"
					dport   = "80"
					proto   = "tcp"
					comment = "Allow HTTP"
				}
			}`),
			Check: resource.ComposeTestCheckFunc(
				ResourceAttributes("proxmox_virtual_environment_firewall_rules.prepending_new_rules", map[string]string{
					"rule.#":         "2",
					"rule.0.type":    "in",
					"rule.0.action":  "ACCEPT",
					"rule.0.dport":   "443",
					"rule.0.proto":   "tcp",
					"rule.0.comment": "Allow HTTPS",
					"rule.0.pos":     "0",
					"rule.1.type":    "in",
					"rule.1.action":  "ACCEPT",
					"rule.1.dport":   "80",
					"rule.1.proto":   "tcp",
					"rule.1.comment": "Allow HTTP",
					"rule.1.pos":     "1",
				}),
			),
		}}},
		{"deleting all rules", []resource.TestStep{{
			Config: te.RenderConfig(`
			resource "proxmox_virtual_environment_firewall_rules" "deleting_all_rules" {
				rule {
					type    = "in"
					action  = "ACCEPT"
					dport   = "443"
					proto   = "tcp"
					comment = "Allow HTTPS"
				}
				rule {
					type    = "in"
					action  = "ACCEPT"
					dport   = "80"
					proto   = "tcp"
					comment = "Allow HTTP"
				}
			}`),
			Check: resource.ComposeTestCheckFunc(
				ResourceAttributes("proxmox_virtual_environment_firewall_rules.deleting_all_rules", map[string]string{
					"rule.#":         "2",
					"rule.0.type":    "in",
					"rule.0.action":  "ACCEPT",
					"rule.0.dport":   "443",
					"rule.0.proto":   "tcp",
					"rule.0.comment": "Allow HTTPS",
					"rule.0.pos":     "0",
					"rule.1.type":    "in",
					"rule.1.action":  "ACCEPT",
					"rule.1.dport":   "80",
					"rule.1.proto":   "tcp",
					"rule.1.comment": "Allow HTTP",
					"rule.1.pos":     "1",
				}),
			),
		}, {
			Config: te.RenderConfig(`
			resource "proxmox_virtual_environment_firewall_rules" "deleting_all_rules" {
			}`),
			Check: resource.ComposeTestCheckFunc(
				ResourceAttributes("proxmox_virtual_environment_firewall_rules.deleting_all_rules", map[string]string{
					"rule.#": "0",
				}),
			),
		}}},
		{"remove rules from the end", []resource.TestStep{{
			Config: te.RenderConfig(`
			resource "proxmox_virtual_environment_firewall_rules" "remove_rules_from_end" {
				rule {
					type    = "in"
					action  = "ACCEPT"
					dport   = "443"
					proto   = "tcp"
					comment = "Allow HTTPS"
				}
				rule {
					type    = "in"
					action  = "ACCEPT"
					dport   = "80"
					proto   = "tcp"
					comment = "Allow HTTP"
				}
			}`),
			Check: resource.ComposeTestCheckFunc(
				ResourceAttributes("proxmox_virtual_environment_firewall_rules.remove_rules_from_end", map[string]string{
					"rule.#":         "2",
					"rule.0.type":    "in",
					"rule.0.action":  "ACCEPT",
					"rule.0.dport":   "443",
					"rule.0.proto":   "tcp",
					"rule.0.comment": "Allow HTTPS",
					"rule.0.pos":     "0",
					"rule.1.type":    "in",
					"rule.1.action":  "ACCEPT",
					"rule.1.dport":   "80",
					"rule.1.proto":   "tcp",
					"rule.1.comment": "Allow HTTP",
					"rule.1.pos":     "1",
				}),
			),
		}, {
			Config: te.RenderConfig(`
			resource "proxmox_virtual_environment_firewall_rules" "remove_rules_from_end" {
				rule {
					type    = "in"
					action  = "ACCEPT"
					dport   = "443"
					proto   = "tcp"
					comment = "Allow HTTPS"
				}
			}`),
			Check: resource.ComposeTestCheckFunc(
				ResourceAttributes("proxmox_virtual_environment_firewall_rules.remove_rules_from_end", map[string]string{
					"rule.#":         "1",
					"rule.0.type":    "in",
					"rule.0.action":  "ACCEPT",
					"rule.0.dport":   "443",
					"rule.0.proto":   "tcp",
					"rule.0.comment": "Allow HTTPS",
					"rule.0.pos":     "0",
				}),
			),
		}}},
		{"remove rules from the beginning", []resource.TestStep{{
			Config: te.RenderConfig(`
			resource "proxmox_virtual_environment_firewall_rules" "remove_rules_from_beginning" {
				rule {
					type    = "in"
					action  = "ACCEPT"
					dport   = "443"
					proto   = "tcp"
					comment = "Allow HTTPS"
				}
				rule {
					type    = "in"
					action  = "ACCEPT"
					dport   = "80"
					proto   = "tcp"
					comment = "Allow HTTP"
				}
			}`),
			Check: resource.ComposeTestCheckFunc(
				ResourceAttributes("proxmox_virtual_environment_firewall_rules.remove_rules_from_beginning", map[string]string{
					"rule.#":         "2",
					"rule.0.type":    "in",
					"rule.0.action":  "ACCEPT",
					"rule.0.dport":   "443",
					"rule.0.proto":   "tcp",
					"rule.0.comment": "Allow HTTPS",
					"rule.0.pos":     "0",
					"rule.1.type":    "in",
					"rule.1.action":  "ACCEPT",
					"rule.1.dport":   "80",
					"rule.1.proto":   "tcp",
					"rule.1.comment": "Allow HTTP",
					"rule.1.pos":     "1",
				}),
			),
		}, {
			Config: te.RenderConfig(`
			resource "proxmox_virtual_environment_firewall_rules" "remove_rules_from_beginning" {
				rule {
					type    = "in"
					action  = "ACCEPT"
					dport   = "80"
					proto   = "tcp"
					comment = "Allow HTTP"
				}
			}`),
			Check: resource.ComposeTestCheckFunc(
				ResourceAttributes("proxmox_virtual_environment_firewall_rules.remove_rules_from_beginning", map[string]string{
					"rule.#":         "1",
					"rule.0.type":    "in",
					"rule.0.action":  "ACCEPT",
					"rule.0.dport":   "80",
					"rule.0.proto":   "tcp",
					"rule.0.comment": "Allow HTTP",
					"rule.0.pos":     "0",
				}),
			),
		}}},
		{"remove rule from the middle", []resource.TestStep{{
			Config: te.RenderConfig(`
			resource "proxmox_virtual_environment_firewall_rules" "remove_rule_from_middle" {
				rule {
					type    = "in"
					action  = "ACCEPT"
					dport   = "443"
					proto   = "tcp"
					comment = "Allow HTTPS"
				}
				rule {
					type    = "in"
					action  = "ACCEPT"
					dport   = "22"
					proto   = "tcp"
					comment = "Allow SSH"
				}
				rule {
					type    = "in"
					action  = "ACCEPT"
					dport   = "80"
					proto   = "tcp"
					comment = "Allow HTTP"
				}
			}`),
			Check: resource.ComposeTestCheckFunc(
				ResourceAttributes("proxmox_virtual_environment_firewall_rules.remove_rule_from_middle", map[string]string{
					"rule.#":         "3",
					"rule.0.type":    "in",
					"rule.0.action":  "ACCEPT",
					"rule.0.dport":   "443",
					"rule.0.proto":   "tcp",
					"rule.0.comment": "Allow HTTPS",
					"rule.0.pos":     "0",
					"rule.1.type":    "in",
					"rule.1.action":  "ACCEPT",
					"rule.1.dport":   "22",
					"rule.1.proto":   "tcp",
					"rule.1.comment": "Allow SSH",
					"rule.1.pos":     "1",
					"rule.2.type":    "in",
					"rule.2.action":  "ACCEPT",
					"rule.2.dport":   "80",
					"rule.2.proto":   "tcp",
					"rule.2.comment": "Allow HTTP",
				}),
			),
		}, {
			Config: te.RenderConfig(`
			resource "proxmox_virtual_environment_firewall_rules" "remove_rule_from_middle" {
				rule {
					type    = "in"
					action  = "ACCEPT"
					dport   = "443"
					proto   = "tcp"
					comment = "Allow HTTPS"
				}
				rule {
					type    = "in"
					action  = "ACCEPT"
					dport   = "80"
					proto   = "tcp"
					comment = "Allow HTTP"
				}
			}`),
			Check: resource.ComposeTestCheckFunc(
				ResourceAttributes("proxmox_virtual_environment_firewall_rules.remove_rule_from_middle", map[string]string{
					"rule.#":         "2",
					"rule.0.type":    "in",
					"rule.0.action":  "ACCEPT",
					"rule.0.dport":   "443",
					"rule.0.proto":   "tcp",
					"rule.0.comment": "Allow HTTPS",
					"rule.0.pos":     "0",
					"rule.1.type":    "in",
					"rule.1.action":  "ACCEPT",
					"rule.1.dport":   "80",
					"rule.1.proto":   "tcp",
					"rule.1.comment": "Allow HTTP",
					"rule.1.pos":     "1",
				}),
			),
		}}},
		{"insert rule in the middle", []resource.TestStep{{
			Config: te.RenderConfig(`
			resource "proxmox_virtual_environment_firewall_rules" "insert_rule_in_middle" {
				rule {
					type    = "in"
					action  = "ACCEPT"
					dport   = "443"
					proto   = "tcp"
					comment = "Allow HTTPS"
				}
				rule {
					type    = "in"
					action  = "ACCEPT"
					dport   = "80"
					proto   = "tcp"
					comment = "Allow HTTP"
				}
			}`),
			Check: resource.ComposeTestCheckFunc(
				ResourceAttributes("proxmox_virtual_environment_firewall_rules.insert_rule_in_middle", map[string]string{
					"rule.#":         "2",
					"rule.0.type":    "in",
					"rule.0.action":  "ACCEPT",
					"rule.0.dport":   "443",
					"rule.0.proto":   "tcp",
					"rule.0.comment": "Allow HTTPS",
					"rule.0.pos":     "0",
					"rule.1.type":    "in",
					"rule.1.action":  "ACCEPT",
					"rule.1.dport":   "80",
					"rule.1.proto":   "tcp",
					"rule.1.comment": "Allow HTTP",
					"rule.1.pos":     "1",
				}),
			),
		}, {
			Config: te.RenderConfig(`
			resource "proxmox_virtual_environment_firewall_rules" "insert_rule_in_middle" {
				rule {
					type    = "in"
					action  = "ACCEPT"
					dport   = "443"
					proto   = "tcp"
					comment = "Allow HTTPS"
				}
				rule {
					type    = "in"
					action  = "ACCEPT"
					dport   = "22"
					proto   = "tcp"
					comment = "Allow SSH"
				}
				rule {
					type    = "in"
					action  = "ACCEPT"
					dport   = "80"
					proto   = "tcp"
					comment = "Allow HTTP"
				}
			}`),
			Check: resource.ComposeTestCheckFunc(
				ResourceAttributes("proxmox_virtual_environment_firewall_rules.insert_rule_in_middle", map[string]string{
					"rule.#":         "3",
					"rule.0.type":    "in",
					"rule.0.action":  "ACCEPT",
					"rule.0.dport":   "443",
					"rule.0.proto":   "tcp",
					"rule.0.comment": "Allow HTTPS",
					"rule.0.pos":     "0",
					"rule.1.type":    "in",
					"rule.1.action":  "ACCEPT",
					"rule.1.dport":   "22",
					"rule.1.proto":   "tcp",
					"rule.1.comment": "Allow SSH",
					"rule.1.pos":     "1",
					"rule.2.type":    "in",
					"rule.2.action":  "ACCEPT",
					"rule.2.dport":   "80",
					"rule.2.proto":   "tcp",
					"rule.2.comment": "Allow HTTP",
					"rule.2.pos":     "2",
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
							plancheck.ResourceActionUpdate),
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

	// NOTE: These tests are not run in parallel because they modify the same
	// shared cluster-level firewall resource.
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: te.AccProviders,
				Steps:                    tt.steps,
			})
		})
	}
}

func TestAccResourceNodeFirewallRules(t *testing.T) {
	te := InitEnvironment(t)

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"create rules", []resource.TestStep{{
			Config: te.RenderConfig(`
			resource "proxmox_virtual_environment_firewall_rules" "node_rules" {
				node_name = "{{.NodeName}}"

				rule {
					type   = "in"
					action = "ACCEPT"
					iface  = "vmbr0"
					dport = "18006"
					proto = "tcp"
					comment = "PVE Admin Interface"
				}
			}`),
			Check: resource.ComposeTestCheckFunc(
				ResourceAttributes("proxmox_virtual_environment_firewall_rules.node_rules", map[string]string{
					"node_name":      te.NodeName,
					"rule.0.type":    "in",
					"rule.0.action":  "ACCEPT",
					"rule.0.iface":   "vmbr0",
					"rule.0.dport":   "18006",
					"rule.0.proto":   "tcp",
					"rule.0.comment": "PVE Admin Interface",
				}),
			),
		}}},
		{"update rules", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_firewall_rules" "node_rules" {
					node_name = "{{.NodeName}}"

					rule {
						type    = "in"
						action  = "ACCEPT"
						comment = "Rule 0"
						dest    = "192.168.3.5"
						dport   = "443"
						proto   = "tcp"
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_firewall_rules.node_rules", map[string]string{
						"node_name":      te.NodeName,
						"rule.#":         "1",
						"rule.0.type":    "in",
						"rule.0.action":  "ACCEPT",
						"rule.0.comment": "Rule 0",
						"rule.0.dest":    "192.168.3.5",
						"rule.0.dport":   "443",
						"rule.0.proto":   "tcp",
					}),
				),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_firewall_rules" "node_rules" {
					node_name = "{{.NodeName}}"

					rule {
						type    = "in"
						action  = "ACCEPT"
						comment = "Rule 0"
						dest    = "192.168.3.5"
						dport   = "443"
						proto   = "tcp"
					}

					rule {
						type    = "in"
						action  = "ACCEPT"
						comment = "Rule 1"
						dest    = "192.168.3.6"
						dport   = "443"
						proto   = "tcp"
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_firewall_rules.node_rules", map[string]string{
						"node_name":      te.NodeName,
						"rule.#":         "2",
						"rule.0.type":    "in",
						"rule.0.action":  "ACCEPT",
						"rule.0.comment": "Rule 0",
						"rule.0.dest":    "192.168.3.5",
						"rule.0.dport":   "443",
						"rule.0.proto":   "tcp",
						"rule.0.pos":     "0",
						"rule.1.type":    "in",
						"rule.1.action":  "ACCEPT",
						"rule.1.comment": "Rule 1",
						"rule.1.dest":    "192.168.3.6",
						"rule.1.dport":   "443",
						"rule.1.proto":   "tcp",
						"rule.1.pos":     "1",
					}),
				),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_firewall_rules" "node_rules" {
					node_name = "{{.NodeName}}"

					rule {
						type    = "in"
						action  = "ACCEPT"
						comment = "Rule 1"
						dest    = "192.168.3.6"
						dport   = "443"
						proto   = "tcp"
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_firewall_rules.node_rules", map[string]string{
						"node_name":      te.NodeName,
						"rule.#":         "1",
						"rule.0.type":    "in",
						"rule.0.action":  "ACCEPT",
						"rule.0.comment": "Rule 1",
						"rule.0.dest":    "192.168.3.6",
						"rule.0.dport":   "443",
						"rule.0.proto":   "tcp",
						"rule.0.pos":     "0",
					}),
				),
			},
		}},
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

func TestAccResourceFirewallIPSetImport(t *testing.T) {
	te := InitEnvironment(t)

	// Generate dynamic VM and container IDs to avoid conflicts
	testVMID := 100000 + rand.Intn(99999)
	testContainerID := 100000 + rand.Intn(99999)

	te.AddTemplateVars(map[string]interface{}{
		"TestVMID":        testVMID,
		"TestContainerID": testContainerID,
	})

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"cluster rules import", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_firewall_ipset" "cluster_ipset" {
					name = "test"

					cidr {
						name    = "192.168.0.0/24"
						comment = "Local IPv4"
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_firewall_ipset.cluster_ipset", map[string]string{
						"name":           "test",
						"cidr.#":         "1",
						"cidr.0.name":    "192.168.0.0/24",
						"cidr.0.comment": "Local IPv4",
					}),
				),
			},
			{
				ResourceName:      "proxmox_virtual_environment_firewall_ipset.cluster_ipset",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     "cluster/test",
			},
		}},
		{"node rules import", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_firewall_rules" "node_rules" {
					node_name = "{{.NodeName}}"
					rule {
						type    = "in"
						action  = "ACCEPT"
						comment = "Allow SSH"
						dport   = "22"
						proto   = "tcp"
					}
					rule {
						type    = "in"
						action  = "ACCEPT"
						comment = "Allow HTTP"
						dport   = "80"
						proto   = "tcp"
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_firewall_rules.node_rules", map[string]string{
						"node_name":      te.NodeName,
						"rule.0.type":    "in",
						"rule.0.action":  "ACCEPT",
						"rule.0.comment": "Allow SSH",
						"rule.0.dport":   "22",
						"rule.0.proto":   "tcp",
					}),
				),
			},
			{
				ResourceName:      "proxmox_virtual_environment_firewall_rules.node_rules",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     fmt.Sprintf("node/%s", te.NodeName),
			},
		}},
		{"vm ipset import", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_firewall_ipset" "vm_ipset" {
					name = "test"

					node_name = "{{.NodeName}}"
					vm_id     = {{.TestVMID}}

					cidr {
						name    = "192.168.0.0/24"
						comment = "Local IPv4"
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_firewall_ipset.vm_ipset", map[string]string{
						"node_name":      te.NodeName,
						"vm_id":          fmt.Sprintf("%d", testVMID),
						"cidr.#":         "1",
						"cidr.0.name":    "192.168.0.0/24",
						"cidr.0.comment": "Local IPv4",
					}),
				),
			},
			{
				ResourceName:      "proxmox_virtual_environment_firewall_ipset.vm_ipset",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     fmt.Sprintf("vm/%s/%d/test", te.NodeName, testVMID),
			},
		}},
		{"container ipset import", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_firewall_ipset" "container_ipset" {
					name = "test"

					node_name     = "{{.NodeName}}"
					container_id  = {{.TestContainerID}}

					cidr {
						name    = "192.168.0.0/24"
						comment = "Local IPv4"
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_firewall_ipset.container_ipset", map[string]string{
						"node_name":      te.NodeName,
						"container_id":   fmt.Sprintf("%d", testContainerID),
						"cidr.#":         "1",
						"cidr.0.name":    "192.168.0.0/24",
						"cidr.0.comment": "Local IPv4",
					}),
				),
			},
			{
				ResourceName:      "proxmox_virtual_environment_firewall_ipset.container_ipset",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     fmt.Sprintf("container/%s/%d/test", te.NodeName, testContainerID),
			},
		}},
		{"invalid import ID", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_firewall_ipset" "test" {
					name = "invalid-import-id-test"

					cidr {
						name    = "192.168.0.0/24"
						comment = "Local IPv4"
					}
				}`),
			},
			{
				ResourceName:      "proxmox_virtual_environment_firewall_ipset.test",
				ImportState:       true,
				ImportStateVerify: false,
				ImportStateId:     "invalid-import-id",
				ExpectError:       regexp.MustCompile(`invalid import ID: .* \(expected: 'cluster/<ipset_name>', 'vm/<node_name>/<vm_id>/<ipset_name>', or 'container/<node_name>/<container_id>/<ipset_name>'\)`),
			},
		}},
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

func TestAccResourceFirewallRulesImport(t *testing.T) {
	te := InitEnvironment(t)

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"cluster rules import", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_firewall_rules" "cluster_rules" {
					rule {
						type    = "in"
						action  = "ACCEPT"
						comment = "Allow SSH"
						dport   = "22"
						proto   = "tcp"
					}
					rule {
						type    = "in"
						action  = "ACCEPT"
						comment = "Allow HTTP"
						dport   = "80"
						proto   = "tcp"
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_firewall_rules.cluster_rules", map[string]string{
						"rule.0.type":    "in",
						"rule.0.action":  "ACCEPT",
						"rule.0.comment": "Allow SSH",
						"rule.0.dport":   "22",
						"rule.0.proto":   "tcp",
					}),
				),
			},
			{
				ResourceName:      "proxmox_virtual_environment_firewall_rules.cluster_rules",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     "cluster",
			},
		}},
		{"vm rules import", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_firewall_rules" "vm_rules" {
					node_name = "{{.NodeName}}"
					vm_id     = 9997
					rule {
						type    = "in"
						action  = "ACCEPT"
						comment = "VM SSH Access"
						dport   = "22"
						proto   = "tcp"
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_firewall_rules.vm_rules", map[string]string{
						"node_name":      te.NodeName,
						"vm_id":          "9997",
						"rule.#":         "1",
						"rule.0.type":    "in",
						"rule.0.action":  "ACCEPT",
						"rule.0.comment": "VM SSH Access",
						"rule.0.dport":   "22",
						"rule.0.proto":   "tcp",
					}),
				),
			},
			{
				ResourceName:      "proxmox_virtual_environment_firewall_rules.vm_rules",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     fmt.Sprintf("vm/%s/9997", te.NodeName),
			},
		}},
		{"container rules import", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_firewall_rules" "container_rules" {
					node_name     = "{{.NodeName}}"
					container_id  = 9998
					rule {
						type    = "in"
						action  = "ACCEPT"
						comment = "Container HTTP Access"
						dport   = "80"
						proto   = "tcp"
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_firewall_rules.container_rules", map[string]string{
						"node_name":      te.NodeName,
						"container_id":   "9998",
						"rule.#":         "1",
						"rule.0.type":    "in",
						"rule.0.action":  "ACCEPT",
						"rule.0.comment": "Container HTTP Access",
						"rule.0.dport":   "80",
						"rule.0.proto":   "tcp",
					}),
				),
			},
			{
				ResourceName:      "proxmox_virtual_environment_firewall_rules.container_rules",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     fmt.Sprintf("container/%s/9998", te.NodeName),
			},
		}},
		{"invalid import ID", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_firewall_rules" "test" {
					rule {
						type   = "in"
						action = "ACCEPT"
					}
				}`),
			},
			{
				ResourceName:      "proxmox_virtual_environment_firewall_rules.test",
				ImportState:       true,
				ImportStateVerify: false,
				ImportStateId:     "invalid-import-id",
				ExpectError:       regexp.MustCompile("expected: 'cluster', 'vm/<node_name>/<vm_id>', or 'container/<node_name>/<container_id>'"),
			},
		}},
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

func TestAccResourceFirewallOptionsImport(t *testing.T) {
	te := InitEnvironment(t)

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"vm options import", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_firewall_options" "vm_options" {
					node_name = "{{.NodeName}}"
					vm_id     = 9999
					enabled   = true
					input_policy = "DROP"
					output_policy = "ACCEPT"
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_firewall_options.vm_options", map[string]string{
						"node_name":     te.NodeName,
						"vm_id":         "9999",
						"enabled":       "true",
						"input_policy":  "DROP",
						"output_policy": "ACCEPT",
					}),
				),
			},
			{
				ResourceName:      "proxmox_virtual_environment_firewall_options.vm_options",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     fmt.Sprintf("vm/%s/9999", te.NodeName),
			},
		}},
		{"container options import", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_firewall_options" "container_options" {
					node_name     = "{{.NodeName}}"
					container_id  = 10000
					enabled      = false
					input_policy = "ACCEPT"
					output_policy = "ACCEPT"
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_firewall_options.container_options", map[string]string{
						"node_name":     te.NodeName,
						"container_id":  "10000",
						"enabled":       "false",
						"input_policy":  "ACCEPT",
						"output_policy": "ACCEPT",
					}),
				),
			},
			{
				ResourceName:      "proxmox_virtual_environment_firewall_options.container_options",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     fmt.Sprintf("container/%s/10000", te.NodeName),
			},
		}},
		{"invalid options import ID", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_firewall_options" "test" {
					node_name = "{{.NodeName}}"
					vm_id     = 10001
					enabled   = true
				}`),
			},
			{
				ResourceName:      "proxmox_virtual_environment_firewall_options.test",
				ImportState:       true,
				ImportStateVerify: false,
				ImportStateId:     "invalid-options-import-id",
				ExpectError:       regexp.MustCompile("expected: 'vm/<node_name>/<vm_id>' or 'container/<node_name>/<container_id>'"),
			},
		}},
		{"missing vm_id and container_id", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_firewall_options" "test" {
					node_name = "{{.NodeName}}"
					enabled   = true
				}`),
				ExpectError: regexp.MustCompile("one of `container_id,vm_id` must be specified"),
			},
		}},
		{"both vm_id and container_id specified", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_firewall_options" "test" {
					node_name    = "{{.NodeName}}"
					vm_id        = 10001
					container_id = 10002
					enabled      = true
				}`),
				ExpectError: regexp.MustCompile("only one of `container_id,vm_id` can be specified"),
			},
		}},
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

// TestAccResourceFirewallRulesWithSecurityGroups tests inserting, prepending,
// and removing rules around security group references to verify bug #2575 is fixed.
// the bug caused security group references to be corrupted when rules were inserted
// before them.
func TestAccResourceFirewallRulesWithSecurityGroups(t *testing.T) {
	te := InitEnvironment(t)

	// generate unique security group names to avoid conflicts between test runs
	// proxmox limits security group names to 18 characters
	suffix := rand.Intn(100000)
	sgFoo := fmt.Sprintf("sg-foo-%05d", suffix)
	sgBar := fmt.Sprintf("sg-bar-%05d", suffix)

	te.AddTemplateVars(map[string]interface{}{
		"SecurityGroupFoo": sgFoo,
		"SecurityGroupBar": sgBar,
	})

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		// insert rule before security groups
		// initial: ssh + https + security_group foo + security_group bar
		// insert: icmp at position 2
		// verify: security groups are not corrupted and shift to positions 3 and 4
		{"insert before security groups", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_cluster_firewall_security_group" "test_sg_foo" {
					name    = "{{.SecurityGroupFoo}}"
					comment = "Test security group foo"
				}

				resource "proxmox_virtual_environment_cluster_firewall_security_group" "test_sg_bar" {
					name    = "{{.SecurityGroupBar}}"
					comment = "Test security group bar"
				}

				resource "proxmox_virtual_environment_firewall_rules" "test_sg_insert" {
					depends_on = [
						proxmox_virtual_environment_cluster_firewall_security_group.test_sg_foo,
						proxmox_virtual_environment_cluster_firewall_security_group.test_sg_bar,
					]

					rule {
						type    = "in"
						action  = "ACCEPT"
						comment = "SSH"
						dport   = "22"
						proto   = "tcp"
					}

					rule {
						type    = "in"
						action  = "ACCEPT"
						comment = "HTTPS"
						dport   = "443"
						proto   = "tcp"
					}

					rule {
						security_group = proxmox_virtual_environment_cluster_firewall_security_group.test_sg_foo.name
						comment        = "Security group foo"
					}

					rule {
						security_group = proxmox_virtual_environment_cluster_firewall_security_group.test_sg_bar.name
						comment        = "Security group bar"
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_firewall_rules.test_sg_insert", map[string]string{
						"rule.#":                "4",
						"rule.0.comment":        "SSH",
						"rule.0.pos":            "0",
						"rule.1.comment":        "HTTPS",
						"rule.1.pos":            "1",
						"rule.2.security_group": sgFoo,
						"rule.2.pos":            "2",
						"rule.3.security_group": sgBar,
						"rule.3.pos":            "3",
					}),
				),
			},
			{
				// insert icmp rule at position 2
				// before fix: security groups would be corrupted
				// after fix: security groups shift correctly to positions 3 and 4
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_cluster_firewall_security_group" "test_sg_foo" {
					name    = "{{.SecurityGroupFoo}}"
					comment = "Test security group foo"
				}

				resource "proxmox_virtual_environment_cluster_firewall_security_group" "test_sg_bar" {
					name    = "{{.SecurityGroupBar}}"
					comment = "Test security group bar"
				}

				resource "proxmox_virtual_environment_firewall_rules" "test_sg_insert" {
					depends_on = [
						proxmox_virtual_environment_cluster_firewall_security_group.test_sg_foo,
						proxmox_virtual_environment_cluster_firewall_security_group.test_sg_bar,
					]

					rule {
						type    = "in"
						action  = "ACCEPT"
						comment = "SSH"
						dport   = "22"
						proto   = "tcp"
					}

					rule {
						type    = "in"
						action  = "ACCEPT"
						comment = "HTTPS"
						dport   = "443"
						proto   = "tcp"
					}

					rule {
						type    = "in"
						action  = "ACCEPT"
						comment = "ICMP"
						proto   = "icmp"
					}

					rule {
						security_group = proxmox_virtual_environment_cluster_firewall_security_group.test_sg_foo.name
						comment        = "Security group foo"
					}

					rule {
						security_group = proxmox_virtual_environment_cluster_firewall_security_group.test_sg_bar.name
						comment        = "Security group bar"
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_firewall_rules.test_sg_insert", map[string]string{
						"rule.#":                "5",
						"rule.0.comment":        "SSH",
						"rule.0.pos":            "0",
						"rule.1.comment":        "HTTPS",
						"rule.1.pos":            "1",
						"rule.2.comment":        "ICMP",
						"rule.2.pos":            "2",
						"rule.3.security_group": sgFoo, // must be foo, not bar!
						"rule.3.pos":            "3",
						"rule.4.security_group": sgBar, // must be bar, not duplicate!
						"rule.4.pos":            "4",
					}),
				),
			},
		}},
		// prepend rule before security groups
		// initial: security_group foo + security_group bar
		// insert: ssh rule at position 0
		// verify: ssh at pos 0, foo at pos 1, bar at pos 2
		{"prepend before security groups", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_cluster_firewall_security_group" "test_sg_foo" {
					name    = "{{.SecurityGroupFoo}}"
					comment = "Test security group foo"
				}

				resource "proxmox_virtual_environment_cluster_firewall_security_group" "test_sg_bar" {
					name    = "{{.SecurityGroupBar}}"
					comment = "Test security group bar"
				}

				resource "proxmox_virtual_environment_firewall_rules" "test_sg_prepend" {
					depends_on = [
						proxmox_virtual_environment_cluster_firewall_security_group.test_sg_foo,
						proxmox_virtual_environment_cluster_firewall_security_group.test_sg_bar,
					]

					rule {
						security_group = proxmox_virtual_environment_cluster_firewall_security_group.test_sg_foo.name
						comment        = "Security group foo"
					}

					rule {
						security_group = proxmox_virtual_environment_cluster_firewall_security_group.test_sg_bar.name
						comment        = "Security group bar"
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_firewall_rules.test_sg_prepend", map[string]string{
						"rule.#":                "2",
						"rule.0.security_group": sgFoo,
						"rule.0.pos":            "0",
						"rule.1.security_group": sgBar,
						"rule.1.pos":            "1",
					}),
				),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_cluster_firewall_security_group" "test_sg_foo" {
					name    = "{{.SecurityGroupFoo}}"
					comment = "Test security group foo"
				}

				resource "proxmox_virtual_environment_cluster_firewall_security_group" "test_sg_bar" {
					name    = "{{.SecurityGroupBar}}"
					comment = "Test security group bar"
				}

				resource "proxmox_virtual_environment_firewall_rules" "test_sg_prepend" {
					depends_on = [
						proxmox_virtual_environment_cluster_firewall_security_group.test_sg_foo,
						proxmox_virtual_environment_cluster_firewall_security_group.test_sg_bar,
					]

					rule {
						type    = "in"
						action  = "ACCEPT"
						comment = "SSH"
						dport   = "22"
						proto   = "tcp"
					}

					rule {
						security_group = proxmox_virtual_environment_cluster_firewall_security_group.test_sg_foo.name
						comment        = "Security group foo"
					}

					rule {
						security_group = proxmox_virtual_environment_cluster_firewall_security_group.test_sg_bar.name
						comment        = "Security group bar"
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_firewall_rules.test_sg_prepend", map[string]string{
						"rule.#":                "3",
						"rule.0.comment":        "SSH",
						"rule.0.pos":            "0",
						"rule.1.security_group": sgFoo,
						"rule.1.pos":            "1",
						"rule.2.security_group": sgBar,
						"rule.2.pos":            "2",
					}),
				),
			},
		}},
		// remove rule before security groups
		// initial: ssh + https + icmp + security_group foo + security_group bar
		// remove icmp
		// verify: foo shifts from pos 3 to pos 2, bar shifts from pos 4 to pos 3
		{"remove before security groups", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_cluster_firewall_security_group" "test_sg_foo" {
					name    = "{{.SecurityGroupFoo}}"
					comment = "Test security group foo"
				}

				resource "proxmox_virtual_environment_cluster_firewall_security_group" "test_sg_bar" {
					name    = "{{.SecurityGroupBar}}"
					comment = "Test security group bar"
				}

				resource "proxmox_virtual_environment_firewall_rules" "test_sg_remove" {
					depends_on = [
						proxmox_virtual_environment_cluster_firewall_security_group.test_sg_foo,
						proxmox_virtual_environment_cluster_firewall_security_group.test_sg_bar,
					]

					rule {
						type    = "in"
						action  = "ACCEPT"
						comment = "SSH"
						dport   = "22"
						proto   = "tcp"
					}

					rule {
						type    = "in"
						action  = "ACCEPT"
						comment = "HTTPS"
						dport   = "443"
						proto   = "tcp"
					}

					rule {
						type    = "in"
						action  = "ACCEPT"
						comment = "ICMP"
						proto   = "icmp"
					}

					rule {
						security_group = proxmox_virtual_environment_cluster_firewall_security_group.test_sg_foo.name
						comment        = "Security group foo"
					}

					rule {
						security_group = proxmox_virtual_environment_cluster_firewall_security_group.test_sg_bar.name
						comment        = "Security group bar"
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_firewall_rules.test_sg_remove", map[string]string{
						"rule.#":                "5",
						"rule.0.comment":        "SSH",
						"rule.0.pos":            "0",
						"rule.1.comment":        "HTTPS",
						"rule.1.pos":            "1",
						"rule.2.comment":        "ICMP",
						"rule.2.pos":            "2",
						"rule.3.security_group": sgFoo,
						"rule.3.pos":            "3",
						"rule.4.security_group": sgBar,
						"rule.4.pos":            "4",
					}),
				),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_cluster_firewall_security_group" "test_sg_foo" {
					name    = "{{.SecurityGroupFoo}}"
					comment = "Test security group foo"
				}

				resource "proxmox_virtual_environment_cluster_firewall_security_group" "test_sg_bar" {
					name    = "{{.SecurityGroupBar}}"
					comment = "Test security group bar"
				}

				resource "proxmox_virtual_environment_firewall_rules" "test_sg_remove" {
					depends_on = [
						proxmox_virtual_environment_cluster_firewall_security_group.test_sg_foo,
						proxmox_virtual_environment_cluster_firewall_security_group.test_sg_bar,
					]

					rule {
						type    = "in"
						action  = "ACCEPT"
						comment = "SSH"
						dport   = "22"
						proto   = "tcp"
					}

					rule {
						type    = "in"
						action  = "ACCEPT"
						comment = "HTTPS"
						dport   = "443"
						proto   = "tcp"
					}

					rule {
						security_group = proxmox_virtual_environment_cluster_firewall_security_group.test_sg_foo.name
						comment        = "Security group foo"
					}

					rule {
						security_group = proxmox_virtual_environment_cluster_firewall_security_group.test_sg_bar.name
						comment        = "Security group bar"
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_firewall_rules.test_sg_remove", map[string]string{
						"rule.#":                "4",
						"rule.0.comment":        "SSH",
						"rule.0.pos":            "0",
						"rule.1.comment":        "HTTPS",
						"rule.1.pos":            "1",
						"rule.2.security_group": sgFoo,
						"rule.2.pos":            "2",
						"rule.3.security_group": sgBar,
						"rule.3.pos":            "3",
					}),
				),
			},
		}},
		// simultaneous insert and delete around security groups
		// initial: ssh + https + security_group foo
		// changes: remove ssh , add icmp between https and foo
		// verify: https at pos 0, icmp at pos 1, foo at pos 2
		{"insert and delete around security groups", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_cluster_firewall_security_group" "test_sg_foo" {
					name    = "{{.SecurityGroupFoo}}"
					comment = "Test security group foo"
				}

				resource "proxmox_virtual_environment_firewall_rules" "test_sg_complex" {
					depends_on = [
						proxmox_virtual_environment_cluster_firewall_security_group.test_sg_foo,
					]

					rule {
						type    = "in"
						action  = "ACCEPT"
						comment = "SSH"
						dport   = "22"
						proto   = "tcp"
					}

					rule {
						type    = "in"
						action  = "ACCEPT"
						comment = "HTTPS"
						dport   = "443"
						proto   = "tcp"
					}

					rule {
						security_group = proxmox_virtual_environment_cluster_firewall_security_group.test_sg_foo.name
						comment        = "Security group foo"
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_firewall_rules.test_sg_complex", map[string]string{
						"rule.#":                "3",
						"rule.0.comment":        "SSH",
						"rule.0.pos":            "0",
						"rule.1.comment":        "HTTPS",
						"rule.1.pos":            "1",
						"rule.2.security_group": sgFoo,
						"rule.2.pos":            "2",
					}),
				),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_cluster_firewall_security_group" "test_sg_foo" {
					name    = "{{.SecurityGroupFoo}}"
					comment = "Test security group foo"
				}

				resource "proxmox_virtual_environment_firewall_rules" "test_sg_complex" {
					depends_on = [
						proxmox_virtual_environment_cluster_firewall_security_group.test_sg_foo,
					]

					rule {
						type    = "in"
						action  = "ACCEPT"
						comment = "HTTPS"
						dport   = "443"
						proto   = "tcp"
					}

					rule {
						type    = "in"
						action  = "ACCEPT"
						comment = "ICMP"
						proto   = "icmp"
					}

					rule {
						security_group = proxmox_virtual_environment_cluster_firewall_security_group.test_sg_foo.name
						comment        = "Security group foo"
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_firewall_rules.test_sg_complex", map[string]string{
						"rule.#":                "3",
						"rule.0.comment":        "HTTPS",
						"rule.0.pos":            "0",
						"rule.1.comment":        "ICMP",
						"rule.1.pos":            "1",
						"rule.2.security_group": sgFoo,
						"rule.2.pos":            "2",
					}),
				),
			},
		}},
	}

	// NOTE: these tests are not run in parallel because they modify
	// cluster-level firewall rules which are shared state.
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: te.AccProviders,
				Steps:                    tt.steps,
			})
		})
	}
}
