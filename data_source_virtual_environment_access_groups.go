/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package main

import (
	"github.com/hashicorp/terraform/helper/schema"
)

const (
	mkDataSourceVirtualEnvironmentAccessGroupsComments = "comments"
	mkDataSourceVirtualEnvironmentAccessGroupsGroupIDs = "group_ids"
)

func dataSourceVirtualEnvironmentAccessGroups() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			mkDataSourceVirtualEnvironmentAccessGroupsComments: &schema.Schema{
				Type:        schema.TypeList,
				Description: "The group comments",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			mkDataSourceVirtualEnvironmentAccessGroupsGroupIDs: &schema.Schema{
				Type:        schema.TypeList,
				Description: "The group ids",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
		Read: dataSourceVirtualEnvironmentAccessGroupsRead,
	}
}

func dataSourceVirtualEnvironmentAccessGroupsRead(d *schema.ResourceData, m interface{}) error {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()

	if err != nil {
		return err
	}

	list, err := veClient.ListAccessGroups()

	if err != nil {
		return err
	}

	comments := make([]interface{}, len(list))
	groupIDs := make([]interface{}, len(list))

	for i, v := range list {
		comments[i] = v.Comment
		groupIDs[i] = v.ID
	}

	d.SetId("access_groups")

	d.Set(mkDataSourceVirtualEnvironmentAccessGroupsComments, comments)
	d.Set(mkDataSourceVirtualEnvironmentAccessGroupsGroupIDs, groupIDs)

	return nil
}
