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

type nodeResolver struct {
	node ssh.ProxmoxNode
}

func (c *nodeResolver) Resolve(_ context.Context, _ string) (ssh.ProxmoxNode, error) {
	return c.node, nil
}

func TestAccResourceFile(t *testing.T) {
	te := initTestEnvironment(t)

	snippetRaw := fmt.Sprintf("snippet-raw-%s.txt", gofakeit.Word())
	snippetURL := "https://raw.githubusercontent.com/yaml/yaml-test-suite/main/src/229Q.yaml"
	snippetFile1 := strings.ReplaceAll(createFile(t, "snippet-file-1-*.yaml", "test snippet 1 - file").Name(), `\`, `/`)
	snippetFile2 := strings.ReplaceAll(createFile(t, "snippet-file-2-*.yaml", "test snippet 2 - file").Name(), `\`, `/`)
	fileISO := strings.ReplaceAll(createFile(t, "file-*.iso", "pretend it is an ISO").Name(), `\`, `/`)

	te.addTemplateVars(map[string]interface{}{
		"SnippetRaw":   snippetRaw,
		"SnippetURL":   snippetURL,
		"SnippetFile1": snippetFile1,
		"SnippetFile2": snippetFile2,
		"FileISO":      fileISO,
	})

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: te.accProviders,
		PreCheck: func() {
			uploadSnippetFile(t, snippetFile2)
			t.Cleanup(func() {
				deleteSnippet(te, filepath.Base(snippetFile1))
				deleteSnippet(te, filepath.Base(snippetFile2))

				_ = os.Remove(snippetFile1)
				_ = os.Remove(snippetFile2)
				_ = os.Remove(fileISO)
			})
		},
		Steps: []resource.TestStep{
			{
				Config: te.renderConfig(`
				resource "proxmox_virtual_environment_file" "test_raw" {
				content_type = "snippets"
				datastore_id = "local"
				node_name    = "{{.NodeName}}"
				source_raw {
					data = <<EOF
				test snippet
					EOF
					file_name = "{{.SnippetRaw}}"
				}
				}`),
				Check: testResourceAttributes("proxmox_virtual_environment_file.test_raw", map[string]string{
					"content_type":           "snippets",
					"file_name":              snippetRaw,
					"source_raw.0.file_name": snippetRaw,
					"source_raw.0.data":      "test snippet\n",
					"id":                     fmt.Sprintf("local:snippets/%s", snippetRaw),
				}),
			},
			{
				Config: te.renderConfig(`
				resource "proxmox_virtual_environment_file" "test" {
					datastore_id = "local"
					node_name    = "{{.NodeName}}"
					source_file {
					  path = "{{.SnippetFile1}}"
					}
				}`),
				Check: testResourceAttributes("proxmox_virtual_environment_file.test", map[string]string{
					"content_type": "snippets",
					"file_name":    filepath.Base(snippetFile1),
					"id":           fmt.Sprintf("local:snippets/%s", filepath.Base(snippetFile1)),
				}),
			},
			{
				Config: te.renderConfig(`
				resource "proxmox_virtual_environment_file" "test" {
					datastore_id = "local"
					node_name    = "{{.NodeName}}"
					source_file {
					  path = "{{.SnippetURL}}"
					}
				}`),
				Check: testResourceAttributes("proxmox_virtual_environment_file.test", map[string]string{
					"content_type": "snippets",
					"file_name":    filepath.Base(snippetURL),
					"id":           fmt.Sprintf("local:snippets/%s", filepath.Base(snippetURL)),
				}),
			},
			{
				Config: te.renderConfig(`
				resource "proxmox_virtual_environment_file" "test" {
					datastore_id = "local"
					node_name    = "{{.NodeName}}"
					source_file {
					  path = "{{.FileISO}}"
					}
				}`),
				Check: testResourceAttributes("proxmox_virtual_environment_file.test", map[string]string{
					"content_type": "iso",
					"file_name":    filepath.Base(fileISO),
					"id":           fmt.Sprintf("local:iso/%s", filepath.Base(fileISO)),
				}),
			},
			{
				Config: te.renderConfig(`
				resource "proxmox_virtual_environment_file" "test" {
				  datastore_id = "local"
				  node_name    = "{{.NodeName}}"
				  source_raw {
					data = <<EOF
				test snippet
					EOF
					file_name = "foo.yaml"
				  }
				  source_file {
					path = "bar.yaml"
				  }
				}`),
				ExpectError: regexp.MustCompile("please specify .* - not both"),
			},
			{
				Config: te.renderConfig(`
				resource "proxmox_virtual_environment_file" "test" {
					datastore_id = "local"
					node_name    = "{{.NodeName}}"
					content_type = "iso"
					source_file {
					  path = "https://github.com"
					}
				}`),
				ExpectError: regexp.MustCompile("failed to determine file name from the URL"),
			},
			{
				Config: te.renderConfig(`
				resource "proxmox_virtual_environment_file" "test" {
				  datastore_id = "local"
				  node_name    = "{{.NodeName}}"
				}`),
				ExpectError: regexp.MustCompile("missing argument"),
			},
			// Do not allow to overwrite the file
			{
				Config: te.renderConfig(`
				resource "proxmox_virtual_environment_file" "test" {
					datastore_id = "local"
					node_name    = "{{.NodeName}}"
					overwrite    = false
					source_file {
					  path = "{{.SnippetFile2}}"
					}
				}`),
				ExpectError: regexp.MustCompile("already exists"),
			},
			// Allow to overwrite the file by default
			{
				Config: te.renderConfig(`
				resource "proxmox_virtual_environment_file" "test" {
					datastore_id = "local"
					node_name    = "{{.NodeName}}"
					source_file {
					  path = "{{.SnippetFile2}}"
					}
				}`),
				Check: testResourceAttributes("proxmox_virtual_environment_file.test", map[string]string{
					"content_type": "snippets",
					"file_name":    filepath.Base(snippetFile2),
					"id":           fmt.Sprintf("local:snippets/%s", filepath.Base(snippetFile2)),
				}),
			},
			// Update testing
			{
				PreConfig: func() {
					deleteSnippet(te, filepath.Base(snippetFile1))
				},
				Config: te.renderConfig(`
				resource "proxmox_virtual_environment_file" "test" {
				  datastore_id = "local"
				  node_name    = "{{.NodeName}}"
				  source_file {
					path = "{{.SnippetFile1}}"
				  }
				}`),
				Check: testResourceAttributes("proxmox_virtual_environment_file.test", map[string]string{
					"content_type": "snippets",
					"file_name":    filepath.Base(snippetFile1),
					"id":           fmt.Sprintf("local:snippets/%s", filepath.Base(snippetFile1)),
				}),
			},
		},
	})
}

func uploadSnippetFile(t *testing.T, fileName string) {
	t.Helper()

	endpoint := utils.GetAnyStringEnv("PROXMOX_VE_ENDPOINT")
	u, err := url.ParseRequestURI(endpoint)
	require.NoError(t, err)

	sshAgent := utils.GetAnyBoolEnv("PROXMOX_VE_SSH_AGENT")
	sshUsername := utils.GetAnyStringEnv("PROXMOX_VE_SSH_USERNAME")
	sshAgentSocket := utils.GetAnyStringEnv("SSH_AUTH_SOCK", "PROXMOX_VE_SSH_AUTH_SOCK")
	sshPrivateKey := utils.GetAnyStringEnv("PROXMOX_VE_SSH_PRIVATE_KEY")
	sshPort := utils.GetAnyIntEnv("PROXMOX_VE_ACC_NODE_SSH_PORT")
	sshClient, err := ssh.NewClient(
		sshUsername, "", sshAgent, sshAgentSocket, sshPrivateKey,
		"", "", "",
		&nodeResolver{
			node: ssh.ProxmoxNode{
				Address: u.Hostname(),
				Port:    int32(sshPort),
			},
		},
	)
	require.NoError(t, err)

	f, err := os.Open(fileName)
	require.NoError(t, err)

	defer f.Close()

	fname := filepath.Base(fileName)
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
