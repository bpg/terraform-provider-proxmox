/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package firewall

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"

	proxmoxfirewall "github.com/bpg/terraform-provider-proxmox/proxmox/firewall"
	proxmoxtypes "github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

type fakeRuleAPI struct {
	rules      []*proxmoxfirewall.RuleGetResponseData
	operations []string
}

func (f *fakeRuleAPI) GetRulesID() string {
	return "fake-rules"
}

func (f *fakeRuleAPI) CreateRule(_ context.Context, body *proxmoxfirewall.RuleCreateRequestBody) error {
	rule := &proxmoxfirewall.RuleGetResponseData{
		BaseRule: body.BaseRule,
		Action:   body.Action,
		Type:     body.Type,
	}
	f.rules = append([]*proxmoxfirewall.RuleGetResponseData{rule}, f.rules...)
	f.operations = append(f.operations, "create:"+body.Action)

	return nil
}

func (f *fakeRuleAPI) GetRule(_ context.Context, position int) (*proxmoxfirewall.RuleGetResponseData, error) {
	if position < 0 || position >= len(f.rules) {
		return nil, proxmoxfirewall.ErrNoRuleAtPosition
	}

	rule := f.rules[position]
	rule.Pos = strconv.Itoa(position)

	return rule, nil
}

func (f *fakeRuleAPI) ListRules(_ context.Context) ([]*proxmoxfirewall.RuleListResponseData, error) {
	result := make([]*proxmoxfirewall.RuleListResponseData, len(f.rules))
	for i := range f.rules {
		result[i] = &proxmoxfirewall.RuleListResponseData{Pos: i}
	}

	return result, nil
}

func (f *fakeRuleAPI) UpdateRule(_ context.Context, position int, body *proxmoxfirewall.RuleUpdateRequestBody) error {
	if position < 0 || position >= len(f.rules) {
		return proxmoxfirewall.ErrNoRuleAtPosition
	}

	if body.MoveTo != nil {
		rule := f.rules[position]
		f.rules = append(f.rules[:position], f.rules[position+1:]...)
		target := min(*body.MoveTo, len(f.rules))
		f.rules = append(f.rules, nil)
		copy(f.rules[target+1:], f.rules[target:])
		f.rules[target] = rule
		f.operations = append(f.operations, fmt.Sprintf("move:%d:%d", position, target))

		return nil
	}

	rule := f.rules[position]
	copyBaseRuleFields(&rule.BaseRule, &body.BaseRule)

	if body.Action != nil {
		rule.Action = *body.Action
	}

	if body.Type != nil {
		rule.Type = *body.Type
	}

	for _, field := range body.Delete {
		deleteBaseRuleField(&rule.BaseRule, field)
	}

	f.operations = append(f.operations, fmt.Sprintf("update:%d", position))

	return nil
}

func (f *fakeRuleAPI) DeleteRule(_ context.Context, position int) error {
	if position < 0 || position >= len(f.rules) {
		return proxmoxfirewall.ErrNoRuleAtPosition
	}

	f.rules = append(f.rules[:position], f.rules[position+1:]...)
	f.operations = append(f.operations, fmt.Sprintf("delete:%d", position))

	return nil
}

func copyBaseRuleFields(target, source *proxmoxfirewall.BaseRule) {
	if source.Comment != nil {
		target.Comment = source.Comment
	}

	if source.Dest != nil {
		target.Dest = source.Dest
	}

	if source.DPort != nil {
		target.DPort = source.DPort
	}

	if source.Enable != nil {
		target.Enable = source.Enable
	}

	if source.ICMPType != nil {
		target.ICMPType = source.ICMPType
	}

	if source.Log != nil {
		target.Log = source.Log
	}

	if source.Macro != nil {
		target.Macro = source.Macro
	}

	if source.Proto != nil {
		target.Proto = source.Proto
	}

	if source.Source != nil {
		target.Source = source.Source
	}

	if source.SPort != nil {
		target.SPort = source.SPort
	}
}

func deleteBaseRuleField(rule *proxmoxfirewall.BaseRule, field string) {
	switch field {
	case "comment":
		rule.Comment = nil
	case "dest":
		rule.Dest = nil
	case "dport":
		rule.DPort = nil
	case "icmp-type":
		rule.ICMPType = nil
	case "log":
		rule.Log = nil
	case "macro":
		rule.Macro = nil
	case "proto":
		rule.Proto = nil
	case "source":
		rule.Source = nil
	case "sport":
		rule.SPort = nil
	}
}

