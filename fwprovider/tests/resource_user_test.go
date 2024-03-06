/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package tests

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceUser(t *testing.T) {
	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"create and update user", []resource.TestStep{{
			Config: `
				resource "proxmox_virtual_environment_user" "user1" {
					  comment  			= "Managed by Terraform"
					  email 			= "user1@pve"
					  enabled 			= true
					  expiration_date 	= "2034-01-01T22:00:00Z"
					  first_name 		= "First"
					  last_name 		= "Last"
					  //password 			= "password"
					  user_id  			= "user1@pve"
				}
				`,
			Check: resource.ComposeTestCheckFunc(
				testResourceAttributes("proxmox_virtual_environment_user.user1", map[string]string{
					"comment":         "Managed by Terraform",
					"email":           "user1@pve",
					"enabled":         "true",
					"expiration_date": "2034-01-01T22:00:00Z",
					"first_name":      "First",
					"last_name":       "Last",
					"user_id":         "user1@pve",
				}),
			),
		}, {
			Config: `
				resource "proxmox_virtual_environment_user" "user1" {
					  enabled 			= false
					  expiration_date 	= "2035-01-01T22:00:00Z"
					  user_id  			= "user1@pve"
					  first_name 		= "First One"
				}
				`,
			Check: testResourceAttributes("proxmox_virtual_environment_user.user1", map[string]string{
				"enabled":         "false",
				"expiration_date": "2035-01-01T22:00:00Z",
				"first_name":      "First One",
				"user_id":         "user1@pve",
			}),
		}}},
	}

	accProviders := testAccMuxProviders(context.Background(), t)

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: accProviders,
				Steps:                    tt.steps,
			})
		})
	}
}
