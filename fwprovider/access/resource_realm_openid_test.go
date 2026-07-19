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
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"

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

func TestAccRealmOpenIDWriteOnly(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)
	realm := test.SafeResourceName("oidc-wo")
	te.AddTemplateVars(map[string]interface{}{
		"Realm": realm,
	})

	tests := []struct {
		name string
		step []resource.TestStep
	}{
		{"client_key_wo keeps the secret out of state", []resource.TestStep{
			{
				Config: te.RenderConfig(`
					resource "proxmox_realm_openid" "test_wo" {
						realm                 = "{{.Realm}}"
						issuer_url            = "https://accounts.google.com"
						client_id             = "test-client-id"
						client_key_wo         = "super-secret-value"
						client_key_wo_version = 1
					}`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_realm_openid.test_wo", map[string]string{
						"realm":                 realm,
						"client_id":             "test-client-id",
						"client_key_wo_version": "1",
					}),
					// The write-only secret must never be persisted, and the legacy
					// client_key mirror must stay unset when client_key_wo is used.
					resource.TestCheckNoResourceAttr("proxmox_realm_openid.test_wo", "client_key_wo"),
					resource.TestCheckNoResourceAttr("proxmox_realm_openid.test_wo", "client_key"),
					// Behavioral proof: the realm was actually created on the server.
					testCheckOpenIDRealmExists(te, realm),
				),
			},
		}},
		{"client_key_wo_version triggers rotation", []resource.TestStep{
			{
				Config: te.RenderConfig(`
					resource "proxmox_realm_openid" "test_rotate" {
						realm                 = "{{.Realm}}-rot"
						issuer_url            = "https://accounts.google.com"
						client_id             = "test-client-id"
						client_key_wo         = "old-secret"
						client_key_wo_version = 1
					}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("proxmox_realm_openid.test_rotate", "client_key_wo_version", "1"),
					resource.TestCheckNoResourceAttr("proxmox_realm_openid.test_rotate", "client_key_wo"),
				),
			},
			{
				// Rotate the secret: a new client_key_wo with a bumped version. The
				// version bump is what produces a diff, since write-only values are
				// invisible to Terraform. Apply must succeed and re-send the secret.
				Config: te.RenderConfig(`
					resource "proxmox_realm_openid" "test_rotate" {
						realm                 = "{{.Realm}}-rot"
						issuer_url            = "https://accounts.google.com"
						client_id             = "test-client-id"
						client_key_wo         = "new-secret"
						client_key_wo_version = 2
					}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("proxmox_realm_openid.test_rotate", "client_key_wo_version", "2"),
					resource.TestCheckNoResourceAttr("proxmox_realm_openid.test_rotate", "client_key_wo"),
				),
			},
		}},
		{"removing client_key_wo deletes the key from PVE", []resource.TestStep{
			{
				Config: te.RenderConfig(`
					resource "proxmox_realm_openid" "test_remove" {
						realm                 = "{{.Realm}}-rm"
						issuer_url            = "https://accounts.google.com"
						client_id             = "test-client-id"
						client_key_wo         = "remove-me"
						client_key_wo_version = 1
					}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("proxmox_realm_openid.test_remove", "client_key_wo_version", "1"),
					testCheckOpenIDRealmClientKey(te, fmt.Sprintf("%s-rm", realm), true),
				),
			},
			{
				Config: te.RenderConfig(`
					resource "proxmox_realm_openid" "test_remove" {
						realm      = "{{.Realm}}-rm"
						issuer_url = "https://accounts.google.com"
						client_id  = "test-client-id"
					}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("proxmox_realm_openid.test_remove", "client_key_wo_version"),
					testCheckOpenIDRealmClientKey(te, fmt.Sprintf("%s-rm", realm), false),
				),
			},
		}},
		{"migrating from client_key to client_key_wo clears it from state", []resource.TestStep{
			{
				Config: te.RenderConfig(`
					resource "proxmox_realm_openid" "test_migrate" {
						realm      = "{{.Realm}}-mig"
						issuer_url = "https://accounts.google.com"
						client_id  = "test-client-id"
						client_key = "legacy-secret"
					}`),
				Check: resource.TestCheckResourceAttr("proxmox_realm_openid.test_migrate", "client_key", "legacy-secret"),
			},
			{
				Config: te.RenderConfig(`
					resource "proxmox_realm_openid" "test_migrate" {
						realm                 = "{{.Realm}}-mig"
						issuer_url            = "https://accounts.google.com"
						client_id             = "test-client-id"
						client_key_wo         = "migrated-secret"
						client_key_wo_version = 1
					}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("proxmox_realm_openid.test_migrate", "client_key_wo_version", "1"),
					resource.TestCheckNoResourceAttr("proxmox_realm_openid.test_migrate", "client_key"),
					resource.TestCheckNoResourceAttr("proxmox_realm_openid.test_migrate", "client_key_wo"),
				),
			},
		}},
		{"client_key_wo_version requires client_key_wo", []resource.TestStep{
			{
				Config: te.RenderConfig(`
					resource "proxmox_realm_openid" "test_noversion" {
						realm                 = "{{.Realm}}-nov"
						issuer_url            = "https://accounts.google.com"
						client_id             = "test-client-id"
						client_key_wo_version = 1
					}`),
				ExpectError: regexp.MustCompile(`Attribute "client_key_wo" must be specified when "client_key_wo_version" is\s+specified`),
			},
		}},
		{"client_key and client_key_wo are mutually exclusive", []resource.TestStep{
			{
				Config: te.RenderConfig(`
					resource "proxmox_realm_openid" "test_conflict" {
						realm         = "{{.Realm}}-cnf"
						issuer_url    = "https://accounts.google.com"
						client_id     = "test-client-id"
						client_key    = "secret-a"
						client_key_wo = "secret-b"
					}`),
				ExpectError: regexp.MustCompile(`These attributes cannot be configured together: \[client_key,client_key_wo\]`),
			},
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource.ParallelTest(t, resource.TestCase{
				// Write-only attributes require Terraform 1.11+.
				TerraformVersionChecks: []tfversion.TerraformVersionCheck{
					tfversion.SkipBelow(tfversion.Version1_11_0),
				},
				ProtoV6ProviderFactories: te.AccProviders,
				Steps:                    tt.step,
			})
		})
	}
}

// testCheckOpenIDRealmExists verifies, via a direct API read, that the realm exists
// on the server and is an OpenID realm. Used to prove a realm configured solely with
// the write-only client_key_wo is actually created, even though the secret never
// lands in Terraform state.
func testCheckOpenIDRealmExists(te *test.Environment, realm string) resource.TestCheckFunc {
	return func(*terraform.State) error {
		data, err := te.AccessClient().GetRealm(context.Background(), realm)
		if err != nil {
			return fmt.Errorf("reading OpenID realm %q: %w", realm, err)
		}

		if data.Type != "openid" {
			return fmt.Errorf("realm %q has type %q, want %q", realm, data.Type, "openid")
		}

		return nil
	}
}

// testCheckOpenIDRealmClientKey verifies, via a direct API read, whether the realm has a
// client-key configured. Used to prove that removing client_key_wo sends delete=client-key.
func testCheckOpenIDRealmClientKey(te *test.Environment, realm string, want bool) resource.TestCheckFunc {
	return func(*terraform.State) error {
		data, err := te.AccessClient().GetRealm(context.Background(), realm)
		if err != nil {
			return fmt.Errorf("reading OpenID realm %q: %w", realm, err)
		}

		has := data.ClientKey != nil && *data.ClientKey != ""
		if has != want {
			if want {
				return fmt.Errorf("realm %q has no client-key; write-only client_key_wo was not stored", realm)
			}

			return fmt.Errorf("client-key still present on PVE after removing client_key_wo")
		}

		return nil
	}
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
						"audiences",
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
						audiences         = "1234567890"
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
					"audiences":         "1234567890",
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
						"audiences",
					}),
				),
			},
		},
	})
}
