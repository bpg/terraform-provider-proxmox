/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package main

import (
	"errors"

	"github.com/hashicorp/terraform/helper/schema"
)

const (
	mkDataSourceVirtualEnvironmentAccessGroupsComments = "comments"
	mkDataSourceVirtualEnvironmentAccessGroupsIDs      = "ids"
)

// dataSourceVirtualEnvironmentAccessGroups retrieves a list of access groups.
func dataSourceVirtualEnvironmentAccessGroups() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			mkDataSourceVirtualEnvironmentAccessGroupsComments: &schema.Schema{
				Type:        schema.TypeList,
				Description: "The group comments",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			mkDataSourceVirtualEnvironmentAccessGroupsIDs: &schema.Schema{
				Type:        schema.TypeList,
				Description: "The group ids",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
		Read: dataSourceVirtualEnvironmentAccessGroupsRead,
	}
}

// dataSourceVirtualEnvironmentAccessGroupsRead retrieves a list of access groups.
func dataSourceVirtualEnvironmentAccessGroupsRead(d *schema.ResourceData, m interface{}) error {
	config := m.(providerConfiguration)

	if config.veClient == nil {
		return errors.New("You must specify the virtual environment details in the provider configuration to use this data source")
	}

	list, err := config.veClient.ListAccessGroups()

	if err != nil {
		return err
	}

	comments := make([]interface{}, len(list))
	ids := make([]interface{}, len(list))

	for i, v := range list {
		comments[i] = v.Comment
		ids[i] = v.ID
	}

	d.SetId("access_groups")

	d.Set(mkDataSourceVirtualEnvironmentAccessGroupsComments, comments)
	d.Set(mkDataSourceVirtualEnvironmentAccessGroupsIDs, ids)

	return nil
}
