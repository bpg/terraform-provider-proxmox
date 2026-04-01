//go:build acceptance || all

//testacc:tier=light
//testacc:resource=replication

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

func TestAccDataSourceReplications(t *testing.T) {
	te := test.InitEnvironment(t)

	skipReplication(t, te)

	imageFileName := gofakeit.LetterN(8) + "-ubuntu-24.04-standard_24.04-2_amd64.tar.zst"

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
		{"read replication data sources with minimum attributes", func() []resource.TestStep {
			cid1 := newCID()
			cid2 := newCID()
			id1 := fmt.Sprintf("%d-1", cid1)
			id2 := fmt.Sprintf("%d-1", cid2)
			te.AddTemplateVars(map[string]any{
				"TestContainerID1": cid1,
				"TestContainerID2": cid2,
			})
			return []resource.TestStep{{
				Config: te.RenderConfig(`

				resource "proxmox_virtual_environment_container" "test_container1" {
					node_name = "{{.NodeName}}"
					vm_id     = {{.TestContainerID1}}

					disk {
						datastore_id = "{{.ZfsDatastoreID}}"
						size         = 4
						mount_options = []
					}

					operating_system {
						template_file_id = "local:vztmpl/{{ .ImageFileName }}"
						type             = "ubuntu"
					}

					initialization {
						hostname = "test"
					}

					started = false
				}

				resource "proxmox_virtual_environment_container" "test_container2" {
					node_name = "{{.NodeName}}"
					vm_id     = {{.TestContainerID2}}

					disk {
						datastore_id = "{{.ZfsDatastoreID}}"
						size         = 4
						mount_options = []
					}

					operating_system {
						template_file_id = "local:vztmpl/{{ .ImageFileName }}"
						type             = "ubuntu"
					}

					initialization {
						hostname = "test"
					}

					started = false
				}

				resource "proxmox_replication" "test_replication1" {
					id     = "${proxmox_virtual_environment_container.test_container1.id}-1"
					target = "{{.Node2Name}}"
					type = "local"
				}

				resource "proxmox_replication" "test_replication2" {
					id     = "${proxmox_virtual_environment_container.test_container2.id}-1"
					target = "{{.Node2Name}}"
					type = "local"
				}

				data "proxmox_replications" "all" {
							depends_on = [
							proxmox_replication.test_replication1,
							proxmox_replication.test_replication2
						]
				}
					`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.proxmox_replications.all", "replications.#"),
					resource.TestCheckTypeSetElemNestedAttrs("data.proxmox_replications.all", "replications.*", map[string]string{
						"id":     id1,
						"target": te.Node2Name,
						"type":   "local",
					}),
					resource.TestCheckTypeSetElemNestedAttrs("data.proxmox_replications.all", "replications.*", map[string]string{
						"id":     id2,
						"target": te.Node2Name,
						"type":   "local",
					}),
				),
			}}
		}()},
		{
			"read replication data sources with all attributes", func() []resource.TestStep {
				cid1 := newCID()
				cid2 := newCID()
				id1 := fmt.Sprintf("%d-1", cid1)
				id2 := fmt.Sprintf("%d-1", cid2)
				guest1 := fmt.Sprintf("%d", cid1)
				guest2 := fmt.Sprintf("%d", cid2)
				te.AddTemplateVars(map[string]any{
					"TestContainerID1": cid1,
					"TestContainerID2": cid2,
				})
				return []resource.TestStep{{
					Config: te.RenderConfig(`

				resource "proxmox_virtual_environment_container" "test_container1" {
					node_name = "{{.NodeName}}"
					vm_id     = {{.TestContainerID1}}

					disk {
						datastore_id = "{{.ZfsDatastoreID}}"
						size         = 4
						mount_options = []
					}

					operating_system {
						template_file_id = "local:vztmpl/{{ .ImageFileName }}"
						type             = "ubuntu"
					}

					initialization {
						hostname = "test"
					}

					started = false
				}

				resource "proxmox_virtual_environment_container" "test_container2" {
					node_name = "{{.NodeName}}"
					vm_id     = {{.TestContainerID2}}

					disk {
						datastore_id = "{{.ZfsDatastoreID}}"
						size         = 4
						mount_options = []
					}

					operating_system {
						template_file_id = "local:vztmpl/{{ .ImageFileName }}"
						type             = "ubuntu"
					}

					initialization {
						hostname = "test"
					}

					started = false
				}

				resource "proxmox_replication" "test_replication1" {
					id     = "${proxmox_virtual_environment_container.test_container1.id}-1"
					target = "{{.Node2Name}}"
					type = "local"
					disable = true
					comment = "comment 123"
					schedule = "*/30"
					rate = 10
				}

				resource "proxmox_replication" "test_replication2" {
					id     = "${proxmox_virtual_environment_container.test_container2.id}-1"
					target = "{{.Node2Name}}"
					type = "local"
					disable = true
					comment = "comment 123"
					schedule = "*/30"
					rate = 10
				}

				data "proxmox_replications" "all" {
							depends_on = [
							proxmox_replication.test_replication1,
							proxmox_replication.test_replication2
						]
				}
					`),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttrSet("data.proxmox_replications.all", "replications.#"),
						resource.TestCheckTypeSetElemNestedAttrs("data.proxmox_replications.all", "replications.*", map[string]string{
							"id":       id1,
							"target":   te.Node2Name,
							"type":     "local",
							"jobnum":   jobnum,
							"guest":    guest1,
							"disable":  "true",
							"comment":  "comment 123",
							"schedule": "*/30",
							"rate":     "10",
						}),
						resource.TestCheckTypeSetElemNestedAttrs("data.proxmox_replications.all", "replications.*", map[string]string{
							"id":       id2,
							"target":   te.Node2Name,
							"type":     "local",
							"jobnum":   jobnum,
							"guest":    guest2,
							"disable":  "true",
							"comment":  "comment 123",
							"schedule": "*/30",
							"rate":     "10",
						}),
					),
				}}
			}(),
		},
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
