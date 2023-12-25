/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	accTestDownloadIsoFileName   = "proxmox_virtual_environment_download_file.iso_image"
	accTestDownloadQcow2FileName = "proxmox_virtual_environment_download_file.qcow2_image"
)

func TestAccResourceDownloadFile(t *testing.T) {
	t.Parallel()

	accProviders := testAccMuxProviders(context.Background(), t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: accProviders,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccResourceDownloadIsoFileCreatedConfig(),
				Check:  testAccResourceDownloadIsoFileCreatedCheck(),
			},
			{
				Config: testAccResourceDownloadQcow2FileCreatedConfig(),
				Check:  testAccResourceDownloadQcow2FileCreatedCheck(),
			},
			// Update testing
			{
				Config: testAccResourceDownloadIsoFileUpdatedConfig(),
				Check:  testAccResourceDownloadIsoFileUpdatedCheck(),
			},
		},
	})
}

func testAccResourceDownloadIsoFileCreatedConfig() string {
	return fmt.Sprintf(`
	resource "proxmox_virtual_environment_download_file" "iso_image" {
		content_type = "iso"
		node_name    = "%s"
		datastore_id = "%s"
		download_url = "https://cdn.githubraw.com/rafsaf/a4b19ea5e3485f8da6ca4acf46d09650/raw/d340ec3ddcef9b907ede02f64b5d3f694da5d081/fake_file.iso"
	  }
	 `, accTestNodeName, accTestStorageName)
}

func testAccResourceDownloadQcow2FileCreatedConfig() string {
	return fmt.Sprintf(`
	resource "proxmox_virtual_environment_download_file" "qcow2_image" {
		content_type = "iso"
		node_name    = "%s"
		datastore_id = "%s"
		download_url = "https://cdn.githubraw.com/rafsaf/036eece601975a3ad632a77fc2809046/raw/10500012fca9b4425b50de67a7258a12cba0c076/fake_file.qcow2"
		allow_unsupported_types = true
	  }
	 `, accTestNodeName, accTestStorageName)
}

func testAccResourceDownloadIsoFileCreatedCheck() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(accTestDownloadIsoFileName, "content_type", "iso"),
		resource.TestCheckResourceAttr(accTestDownloadIsoFileName, "node_name", accTestNodeName),
		resource.TestCheckResourceAttr(accTestDownloadIsoFileName, "datastore_id", accTestStorageName),
		resource.TestCheckResourceAttr(accTestDownloadIsoFileName, "download_url", "https://cdn.githubraw.com/rafsaf/a4b19ea5e3485f8da6ca4acf46d09650/raw/d340ec3ddcef9b907ede02f64b5d3f694da5d081/fake_file.iso"),
		resource.TestCheckResourceAttr(accTestDownloadIsoFileName, "filename", "fake_file.iso"),
		resource.TestCheckResourceAttr(accTestDownloadIsoFileName, "allow_unsupported_types", "false"),
		resource.TestCheckResourceAttr(accTestDownloadIsoFileName, "upload_timeout", "600"),
		resource.TestCheckResourceAttr(accTestDownloadIsoFileName, "size", "512"),
		resource.TestCheckResourceAttr(accTestDownloadIsoFileName, "verify", "true"),
		resource.TestCheckNoResourceAttr(accTestDownloadIsoFileName, "checksum"),
		resource.TestCheckNoResourceAttr(accTestDownloadIsoFileName, "checksum_algorithm"),
		resource.TestCheckNoResourceAttr(accTestDownloadIsoFileName, "compression"),
	)
}

func testAccResourceDownloadQcow2FileCreatedCheck() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(accTestDownloadQcow2FileName, "content_type", "iso"),
		resource.TestCheckResourceAttr(accTestDownloadQcow2FileName, "node_name", accTestNodeName),
		resource.TestCheckResourceAttr(accTestDownloadQcow2FileName, "datastore_id", accTestStorageName),
		resource.TestCheckResourceAttr(accTestDownloadQcow2FileName, "download_url", "https://cdn.githubraw.com/rafsaf/036eece601975a3ad632a77fc2809046/raw/10500012fca9b4425b50de67a7258a12cba0c076/fake_file.qcow2"),
		resource.TestCheckResourceAttr(accTestDownloadQcow2FileName, "filename", "fake_file.qcow2.iso"),
		resource.TestCheckResourceAttr(accTestDownloadQcow2FileName, "allow_unsupported_types", "true"),
		resource.TestCheckResourceAttr(accTestDownloadQcow2FileName, "upload_timeout", "600"),
		resource.TestCheckResourceAttr(accTestDownloadQcow2FileName, "size", "512"),
		resource.TestCheckResourceAttr(accTestDownloadQcow2FileName, "verify", "true"),
		resource.TestCheckNoResourceAttr(accTestDownloadQcow2FileName, "checksum"),
		resource.TestCheckNoResourceAttr(accTestDownloadQcow2FileName, "checksum_algorithm"),
		resource.TestCheckNoResourceAttr(accTestDownloadQcow2FileName, "compression"),
	)
}

func testAccResourceDownloadIsoFileUpdatedConfig() string {
	return fmt.Sprintf(`
	resource "proxmox_virtual_environment_download_file" "iso_image" {
		content_type = "iso"
		node_name    = "%s"
		datastore_id = "%s"
		download_url = "https://cdn.githubraw.com/rafsaf/a4b19ea5e3485f8da6ca4acf46d09650/raw/d340ec3ddcef9b907ede02f64b5d3f694da5d081/fake_file.iso"
		upload_timeout = 10000
	  }
	 `, accTestNodeName, accTestStorageName)
}

func testAccResourceDownloadIsoFileUpdatedCheck() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(accTestDownloadIsoFileName, "content_type", "iso"),
		resource.TestCheckResourceAttr(accTestDownloadIsoFileName, "node_name", accTestNodeName),
		resource.TestCheckResourceAttr(accTestDownloadIsoFileName, "datastore_id", accTestStorageName),
		resource.TestCheckResourceAttr(accTestDownloadIsoFileName, "download_url", "https://cdn.githubraw.com/rafsaf/a4b19ea5e3485f8da6ca4acf46d09650/raw/d340ec3ddcef9b907ede02f64b5d3f694da5d081/fake_file.iso"),
		resource.TestCheckResourceAttr(accTestDownloadIsoFileName, "filename", "fake_file.iso"),
		resource.TestCheckResourceAttr(accTestDownloadIsoFileName, "allow_unsupported_types", "false"),
		resource.TestCheckResourceAttr(accTestDownloadIsoFileName, "upload_timeout", "10000"),
		resource.TestCheckResourceAttr(accTestDownloadIsoFileName, "size", "512"),
		resource.TestCheckResourceAttr(accTestDownloadIsoFileName, "verify", "true"),
		resource.TestCheckNoResourceAttr(accTestDownloadIsoFileName, "checksum"),
		resource.TestCheckNoResourceAttr(accTestDownloadIsoFileName, "checksum_algorithm"),
		resource.TestCheckNoResourceAttr(accTestDownloadIsoFileName, "compression"),
	)
}
