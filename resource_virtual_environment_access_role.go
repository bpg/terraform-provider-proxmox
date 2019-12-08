/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package main

import (
	"strings"

	"github.com/danitso/terraform-provider-proxmox/proxmox"
	"github.com/hashicorp/terraform/helper/schema"
)

const (
	mkResourceVirtualEnvironmentAccessRolePrivileges = "privileges"
	mkResourceVirtualEnvironmentAccessRoleRoleID     = "role_id"
)

func resourceVirtualEnvironmentAccessRole() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			mkResourceVirtualEnvironmentAccessRolePrivileges: &schema.Schema{
				Type:        schema.TypeList,
				Description: "The role privileges",
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			mkResourceVirtualEnvironmentAccessRoleRoleID: &schema.Schema{
				Type:        schema.TypeString,
				Description: "The role id",
				Required:    true,
			},
		},
		Create: resourceVirtualEnvironmentAccessRoleCreate,
		Read:   resourceVirtualEnvironmentAccessRoleRead,
		Update: resourceVirtualEnvironmentAccessRoleUpdate,
		Delete: resourceVirtualEnvironmentAccessRoleDelete,
	}
}

func resourceVirtualEnvironmentAccessRoleCreate(d *schema.ResourceData, m interface{}) error {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()

	if err != nil {
		return err
	}

	privileges := d.Get(mkResourceVirtualEnvironmentAccessRolePrivileges).([]interface{})
	roleID := d.Get(mkResourceVirtualEnvironmentAccessRoleRoleID).(string)

	customPrivileges := make(proxmox.CustomPrivileges, len(privileges))

	for i, v := range privileges {
		customPrivileges[i] = v.(string)
	}

	body := &proxmox.VirtualEnvironmentAccessRoleCreateRequestBody{
		ID:         roleID,
		Privileges: customPrivileges,
	}

	err = veClient.CreateAccessRole(body)

	if err != nil {
		return err
	}

	d.SetId(roleID)

	return resourceVirtualEnvironmentAccessRoleRead(d, m)
}

func resourceVirtualEnvironmentAccessRoleRead(d *schema.ResourceData, m interface{}) error {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()

	if err != nil {
		return err
	}

	roleID := d.Id()
	accessRole, err := veClient.GetAccessRole(roleID)

	if err != nil {
		if strings.Contains(err.Error(), "HTTP 404") {
			d.SetId("")

			return nil
		}

		return err
	}

	d.SetId(roleID)

	d.Set(mkResourceVirtualEnvironmentAccessRolePrivileges, *accessRole)

	return nil
}

func resourceVirtualEnvironmentAccessRoleUpdate(d *schema.ResourceData, m interface{}) error {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()

	if err != nil {
		return err
	}

	privileges := d.Get(mkResourceVirtualEnvironmentAccessRolePrivileges).([]interface{})
	customPrivileges := make(proxmox.CustomPrivileges, len(privileges))

	for i, v := range privileges {
		customPrivileges[i] = v.(string)
	}

	body := &proxmox.VirtualEnvironmentAccessRoleUpdateRequestBody{
		Privileges: customPrivileges,
	}

	roleID := d.Id()
	err = veClient.UpdateAccessRole(roleID, body)

	if err != nil {
		return err
	}

	return resourceVirtualEnvironmentAccessRoleRead(d, m)
}

func resourceVirtualEnvironmentAccessRoleDelete(d *schema.ResourceData, m interface{}) error {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()

	if err != nil {
		return err
	}

	roleID := d.Id()
	err = veClient.DeleteAccessRole(roleID)

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
