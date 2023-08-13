/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package datasource

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/bpg/terraform-provider-proxmox/proxmoxtf"
)

const (
	mkDataSourceVirtualEnvironmentHAGroupID             = "group"
	mkDataSourceVirtualEnvironmentHAGroupComment        = "comment"
	mkDataSourceVirtualEnvironmentHAGroupMembers        = "members"
	mkDataSourceVirtualEnvironmentHAGroupMemberNodeName = "node_name"
	mkDataSourceVirtualEnvironmentHAGroupMemberPriority = "priority"
	mkDataSourceVirtualEnvironmentHAGroupRestricted     = "restricted"
	mkDataSourceVirtualEnvironmentHAGroupNoFailback     = "no_failback"
)

// HAGroup returns a resource that describes a single High Availability group.
func HAGroup() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			mkDataSourceVirtualEnvironmentHAGroupID: {
				Type:        schema.TypeString,
				Description: "The HA group's identifier",
				Required:    true,
			},
			mkDataSourceVirtualEnvironmentHAGroupComment: {
				Type:        schema.TypeString,
				Description: "The HA group's comment",
				Computed:    true,
			},
			mkDataSourceVirtualEnvironmentHAGroupRestricted: {
				Type:        schema.TypeBool,
				Description: "A flag that indicates that other nodes may not be used to run resources associated to this HA group.",
				Computed:    true,
			},
			mkDataSourceVirtualEnvironmentHAGroupNoFailback: {
				Type:        schema.TypeBool,
				Description: "A flag that indicates that failing back to a higher priority node is disabled for this HA group.",
				Computed:    true,
			},
			mkDataSourceVirtualEnvironmentHAGroupMembers: {
				Type:        schema.TypeList,
				Description: "The list of the HA group's members",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mkDataSourceVirtualEnvironmentHAGroupMemberNodeName: {
							Type:        schema.TypeString,
							Description: "The name of the member node",
							Computed:    true,
						},
						mkDataSourceVirtualEnvironmentHAGroupMemberPriority: {
							Type:        schema.TypeInt,
							Description: "The priority assigned to this member node, or -1 if no priority is specified",
							Computed:    true,
						},
					},
				},
			},
		},
		ReadContext: haGroupRead,
	}
}

// Request information about a single High Availability group from the Proxmox API and transform it
// into a Terraform resource.
func haGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	config := m.(proxmoxtf.ProviderConfiguration)

	api, err := config.GetClient()
	if err != nil {
		return diag.FromErr(err)
	}

	groupID := d.Get(mkDataSourceVirtualEnvironmentHAGroupID).(string)

	group, err := api.Cluster().HA().Groups().Get(ctx, groupID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(groupID)

	var comment string
	if group.Comment == nil {
		comment = ""
	} else {
		comment = *group.Comment
	}

	err = d.Set(mkDataSourceVirtualEnvironmentHAGroupComment, comment)
	diags = append(diags, diag.FromErr(err)...)

	err = d.Set(mkDataSourceVirtualEnvironmentHAGroupRestricted, group.Restricted != 0)
	diags = append(diags, diag.FromErr(err)...)

	err = d.Set(mkDataSourceVirtualEnvironmentHAGroupNoFailback, group.NoFailback != 0)
	diags = append(diags, diag.FromErr(err)...)

	members, diags := parseHAGroupMembers(diags, groupID, group.Nodes)
	err = d.Set(mkDataSourceVirtualEnvironmentHAGroupMembers, members)
	diags = append(diags, diag.FromErr(err)...)

	return diags
}

// Parse the list of member nodes. The list is received from the Proxmox API as a string. It must
// be converted into a list of maps that satisfy the schema. Any errors will be added to the list
// of Terraform diagnostics.
func parseHAGroupMembers(diags diag.Diagnostics, groupID string, nodes string) ([]interface{}, diag.Diagnostics) {
	membersIn := strings.Split(nodes, ",")
	membersOut := make([]interface{}, len(membersIn))

	for i, nodeDescStr := range membersIn {
		nodeDesc := strings.Split(nodeDescStr, ":")
		if len(nodeDesc) > 2 {
			membersOut[i] = map[string]interface{}{}

			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary: fmt.Sprintf(
					"Could not parse node specification '%s' in HA group '%s'.",
					nodeDescStr, groupID),
			})

			continue
		}

		memberOut := map[string]interface{}{}
		memberOut[mkDataSourceVirtualEnvironmentHAGroupMemberNodeName] = nodeDesc[0]

		if len(nodeDesc) == 2 {
			prio, err := strconv.Atoi(nodeDesc[1])
			if err == nil {
				memberOut[mkDataSourceVirtualEnvironmentHAGroupMemberPriority] = prio
			} else {
				diags = append(diags, diag.FromErr(err)...)
			}
		} else {
			memberOut[mkDataSourceVirtualEnvironmentHAGroupMemberPriority] = -1
		}

		membersOut[i] = memberOut
	}

	return membersOut, diags
}
