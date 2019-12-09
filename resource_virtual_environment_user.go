/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package main

import (
	"strings"
	"time"

	"github.com/danitso/terraform-provider-proxmox/proxmox"
	"github.com/hashicorp/terraform/helper/schema"
)

const (
	mkResourceVirtualEnvironmentUserComment        = "comment"
	mkResourceVirtualEnvironmentUserEmail          = "email"
	mkResourceVirtualEnvironmentUserEnabled        = "enabled"
	mkResourceVirtualEnvironmentUserExpirationDate = "expiration_date"
	mkResourceVirtualEnvironmentUserFirstName      = "first_name"
	mkResourceVirtualEnvironmentUserGroups         = "groups"
	mkResourceVirtualEnvironmentUserKeys           = "keys"
	mkResourceVirtualEnvironmentUserLastName       = "last_name"
	mkResourceVirtualEnvironmentUserPassword       = "password"
	mkResourceVirtualEnvironmentUserUserID         = "user_id"
)

func resourceVirtualEnvironmentUser() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			mkResourceVirtualEnvironmentUserComment: &schema.Schema{
				Type:        schema.TypeString,
				Description: "The user comment",
				Optional:    true,
				Default:     "",
			},
			mkResourceVirtualEnvironmentUserEmail: &schema.Schema{
				Type:        schema.TypeString,
				Description: "The user's email address",
				Optional:    true,
				Default:     "",
			},
			mkResourceVirtualEnvironmentUserEnabled: &schema.Schema{
				Type:        schema.TypeBool,
				Description: "Whether the user account is enabled",
				Optional:    true,
				Default:     true,
			},
			mkResourceVirtualEnvironmentUserExpirationDate: &schema.Schema{
				Type:        schema.TypeString,
				Description: "The user account's expiration date",
				Optional:    true,
				Default:     time.Unix(0, 0).UTC().Format(time.RFC3339),
			},
			mkResourceVirtualEnvironmentUserFirstName: &schema.Schema{
				Type:        schema.TypeString,
				Description: "The user's first name",
				Optional:    true,
				Default:     "",
			},
			mkResourceVirtualEnvironmentUserGroups: &schema.Schema{
				Type:        schema.TypeSet,
				Description: "The user's groups",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				DefaultFunc: func() (interface{}, error) {
					return []string{}, nil
				},
			},
			mkResourceVirtualEnvironmentUserKeys: &schema.Schema{
				Type:        schema.TypeString,
				Description: "The user's keys",
				Optional:    true,
				Default:     "",
			},
			mkResourceVirtualEnvironmentUserLastName: &schema.Schema{
				Type:        schema.TypeString,
				Description: "The user's last name",
				Optional:    true,
				Default:     "",
			},
			mkResourceVirtualEnvironmentUserPassword: &schema.Schema{
				Type:        schema.TypeString,
				Description: "The user's password",
				Required:    true,
			},
			mkResourceVirtualEnvironmentUserUserID: &schema.Schema{
				Type:        schema.TypeString,
				Description: "The user id",
				Required:    true,
				ForceNew:    true,
			},
		},
		Create: resourceVirtualEnvironmentUserCreate,
		Read:   resourceVirtualEnvironmentUserRead,
		Update: resourceVirtualEnvironmentUserUpdate,
		Delete: resourceVirtualEnvironmentUserDelete,
	}
}

func resourceVirtualEnvironmentUserCreate(d *schema.ResourceData, m interface{}) error {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()

	if err != nil {
		return err
	}

	comment := d.Get(mkResourceVirtualEnvironmentUserComment).(string)
	email := d.Get(mkResourceVirtualEnvironmentUserEmail).(string)
	enabled := proxmox.CustomBool(d.Get(mkResourceVirtualEnvironmentUserEnabled).(bool))
	expirationDate, err := time.Parse(time.RFC3339, d.Get(mkResourceVirtualEnvironmentUserExpirationDate).(string))

	if err != nil {
		return err
	}

	expirationDateCustom := proxmox.CustomTimestamp(expirationDate)
	firstName := d.Get(mkResourceVirtualEnvironmentUserFirstName).(string)
	groups := d.Get(mkResourceVirtualEnvironmentUserGroups).(*schema.Set).List()
	groupsCustom := make([]string, len(groups))

	for i, v := range groups {
		groupsCustom[i] = v.(string)
	}

	keys := d.Get(mkResourceVirtualEnvironmentUserKeys).(string)
	lastName := d.Get(mkResourceVirtualEnvironmentUserLastName).(string)
	password := d.Get(mkResourceVirtualEnvironmentUserPassword).(string)
	userID := d.Get(mkResourceVirtualEnvironmentUserUserID).(string)

	body := &proxmox.VirtualEnvironmentUserCreateRequestBody{
		Comment:        &comment,
		Email:          &email,
		Enabled:        &enabled,
		ExpirationDate: &expirationDateCustom,
		FirstName:      &firstName,
		Groups:         groupsCustom,
		ID:             userID,
		Keys:           &keys,
		LastName:       &lastName,
		Password:       password,
	}

	err = veClient.CreateUser(body)

	if err != nil {
		return err
	}

	d.SetId(userID)

	return resourceVirtualEnvironmentUserRead(d, m)
}

