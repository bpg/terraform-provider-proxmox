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
