//go:build acceptance || all

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package nodes_test

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/test"
	"github.com/bpg/terraform-provider-proxmox/proxmox/helpers/ptr"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/storage"
)

const (
	testOCIImage       = "docker.io/library/hello-world:latest"
	testOCIImageAlpine = "docker.io/library/alpine:latest"
)

func TestAccResourceOCIImage(t *testing.T) {
	te := test.InitEnvironment(t)

	te.AddTemplateVars(map[string]interface{}{
		"OCIImage": testOCIImage,
	})

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"missing reference", []resource.TestStep{{
			Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_oci_image" "test_image" {
					node_name    = "{{.NodeName}}"
					datastore_id = "{{.DatastoreID}}"
					reference    = ""
				}`),
			ExpectError: regexp.MustCompile(`Attribute reference must match OCI image reference regex`),
		}}},
		{"pull OCI image without file_name", []resource.TestStep{{
			Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_oci_image" "test_image" {
					node_name          = "{{.NodeName}}"
					datastore_id       = "{{.DatastoreID}}"
					reference          = "{{.OCIImage}}"
					overwrite_unmanaged = true
				}`),
			Check: resource.ComposeTestCheckFunc(
				test.ResourceAttributes("proxmox_virtual_environment_oci_image.test_image", map[string]string{
					"id":                  "local:vztmpl/hello-world_latest.tar",
					"node_name":           te.NodeName,
					"datastore_id":        te.DatastoreID,
					"reference":           testOCIImage,
					"file_name":           "hello-world_latest.tar",
					"upload_timeout":      "600",
					"overwrite":           "true",
					"overwrite_unmanaged": "true",
				}),
				test.ResourceAttributesSet("proxmox_virtual_environment_oci_image.test_image", []string{
					"size",
				}),
			),
		}}},
		{"pull OCI image with custom file_name", []resource.TestStep{{
			Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_oci_image" "test_custom_name" {
					node_name          = "{{.NodeName}}"
					datastore_id       = "{{.DatastoreID}}"
					reference          = "{{.OCIImage}}"
					file_name          = "custom_hello_world.tar"
					overwrite_unmanaged = true
				}`),
			Check: resource.ComposeTestCheckFunc(
				test.ResourceAttributes("proxmox_virtual_environment_oci_image.test_custom_name", map[string]string{
					"id":                  "local:vztmpl/custom_hello_world.tar",
					"node_name":           te.NodeName,
					"datastore_id":        te.DatastoreID,
					"reference":           testOCIImage,
					"file_name":           "custom_hello_world.tar",
					"upload_timeout":      "600",
					"overwrite":           "true",
					"overwrite_unmanaged": "true",
				}),
				test.ResourceAttributesSet("proxmox_virtual_environment_oci_image.test_custom_name", []string{
					"size",
				}),
			),
		}}},
		{"invalid file_name without .tar extension", []resource.TestStep{{
			Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_oci_image" "test_invalid_name" {
					node_name          = "{{.NodeName}}"
					datastore_id       = "{{.DatastoreID}}"
					reference          = "{{.OCIImage}}"
					file_name          = "invalid_name.txt"
				}`),
			ExpectError: regexp.MustCompile(`file name must end with .tar`),
		}}},
		{"pull & update OCI image", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_oci_image" "test_update" {
					node_name          = "{{.NodeName}}"
					datastore_id       = "{{.DatastoreID}}"
					reference          = "{{.OCIImage}}"
					file_name          = "test_update_image.tar"
					overwrite_unmanaged = true
				}`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_virtual_environment_oci_image.test_update", map[string]string{
						"id":             "local:vztmpl/test_update_image.tar",
						"node_name":      te.NodeName,
						"datastore_id":   te.DatastoreID,
						"reference":      testOCIImage,
						"file_name":      "test_update_image.tar",
						"upload_timeout": "600",
					}),
					test.ResourceAttributesSet("proxmox_virtual_environment_oci_image.test_update", []string{
						"size",
					}),
				),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_oci_image" "test_update" {
					node_name          = "{{.NodeName}}"
					datastore_id       = "{{.DatastoreID}}"
					reference          = "{{.OCIImage}}"
					file_name          = "test_update_image.tar"
					upload_timeout     = 1200
					overwrite_unmanaged = true
				}`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_virtual_environment_oci_image.test_update", map[string]string{
						"id":             "local:vztmpl/test_update_image.tar",
						"node_name":      te.NodeName,
						"datastore_id":   te.DatastoreID,
						"reference":      testOCIImage,
						"file_name":      "test_update_image.tar",
						"upload_timeout": "1200",
					}),
					test.ResourceAttributesSet("proxmox_virtual_environment_oci_image.test_update", []string{
						"size",
					}),
				),
			},
		}},
		{"override behavior", []resource.TestStep{{
			Destroy: false,
			PreConfig: func() {
				ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
				defer cancel()

				fileID := "vztmpl/test_override_image.tar"
				_ = te.NodeStorageClient().DeleteDatastoreFile(ctx, fileID) //nolint: errcheck

				// Pull OCI image outside of Terraform
				filenameWithoutTar := "test_override_image"
				err := te.NodeStorageClient().DownloadOCIImageByReference(ctx, &storage.OCIRegistryPullRequestBody{
					Storage:   ptr.Ptr(te.DatastoreID),
					FileName:  &filenameWithoutTar,
					Reference: ptr.Ptr(testOCIImage),
				})
				require.NoError(t, err)

				t.Cleanup(func() {
					e := te.NodeStorageClient().DeleteDatastoreFile(context.Background(), fileID)
					require.NoError(t, e)
				})
			},
			Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_oci_image" "test_override" {
					node_name           = "{{.NodeName}}"
					datastore_id        = "{{.DatastoreID}}"
					reference           = "{{.OCIImage}}"
					file_name           = "test_override_another.tar"
					overwrite_unmanaged = true
					overwrite           = false
				}`),
			Check: resource.ComposeTestCheckFunc(
				test.ResourceAttributes("proxmox_virtual_environment_oci_image.test_override", map[string]string{
					"id":                  "local:vztmpl/test_override_another.tar",
					"node_name":           te.NodeName,
					"datastore_id":        te.DatastoreID,
					"reference":           testOCIImage,
					"file_name":           "test_override_another.tar",
					"overwrite":           "false",
					"overwrite_unmanaged": "true",
				}),
				test.ResourceAttributesSet("proxmox_virtual_environment_oci_image.test_override", []string{
					"size",
				}),
			),
		}, {
			Destroy: false,
			PreConfig: func() {
				ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
				defer cancel()

				fileID := "vztmpl/test_override_another.tar"

				// Manually pull a different version to simulate external change
				filenameWithoutTar := "test_override_another"
				_ = te.NodeStorageClient().DeleteDatastoreFile(ctx, fileID) //nolint: errcheck

				err := te.NodeStorageClient().DownloadOCIImageByReference(ctx, &storage.OCIRegistryPullRequestBody{
					Storage:   ptr.Ptr(te.DatastoreID),
					FileName:  &filenameWithoutTar,
					Reference: ptr.Ptr(testOCIImage),
				})
				require.NoError(t, err)
			},
			Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_oci_image" "test_override" {
					node_name           = "{{.NodeName}}"
					datastore_id        = "{{.DatastoreID}}"
					reference           = "{{.OCIImage}}"
					file_name           = "test_override_another.tar"
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
				ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
				defer cancel()

				fileID := "vztmpl/test_override_another.tar"

				// Manually pull a different image to simulate external change with different size
				filenameWithoutTar := "test_override_another"
				_ = te.NodeStorageClient().DeleteDatastoreFile(ctx, fileID) //nolint: errcheck

				err := te.NodeStorageClient().DownloadOCIImageByReference(ctx, &storage.OCIRegistryPullRequestBody{
					Storage:   ptr.Ptr(te.DatastoreID),
					FileName:  &filenameWithoutTar,
					Reference: ptr.Ptr(testOCIImageAlpine),
				})
				require.NoError(t, err)
			},
			Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_oci_image" "test_override" {
					node_name           = "{{.NodeName}}"
					datastore_id        = "{{.DatastoreID}}"
					reference           = "{{.OCIImage}}"
					file_name           = "test_override_another.tar"
					overwrite_unmanaged = true
					overwrite           = true
				}`),
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction("proxmox_virtual_environment_oci_image.test_override", plancheck.ResourceActionDestroyBeforeCreate),
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
