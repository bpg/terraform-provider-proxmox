/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package main

import (
	"github.com/hashicorp/terraform/helper/schema"
)

const (
	mkDataSourceVirtualEnvironmentGroupComment = "comment"
	mkDataSourceVirtualEnvironmentGroupID      = "group_id"
	mkDataSourceVirtualEnvironmentGroupMembers = "members"
)

func dataSourceVirtualEnvironmentGroup() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
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
	accessGroup, err := veClient.GetGroup(groupID)

	if err != nil {
		return err
	}

	d.SetId(groupID)

	if accessGroup.Comment != nil {
		d.Set(mkDataSourceVirtualEnvironmentGroupComment, accessGroup.Comment)
	} else {
		d.Set(mkDataSourceVirtualEnvironmentGroupComment, "")
	}

	d.Set(mkDataSourceVirtualEnvironmentGroupMembers, accessGroup.Members)

	return nil
}
