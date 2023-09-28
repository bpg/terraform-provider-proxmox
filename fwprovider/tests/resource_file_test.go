/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package tests

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/ssh"
	"github.com/bpg/terraform-provider-proxmox/utils"
)

const (
	accTestFileName = "proxmox_virtual_environment_file.test"
)

type nodeResolver struct {
	node ssh.ProxmoxNode
}

func (c *nodeResolver) Resolve(_ context.Context, _ string) (ssh.ProxmoxNode, error) {
	return c.node, nil
}

func TestAccResourceFile(t *testing.T) {
	t.Parallel()

	accProviders := testAccMuxProviders(context.Background(), t)

	snippetRaw := fmt.Sprintf("snippet-raw-%s.txt", gofakeit.Word())
	snippetURL := "https://raw.githubusercontent.com/yaml/yaml-test-suite/main/src/229Q.yaml"
	snippetFile1 := createFile(t, "snippet-file-1-*.yaml", "test snippet 1 - file")
	snippetFile2 := createFile(t, "snippet-file-2-*.yaml", "test snippet 2 - file")
	fileISO := createFile(t, "file-*.iso", "pretend it is an ISO")

	endpoint := utils.GetAnyStringEnv("PROXMOX_VE_ENDPOINT")
	u, err := url.ParseRequestURI(endpoint)
	require.NoError(t, err)

	sshUsername := strings.Split(utils.GetAnyStringEnv("PROXMOX_VE_USERNAME"), "@")[0]
	sshAgentSocket := utils.GetAnyStringEnv("SSH_AUTH_SOCK", "PROXMOX_VE_SSH_AUTH_SOCK", "PM_VE_SSH_AUTH_SOCK")

	sshClient, err := ssh.NewClient(
		sshUsername, "", true, sshAgentSocket,
		&nodeResolver{
			node: ssh.ProxmoxNode{
				Address: u.Hostname(),
				Port:    22,
			},
		},
	)
	require.NoError(t, err)

	f, err := os.Open(snippetFile2.Name())
	defer f.Close()
	err = sshClient.NodeUpload(context.Background(), "pve", "/var/lib/vz",
		&api.FileUploadRequest{
			ContentType: "snippets",
			FileName:    filepath.Base(snippetFile2.Name()),
			File:        f,
		})
	require.NoError(t, err)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: accProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceFileSnippetRawCreatedConfig(snippetRaw),
				Check:  testAccResourceFileSnippetRawCreatedCheck(snippetRaw),
			},
			{
				Config: testAccResourceFileCreatedConfig(snippetFile1.Name()),
				Check:  testAccResourceFileCreatedCheck("snippets", snippetFile1.Name()),
			},
			// allow to overwrite the a file by default
			{
				Config: testAccResourceFileCreatedConfig(snippetFile2.Name()),
				Check:  testAccResourceFileCreatedCheck("snippets", snippetFile2.Name()),
			},
			{
				Config: testAccResourceFileCreatedConfig(snippetURL),
				Check:  testAccResourceFileCreatedCheck("snippets", snippetURL),
			},
			{
				Config: testAccResourceFileCreatedConfig(fileISO.Name()),
				Check:  testAccResourceFileCreatedCheck("iso", fileISO.Name()),
			},
			{
				Config:      testAccResourceFileTwoSourcesCreatedConfig(),
				ExpectError: regexp.MustCompile("please specify .* - not both"),
			},
			{
				Config:      testAccResourceFileCreatedConfig("https://github.com", "content_type = \"iso\""),
				ExpectError: regexp.MustCompile("failed to determine file name from the URL"),
			},
			{
				Config:      testAccResourceFileMissingSourceConfig(),
				ExpectError: regexp.MustCompile("missing argument"),
			},
			// do not allow to overwrite the a file
			{
				Config:      testAccResourceFileCreatedConfig(snippetFile2.Name(), "overwrite = false"),
				ExpectError: regexp.MustCompile("already exists"),
			},
			// Update testing
			{
				Config: testAccResourceFileSnippetRawUpdatedConfig(snippetRaw),
				Check:  testAccResourceFileSnippetUpdatedCheck(snippetRaw),
			},
			// ImportState testing
			{
				ResourceName:      accTestFileName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     fmt.Sprintf("pve/local:snippets/%s", filepath.Base(snippetFile2.Name())),
				SkipFunc: func() (bool, error) {
					// TODO: add a file to the snippets directory outside of terraform
					// and then import it here
					return true, nil
				},
			},
		},
	})
}

