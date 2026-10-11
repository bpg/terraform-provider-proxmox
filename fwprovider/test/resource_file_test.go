//go:build acceptance || all

//testacc:tier=heavy
//testacc:resource=file

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package test

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
)

func TestAccResourceFile(t *testing.T) {
	te := InitEnvironment(t)

	snippetRaw := SafeResourceName("snippet-raw") + ".txt"
	snippetFile1 := strings.ReplaceAll(CreateTempFile(t, "snippet-file-1-*.yaml", "test snippet 1 - file").Name(), `\`, `/`)
	snippetFile2 := strings.ReplaceAll(CreateTempFile(t, "snippet-file-2-*.yaml", "test snippet 2 - file").Name(), `\`, `/`)
	fileISO := strings.ReplaceAll(CreateTempFile(t, "file-*.iso", "pretend this is an ISO").Name(), `\`, `/`)
	sftpOverwriteFile := strings.ReplaceAll(CreateTempFile(t, "sftp-overwrite-*.yaml", "sftp overwrite test").Name(), `\`, `/`)
	sftpFile := strings.ReplaceAll(CreateTempFile(t, "sftp-file-*.yaml", "sftp upload").Name(), `\`, `/`)

	fileServer := NewTestFileServer(t)
	if fileServer == nil {
		t.Skip("PROXMOX_VE_ACC_TEST_FILE_SERVER_IP not set - skipping file resource test")
	}

	snippetContent := []byte("test: yaml\nkey: value\n")
	snippetURL := fileServer.AddFile("/229Q.yaml", "229Q.yaml", snippetContent)

	te.AddTemplateVars(map[string]interface{}{
		"SnippetRaw":        snippetRaw,
		"SnippetURL":        snippetURL,
		"SnippetFile1":      snippetFile1,
		"SnippetFile2":      snippetFile2,
		"FileISO":           fileISO,
		"SftpOverwriteFile": sftpOverwriteFile,
		"SftpFile":          sftpFile,
	})

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		PreCheck: func() {
			uploadSnippetFile(te, snippetFile2)
			uploadSnippetFile(te, sftpOverwriteFile)
			t.Cleanup(func() {
				deleteSnippet(te, filepath.Base(snippetFile1))
				deleteSnippet(te, filepath.Base(snippetFile2))
				deleteSnippet(te, filepath.Base(sftpOverwriteFile))

				_ = os.Remove(snippetFile1)
				_ = os.Remove(snippetFile2)
				_ = os.Remove(fileISO)
				_ = os.Remove(sftpOverwriteFile)
				_ = os.Remove(sftpFile)
			})
		},
		Steps: []resource.TestStep{
			{
				Config: te.RenderConfig(`
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
				Check: ResourceAttributes("proxmox_virtual_environment_file.test_raw", map[string]string{
					"content_type":           "snippets",
					"file_name":              snippetRaw,
					"source_raw.0.file_name": snippetRaw,
					"source_raw.0.data":      "test snippet\n",
					"id":                     fmt.Sprintf("local:snippets/%s", snippetRaw),
				}),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_file" "test" {
					datastore_id = "local"
					node_name    = "{{.NodeName}}"
					source_file {
					  path = "{{.SnippetFile1}}"
					}
				}`),
				Check: ResourceAttributes("proxmox_virtual_environment_file.test", map[string]string{
					"content_type": "snippets",
					"file_name":    filepath.Base(snippetFile1),
					"id":           fmt.Sprintf("local:snippets/%s", filepath.Base(snippetFile1)),
				}),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_file" "test" {
					datastore_id = "local"
					node_name    = "{{.NodeName}}"
					source_file {
					  path = "{{.SnippetURL}}"
					}
				}`),
				Check: ResourceAttributes("proxmox_virtual_environment_file.test", map[string]string{
					"content_type": "snippets",
					"file_name":    filepath.Base(snippetURL),
					"id":           fmt.Sprintf("local:snippets/%s", filepath.Base(snippetURL)),
				}),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_file" "test" {
					datastore_id = "local"
					node_name    = "{{.NodeName}}"
					source_file {
					  path = "{{.FileISO}}"
					}
				}`),
				Check: ResourceAttributes("proxmox_virtual_environment_file.test", map[string]string{
					"content_type": "iso",
					"file_name":    filepath.Base(fileISO),
					"id":           fmt.Sprintf("local:iso/%s", filepath.Base(fileISO)),
				}),
			},
			{
				Config: te.RenderConfig(`
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
				Config: te.RenderConfig(`
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
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_file" "test" {
				  datastore_id = "local"
				  node_name    = "{{.NodeName}}"
				}`),
				ExpectError: regexp.MustCompile("missing argument"),
			},
			// Do not allow to overwrite the file
			{
				Config: te.RenderConfig(`
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
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_file" "test" {
					datastore_id = "local"
					node_name    = "{{.NodeName}}"
					source_file {
					  path = "{{.SnippetFile2}}"
					}
				}`),
				Check: ResourceAttributes("proxmox_virtual_environment_file.test", map[string]string{
					"content_type": "snippets",
					"file_name":    filepath.Base(snippetFile2),
					"id":           fmt.Sprintf("local:snippets/%s", filepath.Base(snippetFile2)),
				}),
			},
			// Update testing: no original file
			{
				PreConfig: func() {
					_ = os.Remove(snippetFile2)
					deleteSnippet(te, filepath.Base(snippetFile1))
				},
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_file" "test" {
				  datastore_id = "local"
				  node_name    = "{{.NodeName}}"
				  source_file {
					path = "{{.SnippetFile1}}"
				  }
				}`),
				Check: ResourceAttributes("proxmox_virtual_environment_file.test", map[string]string{
					"content_type": "snippets",
					"file_name":    filepath.Base(snippetFile1),
					"id":           fmt.Sprintf("local:snippets/%s", filepath.Base(snippetFile1)),
				}),
			},
			// Update testing: original file
			{
				PreConfig: func() {
					deleteSnippet(te, filepath.Base(snippetFile1))
				},
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_file" "test" {
				  datastore_id = "local"
				  node_name    = "{{.NodeName}}"
				  source_file {
					path = "{{.SnippetFile1}}"
				  }
				}`),
				Check: ResourceAttributes("proxmox_virtual_environment_file.test", map[string]string{
					"content_type": "snippets",
					"file_name":    filepath.Base(snippetFile1),
					"id":           fmt.Sprintf("local:snippets/%s", filepath.Base(snippetFile1)),
				}),
			},
			// SFTP upload mode: file upload
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_file" "test_sftp_file" {
					datastore_id = "local"
					node_name    = "{{.NodeName}}"
					upload_mode  = "sftp"
					source_file {
					  path = "{{.SftpFile}}"
					}
				}`),
				Check: ResourceAttributes("proxmox_virtual_environment_file.test_sftp_file", map[string]string{
					"content_type": "snippets",
					"file_name":    filepath.Base(sftpFile),
					"id":           fmt.Sprintf("local:snippets/%s", filepath.Base(sftpFile)),
				}),
			},
			// SFTP upload mode: do not allow to overwrite the file
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_file" "test_sftp_overwrite" {
					datastore_id = "local"
					node_name    = "{{.NodeName}}"
					overwrite    = false
					upload_mode  = "sftp"
					source_file {
					  path = "{{.SftpOverwriteFile}}"
					}
				}`),
				ExpectError: regexp.MustCompile("already exists"),
			},
			// SFTP upload mode: allow to overwrite the file by default
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_file" "test_sftp_overwrite" {
					datastore_id = "local"
					node_name    = "{{.NodeName}}"
					upload_mode  = "sftp"
					source_file {
					  path = "{{.SftpOverwriteFile}}"
					}
				}`),
				Check: ResourceAttributes("proxmox_virtual_environment_file.test_sftp_overwrite", map[string]string{
					"content_type": "snippets",
					"file_name":    filepath.Base(sftpOverwriteFile),
					"id":           fmt.Sprintf("local:snippets/%s", filepath.Base(sftpOverwriteFile)),
				}),
			},
		},
	})
}

// TestAccNodeStreamUpload verifies that NodeStreamUpload succeeds when the SSH
// user relies on sudo to write to the snippets directory.
func TestAccNodeStreamUpload(t *testing.T) {
	t.Skip("disabled until NodeStreamUpload can chmod a sudo-written file as a non-root SSH user (#2987)")

	te := InitEnvironment(t)

	client := te.SSHClient()
	if client.Username() == "root" {
		t.Skip("requires a non-root sudo SSH user")
	}

	f := CreateTempFile(t, "stream-upload-*.yaml", "#cloud-config\nruncmd:\n  - echo hello\n")
	fname := filepath.Base(f.Name())

	t.Cleanup(func() {
		err := te.NodeStorageClient().DeleteDatastoreFile(context.Background(), fmt.Sprintf("snippets/%s", fname)).Err()
		if err != nil {
			t.Logf("cleanup: failed to delete snippet %s: %v", fname, err)
		}
	})

	fh, err := os.Open(f.Name())
	require.NoError(t, err)

	t.Cleanup(func() { _ = fh.Close() })

	err = client.NodeStreamUpload(
		context.Background(),
		te.NodeName,
		"/var/lib/vz/",
		&api.FileUploadRequest{
			ContentType: "snippets",
			FileName:    fname,
			File:        fh,
			Mode:        "0700",
		},
	)
	require.NoError(t, err)

	out := te.ExecuteNodeCommands([]string{"stat -c '%a' /var/lib/vz/snippets/" + fname})
	require.Equal(t, "700", strings.TrimSpace(out))
}

func uploadSnippetFile(te *Environment, fileName string) {
	te.t.Helper()

	f, err := os.Open(fileName)
	require.NoError(te.t, err)

	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	err = te.SSHClient().NodeStreamUpload(context.Background(), te.NodeName, "/var/lib/vz/",
		&api.FileUploadRequest{
			ContentType: "snippets",
			FileName:    filepath.Base(fileName),
			File:        f,
		})
	require.NoError(te.t, err)
}

func deleteSnippet(te *Environment, fname string) {
	te.t.Helper()

	err := te.NodeStorageClient().DeleteDatastoreFile(context.Background(), fmt.Sprintf("snippets/%s", fname)).Err()
	require.NoError(te.t, err)
}

// TestAccDeleteDatastoreFileWaitsForTask verifies that deleting a datastore
// file waits for the completion of the async deletion task started by PVE,
// and surfaces the task's error if the deletion fails.
//
// Regression test for #3062: PVE storage DELETE returns 200 with the UPID of
// an async imgdel task. A protected backup is used as the discriminator: the
// imgdel worker fails with "cannot remove protected volume", but the HTTP
// DELETE still returns 200. Without waiting for the task, the failure is
// invisible and the delete appears to succeed while the file remains.
func TestAccDeleteDatastoreFileWaitsForTask(t *testing.T) {
	te := InitEnvironment(t)

	sc := te.NodeStorageClient()

	fname := "vzdump-qemu-999-2026_01_02-03_04_05.vma.zst"
	volid := "backup/" + fname

	t.Cleanup(func() {
		// best effort: un-protect then delete
		_ = te.NodeClient().DoRequest(
			context.Background(), http.MethodPut,
			te.NodeClient().ExpandPath(fmt.Sprintf("storage/%s/content/%s", sc.StorageName, url.PathEscape(volid))),
			&struct {
				Protected string `url:"protected"`
			}{Protected: "0"},
			nil,
		)
		_ = sc.DeleteDatastoreFile(context.Background(), volid).Err()
	})

	// Create a backup file directly in the datastore directory.
	f, err := os.CreateTemp("", fname)
	require.NoError(t, err)

	_, werr := f.WriteString("# backup placeholder\n")
	require.NoError(t, werr)
	require.NoError(t, f.Close())

	t.Cleanup(func() { _ = os.Remove(f.Name()) })

	uf, uerr := os.Open(f.Name())
	require.NoError(t, uerr)

	t.Cleanup(func() { _ = uf.Close() })

	// NodeStreamUpload appends the content-type subdir to remote_dir itself;
	// for backups that on-disk subdir is "dump", so pass "/var/lib/vz" and
	// "dump" as the content type to land in /var/lib/vz/dump/.
	require.NoError(t, te.SSHClient().NodeStreamUpload(context.Background(), te.NodeName, "/var/lib/vz",
		&api.FileUploadRequest{
			ContentType: "dump",
			FileName:    fname,
			File:        uf,
		},
	))

	// Protect it: only backups support the 'protected' attribute, and the
	// imgdel worker refuses to remove protected volumes.
	require.NoError(t, te.NodeClient().DoRequest(
		context.Background(), http.MethodPut,
		te.NodeClient().ExpandPath(fmt.Sprintf("storage/%s/content/%s", sc.StorageName, url.PathEscape(volid))),
		&struct {
			Protected string `url:"protected"`
		}{Protected: "1"},
		nil,
	))

	// With the fix, the delete must fail with the task's error...
	err = sc.DeleteDatastoreFile(context.Background(), volid).Err()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cannot remove protected volume")

	// ...and the file must still exist (checked on the filesystem, since the
	// test datastore may not declare the "backup" content type).
	out, execErr := te.SSHClient().ExecuteNodeCommands(context.Background(), te.NodeName,
		[]string{fmt.Sprintf("test -f /var/lib/vz/dump/%s && echo present", fname)})
	require.NoError(t, execErr)
	assert.Contains(t, string(out), "present", "protected file should still exist after failed delete")
}
