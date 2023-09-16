/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package cluster

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	hagroups "github.com/bpg/terraform-provider-proxmox/proxmox/cluster/ha/groups"
)

// haGroupModel is the model used to represent a High Availability group.
type haGroupModel struct {
	ID         types.String `tfsdk:"id"`          // Identifier used by Terraform
	Group      types.String `tfsdk:"group"`       // HA group name
	Comment    types.String `tfsdk:"comment"`     // Comment, if present
	Nodes      types.Map    `tfsdk:"nodes"`       // Map of member nodes associated with their priorities
	NoFailback types.Bool   `tfsdk:"no_failback"` // Flag that disables failback
	Restricted types.Bool   `tfsdk:"restricted"`  // Flag that prevents execution on other member nodes
}

// Import the contents of a HA group model from the API's response data.
func (m *haGroupModel) importFromAPI(group hagroups.HAGroupGetResponseData) diag.Diagnostics {
	m.Comment = types.StringPointerValue(group.Comment)
	m.NoFailback = group.NoFailback.ToValue()
	m.Restricted = group.Restricted.ToValue()

	return m.parseHAGroupNodes(group.Nodes)
}

// Parse the list of member nodes. The list is received from the Proxmox API as a string. It must
// be converted into a map value. Errors will be returned as Terraform diagnostics.
func (m *haGroupModel) parseHAGroupNodes(nodes string) diag.Diagnostics {
	var diags diag.Diagnostics

	nodesIn := strings.Split(nodes, ",")
	nodesOut := make(map[string]attr.Value)

	for _, nodeDescStr := range nodesIn {
		nodeDesc := strings.Split(nodeDescStr, ":")
		if len(nodeDesc) > 2 {
			diags.AddWarning(
				"Could not parse HA group node",
				fmt.Sprintf("Received group node '%s' for HA group '%s'",
					nodeDescStr, m.Group.ValueString()),
			)

			continue
		}

		priority := types.Int64Null()

		if len(nodeDesc) == 2 {
			prio, err := strconv.Atoi(nodeDesc[1])
			if err == nil {
				priority = types.Int64Value(int64(prio))
			} else {
				diags.AddWarning(
					"Could not parse HA group node priority",
					fmt.Sprintf("Node priority string '%s' for node %s of HA group '%s'",
						nodeDesc[1], nodeDesc[0], m.Group.ValueString()),
				)
			}
		}

		nodesOut[nodeDesc[0]] = priority
	}

	value, mbDiags := types.MapValue(types.Int64Type, nodesOut)
	diags.Append(mbDiags...)

	m.Nodes = value

	return diags
}
