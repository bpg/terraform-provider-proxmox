/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package firewall

import (
	"context"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/bpg/terraform-provider-proxmox/proxmox/firewall"
	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/resource/validators"
	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/structure"
)

const (
	dvSecurityGroup = ""
	dvRuleComment   = ""
	dvRuleDPort     = ""
	dvRuleDest      = ""
	dvRuleEnabled   = true
	dvRuleIface     = ""
	dvRuleLog       = ""
	dvRuleMacro     = ""
	dvRuleProto     = ""
	dvRuleSPort     = ""
	dvRuleSource    = ""
	// MkRule defines the name of the rule resource in the schema.
	MkRule          = "rule"
	mkSecurityGroup = "security_group"
	mkRuleAction    = "action"
	mkRuleComment   = "comment"
	mkRuleDPort     = "dport"
	mkRuleDest      = "dest"
	mkRuleEnabled   = "enabled"
	mkRuleIFace     = "iface"
	mkRuleLog       = "log"
	mkRuleMacro     = "macro"
	mkRulePos       = "pos"
	mkRuleProto     = "proto"
	mkRuleSource    = "source"
	mkRuleSPort     = "sport"
	mkRuleType      = "type"
)

// Rules returns a resource that manages firewall rules.
func Rules() *schema.Resource {
	rule := map[string]*schema.Schema{
		mkRulePos: {
			Type:        schema.TypeInt,
			Description: "Rules position",
			Computed:    true,
		},
		mkSecurityGroup: {
			Type:        schema.TypeString,
			Description: "Security group name",
			Optional:    true,
			ForceNew:    true,
			Default:     dvSecurityGroup,
		},
		mkRuleAction: {
			Type:             schema.TypeString,
			Description:      "Rules action ('ACCEPT', 'DROP', 'REJECT')",
			Optional:         true,
			ValidateDiagFunc: validators.FirewallPolicy(),
		},
		mkRuleType: {
			Type:             schema.TypeString,
			Description:      "Rules type ('in', 'out')",
			Optional:         true,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"in", "out"}, true)),
		},
		mkRuleComment: {
			Type:        schema.TypeString,
			Description: "Rules comment",
			Optional:    true,
			Default:     dvRuleComment,
		},
		mkRuleDest: {
			Type: schema.TypeString,
			Description: "Restrict packet destination address. This can refer to a single IP address, an" +
				" IP set ('+ipsetname') or an IP alias definition. You can also specify an address range " +
				"like '20.34.101.207-201.3.9.99', or a list of IP addresses and networks (entries are " +
				"separated by comma). Please do not mix IPv4 and IPv6 addresses inside such lists.",
			Optional: true,
			Default:  dvRuleDest,
		},
		mkRuleDPort: {
			Type: schema.TypeString,
			Description: "Restrict TCP/UDP destination port. You can use service names or simple numbers " +
				"(0-65535), as defined in '/etc/services'. Port ranges can be specified with '\\d+:\\d+'," +
				" for example '80:85', and you can use comma separated list to match several ports or ranges.",
			Optional: true,
			Default:  dvRuleDPort,
		},
		mkRuleEnabled: {
			Type:        schema.TypeBool,
			Description: "Enable rule",
			Optional:    true,
			Default:     dvRuleEnabled,
		},
		mkRuleIFace: {
			Type: schema.TypeString,
			Description: "Network interface name. You have to use network configuration key names for VMs" +
				" and containers ('net\\d+'). Host related rules can use arbitrary strings.",
			Optional: true,
			Default:  dvRuleIface,
		},
		mkRuleLog: {
			Type: schema.TypeString,
			Description: "Log level for this rule ('emerg', 'alert', 'crit', 'err', 'warning', 'notice'," +
				" 'info', 'debug', 'nolog')",
			Optional: true,
			Default:  dvRuleLog,
		},
		mkRuleMacro: {
			Type:        schema.TypeString,
			Description: "Use predefined standard macro",
			Optional:    true,
			Default:     dvRuleMacro,
		},
		mkRuleProto: {
			Type: schema.TypeString,
			Description: "Restrict packet protocol. You can use protocol names or simple numbers " +
				"(0-255), as defined in '/etc/protocols'.",
			Optional: true,
			Default:  dvRuleProto,
		},
		mkRuleSource: {
			Type: schema.TypeString,
			Description: "Restrict packet source address. This can refer to a single IP address, an" +
				" IP set ('+ipsetname') or an IP alias definition. You can also specify an address range " +
				"like '20.34.101.207-201.3.9.99', or a list of IP addresses and networks (entries are " +
				"separated by comma). Please do not mix IPv4 and IPv6 addresses inside such lists.",
			Optional: true,
			Default:  dvRuleSource,
		},
		mkRuleSPort: {
			Type: schema.TypeString,
			Description: "Restrict TCP/UDP source port. You can use service names or simple numbers " +
				"(0-65535), as defined in '/etc/services'. Port ranges can be specified with '\\d+:\\d+'," +
				" for example '80:85', and you can use comma separated list to match several ports or ranges.",
			Optional: true,
			Default:  dvRuleSPort,
		},
	}

	s := map[string]*schema.Schema{
		MkRule: {
			Type:        schema.TypeList,
			Description: "List of rules",
			Required:    true,
			ForceNew:    true,
			Elem:        &schema.Resource{Schema: rule},
		},
	}

	structure.MergeSchema(s, selectorSchema())

	return &schema.Resource{
		Schema:        s,
		CreateContext: invokeRuleAPI(RulesCreate),
		ReadContext:   invokeRuleAPI(RulesRead),
		UpdateContext: invokeRuleAPI(RulesUpdate),
		DeleteContext: invokeRuleAPI(RulesDelete),
	}
}