func forwardRule(action string) ruleModel {
	return ruleModel{
		Action:        types.StringValue(action),
		Comment:       types.StringNull(),
		Dest:          types.StringNull(),
		DPort:         types.StringNull(),
		Enabled:       types.BoolValue(true),
		ICMPType:      types.StringNull(),
		Log:           types.StringNull(),
		Macro:         types.StringNull(),
		Proto:         types.StringNull(),
		SecurityGroup: types.StringNull(),
		Source:        types.StringNull(),
		SPort:         types.StringNull(),
	}
}

func TestCreateRulesPreservesDesiredOrder(t *testing.T) {
	t.Parallel()

	api := &fakeRuleAPI{}
	rules := []ruleModel{
		forwardRule("ACCEPT"),
		forwardRule("DROP"),
	}

	require.NoError(t, createRules(t.Context(), api, rules))
	require.Len(t, api.rules, 2)
	require.Equal(t, "ACCEPT", api.rules[0].Action)
	require.Equal(t, "DROP", api.rules[1].Action)
}

func TestReconcileRulesCreatesMovesUpdatesAndDeletes(t *testing.T) {
	t.Parallel()

	httpsRule := forwardRule("ACCEPT")
	httpsRule.Comment = types.StringValue("Allow application HTTPS")
	httpsRule.Dest = types.StringValue("10.96.20.0/24")
	httpsRule.DPort = types.StringValue("443")
	httpsRule.Log = types.StringValue("info")
	httpsRule.Proto = types.StringValue("tcp")
	httpsRule.Source = types.StringValue("10.96.10.0/24")

	dropRule := forwardRule("DROP")
	dropRule.Comment = types.StringValue("Deny remaining application traffic")
	dropRule.Dest = types.StringValue("10.96.20.0/24")
	dropRule.Enabled = types.BoolValue(false)
	dropRule.Source = types.StringValue("10.96.10.0/24")

	api := &fakeRuleAPI{}
	require.NoError(t, createRules(t.Context(), api, []ruleModel{httpsRule, dropRule}))
	api.operations = nil

	dnsRule := forwardRule("ACCEPT")
	dnsRule.Comment = types.StringValue("Allow DNS")
	dnsRule.Dest = types.StringValue("10.96.53.53")
	dnsRule.DPort = types.StringValue("53")
	dnsRule.Proto = types.StringValue("udp")
	dnsRule.Source = types.StringValue("10.96.10.0/24")

	updatedHTTPS := httpsRule
	updatedHTTPS.Comment = types.StringValue("Allow production HTTPS")
	updatedHTTPS.Log = types.StringValue("notice")

	updatedDrop := dropRule
	updatedDrop.Comment = types.StringNull()
	updatedDrop.Enabled = types.BoolValue(true)

	desired := []ruleModel{dnsRule, updatedHTTPS, updatedDrop}
	require.NoError(t, reconcileRules(t.Context(), api, []ruleModel{httpsRule, dropRule}, desired))

	result, err := readRules(t.Context(), "fwvmain", api)
	require.NoError(t, err)
	require.Equal(t, desired, result.Rules)
	require.Contains(t, api.operations, "create:ACCEPT")
	require.Contains(t, api.operations, "update:1")
	require.Contains(t, api.operations, "update:2")
}

func TestReconcileRulesRemovesExcessRules(t *testing.T) {
	t.Parallel()

	accept := forwardRule("ACCEPT")
	dns := forwardRule("ACCEPT")
	dns.DPort = types.StringValue("53")
	dns.Proto = types.StringValue("udp")
	drop := forwardRule("DROP")

	api := &fakeRuleAPI{}
	require.NoError(t, createRules(t.Context(), api, []ruleModel{accept, dns, drop}))
	api.operations = nil

	require.NoError(t, reconcileRules(t.Context(), api, []ruleModel{accept, dns, drop}, []ruleModel{accept}))
	require.Len(t, api.rules, 1)
	require.Equal(t, "ACCEPT", api.rules[0].Action)
	require.Equal(t, []string{"update:0", "delete:2", "delete:1"}, api.operations)
}

func TestRuleModelMapsSecurityGroupAndDisabledDefault(t *testing.T) {
	t.Parallel()

	apiRule := &proxmoxfirewall.RuleGetResponseData{
		BaseRule: proxmoxfirewall.BaseRule{
			Comment: new("shared policy"),
			Enable:  nil,
		},
		Action: "monitoring",
		Type:   "group",
	}

	var rule ruleModel
	rule.fromAPI(apiRule)
	require.True(t, rule.Action.IsNull())
	require.Equal(t, "monitoring", rule.SecurityGroup.ValueString())
	require.False(t, rule.Enabled.ValueBool())

	created := rule.toCreateAPI()
	require.Equal(t, "monitoring", created.Action)
	require.Equal(t, "group", created.Type)
	require.Equal(t, proxmoxtypes.CustomBool(false), *created.Enable)
}
