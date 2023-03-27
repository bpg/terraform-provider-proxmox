/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package firewall

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/firewall"
)

const (
	mkRule = "rule"

	mkRuleAction  = "action"
	mkRuleComment = "comment"
	mkRuleDPort   = "dport"
	mkRuleDest    = "dest"
	mkRuleEnable  = "enable"
	mkRuleIFace   = "iface"
	mkRuleLog     = "log"
	mkRuleMacro   = "macro"
	mkRulePos     = "pos"
	mkRuleProto   = "proto"
	mkRuleSource  = "source"
	mkRuleSPort   = "sport"
	mkRuleType    = "type"
)

func RuleSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		mkRulePos: {
			Type:        schema.TypeInt,
			Description: "Rule position",
			Computed:    true,
		},
		mkRuleAction: {
			Type:        schema.TypeString,
			Description: "Rule action ('ACCEPT', 'DROP', 'REJECT')",
			Required:    true,
		},
		mkRuleType: {
			Type:        schema.TypeString,
			Description: "Rule type ('in', 'out')",
			Required:    true,
		},
		mkRuleComment: {
			Type:        schema.TypeString,
			Description: "Rule comment",
			Computed:    true,
		},
		mkRuleDest: {
			Type:        schema.TypeString,
			Description: "Packet destination address",
			Computed:    true,
		},
		mkRuleDPort: {
			Type:        schema.TypeString,
			Description: "TCP/UDP destination port.",
			Computed:    true,
		},
		mkRuleEnable: {
			Type:        schema.TypeBool,
			Description: "Enable rule",
			Computed:    true,
		},
		mkRuleIFace: {
			Type:        schema.TypeString,
			Description: "Network interface name.",
			Computed:    true,
		},
		mkRuleLog: {
			Type:        schema.TypeString,
			Description: "Log level for this rule",
			Computed:    true,
		},
		mkRuleMacro: {
			Type:        schema.TypeString,
			Description: "Use predefined standard macro",
			Computed:    true,
		},
		mkRuleProto: {
			Type:        schema.TypeString,
			Description: "Packet protocol.",
			Computed:    true,
		},
		mkRuleSource: {
			Type:        schema.TypeString,
			Description: "Packet source address.",
			Computed:    true,
		},
		mkRuleSPort: {
			Type:        schema.TypeString,
			Description: "TCP/UDP source port.",
			Computed:    true,
		},
	}
}

func baseRuleToMap(baseRule *firewall.BaseRule, rule map[string]interface{}) {
	if baseRule.Comment != nil {
		rule[mkRuleComment] = *baseRule.Comment
	}
	if baseRule.Dest != nil {
		rule[mkRuleDest] = *baseRule.Dest
	}
	if baseRule.DPort != nil {
		rule[mkRuleDPort] = *baseRule.DPort
	}
	if baseRule.Enable != nil {
		rule[mkRuleEnable] = *baseRule.Enable
	}
	if baseRule.IFace != nil {
		rule[mkRuleIFace] = *baseRule.IFace
	}
	if baseRule.Log != nil {
		rule[mkRuleLog] = *baseRule.Log
	}
	if baseRule.Macro != nil {
		rule[mkRuleMacro] = *baseRule.Macro
	}
	if baseRule.Proto != nil {
		rule[mkRuleProto] = *baseRule.Proto
	}
	if baseRule.Source != nil {
		rule[mkRuleSource] = *baseRule.Source
	}
	if baseRule.SPort != nil {
		rule[mkRuleSPort] = *baseRule.SPort
	}
}
