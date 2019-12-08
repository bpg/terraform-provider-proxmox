/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package main

import (
	"github.com/hashicorp/terraform/helper/schema"
)

const (
	mkDataSourceVirtualEnvironmentRoleID         = "role_id"
	mkDataSourceVirtualEnvironmentRolePrivileges = "privileges"
)

func dataSourceVirtualEnvironmentRole() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			mkDataSourceVirtualEnvironmentRoleID: &schema.Schema{
				Type:        schema.TypeString,
				Description: "The role id",
				Required:    true,
			},
			mkDataSourceVirtualEnvironmentRolePrivileges: &schema.Schema{
				Type:        schema.TypeList,
				Description: "The role privileges",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
		Read: dataSourceVirtualEnvironmentRoleRead,
	}
}

func dataSourceVirtualEnvironmentRoleRead(d *schema.ResourceData, m interface{}) error {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()

	if err != nil {
		return err
	}

	roleID := d.Get(mkDataSourceVirtualEnvironmentRoleID).(string)
	accessRole, err := veClient.GetRole(roleID)

	if err != nil {
		return err
	}

	d.SetId(roleID)
	d.Set(mkDataSourceVirtualEnvironmentRolePrivileges, *accessRole)

	return nil
}
