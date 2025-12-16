//go:build acceptance || all

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/proxmox/helpers/ptr"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/storage"
)

func TestAccDatasourceFile(t *testing.T) {
	t.Parallel()

	te := InitEnvironment(t)

	fileName := gofakeit.Word() + "-test-file.yaml"

	te.AddTemplateVars(map[string]interface{}{
		"TestFileName": fileName,
	})

	err := te.NodeStorageClient().DownloadFileByURL(context.Background(), &storage.DownloadURLPostRequestBody{
		Content:  ptr.Ptr("snippets"),
		FileName: ptr.Ptr(fileName),
		Node:     ptr.Ptr(te.NodeName),
		Storage:  ptr.Ptr("local"),
		URL:      ptr.Ptr("https://raw.githubusercontent.com/yaml/yaml-test-suite/main/src/229Q.yaml"),
	})
	require.NoError(t, err)

	t.Cleanup(func() {
		e := te.NodeStorageClient().DeleteDatastoreFile(context.Background(), fmt.Sprintf("snippets/%s", fileName))
		require.NoError(t, e)
	})

	datasourceName := "data.proxmox_virtual_environment_file.test"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: te.RenderConfig(`
					data "proxmox_virtual_environment_file" "test" {
						node_name    = "{{.NodeName}}"
						datastore_id = "local"
						content_type = "snippets"
						file_name    = "{{.TestFileName}}"
					}
				`),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "node_name", te.NodeName),
					resource.TestCheckResourceAttr(datasourceName, "datastore_id", "local"),
					resource.TestCheckResourceAttr(datasourceName, "content_type", "snippets"),
					resource.TestCheckResourceAttr(datasourceName, "file_name", fileName),
					resource.TestCheckResourceAttr(datasourceName, "id", fmt.Sprintf("local:snippets/%s", fileName)),
					resource.TestCheckResourceAttrSet(datasourceName, "file_size"),
					resource.TestCheckResourceAttrSet(datasourceName, "file_format"),
				),
			},
		},
	})
}

func TestAccDatasourceFileImport(t *testing.T) {
	te := InitEnvironment(t)

	// Test that the import content type is accepted in the schema
	// Since import content type doesn't support API uploads, we test schema validation
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: `
					data "proxmox_virtual_environment_file" "test_import" {
						node_name    = "pve"
						datastore_id = "local"
						content_type = "import"
						file_name    = "non-existent-file.yaml"
					}
				`,
				ExpectError: regexp.MustCompile("File Not Found"),
			},
		},
	})
}

func TestAccDatasourceFileNotFound(t *testing.T) {
	t.Parallel()

	te := InitEnvironment(t)

	nonExistentFileName := "non-existent-" + gofakeit.Word() + ".txt"

	te.AddTemplateVars(map[string]interface{}{
		"NonExistentFileName": nonExistentFileName,
	})

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: te.RenderConfig(`
					data "proxmox_virtual_environment_file" "test_not_found" {
						node_name    = "{{.NodeName}}"
						datastore_id = "local"
						content_type = "snippets"
						file_name    = "{{.NonExistentFileName}}"
					}
				`),
				ExpectError: regexp.MustCompile("File Not Found"),
			},
		},
	})
}

