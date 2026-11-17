//go:build acceptance || all

//testacc:tier=light
//testacc:resource=storage

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package storage_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/test"
)

func TestAccResourceStoragePBSUpdateInPlace(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)

	tests := []struct {
		name       string
		encryption string
		keyCheck   knownvalue.Check
	}{
		{
			"existing encryption key",
			`encryption_key = jsonencode({ data = "dGVzdA==", fingerprint = "11:22:33:44:55:66:77:88:99:00:11:22:33:44:55:66:77:88:99:00:11:22:33:44:55:66:77:88:99:00:11:22" })`,
			knownvalue.Null(),
		},
		{"generated encryption key", `generate_encryption_key = true`, knownvalue.NotNull()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			name := test.SafeResourceName("pbs-update")

			// disable = true skips the PBS connection check, so no PBS server is needed.
			config := func(extra string) string {
				return te.RenderConfig(fmt.Sprintf(`
					resource "proxmox_storage_pbs" "test" {
						id          = %q
						server      = "192.0.2.10"
						datastore   = "store1"
						username    = "backup@pbs"
						password    = "secret"
						fingerprint = "aa:bb:cc:dd:ee:ff:00:11:22:33:44:55:66:77:88:99:aa:bb:cc:dd:ee:ff:00:11:22:33:44:55:66:77:88:99"
						disable     = true
						%s
						%s
					}`, name, tt.encryption, extra))
			}

			resource.ParallelTest(t, resource.TestCase{
				ProtoV6ProviderFactories: te.AccProviders,
				Steps: []resource.TestStep{
					{
						Config: config(""),
						ConfigStateChecks: []statecheck.StateCheck{
							statecheck.ExpectKnownValue("proxmox_storage_pbs.test", tfjsonpath.New("generated_encryption_key"), tt.keyCheck),
						},
					},
					{
						Config: config(`nodes = ["{{.NodeName}}"]`),
						ConfigStateChecks: []statecheck.StateCheck{
							statecheck.ExpectKnownValue("proxmox_storage_pbs.test", tfjsonpath.New("generated_encryption_key"), tt.keyCheck),
							statecheck.ExpectKnownValue("proxmox_storage_pbs.test", tfjsonpath.New("nodes"),
								knownvalue.SetExact([]knownvalue.Check{knownvalue.StringExact(te.NodeName)})),
						},
					},
				},
			})
		})
	}
}
