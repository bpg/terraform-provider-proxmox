/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmoxtf

import (
	"time"

	"github.com/hashicorp/terraform/helper/schema"
)

const (
	mkDataSourceVirtualEnvironmentUserACL            = "acl"
	mkDataSourceVirtualEnvironmentUserACLPath        = "path"
	mkDataSourceVirtualEnvironmentUserACLPropagate   = "propagate"
	mkDataSourceVirtualEnvironmentUserACLRoleID      = "role_id"
	mkDataSourceVirtualEnvironmentUserComment        = "comment"
	mkDataSourceVirtualEnvironmentUserEmail          = "email"
	mkDataSourceVirtualEnvironmentUserEnabled        = "enabled"
	mkDataSourceVirtualEnvironmentUserExpirationDate = "expiration_date"
	mkDataSourceVirtualEnvironmentUserFirstName      = "first_name"
	mkDataSourceVirtualEnvironmentUserGroups         = "groups"
	mkDataSourceVirtualEnvironmentUserKeys           = "keys"
	mkDataSourceVirtualEnvironmentUserLastName       = "last_name"
	mkDataSourceVirtualEnvironmentUserUserID         = "user_id"
)

func dataSourceVirtualEnvironmentUser() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			mkDataSourceVirtualEnvironmentUserACL: &schema.Schema{
				Type:        schema.TypeSet,
				Description: "The access control list",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mkDataSourceVirtualEnvironmentUserACLPath: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The path",
						},
						mkDataSourceVirtualEnvironmentUserACLPropagate: {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether to propagate to child paths",
						},
						mkDataSourceVirtualEnvironmentUserACLRoleID: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The role id",
						},
					},
				},
			},
			mkDataSourceVirtualEnvironmentUserComment: &schema.Schema{
				Type:        schema.TypeString,
				Description: "The user comment",
				Computed:    true,
			},
			mkDataSourceVirtualEnvironmentUserEmail: &schema.Schema{
				Type:        schema.TypeString,
				Description: "The user's email address",
				Computed:    true,
			},
			mkDataSourceVirtualEnvironmentUserEnabled: &schema.Schema{
				Type:        schema.TypeBool,
				Description: "Whether the user account is enabled",
				Computed:    true,
			},
			mkDataSourceVirtualEnvironmentUserExpirationDate: &schema.Schema{
				Type:        schema.TypeString,
				Description: "The user account's expiration date",
				Computed:    true,
			},
			mkDataSourceVirtualEnvironmentUserFirstName: &schema.Schema{
				Type:        schema.TypeString,
				Description: "The user's first name",
				Computed:    true,
			},
			mkDataSourceVirtualEnvironmentUserGroups: &schema.Schema{
				Type:        schema.TypeList,
				Description: "The user's groups",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			mkDataSourceVirtualEnvironmentUserKeys: &schema.Schema{
				Type:        schema.TypeString,
				Description: "The user's keys",
				Computed:    true,
			},
			mkDataSourceVirtualEnvironmentUserLastName: &schema.Schema{
				Type:        schema.TypeString,
				Description: "The user's last name",
				Computed:    true,
			},
			mkDataSourceVirtualEnvironmentUserUserID: &schema.Schema{
				Type:        schema.TypeString,
				Description: "The user id",
				Required:    true,
			},
		},
		Read: dataSourceVirtualEnvironmentUserRead,
	}
}

func dataSourceVirtualEnvironmentUserRead(d *schema.ResourceData, m interface{}) error {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()

	if err != nil {
		return err
	}

	userID := d.Get(mkDataSourceVirtualEnvironmentUserUserID).(string)
	v, err := veClient.GetUser(userID)

	if err != nil {
		return err
	}

	acl, err := veClient.GetACL()

	if err != nil {
		return err
	}

	d.SetId(userID)

	aclParsed := make([]interface{}, 0)

	for _, v := range acl {
		if v.Type == "user" && v.UserOrGroupID == userID {
			aclEntry := make(map[string]interface{})

			aclEntry[mkDataSourceVirtualEnvironmentUserACLPath] = v.Path

			if v.Propagate != nil {
				aclEntry[mkDataSourceVirtualEnvironmentUserACLPropagate] = bool(*v.Propagate)
			} else {
				aclEntry[mkDataSourceVirtualEnvironmentUserACLPropagate] = false
			}

			aclEntry[mkDataSourceVirtualEnvironmentUserACLRoleID] = v.RoleID

			aclParsed = append(aclParsed, aclEntry)
		}
	}

	d.Set(mkDataSourceVirtualEnvironmentUserACL, aclParsed)

	if v.Comment != nil {
		d.Set(mkDataSourceVirtualEnvironmentUserComment, v.Comment)
	} else {
		d.Set(mkDataSourceVirtualEnvironmentUserComment, "")
	}

	if v.Email != nil {
		d.Set(mkDataSourceVirtualEnvironmentUserEmail, v.Email)
	} else {
		d.Set(mkDataSourceVirtualEnvironmentUserEmail, "")
	}

	if v.Enabled != nil {
		d.Set(mkDataSourceVirtualEnvironmentUserEnabled, v.Enabled)
	} else {
		d.Set(mkDataSourceVirtualEnvironmentUserEnabled, true)
	}

	if v.ExpirationDate != nil {
		t := time.Time(*v.ExpirationDate)

		if t.Unix() > 0 {
			d.Set(mkDataSourceVirtualEnvironmentUserExpirationDate, t.UTC().Format(time.RFC3339))
		} else {
			d.Set(mkDataSourceVirtualEnvironmentUserExpirationDate, time.Unix(0, 0).UTC().Format(time.RFC3339))
		}
	} else {
		d.Set(mkDataSourceVirtualEnvironmentUserExpirationDate, time.Unix(0, 0).UTC().Format(time.RFC3339))
	}

	if v.FirstName != nil {
		d.Set(mkDataSourceVirtualEnvironmentUserFirstName, v.FirstName)
	} else {
		d.Set(mkDataSourceVirtualEnvironmentUserFirstName, "")
	}

	if v.Groups != nil {
		d.Set(mkDataSourceVirtualEnvironmentUserGroups, v.Groups)
	} else {
		d.Set(mkDataSourceVirtualEnvironmentUserGroups, "")
	}

	if v.Keys != nil {
		d.Set(mkDataSourceVirtualEnvironmentUserGroups, v.Keys)
	} else {
		d.Set(mkDataSourceVirtualEnvironmentUserGroups, "")
	}

	if v.LastName != nil {
		d.Set(mkDataSourceVirtualEnvironmentUserLastName, v.LastName)
	} else {
		d.Set(mkDataSourceVirtualEnvironmentUserLastName, "")
	}

	return nil
}
