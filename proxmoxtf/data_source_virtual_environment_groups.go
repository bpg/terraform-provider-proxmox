/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmoxtf

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	mkDataSourceVirtualEnvironmentGroupsComments = "comments"
	mkDataSourceVirtualEnvironmentGroupsGroupIDs = "group_ids"
)

func dataSourceVirtualEnvironmentGroups() *schema.Resource {
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
		Read: dataSourceVirtualEnvironmentGroupsRead,
	}
}

func dataSourceVirtualEnvironmentGroupsRead(d *schema.ResourceData, m interface{}) error {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()

	if err != nil {
		return err
	}

	list, err := veClient.ListGroups()

	if err != nil {
		return err
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

	d.Set(mkDataSourceVirtualEnvironmentGroupsComments, comments)
	d.Set(mkDataSourceVirtualEnvironmentGroupsGroupIDs, groupIDs)

	return nil
}
