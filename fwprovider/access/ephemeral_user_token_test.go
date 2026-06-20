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
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/test"
)

func TestAccEphemeralUserToken_AutoRevoke(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)

	tokenName := test.SafeResourceName("ephtoken")

	te.AddTemplateVars(map[string]any{
		"TokenName": tokenName,
	})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				// Ephemeral resource with auto_revoke = true (default).
				// By the time Check runs, terraform apply has completed and Close() has
				// already deleted the token. The assertion proves revocation happened.
				Config: te.RenderConfig(`
					ephemeral "proxmox_user_token" "test" {
						user_id    = "root@pam"
						token_name = "{{.TokenName}}"
						comment    = "ephemeral test token"
					}

					locals {
						# Reference the ephemeral resource to ensure it is opened during apply.
						_token_ref = ephemeral.proxmox_user_token.test.id
					}
				`, test.WithRootUser()),
				Check: func(*terraform.State) error {
					// Close() was called at end of apply → token must be gone.
					_, err := te.AccessClient().GetUserToken(
						context.Background(), "root@pam", tokenName,
					)
					if err == nil {
						return fmt.Errorf("token %q still exists after apply — Close() did not revoke it", tokenName)
					}
					// Any error from the API means the token is gone; that's the expected outcome.
					return nil
				},
			},
		},
	})
}

func TestAccEphemeralUserToken_NoAutoRevoke(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)

	tokenName := test.SafeResourceName("ephtoken")

	te.AddTemplateVars(map[string]any{
		"TokenName": tokenName,
	})

	// Register cleanup before the steps so interrupted runs don't leave stale tokens.
	t.Cleanup(func() {
		_ = te.AccessClient().DeleteUserToken(context.Background(), "root@pam", tokenName)
	})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				// Ephemeral resource with auto_revoke = false.
				// The token must survive after apply because Close() skips deletion.
				Config: te.RenderConfig(`
					ephemeral "proxmox_user_token" "test" {
						user_id     = "root@pam"
						token_name  = "{{.TokenName}}"
						auto_revoke = false
					}

					locals {
						# Reference the ephemeral resource to ensure it is opened during apply.
						_token_ref = ephemeral.proxmox_user_token.test.id
					}
				`, test.WithRootUser()),
				Check: func(*terraform.State) error {
					// Close() ran but skipped deletion → token must still exist.
					_, err := te.AccessClient().GetUserToken(context.Background(), "root@pam", tokenName)
					if err != nil {
						return fmt.Errorf("token %q was deleted unexpectedly — auto_revoke=false should leave it: %w", tokenName, err)
					}

					return nil
				},
			},
		},
	})
}
