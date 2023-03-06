/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package firewall

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/bpg/terraform-provider-proxmox/proxmox"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/firewall"
	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
	"github.com/bpg/terraform-provider-proxmox/proxmoxtf"
)

const (
	dvGroupComment = ""
	dvRuleComment  = ""
	dvRuleDPort    = ""
	dvRuleDest     = ""
	dvRuleEnable   = false
	dvRuleIface    = ""
	dvRuleLog      = "nolog"
	dvRuleMacro    = ""
	dvRuleProto    = ""
	dvRuleSPort    = ""
	dvRuleSource   = ""

	mkGroupName    = "name"
	mkGroupComment = "comment"
	mkGroupRule    = "rule"

	mkRuleAction  = "action"
	mkRuleComment = "comment"
	mkRuleDPort   = "dport"
	mkRuleDest    = "dest"
	mkRuleDigest  = "digest"
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

func SecurityGroup() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			mkGroupName: {
				Type:        schema.TypeString,
				Description: "Security group name",
				Required:    true,
				ForceNew:    false,
			},
			mkGroupComment: {
				Type:        schema.TypeString,
				Description: "Security group comment",
				Optional:    true,
				Default:     dvGroupComment,
			},
			mkGroupRule: {
				Type:        schema.TypeList,
				Description: "List of rules",
				Optional:    true,
				DefaultFunc: func() (interface{}, error) {
					return []interface{}{}, nil
				},
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
							Description: "Restrict packet destination address. This can refer to a single IP address, an" +
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
					},
				},
			},
		},
		CreateContext: securityGroupCreate,
		ReadContext:   securityGroupRead,
		UpdateContext: securityGroupUpdate,
		DeleteContext: securityGroupDelete,
	}
}

func securityGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(proxmoxtf.ProviderConfiguration)
	veClient, err := config.GetVEClient()
	if err != nil {
		return diag.FromErr(err)
	}

	comment := d.Get(mkGroupComment).(string)
	name := d.Get(mkGroupName).(string)

	body := &firewall.GroupCreateRequestBody{
		Comment: &comment,
		Group:   name,
	}

	err = veClient.API().Cluster().Firewall().CreateGroup(ctx, body)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(name)

	var diags diag.Diagnostics

	rules := d.Get(mkGroupRule).([]interface{})
	for i := len(rules) - 1; i >= 0; i-- {
		rule := rules[i].(map[string]interface{})

		ruleBody := firewall.RuleCreateRequestBody{
			Action: rule[mkRuleAction].(string),
			Type:   rule[mkRuleType].(string),
			Group:  &name,
		}

		comment := rule[mkRuleComment].(string)
		if comment != "" {
			ruleBody.Comment = &comment
		}
		dest := rule[mkRuleDest].(string)
		if dest != "" {
			ruleBody.Dest = &dest
		}
		dport := rule[mkRuleDPort].(string)
		if dport != "" {
			ruleBody.DPort = &dport
		}
		enable := rule[mkRuleEnable].(bool)
		if enable {
			enableBool := types.CustomBool(true)
			ruleBody.Enable = &enableBool
		}
		iface := rule[mkRuleIFace].(string)
		if iface != "" {
			ruleBody.IFace = &iface
		}
		log := rule[mkRuleLog].(string)
		if log != "" {
			ruleBody.Log = &log
		}
		macro := rule[mkRuleMacro].(string)
		if macro != "" {
			ruleBody.Macro = &macro
		}
		pos := rule[mkRulePos].(int)
		if pos >= 0 {
			ruleBody.Pos = &pos
		}
		proto := rule[mkRuleProto].(string)
		if proto != "" {
			ruleBody.Proto = &proto
		}
		source := rule[mkRuleSource].(string)
		if source != "" {
			ruleBody.Source = &source
		}
		sport := rule[mkRuleSPort].(string)
		if sport != "" {
			ruleBody.SPort = &sport
		}

		err = veClient.API().Cluster().Firewall().CreateGroupRule(ctx, name, &ruleBody)
		if err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	if diags.HasError() {
		return diags
	}

	// reset rules, we re-read them from the API
	err = d.Set(mkGroupRule, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	return securityGroupRead(ctx, d, m)
}

func securityGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	config := m.(proxmoxtf.ProviderConfiguration)
	veClient, err := config.GetVEClient()
	if err != nil {
		return diag.FromErr(err)
	}

	name := d.Id()

	allGroups, err := veClient.API().Cluster().Firewall().ListGroups(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	for _, v := range allGroups {
		if v.Group == name {
			err = d.Set(mkGroupName, v.Group)
			diags = append(diags, diag.FromErr(err)...)
			err = d.Set(mkGroupComment, v.Comment)
			diags = append(diags, diag.FromErr(err)...)
			break
		}
	}

	rules := d.Get(mkGroupRule).([]interface{})
	//nolint:nestif
	if len(rules) > 0 {
		// We have rules in the state, so we need to read them from the API
		for _, v := range rules {
			ruleMap := v.(map[string]interface{})
			pos := ruleMap[mkRulePos].(int)

			err = readRule(ctx, veClient, name, pos, ruleMap)
			if err != nil {
				diags = append(diags, diag.FromErr(err)...)
			}
		}
	} else {
		ruleIDs, err := veClient.API().Cluster().Firewall().GetGroupRules(ctx, name)
		if err != nil {
			if strings.Contains(err.Error(), "HTTP 404") {
				d.SetId("")
				return nil
			}
			return diag.FromErr(err)
		}
		for _, id := range ruleIDs {
			ruleMap := map[string]interface{}{}
			err = readRule(ctx, veClient, name, id.Pos, ruleMap)
			if err != nil {
				diags = append(diags, diag.FromErr(err)...)
			} else {
				rules = append(rules, ruleMap)
			}
		}
	}

	if diags.HasError() {
		return diags
	}

	err = d.Set(mkGroupRule, rules)
	diags = append(diags, diag.FromErr(err)...)

	return diags
}

func readRule(
	ctx context.Context,
	client *proxmox.VirtualEnvironmentClient,
	group string,
	pos int,
	ruleMap map[string]interface{},
) error {
	rule, err := client.API().Cluster().Firewall().GetGroupRule(ctx, group, pos)
	if err != nil {
		return fmt.Errorf("error reading rule %d for group %s: %w", pos, group, err)
	}

	ruleMap[mkRulePos] = pos
	ruleMap[mkRuleAction] = rule.Action
	ruleMap[mkRuleType] = rule.Type

	if rule.Comment != nil {
		ruleMap[mkRuleComment] = *rule.Comment
	}
	if rule.Dest != nil {
		ruleMap[mkRuleDest] = *rule.Dest
	}
	if rule.DPort != nil {
		ruleMap[mkRuleDPort] = *rule.DPort
	}
	if rule.Enable != nil {
		ruleMap[mkRuleEnable] = *rule.Enable
	}
	if rule.IFace != nil {
		ruleMap[mkRuleIFace] = *rule.IFace
	}
	if rule.Log != nil {
		ruleMap[mkRuleLog] = *rule.Log
	}
	if rule.Macro != nil {
		ruleMap[mkRuleMacro] = *rule.Macro
	}
	if rule.Proto != nil {
		ruleMap[mkRuleProto] = *rule.Proto
	}
	if rule.Source != nil {
		ruleMap[mkRuleSource] = *rule.Source
	}
	if rule.SPort != nil {
		ruleMap[mkRuleSPort] = *rule.SPort
	}

	return nil
}

func securityGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(proxmoxtf.ProviderConfiguration)
	veClient, err := config.GetVEClient()
	if err != nil {
		return diag.FromErr(err)
	}

	comment := d.Get(mkGroupComment).(string)
	newName := d.Get(mkGroupName).(string)
	previousName := d.Id()

	body := &firewall.GroupUpdateRequestBody{
		ReName:  newName,
		Comment: &comment,
	}

	err = veClient.API().Cluster().Firewall().UpdateGroup(ctx, previousName, body)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(newName)

	return securityGroupRead(ctx, d, m)
}

func securityGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(proxmoxtf.ProviderConfiguration)
	veClient, err := config.GetVEClient()
	if err != nil {
		return diag.FromErr(err)
	}

	group := d.Id()

	rules := d.Get(mkGroupRule).([]interface{})
	sort.Slice(rules, func(i, j int) bool {
		ruleI := rules[i].(map[string]interface{})
		ruleJ := rules[j].(map[string]interface{})
		return ruleI[mkRulePos].(int) > ruleJ[mkRulePos].(int)
	})

	for _, v := range rules {
		rule := v.(map[string]interface{})
		pos := rule[mkRulePos].(int)
		err = veClient.API().Cluster().Firewall().DeleteGroupRule(ctx, group, pos)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	err = veClient.API().Cluster().Firewall().DeleteGroup(ctx, group)

	if err != nil {
		if strings.Contains(err.Error(), "HTTP 404") {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
