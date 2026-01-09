//go:build acceptance || all

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package access_test

import (
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
			// Create with minimal required fields
			{
				Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_realm_ldap" "test" {
						realm     = "test-realm.local"
						server1   = "ldap.example.com"
						base_dn   = "dc=example,dc=com"
						user_attr = "uid"
						secure    = false
						verify    = false
						comment   = "Test LDAP realm"
					}
				`),
				Check: test.ResourceAttributes("proxmox_virtual_environment_realm_ldap.test", map[string]string{
					"realm":     "test-realm.local",
					"server1":   "ldap.example.com",
					"base_dn":   "dc=example,dc=com",
					"user_attr": "uid",
					"secure":    "false",
					"verify":    "false",
				}),
			},
			// Import state
			{
				ResourceName:            "proxmox_virtual_environment_realm_ldap.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"bind_password"}, // Password not returned by API
			},
			// Update with optional fields added
			{
				Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_realm_ldap" "test" {
						realm     = "test-realm.local"
						server1   = "ldap2.example.com"
						base_dn   = "dc=example,dc=com"
						user_attr = "uid"
						secure    = false
						verify    = false
						comment   = "Test realm with optionals"
						filter    = "(objectClass=person)"
						group_dn  = "ou=groups,dc=example,dc=com"
					}
				`),
				Check: test.ResourceAttributes("proxmox_virtual_environment_realm_ldap.test", map[string]string{
					"realm":    "test-realm.local",
					"server1":  "ldap2.example.com",
					"base_dn":  "dc=example,dc=com",
					"comment":  "Test realm with optionals",
					"filter":   "(objectClass=person)",
					"group_dn": "ou=groups,dc=example,dc=com",
				}),
			},
			// Remove optional fields to verify proper cleanup
			{
				Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_realm_ldap" "test" {
						realm     = "test-realm.local"
						server1   = "ldap2.example.com"
						base_dn   = "dc=example,dc=com"
						user_attr = "uid"
						secure    = false
						verify    = false
						comment   = "Updated test realm"
					}
				`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_virtual_environment_realm_ldap.test", map[string]string{
						"realm":   "test-realm.local",
						"server1": "ldap2.example.com",
						"base_dn": "dc=example,dc=com",
						"comment": "Updated test realm",
					}),
					test.NoResourceAttributesSet("proxmox_virtual_environment_realm_ldap.test", []string{
						"filter",
						"group_dn",
					}),
				),
			},
		},
	})
}
