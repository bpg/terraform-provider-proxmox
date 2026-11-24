/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package firewall

import (
	"context"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	clusterfirewall "github.com/bpg/terraform-provider-proxmox/proxmox/cluster/firewall"
	"github.com/bpg/terraform-provider-proxmox/proxmox/firewall"
	"github.com/bpg/terraform-provider-proxmox/proxmox/helpers/ptr"
	proxmoxtypes "github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

type securityGroupModel struct {
	ID      types.String `tfsdk:"id"`
	Name    types.String `tfsdk:"name"`
	Comment types.String `tfsdk:"comment"`
	Rules   types.List   `tfsdk:"rule"`
}

type ruleModel struct {
	Pos     types.Int64  `tfsdk:"pos"`
	Action  types.String `tfsdk:"action"`
	Type    types.String `tfsdk:"type"`
	Comment types.String `tfsdk:"comment"`
	Dest    types.String `tfsdk:"dest"`
	DPort   types.String `tfsdk:"dport"`
	Enabled types.Bool   `tfsdk:"enabled"`
	Log     types.String `tfsdk:"log"`
	Macro   types.String `tfsdk:"macro"`
	Proto   types.String `tfsdk:"proto"`
	Source  types.String `tfsdk:"source"`
	SPort   types.String `tfsdk:"sport"`
}

func ruleObjectType() types.ObjectType {
	return types.ObjectType{AttrTypes: map[string]attr.Type{
		"pos":     types.Int64Type,
		"action":  types.StringType,
		"type":    types.StringType,
		"comment": types.StringType,
		"dest":    types.StringType,
		"dport":   types.StringType,
		"enabled": types.BoolType,
		"log":     types.StringType,
		"macro":   types.StringType,
		"proto":   types.StringType,
		"source":  types.StringType,
		"sport":   types.StringType,
	}}
}

// ruleModels decodes the rule list; a null or unknown list yields no rules.
func (m *securityGroupModel) ruleModels(ctx context.Context, diags *diag.Diagnostics) []ruleModel {
	var rules []ruleModel

	if m.Rules.IsNull() || m.Rules.IsUnknown() {
		return rules
	}

	diags.Append(m.Rules.ElementsAs(ctx, &rules, false)...)

	return rules
}

// fromAPI fills the model from PVE data. prior supplies the rules from config or state, so a macro PVE re-cased
// (e.g. "ssh" -> "SSH") keeps the user's spelling instead of producing a diff.
func (m *securityGroupModel) fromAPI(
	ctx context.Context,
	group *clusterfirewall.GroupListResponseData,
	rules []rule,
	prior []ruleModel,
	diags *diag.Diagnostics,
) {
	m.ID = types.StringValue(group.Group)
	m.Name = types.StringValue(group.Group)
	m.Comment = stringOrNull(ptr.Or(group.Comment, ""))

	models := make([]ruleModel, len(rules))

	for i, r := range rules {
		models[i] = r.toModel(int64(i))

		if i < len(prior) && strings.EqualFold(prior[i].Macro.ValueString(), r.Macro) {
			models[i].Macro = prior[i].Macro
		}
	}

	list, d := types.ListValueFrom(ctx, ruleObjectType(), models)
	diags.Append(d...)

	m.Rules = list
}

// rule is a security group rule as PVE stores it. An empty string means the field is unset.
type rule struct {
	Action  string
	Type    string
	Comment string
	Dest    string
	DPort   string
	Enabled bool
	Log     string
	Macro   string
	Proto   string
	Source  string
	SPort   string
}

func ruleFromAPI(d *firewall.RuleGetResponseData) rule {
	return rule{
		Action:  d.Action,
		Type:    d.Type,
		Comment: ptr.Or(d.Comment, ""),
		Dest:    ptr.Or(d.Dest, ""),
		DPort:   ptr.Or(d.DPort, ""),
		Enabled: d.Enable != nil && bool(*d.Enable),
		Log:     ptr.Or(d.Log, ""),
		Macro:   ptr.Or(d.Macro, ""),
		Proto:   ptr.Or(d.Proto, ""),
		Source:  ptr.Or(d.Source, ""),
		SPort:   ptr.Or(d.SPort, ""),
	}
}

func (m ruleModel) toRule() rule {
	return rule{
		Action:  m.Action.ValueString(),
		Type:    m.Type.ValueString(),
		Comment: m.Comment.ValueString(),
		Dest:    m.Dest.ValueString(),
		DPort:   m.DPort.ValueString(),
		Enabled: m.Enabled.ValueBool(),
		Log:     m.Log.ValueString(),
		Macro:   m.Macro.ValueString(),
		Proto:   m.Proto.ValueString(),
		Source:  m.Source.ValueString(),
		SPort:   m.SPort.ValueString(),
	}
}

func (r rule) toModel(pos int64) ruleModel {
	return ruleModel{
		Pos:     types.Int64Value(pos),
		Action:  types.StringValue(r.Action),
		Type:    types.StringValue(r.Type),
		Comment: stringOrNull(r.Comment),
		Dest:    stringOrNull(r.Dest),
		DPort:   stringOrNull(r.DPort),
		Enabled: types.BoolValue(r.Enabled),
		Log:     stringOrNull(r.Log),
		Macro:   stringOrNull(r.Macro),
		Proto:   stringOrNull(r.Proto),
		Source:  stringOrNull(r.Source),
		SPort:   stringOrNull(r.SPort),
	}
}

// signature identifies a rule by the fields that decide which packets it matches. Rules with the same signature are
// interchangeable when reconciling, so a comment, log or enabled change updates a rule in place.
func (r rule) signature() string {
	return strings.Join([]string{
		r.Type, r.Action, r.Dest, r.DPort, r.Source, r.SPort, r.Proto, strings.ToLower(r.Macro),
	}, "\x00")
}

func (r rule) equal(o rule) bool {
	return r.signature() == o.signature() && r.Comment == o.Comment && r.Enabled == o.Enabled && r.Log == o.Log
}

// optionalFields pairs each clearable PVE rule property with its value.
func (r rule) optionalFields() map[string]string {
	return map[string]string{
		"comment": r.Comment,
		"dest":    r.Dest,
		"dport":   r.DPort,
		"log":     r.Log,
		"macro":   r.Macro,
		"proto":   r.Proto,
		"source":  r.Source,
		"sport":   r.SPort,
	}
}

func (r rule) baseRule() firewall.BaseRule {
	// Always sent: PVE stores a rule created without enable as disabled.
	enabled := proxmoxtypes.CustomBool(r.Enabled)

	return firewall.BaseRule{
		Comment: ptrOrNil(r.Comment),
		Dest:    ptrOrNil(r.Dest),
		DPort:   ptrOrNil(r.DPort),
		Enable:  &enabled,
		Log:     ptrOrNil(r.Log),
		Macro:   ptrOrNil(r.Macro),
		Proto:   ptrOrNil(r.Proto),
		Source:  ptrOrNil(r.Source),
		SPort:   ptrOrNil(r.SPort),
	}
}

func (r rule) createBody() *firewall.RuleCreateRequestBody {
	return &firewall.RuleCreateRequestBody{
		BaseRule: r.baseRule(),
		Action:   r.Action,
		Type:     r.Type,
	}
}

func (r rule) updateBody(current rule) *firewall.RuleUpdateRequestBody {
	body := &firewall.RuleUpdateRequestBody{
		BaseRule: r.baseRule(),
		Action:   &r.Action,
		Type:     &r.Type,
	}

	want := r.optionalFields()

	for name, value := range current.optionalFields() {
		if value != "" && want[name] == "" {
			body.Delete = append(body.Delete, name)
		}
	}

	slices.Sort(body.Delete)

	return body
}

type ruleOpKind int

const (
	ruleOpCreate ruleOpKind = iota
	ruleOpMove
	ruleOpUpdate
	ruleOpDelete
)

type ruleOp struct {
	kind    ruleOpKind
	pos     int
	moveTo  int
	rule    rule
	current rule
}

// planRuleOps returns the ordered API calls that turn the current rules into the desired ones. Existing rules are
// moved rather than rewritten position by position, so rules that are kept never drop out of the live ruleset while
// the apply is in progress.
func planRuleOps(current, desired []rule) []ruleOp {
	var ops []ruleOp

	sim := slices.Clone(current)

	available := make(map[string]int, len(current))
	for _, r := range current {
		available[r.signature()]++
	}

	var toCreate []rule

	for _, r := range desired {
		if sig := r.signature(); available[sig] > 0 {
			available[sig]--
		} else {
			toCreate = append(toCreate, r)
		}
	}

	// PVE inserts every new rule at position 0; creating in reverse keeps them in their relative order.
	for _, r := range slices.Backward(toCreate) {
		ops = append(ops, ruleOp{kind: ruleOpCreate, rule: r})
		sim = slices.Insert(sim, 0, r)
	}

	// Only upward moves are issued: PVE lands a rule exactly on moveto when moving up, but on moveto-1 when moving down.
	for target, want := range desired {
		from := findRule(sim, target, want)
		if from == target {
			continue
		}

		ops = append(ops, ruleOp{kind: ruleOpMove, pos: from, moveTo: target})

		r := sim[from]
		sim = slices.Delete(sim, from, from+1)
		sim = slices.Insert(sim, target, r)
	}

	for pos, want := range desired {
		if !sim[pos].equal(want) {
			ops = append(ops, ruleOp{kind: ruleOpUpdate, pos: pos, rule: want, current: sim[pos]})
			sim[pos] = want
		}
	}

	for pos := len(sim) - 1; pos >= len(desired); pos-- {
		ops = append(ops, ruleOp{kind: ruleOpDelete, pos: pos})
	}

	return ops
}

// findRule returns the first rule at or after from that matches want, preferring an exact match over a signature
// match. planRuleOps creates enough rules per signature that a match always exists.
func findRule(rules []rule, from int, want rule) int {
	sig := want.signature()
	match := -1

	for i := from; i < len(rules); i++ {
		if rules[i].equal(want) {
			return i
		}

		if match == -1 && rules[i].signature() == sig {
			match = i
		}
	}

	return match
}

func stringOrNull(s string) types.String {
	if s == "" {
		return types.StringNull()
	}

	return types.StringValue(s)
}

func ptrOrNil(s string) *string {
	if s == "" {
		return nil
	}

	return &s
}
