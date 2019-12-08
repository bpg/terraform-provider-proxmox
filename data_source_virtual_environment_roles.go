/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package main

import (
	"github.com/hashicorp/terraform/helper/schema"
)

const (
	mkDataSourceVirtualEnvironmentRolesPrivileges = "privileges"
	mkDataSourceVirtualEnvironmentRolesRoleIDs    = "role_ids"
	mkDataSourceVirtualEnvironmentRolesSpecial    = "special"
)

func dataSourceVirtualEnvironmentRoles() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			mkDataSourceVirtualEnvironmentRolesPrivileges: &schema.Schema{
				Type:        schema.TypeList,
				Description: "The role privileges",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeList,
					Elem: &schema.Schema{Type: schema.TypeString},
				},
			},
			mkDataSourceVirtualEnvironmentRolesRoleIDs: &schema.Schema{
				Type:        schema.TypeList,
				Description: "The role ids",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			mkDataSourceVirtualEnvironmentRolesSpecial: &schema.Schema{
				Type:        schema.TypeList,
				Description: "Whether the role is special (built-in)",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeBool},
			},
		},
		Read: dataSourceVirtualEnvironmentRolesRead,
	}
}

func dataSourceVirtualEnvironmentRolesRead(d *schema.ResourceData, m interface{}) error {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()

	if err != nil {
		return err
	}

	list, err := veClient.ListRoles()

	if err != nil {
		return err
	}

	privileges := make([]interface{}, len(list))
	roleIDs := make([]interface{}, len(list))
	special := make([]interface{}, len(list))

	for i, v := range list {
		privileges[i] = v.Privileges
		roleIDs[i] = v.ID
		special[i] = v.Special
	}

	d.SetId("access_roles")

	d.Set(mkDataSourceVirtualEnvironmentRolesPrivileges, privileges)
	d.Set(mkDataSourceVirtualEnvironmentRolesRoleIDs, roleIDs)
	d.Set(mkDataSourceVirtualEnvironmentRolesSpecial, special)

	return nil
}
