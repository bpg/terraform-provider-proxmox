//go:build acceptance || all

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package access_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/test"
	"github.com/bpg/terraform-provider-proxmox/proxmox/access"
)

func TestAccAcl_User(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)

	userID := fmt.Sprintf("%s@pve", gofakeit.Username())
	te.AddTemplateVars(map[string]any{
		"UserID": userID,
	})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		CheckDestroy:             nil,
		PreCheck: func() {
			err := te.AccessClient().CreateUser(context.Background(), &access.UserCreateRequestBody{
				ID:       userID,
				Password: gofakeit.Password(true, true, true, true, false, 8),
			})
			require.NoError(t, err)

			t.Cleanup(func() {
				err := te.AccessClient().DeleteUser(context.Background(), userID)
				require.NoError(t, err)
			})
		},
		Steps: []resource.TestStep{
			{
				Config: te.RenderConfig(`resource "proxmox_virtual_environment_acl" "test" {
					user_id = "{{.UserID}}"
					path = "/"
					role_id = "NoAccess"
				}`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_virtual_environment_acl.test", map[string]string{
						"path":      "/",
						"role_id":   "NoAccess",
						"user_id":   userID,
						"propagate": "true",
					}),
					test.NoResourceAttributesSet("proxmox_virtual_environment_acl.test", []string{
						"group_id",
						"token_id",
					}),
				),
			},
			{
				ResourceName:      "proxmox_virtual_environment_acl.test",
				ImportState:       true,
				ImportStateIdFunc: testAccACLImportStateIDFunc(),
				ImportStateVerify: true,
			},
			{
				Config: te.RenderConfig(`resource "proxmox_virtual_environment_acl" "test" {
					user_id = "{{.UserID}}"
					path = "/"
					role_id = "PVEPoolUser"
				}`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_virtual_environment_acl.test", map[string]string{
						"path":      "/",
						"role_id":   "PVEPoolUser",
						"user_id":   userID,
						"propagate": "true",
					}),
					test.NoResourceAttributesSet("proxmox_virtual_environment_acl.test", []string{
						"group_id",
						"token_id",
					}),
				),
			},
		},
	})
}

func TestAccAcl_Validators(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				PlanOnly: true,
				Config: `resource "proxmox_virtual_environment_acl" "test" {
					group_id = "test"
					path = "/"
					role_id = "test"
					token_id = "test"
				}`,
				ExpectError: regexp.MustCompile(`.*Error: Invalid Attribute Combination`),
			},
			{
				PlanOnly: true,
				Config: `resource "proxmox_virtual_environment_acl" "test" {
					path = "/"
					role_id = "test"
					token_id = "test"
					user_id = "test"
				}`,
				ExpectError: regexp.MustCompile(`.*Error: Invalid Attribute Combination`),
			},
			{
				PlanOnly: true,
				Config: `resource "proxmox_virtual_environment_acl" "test" {
					group_id = "test"
					path = "/"
					role_id = "test"
					user_id = "test"
				}`,
				ExpectError: regexp.MustCompile(`.*Error: Invalid Attribute Combination`),
			},
			{
				PlanOnly: true,
				Config: `resource "proxmox_virtual_environment_acl" "test" {
					group_id = "test"
					path = "/"
					role_id = "test"
					token_id = "test"
					user_id = "test"
				}`,
				ExpectError: regexp.MustCompile(`.*Error: Invalid Attribute Combination`),
			},
		},
	})
}

func testAccACLImportStateIDFunc() resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		resourceName := "proxmox_virtual_environment_acl.test"

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		path := rs.Primary.Attributes["path"]

		groupID := rs.Primary.Attributes["group_id"]
		tokenID := rs.Primary.Attributes["token_id"]
		userID := rs.Primary.Attributes["user_id"]
		entityID := groupID + tokenID + userID

		roleID := rs.Primary.Attributes["role_id"]

		return path + "?" + entityID + "?" + roleID, nil
	}
}
