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
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/utils"
)

// pveTicketPattern matches a PVE auth ticket, e.g. "PVE:root@pam:6710ABCD::<signature>".
var pveTicketPattern = regexp.MustCompile(`PVE:[^:\s"]+@[^:\s"]+:[0-9A-F]{8}::`)

var uploadRequestPattern = regexp.MustCompile(`Sending HTTP Request.*multipart/form-data`)

// TestAccProviderDebugLogMasksCredentials verifies that the HTTP request/response logging enabled by TF_LOG=DEBUG does
// not write API credentials to the log, for both token and password authentication.
func TestAccProviderDebugLogMasksCredentials(t *testing.T) {
	te := InitEnvironment(t)

	token := utils.GetAnyStringEnv("PROXMOX_VE_API_TOKEN")
	password := rootPAMPassword()

	_, tokenSecret, _ := strings.Cut(token, "=")
	require.NotEmpty(t, tokenSecret, "PROXMOX_VE_API_TOKEN must be set")

	tests := []struct {
		name       string
		opts       []RenderConfigOption
		authHeader string
		secrets    []string
	}{
		{"api token", nil, "Authorization", []string{tokenSecret}},
		{"password login", []RenderConfigOption{WithRootPassword()}, "Cookie", []string{password}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if slices.Contains(tt.secrets, "") {
				t.Skip("credential not set")
			}

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

// TestAccProviderDebugLogMasksUploadBody verifies that the content of a file uploaded through the API is not written
// to the TF_LOG=DEBUG log.
func TestAccProviderDebugLogMasksUploadBody(t *testing.T) {
	te := InitEnvironment(t)

	content := "password: " + gofakeit.LetterN(16)
	fileName := "test-log-masking-" + gofakeit.LetterN(8) + ".iso"

	te.AddTemplateVars(map[string]any{
		"Content":  content,
		"FileName": fileName,
	})

	logPath := filepath.Join(t.TempDir(), "provider.log")

	t.Setenv("TF_LOG", "DEBUG")
	t.Setenv("TF_LOG_PATH", logPath)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_file" "test" {
					content_type = "iso"
					datastore_id = "{{.DatastoreID}}"
					node_name    = "{{.NodeName}}"
					source_raw {
						data      = "{{.Content}}"
						file_name = "{{.FileName}}"
					}
				}`),
				Check: func(*terraform.State) error {
					contentType := "iso"

					files, err := te.NodeStorageClient().ListDatastoreFiles(t.Context(), &contentType)
					if err != nil {
						return err
					}

					for _, f := range files {
						if f.VolumeID == te.DatastoreID+":iso/"+fileName {
							if f.FileSize != int64(len(content)) {
								return fmt.Errorf("uploaded file has %d bytes, expected %d", f.FileSize, len(content))
							}

							return nil
						}
					}

					return fmt.Errorf("uploaded file %q not found", fileName)
				},
			},
		},
	})

	raw, err := os.ReadFile(logPath)
	require.NoError(t, err)

	log := string(raw)

	require.True(t, uploadRequestPattern.MatchString(log), "upload request was not logged")
	require.False(t, strings.Contains(log, content), "log contains the uploaded file content (%d occurrences)",
		strings.Count(log, content))
}
