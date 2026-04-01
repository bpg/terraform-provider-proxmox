//go:build acceptance || all

//testacc:tier=light
//testacc:resource=access

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

func TestAccRealmOpenIDUsernameClaim(t *testing.T) {
	te := test.InitEnvironment(t)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			// Create with username_claim = "upn" (custom claim used by ADFS/Azure AD)
			{
				Config: te.RenderConfig(`
					resource "proxmox_realm_openid" "test_upn" {
						realm          = "test-upn"
						issuer_url     = "https://accounts.google.com"
						client_id      = "test-client-id"
						username_claim = "upn"
					}
				`),
				Check: test.ResourceAttributes("proxmox_realm_openid.test_upn", map[string]string{
					"realm":          "test-upn",
					"username_claim": "upn",
				}),
			},
			// Import state
			{
				ResourceName:            "proxmox_realm_openid.test_upn",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"client_key"},
			},
		},
	})
}

func TestAccRealmOpenID(t *testing.T) {
	te := test.InitEnvironment(t)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			// Step 1: Create with minimal required fields, verify computed defaults
			{
				Config: te.RenderConfig(`
					resource "proxmox_realm_openid" "test" {
						realm      = "test-oidc"
						issuer_url = "https://accounts.google.com"
						client_id  = "test-client-id"
						comment    = "Test OpenID realm"
					}
				`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_realm_openid.test", map[string]string{
						"realm":      "test-oidc",
						"issuer_url": "https://accounts.google.com",
						"client_id":  "test-client-id",
						"comment":    "Test OpenID realm",
						// Verify computed defaults from PVE API.
						// Note: PVE does not return query_userinfo or scopes for a
						// minimal realm — they are null/unset, not "true"/"email profile"
						// as previously asserted (corrected from PR #2655).
						"autocreate":        "false",
						"default":           "false",
						"groups_autocreate": "false",
						"groups_overwrite":  "false",
						"query_userinfo":    "false",
					}),
					test.NoResourceAttributesSet("proxmox_realm_openid.test", []string{
						"username_claim",
						"groups_claim",
						"prompt",
						"acr_values",
						"client_key",
						"scopes",
					}),
				),
			},
			// Step 2: Import state
			{
				ResourceName:            "proxmox_realm_openid.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"client_key"}, // Not returned by API
			},
			// Step 3: Update with optional fields added
			{
				Config: te.RenderConfig(`
					resource "proxmox_realm_openid" "test" {
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
				Check: test.ResourceAttributes("proxmox_realm_openid.test", map[string]string{
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
					resource "proxmox_realm_openid" "test" {
						realm      = "test-oidc"
						issuer_url = "https://accounts.google.com"
						client_id  = "test-client-id"
						comment    = "Cleaned OpenID realm"
					}
				`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_realm_openid.test", map[string]string{
						"realm":      "test-oidc",
						"issuer_url": "https://accounts.google.com",
						"client_id":  "test-client-id",
						"comment":    "Cleaned OpenID realm",
						// Verify defaults restored after cleanup
						"autocreate":        "false",
						"groups_autocreate": "false",
						"groups_overwrite":  "false",
						"query_userinfo":    "false",
					}),
					test.NoResourceAttributesSet("proxmox_realm_openid.test", []string{
						"groups_claim",
						"prompt",
						"acr_values",
					}),
				),
			},
		},
	})
}
