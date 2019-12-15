/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmoxtf

import (
	"github.com/hashicorp/terraform/helper/schema"
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
			mkDataSourceVirtualEnvironmentGroupACL: &schema.Schema{
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
			mkDataSourceVirtualEnvironmentGroupComment: &schema.Schema{
				Type:        schema.TypeString,
				Description: "The group comment",
				Computed:    true,
			},
			mkDataSourceVirtualEnvironmentGroupID: &schema.Schema{
				Type:        schema.TypeString,
				Description: "The group id",
				Required:    true,
			},
			mkDataSourceVirtualEnvironmentGroupMembers: &schema.Schema{
				Type:        schema.TypeSet,
				Description: "The group members",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
		Read: dataSourceVirtualEnvironmentGroupRead,
	}
}

func dataSourceVirtualEnvironmentGroupRead(d *schema.ResourceData, m interface{}) error {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()

	if err != nil {
		return err
	}

	groupID := d.Get(mkDataSourceVirtualEnvironmentGroupID).(string)
	group, err := veClient.GetGroup(groupID)

	if err != nil {
		return err
	}

	acl, err := veClient.GetACL()

	if err != nil {
		return err
	}

	d.SetId(groupID)

	aclParsed := make([]interface{}, 0)

	for _, v := range acl {
		if v.Type == "group" && v.UserOrGroupID == groupID {
			aclEntry := make(map[string]interface{})

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

	d.Set(mkDataSourceVirtualEnvironmentGroupACL, aclParsed)

	if group.Comment != nil {
		d.Set(mkDataSourceVirtualEnvironmentGroupComment, group.Comment)
	} else {
		d.Set(mkDataSourceVirtualEnvironmentGroupComment, "")
	}

	d.Set(mkDataSourceVirtualEnvironmentGroupMembers, group.Members)

	return nil
}
