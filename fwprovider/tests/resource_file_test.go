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

	"github.com/brianvoe/gofakeit/v7"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/ssh"

	"github.com/bpg/terraform-provider-proxmox/utils"
)

const (
	accTestFileRawName = "proxmox_virtual_environment_file.test_raw"
	accTestFileName    = "proxmox_virtual_environment_file.test"
)

type nodeResolver struct {
	node ssh.ProxmoxNode
}

func (c *nodeResolver) Resolve(_ context.Context, _ string) (ssh.ProxmoxNode, error) {
	return c.node, nil
}

func TestAccResourceFile(t *testing.T) {
	t.Parallel()

	te := initTestEnvironment(t)

	snippetRaw := fmt.Sprintf("snippet-raw-%s.txt", gofakeit.Word())
	snippetURL := "https://raw.githubusercontent.com/yaml/yaml-test-suite/main/src/229Q.yaml"
	snippetFile1 := createFile(t, "snippet-file-1-*.yaml", "test snippet 1 - file")
	snippetFile2 := createFile(t, "snippet-file-2-*.yaml", "test snippet 2 - file")
	fileISO := createFile(t, "file-*.iso", "pretend it is an ISO")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: te.accProviders,
		PreCheck: func() {
			uploadSnippetFile(t, snippetFile2)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccResourceFileSnippetRawCreatedConfig(te, snippetRaw),
				Check:  testAccResourceFileSnippetRawCreatedCheck(snippetRaw),
			},
			{
				Config: testAccResourceFileCreatedConfig(te, snippetFile1.Name()),
				Check:  testAccResourceFileCreatedCheck("snippets", snippetFile1.Name()),
			},
			{
				Config: testAccResourceFileCreatedConfig(te, snippetURL),
				Check:  testAccResourceFileCreatedCheck("snippets", snippetURL),
			},
			{
				Config: testAccResourceFileCreatedConfig(te, fileISO.Name()),
				Check:  testAccResourceFileCreatedCheck("iso", fileISO.Name()),
			},
			{
				Config:      testAccResourceFileTwoSourcesCreatedConfig(te),
				ExpectError: regexp.MustCompile("please specify .* - not both"),
			},
			{
				Config:      testAccResourceFileCreatedConfig(te, "https://github.com", "content_type = \"iso\""),
				ExpectError: regexp.MustCompile("failed to determine file name from the URL"),
			},
			{
				Config:      testAccResourceFileMissingSourceConfig(te),
				ExpectError: regexp.MustCompile("missing argument"),
			},
			// Do not allow to overwrite the file
			{
				Config:      testAccResourceFileCreatedConfig(te, snippetFile2.Name(), "overwrite = false"),
				ExpectError: regexp.MustCompile("already exists"),
			},
			// Allow to overwrite the file by default
			{
				Config: testAccResourceFileCreatedConfig(te, snippetFile2.Name()),
				Check:  testAccResourceFileCreatedCheck("snippets", snippetFile2.Name()),
			},
			// Update testing
			{
				PreConfig: func() {
					deleteSnippet(te, filepath.Base(snippetFile1.Name()))
				},
				Config: testAccResourceFileSnippetUpdateConfig(te, snippetFile1.Name()),
				Check:  testAccResourceFileSnippetUpdatedCheck(snippetFile1.Name()),
			},
			// ImportState testing
			{
				ResourceName:      accTestFileName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     fmt.Sprintf("pve/local:snippets/%s", filepath.Base(snippetFile2.Name())),
				SkipFunc: func() (bool, error) {
					// doesn't work, not sure why
					return true, nil
				},
			},
		},
	})
}

