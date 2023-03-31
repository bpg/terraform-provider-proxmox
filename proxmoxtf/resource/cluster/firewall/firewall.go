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
	"github.com/bpg/terraform-provider-proxmox/proxmoxtf"
)

// const (
// 	dvLogRatelimiEnabled = true
// 	dvLogRatelimitBurst  = 5
// 	dvLogRatelimitRate   = "1/second"
// 	dvPolicyIn           = "DROP"
// 	dvPolicyOut          = "ACCEPT"
//
// 	mkEBTables            = "ebtables"
// 	mkEnabled             = "enabled"
// 	mkLogRatelimit        = "log_ratelimit"
// 	mkLogRatelimitEnabled = "enabled"
// 	mkLogRatelimitBurst   = "burst"
// 	mkLogRatelimitRate    = "rate"
// 	mkPolicyIn            = "input_policy"
// 	mkPolicyOut           = "output_policy"
// )

// func Firewall() *schema.Resource {
// 	return &schema.Resource{
// 		Schema: map[string]*schema.Schema{
// 			mkEBTables: {
// 				Type:        schema.TypeBool,
// 				Description: "Enable ebtables cluster-wide",
// 				Optional:    true,
// 			},
// 			mkEnabled: {
// 				Type:        schema.TypeBool,
// 				Description: "Enable or disable the firewall cluster-wide",
// 				Optional:    true,
// 			},
// 			mkLogRatelimit: {
// 				Type:        schema.TypeList,
// 				Description: "Log ratelimiting settings",
// 				Optional:    true,
// 				DefaultFunc: func() (interface{}, error) {
// 					return []interface{}{
// 						map[string]interface{}{
// 							mkLogRatelimitEnabled: dvLogRatelimiEnabled,
// 							mkLogRatelimitBurst:   dvLogRatelimitBurst,
// 							mkLogRatelimitRate:    dvLogRatelimitRate,
// 						},
// 					}, nil
// 				},
// 				Elem: &schema.Resource{
// 					Schema: map[string]*schema.Schema{
// 						mkLogRatelimitEnabled: {
// 							Type:        schema.TypeBool,
// 							Description: "Enable or disable log ratelimiting",
// 							Optional:    true,
// 							Default:     dvLogRatelimiEnabled,
// 						},
// 						mkLogRatelimitBurst: {
// 							Type:        schema.TypeInt,
// 							Description: "Initial burst of packages which will always get logged before the rate is applied",
// 							Optional:    true,
// 							Default:     dvLogRatelimitBurst,
// 						},
// 						mkLogRatelimitRate: {
// 							Type:             schema.TypeString,
// 							Description:      "Frequency with which the burst bucket gets refilled",
// 							Optional:         true,
// 							Default:          dvLogRatelimitRate,
// 							ValidateDiagFunc: validator.FirewallRate(),
// 						},
// 					},
// 				},
// 				MaxItems: 1,
// 				MinItems: 0,
// 			},
// 			mkPolicyIn: {
// 				Type:             schema.TypeString,
// 				Description:      "Default policy for incoming traffic",
// 				Optional:         true,
// 				Default:          dvPolicyIn,
// 				ValidateDiagFunc: validator.FirewallPolicy(),
// 			},
// 			mkPolicyOut: {
// 				Type:             schema.TypeString,
// 				Description:      "Default policy for outgoing traffic",
// 				Optional:         true,
// 				Default:          dvPolicyOut,
// 				ValidateDiagFunc: validator.FirewallPolicy(),
// 			},
// 			firewall.MkRule: {
// 				Type:        schema.TypeList,
// 				Description: "List of rules",
// 				Optional:    true,
// 				DefaultFunc: func() (interface{}, error) {
// 					return []interface{}{}, nil
// 				},
// 				Elem: &schema.Resource{Schema: firewall.RuleSchema()},
// 			},
// 		},
// 		CreateContext: invokeFirewallAPI(firewallCreate),
// 		ReadContext:   invokeFirewallAPI(firewallRead),
// 		UpdateContext: invokeFirewallAPI(firewallUpdate),
// 		DeleteContext: invokeFirewallAPI(firewallDelete),
// 	}
// }
//
// func firewallCreate(ctx context.Context, api fw.API, d *schema.ResourceData) diag.Diagnostics {
// 	diags := setOptions(ctx, api, d)
// 	if diags.HasError() {
// 		return diags
// 	}
//
// 	diags = firewall.RuleCreate(d, func(body *fw.RuleCreateRequestBody) error {
// 		e := api.CreateRule(ctx, body)
// 		if e != nil {
// 			return fmt.Errorf("error creating rule: %w", e)
// 		}
// 		return nil
// 	})
// 	if diags.HasError() {
// 		return diags
// 	}
//
// 	// reset rules, we re-read them (with proper positions) from the API
// 	err := d.Set(firewall.MkRule, nil)
// 	if err != nil {
// 		return diag.FromErr(err)
// 	}
//
// 	return firewallRead(ctx, api, d)
// }
//
// func setOptions(ctx context.Context, api fw.API, d *schema.ResourceData) diag.Diagnostics {
// 	policyIn := d.Get(mkPolicyIn).(string)
// 	policyOut := d.Get(mkPolicyOut).(string)
// 	body := &fw.OptionsPutRequestBody{
// 		PolicyIn:  &policyIn,
// 		PolicyOut: &policyOut,
// 	}
//
// 	logRatelimit := d.Get(mkLogRatelimit).([]interface{})
// 	if len(logRatelimit) > 0 {
// 		m := logRatelimit[0].(map[string]interface{})
// 		burst := m[mkLogRatelimitBurst].(int)
// 		rate := m[mkLogRatelimitRate].(string)
// 		rl := fw.CustomLogRateLimit{
// 			Enable: types.CustomBool(m[mkLogRatelimitEnabled].(bool)),
// 			Burst:  &burst,
// 			Rate:   &rate,
// 		}
// 		body.LogRateLimit = &rl
// 	}
//
// 	ebtablesBool := types.CustomBool(d.Get(mkEBTables).(bool))
// 	body.EBTables = &ebtablesBool
//
// 	enabledBool := types.CustomBool(d.Get(mkEnabled).(bool))
// 	body.Enable = &enabledBool
//
// 	err := api.SetOptions(ctx, body)
// 	return diag.FromErr(err)
// }
//
// func firewallRead(ctx context.Context, api fw.API, d *schema.ResourceData) diag.Diagnostics {
// 	var diags diag.Diagnostics
//
// 	options, err := api.GetOptions(ctx)
// 	if err != nil {
// 		return diag.FromErr(err)
// 	}
//
// 	if options.EBTables != nil {
// 		err = d.Set(mkEBTables, *options.EBTables)
// 		diags = append(diags, diag.FromErr(err)...)
// 	}
//
// 	if options.Enable != nil {
// 		err = d.Set(mkEnabled, *options.Enable)
// 		diags = append(diags, diag.FromErr(err)...)
// 	}
//
// 	if options.LogRateLimit != nil {
// 		err = d.Set(mkLogRatelimit, []interface{}{
// 			map[string]interface{}{
// 				mkLogRatelimitEnabled: options.LogRateLimit.Enable,
// 				mkLogRatelimitBurst:   *options.LogRateLimit.Burst,
// 				mkLogRatelimitRate:    *options.LogRateLimit.Rate,
// 			},
// 		})
// 		diags = append(diags, diag.FromErr(err)...)
// 	}
//
// 	if options.PolicyIn != nil {
// 		err = d.Set(mkPolicyIn, *options.PolicyIn)
// 		diags = append(diags, diag.FromErr(err)...)
// 	}
//
// 	if options.PolicyOut != nil {
// 		err = d.Set(mkPolicyOut, *options.PolicyOut)
// 		diags = append(diags, diag.FromErr(err)...)
// 	}
//
// 	d.SetId("cluster-firewall")
//
// 	rules := d.Get(firewall.MkRule).([]interface{})
// 	//nolint:nestif
// 	if len(rules) > 0 {
// 		// We have rules in the state, so we need to read them from the API
// 		for _, v := range rules {
// 			ruleMap := v.(map[string]interface{})
// 			pos := ruleMap[firewall.MkRulePos].(int)
//
// 			err = readRule(ctx, api, pos, ruleMap)
// 			if err != nil {
// 				diags = append(diags, diag.FromErr(err)...)
// 			}
// 		}
// 	} else {
// 		ruleIDs, err := api.ListRules(ctx)
// 		if err != nil {
// 			return diag.FromErr(err)
// 		}
// 		for _, id := range ruleIDs {
// 			ruleMap := map[string]interface{}{}
// 			err = readRule(ctx, api, id.Pos, ruleMap)
// 			if err != nil {
// 				diags = append(diags, diag.FromErr(err)...)
// 			} else {
// 				rules = append(rules, ruleMap)
// 			}
// 		}
// 	}
//
// 	if diags.HasError() {
// 		return diags
// 	}
//
// 	err = d.Set(firewall.MkRule, rules)
// 	return diag.FromErr(err)
// }
//
// func readRule(
// 	ctx context.Context,
// 	api fw.API,
// 	pos int,
// 	ruleMap map[string]interface{},
// ) error {
// 	rule, err := api.GetRule(ctx, pos)
// 	if err != nil {
// 		return fmt.Errorf("error reading rule %d: %w", pos, err)
// 	}
//
// 	firewall.BaseRuleToMap(&rule.BaseRule, ruleMap)
//
// 	// pos in the map should be int!
// 	ruleMap[firewall.MkRulePos] = pos
// 	ruleMap[firewall.MkRuleAction] = rule.Action
// 	ruleMap[firewall.MkRuleType] = rule.Type
//
// 	return nil
// }
//
// func firewallUpdate(ctx context.Context, api fw.API, d *schema.ResourceData) diag.Diagnostics {
// 	diags := setOptions(ctx, api, d)
// 	if diags.HasError() {
// 		return diags
// 	}
//
// 	diags = firewall.RuleUpdate(d, func(body *fw.RuleUpdateRequestBody) error {
// 		err := api.UpdateRule(ctx, *body.Pos, body)
// 		if err != nil {
// 			return fmt.Errorf("error updating rule: %w", err)
// 		}
// 		return nil
// 	})
// 	if diags.HasError() {
// 		return diags
// 	}
//
// 	return firewallRead(ctx, api, d)
// }
//
// func firewallDelete(ctx context.Context, api fw.API, d *schema.ResourceData) diag.Diagnostics {
// 	rules := d.Get(firewall.MkRule).([]interface{})
// 	sort.Slice(rules, func(i, j int) bool {
// 		ruleI := rules[i].(map[string]interface{})
// 		ruleJ := rules[j].(map[string]interface{})
// 		return ruleI[firewall.MkRulePos].(int) > ruleJ[firewall.MkRulePos].(int)
// 	})
//
// 	for _, v := range rules {
// 		rule := v.(map[string]interface{})
// 		pos := rule[firewall.MkRulePos].(int)
// 		err := api.DeleteRule(ctx, pos)
// 		if err != nil {
// 			return diag.FromErr(err)
// 		}
// 	}
//
// 	d.SetId("")
//
// 	return nil
// }

