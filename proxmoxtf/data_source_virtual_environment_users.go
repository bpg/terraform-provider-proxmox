/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmoxtf

import (
	"time"

	"github.com/hashicorp/terraform/helper/schema"
)

const (
	mkDataSourceVirtualEnvironmentUsersComments        = "comments"
	mkDataSourceVirtualEnvironmentUsersEmails          = "emails"
	mkDataSourceVirtualEnvironmentUsersEnabled         = "enabled"
	mkDataSourceVirtualEnvironmentUsersExpirationDates = "expiration_dates"
	mkDataSourceVirtualEnvironmentUsersFirstNames      = "first_names"
	mkDataSourceVirtualEnvironmentUsersGroups          = "groups"
	mkDataSourceVirtualEnvironmentUsersKeys            = "keys"
	mkDataSourceVirtualEnvironmentUsersLastNames       = "last_names"
	mkDataSourceVirtualEnvironmentUsersUserIDs         = "user_ids"
)

func dataSourceVirtualEnvironmentUsers() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			mkDataSourceVirtualEnvironmentUsersComments: &schema.Schema{
				Type:        schema.TypeList,
				Description: "The user comments",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			mkDataSourceVirtualEnvironmentUsersEmails: &schema.Schema{
				Type:        schema.TypeList,
				Description: "The users' email addresses",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			mkDataSourceVirtualEnvironmentUsersEnabled: &schema.Schema{
				Type:        schema.TypeList,
				Description: "Whether a user account is enabled",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeBool},
			},
			mkDataSourceVirtualEnvironmentUsersExpirationDates: &schema.Schema{
				Type:        schema.TypeList,
				Description: "The user accounts' expiration dates",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			mkDataSourceVirtualEnvironmentUsersFirstNames: &schema.Schema{
				Type:        schema.TypeList,
				Description: "The users' first names",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			mkDataSourceVirtualEnvironmentUsersGroups: &schema.Schema{
				Type:        schema.TypeList,
				Description: "The users' groups",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeList,
					Elem: &schema.Schema{Type: schema.TypeString},
				},
			},
			mkDataSourceVirtualEnvironmentUsersKeys: &schema.Schema{
				Type:        schema.TypeList,
				Description: "The users' keys",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			mkDataSourceVirtualEnvironmentUsersLastNames: &schema.Schema{
				Type:        schema.TypeList,
				Description: "The users' last names",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			mkDataSourceVirtualEnvironmentUsersUserIDs: &schema.Schema{
				Type:        schema.TypeList,
				Description: "The user ids",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
		Read: dataSourceVirtualEnvironmentUsersRead,
	}
}

func dataSourceVirtualEnvironmentUsersRead(d *schema.ResourceData, m interface{}) error {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()

	if err != nil {
		return err
	}

	list, err := veClient.ListUsers()

	if err != nil {
		return err
	}

	comments := make([]interface{}, len(list))
	emails := make([]interface{}, len(list))
	enabled := make([]interface{}, len(list))
	expirationDates := make([]interface{}, len(list))
	firstNames := make([]interface{}, len(list))
	groups := make([]interface{}, len(list))
	keys := make([]interface{}, len(list))
	lastNames := make([]interface{}, len(list))
	userIDs := make([]interface{}, len(list))

	for i, v := range list {
		if v.Comment != nil {
			comments[i] = v.Comment
		} else {
			comments[i] = ""
		}

		if v.Email != nil {
			emails[i] = v.Email
		} else {
			emails[i] = ""
		}

		if v.Enabled != nil {
			enabled[i] = v.Enabled
		} else {
			enabled[i] = true
		}

		if v.ExpirationDate != nil {
			t := time.Time(*v.ExpirationDate)

			if t.Unix() > 0 {
				expirationDates[i] = t.UTC().Format(time.RFC3339)
			} else {
				expirationDates[i] = time.Unix(0, 0).UTC().Format(time.RFC3339)
			}
		} else {
			expirationDates[i] = time.Unix(0, 0).UTC().Format(time.RFC3339)
		}

		if v.FirstName != nil {
			firstNames[i] = v.FirstName
		} else {
			firstNames[i] = ""
		}

		if v.Groups != nil {
			groups[i] = v.Groups
		} else {
			groups[i] = make([]string, 0)
		}

		if v.Keys != nil {
			keys[i] = v.Keys
		} else {
			keys[i] = ""
		}

		if v.LastName != nil {
			lastNames[i] = v.LastName
		} else {
			lastNames[i] = ""
		}

		userIDs[i] = v.ID
	}

	d.SetId("users")

	d.Set(mkDataSourceVirtualEnvironmentUsersComments, comments)
	d.Set(mkDataSourceVirtualEnvironmentUsersEmails, emails)
	d.Set(mkDataSourceVirtualEnvironmentUsersEnabled, enabled)
	d.Set(mkDataSourceVirtualEnvironmentUsersExpirationDates, expirationDates)
	d.Set(mkDataSourceVirtualEnvironmentUsersFirstNames, firstNames)
	d.Set(mkDataSourceVirtualEnvironmentUsersGroups, groups)
	d.Set(mkDataSourceVirtualEnvironmentUsersKeys, keys)
	d.Set(mkDataSourceVirtualEnvironmentUsersLastNames, lastNames)
	d.Set(mkDataSourceVirtualEnvironmentUsersUserIDs, userIDs)

	return nil
}