func uploadSnippetFile(t *testing.T, file *os.File) {
	t.Helper()

	endpoint := utils.GetAnyStringEnv("PROXMOX_VE_ENDPOINT")
	u, err := url.ParseRequestURI(endpoint)
	require.NoError(t, err)

	sshAgent := utils.GetAnyBoolEnv("PROXMOX_VE_SSH_AGENT")
	sshUsername := utils.GetAnyStringEnv("PROXMOX_VE_SSH_USERNAME")
	sshAgentSocket := utils.GetAnyStringEnv("SSH_AUTH_SOCK", "PROXMOX_VE_SSH_AUTH_SOCK")
	sshPrivateKey := utils.GetAnyStringEnv("PROXMOX_VE_SSH_PRIVATE_KEY")
	sshClient, err := ssh.NewClient(
		sshUsername, "", sshAgent, sshAgentSocket, sshPrivateKey,
		"", "", "",
		&nodeResolver{
			node: ssh.ProxmoxNode{
				Address: u.Hostname(),
				Port:    22,
			},
		},
	)
	require.NoError(t, err)

	f, err := os.Open(file.Name())
	require.NoError(t, err)

	defer f.Close()

	fname := filepath.Base(file.Name())
	err = sshClient.NodeStreamUpload(context.Background(), "pve", "/var/lib/vz/",
		&api.FileUploadRequest{
			ContentType: "snippets",
			FileName:    fname,
			File:        f,
		})
	require.NoError(t, err)
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

func deleteSnippet(te *testEnvironment, fname string) {
	te.t.Helper()

	err := te.nodeStorageClient().DeleteDatastoreFile(context.Background(), fmt.Sprintf("snippets/%s", fname))
	require.NoError(te.t, err)
}

func testAccResourceFileSnippetRawCreatedConfig(te *testEnvironment, fname string) string {
	te.t.Helper()

	return fmt.Sprintf(`%s
resource "proxmox_virtual_environment_file" "test_raw" {
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
	`, te.providerConfig, te.nodeName, fname)
}

func testAccResourceFileCreatedConfig(te *testEnvironment, fname string, extra ...string) string {
	te.t.Helper()

	return fmt.Sprintf(`%s
resource "proxmox_virtual_environment_file" "test" {
  datastore_id = "local"
  node_name    = "%s"
  source_file {
    path = "%s"
  }
  %s
}
	`, te.providerConfig, te.nodeName, strings.ReplaceAll(fname, `\`, `/`), strings.Join(extra, "\n"))
}

func testAccResourceFileTwoSourcesCreatedConfig(te *testEnvironment) string {
	te.t.Helper()

	return fmt.Sprintf(`%s
resource "proxmox_virtual_environment_file" "test" {
  datastore_id = "local"
  node_name    = "%s"
  source_raw {
    data = <<EOF
test snippet
    EOF
	file_name = "foo.yaml"
  }
  source_file {
    path = "bar.yaml"
  }
}
	`, te.providerConfig, te.nodeName)
}

func testAccResourceFileMissingSourceConfig(te *testEnvironment) string {
	te.t.Helper()

	return fmt.Sprintf(`%s
resource "proxmox_virtual_environment_file" "test" {
  datastore_id = "local"
  node_name    = "%s"
}
	`, te.providerConfig, te.nodeName)
}

func testAccResourceFileSnippetRawCreatedCheck(fname string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(accTestFileRawName, "content_type", "snippets"),
		resource.TestCheckResourceAttr(accTestFileRawName, "file_name", fname),
		resource.TestCheckResourceAttr(accTestFileRawName, "source_raw.0.file_name", fname),
		resource.TestCheckResourceAttr(accTestFileRawName, "source_raw.0.data", "test snippet\n"),
		resource.TestCheckResourceAttr(accTestFileRawName, "id", fmt.Sprintf("local:snippets/%s", fname)),
	)
}

func testAccResourceFileCreatedCheck(ctype string, fname string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(accTestFileName, "content_type", ctype),
		resource.TestCheckResourceAttr(accTestFileName, "file_name", filepath.Base(fname)),
		resource.TestCheckResourceAttr(accTestFileName, "id", fmt.Sprintf("local:%s/%s", ctype, filepath.Base(fname))),
	)
}

func testAccResourceFileSnippetUpdateConfig(te *testEnvironment, fname string) string {
	te.t.Helper()

	return fmt.Sprintf(`%s
resource "proxmox_virtual_environment_file" "test" {
  datastore_id = "local"
  node_name    = "%s"
  source_file {
    path = "%s"
  }
}
	`, te.providerConfig, te.nodeName, strings.ReplaceAll(fname, `\`, `/`))
}

func testAccResourceFileSnippetUpdatedCheck(fname string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(accTestFileName, "content_type", "snippets"),
		resource.TestCheckResourceAttr(accTestFileName, "file_name", filepath.Base(fname)),
		resource.TestCheckResourceAttr(accTestFileName, "id", fmt.Sprintf("local:%s/%s", "snippets", filepath.Base(fname))),
	)
}
