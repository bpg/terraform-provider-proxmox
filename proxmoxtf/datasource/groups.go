/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package datasource

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/bpg/terraform-provider-proxmox/proxmoxtf"
)

const (
	mkDataSourceVirtualEnvironmentGroupsComments = "comments"
	mkDataSourceVirtualEnvironmentGroupsGroupIDs = "group_ids"
)

// Groups returns a resource for the Proxmox user groups
func Groups() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			mkDataSourceVirtualEnvironmentGroupsComments: {
				Type:        schema.TypeList,
				Description: "The group comments",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			mkDataSourceVirtualEnvironmentGroupsGroupIDs: {
				Type:        schema.TypeList,
				Description: "The group ids",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
		ReadContext: groupsRead,
	}
}

func groupsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	config := m.(proxmoxtf.ProviderConfiguration)
	veClient, err := config.GetVEClient()
	if err != nil {
		return diag.FromErr(err)
	}

	list, err := veClient.API().Access().ListGroups(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	comments := make([]interface{}, len(list))
	groupIDs := make([]interface{}, len(list))

	for i, v := range list {
		if v.Comment != nil {
			comments[i] = v.Comment
		} else {
			comments[i] = ""
		}

		groupIDs[i] = v.ID
	}

	d.SetId("groups")

	err = d.Set(mkDataSourceVirtualEnvironmentGroupsComments, comments)
	diags = append(diags, diag.FromErr(err)...)
	err = d.Set(mkDataSourceVirtualEnvironmentGroupsGroupIDs, groupIDs)
	diags = append(diags, diag.FromErr(err)...)

	return diags
}
