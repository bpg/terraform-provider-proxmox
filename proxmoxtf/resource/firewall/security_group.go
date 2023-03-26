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

	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/firewall"
)

const (
	dvGroupComment = ""

	mkGroupName    = "name"
	mkGroupComment = "comment"
)

func SecurityGroupSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
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
	}
}

func SecurityGroupCreate(ctx context.Context, fw *firewall.API, d *schema.ResourceData) diag.Diagnostics {
	comment := d.Get(mkGroupComment).(string)
	name := d.Get(mkGroupName).(string)

	body := &firewall.GroupCreateRequestBody{
		Comment: &comment,
		Group:   name,
	}

	err := fw.CreateGroup(ctx, body)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(name)

	diags := ruleCreate(d, func(body *firewall.RuleCreateRequestBody) error {
		body.Group = &name
		e := fw.CreateGroupRule(ctx, name, body)
		if e != nil {
			return fmt.Errorf("error creating rule: %w", e)
		}
		return nil
	})
	if diags.HasError() {
		return diags
	}

	// reset rules, we re-read them (with proper positions) from the API
	err = d.Set(mkRule, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	return SecurityGroupRead(ctx, fw, d)
}

func SecurityGroupRead(ctx context.Context, fw *firewall.API, d *schema.ResourceData) diag.Diagnostics {
	var diags diag.Diagnostics

	name := d.Id()

	allGroups, err := fw.ListGroups(ctx)
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

			err = readGroupRule(ctx, fw, name, pos, ruleMap)
			if err != nil {
				diags = append(diags, diag.FromErr(err)...)
			}
		}
	} else {
		ruleIDs, err := fw.GetGroupRules(ctx, name)
		if err != nil {
			if strings.Contains(err.Error(), "no such security group") {
				d.SetId("")
				return nil
			}
			return diag.FromErr(err)
		}
		for _, id := range ruleIDs {
			ruleMap := map[string]interface{}{}
			err = readGroupRule(ctx, fw, name, id.Pos, ruleMap)
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
	fw *firewall.API,
	group string,
	pos int,
	ruleMap map[string]interface{},
) error {
	rule, err := fw.GetGroupRule(ctx, group, pos)
	if err != nil {
		if strings.Contains(err.Error(), "no such security group") {
			return nil
		}
		return fmt.Errorf("error reading rule %d for group %s: %w", pos, group, err)
	}

	baseRuleToMap(&rule.BaseRule, ruleMap)

	// pos in the map should be int!
	ruleMap[mkRulePos] = pos
	ruleMap[mkRuleAction] = rule.Action
	ruleMap[mkRuleType] = rule.Type

	return nil
}

func SecurityGroupUpdate(ctx context.Context, fw *firewall.API, d *schema.ResourceData) diag.Diagnostics {
	comment := d.Get(mkGroupComment).(string)
	newName := d.Get(mkGroupName).(string)
	previousName := d.Id()

	body := &firewall.GroupUpdateRequestBody{
		Group:   newName,
		ReName:  &previousName,
		Comment: &comment,
	}

	err := fw.UpdateGroup(ctx, body)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(newName)

	diags := ruleUpdate(d, func(body *firewall.RuleUpdateRequestBody) error {
		body.Group = &newName
		e := fw.UpdateGroupRule(ctx, newName, *body.Pos, body)
		return fmt.Errorf("error updating rule: %w", e)
	})
	if diags.HasError() {
		return diags
	}

	return SecurityGroupRead(ctx, fw, d)
}

func SecurityGroupDelete(ctx context.Context, fw *firewall.API, d *schema.ResourceData) diag.Diagnostics {
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
		err := fw.DeleteGroupRule(ctx, group, pos)
		if err != nil {
			if strings.Contains(err.Error(), "no such security group") {
				break
			}
			return diag.FromErr(err)
		}
	}

	err := fw.DeleteGroup(ctx, group)
	if err != nil {
		if strings.Contains(err.Error(), "no such security group") {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
