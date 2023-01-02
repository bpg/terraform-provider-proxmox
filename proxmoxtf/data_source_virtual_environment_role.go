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
	mkDataSourceVirtualEnvironmentRoleID         = "role_id"
	mkDataSourceVirtualEnvironmentRolePrivileges = "privileges"
)

func dataSourceVirtualEnvironmentRole() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			mkDataSourceVirtualEnvironmentRoleID: {
				Type:        schema.TypeString,
				Description: "The role id",
				Required:    true,
			},
			mkDataSourceVirtualEnvironmentRolePrivileges: {
				Type:        schema.TypeSet,
				Description: "The role privileges",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
		ReadContext: dataSourceVirtualEnvironmentRoleRead,
	}
}

func dataSourceVirtualEnvironmentRoleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()
	if err != nil {
		return diag.FromErr(err)
	}

	roleID := d.Get(mkDataSourceVirtualEnvironmentRoleID).(string)
	accessRole, err := veClient.GetRole(ctx, roleID)
	if err != nil {
		return diag.FromErr(err)
	}

	privileges := schema.NewSet(schema.HashString, []interface{}{})

	if *accessRole != nil {
		for _, v := range *accessRole {
			privileges.Add(v)
		}
	}

	d.SetId(roleID)

	err = d.Set(mkDataSourceVirtualEnvironmentRolePrivileges, privileges)

	return diag.FromErr(err)
}
