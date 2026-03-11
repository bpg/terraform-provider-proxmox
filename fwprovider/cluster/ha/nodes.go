/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package ha

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// parseNodePriorities parses a comma-separated nodes string (e.g. "node1:2,node2:1,node3")
// into a Terraform map of node names to optional integer priorities. The resourceName is used
// for diagnostic messages.
func parseNodePriorities(nodes, resourceName string) (types.Map, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	nodesIn := strings.Split(nodes, ",")
	nodesOut := make(map[string]attr.Value)

	for _, nodeDescStr := range nodesIn {
		nodeDesc := strings.Split(nodeDescStr, ":")
		if len(nodeDesc) > 2 {
			diags.AddWarning(
				"Could not parse HA node entry",
				fmt.Sprintf("Received node entry '%s' for HA resource '%s'",
					nodeDescStr, resourceName),
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
					"Could not parse HA node priority",
					fmt.Sprintf("Node priority string '%s' for node %s of HA resource '%s'",
						nodeDesc[1], nodeDesc[0], resourceName),
				)
			}
		}

		nodesOut[nodeDesc[0]] = priority
	}

	value, mbDiags := types.MapValue(types.Int64Type, nodesOut)
	diags.Append(mbDiags...)

	return value, diags
}
