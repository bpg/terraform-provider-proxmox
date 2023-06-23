/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package resource

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/bpg/terraform-provider-proxmox/internal/types"
	"github.com/bpg/terraform-provider-proxmox/proxmox/access"
	"github.com/bpg/terraform-provider-proxmox/proxmoxtf"
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

// User returns a resource that manages a user in the Proxmox VE access control list.
func User() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			mkResourceVirtualEnvironmentUserACL: {
				Type:        schema.TypeSet,
				Description: "The access control list",
				Optional:    true,
				DefaultFunc: func() (interface{}, error) {
					return []interface{}{}, nil
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
			mkResourceVirtualEnvironmentUserComment: {
				Type:        schema.TypeString,
				Description: "The user comment",
				Optional:    true,
				Default:     dvResourceVirtualEnvironmentUserComment,
			},
			mkResourceVirtualEnvironmentUserEmail: {
				Type:        schema.TypeString,
				Description: "The user's email address",
				Optional:    true,
				Default:     dvResourceVirtualEnvironmentUserEmail,
			},
			mkResourceVirtualEnvironmentUserEnabled: {
				Type:        schema.TypeBool,
				Description: "Whether the user account is enabled",
				Optional:    true,
				Default:     dvResourceVirtualEnvironmentUserEnabled,
			},
			mkResourceVirtualEnvironmentUserExpirationDate: {
				Type:         schema.TypeString,
				Description:  "The user account's expiration date",
				Optional:     true,
				Default:      time.Unix(0, 0).UTC().Format(time.RFC3339),
				ValidateFunc: validation.IsRFC3339Time,
			},
			mkResourceVirtualEnvironmentUserFirstName: {
				Type:        schema.TypeString,
				Description: "The user's first name",
				Optional:    true,
				Default:     dvResourceVirtualEnvironmentUserFirstName,
			},
			mkResourceVirtualEnvironmentUserGroups: {
				Type:        schema.TypeSet,
				Description: "The user's groups",
				Optional:    true,
				DefaultFunc: func() (interface{}, error) {
					return []string{}, nil
				},
				Elem: &schema.Schema{Type: schema.TypeString},
			},
			mkResourceVirtualEnvironmentUserKeys: {
				Type:        schema.TypeString,
				Description: "The user's keys",
				Optional:    true,
				Default:     dvResourceVirtualEnvironmentUserKeys,
			},
			mkResourceVirtualEnvironmentUserLastName: {
				Type:        schema.TypeString,
				Description: "The user's last name",
				Optional:    true,
				Default:     dvResourceVirtualEnvironmentUserLastName,
			},
			mkResourceVirtualEnvironmentUserPassword: {
				Type:        schema.TypeString,
				Description: "The user's password",
				Required:    true,
			},
			mkResourceVirtualEnvironmentUserUserID: {
				Type:        schema.TypeString,
				Description: "The user id",
				Required:    true,
				ForceNew:    true,
			},
		},
		CreateContext: userCreate,
		ReadContext:   userRead,
		UpdateContext: userUpdate,
		DeleteContext: userDelete,
	}
}

func userCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(proxmoxtf.ProviderConfiguration)
	api, err := config.GetClient()
	if err != nil {
		return diag.FromErr(err)
	}

	comment := d.Get(mkResourceVirtualEnvironmentUserComment).(string)
	email := d.Get(mkResourceVirtualEnvironmentUserEmail).(string)
	enabled := types.CustomBool(d.Get(mkResourceVirtualEnvironmentUserEnabled).(bool))
	expirationDate, err := time.Parse(
		time.RFC3339,
		d.Get(mkResourceVirtualEnvironmentUserExpirationDate).(string),
	)
	if err != nil {
		return diag.FromErr(err)
	}

	expirationDateCustom := types.CustomTimestamp(expirationDate)
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

	body := &access.UserCreateRequestBody{
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

	err = api.Access().CreateUser(ctx, body)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(userID)

	aclParsed := d.Get(mkResourceVirtualEnvironmentUserACL).(*schema.Set).List()

	for _, v := range aclParsed {
		aclDelete := types.CustomBool(false)
		aclEntry := v.(map[string]interface{})
		aclPropagate := types.CustomBool(
			aclEntry[mkResourceVirtualEnvironmentUserACLPropagate].(bool),
		)

		aclBody := &access.ACLUpdateRequestBody{
			Delete:    &aclDelete,
			Path:      aclEntry[mkResourceVirtualEnvironmentUserACLPath].(string),
			Propagate: &aclPropagate,
			Roles:     []string{aclEntry[mkResourceVirtualEnvironmentUserACLRoleID].(string)},
			Users:     []string{userID},
		}

		err := api.Access().UpdateACL(ctx, aclBody)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return userRead(ctx, d, m)
}

func userRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(proxmoxtf.ProviderConfiguration)
	api, err := config.GetClient()
	if err != nil {
		return diag.FromErr(err)
	}

	userID := d.Id()
	user, err := api.Access().GetUser(ctx, userID)
	if err != nil {
		if strings.Contains(err.Error(), "HTTP 404") {
			d.SetId("")

			return nil
		}
		return diag.FromErr(err)
	}

	acl, err := api.Access().GetACL(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	var aclParsed []interface{}

	for _, v := range acl {
		if v.Type == "user" && v.UserOrGroupID == userID {
			aclEntry := map[string]interface{}{}

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

	var diags diag.Diagnostics

	err = d.Set(mkResourceVirtualEnvironmentUserACL, aclParsed)
	diags = append(diags, diag.FromErr(err)...)

	if user.Comment != nil {
		err = d.Set(mkResourceVirtualEnvironmentUserComment, user.Comment)
	} else {
		err = d.Set(mkResourceVirtualEnvironmentUserComment, "")
	}
	diags = append(diags, diag.FromErr(err)...)

	if user.Email != nil {
		err = d.Set(mkResourceVirtualEnvironmentUserEmail, user.Email)
	} else {
		err = d.Set(mkResourceVirtualEnvironmentUserEmail, "")
	}
	diags = append(diags, diag.FromErr(err)...)

	if user.Enabled != nil {
		err = d.Set(mkResourceVirtualEnvironmentUserEnabled, user.Enabled)
	} else {
		err = d.Set(mkResourceVirtualEnvironmentUserEnabled, true)
	}
	diags = append(diags, diag.FromErr(err)...)

	if user.ExpirationDate != nil {
		err = d.Set(
			mkResourceVirtualEnvironmentUserExpirationDate,
			time.Time(*user.ExpirationDate).Format(time.RFC3339),
		)
	} else {
		err = d.Set(mkResourceVirtualEnvironmentUserExpirationDate, time.Unix(0, 0).UTC().Format(time.RFC3339))
	}
	diags = append(diags, diag.FromErr(err)...)

	if user.FirstName != nil {
		err = d.Set(mkResourceVirtualEnvironmentUserFirstName, user.FirstName)
	} else {
		err = d.Set(mkResourceVirtualEnvironmentUserFirstName, "")
	}
	diags = append(diags, diag.FromErr(err)...)

	groups := schema.NewSet(schema.HashString, []interface{}{})

	if user.Groups != nil {
		for _, v := range *user.Groups {
			groups.Add(v)
		}
	}

	err = d.Set(mkResourceVirtualEnvironmentUserGroups, groups)
	diags = append(diags, diag.FromErr(err)...)

	if user.Keys != nil {
		err = d.Set(mkResourceVirtualEnvironmentUserKeys, user.Keys)
	} else {
		err = d.Set(mkResourceVirtualEnvironmentUserKeys, "")
	}
	diags = append(diags, diag.FromErr(err)...)

	if user.LastName != nil {
		err = d.Set(mkResourceVirtualEnvironmentUserLastName, user.LastName)
	} else {
		err = d.Set(mkResourceVirtualEnvironmentUserLastName, "")
	}
	diags = append(diags, diag.FromErr(err)...)

	return diags
}

func userUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(proxmoxtf.ProviderConfiguration)
	api, err := config.GetClient()
	if err != nil {
		return diag.FromErr(err)
	}

	comment := d.Get(mkResourceVirtualEnvironmentUserComment).(string)
	email := d.Get(mkResourceVirtualEnvironmentUserEmail).(string)
	enabled := types.CustomBool(d.Get(mkResourceVirtualEnvironmentUserEnabled).(bool))
	expirationDate, err := time.Parse(
		time.RFC3339,
		d.Get(mkResourceVirtualEnvironmentUserExpirationDate).(string),
	)
	if err != nil {
		return diag.FromErr(err)
	}

	expirationDateCustom := types.CustomTimestamp(expirationDate)
	firstName := d.Get(mkResourceVirtualEnvironmentUserFirstName).(string)
	groups := d.Get(mkResourceVirtualEnvironmentUserGroups).(*schema.Set).List()
	groupsCustom := make([]string, len(groups))

	for i, v := range groups {
		groupsCustom[i] = v.(string)
	}

	keys := d.Get(mkResourceVirtualEnvironmentUserKeys).(string)
	lastName := d.Get(mkResourceVirtualEnvironmentUserLastName).(string)

	body := &access.UserUpdateRequestBody{
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
	err = api.Access().UpdateUser(ctx, userID, body)
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange(mkResourceVirtualEnvironmentUserPassword) {
		password := d.Get(mkResourceVirtualEnvironmentUserPassword).(string)
		err = api.Access().ChangeUserPassword(ctx, userID, password)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	aclArgOld, aclArg := d.GetChange(mkResourceVirtualEnvironmentUserACL)
	aclParsedOld := aclArgOld.(*schema.Set).List()

	for _, v := range aclParsedOld {
		aclDelete := types.CustomBool(true)
		aclEntry := v.(map[string]interface{})
		aclPropagate := types.CustomBool(
			aclEntry[mkResourceVirtualEnvironmentUserACLPropagate].(bool),
		)

		aclBody := &access.ACLUpdateRequestBody{
			Delete:    &aclDelete,
			Path:      aclEntry[mkResourceVirtualEnvironmentUserACLPath].(string),
			Propagate: &aclPropagate,
			Roles:     []string{aclEntry[mkResourceVirtualEnvironmentUserACLRoleID].(string)},
			Users:     []string{userID},
		}

		err := api.Access().UpdateACL(ctx, aclBody)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	aclParsed := aclArg.(*schema.Set).List()

	for _, v := range aclParsed {
		aclDelete := types.CustomBool(false)
		aclEntry := v.(map[string]interface{})
		aclPropagate := types.CustomBool(
			aclEntry[mkResourceVirtualEnvironmentUserACLPropagate].(bool),
		)

		aclBody := &access.ACLUpdateRequestBody{
			Delete:    &aclDelete,
			Path:      aclEntry[mkResourceVirtualEnvironmentUserACLPath].(string),
			Propagate: &aclPropagate,
			Roles:     []string{aclEntry[mkResourceVirtualEnvironmentUserACLRoleID].(string)},
			Users:     []string{userID},
		}

		err := api.Access().UpdateACL(ctx, aclBody)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return userRead(ctx, d, m)
}

func userDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(proxmoxtf.ProviderConfiguration)
	api, err := config.GetClient()
	if err != nil {
		return diag.FromErr(err)
	}

	aclParsed := d.Get(mkResourceVirtualEnvironmentUserACL).(*schema.Set).List()
	userID := d.Id()

	for _, v := range aclParsed {
		aclDelete := types.CustomBool(true)
		aclEntry := v.(map[string]interface{})
		aclPropagate := types.CustomBool(
			aclEntry[mkResourceVirtualEnvironmentUserACLPropagate].(bool),
		)

		aclBody := &access.ACLUpdateRequestBody{
			Delete:    &aclDelete,
			Path:      aclEntry[mkResourceVirtualEnvironmentUserACLPath].(string),
			Propagate: &aclPropagate,
			Roles:     []string{aclEntry[mkResourceVirtualEnvironmentUserACLRoleID].(string)},
			Users:     []string{userID},
		}

		err = api.Access().UpdateACL(ctx, aclBody)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	err = api.Access().DeleteUser(ctx, userID)

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
