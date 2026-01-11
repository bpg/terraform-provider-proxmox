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

func TestAccRealmSync(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: te.RenderConfig(testAccRealmLDAPWithSyncConfig("test-realm-sync.local", "ldap.example.com", "dc=example,dc=com")),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("proxmox_virtual_environment_realm_sync.sync", "realm", "test-realm-sync.local"),
					resource.TestCheckResourceAttr("proxmox_virtual_environment_realm_sync.sync", "scope", "users"),
				),
			},
		},
	})
}

func testAccRealmLDAPWithSyncConfig(realm, server, baseDN string) string {
	return fmt.Sprintf(`
resource "proxmox_virtual_environment_realm_ldap" "test" {
  realm    = "%s"
  server1  = "%s"
  base_dn  = "%s"
  user_attr = "uid"
  comment  = "Test LDAP realm for sync"
}

resource "proxmox_virtual_environment_realm_sync" "sync" {
  realm = proxmox_virtual_environment_realm_ldap.test.realm
  scope = "users"
}
`, realm, server, baseDN)
}
