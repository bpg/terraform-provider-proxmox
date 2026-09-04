//go:build acceptance || all

//testacc:tier=light
//testacc:resource=misc

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccProviderAPIHeaders verifies that user-defined headers do not disturb a real API request.
// That the headers actually reach the wire is covered by the unit tests in proxmox/api, which can
// assert on the received request, and by the mitmproxy workflow.
func TestAccProviderAPIHeaders(t *testing.T) {
	te := InitEnvironment(t)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: te.RenderConfig(
					`data "proxmox_virtual_environment_version" "test" {}`,
					WithAPIHeaders(map[string]string{"X-Acc-Test": "1"}),
				),
				Check: resource.TestCheckResourceAttrSet(
					"data.proxmox_virtual_environment_version.test", "version",
				),
			},
		},
	})
}

// TestAccProviderAPIHeadersReserved verifies that a header managed by the provider is rejected
// during configuration rather than silently dropped or overriding the authentication.
func TestAccProviderAPIHeadersReserved(t *testing.T) {
	te := InitEnvironment(t)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: te.RenderConfig(
					`data "proxmox_virtual_environment_version" "test" {}`,
					WithAPIHeaders(map[string]string{"Authorization": "Bearer nope"}),
				),
				ExpectError: regexp.MustCompile(`managed by the provider`),
			},
		},
	})
}

// TestAccProviderAPIHeadersEmptyOverridesEnv verifies that an explicit empty map wins over the
// environment variable. A reserved header in the variable would fail configuration, so a successful
// read proves the configuration took precedence.
func TestAccProviderAPIHeadersEmptyOverridesEnv(t *testing.T) {
	te := InitEnvironment(t)

	t.Setenv("PROXMOX_VE_API_HEADERS", "Authorization=Bearer nope")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: te.RenderConfig(
					`data "proxmox_virtual_environment_version" "test" {}`,
					WithAPIHeaders(map[string]string{}),
				),
				Check: resource.TestCheckResourceAttrSet(
					"data.proxmox_virtual_environment_version.test", "version",
				),
			},
		},
	})
}
