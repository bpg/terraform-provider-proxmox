/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package datasource

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/bpg/terraform-provider-proxmox/proxmoxtf"
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

// User returns a resource for a single Proxmox user.
func User() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			mkDataSourceVirtualEnvironmentUserACL: {
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
			mkDataSourceVirtualEnvironmentUserComment: {
				Type:        schema.TypeString,
				Description: "The user comment",
				Computed:    true,
			},
			mkDataSourceVirtualEnvironmentUserEmail: {
				Type:        schema.TypeString,
				Description: "The user's email address",
				Computed:    true,
			},
			mkDataSourceVirtualEnvironmentUserEnabled: {
				Type:        schema.TypeBool,
				Description: "Whether the user account is enabled",
				Computed:    true,
			},
			mkDataSourceVirtualEnvironmentUserExpirationDate: {
				Type:        schema.TypeString,
				Description: "The user account's expiration date",
				Computed:    true,
			},
			mkDataSourceVirtualEnvironmentUserFirstName: {
				Type:        schema.TypeString,
				Description: "The user's first name",
				Computed:    true,
			},
			mkDataSourceVirtualEnvironmentUserGroups: {
				Type:        schema.TypeList,
				Description: "The user's groups",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			mkDataSourceVirtualEnvironmentUserKeys: {
				Type:        schema.TypeString,
				Description: "The user's keys",
				Computed:    true,
			},
			mkDataSourceVirtualEnvironmentUserLastName: {
				Type:        schema.TypeString,
				Description: "The user's last name",
				Computed:    true,
			},
			mkDataSourceVirtualEnvironmentUserUserID: {
				Type:        schema.TypeString,
				Description: "The user id",
				Required:    true,
			},
		},
		ReadContext: userRead,
	}
}

func userRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	config := m.(proxmoxtf.ProviderConfiguration)
	api, err := config.GetClient()
	if err != nil {
		return diag.FromErr(err)
	}

	userID := d.Get(mkDataSourceVirtualEnvironmentUserUserID).(string)
	v, err := api.Access().GetUser(ctx, userID)
	if err != nil {
		return diag.FromErr(err)
	}

	acl, err := api.Access().GetACL(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(userID)

	var aclParsed []interface{}

	for _, v := range acl {
		if v.Type == "user" && v.UserOrGroupID == userID {
			aclEntry := map[string]interface{}{}

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

	err = d.Set(mkDataSourceVirtualEnvironmentUserACL, aclParsed)
	diags = append(diags, diag.FromErr(err)...)

	if v.Comment != nil {
		err = d.Set(mkDataSourceVirtualEnvironmentUserComment, v.Comment)
	} else {
		err = d.Set(mkDataSourceVirtualEnvironmentUserComment, "")
	}
	diags = append(diags, diag.FromErr(err)...)

	if v.Email != nil {
		err = d.Set(mkDataSourceVirtualEnvironmentUserEmail, v.Email)
	} else {
		err = d.Set(mkDataSourceVirtualEnvironmentUserEmail, "")
	}
	diags = append(diags, diag.FromErr(err)...)

	if v.Enabled != nil {
		err = d.Set(mkDataSourceVirtualEnvironmentUserEnabled, v.Enabled)
	} else {
		err = d.Set(mkDataSourceVirtualEnvironmentUserEnabled, true)
	}
	diags = append(diags, diag.FromErr(err)...)

	if v.ExpirationDate != nil {
		t := time.Time(*v.ExpirationDate)
		if t.Unix() > 0 {
			err = d.Set(
				mkDataSourceVirtualEnvironmentUserExpirationDate,
				t.UTC().Format(time.RFC3339),
			)
		} else {
			err = d.Set(mkDataSourceVirtualEnvironmentUserExpirationDate, time.Unix(0, 0).UTC().Format(time.RFC3339))
		}
	} else {
		err = d.Set(mkDataSourceVirtualEnvironmentUserExpirationDate, time.Unix(0, 0).UTC().Format(time.RFC3339))
	}
	diags = append(diags, diag.FromErr(err)...)

	if v.FirstName != nil {
		err = d.Set(mkDataSourceVirtualEnvironmentUserFirstName, v.FirstName)
	} else {
		err = d.Set(mkDataSourceVirtualEnvironmentUserFirstName, "")
	}
	diags = append(diags, diag.FromErr(err)...)

	if v.Groups != nil {
		err = d.Set(mkDataSourceVirtualEnvironmentUserGroups, v.Groups)
	} else {
		err = d.Set(mkDataSourceVirtualEnvironmentUserGroups, []string{})
	}
	diags = append(diags, diag.FromErr(err)...)

	if v.Keys != nil {
		err = d.Set(mkDataSourceVirtualEnvironmentUsersKeys, v.Keys)
	} else {
		err = d.Set(mkDataSourceVirtualEnvironmentUsersKeys, "")
	}
	diags = append(diags, diag.FromErr(err)...)

	if v.LastName != nil {
		err = d.Set(mkDataSourceVirtualEnvironmentUserLastName, v.LastName)
	} else {
		err = d.Set(mkDataSourceVirtualEnvironmentUserLastName, "")
	}
	diags = append(diags, diag.FromErr(err)...)

	return diags
}
