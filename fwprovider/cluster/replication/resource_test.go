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
	"math/rand"
	"regexp"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/test"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/storage"
)

// skipReplication validates if replication requirements are met for test t
// otherwise skip t
func skipReplication(t *testing.T, te *test.Environment) {
	if te.Node2Name == "" {
		t.Skip("Skipping replication test, Node2Name not defined")
	}

	if te.ZfsDatastoreID == "" {
		t.Skip("Skipping replication test, ZfsDatastoreID not defined")
	}
}

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

	skipReplication(t, te)

	imageFileName := gofakeit.LetterN(8) + "-ubuntu-24.04-standard_24.04-2_amd64.tar.zst"

	jobnum := "1"

	te.AddTemplateVars(map[string]any{
		"JobNum":        jobnum,
		"ImageFileName": imageFileName,
	})

	err := te.NodeStorageClient().DownloadFileByURL(context.Background(), &storage.DownloadURLPostRequestBody{
		Content:  new("vztmpl"),
		FileName: &imageFileName,
		Node:     &te.NodeName,
		Storage:  &te.DatastoreID,
		URL:      new(fmt.Sprintf("%s/images/system/ubuntu-24.04-standard_24.04-2_amd64.tar.zst", te.ContainerImagesServer)),
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
			"create and update minimal replication", func() []resource.TestStep {
				cid := newCID()
				id := fmt.Sprintf("%d-%s", cid, jobnum)
				guest := fmt.Sprintf("%d", cid)
				return []resource.TestStep{
					{
						Config: renderConfigWithCT(te, cid, `

						resource "proxmox_replication" "test_replication" {
							id     = "${proxmox_virtual_environment_container.test_container.id}-{{.JobNum}}"
							target = "{{.Node2Name}}"
							type = "local"
						}`),
						Check: resource.ComposeTestCheckFunc(
							test.ResourceAttributes("proxmox_replication.test_replication", map[string]string{
								"id":     id,
								"target": te.Node2Name,
								"type":   "local",
								"jobnum": jobnum,
								"guest":  guest,
							}),
						),
					},
					{
						Config: renderConfigWithCT(te, cid, `
						resource "proxmox_replication" "test_replication" {
							id     = "${proxmox_virtual_environment_container.test_container.id}-{{.JobNum}}"
							target = "{{.Node2Name}}"
							type = "local"
							disable = true
						}
						`),
						Check: resource.ComposeTestCheckFunc(
							test.ResourceAttributes("proxmox_replication.test_replication", map[string]string{
								"id":      id,
								"target":  te.Node2Name,
								"type":    "local",
								"jobnum":  jobnum,
								"guest":   guest,
								"disable": "true",
							}),
						),
					},
					{
						Config: renderConfigWithCT(te, cid, `
						resource "proxmox_replication" "test_replication" {
							id     = "${proxmox_virtual_environment_container.test_container.id}-{{.JobNum}}"
							target = "{{.Node2Name}}"
							type = "local"
							disable = true
							comment = "comment 123"
						}
						`),
						Check: resource.ComposeTestCheckFunc(
							test.ResourceAttributes("proxmox_replication.test_replication", map[string]string{
								"id":      id,
								"target":  te.Node2Name,
								"type":    "local",
								"jobnum":  jobnum,
								"guest":   guest,
								"disable": "true",
								"comment": "comment 123",
							}),
						),
					},
					{
						Config: renderConfigWithCT(te, cid, `

						resource "proxmox_replication" "test_replication" {
							id     = "${proxmox_virtual_environment_container.test_container.id}-{{.JobNum}}"
							target = "{{.Node2Name}}"
							type = "local"
							disable = true
							comment = "comment 123"
							schedule = "*/30"
						}
						`),
						Check: resource.ComposeTestCheckFunc(
							test.ResourceAttributes("proxmox_replication.test_replication", map[string]string{
								"id":       id,
								"target":   te.Node2Name,
								"type":     "local",
								"jobnum":   jobnum,
								"guest":    guest,
								"disable":  "true",
								"comment":  "comment 123",
								"schedule": `^\*/30$`,
							}),
						),
					},
					{
						Config: renderConfigWithCT(te, cid, `
						resource "proxmox_replication" "test_replication" {
							id     = "${proxmox_virtual_environment_container.test_container.id}-{{.JobNum}}"
							target = "{{.Node2Name}}"
							type = "local"
							disable = true
							comment = "comment 123"
							schedule = "*/30"
							rate = 10
						}
						`),
						Check: resource.ComposeTestCheckFunc(
							test.ResourceAttributes("proxmox_replication.test_replication", map[string]string{
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
						ResourceName:      "proxmox_replication.test_replication",
						ImportState:       true,
						ImportStateVerify: true,
					},
				}
			}(),
		},
		{"create disabled replication", []resource.TestStep{
			func() resource.TestStep {
				cid := newCID()
				id := fmt.Sprintf("%d-%s", cid, jobnum)
				guest := fmt.Sprintf("%d", cid)
				return resource.TestStep{
					Config: renderConfigWithCT(te, cid, `

				resource "proxmox_replication" "test_replication" {
					id     = "${proxmox_virtual_environment_container.test_container.id}-{{.JobNum}}"
					target = "{{.Node2Name}}"
					type = "local"
					disable = true
				}
					`),
					Check: resource.ComposeTestCheckFunc(
						test.ResourceAttributes("proxmox_replication.test_replication", map[string]string{
							"id":      id,
							"target":  te.Node2Name,
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
			"create scheduled replication", []resource.TestStep{
				func() resource.TestStep {
					cid := newCID()
					id := fmt.Sprintf("%d-%s", cid, jobnum)
					guest := fmt.Sprintf("%d", cid)
					return resource.TestStep{
						Config: renderConfigWithCT(te, cid, `

				resource "proxmox_replication" "test_replication" {
					id     = "${proxmox_virtual_environment_container.test_container.id}-{{.JobNum}}"
					target = "{{.Node2Name}}"
					type = "local"
					schedule = "*/30"
				}
					`),
						Check: resource.ComposeTestCheckFunc(
							test.ResourceAttributes("proxmox_replication.test_replication", map[string]string{
								"id":       id,
								"target":   te.Node2Name,
								"type":     "local",
								"jobnum":   jobnum,
								"guest":    guest,
								"schedule": `^\*/30$`,
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

				resource "proxmox_replication" "test_replication" {
					id     = "${proxmox_virtual_environment_container.test_container.id}-{{.JobNum}}"
					target = "{{.Node2Name}}"
					type = "local"
					rate = 10
				}
					`),
					Check: resource.ComposeTestCheckFunc(
						test.ResourceAttributes("proxmox_replication.test_replication", map[string]string{
							"id":     id,
							"target": te.Node2Name,
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

				resource "proxmox_replication" "test_replication" {
					id     = "${proxmox_virtual_environment_container.test_container.id}-{{.JobNum}}"
					target = "{{.Node2Name}}"
					type = "local"
					comment = "comment 123"
				}
					`),
					Check: resource.ComposeTestCheckFunc(
						test.ResourceAttributes("proxmox_replication.test_replication", map[string]string{
							"id":      id,
							"target":  te.Node2Name,
							"type":    "local",
							"jobnum":  jobnum,
							"guest":   guest,
							"comment": "comment 123",
						}),
					),
				}
			}(),
		}},
		{
			"replication fields deletion", func() []resource.TestStep {
				cid := newCID()
				id := fmt.Sprintf("%d-%s", cid, jobnum)
				guest := fmt.Sprintf("%d", cid)
				return []resource.TestStep{
					{
						Config: renderConfigWithCT(te, cid, `
						resource "proxmox_replication" "test_replication" {
							id     = "${proxmox_virtual_environment_container.test_container.id}-{{.JobNum}}"
							target = "{{.Node2Name}}"
							type = "local"
							disable = true
							comment = "comment 123"
							schedule = "*/30"
							rate = 10
						}
										`),
						Check: resource.ComposeTestCheckFunc(
							test.ResourceAttributes("proxmox_replication.test_replication", map[string]string{
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
					},
					{
						Config: renderConfigWithCT(te, cid, `
						resource "proxmox_replication" "test_replication" {
							id     = "${proxmox_virtual_environment_container.test_container.id}-{{.JobNum}}"
							target = "{{.Node2Name}}"
							type = "local"
							# removed disable, comment, schedule, rate
						}
										`),
						Check: resource.ComposeTestCheckFunc(
							test.ResourceAttributes("proxmox_replication.test_replication", map[string]string{
								"id":     id,
								"target": te.Node2Name,
								"type":   "local",
								"jobnum": jobnum,
								"guest":  guest,
							}),
						),
					},
				}
			}(),
		},
		{
			"replication fields deletion and re-addition", func() []resource.TestStep {
				cid := newCID()
				id := fmt.Sprintf("%d-%s", cid, jobnum)
				guest := fmt.Sprintf("%d", cid)
				return []resource.TestStep{
					{
						Config: renderConfigWithCT(te, cid, `
						resource "proxmox_replication" "test_replication" {
							id     = "${proxmox_virtual_environment_container.test_container.id}-{{.JobNum}}"
							target = "{{.Node2Name}}"
							type = "local"
							disable = true
							comment = "comment 123"
							schedule = "*/30"
							rate = 10
						}`),
						Check: resource.ComposeTestCheckFunc(
							test.ResourceAttributes("proxmox_replication.test_replication", map[string]string{
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
					},
					{
						Config: renderConfigWithCT(te, cid, `
						resource "proxmox_replication" "test_replication" {
							id     = "${proxmox_virtual_environment_container.test_container.id}-{{.JobNum}}"
							target = "{{.Node2Name}}"
							type = "local"
							# removed disable, comment, schedule, rate
						}`),
						Check: resource.ComposeTestCheckFunc(
							test.ResourceAttributes("proxmox_replication.test_replication", map[string]string{
								"id":     id,
								"target": te.Node2Name,
								"type":   "local",
								"jobnum": jobnum,
								"guest":  guest,
							}),
						),
					},
					{
						Config: renderConfigWithCT(te, cid, `
						resource "proxmox_replication" "test_replication" {
							id     = "${proxmox_virtual_environment_container.test_container.id}-{{.JobNum}}"
							target = "{{.Node2Name}}"
							type = "local"
							disable = true
							comment = "comment 123"
							schedule = "*/30"
							rate = 10
						}`),
						Check: resource.ComposeTestCheckFunc(
							test.ResourceAttributes("proxmox_replication.test_replication", map[string]string{
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
					},
				}
			}(),
		},
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

func TestUnitResourceReplication_Validators(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				// id must be <GUEST>-<JOBNUM>
				PlanOnly: true,
				Config: `
				resource "proxmox_replication" "test" {
				id     = "invalidid"
				target = "pve2"
				type   = "local"
				}`,
				ExpectError: regexp.MustCompile(`id must be <GUEST>-<JOBNUM>`),
			},
			{
				// id must be <GUEST>-<JOBNUM>
				PlanOnly: true,
				Config: `
				resource "proxmox_replication" "test" {
				id     = "a-b"
				target = "pve2"
				type   = "local"
				}`,
				ExpectError: regexp.MustCompile(`id must be <GUEST>-<JOBNUM>`),
			},
			{
				// id must be <GUEST>-<JOBNUM>
				PlanOnly: true,
				Config: `
				resource "proxmox_replication" "test" {
				id     = "-"
				target = "pve2"
				type   = "local"
				}`,
				ExpectError: regexp.MustCompile(`id must be <GUEST>-<JOBNUM>`),
			},
			{
				// type must be "local"
				PlanOnly: true,
				Config: `
				resource "proxmox_replication" "test" {
				id     = "100-1"
				target = "pve2"
				type   = "remote"
				}
				`,
				ExpectError: regexp.MustCompile(`Attribute type value must be one of`),
			},
			{
				// missing required attribute: id
				PlanOnly: true,
				Config: `
				resource "proxmox_replication" "test" {
				target = "pve2"
				type   = "local"
				}
				`,
				ExpectError: regexp.MustCompile(`The argument "id" is required, but no definition was found.`),
			},
			{
				// missing required attribute: target
				PlanOnly: true,
				Config: `
				resource "proxmox_replication" "test" {
				id   = "100-1"
				type = "local"
				}
				`,
				ExpectError: regexp.MustCompile(`The argument "target" is required, but no definition was found.`),
			},
			{
				// missing required attribute: type
				PlanOnly: true,
				Config: `
				resource "proxmox_replication" "test" {
				id     = "100-1"
				target = "pve2"
				}
				`,
				ExpectError: regexp.MustCompile(`The argument "type" is required, but no definition was found.`),
			},
		},
	})
}
