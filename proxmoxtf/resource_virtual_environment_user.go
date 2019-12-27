/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmoxtf

import (
	"strings"
	"time"

	"github.com/danitso/terraform-provider-proxmox/proxmox"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

const (
	dvResourceVirtualEnvironmentUserComment   = ""
	dvResourceVirtualEnvironmentUserEmail     = ""
	dvResourceVirtualEnvironmentUserEnabled   = true
	dvResourceVirtualEnvironmentUserFirstName = ""
	dvResourceVirtualEnvironmentUserKeys      = ""
	dvResourceVirtualEnvironmentUserLastName  = ""

	mkResourceVirtualEnvironmentUserACL            = "acl"
	mkResourceVirtualEnvironmentUserACLPath        = "path"
	mkResourceVirtualEnvironmentUserACLPropagate   = "propagate"
	mkResourceVirtualEnvironmentUserACLRoleID      = "role_id"
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
			mkResourceVirtualEnvironmentUserACL: &schema.Schema{
				Type:        schema.TypeSet,
				Description: "The access control list",
				Optional:    true,
				DefaultFunc: func() (interface{}, error) {
					return make([]interface{}, 0), nil
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mkResourceVirtualEnvironmentUserACLPath: {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The path",
						},
						mkResourceVirtualEnvironmentUserACLPropagate: {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to propagate to child paths",
							Default:     false,
						},
						mkResourceVirtualEnvironmentUserACLRoleID: {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The role id",
						},
					},
				},
			},
			mkResourceVirtualEnvironmentUserComment: &schema.Schema{
				Type:        schema.TypeString,
				Description: "The user comment",
				Optional:    true,
				Default:     dvResourceVirtualEnvironmentUserComment,
			},
			mkResourceVirtualEnvironmentUserEmail: &schema.Schema{
				Type:        schema.TypeString,
				Description: "The user's email address",
				Optional:    true,
				Default:     dvResourceVirtualEnvironmentUserEmail,
			},
			mkResourceVirtualEnvironmentUserEnabled: &schema.Schema{
				Type:        schema.TypeBool,
				Description: "Whether the user account is enabled",
				Optional:    true,
				Default:     dvResourceVirtualEnvironmentUserEnabled,
			},
			mkResourceVirtualEnvironmentUserExpirationDate: &schema.Schema{
				Type:         schema.TypeString,
				Description:  "The user account's expiration date",
				Optional:     true,
				Default:      time.Unix(0, 0).UTC().Format(time.RFC3339),
				ValidateFunc: validation.ValidateRFC3339TimeString,
			},
			mkResourceVirtualEnvironmentUserFirstName: &schema.Schema{
				Type:        schema.TypeString,
				Description: "The user's first name",
				Optional:    true,
				Default:     dvResourceVirtualEnvironmentUserFirstName,
			},
			mkResourceVirtualEnvironmentUserGroups: &schema.Schema{
				Type:        schema.TypeSet,
				Description: "The user's groups",
				Optional:    true,
				DefaultFunc: func() (interface{}, error) {
					return []string{}, nil
				},
				Elem: &schema.Schema{Type: schema.TypeString},
			},
			mkResourceVirtualEnvironmentUserKeys: &schema.Schema{
				Type:        schema.TypeString,
				Description: "The user's keys",
				Optional:    true,
				Default:     dvResourceVirtualEnvironmentUserKeys,
			},
			mkResourceVirtualEnvironmentUserLastName: &schema.Schema{
				Type:        schema.TypeString,
				Description: "The user's last name",
				Optional:    true,
				Default:     dvResourceVirtualEnvironmentUserLastName,
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

	aclParsed := d.Get(mkResourceVirtualEnvironmentUserACL).(*schema.Set).List()

	for _, v := range aclParsed {
		aclDelete := proxmox.CustomBool(false)
		aclEntry := v.(map[string]interface{})
		aclPropagate := proxmox.CustomBool(aclEntry[mkResourceVirtualEnvironmentUserACLPropagate].(bool))

		aclBody := &proxmox.VirtualEnvironmentACLUpdateRequestBody{
			Delete:    &aclDelete,
			Path:      aclEntry[mkResourceVirtualEnvironmentUserACLPath].(string),
			Propagate: &aclPropagate,
			Roles:     []string{aclEntry[mkResourceVirtualEnvironmentUserACLRoleID].(string)},
			Users:     []string{userID},
		}

		err := veClient.UpdateACL(aclBody)

		if err != nil {
			return err
		}
	}

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

	acl, err := veClient.GetACL()

	if err != nil {
		return err
	}

	aclParsed := make([]interface{}, 0)

	for _, v := range acl {
		if v.Type == "user" && v.UserOrGroupID == userID {
			aclEntry := make(map[string]interface{})

			aclEntry[mkResourceVirtualEnvironmentUserACLPath] = v.Path

			if v.Propagate != nil {
				aclEntry[mkResourceVirtualEnvironmentUserACLPropagate] = bool(*v.Propagate)
			} else {
				aclEntry[mkResourceVirtualEnvironmentUserACLPropagate] = false
			}

			aclEntry[mkResourceVirtualEnvironmentUserACLRoleID] = v.RoleID

			aclParsed = append(aclParsed, aclEntry)
		}
	}

	d.Set(mkResourceVirtualEnvironmentUserACL, aclParsed)

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

	aclArgOld, aclArg := d.GetChange(mkResourceVirtualEnvironmentUserACL)
	aclParsedOld := aclArgOld.(*schema.Set).List()

	for _, v := range aclParsedOld {
		aclDelete := proxmox.CustomBool(true)
		aclEntry := v.(map[string]interface{})
		aclPropagate := proxmox.CustomBool(aclEntry[mkResourceVirtualEnvironmentUserACLPropagate].(bool))

		aclBody := &proxmox.VirtualEnvironmentACLUpdateRequestBody{
			Delete:    &aclDelete,
			Path:      aclEntry[mkResourceVirtualEnvironmentUserACLPath].(string),
			Propagate: &aclPropagate,
			Roles:     []string{aclEntry[mkResourceVirtualEnvironmentUserACLRoleID].(string)},
			Users:     []string{userID},
		}

		err := veClient.UpdateACL(aclBody)

		if err != nil {
			return err
		}
	}

	aclParsed := aclArg.(*schema.Set).List()

	for _, v := range aclParsed {
		aclDelete := proxmox.CustomBool(false)
		aclEntry := v.(map[string]interface{})
		aclPropagate := proxmox.CustomBool(aclEntry[mkResourceVirtualEnvironmentUserACLPropagate].(bool))

		aclBody := &proxmox.VirtualEnvironmentACLUpdateRequestBody{
			Delete:    &aclDelete,
			Path:      aclEntry[mkResourceVirtualEnvironmentUserACLPath].(string),
			Propagate: &aclPropagate,
			Roles:     []string{aclEntry[mkResourceVirtualEnvironmentUserACLRoleID].(string)},
			Users:     []string{userID},
		}

		err := veClient.UpdateACL(aclBody)

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

	aclParsed := d.Get(mkResourceVirtualEnvironmentUserACL).(*schema.Set).List()
	userID := d.Id()

	for _, v := range aclParsed {
		aclDelete := proxmox.CustomBool(true)
		aclEntry := v.(map[string]interface{})
		aclPropagate := proxmox.CustomBool(aclEntry[mkResourceVirtualEnvironmentUserACLPropagate].(bool))

		aclBody := &proxmox.VirtualEnvironmentACLUpdateRequestBody{
			Delete:    &aclDelete,
			Path:      aclEntry[mkResourceVirtualEnvironmentUserACLPath].(string),
			Propagate: &aclPropagate,
			Roles:     []string{aclEntry[mkResourceVirtualEnvironmentUserACLRoleID].(string)},
			Users:     []string{userID},
		}

		err := veClient.UpdateACL(aclBody)

		if err != nil {
			return err
		}
	}

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
