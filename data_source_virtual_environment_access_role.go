/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package main

import (
	"github.com/hashicorp/terraform/helper/schema"
)

const (
	mkDataSourceVirtualEnvironmentAccessRoleID         = "role_id"
	mkDataSourceVirtualEnvironmentAccessRolePrivileges = "privileges"
)

func dataSourceVirtualEnvironmentAccessRole() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			mkDataSourceVirtualEnvironmentAccessRoleID: &schema.Schema{
				Type:        schema.TypeString,
				Description: "The role id",
				Required:    true,
			},
			mkDataSourceVirtualEnvironmentAccessRolePrivileges: &schema.Schema{
				Type:        schema.TypeList,
				Description: "The role privileges",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
		Read: dataSourceVirtualEnvironmentAccessRoleRead,
	}
}

func dataSourceVirtualEnvironmentAccessRoleRead(d *schema.ResourceData, m interface{}) error {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()

	if err != nil {
		return err
	}

	roleID := d.Get(mkDataSourceVirtualEnvironmentAccessRoleID).(string)
	accessRole, err := veClient.GetAccessRole(roleID)

	if err != nil {
		return err
	}

	d.SetId(roleID)
	d.Set(mkDataSourceVirtualEnvironmentAccessRolePrivileges, *accessRole)

	return nil
}
