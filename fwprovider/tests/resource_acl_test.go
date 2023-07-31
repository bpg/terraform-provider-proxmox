/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package tests

import (
	"context"
	"fmt"
	"net/url"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccAcl_User(t *testing.T) {
	resourceName := "proxmox_virtual_environment_acl.test"

	accProviders := testAccMuxProviders(context.Background(), t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: accProviders,
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigACLUser("provider_test@pve", "NoAccess"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "path", "/"),
					resource.TestCheckResourceAttr(resourceName, "role_id", "NoAccess"),
					resource.TestCheckResourceAttr(resourceName, "user_id", "provider_test@pve"),
					resource.TestCheckResourceAttr(resourceName, "propagate", "true"),
					resource.TestCheckNoResourceAttr(resourceName, "group_id"),
					resource.TestCheckNoResourceAttr(resourceName, "token_id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccACLImportStateIDFunc(),
				ImportStateVerify: true,
			},
			{
				Config: testAccConfigACLUser("provider_test@pve", "PVEPoolUser"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "path", "/"),
					resource.TestCheckResourceAttr(resourceName, "role_id", "PVEPoolUser"),
					resource.TestCheckResourceAttr(resourceName, "user_id", "provider_test@pve"),
					resource.TestCheckResourceAttr(resourceName, "propagate", "true"),
					resource.TestCheckNoResourceAttr(resourceName, "group_id"),
					resource.TestCheckNoResourceAttr(resourceName, "token_id"),
				),
			},
		},
	})
}

func TestAccAcl_Validators(t *testing.T) {
	t.Parallel()

	accProviders := testAccMuxProviders(context.Background(), t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: accProviders,
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				PlanOnly:    true,
				Config:      testAccConfigACLValidators("test", "test", ""),
				ExpectError: regexp.MustCompile(`.*Error: Invalid Attribute Combination`),
			},
			{
				PlanOnly:    true,
				Config:      testAccConfigACLValidators("", "test", "test"),
				ExpectError: regexp.MustCompile(`.*Error: Invalid Attribute Combination`),
			},
			{
				PlanOnly:    true,
				Config:      testAccConfigACLValidators("test", "", "test"),
				ExpectError: regexp.MustCompile(`.*Error: Invalid Attribute Combination`),
			},
			{
				PlanOnly:    true,
				Config:      testAccConfigACLValidators("test", "test", "test"),
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

		roleID := rs.Primary.Attributes["role_id"]

		v := url.Values{
			"entity_id": []string{groupID + tokenID + userID},
			"role_id":   []string{roleID},
		}

		return path + "?" + v.Encode(), nil
	}
}

func testAccConfigACLUser(userID string, roleID string) string {
	return fmt.Sprintf(`
resource "proxmox_virtual_environment_acl" "test" {
  user_id = %q
  path = "/"
  role_id = %q
}
`, userID, roleID)
}

func testAccConfigACLValidators(groupID string, tokenID string, userID string) string {
	var groupAttr string
	if groupID != "" {
		groupAttr = fmt.Sprintf("\n  group_id = %q", groupID)
	}

	var tokenAttr string
	if tokenID != "" {
		tokenAttr = fmt.Sprintf("\n  token_id = %q", tokenID)
	}

	var userAttr string
	if userID != "" {
		userAttr = fmt.Sprintf("\n  user_id = %q", userID)
	}

	return fmt.Sprintf(`
resource "proxmox_virtual_environment_acl" "test" {%v%v%v
  path = "/"
  role_id = "test"
}
`, groupAttr, tokenAttr, userAttr)
}
