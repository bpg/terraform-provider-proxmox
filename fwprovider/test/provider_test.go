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
	"net/url"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/utils"
)

func TestAccProviderSSHNodeAddressSource(t *testing.T) {
	te := InitEnvironment(t)

	nodeName := utils.GetAnyStringEnv("PROXMOX_VE_ACC_NODE_NAME")
	if nodeName == "" {
		nodeName = "pve"
	}

	nodeAddress := utils.GetAnyStringEnv("PROXMOX_VE_ACC_NODE_SSH_ADDRESS")
	if nodeAddress == "" {
		endpoint := utils.GetAnyStringEnv("PROXMOX_VE_ENDPOINT")

		u, err := url.Parse(endpoint)
		require.NoError(t, err)

		nodeAddress = u.Hostname()
	}

	nodePort := utils.GetAnyStringEnv("PROXMOX_VE_ACC_NODE_SSH_PORT")
	if nodePort == "" {
		nodePort = "22"
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				// Verify attribute is accepted alongside node overrides
				Config: fmt.Sprintf(`
					provider "proxmox" {
						ssh {
							node_address_source = "dns"
							node {
								name    = %q
								address = %q
								port    = %s
							}
						}
					}

					data "proxmox_virtual_environment_version" "test" {}
				`, nodeName, nodeAddress, nodePort),
				Check: resource.TestCheckResourceAttrSet("data.proxmox_virtual_environment_version.test", "version"),
			},
			{
				// Verify DNS resolution actually works without node overrides
				Config: `
					provider "proxmox" {
						ssh {
							node_address_source = "dns"
						}
					}

					data "proxmox_virtual_environment_version" "test" {}
				`,
				Check: resource.TestCheckResourceAttrSet("data.proxmox_virtual_environment_version.test", "version"),
			},
		},
	})
}

func TestAccProviderCloudflareAccess(t *testing.T) {
	te := InitEnvironment(t)

	const versionDataSource = `
		data "proxmox_virtual_environment_version" "test" {}
	`

	versionRead := resource.TestCheckResourceAttrSet("data.proxmox_virtual_environment_version.test", "version")

	// The PVE test endpoint ignores the CF-Access-* headers, so every complete configuration must still read the
	// version data source. Environment variables are set per scenario, hence resource.Test instead of ParallelTest.
	tests := []struct {
		name string
		env  map[string]string
		step resource.TestStep
	}{
		{
			name: "both set in block",
			step: resource.TestStep{
				Config: `
					provider "proxmox" {
						cloudflare_access {
							client_id     = "test-client-id.access"
							client_secret = "test-client-secret"
						}
					}` + versionDataSource,
				Check: versionRead,
			},
		},
		{
			name: "empty block disables cloudflare access",
			step: resource.TestStep{
				Config: `
					provider "proxmox" {
						cloudflare_access {}
					}` + versionDataSource,
				Check: versionRead,
			},
		},
		{
			name: "only client_id in block",
			step: resource.TestStep{
				Config: `
					provider "proxmox" {
						cloudflare_access {
							client_id = "test-client-id.access"
						}
					}` + versionDataSource,
				ExpectError: regexp.MustCompile(`(?s)requires both client_id and\s+client_secret`),
			},
		},
		{
			name: "client_id in block, secret from env",
			env:  map[string]string{"PROXMOX_VE_CF_ACCESS_CLIENT_SECRET": "test-client-secret"},
			step: resource.TestStep{
				Config: `
					provider "proxmox" {
						cloudflare_access {
							client_id = "test-client-id.access"
						}
					}` + versionDataSource,
				Check: versionRead,
			},
		},
		{
			name: "both from env",
			env: map[string]string{
				"PROXMOX_VE_CF_ACCESS_CLIENT_ID":     "test-client-id.access",
				"PROXMOX_VE_CF_ACCESS_CLIENT_SECRET": "test-client-secret",
			},
			step: resource.TestStep{
				Config: `
					provider "proxmox" {}` + versionDataSource,
				Check: versionRead,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.env {
				t.Setenv(k, v)
			}

			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: te.AccProviders,
				Steps:                    []resource.TestStep{tt.step},
			})
		})
	}
}
