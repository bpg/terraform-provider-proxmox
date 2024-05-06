/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceUser(t *testing.T) {
	t.Parallel()

	te := initTestEnvironment(t)

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"create and update user", []resource.TestStep{
			{
				Config: `resource "proxmox_virtual_environment_user" "user1" {
					  comment  			= "Managed by Terraform"
					  email 			= "user1@pve"
					  enabled 			= true
					  expiration_date 	= "2034-01-01T22:00:00Z"
					  first_name 		= "First"
					  last_name 		= "Last"
					  user_id  			= "user1@pve"
				}`,
				Check: testResourceAttributes("proxmox_virtual_environment_user.user1", map[string]string{
					"comment":         "Managed by Terraform",
					"email":           "user1@pve",
					"enabled":         "true",
					"expiration_date": "2034-01-01T22:00:00Z",
					"first_name":      "First",
					"last_name":       "Last",
					"user_id":         "user1@pve",
				}),
			},
			{
				Config: `resource "proxmox_virtual_environment_user" "user1" {
					  enabled 			= false
					  expiration_date 	= "2035-01-01T22:00:00Z"
					  user_id  			= "user1@pve"
					  first_name 		= "First One"
				}`,
				Check: testResourceAttributes("proxmox_virtual_environment_user.user1", map[string]string{
					"enabled":         "false",
					"expiration_date": "2035-01-01T22:00:00Z",
					"first_name":      "First One",
					"user_id":         "user1@pve",
				}),
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

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"create and update user token", []resource.TestStep{
			{
				Config: `resource "proxmox_virtual_environment_user" "user1" {
					comment  			= "Managed by Terraform"
					email 			= "user1@pve"
					enabled 			= true
					expiration_date 	= "2034-01-01T22:00:00Z"
					first_name 		= "First"
					last_name 		= "Last"
					user_id  			= "user1@pve"
				}
				resource "proxmox_virtual_environment_user_token" "user1_token" {
					comment  			= "Managed by Terraform"
					expiration_date 	= "2034-01-01T22:00:00Z"
					id 					= "tk1"
					user_id  			= proxmox_virtual_environment_user.user1.user_id
				}
				`,
				Check: testResourceAttributes("proxmox_virtual_environment_user_token.user1_token", map[string]string{
					"comment": "Managed by Terraform",
					"user_id": "user1@pve",
					"value":   `user1@pve!tk1=.*`,
				}),
			},
			//{
			//	Config: `resource "proxmox_virtual_environment_user" "user1" {
			//		  enabled 			= false
			//		  expiration_date 	= "2035-01-01T22:00:00Z"
			//		  user_id  			= "user1@pve"
			//		  first_name 		= "First One"
			//	}`,
			//	Check: testResourceAttributes("proxmox_virtual_environment_user.user1", map[string]string{
			//		"enabled":         "false",
			//		"expiration_date": "2035-01-01T22:00:00Z",
			//		"first_name":      "First One",
			//		"user_id":         "user1@pve",
			//	}),
			//},
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
