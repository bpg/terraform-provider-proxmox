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

func TestAccRealmOpenID(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			// Step 1: Create with minimal required fields, verify computed defaults
			{
				Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_realm_openid" "test" {
						realm      = "test-oidc"
						issuer_url = "https://accounts.google.com"
						client_id  = "test-client-id"
						comment    = "Test OpenID realm"
					}
				`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_virtual_environment_realm_openid.test", map[string]string{
						"realm":      "test-oidc",
						"issuer_url": "https://accounts.google.com",
						"client_id":  "test-client-id",
						"comment":    "Test OpenID realm",
						// Verify computed defaults
						"autocreate":        "false",
						"default":           "false",
						"groups_autocreate": "false",
						"groups_overwrite":  "false",
						"query_userinfo":    "true",
						"scopes":            "email profile",
					}),
					test.NoResourceAttributesSet("proxmox_virtual_environment_realm_openid.test", []string{
						"username_claim",
						"groups_claim",
						"prompt",
						"acr_values",
						"client_key",
					}),
				),
			},
			// Step 2: Import state
			{
				ResourceName:            "proxmox_virtual_environment_realm_openid.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"client_key"}, // Not returned by API
			},
			// Step 3: Update with optional fields added
			{
				Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_realm_openid" "test" {
						realm             = "test-oidc"
						issuer_url        = "https://accounts.google.com"
						client_id         = "test-client-id"
						autocreate        = true
						groups_claim      = "groups"
						groups_autocreate = true
						groups_overwrite  = true
						scopes            = "openid email profile"
						query_userinfo    = false
						comment           = "Updated OpenID realm"
					}
				`),
				Check: test.ResourceAttributes("proxmox_virtual_environment_realm_openid.test", map[string]string{
					"realm":             "test-oidc",
					"issuer_url":        "https://accounts.google.com",
					"client_id":         "test-client-id",
					"autocreate":        "true",
					"groups_claim":      "groups",
					"groups_autocreate": "true",
					"groups_overwrite":  "true",
					"scopes":            "openid email profile",
					"query_userinfo":    "false",
					"comment":           "Updated OpenID realm",
				}),
			},
			// Step 4: Remove optional fields to verify proper cleanup via delete parameter
			{
				Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_realm_openid" "test" {
						realm      = "test-oidc"
						issuer_url = "https://accounts.google.com"
						client_id  = "test-client-id"
						comment    = "Cleaned OpenID realm"
					}
				`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_virtual_environment_realm_openid.test", map[string]string{
						"realm":      "test-oidc",
						"issuer_url": "https://accounts.google.com",
						"client_id":  "test-client-id",
						"comment":    "Cleaned OpenID realm",
						// Verify defaults restored after cleanup
						"autocreate":        "false",
						"groups_autocreate": "false",
						"groups_overwrite":  "false",
						"query_userinfo":    "true",
						"scopes":            "email profile",
					}),
					test.NoResourceAttributesSet("proxmox_virtual_environment_realm_openid.test", []string{
						"groups_claim",
						"prompt",
						"acr_values",
					}),
				),
			},
		},
	})
}