// RulesCreate creates new firewall rules.
func RulesCreate(ctx context.Context, api firewall.Rule, d *schema.ResourceData) diag.Diagnostics {
	var diags diag.Diagnostics

	rules := d.Get(MkRule).([]interface{})

	for i := len(rules) - 1; i >= 0; i-- {
		var ruleBody firewall.RuleCreateRequestBody

		rule := rules[i].(map[string]interface{})

		sg := rule[mkSecurityGroup].(string)
		if sg != "" {
			// this is a special case of security group insertion
			ruleBody = firewall.RuleCreateRequestBody{
				Action:   sg,
				Type:     "group",
				BaseRule: *mapToSecurityGroupBaseRule(rule),
			}
		} else {
			a := rule[mkRuleAction].(string)
			t := rule[mkRuleType].(string)

			if a == "" || t == "" {
				diags = append(diags, diag.Errorf("Either '%s' OR both '%s' and '%s' must be defined for the rule #%d",
					mkSecurityGroup, mkRuleAction, mkRuleType, i)...)
				continue
			}

			ruleBody = firewall.RuleCreateRequestBody{
				Action:   a,
				Type:     t,
				BaseRule: *mapToBaseRule(rule),
			}
		}

		err := api.CreateRule(ctx, &ruleBody)
		diags = append(diags, diag.FromErr(err)...)
	}

	if diags.HasError() {
		return diags
	}

	d.SetId(api.GetRulesID())

	return RulesRead(ctx, api, d)
}

// RulesRead gets updates the tf state.
//
//nolint:revive,unused-parameter //bypass unused-parameter warning as ctx is needed for invokeRuleAPI.
func RulesRead(ctx context.Context, api firewall.Rule, d *schema.ResourceData) diag.Diagnostics {
	var diags diag.Diagnostics

	rules := d.Get(MkRule).([]interface{})
	newRules := make([]map[string]interface{}, 0)

	c := 0

	for _, rule := range rules {
		ruleMap := rule.(map[string]interface{})
		id := c
		ruleMap[mkRulePos] = id

		newRules = append(newRules, ruleMap)
		c++
	}

	err := d.Set(MkRule, newRules)
	diags = append(diags, diag.FromErr(err)...)

	return diags
}

// RulesUpdate updates rules.
func RulesUpdate(ctx context.Context, api firewall.Rule, d *schema.ResourceData) diag.Diagnostics {
	var diags diag.Diagnostics

	rules := d.Get(MkRule).([]interface{})
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

		err := api.UpdateRule(ctx, pos, &ruleBody)
		if err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	if diags.HasError() {
		return diags
	}

	return RulesRead(ctx, api, d)
}

// RulesDelete deletes all rules.
func RulesDelete(ctx context.Context, api firewall.Rule, d *schema.ResourceData) diag.Diagnostics {
	var diags diag.Diagnostics

	rules := d.Get(MkRule).([]interface{})
	sort.Slice(rules, func(i, j int) bool {
		ruleI := rules[i].(map[string]interface{})
		ruleJ := rules[j].(map[string]interface{})

		return ruleI[mkRulePos].(int) > ruleJ[mkRulePos].(int)
	})

	for _, v := range rules {
		rule := v.(map[string]interface{})
		pos := rule[mkRulePos].(int)

		_, err := api.GetRule(ctx, pos)
		if err != nil {
			// if the rule is not found / can't be retrieved, we can safely ignore it
			continue
		}

		err = api.DeleteRule(ctx, pos)
		diags = append(diags, diag.FromErr(err)...)
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

	enableBool := types.CustomBool(rule[mkRuleEnabled].(bool))
	baseRule.Enable = &enableBool

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

func mapToSecurityGroupBaseRule(rule map[string]interface{}) *firewall.BaseRule {
	baseRule := &firewall.BaseRule{}

	comment := rule[mkRuleComment].(string)
	if comment != "" {
		baseRule.Comment = &comment
	}

	enableBool := types.CustomBool(rule[mkRuleEnabled].(bool))
	baseRule.Enable = &enableBool

	iface := rule[mkRuleIFace].(string)
	if iface != "" {
		baseRule.IFace = &iface
	}

	return baseRule
}

func invokeRuleAPI(
	f func(context.Context, firewall.Rule, *schema.ResourceData) diag.Diagnostics,
) func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		return selectFirewallAPI(func(ctx context.Context, api firewall.API, data *schema.ResourceData) diag.Diagnostics {
			return f(ctx, api, data)
		})(ctx, d, m)
	}
}
