/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmoxtf

import (
	"strings"

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
		Create: resourceVirtualEnvironmentRoleCreate,
		Read:   resourceVirtualEnvironmentRoleRead,
		Update: resourceVirtualEnvironmentRoleUpdate,
		Delete: resourceVirtualEnvironmentRoleDelete,
	}
}

func resourceVirtualEnvironmentRoleCreate(d *schema.ResourceData, m interface{}) error {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()

	if err != nil {
		return err
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

	err = veClient.CreateRole(body)

	if err != nil {
		return err
	}

	d.SetId(roleID)

	return resourceVirtualEnvironmentRoleRead(d, m)
}

func resourceVirtualEnvironmentRoleRead(d *schema.ResourceData, m interface{}) error {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()

	if err != nil {
		return err
	}

	roleID := d.Id()
	role, err := veClient.GetRole(roleID)

	if err != nil {
		if strings.Contains(err.Error(), "HTTP 404") {
			d.SetId("")

			return nil
		}

		return err
	}

	privileges := schema.NewSet(schema.HashString, []interface{}{})

	if *role != nil {
		for _, v := range *role {
			privileges.Add(v)
		}
	}

	d.Set(mkResourceVirtualEnvironmentRolePrivileges, privileges)

	return nil
}

func resourceVirtualEnvironmentRoleUpdate(d *schema.ResourceData, m interface{}) error {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()

	if err != nil {
		return err
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

	err = veClient.UpdateRole(roleID, body)

	if err != nil {
		return err
	}

	return resourceVirtualEnvironmentRoleRead(d, m)
}

func resourceVirtualEnvironmentRoleDelete(d *schema.ResourceData, m interface{}) error {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()

	if err != nil {
		return err
	}

	roleID := d.Id()
	err = veClient.DeleteRole(roleID)

	if err != nil {
		if strings.Contains(err.Error(), "HTTP 404") {
			d.SetId("")

			return nil
		}

		return err
	}

	d.SetId("")

	return nil
}
