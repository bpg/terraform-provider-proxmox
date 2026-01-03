//go:build acceptance || all

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package access_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/test"
)

func TestAccRealmLDAP(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: te.RenderConfig(testAccRealmLDAPConfig("test-realm.local", "ldap.example.com", "dc=example,dc=com")),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("proxmox_virtual_environment_realm_ldap.test", "realm", "test-realm.local"),
					resource.TestCheckResourceAttr("proxmox_virtual_environment_realm_ldap.test", "server1", "ldap.example.com"),
					resource.TestCheckResourceAttr("proxmox_virtual_environment_realm_ldap.test", "base_dn", "dc=example,dc=com"),
					resource.TestCheckResourceAttr("proxmox_virtual_environment_realm_ldap.test", "user_attr", "uid"),
					resource.TestCheckResourceAttr("proxmox_virtual_environment_realm_ldap.test", "secure", "false"),
					resource.TestCheckResourceAttr("proxmox_virtual_environment_realm_ldap.test", "verify", "false"),
				),
			},
			{
				ResourceName:            "proxmox_virtual_environment_realm_ldap.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"bind_password"}, // Password not returned by API
			},
			{
				Config: te.RenderConfig(testAccRealmLDAPConfigUpdate("test-realm.local", "ldap2.example.com", "dc=example,dc=com")),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("proxmox_virtual_environment_realm_ldap.test", "realm", "test-realm.local"),
					resource.TestCheckResourceAttr("proxmox_virtual_environment_realm_ldap.test", "server1", "ldap2.example.com"),
					resource.TestCheckResourceAttr("proxmox_virtual_environment_realm_ldap.test", "base_dn", "dc=example,dc=com"),
					resource.TestCheckResourceAttr("proxmox_virtual_environment_realm_ldap.test", "comment", "Updated test realm"),
				),
			},
		},
	})
}

func testAccRealmLDAPConfig(realm, server, baseDN string) string {
	return fmt.Sprintf(`
resource "proxmox_virtual_environment_realm_ldap" "test" {
  realm    = "%s"
  server1  = "%s"
  base_dn  = "%s"
  user_attr = "uid"
  secure   = false
  verify   = false
  comment  = "Test LDAP realm created by Terraform"
}
`, realm, server, baseDN)
}

func testAccRealmLDAPConfigUpdate(realm, server, baseDN string) string {
	return fmt.Sprintf(`
resource "proxmox_virtual_environment_realm_ldap" "test" {
  realm    = "%s"
  server1  = "%s"
  base_dn  = "%s"
  user_attr = "uid"
  secure   = false
  verify   = false
  comment  = "Updated test realm"
}
`, realm, server, baseDN)
}
