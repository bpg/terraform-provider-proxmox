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
	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/storage"
	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

const (
	fakeFileISO   = "https://cdn.githubraw.com/rafsaf/a4b19ea5e3485f8da6ca4acf46d09650/raw/d340ec3ddcef9b907ede02f64b5d3f694da5d081/fake_file.iso"
	fakeFileQCOW2 = "https://cdn.githubraw.com/rafsaf/036eece601975a3ad632a77fc2809046/raw/10500012fca9b4425b50de67a7258a12cba0c076/fake_file.qcow2"
)

func TestAccResourceDownloadFile(t *testing.T) {
	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"download iso file", []resource.TestStep{{
			Config: fmt.Sprintf(`
				resource "proxmox_virtual_environment_download_file" "iso_image" {
					content_type = "iso"
					node_name    = "%s"
					datastore_id = "%s"
					url          = "%s"
				  }
				 `, accTestNodeName, accTestStorageName, fakeFileISO),
			Check: resource.ComposeTestCheckFunc(
				testResourceAttributes("proxmox_virtual_environment_download_file.iso_image", map[string]string{
					"id":             "local:iso/fake_file.iso",
					"node_name":      accTestNodeName,
					"datastore_id":   accTestStorageName,
					"url":            fakeFileISO,
					"file_name":      "fake_file.iso",
					"upload_timeout": "600",
					"size":           "3",
					"verify":         "true",
				}),
				testNoResourceAttributes("proxmox_virtual_environment_download_file.iso_image", []string{
					"checksum",
					"checksum_algorithm",
					"decompression_algorithm",
				}),
			),
		}}},
		{"download qcow2 file", []resource.TestStep{{
			Config: fmt.Sprintf(`
				resource "proxmox_virtual_environment_download_file" "qcow2_image" {
					content_type       = "iso"
					node_name          = "%s"
					datastore_id       = "%s"
					file_name          = "fake_qcow2_file.img"
					url                =  "%s"
					checksum           = "688787d8ff144c502c7f5cffaafe2cc588d86079f9de88304c26b0cb99ce91c6"
					checksum_algorithm = "sha256"
				  }
				 `, accTestNodeName, accTestStorageName, fakeFileQCOW2),
			Check: resource.ComposeTestCheckFunc(
				testResourceAttributes("proxmox_virtual_environment_download_file.qcow2_image", map[string]string{
					"id":                 "local:iso/fake_qcow2_file.img",
					"content_type":       "iso",
					"node_name":          accTestNodeName,
					"datastore_id":       accTestStorageName,
					"url":                fakeFileQCOW2,
					"file_name":          "fake_qcow2_file.img",
					"upload_timeout":     "600",
					"size":               "3",
					"verify":             "true",
					"checksum":           "688787d8ff144c502c7f5cffaafe2cc588d86079f9de88304c26b0cb99ce91c6",
					"checksum_algorithm": "sha256",
				}),
				testNoResourceAttributes("proxmox_virtual_environment_download_file.qcow2_image", []string{
					"decompression_algorithm",
				}),
			),
		}}},
		{"update file", []resource.TestStep{{
			Config: fmt.Sprintf(`
				resource "proxmox_virtual_environment_download_file" "iso_image" {
					content_type   = "iso"
					node_name      = "%s"
					datastore_id   = "%s"
					file_name      = "fake_iso_file.img"
					url            = "%s"
					upload_timeout = 10000
				  }
				 `, accTestNodeName, accTestStorageName, fakeFileISO),
			Check: resource.ComposeTestCheckFunc(
				testResourceAttributes("proxmox_virtual_environment_download_file.iso_image", map[string]string{
					"id":             "local:iso/fake_iso_file.img",
					"content_type":   "iso",
					"node_name":      accTestNodeName,
					"datastore_id":   accTestStorageName,
					"url":            fakeFileISO,
					"file_name":      "fake_iso_file.img",
					"upload_timeout": "10000",
					"size":           "3",
					"verify":         "true",
				}),
				testNoResourceAttributes("proxmox_virtual_environment_download_file.iso_image", []string{
					"checksum",
					"checksum_algorithm",
					"decompression_algorithm",
				}),
			),
		}}},
		{"override unmanaged file", []resource.TestStep{{
			PreConfig: func() {
				err := getNodeStorageClient().DownloadFileByURL(context.Background(), &storage.DownloadURLPostRequestBody{
					Content:  types.StrPtr("iso"),
					FileName: types.StrPtr("fake_file.iso"),
					Node:     types.StrPtr(accTestNodeName),
					Storage:  types.StrPtr(accTestStorageName),
					URL:      types.StrPtr(fakeFileISO),
				}, 600)
				require.NoError(t, err)
				t.Cleanup(func() {
					err := getNodeStorageClient().DeleteDatastoreFile(context.Background(), "iso/fake_file.iso")
					require.NoError(t, err)
				})
			},
			Config: fmt.Sprintf(`
				resource "proxmox_virtual_environment_download_file" "iso_image" {
					content_type        = "iso"
					node_name           = "%s"
					datastore_id        = "%s"
					url 		        = "%s"
					overwrite_unmanaged = true
				  }
				 `, accTestNodeName, accTestStorageName, fakeFileISO),
			Check: resource.ComposeTestCheckFunc(
				testResourceAttributes("proxmox_virtual_environment_download_file.iso_image", map[string]string{
					"id":           "local:iso/fake_file.iso",
					"content_type": "iso",
					"node_name":    accTestNodeName,
					"datastore_id": accTestStorageName,
					"url":          fakeFileISO,
					"file_name":    "fake_file.iso",
					"size":         "3",
					"verify":       "true",
				}),
				testNoResourceAttributes("proxmox_virtual_environment_download_file.iso_image", []string{
					"checksum",
					"checksum_algorithm",
					"decompression_algorithm",
				}),
			),
		}}},
	}

	accProviders := testAccMuxProviders(context.Background(), t)

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: accProviders,
				Steps:                    tt.steps,
			})
		})
	}
}
