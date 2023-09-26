/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package tests

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/require"
)

const (
	accTestFileName = "proxmox_virtual_environment_file.test"
)

func TestAccResourceFile(t *testing.T) {
	t.Parallel()

	accProviders := testAccMuxProviders(context.Background(), t)

	snippetRaw := fmt.Sprintf("snippet-raw-%s.txt", gofakeit.Word())
	snippetURL := "https://raw.githubusercontent.com/yaml/yaml-test-suite/main/src/229Q.yaml"

	snippetFile, err := os.CreateTemp("", "snippet-file-*.yaml")
	require.NoError(t, err)

	defer snippetFile.Close()

	_, err = snippetFile.WriteString("test snippet - file\n")
	require.NoError(t, err)

	snippetFileISO, err := os.CreateTemp("", "snippet-file-*.iso")
	require.NoError(t, err)

	defer snippetFile.Close()

	_, err = snippetFile.WriteString("pretend it is an ISO")
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = os.Remove(snippetFile.Name())
		_ = os.Remove(snippetFileISO.Name())
	})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: accProviders,
		Steps: []resource.TestStep{
			// Upload a snippet file from a raw source
			{
				Config: testAccResourceFileSnippetRawCreatedConfig(snippetRaw),
				Check:  testAccResourceFileSnippetRawCreatedCheck(snippetRaw),
				// RefreshState: true,
			},
			{
				Config:  testAccResourceFileCreatedConfig(snippetFile.Name()),
				Check:   testAccResourceFileCreatedCheck("snippets", snippetFile.Name()),
				Destroy: false,
			},
			{
				Config: testAccResourceFileCreatedConfig(snippetURL),
				Check:  testAccResourceFileCreatedCheck("snippets", snippetURL),
			},
			{
				Config: testAccResourceFileCreatedConfig(snippetFileISO.Name()),
				Check:  testAccResourceFileCreatedCheck("iso", snippetFileISO.Name()),
			},
			{
				Config:      testAccResourceFileWrongSourceCreatedConfig(),
				ExpectError: regexp.MustCompile("please specify .* - not both"),
			},
			// // ImportState testing
			// {
			// 	ResourceName:      accTestFileName,
			// 	ImportState:       true,
			// 	ImportStateVerify: true,
			// 	ImportStateId:     fmt.Sprintf("pve/local:snippets/%s", filepath.Base(snippetFile.Name())),
			// },
			// Update testing
			{
				Config: testAccResourceFileSnippetRawUpdatedConfig(snippetRaw),
				Check:  testAccResourceFileSnippetUpdatedCheck(snippetRaw),
			},
		},
	})
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

func testAccResourceFileCreatedConfig(fname string) string {
	return fmt.Sprintf(`
resource "proxmox_virtual_environment_file" "test" {
  datastore_id = "local"
  node_name    = "%s"

  source_file {
    path = "%s"
  }
}
	`, accTestNodeName, fname)
}

func testAccResourceFileWrongSourceCreatedConfig() string {
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
