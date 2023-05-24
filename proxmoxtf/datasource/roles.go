/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package datasource

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/bpg/terraform-provider-proxmox/proxmoxtf"
)

const (
	mkDataSourceVirtualEnvironmentRolesPrivileges = "privileges"
	mkDataSourceVirtualEnvironmentRolesRoleIDs    = "role_ids"
	mkDataSourceVirtualEnvironmentRolesSpecial    = "special"
)

// Roles returns a resource for the user roles.
func Roles() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			mkDataSourceVirtualEnvironmentRolesPrivileges: {
				Type:        schema.TypeList,
				Description: "The role privileges",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeSet,
					Elem: &schema.Schema{Type: schema.TypeString},
				},
			},
			mkDataSourceVirtualEnvironmentRolesRoleIDs: {
				Type:        schema.TypeList,
				Description: "The role ids",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			mkDataSourceVirtualEnvironmentRolesSpecial: {
				Type:        schema.TypeList,
				Description: "Whether the role is special (built-in)",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeBool},
			},
		},
		ReadContext: rolesRead,
	}
}

func rolesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	config := m.(proxmoxtf.ProviderConfiguration)
	api, err := config.GetAPI()
	if err != nil {
		return diag.FromErr(err)
	}

	list, err := api.Access().ListRoles(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	privileges := make([]interface{}, len(list))
	roleIDs := make([]interface{}, len(list))
	special := make([]interface{}, len(list))

	for i, v := range list {
		if v.Privileges != nil {
			p := schema.NewSet(schema.HashString, []interface{}{})

			for _, v := range *v.Privileges {
				p.Add(v)
			}

			privileges[i] = p
		} else {
			privileges[i] = map[string]interface{}{}
		}

		roleIDs[i] = v.ID

		if v.Special != nil {
			special[i] = v.Special
		} else {
			special[i] = false
		}
	}

	d.SetId("roles")

	err = d.Set(mkDataSourceVirtualEnvironmentRolesPrivileges, privileges)
	diags = append(diags, diag.FromErr(err)...)
	err = d.Set(mkDataSourceVirtualEnvironmentRolesRoleIDs, roleIDs)
	diags = append(diags, diag.FromErr(err)...)
	err = d.Set(mkDataSourceVirtualEnvironmentRolesSpecial, special)
	diags = append(diags, diag.FromErr(err)...)

	return diags
}
