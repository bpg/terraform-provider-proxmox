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
	"math/rand"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/proxmox/helpers/ptr"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/storage"
)

const (
	accTestContainerName = "proxmox_virtual_environment_container.test_container"
)

func TestAccResourceContainer(t *testing.T) {
	t.Parallel()

	te := InitEnvironment(t)

	imageFileName := gofakeit.Word() + "-ubuntu-23.04-standard_23.04-1_amd64.tar.zst"
	accTestContainerID := 100000 + rand.Intn(99999)
	accTestContainerIDClone := 100000 + rand.Intn(99999)

	te.AddTemplateVars(map[string]interface{}{
		"ImageFileName":        imageFileName,
		"TestContainerID":      accTestContainerID,
		"TestContainerIDClone": accTestContainerIDClone,
	})

	err := te.NodeStorageClient().DownloadFileByURL(context.Background(), &storage.DownloadURLPostRequestBody{
		Content:  ptr.Ptr("vztmpl"),
		FileName: ptr.Ptr(imageFileName),
		Node:     ptr.Ptr(te.NodeName),
		Storage:  ptr.Ptr(te.DatastoreID),
		URL:      ptr.Ptr(fmt.Sprintf("%s/images/system/ubuntu-23.04-standard_23.04-1_amd64.tar.zst", te.ContainerImagesServer)),
	})
	require.NoError(t, err)

	t.Cleanup(func() {
		e := te.NodeStorageClient().DeleteDatastoreFile(context.Background(), fmt.Sprintf("vztmpl/%s", imageFileName))
		require.NoError(t, e)
	})

	tests := []struct {
		name string
		step []resource.TestStep
	}{
		{"crete and start container", []resource.TestStep{{
			Config: te.RenderConfig(`
			resource "proxmox_virtual_environment_container" "test_container" {
				node_name = "{{.NodeName}}"
				vm_id     = {{.TestContainerID}}
				disk {
					datastore_id = "local-lvm"
					size         = 4
				}
				mount_point {
					volume = "local-lvm"
					size   = "4G"
					path   = "mnt/local"
				}
				description = <<-EOT
					my
					description
					value
				EOT
				initialization {
					hostname = "test"
					ip_config {
						ipv4 {
						  address = "dhcp"
						}
					}
				}
				network_interface {
					name = "vmbr0"
				}
				operating_system {
					template_file_id = "local:vztmpl/{{.ImageFileName}}"
					type             = "ubuntu"
				}
			}`),
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr(accTestContainerName, "description", "my\ndescription\nvalue\n"),
				func(*terraform.State) error {
					ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
					defer cancel()

					err := te.NodeClient().Container(accTestContainerID).WaitForContainerStatus(ctx, "running")
					require.NoError(te.t, err, "container did not start")

					return nil
				},
			),
		}}},
		{"update mount points", []resource.TestStep{
			{
				SkipFunc: func() (bool, error) {
					// mount point deletion: https://github.com/bpg/terraform-provider-proxmox/issues/1392
					return true, nil
				},
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_container" "test_container" {
				    node_name = "{{.NodeName}}"
				    vm_id     = {{.RandomVMID}}
				  	started   = false
				    disk {
						datastore_id = "local-lvm"
						size         = 4
					}
				    mount_point {
						volume = "local-lvm"
						size   = "4G"
						path   = "mnt/local1"
				    }
					mount_point {
						volume = "local"
						size   = "8G"
						path   = "mnt/local2"
				    }
				    initialization {
				  		hostname = "test"
						ip_config {
						  	ipv4 {
								address = "dhcp"
						  	}
						}
				    }
				    network_interface {
				  	    name = "vmbr0"
				    }
				    operating_system {
						template_file_id = "local:vztmpl/{{.ImageFileName}}"
						type             = "ubuntu"
				    }
				}`),
				Check: ResourceAttributes("proxmox_virtual_environment_container.test_container", map[string]string{
					"mount_point.#": "2",
				}),
			},
			{
				SkipFunc: func() (bool, error) {
					// mount point deletion: https://github.com/bpg/terraform-provider-proxmox/issues/1392
					return true, nil
				},
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_container" "test_container" {
				    node_name = "{{.NodeName}}"
				  	started   = false
				    disk {
						datastore_id = "local-lvm"
						size         = 4
					}
				    mount_point {
						volume = "local-lvm"
						size   = "4G"
						path   = "mnt/local1"
				    }
				    initialization {
				  		hostname = "test"
						ip_config {
						  	ipv4 {
								address = "dhcp"
						  	}
						}
				    }
				    network_interface {
				  	    name = "vmbr0"
				    }
				    operating_system {
						template_file_id = "local:vztmpl/{{.ImageFileName}}"
						type             = "ubuntu"
				    }
				}`),
				Check: ResourceAttributes("proxmox_virtual_environment_container.test_container", map[string]string{
					"mount_point.#": "2",
				}),
			},
		}},
		{"clone container", []resource.TestStep{{
			Config: te.RenderConfig(`
			resource "proxmox_virtual_environment_container" "test_container" {
				node_name = "{{.NodeName}}"
				vm_id     = {{.RandomVMID}}
				template  = true
				disk {
					datastore_id = "local-lvm"
					size         = 4
				}
				mount_point {
					volume = "local-lvm"
					size   = "4G"
					path   = "mnt/local"
				}
				initialization {
					hostname = "test"
					ip_config {
						ipv4 {
						  address = "dhcp"
						}
					}
				}
				network_interface {
					name = "vmbr0"
				}
				operating_system {
					template_file_id = "local:vztmpl/{{.ImageFileName}}"
					type             = "ubuntu"
				}
			}
			resource "proxmox_virtual_environment_container" "test_container_clone" {
				depends_on = [proxmox_virtual_environment_container.test_container]
				node_name  = "{{.NodeName}}"
				vm_id      = {{.TestContainerIDClone}}
				
				clone {
					vm_id = proxmox_virtual_environment_container.test_container.id
				}
				
				initialization {
					hostname = "test-clone"
				}
			}`),
			Check: resource.ComposeTestCheckFunc(
				func(*terraform.State) error {
					ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
					defer cancel()

					err := te.NodeClient().Container(accTestContainerIDClone).WaitForContainerStatus(ctx, "running")
					require.NoError(te.t, err, "container did not start")

					return nil
				},
			),
		}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource.ParallelTest(t, resource.TestCase{
				ProtoV6ProviderFactories: te.AccProviders,
				Steps:                    tt.step,
			})
		})
	}
}
