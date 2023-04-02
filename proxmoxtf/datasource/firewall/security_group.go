/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package firewall

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/firewall"
)

const (
	mkSecurityGroupName    = "name"
	mkSecurityGroupComment = "comment"
	mkRules                = "rules"
)

func SecurityGroupSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		mkSecurityGroupName: {
			Type:        schema.TypeString,
			Description: "Security group name",
			Required:    true,
		},
		mkSecurityGroupComment: {
			Type:        schema.TypeString,
			Description: "Security group comment",
			Computed:    true,
		},
		mkRules: {
			Type:        schema.TypeList,
			Description: "List of rules",
			Computed:    true,
			Elem:        &schema.Resource{Schema: RuleSchema()},
		},
	}
}

func SecurityGroupRead(ctx context.Context, api firewall.SecurityGroup, d *schema.ResourceData) diag.Diagnostics {
	var diags diag.Diagnostics

	name := d.Get(mkSecurityGroupName).(string)

	allGroups, err := api.ListGroups(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	for _, v := range allGroups {
		if v.Group == name {
			err = d.Set(mkSecurityGroupName, v.Group)
			diags = append(diags, diag.FromErr(err)...)
			err = d.Set(mkSecurityGroupComment, v.Comment)
			diags = append(diags, diag.FromErr(err)...)
			break
		}
	}

	// rules := d.Get(mkRules).([]interface{})
	// ruleIDs, err := fw.ListGroupRules(ctx, name)
	// if err != nil {
	// 	if strings.Contains(err.Error(), "no such security group") {
	// 		d.SetId("")
	// 		return nil
	// 	}
	// 	return diag.FromErr(err)
	// }
	// for _, id := range ruleIDs {
	// 	ruleMap := map[string]interface{}{}
	// 	err = readGroupRule(ctx, fw, name, id.Pos, ruleMap)
	// 	if err != nil {
	// 		diags = append(diags, diag.FromErr(err)...)
	// 	} else {
	// 		rules = append(rules, ruleMap)
	// 	}
	// }

	// if diags.HasError() {
	// 	return diags
	// }

	// err = d.Set(mkRules, rules)
	// diags = append(diags, diag.FromErr(err)...)

	d.SetId(name)

	return diags
}

// func readGroupRule(
// 	ctx context.Context,
// 	fw firewall.API,
// 	group string,
// 	pos int,
// 	ruleMap map[string]interface{},
// ) error {
// 	rule, err := fw.GetGroupRule(ctx, group, pos)
// 	if err != nil {
// 		if strings.Contains(err.Error(), "no such security group") {
// 			return nil
// 		}
// 		return fmt.Errorf("error reading rule %d for group %s: %w", pos, group, err)
// 	}
//
// 	baseRuleToMap(&rule.BaseRule, ruleMap)
//
// 	// pos in the map should be int!
// 	ruleMap[mkRulePos] = pos
// 	ruleMap[mkRuleAction] = rule.Action
// 	ruleMap[mkRuleType] = rule.Type
//
// 	return nil
// }