// func FirewallAlias() *schema.Resource {
// 	return &schema.Resource{
// 		Schema:        firewall.AliasSchema(),
// 		CreateContext: invokeFirewallAPI(firewall.AliasCreate),
// 		ReadContext:   invokeFirewallAPI(firewall.AliasRead),
// 		UpdateContext: invokeFirewallAPI(firewall.AliasUpdate),
// 		DeleteContext: invokeFirewallAPI(firewall.AliasDelete),
// 	}
// }
//
// func FirewallIPSet() *schema.Resource {
// 	return &schema.Resource{
// 		Schema:        firewall.IPSetSchema(),
// 		CreateContext: invokeFirewallAPI(firewall.IPSetCreate),
// 		ReadContext:   invokeFirewallAPI(firewall.IPSetRead),
// 		UpdateContext: invokeFirewallAPI(firewall.IPSetUpdate),
// 		DeleteContext: invokeFirewallAPI(firewall.IPSetDelete),
// 	}
// }
//
// func FirewallSecurityGroup() *schema.Resource {
// 	return &schema.Resource{
// 		Schema:        firewall.SecurityGroupSchema(),
// 		CreateContext: invokeFirewallAPI(firewall.SecurityGroupCreate),
// 		ReadContext:   invokeFirewallAPI(firewall.SecurityGroupRead),
// 		UpdateContext: invokeFirewallAPI(firewall.SecurityGroupUpdate),
// 		DeleteContext: invokeFirewallAPI(firewall.SecurityGroupDelete),
// 	}
// }

// func invokeFirewallAPI(
// 	f func(context.Context, fw.API, *schema.ResourceData) diag.Diagnostics,
// ) func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics {
// 	return func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
// 		config := m.(proxmoxtf.ProviderConfiguration)
// 		veClient, err := config.GetVEClient()
// 		if err != nil {
// 			return diag.FromErr(err)
// 		}
//
// 		return f(ctx, veClient.API().Cluster().Firewall(), d)
// 	}
// }

func selectFirewallAPI(
	f func(context.Context, firewall.API, *schema.ResourceData) diag.Diagnostics,
) func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		config := m.(proxmoxtf.ProviderConfiguration)
		veClient, err := config.GetVEClient()
		if err != nil {
			return diag.FromErr(err)
		}

		api := veClient.API().Cluster().Firewall()

		return f(ctx, api, d)
	}
}
