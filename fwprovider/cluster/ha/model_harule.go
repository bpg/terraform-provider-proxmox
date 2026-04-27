/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package ha

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/attribute"
	harules "github.com/bpg/terraform-provider-proxmox/proxmox/cluster/ha/rules"
)

const (
	// RuleTypeNodeAffinity is the HA rule type for node affinity.
	RuleTypeNodeAffinity = "node-affinity"
	// RuleTypeResourceAffinity is the HA rule type for resource affinity.
	RuleTypeResourceAffinity = "resource-affinity"
)

// RuleModel is the model used to represent a High Availability rule.
type RuleModel struct {
	ID        types.String `tfsdk:"id"`        // Identifier used by Terraform
	Rule      types.String `tfsdk:"rule"`      // HA rule identifier
	Type      types.String `tfsdk:"type"`      // HA rule type (node-affinity or resource-affinity)
	Comment   types.String `tfsdk:"comment"`   // Comment, if present
	Disable   types.Bool   `tfsdk:"disable"`   // Whether the rule is disabled
	Resources types.Set    `tfsdk:"resources"` // Set of HA resource IDs (e.g. vm:100, ct:101)
	Nodes     types.Map    `tfsdk:"nodes"`     // Map of node names to priorities (node-affinity only)
	Strict    types.Bool   `tfsdk:"strict"`    // Whether the node affinity is strict (node-affinity only)
	Affinity  types.String `tfsdk:"affinity"`  // positive or negative (resource-affinity only)
}

// ImportFromAPI imports the contents of a HA rule model from the API's response data.
func (m *RuleModel) ImportFromAPI(rule harules.HARuleGetResponseData) diag.Diagnostics {
	diags := diag.Diagnostics{}

	m.Rule = types.StringValue(rule.Rule)
	m.Type = types.StringValue(rule.Type)
	m.Comment = types.StringPointerValue(rule.Comment)
	m.Disable = rule.Disable.ToValue()

	// Parse resources string into a set.
	resDiags := m.parseResources(rule.Resources)
	diags.Append(resDiags...)

	// Type-specific fields.
	switch rule.Type {
	case RuleTypeNodeAffinity:
		if rule.Nodes != nil {
			nodeDiags := m.parseNodes(*rule.Nodes)
			diags.Append(nodeDiags...)
		} else {
			m.Nodes = types.MapNull(types.Int64Type)
		}

		m.Strict = attribute.BoolValueFromCustomBoolPtr(rule.Strict)

		m.Affinity = types.StringNull()

	case RuleTypeResourceAffinity:
		m.Nodes = types.MapNull(types.Int64Type)
		// strict has Computed+Default(false) in schema, so use false (not null)
		// to avoid perpetual plan diffs for resource-affinity rules.
		m.Strict = types.BoolValue(false)
		m.Affinity = types.StringPointerValue(rule.Affinity)
	}

	return diags
}

// parseResources parses a comma-separated resource string into a Terraform set.
func (m *RuleModel) parseResources(resources string) diag.Diagnostics {
	parts := strings.Split(resources, ",")
	elements := make([]attr.Value, 0, len(parts))

	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			elements = append(elements, types.StringValue(trimmed))
		}
	}

	value, diags := types.SetValue(types.StringType, elements)
	m.Resources = value

	return diags
}

// parseNodes delegates to the shared parseNodePriorities helper.
func (m *RuleModel) parseNodes(nodes string) diag.Diagnostics {
	value, diags := parseNodePriorities(nodes, m.Rule.ValueString())
	m.Nodes = value

	return diags
}