// TestAccDatasourceFileContentTypeFiltering verifies server-side content type filtering works correctly.
func TestAccDatasourceFileContentTypeFiltering(t *testing.T) {
	te := InitEnvironment(t)

	vztmplFileName := gofakeit.Word() + "-template.tar.zst"
	isoFileName := gofakeit.Word() + "-test.iso"

	te.AddTemplateVars(map[string]interface{}{
		"VZTmplFileName": vztmplFileName,
		"ISOFileName":    isoFileName,
	})

	// Upload a vztmpl file (container template)
	err := te.NodeStorageClient().DownloadFileByURL(context.Background(), &storage.DownloadURLPostRequestBody{
		Content:  ptr.Ptr("vztmpl"),
		FileName: ptr.Ptr(vztmplFileName),
		Node:     ptr.Ptr(te.NodeName),
		Storage:  ptr.Ptr("local"),
		URL:      ptr.Ptr("http://download.proxmox.com/images/system/alpine-3.19-default_20240207_amd64.tar.xz"),
	})
	require.NoError(t, err)

	// Upload an ISO file (small test ISO)
	err = te.NodeStorageClient().DownloadFileByURL(context.Background(), &storage.DownloadURLPostRequestBody{
		Content:  ptr.Ptr("iso"),
		FileName: ptr.Ptr(isoFileName),
		Node:     ptr.Ptr(te.NodeName),
		Storage:  ptr.Ptr("local"),
		URL:      ptr.Ptr("https://boot.netboot.xyz/ipxe/netboot.xyz.iso"),
	})
	require.NoError(t, err)

	t.Cleanup(func() {
		e := te.NodeStorageClient().DeleteDatastoreFile(context.Background(), fmt.Sprintf("vztmpl/%s", vztmplFileName))
		require.NoError(t, e)
		e = te.NodeStorageClient().DeleteDatastoreFile(context.Background(), fmt.Sprintf("iso/%s", isoFileName))
		require.NoError(t, e)
	})

	datasourceVZTmplName := "data.proxmox_virtual_environment_file.test_vztmpl"
	datasourceISOName := "data.proxmox_virtual_environment_file.test_iso"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: te.RenderConfig(`
					data "proxmox_virtual_environment_file" "test_vztmpl" {
						node_name    = "{{.NodeName}}"
						datastore_id = "local"
						content_type = "vztmpl"
						file_name    = "{{.VZTmplFileName}}"
					}

					data "proxmox_virtual_environment_file" "test_iso" {
						node_name    = "{{.NodeName}}"
						datastore_id = "local"
						content_type = "iso"
						file_name    = "{{.ISOFileName}}"
					}
				`),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify vztmpl file is found with correct content type
					resource.TestCheckResourceAttr(datasourceVZTmplName, "node_name", te.NodeName),
					resource.TestCheckResourceAttr(datasourceVZTmplName, "datastore_id", "local"),
					resource.TestCheckResourceAttr(datasourceVZTmplName, "content_type", "vztmpl"),
					resource.TestCheckResourceAttr(datasourceVZTmplName, "file_name", vztmplFileName),
					resource.TestCheckResourceAttr(datasourceVZTmplName, "id", fmt.Sprintf("local:vztmpl/%s", vztmplFileName)),
					resource.TestCheckResourceAttrSet(datasourceVZTmplName, "file_size"),

					// Verify ISO file is found with correct content type
					resource.TestCheckResourceAttr(datasourceISOName, "node_name", te.NodeName),
					resource.TestCheckResourceAttr(datasourceISOName, "datastore_id", "local"),
					resource.TestCheckResourceAttr(datasourceISOName, "content_type", "iso"),
					resource.TestCheckResourceAttr(datasourceISOName, "file_name", isoFileName),
					resource.TestCheckResourceAttr(datasourceISOName, "id", fmt.Sprintf("local:iso/%s", isoFileName)),
					resource.TestCheckResourceAttrSet(datasourceISOName, "file_size"),
				),
			},
		},
	})
}

// TestAccDatasourceFileContentTypeMismatch verifies that filtering by wrong content type returns not found.
func TestAccDatasourceFileContentTypeMismatch(t *testing.T) {
	te := InitEnvironment(t)

	isoFileName := gofakeit.Word() + "-mismatch-test.iso"

	te.AddTemplateVars(map[string]interface{}{
		"ISOFileName": isoFileName,
	})

	// Upload an ISO file
	err := te.NodeStorageClient().DownloadFileByURL(context.Background(), &storage.DownloadURLPostRequestBody{
		Content:  ptr.Ptr("iso"),
		FileName: ptr.Ptr(isoFileName),
		Node:     ptr.Ptr(te.NodeName),
		Storage:  ptr.Ptr("local"),
		URL:      ptr.Ptr("https://boot.netboot.xyz/ipxe/netboot.xyz.iso"),
	})
	require.NoError(t, err)

	t.Cleanup(func() {
		e := te.NodeStorageClient().DeleteDatastoreFile(context.Background(), fmt.Sprintf("iso/%s", isoFileName))
		require.NoError(t, e)
	})

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				// Try to find the ISO file with wrong content type (vztmpl)
				Config: te.RenderConfig(`
					data "proxmox_virtual_environment_file" "test_mismatch" {
						node_name    = "{{.NodeName}}"
						datastore_id = "local"
						content_type = "vztmpl"
						file_name    = "{{.ISOFileName}}"
					}
				`),
				ExpectError: regexp.MustCompile("File Not Found"),
			},
		},
	})
}
