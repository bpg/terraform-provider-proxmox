//go:build acceptance || all

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package replication_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/test"
	"github.com/bpg/terraform-provider-proxmox/proxmox/helpers/ptr"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/storage"
)

func TestAccDataSourceReplication(t *testing.T) {
	te := test.InitEnvironment(t)

	skipReplication(t, te)

	imageFileName := gofakeit.Word() + "-ubuntu-24.04-standard_24.04-2_amd64.tar.zst"

	jobnum := "1"

	te.AddTemplateVars(map[string]any{
		"JobNum":        jobnum,
		"ImageFileName": imageFileName,
	})

	err := te.NodeStorageClient().DownloadFileByURL(context.Background(), &storage.DownloadURLPostRequestBody{
		Content:  ptr.Ptr("vztmpl"),
		FileName: &imageFileName,
		Node:     &te.NodeName,
		Storage:  &te.DatastoreID,
		URL:      ptr.Ptr(fmt.Sprintf("%s/images/system/ubuntu-24.04-standard_24.04-2_amd64.tar.zst", te.ContainerImagesServer)),
	})
	require.NoError(t, err)

	t.Cleanup(func() {
		e := te.NodeStorageClient().DeleteDatastoreFile(context.Background(), fmt.Sprintf("vztmpl/%s", imageFileName))
		require.NoError(t, e)
	})

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"read replication data source with minimum attributes", []resource.TestStep{
			func() resource.TestStep {
				cid := newCID()
				id := fmt.Sprintf("%d-%s", cid, jobnum)
				guest := fmt.Sprintf("%d", cid)
				return resource.TestStep{
					Config: renderConfigWithCT(te, cid, `
				
			resource "proxmox_virtual_environment_replication" "test_replication" {
				id     = "${proxmox_virtual_environment_container.test_container.id}-{{.JobNum}}"
				target = "{{.Node2Name}}"
				type = "local"
			}

			data "proxmox_virtual_environment_replication" "test_replication_data" {
				id = proxmox_virtual_environment_replication.test_replication.id
			}
				`),
					Check: resource.ComposeTestCheckFunc(
						test.ResourceAttributes("data.proxmox_virtual_environment_replication.test_replication_data", map[string]string{
							"id":     id,
							"target": te.Node2Name,
							"type":   "local",
							"jobnum": jobnum,
							"guest":  guest,
						}),
					),
				}
			}(),
		}},
		{"read replication data source with all attributes", []resource.TestStep{
			func() resource.TestStep {
				cid := newCID()
				id := fmt.Sprintf("%d-%s", cid, jobnum)
				guest := fmt.Sprintf("%d", cid)
				return resource.TestStep{
					Config: renderConfigWithCT(te, cid, `
				
			resource "proxmox_virtual_environment_replication" "test_replication" {
				id     = "${proxmox_virtual_environment_container.test_container.id}-{{.JobNum}}"
				target = "{{.Node2Name}}"
				type = "local"
				disable = true
				comment = "comment 123"
				schedule = "*/30"
				rate = 10
			}

			data "proxmox_virtual_environment_replication" "test_replication_data" {
				id = proxmox_virtual_environment_replication.test_replication.id
			}
				`),
					Check: resource.ComposeTestCheckFunc(
						test.ResourceAttributes("data.proxmox_virtual_environment_replication.test_replication_data", map[string]string{
							"id":       id,
							"target":   te.Node2Name,
							"type":     "local",
							"jobnum":   jobnum,
							"guest":    guest,
							"disable":  "true",
							"comment":  "comment 123",
							"schedule": `^\*/30$`,
							"rate":     "10",
						}),
					),
				}
			}(),
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: te.AccProviders,
				Steps:                    tt.steps,
			})
		})
	}
}