func createFile(t *testing.T, namePattern string, content string) *os.File {
	t.Helper()

	f, err := os.CreateTemp("", namePattern)
	require.NoError(t, err)

	_, err = f.WriteString(content)
	require.NoError(t, err)

	defer f.Close()

	t.Cleanup(func() {
		_ = os.Remove(f.Name())
	})

	return f
}

func testAccResourceFileSnippetRawCreatedConfig(fname string) string {
	return fmt.Sprintf(`
resource "proxmox_virtual_environment_file" "test" {
  content_type = "snippets"
  datastore_id = "local"
  node_name    = "%s"
  source_raw {
    data = <<EOF
test snippet
    EOF
    file_name = "%s"
  }
}
	`, accTestNodeName, fname)
}

func testAccResourceFileCreatedConfig(fname string, extra ...string) string {
	return fmt.Sprintf(`
resource "proxmox_virtual_environment_file" "test" {
  datastore_id = "local"
  node_name    = "%s"
  source_file {
    path = "%s"
  }
  %s
}
	`, accTestNodeName, fname, strings.Join(extra, "\n"))
}

func testAccResourceFileTwoSourcesCreatedConfig() string {
	return fmt.Sprintf(`
resource "proxmox_virtual_environment_file" "test" {
  datastore_id = "local"
  node_name    = "%s"
  source_raw {
    data = <<EOF
test snippet
    EOF
	file_name = "foo.txt"
  }
  source_file {
    path = "bar.txt"
  }
}
	`, accTestNodeName)
}

func testAccResourceFileMissingSourceConfig() string {
	return fmt.Sprintf(`
resource "proxmox_virtual_environment_file" "test" {
  datastore_id = "local"
  node_name    = "%s"
}
	`, accTestNodeName)
}

func testAccResourceFileSnippetRawCreatedCheck(fname string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(accTestFileName, "content_type", "snippets"),
		resource.TestCheckResourceAttr(accTestFileName, "file_name", fname),
		resource.TestCheckResourceAttr(accTestFileName, "source_raw.0.file_name", fname),
		resource.TestCheckResourceAttr(accTestFileName, "source_raw.0.data", "test snippet\n"),
		resource.TestCheckResourceAttr(accTestFileName, "id", fmt.Sprintf("local:snippets/%s", fname)),
	)
}

func testAccResourceFileCreatedCheck(ctype string, fname string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(accTestFileName, "content_type", ctype),
		resource.TestCheckResourceAttr(accTestFileName, "file_name", filepath.Base(fname)),
		resource.TestCheckResourceAttr(accTestFileName, "id", fmt.Sprintf("local:%s/%s", ctype, filepath.Base(fname))),
	)
}

func testAccResourceFileSnippetRawUpdatedConfig(fname string) string {
	return fmt.Sprintf(`
resource "proxmox_virtual_environment_file" "test" {
  content_type = "snippets"
  datastore_id = "local"
  node_name    = "%s"

  source_raw {
    data = <<EOF
test snippet - updated
    EOF

    file_name = "%s"
  }
}
	`, accTestNodeName, fname)
}

func testAccResourceFileSnippetUpdatedCheck(fname string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(accTestFileName, "content_type", "snippets"),
		// resource.TestCheckResourceAttr(accTestFileName, "file_name", fname),
		resource.TestCheckResourceAttr(accTestFileName, "source_raw.0.file_name", fname),
		resource.TestCheckResourceAttr(accTestFileName, "source_raw.0.data", "test snippet - updated\n"),
		resource.TestCheckResourceAttr(accTestFileName, "id", fmt.Sprintf("local:snippets/%s", fname)),
	)
}
