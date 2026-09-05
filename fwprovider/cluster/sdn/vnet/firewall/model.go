/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package firewall

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/attribute"
	proxmoxfirewall "github.com/bpg/terraform-provider-proxmox/proxmox/firewall"
)

type model struct {
	ID    types.String `tfsdk:"id"`
	VNet  types.String `tfsdk:"vnet"`
	Rules []ruleModel  `tfsdk:"rules"`
}

type ruleModel struct {
	Action        types.String `tfsdk:"action"`
	Comment       types.String `tfsdk:"comment"`
	Dest          types.String `tfsdk:"dest"`
	DPort         types.String `tfsdk:"dport"`
	Enabled       types.Bool   `tfsdk:"enabled"`
	ICMPType      types.String `tfsdk:"icmp_type"`
	Log           types.String `tfsdk:"log"`
	Macro         types.String `tfsdk:"macro"`
	Proto         types.String `tfsdk:"proto"`
	SecurityGroup types.String `tfsdk:"security_group"`
	Source        types.String `tfsdk:"source"`
	SPort         types.String `tfsdk:"sport"`
}

func (m *model) fromAPI(vnet string, rules []*proxmoxfirewall.RuleGetResponseData) {
	m.ID = types.StringValue(vnet)
	m.VNet = types.StringValue(vnet)
	m.Rules = make([]ruleModel, len(rules))

	for i, rule := range rules {
		m.Rules[i].fromAPI(rule)
	}
}

func (m *ruleModel) fromAPI(rule *proxmoxfirewall.RuleGetResponseData) {
	if rule.Type == "group" {
		m.Action = types.StringNull()
		m.SecurityGroup = types.StringValue(rule.Action)
	} else {
		m.Action = types.StringValue(rule.Action)
		m.SecurityGroup = types.StringNull()
	}

	m.Comment = types.StringPointerValue(rule.Comment)
	m.Dest = types.StringPointerValue(rule.Dest)
	m.DPort = types.StringPointerValue(rule.DPort)
	m.ICMPType = types.StringPointerValue(rule.ICMPType)
	m.Log = types.StringPointerValue(rule.Log)
	m.Macro = types.StringPointerValue(rule.Macro)
	m.Proto = types.StringPointerValue(rule.Proto)
	m.Source = types.StringPointerValue(rule.Source)
	m.SPort = types.StringPointerValue(rule.SPort)

	if enabled := rule.Enable.PointerBool(); enabled != nil {
		m.Enabled = types.BoolValue(*enabled)
	} else {
		m.Enabled = types.BoolValue(false)
	}
}

func (m *ruleModel) toCreateAPI() *proxmoxfirewall.RuleCreateRequestBody {
	action, ruleType := m.actionAndType()

	return &proxmoxfirewall.RuleCreateRequestBody{
		BaseRule: m.toAPIBaseRule(),
		Action:   action,
		Type:     ruleType,
	}
}

func (m *ruleModel) toUpdateAPI(previous ruleModel) *proxmoxfirewall.RuleUpdateRequestBody {
	action, ruleType := m.actionAndType()
	toDelete := make([]string, 0)

	attribute.CheckDelete(m.Comment, previous.Comment, &toDelete, "comment")
	attribute.CheckDelete(m.Dest, previous.Dest, &toDelete, "dest")
	attribute.CheckDelete(m.DPort, previous.DPort, &toDelete, "dport")
	attribute.CheckDelete(m.ICMPType, previous.ICMPType, &toDelete, "icmp-type")
	attribute.CheckDelete(m.Log, previous.Log, &toDelete, "log")
	attribute.CheckDelete(m.Macro, previous.Macro, &toDelete, "macro")
	attribute.CheckDelete(m.Proto, previous.Proto, &toDelete, "proto")
	attribute.CheckDelete(m.Source, previous.Source, &toDelete, "source")
	attribute.CheckDelete(m.SPort, previous.SPort, &toDelete, "sport")

	return &proxmoxfirewall.RuleUpdateRequestBody{
		BaseRule: m.toAPIBaseRule(),
		Action:   &action,
		Type:     &ruleType,
		Delete:   toDelete,
	}
}

func (m *ruleModel) toAPIBaseRule() proxmoxfirewall.BaseRule {
	return proxmoxfirewall.BaseRule{
		Comment:  attribute.StringPtrFromValue(m.Comment),
		Dest:     attribute.StringPtrFromValue(m.Dest),
		DPort:    attribute.StringPtrFromValue(m.DPort),
		Enable:   attribute.CustomBoolPtrFromValue(m.Enabled),
		ICMPType: attribute.StringPtrFromValue(m.ICMPType),
		Log:      attribute.StringPtrFromValue(m.Log),
		Macro:    attribute.StringPtrFromValue(m.Macro),
		Proto:    attribute.StringPtrFromValue(m.Proto),
		Source:   attribute.StringPtrFromValue(m.Source),
		SPort:    attribute.StringPtrFromValue(m.SPort),
	}
}

func (m *ruleModel) actionAndType() (string, string) {
	if attribute.IsDefined(m.SecurityGroup) {
		return m.SecurityGroup.ValueString(), "group"
	}

	return m.Action.ValueString(), "forward"
}

func (m *ruleModel) signature() string {
	if attribute.IsDefined(m.SecurityGroup) {
		return strings.Join([]string{"group", m.SecurityGroup.ValueString()}, "\x00")
	}

	return strings.Join([]string{
		"forward",
		m.Action.ValueString(),
		m.Dest.ValueString(),
		m.DPort.ValueString(),
		m.ICMPType.ValueString(),
		m.Macro.ValueString(),
		m.Proto.ValueString(),
		m.Source.ValueString(),
		m.SPort.ValueString(),
	}, "\x00")
}

func apiRuleSignature(rule *proxmoxfirewall.RuleGetResponseData) string {
	if rule.Type == "group" {
		return strings.Join([]string{"group", rule.Action}, "\x00")
	}

	return strings.Join([]string{
		"forward",
		rule.Action,
		stringValue(rule.Dest),
		stringValue(rule.DPort),
		stringValue(rule.ICMPType),
		stringValue(rule.Macro),
		stringValue(rule.Proto),
		stringValue(rule.Source),
		stringValue(rule.SPort),
	}, "\x00")
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}

	return *value
}
