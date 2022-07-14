/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmoxtf

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	mkDataSourceVirtualEnvironmentGroupACL          = "acl"
	mkDataSourceVirtualEnvironmentGroupACLPath      = "path"
	mkDataSourceVirtualEnvironmentGroupACLPropagate = "propagate"
	mkDataSourceVirtualEnvironmentGroupACLRoleID    = "role_id"
	mkDataSourceVirtualEnvironmentGroupComment      = "comment"
	mkDataSourceVirtualEnvironmentGroupID           = "group_id"
	mkDataSourceVirtualEnvironmentGroupMembers      = "members"
)

func dataSourceVirtualEnvironmentGroup() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			mkDataSourceVirtualEnvironmentGroupACL: {
				Type:        schema.TypeSet,
				Description: "The access control list",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mkDataSourceVirtualEnvironmentGroupACLPath: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The path",
						},
						mkDataSourceVirtualEnvironmentGroupACLPropagate: {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether to propagate to child paths",
						},
						mkDataSourceVirtualEnvironmentGroupACLRoleID: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The role id",
						},
					},
				},
			},
			mkDataSourceVirtualEnvironmentGroupComment: {
				Type:        schema.TypeString,
				Description: "The group comment",
				Computed:    true,
			},
			mkDataSourceVirtualEnvironmentGroupID: {
				Type:        schema.TypeString,
				Description: "The group id",
				Required:    true,
			},
			mkDataSourceVirtualEnvironmentGroupMembers: {
				Type:        schema.TypeSet,
				Description: "The group members",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
		ReadContext: dataSourceVirtualEnvironmentGroupRead,
	}
}

func dataSourceVirtualEnvironmentGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()
	if err != nil {
		return diag.FromErr(err)
	}

	groupID := d.Get(mkDataSourceVirtualEnvironmentGroupID).(string)
	group, err := veClient.GetGroup(ctx, groupID)
	if err != nil {
		return diag.FromErr(err)
	}

	acl, err := veClient.GetACL(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(groupID)

	var aclParsed []interface{}

	for _, v := range acl {
		if v.Type == "group" && v.UserOrGroupID == groupID {
			aclEntry := map[string]interface{}{}

			aclEntry[mkDataSourceVirtualEnvironmentGroupACLPath] = v.Path

			if v.Propagate != nil {
				aclEntry[mkDataSourceVirtualEnvironmentGroupACLPropagate] = bool(*v.Propagate)
			} else {
				aclEntry[mkDataSourceVirtualEnvironmentGroupACLPropagate] = false
			}

			aclEntry[mkDataSourceVirtualEnvironmentGroupACLRoleID] = v.RoleID

			aclParsed = append(aclParsed, aclEntry)
		}
	}

	err = d.Set(mkDataSourceVirtualEnvironmentGroupACL, aclParsed)
	diags = append(diags, diag.FromErr(err)...)

	if group.Comment != nil {
		err = d.Set(mkDataSourceVirtualEnvironmentGroupComment, group.Comment)
	} else {
		err = d.Set(mkDataSourceVirtualEnvironmentGroupComment, "")
	}
	diags = append(diags, diag.FromErr(err)...)

	err = d.Set(mkDataSourceVirtualEnvironmentGroupMembers, group.Members)
	diags = append(diags, diag.FromErr(err)...)

	return diags
}
