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
	"path/filepath"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestAccDatasourceFiles(t *testing.T) {
	te := InitEnvironment(t)

	// Upload a snippet file so we have at least one file to list
	snippetFile := CreateTempFile(t, "files-ds-test-*.yaml", "test: yaml\nkey: value\n")
	uploadSnippetFile(t, snippetFile.Name())

	fileName := filepath.Base(snippetFile.Name())

	te.AddTemplateVars(map[string]interface{}{
		"TestFileName": fileName,
	})

	t.Cleanup(func() {
		e := te.NodeStorageClient().DeleteDatastoreFile(
			context.Background(), fmt.Sprintf("snippets/%s", fileName))
		require.NoError(t, e)
	})

	datasourceName := "data.proxmox_virtual_environment_files.test"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: te.RenderConfig(`
					data "proxmox_virtual_environment_files" "test" {
						node_name    = "{{.NodeName}}"
						datastore_id = "local"
						content_type = "snippets"
					}
				`),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "node_name", te.NodeName),
					resource.TestCheckResourceAttr(datasourceName, "datastore_id", "local"),
					resource.TestCheckResourceAttr(datasourceName, "content_type", "snippets"),
					// At least one file should exist (the one we uploaded)
					resource.TestCheckResourceAttrSet(datasourceName, "files.#"),
				),
			},
		},
	})
}

func TestAccDatasourceFilesNoFilter(t *testing.T) {
	te := InitEnvironment(t)

	// Upload a snippet so there's at least one file
	snippetFile := CreateTempFile(t, "files-ds-nofilter-*.yaml", "test: yaml\n")
	uploadSnippetFile(t, snippetFile.Name())

	fileName := filepath.Base(snippetFile.Name())

	t.Cleanup(func() {
		e := te.NodeStorageClient().DeleteDatastoreFile(
			context.Background(), fmt.Sprintf("snippets/%s", fileName))
		require.NoError(t, e)
	})

	datasourceName := "data.proxmox_virtual_environment_files.test_no_filter"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: te.RenderConfig(`
					data "proxmox_virtual_environment_files" "test_no_filter" {
						node_name    = "{{.NodeName}}"
						datastore_id = "local"
					}
				`),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "node_name", te.NodeName),
					resource.TestCheckResourceAttr(datasourceName, "datastore_id", "local"),
					resource.TestCheckResourceAttrSet(datasourceName, "files.#"),
				),
			},
		},
	})
}

func TestAccDatasourceFilesEmptyResult(t *testing.T) {
	te := InitEnvironment(t)

	datasourceName := "data.proxmox_virtual_environment_files.test_empty"

	// "import" content type is unlikely to have files on a default local datastore
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: te.RenderConfig(`
					data "proxmox_virtual_environment_files" "test_empty" {
						node_name    = "{{.NodeName}}"
						datastore_id = "local"
						content_type = "import"
					}
				`),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "node_name", te.NodeName),
					resource.TestCheckResourceAttr(datasourceName, "datastore_id", "local"),
					resource.TestCheckResourceAttr(datasourceName, "content_type", "import"),
					// Empty list — no error
					resource.TestCheckResourceAttr(datasourceName, "files.#", "0"),
				),
			},
		},
	})
}
