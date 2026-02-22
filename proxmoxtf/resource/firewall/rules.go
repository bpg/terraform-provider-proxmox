/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package firewall

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

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
	MkRule = "rule"

	mkSecurityGroup = "security_group"

	mkRuleAction  = "action"
	mkRuleComment = "comment"
	mkRuleDPort   = "dport"
	mkRuleDest    = "dest"
	mkRuleEnabled = "enabled"
	mkRuleIFace   = "iface"
	mkRuleLog     = "log"
	mkRuleMacro   = "macro"
	mkRulePos     = "pos"
	mkRuleProto   = "proto"
	mkRuleSource  = "source"
	mkRuleSPort   = "sport"
	mkRuleType    = "type"
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
			Description:      "Rules type ('in', 'out', 'forward')",
			Optional:         true,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"in", "out", "forward"}, true)),
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
			Optional:    true,
			DefaultFunc: func() (any, error) {
				return make([]any, 0), nil
			},
			Elem: &schema.Resource{Schema: rule},
		},
	}

	structure.MergeSchema(s, selectorSchema())

	return &schema.Resource{
		Schema:        s,
		CreateContext: invokeRuleAPI(RulesCreate),
		ReadContext:   invokeRuleAPI(RulesRead),
		UpdateContext: invokeRuleAPI(RulesUpdate),
		DeleteContext: invokeRuleAPI(RulesDelete),
		Importer: &schema.ResourceImporter{
			StateContext: RulesImport,
		},
	}
}

// RulesImport imports firewall rules.
func RulesImport(_ context.Context, d *schema.ResourceData, m any) ([]*schema.ResourceData, error) {
	id := d.Id()

	switch {
	case id == "cluster":
	case strings.HasPrefix(id, "node/"):
		parts := strings.SplitN(id, "/", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid node import ID format: %s (expected: node/<node_name>)", id)
		}

		nodeName := parts[1]
		if nodeName == "" {
			return nil, fmt.Errorf("invalid node import ID: node name cannot be empty in %s", id)
		}

		err := d.Set(mkSelectorNodeName, nodeName)
		if err != nil {
			return nil, fmt.Errorf("failed setting node name during import: %w", err)
		}
	case strings.HasPrefix(id, "vm/"):
		parts := strings.SplitN(id, "/", 3)
		if len(parts) != 3 {
			return nil, fmt.Errorf("invalid VM import ID format: %s (expected: vm/<node_name>/<vm_id>)", id)
		}

		nodeName := parts[1]
		if nodeName == "" {
			return nil, fmt.Errorf("invalid VM import ID: node name cannot be empty in %s", id)
		}

		vmID, err := strconv.Atoi(parts[2])
		if err != nil {
			return nil, fmt.Errorf("invalid VM import ID: VM ID must be a number in %s: %w", id, err)
		}

		err = d.Set(mkSelectorNodeName, nodeName)
		if err != nil {
			return nil, fmt.Errorf("failed setting node name during import: %w", err)
		}

		err = d.Set(mkSelectorVMID, vmID)
		if err != nil {
			return nil, fmt.Errorf("failed setting VM ID during import: %w", err)
		}
	case strings.HasPrefix(id, "container/"):
		parts := strings.SplitN(id, "/", 3)

		if len(parts) != 3 {
			return nil, fmt.Errorf("invalid container import ID format: %s (expected: container/<node_name>/<container_id>)", id)
		}

		nodeName := parts[1]
		if nodeName == "" {
			return nil, fmt.Errorf("invalid container import ID: node name cannot be empty in %s", id)
		}

		containerID, err := strconv.Atoi(parts[2])
		if err != nil {
			return nil, fmt.Errorf("invalid container import ID: container ID must be a number in %s: %w", id, err)
		}

		err = d.Set(mkSelectorNodeName, nodeName)
		if err != nil {
			return nil, fmt.Errorf("failed setting node name during import: %w", err)
		}

		err = d.Set(mkSelectorContainerID, containerID)
		if err != nil {
			return nil, fmt.Errorf("failed setting container ID during import: %w", err)
		}
	default:
		return nil, fmt.Errorf("invalid import ID: %s (expected: 'cluster', 'vm/<node_name>/<vm_id>', or 'container/<node_name>/<container_id>')", id)
	}

	api, err := firewallAPIFor(d, m)
	if err != nil {
		return nil, err
	}

	d.SetId(api.GetRulesID())

	return []*schema.ResourceData{d}, nil
}

