/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package firewall

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/bpg/terraform-provider-proxmox/proxmox/firewall"
	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

const (
	dvRuleComment = ""
	dvRuleDPort   = ""
	dvRuleDest    = ""
	dvRuleEnable  = true
	dvRuleIface   = ""
	dvRuleLog     = "nolog"
	dvRuleMacro   = ""
	dvRuleProto   = ""
	dvRuleSPort   = ""
	dvRuleSource  = ""

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
			ForceNew:    false,
		},
		mkRuleType: {
			Type:        schema.TypeString,
			Description: "Rule type ('in', 'out')",
			Required:    true,
			ForceNew:    false,
		},
		mkRuleComment: {
			Type:        schema.TypeString,
			Description: "Rule comment",
			Optional:    true,
			Default:     dvRuleComment,
			ForceNew:    false,
		},
		mkRuleDest: {
			Type: schema.TypeString,
			Description: "Restrict packet destination address. This can refer to a single IP address, an" +
				" IP set ('+ipsetname') or an IP alias definition. You can also specify an address range " +
				"like '20.34.101.207-201.3.9.99', or a list of IP addresses and networks (entries are " +
				"separated by comma). Please do not mix IPv4 and IPv6 addresses inside such lists.",
			Optional: true,
			Default:  dvRuleDest,
			ForceNew: false,
		},
		mkRuleDPort: {
			Type: schema.TypeString,
			Description: "Restrict TCP/UDP destination port. You can use service names or simple numbers " +
				"(0-65535), as defined in '/etc/services'. Port ranges can be specified with '\\d+:\\d+'," +
				" for example '80:85', and you can use comma separated list to match several ports or ranges.",
			Optional: true,
			Default:  dvRuleDPort,
			ForceNew: false,
		},
		mkRuleEnable: {
			Type:        schema.TypeBool,
			Description: "Enable rule",
			Optional:    true,
			Default:     dvRuleEnable,
			ForceNew:    false,
		},
		mkRuleIFace: {
			Type: schema.TypeString,
			Description: "Network interface name. You have to use network configuration key names for VMs" +
				" and containers ('net\\d+'). Host related rules can use arbitrary strings.",
			Optional: true,
			Default:  dvRuleIface,
			ForceNew: false,
		},
		mkRuleLog: {
			Type: schema.TypeString,
			Description: "Log level for this rule ('emerg', 'alert', 'crit', 'err', 'warning', 'notice'," +
				" 'info', 'debug', 'nolog')",
			Optional: true,
			Default:  dvRuleLog,
			ForceNew: false,
		},
		mkRuleMacro: {
			Type:        schema.TypeString,
			Description: "Use predefined standard macro",
			Optional:    true,
			Default:     dvRuleMacro,
			ForceNew:    false,
		},
		mkRuleProto: {
			Type: schema.TypeString,
			Description: "Restrict packet protocol. You can use protocol names or simple numbers " +
				"(0-255), as defined in '/etc/protocols'.",
			Optional: true,
			Default:  dvRuleProto,
			ForceNew: false,
		},
		mkRuleSource: {
			Type: schema.TypeString,
			Description: "Restrict packet source address. This can refer to a single IP address, an" +
				" IP set ('+ipsetname') or an IP alias definition. You can also specify an address range " +
				"like '20.34.101.207-201.3.9.99', or a list of IP addresses and networks (entries are " +
				"separated by comma). Please do not mix IPv4 and IPv6 addresses inside such lists.",
			Optional: true,
			Default:  dvRuleSource,
			ForceNew: false,
		},
		mkRuleSPort: {
			Type: schema.TypeString,
			Description: "Restrict TCP/UDP source port. You can use service names or simple numbers " +
				"(0-65535), as defined in '/etc/services'. Port ranges can be specified with '\\d+:\\d+'," +
				" for example '80:85', and you can use comma separated list to match several ports or ranges.",
			Optional: true,
			Default:  dvRuleSPort,
			ForceNew: false,
		},
	}
}

func ruleCreate(d *schema.ResourceData, apiCaller func(*firewall.RuleCreateRequestBody) error) diag.Diagnostics {
	var diags diag.Diagnostics

	rules := d.Get(mkRule).([]interface{})
	for i := len(rules) - 1; i >= 0; i-- {
		rule := rules[i].(map[string]interface{})

		ruleBody := firewall.RuleCreateRequestBody{
			BaseRule: *mapToBaseRule(rule),
			Action:   rule[mkRuleAction].(string),
			Type:     rule[mkRuleType].(string),
		}

		err := apiCaller(&ruleBody)
		if err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	return diags
}

func ruleUpdate(d *schema.ResourceData, apiCaller func(*firewall.RuleUpdateRequestBody) error) diag.Diagnostics {
	var diags diag.Diagnostics

	rules := d.Get(mkRule).([]interface{})
	for i := len(rules) - 1; i >= 0; i-- {
		rule := rules[i].(map[string]interface{})

		ruleBody := firewall.RuleUpdateRequestBody{
			BaseRule: *mapToBaseRule(rule),
		}

		pos := rule[mkRulePos].(int)
		if pos >= 0 {
			ruleBody.Pos = &pos
		}
		action := rule[mkRuleAction].(string)
		if action != "" {
			ruleBody.Action = &action
		}
		rType := rule[mkRuleType].(string)
		if rType != "" {
			ruleBody.Type = &rType
		}

		err := apiCaller(&ruleBody)
		if err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	return diags
}

func mapToBaseRule(rule map[string]interface{}) *firewall.BaseRule {
	baseRule := &firewall.BaseRule{}

	comment := rule[mkRuleComment].(string)
	if comment != "" {
		baseRule.Comment = &comment
	}
	dest := rule[mkRuleDest].(string)
	if dest != "" {
		baseRule.Dest = &dest
	}
	dport := rule[mkRuleDPort].(string)
	if dport != "" {
		baseRule.DPort = &dport
	}
	enable := rule[mkRuleEnable].(bool)
	if enable {
		enableBool := types.CustomBool(true)
		baseRule.Enable = &enableBool
	}
	iface := rule[mkRuleIFace].(string)
	if iface != "" {
		baseRule.IFace = &iface
	}
	log := rule[mkRuleLog].(string)
	if log != "" {
		baseRule.Log = &log
	}
	macro := rule[mkRuleMacro].(string)
	if macro != "" {
		baseRule.Macro = &macro
	}
	proto := rule[mkRuleProto].(string)
	if proto != "" {
		baseRule.Proto = &proto
	}
	source := rule[mkRuleSource].(string)
	if source != "" {
		baseRule.Source = &source
	}
	sport := rule[mkRuleSPort].(string)
	if sport != "" {
		baseRule.SPort = &sport
	}

	return baseRule
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