func resourceVirtualEnvironmentUserRead(d *schema.ResourceData, m interface{}) error {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()

	if err != nil {
		return err
	}

	userID := d.Id()
	user, err := veClient.GetUser(userID)

	if err != nil {
		if strings.Contains(err.Error(), "HTTP 404") {
			d.SetId("")

			return nil
		}

		return err
	}

	d.SetId(userID)

	if user.Comment != nil {
		d.Set(mkResourceVirtualEnvironmentUserComment, user.Comment)
	} else {
		d.Set(mkResourceVirtualEnvironmentUserComment, "")
	}

	if user.Email != nil {
		d.Set(mkResourceVirtualEnvironmentUserEmail, user.Email)
	} else {
		d.Set(mkResourceVirtualEnvironmentUserEmail, "")
	}

	if user.Enabled != nil {
		d.Set(mkResourceVirtualEnvironmentUserEnabled, user.Enabled)
	} else {
		d.Set(mkResourceVirtualEnvironmentUserEnabled, true)
	}

	if user.ExpirationDate != nil {
		d.Set(mkResourceVirtualEnvironmentUserExpirationDate, time.Time(*user.ExpirationDate).Format(time.RFC3339))
	} else {
		d.Set(mkResourceVirtualEnvironmentUserExpirationDate, time.Unix(0, 0).UTC().Format(time.RFC3339))
	}

	if user.FirstName != nil {
		d.Set(mkResourceVirtualEnvironmentUserFirstName, user.FirstName)
	} else {
		d.Set(mkResourceVirtualEnvironmentUserFirstName, "")
	}

	groups := schema.NewSet(schema.HashString, make([]interface{}, 0))

	if user.Groups != nil {
		for _, v := range *user.Groups {
			groups.Add(v)
		}
	}

	d.Set(mkResourceVirtualEnvironmentUserGroups, groups)

	if user.Keys != nil {
		d.Set(mkResourceVirtualEnvironmentUserKeys, user.Keys)
	} else {
		d.Set(mkResourceVirtualEnvironmentUserKeys, "")
	}

	if user.LastName != nil {
		d.Set(mkResourceVirtualEnvironmentUserLastName, user.LastName)
	} else {
		d.Set(mkResourceVirtualEnvironmentUserLastName, "")
	}

	return nil
}

func resourceVirtualEnvironmentUserUpdate(d *schema.ResourceData, m interface{}) error {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()

	if err != nil {
		return err
	}

	comment := d.Get(mkResourceVirtualEnvironmentUserComment).(string)
	email := d.Get(mkResourceVirtualEnvironmentUserEmail).(string)
	enabled := proxmox.CustomBool(d.Get(mkResourceVirtualEnvironmentUserEnabled).(bool))
	expirationDate, err := time.Parse(time.RFC3339, d.Get(mkResourceVirtualEnvironmentUserExpirationDate).(string))

	if err != nil {
		return err
	}

	expirationDateCustom := proxmox.CustomTimestamp(expirationDate)
	firstName := d.Get(mkResourceVirtualEnvironmentUserFirstName).(string)
	groups := d.Get(mkResourceVirtualEnvironmentUserGroups).(*schema.Set).List()
	groupsCustom := make([]string, len(groups))

	for i, v := range groups {
		groupsCustom[i] = v.(string)
	}

	keys := d.Get(mkResourceVirtualEnvironmentUserKeys).(string)
	lastName := d.Get(mkResourceVirtualEnvironmentUserLastName).(string)

	body := &proxmox.VirtualEnvironmentUserUpdateRequestBody{
		Comment:        &comment,
		Email:          &email,
		Enabled:        &enabled,
		ExpirationDate: &expirationDateCustom,
		FirstName:      &firstName,
		Groups:         groupsCustom,
		Keys:           &keys,
		LastName:       &lastName,
	}

	userID := d.Id()
	err = veClient.UpdateUser(userID, body)

	if err != nil {
		return err
	}

	if d.HasChange(mkResourceVirtualEnvironmentUserPassword) {
		password := d.Get(mkResourceVirtualEnvironmentUserPassword).(string)
		err = veClient.ChangeUserPassword(userID, password)

		if err != nil {
			return err
		}
	}

	return resourceVirtualEnvironmentUserRead(d, m)
}

func resourceVirtualEnvironmentUserDelete(d *schema.ResourceData, m interface{}) error {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()

	if err != nil {
		return err
	}

	userID := d.Id()
	err = veClient.DeleteUser(userID)

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
