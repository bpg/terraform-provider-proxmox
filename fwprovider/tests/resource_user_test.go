/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/proxmox/access"
)

func TestAccResourceUser(t *testing.T) {
	t.Parallel()

	te := initTestEnvironment(t)

	username := fmt.Sprintf("%s@pve", gofakeit.Username())
	te.addTemplateVars(map[string]any{
		"Username": username,
	})

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"create and update user", []resource.TestStep{
			{
				Config: te.renderConfig(`resource "proxmox_virtual_environment_user" "user" {
					  comment  			= "Managed by Terraform"
					  email 			= "{{.Username}}"
					  enabled 			= true
					  expiration_date 	= "2034-01-01T22:00:00Z"
					  first_name 		= "First"
					  last_name 		= "Last"
					  user_id  			= "{{.Username}}"
				}`),
				Check: testResourceAttributes("proxmox_virtual_environment_user.user", map[string]string{
					"comment":         "Managed by Terraform",
					"email":           username,
					"enabled":         "true",
					"expiration_date": "2034-01-01T22:00:00Z",
					"first_name":      "First",
					"last_name":       "Last",
					"user_id":         username,
				}),
			},
			{
				Config: te.renderConfig(`resource "proxmox_virtual_environment_user" "user" {
					  enabled 			= false
					  expiration_date 	= "2035-01-01T22:00:00Z"
					  user_id  			= "{{.Username}}"
					  first_name 		= "First One"
				}`),
				Check: testResourceAttributes("proxmox_virtual_environment_user.user", map[string]string{
					"enabled":         "false",
					"expiration_date": "2035-01-01T22:00:00Z",
					"first_name":      "First One",
					"user_id":         username,
				}),
			},
			{
				ResourceName:      "proxmox_virtual_environment_user.user",
				ImportState:       true,
				ImportStateVerify: true,
			},
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: te.accProviders,
				Steps:                    tt.steps,
			})
		})
	}
}

func TestAccResourceUserToken(t *testing.T) {
	t.Parallel()

	te := initTestEnvironment(t)
	username := fmt.Sprintf("%s@pve", gofakeit.Username())
	tokenName := gofakeit.Word()

	te.addTemplateVars(map[string]any{
		"Username":  username,
		"TokenName": tokenName,
	})

	tests := []struct {
		name     string
		preCheck func()
		steps    []resource.TestStep
	}{
		{
			"create and update user token",
			func() {
				err := te.accessClient().CreateUser(context.Background(), &access.UserCreateRequestBody{
					ID:       username,
					Password: gofakeit.Password(true, true, true, true, false, 8),
				})
				require.NoError(t, err)

				t.Cleanup(func() {
					err := te.accessClient().DeleteUser(context.Background(), username)
					require.NoError(t, err)
				})
			},
			[]resource.TestStep{
				{
					Config: te.renderConfig(`resource "proxmox_virtual_environment_user_token" "user_token" {
					comment  			= "Managed by Terraform"
					expiration_date 	= "2034-01-01T22:00:00Z"
					token_name 			= "{{.TokenName}}"
					user_id  			= "{{.Username}}"
				}`),
					Check: testResourceAttributes("proxmox_virtual_environment_user_token.user_token", map[string]string{
						"comment":         "Managed by Terraform",
						"expiration_date": "2034-01-01T22:00:00Z",
						"id":              fmt.Sprintf("%s!%s", username, tokenName),
						"user_id":         username,
						"value":           fmt.Sprintf("%s!%s=.*", username, tokenName),
					}),
				},
				{
					Config: te.renderConfig(`resource "proxmox_virtual_environment_user_token" "user_token" {
					comment  			  = "Managed by Terraform 2"
					expiration_date 	  = "2033-01-01T01:01:01Z"
					privileges_separation = false
					token_name 			  = "{{.TokenName}}"
					user_id  			  = "{{.Username}}"
				}`),
					Check: resource.ComposeTestCheckFunc(
						testResourceAttributes("proxmox_virtual_environment_user_token.user_token", map[string]string{
							"comment":               "Managed by Terraform 2",
							"expiration_date":       "2033-01-01T01:01:01Z",
							"privileges_separation": "false",
							"token_name":            tokenName,
							"user_id":               username,
						}),
						testNoResourceAttributesSet("proxmox_virtual_environment_user_token.user_token", []string{
							"value",
						}),
					),
				},
				{
					ResourceName:      "proxmox_virtual_environment_user_token.user_token",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: te.accProviders,
				PreCheck:                 tt.preCheck,
				Steps:                    tt.steps,
			})
		})
	}
}
