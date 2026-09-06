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
	"fmt"
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

// TestAccProviderAPIHeadersFromEnv verifies that both providers read PROXMOX_VE_API_HEADERS when the
// attribute is absent: a custom header must not disturb the request and a reserved one must be rejected.
func TestAccProviderAPIHeadersFromEnv(t *testing.T) {
	te := InitEnvironment(t)

	tests := []struct {
		name  string
		value string
		step  resource.TestStep
	}{
		{
			name:  "custom header",
			value: "X-Acc-Test=1",
			step: resource.TestStep{
				Config: te.RenderConfig(`data "proxmox_virtual_environment_version" "test" {}`),
				Check: resource.TestCheckResourceAttrSet(
					"data.proxmox_virtual_environment_version.test", "version",
				),
			},
		},
		{
			name:  "reserved header",
			value: "Authorization=Bearer nope",
			step: resource.TestStep{
				Config:      te.RenderConfig(`data "proxmox_virtual_environment_version" "test" {}`),
				ExpectError: regexp.MustCompile(`managed by the provider`),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("PROXMOX_VE_API_HEADERS", tt.value)

			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: te.AccProviders,
				Steps:                    []resource.TestStep{tt.step},
			})
		})
	}
}

// TestAccProviderAPIHeadersUnknownOrNullValue verifies that a header value that is unknown at plan time,
// or null, fails configuration with a diagnostic naming the attribute rather than a generic conversion
// error.
func TestAccProviderAPIHeadersUnknownOrNullValue(t *testing.T) {
	te := InitEnvironment(t)

	tests := []struct {
		name    string
		headers string
		wantErr *regexp.Regexp
	}{
		{
			name:    "unknown value",
			headers: `{ "X-Acc-Test" = terraform_data.header.output }`,
			wantErr: regexp.MustCompile(`(?s)Unknown Proxmox VE API\s+Headers`),
		},
		{
			name:    "null value",
			headers: `{ "X-Acc-Test" = null }`,
			wantErr: regexp.MustCompile(`(?s)"X-Acc-Test"\s+is\s+null`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource.ParallelTest(t, resource.TestCase{
				ProtoV6ProviderFactories: te.AccProviders,
				Steps: []resource.TestStep{
					{
						Config: fmt.Sprintf(`
							resource "terraform_data" "header" {
								input = "1"
							}

							provider "proxmox" {
								api_headers = %s
							}

							data "proxmox_virtual_environment_version" "test" {}
						`, tt.headers),
						ExpectError: tt.wantErr,
					},
				},
			})
		})
	}
}
