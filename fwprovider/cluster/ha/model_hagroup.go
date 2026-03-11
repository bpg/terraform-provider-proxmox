/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package ha

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	hagroups "github.com/bpg/terraform-provider-proxmox/proxmox/cluster/ha/groups"
)

// GroupModel is the model used to represent a High Availability group.
type GroupModel struct {
	ID         types.String `tfsdk:"id"`          // Identifier used by Terraform
	Group      types.String `tfsdk:"group"`       // HA group name
	Comment    types.String `tfsdk:"comment"`     // Comment, if present
	Nodes      types.Map    `tfsdk:"nodes"`       // Map of member nodes associated with their priorities
	NoFailback types.Bool   `tfsdk:"no_failback"` // Flag that disables failback
	Restricted types.Bool   `tfsdk:"restricted"`  // Flag that prevents execution on other member nodes
}

// ImportFromAPI imports the contents of a HA group model from the API's response data.
func (m *GroupModel) ImportFromAPI(group hagroups.HAGroupGetResponseData) diag.Diagnostics {
	m.Comment = types.StringPointerValue(group.Comment)
	m.NoFailback = group.NoFailback.ToValue()
	m.Restricted = group.Restricted.ToValue()

	return m.parseHAGroupNodes(group.Nodes)
}

// parseHAGroupNodes delegates to the shared parseNodePriorities helper.
func (m *GroupModel) parseHAGroupNodes(nodes string) diag.Diagnostics {
	value, diags := parseNodePriorities(nodes, m.Group.ValueString())
	m.Nodes = value

	return diags
}
