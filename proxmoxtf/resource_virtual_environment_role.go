/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmoxtf

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/bpg/terraform-provider-proxmox/proxmox"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	mkResourceVirtualEnvironmentRolePrivileges = "privileges"
	mkResourceVirtualEnvironmentRoleRoleID     = "role_id"
)

func resourceVirtualEnvironmentRole() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			mkResourceVirtualEnvironmentRolePrivileges: {
				Type:        schema.TypeSet,
				Description: "The role privileges",
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			mkResourceVirtualEnvironmentRoleRoleID: {
				Type:        schema.TypeString,
				Description: "The role id",
				Required:    true,
				ForceNew:    true,
			},
		},
		CreateContext: resourceVirtualEnvironmentRoleCreate,
		ReadContext:   resourceVirtualEnvironmentRoleRead,
		UpdateContext: resourceVirtualEnvironmentRoleUpdate,
		DeleteContext: resourceVirtualEnvironmentRoleDelete,
	}
}

func resourceVirtualEnvironmentRoleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()
	if err != nil {
		return diag.FromErr(err)
	}

	privileges := d.Get(mkResourceVirtualEnvironmentRolePrivileges).(*schema.Set).List()
	customPrivileges := make(proxmox.CustomPrivileges, len(privileges))
	roleID := d.Get(mkResourceVirtualEnvironmentRoleRoleID).(string)

	for i, v := range privileges {
		customPrivileges[i] = v.(string)
	}

	body := &proxmox.VirtualEnvironmentRoleCreateRequestBody{
		ID:         roleID,
		Privileges: customPrivileges,
	}

	err = veClient.CreateRole(ctx, body)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(roleID)

	return resourceVirtualEnvironmentRoleRead(ctx, d, m)
}

func resourceVirtualEnvironmentRoleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()
	if err != nil {
		return diag.FromErr(err)
	}

	roleID := d.Id()
	role, err := veClient.GetRole(ctx, roleID)

	if err != nil {
		if strings.Contains(err.Error(), "HTTP 404") {
			d.SetId("")

			return nil
		}
		return diag.FromErr(err)
	}

	privileges := schema.NewSet(schema.HashString, []interface{}{})

	if *role != nil {
		for _, v := range *role {
			privileges.Add(v)
		}
	}

	err = d.Set(mkResourceVirtualEnvironmentRolePrivileges, privileges)
	return diag.FromErr(err)
}

func resourceVirtualEnvironmentRoleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()
	if err != nil {
		return diag.FromErr(err)
	}

	privileges := d.Get(mkResourceVirtualEnvironmentRolePrivileges).(*schema.Set).List()
	customPrivileges := make(proxmox.CustomPrivileges, len(privileges))
	roleID := d.Id()

	for i, v := range privileges {
		customPrivileges[i] = v.(string)
	}

	body := &proxmox.VirtualEnvironmentRoleUpdateRequestBody{
		Privileges: customPrivileges,
	}

	err = veClient.UpdateRole(ctx, roleID, body)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceVirtualEnvironmentRoleRead(ctx, d, m)
}

func resourceVirtualEnvironmentRoleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()
	if err != nil {
		return diag.FromErr(err)
	}

	roleID := d.Id()
	err = veClient.DeleteRole(ctx, roleID)

	if err != nil {
		if strings.Contains(err.Error(), "HTTP 404") {
			d.SetId("")

			return nil
		}
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
