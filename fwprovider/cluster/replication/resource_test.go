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
	"math/rand"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/test"
	"github.com/bpg/terraform-provider-proxmox/proxmox/helpers/ptr"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/storage"
)

func newCID() int {
	return 100000 + rand.Intn(99999)
}

func renderConfigWithCT(te *test.Environment, cid int, replicationBlock string) string {
	te.AddTemplateVars(map[string]any{
		"TestContainerID": cid,
	})

	return te.RenderConfig(fmt.Sprintf(`
        resource "proxmox_virtual_environment_container" "test_container" {
            node_name = "{{.NodeName}}"
            vm_id     = {{.TestContainerID}}

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

        %s
    `, replicationBlock))
}

func TestAccResourceReplication(t *testing.T) {
	te := test.InitEnvironment(t)

	imageFileName := gofakeit.Word() + "-ubuntu-24.04-standard_24.04-2_amd64.tar.zst"

	target := "pve02"
	jobnum := "1"

	te.AddTemplateVars(map[string]any{
		"Target":        target,
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
		{
			"create and update minimal replication", []resource.TestStep{
				func() resource.TestStep {
					cid := newCID()
					id := fmt.Sprintf("%d-%s", cid, jobnum)
					guest := fmt.Sprintf("%d", cid)
					return resource.TestStep{
						Config: renderConfigWithCT(te, cid, `

			resource "proxmox_virtual_environment_replication" "test_replication" {
				id     = "${proxmox_virtual_environment_container.test_container.id}-{{.JobNum}}"
				target = "{{.Target}}"
				type = "local"
			}`),
						Check: resource.ComposeTestCheckFunc(
							test.ResourceAttributes("proxmox_virtual_environment_replication.test_replication", map[string]string{
								"id":     id,
								"target": target,
								"type":   "local",
								"jobnum": jobnum,
								"guest":  guest,
							}),
						),
					}
				}(),
				func() resource.TestStep {
					cid := newCID()
					id := fmt.Sprintf("%d-%s", cid, jobnum)
					guest := fmt.Sprintf("%d", cid)
					return resource.TestStep{
						Config: renderConfigWithCT(te, cid, `

			resource "proxmox_virtual_environment_replication" "test_replication" {
				id     = "${proxmox_virtual_environment_container.test_container.id}-{{.JobNum}}"
				target = "{{.Target}}"
				type = "local"
				disable = true
			}
				`),
						Check: resource.ComposeTestCheckFunc(
							test.ResourceAttributes("proxmox_virtual_environment_replication.test_replication", map[string]string{
								"id":      id,
								"target":  target,
								"type":    "local",
								"jobnum":  jobnum,
								"guest":   guest,
								"disable": "true",
							}),
						),
					}
				}(),
				func() resource.TestStep {
					cid := newCID()
					id := fmt.Sprintf("%d-%s", cid, jobnum)
					guest := fmt.Sprintf("%d", cid)
					return resource.TestStep{
						Config: renderConfigWithCT(te, cid, `

				resource "proxmox_virtual_environment_replication" "test_replication" {
					id     = "${proxmox_virtual_environment_container.test_container.id}-{{.JobNum}}"
					target = "{{.Target}}"
					type = "local"
					disable = true
					comment = "comment 123"
				}
					`),
						Check: resource.ComposeTestCheckFunc(
							test.ResourceAttributes("proxmox_virtual_environment_replication.test_replication", map[string]string{
								"id":      id,
								"target":  target,
								"type":    "local",
								"jobnum":  jobnum,
								"guest":   guest,
								"disable": "true",
								"comment": "comment 123",
							}),
						),
					}
				}(),
				func() resource.TestStep {
					cid := newCID()
					id := fmt.Sprintf("%d-%s", cid, jobnum)
					guest := fmt.Sprintf("%d", cid)
					return resource.TestStep{
						Config: renderConfigWithCT(te, cid, `

				resource "proxmox_virtual_environment_replication" "test_replication" {
					id     = "${proxmox_virtual_environment_container.test_container.id}-{{.JobNum}}"
					target = "{{.Target}}"
					type = "local"
					disable = true
					comment = "comment 123"
					schedule = "*/30"
				}
					`),
						Check: resource.ComposeTestCheckFunc(
							test.ResourceAttributes("proxmox_virtual_environment_replication.test_replication", map[string]string{
								"id":       id,
								"target":   target,
								"type":     "local",
								"jobnum":   jobnum,
								"guest":    guest,
								"disable":  "true",
								"comment":  "comment 123",
								"schedule": "^*/30$",
							}),
						),
					}
				}(),
				func() resource.TestStep {
					cid := newCID()
					id := fmt.Sprintf("%d-%s", cid, jobnum)
					guest := fmt.Sprintf("%d", cid)
					return resource.TestStep{
						Config: renderConfigWithCT(te, cid, `

				resource "proxmox_virtual_environment_replication" "test_replication" {
					id     = "${proxmox_virtual_environment_container.test_container.id}-{{.JobNum}}"
					target = "{{.Target}}"
					type = "local"
					disable = true
					comment = "comment 123"
					schedule = "*/30"
					rate = 10
				}
					`),
						Check: resource.ComposeTestCheckFunc(
							test.ResourceAttributes("proxmox_virtual_environment_replication.test_replication", map[string]string{
								"id":       id,
								"target":   target,
								"type":     "local",
								"jobnum":   jobnum,
								"guest":    guest,
								"disable":  "true",
								"comment":  "comment 123",
								"schedule": "^*/30$",
								"rate":     "10",
							}),
						),
						ResourceName:      "proxmox_virtual_environment_replication.test_replication",
						ImportState:       true,
						ImportStateVerify: true,
					}
				}(),
			},
		},
		{"create disabled replication", []resource.TestStep{
			func() resource.TestStep {
				cid := newCID()
				id := fmt.Sprintf("%d-%s", cid, jobnum)
				guest := fmt.Sprintf("%d", cid)
				return resource.TestStep{
					Config: renderConfigWithCT(te, cid, `

				resource "proxmox_virtual_environment_replication" "test_replication" {
					id     = "${proxmox_virtual_environment_container.test_container.id}-{{.JobNum}}"
					target = "{{.Target}}"
					type = "local"
					disable = true
				}
					`),
					Check: resource.ComposeTestCheckFunc(
						test.ResourceAttributes("proxmox_virtual_environment_replication.test_replication", map[string]string{
							"id":      id,
							"target":  target,
							"type":    "local",
							"jobnum":  jobnum,
							"guest":   guest,
							"disable": "true",
						}),
					),
				}
			}(),
		}},
		{
			"create schduled replication", []resource.TestStep{
				func() resource.TestStep {
					cid := newCID()
					id := fmt.Sprintf("%d-%s", cid, jobnum)
					guest := fmt.Sprintf("%d", cid)
					return resource.TestStep{
						Config: renderConfigWithCT(te, cid, `

				resource "proxmox_virtual_environment_replication" "test_replication" {
					id     = "${proxmox_virtual_environment_container.test_container.id}-{{.JobNum}}"
					target = "{{.Target}}"
					type = "local"
					schedule = "*/30"
				}
					`),
						Check: resource.ComposeTestCheckFunc(
							test.ResourceAttributes("proxmox_virtual_environment_replication.test_replication", map[string]string{
								"id":       id,
								"target":   target,
								"type":     "local",
								"jobnum":   jobnum,
								"guest":    guest,
								"schedule": "^*/30$",
							}),
						),
					}
				}(),
			},
		},
		{"create rated replication", []resource.TestStep{
			func() resource.TestStep {
				cid := newCID()
				id := fmt.Sprintf("%d-%s", cid, jobnum)
				guest := fmt.Sprintf("%d", cid)
				return resource.TestStep{
					Config: renderConfigWithCT(te, cid, `

				resource "proxmox_virtual_environment_replication" "test_replication" {
					id     = "${proxmox_virtual_environment_container.test_container.id}-{{.JobNum}}"
					target = "{{.Target}}"
					type = "local"
					rate = 10
				}
					`),
					Check: resource.ComposeTestCheckFunc(
						test.ResourceAttributes("proxmox_virtual_environment_replication.test_replication", map[string]string{
							"id":     id,
							"target": target,
							"type":   "local",
							"jobnum": jobnum,
							"guest":  guest,
							"rate":   "10",
						}),
					),
				}
			}(),
		}},
		{"create commented replication", []resource.TestStep{
			func() resource.TestStep {
				cid := newCID()
				id := fmt.Sprintf("%d-%s", cid, jobnum)
				guest := fmt.Sprintf("%d", cid)
				return resource.TestStep{
					Config: renderConfigWithCT(te, cid, `

				resource "proxmox_virtual_environment_replication" "test_replication" {
					id     = "${proxmox_virtual_environment_container.test_container.id}-{{.JobNum}}"
					target = "{{.Target}}"
					type = "local"
					comment = "comment 123"
				}
					`),
					Check: resource.ComposeTestCheckFunc(
						test.ResourceAttributes("proxmox_virtual_environment_replication.test_replication", map[string]string{
							"id":      id,
							"target":  target,
							"type":    "local",
							"jobnum":  jobnum,
							"guest":   guest,
							"comment": "comment 123",
						}),
					),
				}
			}(),
		}},
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
