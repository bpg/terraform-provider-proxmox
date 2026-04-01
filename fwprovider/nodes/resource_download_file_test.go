//go:build acceptance || all

//testacc:tier=heavy
//testacc:resource=file

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package nodes_test

import (
	"context"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/test"
	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/helpers/ptr"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/storage"
	"github.com/bpg/terraform-provider-proxmox/proxmox/ssh"
	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
	"github.com/bpg/terraform-provider-proxmox/utils"
)

func TestAccResourceDownloadFile(t *testing.T) {
	te := test.InitEnvironment(t)

	fileServer := test.NewTestFileServer(t)
	if fileServer == nil {
		t.Skip("PROXMOX_VE_ACC_TEST_FILE_SERVER_IP not set - skipping download file test")
	}

	content := []byte("asd")
	fakeFileISO := fileServer.AddFile("/fake_file.iso", "fake_file.iso", content)
	fakeFileQCOW2 := fileServer.AddFile("/fake_file.qcow2", "fake_file.qcow2", content)
	checksum := fileServer.GetFileSHA256("/fake_file.iso")

	te.AddTemplateVars(map[string]interface{}{
		"FakeFileISO":   fakeFileISO,
		"FakeFileQCOW2": fakeFileQCOW2,
		"Checksum":      checksum,
	})

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"missing url", []resource.TestStep{{
			Config: te.RenderConfig(`
				resource "proxmox_download_file" "qcow2_image" {
					content_type       = "iso"
					node_name          = "{{.NodeName}}"
					datastore_id       = "{{.DatastoreID}}"
					file_name          = "fake_qcow2_file.img"
					url                =  ""
				  }`),
			ExpectError: regexp.MustCompile(`Attribute url must match HTTP URL regex`),
		}}},
		{"download qcow2 file to iso storage", []resource.TestStep{{
			Config: te.RenderConfig(`
				resource "proxmox_download_file" "qcow2_image" {
					content_type       = "iso"
					node_name          = "{{.NodeName}}"
					datastore_id       = "{{.DatastoreID}}"
					file_name          = "fake_qcow2_file.img"
					url                =  "{{.FakeFileQCOW2}}"
					checksum           = "{{.Checksum}}"
					checksum_algorithm = "sha256"
					overwrite_unmanaged = true
				  }`),
			Check: resource.ComposeTestCheckFunc(
				test.ResourceAttributes("proxmox_download_file.qcow2_image", map[string]string{
					"id":                 "local:iso/fake_qcow2_file.img",
					"content_type":       "iso",
					"node_name":          te.NodeName,
					"datastore_id":       te.DatastoreID,
					"url":                fakeFileQCOW2,
					"file_name":          "fake_qcow2_file.img",
					"upload_timeout":     "600",
					"size":               "3",
					"verify":             "true",
					"checksum":           checksum,
					"checksum_algorithm": "sha256",
				}),
				test.NoResourceAttributesSet("proxmox_download_file.qcow2_image", []string{
					"decompression_algorithm",
				}),
			),
		}}},
		{"download qcow2 file to import storage", []resource.TestStep{{
			Config: te.RenderConfig(`
				resource "proxmox_download_file" "qcow2_image" {
					content_type       = "import"
					node_name          = "{{.NodeName}}"
					datastore_id       = "{{.DatastoreID}}"
					file_name          = "fake_qcow2_file.qcow2"
					url                =  "{{.FakeFileQCOW2}}"
					checksum           = "{{.Checksum}}"
					checksum_algorithm = "sha256"
					overwrite_unmanaged = true
				  }`),
			// the details says "Image is not in qcow2 format", but we can't assert that
			ExpectError: regexp.MustCompile(`Error downloading file from url`),
		}}},
		{"download & update iso file", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_download_file" "iso_image" {
					content_type = "iso"
					node_name    = "{{.NodeName}}"
					datastore_id = "{{.DatastoreID}}"
					url          = "{{.FakeFileISO}}"
					overwrite_unmanaged = true
				  }`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_download_file.iso_image", map[string]string{
						"id":             "local:iso/fake_file.iso",
						"node_name":      te.NodeName,
						"datastore_id":   te.DatastoreID,
						"url":            fakeFileISO,
						"file_name":      "fake_file.iso",
						"upload_timeout": "600",
						"size":           "3",
						"verify":         "true",
					}),
					test.NoResourceAttributesSet("proxmox_download_file.iso_image", []string{
						"checksum",
						"checksum_algorithm",
						"decompression_algorithm",
					}),
				),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_download_file" "iso_image" {
					content_type   = "iso"
					node_name      = "{{.NodeName}}"
					datastore_id   = "{{.DatastoreID}}"
					file_name      = "fake_iso_file.img"
					url            = "{{.FakeFileISO}}"
					upload_timeout = 10000
					overwrite_unmanaged = true
				  }`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_download_file.iso_image", map[string]string{
						"id":             "local:iso/fake_iso_file.img",
						"content_type":   "iso",
						"node_name":      te.NodeName,
						"datastore_id":   te.DatastoreID,
						"url":            fakeFileISO,
						"file_name":      "fake_iso_file.img",
						"upload_timeout": "10000",
						"size":           "3",
						"verify":         "true",
					}),
					test.NoResourceAttributesSet("proxmox_download_file.iso_image", []string{
						"checksum",
						"checksum_algorithm",
						"decompression_algorithm",
					}),
				),
			},
		}},
		{"override file", []resource.TestStep{{
			Destroy: false,
			PreConfig: func() {
				ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
				defer cancel()

				_ = te.NodeStorageClient().DeleteDatastoreFile(ctx, "iso/fake_file.iso") //nolint: errcheck

				err := te.NodeStorageClient().DownloadFileByURL(ctx, &storage.DownloadURLPostRequestBody{
					Content:  ptr.Ptr("iso"),
					FileName: ptr.Ptr("fake_file.iso"),
					Node:     ptr.Ptr(te.NodeName),
					Storage:  ptr.Ptr(te.DatastoreID),
					URL:      ptr.Ptr(fakeFileISO),
					Verify:   ptr.Ptr(types.CustomBool(false)),
				})
				require.NoError(t, err)

				t.Cleanup(func() {
					e := te.NodeStorageClient().DeleteDatastoreFile(context.Background(), "iso/fake_file.iso")
					require.NoError(t, e)
				})
			},
			Config: te.RenderConfig(`
				resource "proxmox_download_file" "iso_image3" {
					content_type        = "iso"
					node_name           = "{{.NodeName}}"
					datastore_id        = "{{.DatastoreID}}"
					url 		        = "{{.FakeFileISO}}"
					file_name           = "fake_iso_file3.iso"
					overwrite_unmanaged = true
					overwrite           = false
				  }`),
			Check: resource.ComposeTestCheckFunc(
				test.ResourceAttributes("proxmox_download_file.iso_image3", map[string]string{
					"id":           "local:iso/fake_iso_file3.iso",
					"content_type": "iso",
					"node_name":    te.NodeName,
					"datastore_id": te.DatastoreID,
					"url":          fakeFileISO,
					"file_name":    "fake_iso_file3.iso",
					"size":         "3",
					"verify":       "true",
				}),
				test.NoResourceAttributesSet("proxmox_download_file.iso_image3", []string{
					"checksum",
					"checksum_algorithm",
					"decompression_algorithm",
				}),
			),
		}, {
			Destroy: false,
			PreConfig: func() {
				isoFile := strings.ReplaceAll(createFile(t, "fake_iso_file3.iso", "updated iso").Name(), `\`, `/`)
				uploadIsoFile(t, isoFile)
			},
			Config: te.RenderConfig(`
				resource "proxmox_download_file" "iso_image3" {
					content_type        = "iso"
					node_name           = "{{.NodeName}}"
					datastore_id        = "{{.DatastoreID}}"
					url 		        = "{{.FakeFileISO}}"
					file_name           = "fake_iso_file3.iso"
					overwrite_unmanaged = true
					overwrite           = false
				}`),
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					plancheck.ExpectEmptyPlan(),
				},
			},
		}, {
			PreConfig: func() {
				isoFile := strings.ReplaceAll(createFile(t, "fake_iso_file3.iso", "updated iso again").Name(), `\`, `/`)
				uploadIsoFile(t, isoFile)
			},
			Config: te.RenderConfig(`
				resource "proxmox_download_file" "iso_image3" {
					content_type        = "iso"
					node_name           = "{{.NodeName}}"
					datastore_id        = "{{.DatastoreID}}"
					url 		        = "{{.FakeFileISO}}"
					file_name           = "fake_iso_file3.iso"
					overwrite_unmanaged = true
					overwrite           = true
				}`),
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction("proxmox_download_file.iso_image3", plancheck.ResourceActionDestroyBeforeCreate),
				},
			},
		}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource.ParallelTest(t, resource.TestCase{
				ProtoV6ProviderFactories: te.AccProviders,
				Steps:                    tt.steps,
			})
		})
	}
}

func uploadIsoFile(t *testing.T, fileName string) {
	t.Helper()

	endpoint := utils.GetAnyStringEnv("PROXMOX_VE_ENDPOINT")
	u, err := url.ParseRequestURI(endpoint)
	require.NoError(t, err)

	sshAgent := utils.GetAnyBoolEnv("PROXMOX_VE_SSH_AGENT")
	sshUsername := utils.GetAnyStringEnv("PROXMOX_VE_SSH_USERNAME")
	sshPassword := utils.GetAnyStringEnv("PROXMOX_VE_SSH_PASSWORD")
	sshAgentSocket := utils.GetAnyStringEnv("SSH_AUTH_SOCK", "PROXMOX_VE_SSH_AUTH_SOCK")
	sshAgentForwarding := utils.GetAnyBoolEnv("PROXMOX_VE_SSH_AGENT_FORWARDING")
	sshPrivateKey := utils.GetAnyStringEnv("PROXMOX_VE_SSH_PRIVATE_KEY")
	sshPort := utils.GetAnyIntEnv("PROXMOX_VE_ACC_NODE_SSH_PORT")
	sshClient, err := ssh.NewClient(
		sshUsername, sshPassword, sshAgent, sshAgentSocket, sshAgentForwarding, sshPrivateKey,
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

	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	fname := filepath.Base(fileName)
	err = sshClient.NodeStreamUpload(context.Background(), "pve", "/var/lib/vz/template",
		&api.FileUploadRequest{
			ContentType: "iso",
			FileName:    fname,
			File:        f,
		})
	require.NoError(t, err)
}

type nodeResolver struct {
	node ssh.ProxmoxNode
}

func (c *nodeResolver) Resolve(_ context.Context, _ string) (ssh.ProxmoxNode, error) {
	return c.node, nil
}

func createFile(t *testing.T, namePattern string, content string) *os.File {
	t.Helper()

	f, err := os.Create(path.Join(os.TempDir(), namePattern))
	require.NoError(t, err)

	_, err = f.WriteString(content)
	require.NoError(t, err)

	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	t.Cleanup(func() {
		_ = os.Remove(f.Name())
	})

	return f
}