// RulesCreate creates new firewall rules.
func RulesCreate(ctx context.Context, api firewall.Rule, d *schema.ResourceData) diag.Diagnostics {
	diags := diag.Diagnostics{}

	existingRules, err := api.ListRules(ctx)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	if len(existingRules) > 0 {
		diags = append(diags, diag.Errorf("Existing rules detected. Aborting...")...)
		return diags
	}

	rules := d.Get(MkRule).([]any)

	for i := len(rules) - 1; i >= 0; i-- {
		rule := rules[i].(map[string]any)

		ruleBody, err := mapToRuleCreateRequestBody(rule)
		if err != nil {
			diags = append(diags, diag.FromErr(err)...)
			continue
		}

		err = api.CreateRule(ctx, ruleBody)
		diags = append(diags, diag.FromErr(err)...)
	}

	if diags.HasError() {
		return diags
	}

	// reset rules, we re-read them (with proper positions) from the API
	err = d.Set(MkRule, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(api.GetRulesID())

	return RulesRead(ctx, api, d)
}

// RulesRead reads rules from the API and updates the state.
func RulesRead(ctx context.Context, api firewall.Rule, d *schema.ResourceData) diag.Diagnostics {
	diags := diag.Diagnostics{}

	readRule := func(pos int, ruleMap map[string]any) error {
		rule, err := api.GetRule(ctx, pos)
		if err != nil {
			if errors.Is(err, firewall.ErrNoRuleAtPosition) {
				return firewall.ErrNoRuleAtPosition
			}

			return fmt.Errorf("error reading rule %d : %w", pos, err)
		}

		ruleMap[mkRulePos] = pos

		if rule.Type == "group" {
			ruleMap[mkSecurityGroup] = rule.Action
			securityGroupBaseRuleToMap(&rule.BaseRule, ruleMap)
		} else {
			ruleMap[mkRuleAction] = rule.Action
			ruleMap[mkRuleType] = rule.Type
			baseRuleToMap(&rule.BaseRule, ruleMap)
		}

		return nil
	}

	ruleIDs, err := api.ListRules(ctx)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	rules := make([]map[string]any, 0)

	for _, id := range ruleIDs {
		ruleMap := map[string]any{}

		err = readRule(id.Pos, ruleMap)
		if err != nil {
			if !errors.Is(err, firewall.ErrNoRuleAtPosition) {
				diags = append(diags, diag.FromErr(err)...)
			}
		} else if len(ruleMap) > 0 {
			rules = append(rules, ruleMap)
		}
	}

	if diags.HasError() {
		return diags
	}

	err = d.Set(MkRule, rules)
	diags = append(diags, diag.FromErr(err)...)

	return diags
}

// RulesUpdate updates rules.
func RulesUpdate(ctx context.Context, api firewall.Rule, d *schema.ResourceData) diag.Diagnostics {
	var diags diag.Diagnostics

	oldRules, newRules := d.GetChange(MkRule)
	oldRulesList := oldRules.([]any)
	newRulesList := newRules.([]any)

	// build per-signature queues for old rules to handle duplicate identities correctly.
	// multiple rules can share the same signature (identity fields only, excludes comment/enabled/log).
	oldBySig := make(map[string][]map[string]any)

	for _, rule := range oldRulesList {
		ruleMap := rule.(map[string]any)
		sig := computeRuleSignature(rule)
		oldBySig[sig] = append(oldBySig[sig], ruleMap)
	}

	// match new rules to old rules by consuming from per-signature queues.
	// unmatched new rules need to be created.
	sigConsumed := make(map[string]int)
	newToOldRuleMap := make(map[int]map[string]any)

	for i, rule := range newRulesList {
		sig := computeRuleSignature(rule)
		ruleMap := rule.(map[string]any)
		idx := sigConsumed[sig]
		sigConsumed[sig]++

		if idx < len(oldBySig[sig]) {
			newToOldRuleMap[i] = oldBySig[sig][idx]
		} else {
			ruleBody, err := mapToRuleCreateRequestBody(ruleMap)
			if err != nil {
				diags = append(diags, diag.Errorf("could not create rule: %v", err)...)
				return diags
			}

			err = api.CreateRule(ctx, ruleBody)
			if err != nil {
				diags = append(diags, diag.Errorf("could not create rule: %v", err)...)
				return diags
			}
		}
	}

	// reposition all rules to target positions.
	// we re-read from API after each move to get accurate positions.
	// skip positions below targetPos to avoid re-matching already-finalized positions,
	// which is critical when multiple rules share the same signature.
	// this is O(n^2) but correct for small rule sets.
	for targetPos := range newRulesList {
		rule := newRulesList[targetPos]
		sig := computeRuleSignature(rule)

		// re-read current positions from API
		currentRules, err := api.ListRules(ctx)
		if err != nil {
			diags = append(diags, diag.FromErr(err)...)
			return diags
		}

		// find current position of this rule by signature
		currentPos := -1

		for _, cr := range currentRules {
			if cr.Pos < targetPos {
				continue
			}

			r, err := api.GetRule(ctx, cr.Pos)
			if err != nil {
				if errors.Is(err, firewall.ErrNoRuleAtPosition) {
					continue
				}

				diags = append(diags, diag.Errorf("could not read rule at pos %d: %v", cr.Pos, err)...)

				return diags
			}

			if computeRuleSignatureFromAPI(r) == sig {
				currentPos = cr.Pos
				break
			}
		}

		if currentPos == -1 {
			diags = append(diags, diag.Errorf(
				"could not find rule with signature %q at or after position %d during repositioning", sig, targetPos)...)

			return diags
		}

		if currentPos == targetPos {
			continue
		}

		err = api.UpdateRule(ctx, currentPos, &firewall.RuleUpdateRequestBody{
			MoveTo: &targetPos,
		})
		if err != nil {
			diags = append(diags, diag.Errorf("could not move rule from pos %d to %d: %v", currentPos, targetPos, err)...)
			return diags
		}
	}

	// update non-positional attributes for matched rules.
	// newly created rules are skipped since they already have correct attributes.
	for i, newRule := range newRulesList {
		newRuleMap := newRule.(map[string]any)

		oldRuleMap, exists := newToOldRuleMap[i]
		if !exists {
			continue
		}

		isSG := newRuleMap[mkSecurityGroup].(string) != ""

		var ruleBody firewall.RuleUpdateRequestBody

		if isSG {
			ruleBody.BaseRule = *mapToSecurityGroupBaseRule(newRuleMap)
		} else {
			ruleBody.BaseRule = *mapToBaseRule(newRuleMap)

			if action := newRuleMap[mkRuleAction].(string); action != "" {
				ruleBody.Action = &action
			}

			if rType := newRuleMap[mkRuleType].(string); rType != "" {
				ruleBody.Type = &rType
			}
		}

		var fieldsToDelete []string

		var fields []string
		if isSG {
			fields = []string{mkRuleComment, mkRuleIFace}
		} else {
			fields = []string{
				mkRuleComment,
				mkRuleDPort,
				mkRuleDest,
				mkRuleIFace,
				mkRuleLog,
				mkRuleMacro,
				mkRuleProto,
				mkRuleSource,
				mkRuleSPort,
			}
		}

		for _, field := range fields {
			if newRuleMap[field].(string) == "" && oldRuleMap[field].(string) != "" {
				fieldsToDelete = append(fieldsToDelete, field)
			}
		}

		if len(fieldsToDelete) > 0 {
			ruleBody.Delete = fieldsToDelete
		}

		err := api.UpdateRule(ctx, i, &ruleBody)
		if err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	// delete excess rules.
	// after repositioning, desired rules occupy positions 0..N-1.
	// any rules at positions N..M-1 are excess and should be deleted in reverse order.
	currentRulesForDelete, err := api.ListRules(ctx)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	for i := len(currentRulesForDelete) - 1; i >= len(newRulesList); i-- {
		pos := currentRulesForDelete[i].Pos

		err := api.DeleteRule(ctx, pos)
		if err != nil {
			diags = append(diags, diag.Errorf("could not delete rule at pos %d: %v", pos, err)...)
			return diags
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

	rules := d.Get(MkRule).([]any)
	sort.Slice(rules, func(i, j int) bool {
		ruleI := rules[i].(map[string]any)
		ruleJ := rules[j].(map[string]any)

		return ruleI[mkRulePos].(int) > ruleJ[mkRulePos].(int)
	})

	for _, v := range rules {
		rule := v.(map[string]any)
		pos := rule[mkRulePos].(int)

		_, err := api.GetRule(ctx, pos)
		if err != nil {
			if errors.Is(err, firewall.ErrNoRuleAtPosition) {
				continue
			}

			diags = append(diags, diag.Errorf("error checking rule at pos %d before deletion: %v", pos, err)...)

			return diags
		}

		err = api.DeleteRule(ctx, pos)
		diags = append(diags, diag.FromErr(err)...)
	}

	return diags
}

func mapToRuleCreateRequestBody(rule map[string]any) (*firewall.RuleCreateRequestBody, error) {
	var body firewall.RuleCreateRequestBody

	sg := rule[mkSecurityGroup].(string)
	if sg != "" {
		// this is a special case of security group insertion
		body = firewall.RuleCreateRequestBody{
			Action:   sg,
			Type:     "group",
			BaseRule: *mapToSecurityGroupBaseRule(rule),
		}
	} else {
		a := rule[mkRuleAction].(string)
		t := rule[mkRuleType].(string)

		if a == "" || t == "" {
			return nil, fmt.Errorf("either '%s' OR both '%s' and '%s' must be defined", mkSecurityGroup, mkRuleAction, mkRuleType)
		}

		body = firewall.RuleCreateRequestBody{
			Action:   a,
			Type:     t,
			BaseRule: *mapToBaseRule(rule),
		}
	}

	return &body, nil
}

func mapToBaseRule(rule map[string]any) *firewall.BaseRule {
	baseRule := &firewall.BaseRule{}

	comment := rule[mkRuleComment].(string)
	baseRule.Comment = &comment

	dest := rule[mkRuleDest].(string)
	baseRule.Dest = &dest

	dport := rule[mkRuleDPort].(string)
	baseRule.DPort = &dport

	enableBool := types.CustomBool(rule[mkRuleEnabled].(bool))
	baseRule.Enable = &enableBool

	macro := rule[mkRuleMacro].(string)
	baseRule.Macro = &macro

	proto := rule[mkRuleProto].(string)
	baseRule.Proto = &proto

	source := rule[mkRuleSource].(string)
	baseRule.Source = &source

	sport := rule[mkRuleSPort].(string)
	baseRule.SPort = &sport

	iface := rule[mkRuleIFace].(string)
	if iface != "" {
		baseRule.IFace = &iface
	}

	log := rule[mkRuleLog].(string)
	if log != "" {
		baseRule.Log = &log
	}

	return baseRule
}

func mapToSecurityGroupBaseRule(rule map[string]any) *firewall.BaseRule {
	baseRule := &firewall.BaseRule{}

	comment := rule[mkRuleComment].(string)
	baseRule.Comment = &comment

	enableBool := types.CustomBool(rule[mkRuleEnabled].(bool))
	baseRule.Enable = &enableBool

	iface := rule[mkRuleIFace].(string)
	if iface != "" {
		baseRule.IFace = &iface
	}

	return baseRule
}

func baseRuleToMap(baseRule *firewall.BaseRule, rule map[string]any) {
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
		rule[mkRuleEnabled] = *baseRule.Enable
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

func securityGroupBaseRuleToMap(baseRule *firewall.BaseRule, rule map[string]any) {
	if baseRule.Comment != nil {
		rule[mkRuleComment] = *baseRule.Comment
	}

	if baseRule.Enable != nil {
		rule[mkRuleEnabled] = *baseRule.Enable
	}

	if baseRule.IFace != nil {
		rule[mkRuleIFace] = *baseRule.IFace
	}
}

// computeRuleSignature generates a signature for a rule based on its identity fields.
// for regular rules: type, action, dest, dport, source, sport, proto, macro, iface.
// for security group rules: group name and iface.
// non-identity fields (comment, enabled, log) are excluded.
// rules with the same signature are considered the same logical rule for matching
// during updates. changing an identity field produces a new signature, causing the
// algorithm to create a new rule and delete the old one (not an in-place update).
func computeRuleSignature(rule any) string {
	ruleMap := rule.(map[string]any)

	// security group rules
	// identified by group name and interface
	if sg := ruleMap[mkSecurityGroup].(string); sg != "" {
		return strings.Join([]string{"group", sg, ruleMap[mkRuleIFace].(string)}, ":")
	}

	fields := []string{
		"rule",
		ruleMap[mkRuleType].(string),
		ruleMap[mkRuleAction].(string),
		ruleMap[mkRuleDest].(string),
		ruleMap[mkRuleDPort].(string),
		ruleMap[mkRuleSource].(string),
		ruleMap[mkRuleSPort].(string),
		ruleMap[mkRuleProto].(string),
		ruleMap[mkRuleMacro].(string),
		ruleMap[mkRuleIFace].(string),
	}

	return strings.Join(fields, ":")
}

func strOrEmpty(s *string) string {
	if s != nil {
		return *s
	}

	return ""
}

// computeRuleSignatureFromAPI generates a signature from API response data.
// for security group rules, the Proxmox API stores the group name in the Action field
// when Type is "group" (see mapToRuleCreateRequestBody and RulesRead).
func computeRuleSignatureFromAPI(rule *firewall.RuleGetResponseData) string {
	if rule.Type == "group" {
		return strings.Join([]string{"group", rule.Action, strOrEmpty(rule.IFace)}, ":")
	}

	fields := []string{
		"rule",
		rule.Type,
		rule.Action,
		strOrEmpty(rule.Dest),
		strOrEmpty(rule.DPort),
		strOrEmpty(rule.Source),
		strOrEmpty(rule.SPort),
		strOrEmpty(rule.Proto),
		strOrEmpty(rule.Macro),
		strOrEmpty(rule.IFace),
	}

	return strings.Join(fields, ":")
}

func invokeRuleAPI(
	f func(context.Context, firewall.Rule, *schema.ResourceData) diag.Diagnostics,
) func(context.Context, *schema.ResourceData, any) diag.Diagnostics {
	return func(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
		return selectFirewallAPI(func(ctx context.Context, api firewall.API, data *schema.ResourceData) diag.Diagnostics {
			return f(ctx, api, data)
		})(ctx, d, m)
	}
}
