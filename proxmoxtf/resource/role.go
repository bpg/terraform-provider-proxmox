/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package resource

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/bpg/terraform-provider-proxmox/proxmox/access"
	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
	"github.com/bpg/terraform-provider-proxmox/proxmoxtf"
)

const (
	mkResourceVirtualEnvironmentRolePrivileges = "privileges"
	mkResourceVirtualEnvironmentRoleRoleID     = "role_id"
)

// Role returns a resource that manages roles.
func Role() *schema.Resource {
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
		CreateContext: roleCreate,
		ReadContext:   roleRead,
		UpdateContext: roleUpdate,
		DeleteContext: roleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
				roleID := d.Id()

				err := d.Set(mkResourceVirtualEnvironmentRoleRoleID, roleID)
				if err != nil {
					return nil, fmt.Errorf("failed setting state during import: %w", err)
				}

				return []*schema.ResourceData{d}, nil
			},
		},
	}
}

func roleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(proxmoxtf.ProviderConfiguration)
	client, err := config.GetClient()
	if err != nil {
		return diag.FromErr(err)
	}

	privileges := d.Get(mkResourceVirtualEnvironmentRolePrivileges).(*schema.Set).List()
	customPrivileges := make(types.CustomPrivileges, len(privileges))
	roleID := d.Get(mkResourceVirtualEnvironmentRoleRoleID).(string)

	for i, v := range privileges {
		customPrivileges[i] = v.(string)
	}

	body := &access.RoleCreateRequestBody{
		ID:         roleID,
		Privileges: customPrivileges,
	}

	err = client.Access().CreateRole(ctx, body)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(roleID)

	return roleRead(ctx, d, m)
}

func roleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(proxmoxtf.ProviderConfiguration)
	client, err := config.GetClient()
	if err != nil {
		return diag.FromErr(err)
	}

	roleID := d.Id()
	role, err := client.Access().GetRole(ctx, roleID)
	if err != nil {
		if errors.Is(err, api.ErrResourceDoesNotExist) {
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

func roleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(proxmoxtf.ProviderConfiguration)
	client, err := config.GetClient()
	if err != nil {
		return diag.FromErr(err)
	}

	privileges := d.Get(mkResourceVirtualEnvironmentRolePrivileges).(*schema.Set).List()
	customPrivileges := make(types.CustomPrivileges, len(privileges))
	roleID := d.Id()

	for i, v := range privileges {
		customPrivileges[i] = v.(string)
	}

	body := &access.RoleUpdateRequestBody{
		Privileges: customPrivileges,
	}

	err = client.Access().UpdateRole(ctx, roleID, body)
	if err != nil {
		return diag.FromErr(err)
	}

	return roleRead(ctx, d, m)
}

func roleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(proxmoxtf.ProviderConfiguration)
	client, err := config.GetClient()
	if err != nil {
		return diag.FromErr(err)
	}

	roleID := d.Id()

	err = client.Access().DeleteRole(ctx, roleID)
	if err != nil && !errors.Is(err, api.ErrResourceDoesNotExist) {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
