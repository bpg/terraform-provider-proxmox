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
	"github.com/bpg/terraform-provider-proxmox/proxmoxtf"
)

const (
	dvGroupComment = ""

	mkGroupName    = "name"
	mkGroupComment = "comment"
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
			mkRule: {
				Type:        schema.TypeList,
				Description: "List of rules",
				Optional:    true,
				DefaultFunc: func() (interface{}, error) {
					return []interface{}{}, nil
				},
				ForceNew: true,
				Elem:     Rule(),
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

	diags := ruleCreate(d, func(body *firewall.RuleCreateRequestBody) error {
		body.Group = &name
		return veClient.API().Cluster().Firewall().CreateGroupRule(ctx, name, body)
	})
	if diags.HasError() {
		return diags
	}

	// reset rules, we re-read them (with proper positions) from the API
	err = d.Set(mkRule, nil)
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

	rules := d.Get(mkRule).([]interface{})
	//nolint:nestif
	if len(rules) > 0 {
		// We have rules in the state, so we need to read them from the API
		for _, v := range rules {
			ruleMap := v.(map[string]interface{})
			pos := ruleMap[mkRulePos].(int)

			err = readGroupRule(ctx, veClient, name, pos, ruleMap)
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
			err = readGroupRule(ctx, veClient, name, id.Pos, ruleMap)
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

	err = d.Set(mkRule, rules)
	diags = append(diags, diag.FromErr(err)...)

	return diags
}

func readGroupRule(
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

	baseRuleToMap(&rule.BaseRule, ruleMap)

	// pos in the map should be int!
	ruleMap[mkRulePos] = pos
	ruleMap[mkRuleAction] = rule.Action
	ruleMap[mkRuleType] = rule.Type

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
		Group:   newName,
		ReName:  &previousName,
		Comment: &comment,
	}

	err = veClient.API().Cluster().Firewall().UpdateGroup(ctx, body)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(newName)

	diags := ruleUpdate(d, func(body *firewall.RuleUpdateRequestBody) error {
		body.Group = &newName
		return veClient.API().Cluster().Firewall().UpdateGroupRule(ctx, newName, *body.Pos, body)
	})
	if diags.HasError() {
		return diags
	}

	return securityGroupRead(ctx, d, m)
}

func securityGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(proxmoxtf.ProviderConfiguration)
	veClient, err := config.GetVEClient()
	if err != nil {
		return diag.FromErr(err)
	}

	group := d.Id()

	rules := d.Get(mkRule).([]interface{})
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
