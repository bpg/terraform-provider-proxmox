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
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/utils"
)

// pveTicketPattern matches a PVE auth ticket, e.g. "PVE:root@pam:6710ABCD::<signature>".
var pveTicketPattern = regexp.MustCompile(`PVE:[^:\s"]+@[^:\s"]+:[0-9A-F]{8}::`)

// TestAccProviderDebugLogMasksCredentials verifies that the HTTP request/response logging enabled by TF_LOG=DEBUG does
// not write API credentials to the log, for both token and password authentication.
func TestAccProviderDebugLogMasksCredentials(t *testing.T) {
	te := InitEnvironment(t)

	token := utils.GetAnyStringEnv("PROXMOX_VE_API_TOKEN")
	password := utils.GetAnyStringEnv("PROXMOX_VE_PASSWORD")

	_, tokenSecret, _ := strings.Cut(token, "=")
	require.NotEmpty(t, tokenSecret, "PROXMOX_VE_API_TOKEN must be set")
	require.NotEmpty(t, password, "PROXMOX_VE_PASSWORD must be set")

	tests := []struct {
		name       string
		opts       []RenderConfigOption
		authHeader string
		secrets    []string
	}{
		{"api token", nil, "Authorization", []string{tokenSecret}},
		{"password login", []RenderConfigOption{WithRootUser()}, "Cookie", []string{password}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logPath := filepath.Join(t.TempDir(), "provider.log")

			t.Setenv("TF_LOG", "DEBUG")
			t.Setenv("TF_LOG_PATH", logPath)

			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: te.AccProviders,
				Steps: []resource.TestStep{
					{
						Config: te.RenderConfig(`data "proxmox_virtual_environment_version" "test" {}`, tt.opts...),
						Check: resource.TestCheckResourceAttrSet(
							"data.proxmox_virtual_environment_version.test", "version",
						),
					},
				},
			})

			raw, err := os.ReadFile(logPath)
			require.NoError(t, err)

			log := string(raw)

			// Boolean checks only: testify's (Not)Contains prints the whole log, secrets included, on failure.
			require.True(t, strings.Contains(log, "Sending HTTP Request"), "HTTP requests were not logged")
			require.True(t, strings.Contains(log, tt.authHeader+"="), "auth header field was not logged")

			for _, s := range tt.secrets {
				require.False(t, strings.Contains(log, s), "log contains a credential (%d occurrences)", strings.Count(log, s))
			}

			require.False(t, pveTicketPattern.MatchString(log), "log contains a PVE auth ticket")
		})
	}
}
