/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package main

import (
	"errors"

	"github.com/hashicorp/terraform/helper/schema"
)

const (
	mkDataSourceVirtualEnvironmentAccessGroupComment = "comment"
	mkDataSourceVirtualEnvironmentAccessGroupID      = "id"
	mkDataSourceVirtualEnvironmentAccessGroupMembers = "members"
)

func dataSourceVirtualEnvironmentAccessGroup() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			mkDataSourceVirtualEnvironmentAccessGroupComment: &schema.Schema{
				Type:        schema.TypeString,
				Description: "The group comment",
				Computed:    true,
			},
			mkDataSourceVirtualEnvironmentAccessGroupID: &schema.Schema{
				Type:        schema.TypeString,
				Description: "The group id",
				Required:    true,
			},
			mkDataSourceVirtualEnvironmentAccessGroupMembers: &schema.Schema{
				Type:        schema.TypeList,
				Description: "The group members",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
		Read: dataSourceVirtualEnvironmentAccessGroupRead,
	}
}

func dataSourceVirtualEnvironmentAccessGroupRead(d *schema.ResourceData, m interface{}) error {
	config := m.(providerConfiguration)

	if config.veClient == nil {
		return errors.New("You must specify the virtual environment details in the provider configuration to use this data source")
	}

	groupID := d.Get(mkDataSourceVirtualEnvironmentAccessGroupID).(string)
	accessGroup, err := config.veClient.GetAccessGroup(groupID)

	if err != nil {
		return err
	}

	d.SetId(groupID)

	d.Set(mkDataSourceVirtualEnvironmentAccessGroupComment, accessGroup.Comment)
	d.Set(mkDataSourceVirtualEnvironmentAccessGroupMembers, accessGroup.Members)

	return nil
}
