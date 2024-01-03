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
		url = "https://cdn.githubraw.com/rafsaf/a4b19ea5e3485f8da6ca4acf46d09650/raw/d340ec3ddcef9b907ede02f64b5d3f694da5d081/fake_file.iso"
	  }
	 `, accTestNodeName, accTestStorageName)
}

func testAccResourceDownloadQcow2FileCreatedConfig() string {
	return fmt.Sprintf(`
	resource "proxmox_virtual_environment_download_file" "qcow2_image" {
		content_type       = "iso"
		node_name          = "%s"
		datastore_id       = "%s"
		file_name          = "fake_qcow2_file.img"
		url                = "https://cdn.githubraw.com/rafsaf/036eece601975a3ad632a77fc2809046/raw/10500012fca9b4425b50de67a7258a12cba0c076/fake_file.qcow2"
		checksum           = "688787d8ff144c502c7f5cffaafe2cc588d86079f9de88304c26b0cb99ce91c6"
		checksum_algorithm = "sha256"
	  }
	 `, accTestNodeName, accTestStorageName)
}

func testAccResourceDownloadIsoFileCreatedCheck() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(accTestDownloadIsoFileName, "id", "local:iso/fake_file.iso"),
		resource.TestCheckResourceAttr(accTestDownloadIsoFileName, "node_name", accTestNodeName),
		resource.TestCheckResourceAttr(accTestDownloadIsoFileName, "datastore_id", accTestStorageName),
		resource.TestCheckResourceAttr(accTestDownloadIsoFileName, "url", "https://cdn.githubraw.com/rafsaf/a4b19ea5e3485f8da6ca4acf46d09650/raw/d340ec3ddcef9b907ede02f64b5d3f694da5d081/fake_file.iso"),
		resource.TestCheckResourceAttr(accTestDownloadIsoFileName, "file_name", "fake_file.iso"),
		resource.TestCheckResourceAttr(accTestDownloadIsoFileName, "upload_timeout", "600"),
		resource.TestCheckResourceAttr(accTestDownloadIsoFileName, "size", "3"),
		resource.TestCheckResourceAttr(accTestDownloadIsoFileName, "verify", "true"),
		resource.TestCheckNoResourceAttr(accTestDownloadIsoFileName, "checksum"),
		resource.TestCheckNoResourceAttr(accTestDownloadIsoFileName, "checksum_algorithm"),
		resource.TestCheckNoResourceAttr(accTestDownloadIsoFileName, "decompression_algorithm"),
	)
}

func testAccResourceDownloadQcow2FileCreatedCheck() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(accTestDownloadQcow2FileName, "id", "local:iso/fake_qcow2_file.img"),
		resource.TestCheckResourceAttr(accTestDownloadQcow2FileName, "content_type", "iso"),
		resource.TestCheckResourceAttr(accTestDownloadQcow2FileName, "node_name", accTestNodeName),
		resource.TestCheckResourceAttr(accTestDownloadQcow2FileName, "datastore_id", accTestStorageName),
		resource.TestCheckResourceAttr(accTestDownloadQcow2FileName, "url", "https://cdn.githubraw.com/rafsaf/036eece601975a3ad632a77fc2809046/raw/10500012fca9b4425b50de67a7258a12cba0c076/fake_file.qcow2"),
		resource.TestCheckResourceAttr(accTestDownloadQcow2FileName, "file_name", "fake_qcow2_file.img"),
		resource.TestCheckResourceAttr(accTestDownloadQcow2FileName, "upload_timeout", "600"),
		resource.TestCheckResourceAttr(accTestDownloadQcow2FileName, "size", "3"),
		resource.TestCheckResourceAttr(accTestDownloadQcow2FileName, "verify", "true"),
		resource.TestCheckResourceAttr(accTestDownloadQcow2FileName, "checksum", "688787d8ff144c502c7f5cffaafe2cc588d86079f9de88304c26b0cb99ce91c6"),
		resource.TestCheckResourceAttr(accTestDownloadQcow2FileName, "checksum_algorithm", "sha256"),
		resource.TestCheckNoResourceAttr(accTestDownloadQcow2FileName, "decompression_algorithm"),
	)
}

func testAccResourceDownloadIsoFileUpdatedConfig() string {
	return fmt.Sprintf(`
	resource "proxmox_virtual_environment_download_file" "iso_image" {
		content_type = "iso"
		node_name    = "%s"
		datastore_id = "%s"
		file_name    = "fake_iso_file.img"
		url = "https://cdn.githubraw.com/rafsaf/a4b19ea5e3485f8da6ca4acf46d09650/raw/d340ec3ddcef9b907ede02f64b5d3f694da5d081/fake_file.iso"
		upload_timeout = 10000
	  }
	 `, accTestNodeName, accTestStorageName)
}

func testAccResourceDownloadIsoFileUpdatedCheck() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(accTestDownloadIsoFileName, "id", "local:iso/fake_iso_file.img"),
		resource.TestCheckResourceAttr(accTestDownloadIsoFileName, "content_type", "iso"),
		resource.TestCheckResourceAttr(accTestDownloadIsoFileName, "node_name", accTestNodeName),
		resource.TestCheckResourceAttr(accTestDownloadIsoFileName, "datastore_id", accTestStorageName),
		resource.TestCheckResourceAttr(accTestDownloadIsoFileName, "url", "https://cdn.githubraw.com/rafsaf/a4b19ea5e3485f8da6ca4acf46d09650/raw/d340ec3ddcef9b907ede02f64b5d3f694da5d081/fake_file.iso"),
		resource.TestCheckResourceAttr(accTestDownloadIsoFileName, "file_name", "fake_iso_file.img"),
		resource.TestCheckResourceAttr(accTestDownloadIsoFileName, "upload_timeout", "10000"),
		resource.TestCheckResourceAttr(accTestDownloadIsoFileName, "size", "3"),
		resource.TestCheckResourceAttr(accTestDownloadIsoFileName, "verify", "true"),
		resource.TestCheckNoResourceAttr(accTestDownloadIsoFileName, "checksum"),
		resource.TestCheckNoResourceAttr(accTestDownloadIsoFileName, "checksum_algorithm"),
		resource.TestCheckNoResourceAttr(accTestDownloadIsoFileName, "decompression_algorithm"),
	)
}
